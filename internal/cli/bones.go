package cli

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"

	"github.com/SethCurry/stax/internal/bones"
	"github.com/SethCurry/stax/internal/etl"
	"github.com/SethCurry/stax/pkg/scryfall"
)

func dataDirectory() (string, error) {
	userHome, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	dataPath := filepath.Join(userHome, ".local", "share", "stax")

	if err = os.MkdirAll(dataPath, 0o755); err != nil {
		return "", fmt.Errorf("failed to create data directory: %w", err)
	}

	return dataPath, nil
}

func connectToDatabase(ctx context.Context, logger *zap.Logger, disableJournal bool) (*bones.Client, error) {
	dataDir, err := dataDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "oracle.sqlite3?_fk=1&cache=shared")
	logger.Info("connecting to database", zap.String("path", dbPath))

	sqlConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if disableJournal {
		db := sqlConn.DB()
		_, err := db.Exec("PRAGMA journal_mode = OFF")
		if err != nil {
			return nil, fmt.Errorf("failed to set journal_mode=OFF: %w", err)
		}
	}

	conn := bones.NewClient(bones.Driver(sqlConn))

	if err = conn.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return conn, nil
}

// BonesCmd is a command group for interacting with the bones database of MTG cards.
type BonesCmd struct {
	Load  BonesLoadCmd  `cmd:"" help:"Load cards from the Scryfall API."`
	Reset BonesResetCmd `cmd:"" help:"Reset the database."`
}

type BonesResetCmd struct{}

func (r *BonesResetCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	dataDir, err := dataDirectory()
	if err != nil {
		return fmt.Errorf("failed to get data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "bones.sqlite3")

	logger.Info("removing database file", zap.String("path", dbPath))

	err = os.Remove(dbPath)
	if err != nil {
		return fmt.Errorf("failed to remove database file: %w", err)
	}

	return nil
}

type BonesLoadCmd struct {
	DataFile string `arg:"" optional:"" help:"The path to the Scryfall bulk data file."`
	HTTP     bool   `name:"http" help:"Use the HTTP client instead of the default client."`
}

func (r *BonesLoadCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	sfall := scryfall.NewClient(nil)

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
		fd, err := os.Open(r.DataFile)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		defer fd.Close()

		reader, err = scryfall.NewBulkReader[scryfall.Card](fd)
		if err != nil {
			return fmt.Errorf("failed to create bulk reader: %w", err)
		}
	}

	dbClient, err := connectToDatabase(ctx.Context, logger, true)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	defer dbClient.Close()

	fmt.Println("Starting ingest; this will take several minutes.")

	err = etl.ScryfallCards(ctx.Context, logger, dbClient, reader)
	if err != nil {
		return fmt.Errorf("failed to load cards from Scryfall: %w", err)
	}

	return nil
}
