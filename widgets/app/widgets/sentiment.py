from app.widgets import BaseWidget


class SentimentWidget(BaseWidget):
    """A simple sentiment analysis widget using keyword matching."""

    id = "sentiment"
    name = "Sentiment Analysis"
    description = "Analyze the sentiment of a text snippet. Enter text and get a positive/negative/neutral result."

    POSITIVE_WORDS = {
        "good",
        "great",
        "excellent",
        "amazing",
        "wonderful",
        "fantastic",
        "happy",
        "love",
        "best",
        "awesome",
        "brilliant",
        "joy",
        "perfect",
        "beautiful",
        "delightful",
        "superb",
        "outstanding",
        "nice",
        "fun",
    }
    NEGATIVE_WORDS = {
        "bad",
        "terrible",
        "awful",
        "horrible",
        "worst",
        "hate",
        "ugly",
        "sad",
        "angry",
        "poor",
        "disappointing",
        "dreadful",
        "miserable",
        "boring",
        "annoying",
        "broken",
        "fail",
        "wrong",
        "painful",
    }

    async def process(self, input_data: dict) -> dict:
        text = input_data.get("text", "").lower()
        words = text.split()

        pos = sum(1 for w in words if w.strip(".,!?;:") in self.POSITIVE_WORDS)
        neg = sum(1 for w in words if w.strip(".,!?;:") in self.NEGATIVE_WORDS)

        if pos > neg:
            sentiment = "positive"
            confidence = min(0.5 + (pos - neg) * 0.1, 0.99)
        elif neg > pos:
            sentiment = "negative"
            confidence = min(0.5 + (neg - pos) * 0.1, 0.99)
        else:
            sentiment = "neutral"
            confidence = 0.5

        emoji = {"positive": "😊", "negative": "😞", "neutral": "😐"}[sentiment]

        return {
            "sentiment": sentiment,
            "confidence": round(confidence, 2),
            "emoji": emoji,
            "positive_count": pos,
            "negative_count": neg,
        }

    def render_form(self) -> str:
        return """
        <div>
            <textarea id="input-text" rows="4" placeholder="Enter text to analyze…"
                style="width:100%;padding:8px;border:2px solid #2D2D2D;font-family:Inter,sans-serif;
                font-size:14px;resize:vertical;background:#FAFAF8;"></textarea>
            <button onclick="analyze()" style="margin-top:8px;padding:8px 16px;border:2px solid #2D2D2D;
                background:#E85D3A;color:white;font-family:'JetBrains Mono',monospace;font-size:13px;
                cursor:pointer;">Analyze</button>
            <div id="result" style="margin-top:12px;padding:12px;border:2px solid #2D2D2D;
                display:none;font-family:Inter,sans-serif;"></div>
        </div>
        <script>
        async function analyze() {
            const text = document.getElementById('input-text').value;
            if (!text.trim()) return;
            const resultDiv = document.getElementById('result');
            resultDiv.style.display = 'block';
            resultDiv.innerHTML = 'Analyzing…';
            try {
                const resp = await fetch(window.location.pathname.replace('/embed', '/process'), {
                    method: 'POST',
                    headers: {'Content-Type': 'application/json'},
                    body: JSON.stringify({text})
                });
                const data = await resp.json();
                const colors = {positive: '#2D8A4E', negative: '#E85D3A', neutral: '#6B6B6B'};
                resultDiv.innerHTML = `
                    <div style="font-size:32px;margin-bottom:4px;">${data.emoji}</div>
                    <div style="font-size:18px;font-weight:bold;color:${colors[data.sentiment]};
                        font-family:'JetBrains Mono',monospace;">${data.sentiment.toUpperCase()}</div>
                    <div style="font-size:13px;color:#6B6B6B;margin-top:4px;">
                        Confidence: ${Math.round(data.confidence * 100)}%
                    </div>`;
            } catch {
                resultDiv.innerHTML = '<span style="color:#E85D3A;">Error processing text.</span>';
            }
        }
        </script>
        """
