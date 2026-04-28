# Architecture

## High-Level Overview

```
┌─────────────────────────────────────────────────┐
│                   Client (Browser)              │
│                                                 │
│  React + Vite + Tailwind CSS + React Router     │
│  (SPA with client-side routing)                 │
└──────────────┬──────────────────┬───────────────┘
               │ REST API         │ Widget embed
               ▼                  ▼
┌──────────────────────┐  ┌──────────────────────┐
│   Backend API (Go)   │  │  Widget Service (Py) │
│                      │  │                      │
│  Chi router          │  │  FastAPI             │
│  Business logic      │  │  HuggingFace models  │
│  Auth middleware      │  │  Notebook runner     │
│  Content management  │  │                      │
└──────────┬───────────┘  └──────────────────────┘
           │
           ▼
┌──────────────────────┐
│   Database           │
│   SQLite (dev)       │
│   PostgreSQL (prod)  │
└──────────────────────┘
```

## Tech Stack

| Layer | Technology | Rationale |
|---|---|---|
| **Frontend** | React 18, Vite, Tailwind CSS, React Router | Fast dev builds, utility-first CSS, SPA routing |
| **Backend API** | Go (Chi router) | Fast, strongly typed, single binary deployment |
| **Widget Service** | Python (FastAPI) | Native HuggingFace/ML ecosystem, notebook support |
| **Database** | SQLite (dev) / PostgreSQL (prod) | Zero-config locally, scalable in production |
| **Auth** | Basic auth (phase 1), OAuth2 SSO (phase 2) | Progressive security complexity |

### Why Go over Node for the backend?

The user mentioned Node and Go. We use **Go** for the API server because:
- Single binary deployment (simpler ops)
- Excellent performance for API workloads
- Strong typing catches bugs at compile time
- Node remains in the stack via the Vite frontend toolchain

Node.js is used for the **frontend build toolchain** (Vite, Tailwind, dev server).

## Project Structure

```
neo/
├── CLAUDE.md
├── README.md
├── docs/
│   └── prepare/              # This implementation plan
│
├── frontend/                 # React SPA
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── index.html
│   ├── public/
│   │   └── fonts/
│   ├── src/
│   │   ├── main.tsx
│   │   ├── App.tsx
│   │   ├── routes/           # Page-level route components
│   │   │   ├── Home.tsx
│   │   │   ├── Post.tsx
│   │   │   ├── Admin.tsx
│   │   │   └── NotFound.tsx
│   │   ├── components/       # Reusable UI components
│   │   │   ├── layout/
│   │   │   │   ├── Header.tsx
│   │   │   │   ├── Footer.tsx
│   │   │   │   └── Sidebar.tsx
│   │   │   ├── blog/
│   │   │   │   ├── PostCard.tsx
│   │   │   │   ├── PostList.tsx
│   │   │   │   ├── TagPill.tsx
│   │   │   │   └── CommentSection.tsx
│   │   │   ├── admin/
│   │   │   │   ├── PostEditor.tsx
│   │   │   │   └── PageManager.tsx
│   │   │   └── widgets/
│   │   │       └── WidgetEmbed.tsx
│   │   ├── hooks/            # Custom React hooks
│   │   ├── services/         # API client functions
│   │   │   └── api.ts
│   │   ├── themes/           # Theme definitions
│   │   │   ├── tokens.ts     # Design token types
│   │   │   ├── retro.ts      # Default retro theme
│   │   │   └── ThemeProvider.tsx
│   │   ├── types/            # Shared TypeScript types
│   │   └── utils/            # Utility functions
│   └── tests/
│       ├── components/
│       └── setup.ts
│
├── backend/                  # Go API server
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go       # Entrypoint
│   ├── internal/
│   │   ├── config/           # Configuration loading
│   │   │   └── config.go
│   │   ├── handler/          # HTTP handlers (controllers)
│   │   │   ├── posts.go
│   │   │   ├── comments.go
│   │   │   ├── admin.go
│   │   │   └── health.go
│   │   ├── middleware/       # Auth, CORS, logging
│   │   │   ├── auth.go
│   │   │   └── cors.go
│   │   ├── model/            # Domain models
│   │   │   ├── post.go
│   │   │   └── comment.go
│   │   ├── service/          # Business logic
│   │   │   ├── post_service.go
│   │   │   └── comment_service.go
│   │   ├── repository/       # Data access layer
│   │   │   ├── post_repo.go
│   │   │   └── comment_repo.go
│   │   └── database/         # DB connection and migrations
│   │       ├── db.go
│   │       └── migrations/
│   └── tests/
│       ├── handler/
│       ├── service/
│       └── repository/
│
├── widgets/                  # Python widget service
│   ├── pyproject.toml
│   ├── app/
│   │   ├── main.py           # FastAPI entrypoint
│   │   ├── routes/
│   │   └── services/
│   └── tests/
│
├── docker-compose.yml        # Local dev orchestration
├── Dockerfile.frontend
├── Dockerfile.backend
├── Dockerfile.widgets
└── .env.example
```

## Separation of Concerns

### Business Logic vs Infrastructure

```
handler/     → HTTP concerns only (parse request, return response)
service/     → Business rules (validation, orchestration)
repository/  → Data access (SQL queries, DB operations)
config/      → Environment-aware configuration
middleware/  → Cross-cutting concerns (auth, logging, CORS)
```

Handlers never touch the database directly. Services never know about HTTP. Repositories never contain business rules.

### Frontend Separation

```
routes/      → Page-level components bound to URL paths
components/  → Reusable, stateless UI components
services/    → API client layer (fetch wrappers)
hooks/       → Stateful logic extracted from components
themes/      → Visual tokens and theme provider
```

Components receive data via props. API calls live in `services/`. State management uses React Context + hooks (no Redux needed at this scale).

## API Design

RESTful JSON API. All endpoints prefixed with `/api/v1/`.

### Public Endpoints

```
GET    /api/v1/posts              # List posts (paginated, filterable by tag)
GET    /api/v1/posts/:slug        # Get single post
GET    /api/v1/posts/:slug/comments  # List comments for a post
POST   /api/v1/posts/:slug/comments  # Add a comment (rate-limited)
GET    /api/v1/tags               # List all tags
```

### Admin Endpoints (authenticated)

```
POST   /api/v1/admin/posts        # Create post
PUT    /api/v1/admin/posts/:slug  # Update post
DELETE /api/v1/admin/posts/:slug  # Delete post
GET    /api/v1/admin/pages        # List pages
POST   /api/v1/admin/pages        # Create page
PUT    /api/v1/admin/pages/:slug  # Update page
DELETE /api/v1/admin/pages/:slug  # Delete page
```

### Health

```
GET    /api/v1/health             # Service health check
```

## Database Schema (Initial)

```sql
CREATE TABLE posts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    slug        TEXT UNIQUE NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    content_type TEXT DEFAULT 'markdown',  -- markdown, html, notebook
    tags        TEXT,                       -- JSON array of strings
    published   BOOLEAN DEFAULT FALSE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE comments (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id     INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    author_name TEXT NOT NULL,
    content     TEXT NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE pages (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    slug        TEXT UNIQUE NOT NULL,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    sort_order  INTEGER DEFAULT 0,
    published   BOOLEAN DEFAULT FALSE,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Configuration

Environment-based config with `.env` files. The backend reads config from environment variables with sensible defaults.

```
# .env.example
PORT=8080
DATABASE_URL=sqlite://neo.db
ADMIN_USERNAME=admin
ADMIN_PASSWORD=changeme
CORS_ORIGINS=http://localhost:5173
WIDGET_SERVICE_URL=http://localhost:8000
```

No config in business logic. All external values injected via config layer.
