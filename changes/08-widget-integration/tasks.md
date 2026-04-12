# Tasks: Python Widget Integration

## Widget Service

1. [x] Define `BaseWidget` abstract class with process/render interface
2. [x] Implement `WidgetRegistry` service
3. [x] Write tests for registry (register, get, list)
4. [x] Create embed HTML template (Jinja2)
5. [x] Write tests for registry API routes (GET /widgets, GET /widgets/:id)
6. [x] Implement registry API routes
7. [x] Write tests for embed route (GET /widgets/:id/embed returns HTML)
8. [x] Implement embed route with template rendering
9. [x] Create sample sentiment analysis widget
10. [x] Write tests for widget process endpoint
11. [x] Implement widget process endpoint (POST /widgets/:id/process)
12. [x] Update Dockerfile.widgets for production

## Frontend

13. [x] Write tests for `WidgetEmbed` component (loading, loaded, error states)
14. [x] Implement `WidgetEmbed` component with iframe
15. [x] Update `PostContent` to handle widget post type
16. [x] Add widget ID field to admin post editor for widget type
17. [x] Implement `/widgets/:id` standalone view route
18. [x] Verify end-to-end: create widget post in admin → view renders iframe → widget works
