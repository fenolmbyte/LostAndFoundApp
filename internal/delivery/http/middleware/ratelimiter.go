package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func RateLimitByUserID(redisClient *redis.Client, maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			key := fmt.Sprintf("rate_limit:%s:%s", r.URL.Path, authHeader)

			count, err := redisClient.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				http.Error(w, "internal rate limiter error", http.StatusInternalServerError)
				return
			}

			if count >= maxRequests {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(map[string]string{
					"error":   "rate_limit_exceeded",
					"message": "Too many requests. Please slow down.",
				})
				return
			}

			pipe := redisClient.TxPipeline()
			pipe.Incr(ctx, key)
			pipe.Expire(ctx, key, window)
			_, _ = pipe.Exec(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
