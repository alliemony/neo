import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { LikeButton } from './LikeButton';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('LikeButton', () => {
  it('renders the initial count', () => {
    render(<LikeButton slug="test" initialCount={12} />);
    expect(screen.getByText('12')).toBeInTheDocument();
  });

  it('shows optimistic update on click', async () => {
    const user = userEvent.setup();
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ like_count: 13 }),
    });

    render(<LikeButton slug="test" initialCount={12} />);

    await user.click(screen.getByRole('button'));

    // Should optimistically show 13.
    await waitFor(() => {
      expect(screen.getByText('13')).toBeInTheDocument();
    });
  });

  it('disables after clicking', async () => {
    const user = userEvent.setup();
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve({ like_count: 13 }),
    });

    render(<LikeButton slug="test" initialCount={12} />);

    const button = screen.getByRole('button');
    await user.click(button);

    await waitFor(() => {
      expect(button).toBeDisabled();
    });
  });

  it('reverts on API failure', async () => {
    const user = userEvent.setup();
    mockFetch.mockResolvedValueOnce({
      ok: false,
      status: 500,
    });

    render(<LikeButton slug="test" initialCount={12} />);

    await user.click(screen.getByRole('button'));

    // Should revert to 12 after failure.
    await waitFor(() => {
      expect(screen.getByText('12')).toBeInTheDocument();
    });

    // Button should be re-enabled.
    expect(screen.getByRole('button')).not.toBeDisabled();
  });
});
