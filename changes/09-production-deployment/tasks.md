# Tasks: Production Deployment

## Backend

1. [x] Implement static file serving with SPA fallback in Go backend
2. [x] Write tests for static file serving (serves files, SPA fallback to index.html)
3. [x] Add PostgreSQL driver support (`pgx`) alongside SQLite
4. [x] Create dialect-aware migrations (SQLite vs PostgreSQL)
5. [x] Write tests for database driver selection based on DATABASE_URL
6. [x] Implement RSS feed handler (GET /api/v1/feed.xml)
7. [x] Write tests for RSS feed (valid XML, contains recent posts)

## Frontend

8. [x] Install and configure `react-helmet-async`
9. [x] Add Open Graph meta tags to post view, page view, and home page
10. [ ] Verify meta tags render correctly via view-source

## Infrastructure

11. [x] Create unified production Dockerfile (multi-stage: frontend + backend)
12. [x] Create production Dockerfile for widget service
13. [x] Create `fly.toml` (or equivalent platform config)
14. [x] Document deployment steps in README
15. [ ] Run Lighthouse audit and address any scores below 90
16. [ ] Run accessibility audit and fix issues
17. [ ] Test full deployment: build → deploy → verify health → verify site
