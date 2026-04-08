import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect } from 'vitest';
import { Header } from './Header';

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('Header', () => {
  it('renders the site name', () => {
    renderWithRouter(<Header />);
    expect(screen.getByText('neo')).toBeInTheDocument();
  });

  it('renders blog navigation link', () => {
    renderWithRouter(<Header />);
    expect(screen.getByText('blog')).toBeInTheDocument();
  });

  it('links site name to home', () => {
    renderWithRouter(<Header />);
    const link = screen.getByText('neo').closest('a');
    expect(link).toHaveAttribute('href', '/');
  });
});
