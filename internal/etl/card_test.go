package etl

import (
	"context"
	"testing"

	"github.com/SethCurry/stax/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindCardByName(t *testing.T) {
	db := testutils.NewDB(t)

	ctx := context.Background()

	tx, err := db.Tx(ctx)
	require.NoError(t, err)

	_, err = createCard(ctx, tx, "cardName", "someOracleID", 0)
	require.NoError(t, err)

	foundCard, err := findCardByName(ctx, tx, "cardName")
	require.NoError(t, err)
	assert.NotNil(t, foundCard)

	foundCard, err = findCardByName(ctx, tx, "notExist")
	assert.Error(t, err)
	assert.Nil(t, foundCard)

	err = tx.Commit()
	require.NoError(t, err)
}
