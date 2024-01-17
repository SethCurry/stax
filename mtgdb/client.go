package mtgdb

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	// sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

func Open(driverName, dataSourceName string) (*Client, error) {
	conn, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("sqlx failed to open database: %w", err)
	}

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	switch driverName {
	case "sqlite3":
		// enable foreign keys
		_, err := conn.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			conn.Close()

			return nil, fmt.Errorf("failed to enable foreign keys on sqlite database: %w", err)
		}
	case "postgres":
		// use correct variable format for postgres
		queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	}

	client := &Client{
		conn: conn,
	}

	client.baseClient = &baseClient{
		client:       client,
		queryBuilder: queryBuilder,
	}

	return client, nil
}

// Client is the top-level interface for interacting with the database.
// It primarily returns other model-specific structs.
type Client struct {
	conn *sqlx.DB
	*baseClient
}

func (c *Client) MigrateSchema(ctx context.Context) error {
	_, err := c.conn.ExecContext(ctx, string(schema))
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	return nil
}

func (c *Client) Artists() *ArtistClient {
	return newArtistClient(c.baseClient)
}

func (c *Client) Sets() *SetClient {
	return newSetClient(c.baseClient)
}

func (c *Client) Cards() *CardClient {
	return newCardClient(c.baseClient)
}

func (c *Client) Rulings() *RulingClient {
	return newRulingClient(c.baseClient)
}

func (c *Client) CardFaces() *CardFaceClient {
	return newCardFaceClient(c.baseClient)
}
