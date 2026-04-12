import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { getTags } from "../../services/api";
import type { TagCount } from "../../types/post";

export function TagCloud() {
  const [tags, setTags] = useState<TagCount[]>([]);

  useEffect(() => {
    getTags()
      .then(setTags)
      .catch(() => setTags([]));
  }, []);

  if (tags.length === 0) return null;

  const maxCount = Math.max(...tags.map((t) => t.count));

  return (
    <div className="border-2 border-border bg-surface p-4">
      <h3 className="font-heading text-sm font-bold mb-3 uppercase tracking-wider">
        Tags
      </h3>
      <div className="flex flex-wrap gap-2">
        {tags.map((tag) => {
          // Scale font size between 0.75rem and 1.25rem based on count
          const scale = maxCount > 1 ? (tag.count - 1) / (maxCount - 1) : 0;
          const fontSize = 0.75 + scale * 0.5;
          const fontWeight = scale > 0.5 ? 700 : 400;

          return (
            <Link
              key={tag.name}
              to={`/tag/${tag.name}`}
              style={{ fontSize: `${fontSize}rem`, fontWeight }}
              className="text-text-secondary hover:text-accent transition-colors"
            >
              #{tag.name}
            </Link>
          );
        })}
      </div>
    </div>
  );
}
