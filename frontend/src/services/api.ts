import type { Post, PostListResponse, TagCount, GetPostsParams } from '../types/post';

const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080';

async function fetchJSON<T>(url: string): Promise<T> {
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  return response.json();
}

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
