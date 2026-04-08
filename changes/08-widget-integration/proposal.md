# Proposal: Python Widget Integration

## Why

A key differentiator of this personal site is the ability to embed interactive Python-based widgets (HuggingFace models, data visualizations, notebook-style demos) directly in blog posts. This requires a separate Python service and a frontend embed mechanism.

## What

Build the Python widget service (FastAPI), a widget registry, an iframe-based embed component for the frontend, and admin support for the widget post type.

## What Changes

### Widget Service (Python)
- Expand the FastAPI service beyond the health check
- Create a widget registry (list, get by ID)
- Build a sample widget (e.g., text sentiment analysis using a HuggingFace model)
- Serve widget UIs as self-contained HTML pages (embeddable via iframe)

### Frontend
- Create `WidgetEmbed` component (responsive iframe with loading state)
- Add widget post type support in the post editor
- Add `/widgets/:id` route for standalone widget viewing

### Backend (Go)
- No changes needed -- widgets are served by the Python service and embedded via iframe

## Approach

Widgets are self-contained Python endpoints that serve both an API and a simple HTML UI. The frontend embeds them via iframes, which provides security isolation and independence. The widget service can be deployed separately and scaled independently.
