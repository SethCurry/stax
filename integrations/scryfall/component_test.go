package scryfall_test

import (
	"testing"

	"github.com/SethCurry/gofall"
	"github.com/stretchr/testify/assert"
)

func Test_Component_Equality(t *testing.T) {
	t.Parallel()

	//nolint:staticcheck
	assert.True(t, gofall.ComponentMeldPart() == gofall.ComponentMeldPart())
}

func Test_Component_Inequality(t *testing.T) {
	t.Parallel()

	assert.True(t, gofall.ComponentMeldPart() != gofall.ComponentMeldResult())
}
