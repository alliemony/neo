package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateJWT_AndValidate(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	user := OAuthUser{
		ID:        "12345",
		Username:  "testuser",
		AvatarURL: "https://example.com/avatar.png",
		Provider:  "github",
	}

	tokenStr, err := CreateJWT(secret, user, 24*time.Hour)
	if err != nil {
		t.Fatalf("CreateJWT error: %v", err)
	}

	claims, err := ValidateJWT(tokenStr, secret)
	if err != nil {
		t.Fatalf("ValidateJWT error: %v", err)
	}

	if claims["sub"] != "12345" {
		t.Errorf("sub = %v, want 12345", claims["sub"])
	}
	if claims["username"] != "testuser" {
		t.Errorf("username = %v, want testuser", claims["username"])
	}
	if claims["avatar_url"] != "https://example.com/avatar.png" {
		t.Errorf("avatar_url = %v, want https://example.com/avatar.png", claims["avatar_url"])
	}
	if claims["provider"] != "github" {
		t.Errorf("provider = %v, want github", claims["provider"])
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	user := OAuthUser{ID: "1", Username: "u", AvatarURL: "", Provider: "github"}

	tokenStr, err := CreateJWT(secret, user, -1*time.Hour) // already expired
	if err != nil {
		t.Fatalf("CreateJWT error: %v", err)
	}

	_, err = ValidateJWT(tokenStr, secret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	secret := []byte("correct-secret-key-32-bytes!!!")
	wrongSecret := []byte("wrong-secret-key-32-bytes!!!xx")
	user := OAuthUser{ID: "1", Username: "u", AvatarURL: "", Provider: "github"}

	tokenStr, err := CreateJWT(secret, user, 24*time.Hour)
	if err != nil {
		t.Fatalf("CreateJWT error: %v", err)
	}

	_, err = ValidateJWT(tokenStr, wrongSecret)
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_TamperedToken(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	user := OAuthUser{ID: "1", Username: "u", AvatarURL: "", Provider: "github"}

	tokenStr, err := CreateJWT(secret, user, 24*time.Hour)
	if err != nil {
		t.Fatalf("CreateJWT error: %v", err)
	}

	// Tamper with the payload
	tampered := tokenStr[:len(tokenStr)-5] + "xxxxx"
	_, err = ValidateJWT(tampered, secret)
	if err == nil {
		t.Error("expected error for tampered token, got nil")
	}
}

func TestOAuthMiddleware_ValidSession(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	oauth := NewOAuthAuthenticator(OAuthConfig{SessionSecret: secret})
	user := OAuthUser{ID: "42", Username: "admin", AvatarURL: "https://example.com/a.png", Provider: "github"}

	tokenStr, _ := CreateJWT(secret, user, 24*time.Hour)

	handler := oauth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := CurrentUser(r)
		if u == nil {
			t.Error("expected user in context, got nil")
			return
		}
		if u.Username != "admin" {
			t.Errorf("username = %q, want admin", u.Username)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/v1/admin/posts", nil)
	req.AddCookie(&http.Cookie{Name: "neo_session", Value: tokenStr})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestOAuthMiddleware_MissingCookie(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	oauth := NewOAuthAuthenticator(OAuthConfig{SessionSecret: secret})

	handler := oauth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/v1/admin/posts", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestOAuthMiddleware_ExpiredSession(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	oauth := NewOAuthAuthenticator(OAuthConfig{SessionSecret: secret})
	user := OAuthUser{ID: "1", Username: "u", AvatarURL: "", Provider: "github"}

	tokenStr, _ := CreateJWT(secret, user, -1*time.Hour) // expired

	handler := oauth.Middleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	req := httptest.NewRequest("GET", "/api/v1/admin/posts", nil)
	req.AddCookie(&http.Cookie{Name: "neo_session", Value: tokenStr})
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestIsUserAllowed(t *testing.T) {
	allowed := []string{"alice", "bob", "charlie"}

	if !IsUserAllowed("alice", allowed) {
		t.Error("expected alice to be allowed")
	}
	if !IsUserAllowed("Bob", allowed) {
		t.Error("expected Bob (case-insensitive) to be allowed")
	}
	if IsUserAllowed("stranger", allowed) {
		t.Error("expected stranger to not be allowed")
	}
	if IsUserAllowed("", allowed) {
		t.Error("expected empty string to not be allowed")
	}
}

func TestGenerateState(t *testing.T) {
	s1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState error: %v", err)
	}
	s2, _ := GenerateState()

	if len(s1) != 32 {
		t.Errorf("state length = %d, want 32 hex chars", len(s1))
	}
	if s1 == s2 {
		t.Error("two generated states should be different")
	}
}
