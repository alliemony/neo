from fastapi import APIRouter, HTTPException
from fastapi.responses import HTMLResponse

from app.services import WidgetRegistry

router = APIRouter()
registry = WidgetRegistry()


def init_registry(reg: WidgetRegistry) -> None:
    """Set the registry instance used by routes."""
    global registry
    registry = reg


@router.get("/widgets")
async def list_widgets():
    return registry.list_all()


@router.get("/widgets/{widget_id}")
async def get_widget(widget_id: str):
    widget = registry.get(widget_id)
    if not widget:
        raise HTTPException(status_code=404, detail="Widget not found")
    return widget.metadata()


@router.get("/widgets/{widget_id}/embed", response_class=HTMLResponse)
async def embed_widget(widget_id: str):
    widget = registry.get(widget_id)
    if not widget:
        raise HTTPException(status_code=404, detail="Widget not found")

    form_html = widget.render_form()
    html = f"""<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{widget.name}</title>
    <style>
        * {{ margin: 0; padding: 0; box-sizing: border-box; }}
        body {{
            font-family: Inter, system-ui, sans-serif;
            background: #FAFAF8;
            color: #1A1A1A;
            padding: 16px;
            line-height: 1.6;
        }}
        h2 {{
            font-family: 'JetBrains Mono', monospace;
            font-size: 16px;
            margin-bottom: 4px;
        }}
        .desc {{
            font-size: 13px;
            color: #6B6B6B;
            margin-bottom: 12px;
        }}
    </style>
</head>
<body>
    <h2>{widget.name}</h2>
    <p class="desc">{widget.description}</p>
    {form_html}
</body>
</html>"""
    return HTMLResponse(content=html)


@router.post("/widgets/{widget_id}/process")
async def process_widget(widget_id: str, input_data: dict):
    widget = registry.get(widget_id)
    if not widget:
        raise HTTPException(status_code=404, detail="Widget not found")
    result = await widget.process(input_data)
    return result
