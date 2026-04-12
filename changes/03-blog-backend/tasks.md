# Tasks: Blog Backend (Posts & Tags API)

1. [ ] Create `database/db.go` with SQLite connection and migration runner
2. [ ] Create posts table migration SQL
3. [ ] Write tests for `PostRepository` (Create, GetBySlug, List, ListByTag, AllTags)
4. [ ] Implement `PostRepository` with SQL queries
5. [ ] Write tests for `PostService` (Create with slug gen, validation, list published, filter by tag)
6. [ ] Implement `PostService` with business logic
7. [ ] Implement `slugify()` utility with tests
8. [ ] Write tests for `PostHandler` (List, GetBySlug, ListTags -- using httptest)
9. [ ] Implement `PostHandler` HTTP handlers
10. [ ] Add CORS middleware configuration
11. [ ] Register post routes in `cmd/server/main.go`
12. [ ] Create seed data loader for development
13. [ ] Write integration test: start server, hit API, verify JSON response
14. [ ] Verify all `go test ./...` pass
