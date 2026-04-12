import { useSearchParams } from "react-router-dom";
import { Layout } from "../components/layout/Layout";
import { Sidebar } from "../components/layout/Sidebar";
import { PostList } from "../components/blog/PostList";
import { TagCloud } from "../components/blog/TagCloud";
import { SEO } from "../components/SEO";
import { usePosts } from "../hooks/usePosts";

export function Home() {
  const [searchParams] = useSearchParams();
  const page = Number(searchParams.get("page")) || 1;

  const { posts, total, loading, error } = usePosts({ page });

  const sidebar = (
    <Sidebar>
      <TagCloud />
    </Sidebar>
  );

  return (
    <Layout sidebar={sidebar}>
      <SEO path="/" />
      <h1 className="font-heading text-3xl font-bold mb-2">neo</h1>
      <p className="text-text-secondary mb-8">personal web garden</p>
      {loading && <p className="text-text-secondary">Loading posts…</p>}
      {error && <p className="text-accent">Failed to load posts.</p>}
      {!loading && !error && (
        <PostList posts={posts} total={total} page={page} />
      )}
    </Layout>
  );
}
