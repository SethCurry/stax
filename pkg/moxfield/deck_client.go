package moxfield

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type DeckClient struct {
	client *http.Client
}

func (d *DeckClient) Get(ctx context.Context, deckID string) (*GetDeckResponse, error) {
	reqURL := "https://api2.moxfield.com/v3/decks/all/" + deckID

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var ret GetDeckResponse

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &ret, nil
}

func (d *DeckClient) GetDeckList(ctx context.Context, deckID string) ([]DeckListLine, error) {
	deck, err := d.Get(ctx, deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	deckBody, err := d.downloadDeckList(ctx, deckID, deck.ExportID)
	if err != nil {
		return nil, fmt.Errorf("failed to download deck list: %w", err)
	}

	deckLines := strings.Split(string(deckBody), "\n")

	var ret []DeckListLine

	for _, line := range deckLines {
		if line == "" {
			continue
		}

		parsed, err := parseDeckListLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed trying to parse deck list line: %w", err)
		}

		ret = append(ret, *parsed)
	}

	return ret, nil
}

func (d *DeckClient) downloadDeckList(ctx context.Context, deckID, exportID string) ([]byte, error) {
	reqURL := fmt.Sprintf("https://api2.moxfield.com/v2/decks/all/%s/export", deckID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	query := req.URL.Query()
	query.Add("arenaOnly", "false")
	query.Add("format", "full")
	query.Add("includeFinish", "true")
	query.Add("exportId", exportID)

	req.URL.RawQuery = query.Encode()

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

func parseDeckListLine(line string) (*DeckListLine, error) {
	ssv := strings.Split(line, " ")

	if len(ssv) < 3 {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedStatusCode, line)
	}

	offset := 0

	quantityStr := ssv[0]
	collectorNumber := ssv[len(ssv)-1]

	if strings.ToLower(collectorNumber) == "*f*" {
		offset = 1
		collectorNumber = ssv[len(ssv)-2]
	}

	setCode := ssv[len(ssv)-2-offset]
	setCode = strings.Trim(setCode, "()")
	name := strings.Join(ssv[1:len(ssv)-2-offset], " ")

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %w", err)
	}

	return &DeckListLine{
		Quantity:        quantity,
		Name:            name,
		Set:             setCode,
		CollectorNumber: collectorNumber,
	}, nil
}

type DeckListLine struct {
	Quantity        int
	Name            string
	Set             string
	CollectorNumber string
}

type GetDeckResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ExportID string `json:"exportId"`
}
