# CLAUDE.md - Neo Project

## Project Overview

Neo is a personal web garden: a blog, notebook, and widget platform with a retro aesthetic. It is a monorepo with three services: a React frontend, a Go backend API, and a Python widget service.

## Architecture

- **Frontend**: `frontend/` -- React + Vite + Tailwind CSS + TypeScript
- **Backend**: `backend/` -- Go with Chi router, clean architecture (handler → service → repository)
- **Widgets**: `widgets/` -- Python FastAPI service for ML/HuggingFace widget embedding
- **Database**: SQLite (dev), PostgreSQL (prod)

## Coding Principles

### KISS (Keep It Simple)
- Simplest solution that works. No premature abstractions.
- Three similar lines > one clever abstraction used three times.
- Don't design for hypothetical future requirements.

### Separation of Concerns
- **Handlers** deal with HTTP only (parse request, return response)
- **Services** contain business logic (validation, orchestration)
- **Repositories** manage data access (SQL, DB operations)
- **Config** is injected, never hardcoded in business logic
- Frontend components receive data via props; API calls live in `services/`

### TDD (Test-Driven Development)
- Every feature starts with a failing test.
- Every bug fix starts with a test that reproduces the bug.
- Red → Green → Refactor cycle.
- Tests are co-located with source files (`Component.test.tsx`, `service_test.go`).

### Architecture Boundaries
- Frontend and backend are fully independent -- they communicate only via REST API.
- The widget service is an optional sidecar -- the main site works without it.
- No shared code between frontend and backend (types are duplicated intentionally).

## Commands

### Development
```bash
# Frontend
cd frontend && npm install && npm run dev     # Dev server at :5173
cd frontend && npm test                        # Run Vitest
cd frontend && npm run build                   # Production build

# Backend
cd backend && go run ./cmd/server              # Dev server at :8080
cd backend && go test ./...                    # Run all Go tests

# Widgets
cd widgets && pip install -e . && uvicorn app.main:app --reload  # Dev at :8000
cd widgets && pytest                           # Run pytest

# All services
docker compose up --build                      # Full stack via Docker
make test                                      # Run all tests
```

### Linting & Formatting
```bash
cd frontend && npm run lint                    # ESLint
cd frontend && npm run format                  # Prettier
cd backend && gofmt -w .                       # Go formatting
cd backend && go vet ./...                     # Go static analysis
```

## File Naming Conventions

- **Frontend**: PascalCase for components (`PostCard.tsx`), camelCase for utils/hooks (`usePosts.ts`)
- **Backend (Go)**: snake_case for files (`post_service.go`), PascalCase for exports
- **Tests**: Same name as source with `.test.tsx` / `_test.go` / `test_*.py` suffix
- **Docs**: kebab-case (`design-system.md`)

## Theme System

Themes are defined as TypeScript token objects in `frontend/src/themes/`. Adding a new theme:
1. Create a new file exporting a `ThemeTokens` object
2. Register it in `ThemeProvider`
3. Done -- CSS variables update automatically, Tailwind classes follow

## Git Conventions

- Branch naming: `feature/description`, `fix/description`, `docs/description`
- Commit messages: imperative mood, concise ("Add post list component", "Fix tag filtering query")
- One logical change per commit

## Environment

- `.env.example` is the template (committed)
- `.env` is local config (gitignored)
- Never commit secrets or credentials

## API Patterns

- All API endpoints prefixed with `/api/v1/`
- Admin endpoints under `/api/v1/admin/` (authenticated)
- JSON request/response bodies
- Standard HTTP status codes (200, 201, 400, 401, 404, 500)
- Error responses: `{"error": "description"}`

## When Adding Features

1. Write failing tests first (TDD)
2. Implement the minimum to pass
3. Refactor while tests stay green
4. Ensure no regressions (`make test`)
5. Keep changes focused -- one feature per PR

## When Fixing Bugs

1. Write a test that reproduces the bug
2. Verify the test fails
3. Fix the bug
4. Verify the test passes
5. Check for similar bugs nearby
