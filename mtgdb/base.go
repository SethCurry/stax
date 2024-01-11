package mtgdb

import (
	"github.com/Masterminds/squirrel"
)

type baseClient struct {
	client       *Client
	queryBuilder squirrel.StatementBuilderType
}
