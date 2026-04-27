import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { CommentSection } from './CommentSection';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('CommentSection', () => {
  it('renders comments after loading', async () => {
    const comments = [
      { id: 1, post_id: 1, author_name: 'alice', content: 'Great post!', created_at: '2026-04-10T10:00:00Z' },
      { id: 2, post_id: 1, author_name: 'bob', content: 'Thanks!', created_at: '2026-04-10T11:00:00Z' },
    ];
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(comments),
    });

    render(<CommentSection slug="test-post" />);

    await waitFor(() => {
      expect(screen.getByText('alice')).toBeInTheDocument();
      expect(screen.getByText('Great post!')).toBeInTheDocument();
      expect(screen.getByText('bob')).toBeInTheDocument();
    });
  });

  it('shows comment count in heading', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([
        { id: 1, post_id: 1, author_name: 'alice', content: 'Nice!', created_at: '2026-01-01T00:00:00Z' },
      ]),
    });

    render(<CommentSection slug="test-post" />);

    await waitFor(() => {
      expect(screen.getByText(/comments\s*\(1\)/i)).toBeInTheDocument();
    });
  });

  it('shows the comment form', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    render(<CommentSection slug="test-post" />);

    await waitFor(() => {
      expect(screen.getByPlaceholderText(/name/i)).toBeInTheDocument();
      expect(screen.getByPlaceholderText(/comment/i)).toBeInTheDocument();
    });
  });

  it('shows empty state when no comments', async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve([]),
    });

    render(<CommentSection slug="test-post" />);

    await waitFor(() => {
      expect(screen.getByText(/no comments yet/i)).toBeInTheDocument();
    });
  });
});
