import { useState, useEffect, useCallback } from 'react';
import { getComments, createComment } from '../services/api';
import type { Comment, CreateCommentInput } from '../types/post';

interface UseCommentsResult {
  comments: Comment[];
  loading: boolean;
  error: string | null;
  addComment: (input: CreateCommentInput) => Promise<void>;
}

export function useComments(slug: string | undefined): UseCommentsResult {
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!slug) return;

    let cancelled = false;
    setLoading(true);
    setError(null);

    getComments(slug)
      .then((data) => {
        if (!cancelled) {
          setComments(data || []);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err.message);
          setComments([]);
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, [slug]);

  const addComment = useCallback(
    async (input: CreateCommentInput) => {
      if (!slug) return;
      const comment = await createComment(slug, input);
      setComments((prev) => [...prev, comment]);
    },
    [slug],
  );

  return { comments, loading, error, addComment };
}
