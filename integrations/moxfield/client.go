package moxfield

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func NewClient(baseClient *http.Client) *Client {
	return &Client{
		httpClient: baseClient,
	}
}

type Client struct {
	httpClient *http.Client
}

type ListUserDecksRequest struct {
	PageSize   int
	PageNumber int
}

type UserDeck struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HasPrimer bool   `json:"hasPrimer"`
	Format    string `json:"format"`
	PublicURL string `json:"publicUrl"`
	PublicID  string `json:"publicId"`
	IsLegal   bool   `json:"isLegal"`
}

type ListUserDecksResponse struct {
	PageNumber   int        `json:"pageNumber"`
	PageSize     int        `json:"pageSize"`
	TotalPages   int        `json:"totalPages"`
	TotalResults int        `json:"totalResults"`
	Data         []UserDeck `json:"data"`
}

type GetDeckResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ExportID string `json:"exportId"`
}

func (c *Client) GetDeck(deckID string) (*GetDeckResponse, error) {
	reqUrl := fmt.Sprintf("https://api2.moxfield.com/v3/decks/all/%s", deckID)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

func (c *Client) GetDeckList(deckID string) ([]byte, error) {
	deck, err := c.GetDeck(deckID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deck: %w", err)
	}

	return c.DownloadDeckList(deckID, deck.ExportID)
}

func (c *Client) DownloadDeckList(deckID, exportId string) ([]byte, error) {
	reqUrl := fmt.Sprintf("https://api2.moxfield.com/v2/decks/all/%s/export", deckID)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()
	q.Add("arenaOnly", "false")
	q.Add("format", "full")
	q.Add("includeFinish", "true")
	q.Add("exportId", exportId)

	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

func (c *Client) ListUserDecks(username string, options ListUserDecksRequest) (*ListUserDecksResponse, error) {
	reqUrl := fmt.Sprintf("https://api2.moxfield.com/v2/users/%s/decks", username)

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	q := req.URL.Query()

	pageSize := 12
	if options.PageSize > 0 {
		pageSize = options.PageSize
	}
	q.Add("pageSize", strconv.Itoa(pageSize))

	pageNum := 1
	if options.PageNumber > 0 {
		pageNum = options.PageNumber
	}
	q.Add("pageNumber", strconv.Itoa(pageNum))

	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
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

func ParseDeckListLine(line string) (*DeckListLine, error) {
	ssv := strings.Split(line, " ")

	if len(ssv) < 3 {
		return nil, fmt.Errorf("invalid deck list line: %s", line)
	}

	offset := 0

	quantityStr := ssv[0]
	collectorNumber := ssv[len(ssv)-1]
	if strings.ToLower(collectorNumber) == "*f*" {
		offset = 1
		collectorNumber = ssv[len(ssv)-2]
	}

	collectorNumber = strings.Replace(collectorNumber, "★", "*", 0)
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
