import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { AuthProvider, useAuth } from './AuthContext';

const mockFetch = vi.fn();

beforeEach(() => {
  vi.stubGlobal('fetch', mockFetch);
});

afterEach(() => {
  vi.restoreAllMocks();
});

function TestConsumer() {
  const { credentials, login, logout, authHeader } = useAuth();
  return (
    <div>
      <span data-testid="status">{credentials ? 'logged-in' : 'logged-out'}</span>
      <span data-testid="header">{authHeader() || 'none'}</span>
      <button onClick={() => login('admin', 'secret')}>Login</button>
      <button onClick={logout}>Logout</button>
    </div>
  );
}

describe('AuthContext', () => {
  it('starts logged out', () => {
    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );
    expect(screen.getByTestId('status')).toHaveTextContent('logged-out');
    expect(screen.getByTestId('header')).toHaveTextContent('none');
  });

  it('logs in on successful auth', async () => {
    mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    await userEvent.click(screen.getByText('Login'));

    await waitFor(() => {
      expect(screen.getByTestId('status')).toHaveTextContent('logged-in');
    });

    expect(screen.getByTestId('header')).toHaveTextContent('Basic');
  });

  it('stays logged out on failed auth', async () => {
    mockFetch.mockResolvedValueOnce({ ok: false, status: 401 });

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    await userEvent.click(screen.getByText('Login'));

    await waitFor(() => {
      expect(mockFetch).toHaveBeenCalled();
    });

    expect(screen.getByTestId('status')).toHaveTextContent('logged-out');
  });

  it('clears credentials on logout', async () => {
    mockFetch.mockResolvedValueOnce({ ok: true, json: () => Promise.resolve({}) });

    render(
      <AuthProvider>
        <TestConsumer />
      </AuthProvider>,
    );

    await userEvent.click(screen.getByText('Login'));
    await waitFor(() => {
      expect(screen.getByTestId('status')).toHaveTextContent('logged-in');
    });

    await userEvent.click(screen.getByText('Logout'));
    expect(screen.getByTestId('status')).toHaveTextContent('logged-out');
  });
});
