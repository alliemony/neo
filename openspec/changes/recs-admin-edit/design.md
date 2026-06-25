## Context

Neo's Recs page is entirely static — a `const RECS: RecSection[]` array in `Recs.tsx`. Posts and pages are fully database-backed with admin CRUD. The backend uses Go + Chi, clean architecture (handler → service → repository), and in-code migration maps in `db.go` (not SQL files). The next available migration version is `006`. The frontend calls `/api/v1/*` for public data and `/api/v1/admin/*` for mutations, using a session cookie for admin auth.

The existing 13 hardcoded recs have these fields per item: `name`, `href` (optional), `section` (the group title: "Tools", "Reading", "Misc"), `tag`, `tagBg`, `tagColor`, `desc`.

## Goals / Non-Goals

**Goals:**
- Database-backed recs with full CRUD from the admin dashboard
- Public `/recs` page visually identical to today — just dynamic
- Seed existing hardcoded recs into the DB on first deploy so no content is lost
- Same inline form panel UX being built for music-recs

**Non-Goals:**
- Reordering via drag-and-drop (sort_order field present but no drag UI in v1)
- Rich text / markdown descriptions (plain text only, same as today)
- Tag colour picker (admin types hex values; same fields as today)

## Decisions

### 1. Data model: single `recs` table, section as a plain text column

```sql
CREATE TABLE recs (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL,
  href        TEXT,
  section     TEXT NOT NULL,   -- e.g. "Tools", "Reading", "Misc"
  tag         TEXT NOT NULL,
  tag_bg      TEXT NOT NULL,   -- hex colour e.g. "#edf5f3"
  tag_color   TEXT NOT NULL,   -- hex colour e.g. "#1a7060"
  description TEXT NOT NULL,
  published   BOOLEAN NOT NULL DEFAULT 0,
  sort_order  INTEGER NOT NULL DEFAULT 0,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Alternative considered**: a separate `rec_sections` table. Rejected — sections are just labels; a separate table adds a join for no real benefit on a personal site.

### 2. Grouping happens in the API response, not the DB

`GET /api/v1/recs` returns `{ "sections": [ { "title": "Tools", "items": [...] } ] }` — the backend groups rows by `section` in Go, ordered by `sort_order ASC, created_at ASC` within each group. Section order follows the first appearance of each section name.

**Alternative**: return a flat array and group client-side. Rejected — the current frontend groups by section; keeping that contract in the API makes the frontend simpler.

### 3. Seed existing hardcoded data via the existing `database.Seed()` function

`db.go` already has a `Seed()` function called at startup. Extend it to insert the 13 existing recs if the `recs` table is empty. This is idempotent and survives restarts.

**Alternative**: embed seed data in the migration itself. Rejected — migrations should be schema-only; seed data in a migration can't be easily changed later.

### 4. Admin UI: inline form panel (same pattern as music-recs)

Fields: name*, section*, tag*, tag_bg*, tag_color*, description*, href (optional), published toggle, sort_order. The section field is a text input (not a dropdown) so new sections can be added freely.

### 5. Public response only includes published recs

Unpublished recs are invisible to the public endpoint. The admin endpoint returns all recs regardless of published status.

## Risks / Trade-offs

- **Seed running on every cold start** → Seed is guarded by `WHERE NOT EXISTS` / count check so it's a no-op after first run. No performance concern.
- **Section names are free-text** → Typos create phantom sections. Mitigation: admin UI will show a datalist of existing section names as suggestions.
- **No rollback for seed data** → If the migration is rolled back, seeded rows disappear. Acceptable — this is additive only and a personal site.

## Migration Plan

1. Add `006_create_recs` to both `sqliteMigrations` and `postgresMigrations` maps in `db.go`
2. Extend `Seed()` in `db.go` to insert the 13 existing recs if the table is empty
3. Deploy backend — migration runs, seed runs once
4. Deploy frontend — `Recs.tsx` now fetches from API; identical appearance
5. No rollback needed — additive migration, no existing tables touched
