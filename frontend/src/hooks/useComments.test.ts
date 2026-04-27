import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useComments } from './useComments';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('useComments', () => {
  it('fetches and returns comments', async () => {
    const comments = [
      { id: 1, post_id: 1, author_name: 'alice', content: 'Nice!', created_at: '2026-01-01T00:00:00Z' },
    ];
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(comments),
    });

    const { result } = renderHook(() => useComments('test-post'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.comments).toHaveLength(1);
    expect(result.current.comments[0]!.author_name).toBe('alice');
  });

  it('adds a comment via addComment', async () => {
    // Initial fetch.
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    const { result } = renderHook(() => useComments('test-post'));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Mock the POST response.
    const newComment = {
      id: 2, post_id: 1, author_name: 'bob', content: 'Hello!',
      created_at: '2026-01-01T01:00:00Z',
    };
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(newComment),
    });

    await act(async () => {
      await result.current.addComment({ author_name: 'bob', content: 'Hello!' });
    });

    expect(result.current.comments).toHaveLength(1);
    expect(result.current.comments[0]!.author_name).toBe('bob');
  });

  it('starts in loading state', () => {
    mockFetch.mockReturnValue(new Promise(() => {}));
    const { result } = renderHook(() => useComments('test-post'));
    expect(result.current.loading).toBe(true);
  });
});
