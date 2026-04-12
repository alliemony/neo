# Design: Project Foundation

## Technical Approach

### Frontend (`frontend/`)

Initialize with `npm create vite@latest` using the React+TypeScript template, then add Tailwind CSS v4.

```
frontend/
├── package.json          # vite, react, react-dom, tailwind, vitest, @testing-library/react
├── vite.config.ts        # Vitest inline config
├── tailwind.config.ts    # Extended with theme CSS variables (empty for now)
├── tsconfig.json
├── index.html
├── src/
│   ├── main.tsx          # ReactDOM.createRoot entry
│   ├── App.tsx           # Placeholder with "neo" text
│   └── App.test.tsx      # First test: renders without crashing
└── tests/
    └── setup.ts          # jsdom setup, @testing-library/jest-dom matchers
```

Key dependencies:
- `react`, `react-dom` (18+)
- `react-router-dom` (7+)
- `tailwindcss`, `@tailwindcss/vite`
- `vitest`, `@testing-library/react`, `@testing-library/jest-dom`, `jsdom`

### Backend (`backend/`)

Standard Go project layout using `internal/` for private packages.

```
backend/
├── go.mod                # module github.com/alliemony/neo/backend
├── cmd/
│   └── server/
│       └── main.go       # Creates router, registers health handler, starts server
├── internal/
│   ├── config/
│   │   └── config.go     # Reads PORT, DATABASE_URL, etc. from env with defaults
│   └── handler/
│       ├── health.go     # GET /api/v1/health → {"status": "ok"}
│       └── health_test.go
```

Key dependencies:
- `github.com/go-chi/chi/v5` (router)
- `github.com/go-chi/cors` (CORS middleware)

### Widget Service (`widgets/`)

Minimal FastAPI setup.

```
widgets/
├── pyproject.toml        # fastapi, uvicorn, pytest, httpx
├── app/
│   ├── __init__.py
│   └── main.py           # FastAPI app with /health endpoint
└── tests/
    ├── __init__.py
    └── test_health.py    # Tests /health returns 200
```

### Infrastructure

```
docker-compose.yml        # 3 services: frontend, backend, widgets
Dockerfile.frontend       # node:20-alpine, npm ci, npm run dev
Dockerfile.backend        # golang:1.22-alpine, go run ./cmd/server
Dockerfile.widgets        # python:3.11-slim, pip install, uvicorn
.env.example              # PORT, DATABASE_URL, ADMIN_USERNAME, etc.
Makefile                  # dev, test, build, clean targets
```

### Dev vs Prod Docker Strategy

For local dev, Dockerfiles run the dev servers with hot reload. Production Dockerfiles (added later in 09-production-deployment) use multi-stage builds.

## Dependencies Between Services

None at this stage. Each service starts independently. The frontend has a `VITE_API_URL` env var but no actual API calls yet.
