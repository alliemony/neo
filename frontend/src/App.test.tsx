import { render, screen, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { HelmetProvider } from "react-helmet-async";
import App from "./App";

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
  mockFetch.mockImplementation((url: string) => {
    if (url.includes("/api/v1/auth/mode")) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ mode: "basic" }),
      });
    }
    if (url.includes("/api/v1/pages")) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve([]),
      });
    }
    if (url.includes("/api/v1/tags")) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve([]),
      });
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

describe("App", () => {
  it("renders without crashing", async () => {
    render(
      <HelmetProvider>
        <App />
      </HelmetProvider>,
    );
    await waitFor(() => {
      expect(screen.getByText("neo", { selector: "h1" })).toBeInTheDocument();
    });
  });

  it("shows the site tagline", async () => {
    render(
      <HelmetProvider>
        <App />
      </HelmetProvider>,
    );
    await waitFor(() => {
      expect(screen.getByText("personal web garden")).toBeInTheDocument();
    });
  });

  it("includes header and footer", async () => {
    render(
      <HelmetProvider>
        <App />
      </HelmetProvider>,
    );
    await waitFor(() => {
      expect(screen.getByText(/built with care/)).toBeInTheDocument();
    });
  });
});
