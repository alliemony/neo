package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter returns middleware that limits requests per IP.
// window is how long to remember a request; if a second request arrives within the window, it's rejected.
func RateLimiter(window time.Duration) func(http.Handler) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]time.Time)

	// Periodically clean up expired entries.
	go func() {
		for {
			time.Sleep(window * 2)
			mu.Lock()
			now := time.Now()
			for ip, t := range clients {
				if now.Sub(t) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr

			mu.Lock()
			last, exists := clients[ip]
			if exists && time.Since(last) < window {
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"too many requests"}`))
				return
			}
			clients[ip] = time.Now()
			mu.Unlock()

			next.ServeHTTP(w, r)
		})
	}
}
