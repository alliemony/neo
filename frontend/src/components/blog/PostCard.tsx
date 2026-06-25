import { Link } from "react-router-dom";
import type { Post } from "../../types/post";
import { TagPill } from "./TagPill";

interface PostCardProps {
  post: Post;
}

function formatDate(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleDateString("en-US", { year: "numeric", month: "long", day: "numeric" });
}

function truncate(text: string, max = 180): string {
  const plain = text.replace(/<[^>]+>/g, "");
  if (plain.length <= max) return plain;
  return plain.slice(0, max).trimEnd() + "…";
}

export function PostCard({ post }: PostCardProps) {
  return (
    <article className="py-4 border-t border-border first:border-t-0">
      <Link to={`/blog/${post.slug}`} className="block no-underline group">
        <div className="text-[14px] text-text-secondary mb-1">{formatDate(post.created_at)}</div>
        <h2 className="font-semibold text-[20px] text-text-primary mb-1.5 leading-snug transition-colors duration-150 group-hover:text-accent">
          {post.title}
        </h2>
        <p className="text-[16px] text-text-secondary leading-relaxed mb-2.5">
          {truncate(post.content)}
        </p>
      </Link>
      {post.tags.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {post.tags.map((tag) => (
            <TagPill key={tag} tag={tag} />
          ))}
        </div>
      )}
    </article>
  );
}
