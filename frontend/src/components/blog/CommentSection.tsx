import type { Comment } from "../../types/post";
import { useComments } from "../../hooks/useComments";
import { CommentForm } from "./CommentForm";
import { formatRelativeTime } from "../../utils/time";

interface CommentSectionProps {
  slug: string;
}

export function CommentSection({ slug }: CommentSectionProps) {
  const { comments, loading, error, addComment } = useComments(slug);

  return (
    <div className="border-2 border-border bg-surface p-4">
      <h3 className="font-heading text-lg font-bold mb-4">
        Comments {comments.length > 0 && `(${comments.length})`}
      </h3>

      {loading && (
        <p className="text-text-secondary text-sm">Loading comments…</p>
      )}
      {error && <p className="text-accent text-sm">Failed to load comments.</p>}

      {!loading && !error && (
        <>
          {comments.length === 0 ? (
            <p className="text-text-secondary text-sm mb-4">
              No comments yet. Be the first!
            </p>
          ) : (
            <div className="space-y-4 mb-6">
              {comments.map((comment: Comment) => (
                <div
                  key={comment.id}
                  className="border-b border-border pb-3 last:border-0"
                >
                  <div className="flex justify-between items-baseline mb-1">
                    <span className="font-heading text-sm font-bold">
                      {comment.author_name}
                    </span>
                    <time
                      dateTime={comment.created_at}
                      className="text-xs text-text-secondary"
                    >
                      {formatRelativeTime(comment.created_at)}
                    </time>
                  </div>
                  <p className="text-sm text-text-primary">{comment.content}</p>
                </div>
              ))}
            </div>
          )}

          <CommentForm onSubmit={addComment} />
        </>
      )}
    </div>
  );
}
