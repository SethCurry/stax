package magicdb

import "context"

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
		return nil, err
	}

	var id int

	err = s.conn.QueryRowContext(ctx, query, queryVars...).Scan(&id)
	if err != nil {
		return nil, err
	}

	set := newSet(s.baseClient, id, name, code)
	return set, nil
}

func (s *SetClient) Get(ctx context.Context, id int) (*Set, error) {
	builder := s.queryBuilder.Select("name", "code").From(tableSet).Where("id = ?", id)

	query, queryVars, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var name string
	var code string

	err = s.conn.QueryRowContext(ctx, query, queryVars...).Scan(&name, &code)
	if err != nil {
		return nil, err
	}

	set := newSet(s.baseClient, id, name, code)
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
