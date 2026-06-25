## ADDED Requirements

### Requirement: Database migration creates music_recs table
The system SHALL include a migration that creates the `music_recs` table with columns: id, title, artist, album, cover_url, spotify_url, apple_url, note, published, sort_order, created_at, updated_at.

#### Scenario: Migration runs on startup
- **WHEN** the backend starts with a fresh or existing database
- **THEN** the `music_recs` table exists and is queryable

---

### Requirement: Public endpoint lists published recs
`GET /api/v1/music` SHALL return all recs where `published = true`, ordered by `sort_order ASC, created_at DESC`. Response: `{ "recs": [ MusicRec, ... ] }`.

#### Scenario: Returns only published recs
- **WHEN** `GET /api/v1/music` is called and some recs are drafts
- **THEN** the response contains only published recs

#### Scenario: Returns empty list when no published recs
- **WHEN** `GET /api/v1/music` is called and no recs are published
- **THEN** the response is `{ "recs": [] }` with status 200

---

### Requirement: Admin endpoint lists all recs
`GET /api/v1/admin/music` SHALL return all recs regardless of published status, ordered by `sort_order ASC, created_at DESC`. Requires valid session.

#### Scenario: Returns all recs for admin
- **WHEN** an authenticated admin calls `GET /api/v1/admin/music`
- **THEN** the response contains both published and draft recs

#### Scenario: Unauthorized access is rejected
- **WHEN** `GET /api/v1/admin/music` is called without a valid session cookie
- **THEN** the response is 401 Unauthorized

---

### Requirement: Admin endpoint creates a rec
`POST /api/v1/admin/music` SHALL accept `{ title, artist, album?, cover_url?, spotify_url?, apple_url?, note?, published?, sort_order? }` and return the created rec with status 201. title and artist are required.

#### Scenario: Valid create request succeeds
- **WHEN** an authenticated admin POSTs a body with title and artist
- **THEN** the response is 201 with the created rec including its new id

#### Scenario: Missing required fields returns 400
- **WHEN** an authenticated admin POSTs a body missing title or artist
- **THEN** the response is 400 Bad Request with an error message

---

### Requirement: Admin endpoint updates a rec
`PUT /api/v1/admin/music/{id}` SHALL accept the same body as create (all fields optional) and return the updated rec.

#### Scenario: Valid update request succeeds
- **WHEN** an authenticated admin PUTs a partial or full body for an existing rec id
- **THEN** the response is 200 with the updated rec

#### Scenario: Update on non-existent id returns 404
- **WHEN** an authenticated admin PUTs to an id that does not exist
- **THEN** the response is 404 Not Found

---

### Requirement: Admin endpoint deletes a rec
`DELETE /api/v1/admin/music/{id}` SHALL permanently remove the rec and return 204 No Content.

#### Scenario: Valid delete succeeds
- **WHEN** an authenticated admin DELETEs an existing rec id
- **THEN** the response is 204 and the rec no longer appears in any list endpoint

#### Scenario: Delete on non-existent id returns 404
- **WHEN** an authenticated admin DELETEs an id that does not exist
- **THEN** the response is 404 Not Found
