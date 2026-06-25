## Why

The Recs page is hardcoded in the frontend — adding, editing, or removing a recommendation requires a code change and redeploy. Moving recs to the database and wiring up admin CRUD makes the site a true CMS for all content types, consistent with how posts and pages already work.

## What Changes

- New `recs` table in the database storing all recommendation fields (name, href, section, tag colours, description, published, sort_order)
- The existing hardcoded rec data is seeded into the DB via a one-time seed so no content is lost
- **BREAKING** (frontend-only): `Recs.tsx` replaces the hardcoded `RECS` constant with an API fetch from `GET /api/v1/recs`
- New public endpoint `GET /api/v1/recs` returns published recs grouped by section
- New admin endpoints `GET/POST /api/v1/admin/recs` and `PUT/DELETE /api/v1/admin/recs/{id}` for CRUD
- Admin dashboard gains a Recs section — inline form panel to create/edit, publish toggle, delete with confirm — same UX pattern as the Music recs feature

## Capabilities

### New Capabilities

- `recs-api`: Backend migration, repository, service, and handler for recs CRUD
- `recs-admin`: Admin dashboard section for managing recs (create, edit, publish, delete)

### Modified Capabilities

- `recs-public`: The public `/recs` page changes from static to API-driven — same visual output, dynamic data source

## Impact

- **Backend**: `db.go` gains migration `006_create_recs` (SQLite + Postgres), new files `recs_repo.go`, `recs_service.go`, `recs_handler.go`, routes wired in `main.go`
- **Frontend**: `Recs.tsx` switches from hardcoded constant to `useEffect` fetch; new `recsApi.ts` service; `AdminDashboard.tsx` gains Recs section; `MusicRec`-style inline form added
- **Seed data**: existing 13 hardcoded recs seeded via a DB seed function so content survives the migration
- **No external dependencies added**
