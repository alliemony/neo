# neo

Personal web garden — a blog, notebook, and widget platform with a retro aesthetic.

## Architecture

- **Frontend**: React + Vite + Tailwind CSS + TypeScript
- **Backend**: Go with Chi router (handler → service → repository)
- **Widgets**: Python FastAPI service for ML widget embedding
- **Database**: SQLite (dev + lightweight production), PostgreSQL (optional larger production)

## Quick Start

```bash
# Run everything with Docker Compose
docker compose up --build

# Or run services individually:
cd frontend && npm install && npm run dev     # :5173
cd backend && go run ./cmd/server             # :8080
cd widgets && pip install -e . && uvicorn app.main:app --reload  # :8000
```

## Development

```bash
# Run all tests
make test

# Individual test suites
cd frontend && npm test
cd backend && go test ./...
cd widgets && python -m pytest tests/
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Backend server port |
| `DATABASE_URL` | `sqlite://neo.db` | Database connection (`sqlite://` for file-based SQLite or `postgres://` for hosted Postgres) |
| `ADMIN_USERNAME` | `admin` | Admin panel username |
| `ADMIN_PASSWORD` | `changeme` | Admin panel password |
| `CORS_ORIGINS` | `http://localhost:5173` | Allowed CORS origins |
| `WIDGET_SERVICE_URL` | `http://localhost:8000` | Widget service URL |
| `AUTH_MODE` | `basic` | Admin auth mode (`basic` or `oauth`) |
| `GITHUB_CLIENT_ID` | (empty) | GitHub OAuth app client ID |
| `GITHUB_CLIENT_SECRET` | (empty) | GitHub OAuth app client secret |
| `OAUTH_ALLOWED_USERS` | (empty) | Comma-separated GitHub usernames allowed into admin |
| `SESSION_SECRET` | `change-me-to-a-random-32-byte-key` | Secret used to sign OAuth session JWTs |
| `BASE_URL` | `http://localhost:8080` | Public backend base URL used for OAuth callback URLs |
| `SERVE_STATIC` | (empty) | Set to `true` to serve frontend from backend |
| `STATIC_DIR` | `./static` | Path to built frontend files |

## Admin Authentication

### Basic auth for local development

Leave `AUTH_MODE=basic` (the default), then log in with `ADMIN_USERNAME` and `ADMIN_PASSWORD`.

### GitHub OAuth for staging or production

1. Create a GitHub OAuth app in GitHub Developer Settings.
2. Set the homepage URL to your site or backend origin, for example `http://localhost:8080` in local testing.
3. Set the authorization callback URL to `${BASE_URL}/api/v1/auth/callback`, for example `http://localhost:8080/api/v1/auth/callback`.
4. Set `AUTH_MODE=oauth` in `.env`.
5. Set `GITHUB_CLIENT_ID` and `GITHUB_CLIENT_SECRET` from the GitHub app.
6. Set `OAUTH_ALLOWED_USERS` to a comma-separated allowlist of GitHub usernames.
7. Set `SESSION_SECRET` to a long random secret.
8. Restart the backend and use the `Sign in with GitHub` button on `/admin/login`.

The backend keeps the OAuth session in an HTTP-only cookie. Switching `AUTH_MODE` back to `basic` preserves the existing username/password flow for local development.

## Frontend API Mocks

MSW handlers are available for tests and optional local UI work.

- Tests start the MSW server automatically from [frontend/src/test-setup.ts](/Users/adamgoh/Documents/GitHub/neo/frontend/src/test-setup.ts).
- To enable browser-side API mocks during `npm run dev`, start the frontend with `VITE_ENABLE_API_MOCKS=true`.

## Production Deployment

### Fly.io

```bash
# Install flyctl and authenticate
fly auth login

# Create the app (first time only)
fly apps create neo

# Set secrets
fly secrets set ADMIN_USERNAME=admin ADMIN_PASSWORD='<bcrypt-hash>' DATABASE_URL='postgres://...'

# Deploy
fly deploy
```

The production Dockerfile (`Dockerfile.production`) builds the frontend and backend into a single image with static file serving enabled.

> For a lower-cost deployment that avoids a separate hosted database, set `DATABASE_URL=sqlite:///data/neo.db` and mount a persistent volume for `/data`.
>
### Docker (self-hosted)

```bash
# Build the production image
docker build -f Dockerfile.production -t neo .

# Run with environment variables
docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@db:5432/neo \
  -e ADMIN_USERNAME=admin \
  -e ADMIN_PASSWORD=changeme \
  neo
```

## Features

- Blog with markdown posts, tags, and comments
- Like button on posts
- Static pages with customizable navigation
- Admin panel for content management
- Interactive ML widgets (sentiment analysis)
- RSS feed at `/api/v1/feed.xml`
- SEO with Open Graph meta tags
- Retro theme system (dark/light modes)
