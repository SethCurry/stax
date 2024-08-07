package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/SethCurry/stax/pkg/moxfield"
	"github.com/SethCurry/stax/pkg/stax"
	"go.uber.org/zap"
)

type MoxfieldCmd struct {
	ExportUser MoxfieldExportUserCmd `cmd:"" help:"Export all of a user's decks"`
	HTML       RulesHTMLCmd          `cmd:"" help:"Generate stand-alone HTML documentation of the rules."`
}

type MoxfieldExportUserCmd struct {
	Username        string `arg:"username" help:"The username on Moxfield to export decks for."`
	OutputDirectory string `arg:"output-directory" help:"The directory to export the decklists to."`
}

func (m *MoxfieldExportUserCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	client := moxfield.NewClient(nil)

	decks, err := client.Users.ListUserDecks(context.Background(), m.Username, moxfield.ListUserDecksRequest{})
	if err != nil {
		return err
	}

	for _, v := range decks.Data {
		logger := logger.With(zap.String("deck_name", v.Name))

		lines, err := client.Decks.GetDeckList(context.Background(), v.ID)
		if err != nil {
			logger.Error("failed to get deck list", zap.Error(err))
			continue
		}

		outputPath := filepath.Join(m.OutputDirectory, v.Name+".txt")
		logger = logger.With(zap.String("file_path", outputPath))

		fd, err := os.Create(filepath.Join(m.OutputDirectory, v.Name+".txt"))
		if err != nil {
			logger.Error("failed to create moxfield export file", zap.Error(err))
		}
		defer fd.Close()

		writer := stax.NewMTGODecklistWriter(fd)

		for _, l := range lines {
			err = writer.AddCard(l.Name, l.Quantity)
			if err != nil {
				logger.Error("failed to add card to decklist", zap.Error(err))

				return fmt.Errorf("failed to add card to decklist: %w", err)
			}
		}
	}

	return nil
}
