package mtgdb

import (
	"context"
	"fmt"
)

func newArtistClient(base *baseClient) *ArtistClient {
	return &ArtistClient{
		baseClient: base,
	}
}

type ArtistClient struct {
	*baseClient
}

func (a *ArtistClient) Create(ctx context.Context, name string) (*Artist, error) {
	builder := a.queryBuilder.Insert("artist").Columns("name").Values(name).Suffix("RETURNING id")

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query to insert artist: %w", err)
	}

	var artistID int

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&artistID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new artist: %w", err)
	}

	artist := newArtist(a.baseClient, artistID, name)

	return artist, nil
}

func (a *ArtistClient) Get(ctx context.Context, artistID int) (*Artist, error) {
	builder := a.queryBuilder.Select("name").From("artist").Where("id = ?", artistID)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate query to get artist by ID: %w", err)
	}

	var name string

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&name)
	if err != nil {
		return nil, fmt.Errorf("failed to get artist %d by ID: %w", artistID, err)
	}

	artist := newArtist(a.baseClient, artistID, name)

	return artist, nil
}

func (a *ArtistClient) GetByName(ctx context.Context, name string) (*Artist, error) {
	builder := a.queryBuilder.Select("id").From("artist").Where("name = ?", name)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query to get artist %q by name: %w", name, err)
	}

	var artistID int

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&artistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get artist %q by name: %w", name, err)
	}

	artist := newArtist(a.baseClient, artistID, name)

	return artist, nil
}

func newArtist(base *baseClient, id int, name string) *Artist {
	return &Artist{
		ID:         id,
		Name:       name,
		baseClient: base,
	}
}

type Artist struct {
	ID   int
	Name string
	*baseClient
}
