# Tasks: Comments & Interactions

## Backend

1. [x] Create comments table migration SQL
2. [x] Add `like_count` column to posts table migration
3. [x] Write tests for `CommentRepository` (Create, ListByPostSlug)
4. [x] Implement `CommentRepository`
5. [x] Write tests for `CommentService` (validation: empty name, empty content, max length)
6. [x] Implement `CommentService`
7. [x] Write tests for comment HTTP handlers (GET list, POST create, 404 on bad slug)
8. [x] Implement comment HTTP handlers
9. [x] Write tests for like endpoint (increment, 404 on bad slug)
10. [x] Implement like handler (`POST /api/v1/posts/:slug/like`)
11. [x] Write tests for rate limiter middleware
12. [x] Implement rate limiter middleware
13. [x] Register comment and like routes

## Frontend

14. [x] Add comment and like API client methods
15. [x] Write tests for `useComments` hook
16. [x] Implement `useComments` hook
17. [x] Write tests for `CommentSection` component (renders comments, shows form)
18. [x] Implement `CommentSection` component
19. [x] Write tests for `CommentForm` (submit, validation display)
20. [x] Implement `CommentForm` component
21. [x] Write tests for `LikeButton` (optimistic update, revert on error)
22. [x] Implement `LikeButton` component
23. [x] Integrate CommentSection into PostView sidebar
24. [ ] Verify full flow: view post → add comment → like → see updates
