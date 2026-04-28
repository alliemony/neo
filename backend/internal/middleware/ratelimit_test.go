package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_AllowsNormalUsage(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 5)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRateLimiter_BlocksRapidSubmissions(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 2)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for i := range 2 {
		req := httptest.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("request %d: status = %d, want %d", i, w.Code, http.StatusOK)
		}
	}

	// Third request should be blocked.
	req := httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("status = %d, want %d", w.Code, http.StatusTooManyRequests)
	}
}

func TestRateLimiter_DifferentIPsAreIndependent(t *testing.T) {
	rl := NewRateLimiter(time.Minute, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First IP uses its quota.
	req := httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("first IP: status = %d, want %d", w.Code, http.StatusOK)
	}

	// Second IP should still be allowed.
	req = httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.2:5678"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("second IP: status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	rl := NewRateLimiter(50*time.Millisecond, 1)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("first request: status = %d, want %d", w.Code, http.StatusOK)
	}

	// Wait for window to expire.
	time.Sleep(60 * time.Millisecond)

	req = httptest.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("after window: status = %d, want %d", w.Code, http.StatusOK)
	}
}
