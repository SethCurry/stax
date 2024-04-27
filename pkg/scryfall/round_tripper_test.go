package scryfall

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RateLimiter(t *testing.T) {
	t.Parallel()

	limiter := newRateLimiter(time.Second, 1)

	for i := 0; i < 5; i++ {
		ok := limiter.AddEvent()
		if i == 0 {
			assert.True(t, ok)
		} else {
			assert.False(t, ok)
		}
	}

	time.Sleep(time.Second)

	ok := limiter.AddEvent()
	assert.True(t, ok)
	ok = limiter.AddEvent()
	assert.False(t, ok)
}
