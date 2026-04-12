# Proposal: Blog Backend (Posts & Tags API)

## Why

The blog is the core feature of the site. The backend must provide a RESTful API for creating, reading, and filtering posts before the frontend can display them. This change covers the full backend stack for posts: database schema, repository, service, and HTTP handlers.

## What

Implement the posts and tags data layer and API endpoints. This includes database migrations, the post model, repository for data access, service for business logic, and HTTP handlers for the public API.

## What Changes

- Create database connection and migration infrastructure
- Define the `posts` table schema with slug, title, content, tags, timestamps
- Implement `PostRepository` for CRUD operations
- Implement `PostService` for business logic (slug generation, validation, tag filtering)
- Implement HTTP handlers for `GET /api/v1/posts`, `GET /api/v1/posts/:slug`, `GET /api/v1/tags`
- Add seed data for development (3-5 sample posts)
- CORS middleware for frontend communication

## Approach

Clean architecture: handlers depend on services, services depend on repositories, repositories depend on the database. Each layer is testable independently via interfaces. Start with SQLite for simplicity.
