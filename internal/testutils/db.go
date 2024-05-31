package testutils

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/SethCurry/stax/internal/bones"
)

func NewDB(t *testing.T) *bones.Client {
	conn, err := bones.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	if err := conn.Schema.Create(context.Background()); err != nil {
		t.Fatalf("failed to create test database schema: %v", err)
	}

	return conn
}
