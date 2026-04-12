import type {
  Post,
  PostListResponse,
  TagCount,
  GetPostsParams,
  Comment,
  CreateCommentInput,
  Page,
  CreatePostInput,
  UpdatePostInput,
  CreatePageInput,
  UpdatePageInput,
} from '../types/post';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, options);
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  return response.json();
}

function authHeaders(credentials?: { username: string; password: string }): HeadersInit {
  if (!credentials) return { 'Content-Type': 'application/json' };
  return {
    'Content-Type': 'application/json',
    Authorization: `Basic ${btoa(`${credentials.username}:${credentials.password}`)}`,
  };
}

// Admin fetch helper — uses cookies when no credentials provided (OAuth mode),
// falls back to Basic Auth header when credentials exist (basic mode).
function adminFetchOptions(
  credentials: { username: string; password: string } | null,
  method?: string,
  body?: unknown,
): RequestInit {
  const opts: RequestInit = {
    credentials: 'include', // always send cookies
  };
  if (method) opts.method = method;
  if (credentials) {
    opts.headers = authHeaders(credentials);
  } else {
    opts.headers = { 'Content-Type': 'application/json' };
  }
  if (body !== undefined) {
    opts.body = JSON.stringify(body);
  }
  return opts;
}

// Public post endpoints
export async function getPosts(params?: GetPostsParams): Promise<PostListResponse> {
  const searchParams = new URLSearchParams();
  if (params?.tag) searchParams.set('tag', params.tag);
  if (params?.page) searchParams.set('page', String(params.page));
  if (params?.per_page) searchParams.set('per_page', String(params.per_page));

  const qs = searchParams.toString();
  const url = `${API_BASE}/api/v1/posts${qs ? `?${qs}` : ''}`;
  return fetchJSON<PostListResponse>(url);
}

export async function getPostBySlug(slug: string): Promise<Post> {
  return fetchJSON<Post>(`${API_BASE}/api/v1/posts/${slug}`);
}

export async function getTags(): Promise<TagCount[]> {
  return fetchJSON<TagCount[]>(`${API_BASE}/api/v1/tags`);
}

// Comment endpoints
export async function getComments(slug: string): Promise<Comment[]> {
  return fetchJSON<Comment[]>(`${API_BASE}/api/v1/posts/${slug}/comments`);
}

export async function createComment(slug: string, input: CreateCommentInput): Promise<Comment> {
  return fetchJSON<Comment>(`${API_BASE}/api/v1/posts/${slug}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(input),
  });
}

// Like endpoint
export async function likePost(slug: string): Promise<{ like_count: number }> {
  return fetchJSON<{ like_count: number }>(`${API_BASE}/api/v1/posts/${slug}/like`, {
    method: 'POST',
  });
}

// Public page endpoints
export async function getPages(): Promise<Page[]> {
  return fetchJSON<Page[]>(`${API_BASE}/api/v1/pages`);
}

export async function getPageBySlug(slug: string): Promise<Page> {
  return fetchJSON<Page>(`${API_BASE}/api/v1/pages/${slug}`);
}

// Admin endpoints — credentials is null in OAuth mode (uses cookies)
export async function adminGetPosts(
  credentials: { username: string; password: string } | null,
): Promise<PostListResponse> {
  return fetchJSON<PostListResponse>(
    `${API_BASE}/api/v1/admin/posts`,
    adminFetchOptions(credentials),
  );
}

export async function adminCreatePost(
  credentials: { username: string; password: string } | null,
  input: CreatePostInput,
): Promise<Post> {
  return fetchJSON<Post>(
    `${API_BASE}/api/v1/admin/posts`,
    adminFetchOptions(credentials, 'POST', input),
  );
}

export async function adminUpdatePost(
  credentials: { username: string; password: string } | null,
  slug: string,
  input: UpdatePostInput,
): Promise<Post> {
  return fetchJSON<Post>(
    `${API_BASE}/api/v1/admin/posts/${slug}`,
    adminFetchOptions(credentials, 'PUT', input),
  );
}

export async function adminDeletePost(
  credentials: { username: string; password: string } | null,
  slug: string,
): Promise<void> {
  const response = await fetch(
    `${API_BASE}/api/v1/admin/posts/${slug}`,
    adminFetchOptions(credentials, 'DELETE'),
  );
  if (!response.ok) throw new Error(`API error: ${response.status}`);
}

export async function adminGetPages(
  credentials: { username: string; password: string } | null,
): Promise<Page[]> {
  return fetchJSON<Page[]>(
    `${API_BASE}/api/v1/admin/pages`,
    adminFetchOptions(credentials),
  );
}

export async function adminCreatePage(
  credentials: { username: string; password: string } | null,
  input: CreatePageInput,
): Promise<Page> {
  return fetchJSON<Page>(
    `${API_BASE}/api/v1/admin/pages`,
    adminFetchOptions(credentials, 'POST', input),
  );
}

export async function adminUpdatePage(
  credentials: { username: string; password: string } | null,
  slug: string,
  input: UpdatePageInput,
): Promise<Page> {
  return fetchJSON<Page>(
    `${API_BASE}/api/v1/admin/pages/${slug}`,
    adminFetchOptions(credentials, 'PUT', input),
  );
}

export async function adminDeletePage(
  credentials: { username: string; password: string } | null,
  slug: string,
): Promise<void> {
  const response = await fetch(
    `${API_BASE}/api/v1/admin/pages/${slug}`,
    adminFetchOptions(credentials, 'DELETE'),
  );
  if (!response.ok) throw new Error(`API error: ${response.status}`);
}
