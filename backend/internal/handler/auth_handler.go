package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alliemony/neo/backend/internal/middleware"
)

const (
	githubAuthorizeURL = "https://github.com/login/oauth/authorize"
	githubTokenURL     = "https://github.com/login/oauth/access_token"
	githubUserURL      = "https://api.github.com/user"
	sessionExpiry      = 24 * time.Hour
)

// AuthHandler handles OAuth login, callback, logout, and session endpoints.
type AuthHandler struct {
	oauth *middleware.OAuthAuthenticator
	// Overridable URLs for testing
	tokenEndpoint string
	userEndpoint  string
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(oauth *middleware.OAuthAuthenticator) *AuthHandler {
	return &AuthHandler{
		oauth:         oauth,
		tokenEndpoint: githubTokenURL,
		userEndpoint:  githubUserURL,
	}
}

// Login redirects the user to GitHub's OAuth authorization page.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state, err := middleware.GenerateState()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate state")
		return
	}

	// Store state in a short-lived cookie for CSRF protection
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600, // 10 minutes
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.oauth.Config.Secure,
	})

	url := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=read:user&state=%s",
		githubAuthorizeURL,
		h.oauth.Config.ClientID,
		h.oauth.Config.RedirectURL,
		state,
	)

	http.Redirect(w, r, url, http.StatusFound)
}

// Callback handles the OAuth callback from GitHub.
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	// 1. Verify state parameter
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value == "" {
		writeError(w, http.StatusForbidden, "missing state cookie")
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState != stateCookie.Value {
		writeError(w, http.StatusForbidden, "invalid state parameter")
		return
	}

	// Clear the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// 2. Exchange code for access token
	code := r.URL.Query().Get("code")
	if code == "" {
		writeError(w, http.StatusBadRequest, "missing code parameter")
		return
	}

	tokenResp, err := middleware.ExchangeCodeForToken(code, h.oauth.Config, h.tokenEndpoint)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to exchange code")
		return
	}

	// 3. Fetch GitHub user profile
	ghUser, err := middleware.FetchGitHubUser(tokenResp.AccessToken, h.userEndpoint)
	if err != nil {
		writeError(w, http.StatusBadGateway, "failed to fetch user")
		return
	}

	// 4. Check allowlist
	if !middleware.IsUserAllowed(ghUser.Login, h.oauth.Config.AllowedUsers) {
		writeError(w, http.StatusForbidden, "user not authorized")
		return
	}

	// 5. Create JWT session
	user := middleware.OAuthUser{
		ID:        fmt.Sprintf("%d", ghUser.ID),
		Username:  ghUser.Login,
		AvatarURL: ghUser.AvatarURL,
		Provider:  "github",
	}

	tokenStr, err := middleware.CreateJWT(h.oauth.Config.SessionSecret, user, sessionExpiry)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	// 6. Set session cookie and redirect to admin
	http.SetCookie(w, &http.Cookie{
		Name:     "neo_session",
		Value:    tokenStr,
		Path:     "/",
		MaxAge:   int(sessionExpiry.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.oauth.Config.Secure,
	})

	http.Redirect(w, r, "/admin", http.StatusFound)
}

// Me returns the current authenticated user's info.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("neo_session")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	claims, err := middleware.ValidateJWT(cookie.Value, h.oauth.Config.SessionSecret)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"username":   claims["username"],
		"avatar_url": claims["avatar_url"],
		"provider":   claims["provider"],
	})
}

// Logout clears the session cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "neo_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.oauth.Config.Secure,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// AuthMode returns the current auth mode (oauth or basic).
func (h *AuthHandler) AuthMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mode": "oauth"})
}
