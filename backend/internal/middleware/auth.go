package middleware

import (
	"crypto/subtle"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Authenticator is implemented by types that can produce an HTTP middleware.
type Authenticator interface {
	Middleware() func(http.Handler) http.Handler
}

// basicAuth implements Authenticator using HTTP Basic Auth.
type basicAuth struct {
	username string
	password string
}

// NewBasicAuth returns an Authenticator that validates Basic auth credentials.
// password should be a bcrypt hash; plain-text is accepted as a fallback (dev only).
func NewBasicAuth(username, password string) Authenticator {
	return &basicAuth{username: username, password: password}
}

func (a *basicAuth) Middleware() func(http.Handler) http.Handler {
	return BasicAuth(a.username, a.password)
}

// BasicAuth returns middleware that validates Basic auth credentials.
// passwordHash should be a bcrypt hash. If empty, raw password comparison is used (dev only).
func BasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			if !ok {
				w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			usernameMatch := subtle.ConstantTimeCompare([]byte(u), []byte(username)) == 1

			// Try bcrypt first; if password isn't a valid bcrypt hash, fall back to constant-time compare.
			var passwordMatch bool
			if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(p)); err == nil {
				passwordMatch = true
			} else {
				passwordMatch = subtle.ConstantTimeCompare([]byte(p), []byte(password)) == 1
			}

			if !usernameMatch || !passwordMatch {
				w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
