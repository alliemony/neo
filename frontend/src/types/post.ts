export interface Post {
  id: number;
  slug: string;
  title: string;
  content: string;
  content_type: string;
  tags: string[];
  published: boolean;
  like_count: number;
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

export interface Comment {
  id: number;
  post_id: number;
  author_name: string;
  content: string;
  created_at: string;
}

export interface Page {
  id: number;
  slug: string;
  title: string;
  content: string;
  published: boolean;
  sort_order: number;
  created_at: string;
  updated_at: string;
}

export interface CreateCommentInput {
  author_name: string;
  content: string;
}

export interface CreatePostInput {
  title: string;
  content: string;
  content_type: string;
  tags: string[];
  published: boolean;
}

export interface UpdatePostInput {
  title?: string;
  content?: string;
  content_type?: string;
  tags?: string[];
  published?: boolean;
}

export interface CreatePageInput {
  title: string;
  content: string;
  published: boolean;
  sort_order: number;
}

export interface UpdatePageInput {
  title?: string;
  content?: string;
  published?: boolean;
  sort_order?: number;
}
