package cmd

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SethCurry/stax/integrations/moxfield"
	"github.com/SethCurry/stax/integrations/xmage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// moxfieldExportCmd represents the moxfieldExport command
var moxfieldExportCmd = &cobra.Command{
	Use:   "export",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logger.Fatal("must provide a single username")
		}

		username := args[0]

		outputDir, err := cmd.Flags().GetString("output")
		if err != nil {
			logger.Fatal("failed to read output flag", zap.Error(err))
		}

		if _, err = os.Stat(outputDir); os.IsNotExist(err) {
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				logger.Fatal("failed to create output directory", zap.Error(err))
			}
		}

		client := moxfield.NewClient(http.DefaultClient)

		decks, err := client.ListUserDecks(username, moxfield.ListUserDecksRequest{
			PageSize: 100,
		})
		if err != nil {
			logger.Fatal("failed to list user decks", zap.Error(err))
		}

		for _, deck := range decks.Data {
			if deck.Format != "commander" {
				continue
			}

			if !deck.IsLegal {
				logger.Warn("deck is not legal", zap.String("deck_id", deck.PublicID), zap.String("name", deck.Name))
				continue
			}

			logger := logger.With(zap.String("deck_id", deck.PublicID))
			logger.Info("deck", zap.String("name", deck.Name), zap.String("id", deck.ID))
			deckList, err := client.GetDeckList(deck.PublicID)
			if err != nil {
				logger.Fatal("failed to export deck list", zap.Error(err))
			}
			deckLines := strings.Split(string(deckList), "\n")

			cards := []xmage.DeckCard{}
			sideboard := []xmage.DeckCard{}

			for idx, line := range deckLines {
				if line == "" {
					continue
				}

				parsed, err := moxfield.ParseDeckListLine(line)
				if err != nil {
					logger.Error("failed trying to parse deck list line", zap.String("line", line), zap.Error(err))
					continue
				}

				newCard := xmage.DeckCard{
					Name:            parsed.Name,
					Quantity:        parsed.Quantity,
					CollectorNumber: parsed.CollectorNumber,
					SetCode:         parsed.Set,
				}

				if idx == 0 {
					sideboard = append(sideboard, newCard)
				} else {
					cards = append(cards, newCard)
				}
			}

			xmageDeck := xmage.NewDeckList(deck.Name, cards, sideboard)
			marshalled := xmageDeck.MarshalDck()
			outputFile := filepath.Join(outputDir, deck.Name+".dck")
			outputFd, err := os.Create(outputFile)
			if err != nil {
				logger.Fatal("failed to create output file", zap.String("path", outputFile), zap.Error(err))
			}
			outputFd.Write([]byte(marshalled))
			outputFd.Close()
			logger.Info("saved deck", zap.String("name", deck.Name), zap.String("path", outputFile))
			time.Sleep(time.Second)
		}
	},
}

func init() {
	moxfieldCmd.AddCommand(moxfieldExportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moxfieldExportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moxfieldExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	moxfieldExportCmd.Flags().StringP("output", "o", "stax-moxfield", "The directory to output the files into.")
}
