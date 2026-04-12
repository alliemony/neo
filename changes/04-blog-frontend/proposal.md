# Proposal: Blog Frontend (Feed & Post Views)

## Why

With the backend API serving posts and the design system in place, the frontend needs to consume the API and render the blog experience. This is the user-facing heart of the site -- the Tumblr-inspired post feed and single post view.

## What

Build the frontend components and routes for the blog: the post feed (home page), single post view, tag filtering, and markdown rendering. All components use the design system tokens and follow the retro aesthetic.

## What Changes

- Create API client service (`services/api.ts`) for fetching posts and tags
- Create `PostCard` component with title, excerpt, tags, timestamp, interaction counts
- Create `PostList` component (paginated feed)
- Create `TagPill` component for tag badges
- Create `Home` route (post feed)
- Create `Post` route (single post view with markdown rendering)
- Create `Tag` route (filtered feed by tag)
- Add markdown rendering (code syntax highlighting included)
- Implement relative timestamps ("2 hours ago")
- Responsive layout for all views

## Approach

Components are stateless and receive data via props. Data fetching happens in route components or custom hooks. The API client is a thin fetch wrapper. Markdown is rendered client-side with syntax highlighting for code blocks.
