import { Link } from 'react-router-dom';
import type { Post } from '../../types/post';
import { PostCard } from './PostCard';

interface PostListProps {
  posts: Post[];
  total: number;
  page?: number;
  perPage?: number;
}

export function PostList({ posts, total, page = 1, perPage = 10 }: PostListProps) {
  if (posts.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-text-secondary text-lg">No posts yet.</p>
      </div>
    );
  }

  const totalPages = Math.ceil(total / perPage);

  return (
    <div className="space-y-6">
      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}
      {totalPages > 1 && (
        <nav className="flex items-center justify-center gap-4 pt-4">
          {page > 1 && (
            <Link
              to={`?page=${page - 1}`}
              className="text-accent hover:underline"
            >
              ← Prev
            </Link>
          )}
          <span className="text-text-secondary text-sm">
            Page {page} of {totalPages}
          </span>
          {page < totalPages && (
            <Link
              to={`?page=${page + 1}`}
              className="text-accent hover:underline"
            >
              Next →
            </Link>
          )}
        </nav>
      )}
    </div>
  );
}
