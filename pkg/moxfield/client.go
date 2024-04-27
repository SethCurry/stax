package moxfield

import (
	"errors"
	"net/http"
)

// ErrUnexpectedStatusCode is returned when the server
// returns a status code that is not 200.
var ErrUnexpectedStatusCode = errors.New("unexpected status code")

// NewClient creates a new *Client.
// If baseClient is nil, http.DefaultClient is used.
func NewClient(baseClient *http.Client) *Client {
	if baseClient == nil {
		baseClient = http.DefaultClient
	}

	return &Client{
		Decks: &DeckClient{
			client: baseClient,
		},
		Users: &UserClient{
			client: baseClient,
		},
	}
}

// Client is the root client that other resource-specific clients
// are attached to.
type Client struct {
	Decks *DeckClient
	Users *UserClient
}
