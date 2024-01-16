package mtgdb

import (
	"context"
	"fmt"
)

const tableCard = "card"

func newCardClient(base *baseClient) *CardClient {
	return &CardClient{
		baseClient: base,
	}
}

type CardClient struct {
	*baseClient
}

func (c *CardClient) Create(
	ctx context.Context,
	name string,
	oracleID string,
	colorIdentity string,
) (*Card, error) {
	builder := c.queryBuilder.Insert(tableCard).Columns(
		"name",
		"oracle_id",
		"color_identity").Values(
		name,
		oracleID,
		colorIdentity).Suffix("RETURNING id")

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL query to create card: %w", err)
	}

	var cardID int

	err = c.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new card in database: %w", err)
	}

	card := newCard(
		c.baseClient,
		cardID,
		name,
		oracleID,
		colorIdentity,
	)

	return card, nil
}

func newCard(
	base *baseClient,
	cardID int,
	name string,
	oracleID string,
	colorIdentity string,
) *Card {
	return &Card{
		baseClient:    base,
		ID:            cardID,
		Name:          name,
		OracleID:      oracleID,
		ColorIdentity: colorIdentity,
	}
}

type Card struct {
	*baseClient
	ID            int
	Name          string
	OracleID      string
	ColorIdentity string
}
