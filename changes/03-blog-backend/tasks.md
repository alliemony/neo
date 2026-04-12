# Tasks: Blog Backend (Posts & Tags API)

1. [x] Create `database/db.go` with SQLite connection and migration runner
2. [x] Create posts table migration SQL
3. [x] Write tests for `PostRepository` (Create, GetBySlug, List, ListByTag, AllTags)
4. [x] Implement `PostRepository` with SQL queries
5. [x] Write tests for `PostService` (Create with slug gen, validation, list published, filter by tag)
6. [x] Implement `PostService` with business logic
7. [x] Implement `slugify()` utility with tests
8. [x] Write tests for `PostHandler` (List, GetBySlug, ListTags -- using httptest)
9. [x] Implement `PostHandler` HTTP handlers
10. [x] Add CORS middleware configuration
11. [x] Register post routes in `cmd/server/main.go`
12. [x] Create seed data loader for development
13. [ ] Write integration test: start server, hit API, verify JSON response
14. [x] Verify all `go test ./...` pass
