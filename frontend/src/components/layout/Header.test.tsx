import { render, screen, waitFor } from "@testing-library/react";
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
  it("renders the site name", async () => {
    renderWithRouter(<Header />);
    await waitFor(() => {
      expect(screen.getByText("neo")).toBeInTheDocument();
    });
  });

  it("renders blog navigation link", async () => {
    renderWithRouter(<Header />);
    await waitFor(() => {
      expect(screen.getByText("blog")).toBeInTheDocument();
    });
  });

  it("links site name to home", async () => {
    renderWithRouter(<Header />);
    await waitFor(() => {
      const link = screen.getByText("neo").closest("a");
      expect(link).toHaveAttribute("href", "/");
    });
  });
});
