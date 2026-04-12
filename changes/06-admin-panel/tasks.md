# Tasks: Admin Panel

## Backend

1. [x] Define `Authenticator` interface with `Middleware()` and `CurrentUser()` methods
2. [x] Write tests for `BasicAuthenticator` (valid creds, invalid creds, missing header)
3. [x] Implement `BasicAuthenticator` with bcrypt password verification
4. [x] Write tests for admin post handlers (Create, Update, Delete -- with auth)
5. [x] Implement admin post handlers
6. [x] Create pages table migration SQL
7. [x] Write tests for `PageRepository` and `PageService`
8. [x] Implement `PageRepository` and `PageService`
9. [x] Write tests for admin page handlers (Create, Update, Delete)
10. [x] Implement admin page handlers
11. [x] Register all admin routes behind auth middleware

## Frontend

12. [x] Write tests for admin login page (submit, error display)
13. [x] Implement admin login page with credential form
14. [x] Create auth context for storing credentials in memory
15. [x] Write tests for admin dashboard (lists posts and pages, status indicators)
16. [x] Implement admin dashboard component
17. [x] Write tests for post editor (field population, preview rendering, save)
18. [x] Implement post editor with live markdown preview
19. [x] Write tests for tag autocomplete input
20. [x] Implement tag autocomplete input component
21. [x] Implement publish/draft toggle
22. [x] Implement delete confirmation dialog
23. [x] Implement page editor (reuse post editor patterns)
24. [ ] Add protected admin routes to React Router
25. [ ] Verify full admin flow: login → create post → edit → publish → delete
