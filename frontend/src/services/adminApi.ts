import type {
  Post,
  PostListResponse,
  CreatePostInput,
  UpdatePostInput,
  Page,
  CreatePageInput,
  UpdatePageInput,
} from '../types/post';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function adminFetch<T>(url: string, authHeader: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    Authorization: authHeader,
    ...((init?.headers as Record<string, string>) || {}),
  };

  const response = await fetch(url, { ...init, headers });
  if (response.status === 401) {
    throw new Error('unauthorized');
  }
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  if (response.status === 204) {
    return undefined as T;
  }
  return response.json();
}

export function createAdminApi(authHeader: string) {
  const base = `${API_BASE}/api/v1/admin`;

  return {
    listPosts: (page = 1, perPage = 20): Promise<PostListResponse> =>
      adminFetch(`${base}/posts?page=${page}&per_page=${perPage}`, authHeader),

    createPost: (input: CreatePostInput): Promise<Post> =>
      adminFetch(`${base}/posts`, authHeader, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      }),

    updatePost: (slug: string, input: UpdatePostInput): Promise<Post> =>
      adminFetch(`${base}/posts/${slug}`, authHeader, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      }),

    deletePost: (slug: string): Promise<void> =>
      adminFetch(`${base}/posts/${slug}`, authHeader, { method: 'DELETE' }),

    listPages: (): Promise<Page[]> =>
      adminFetch(`${base}/pages`, authHeader),

    createPage: (input: CreatePageInput): Promise<Page> =>
      adminFetch(`${base}/pages`, authHeader, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      }),

    updatePage: (slug: string, input: UpdatePageInput): Promise<Page> =>
      adminFetch(`${base}/pages/${slug}`, authHeader, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(input),
      }),

    deletePage: (slug: string): Promise<void> =>
      adminFetch(`${base}/pages/${slug}`, authHeader, { method: 'DELETE' }),
  };
}
