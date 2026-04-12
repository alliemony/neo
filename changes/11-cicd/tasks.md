# Tasks: CI/CD with GitHub Actions

## Workflows

1. [x] Create `.github/workflows/ci.yml` — test frontend, backend, widgets in parallel
2. [x] Create `.github/workflows/build.yml` — build production Docker image on main
3. [x] Create `.github/workflows/deploy.yml` — deploy to Fly.io (fails without secrets)

## Verification

4. [ ] Push to a branch and confirm CI runs all tests
5. [ ] Merge to main and confirm build + deploy workflows trigger
6. [ ] Set `FLY_API_TOKEN` and confirm deploy succeeds
