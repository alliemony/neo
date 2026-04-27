import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { Layout } from '../../components/layout/Layout';
import { useAuth } from '../../contexts/AuthContext';
import { createAdminApi } from '../../services/adminApi';
import type { Post, Page } from '../../types/post';

export function AdminDashboard() {
  const { authHeader, logout } = useAuth();
  const [posts, setPosts] = useState<Post[]>([]);
  const [pages, setPages] = useState<Page[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteTarget, setDeleteTarget] = useState<{ type: 'post' | 'page'; slug: string } | null>(null);

  const loadData = useCallback(async () => {
    const header = authHeader();
    if (!header) return;
    const api = createAdminApi(header);

    try {
      const [postResult, pageResult] = await Promise.all([api.listPosts(), api.listPages()]);
      setPosts(postResult.posts);
      setPages(pageResult);
    } catch {
      // If unauthorized, log out
    }
    setLoading(false);
  }, [authHeader]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  async function handleDelete() {
    if (!deleteTarget) return;
    const header = authHeader();
    if (!header) return;
    const api = createAdminApi(header);

    try {
      if (deleteTarget.type === 'post') {
        await api.deletePost(deleteTarget.slug);
      } else {
        await api.deletePage(deleteTarget.slug);
      }
      setDeleteTarget(null);
      loadData();
    } catch {
      // handle error
    }
  }

  if (loading) {
    return (
      <Layout>
        <p className="text-text-secondary">Loading…</p>
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="space-y-8">
        <div className="flex items-center justify-between">
          <h1 className="font-heading text-2xl font-bold">Admin Dashboard</h1>
          <button onClick={logout} className="text-sm text-text-secondary hover:text-accent">
            Log out
          </button>
        </div>

        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-heading text-xl font-bold">Posts</h2>
            <Link to="/admin/posts/new" className="bg-accent text-white px-3 py-1 text-sm font-bold">
              New Post
            </Link>
          </div>
          {posts.length === 0 ? (
            <p className="text-text-secondary">No posts yet.</p>
          ) : (
            <div className="space-y-2">
              {posts.map((post) => (
                <div
                  key={post.id}
                  className="border-[length:var(--border-width)] border-border bg-surface p-3 flex items-center justify-between"
                >
                  <div>
                    <span className="font-bold">{post.title}</span>
                    <span
                      className={`ml-2 text-xs px-2 py-0.5 ${
                        post.published ? 'bg-green-800 text-green-200' : 'bg-yellow-800 text-yellow-200'
                      }`}
                    >
                      {post.published ? 'Published' : 'Draft'}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Link to={`/admin/posts/${post.slug}`} className="text-accent text-sm hover:underline">
                      Edit
                    </Link>
                    <button
                      onClick={() => setDeleteTarget({ type: 'post', slug: post.slug })}
                      className="text-red-400 text-sm hover:underline"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <section>
          <div className="flex items-center justify-between mb-4">
            <h2 className="font-heading text-xl font-bold">Pages</h2>
            <Link to="/admin/pages/new" className="bg-accent text-white px-3 py-1 text-sm font-bold">
              New Page
            </Link>
          </div>
          {pages.length === 0 ? (
            <p className="text-text-secondary">No pages yet.</p>
          ) : (
            <div className="space-y-2">
              {pages.map((page) => (
                <div
                  key={page.id}
                  className="border-[length:var(--border-width)] border-border bg-surface p-3 flex items-center justify-between"
                >
                  <div>
                    <span className="font-bold">{page.title}</span>
                    <span
                      className={`ml-2 text-xs px-2 py-0.5 ${
                        page.published ? 'bg-green-800 text-green-200' : 'bg-yellow-800 text-yellow-200'
                      }`}
                    >
                      {page.published ? 'Published' : 'Draft'}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Link to={`/admin/pages/${page.slug}`} className="text-accent text-sm hover:underline">
                      Edit
                    </Link>
                    <button
                      onClick={() => setDeleteTarget({ type: 'page', slug: page.slug })}
                      className="text-red-400 text-sm hover:underline"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        {deleteTarget && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
            <div className="bg-surface border-[length:var(--border-width)] border-border p-6 max-w-sm">
              <p className="mb-4">
                Are you sure you want to delete this {deleteTarget.type}?
              </p>
              <div className="flex gap-2 justify-end">
                <button
                  onClick={() => setDeleteTarget(null)}
                  className="px-3 py-1 text-sm border-[length:var(--border-width)] border-border"
                >
                  Cancel
                </button>
                <button
                  onClick={handleDelete}
                  className="px-3 py-1 text-sm bg-red-600 text-white font-bold"
                >
                  Confirm Delete
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
}
