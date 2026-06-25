## Why

Neo is a personal web garden but currently has no way to share music taste — a natural expression of personality alongside posts and notebook entries. Adding a Music tab fills this gap and gives the site more character.

## What Changes

- New public **Music** tab at `/music` listing curated song/album recommendations
- Each rec displays title, artist, optional album art, and streaming service pill links (Spotify, Apple Music) that can be toggled on or off per rec
- Expandable lyrics section per rec, fetched client-side from a free lyrics API (lyrics.ovh)
- Admin dashboard gains a **Music** section to create, edit, publish/unpublish, and delete recs
- New `music_recs` table in the database with a matching backend CRUD API

## Capabilities

### New Capabilities

- `music-recs-public`: Public music tab displaying published recommendations with streaming pills and expandable lyrics
- `music-recs-admin`: Admin CRUD interface for managing music recommendations
- `music-recs-api`: Backend REST endpoints and database schema for music recommendations

### Modified Capabilities

<!-- none -->

## Impact

- **Backend**: new migration, new repository/service/handler files, two new route groups (`GET /api/v1/music` public, `/api/v1/admin/music` protected)
- **Frontend**: new route `/music`, new admin sub-section, new `musicApi` service, new `MusicRec` type
- **Database**: additive migration — no existing tables touched
- **Dependencies**: no new backend deps; frontend fetches lyrics.ovh directly (no SDK needed)
