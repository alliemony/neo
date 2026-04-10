import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import { getPosts, getPostBySlug, getTags } from './api';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('getPosts', () => {
  it('fetches all posts', async () => {
    const data = { posts: [{ id: 1, title: 'Test' }], total: 1 };
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(data),
    });

    const result = await getPosts();
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/posts'),
    );
    expect(result).toEqual(data);
  });

  it('passes tag filter as query param', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ posts: [], total: 0 }),
    });

    await getPosts({ tag: 'python' });
    const url = mockFetch.mock.calls[0]![0] as string;
    expect(url).toContain('tag=python');
  });

  it('passes page param', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ posts: [], total: 0 }),
    });

    await getPosts({ page: 2 });
    const url = mockFetch.mock.calls[0]![0] as string;
    expect(url).toContain('page=2');
  });

  it('throws on API error', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
    });

    await expect(getPosts()).rejects.toThrow('API error: 500');
  });
});

describe('getPostBySlug', () => {
  it('fetches a single post by slug', async () => {
    const post = { id: 1, slug: 'hello-world', title: 'Hello World' };
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(post),
    });

    const result = await getPostBySlug('hello-world');
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/posts/hello-world'),
    );
    expect(result).toEqual(post);
  });

  it('throws on 404', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 404,
      statusText: 'Not Found',
    });

    await expect(getPostBySlug('nonexistent')).rejects.toThrow('API error: 404');
  });
});

describe('getTags', () => {
  it('fetches all tags', async () => {
    const tags = [{ name: 'python', count: 2 }];
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(tags),
    });

    const result = await getTags();
    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/tags'),
    );
    expect(result).toEqual(tags);
  });
});
