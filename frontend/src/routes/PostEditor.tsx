import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { MarkdownContent } from "../components/blog/MarkdownContent";
import { useAuth } from "../hooks/useAuth";
import {
  getPostBySlug,
  adminCreatePost,
  adminUpdatePost,
  getTags,
} from "../services/api";
import type { TagCount } from "../types/post";

export function PostEditor() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { credentials, authMode, isAuthenticated } = useAuth();
  const isNew = !slug || slug === "new";
  const authCreds = authMode === "oauth" ? null : credentials;

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tagInput, setTagInput] = useState("");
  const [tags, setTags] = useState<string[]>([]);
  const [published, setPublished] = useState(false);
  const [saving, setSaving] = useState(false);
  const [allTags, setAllTags] = useState<TagCount[]>([]);
  const [suggestions, setSuggestions] = useState<string[]>([]);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/admin/login");
      return;
    }

    getTags()
      .then(setAllTags)
      .catch(() => {});

    if (!isNew && slug) {
      getPostBySlug(slug)
        .then((post) => {
          setTitle(post.title);
          setContent(post.content);
          setTags(post.tags);
          setPublished(post.published);
        })
        .catch(() => navigate("/admin"));
    }
  }, [slug, isNew, isAuthenticated, navigate]);

  useEffect(() => {
    if (tagInput.length > 0) {
      const filtered = allTags
        .map((t) => t.name)
        .filter(
          (name) =>
            name.toLowerCase().includes(tagInput.toLowerCase()) &&
            !tags.includes(name),
        );
      setSuggestions(filtered);
    } else {
      setSuggestions([]);
    }
  }, [tagInput, allTags, tags]);

  const addTag = (tag: string) => {
    const trimmed = tag.trim().toLowerCase();
    if (trimmed && !tags.includes(trimmed)) {
      setTags((prev) => [...prev, trimmed]);
    }
    setTagInput("");
    setSuggestions([]);
  };

  const removeTag = (tag: string) => {
    setTags((prev) => prev.filter((t) => t !== tag));
  };

  const handleSave = async () => {
    if (!isAuthenticated || !title.trim()) return;
    setSaving(true);

    try {
      if (isNew) {
        await adminCreatePost(authCreds, {
          title,
          content,
          content_type: "markdown",
          tags,
          published,
        });
      } else if (slug) {
        await adminUpdatePost(authCreds, slug, {
          title,
          content,
          tags,
          published,
        });
      }
      navigate("/admin");
    } catch {
      alert("Failed to save post.");
    } finally {
      setSaving(false);
    }
  };

  if (!isAuthenticated) return null;

  return (
    <Layout>
      <div className="mb-6">
        <h1 className="font-heading text-2xl font-bold">
          {isNew ? "New Post" : "Edit Post"}
        </h1>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Editor */}
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

          <div>
            <label className="block text-sm font-heading mb-1">Tags</label>
            <div className="flex flex-wrap gap-1 mb-2">
              {tags.map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center gap-1 bg-tag-bg px-2 py-0.5 text-xs font-heading"
                >
                  #{tag}
                  <button
                    onClick={() => removeTag(tag)}
                    className="text-text-secondary hover:text-accent"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
            <div className="relative">
              <input
                type="text"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") {
                    e.preventDefault();
                    addTag(tagInput);
                  }
                }}
                placeholder="Add tag…"
                className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary focus:outline-none focus:border-accent"
              />
              {suggestions.length > 0 && (
                <div className="absolute z-10 w-full border-2 border-border bg-surface mt-[-2px]">
                  {suggestions.map((s) => (
                    <button
                      key={s}
                      onClick={() => addTag(s)}
                      className="block w-full text-left px-3 py-1 text-sm hover:bg-tag-bg"
                    >
                      #{s}
                    </button>
                  ))}
                </div>
              )}
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

        {/* Preview */}
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
