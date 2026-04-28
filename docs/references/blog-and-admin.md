# Blog System & Admin Panel

## Blog System

### Post Types

| Type | Description | Rendering |
|---|---|---|
| **Text** | Standard blog post (Markdown) | Rendered via markdown-it or remark |
| **Notebook** | Jupyter-style notebook content | Rendered via custom notebook viewer |
| **Link** | Shared link with commentary | Link card with preview + author notes |
| **Widget** | Embedded Python/HuggingFace widget | Iframe embed with fallback |

### Tumblr-Inspired UX Patterns

The blog takes cues from Tumblr's social post experience while maintaining the retro aesthetic:

#### Post Feed

- **Chronological feed** of post cards on the main page
- **Infinite scroll** or paginated navigation (user preference via config)
- Each post card shows: author mark, timestamp, title, excerpt, tags, interaction counts
- Posts are visually distinct by type (text posts have a body preview, link posts show URL card, widget posts show an interactive embed)

#### Tags

- Tags displayed as **pill-shaped badges** at the bottom of each post card
- Clicking a tag filters the feed to that tag
- Tag cloud in sidebar shows all tags with post counts
- Tags are stored as a JSON array in the database, indexed for query performance

```
[machine-learning] [python] [tutorial]
```

#### Interactions

- **Like/heart**: Anonymous likes stored by post ID (count only, no user tracking)
- **Comment count**: Shows number of comments, links to comment section
- **Share**: Copy link to clipboard

#### Comment Section

Comments appear in the **right sidebar** when viewing a single post (on desktop). On mobile, they appear below the post content.

- **Linear thread** (not nested) -- keeps it simple
- Each comment shows: author name, timestamp, content
- Comments require a display name (no account needed)
- Rate-limited to prevent spam (IP-based, configurable window)
- Optional: basic profanity filter or manual approval queue

```
┌─ Comments (3) ──────────────┐
│                              │
│  alice · 2h ago              │
│  Great post! I've been       │
│  looking for this.           │
│                              │
│  ────────────────────────    │
│                              │
│  bob · 1h ago                │
│  Have you tried using...     │
│                              │
│  ────────────────────────    │
│                              │
│  [Name: _______________]     │
│  [Comment: ____________]     │
│  [____________ Submit ▸]     │
│                              │
└──────────────────────────────┘
```

#### Timestamps

- Relative timestamps for recent posts ("2 hours ago", "yesterday")
- Absolute timestamps for older posts ("Mar 15, 2026")
- Full ISO timestamp available on hover (via `<time>` element)

### Single Post View

```
/blog/:slug

┌─────────────────────────────────┬──────────────────────┐
│                                 │                      │
│  POST TITLE                     │  Comments (3)        │
│  author · Mar 15, 2026          │                      │
│                                 │  alice · 2h ago      │
│  Full post content rendered     │  Great post!...      │
│  from markdown with syntax      │                      │
│  highlighting for code blocks   │  ──────────────      │
│  and proper typography.         │                      │
│                                 │  bob · 1h ago        │
│  [tag-one] [tag-two]            │  Have you tried...   │
│                                 │                      │
│  ♡ 12    ↩ 3     ⊕ share       │  [Add comment...]    │
│                                 │                      │
├─────────────────────────────────┤                      │
│  ← Previous  ·  Next →         │                      │
└─────────────────────────────────┴──────────────────────┘
```

### URL Structure

```
/                       → Home / post feed
/blog/:slug             → Single post view
/tag/:tag               → Posts filtered by tag
/page/:slug             → Static page (about, projects, etc.)
/widgets/:id            → Widget embed page
/admin                  → Admin dashboard (protected)
```

## Admin Panel

### Overview

A minimal, functional admin interface for content management. Secured with basic auth initially, designed to accept OAuth2 SSO as a future upgrade.

### Authentication

#### Phase 1: Basic Auth

- HTTP Basic Authentication over HTTPS
- Credentials stored as hashed values in environment variables
- Session cookie after initial auth (configurable expiry)
- All `/admin` and `/api/v1/admin/*` routes protected by auth middleware

```go
// middleware/auth.go
func BasicAuth(cfg config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user, pass, ok := r.BasicAuth()
            if !ok || !validateCredentials(user, pass, cfg) {
                w.Header().Set("WWW-Authenticate", `Basic realm="admin"`)
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

#### Phase 2: OAuth2 SSO (Future)

The auth middleware is abstracted behind an interface so the implementation can be swapped:

```go
type Authenticator interface {
    Middleware() func(http.Handler) http.Handler
    CurrentUser(r *http.Request) (*User, error)
}

// BasicAuthenticator implements Authenticator (phase 1)
// OAuthAuthenticator implements Authenticator (phase 2)
```

### Admin Dashboard

```
/admin

┌──────────────────────────────────────────────────┐
│  neo admin                          [← Back to site] │
├──────────────────────────────────────────────────┤
│                                                  │
│  ┌─ Posts ───────────────────────────────────┐   │
│  │                                           │   │
│  │  [+ New Post]                             │   │
│  │                                           │   │
│  │  ● Published  My First Post    [Edit][Del]│   │
│  │  ○ Draft      Work in Progress [Edit][Del]│   │
│  │  ● Published  Hello World      [Edit][Del]│   │
│  │                                           │   │
│  └───────────────────────────────────────────┘   │
│                                                  │
│  ┌─ Pages ───────────────────────────────────┐   │
│  │                                           │   │
│  │  [+ New Page]                             │   │
│  │                                           │   │
│  │  About                         [Edit][Del]│   │
│  │  Projects                      [Edit][Del]│   │
│  │                                           │   │
│  └───────────────────────────────────────────┘   │
│                                                  │
└──────────────────────────────────────────────────┘
```

### Post Editor

- **Title field**: Plain text input
- **Slug**: Auto-generated from title, editable
- **Content**: Markdown textarea with live preview side-by-side
- **Tags**: Comma-separated input with autocomplete from existing tags
- **Type**: Dropdown selector (text, notebook, link, widget)
- **Publish toggle**: Draft / Published
- **Save**: Saves as draft or publishes immediately

```
/admin/posts/new
/admin/posts/:slug/edit

┌──────────────────────────────────────────────────┐
│  Edit Post                      [Save Draft] [Publish] │
├──────────────┬───────────────────────────────────┤
│              │                                   │
│  Title:      │  Preview:                         │
│  [________]  │                                   │
│              │  # My Post Title                  │
│  Slug:       │                                   │
│  [my-post]   │  This is the rendered preview     │
│              │  of the markdown content...       │
│  Tags:       │                                   │
│  [ml, py]    │  ```python                        │
│              │  print("hello")                   │
│  Type:       │  ```                              │
│  [Text    ▾] │                                   │
│              │                                   │
│  Content:    │                                   │
│  [# My Post] │                                   │
│  [Title    ] │                                   │
│  [         ] │                                   │
│  [This is  ] │                                   │
│  [the...   ] │                                   │
│              │                                   │
└──────────────┴───────────────────────────────────┘
```

### Admin Features Summary

| Feature | Phase 1 | Future |
|---|---|---|
| Create/edit/delete posts | Yes | - |
| Create/edit/delete pages | Yes | - |
| Markdown editor with preview | Yes | - |
| Tag management | Yes | - |
| Publish/draft toggle | Yes | - |
| Basic auth | Yes | OAuth2 SSO |
| Comment moderation | - | Yes |
| Media uploads | - | Yes |
| Analytics dashboard | - | Yes |
