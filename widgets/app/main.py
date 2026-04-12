from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.routes import router, init_registry
from app.services import WidgetRegistry
from app.widgets.sentiment import SentimentWidget

app = FastAPI(title="Neo Widgets", version="0.1.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

# Set up widget registry
registry = WidgetRegistry()
registry.register(SentimentWidget())
init_registry(registry)

app.include_router(router)


@app.get("/health")
async def health():
    return {"status": "ok"}
