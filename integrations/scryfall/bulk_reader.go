package scryfall

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

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
