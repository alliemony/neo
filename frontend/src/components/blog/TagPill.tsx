import { Link } from 'react-router-dom';

interface TagPillProps {
  tag: string;
  active?: boolean;
}

export function TagPill({ tag, active }: TagPillProps) {
  return (
    <Link
      to={`/tag/${tag}`}
      className={`inline-block px-2 py-0.5 text-xs font-body border-[length:var(--border-width)] border-border ${
        active
          ? 'bg-accent text-white'
          : 'bg-tag-bg text-text-secondary hover:bg-accent hover:text-white'
      } transition-colors`}
    >
      {tag}
    </Link>
  );
}
