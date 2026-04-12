from app.widgets import BaseWidget


class WidgetRegistry:
    """Registry of available widgets."""

    def __init__(self):
        self._widgets: dict[str, BaseWidget] = {}

    def register(self, widget: BaseWidget) -> None:
        self._widgets[widget.id] = widget

    def get(self, widget_id: str) -> BaseWidget | None:
        return self._widgets.get(widget_id)

    def list_all(self) -> list[dict]:
        return [w.metadata() for w in self._widgets.values()]
