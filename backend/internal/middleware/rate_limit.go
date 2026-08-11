package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lanxre/kyokusulib/internal/utils/response"
	"github.com/redis/go-redis/v9"
)

var rateLimitScript = redis.NewScript(`
	local current = redis.call('INCR', KEYS[1])
	if current == 1 then
	    redis.call('PEXPIRE', KEYS[1], ARGV[1])
	end
	return current
`)

func RateLimitMiddleware(redis *redis.Client, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}

			key := fmt.Sprintf("rate_limit:%s:%s", r.URL.Path, ip)
			ctx := r.Context()

			count, err := rateLimitScript.Run(ctx, redis, []string{key}, window.Milliseconds()).Int()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if count > limit {
				response.Error(w, http.StatusTooManyRequests, "Too many requests. Please try again later.")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
