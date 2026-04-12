import { useState, useEffect } from "react";
import { Link, useNavigate } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { useAuth } from "../hooks/useAuth";
import {
  adminGetPosts,
  adminGetPages,
  adminDeletePost,
  adminDeletePage,
} from "../services/api";
import type { Post, Page } from "../types/post";

export function AdminDashboard() {
  const { credentials, user, authMode, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [posts, setPosts] = useState<Post[]>([]);
  const [pages, setPages] = useState<Page[]>([]);
  const [loading, setLoading] = useState(true);

  // For OAuth mode, use null credentials (cookies handle it); for basic, use credentials.
  const authCreds = authMode === "oauth" ? null : credentials;

  useEffect(() => {
    if (!isAuthenticated) {
      navigate("/admin/login");
      return;
    }

    Promise.all([adminGetPosts(authCreds), adminGetPages(authCreds)])
      .then(([postData, pageData]) => {
        setPosts(postData.posts || []);
        setPages(pageData || []);
        setLoading(false);
      })
      .catch(() => {
        setLoading(false);
      });
  }, [isAuthenticated, authCreds, navigate]);

  const handleDeletePost = async (slug: string) => {
    if (!confirm(`Delete post "${slug}"?`)) return;
    try {
      await adminDeletePost(authCreds, slug);
      setPosts((prev) => prev.filter((p) => p.slug !== slug));
    } catch {
      alert("Failed to delete post.");
    }
  };

  const handleDeletePage = async (slug: string) => {
    if (!confirm(`Delete page "${slug}"?`)) return;
    try {
      await adminDeletePage(authCreds, slug);
      setPages((prev) => prev.filter((p) => p.slug !== slug));
    } catch {
      alert("Failed to delete page.");
    }
  };

  if (!isAuthenticated) return null;

  const displayName = authMode === "oauth" && user ? user.username : "Admin";
  const avatarUrl = authMode === "oauth" && user ? user.avatar_url : null;

  return (
    <Layout>
      <div className="flex justify-between items-center mb-8">
        <h1 className="font-heading text-2xl font-bold">Admin Dashboard</h1>
        <div className="flex items-center gap-3">
          {avatarUrl && (
            <img
              src={avatarUrl}
              alt={displayName}
              className="w-6 h-6 border border-border"
            />
          )}
          <span className="text-sm text-text-secondary font-heading">
            {displayName}
          </span>
          <button
            onClick={() => {
              logout();
              navigate("/admin/login");
            }}
            className="border-2 border-border px-3 py-1 text-sm font-heading hover:border-accent"
          >
            Logout
          </button>
        </div>
      </div>

      {loading ? (
        <p className="text-text-secondary">Loading…</p>
      ) : (
        <>
          {/* Posts Section */}
          <section className="mb-10">
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Posts ({posts.length})
              </h2>
              <Link
                to="/admin/posts/new"
                className="border-2 border-border bg-accent text-white px-3 py-1 text-sm font-heading hover:opacity-90 no-underline"
              >
                + New Post
              </Link>
            </div>

            <div className="space-y-2">
              {posts.map((post) => (
                <div
                  key={post.slug}
                  className="flex items-center justify-between border-2 border-border bg-surface p-3"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`inline-block w-2 h-2 rounded-full ${
                        post.published ? "bg-success" : "bg-text-secondary"
                      }`}
                      title={post.published ? "Published" : "Draft"}
                    />
                    <span className="font-heading text-sm">{post.title}</span>
                    {!post.published && (
                      <span className="text-xs text-text-secondary">
                        (draft)
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Link
                      to={`/admin/posts/${post.slug}/edit`}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent no-underline"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDeletePost(post.slug)}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
              {posts.length === 0 && (
                <p className="text-text-secondary text-sm">No posts yet.</p>
              )}
            </div>
          </section>

          {/* Pages Section */}
          <section>
            <div className="flex justify-between items-center mb-4">
              <h2 className="font-heading text-xl font-bold">
                Pages ({pages.length})
              </h2>
              <Link
                to="/admin/pages/new"
                className="border-2 border-border bg-accent text-white px-3 py-1 text-sm font-heading hover:opacity-90 no-underline"
              >
                + New Page
              </Link>
            </div>

            <div className="space-y-2">
              {pages.map((page) => (
                <div
                  key={page.slug}
                  className="flex items-center justify-between border-2 border-border bg-surface p-3"
                >
                  <div className="flex items-center gap-3">
                    <span
                      className={`inline-block w-2 h-2 rounded-full ${
                        page.published ? "bg-success" : "bg-text-secondary"
                      }`}
                      title={page.published ? "Published" : "Draft"}
                    />
                    <span className="font-heading text-sm">{page.title}</span>
                    <span className="text-xs text-text-secondary">
                      #{page.sort_order}
                    </span>
                    {!page.published && (
                      <span className="text-xs text-text-secondary">
                        (draft)
                      </span>
                    )}
                  </div>
                  <div className="flex gap-2">
                    <Link
                      to={`/admin/pages/${page.slug}/edit`}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent no-underline"
                    >
                      Edit
                    </Link>
                    <button
                      onClick={() => handleDeletePage(page.slug)}
                      className="border border-border px-2 py-1 text-xs font-heading hover:border-accent text-accent"
                    >
                      Delete
                    </button>
                  </div>
                </div>
              ))}
              {pages.length === 0 && (
                <p className="text-text-secondary text-sm">No pages yet.</p>
              )}
            </div>
          </section>
        </>
      )}
    </Layout>
  );
}
