package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is a private type for context keys in this package.
type contextKey string

const userContextKey contextKey = "oauth_user"

// OAuthUser represents an authenticated user from the JWT claims.
type OAuthUser struct {
	ID        string `json:"sub"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Provider  string `json:"provider"`
}

// OAuthConfig holds configuration for the OAuth authenticator.
type OAuthConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AllowedUsers  []string
	SessionSecret []byte
	Secure        bool // set true in production (HTTPS)
}

// OAuthAuthenticator implements session-based JWT auth after GitHub OAuth.
type OAuthAuthenticator struct {
	Config OAuthConfig
}

// NewOAuthAuthenticator creates a new OAuthAuthenticator.
func NewOAuthAuthenticator(cfg OAuthConfig) *OAuthAuthenticator {
	return &OAuthAuthenticator{Config: cfg}
}

// Middleware returns a chi-compatible middleware that validates JWT session cookies.
func (o *OAuthAuthenticator) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("neo_session")
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			claims, err := ValidateJWT(cookie.Value, o.Config.SessionSecret)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			user := OAuthUser{
				ID:        claims["sub"].(string),
				Username:  claims["username"].(string),
				AvatarURL: claims["avatar_url"].(string),
				Provider:  claims["provider"].(string),
			}

			ctx := context.WithValue(r.Context(), userContextKey, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CurrentUser returns the authenticated user from the request context.
func CurrentUser(r *http.Request) *OAuthUser {
	user, _ := r.Context().Value(userContextKey).(*OAuthUser)
	return user
}

// --- JWT Helpers ---

// CreateJWT creates a signed JWT with the given claims.
func CreateJWT(secret []byte, user OAuthUser, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":        user.ID,
		"username":   user.Username,
		"avatar_url": user.AvatarURL,
		"provider":   user.Provider,
		"iat":        now.Unix(),
		"exp":        now.Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ValidateJWT parses and validates a JWT string, returning its claims.
func ValidateJWT(tokenStr string, secret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// --- GitHub OAuth Helpers ---

// GitHubTokenResponse is the response from GitHub's token endpoint.
type GitHubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubUser is the response from GitHub's user API.
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

// ExchangeCodeForToken exchanges an authorization code for an access token.
// tokenEndpoint is injectable for testing.
func ExchangeCodeForToken(code string, cfg OAuthConfig, tokenEndpoint string) (*GitHubTokenResponse, error) {
	data := url.Values{
		"client_id":     {cfg.ClientID},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"redirect_uri":  {cfg.RedirectURL},
	}

	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read token response: %w", err)
	}

	var tokenResp GitHubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("no access_token in response")
	}

	return &tokenResp, nil
}

// FetchGitHubUser fetches the authenticated user's profile.
// userEndpoint is injectable for testing.
func FetchGitHubUser(accessToken string, userEndpoint string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", userEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user fetch request: %w", err)
	}
	defer resp.Body.Close()

	var user GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("decode user response: %w", err)
	}

	return &user, nil
}

// GenerateState creates a random state string for CSRF protection.
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// IsUserAllowed checks if a username is in the allowed users list.
func IsUserAllowed(username string, allowed []string) bool {
	for _, u := range allowed {
		if strings.EqualFold(strings.TrimSpace(u), username) {
			return true
		}
	}
	return false
}
