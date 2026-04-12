import pytest
from httpx import ASGITransport, AsyncClient

from app.main import app


@pytest.mark.asyncio
async def test_list_widgets():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/widgets")
    assert response.status_code == 200
    data = response.json()
    assert len(data) >= 1
    assert any(w["id"] == "sentiment" for w in data)


@pytest.mark.asyncio
async def test_get_widget_by_id():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/widgets/sentiment")
    assert response.status_code == 200
    data = response.json()
    assert data["id"] == "sentiment"
    assert data["name"] == "Sentiment Analysis"


@pytest.mark.asyncio
async def test_get_widget_not_found():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/widgets/nonexistent")
    assert response.status_code == 404


@pytest.mark.asyncio
async def test_embed_returns_html():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.get("/widgets/sentiment/embed")
    assert response.status_code == 200
    assert "text/html" in response.headers["content-type"]
    assert "Sentiment Analysis" in response.text


@pytest.mark.asyncio
async def test_process_sentiment_positive():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.post(
            "/widgets/sentiment/process",
            json={"text": "This is great and amazing work!"},
        )
    assert response.status_code == 200
    data = response.json()
    assert data["sentiment"] == "positive"
    assert data["confidence"] > 0.5


@pytest.mark.asyncio
async def test_process_sentiment_negative():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.post(
            "/widgets/sentiment/process",
            json={"text": "This is terrible and horrible."},
        )
    assert response.status_code == 200
    data = response.json()
    assert data["sentiment"] == "negative"


@pytest.mark.asyncio
async def test_process_sentiment_neutral():
    transport = ASGITransport(app=app)
    async with AsyncClient(transport=transport, base_url="http://test") as client:
        response = await client.post(
            "/widgets/sentiment/process",
            json={"text": "The sky is blue."},
        )
    assert response.status_code == 200
    data = response.json()
    assert data["sentiment"] == "neutral"
