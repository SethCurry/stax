package stax

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTGODecklistWriter(t *testing.T) {
	type writeCard struct {
		Name  string
		Count int
	}
	testCases := []struct {
		name     string
		cards    []writeCard
		expected string
	}{
		{
			name: "basic test",
			cards: []writeCard{
				{
					Name:  "Jace",
					Count: 1,
				},
				{
					Name:  "Kruphix",
					Count: 2,
				},
				{
					Name:  "Somebody",
					Count: 4,
				},
			},
			expected: "1 Jace\n2 Kruphix\n4 Somebody\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer([]byte{})

			writer := NewMTGODecklistWriter(buf)

			for _, v := range tc.cards {
				err := writer.AddCard(v.Name, v.Count)
				require.NoError(t, err)
			}

			assert.Equal(t, tc.expected, buf.String())
		})
	}
}
