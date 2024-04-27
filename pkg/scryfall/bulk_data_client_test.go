package scryfall_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/SethCurry/stax/pkg/scryfall"
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

func TestBulkReader(t *testing.T) {
	t.Parallel()
	//nolint:lll
	testData := `[{"object":"card","id":"0000579f-7b35-4ed3-b44c-db2a538066fe","oracle_id":"44623693-51d6-49ad-8cd7-140505caf02f","multiverse_ids":[109722],"mtgo_id":25527,"mtgo_foil_id":25528,"tcgplayer_id":14240,"cardmarket_id":13850,"name":"Fury Sliver","lang":"en","released_at":"2006-10-06","uri":"https://api.scryfall.com/cards/0000579f-7b35-4ed3-b44c-db2a538066fe","scryfall_uri":"https://scryfall.com/card/tsp/157/fury-sliver?utm_source=api","layout":"normal","highres_image":true,"image_status":"highres_scan","image_uris":{"small":"https://cards.scryfall.io/small/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.jpg?1562894979","normal":"https://cards.scryfall.io/normal/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.jpg?1562894979","large":"https://cards.scryfall.io/large/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.jpg?1562894979","png":"https://cards.scryfall.io/png/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.png?1562894979","art_crop":"https://cards.scryfall.io/art_crop/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.jpg?1562894979","border_crop":"https://cards.scryfall.io/border_crop/front/0/0/0000579f-7b35-4ed3-b44c-db2a538066fe.jpg?1562894979"},"mana_cost":"{5}{R}","cmc":6.0,"type_line":"Creature — Sliver","oracle_text":"All Sliver creatures have double strike.","power":"3","toughness":"3","colors":["R"],"color_identity":["R"],"keywords":[],"legalities":{"standard":"not_legal","future":"not_legal","historic":"not_legal","gladiator":"not_legal","pioneer":"not_legal","explorer":"not_legal","modern":"legal","legacy":"legal","pauper":"not_legal","vintage":"legal","penny":"not_legal","commander":"legal","oathbreaker":"legal","brawl":"not_legal","historicbrawl":"not_legal","alchemy":"not_legal","paupercommander":"restricted","duel":"legal","oldschool":"not_legal","premodern":"not_legal","predh":"legal"},"games":["paper","mtgo"],"reserved":false,"foil":true,"nonfoil":true,"finishes":["nonfoil","foil"],"oversized":false,"promo":false,"reprint":false,"variation":false,"set_id":"c1d109bc-ffd8-428f-8d7d-3f8d7e648046","set":"tsp","set_name":"Time Spiral","set_type":"expansion","set_uri":"https://api.scryfall.com/sets/c1d109bc-ffd8-428f-8d7d-3f8d7e648046","set_search_uri":"https://api.scryfall.com/cards/search?order=set&q=e%3Atsp&unique=prints","scryfall_set_uri":"https://scryfall.com/sets/tsp?utm_source=api","rulings_uri":"https://api.scryfall.com/cards/0000579f-7b35-4ed3-b44c-db2a538066fe/rulings","prints_search_uri":"https://api.scryfall.com/cards/search?order=released&q=oracleid%3A44623693-51d6-49ad-8cd7-140505caf02f&unique=prints","collector_number":"157","digital":false,"rarity":"uncommon","flavor_text":"\"A rift opened, and our arrows were abruptly stilled. To move was to push the world. But the sliver's claw still twitched, red wounds appeared in Thed's chest, and ribbons of blood hung in the air.\"\n—Adom Capashen, Benalish hero","card_back_id":"0aeebaf5-8c7d-4636-9e82-8c27447861f7","artist":"Paolo Parente","artist_ids":["d48dd097-720d-476a-8722-6a02854ae28b"],"illustration_id":"2fcca987-364c-4738-a75b-099d8a26d614","border_color":"black","frame":"2003","full_art":false,"textless":false,"booster":true,"story_spotlight":false,"edhrec_rank":6467,"penny_rank":11462,"prices":{"usd":"0.44","usd_foil":"4.48","usd_etched":null,"eur":"0.10","eur_foil":"0.25","tix":"0.03"},"related_uris":{"gatherer":"https://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=109722","tcgplayer_infinite_articles":"https://infinite.tcgplayer.com/search?contentMode=article&game=magic&partner=scryfall&q=Fury+Sliver&utm_campaign=affiliate&utm_medium=api&utm_source=scryfall","tcgplayer_infinite_decks":"https://infinite.tcgplayer.com/search?contentMode=deck&game=magic&partner=scryfall&q=Fury+Sliver&utm_campaign=affiliate&utm_medium=api&utm_source=scryfall","edhrec":"https://edhrec.com/route/?cc=Fury+Sliver"},"purchase_uris":{"tcgplayer":"https://www.tcgplayer.com/product/14240?page=1&utm_campaign=affiliate&utm_medium=api&utm_source=scryfall","cardmarket":"https://www.cardmarket.com/en/Magic/Products/Search?referrer=scryfall&searchString=Fury+Sliver&utm_campaign=card_prices&utm_medium=text&utm_source=scryfall","cardhoarder":"https://www.cardhoarder.com/cards/25527?affiliate_id=scryfall&ref=card-profile&utm_campaign=affiliate&utm_medium=card&utm_source=scryfall"}}]`
	reader := bytes.NewBufferString(testData)

	bulkReader, err := scryfall.NewBulkReader[scryfall.Card](reader)
	require.NoError(t, err)

	card, err := bulkReader.Next()
	require.NoError(t, err)

	assert.Equal(t, "44623693-51d6-49ad-8cd7-140505caf02f", card.OracleID)

	card, err = bulkReader.Next()
	assert.Equal(t, err, io.EOF)
	assert.Nil(t, card)
}

func TestBulkReader_Cards(t *testing.T) {
	t.Parallel()

	testFd, err := os.Open("test/cards.json")
	require.NoError(t, err)

	defer testFd.Close()

	bulkReader, err := scryfall.NewBulkReader[scryfall.Card](testFd)
	require.NoError(t, err)

	numCards := 0

	for {
		card, err := bulkReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)
		require.NotNil(t, card)

		numCards++
	}

	assert.Equal(t, 10, numCards)
}

func TestBulkReader_Rulings(t *testing.T) {
	t.Parallel()

	testFd, err := os.Open("test/rulings.json")
	require.NoError(t, err)

	defer testFd.Close()

	bulkReader, err := scryfall.NewBulkReader[scryfall.Ruling](testFd)
	require.NoError(t, err)

	numRulings := 0

	for {
		ruling, err := bulkReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)
		require.NotNil(t, ruling)

		numRulings++
	}

	assert.Equal(t, 10, numRulings)
}
