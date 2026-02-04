package middleware

import (
	"fmt"
	"net"
	"net/http"
	"distributed_rate_limiter/internal/limiter"
)

type RateLimitMiddleware struct {
	manager	limiter.Limiter
	next	http.Handler
}

func (m *RateLimitMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	} 
	fmt.Printf("Rate Limit Check | Original: %s | Key Used: %s\n", r.RemoteAddr, host)

	// check if tokens available for this ip
	if !m.manager.Allow(host) {
		fmt.Printf(">>> BLOCKED: %s\n", host)
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	// test passed, serve next request
	m.next.ServeHTTP(w, r)
}

func NewRateLimiter(mgr limiter.Limiter, nextToCall http.Handler) http.Handler {
	return &RateLimitMiddleware{
		manager:	mgr,
		next:		nextToCall,	
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