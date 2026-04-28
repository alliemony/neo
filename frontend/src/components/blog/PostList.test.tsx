import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { describe, it, expect } from "vitest";
import { PostList } from "./PostList";
import type { Post } from "../../types/post";

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

const posts: Post[] = [
  {
    id: 1,
    slug: "post-one",
    title: "Post One",
    content: "Content one",
    content_type: "markdown",
    tags: ["go"],
    published: true,
    like_count: 0,
    created_at: "2026-04-10T10:00:00Z",
    updated_at: "2026-04-10T10:00:00Z",
  },
  {
    id: 2,
    slug: "post-two",
    title: "Post Two",
    content: "Content two",
    content_type: "markdown",
    tags: ["python"],
    published: true,
    like_count: 0,
    created_at: "2026-04-09T10:00:00Z",
    updated_at: "2026-04-09T10:00:00Z",
  },
];

describe("PostList", () => {
  it("renders a list of post cards", () => {
    renderWithRouter(<PostList posts={posts} total={2} />);
    expect(screen.getByText("Post One")).toBeInTheDocument();
    expect(screen.getByText("Post Two")).toBeInTheDocument();
  });

  it("shows empty state when no posts", () => {
    renderWithRouter(<PostList posts={[]} total={0} />);
    expect(screen.getByText(/no posts yet/i)).toBeInTheDocument();
  });

  it("shows pagination info when total exceeds displayed", () => {
    renderWithRouter(
      <PostList posts={posts} total={15} page={1} perPage={10} />,
    );
    expect(screen.getByText(/page 1/i)).toBeInTheDocument();
  });

  it("does not show pagination when all posts fit on one page", () => {
    renderWithRouter(
      <PostList posts={posts} total={2} page={1} perPage={10} />,
    );
    expect(screen.queryByText(/page/i)).not.toBeInTheDocument();
  });
});
