package common_test

import (
	"testing"

	"github.com/SethCurry/stax/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	testCases := []struct {
		name     string
		items    []string
		fn       func(string) bool
		expected []bool
	}{
		{
			name:  "basic string comparison",
			items: []string{"one", "two", "three"},
			fn: func(item string) bool {
				return item == "two"
			},
			expected: []bool{false, true, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := common.Map(tc.items, tc.fn)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name     string
		items    []string
		fn       func(string) bool
		expected []string
	}{
		{
			name:  "basic string comparison",
			items: []string{"one", "two", "three"},
			fn: func(item string) bool {
				return item == "two"
			},
			expected: []string{"two"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := common.Filter(tc.items, tc.fn)

			assert.Equal(t, tc.expected, result)
		})
	}
}
