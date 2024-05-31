package etl

import (
	"context"
	"fmt"

	"github.com/SethCurry/stax/internal/bones"
	"github.com/SethCurry/stax/internal/bones/card"
	"github.com/SethCurry/stax/pkg/scryfall"
	"go.uber.org/zap"
)

func findCardByName(
	ctx context.Context,
	tx *bones.Tx,
	cardName string,
) (*bones.Card, error) {
	return tx.Card.Query().Where(card.NameEQ(cardName)).Only(ctx)
}

func createCard(
	ctx context.Context,
	tx *bones.Tx,
	cardName string,
	oracleID string,
	colorIdentity uint8,
) (*bones.Card, error) {
	return tx.Card.Create().
		SetName(cardName).
		SetOracleID(oracleID).
		SetColorIdentity(colorIdentity).
		Save(ctx)
}

func getOrCreateCard(
	ctx context.Context,
	logger *zap.Logger,
	db *bones.Tx,
	row *scryfall.Card,
	isFresh bool,
) (*bones.Card, error) {
	logger = logger.With(
		zap.String("card_name", row.Name),
		zap.String("oracle_id", row.OracleID),
		zap.Strings("colors", row.Colors),
	)

	if !isFresh {
		existingCard, err := findCardByName(ctx, db, row.Name)
		if err == nil {
			logger.Debug("card already exists")
			return existingCard, nil
		}

		if !bones.IsNotFound(err) {
			logger.Error("failed to query for existing card", zap.Error(err))
			return nil, fmt.Errorf("failed to query for existing card: %w", err)
		}
	}

	// TODO write function to convert row.Colors to a bitfield
	newCard, err := createCard(ctx, db, row.Name, row.OracleID, 0)
	if err != nil {
		logger.Error("failed to create new card", zap.Error(err))
		return nil, fmt.Errorf("failed to create new card: %w", err)
	}

	logger.Info("created new card")

	return newCard, nil
}
