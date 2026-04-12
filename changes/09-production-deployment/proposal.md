# Proposal: Production Deployment

## Why

The site needs to move from local development to a publicly accessible deployment. This includes production builds, database migration to PostgreSQL, HTTPS, SEO basics, and platform configuration.

## What

Prepare all three services for production deployment: multi-stage Docker builds, Go backend serving static frontend files, PostgreSQL support, SEO meta tags, RSS feed, and deployment configuration for a common platform.

## What Changes

- Production multi-stage Dockerfiles for all services
- Go backend serves pre-built frontend static files (single origin)
- Database abstraction supports both SQLite and PostgreSQL
- Environment-based database driver selection
- Meta tags and Open Graph tags for social sharing
- RSS feed endpoint (`GET /api/v1/feed.xml`)
- Platform deployment config (Fly.io or similar)
- HTTPS/TLS documentation
- Performance and accessibility audits

## Approach

In production, the Go backend serves everything: the API, the static frontend, and acts as a reverse proxy to the widget service if needed. This simplifies deployment to a single primary service plus the optional Python widget sidecar.
