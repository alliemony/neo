import { useParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { MarkdownContent } from "../components/blog/MarkdownContent";
import { WidgetEmbed } from "../components/blog/WidgetEmbed";
import { TagPill } from "../components/blog/TagPill";
import { LikeButton } from "../components/blog/LikeButton";
import { CommentSection } from "../components/blog/CommentSection";
import { SEO } from "../components/SEO";
import { getPostBySlug } from "../services/api";
import type { Post } from "../types/post";

function formatDate(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" });
}

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
      .then((data) => { setPost(data); setLoading(false); })
      .catch((err) => {
        if (err.message.includes("404")) setNotFound(true);
        setLoading(false);
      });
  }, [slug]);

  if (loading) {
    return <Layout><p className="text-text-secondary">Loading…</p></Layout>;
  }

  if (notFound || !post) {
    return (
      <Layout>
        <div className="text-center py-20">
          <h1 className="text-6xl font-bold mb-4">404</h1>
          <p className="text-text-secondary mb-4">Post not found.</p>
          <Link to="/blog" className="text-accent hover:text-accent-alt">← Back to blog</Link>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <SEO
        title={post.title}
        description={post.content.slice(0, 160)}
        path={`/blog/${post.slug}`}
      />
      <article>
        <div className="mb-6">
          <Link to="/blog" className="text-[15px] text-accent hover:text-accent-alt no-underline">← all posts</Link>
        </div>
        <header className="mb-8">
          <time className="text-[14px] text-text-secondary block mb-2">{formatDate(post.created_at)}</time>
          <h1 className="text-[26px] font-bold text-text-primary leading-snug mb-3">{post.title}</h1>
          {post.tags.length > 0 && (
            <div className="flex flex-wrap gap-1.5">
              {post.tags.map((tag) => (
                <TagPill key={tag} tag={tag} />
              ))}
            </div>
          )}
        </header>
        <div className="text-[17px] leading-[1.75] text-text-secondary">
          {post.content_type === "widget" ? (
            <WidgetEmbed widgetId={post.content} />
          ) : (
            <MarkdownContent content={post.content} />
          )}
        </div>
        <div className="mt-8 flex items-center gap-4">
          <LikeButton slug={post.slug} initialCount={post.like_count} />
        </div>
        <div className="mt-10">
          <CommentSection slug={post.slug} />
        </div>
      </article>
    </Layout>
  );
}
