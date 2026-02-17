package limiter

import (
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu           sync.Mutex
	state        string
	failures     int
	threshold    int
	resetTimeout time.Duration
	lastFailure  time.Time
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.state == "closed" || cb.state == "half-open" {
		return true
	}
	if cb.state == "open" {
		if time.Since(cb.lastFailure) >= cb.resetTimeout {
			cb.state = "half-open"
			return true
		} else {
			return false
		}
	}
	return false
}

func (cb *CircuitBreaker) RecordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if success {
		cb.state = "closed"
		cb.failures = 0
	} else {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.threshold {
			cb.state = "open"
		}
	}
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:        "closed",
		threshold:    threshold,
		resetTimeout: resetTimeout,
	}
}
