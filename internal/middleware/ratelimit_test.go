package middleware

import (
	"testing"
	"net/http"
	"net/http/httptest"
)

// create mock for Limiter
type mockLimiter struct {
	allowed	bool
	lastIP	string
}

func (m *mockLimiter) Allow(ip string) bool {
	m.lastIP = ip
	return m.allowed
}

func TestRateLimiterAllowed(t *testing.T) {
	mockLmtr := &mockLimiter{allowed: true}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	newRateLmtr := NewRateLimiter(mockLmtr, next)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder() // captures what handler writes

	// request is allowed
	newRateLmtr.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", rec.Code)
	}
} 

func TestRateLimiterBlocked(t *testing.T) {
	mockLmtr := &mockLimiter{allowed: false}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	newRateLmtr := NewRateLimiter(mockLmtr, next)

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	// request is blocked
	newRateLmtr.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Errorf("Expected 429, got %d", rec.Code)
	}
}

func TestRateLimiterIPExtraction(t *testing.T) {
	mockLmtr := &mockLimiter{allowed: true}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	newRateLmtr := NewRateLimiter(mockLmtr, next)


	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	rec := httptest.NewRecorder()

	// request is allowed
	newRateLmtr.ServeHTTP(rec, req)

	if mockLmtr.lastIP != "192.168.1.1" {
		t.Errorf("Expected 192.168.1.1, got %s", mockLmtr.lastIP)
	}
}