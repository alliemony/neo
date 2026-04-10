package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Seed inserts sample posts if the database is empty.
func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	seeds := []struct {
		slug        string
		title       string
		content     string
		tags        []string
		published   bool
		createdAt   time.Time
	}{
		{
			slug:  "welcome-to-neo",
			title: "Welcome to Neo",
			content: `# Welcome to Neo

This is the first post on my personal web garden. Neo is a place for thoughts, experiments, and things I'm building.

## What to expect

- **Blog posts** about programming, design, and technology
- **Notebooks** with code experiments and tutorials
- **Widgets** powered by machine learning models

Stay tuned for more.`,
			tags:      []string{"meta", "intro"},
			published: true,
			createdAt: time.Now().Add(-72 * time.Hour),
		},
		{
			slug:  "building-with-go-and-react",
			title: "Building with Go and React",
			content: `# Building with Go and React

This site is built with a Go backend and a React frontend. Here's why I chose this stack.

## Go for the backend

Go is fast, simple, and has an excellent standard library. The ` + "`net/http`" + ` package gives you everything you need for a web server. Combined with Chi for routing, it's a joy to work with.

` + "```go" + `
func main() {
    r := chi.NewRouter()
    r.Get("/api/v1/posts", postHandler.List)
    http.ListenAndServe(":8080", r)
}
` + "```" + `

## React for the frontend

React with TypeScript and Tailwind CSS makes building UIs fast and maintainable. Vite keeps the dev experience snappy.

The retro aesthetic is achieved with CSS custom properties and a theme system that lets you swap between light and dark modes.`,
			tags:      []string{"go", "react", "architecture"},
			published: true,
			createdAt: time.Now().Add(-48 * time.Hour),
		},
		{
			slug:  "retro-web-design",
			title: "Retro Web Design in 2026",
			content: `# Retro Web Design in 2026

There's something appealing about the stark, functional aesthetic of early web design. No rounded corners, no shadows, just content and clear typography.

## Design principles

1. **Stark borders** — 2px solid lines, no border-radius
2. **Monospace headings** — JetBrains Mono for that terminal feel
3. **Warm palette** — off-white backgrounds, muted accents
4. **Readable body text** — Inter at comfortable line heights

The goal isn't nostalgia — it's clarity. Every pixel serves a purpose.`,
			tags:      []string{"design", "css", "retro"},
			published: true,
			createdAt: time.Now().Add(-24 * time.Hour),
		},
		{
			slug:  "python-ml-widgets",
			title: "Embedding ML Widgets with Python",
			content: `# Embedding ML Widgets with Python

One of the coolest features of Neo is the ability to embed interactive ML widgets powered by HuggingFace models.

## How it works

The widget service is a Python FastAPI sidecar that:

1. Loads HuggingFace models on demand
2. Exposes them as simple REST endpoints
3. Gets embedded in blog posts via iframes

This keeps the main site fast while still offering interactive ML demos.`,
			tags:      []string{"python", "ml", "widgets"},
			published: true,
			createdAt: time.Now().Add(-12 * time.Hour),
		},
		{
			slug:    "draft-upcoming-features",
			title:   "Upcoming Features",
			content: "This is a draft post about upcoming features. Not published yet.",
			tags:    []string{"meta"},
			published: false,
			createdAt: time.Now(),
		},
	}

	for _, s := range seeds {
		tagsJSON, _ := json.Marshal(s.tags)
		_, err := db.Exec(
			`INSERT INTO posts (slug, title, content, content_type, tags, published, created_at, updated_at)
			 VALUES (?, ?, ?, 'markdown', ?, ?, ?, ?)`,
			s.slug, s.title, s.content, string(tagsJSON), s.published, s.createdAt, s.createdAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
