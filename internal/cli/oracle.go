package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"

	"github.com/SethCurry/stax/internal/oracle/etl"
	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/SethCurry/stax/pkg/scryfall"
)

func dataDirectory() (string, error) {
	userHome, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	dataPath := filepath.Join(userHome, ".local", ".share", "stax")

	if err = os.MkdirAll(dataPath, 0o755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataPath, nil
}

func connectToDatabase(ctx context.Context) (*oracledb.Client, error) {
	dataDir, err := dataDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	conn, err := oracledb.Open("sqlite3", filepath.Join(dataDir, "oracle.sqlite3?_fk=1&cache=shared"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = conn.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return conn, nil
}

// OracleCmd is a command group for interacting with the Oracle database of MTG cards.
type OracleCmd struct {
	Load  OracleLoadCmd  `cmd:"" help:"Load cards from the Scryfall API."`
	Reset OracleResetCmd `cmd:"" help:"Reset the Oracle database."`
}

type OracleResetCmd struct{}

func (r *OracleResetCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	dataDir, err := dataDirectory()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "oracle.sqlite3")

	logger.Info("removing database file", zap.String("path", dbPath))

	err = os.Remove(dbPath)
	if err != nil {
		return fmt.Errorf("failed to remove database file: %w", err)
	}

	return nil
}

type OracleLoadCmd struct {
	Args []string `arg:"" optional:"" help:"The path to the Scryfall bulk data file."`
	HTTP bool     `name:"http" help:"Use the HTTP client instead of the default client."`
}

func (r *OracleLoadCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	sfall := scryfall.NewClient(nil)

	if len(r.Args) > 1 {
		return errors.New("too many arguments")
	}

	var reader *scryfall.BulkReader[scryfall.Card]

	if r.HTTP {
		currentBulkFiles, err := sfall.BulkData.ListSources(ctx.Context)
		if err != nil {
			return fmt.Errorf("failed to get current bulk files from Scryfall: %w", err)
		}

		resp, err := http.Get(currentBulkFiles.DefaultCards.DownloadURI)
		if err != nil {
			return fmt.Errorf("failed to download default cards file: %w", err)
		}

		defer resp.Body.Close()

		reader, err = scryfall.NewBulkReader[scryfall.Card](resp.Body)
		if err != nil {
			return fmt.Errorf("failed to create bulk reader: %w", err)
		}
	} else {
		if len(r.Args) != 1 {
			return errors.New("missing argument: path to Scryfall bulk data file")
		}

		fd, err := os.Open(r.Args[0])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		defer fd.Close()

		reader, err = scryfall.NewBulkReader[scryfall.Card](fd)
		if err != nil {
			return fmt.Errorf("failed to create bulk reader: %w", err)
		}
	}

	dbClient, err := connectToDatabase(ctx.Context)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer dbClient.Close()

	err = etl.ScryfallCards(ctx.Context, logger, dbClient, reader)
	if err != nil {
		return fmt.Errorf("failed to load cards from Scryfall: %w", err)
	}

	return nil
}
