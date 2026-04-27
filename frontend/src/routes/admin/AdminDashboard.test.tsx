import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { AuthProvider } from '../../contexts/AuthContext';
import { AdminDashboard } from './AdminDashboard';

const mockFetch = vi.fn();

function mockApi(overrides: Record<string, unknown> = {}) {
  mockFetch.mockImplementation((url: string) => {
    if (url.includes('/api/v1/admin/posts') && !url.includes('comments')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(overrides.adminPosts ?? { posts: [], total: 0 }),
      });
    }
    if (url.includes('/api/v1/admin/pages')) {
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve(overrides.adminPages ?? []),
      });
    }
    if (url.includes('/api/v1/pages')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    if (url.includes('/api/v1/tags')) {
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    }
    return Promise.resolve({ ok: true, json: () => Promise.resolve({}) });
  });
}

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

function renderDashboard() {
  return render(
    <MemoryRouter>
      <AuthProvider initialCredentials={{ username: 'admin', password: 'secret' }}>
        <AdminDashboard />
      </AuthProvider>
    </MemoryRouter>,
  );
}

describe('AdminDashboard', () => {
  it('shows loading state initially', () => {
    mockFetch.mockReturnValue(new Promise(() => {}));
    renderDashboard();
    expect(screen.getByText(/loading/i)).toBeInTheDocument();
  });

  it('renders posts and pages after loading', async () => {
    mockApi({
      adminPosts: {
        posts: [
          { id: 1, slug: 'post-one', title: 'Post One', published: true, tags: [], created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
          { id: 2, slug: 'draft-post', title: 'Draft Post', published: false, tags: [], created_at: '2026-01-02T00:00:00Z', updated_at: '2026-01-02T00:00:00Z' },
        ],
        total: 2,
      },
      adminPages: [
        { id: 1, slug: 'about', title: 'About', published: true, created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
      ],
    });

    renderDashboard();

    await waitFor(() => {
      expect(screen.getByText('Post One')).toBeInTheDocument();
    });

    expect(screen.getByText('Draft Post')).toBeInTheDocument();
    expect(screen.getByText('About')).toBeInTheDocument();
  });

  it('shows published and draft status indicators', async () => {
    mockApi({
      adminPosts: {
        posts: [
          { id: 1, slug: 'pub', title: 'Published Post', published: true, tags: [], created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
          { id: 2, slug: 'draft', title: 'Draft Post', published: false, tags: [], created_at: '2026-01-02T00:00:00Z', updated_at: '2026-01-02T00:00:00Z' },
        ],
        total: 2,
      },
    });

    renderDashboard();

    await waitFor(() => {
      expect(screen.getByText('Published Post')).toBeInTheDocument();
    });

    const statuses = screen.getAllByText(/published|draft/i);
    expect(statuses.length).toBeGreaterThanOrEqual(2);
  });

  it('has edit and delete buttons for posts', async () => {
    mockApi({
      adminPosts: {
        posts: [
          { id: 1, slug: 'test', title: 'Test Post', published: true, tags: [], created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' },
        ],
        total: 1,
      },
    });

    renderDashboard();

    await waitFor(() => {
      expect(screen.getByText('Test Post')).toBeInTheDocument();
    });

    expect(screen.getByRole('link', { name: /edit/i })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /delete/i })).toBeInTheDocument();
  });

  it('shows delete confirmation and deletes on confirm', async () => {
    let deleteCallCount = 0;
    mockFetch.mockImplementation((url: string, init?: RequestInit) => {
      if (init?.method === 'DELETE') {
        deleteCallCount++;
        return Promise.resolve({ ok: true, status: 204 });
      }
      if (url.includes('/api/v1/admin/posts')) {
        const posts = deleteCallCount === 0
          ? [{ id: 1, slug: 'test', title: 'Test Post', published: true, tags: [], created_at: '2026-01-01T00:00:00Z', updated_at: '2026-01-01T00:00:00Z' }]
          : [];
        const total = deleteCallCount === 0 ? 1 : 0;
        return Promise.resolve({ ok: true, json: () => Promise.resolve({ posts, total }) });
      }
      if (url.includes('/api/v1/admin/pages')) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      if (url.includes('/api/v1/pages')) {
        return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
      }
      return Promise.resolve({ ok: true, json: () => Promise.resolve([]) });
    });

    renderDashboard();

    await waitFor(() => {
      expect(screen.getByText('Test Post')).toBeInTheDocument();
    });

    await userEvent.click(screen.getByRole('button', { name: /delete/i }));
    expect(screen.getByText(/are you sure/i)).toBeInTheDocument();

    await userEvent.click(screen.getByRole('button', { name: /confirm/i }));

    await waitFor(() => {
      expect(screen.queryByText('Test Post')).not.toBeInTheDocument();
    });
  });
});
