package limiter

import (
	"testing"
	"time"
)

func TestNewCircuitBreakerAllowed(t *testing.T) {
	circuitBreaker := NewCircuitBreaker(3, 30*time.Second)

	allowed := circuitBreaker.Allow()
	if !allowed {
		t.Errorf("Request should be allowed.")
	}
}

func TestOpenAfterThreshold(t *testing.T) {
	circuitBreaker := NewCircuitBreaker(3, 100*time.Millisecond)

	for cntr := 0; cntr < circuitBreaker.threshold; cntr++ {
		circuitBreaker.RecordResult(false)
	}
	allowed := circuitBreaker.Allow()
	if allowed {
		t.Errorf("Request should not be allowed")
	}
}

func TestHalfOpenAfterTimeout(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	for i := 0; i < cb.threshold; i++ {
		cb.RecordResult(false)
	}
	time.Sleep(101 * time.Millisecond)
	if !cb.Allow() {
		t.Errorf("Request should be allowed.")
	}
}

func TestClosesOnSuccess(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	for i := 0; i < cb.threshold; i++ {
		cb.RecordResult(false)
	}
	time.Sleep(101 * time.Millisecond)
	// now half open
	cb.RecordResult(true)
	// closed
	if !cb.Allow() {
		t.Errorf("Request should be allowed.")
	}
}

func TestReopensAfterFailureInHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	for i := 0; i < cb.threshold; i++ {
		cb.RecordResult(false)
	}
	time.Sleep(101 * time.Millisecond)
	cb.RecordResult(false)
	if cb.Allow() {
		t.Errorf("Request should not be allowed.")
	}
}
