package scryfall

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CardClient struct {
	client *http.Client
}

func (c *CardClient) Named(ctx context.Context, cardName string) (*Card, error) {
	var card Card

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/named", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("fuzzy", cardName)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response: %w", err)
	}

	err = json.Unmarshal(bodyBytes, &card)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &card, nil
}
