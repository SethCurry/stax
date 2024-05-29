package ruleparser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseFile(t *testing.T) {
	testFiles, err := os.ReadDir("testdata/")
	require.NoError(t, err, "failed to list test data files")

	for _, v := range testFiles {
		t.Run("parse-"+v.Name(), func(t *testing.T) {
			parsed, err := ParseFile("testdata/" + v.Name())
			require.NoError(t, err)
			assert.Len(t, parsed.Sections, 9)
		})
	}
}
