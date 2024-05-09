package ruleparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFile(t *testing.T) {
	parsed, err := ParseFile("testdata/20240308.txt")

	require.NoError(t, err)

	assert.Len(t, parsed.Sections, 9)
}
