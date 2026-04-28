import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { AuthProvider } from "../../contexts/AuthContext";
import { PostEditor } from "./PostEditor";

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
  mockFetch.mockImplementation((url: string) => {
    if (url.includes("/api/v1/pages")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    if (url.includes("/api/v1/tags")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    return Promise.resolve({
      ok: true,
      json: () => Promise.resolve({ posts: [], total: 0 }),
    });
  });
});

afterEach(() => {
  vi.restoreAllMocks();
});

function renderEditor(slug?: string) {
  const path = slug ? `/admin/posts/${slug}` : "/admin/posts/new";
  return render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider
        initialCredentials={{ username: "admin", password: "secret" }}
      >
        <Routes>
          <Route path="/admin/posts/new" element={<PostEditor />} />
          <Route path="/admin/posts/:slug" element={<PostEditor />} />
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  );
}

describe("PostEditor", () => {
  it("renders empty form for new post", async () => {
    renderEditor();
    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toHaveValue("");
      expect(screen.getByLabelText(/content/i)).toHaveValue("");
    });
  });

  it("populates fields when editing existing post", async () => {
    mockFetch.mockImplementation((url: string) => {
      if (url.includes("/api/v1/pages")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (url.includes("/api/v1/tags")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (url.includes("/api/v1/admin/posts")) {
        return Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve({
              posts: [
                {
                  id: 1,
                  slug: "hello",
                  title: "Hello World",
                  content: "# Hello",
                  content_type: "markdown",
                  tags: ["go"],
                  published: true,
                  like_count: 0,
                  comment_count: 0,
                  created_at: "2026-01-01T00:00:00Z",
                  updated_at: "2026-01-01T00:00:00Z",
                },
              ],
              total: 1,
            }),
        });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
    });

    renderEditor("hello");

    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toHaveValue("Hello World");
    });

    expect(screen.getByLabelText(/content/i)).toHaveValue("# Hello");
  });

  it("shows live preview of markdown content", async () => {
    renderEditor();

    await userEvent.type(screen.getByLabelText(/content/i), "**bold text**");

    await waitFor(() => {
      expect(screen.getByText("bold text")).toBeInTheDocument();
    });
  });

  it("has a publish toggle", async () => {
    renderEditor();
    await waitFor(() => {
      expect(screen.getByLabelText(/published/i)).toBeInTheDocument();
    });
  });

  it("saves new post on submit", async () => {
    let createCalled = false;
    mockFetch.mockImplementation((url: string, init?: RequestInit) => {
      if (url.includes("/api/v1/pages")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (url.includes("/api/v1/tags")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (init?.method === "POST") {
        createCalled = true;
        return Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve({
              id: 1,
              slug: "new-post",
              title: "New Post",
              content: "content",
              tags: [],
              published: false,
            }),
        });
      }
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ posts: [], total: 0 }),
      });
    });

    renderEditor();

    await userEvent.type(screen.getByLabelText(/title/i), "New Post");
    await userEvent.type(screen.getByLabelText(/content/i), "content");
    await userEvent.click(screen.getByRole("button", { name: /save/i }));

    await waitFor(() => {
      expect(createCalled).toBe(true);
    });
  });
});
