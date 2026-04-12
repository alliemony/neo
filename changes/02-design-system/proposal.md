# Proposal: Design System & Theming

## Why

The site's retro aesthetic is a core part of its identity. Building the theme system early ensures every component built afterward automatically inherits the visual language. A swappable theme architecture also makes it trivial to add dark mode or entirely new themes later.

## What

Create a design token system, a ThemeProvider React context, Tailwind CSS integration via CSS custom properties, and the base layout components (Header, Footer, Sidebar) styled with the default "retro terminal" theme.

## What Changes

- Define `ThemeTokens` TypeScript interface (colors, fonts, spacing, effects)
- Implement the default retro theme and a dark variant
- Create `ThemeProvider` React context that injects tokens as CSS custom properties
- Configure `tailwind.config.ts` to consume CSS variables so utility classes are theme-aware
- Load web fonts (JetBrains Mono, Inter)
- Build base layout components: `Header`, `Footer`, `Sidebar`, `Layout` wrapper
- Responsive shell: two-column (desktop) → single-column (mobile)

## Approach

Tokens are plain TypeScript objects. The ThemeProvider maps them to CSS variables on `<html>`. Tailwind reads those variables in `theme.extend`. This means `className="text-accent border-border bg-surface"` automatically adapts to the active theme with zero prop threading.
