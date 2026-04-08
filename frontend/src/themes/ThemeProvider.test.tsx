import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect } from 'vitest';
import { ThemeProvider, useTheme } from './ThemeProvider';

function ThemeDisplay() {
  const { theme, setTheme, availableThemes } = useTheme();
  return (
    <div>
      <span data-testid="current-theme">{theme}</span>
      <span data-testid="available">{availableThemes.join(',')}</span>
      <button onClick={() => setTheme('dark')}>Switch to dark</button>
    </div>
  );
}

describe('ThemeProvider', () => {
  it('provides default retro theme', () => {
    render(
      <ThemeProvider>
        <ThemeDisplay />
      </ThemeProvider>,
    );
    expect(screen.getByTestId('current-theme')).toHaveTextContent('retro');
  });

  it('lists available themes', () => {
    render(
      <ThemeProvider>
        <ThemeDisplay />
      </ThemeProvider>,
    );
    expect(screen.getByTestId('available')).toHaveTextContent('retro,dark');
  });

  it('sets CSS variables on document root', () => {
    render(
      <ThemeProvider>
        <ThemeDisplay />
      </ThemeProvider>,
    );
    const root = document.documentElement;
    expect(root.style.getPropertyValue('--color-bg')).toBe('#FAFAF8');
    expect(root.style.getPropertyValue('--color-accent')).toBe('#E85D3A');
    expect(root.style.getPropertyValue('--font-heading')).toBe('"JetBrains Mono", monospace');
    expect(root.style.getPropertyValue('--border-width')).toBe('2px');
  });

  it('updates CSS variables when theme changes', async () => {
    const user = userEvent.setup();
    render(
      <ThemeProvider>
        <ThemeDisplay />
      </ThemeProvider>,
    );

    await act(async () => {
      await user.click(screen.getByText('Switch to dark'));
    });

    expect(screen.getByTestId('current-theme')).toHaveTextContent('dark');
    const root = document.documentElement;
    expect(root.style.getPropertyValue('--color-bg')).toBe('#1A1A1A');
    expect(root.style.getPropertyValue('--color-accent')).toBe('#FF7B5C');
  });

  it('throws when useTheme is used outside provider', () => {
    // Suppress console.error for this test
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {});
    expect(() => render(<ThemeDisplay />)).toThrow(
      'useTheme must be used within a ThemeProvider',
    );
    spy.mockRestore();
  });
});
