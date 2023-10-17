package stax_test

import (
	"testing"

	"github.com/SethCurry/stax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseTypeLine(t *testing.T) {
	parsed, err := stax.ParseTypeLine("Instant")
	require.NoError(t, err)

	assert.Equal(t, 0, len(parsed.Subtypes))
	assert.Equal(t, 0, len(parsed.Supertypes))
	assert.Equal(t, 1, len(parsed.Types))
	assert.Equal(t, stax.Type("Instant"), parsed.Types[0])
}

func Test_ParseTypeLine_MultiType(t *testing.T) {
	parsed, err := stax.ParseTypeLine("Legendary Artifact Creature —  Sliver")
	require.NoError(t, err)

	assert.Equal(t, 1, len(parsed.Supertypes))
	assert.Equal(t, stax.Supertype("Legendary"), parsed.Supertypes[0])

	assert.Equal(t, 1, len(parsed.Subtypes))
	assert.Equal(t, stax.Subtype("Sliver"), parsed.Subtypes[0])

	assert.Equal(t, 2, len(parsed.Types))
	assert.Contains(t, parsed.Types, stax.Type("Artifact"))
	assert.Contains(t, parsed.Types, stax.Type("Creature"))
}

func Test_ParseTypeLine_Token(t *testing.T) {
	parsed, err := stax.ParseTypeLine("Token Artifact Creature — Thopter")
	require.NoError(t, err)

	assert.Equal(t, 1, len(parsed.Supertypes))
	assert.Equal(t, 2, len(parsed.Types))
	assert.Equal(t, 1, len(parsed.Subtypes))

	assert.Equal(t, stax.Supertype("Token"), parsed.Supertypes[0])
}
