# Spec: Project Foundation

## Purpose

Provide a buildable, testable, and runnable project skeleton for all three services (frontend, backend, widgets) with local development infrastructure.

## Requirements

### Requirement: Frontend project SHALL initialize with Vite, React 18, TypeScript, and Tailwind CSS

The frontend project must be a working Vite+React application with TypeScript and Tailwind CSS configured. It must start a dev server and produce a production build.

#### Scenario: Frontend dev server starts successfully

- **GIVEN:** The frontend project is initialized with all dependencies installed
- **WHEN:** The developer runs `npm run dev` from `frontend/`
- **THEN:** A development server starts on port 5173 and serves a page without errors

#### Scenario: Frontend production build succeeds

- **GIVEN:** The frontend project has all dependencies installed
- **WHEN:** The developer runs `npm run build` from `frontend/`
- **THEN:** Static assets are generated in `frontend/dist/` without errors

### Requirement: Backend project SHALL initialize as a Go module with Chi router and health endpoint

The backend must be a Go module that compiles and runs an HTTP server with at least a `/api/v1/health` endpoint.

#### Scenario: Backend server starts and responds to health check

- **GIVEN:** The Go module is initialized with dependencies downloaded
- **WHEN:** The developer runs `go run ./cmd/server` from `backend/`
- **THEN:** An HTTP server starts on port 8080
- **AND:** `GET /api/v1/health` returns `200 OK` with `{"status": "ok"}`

### Requirement: Widget service SHALL initialize as a Python FastAPI project with health endpoint

The widget service must be a FastAPI application with a health check endpoint.

#### Scenario: Widget service starts and responds to health check

- **GIVEN:** The Python dependencies are installed from `pyproject.toml`
- **WHEN:** The developer runs `uvicorn app.main:app` from `widgets/`
- **THEN:** An HTTP server starts on port 8000
- **AND:** `GET /health` returns `200 OK` with `{"status": "ok"}`

### Requirement: Docker Compose SHALL orchestrate all three services for local development

A `docker-compose.yml` at the repo root must define services for frontend, backend, and widgets.

#### Scenario: All services start via Docker Compose

- **GIVEN:** Docker and Docker Compose are installed
- **WHEN:** The developer runs `docker compose up --build`
- **THEN:** All three services start and their health endpoints respond

### Requirement: Each service SHALL have a test framework configured with at least one passing test

Test infrastructure must be in place for TDD from the start.

#### Scenario: Frontend tests pass

- **GIVEN:** Vitest and React Testing Library are configured
- **WHEN:** The developer runs `npm test` from `frontend/`
- **THEN:** At least one test runs and passes

#### Scenario: Backend tests pass

- **GIVEN:** Go test files exist in the backend
- **WHEN:** The developer runs `go test ./...` from `backend/`
- **THEN:** At least one test runs and passes

#### Scenario: Widget tests pass

- **GIVEN:** pytest is installed and test files exist
- **WHEN:** The developer runs `pytest` from `widgets/`
- **THEN:** At least one test runs and passes

### Requirement: Environment configuration SHALL use .env files with a committed example template

#### Scenario: Environment example file documents all required variables

- **GIVEN:** The `.env.example` file exists at the repo root
- **WHEN:** The developer copies it to `.env`
- **THEN:** All three services can start with the default values

### Requirement: A Makefile SHOULD provide common dev commands

#### Scenario: Make targets run expected commands

- **GIVEN:** The Makefile exists at the repo root
- **WHEN:** The developer runs `make test`
- **THEN:** Tests for all three services execute sequentially
