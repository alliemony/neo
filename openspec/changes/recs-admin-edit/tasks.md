## 1. Database Migration & Seed

- [x] 1.1 Add `"006_create_recs"` to `sqliteMigrations` in `backend/internal/database/db.go` with the `CREATE TABLE recs` DDL (id, name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at)
- [x] 1.2 Add the same migration to `postgresMigrations` using `SERIAL`, `BOOLEAN`, and `TIMESTAMP` types
- [x] 1.3 Extend `database.Seed()` in `db.go` to insert the 13 existing hardcoded recs (from `Recs.tsx`) if `SELECT COUNT(*) FROM recs` returns 0 — all seeded as `published = true`
- [x] 1.4 Verify migration and seed run cleanly: `go run ./cmd/server`, confirm rows exist

## 2. Backend — Repository

- [x] 2.1 Create `backend/internal/repository/recs_repo.go` defining `Rec` struct and `RecsRepo` interface with: `ListAll()`, `ListPublished()`, `GetByID(id)`, `Create(rec)`, `Update(id, rec)`, `Delete(id)`
- [x] 2.2 Implement the SQLite/Postgres-compatible `recsRepo` struct (use `database.CurrentDialect` for any dialect-specific SQL if needed)
- [x] 2.3 Write `backend/internal/repository/recs_repo_test.go` covering: ListPublished filters drafts, Create roundtrip, Update, Delete, GetByID not-found returns error

## 3. Backend — Service

- [x] 3.1 Create `backend/internal/service/recs_service.go` with `RecsService` wrapping the repo; validate that name, section, tag, tag_bg, tag_color, and description are non-empty on Create/Update
- [x] 3.2 Write `backend/internal/service/recs_service_test.go` covering validation errors and happy paths (mock repo)

## 4. Backend — Handler & Routes

- [x] 4.1 Create `backend/internal/handler/recs_handler.go` with handlers: `ListPublic` (grouped by section), `AdminList`, `AdminCreate`, `AdminUpdate`, `AdminDelete`
- [x] 4.2 Write `backend/internal/handler/recs_handler_test.go` covering: public grouping, 400 on missing fields, 404 on unknown id, 401 on unauthenticated admin
- [x] 4.3 Register routes in `backend/cmd/server/main.go`: `GET /api/v1/recs` (public) and `GET/POST /api/v1/admin/recs`, `PUT/DELETE /api/v1/admin/recs/{id}` (protected)

## 5. Frontend — Types & API Service

- [x] 5.1 Add `Rec` and `RecSection` types to `frontend/src/types/` (or extend existing `post.ts`) matching backend fields (name, href, section, tag, tagBg, tagColor, description, published, sortOrder, id)
- [x] 5.2 Create `frontend/src/services/recsApi.ts` with `listRecs()` (public, returns `{ sections }`) and admin functions: `adminListRecs()`, `adminCreateRec()`, `adminUpdateRec()`, `adminDeleteRec()`

## 6. Frontend — Public Recs Page

- [x] 6.1 Replace the hardcoded `RECS` constant in `frontend/src/routes/Recs.tsx` with a `useEffect` + `useState` that calls `listRecs()`
- [x] 6.2 Add loading state (spinner or skeleton) while fetch is in flight
- [x] 6.3 Add empty state message when `sections` is empty
- [x] 6.4 Add error state message on fetch failure
- [ ] 6.5 Verify the rendered output is visually identical to the previous hardcoded version (tag pill colours, links, descriptions)

## 7. Frontend — Admin Recs Section

- [x] 7.1 Add a Recs section to `frontend/src/routes/admin/AdminDashboard.tsx` that fetches all recs via `adminListRecs()` and lists them with name, section badge, published badge, Edit and Delete buttons
- [x] 7.2 Build the inline add/edit form panel: fields for name*, section* (text + `<datalist>` of existing sections), tag*, tag_bg*, tag_color*, description*, href, published toggle, sort_order — validate required fields before submit
- [x] 7.3 Wire Create: "+ New Rec" button opens blank form; submit calls `adminCreateRec()` and refreshes list
- [x] 7.4 Wire Edit: Edit button pre-fills form with rec data; submit calls `adminUpdateRec()` and refreshes list
- [x] 7.5 Wire publish toggle: clicking the published badge calls `adminUpdateRec()` with flipped `published` and updates list
- [x] 7.6 Wire Delete: Delete button opens existing confirm modal; confirm calls `adminDeleteRec()` and refreshes list

## 8. Verification & Docs

- [x] 8.1 Run `cd backend && go test ./...` — all tests pass
- [x] 8.2 Run `cd frontend && npm test` — all tests pass
- [ ] 8.3 Deploy to Fly, confirm seed data appears on `/recs`, add a new rec via admin, toggle published, edit it, delete it
- [x] 8.4 Update `docs/references/architecture.md` with `GET /api/v1/recs` and `/api/v1/admin/recs` route groups
