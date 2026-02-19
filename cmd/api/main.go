package main

import (
	"fmt"
	"net/http"
	"os"
	"rate-limiter/internal/middleware"

	"github.com/redis/go-redis/v9"
)

func main() {
	// 1. Config
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// 2. Dependencies
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// 3. Router & Middleware
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Distributed Rate Limiter Active</h1>")
	})

	// Wrap the handler
	mux.Handle("/", middleware.RateLimit(rdb)(finalHandler))

	// 4. Start
	fmt.Printf("🚀 Server running on :8080 (Redis: %s)\n", redisAddr)
	http.ListenAndServe(":8080", mux)
}
