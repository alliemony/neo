import { Link } from 'react-router-dom';
import { getTagColor } from '../../utils/tagColor';

interface TagPillProps {
  tag: string;
  active?: boolean;
  onClick?: () => void;
  asSpan?: boolean;
}

export function TagPill({ tag, active, onClick, asSpan }: TagPillProps) {
  const color = getTagColor(tag);
  const style = active
    ? { background: color.text, color: '#fff' }
    : { background: color.bg, color: color.text };

  const className =
    "inline-flex items-center text-[13px] font-medium px-2.5 py-0.5 rounded-full no-underline transition-all duration-150 cursor-pointer select-none leading-relaxed hover:brightness-95";

  if (asSpan || onClick) {
    return (
      <span className={className} style={style} onClick={onClick}>
        {tag}
      </span>
    );
  }

  return (
    <Link
      to={`/blog?tag=${encodeURIComponent(tag)}`}
      className={className}
      style={style}
    >
      {tag}
    </Link>
  );
}
