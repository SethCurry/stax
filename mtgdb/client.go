package mtgdb

import (
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

	return &Client{
		conn:         conn,
		queryBuilder: queryBuilder,
	}, nil
}

// Client is the top-level interface for interacting with the database.
// It primarily returns other model-specific structs.
type Client struct {
	conn         *sqlx.DB
	queryBuilder squirrel.StatementBuilderType
	*baseClient
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
