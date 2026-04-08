# Tasks: Production Deployment

## Backend

1. [ ] Implement static file serving with SPA fallback in Go backend
2. [ ] Write tests for static file serving (serves files, SPA fallback to index.html)
3. [ ] Add PostgreSQL driver support (`pgx`) alongside SQLite
4. [ ] Create dialect-aware migrations (SQLite vs PostgreSQL)
5. [ ] Write tests for database driver selection based on DATABASE_URL
6. [ ] Implement RSS feed handler (GET /api/v1/feed.xml)
7. [ ] Write tests for RSS feed (valid XML, contains recent posts)

## Frontend

8. [ ] Install and configure `react-helmet-async`
9. [ ] Add Open Graph meta tags to post view, page view, and home page
10. [ ] Verify meta tags render correctly via view-source

## Infrastructure

11. [ ] Create unified production Dockerfile (multi-stage: frontend + backend)
12. [ ] Create production Dockerfile for widget service
13. [ ] Create `fly.toml` (or equivalent platform config)
14. [ ] Document deployment steps in README
15. [ ] Run Lighthouse audit and address any scores below 90
16. [ ] Run accessibility audit and fix issues
17. [ ] Test full deployment: build → deploy → verify health → verify site
