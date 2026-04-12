import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { MarkdownContent } from "../components/blog/MarkdownContent";
import { useAuth } from "../hooks/useAuth";
import {
  getPageBySlug,
  adminCreatePage,
  adminUpdatePage,
} from "../services/api";

export function PageEditor() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { credentials, authMode, isAuthenticated } = useAuth();
  const isNew = !slug || slug === "new";
  const authCreds = authMode === "oauth" ? null : credentials;

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [published, setPublished] = useState(false);
  const [sortOrder, setSortOrder] = useState(0);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/admin/login");
      return;
    }

    if (!isNew && slug) {
      getPageBySlug(slug)
        .then((page) => {
          setTitle(page.title);
          setContent(page.content);
          setPublished(page.published);
          setSortOrder(page.sort_order);
        })
        .catch(() => navigate("/admin"));
    }
  }, [slug, isNew, isAuthenticated, navigate]);

  const handleSave = async () => {
    if (!isAuthenticated || !title.trim()) return;
    setSaving(true);

    try {
      if (isNew) {
        await adminCreatePage(authCreds, {
          title,
          content,
          published,
          sort_order: sortOrder,
        });
      } else if (slug) {
        await adminUpdatePage(authCreds, slug, {
          title,
          content,
          published,
          sort_order: sortOrder,
        });
      }
      navigate("/admin");
    } catch {
      alert("Failed to save page.");
    } finally {
      setSaving(false);
    }
  };

  if (!isAuthenticated) return null;

  return (
    <Layout>
      <div className="mb-6">
        <h1 className="font-heading text-2xl font-bold">
          {isNew ? "New Page" : "Edit Page"}
        </h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-heading mb-1">Title</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary focus:outline-none focus:border-accent"
            />
          </div>

          <div>
            <label className="block text-sm font-heading mb-1">
              Content (Markdown)
            </label>
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              rows={20}
              className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-code text-text-primary focus:outline-none focus:border-accent resize-y"
            />
          </div>

          <div className="flex gap-4">
            <div>
              <label className="block text-sm font-heading mb-1">
                Sort Order
              </label>
              <input
                type="number"
                value={sortOrder}
                onChange={(e) => setSortOrder(Number(e.target.value))}
                className="w-20 border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary focus:outline-none focus:border-accent"
              />
            </div>
          </div>

          <div className="flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm font-heading cursor-pointer">
              <input
                type="checkbox"
                checked={published}
                onChange={(e) => setPublished(e.target.checked)}
                className="w-4 h-4"
              />
              Published
            </label>

            <button
              onClick={handleSave}
              disabled={saving || !title.trim()}
              className="border-2 border-border bg-accent text-white px-4 py-2 text-sm font-heading hover:opacity-90 disabled:opacity-50"
            >
              {saving ? "Saving…" : "Save"}
            </button>

            <button
              onClick={() => navigate("/admin")}
              className="border-2 border-border px-4 py-2 text-sm font-heading hover:border-accent"
            >
              Cancel
            </button>
          </div>
        </div>

        <div className="border-2 border-border bg-surface p-4">
          <h2 className="font-heading text-sm font-bold mb-3 uppercase tracking-wider text-text-secondary">
            Preview
          </h2>
          <h3 className="font-heading text-xl font-bold mb-4">
            {title || "Untitled"}
          </h3>
          <MarkdownContent content={content} />
        </div>
      </div>
    </Layout>
  );
}
