import { useParams, useSearchParams, Link } from 'react-router-dom';
import { Layout } from '../components/layout/Layout';
import { PostList } from '../components/blog/PostList';
import { usePosts } from '../hooks/usePosts';

export function TagFeed() {
  const { tag } = useParams<{ tag: string }>();
  const [searchParams] = useSearchParams();
  const page = Number(searchParams.get('page')) || 1;

  const { posts, total, loading, error } = usePosts({ tag, page });

  return (
    <Layout>
      <div className="mb-8">
        <Link to="/" className="text-accent hover:underline text-sm">
          ← All posts
        </Link>
        <h1 className="font-heading text-2xl font-bold mt-2">
          Posts tagged <span className="text-accent">#{tag}</span>
        </h1>
      </div>
      {loading && <p className="text-text-secondary">Loading posts…</p>}
      {error && <p className="text-accent">Failed to load posts.</p>}
      {!loading && !error && (
        <PostList posts={posts} total={total} page={page} />
      )}
    </Layout>
  );
}
