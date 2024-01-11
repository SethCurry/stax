package magicdb

import "context"

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
		return nil, err
	}

	var id int

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&id)
	if err != nil {
		return nil, err
	}

	artist := newArtist(a.baseClient, id, name)
	return artist, nil
}

func (a *ArtistClient) Get(ctx context.Context, id int) (*Artist, error) {
	builder := a.queryBuilder.Select("name").From("artist").Where("id = ?", id)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var name string

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&name)
	if err != nil {
		return nil, err
	}

	artist := newArtist(a.baseClient, id, name)
	return artist, nil
}

func (a *ArtistClient) GetByName(ctx context.Context, name string) (*Artist, error) {
	builder := a.queryBuilder.Select("id").From("artist").Where("name = ?", name)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var id int

	err = a.conn.QueryRowContext(ctx, query, queryVars...).Scan(&id)
	if err != nil {
		return nil, err
	}

	artist := newArtist(a.baseClient, id, name)
	return artist, nil
}

func newArtist(base *baseClient, id int, name string) *Artist {
	return &Artist{
		ID:         id,
		baseClient: base,
	}
}

type Artist struct {
	ID   int
	Name string
	*baseClient
}
