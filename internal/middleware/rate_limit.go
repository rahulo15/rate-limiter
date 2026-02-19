package middleware

import (
	"fmt"
	"net/http"
	"os"
	"rate-limiter/internal/limiter"

	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.URL.Query().Get("user")
			if userID == "" {
				userID = r.RemoteAddr
			}

			podName, _ := os.Hostname()

			allowed := limiter.AllowRequest(rdb, userID, 5, 5.0/60.0)

			if allowed {
				fmt.Printf("[%s] ✅ Allowed request for %s\n", podName, userID)
				next.ServeHTTP(w, r)
			} else {
				fmt.Printf("[%s] 🚫 Rate Limited %s\n", podName, userID)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "Too many requests"}`))
			}
		})
	}
}
