package etl

import (
	"context"
	"fmt"
	"io"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
	"github.com/SethCurry/stax/internal/oracle/oracledb/artist"
	"github.com/SethCurry/stax/internal/oracle/oracledb/set"
	"github.com/SethCurry/stax/pkg/scryfall"
	"go.uber.org/zap"
)

func ScryfallCards(ctx context.Context, logger *zap.Logger, db *oracledb.Client, reader *scryfall.BulkReader[scryfall.Card]) error {
	for {
		card, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				logger.Debug("got EOF")
				return nil
			}

			return err
		}

		if card == nil {
			logger.Debug("card is nil, breaking")
			break
		}

		err = scryfallCardIngestor(ctx, logger.With(zap.String("card_name", card.Name)), db, card)
		if err != nil {
			return err
		}
	}
	return nil
}

func scryfallCardIngestor(ctx context.Context, logger *zap.Logger, db *oracledb.Client, row *scryfall.Card) error {
	cardArtist, err := getOrCreateCardArtist(ctx, logger, db, row.Artist)
	if err != nil {
		return fmt.Errorf("failed to get or create card artist: %w", err)
	}

	cardSet, err := getOrCreateSet(ctx, logger, db, row.SetName, row.SetCode)
	if err != nil {
		return fmt.Errorf("failed to get or create card set: %w", err)
	}

	return nil
}

func getOrCreateSet(ctx context.Context, logger *zap.Logger, db *oracledb.Client, setName string, setCode string) (*oracledb.Set, error) {
	existingSet, err := db.Set.Query().Where(set.NameEQ(setName)).Only(ctx)
	if err == nil {
		return existingSet, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing set: %w", err)
	}

	newSet, err := db.Set.Create().SetName(setName).SetCode(setCode).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new set: %w", err)
	}

	return newSet, nil
}

func getOrCreateCardArtist(ctx context.Context, logger *zap.Logger, db *oracledb.Client, artistName string) (*oracledb.Artist, error) {
	existingArtist, err := db.Artist.Query().Where(artist.NameEQ(artistName)).Only(ctx)
	if err == nil {
		return existingArtist, nil
	}

	if !oracledb.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query for existing artist: %w", err)
	}

	newArtist, err := db.Artist.Create().SetName(artistName).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create new artist: %w", err)
	}

	return newArtist, nil
}
