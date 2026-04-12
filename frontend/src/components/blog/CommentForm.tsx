import { useState } from "react";
import type { CreateCommentInput } from "../../types/post";

interface CommentFormProps {
  onSubmit: (input: CreateCommentInput) => Promise<void>;
}

export function CommentForm({ onSubmit }: CommentFormProps) {
  const [authorName, setAuthorName] = useState("");
  const [content, setContent] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!authorName.trim() || !content.trim()) {
      setError("Name and comment are required.");
      return;
    }

    setSubmitting(true);
    setError(null);

    try {
      await onSubmit({
        author_name: authorName.trim(),
        content: content.trim(),
      });
      setAuthorName("");
      setContent("");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to submit comment.",
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <div>
        <input
          type="text"
          placeholder="Your name"
          value={authorName}
          onChange={(e) => setAuthorName(e.target.value)}
          className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary placeholder:text-text-secondary focus:outline-none focus:border-accent"
          maxLength={100}
        />
      </div>
      <div>
        <textarea
          placeholder="Write a comment…"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          rows={3}
          className="w-full border-2 border-border bg-surface px-3 py-2 text-sm font-body text-text-primary placeholder:text-text-secondary focus:outline-none focus:border-accent resize-y"
          maxLength={2000}
        />
      </div>
      {error && <p className="text-accent text-xs">{error}</p>}
      <button
        type="submit"
        disabled={submitting}
        className="border-2 border-border bg-accent text-white px-4 py-2 text-sm font-heading hover:opacity-90 disabled:opacity-50"
      >
        {submitting ? "Sending…" : "Comment"}
      </button>
    </form>
  );
}
