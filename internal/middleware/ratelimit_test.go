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

func TestRateLimiter(t *testing.T) {
	tests := []struct {
		name		string
		allowed		bool
		remoteAddr	string
		wantCode	int
		wantIP		string
	}{
		{
			name:		"allowed request returns 200",
			allowed:	true,
			remoteAddr: "192.168.1.1:12345",
			wantCode:	http.StatusOK,
			wantIP:		"192.168.1.1",
		},
		{
			name:		"blocked request returns 429",
			allowed:	false,
			remoteAddr: "192.168.1.1:12345",
			wantCode:	http.StatusTooManyRequests,
			wantIP:		"192.168.1.1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockLmtr := &mockLimiter{allowed: tc.allowed}
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			newRateLmtr := NewRateLimiter(mockLmtr, next)

			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tc.remoteAddr
			rec := httptest.NewRecorder() //captures what handler writes
			
			newRateLmtr.ServeHTTP(rec, req)
			if rec.Code != tc.wantCode {
				t.Errorf("Expected %d, got %d", tc.wantCode, rec.Code)
			}
			
			if mockLmtr.lastIP != tc.wantIP {
				t.Errorf("Expected %s, got %s", tc.wantIP, mockLmtr.lastIP)
			}
		})
	}
}