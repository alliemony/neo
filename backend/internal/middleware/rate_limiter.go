package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"
)

// rateLimiter tracks request timestamps per IP using a sliding window.
type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string][]time.Time
	window time.Duration
	limit  int
}

// NewRateLimiter creates a rate limiter that allows limit requests per window duration.
func NewRateLimiter(window time.Duration, limit int) *rateLimiter {
	return &rateLimiter{
		hits:   make(map[string][]time.Time),
		window: window,
		limit:  limit,
	}
}

// Middleware returns an HTTP handler that enforces the rate limit.
func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r.RemoteAddr)
		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "too many requests"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	existing := rl.hits[ip]
	var valid []time.Time
	for _, t := range existing {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.hits[ip] = valid
		return false
	}

	rl.hits[ip] = append(valid, now)
	return true
}

func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

// RateLimiter returns middleware using a sliding window with a default limit of 10 requests.
// Use NewRateLimiter for configurable limits.
func RateLimiter(window time.Duration) func(http.Handler) http.Handler {
	rl := NewRateLimiter(window, 10)
	return rl.Middleware
}
