package mtgdb

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type baseClient struct {
	conn         *sqlx.DB
	queryBuilder squirrel.StatementBuilderType
}
