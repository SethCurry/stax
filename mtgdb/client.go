package magicdb

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Open(driverName, dataSourceName string) (*Client, error) {
	conn, err := sqlx.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	switch driverName {
	case "sqlite3":
		_, err := conn.Exec("PRAGMA foreign_keys = ON;")
		if err != nil {
			conn.Close()
			return nil, err
		}
	case "postgres":
		queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	}
	return &Client{
		&baseClient{
			conn:         conn,
			queryBuilder: queryBuilder,
		},
	}, nil
}

// Client is the top-level interface for interacting with the database.
// It primarily returns other model-specific structs.
type Client struct {
	*baseClient
}

func (c *Client) Artists() *ArtistClient {
	return newArtistClient(c.baseClient)
}
