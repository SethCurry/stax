package scryfall_test

import (
	"context"
	"testing"

	"github.com/SethCurry/stax/pkg/scryfall"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Client_Card_Named(t *testing.T) {
	t.Parallel()

	client := scryfall.NewClient(nil)

	card, err := client.Card.Named(context.Background(), "Black Lotus")
	require.NoError(t, err)

	assert.Equal(t, "bd8fa327-dd41-4737-8f19-2cf5eb1f7cdd", card.ID)
	assert.Equal(t, "Black Lotus", card.Name)
}

func Test_Client_Card_Search(t *testing.T) {
	t.Parallel()

	client := scryfall.NewClient(nil)

	cardPager, err := client.Card.Search(context.Background(), "Black Lotus", scryfall.CardSearchOptions{})
	require.NoError(t, err)

	cards, err := cardPager.Next(context.Background())
	require.NoError(t, err)

	assert.Len(t, cards, 2)
}

func Test_Client_Card_Autocomplete(t *testing.T) {
	t.Parallel()

	client := scryfall.NewClient(nil)

	autocomplete, err := client.Card.Autocomplete(context.Background(), "Black Lotu")
	require.NoError(t, err)

	assert.Len(t, autocomplete, 1)
}
