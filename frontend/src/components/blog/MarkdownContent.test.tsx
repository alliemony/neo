import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { MarkdownContent } from './MarkdownContent';

describe('MarkdownContent', () => {
  it('renders markdown headings', () => {
    render(<MarkdownContent content="# Hello World" />);
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent(
      'Hello World',
    );
  });

  it('renders markdown paragraphs', () => {
    render(<MarkdownContent content="This is a paragraph." />);
    expect(screen.getByText('This is a paragraph.')).toBeInTheDocument();
  });

  it('renders markdown links', () => {
    render(<MarkdownContent content="[Click here](https://example.com)" />);
    const link = screen.getByRole('link', { name: 'Click here' });
    expect(link).toHaveAttribute('href', 'https://example.com');
  });

  it('renders code blocks', () => {
    const md = '```js\nconsole.log("hello");\n```';
    const { container } = render(<MarkdownContent content={md} />);
    const codeEl = container.querySelector('code');
    expect(codeEl).toBeInTheDocument();
    expect(codeEl?.textContent).toContain('console.log("hello");');
  });

  it('renders lists', () => {
    render(<MarkdownContent content={'- Item one\n- Item two'} />);
    expect(screen.getByText('Item one')).toBeInTheDocument();
    expect(screen.getByText('Item two')).toBeInTheDocument();
  });
});
