# Spec: Admin Panel

## Purpose

Provide an authenticated content management interface for creating, editing, publishing, and deleting blog posts and static pages.

## Requirements

### Requirement: An Authenticator interface SHALL abstract the auth mechanism to allow future swapping

#### Scenario: BasicAuthenticator implements the Authenticator interface

- **GIVEN:** The `Authenticator` interface defines `Middleware()` and `CurrentUser()` methods
- **WHEN:** `BasicAuthenticator` is instantiated with config
- **THEN:** It satisfies the `Authenticator` interface
- **AND:** A future `OAuthAuthenticator` can replace it without changing handlers

### Requirement: Basic auth middleware SHALL protect all admin API routes

#### Scenario: Valid credentials grant access

- **GIVEN:** `ADMIN_USERNAME=admin` and `ADMIN_PASSWORD` hash is configured
- **WHEN:** A request to `/api/v1/admin/posts` includes `Authorization: Basic <valid>` header
- **THEN:** The request proceeds to the handler

#### Scenario: Missing credentials return 401

- **GIVEN:** No Authorization header is present
- **WHEN:** A request to `/api/v1/admin/posts` is made
- **THEN:** Response status is 401 with `WWW-Authenticate: Basic` header

#### Scenario: Invalid credentials return 401

- **GIVEN:** Wrong username or password in the Authorization header
- **WHEN:** A request to `/api/v1/admin/posts` is made
- **THEN:** Response status is 401

### Requirement: Admin credentials SHALL be stored as hashed values, never plaintext

#### Scenario: Password is compared against bcrypt hash

- **GIVEN:** `ADMIN_PASSWORD` env var contains a bcrypt hash
- **WHEN:** A login attempt is made with the raw password
- **THEN:** The password is verified against the bcrypt hash

### Requirement: POST /api/v1/admin/posts SHALL create a new blog post

#### Scenario: Create a draft post

- **GIVEN:** Authenticated admin session
- **WHEN:** `POST /api/v1/admin/posts` with `{"title": "My Post", "content": "...", "tags": ["go"], "published": false}`
- **THEN:** Response status is 201 with the created post including auto-generated slug

### Requirement: PUT /api/v1/admin/posts/:slug SHALL update an existing post

#### Scenario: Update post content

- **GIVEN:** Post "hello-world" exists and admin is authenticated
- **WHEN:** `PUT /api/v1/admin/posts/hello-world` with updated content
- **THEN:** Response status is 200 with the updated post
- **AND:** `updated_at` timestamp is refreshed

### Requirement: DELETE /api/v1/admin/posts/:slug SHALL delete a post and its comments

#### Scenario: Delete cascades to comments

- **GIVEN:** Post "hello-world" has 5 comments
- **WHEN:** `DELETE /api/v1/admin/posts/hello-world` is called
- **THEN:** Response status is 204
- **AND:** The post and all 5 comments are removed

### Requirement: Admin page CRUD endpoints SHALL follow the same pattern as posts

#### Scenario: Create, update, delete pages

- **GIVEN:** Admin is authenticated
- **WHEN:** Pages are managed via `/api/v1/admin/pages` endpoints
- **THEN:** Pages support the same CRUD operations as posts (create, read, update, delete)

### Requirement: Admin login page SHALL collect credentials and establish a session

#### Scenario: Successful login

- **GIVEN:** The admin visits `/admin`
- **WHEN:** They enter valid credentials and submit
- **THEN:** They are redirected to the admin dashboard
- **AND:** Subsequent API calls include auth credentials

#### Scenario: Failed login shows error

- **GIVEN:** The admin visits `/admin`
- **WHEN:** They enter invalid credentials
- **THEN:** An error message is displayed and they remain on the login page

### Requirement: Admin dashboard SHALL list all posts and pages with management actions

#### Scenario: Dashboard shows posts with status

- **GIVEN:** Admin is logged in and 5 posts exist (3 published, 2 drafts)
- **WHEN:** The admin dashboard loads
- **THEN:** All 5 posts are listed with their published/draft status
- **AND:** Each post has Edit and Delete action buttons

### Requirement: Post editor SHALL provide a markdown textarea with live preview

#### Scenario: Live preview updates as user types

- **GIVEN:** The admin is editing a post
- **WHEN:** They type markdown in the content textarea
- **THEN:** A rendered preview updates in real-time in an adjacent panel

#### Scenario: Editor populates fields for existing post

- **GIVEN:** The admin clicks Edit on an existing post
- **WHEN:** The editor loads
- **THEN:** Title, slug, content, tags, and publish status are pre-filled

### Requirement: Tag input SHALL autocomplete from existing tags

#### Scenario: Tag suggestions appear as user types

- **GIVEN:** Tags "python", "go", "tutorial" exist in the database
- **WHEN:** The admin types "py" in the tag input
- **THEN:** "python" appears as a suggestion

### Requirement: Publish toggle SHALL switch between draft and published states

#### Scenario: Publishing a draft

- **GIVEN:** A post is in draft status
- **WHEN:** The admin toggles the publish switch and saves
- **THEN:** The post becomes publicly visible on the blog

### Requirement: Delete action SHALL require confirmation before executing

#### Scenario: Delete requires confirmation

- **GIVEN:** The admin clicks Delete on a post
- **WHEN:** A confirmation dialog appears
- **THEN:** The post is only deleted if the admin confirms
