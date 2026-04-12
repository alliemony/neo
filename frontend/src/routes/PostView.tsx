import { useParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { Sidebar } from "../components/layout/Sidebar";
import { MarkdownContent } from "../components/blog/MarkdownContent";
import { WidgetEmbed } from "../components/blog/WidgetEmbed";
import { TagPill } from "../components/blog/TagPill";
import { LikeButton } from "../components/blog/LikeButton";
import { CommentSection } from "../components/blog/CommentSection";
import { SEO } from "../components/SEO";
import { formatRelativeTime } from "../utils/time";
import { getPostBySlug } from "../services/api";
import type { Post } from "../types/post";

export function PostView() {
  const { slug } = useParams<{ slug: string }>();
  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    if (!slug) return;

    setLoading(true);
    setNotFound(false);

    getPostBySlug(slug)
      .then((data) => {
        setPost(data);
        setLoading(false);
      })
      .catch((err) => {
        if (err.message.includes("404")) {
          setNotFound(true);
        }
        setLoading(false);
      });
  }, [slug]);

  if (loading) {
    return (
      <Layout>
        <p className="text-text-secondary">Loading…</p>
      </Layout>
    );
  }

  if (notFound || !post) {
    return (
      <Layout>
        <div className="text-center py-20">
          <h1 className="font-heading text-6xl mb-4">404</h1>
          <p className="text-text-secondary mb-4">Post not found.</p>
          <Link to="/" className="text-accent hover:underline">
            ← Back to home
          </Link>
        </div>
      </Layout>
    );
  }

  const sidebarContent = (
    <Sidebar>
      <CommentSection slug={post.slug} />
    </Sidebar>
  );

  return (
    <Layout sidebar={sidebarContent}>
      <SEO
        title={post.title}
        description={post.content.slice(0, 160)}
        path={`/blog/${post.slug}`}
      />
      <article>
        <header className="mb-8">
          <h1 className="font-heading text-3xl font-bold mb-2">{post.title}</h1>
          <div className="flex flex-wrap items-center gap-2 mb-4">
            {post.tags.map((tag) => (
              <TagPill key={tag} tag={tag} />
            ))}
            <time
              dateTime={post.created_at}
              title={new Date(post.created_at).toISOString()}
              className="text-xs text-text-secondary ml-auto"
            >
              {formatRelativeTime(post.created_at)}
            </time>
          </div>
        </header>
        {post.content_type === "widget" ? (
          <WidgetEmbed widgetId={post.content} />
        ) : (
          <MarkdownContent content={post.content} />
        )}
        <div className="mt-8 flex items-center gap-4">
          <LikeButton slug={post.slug} initialCount={post.like_count} />
        </div>

        {/* Mobile: comments below content */}
        <div className="mt-8 lg:hidden">
          <CommentSection slug={post.slug} />
        </div>
      </article>
    </Layout>
  );
}
