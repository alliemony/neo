# Proposal: Project Foundation

## Why

The project has no code yet -- only a README and .gitignore. Before any feature work can begin, we need initialized, buildable, and testable frontend, backend, and widget service projects with a working local development environment.

## What

Initialize the three service projects with their respective toolchains, create shared infrastructure (Docker Compose, environment config, Makefile), and establish the test frameworks so that every subsequent change starts with TDD from day one.

## What Changes

- Create `frontend/` with Vite + React + TypeScript + Tailwind CSS
- Create `backend/` with Go module, Chi router, and a health endpoint
- Create `widgets/` with Python FastAPI and a health endpoint
- Create `docker-compose.yml` to orchestrate all three services
- Create `.env.example` with all required environment variables
- Create a top-level `Makefile` for common dev commands
- Configure test frameworks: Vitest (frontend), Go testing (backend), pytest (widgets)
- Each service has at least one passing test (health check)

## Approach

Scaffold each project independently so they can be built and tested in isolation. Docker Compose is optional for local dev -- each service can also be started directly. The foundation is intentionally minimal: no business logic, just the skeleton.
