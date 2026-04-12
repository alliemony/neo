import { useSearchParams } from 'react-router-dom';
import { Layout } from '../components/layout/Layout';
import { PostList } from '../components/blog/PostList';
import { usePosts } from '../hooks/usePosts';

export function Home() {
  const [searchParams] = useSearchParams();
  const page = Number(searchParams.get('page')) || 1;

  const { posts, total, loading, error } = usePosts({ page });

  return (
    <Layout>
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
