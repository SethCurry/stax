package scryfall

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

func NewRateLimiter(window time.Duration, max int) *RateLimiter {
	return &RateLimiter{
		window: window,
		max:    max,
		events: []time.Time{},
		lock:   sync.RWMutex{},
	}
}

type RateLimiter struct {
	window time.Duration
	max    int
	events []time.Time
	lock   sync.RWMutex
}

func (r *RateLimiter) AddEvent() bool {
	r.Clean()
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.events) >= r.max {
		return false
	}

	r.events = append(r.events, time.Now())

	return true
}

func (r *RateLimiter) Clean() {
	r.lock.Lock()
	defer r.lock.Unlock()

	now := time.Now()

	indexesToRemove := []int{}

	for i, event := range r.events {
		if event.Before(now.Add(-r.window)) {
			indexesToRemove = append(indexesToRemove, i)
		}
	}

	for i, idx := range indexesToRemove {
		r.events = append(r.events[:idx-i], r.events[idx-i+1:]...)
	}
}

var ErrTimeoutFromLimiter = errors.New("timed out while waiting for available request in rate limiter")

func NewRoundTripperError(inner error) *RoundTripperError {
	return &RoundTripperError{Inner: inner}
}

type RoundTripperError struct {
	Inner error
}

func (r *RoundTripperError) Error() string {
	return "round tripper failed: " + r.Inner.Error()
}

type roundTripper struct {
	inner      http.RoundTripper
	limiter    *RateLimiter
	maxRetries int
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	numAttempts := 0
	lastSleep := time.Second

	for numAttempts < r.maxRetries {
		numAttempts++

		ok := r.limiter.AddEvent()
		if !ok {
			time.Sleep(lastSleep)
			lastSleep *= 2

			continue
		}

		resp, err := r.inner.RoundTrip(req)
		if err != nil {
			return nil, NewRoundTripperError(err)
		}

		return resp, nil
	}

	return nil, ErrTimeoutFromLimiter
}

// NewClient creates a new Client.
func NewClient(startingClient *http.Client) *Client {
	defaultMaxRetries := 5
	defaultMaxRequests := 5
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
			limiter:    NewRateLimiter(defaultWindow, defaultMaxRequests),
			inner:      transport,
		},
		CheckRedirect: startingClient.CheckRedirect,
		Jar:           startingClient.Jar,
		Timeout:       startingClient.Timeout,
	}

	return &Client{
		Card:     &CardClient{client: httpClient},
		BulkData: &BulkDataClient{client: httpClient},
	}
}

// Client is the interface to Scryfall's API.
type Client struct {
	Card     *CardClient
	BulkData *BulkDataClient
}
