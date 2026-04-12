import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { Layout } from "./Layout";
import { Sidebar } from "./Sidebar";

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
  it("renders children in main content area", () => {
    renderWithRouter(
      <Layout>
        <p>Main content</p>
      </Layout>,
    );
    expect(screen.getByText("Main content")).toBeInTheDocument();
  });

  it("renders header and footer", () => {
    renderWithRouter(
      <Layout>
        <p>Content</p>
      </Layout>,
    );
    expect(screen.getByText("neo")).toBeInTheDocument();
    expect(screen.getByText(/built with care/)).toBeInTheDocument();
  });

  it("renders sidebar when provided", () => {
    renderWithRouter(
      <Layout
        sidebar={
          <Sidebar>
            <p>Sidebar content</p>
          </Sidebar>
        }
      >
        <p>Main content</p>
      </Layout>,
    );
    expect(screen.getByText("Sidebar content")).toBeInTheDocument();
    expect(screen.getByText("Main content")).toBeInTheDocument();
  });

  it("renders without sidebar", () => {
    renderWithRouter(
      <Layout>
        <p>Only main</p>
      </Layout>,
    );
    expect(screen.getByText("Only main")).toBeInTheDocument();
    expect(screen.queryByRole("complementary")).not.toBeInTheDocument();
  });
});
