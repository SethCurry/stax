package mtgdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/SethCurry/stax/mtgdb"
	"github.com/stretchr/testify/require"
)

func tempDB(t *testing.T, name string) *mtgdb.Client {
	conn, err := mtgdb.Open("sqlite3", fmt.Sprintf("file:%s?mode=memory&cache=shared", name))
	require.NoError(t, err)

	err = conn.MigrateSchema(context.Background())
	require.NoError(t, err)

	return conn
}

func Test_Migrate(t *testing.T) {
	tempDB(t, "Test_Migrate")
}
