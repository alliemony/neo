# Design: Production Deployment

## Technical Approach

### Static File Serving in Go

```go
// cmd/server/main.go
if cfg.ServeStatic {
    // Serve frontend static files
    staticDir := http.Dir(cfg.StaticDir) // defaults to "./static"
    fileServer := http.FileServer(staticDir)

    // SPA fallback: serve index.html for non-file, non-API routes
    r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
        // Try to serve the file; if not found, serve index.html
        path := filepath.Join(cfg.StaticDir, r.URL.Path)
        if _, err := os.Stat(path); os.IsNotExist(err) {
            http.ServeFile(w, r, filepath.Join(cfg.StaticDir, "index.html"))
            return
        }
        fileServer.ServeHTTP(w, r)
    })
}
```

### Production Dockerfile (Unified)

```dockerfile
# Stage 1: Build frontend
FROM node:20-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend-build
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ .
RUN CGO_ENABLED=1 go build -o server ./cmd/server

# Stage 3: Production image
FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite
COPY --from=backend-build /app/backend/server /usr/local/bin/server
COPY --from=frontend-build /app/frontend/dist /app/static
ENV SERVE_STATIC=true
ENV STATIC_DIR=/app/static
EXPOSE 8080
CMD ["server"]
```

### Database Driver Selection

```go
// database/db.go
func New(databaseURL string) (*sql.DB, error) {
    if strings.HasPrefix(databaseURL, "postgres") {
        return sql.Open("pgx", databaseURL)
    }
    return sql.Open("sqlite", strings.TrimPrefix(databaseURL, "sqlite://"))
}
```

Migrations use dialect-aware SQL (e.g., `AUTOINCREMENT` vs `SERIAL`).

### RSS Feed

```go
// handler/feed.go
func (h *FeedHandler) RSS(w http.ResponseWriter, r *http.Request) {
    posts, _, _ := h.service.ListPublished(ListOptions{PerPage: 20})
    feed := generateRSS(posts, siteURL)
    w.Header().Set("Content-Type", "application/rss+xml")
    w.Write([]byte(feed))
}
```

### SEO Meta Tags

Frontend uses `react-helmet-async` for dynamic meta tags:

```tsx
function PostHead({ post }: { post: Post }) {
  return (
    <Helmet>
      <title>{post.title} | neo</title>
      <meta property="og:title" content={post.title} />
      <meta property="og:description" content={excerpt(post.content, 160)} />
      <meta property="og:url" content={`${SITE_URL}/blog/${post.slug}`} />
      <meta property="og:type" content="article" />
    </Helmet>
  );
}
```

### Fly.io Configuration

```toml
# fly.toml
app = "neo-site"
primary_region = "iad"

[build]
  dockerfile = "Dockerfile.production"

[http_service]
  internal_port = 8080
  force_https = true

[env]
  SERVE_STATIC = "true"
  STATIC_DIR = "/app/static"
```

### Deployment Flow

```
1. git push to main
2. CI runs tests (frontend + backend)
3. CI builds production Docker image
4. flyctl deploy (or platform equivalent)
5. Platform pulls image, starts container
6. Health check passes → traffic routed
```
