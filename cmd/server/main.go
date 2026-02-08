package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"distributed_rate_limiter/internal/config"
	"distributed_rate_limiter/internal/limiter"
	"distributed_rate_limiter/internal/middleware"
)

func main() {
	cfg := config.Load()
	//local
	//mgr := limiter.NewManager(5,1)
	//redis
	mgr := limiter.NewRedisManager(cfg.RedisAddr, cfg.RateLimit)

	// this is the resource, the user wants to access and is protected by the rate limiter
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format("15:04:05")
		fmt.Fprintf(w, "Success! Request processed at %s \n", currentTime)
		fmt.Fprintln(w, "You are seeing this because you are in the rate limit!")
	})

	wrappedServer := middleware.NewRateLimiter(mgr, finalHandler)

	port := cfg.Port
	slog.Info("server starting on http://localhost", "port", port)
	slog.Info("try refreshing your browser quickly to trigger the limit...")

	err := http.ListenAndServe(port, wrappedServer)
	if err != nil {
		slog.Error("server failed to start", "error", err)
	}
}