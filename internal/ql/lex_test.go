package ql

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLex(t *testing.T) {
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
		{
			name:  "quoted literal",
			query: `name="Some Long Name"`,
			expected: []Token{
				{
					Family: FamilyLiteral,
					Value:  "name",
				},
				{
					Family: FamilyOperator,
					Value:  "=",
				},
				{
					Family: FamilyLiteral,
					Value:  "Some Long Name",
				},
			},
		},
		{
			name:  "has keyword",
			query: "name=something OR name=other",
			expected: []Token{
				{
					Family: FamilyLiteral,
					Value:  "name",
				},
				{
					Family: FamilyOperator,
					Value:  "=",
				},
				{
					Family: FamilyLiteral,
					Value:  "something",
				},
				{
					Family: FamilyKeyword,
					Value:  "OR",
				},
				{
					Family: FamilyLiteral,
					Value:  "name",
				},
				{
					Family: FamilyOperator,
					Value:  "=",
				},
				{
					Family: FamilyLiteral,
					Value:  "other",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := lexString(tc.query)

			require.NoError(t, err)

			assert.Equal(t, tc.expected, got)
		})
	}
}
