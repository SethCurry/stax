package scryfall_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/SethCurry/stax/integrations/scryfall"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBulkData(t *testing.T) {
	t.Parallel()

	client := scryfall.NewClient(nil)

	sources, err := client.BulkData.ListSources(context.Background())
	require.NoError(t, err)

	assert.NotNil(t, sources.AllCards)
	assert.NotNil(t, sources.DefaultCards)
	assert.NotNil(t, sources.OracleCards)
	assert.NotNil(t, sources.Rulings)
	assert.NotNil(t, sources.UniqueArtwork)

	t.Run("DefaultCards", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, sources.DefaultCards.DownloadURI, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bulkReader, err := scryfall.NewBulkReader[scryfall.Card](resp.Body)
		require.NoError(t, err)

		for {
			card, err := bulkReader.Next()
			if err != nil && errors.Is(err, io.EOF) {
				break
			}

			require.NoError(t, err)
			assert.NotNil(t, card)
		}
	})

	t.Run("Rulings", func(t *testing.T) {
		t.Parallel()

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, sources.Rulings.DownloadURI, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		bulkReader, err := scryfall.NewBulkReader[scryfall.Card](resp.Body)
		require.NoError(t, err)

		for {
			rule, err := bulkReader.Next()
			if err != nil && errors.Is(err, io.EOF) {
				break
			}

			require.NoError(t, err)
			assert.NotNil(t, rule)
		}
	})
}
