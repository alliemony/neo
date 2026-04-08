# Proposal: Static Pages & Navigation

## Why

A personal website needs more than just blog posts. Static pages (About, Projects, etc.) provide fixed content, and the navigation system ties everything together. The tag cloud in the sidebar helps with content discovery.

## What

Implement the public pages API and frontend rendering, dynamic navigation based on published pages, a proper footer, 404 handling, and the tag cloud sidebar widget.

## What Changes

### Backend
- Add `GET /api/v1/pages/:slug` public endpoint
- Add `GET /api/v1/pages` for listing published pages (used by navigation)

### Frontend
- Create page rendering component (reuses markdown renderer)
- Build dynamic navigation header that reads from pages API
- Create tag cloud sidebar widget
- Create footer component with site info
- Implement proper 404 page
- Add `/page/:slug` route

## Approach

Pages reuse the same markdown rendering as blog posts. Navigation links are fetched from the pages API on initial load. The tag cloud uses the existing `/api/v1/tags` endpoint.
