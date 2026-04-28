import { renderHook, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { usePages } from './usePages';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('usePages', () => {
  it('fetches and returns published pages', async () => {
    const pages = [
      { id: 1, slug: 'about', title: 'About', content: 'a', content_type: 'markdown', published: true, sort_order: 1 },
      { id: 2, slug: 'projects', title: 'Projects', content: 'b', content_type: 'markdown', published: true, sort_order: 2 },
    ];
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(pages),
    });

    const { result } = renderHook(() => usePages());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.pages).toHaveLength(2);
    expect(result.current.pages[0]!.title).toBe('About');
  });

  it('starts in loading state', () => {
    mockFetch.mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => usePages());
    expect(result.current.loading).toBe(true);
  });

  it('handles fetch errors', async () => {
    mockFetch.mockRejectedValueOnce(new Error('network error'));

    const { result } = renderHook(() => usePages());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.pages).toHaveLength(0);
  });
});
