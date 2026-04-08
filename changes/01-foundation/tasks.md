# Tasks: Project Foundation

1. [ ] Initialize `frontend/` with Vite + React + TypeScript template
2. [ ] Add Tailwind CSS to frontend with base config
3. [ ] Configure Vitest with React Testing Library and jsdom
4. [ ] Write first frontend test (`App.test.tsx` -- renders without crashing)
5. [ ] Verify `npm run dev` and `npm test` both work
6. [ ] Initialize `backend/go.mod` and install Chi router
7. [ ] Create `cmd/server/main.go` with server startup and health route
8. [ ] Create `internal/config/config.go` for environment-based configuration
9. [ ] Create `internal/handler/health.go` with GET `/api/v1/health`
10. [ ] Write `internal/handler/health_test.go` using httptest
11. [ ] Verify `go run ./cmd/server` and `go test ./...` both work
12. [ ] Initialize `widgets/pyproject.toml` with FastAPI, uvicorn, pytest, httpx
13. [ ] Create `app/main.py` with FastAPI app and `/health` endpoint
14. [ ] Write `tests/test_health.py` using httpx AsyncClient
15. [ ] Verify `uvicorn app.main:app` and `pytest` both work
16. [ ] Create `.env.example` with all environment variables documented
17. [ ] Create `Dockerfile.frontend` (dev mode)
18. [ ] Create `Dockerfile.backend` (dev mode)
19. [ ] Create `Dockerfile.widgets` (dev mode)
20. [ ] Create `docker-compose.yml` with all three services
21. [ ] Create `Makefile` with dev, test, build, clean targets
22. [ ] Verify `docker compose up --build` starts all services
23. [ ] Verify `make test` runs all test suites
