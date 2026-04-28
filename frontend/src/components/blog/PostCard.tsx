import { Link } from "react-router-dom";
import type { Post } from "../../types/post";
import { TagPill } from "./TagPill";
import { formatRelativeTime } from "../../utils/time";

interface PostCardProps {
  post: Post;
}

function excerpt(content: string, maxLen = 200): string {
  if (content.length <= maxLen) return content;
  return content.slice(0, maxLen).trimEnd() + "…";
}

export function PostCard({ post }: PostCardProps) {
  return (
    <article className="border-[length:var(--border-width)] border-border bg-surface p-[var(--card-padding)]">
      <h2 className="font-heading text-xl font-bold mb-2">
        <Link to={`/blog/${post.slug}`} className="hover:text-accent">
          {post.title}
        </Link>
      </h2>
      <p className="text-text-secondary text-sm mb-3">
        {excerpt(post.content)}
      </p>
      <div className="flex flex-wrap items-center gap-2">
        {post.tags.map((tag) => (
          <TagPill key={tag} tag={tag} />
        ))}
        <time
          role="time"
          dateTime={post.created_at}
          title={new Date(post.created_at).toISOString()}
          className="text-xs text-text-secondary ml-auto"
        >
          {formatRelativeTime(post.created_at)}
        </time>
      </div>
    </article>
  );
}
