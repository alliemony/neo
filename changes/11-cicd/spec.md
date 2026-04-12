# Spec: CI/CD with GitHub Actions

## Files
- `.github/workflows/ci.yml` — Test all services
- `.github/workflows/build.yml` — Build production Docker image
- `.github/workflows/deploy.yml` — Deploy to Fly.io

## CI Workflow (`ci.yml`)
- Runs on: `push` (all branches), `pull_request` (to `main`)
- Three parallel jobs: `test-frontend`, `test-backend`, `test-widgets`
- Each job checks out code, sets up the runtime, installs deps, runs tests

## Build Workflow (`build.yml`)
- Runs on: `push` to `main`
- Depends on CI passing (uses `workflow_run`)
- Builds `Dockerfile.production` using Docker Buildx
- Tags image as `ghcr.io/<owner>/neo:latest` and `ghcr.io/<owner>/neo:<sha>`

## Deploy Workflow (`deploy.yml`)
- Runs on: `push` to `main`
- Depends on build passing (uses `workflow_run`)
- Checks `FLY_API_TOKEN` secret existence
- Runs `flyctl deploy --remote-only`
- Fails with actionable error message if token not configured

## Test Matrix
| Job | Runtime | Command |
|-----|---------|---------|
| test-frontend | Node 20 | `npm ci && npm test` |
| test-backend | Go 1.25 | `go test ./...` |
| test-widgets | Python 3.11 | `pip install -e ".[dev]" && pytest` |
