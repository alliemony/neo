import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi } from 'vitest';
import { TagInput } from './TagInput';

describe('TagInput', () => {
  it('renders existing tags as removable pills', () => {
    render(<TagInput value={['go', 'python']} onChange={() => {}} existingTags={[]} />);
    expect(screen.getByText('go')).toBeInTheDocument();
    expect(screen.getByText('python')).toBeInTheDocument();
  });

  it('shows suggestions when typing', async () => {
    render(<TagInput value={[]} onChange={() => {}} existingTags={['python', 'go', 'tutorial']} />);

    await userEvent.type(screen.getByPlaceholderText(/add tag/i), 'py');

    expect(screen.getByText('python')).toBeInTheDocument();
    expect(screen.queryByText('go')).not.toBeInTheDocument();
  });

  it('adds tag from suggestion', async () => {
    const onChange = vi.fn();
    render(<TagInput value={[]} onChange={onChange} existingTags={['python', 'go']} />);

    await userEvent.type(screen.getByPlaceholderText(/add tag/i), 'py');
    await userEvent.click(screen.getByText('python'));

    expect(onChange).toHaveBeenCalledWith(['python']);
  });

  it('adds custom tag on Enter', async () => {
    const onChange = vi.fn();
    render(<TagInput value={[]} onChange={onChange} existingTags={[]} />);

    await userEvent.type(screen.getByPlaceholderText(/add tag/i), 'newtag{Enter}');

    expect(onChange).toHaveBeenCalledWith(['newtag']);
  });

  it('removes tag on click', async () => {
    const onChange = vi.fn();
    render(<TagInput value={['go', 'python']} onChange={onChange} existingTags={[]} />);

    await userEvent.click(screen.getAllByRole('button', { name: /×/ })[0]!);

    expect(onChange).toHaveBeenCalledWith(['python']);
  });

  it('does not suggest already selected tags', async () => {
    render(<TagInput value={['python']} onChange={() => {}} existingTags={['python', 'go']} />);

    await userEvent.type(screen.getByPlaceholderText(/add tag/i), 'p');

    expect(screen.queryByRole('option')).not.toBeInTheDocument();
  });
});
