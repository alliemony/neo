# Design: CI/CD with GitHub Actions

## Workflows

### 1. CI — `ci.yml`
- **Trigger**: Push to any branch, pull requests to `main`
- **Jobs**:
  - `test-frontend`: Node 20, `npm ci`, `npm test`
  - `test-backend`: Go 1.25, `go test ./...`
  - `test-widgets`: Python 3.11, `pip install -e ".[dev]"`, `pytest`
- All three jobs run in parallel

### 2. Build — `build.yml`
- **Trigger**: Push to `main` (only after CI passes)
- **Job**: Build `Dockerfile.production` to verify the production image builds cleanly
- Optionally push to GHCR if `GHCR_PUSH` secret is set

### 3. Deploy — `deploy.yml`
- **Trigger**: Push to `main` (only after build succeeds)
- **Job**: Deploy to Fly.io using `superfly/flyctl-actions`
- **Requires**: `FLY_API_TOKEN` secret — will fail with a clear message if not set

## Secrets Required
| Secret | Purpose |
|--------|---------|
| `FLY_API_TOKEN` | Fly.io deploy token |

## Decisions
- Use separate workflow files for clarity
- CI runs on all branches; build and deploy only on `main`
- Deploy step explicitly checks for token and prints guidance if missing
