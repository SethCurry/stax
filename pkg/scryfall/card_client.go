package scryfall

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CardClient struct {
	client *http.Client
}

// Named searches for a card by its name.
// Returns an error if more than one card is found.
func (c *CardClient) Named(ctx context.Context, cardName string) (*Card, error) {
	var card Card

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/named", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("fuzzy", cardName)
	req.URL.RawQuery = q.Encode()

	err = doRequest(c.client, req, &card)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return &card, nil
}

type List[T any] struct {
	Object     Object   `json:"object"`
	Data       []T      `json:"data"`
	HasMore    bool     `json:"has_more"`
	NextPage   string   `json:"next_page"`
	TotalCards int      `json:"total_cards"`
	Warnings   []string `json:"warnings"`
}

type CardSearchPager struct {
	client   *http.Client
	nextPage string
	done     bool
}

func (c *CardSearchPager) HasMore() bool {
	return !c.done
}

func (c *CardSearchPager) Next(ctx context.Context) ([]Card, error) {
	if c.done {
		return nil, io.EOF
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.nextPage, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Card]

	err = doRequest(c.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	c.nextPage = list.NextPage
	c.done = !list.HasMore

	return list.Data, nil
}

// CardSearchOptions provide optional parameters for searching for cards.
type CardSearchOptions struct {
	Unique            string
	Order             string
	Direction         string
	IncludeExtras     bool
	IncludeVariations bool
}

func (c *CardClient) Search(ctx context.Context, query string, opts CardSearchOptions) (*CardSearchPager, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/search", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	if opts.Unique != "" {
		q.Add("unique", opts.Unique)
	}
	if opts.Order != "" {
		q.Add("order", opts.Order)
	}
	if opts.Direction != "" {
		q.Add("dir", opts.Direction)
	}

	req.URL.RawQuery = q.Encode()

	pager := &CardSearchPager{
		client:   c.client,
		nextPage: req.URL.String(),
	}

	return pager, nil
}

type autocompleteResponse struct {
	Data []string `json:"data"`
}

func (c *CardClient) Autocomplete(ctx context.Context, query string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/autocomplete", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()

	var autocomplete autocompleteResponse

	err = doRequest(c.client, req, &autocomplete)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return autocomplete.Data, nil
}

type RandomCardOptions struct {
	Query   string
	Face    string
	Version string
	Pretty  bool
}

func (c *CardClient) Random(ctx context.Context, opts RandomCardOptions) (*Card, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/random", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()

	if opts.Query != "" {
		q.Add("q", opts.Query)
	}

	if opts.Face != "" {
		q.Add("face", opts.Face)
	}

	if opts.Version != "" {
		q.Add("version", opts.Version)
	}

	if opts.Pretty {
		q.Add("pretty", "true")
	}

	req.URL.RawQuery = q.Encode()

	var card Card

	err = doRequest(c.client, req, &card)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return &card, nil
}

type CardIdentifier struct {
	ID              string `json:"id"`
	MTGOID          int    `json:"mtgo_id"`
	MultiverseID    int    `json:"multiverse_id"`
	OracleID        string `json:"oracle_id"`
	IllustrationID  string `json:"illustration_id"`
	Name            string `json:"name"`
	Set             string `json:"set"`
	CollectorNumber string `json:"collector_number"`
}

type collectionRequest struct {
	Identifiers []CardIdentifier `json:"identifiers"`
}

func (c *CardClient) Collection(ctx context.Context, identifiers []CardIdentifier) ([]Card, error) {
	marshalled, err := json.Marshal(collectionRequest{Identifiers: identifiers})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal identifiers: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/random", bytes.NewBuffer(marshalled))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Card]

	err = doRequest(c.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}
