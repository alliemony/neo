# Neo - Requirements Plan (OpenSpec)

Each directory under `changes/` is a self-contained requirement set following the [OpenSpec](https://openspec.dev/) structure:

```
changes/<name>/
├── proposal.md   # Why, What, Approach
├── spec.md       # Requirements (SHALL/MUST/SHOULD) + Scenarios (GIVEN/WHEN/THEN)
├── design.md     # Technical architecture and approach
└── tasks.md      # Numbered implementation checklist
```

## Requirement Sets

Implementation order follows the numbering. Each set depends on the ones before it.

| # | Requirement Set | Summary | Depends On |
|---|---|---|---|
| **01** | [foundation](./01-foundation/) | Project scaffolding, Vite+React, Go API, FastAPI, Docker, test infra | - |
| **02** | [design-system](./02-design-system/) | Theme tokens, ThemeProvider, Tailwind integration, retro aesthetic, layout shell | 01 |
| **03** | [blog-backend](./03-blog-backend/) | Posts DB schema, repository, service, REST API, tags, seed data | 01 |
| **04** | [blog-frontend](./04-blog-frontend/) | PostCard, PostList, markdown rendering, feed/post/tag routes | 02, 03 |
| **05** | [comments-and-interactions](./05-comments-and-interactions/) | Comments backend+frontend, likes, rate limiting, sidebar comments | 03, 04 |
| **06** | [admin-panel](./06-admin-panel/) | Basic auth, admin dashboard, post/page editor, tag autocomplete | 03, 04 |
| **07** | [pages-and-navigation](./07-pages-and-navigation/) | Static pages, dynamic nav, tag cloud, footer, 404 | 02, 06 |
| **08** | [widget-integration](./08-widget-integration/) | Python widget service, widget registry, iframe embed, widget post type | 01, 04 |
| **09** | [production-deployment](./09-production-deployment/) | Multi-stage Docker, PostgreSQL, static serving, RSS, SEO, platform deploy | All |

## Dependency Graph

```
01-foundation
├── 02-design-system
│   ├── 04-blog-frontend ← 03-blog-backend
│   │   ├── 05-comments-and-interactions
│   │   └── 08-widget-integration
│   └── 07-pages-and-navigation ← 06-admin-panel ← 03, 04
├── 03-blog-backend
└── 08-widget-integration (Python service)

09-production-deployment (after all above)
```

## Parallelization Opportunities

These can be worked on simultaneously:
- **02 + 03**: Design system and blog backend are independent (both depend only on 01)
- **05 + 06 + 08**: Comments, admin, and widgets are independent of each other (all depend on 04)

## Spec Conventions

- **SHALL / MUST**: Mandatory requirement. Implementation is not complete without it.
- **SHOULD**: Strongly recommended. Omit only with documented justification.
- **MAY**: Optional. Nice to have.
- **Scenarios**: Follow BDD format -- GIVEN (context), WHEN (action), THEN (expected outcome).
- **Tasks**: Numbered checklists. Each task is a single commit-sized unit of work. Tests are written before implementation (TDD).
