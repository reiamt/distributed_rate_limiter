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
	circuitBreaker := NewCircuitBreaker(3, 30*time.Second)

	for cntr := 0; cntr < circuitBreaker.threshold; cntr++ {
		circuitBreaker.RecordResult(false)
	}
	allowed := circuitBreaker.Allow()
	if allowed {
		t.Errorf("Request should not be allowed")
	}
}
