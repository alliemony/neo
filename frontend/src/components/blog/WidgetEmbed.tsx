import { useState } from "react";

interface WidgetEmbedProps {
  widgetId: string;
  widgetServiceUrl?: string;
}

export function WidgetEmbed({
  widgetId,
  widgetServiceUrl = import.meta.env.VITE_WIDGET_URL || "http://localhost:8000",
}: WidgetEmbedProps) {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const embedUrl = `${widgetServiceUrl}/widgets/${widgetId}/embed`;

  return (
    <div className="border-2 border-border bg-surface">
      {loading && !error && (
        <div className="p-4 text-text-secondary text-sm font-heading">
          Loading widget…
        </div>
      )}
      {error && (
        <div className="p-4 text-accent text-sm font-heading">
          Widget unavailable. The widget service may not be running.
        </div>
      )}
      <iframe
        src={embedUrl}
        title={`Widget: ${widgetId}`}
        className={`w-full border-0 ${loading && !error ? "h-0" : "min-h-[300px]"}`}
        onLoad={() => setLoading(false)}
        onError={() => {
          setLoading(false);
          setError(true);
        }}
        sandbox="allow-scripts allow-same-origin"
        style={{ display: error ? "none" : "block" }}
      />
    </div>
  );
}
