import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Layout } from '../../components/layout/Layout';
import { MarkdownContent } from '../../components/blog/MarkdownContent';
import { TagInput } from '../../components/admin/TagInput';
import { useAuth } from '../../contexts/AuthContext';
import { createAdminApi } from '../../services/adminApi';
import { getTags } from '../../services/api';

export function PostEditor() {
  const { slug } = useParams<{ slug: string }>();
  const navigate = useNavigate();
  const { authHeader } = useAuth();
  const isNew = !slug;

  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [tags, setTags] = useState<string[]>([]);
  const [published, setPublished] = useState(false);
  const [existingTags, setExistingTags] = useState<string[]>([]);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    getTags()
      .then((tagCounts) => setExistingTags(tagCounts.map((t) => t.name)))
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (!slug) return;
    const header = authHeader();
    if (!header) return;
    const api = createAdminApi(header);

    api.listPosts().then((result) => {
      const post = result.posts.find((p) => p.slug === slug);
      if (post) {
        setTitle(post.title);
        setContent(post.content);
        setTags(post.tags);
        setPublished(post.published);
      }
    }).catch(() => {});
  }, [slug, authHeader]);

  async function handleSave(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setSaving(true);

    const header = authHeader();
    if (!header) return;
    const api = createAdminApi(header);

    try {
      if (isNew) {
        await api.createPost({ title, content, tags, published });
      } else {
        await api.updatePost(slug, { title, content, tags, published });
      }
      navigate('/admin');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Save failed');
    } finally {
      setSaving(false);
    }
  }

  return (
    <Layout>
      <h1 className="font-heading text-2xl font-bold mb-6">
        {isNew ? 'New Post' : 'Edit Post'}
      </h1>
      <form onSubmit={handleSave}>
        <div className="lg:grid lg:grid-cols-2 lg:gap-6">
          <div className="space-y-4">
            <div>
              <label htmlFor="title" className="block text-sm font-medium mb-1">
                Title
              </label>
              <input
                id="title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full border-[length:var(--border-width)] border-border bg-surface p-2 text-text"
                required
              />
            </div>
            <div>
              <label htmlFor="content" className="block text-sm font-medium mb-1">
                Content
              </label>
              <textarea
                id="content"
                value={content}
                onChange={(e) => setContent(e.target.value)}
                rows={16}
                className="w-full border-[length:var(--border-width)] border-border bg-surface p-2 text-text font-mono text-sm"
              />
            </div>
            <div>
              <label className="block text-sm font-medium mb-1">Tags</label>
              <TagInput value={tags} onChange={setTags} existingTags={existingTags} />
            </div>
            <div className="flex items-center gap-2">
              <input
                id="published"
                type="checkbox"
                checked={published}
                onChange={(e) => setPublished(e.target.checked)}
              />
              <label htmlFor="published" className="text-sm">
                Published
              </label>
            </div>
            {error && <p className="text-red-500 text-sm">{error}</p>}
            <div className="flex gap-2">
              <button
                type="submit"
                disabled={saving}
                className="bg-accent text-white px-4 py-2 font-bold hover:opacity-90 disabled:opacity-50"
              >
                {saving ? 'Saving…' : 'Save'}
              </button>
              <button
                type="button"
                onClick={() => navigate('/admin')}
                className="border-[length:var(--border-width)] border-border px-4 py-2 text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
          <div className="mt-6 lg:mt-0">
            <p className="text-sm font-medium mb-1">Preview</p>
            <div className="border-[length:var(--border-width)] border-border bg-surface p-4 min-h-[300px]">
              {content ? (
                <MarkdownContent content={content} />
              ) : (
                <p className="text-text-secondary text-sm">Start typing to see preview…</p>
              )}
            </div>
          </div>
        </div>
      </form>
    </Layout>
  );
}
