export interface Post {
  id: number;
  slug: string;
  title: string;
  content: string;
  content_type: string;
  tags: string[];
  published: boolean;
  created_at: string;
  updated_at: string;
}

export interface PostListResponse {
  posts: Post[];
  total: number;
}

export interface TagCount {
  name: string;
  count: number;
}

export interface GetPostsParams {
  tag?: string;
  page?: number;
  per_page?: number;
}
