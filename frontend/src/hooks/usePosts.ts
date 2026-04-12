import { useState, useEffect } from 'react';
import { getPosts } from '../services/api';
import type { Post, GetPostsParams } from '../types/post';

interface UsePostsResult {
  posts: Post[];
  total: number;
  loading: boolean;
  error: string | null;
}

export function usePosts(params?: GetPostsParams): UsePostsResult {
  const [posts, setPosts] = useState<Post[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    setLoading(true);
    setError(null);

    getPosts(params)
      .then((data) => {
        if (!cancelled) {
          setPosts(data.posts);
          setTotal(data.total);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err.message);
          setPosts([]);
          setTotal(0);
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [params?.tag, params?.page, params?.per_page]);

  return { posts, total, loading, error };
}
