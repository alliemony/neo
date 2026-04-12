# Design: OAuth Login for Admin Page

## Technical Approach

### OAuth Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Frontend │     │ Backend  │     │  GitHub  │
└────┬─────┘     └────┬─────┘     └────┬─────┘
     │                │                │
     │ Click Login    │                │
     ├───────────────►│                │
     │                │  302 Redirect  │
     │◄───────────────┤                │
     │                │                │
     │  Redirect to GitHub             │
     ├────────────────────────────────►│
     │                │                │
     │         User authorizes app     │
     │◄────────────────────────────────┤
     │  Redirect to /auth/callback     │
     │                │                │
     │  code + state  │                │
     ├───────────────►│                │
     │                │  Exchange code │
     │                ├───────────────►│
     │                │  Access token  │
     │                │◄───────────────┤
     │                │  Get user info │
     │                ├───────────────►│
     │                │  User profile  │
     │                │◄───────────────┤
     │                │                │
     │  Set-Cookie: session=JWT        │
     │  302 → /admin  │                │
     │◄───────────────┤                │
     │                │                │
```

### Backend Architecture

```go
// middleware/oauth.go

type OAuthConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    AllowedUsers []string
    SessionSecret []byte
}

type OAuthAuthenticator struct {
    config OAuthConfig
}

// Implements Authenticator interface — same as BasicAuthenticator
func (o *OAuthAuthenticator) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            cookie, err := r.Cookie("session")
            if err != nil {
                http.Error(w, `{"error":"unauthorized"}`, 401)
                return
            }
            claims, err := validateJWT(cookie.Value, o.config.SessionSecret)
            if err != nil {
                http.Error(w, `{"error":"unauthorized"}`, 401)
                return
            }
            ctx := context.WithValue(r.Context(), userKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Auth Routes

```go
// handler/auth_handler.go

type AuthHandler struct {
    oauth *OAuthAuthenticator
}

// GET /api/v1/auth/login — redirect to GitHub
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    state := generateRandomState()
    setStateCookie(w, state)
    url := fmt.Sprintf(
        "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=read:user&state=%s",
        h.oauth.config.ClientID, h.oauth.config.RedirectURL, state,
    )
    http.Redirect(w, r, url, http.StatusFound)
}

// GET /api/v1/auth/callback — exchange code, create session
func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
    // 1. Verify state matches cookie
    // 2. Exchange code for access token
    // 3. Fetch GitHub user profile
    // 4. Check username against allowlist
    // 5. Create JWT, set cookie, redirect to /admin
}

// POST /api/v1/auth/logout — clear session cookie
// GET /api/v1/auth/me — return current user from JWT
```

### JWT Structure

```json
{
  "sub": "12345",
  "username": "adamgoh",
  "avatar_url": "https://avatars.githubusercontent.com/u/12345",
  "provider": "github",
  "exp": 1700000000,
  "iat": 1699913600
}
```

### Authenticator Selection in main.go

```go
var authenticator middleware.Authenticator

switch cfg.AuthMode {
case "oauth":
    authenticator = middleware.NewOAuthAuthenticator(middleware.OAuthConfig{
        ClientID:     cfg.GitHubClientID,
        ClientSecret: cfg.GitHubClientSecret,
        RedirectURL:  cfg.BaseURL + "/api/v1/auth/callback",
        AllowedUsers: strings.Split(cfg.OAuthAllowedUsers, ","),
        SessionSecret: []byte(cfg.SessionSecret),
    })
default:
    authenticator = middleware.NewBasicAuthenticator(cfg.AdminUsername, cfg.AdminPassword)
}
```

### Frontend Changes

```tsx
// Current: AdminLogin.tsx with username/password form
// New: AdminLogin.tsx with OAuth button

function AdminLogin() {
  return (
    <Layout>
      <div className="max-w-sm mx-auto mt-20 text-center">
        <h1 className="text-2xl font-bold mb-8">Admin Login</h1>
        <a
          href="/api/v1/auth/login"
          className="inline-flex items-center gap-2 px-6 py-3 bg-gray-900 text-white rounded-lg hover:bg-gray-800"
        >
          <GitHubIcon />
          Sign in with GitHub
        </a>
      </div>
    </Layout>
  );
}
```

### Auth Context Updates

```tsx
// Current: stores username/password, sends Basic Auth header
// New: session cookie is sent automatically, use /auth/me to check status

function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch('/api/v1/auth/me', { credentials: 'include' })
      .then(res => res.ok ? res.json() : null)
      .then(setUser)
      .finally(() => setLoading(false));
  }, []);

  // No Authorization header needed — cookie is automatic
}
```

### Security Considerations

- **CSRF Protection**: OAuth state parameter verified against cookie
- **Cookie Security**: `HttpOnly`, `Secure` (HTTPS only), `SameSite=Lax`
- **Token Expiry**: 24-hour JWT expiry, no refresh tokens (re-login required)
- **Allowlist**: Only configured GitHub users can access admin, even with valid OAuth
- **No localStorage**: Tokens never exposed to JavaScript
- **SESSION_SECRET**: Must be a strong random value (≥32 bytes), provided via env var

### Dependencies

**Backend:**
- `github.com/golang-jwt/jwt/v5` — JWT creation and validation

**Frontend:**
- No new dependencies (cookie-based auth is native)
