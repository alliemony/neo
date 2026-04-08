# Tasks: Python Widget Integration

## Widget Service

1. [ ] Define `BaseWidget` abstract class with process/render interface
2. [ ] Implement `WidgetRegistry` service
3. [ ] Write tests for registry (register, get, list)
4. [ ] Create embed HTML template (Jinja2)
5. [ ] Write tests for registry API routes (GET /widgets, GET /widgets/:id)
6. [ ] Implement registry API routes
7. [ ] Write tests for embed route (GET /widgets/:id/embed returns HTML)
8. [ ] Implement embed route with template rendering
9. [ ] Create sample sentiment analysis widget
10. [ ] Write tests for widget process endpoint
11. [ ] Implement widget process endpoint (POST /widgets/:id/process)
12. [ ] Update Dockerfile.widgets for production

## Frontend

13. [ ] Write tests for `WidgetEmbed` component (loading, loaded, error states)
14. [ ] Implement `WidgetEmbed` component with iframe
15. [ ] Update `PostContent` to handle widget post type
16. [ ] Add widget ID field to admin post editor for widget type
17. [ ] Implement `/widgets/:id` standalone view route
18. [ ] Verify end-to-end: create widget post in admin → view renders iframe → widget works
