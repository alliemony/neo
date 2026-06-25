import { render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { Layout } from "./Layout";

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

describe("Layout", () => {
  it("renders children in main content area", async () => {
    renderWithRouter(
      <Layout>
        <p>Main content</p>
      </Layout>,
    );
    await waitFor(() => {
      expect(screen.getByText("Main content")).toBeInTheDocument();
    });
  });

  it("renders header and footer", async () => {
    renderWithRouter(
      <Layout>
        <p>Content</p>
      </Layout>,
    );
    await waitFor(() => {
      expect(screen.getByText("allieg.dev")).toBeInTheDocument();
      expect(screen.getByText(/made with too much coffee/)).toBeInTheDocument();
    });
  });

  it("renders without sidebar", async () => {
    renderWithRouter(
      <Layout>
        <p>Only main</p>
      </Layout>,
    );
    await waitFor(() => {
      expect(screen.getByText("Only main")).toBeInTheDocument();
      expect(screen.queryByRole("complementary")).not.toBeInTheDocument();
    });
  });
});
