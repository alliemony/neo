import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import App from './App';

describe('App', () => {
  it('renders without crashing', () => {
    render(<App />);
    expect(screen.getByText('neo', { selector: 'h1' })).toBeInTheDocument();
  });

  it('shows the site tagline', () => {
    render(<App />);
    expect(screen.getByText('personal web garden')).toBeInTheDocument();
  });

  it('includes header and footer', () => {
    render(<App />);
    expect(screen.getByText(/built with care/)).toBeInTheDocument();
  });
});
