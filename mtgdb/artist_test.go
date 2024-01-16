package mtgdb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Artist_Create(t *testing.T) {
	db := tempDB(t, "Test_Artist_Create")

	artistName := "Seth Curry"

	artist, err := db.Artists().Create(context.Background(), artistName)
	require.NoError(t, err)

	assert.Equal(t, artistName, artist.Name)
	assert.Equal(t, 1, artist.ID)

	gotByID, err := db.Artists().Get(context.Background(), artist.ID)
	require.NoError(t, err)

	assert.Equal(t, artistName, gotByID.Name)

	gotByName, err := db.Artists().GetByName(context.Background(), artistName)
	require.NoError(t, err)
	assert.Equal(t, 1, gotByName.ID)
}
