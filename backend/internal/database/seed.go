package database

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Seed inserts sample posts, pages, and comments if the database is empty.
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
		contentType string
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
		{
			slug:        "try-sentiment-analysis",
			title:       "Try: Sentiment Analysis Widget",
			content:     "sentiment",
			contentType: "widget",
			tags:        []string{"widgets", "ml", "interactive"},
			published:   true,
			createdAt:   time.Now().Add(-6 * time.Hour),
		},
	}

	for _, s := range seeds {
		tagsJSON, _ := json.Marshal(s.tags)
		ct := s.contentType
		if ct == "" {
			ct = "markdown"
		}
		_, err := db.Exec(
			`INSERT INTO posts (slug, title, content, content_type, tags, published, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			s.slug, s.title, s.content, ct, string(tagsJSON), s.published, s.createdAt, s.createdAt,
		)
		if err != nil {
			return err
		}
	}

	// Seed pages
	pages := []struct {
		slug      string
		title     string
		content   string
		published bool
		sortOrder int
	}{
		{
			slug:  "about",
			title: "About",
			content: `# About Neo

Neo is a personal web garden — a space for writing, experimentation, and creative projects.

## Who am I?

I'm a developer who loves building things on the web. This site is my little corner of the internet where I share thoughts about programming, design, and technology.

## The stack

- **Frontend**: React + TypeScript + Tailwind CSS
- **Backend**: Go with Chi router
- **Widgets**: Python FastAPI with HuggingFace models
- **Database**: SQLite (dev) / PostgreSQL (prod)

## Philosophy

I believe in simple, functional design. No unnecessary complexity. Every line of code should earn its place.`,
			published: true,
			sortOrder: 1,
		},
		{
			slug:  "projects",
			title: "Projects",
			content: `# Projects

A collection of things I've built or am currently working on.

## Neo (this site)
A personal web garden with blog, notebooks, and ML widgets. Built with Go, React, and Python.

## CLI Tools
Various command-line utilities for productivity and automation.

## Open Source Contributions
Patches and features contributed to projects I use and care about.

*More projects coming soon.*`,
			published: true,
			sortOrder: 2,
		},
		{
			slug:  "contact",
			title: "Contact",
			content: `# Contact

Want to get in touch? Here's how:

- **Email**: hello@example.com
- **GitHub**: github.com/example
- **Twitter/X**: @example

I'm always happy to chat about programming, design, or interesting projects. Feel free to reach out!`,
			published: true,
			sortOrder: 3,
		},
		{
			slug:      "secret-draft",
			title:     "Secret Draft Page",
			content:   "This page is not published yet. Only admins can see it.",
			published: false,
			sortOrder: 99,
		},
	}

	for _, p := range pages {
		now := time.Now()
		_, err := db.Exec(
			`INSERT INTO pages (slug, title, content, published, sort_order, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.slug, p.title, p.content, p.published, p.sortOrder, now, now,
		)
		if err != nil {
			return err
		}
	}

	// Seed comments on posts
	// Get post IDs first
	var welcomeID, goReactID, retroID int64
	db.QueryRow("SELECT id FROM posts WHERE slug = 'welcome-to-neo'").Scan(&welcomeID)
	db.QueryRow("SELECT id FROM posts WHERE slug = 'building-with-go-and-react'").Scan(&goReactID)
	db.QueryRow("SELECT id FROM posts WHERE slug = 'retro-web-design'").Scan(&retroID)

	comments := []struct {
		postID     int64
		authorName string
		content    string
		createdAt  time.Time
	}{
		{welcomeID, "Alice", "Welcome to the internet! Love the retro vibes.", time.Now().Add(-60 * time.Hour)},
		{welcomeID, "Bob", "Looking forward to more posts. The ML widgets sound really cool!", time.Now().Add(-55 * time.Hour)},
		{welcomeID, "Charlie", "Clean design. Reminds me of the early web days. 👍", time.Now().Add(-50 * time.Hour)},
		{goReactID, "Dave", "Great stack choice! Go + React is a solid combo.", time.Now().Add(-40 * time.Hour)},
		{goReactID, "Eve", "How does the Chi router compare to Gin? Been thinking about switching.", time.Now().Add(-35 * time.Hour)},
		{retroID, "Frank", "Total agree on monospace headings. JetBrains Mono is 🔥", time.Now().Add(-20 * time.Hour)},
		{retroID, "Grace", "The border-radius: 0 approach is bold. I love it.", time.Now().Add(-15 * time.Hour)},
	}

	for _, c := range comments {
		if c.postID == 0 {
			continue
		}
		_, err := db.Exec(
			`INSERT INTO comments (post_id, author_name, content, created_at)
			 VALUES (?, ?, ?, ?)`,
			c.postID, c.authorName, c.content, c.createdAt,
		)
		if err != nil {
			return err
		}
	}

	// Set some like counts on posts
	db.Exec("UPDATE posts SET like_count = 12 WHERE slug = 'welcome-to-neo'")
	db.Exec("UPDATE posts SET like_count = 8 WHERE slug = 'building-with-go-and-react'")
	db.Exec("UPDATE posts SET like_count = 15 WHERE slug = 'retro-web-design'")
	db.Exec("UPDATE posts SET like_count = 5 WHERE slug = 'python-ml-widgets'")

	return nil
}
