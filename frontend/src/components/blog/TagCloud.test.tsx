import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { TagCloud } from "./TagCloud";

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal("fetch", mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe("TagCloud", () => {
  it("renders tags as links", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () =>
        Promise.resolve([
          { name: "python", count: 10 },
          { name: "go", count: 5 },
          { name: "tutorial", count: 2 },
        ]),
    });

    render(
      <MemoryRouter>
        <TagCloud />
      </MemoryRouter>,
    );

    const pythonLink = await screen.findByRole("link", { name: /python/ });
    expect(pythonLink).toBeInTheDocument();
    expect(pythonLink).toHaveAttribute("href", "/tag/python");

    expect(screen.getByRole("link", { name: /go/ })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /tutorial/ })).toBeInTheDocument();
  });

  it("shows heading", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([{ name: "go", count: 1 }]),
    });

    render(
      <MemoryRouter>
        <TagCloud />
      </MemoryRouter>,
    );

    expect(await screen.findByText("Tags")).toBeInTheDocument();
  });

  it("applies larger size to more popular tags", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () =>
        Promise.resolve([
          { name: "popular", count: 10 },
          { name: "rare", count: 1 },
        ]),
    });

    render(
      <MemoryRouter>
        <TagCloud />
      </MemoryRouter>,
    );

    const popular = await screen.findByRole("link", { name: /popular/ });
    const rare = screen.getByRole("link", { name: /rare/ });

    const popularSize = parseFloat(popular.style.fontSize);
    const rareSize = parseFloat(rare.style.fontSize);
    expect(popularSize).toBeGreaterThan(rareSize);
  });
});
