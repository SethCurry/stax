package ql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenize(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected []Token
	}{
		{
			name:  "single query item",
			query: "colors<RWB",
			expected: []Token{
				{
					Family: FamilyLiteral,
					Value:  "colors",
				},
				{
					Family: FamilyOperator,
					Value:  "<",
				},
				{
					Family: FamilyLiteral,
					Value:  "RWB",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Tokenize(tc.query)

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
