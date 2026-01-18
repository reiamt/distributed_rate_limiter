package limiter

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity	float64		// max tokens the bucket can hold
	rate		float64		// tokens added per sec
	tokens		float64		// current available tokens
	lastTick	time.Time	// last update
	mu			sync.Mutex
}

func NewTokenBucket(capacity, rate float64) *TokenBucket {
	return &TokenBucket{
		capacity:	capacity,
		rate: 		rate,
		tokens: 	capacity,
		lastTick:	time.Now(),
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastTick).Seconds()
	tokensToAdd := elapsed * tb.rate

	tb.tokens += tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
	tb.lastTick = now

	if tb.tokens > 1 {
		tb.tokens--
		return true
	}

	return false
}