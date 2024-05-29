package etl

import (
	"context"
	"testing"

	"github.com/SethCurry/stax/internal/oracle/oracledb"
)

func NewTestDB(t *testing.T) *oracledb.Client {
	conn, err := oracledb.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := conn.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to create test database schema: %v", err)
	}

	return conn
}
