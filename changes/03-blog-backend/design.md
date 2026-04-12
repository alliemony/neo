# Design: Blog Backend (Posts & Tags API)

## Technical Approach

### Layer Architecture

```
HTTP Request
    ↓
handler/posts.go     (parse params, call service, return JSON)
    ↓
service/post_service.go  (validate, generate slug, orchestrate)
    ↓
repository/post_repo.go  (SQL queries, database operations)
    ↓
database/db.go       (connection, migrations)
```

### Database

Use `database/sql` with `modernc.org/sqlite` (pure Go SQLite, no CGO needed).

```go
// database/db.go
func New(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("sqlite", databaseURL)
    if err != nil {
        return nil, err
    }
    if err := migrate(db); err != nil {
        return nil, err
    }
    return db, nil
}
```

Migrations are embedded SQL files executed in order on startup.

### Model

```go
// model/post.go
type Post struct {
    ID          int64     `json:"id"`
    Slug        string    `json:"slug"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    ContentType string    `json:"content_type"`
    Tags        []string  `json:"tags"`
    Published   bool      `json:"published"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

Tags stored as JSON text in SQLite, parsed to `[]string` in Go.

### Repository Interface

```go
type PostRepository interface {
    Create(post *Post) error
    GetBySlug(slug string) (*Post, error)
    Update(post *Post) error
    Delete(slug string) error
    List(opts ListOptions) ([]Post, int, error)
    ListByTag(tag string, opts ListOptions) ([]Post, int, error)
    AllTags() ([]TagCount, error)
}
```

This interface enables testing services with in-memory fakes.

### Service

```go
type PostService struct {
    repo PostRepository
}

func (s *PostService) Create(input CreatePostInput) (*Post, error) {
    if input.Title == "" {
        return nil, ErrTitleRequired
    }
    slug := slugify(input.Title)
    existing, _ := s.repo.GetBySlug(slug)
    if existing != nil {
        return nil, ErrSlugExists
    }
    post := &Post{
        Slug:    slug,
        Title:   input.Title,
        Content: input.Content,
        Tags:    input.Tags,
    }
    return post, s.repo.Create(post)
}
```

### Handlers

```go
// handler/posts.go
func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
    tag := r.URL.Query().Get("tag")
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    if page < 1 { page = 1 }

    var posts []model.Post
    var total int
    if tag != "" {
        posts, total, _ = h.service.ListByTag(tag, opts)
    } else {
        posts, total, _ = h.service.ListPublished(opts)
    }
    writeJSON(w, 200, map[string]any{"posts": posts, "total": total})
}
```

### Slug Generation

```go
func slugify(title string) string {
    s := strings.ToLower(title)
    s = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(s, "-")
    s = strings.Trim(s, "-")
    return s
}
```

### Routing

```go
r := chi.NewRouter()
r.Use(middleware.CORS(cfg))

r.Route("/api/v1", func(r chi.Router) {
    r.Get("/health", healthHandler.Health)
    r.Get("/posts", postHandler.List)
    r.Get("/posts/{slug}", postHandler.GetBySlug)
    r.Get("/tags", postHandler.ListTags)
})
```
