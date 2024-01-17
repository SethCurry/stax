package mtgdb

import (
	"context"
	"fmt"
)

/*

CREATE TABLE IF NOT EXISTS card_face (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    flavor_text TEXT NOT NULL,
    oracle_text TEXT NOT NULL,
    lang TEXT NOT NULL,
    cmc float,
    power_ TEXT,
    toughness TEXT,
    loyalty TEXT,
    mana_cost TEXT,
    type_line TEXT,
    colors TEXT
);

*/

func newCardFaceClient(base *baseClient) *CardFaceClient {
	return &CardFaceClient{
		baseClient: base,
	}
}

type CardFaceClient struct {
	*baseClient
}

// Create creates a new card face in the database.
func (c *CardFaceClient) Create(
	ctx context.Context,
	name string,
	flavorText string,
	oracleText string,
	lang string,
	cmc float32,
	power string,
	toughness string,
	loyalty string,
	manaCost string,
	typeLine string,
	colors string,
) (*CardFace, error) {
	builder := c.queryBuilder.
		Insert("card_face").
		Columns(
			"name",
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
		).
		Values(
			name,
			flavorText,
			oracleText,
			lang,
			cmc,
			power,
			toughness,
			loyalty,
		).
		Suffix("RETURNING id")

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query to insert card face: %w", err)
	}

	var cardFaceID int

	err = c.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&cardFaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new card face: %w", err)
	}

	cardFace := newCardFace(
		c.baseClient,
		cardFaceID,
		name,
		flavorText,
		oracleText,
		lang,
		cmc,
		power,
		toughness,
		loyalty,
		manaCost,
		typeLine,
		colors)

	return cardFace, nil
}

// Get gets a card face by ID.
//
//nolint:funlen
func (c *CardFaceClient) Get(ctx context.Context, cardFaceID int) (*CardFace, error) {
	builder := c.queryBuilder.
		Select("name",
			"flavor_text",
			"oracle_text",
			"lang",
			"cmc",
			"power_",
			"toughness",
			"loyalty",
			"mana_cost",
			"type_line",
			"colors").
		From("card_face").
		Where("id = ?", cardFaceID)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query to get card face by ID: %w", err)
	}

	var name,
		flavorText,
		oracleText,
		lang,
		power,
		toughness,
		loyalty,
		manaCost,
		typeLine,
		colors string

	var cmc float32

	err = c.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&name,
		&flavorText,
		&oracleText,
		&lang,
		&cmc,
		&power,
		&toughness,
		&loyalty,
		&manaCost,
		&typeLine,
		&colors)
	if err != nil {
		return nil, fmt.Errorf("failed to get card face %d by ID: %w", cardFaceID, err)
	}

	cardFace := newCardFace(c.baseClient,
		cardFaceID,
		name,
		flavorText,
		oracleText,
		lang,
		cmc,
		power,
		toughness,
		loyalty,
		manaCost,
		typeLine,
		colors)

	return cardFace, nil
}

func newCardFace(
	base *baseClient,
	cardFaceID int,
	name string,
	flavorText string,
	oracleText string,
	lang string,
	cmc float32,
	power string,
	toughness string,
	loyalty string,
	manaCost string,
	typeLine string,
	colors string,
) *CardFace {
	return &CardFace{
		ID:         cardFaceID,
		Name:       name,
		FlavorText: flavorText,
		OracleText: oracleText,
		Language:   lang,
		CMC:        cmc,
		Power:      power,
		Toughness:  toughness,
		Loyalty:    loyalty,
		ManaCost:   manaCost,
		TypeLine:   typeLine,
		Colors:     colors,
		baseClient: base,
	}
}

type CardFace struct {
	ID         int
	Name       string
	FlavorText string
	OracleText string
	Language   string
	CMC        float32
	Power      string
	Toughness  string
	Loyalty    string
	ManaCost   string
	TypeLine   string
	Colors     string
	*baseClient
}
