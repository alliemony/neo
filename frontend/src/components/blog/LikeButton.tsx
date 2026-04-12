import { useState } from "react";
import { likePost } from "../../services/api";

interface LikeButtonProps {
  slug: string;
  initialCount: number;
}

export function LikeButton({ slug, initialCount }: LikeButtonProps) {
  const [count, setCount] = useState(initialCount);
  const [liked, setLiked] = useState(false);

  const handleLike = async () => {
    if (liked) return;

    // Optimistic update
    setCount((prev) => prev + 1);
    setLiked(true);

    try {
      const result = await likePost(slug);
      setCount(result.like_count);
    } catch {
      // Revert on error
      setCount((prev) => prev - 1);
      setLiked(false);
    }
  };

  return (
    <button
      onClick={handleLike}
      disabled={liked}
      className={`inline-flex items-center gap-1.5 border-2 border-border px-3 py-1.5 text-sm font-heading transition-colors ${
        liked
          ? "bg-accent text-white border-accent"
          : "bg-surface text-text-primary hover:border-accent"
      }`}
      title={liked ? "Liked!" : "Like this post"}
    >
      <span>{liked ? "♥" : "♡"}</span>
      <span>{count}</span>
    </button>
  );
}
