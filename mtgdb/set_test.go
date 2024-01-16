package mtgdb_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Set(t *testing.T) {
	db := tempDB(t, "Test_Set")

	setName := "Throne of Eldraine"
	setCode := "ELD"

	set, err := db.Sets().Create(context.Background(), setName, setCode)
	require.NoError(t, err)

	assert.Equal(t, setName, set.Name)
	assert.Equal(t, setCode, set.Code)
	assert.Equal(t, 1, set.ID)

	gotByID, err := db.Sets().Get(context.Background(), set.ID)
	require.NoError(t, err)

	assert.Equal(t, setName, gotByID.Name)
	assert.Equal(t, setCode, gotByID.Code)
}
