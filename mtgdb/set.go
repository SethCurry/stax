package mtgdb

import (
	"context"
	"fmt"
)

const tableSet = "set_"

func newSetClient(base *baseClient) *SetClient {
	return &SetClient{
		baseClient: base,
	}
}

type SetClient struct {
	*baseClient
}

func (s *SetClient) Create(ctx context.Context, name string, code string) (*Set, error) {
	builder := s.queryBuilder.Insert(tableSet).Columns("name", "code").Values(name, code).Suffix("RETURNING id")

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate SQL query to create set: %w", err)
	}

	var setID int

	err = s.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&setID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new set in database: %w", err)
	}

	set := newSet(s.baseClient, setID, name, code)

	return set, nil
}

func (s *SetClient) Get(ctx context.Context, setID int) (*Set, error) {
	builder := s.queryBuilder.Select("name", "code").From(tableSet).Where("id = ?", setID)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query to get set %d by ID: %w", setID, err)
	}

	var name, code string

	err = s.client.conn.QueryRowContext(ctx, query, queryVars...).Scan(&name, &code)
	if err != nil {
		return nil, fmt.Errorf("failed to get set %d by ID from database: %w", setID, err)
	}

	set := newSet(s.baseClient, setID, name, code)

	return set, nil
}

func newSet(base *baseClient, id int, name string, code string) *Set {
	return &Set{
		baseClient: base,
		ID:         id,
		Name:       name,
		Code:       code,
	}
}

type Set struct {
	*baseClient
	ID   int
	Name string
	Code string
}
