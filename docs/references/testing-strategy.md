# Testing Strategy

## Philosophy

**Test-Driven Development (TDD)** is the default workflow:

1. **Red**: Write a failing test that describes the desired behavior
2. **Green**: Write the minimum code to make the test pass
3. **Refactor**: Clean up while keeping tests green

Tests are not an afterthought. Every feature begins with a test. Every bug fix begins with a test that reproduces the bug.

## Test Pyramid

```
         ╱  E2E  ╲           Few, slow, high confidence
        ╱──────────╲
       ╱ Integration ╲       Moderate count, API + DB
      ╱────────────────╲
     ╱   Unit Tests     ╲    Many, fast, isolated
    ╱────────────────────╲
```

| Level | Frontend | Backend | Widgets |
|---|---|---|---|
| **Unit** | Component rendering, hooks, utils | Service logic, model validation | Route handlers, services |
| **Integration** | API service layer mocks, theme provider | Handler + service + repository (in-memory DB) | API endpoint tests |
| **E2E** | Playwright (critical user flows) | - | - |

## Frontend Testing

### Tools

- **Vitest** -- Test runner (Vite-native, fast)
- **React Testing Library** -- Component testing (user-centric)
- **MSW (Mock Service Worker)** -- API mocking
- **Playwright** -- E2E browser tests

### Structure

```
frontend/
├── src/
│   ├── components/
│   │   └── blog/
│   │       ├── PostCard.tsx
│   │       └── PostCard.test.tsx      # Co-located unit test
│   ├── hooks/
│   │   ├── usePosts.ts
│   │   └── usePosts.test.ts
│   └── services/
│       ├── api.ts
│       └── api.test.ts
├── tests/
│   ├── setup.ts                       # Global test setup
│   ├── mocks/                         # MSW handlers
│   │   └── handlers.ts
│   └── e2e/                           # Playwright tests
│       ├── blog-feed.spec.ts
│       └── admin-login.spec.ts
```

### Example: TDD a PostCard Component

```typescript
// 1. RED: Write the test first
// PostCard.test.tsx
import { render, screen } from '@testing-library/react';
import { PostCard } from './PostCard';

const mockPost = {
  slug: 'hello-world',
  title: 'Hello World',
  content: 'This is my first post.',
  tags: ['intro', 'hello'],
  createdAt: '2026-03-15T10:00:00Z',
  commentCount: 3,
  likeCount: 12,
};

describe('PostCard', () => {
  it('renders the post title', () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText('Hello World')).toBeInTheDocument();
  });

  it('renders all tags as pills', () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText('intro')).toBeInTheDocument();
    expect(screen.getByText('hello')).toBeInTheDocument();
  });

  it('shows comment and like counts', () => {
    render(<PostCard post={mockPost} />);
    expect(screen.getByText('3')).toBeInTheDocument();
    expect(screen.getByText('12')).toBeInTheDocument();
  });

  it('displays a relative timestamp', () => {
    render(<PostCard post={mockPost} />);
    // Timestamp element exists with datetime attribute
    const time = screen.getByRole('time');
    expect(time).toHaveAttribute('datetime', '2026-03-15T10:00:00Z');
  });
});

// 2. GREEN: Implement PostCard to pass all tests
// 3. REFACTOR: Extract shared utilities, clean up styles
```

### Component Testing Conventions

- Test **behavior**, not implementation details
- Use `screen.getByRole`, `getByText`, `getByLabelText` -- not `getByTestId`
- Test user interactions via `userEvent` (click, type, etc.)
- Each test file lives next to the component it tests

## Backend Testing (Go)

### Tools

- **Go standard `testing` package** -- Unit and integration tests
- **`httptest`** -- HTTP handler testing
- **`testify`** -- Assertions (optional, for readability)
- **In-memory SQLite** -- Integration tests with real SQL

### Structure

```
backend/
├── internal/
│   ├── service/
│   │   ├── post_service.go
│   │   └── post_service_test.go       # Co-located
│   ├── handler/
│   │   ├── posts.go
│   │   └── posts_test.go
│   └── repository/
│       ├── post_repo.go
│       └── post_repo_test.go
```

### Example: TDD a Post Service

```go
// post_service_test.go

func TestCreatePost_ValidInput(t *testing.T) {
    repo := NewInMemoryPostRepo()
    svc := NewPostService(repo)

    post, err := svc.Create(CreatePostInput{
        Title:   "Hello World",
        Content: "First post content",
        Tags:    []string{"intro"},
    })

    assert.NoError(t, err)
    assert.Equal(t, "hello-world", post.Slug)
    assert.Equal(t, "Hello World", post.Title)
    assert.False(t, post.Published)
}

func TestCreatePost_DuplicateSlug(t *testing.T) {
    repo := NewInMemoryPostRepo()
    svc := NewPostService(repo)

    _, _ = svc.Create(CreatePostInput{Title: "Hello World", Content: "first"})
    _, err := svc.Create(CreatePostInput{Title: "Hello World", Content: "second"})

    assert.ErrorIs(t, err, ErrSlugExists)
}

func TestCreatePost_EmptyTitle(t *testing.T) {
    repo := NewInMemoryPostRepo()
    svc := NewPostService(repo)

    _, err := svc.Create(CreatePostInput{Title: "", Content: "content"})

    assert.ErrorIs(t, err, ErrTitleRequired)
}
```

### Backend Testing Conventions

- **Services**: Test with in-memory repository implementations
- **Handlers**: Test with `httptest.NewRecorder()` and mock services
- **Repositories**: Test against an in-memory SQLite database
- Table-driven tests for input variations
- No mocking frameworks -- use interfaces and simple test doubles

## Widget Service Testing (Python)

### Tools

- **pytest** -- Test runner
- **httpx** -- Async test client for FastAPI
- **pytest-asyncio** -- Async test support

### Example

```python
# tests/test_routes.py
import pytest
from httpx import AsyncClient
from app.main import app

@pytest.mark.asyncio
async def test_health_check():
    async with AsyncClient(app=app, base_url="http://test") as client:
        response = await client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "ok"}
```

## Coverage Targets

| Layer | Minimum Coverage | Focus |
|---|---|---|
| Frontend components | 80% | Rendering, user interactions |
| Frontend services | 90% | API calls, error handling |
| Backend services | 90% | Business logic, edge cases |
| Backend handlers | 80% | Request/response mapping |
| Backend repositories | 80% | Query correctness |
| Widget service | 70% | Route handlers |

Coverage is measured but not enforced via CI gates initially. The goal is high-quality tests, not gaming a number.

## Test Commands

```bash
# Frontend
cd frontend && npm test              # Run all tests
cd frontend && npm test -- --watch   # Watch mode
cd frontend && npm run test:coverage # With coverage report

# Backend
cd backend && go test ./...           # Run all tests
cd backend && go test -v ./...        # Verbose
cd backend && go test -cover ./...    # With coverage

# Widgets
cd widgets && pytest                  # Run all tests
cd widgets && pytest --cov=app        # With coverage

# All
make test                             # Run everything
```
