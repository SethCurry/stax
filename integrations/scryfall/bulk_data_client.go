package scryfall

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type BulkDataClient struct {
	client *http.Client
}

var ErrUnrecognizedBulkDataType = errors.New("unrecognized bulk data type")

func getBulkDataSources(data []BulkDataSource) (*BulkDataSources, error) {
	var ret BulkDataSources

	for _, item := range data {
		dataCopy := item

		switch item.Type {
		case "oracle_cards":
			ret.OracleCards = &dataCopy
		case "unique_artwork":
			ret.UniqueArtwork = &dataCopy
		case "default_cards":
			ret.DefaultCards = &dataCopy
		case "all_cards":
			ret.AllCards = &dataCopy
		case "rulings":
			ret.Rulings = &dataCopy
		default:
			return nil, fmt.Errorf("%w: \"%s\"", ErrUnrecognizedBulkDataType, item.Type)
		}
	}

	return &ret, nil
}

// ListSources lists all available bulk data sources.
func (b *BulkDataClient) ListSources(ctx context.Context) (*BulkDataSources, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.scryfall.com/bulk-data", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request to get bulk data list: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request to get bulk data list: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var dataList BulkDataSourcesList

	err = json.Unmarshal(body, &dataList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bulk data list response: %w", err)
	}

	return getBulkDataSources(dataList.Data)
}

type BulkDataSource struct {
	Object      Object `json:"object"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
	//nolint:tagliatelle
	UpdatedAt string `json:"updated_at"`
	//nolint:tagliatelle
	DownloadURI string `json:"download_uri"`
	//nolint:tagliatelle
	ContentType string `json:"content_type"`
	//nolint:tagliatelle
	ContentEncoding string `json:"content_encoding"`
}

type BulkDataSourcesList struct {
	Object string `json:"object"`
	//nolint:tagliatelle
	HasMore bool             `json:"has_more"`
	Data    []BulkDataSource `json:"data"`
}

type BulkDataSources struct {
	// one entry per oracle ID
	OracleCards   *BulkDataSource
	UniqueArtwork *BulkDataSource
	// includes cards in English, or another language if English not available
	DefaultCards *BulkDataSource
	// all cards in all languages
	AllCards *BulkDataSource
	Rulings  *BulkDataSource
}
