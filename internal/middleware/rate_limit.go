package middleware

import (
	"app_aggregator/internal/ratelimit"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"strings"
)

func RateLimitWithLogger(cfg *ratelimit.Config, logger *slog.Logger) func(http.Handler) http.Handler {
	limiter := ratelimit.NewSlidingWindowLimiter(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)
			if !limiter.IsAllowed(ip) {
				logger.Warn("Request blocked by rate limiter",
					slog.String("ip", ip),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.String("user_agent", r.UserAgent()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(cfg.BlockDuration.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)

				resp := map[string]interface{}{
					"error":       "Rate limit exceeded",
					"retry_after": int(cfg.BlockDuration.Seconds()),
					"limit":       cfg.RequestsPerMinute,
					"window_size": cfg.WindowSize.String(),
					"stats":       limiter.GetStats(),
					"ip":          ip,
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			remaining := limiter.GetRemaining(ip)
			if remaining <= 10 {
				logger.Warn("Rate limit is close to the limit",
					slog.String("ip", ip),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("remaining", remaining),
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		ips := strings.Split(ip, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func getKey(r *http.Request) string {
	return getIP(r)
}
