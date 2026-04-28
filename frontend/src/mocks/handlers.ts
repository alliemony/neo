import { http, HttpResponse } from "msw";

const samplePosts = [
  {
    id: 1,
    slug: "hello-world",
    title: "Hello World",
    content: "# Hello World\n\nFirst mocked post.",
    content_type: "markdown",
    tags: ["python", "intro"],
    published: true,
    like_count: 3,
    created_at: "2026-04-28T00:00:00Z",
    updated_at: "2026-04-28T00:00:00Z",
  },
  {
    id: 2,
    slug: "widget-demo",
    title: "Widget Demo",
    content: "Widget post body",
    content_type: "widget",
    tags: ["widgets"],
    published: true,
    like_count: 1,
    created_at: "2026-04-27T00:00:00Z",
    updated_at: "2026-04-27T00:00:00Z",
  },
];

const samplePages = [
  {
    id: 1,
    slug: "about",
    title: "About",
    content: "# About\n\nMocked page content.",
    published: true,
    sort_order: 1,
    created_at: "2026-04-28T00:00:00Z",
    updated_at: "2026-04-28T00:00:00Z",
  },
];

const sampleTags = [
  { name: "python", count: 1 },
  { name: "intro", count: 1 },
  { name: "widgets", count: 1 },
];

export const handlers = [
  http.get("*/api/v1/auth/mode", () => HttpResponse.json({ mode: "basic" })),
  http.get("*/api/v1/auth/me", () => HttpResponse.json(null, { status: 401 })),
  http.post("*/api/v1/auth/logout", () => new HttpResponse(null, { status: 204 })),

  http.get("*/api/v1/posts", ({ request }) => {
    const url = new URL(request.url);
    const tag = url.searchParams.get("tag");
    const posts = tag
      ? samplePosts.filter((post) => post.tags.includes(tag))
      : samplePosts;

    return HttpResponse.json({ posts, total: posts.length });
  }),

  http.get("*/api/v1/posts/:slug", ({ params }) => {
    const post = samplePosts.find((item) => item.slug === params.slug);
    if (!post) {
      return HttpResponse.json({ error: "not found" }, { status: 404 });
    }

    return HttpResponse.json(post);
  }),

  http.get("*/api/v1/tags", () => HttpResponse.json(sampleTags)),
  http.get("*/api/v1/pages", () => HttpResponse.json(samplePages)),

  http.get("*/api/v1/pages/:slug", ({ params }) => {
    const page = samplePages.find((item) => item.slug === params.slug);
    if (!page) {
      return HttpResponse.json({ error: "not found" }, { status: 404 });
    }

    return HttpResponse.json(page);
  }),

  http.get("*/api/v1/posts/:slug/comments", () => HttpResponse.json([])),
  http.post("*/api/v1/posts/:slug/comments", async ({ request }) => {
    const input = (await request.json()) as { author_name?: string; content?: string };
    return HttpResponse.json(
      {
        id: 1,
        post_id: 1,
        author_name: input.author_name ?? "Guest",
        content: input.content ?? "",
        created_at: "2026-04-28T00:00:00Z",
      },
      { status: 201 },
    );
  }),
  http.post("*/api/v1/posts/:slug/like", () => HttpResponse.json({ like_count: 4 })),

  http.get("*/api/v1/admin/posts", () =>
    HttpResponse.json({ posts: samplePosts, total: samplePosts.length }),
  ),
  http.get("*/api/v1/admin/pages", () => HttpResponse.json(samplePages)),
  http.post("*/api/v1/admin/posts", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ id: 3, slug: "new-post", ...body }, { status: 201 });
  }),
  http.put("*/api/v1/admin/posts/:slug", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ slug: params.slug, ...body });
  }),
  http.delete("*/api/v1/admin/posts/:slug", () => new HttpResponse(null, { status: 204 })),
  http.post("*/api/v1/admin/pages", async ({ request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ id: 2, slug: "new-page", ...body }, { status: 201 });
  }),
  http.put("*/api/v1/admin/pages/:slug", async ({ params, request }) => {
    const body = (await request.json()) as Record<string, unknown>;
    return HttpResponse.json({ slug: params.slug, ...body });
  }),
  http.delete("*/api/v1/admin/pages/:slug", () => new HttpResponse(null, { status: 204 })),
];