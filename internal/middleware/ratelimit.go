package middleware

import (
	"distributed_rate_limiter/internal/limiter"
	"distributed_rate_limiter/internal/metrics"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type RateLimitMiddleware struct {
	manager limiter.Limiter
	next    http.Handler
}

func (m *RateLimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var err error
	host := r.Header.Get("X-Real-IP")
	if host == "" {
		host, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
	}
	slog.Info("rate limit check", "host", host)

	// check if tokens available for this ip
	result := m.manager.Allow(host)
	metrics.RequestDuration.Observe(time.Since(start).Seconds())
	if !result.Allowed {
		metrics.RequestsTotal.WithLabelValues("blocked").Inc()
		slog.Warn("blocked", "ip", host)
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
		w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt, 10))
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}
	metrics.RequestsTotal.WithLabelValues("allowed").Inc()

	w.Header().Set("X-RateLimit-Limit", strconv.Itoa(result.Limit))
	w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(result.Remaining))
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(result.ResetAt, 10))

	// test passed, serve next request
	m.next.ServeHTTP(w, r)
}

func NewRateLimiter(mgr limiter.Limiter, nextToCall http.Handler) http.Handler {
	return &RateLimitMiddleware{
		manager: mgr,
		next:    nextToCall,
	}
}

// wrapped version of the above
// func RateLimiter(m *limiter.Manager) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
// 			ip := r.RemoteAddr
// 			bucket := m.manager.GetBucket(ip)

// 			if !bucket.Allow() {
// 				http.Error(w, "Too many requests", http.StatusTooManyRequests)
// 				return
// 			}

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
