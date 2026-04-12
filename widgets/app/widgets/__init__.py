from abc import ABC, abstractmethod


class BaseWidget(ABC):
    """Abstract base class for all widgets."""

    id: str
    name: str
    description: str

    @abstractmethod
    async def process(self, input_data: dict) -> dict:
        """Process widget input and return result."""

    @abstractmethod
    def render_form(self) -> str:
        """Return HTML form for widget input."""

    def metadata(self) -> dict:
        return {
            "id": self.id,
            "name": self.name,
            "description": self.description,
        }
