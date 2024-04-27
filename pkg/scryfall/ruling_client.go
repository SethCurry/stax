package scryfall

import (
	"context"
	"fmt"
	"net/http"
)

type RulingClient struct {
	client *http.Client
}

// ByScryfallID fetches all of the rulings for a card by its Scryfall ID.
//
// https://scryfall.com/docs/api/rulings/id
func (r *RulingClient) ByScryfallID(ctx context.Context, id string) ([]Ruling, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/"+id+"/rulings", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Ruling]

	err = doRequest(r.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}

func (r *RulingClient) ByMultiverseID(ctx context.Context, id int) ([]Ruling, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/multiverse/"+fmt.Sprint(id)+"/rulings", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Ruling]

	err = doRequest(r.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}

func (r *RulingClient) ByMTGOID(ctx context.Context, id int) ([]Ruling, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/mtgo/"+fmt.Sprint(id)+"/rulings", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Ruling]

	err = doRequest(r.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}

func (r *RulingClient) ByArenaID(ctx context.Context, id int) ([]Ruling, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/arena/"+fmt.Sprint(id)+"/rulings", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Ruling]

	err = doRequest(r.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}

func (r *RulingClient) ByCodeAndNumber(ctx context.Context, code string, number string) ([]Ruling, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/cards/"+code+"/"+number+"/rulings", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	var list List[Ruling]

	err = doRequest(r.client, req, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return list.Data, nil
}
