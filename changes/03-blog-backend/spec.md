# Spec: Blog Backend (Posts & Tags API)

## Purpose

Provide the backend data layer and REST API for blog posts and tags, enabling the frontend to fetch, display, and filter published content.

## Requirements

### Requirement: The database SHALL have a posts table with slug, title, content, tags, published flag, and timestamps

#### Scenario: Posts table is created via migration

- **GIVEN:** The backend starts for the first time
- **WHEN:** Database migrations run
- **THEN:** A `posts` table exists with columns: id, slug (unique), title, content, content_type, tags (JSON text), published (boolean), created_at, updated_at

### Requirement: PostRepository SHALL provide CRUD operations for posts

#### Scenario: Create a post

- **GIVEN:** A valid post input with title "Hello World" and content "First post"
- **WHEN:** `repo.Create(post)` is called
- **THEN:** A new row is inserted with an auto-generated slug "hello-world"
- **AND:** The returned post has a non-zero ID and timestamps

#### Scenario: Get post by slug

- **GIVEN:** A post exists with slug "hello-world"
- **WHEN:** `repo.GetBySlug("hello-world")` is called
- **THEN:** The full post record is returned

#### Scenario: Get post by slug returns not found

- **GIVEN:** No post exists with slug "nonexistent"
- **WHEN:** `repo.GetBySlug("nonexistent")` is called
- **THEN:** An `ErrNotFound` error is returned

#### Scenario: List published posts with pagination

- **GIVEN:** 15 published posts and 3 draft posts exist
- **WHEN:** `repo.List(page=1, perPage=10, publishedOnly=true)` is called
- **THEN:** 10 published posts are returned ordered by created_at descending
- **AND:** Draft posts are excluded

### Requirement: PostService SHALL validate inputs and generate slugs

#### Scenario: Slug is generated from title

- **GIVEN:** A create request with title "My Great Post!"
- **WHEN:** `service.Create(input)` is called
- **THEN:** The slug is set to "my-great-post"

#### Scenario: Duplicate slug is rejected

- **GIVEN:** A post with slug "hello-world" already exists
- **WHEN:** `service.Create(input)` is called with title "Hello World"
- **THEN:** An `ErrSlugExists` error is returned

#### Scenario: Empty title is rejected

- **GIVEN:** A create request with an empty title
- **WHEN:** `service.Create(input)` is called
- **THEN:** An `ErrTitleRequired` error is returned

### Requirement: Posts SHALL be filterable by tag

#### Scenario: Filter posts by tag

- **GIVEN:** 3 posts tagged "python" and 2 posts tagged "go"
- **WHEN:** `repo.ListByTag("python", page=1, perPage=10)` is called
- **THEN:** Only the 3 "python" posts are returned

### Requirement: GET /api/v1/posts SHALL return paginated published posts

#### Scenario: List posts returns JSON array

- **GIVEN:** Published posts exist in the database
- **WHEN:** `GET /api/v1/posts` is called
- **THEN:** Response status is 200
- **AND:** Body is a JSON object with `posts` array and `total` count

#### Scenario: Filter by tag via query parameter

- **GIVEN:** Posts with various tags exist
- **WHEN:** `GET /api/v1/posts?tag=python` is called
- **THEN:** Only posts tagged "python" are returned

### Requirement: GET /api/v1/posts/:slug SHALL return a single post

#### Scenario: Valid slug returns post

- **GIVEN:** A published post with slug "hello-world" exists
- **WHEN:** `GET /api/v1/posts/hello-world` is called
- **THEN:** Response status is 200 with the full post as JSON

#### Scenario: Invalid slug returns 404

- **GIVEN:** No post with slug "nonexistent" exists
- **WHEN:** `GET /api/v1/posts/nonexistent` is called
- **THEN:** Response status is 404 with `{"error": "not found"}`

### Requirement: GET /api/v1/tags SHALL return all unique tags with post counts

#### Scenario: Tags endpoint returns aggregated tag list

- **GIVEN:** Posts exist with tags ["python", "ml"] and ["python", "tutorial"]
- **WHEN:** `GET /api/v1/tags` is called
- **THEN:** Response is `[{"name": "python", "count": 2}, {"name": "ml", "count": 1}, {"name": "tutorial", "count": 1}]`

### Requirement: CORS middleware SHALL allow requests from the configured frontend origin

#### Scenario: Frontend origin is allowed

- **GIVEN:** `CORS_ORIGINS` is set to `http://localhost:5173`
- **WHEN:** A request with `Origin: http://localhost:5173` arrives
- **THEN:** Response includes `Access-Control-Allow-Origin: http://localhost:5173`

### Requirement: Seed data SHOULD be available for development

#### Scenario: Dev seed creates sample posts

- **GIVEN:** The database is empty and the environment is development
- **WHEN:** The server starts (or a seed command is run)
- **THEN:** 3-5 sample posts with varied tags exist in the database
