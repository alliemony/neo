# Design: Admin Panel

## Technical Approach

### Auth Architecture

```go
// middleware/auth.go
type Authenticator interface {
    Middleware() func(http.Handler) http.Handler
    CurrentUser(r *http.Request) (*User, error)
}

type BasicAuthenticator struct {
    username     string
    passwordHash string // bcrypt
}

func (a *BasicAuthenticator) Middleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, pass, ok := r.BasicAuth()
            if !ok || user != a.username || !checkHash(pass, a.passwordHash) {
                w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
                http.Error(w, `{"error":"unauthorized"}`, 401)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Admin Routes

```go
r.Route("/api/v1/admin", func(r chi.Router) {
    r.Use(authenticator.Middleware())

    r.Post("/posts", adminHandler.CreatePost)
    r.Put("/posts/{slug}", adminHandler.UpdatePost)
    r.Delete("/posts/{slug}", adminHandler.DeletePost)

    r.Get("/pages", adminHandler.ListPages)
    r.Post("/pages", adminHandler.CreatePage)
    r.Put("/pages/{slug}", adminHandler.UpdatePage)
    r.Delete("/pages/{slug}", adminHandler.DeletePage)
})
```

### Frontend Auth Flow

```
1. User visits /admin
2. If no stored credentials → show login form
3. User enters username/password
4. Credentials stored in memory (not localStorage for security)
5. All admin API calls include Authorization: Basic header
6. On 401 response → redirect to login
```

### Post Editor Layout

```
┌─────────────────┬─────────────────┐
│  Input Fields    │  Live Preview   │
│                 │                  │
│  Title: [____]  │  # My Post      │
│  Slug: [____]   │                  │
│  Tags: [____]   │  Content here... │
│  Type: [____]   │                  │
│                 │                  │
│  [Markdown      │                  │
│   textarea      │                  │
│   ............] │                  │
│                 │                  │
│  ☐ Published    │                  │
│                 │                  │
│  [Save Draft]   │                  │
│  [Publish]      │                  │
└─────────────────┴─────────────────┘
```

The editor uses a controlled textarea with `onChange` debouncing the preview render (150ms).

### Tag Autocomplete

Simple dropdown filtered by input text:

```tsx
function TagInput({ value, onChange, existingTags }: Props) {
  const [input, setInput] = useState('');
  const suggestions = existingTags.filter(t =>
    t.toLowerCase().startsWith(input.toLowerCase()) && !value.includes(t)
  );
  // render input + dropdown
}
```

Tags are stored as comma-separated in the input, converted to an array on save.
