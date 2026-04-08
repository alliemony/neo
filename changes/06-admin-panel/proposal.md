# Proposal: Admin Panel

## Why

The site owner needs a way to create, edit, and manage blog posts and pages without touching the database directly. A basic admin panel with authentication makes content management practical from day one.

## What

Build an authenticated admin panel with a post editor (markdown with live preview), page management, and publish/draft controls. Auth starts with HTTP Basic Auth, designed to be swapped for OAuth2 SSO later.

## What Changes

### Backend
- Implement `Authenticator` interface with `BasicAuthenticator` as first implementation
- Create basic auth middleware reading credentials from env vars (hashed)
- Add admin API routes: CRUD for posts and pages
- Protect all `/api/v1/admin/*` routes with auth middleware

### Frontend
- Create admin login page
- Create admin dashboard (list posts and pages)
- Create post editor with markdown textarea + live preview
- Create page editor
- Implement tag input with autocomplete from existing tags
- Implement publish/draft toggle
- Implement delete confirmation

## Approach

Auth is behind an interface (`Authenticator`) so the basic auth implementation can be replaced with OAuth2 without changing handler code. The admin frontend is part of the same SPA, just protected routes. Session is maintained via a cookie after initial basic auth.
