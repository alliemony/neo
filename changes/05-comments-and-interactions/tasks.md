# Tasks: Comments & Interactions

## Backend

1. [ ] Create comments table migration SQL
2. [ ] Add `like_count` column to posts table migration
3. [ ] Write tests for `CommentRepository` (Create, ListByPostSlug)
4. [ ] Implement `CommentRepository`
5. [ ] Write tests for `CommentService` (validation: empty name, empty content, max length)
6. [ ] Implement `CommentService`
7. [ ] Write tests for comment HTTP handlers (GET list, POST create, 404 on bad slug)
8. [ ] Implement comment HTTP handlers
9. [ ] Write tests for like endpoint (increment, 404 on bad slug)
10. [ ] Implement like handler (`POST /api/v1/posts/:slug/like`)
11. [ ] Write tests for rate limiter middleware
12. [ ] Implement rate limiter middleware
13. [ ] Register comment and like routes

## Frontend

14. [ ] Add comment and like API client methods
15. [ ] Write tests for `useComments` hook
16. [ ] Implement `useComments` hook
17. [ ] Write tests for `CommentSection` component (renders comments, shows form)
18. [ ] Implement `CommentSection` component
19. [ ] Write tests for `CommentForm` (submit, validation display)
20. [ ] Implement `CommentForm` component
21. [ ] Write tests for `LikeButton` (optimistic update, revert on error)
22. [ ] Implement `LikeButton` component
23. [ ] Integrate CommentSection into PostView sidebar
24. [ ] Verify full flow: view post → add comment → like → see updates
