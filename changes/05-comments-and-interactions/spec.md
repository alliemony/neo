# Spec: Comments & Interactions

## Purpose

Enable readers to engage with blog posts through comments and likes, with spam prevention via rate limiting.

## Requirements

### Requirement: The database SHALL have a comments table linked to posts

#### Scenario: Comments table is created via migration

- **GIVEN:** The backend runs migrations
- **WHEN:** The comments migration executes
- **THEN:** A `comments` table exists with columns: id, post_id (FK to posts), author_name, content, created_at
- **AND:** Deleting a post cascades to delete its comments

### Requirement: The posts table SHALL include a like_count column

#### Scenario: Like count is tracked per post

- **GIVEN:** A post exists with like_count 0
- **WHEN:** The like endpoint is called
- **THEN:** The post's like_count increments to 1

### Requirement: CommentRepository SHALL provide Create and ListByPost operations

#### Scenario: Create a comment

- **GIVEN:** A post with id 1 exists
- **WHEN:** `repo.Create(comment{PostID: 1, AuthorName: "alice", Content: "Great post!"})` is called
- **THEN:** A new comment row is inserted with a timestamp

#### Scenario: List comments for a post

- **GIVEN:** Post "hello-world" has 3 comments
- **WHEN:** `repo.ListByPostSlug("hello-world")` is called
- **THEN:** 3 comments are returned ordered by created_at ascending (oldest first)

### Requirement: CommentService SHALL validate comment inputs

#### Scenario: Empty author name is rejected

- **GIVEN:** A comment with empty author_name
- **WHEN:** `service.Create(input)` is called
- **THEN:** An `ErrAuthorRequired` error is returned

#### Scenario: Empty content is rejected

- **GIVEN:** A comment with empty content
- **WHEN:** `service.Create(input)` is called
- **THEN:** An `ErrContentRequired` error is returned

#### Scenario: Content exceeding max length is rejected

- **GIVEN:** A comment with content longer than 2000 characters
- **WHEN:** `service.Create(input)` is called
- **THEN:** An `ErrContentTooLong` error is returned

### Requirement: GET /api/v1/posts/:slug/comments SHALL return all comments for a post

#### Scenario: Comments are returned in chronological order

- **GIVEN:** Post "hello-world" has comments created at t1, t2, t3
- **WHEN:** `GET /api/v1/posts/hello-world/comments` is called
- **THEN:** Response status is 200 with comments ordered oldest-first

### Requirement: POST /api/v1/posts/:slug/comments SHALL create a new comment

#### Scenario: Valid comment is created

- **GIVEN:** Post "hello-world" exists
- **WHEN:** `POST /api/v1/posts/hello-world/comments` with `{"author_name": "alice", "content": "Nice!"}` is called
- **THEN:** Response status is 201 with the created comment

#### Scenario: Comment on nonexistent post returns 404

- **GIVEN:** No post with slug "nonexistent" exists
- **WHEN:** `POST /api/v1/posts/nonexistent/comments` is called
- **THEN:** Response status is 404

### Requirement: POST /api/v1/posts/:slug/like SHALL increment the post's like count

#### Scenario: Like increments counter

- **GIVEN:** Post "hello-world" has like_count 5
- **WHEN:** `POST /api/v1/posts/hello-world/like` is called
- **THEN:** Response status is 200 with `{"like_count": 6}`

### Requirement: Rate limiting middleware SHALL restrict comment creation frequency per IP

#### Scenario: Normal usage is allowed

- **GIVEN:** A client has not submitted a comment recently
- **WHEN:** They submit a comment
- **THEN:** The comment is created successfully

#### Scenario: Rapid submissions are blocked

- **GIVEN:** A client submitted a comment within the last 30 seconds (configurable)
- **WHEN:** They submit another comment
- **THEN:** Response status is 429 with `{"error": "too many requests"}`

### Requirement: CommentSection component SHALL display in the sidebar on desktop and below content on mobile

#### Scenario: Desktop renders comments in sidebar

- **GIVEN:** Viewport width >= 1024px and post has comments
- **WHEN:** Single post view renders
- **THEN:** Comments appear in the right sidebar column

#### Scenario: Mobile renders comments below content

- **GIVEN:** Viewport width < 768px
- **WHEN:** Single post view renders
- **THEN:** Comments appear below the post content

### Requirement: Comment form SHALL accept a display name and comment text

#### Scenario: User submits a comment

- **GIVEN:** The comment form is displayed with name and content fields
- **WHEN:** The user fills in both fields and clicks submit
- **THEN:** The comment is sent to the API and appears in the list without page reload

#### Scenario: Validation errors are shown inline

- **GIVEN:** The user submits the form with an empty name
- **WHEN:** The API returns a validation error
- **THEN:** An error message is shown near the form

### Requirement: Like button SHALL update optimistically

#### Scenario: Like button shows immediate feedback

- **GIVEN:** A post has 12 likes
- **WHEN:** The user clicks the heart/like button
- **THEN:** The count immediately shows 13 (before API response)
- **AND:** If the API call fails, the count reverts to 12
