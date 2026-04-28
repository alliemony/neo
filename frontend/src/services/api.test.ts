import { describe, it, expect } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "../mocks/server";
import { getPosts, getPostBySlug, getTags } from "./api";

describe("getPosts", () => {
  it("fetches all posts", async () => {
    const result = await getPosts();

    expect(result.total).toBe(2);
    expect(result.posts[0]?.title).toBe("Hello World");
  });

  it("passes tag filter as query param", async () => {
    let requestedTag: string | null = null;

    server.use(
      http.get("*/api/v1/posts", ({ request }) => {
        requestedTag = new URL(request.url).searchParams.get("tag");
        return HttpResponse.json({ posts: [], total: 0 });
      }),
    );

    await getPosts({ tag: "python" });
    expect(requestedTag).toBe("python");
  });

  it("passes page param", async () => {
    let requestedPage: string | null = null;

    server.use(
      http.get("*/api/v1/posts", ({ request }) => {
        requestedPage = new URL(request.url).searchParams.get("page");
        return HttpResponse.json({ posts: [], total: 0 });
      }),
    );

    await getPosts({ page: 2 });
    expect(requestedPage).toBe("2");
  });

  it("throws on API error", async () => {
    server.use(
      http.get("*/api/v1/posts", () =>
        HttpResponse.json({ error: "boom" }, { status: 500 }),
      ),
    );

    await expect(getPosts()).rejects.toThrow("API error: 500");
  });
});

describe("getPostBySlug", () => {
  it("fetches a single post by slug", async () => {
    const result = await getPostBySlug("hello-world");

    expect(result.slug).toBe("hello-world");
    expect(result.title).toBe("Hello World");
  });

  it("throws on 404", async () => {
    await expect(getPostBySlug("nonexistent")).rejects.toThrow("API error: 404");
  });
});

describe("getTags", () => {
  it("fetches all tags", async () => {
    const result = await getTags();

    expect(result).toEqual([
      { name: "python", count: 1 },
      { name: "intro", count: 1 },
      { name: "widgets", count: 1 },
    ]);
  });
});
