import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { describe, it, expect } from 'vitest';
import { TagPill } from './TagPill';

function renderWithRouter(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>);
}

describe('TagPill', () => {
  it('renders the tag name', () => {
    renderWithRouter(<TagPill tag="python" />);
    expect(screen.getByText('python')).toBeInTheDocument();
  });

  it('links to the tag filter page', () => {
    renderWithRouter(<TagPill tag="python" />);
    const link = screen.getByRole('link', { name: 'python' });
    expect(link).toHaveAttribute('href', '/tag/python');
  });

  it('applies active styling when active', () => {
    renderWithRouter(<TagPill tag="python" active />);
    const link = screen.getByRole('link', { name: 'python' });
    expect(link.className).toContain('bg-accent');
  });
});
