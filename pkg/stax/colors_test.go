package stax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorFieldHasWhite(t *testing.T) {
	testCases := []struct {
		name       string
		colorField uint8
		expected   bool
	}{
		{
			name:       "has white",
			colorField: 0b10101,
			expected:   true,
		},
		{
			name:       "does not have white",
			colorField: 0b01110,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ColorField(tc.colorField).HasWhite()

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestColorFieldSetWhite(t *testing.T) {
	field := ColorField(0)
	fieldRef := &field

	assert.Equal(t, false, fieldRef.HasWhite())

	fieldRef.SetWhite(true)

	assert.Equal(t, true, fieldRef.HasWhite())

	fieldRef.SetWhite(false)

	assert.Equal(t, false, fieldRef.HasWhite())
}

func TestColorFieldSetBlue(t *testing.T) {
	field := ColorField(0)
	fieldRef := &field

	assert.Equal(t, false, fieldRef.HasBlue())

	fieldRef.SetBlue(true)

	assert.Equal(t, true, fieldRef.HasBlue())

	fieldRef.SetBlue(false)

	assert.Equal(t, false, fieldRef.HasBlue())
}

func TestColorFieldSetBlack(t *testing.T) {
	field := ColorField(0)
	fieldRef := &field

	assert.Equal(t, false, fieldRef.HasBlack())

	fieldRef.SetBlack(true)

	assert.Equal(t, true, fieldRef.HasBlack())

	fieldRef.SetBlack(false)

	assert.Equal(t, false, fieldRef.HasBlack())
}

func TestColorFieldSetRed(t *testing.T) {
	field := ColorField(0)
	fieldRef := &field

	assert.Equal(t, false, fieldRef.HasRed())

	fieldRef.SetRed(true)

	assert.Equal(t, true, fieldRef.HasRed())

	fieldRef.SetRed(false)

	assert.Equal(t, false, fieldRef.HasRed())
}

func TestColorFieldSetGreen(t *testing.T) {
	field := ColorField(0)
	fieldRef := &field

	assert.Equal(t, false, fieldRef.HasGreen())

	fieldRef.SetGreen(true)

	assert.Equal(t, true, fieldRef.HasGreen())

	fieldRef.SetGreen(false)

	assert.Equal(t, false, fieldRef.HasGreen())
}
