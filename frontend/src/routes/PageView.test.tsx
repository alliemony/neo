import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { PageView } from './PageView';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

function renderPageView(slug: string) {
  return render(
    <MemoryRouter initialEntries={[`/page/${slug}`]}>
      <Routes>
        <Route path="/page/:slug" element={<PageView />} />
      </Routes>
    </MemoryRouter>,
  );
}

describe('PageView', () => {
  it('renders page title and content', async () => {
    mockFetch
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) })
      .mockResolvedValueOnce({
        ok: true,
        json: () =>
          Promise.resolve({
            id: 1,
            slug: 'about',
            title: 'About',
            content: '# About Me\n\nHello world',
            content_type: 'markdown',
            published: true,
          }),
      });

    renderPageView('about');

    await waitFor(() => {
      expect(screen.getByText('About')).toBeInTheDocument();
    });

    expect(screen.getByText('About Me')).toBeInTheDocument();
  });

  it('shows 404 for non-existent page', async () => {
    mockFetch
      .mockResolvedValueOnce({ ok: true, json: () => Promise.resolve([]) })
      .mockResolvedValueOnce({ ok: false, status: 404 });

    renderPageView('nonexistent');

    await waitFor(() => {
      expect(screen.getByText('404')).toBeInTheDocument();
    });
  });

  it('shows loading state', () => {
    mockFetch.mockReturnValue(new Promise(() => {}));
    renderPageView('about');
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });
});
