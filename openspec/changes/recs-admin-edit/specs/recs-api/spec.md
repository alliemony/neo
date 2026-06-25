## ADDED Requirements

### Requirement: Migration creates recs table
The system SHALL include migration `006_create_recs` (SQLite and Postgres variants) that creates the `recs` table with columns: id, name, href, section, tag, tag_bg, tag_color, description, published, sort_order, created_at, updated_at.

#### Scenario: Migration runs on startup
- **WHEN** the backend starts against a database without the `recs` table
- **THEN** the table is created and queryable before any request is served

---

### Requirement: Existing hardcoded recs are seeded into the database
The `Seed()` function SHALL insert the 13 existing recs (Tools, Reading, Misc sections) if the `recs` table is empty, so no content is lost after migration.

#### Scenario: Seed is idempotent
- **WHEN** the backend starts multiple times against a non-empty `recs` table
- **THEN** no duplicate rows are inserted

#### Scenario: Seed inserts all recs on first run
- **WHEN** the backend starts against an empty `recs` table
- **THEN** all 13 recs are present and published

---

### Requirement: Public endpoint returns published recs grouped by section
`GET /api/v1/recs` SHALL return `{ "sections": [ { "title": string, "items": [ Rec, ... ] } ] }` containing only published recs, each section ordered by sort_order ASC then created_at ASC within. Section order follows first appearance.

#### Scenario: Returns grouped published recs
- **WHEN** `GET /api/v1/recs` is called and published recs exist across multiple sections
- **THEN** the response groups them by section with correct ordering

#### Scenario: Returns empty sections array when no published recs
- **WHEN** `GET /api/v1/recs` is called and no recs are published
- **THEN** the response is `{ "sections": [] }` with status 200

---

### Requirement: Admin endpoint lists all recs
`GET /api/v1/admin/recs` SHALL return `{ "recs": [ Rec, ... ] }` — all recs regardless of published status, ordered by sort_order ASC then created_at ASC. Requires valid session.

#### Scenario: Authenticated admin receives all recs
- **WHEN** an authenticated admin calls `GET /api/v1/admin/recs`
- **THEN** the response includes both published and draft recs

#### Scenario: Unauthenticated request is rejected
- **WHEN** `GET /api/v1/admin/recs` is called without a valid session cookie
- **THEN** the response is 401 Unauthorized

---

### Requirement: Admin endpoint creates a rec
`POST /api/v1/admin/recs` SHALL accept `{ name, section, tag, tag_bg, tag_color, description, href?, published?, sort_order? }` and return the created rec with status 201. name, section, tag, tag_bg, tag_color, and description are required.

#### Scenario: Valid create returns 201 with new rec
- **WHEN** an authenticated admin POSTs all required fields
- **THEN** the response is 201 with the full rec including its new id

#### Scenario: Missing required field returns 400
- **WHEN** an authenticated admin POSTs a body missing any required field
- **THEN** the response is 400 Bad Request

---

### Requirement: Admin endpoint updates a rec
`PUT /api/v1/admin/recs/{id}` SHALL accept any subset of rec fields and return the updated rec with status 200.

#### Scenario: Valid update returns 200
- **WHEN** an authenticated admin PUTs a partial body for an existing id
- **THEN** the response is 200 with the updated rec

#### Scenario: Unknown id returns 404
- **WHEN** an authenticated admin PUTs to an id that does not exist
- **THEN** the response is 404 Not Found

---

### Requirement: Admin endpoint deletes a rec
`DELETE /api/v1/admin/recs/{id}` SHALL permanently remove the rec and return 204 No Content.

#### Scenario: Valid delete returns 204
- **WHEN** an authenticated admin DELETEs an existing id
- **THEN** the response is 204 and the rec no longer appears in any list

#### Scenario: Unknown id returns 404
- **WHEN** an authenticated admin DELETEs an id that does not exist
- **THEN** the response is 404 Not Found
