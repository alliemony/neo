# Spec: Blog Frontend (Feed & Post Views)

## Purpose

Render the blog post feed, single post view, and tag filtering in the browser using the retro design system, with data fetched from the backend API.

## Requirements

### Requirement: An API client service SHALL provide typed functions for fetching posts and tags

#### Scenario: Fetch all posts

- **GIVEN:** The backend API is running with published posts
- **WHEN:** `api.getPosts()` is called
- **THEN:** A typed `{ posts: Post[], total: number }` response is returned

#### Scenario: Fetch posts by tag

- **GIVEN:** The backend has posts tagged "python"
- **WHEN:** `api.getPosts({ tag: "python" })` is called
- **THEN:** Only posts with the "python" tag are returned

#### Scenario: API error is handled gracefully

- **GIVEN:** The backend API is unreachable
- **WHEN:** `api.getPosts()` is called
- **THEN:** An error is thrown that callers can catch and display

### Requirement: PostCard SHALL display title, excerpt, tags, timestamp, and interaction counts

#### Scenario: PostCard renders all required elements

- **GIVEN:** A post object with title, content, tags, createdAt, likeCount, commentCount
- **WHEN:** `<PostCard post={post} />` is rendered
- **THEN:** The title is visible as a link to the full post
- **AND:** A content excerpt (first ~200 chars) is shown
- **AND:** Tags are rendered as clickable pill badges
- **AND:** A relative timestamp is displayed
- **AND:** Like count and comment count are shown

#### Scenario: PostCard has retro styling

- **GIVEN:** The retro theme is active
- **WHEN:** A PostCard renders
- **THEN:** It has a 2px solid border, no border-radius, no box-shadow
- **AND:** The title uses the monospace heading font

### Requirement: TagPill SHALL be a clickable badge that links to the tag filter page

#### Scenario: Clicking a tag navigates to filtered view

- **GIVEN:** A TagPill with tag name "python" is rendered
- **WHEN:** The user clicks the tag
- **THEN:** The browser navigates to `/tag/python`

### Requirement: The home route (/) SHALL display a paginated feed of published posts

#### Scenario: Home page shows post feed

- **GIVEN:** The backend has 15 published posts
- **WHEN:** The user visits `/`
- **THEN:** The first page of posts is displayed as PostCard components
- **AND:** Pagination controls are visible if more posts exist

#### Scenario: Empty state is handled

- **GIVEN:** No published posts exist
- **WHEN:** The user visits `/`
- **THEN:** A friendly "No posts yet" message is displayed

### Requirement: The post route (/blog/:slug) SHALL render a single post with full markdown content

#### Scenario: Single post view renders full content

- **GIVEN:** A post with slug "hello-world" exists with markdown content
- **WHEN:** The user visits `/blog/hello-world`
- **THEN:** The full post title, rendered markdown content, tags, and timestamp are displayed

#### Scenario: Code blocks have syntax highlighting

- **GIVEN:** A post contains a fenced code block with language annotation
- **WHEN:** The post is rendered
- **THEN:** The code block has syntax highlighting with the theme's code background color

#### Scenario: Non-existent post shows 404

- **GIVEN:** No post with slug "nonexistent" exists
- **WHEN:** The user visits `/blog/nonexistent`
- **THEN:** A 404 page is displayed

### Requirement: The tag route (/tag/:tag) SHALL display posts filtered by that tag

#### Scenario: Tag page shows filtered feed

- **GIVEN:** 3 posts are tagged "python"
- **WHEN:** The user visits `/tag/python`
- **THEN:** Only the 3 "python" posts are displayed
- **AND:** The active tag is visually indicated

### Requirement: Timestamps SHALL display as relative for recent posts and absolute for older ones

#### Scenario: Recent post shows relative time

- **GIVEN:** A post was created 2 hours ago
- **WHEN:** The PostCard renders
- **THEN:** The timestamp shows "2 hours ago"

#### Scenario: Old post shows absolute date

- **GIVEN:** A post was created 30 days ago
- **WHEN:** The PostCard renders
- **THEN:** The timestamp shows a formatted date like "Mar 9, 2026"

#### Scenario: Full timestamp is available on hover

- **GIVEN:** Any post is rendered
- **WHEN:** The user hovers over the timestamp
- **THEN:** The full ISO datetime is shown via the `<time>` element's title attribute

### Requirement: All views SHALL be responsive across desktop, tablet, and mobile breakpoints

#### Scenario: Post feed adapts to mobile

- **GIVEN:** The viewport is < 768px
- **WHEN:** The home page renders
- **THEN:** PostCards are full-width and the sidebar is below the main content
