# neo

Personal web garden — a blog, notebook, and widget platform with a retro aesthetic.

## Architecture

- **Frontend**: React + Vite + Tailwind CSS + TypeScript
- **Backend**: Go with Chi router (handler → service → repository)
- **Widgets**: Python FastAPI service for ML widget embedding
- **Database**: SQLite (dev), PostgreSQL (prod)

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
| `DATABASE_URL` | `sqlite://neo.db` | Database connection (`sqlite://` or `postgres://`) |
| `ADMIN_USERNAME` | `admin` | Admin panel username |
| `ADMIN_PASSWORD` | `changeme` | Admin panel password |
| `CORS_ORIGINS` | `http://localhost:5173` | Allowed CORS origins |
| `WIDGET_SERVICE_URL` | `http://localhost:8000` | Widget service URL |
| `SERVE_STATIC` | (empty) | Set to `true` to serve frontend from backend |
| `STATIC_DIR` | `./static` | Path to built frontend files |

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
