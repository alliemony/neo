## Context

Neo currently has posts, pages, and widgets. The stack is Go (Chi) backend + React/Vite frontend + SQLite in dev / Postgres in prod. Auth is GitHub OAuth with a session cookie. Admin routes are protected by the `OAuthAuthenticator` middleware. The frontend calls `GET /api/v1/*` for public data and `/api/v1/admin/*` for mutations.

## Goals / Non-Goals

**Goals:**
- Public `/music` page showing published recs with streaming pills and an expandable lyrics section
- Admin CRUD for music recs (create, edit, toggle published, delete)
- Lyrics fetched client-side from a free API — no server-side storage of lyrics
- Streaming service links are optional per rec (omit pill if URL is empty)

**Non-Goals:**
- Spotify/Apple Music OAuth or playback embedding (iframes / SDKs) — links only
- Server-side lyrics caching or storage
- Bulk import, playlist management, or tags on recs

## Decisions

### 1. Lyrics API: lyrics.ovh

lyrics.ovh offers a simple `GET https://api.lyrics.ovh/v1/{artist}/{title}` with no auth or rate-limit for personal use. The response is `{ lyrics: "..." }` or an error. Fetched lazily on accordion expand.

**Alternative considered**: Musixmatch — requires API key and has a 30% lyrics cap on free tier. Not worth the complexity for a personal site.

### 2. Data model: single `music_recs` table

```sql
CREATE TABLE music_recs (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  title        TEXT NOT NULL,
  artist       TEXT NOT NULL,
  album        TEXT,
  cover_url    TEXT,
  spotify_url  TEXT,
  apple_url    TEXT,
  note         TEXT,      -- optional personal note shown on the card
  published    BOOLEAN NOT NULL DEFAULT 0,
  sort_order   INTEGER NOT NULL DEFAULT 0,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Alternative**: storing streaming links in a separate junction table for extensibility. Rejected — there are exactly two services; a flat schema is simpler and easy to extend later if needed.

### 3. Lyrics fetched client-side only

Storing lyrics server-side raises copyright concerns. Fetching on demand keeps the backend clean and means we never cache lyrics.

### 4. Sort order: explicit integer, not created_at

Lets the admin reorder recs without changing timestamps. Default 0 means newest-first when sort_order is tied (ORDER BY sort_order ASC, created_at DESC).

### 5. Admin UI: inline form panel (not a separate editor route)

Posts and pages need a full markdown editor, so they use a dedicated route. Music recs have ~6 fields — an inline slide-down form in the admin dashboard is sufficient and keeps things simple.

## Risks / Trade-offs

- **lyrics.ovh reliability** → API is community-run and may have downtime or missing songs. Mitigation: show a graceful "Lyrics unavailable" message on error; the expand still works.
- **SQLite migration ordering** → Must ensure the migration runs before any query touches the table. Mitigation: use the existing sequential migration system; new file gets the next sequence number.
- **No reorder UI** → sort_order is in the schema but there's no drag-and-drop in v1. Mitigation: admin can set it via the form field; good enough for a personal site.

## Migration Plan

1. Add migration file `005_music_recs.sql` (or next available number)
2. Deploy backend — migration runs on startup
3. Deploy frontend — `/music` route goes live (empty state until recs are added)
4. No rollback needed — additive only, no existing tables modified
