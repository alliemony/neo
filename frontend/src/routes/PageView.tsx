import { useParams, Link } from "react-router-dom";
import { useState, useEffect } from "react";
import { Layout } from "../components/layout/Layout";
import { MarkdownContent } from "../components/blog/MarkdownContent";
import { SEO } from "../components/SEO";
import { getPageBySlug } from "../services/api";
import type { Page } from "../types/post";

export function PageView() {
  const { slug } = useParams<{ slug: string }>();
  const [page, setPage] = useState<Page | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    if (!slug) return;

    setLoading(true);
    setNotFound(false);

    getPageBySlug(slug)
      .then((data) => {
        setPage(data);
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

  if (notFound || !page) {
    return (
      <Layout>
        <div className="text-center py-20">
          <h1 className="font-heading text-6xl mb-4">404</h1>
          <p className="text-text-secondary mb-4">Page not found.</p>
          <Link to="/" className="text-accent hover:underline">
            ← Back to home
          </Link>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      <SEO
        title={page.title}
        description={page.content.slice(0, 160)}
        path={`/page/${page.slug}`}
      />
      <article>
        <header className="mb-8">
          <h1 className="font-heading text-3xl font-bold">{page.title}</h1>
        </header>
        <MarkdownContent content={page.content} />
      </article>
    </Layout>
  );
}
