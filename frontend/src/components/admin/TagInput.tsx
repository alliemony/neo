import { useState } from 'react';

interface TagInputProps {
  value: string[];
  onChange: (tags: string[]) => void;
  existingTags: string[];
}

export function TagInput({ value, onChange, existingTags }: TagInputProps) {
  const [input, setInput] = useState('');

  const suggestions = input.length > 0
    ? existingTags.filter(
        (t) => t.toLowerCase().startsWith(input.toLowerCase()) && !value.includes(t),
      )
    : [];

  function addTag(tag: string) {
    const trimmed = tag.trim().toLowerCase();
    if (trimmed && !value.includes(trimmed)) {
      onChange([...value, trimmed]);
    }
    setInput('');
  }

  function removeTag(tag: string) {
    onChange(value.filter((t) => t !== tag));
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (input.trim()) {
        addTag(input);
      }
    }
  }

  return (
    <div>
      <div className="flex flex-wrap gap-1 mb-1">
        {value.map((tag) => (
          <span
            key={tag}
            className="inline-flex items-center gap-1 bg-accent/20 text-accent text-xs px-2 py-0.5"
          >
            {tag}
            <button type="button" onClick={() => removeTag(tag)} aria-label="×">
              ×
            </button>
          </span>
        ))}
      </div>
      <div className="relative">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="Add tag…"
          className="w-full border-[length:var(--border-width)] border-border bg-surface p-2 text-text text-sm"
        />
        {suggestions.length > 0 && (
          <ul className="absolute z-10 w-full bg-surface border-[length:var(--border-width)] border-border mt-[-1px]">
            {suggestions.map((s) => (
              <li key={s}>
                <button
                  type="button"
                  role="option"
                  onClick={() => addTag(s)}
                  className="w-full text-left px-2 py-1 text-sm hover:bg-accent/10"
                >
                  {s}
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
