package scryfall

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

func newRoundTripperError(inner error) *RoundTripperError {
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
	limiter    *rateLimiter
	maxRetries int
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	numAttempts := 0
	lastSleep := time.Second

	for numAttempts <= r.maxRetries {
		numAttempts++

		ok := r.limiter.AddEvent()
		if !ok {
			time.Sleep(lastSleep)
			lastSleep *= 2

			continue
		}

		resp, err := r.inner.RoundTrip(req)
		if err != nil {
			return nil, newRoundTripperError(err)
		}

		return resp, nil
	}

	return nil, ErrTimeoutFromLimiter
}

var ErrTimeoutFromLimiter = errors.New("timed out while waiting for available request in rate limiter")

func newRateLimiter(window time.Duration, max int) *rateLimiter {
	return &rateLimiter{
		window: window,
		max:    max,
		events: []time.Time{},
		lock:   sync.RWMutex{},
	}
}

type rateLimiter struct {
	window time.Duration
	max    int
	events []time.Time
	lock   sync.RWMutex
}

func (r *rateLimiter) AddEvent() bool {
	r.Clean()
	r.lock.Lock()
	defer r.lock.Unlock()

	if len(r.events) >= r.max {
		return false
	}

	r.events = append(r.events, time.Now())

	return true
}

func (r *rateLimiter) Clean() {
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
