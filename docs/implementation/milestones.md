# Implementation Milestones

## Phase 0: Foundation

**Goal**: Project scaffolding, dev environment, and first green test.

- [ ] Initialize frontend project (Vite + React + TypeScript + Tailwind)
- [ ] Initialize backend project (Go module + Chi router)
- [ ] Set up theme token system and `ThemeProvider`
- [ ] Create base layout components (Header, Footer, Sidebar)
- [ ] Set up Vitest + React Testing Library
- [ ] Set up Go test infrastructure
- [ ] Write first passing test for each layer
- [ ] Create `docker-compose.yml` for local dev
- [ ] Create `.env.example`

**Deliverable**: `npm run dev` and `go run ./cmd/server` both start successfully. At least one test passes per layer.

---

## Phase 1: Blog (Read-Only)

**Goal**: Public blog feed and single post view with the retro theme applied.

- [ ] Database schema + migrations (posts table)
- [ ] Post repository (CRUD operations)
- [ ] Post service (list, get by slug, filter by tag)
- [ ] API handlers: `GET /api/v1/posts`, `GET /api/v1/posts/:slug`
- [ ] Seed data: 3-5 sample posts
- [ ] `PostCard` component with tags, timestamp, interactions
- [ ] `PostList` component (feed page)
- [ ] Single post view with full markdown rendering
- [ ] Tag filtering (`/tag/:tag`)
- [ ] Apply retro theme: borders, typography, colors
- [ ] Responsive layout (desktop two-column, mobile single)

**Deliverable**: Visit the site, see a styled feed of posts, click into a post, filter by tag.

---

## Phase 2: Comments & Interactions

**Goal**: Users can interact with posts via likes and comments.

- [ ] Database schema: comments table
- [ ] Comment repository + service
- [ ] API: `GET /POST /api/v1/posts/:slug/comments`
- [ ] Like counter API (anonymous, count-only)
- [ ] `CommentSection` component (sidebar on desktop)
- [ ] Comment form with name + content fields
- [ ] Rate limiting middleware for comment submissions
- [ ] Like button with optimistic UI update

**Deliverable**: Users can like posts and leave comments. Comments appear in the sidebar.

---

## Phase 3: Admin Panel

**Goal**: Authenticated admin can create, edit, and manage posts and pages.

- [ ] Basic auth middleware
- [ ] Auth configuration via environment variables
- [ ] Admin login page
- [ ] Admin dashboard (list posts + pages)
- [ ] Post editor with markdown preview
- [ ] Create / edit / delete posts via admin API
- [ ] Page management (create / edit / delete)
- [ ] Publish / draft toggle
- [ ] Tag input with autocomplete

**Deliverable**: Admin logs in, creates a post with tags, publishes it, sees it on the public site.

---

## Phase 4: Static Pages & Navigation

**Goal**: Support for static pages (about, projects, etc.) and polished navigation.

- [ ] Pages API: `GET /api/v1/pages/:slug`
- [ ] Page rendering component
- [ ] Navigation header with dynamic page links
- [ ] Footer with site info
- [ ] 404 page
- [ ] Tag cloud in sidebar

**Deliverable**: Site has proper navigation, static pages, and a complete layout.

---

## Phase 5: Python Widget Integration

**Goal**: Python-based widgets and HuggingFace models can be embedded in posts.

- [ ] Initialize widget service (FastAPI + pyproject.toml)
- [ ] Health check endpoint
- [ ] Widget registry (list available widgets)
- [ ] `WidgetEmbed` React component (iframe-based)
- [ ] Widget post type in admin editor
- [ ] Sample widget: text generation or simple ML demo
- [ ] Docker setup for widget service

**Deliverable**: A blog post can embed an interactive Python widget that runs server-side.

---

## Phase 6: Polish & Production

**Goal**: Production-ready deployment, performance, and polish.

- [ ] Go backend serves frontend static files in production mode
- [ ] Production Dockerfile (multi-stage builds)
- [ ] Fly.io (or chosen platform) deployment config
- [ ] PostgreSQL migration path
- [ ] HTTPS / TLS configuration
- [ ] Meta tags and Open Graph for social sharing
- [ ] RSS feed endpoint
- [ ] Performance audit (Lighthouse)
- [ ] Accessibility audit

**Deliverable**: Site is deployed and publicly accessible at a custom domain.

---

## Future (Not in initial scope)

- OAuth2 SSO for admin
- Comment moderation queue
- Media uploads (images, files)
- Notebook viewer (Jupyter-style)
- Analytics dashboard
- Search functionality
- Dark mode toggle (theme system already supports it)
- Email notifications for comments
