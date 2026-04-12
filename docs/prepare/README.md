# Neo - Implementation Plan

> Personal web garden: blogging, notebooks, web apps, and widget deployments.

## Vision

A fast-loading personal website with a retro aesthetic (inspired by radicle.xyz) that serves as a blog, notebook, and platform for interactive web apps and HuggingFace widget deployments. The blog experience draws from Tumblr's social UX patterns -- card-based posts, reactions, tags, and inline comments.

## Documents

| Document | Description |
|---|---|
| [architecture.md](./architecture.md) | System architecture, tech stack, project structure |
| [design-system.md](./design-system.md) | Visual design, theme system, retro aesthetic |
| [blog-and-admin.md](./blog-and-admin.md) | Blog features, Tumblr-inspired UX, admin panel |
| [deployment.md](./deployment.md) | Deployment strategy, Docker, CI/CD |
| [testing-strategy.md](./testing-strategy.md) | TDD approach, test structure, coverage |
| [milestones.md](./milestones.md) | Phased implementation roadmap |

## Guiding Principles

1. **KISS** -- Simple solutions first. No premature abstractions.
2. **Separation of concerns** -- Frontend, backend API, and widget services are independent.
3. **TDD** -- Tests are written before implementation.
4. **Theme-first design** -- UI built on swappable theme tokens from day one.
5. **Progressive complexity** -- Start with SQLite + basic auth, evolve to PostgreSQL + OAuth2.
