package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SethCurry/stax/internal/console"
	"github.com/SethCurry/stax/pkg/scryfall"
	"go.uber.org/zap"
)

// ScryfallCmd is a command group for interacting with the Scryfall API.
type ScryfallCmd struct {
	Search  ScryfallSearchCmd  `cmd:"" help:"Search for cards"`
	Rulings ScryfallRulingsCmd `cmd:"" help:"Get rulings for a card"`
}

// ScryfallSearchCmd is the implementation of "stax scryfall search".
type ScryfallSearchCmd struct {
	Query  []string `arg:"" help:"The search query."`
	Format string   `name:"format" short:"f" help:"The output format." enum:"table,json" default:"table"`
}

func (s *ScryfallSearchCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	query := strings.Join(s.Query, " ")

	logger.Debug("searching for cards", zap.String("query", query))

	client := scryfall.NewClient(nil)

	pager, err := client.Card.Search(context.Background(), query, scryfall.CardSearchOptions{})
	if err != nil {
		return err
	}

	var writeFunc func(*scryfall.Card) error

	switch s.Format {
	case "json", "j":
		logger.Debug("using JSON output")

		os.Stdout.Write([]byte("[\n"))
		defer os.Stdout.Write([]byte("]\n"))

		writeFunc = console.ScryfallCardJSON
	case "table", "t", "":
		writer := console.NewScryfallCardTable(os.Stdout)
		defer writer.Flush()

		writeFunc = func(card *scryfall.Card) error {
			return writer.Write(card)
		}
	}

	pageNum := 0

	for pager.HasMore() {
		pageNum++

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		logger.Debug("retrieving a page of results", zap.Int("page_number", pageNum))

		cards, err := pager.Next(ctx)
		if err != nil {
			return err
		}

		for _, card := range cards {
			err := writeFunc(&card)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type ScryfallRulingsCmd struct {
	CardName []string `arg:"" help:"The name of the card"`
}

func (s *ScryfallRulingsCmd) Run(ctx *Context) error {
	logger := ctx.Logger

	client := scryfall.NewClient(nil)

	cardName := strings.Join(s.CardName, " ")

	foundCard, err := client.Card.Named(context.Background(), cardName)
	if err != nil {
		return err
	}

	logger.Debug("found card", zap.String("name", foundCard.Name))

	rulings, err := client.Rulings.ByScryfallID(context.Background(), foundCard.ID)
	if err != nil {
		return err
	}

	for _, ruling := range rulings {
		fmt.Printf("%s (%s)\n%s\n\n", ruling.PublishedAt, ruling.Source, ruling.Comment)
	}

	return nil
}
