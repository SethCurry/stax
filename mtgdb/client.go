package mtgdb

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	// sqlite3 driver.
	_ "github.com/mattn/go-sqlite3"
)

// Open opens a new connection to the database.  Only "sqlite3" is currently
// supported.
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

// MigrateSchema attempts to migrate the database schema to the latest version.
// This is only done with "CREATE IF NOT EXISTS" queries, so it will not update
// the schema if it has already been created.
//
// You should pin to a release if you want to use a consistent schema with consistent data.
// If you are okay with wiping the database and re-loading data into it, you can be more flexible.
func (c *Client) MigrateSchema(ctx context.Context) error {
	_, err := c.conn.ExecContext(ctx, string(schema))
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	return nil
}

// Artists returns a client for interfacing with Artist objects.
func (c *Client) Artists() *ArtistClient {
	return newArtistClient(c.baseClient)
}

// Sets returns a client for interfacing with Set objects.
func (c *Client) Sets() *SetClient {
	return newSetClient(c.baseClient)
}

// Cards returns a client for interfacing with Card objects.
func (c *Client) Cards() *CardClient {
	return newCardClient(c.baseClient)
}

// Rulings returns a client for interfacing with Ruling objects.
func (c *Client) Rulings() *RulingClient {
	return newRulingClient(c.baseClient)
}

// CardFaces returns a client for interfacing with CardFaces.
//
// Note that Card objects lack fields like oracle text and mana cost.
// You can find those items here on CardFaces, which are linked to a Card.
func (c *Client) CardFaces() *CardFaceClient {
	return newCardFaceClient(c.baseClient)
}
