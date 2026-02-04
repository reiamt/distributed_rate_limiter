package main

import (
	"fmt"
	"net/http"
	"time"

	"distributed_rate_limiter/internal/limiter"
	"distributed_rate_limiter/internal/middleware"
)

func main() {
	//local
	//mgr := limiter.NewManager(5,1)
	//redis
	mgr := limiter.NewRedisManager("localhost:6379", 5)

	// this is the resource, the user wants to access and is protected by the rate limiter
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().Format("15:04:05")
		fmt.Fprintf(w, "Success! Request processed at %s \n", currentTime)
		fmt.Fprintln(w, "You are seeing this because you are in the rate limit!")
	})

	wrappedServer := middleware.NewRateLimiter(mgr, finalHandler)

	port := ":8080"
	fmt.Printf("Server is starting on http://localhost%s\n", port)
	fmt.Println("Try refreshing your browser quickly to trigger the limit...")

	err := http.ListenAndServe(port, wrappedServer)
	if err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}