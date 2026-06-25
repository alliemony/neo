## 1. Database & Migration

- [x] 1.1 Find the current highest migration number in `backend/internal/database/migrations/` and create `00N_music_recs.sql` with the `CREATE TABLE music_recs` DDL from design.md
- [x] 1.2 Verify migration runs cleanly on a fresh `go run ./cmd/server` (table appears in DB)

## 2. Backend — Repository

- [x] 2.1 Create `backend/internal/repository/music_repo.go` with `MusicRec` struct and `MusicRepo` interface: `List(publishedOnly bool)`, `GetByID(id int)`, `Create(rec)`, `Update(id, rec)`, `Delete(id)`
- [x] 2.2 Implement the SQLite-backed `musicRepo` struct satisfying that interface
- [x] 2.3 Write `backend/internal/repository/music_repo_test.go` covering List (published filter), Create, Update, Delete, and GetByID not-found

## 3. Backend — Service

- [x] 3.1 Create `backend/internal/service/music_service.go` with `MusicService` wrapping the repo; validate that title and artist are non-empty on Create/Update
- [x] 3.2 Write `backend/internal/service/music_service_test.go` covering validation errors and happy paths

## 4. Backend — Handler & Routes

- [x] 4.1 Create `backend/internal/handler/music_handler.go` with handlers: `ListPublic`, `AdminList`, `AdminCreate`, `AdminUpdate`, `AdminDelete`
- [x] 4.2 Write `backend/internal/handler/music_handler_test.go` covering each handler (400 on missing fields, 404 on unknown id, 401 on unauthenticated admin routes)
- [x] 4.3 Register routes in `backend/cmd/server/main.go`: `GET /api/v1/music` (public) and `GET/POST /api/v1/admin/music`, `PUT/DELETE /api/v1/admin/music/{id}` (protected)

## 5. Frontend — Types & API Service

- [x] 5.1 Add `MusicRec` type to `frontend/src/types/post.ts` (or a new `music.ts`) matching the backend struct fields
- [x] 5.2 Create `frontend/src/services/musicApi.ts` with `listMusic()` (public) and admin functions `adminListMusic()`, `adminCreateMusic()`, `adminUpdateMusic()`, `adminDeleteMusic()`

## 6. Frontend — Public Music Page

- [x] 6.1 Create `frontend/src/routes/Music.tsx` that fetches published recs and renders a card list with empty state
- [x] 6.2 Implement streaming service pills (Spotify / Apple Music) — render only when URL is non-empty; opens in new tab
- [x] 6.3 Implement collapsible lyrics section: fetches from `https://api.lyrics.ovh/v1/{artist}/{title}` on first open, caches result in component state, shows loading spinner and graceful error message
- [x] 6.4 Render optional `note` field on the card when present
- [x] 6.5 Add the Music route to `frontend/src/App.tsx` at path `/music`
- [x] 6.6 Add Music link to the site navigation

## 7. Frontend — Admin Music Section

- [x] 7.1 Add a Music section to `frontend/src/routes/admin/AdminDashboard.tsx` that lists all recs (title, artist, published badge) fetched from `adminListMusic()`
- [x] 7.2 Build the inline add/edit form panel with fields: title*, artist*, album, cover URL, Spotify URL, Apple Music URL, note, published toggle, sort order — validates title and artist required before submit
- [x] 7.3 Wire Create: submitting a blank form opens it; submitting the filled form calls `adminCreateMusic()` and refreshes the list
- [x] 7.4 Wire Edit: clicking Edit pre-fills the form with the rec's data; saving calls `adminUpdateMusic()`
- [x] 7.5 Wire publish toggle: clicking the badge calls `adminUpdateMusic()` with flipped `published` and updates optimistically
- [x] 7.6 Wire Delete: clicking Delete opens the existing confirm modal; confirming calls `adminDeleteMusic()` and refreshes

## 8. Polish & Verification

- [x] 8.1 Run `cd backend && go test ./...` — all tests pass
- [x] 8.2 Run `cd frontend && npm test` — all tests pass
- [ ] 8.3 Deploy to Fly and manually verify: add a rec in admin, toggle published, view it on `/music`, expand lyrics, click a streaming pill
- [ ] 8.4 Update `docs/references/architecture.md` with the new `/music` and `/api/v1/admin/music` route groups
