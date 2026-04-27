import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { CommentForm } from './CommentForm';

describe('CommentForm', () => {
  it('renders name and content fields', () => {
    render(<CommentForm onSubmit={vi.fn()} />);
    expect(screen.getByPlaceholderText(/name/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/comment/i)).toBeInTheDocument();
  });

  it('calls onSubmit with form data', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<CommentForm onSubmit={onSubmit} />);

    await user.type(screen.getByPlaceholderText(/name/i), 'alice');
    await user.type(screen.getByPlaceholderText(/comment/i), 'Great post!');
    await user.click(screen.getByRole('button', { name: /submit/i }));

    expect(onSubmit).toHaveBeenCalledWith({
      author_name: 'alice',
      content: 'Great post!',
    });
  });

  it('clears the form after successful submit', async () => {
    const user = userEvent.setup();
    const onSubmit = vi.fn().mockResolvedValue(undefined);

    render(<CommentForm onSubmit={onSubmit} />);

    await user.type(screen.getByPlaceholderText(/name/i), 'alice');
    await user.type(screen.getByPlaceholderText(/comment/i), 'Great post!');
    await user.click(screen.getByRole('button', { name: /submit/i }));

    expect(screen.getByPlaceholderText(/comment/i)).toHaveValue('');
  });

  it('shows error message when provided', () => {
    render(<CommentForm onSubmit={vi.fn()} error="Something went wrong" />);
    expect(screen.getByText('Something went wrong')).toBeInTheDocument();
  });
});
