package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return string(hash)
}

func TestBasicAuthImplementsAuthenticator(t *testing.T) {
	var _ Authenticator = NewBasicAuth("admin", "hash")
}

func TestBasicAuth_ValidCredentials(t *testing.T) {
	hash := hashPassword(t, "secret")
	auth := NewBasicAuth("admin", hash)

	handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	req.SetBasicAuth("admin", "secret")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Fatalf("expected body 'ok', got %q", rr.Body.String())
	}
}

func TestBasicAuth_MissingCredentials(t *testing.T) {
	hash := hashPassword(t, "secret")
	auth := NewBasicAuth("admin", hash)

	handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
	if rr.Header().Get("WWW-Authenticate") != `Basic realm="admin"` {
		t.Fatalf("expected WWW-Authenticate header, got %q", rr.Header().Get("WWW-Authenticate"))
	}
}

func TestBasicAuth_InvalidCredentials(t *testing.T) {
	hash := hashPassword(t, "secret")
	auth := NewBasicAuth("admin", hash)

	handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	req.SetBasicAuth("admin", "wrongpassword")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestBasicAuth_WrongUsername(t *testing.T) {
	hash := hashPassword(t, "secret")
	auth := NewBasicAuth("admin", hash)

	handler := auth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/admin/posts", nil)
	req.SetBasicAuth("hacker", "secret")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}
