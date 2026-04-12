# Design: Static Pages & Navigation

## Technical Approach

### Backend

Pages endpoints reuse the existing service/repository pattern:

```go
r.Route("/api/v1", func(r chi.Router) {
    // ... existing routes
    r.Get("/pages", pageHandler.ListPublished)
    r.Get("/pages/{slug}", pageHandler.GetBySlug)
})
```

### Navigation Component

Fetches pages on mount and caches them:

```tsx
function Header() {
  const { pages } = usePages(); // fetches GET /api/v1/pages once

  return (
    <header className="border-b-2 border-border">
      <nav className="max-w-6xl mx-auto px-4 py-4 flex justify-between items-center">
        <Link to="/" className="font-heading text-2xl font-bold">neo</Link>
        <div className="flex gap-6 font-heading text-sm">
          <Link to="/">blog</Link>
          {pages?.map(p => (
            <Link key={p.slug} to={`/page/${p.slug}`}>{p.title.toLowerCase()}</Link>
          ))}
        </div>
      </nav>
    </header>
  );
}
```

### Tag Cloud

```tsx
function TagCloud() {
  const { tags } = useTags();
  const maxCount = Math.max(...tags.map(t => t.count));

  return (
    <div className="border-2 border-border p-4">
      <h3 className="font-heading text-lg mb-3">Tags</h3>
      <div className="flex flex-wrap gap-2">
        {tags.map(t => (
          <Link
            key={t.name}
            to={`/tag/${t.name}`}
            className="text-accent hover:underline"
            style={{ fontSize: `${0.75 + (t.count / maxCount) * 0.5}rem` }}
          >
            {t.name}
          </Link>
        ))}
      </div>
    </div>
  );
}
```

### Routes Update

```tsx
<Routes>
  <Route path="/" element={<Home />} />
  <Route path="/blog/:slug" element={<PostView />} />
  <Route path="/tag/:tag" element={<TagFeed />} />
  <Route path="/page/:slug" element={<PageView />} />
  <Route path="/admin/*" element={<AdminRoutes />} />
  <Route path="*" element={<NotFound />} />
</Routes>
```

### 404 Page

```tsx
function NotFound() {
  return (
    <Layout>
      <div className="text-center py-20">
        <h1 className="font-heading text-6xl mb-4">404</h1>
        <p className="text-text-secondary mb-8">Page not found.</p>
        <Link to="/" className="text-accent hover:underline">← back home</Link>
      </div>
    </Layout>
  );
}
```
