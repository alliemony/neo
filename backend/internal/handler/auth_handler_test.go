package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/alliemony/neo/backend/internal/middleware"
)

func testOAuthConfig(secret []byte) middleware.OAuthConfig {
	return middleware.OAuthConfig{
		ClientID:      "test-client-id",
		ClientSecret:  "test-client-secret",
		RedirectURL:   "http://localhost:8080/api/v1/auth/callback",
		AllowedUsers:  []string{"admin", "allowed-user"},
		SessionSecret: secret,
	}
}

func setupAuthHandler(t *testing.T, secret []byte) *AuthHandler {
	t.Helper()
	cfg := testOAuthConfig(secret)
	oauth := middleware.NewOAuthAuthenticator(cfg)
	return NewAuthHandler(oauth)
}

func TestAuthHandler_Login_RedirectsToGitHub(t *testing.T) {
	h := setupAuthHandler(t, []byte("test-secret-key-32-bytes-long!!!"))

	r := chi.NewRouter()
	r.Get("/api/v1/auth/login", h.Login)

	req := httptest.NewRequest("GET", "/api/v1/auth/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusFound)
	}

	location := w.Header().Get("Location")
	if location == "" {
		t.Fatal("Location header is empty")
	}

	// Should redirect to GitHub authorize URL
	if !contains(location, "github.com/login/oauth/authorize") {
		t.Errorf("location = %q, should contain github authorize URL", location)
	}
	if !contains(location, "client_id=test-client-id") {
		t.Errorf("location = %q, should contain client_id", location)
	}
	if !contains(location, "scope=read:user") {
		t.Errorf("location = %q, should contain scope", location)
	}
	if !contains(location, "state=") {
		t.Errorf("location = %q, should contain state", location)
	}

	// Should set oauth_state cookie
	cookies := w.Result().Cookies()
	var stateCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "oauth_state" {
			stateCookie = c
			break
		}
	}
	if stateCookie == nil {
		t.Fatal("oauth_state cookie not set")
	}
	if !stateCookie.HttpOnly {
		t.Error("oauth_state cookie should be HttpOnly")
	}
}

func TestAuthHandler_Callback_StateMismatch(t *testing.T) {
	h := setupAuthHandler(t, []byte("test-secret-key-32-bytes-long!!!"))

	r := chi.NewRouter()
	r.Get("/api/v1/auth/callback", h.Callback)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=abc&state=wrong", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "correct"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "invalid state parameter" {
		t.Errorf("error = %q, want 'invalid state parameter'", resp["error"])
	}
}

func TestAuthHandler_Callback_MissingState(t *testing.T) {
	h := setupAuthHandler(t, []byte("test-secret-key-32-bytes-long!!!"))

	r := chi.NewRouter()
	r.Get("/api/v1/auth/callback", h.Callback)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=abc&state=xyz", nil)
	// No state cookie
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestAuthHandler_Callback_MissingCode(t *testing.T) {
	h := setupAuthHandler(t, []byte("test-secret-key-32-bytes-long!!!"))

	r := chi.NewRouter()
	r.Get("/api/v1/auth/callback", h.Callback)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Callback_UserNotAllowed(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")

	// Set up mock GitHub servers
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.GitHubTokenResponse{AccessToken: "mock-token"})
	}))
	defer tokenServer.Close()

	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.GitHubUser{ID: 999, Login: "stranger", AvatarURL: "https://example.com/stranger.png"})
	}))
	defer userServer.Close()

	h := setupAuthHandler(t, secret)
	h.tokenEndpoint = tokenServer.URL
	h.userEndpoint = userServer.URL

	r := chi.NewRouter()
	r.Get("/api/v1/auth/callback", h.Callback)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=valid&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}

	var resp map[string]string
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["error"] != "user not authorized" {
		t.Errorf("error = %q, want 'user not authorized'", resp["error"])
	}
}

func TestAuthHandler_Callback_Success(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")

	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.GitHubTokenResponse{AccessToken: "mock-token"})
	}))
	defer tokenServer.Close()

	userServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mock-token" {
			t.Errorf("expected Bearer mock-token, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(middleware.GitHubUser{ID: 42, Login: "admin", AvatarURL: "https://example.com/admin.png"})
	}))
	defer userServer.Close()

	h := setupAuthHandler(t, secret)
	h.tokenEndpoint = tokenServer.URL
	h.userEndpoint = userServer.URL

	r := chi.NewRouter()
	r.Get("/api/v1/auth/callback", h.Callback)

	req := httptest.NewRequest("GET", "/api/v1/auth/callback?code=valid-code&state=xyz", nil)
	req.AddCookie(&http.Cookie{Name: "oauth_state", Value: "xyz"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusFound)
	}

	location := w.Header().Get("Location")
	if location != "/admin" {
		t.Errorf("location = %q, want /admin", location)
	}

	// Should set neo_session cookie
	cookies := w.Result().Cookies()
	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "neo_session" {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("neo_session cookie not set")
	}
	if !sessionCookie.HttpOnly {
		t.Error("neo_session cookie should be HttpOnly")
	}

	// Validate the JWT in the cookie
	claims, err := middleware.ValidateJWT(sessionCookie.Value, secret)
	if err != nil {
		t.Fatalf("validate session JWT: %v", err)
	}
	if claims["username"] != "admin" {
		t.Errorf("username = %v, want admin", claims["username"])
	}
	if claims["sub"] != "42" {
		t.Errorf("sub = %v, want 42", claims["sub"])
	}
}

func TestAuthHandler_Me_Authenticated(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	h := setupAuthHandler(t, secret)

	user := middleware.OAuthUser{ID: "42", Username: "admin", AvatarURL: "https://example.com/a.png", Provider: "github"}
	tokenStr, _ := middleware.CreateJWT(secret, user, 24*time.Hour)

	r := chi.NewRouter()
	r.Get("/api/v1/auth/me", h.Me)

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "neo_session", Value: tokenStr})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["username"] != "admin" {
		t.Errorf("username = %v, want admin", resp["username"])
	}
	if resp["provider"] != "github" {
		t.Errorf("provider = %v, want github", resp["provider"])
	}
}

func TestAuthHandler_Me_NoSession(t *testing.T) {
	h := setupAuthHandler(t, []byte("test-secret-key-32-bytes-long!!!"))

	r := chi.NewRouter()
	r.Get("/api/v1/auth/me", h.Me)

	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	secret := []byte("test-secret-key-32-bytes-long!!!")
	h := setupAuthHandler(t, secret)

	r := chi.NewRouter()
	r.Post("/api/v1/auth/logout", h.Logout)

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	user := middleware.OAuthUser{ID: "1", Username: "u", AvatarURL: "", Provider: "github"}
	tokenStr, _ := middleware.CreateJWT(secret, user, 24*time.Hour)
	req.AddCookie(&http.Cookie{Name: "neo_session", Value: tokenStr})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	// Session cookie should be cleared
	cookies := w.Result().Cookies()
	for _, c := range cookies {
		if c.Name == "neo_session" && c.MaxAge != -1 {
			t.Errorf("neo_session MaxAge = %d, want -1", c.MaxAge)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	return fmt.Sprintf("%s", s) != "" && len(s) >= len(sub) && searchStr(s, sub)
}

func searchStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
