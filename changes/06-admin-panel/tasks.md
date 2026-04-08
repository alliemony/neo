# Tasks: Admin Panel

## Backend

1. [ ] Define `Authenticator` interface with `Middleware()` and `CurrentUser()` methods
2. [ ] Write tests for `BasicAuthenticator` (valid creds, invalid creds, missing header)
3. [ ] Implement `BasicAuthenticator` with bcrypt password verification
4. [ ] Write tests for admin post handlers (Create, Update, Delete -- with auth)
5. [ ] Implement admin post handlers
6. [ ] Create pages table migration SQL
7. [ ] Write tests for `PageRepository` and `PageService`
8. [ ] Implement `PageRepository` and `PageService`
9. [ ] Write tests for admin page handlers (Create, Update, Delete)
10. [ ] Implement admin page handlers
11. [ ] Register all admin routes behind auth middleware

## Frontend

12. [ ] Write tests for admin login page (submit, error display)
13. [ ] Implement admin login page with credential form
14. [ ] Create auth context for storing credentials in memory
15. [ ] Write tests for admin dashboard (lists posts and pages, status indicators)
16. [ ] Implement admin dashboard component
17. [ ] Write tests for post editor (field population, preview rendering, save)
18. [ ] Implement post editor with live markdown preview
19. [ ] Write tests for tag autocomplete input
20. [ ] Implement tag autocomplete input component
21. [ ] Implement publish/draft toggle
22. [ ] Implement delete confirmation dialog
23. [ ] Implement page editor (reuse post editor patterns)
24. [ ] Add protected admin routes to React Router
25. [ ] Verify full admin flow: login → create post → edit → publish → delete
