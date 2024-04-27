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
	// The type of the object.  You can generally ignore
	// this, but it can be useful for debugging.
	Object Object `json:"object"`

	// The ID of the bulk data source as a UUID.
	ID   string `json:"id"`
	Type string `json:"type"`

	// URI is the URI for the browser view of this source.
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Size        int64  `json:"size"`
	UpdatedAt   string `json:"updated_at"`

	// DownloadURI is the URI to download a copy of the data
	// for this data source.
	DownloadURI     string `json:"download_uri"`
	ContentType     string `json:"content_type"`
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

// ErrFirstTokenNotDelim is returned when the first token in the JSON stream is not a delimeter.
var ErrFirstTokenNotDelim = errors.New("first token is not a delimeter")

// ErrFirstTokenNotOpenBracket is returned when the first token in the JSON stream is not an open bracket ([).
var ErrFirstTokenNotOpenBracket = errors.New("first token is not an open bracket")

// NewBulkReader creates a BulkReader[T] from a given io.Reader.
// The BulkReader will not close the underlying reader.
func NewBulkReader[T any](src io.Reader) (*BulkReader[T], error) {
	decoder := json.NewDecoder(src)

	firstToken, err := decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to read first token: %w", err)
	}

	firstDelim, ok := firstToken.(json.Delim)
	if !ok {
		return nil, ErrFirstTokenNotDelim
	}

	if firstDelim.String() != "[" {
		return nil, ErrFirstTokenNotOpenBracket
	}

	return &BulkReader[T]{
		decoder: decoder,
	}, nil
}

// BulkReader[T] is a generic reader for lists of data encoded as JSON.
type BulkReader[T any] struct {
	//nolint:structcheck
	decoder *json.Decoder
}

// Next returns the next item in the reader.  It returns io.EOF if there are no more items.
// It returns a non-EOF error if the next item could not be parsed.
func (b *BulkReader[T]) Next() (*T, error) {
	var ret T

	if !b.decoder.More() {
		return nil, io.EOF
	}

	if err := b.decoder.Decode(&ret); err != nil {
		return nil, fmt.Errorf("failed to parse JSON for bulk: %w", err)
	}

	return &ret, nil
}
