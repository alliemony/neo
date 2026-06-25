import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { AuthProvider } from "../hooks/useAuth";
import { AdminDashboard } from "./AdminDashboard";

const mockFetch = vi.fn();

interface Overrides {
  recs?: unknown;
  music?: unknown;
}

function mockApi(overrides: Overrides = {}) {
  mockFetch.mockImplementation((url: string, init?: RequestInit) => {
    if (url.includes("/api/v1/auth/mode")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve({ mode: "oauth" }) });
    }
    if (url.includes("/api/v1/auth/me")) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ username: "allie", avatar_url: "", provider: "github" }),
      });
    }
    if (url.includes("/api/v1/admin/recs")) {
      if (init?.method === "POST" || init?.method === "PUT") {
        return Promise.resolve({ ok: true, json: () => Promise.resolve({ id: 99 }) });
      }
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(overrides.recs ?? { recs: [] }),
      });
    }
    if (url.includes("/api/v1/admin/music")) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(overrides.music ?? { recs: [] }),
      });
    }
    if (url.includes("/api/v1/admin/posts")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve({ posts: [], total: 0 }) });
    }
    if (url.includes("/api/v1/admin/pages")) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
  });
}

// Form fields use sibling <label>/<input> without htmlFor association, so query
// the input via its label's containing div.
function fieldInput(labelText: RegExp): HTMLInputElement {
  const label = screen.getByText(labelText);
  const input = label.parentElement?.querySelector("input, textarea");
  if (!input) throw new Error(`no input for label ${labelText}`);
  return input as HTMLInputElement;
}

function renderDashboard() {
  return render(
    <BrowserRouter>
      <AuthProvider>
        <AdminDashboard />
      </AuthProvider>
    </BrowserRouter>,
  );
}

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

const sampleRec = {
  id: 1, name: "Ollama", href: "https://ollama.ai", section: "Tools", tag: "CLI",
  tag_bg: "#edf5f3", tag_color: "#1a7060", description: "Run LLMs locally.",
  published: true, sort_order: 0, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z",
};

const sampleMusic = {
  id: 1, title: "Kind of Blue", artist: "Miles Davis", album: null, cover_url: null,
  spotify_url: null, apple_url: null, note: null, published: true, sort_order: 0,
  created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z",
};

describe("AdminDashboard recs & music", () => {
  it("renders Recs and Music sections with loaded data", async () => {
    mockApi({ recs: { recs: [sampleRec] }, music: { recs: [sampleMusic] } });
    renderDashboard();

    await waitFor(() => {
      expect(screen.getByText("Ollama")).toBeInTheDocument();
    });
    expect(screen.getByText("Kind of Blue")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /new rec/i })).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /new music rec/i })).toBeInTheDocument();
  });

  it("opens the rec form and creates a rec", async () => {
    mockApi({ recs: { recs: [] }, music: { recs: [] } });
    renderDashboard();

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /new rec/i })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole("button", { name: /new rec/i }));
    await userEvent.type(fieldInput(/^Name \*/), "Zed");
    await userEvent.type(fieldInput(/^Section \*/), "Tools");
    await userEvent.type(fieldInput(/^Tag \*/), "Editor");
    await userEvent.type(fieldInput(/^Description \*/), "Fast editor.");

    await userEvent.click(screen.getByRole("button", { name: /^create$/i }));

    await waitFor(() => {
      const postCalls = mockFetch.mock.calls.filter(
        ([url, init]) => url.includes("/api/v1/admin/recs") && init?.method === "POST",
      );
      expect(postCalls.length).toBe(1);
    });
  });

  it("opens the music form on New Music Rec", async () => {
    mockApi({ recs: { recs: [] }, music: { recs: [] } });
    renderDashboard();

    await waitFor(() => {
      expect(screen.getByRole("button", { name: /new music rec/i })).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole("button", { name: /new music rec/i }));
    expect(screen.getByText(/^Title \*/)).toBeInTheDocument();
    expect(screen.getByText(/^Artist \*/)).toBeInTheDocument();
  });
});
