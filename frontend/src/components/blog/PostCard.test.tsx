import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect, vi, afterEach } from "vitest";
import { PostCard } from "./PostCard";
import type { Post } from "../../types/post";

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

const samplePost: Post = {
  id: 1,
  slug: "hello-world",
  title: "Hello World",
  content:
    "This is the first post with some content that should get truncated after a certain length to show as an excerpt in the card view.",
  content_type: "markdown",
  tags: ["go", "intro"],
  published: true,
  created_at: "2026-04-10T10:00:00Z",
  updated_at: "2026-04-10T10:00:00Z",
};

describe("PostCard", () => {
  afterEach(() => {
    vi.useRealTimers();
  });

  it("renders the post title as a link", () => {
    renderWithRouter(<PostCard post={samplePost} />);
    const link = screen.getByRole("link", { name: "Hello World" });
    expect(link).toHaveAttribute("href", "/blog/hello-world");
  });

  it("renders a content excerpt", () => {
    renderWithRouter(<PostCard post={samplePost} />);
    expect(screen.getByText(/This is the first post/)).toBeInTheDocument();
  });

  it("renders tags as pill badges", () => {
    renderWithRouter(<PostCard post={samplePost} />);
    expect(screen.getByText("go")).toBeInTheDocument();
    expect(screen.getByText("intro")).toBeInTheDocument();
  });

  it("renders a relative timestamp", () => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2026-04-10T12:00:00Z"));

    renderWithRouter(<PostCard post={samplePost} />);
    expect(screen.getByText("2 hours ago")).toBeInTheDocument();
  });

  it("uses a time element with full datetime", () => {
    renderWithRouter(<PostCard post={samplePost} />);
    const timeEl = screen.getByRole("time");
    expect(timeEl).toHaveAttribute("dateTime", samplePost.created_at);
  });

  it("renders with retro styling (border, no radius)", () => {
    const { container } = renderWithRouter(<PostCard post={samplePost} />);
    const card = container.querySelector("article");
    expect(card).toBeInTheDocument();
    expect(card?.className).toContain("border-border");
  });
});
