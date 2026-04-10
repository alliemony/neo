import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import App from './App';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
  mockFetch.mockResolvedValue({
    ok: true,
    json: () => Promise.resolve({ posts: [], total: 0 }),
  });
});

afterEach(() => {
  vi.restoreAllMocks();
});

describe('App', () => {
  it('renders without crashing', async () => {
    render(<App />);
    await waitFor(() => {
      expect(screen.getByText('neo', { selector: 'h1' })).toBeInTheDocument();
    });
  });

  it('shows the site tagline', async () => {
    render(<App />);
    await waitFor(() => {
      expect(screen.getByText('personal web garden')).toBeInTheDocument();
    });
  });

  it('includes header and footer', async () => {
    render(<App />);
    await waitFor(() => {
      expect(screen.getByText(/built with care/)).toBeInTheDocument();
    });
  });
});
