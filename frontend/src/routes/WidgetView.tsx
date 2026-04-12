import { useParams, Link } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { WidgetEmbed } from "../components/blog/WidgetEmbed";

export function WidgetView() {
  const { id } = useParams<{ id: string }>();

  if (!id) {
    return (
      <Layout>
        <p className="text-text-secondary">No widget specified.</p>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="mb-4">
        <Link to="/" className="text-accent hover:underline text-sm">
          ← Back to home
        </Link>
      </div>
      <WidgetEmbed widgetId={id} />
    </Layout>
  );
}
