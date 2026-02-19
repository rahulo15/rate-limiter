package middleware

import (
	"net/http"
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

			if !limiter.AllowRequest(rdb, userID, 5, 5.0/60.0) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "Too many requests"}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
