package mtgdb

import (
	"context"
	"fmt"
	"time"
)

func newPrintingClient(base *baseClient) *PrintingClient {
	return &PrintingClient{
		baseClient: base,
	}
}

type PrintingClient struct {
	*baseClient
}

// Create creates a new ruling in the database.
func (p *PrintingClient) Create(
	ctx context.Context,
	cardID int,
	rulingText string,
	rulingDate time.Time,
) (*Ruling, error) {
	rulingString := rulingDate.Format("2006/01/02")

	builder := r.queryBuilder.
		Insert("ruling").
		Columns("card_id", "ruling_text", "ruling_date").
		Values(cardID, rulingText, rulingString).
		Suffix("RETURNING id")

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query to insert ruling: %w", err)
	}

	var rulingID int

	err = r.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&rulingID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new ruling: %w", err)
	}

	ruling := newRuling(r.baseClient, rulingID, cardID, rulingText, rulingString)

	return ruling, nil
}

// Get gets a ruling by ID.
func (r *RulingClient) Get(ctx context.Context, rulingID int) (*Ruling, error) {
	builder := r.queryBuilder.
		Select("card_id", "ruling_text", "ruling_date").
		From("ruling").
		Where("id = ?", rulingID)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query to get ruling by ID: %w", err)
	}

	var cardID int

	var rulingText, rulingDate string

	err = r.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&rulingID, &cardID, &rulingText, &rulingDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get ruling %d by ID: %w", rulingID, err)
	}

	ruling := newRuling(r.baseClient, rulingID, cardID, rulingText, rulingDate)

	return ruling, nil
}

func newRuling(base *baseClient, id int, cardID int, text string, rulingDate string) *Ruling {
	return &Ruling{
		ID:         id,
		CardID:     cardID,
		Text:       text,
		Date:       rulingDate,
		baseClient: base,
	}
}

type Ruling struct {
	ID     int
	CardID int
	Text   string
	Date   string
	*baseClient
}
