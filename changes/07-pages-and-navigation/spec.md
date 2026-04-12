# Spec: Static Pages & Navigation

## Purpose

Support static content pages, dynamic site navigation, a tag cloud for content discovery, and polished 404 handling.

## Requirements

### Requirement: GET /api/v1/pages SHALL return published pages ordered by sort_order

#### Scenario: Pages list returns published pages

- **GIVEN:** 3 published pages (About sort=1, Projects sort=2, Contact sort=3) and 1 draft page
- **WHEN:** `GET /api/v1/pages` is called
- **THEN:** Response is 200 with 3 pages ordered by sort_order ascending
- **AND:** The draft page is excluded

### Requirement: GET /api/v1/pages/:slug SHALL return a single published page

#### Scenario: Valid page slug returns content

- **GIVEN:** A published page with slug "about" exists
- **WHEN:** `GET /api/v1/pages/about` is called
- **THEN:** Response is 200 with the page title and markdown content

#### Scenario: Draft page returns 404 to public

- **GIVEN:** A page with slug "secret" exists but is unpublished
- **WHEN:** `GET /api/v1/pages/secret` is called without admin auth
- **THEN:** Response is 404

### Requirement: Navigation header SHALL dynamically render links from published pages

#### Scenario: Navigation reflects current pages

- **GIVEN:** Published pages "About" and "Projects" exist
- **WHEN:** The header component renders
- **THEN:** Links to `/page/about` and `/page/projects` appear in the navigation
- **AND:** A "Blog" link to `/` is always present

#### Scenario: Navigation updates when pages change

- **GIVEN:** Admin publishes a new page "Contact"
- **WHEN:** The site is reloaded
- **THEN:** A "Contact" link appears in the navigation

### Requirement: Page view route (/page/:slug) SHALL render page content with markdown

#### Scenario: Page renders with same typography as posts

- **GIVEN:** The "about" page contains markdown content
- **WHEN:** The user visits `/page/about`
- **THEN:** The page title and rendered markdown are displayed
- **AND:** The retro theme styles are applied

### Requirement: Tag cloud sidebar widget SHALL display all tags with relative sizing

#### Scenario: Tags are displayed with visual weight

- **GIVEN:** Tags exist with varying post counts (python: 10, go: 5, tutorial: 2)
- **WHEN:** The tag cloud renders in the sidebar
- **THEN:** Tags are displayed as clickable links
- **AND:** More popular tags appear visually larger or bolder

### Requirement: 404 page SHALL display a styled not-found message

#### Scenario: Unknown route shows 404

- **GIVEN:** The user navigates to `/nonexistent-path`
- **WHEN:** React Router matches no route
- **THEN:** A styled 404 page is shown with a link back to home

### Requirement: Footer SHALL display site information

#### Scenario: Footer renders consistently

- **GIVEN:** Any page on the site
- **WHEN:** The footer renders
- **THEN:** It displays a site attribution line and year
