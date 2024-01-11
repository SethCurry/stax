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
	rarity string,
	flavorText string,
	oracleText string,
	language string,
	cmc float32,
	power string,
	toughness string,
	loyalty string,
	manaCost string,
	typeLine string,
	colors string,
	colorIdentity string,
) (*Card, error) {
	builder := c.queryBuilder.Insert(tableCard).Columns(
		"name",
		"oracle_id",
		"rarity",
		"flavor_text",
		"oracle_text",
		"lang",
		"cmc",
		"power_",
		"toughness",
		"loyalty",
		"mana_cost",
		"type_line",
		"colors",
		"color_identity").Values(
		name,
		oracleID,
		rarity,
		flavorText,
		oracleText,
		language,
		cmc,
		power,
		toughness,
		loyalty,
		manaCost,
		typeLine,
		colors,
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
		rarity,
		flavorText,
		oracleText,
		language,
		cmc,
		power,
		toughness,
		loyalty,
		manaCost,
		typeLine,
		colors,
		colorIdentity,
	)

	return card, nil
}

func newCard(
	base *baseClient,
	cardID int,
	name string,
	oracleID string,
	rarity string,
	flavorText string,
	oracleText string,
	language string,
	cmc float32,
	power string,
	toughness string,
	loyalty string,
	manaCost string,
	typeLine string,
	colors string,
	colorIdentity string,
) *Card {
	return &Card{
		baseClient:    base,
		ID:            cardID,
		Name:          name,
		OracleID:      oracleID,
		Rarity:        rarity,
		FlavorText:    flavorText,
		OracleText:    oracleText,
		Language:      language,
		CMC:           cmc,
		Power:         power,
		Toughness:     toughness,
		Loyalty:       loyalty,
		ManaCost:      manaCost,
		TypeLine:      typeLine,
		Colors:        colors,
		ColorIdentity: colorIdentity,
	}
}

type Card struct {
	*baseClient
	ID            int
	Name          string
	OracleID      string
	Rarity        string
	FlavorText    string
	OracleText    string
	Language      string
	CMC           float32
	Power         string
	Toughness     string
	Loyalty       string
	ManaCost      string
	TypeLine      string
	Colors        string
	ColorIdentity string
}
