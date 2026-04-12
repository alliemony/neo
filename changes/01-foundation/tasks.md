# Tasks: Project Foundation

1. [x] Initialize `frontend/` with Vite + React + TypeScript template
2. [x] Add Tailwind CSS to frontend with base config
3. [x] Configure Vitest with React Testing Library and jsdom
4. [x] Write first frontend test (`App.test.tsx` -- renders without crashing)
5. [x] Verify `npm run dev` and `npm test` both work
6. [x] Initialize `backend/go.mod` and install Chi router
7. [x] Create `cmd/server/main.go` with server startup and health route
8. [x] Create `internal/config/config.go` for environment-based configuration
9. [x] Create `internal/handler/health.go` with GET `/api/v1/health`
10. [x] Write `internal/handler/health_test.go` using httptest
11. [x] Verify `go run ./cmd/server` and `go test ./...` both work
12. [x] Initialize `widgets/pyproject.toml` with FastAPI, uvicorn, pytest, httpx
13. [x] Create `app/main.py` with FastAPI app and `/health` endpoint
14. [x] Write `tests/test_health.py` using httpx AsyncClient
15. [x] Verify `uvicorn app.main:app` and `pytest` both work
16. [x] Create `.env.example` with all environment variables documented
17. [x] Create `Dockerfile.frontend` (dev mode)
18. [x] Create `Dockerfile.backend` (dev mode)
19. [x] Create `Dockerfile.widgets` (dev mode)
20. [x] Create `docker-compose.yml` with all three services
21. [x] Create `Makefile` with dev, test, build, clean targets
22. [x] Verify `docker compose up --build` starts all services
23. [x] Verify `make test` runs all test suites
