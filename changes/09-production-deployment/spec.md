# Spec: Production Deployment

## Purpose

Make the site production-ready with optimized builds, PostgreSQL support, SEO, RSS, and platform deployment configuration.

## Requirements

### Requirement: Go backend SHALL serve pre-built frontend static files in production mode

#### Scenario: Static files served from embedded or filesystem directory

- **GIVEN:** Frontend is built to `frontend/dist/`
- **WHEN:** The backend starts with `SERVE_STATIC=true`
- **THEN:** Requests not matching `/api/*` serve files from the static directory
- **AND:** Client-side routing works (all non-file paths return `index.html`)

### Requirement: Production Dockerfiles SHALL use multi-stage builds for minimal image size

#### Scenario: Backend Docker image contains only the binary

- **GIVEN:** The backend Dockerfile uses a multi-stage build
- **WHEN:** The image is built
- **THEN:** The final image is based on `alpine` with only the compiled Go binary
- **AND:** The image size is under 50MB

#### Scenario: Frontend is built and copied into backend image

- **GIVEN:** The production deployment builds frontend first
- **WHEN:** The backend image is built
- **THEN:** The `frontend/dist/` files are embedded or copied into the backend image

### Requirement: Database layer SHALL support both SQLite and PostgreSQL via configuration

#### Scenario: SQLite is used in development

- **GIVEN:** `DATABASE_URL=sqlite://neo.db`
- **WHEN:** The backend starts
- **THEN:** SQLite is used as the database driver

#### Scenario: PostgreSQL is used in production

- **GIVEN:** `DATABASE_URL=postgres://user:pass@host:5432/neo`
- **WHEN:** The backend starts
- **THEN:** PostgreSQL is used as the database driver
- **AND:** The same migrations run against PostgreSQL

### Requirement: All pages SHALL include Open Graph meta tags for social sharing

#### Scenario: Post page has OG meta tags

- **GIVEN:** A published post with title and excerpt
- **WHEN:** The page is shared on social media
- **THEN:** The correct title, description, and URL appear in the social preview

### Requirement: An RSS feed SHALL be available at /api/v1/feed.xml

#### Scenario: RSS feed contains recent posts

- **GIVEN:** 10 published posts exist
- **WHEN:** `GET /api/v1/feed.xml` is called
- **THEN:** Response is valid RSS 2.0 XML with the 10 most recent posts
- **AND:** Each item has title, link, description, and pubDate

### Requirement: Platform deployment config SHALL enable single-command deployment

#### Scenario: Deploy to platform via CLI

- **GIVEN:** Platform config file exists (e.g., `fly.toml`)
- **WHEN:** The deploy command is run (e.g., `flyctl deploy`)
- **THEN:** The application is built, deployed, and accessible at the configured URL

### Requirement: The site SHALL score above 90 on Lighthouse performance audit

#### Scenario: Performance meets threshold

- **GIVEN:** The production site is deployed
- **WHEN:** A Lighthouse audit is run
- **THEN:** Performance score is above 90
- **AND:** Accessibility score is above 90

### Requirement: HTTPS SHALL be enforced in production

#### Scenario: HTTP redirects to HTTPS

- **GIVEN:** The site is deployed with TLS configured
- **WHEN:** A user accesses the site via HTTP
- **THEN:** They are redirected to HTTPS
