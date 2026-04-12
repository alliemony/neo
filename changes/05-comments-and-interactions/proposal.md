# Proposal: Comments & Interactions

## Why

User engagement is what makes a blog feel alive. Comments and likes give readers a way to interact with content. The comment section in the sidebar (Tumblr-inspired) is a key part of the single post experience.

## What

Implement the full comment system (backend + frontend) and the post like/heart feature. Includes rate limiting to prevent spam.

## What Changes

### Backend
- Create `comments` table migration
- Implement `CommentRepository` and `CommentService`
- Add `GET /api/v1/posts/:slug/comments` and `POST /api/v1/posts/:slug/comments`
- Add `POST /api/v1/posts/:slug/like` endpoint (anonymous, count-only)
- Add `like_count` field to posts table
- Implement IP-based rate limiting middleware for comment submissions

### Frontend
- Create `CommentSection` component (displays in sidebar on desktop, below post on mobile)
- Create comment submission form (name + content fields)
- Create like/heart button with optimistic UI update
- Integrate comments into single post view
- Update PostCard to show live like/comment counts

## Approach

Comments are a linear thread (not nested) for simplicity. Likes are anonymous (no user tracking, just a counter). Rate limiting uses a simple in-memory token bucket per IP, configurable via env vars.
