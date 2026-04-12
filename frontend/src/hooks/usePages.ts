import { useState, useEffect } from 'react';
import { getPages } from '../services/api';
import type { Page } from '../types/post';

interface UsePagesResult {
  pages: Page[];
  loading: boolean;
  error: string | null;
}

export function usePages(): UsePagesResult {
  const [pages, setPages] = useState<Page[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);

    getPages()
      .then((data) => {
        if (!cancelled) {
          setPages(data || []);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err.message);
          setPages([]);
          setLoading(false);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  return { pages, loading, error };
}
