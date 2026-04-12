# Design: Blog Frontend (Feed & Post Views)

## Technical Approach

### API Client

```typescript
// services/api.ts
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const api = {
  async getPosts(opts?: { tag?: string; page?: number }) {
    const params = new URLSearchParams();
    if (opts?.tag) params.set('tag', opts.tag);
    if (opts?.page) params.set('page', String(opts.page));
    const res = await fetch(`${API_URL}/api/v1/posts?${params}`);
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json() as Promise<{ posts: Post[]; total: number }>;
  },

  async getPost(slug: string) {
    const res = await fetch(`${API_URL}/api/v1/posts/${slug}`);
    if (res.status === 404) return null;
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json() as Promise<Post>;
  },

  async getTags() {
    const res = await fetch(`${API_URL}/api/v1/tags`);
    if (!res.ok) throw new Error(`API error: ${res.status}`);
    return res.json() as Promise<TagCount[]>;
  },
};
```

### Custom Hooks

```typescript
// hooks/usePosts.ts
function usePosts(opts?: { tag?: string; page?: number }) {
  const [data, setData] = useState<{ posts: Post[]; total: number } | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    api.getPosts(opts).then(setData).catch(setError).finally(() => setLoading(false));
  }, [opts?.tag, opts?.page]);

  return { data, loading, error };
}
```

### Markdown Rendering

Use `react-markdown` with `remark-gfm` (GitHub Flavored Markdown) and `rehype-highlight` for syntax highlighting.

```tsx
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';

function MarkdownContent({ content }: { content: string }) {
  return (
    <ReactMarkdown
      remarkPlugins={[remarkGfm]}
      rehypePlugins={[rehypeHighlight]}
      className="prose prose-retro"
    >
      {content}
    </ReactMarkdown>
  );
}
```

### Relative Timestamps

Simple utility -- no need for a full date library:

```typescript
function formatRelativeTime(dateStr: string): string {
  const seconds = Math.floor((Date.now() - new Date(dateStr).getTime()) / 1000);
  if (seconds < 60) return 'just now';
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`;
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
  if (seconds < 604800) return `${Math.floor(seconds / 86400)} days ago`;
  return new Date(dateStr).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
}
```

### Routes

```tsx
// App.tsx
<Routes>
  <Route path="/" element={<Home />} />
  <Route path="/blog/:slug" element={<PostView />} />
  <Route path="/tag/:tag" element={<TagFeed />} />
  <Route path="*" element={<NotFound />} />
</Routes>
```

### Testing Strategy

- **PostCard**: Render with mock data, verify all elements present
- **API client**: Mock fetch with MSW, test success and error paths
- **usePosts hook**: Test with renderHook + MSW
- **Timestamp utility**: Unit test with fixed dates
