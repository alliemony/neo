# Design: Python Widget Integration

## Technical Approach

### Widget Service Architecture

```
widgets/
├── app/
│   ├── main.py              # FastAPI app, mounts widget routes
│   ├── routes/
│   │   ├── registry.py      # GET /widgets, GET /widgets/:id
│   │   └── embed.py         # GET /widgets/:id/embed
│   ├── services/
│   │   └── widget_registry.py  # Widget discovery and metadata
│   ├── widgets/             # Individual widget implementations
│   │   ├── base.py          # Abstract base widget class
│   │   └── sentiment.py     # Sample: sentiment analysis
│   └── templates/
│       └── embed.html       # Jinja2 template for widget HTML page
└── tests/
```

### Widget Base Class

```python
# widgets/base.py
from abc import ABC, abstractmethod

class BaseWidget(ABC):
    id: str
    name: str
    description: str

    @abstractmethod
    async def process(self, input_data: dict) -> dict:
        """Process widget input and return result."""

    @abstractmethod
    def render_form(self) -> str:
        """Return HTML form for widget input."""
```

### Widget Registry

```python
# services/widget_registry.py
class WidgetRegistry:
    def __init__(self):
        self._widgets: dict[str, BaseWidget] = {}

    def register(self, widget: BaseWidget):
        self._widgets[widget.id] = widget

    def get(self, widget_id: str) -> BaseWidget | None:
        return self._widgets.get(widget_id)

    def list_all(self) -> list[dict]:
        return [{"id": w.id, "name": w.name, "description": w.description}
                for w in self._widgets.values()]
```

### Embed HTML Template

Self-contained HTML page served for each widget:

```html
<!-- templates/embed.html -->
<!DOCTYPE html>
<html>
<head>
  <style>
    /* Inline styles matching retro theme */
    body { font-family: 'Inter', sans-serif; margin: 1rem; }
    input, textarea { border: 2px solid #2D2D2D; padding: 0.5rem; }
    button { background: #E85D3A; color: white; border: none; padding: 0.5rem 1rem; cursor: pointer; }
  </style>
</head>
<body>
  {{ widget_form | safe }}
  <div id="result"></div>
  <script>
    // Inline JS for form submission + result display
    // Posts to /widgets/{{ widget_id }}/process
  </script>
</body>
</html>
```

### Frontend WidgetEmbed Component

```tsx
function WidgetEmbed({ widgetId }: { widgetId: string }) {
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState(false);
  const widgetUrl = `${WIDGET_SERVICE_URL}/widgets/${widgetId}/embed`;

  return (
    <div className="border-2 border-border">
      {!loaded && !error && <div className="p-4 text-text-secondary">Loading widget...</div>}
      {error && <div className="p-4 text-text-secondary">Widget unavailable</div>}
      <iframe
        src={widgetUrl}
        className={`w-full ${loaded ? '' : 'hidden'}`}
        style={{ minHeight: '300px', border: 'none' }}
        onLoad={() => setLoaded(true)}
        onError={() => setError(true)}
        sandbox="allow-scripts allow-forms"
        title={`Widget: ${widgetId}`}
      />
    </div>
  );
}
```

### Post Type Integration

Posts with `content_type: "widget"` store the widget ID in the content field. The post view checks the type and renders accordingly:

```tsx
function PostContent({ post }: { post: Post }) {
  if (post.contentType === 'widget') {
    return <WidgetEmbed widgetId={post.content} />;
  }
  return <MarkdownContent content={post.content} />;
}
```
