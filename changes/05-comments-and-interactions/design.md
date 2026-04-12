# Design: Comments & Interactions

## Technical Approach

### Backend

#### Rate Limiter

In-memory token bucket using a sync.Map for IP tracking:

```go
// middleware/ratelimit.go
type RateLimiter struct {
    requests sync.Map // IP -> []time.Time
    window   time.Duration
    limit    int
}

func NewRateLimiter(window time.Duration, limit int) *RateLimiter {
    return &RateLimiter{window: window, limit: limit}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        if !rl.allow(ip) {
            http.Error(w, `{"error":"too many requests"}`, http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Applied only to `POST /api/v1/posts/{slug}/comments`.

#### Like Implementation

Simple atomic increment in the repository:

```go
func (r *PostRepo) IncrementLikeCount(slug string) (int, error) {
    result, err := r.db.Exec(
        "UPDATE posts SET like_count = like_count + 1 WHERE slug = ?", slug)
    // return new count via separate query
}
```

### Frontend

#### CommentSection Component

```tsx
function CommentSection({ slug }: { slug: string }) {
  const { comments, loading } = useComments(slug);
  const [submitting, setSubmitting] = useState(false);

  return (
    <div className="border-2 border-border p-4">
      <h3 className="font-heading text-lg mb-4">
        Comments ({comments?.length ?? 0})
      </h3>
      {comments?.map(c => <CommentItem key={c.id} comment={c} />)}
      <CommentForm slug={slug} onSubmitted={refetch} />
    </div>
  );
}
```

#### Like Button with Optimistic Update

```tsx
function LikeButton({ slug, initialCount }: { slug: string; initialCount: number }) {
  const [count, setCount] = useState(initialCount);
  const [liked, setLiked] = useState(false);

  const handleLike = async () => {
    setCount(c => c + 1);  // optimistic
    setLiked(true);
    try {
      const result = await api.likePost(slug);
      setCount(result.like_count);  // reconcile with server
    } catch {
      setCount(c => c - 1);  // revert on failure
      setLiked(false);
    }
  };

  return (
    <button onClick={handleLike} disabled={liked}>
      {liked ? '♥' : '♡'} {count}
    </button>
  );
}
```

### Integration in Post View

```tsx
// routes/PostView.tsx
function PostView() {
  const { slug } = useParams();
  const { post } = usePost(slug);

  return (
    <Layout sidebar={<CommentSection slug={slug} />}>
      <article>
        <h1 className="font-heading text-3xl">{post.title}</h1>
        <MarkdownContent content={post.content} />
        <div className="flex gap-4 mt-4">
          <LikeButton slug={slug} initialCount={post.likeCount} />
          <span>↩ {post.commentCount} comments</span>
        </div>
      </article>
    </Layout>
  );
}
```
