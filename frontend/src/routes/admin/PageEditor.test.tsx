import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { AuthProvider } from "../../contexts/AuthContext";
import { PageEditor } from "./PageEditor";

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
  mockFetch.mockImplementation((url: string) => {
    if (url.includes("/api/v1/pages") && !url.includes("/admin/")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    if (url.includes("/api/v1/admin/pages")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
  });
});

afterEach(() => {
  vi.restoreAllMocks();
});

function renderEditor(slug?: string) {
  const path = slug ? `/admin/pages/${slug}` : "/admin/pages/new";
  return render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider
        initialCredentials={{ username: "admin", password: "secret" }}
      >
        <Routes>
          <Route path="/admin/pages/new" element={<PageEditor />} />
          <Route path="/admin/pages/:slug" element={<PageEditor />} />
        </Routes>
      </AuthProvider>
    </MemoryRouter>,
  );
}

describe("PageEditor", () => {
  it("renders empty form for new page", async () => {
    renderEditor();
    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toHaveValue("");
      expect(screen.getByLabelText(/content/i)).toHaveValue("");
    });
  });

  it("populates fields when editing existing page", async () => {
    mockFetch.mockImplementation((url: string) => {
      if (url.includes("/api/v1/admin/pages")) {
        return Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve([
              {
                id: 1,
                slug: "about",
                title: "About",
                content: "# About page",
                content_type: "markdown",
                published: true,
                sort_order: 1,
                created_at: "2026-01-01T00:00:00Z",
                updated_at: "2026-01-01T00:00:00Z",
              },
            ]),
        });
      }
      if (url.includes("/api/v1/pages")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
    });

    renderEditor("about");

    await waitFor(() => {
      expect(screen.getByLabelText(/title/i)).toHaveValue("About");
    });

    expect(screen.getByLabelText(/content/i)).toHaveValue("# About page");
  });

  it("shows live preview", async () => {
    renderEditor();

    await userEvent.type(screen.getByLabelText(/content/i), "# Page Title");

    await waitFor(() => {
      expect(screen.getByText("Page Title")).toBeInTheDocument();
    });
  });

  it("saves new page on submit", async () => {
    let createCalled = false;
    mockFetch.mockImplementation((url: string, init?: RequestInit) => {
      if (url.includes("/api/v1/pages") && !url.includes("/admin/")) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (init?.method === "POST") {
        createCalled = true;
        return Promise.resolve({
          ok: true,
          json: () =>
            Promise.resolve({
              id: 1,
              slug: "new-page",
              title: "New Page",
              content: "content",
              content_type: "markdown",
              published: false,
              sort_order: 0,
            }),
        });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    });

    renderEditor();

    await userEvent.type(screen.getByLabelText(/title/i), "New Page");
    await userEvent.type(screen.getByLabelText(/content/i), "content");
    await userEvent.click(screen.getByRole("button", { name: /save/i }));

    await waitFor(() => {
      expect(createCalled).toBe(true);
    });
  });
});
