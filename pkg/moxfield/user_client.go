package moxfield

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type UserClient struct {
	client *http.Client
}

func (u *UserClient) ListUserDecks(
	ctx context.Context,
	username string,
	options ListUserDecksRequest,
) (*ListUserDecksResponse, error) {
	reqURL := fmt.Sprintf("https://api2.moxfield.com/v2/users/%s/decks", username)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	query := req.URL.Query()

	pageSize := 12
	if options.PageSize > 0 {
		pageSize = options.PageSize
	}

	query.Add("pageSize", strconv.Itoa(pageSize))

	pageNum := 1
	if options.PageNumber > 0 {
		pageNum = options.PageNumber
	}

	query.Add("pageNumber", strconv.Itoa(pageNum))

	req.URL.RawQuery = query.Encode()

	resp, err := u.client.Do(req)
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

	var ret ListUserDecksResponse

	err = json.Unmarshal(body, &ret)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &ret, nil
}

type UserDeck struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Format    string `json:"format"`
	PublicURL string `json:"publicUrl"`
	PublicID  string `json:"publicId"`
	HasPrimer bool   `json:"hasPrimer"`
	IsLegal   bool   `json:"isLegal"`
}

type ListUserDecksRequest struct {
	PageSize   int
	PageNumber int
}

type ListUserDecksResponse struct {
	PageNumber   int        `json:"pageNumber"`
	PageSize     int        `json:"pageSize"`
	TotalPages   int        `json:"totalPages"`
	TotalResults int        `json:"totalResults"`
	Data         []UserDeck `json:"data"`
}
