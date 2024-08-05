package scryfall

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NewClient creates a new Client.
func NewClient(startingClient *http.Client) *Client {
	defaultMaxRetries := 5
	defaultMaxRequests := 10
	defaultWindow := time.Second
	defaultTimeoutSeconds := 30

	if startingClient == nil {
		startingClient = &http.Client{
			Timeout:       time.Second * time.Duration(defaultTimeoutSeconds),
			CheckRedirect: nil,
			Transport:     nil,
			Jar:           nil,
		}
	}

	transport := http.DefaultTransport
	if startingClient.Transport != nil {
		transport = startingClient.Transport
	}

	httpClient := &http.Client{
		Transport: &roundTripper{
			maxRetries: defaultMaxRetries,
			limiter:    newRateLimiter(defaultWindow, defaultMaxRequests),
			inner:      transport,
		},
		CheckRedirect: startingClient.CheckRedirect,
		Jar:           startingClient.Jar,
		Timeout:       startingClient.Timeout,
	}

	return &Client{
		Card:     &CardClient{client: httpClient},
		BulkData: &BulkDataClient{client: httpClient},
		Rulings:  &RulingClient{client: httpClient},
	}
}

// Client is the interface to Scryfall's API.
type Client struct {
	Card     *CardClient
	BulkData *BulkDataClient
	Rulings  *RulingClient
}

type APIError struct {
	// The HTTP status code of the response
	Status int `json:"status"`

	// A machine-friendly code for the error condition
	Code string `json:"code"`

	Details string `json:"details"`

	// A computer-friendly string specifying more details about the error condition.
	// E.g. for a 404 it might return "ambiguous" if the request refers to multiple
	// cards potentially.
	//
	// This field can be empty.
	Type string `json:"type"`

	// A series of human-readable errors.
	Warnings []string `json:"warnings"`
}

func (a APIError) Error() string {
	return fmt.Sprintf("API Error: %s: %s", a.Code, strings.Join(a.Warnings, " | "))
}

func doRequest(client *http.Client, req *http.Request, into interface{}) error {
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := decoder.Decode(&apiErr); err != nil {
			return fmt.Errorf("failed to decode respone: %w", err)
		}

		return &apiErr
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := decoder.Decode(into); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
