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

  it('links to the blog filtered by tag', () => {
    renderWithRouter(<TagPill tag="python" />);
    const link = screen.getByRole('link', { name: 'python' });
    expect(link).toHaveAttribute('href', '/blog?tag=python');
  });

  it('applies active styling when active', () => {
    renderWithRouter(<TagPill tag="python" active />);
    const link = screen.getByRole('link', { name: 'python' });
    expect(link).toBeInTheDocument();
    // active uses inline style with inverted colors, not bg-accent class
    expect(link.style.color).toBe('rgb(255, 255, 255)');
  });

  it('renders as span when asSpan is true', () => {
    renderWithRouter(<TagPill tag="python" asSpan />);
    expect(screen.queryByRole('link')).not.toBeInTheDocument();
    expect(screen.getByText('python').tagName).toBe('SPAN');
  });
});
