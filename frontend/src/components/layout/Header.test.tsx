import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { Header } from "./Header";

beforeEach(() => {
  vi.stubGlobal(
    "fetch",
    vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve([]),
    }),
  );
});

afterEach(() => {
  vi.restoreAllMocks();
});

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe("Header", () => {
  it("renders the site name", () => {
    renderWithRouter(<Header />);
    expect(screen.getByText("neo")).toBeInTheDocument();
  });

  it("renders blog navigation link", () => {
    renderWithRouter(<Header />);
    expect(screen.getByText("blog")).toBeInTheDocument();
  });

  it("links site name to home", () => {
    renderWithRouter(<Header />);
    const link = screen.getByText("neo").closest("a");
    expect(link).toHaveAttribute("href", "/");
  });
});
