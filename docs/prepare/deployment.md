# Deployment

## Local Development

### Prerequisites

- **Node.js** >= 20 (frontend toolchain)
- **Go** >= 1.22 (backend API)
- **Python** >= 3.11 (widget service)
- **Docker** + **Docker Compose** (optional, for full-stack local run)

### Running Locally (without Docker)

```bash
# Terminal 1: Frontend dev server
cd frontend && npm install && npm run dev
# → http://localhost:5173

# Terminal 2: Backend API
cd backend && go run ./cmd/server
# → http://localhost:8080

# Terminal 3: Widget service (optional)
cd widgets && pip install -e . && uvicorn app.main:app --reload
# → http://localhost:8000
```

### Running Locally (with Docker Compose)

```bash
docker compose up --build
# Frontend → http://localhost:5173
# Backend  → http://localhost:8080
# Widgets  → http://localhost:8000
```

## Docker Configuration

### docker-compose.yml

```yaml
services:
  frontend:
    build:
      context: ./frontend
      dockerfile: ../Dockerfile.frontend
    ports:
      - "5173:5173"
    volumes:
      - ./frontend/src:/app/src   # Hot reload
    environment:
      - VITE_API_URL=http://localhost:8080

  backend:
    build:
      context: ./backend
      dockerfile: ../Dockerfile.backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=sqlite:///data/neo.db
      - ADMIN_USERNAME=${ADMIN_USERNAME}
      - ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - CORS_ORIGINS=http://localhost:5173
    volumes:
      - db-data:/data

  widgets:
    build:
      context: ./widgets
      dockerfile: ../Dockerfile.widgets
    ports:
      - "8000:8000"

volumes:
  db-data:
```

### Dockerfile.frontend

```dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

### Dockerfile.backend

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o server ./cmd/server

FROM alpine:latest
RUN apk add --no-cache sqlite
COPY --from=build /app/server /usr/local/bin/server
EXPOSE 8080
CMD ["server"]
```

### Dockerfile.widgets

```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY pyproject.toml ./
RUN pip install --no-cache-dir .
COPY app/ ./app/
EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

## Production Deployment

### Target Platforms

The project is designed to deploy on common platforms:

| Platform | Frontend | Backend | Widgets |
|---|---|---|---|
| **Fly.io** | Static asset serving via backend | Fly app (Go binary) | Fly app (Python) |
| **Railway** | Railway static site | Railway service | Railway service |
| **Render** | Static site | Web service | Web service |
| **VPS (manual)** | Nginx serving static | systemd service | systemd service |

### Recommended: Fly.io

Fly.io is a good fit because:
- Deploys Docker containers natively
- Supports multiple services in one project
- Global edge deployment
- Simple CLI (`flyctl`)
- Affordable for personal projects

### Production Architecture

```
                    ┌────────────────┐
                    │   CDN / Edge   │
                    │  (static assets)│
                    └───────┬────────┘
                            │
                    ┌───────▼────────┐
                    │  Reverse Proxy │
                    │  (Caddy/Nginx) │
                    └──┬─────────┬───┘
                       │         │
           ┌───────────▼──┐  ┌──▼───────────┐
           │  Go Backend  │  │  Python       │
           │  /api/*      │  │  /widgets/*   │
           └──────┬───────┘  └──────────────┘
                  │
           ┌──────▼───────┐
           │  PostgreSQL   │
           └──────────────┘
```

In production:
- Frontend is pre-built static files served from CDN or the Go backend itself
- Go backend serves both the API and the static frontend files
- PostgreSQL replaces SQLite
- HTTPS via platform TLS or Caddy auto-TLS

### Environment Management

```
.env.example        → Checked into git (template)
.env                → Local development (gitignored)
.env.production     → Production values (gitignored, set in platform)
```

Key production environment variables:
```
PORT=8080
DATABASE_URL=postgres://user:pass@host:5432/neo
ADMIN_USERNAME=<set-in-platform>
ADMIN_PASSWORD=<bcrypt-hash>
CORS_ORIGINS=https://yourdomain.com
WIDGET_SERVICE_URL=https://widgets.yourdomain.com
```

## CI/CD (Future)

GitHub Actions pipeline:

```
push to main → lint → test → build → deploy
```

```yaml
# .github/workflows/deploy.yml (future)
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
      - run: cd frontend && npm ci && npm test

  test-backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: cd backend && go test ./...

  deploy:
    needs: [test-frontend, test-backend]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy
```

## Makefile

A top-level `Makefile` orchestrates common tasks:

```makefile
.PHONY: dev test build clean

dev:              ## Start all services locally
	docker compose up --build

dev-frontend:     ## Start frontend only
	cd frontend && npm run dev

dev-backend:      ## Start backend only
	cd backend && go run ./cmd/server

dev-widgets:      ## Start widget service only
	cd widgets && uvicorn app.main:app --reload

test:             ## Run all tests
	cd frontend && npm test
	cd backend && go test ./...
	cd widgets && pytest

build:            ## Build all for production
	cd frontend && npm run build
	cd backend && CGO_ENABLED=1 go build -o bin/server ./cmd/server

clean:            ## Clean build artifacts
	rm -rf frontend/dist backend/bin
```
