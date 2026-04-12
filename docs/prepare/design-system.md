# Design System

## Aesthetic Direction

**Retro-modern**: Inspired by radicle.xyz's technical minimalism blended with warmth and approachability. The site should feel like a well-organized personal workshop -- clean, purposeful, with personality in the details.

### Core Visual Characteristics

- **Monospace typography** for headers and code, paired with a clean sans-serif for body text
- **Stark borders and box outlines** instead of soft shadows (CSS `border`, not `box-shadow`)
- **Flat color blocks** -- no gradients, no skeuomorphism
- **Grid-based layouts** with generous whitespace
- **High contrast** for readability
- **Subtle color accents** for warmth (a retro palette with muted tones + one vibrant accent)

## Color Palette

### Default "Retro Terminal" Theme

```
Background:     #FAFAF8  (warm off-white)
Surface:        #FFFFFF  (cards, panels)
Border:         #2D2D2D  (stark, visible borders)
Text Primary:   #1A1A1A  (near-black)
Text Secondary: #6B6B6B  (muted gray)
Accent:         #E85D3A  (warm orange-red, for links and highlights)
Accent Alt:     #3A7CE8  (cool blue, for secondary actions)
Tag Background: #F0EDE6  (warm beige for tag pills)
Success:        #2D8A4E  (green for positive states)
Code Background:#F5F2EB  (warm light beige for code blocks)
```

### Dark Variant

```
Background:     #1A1A1A
Surface:        #242424
Border:         #404040
Text Primary:   #E8E6E3
Text Secondary: #999999
Accent:         #FF7B5C
Accent Alt:     #5C9BFF
Tag Background: #2D2A26
Code Background:#2A2725
```

## Typography

```
Headings:    "JetBrains Mono", "IBM Plex Mono", monospace
Body:        "Inter", "IBM Plex Sans", system-ui, sans-serif
Code:        "JetBrains Mono", "Fira Code", monospace

Scale:
  xs:    0.75rem   (12px)
  sm:    0.875rem  (14px)
  base:  1rem      (16px)
  lg:    1.125rem  (18px)
  xl:    1.25rem   (20px)
  2xl:   1.5rem    (24px)
  3xl:   2rem      (32px)
  4xl:   2.5rem    (40px)

Line height:
  body:  1.7 (generous for readability)
  heading: 1.2
```

## Theme System Architecture

Themes are implemented as **CSS custom properties** (variables) controlled by a React context. This allows runtime theme switching without rebuilding.

### Theme Token Structure

```typescript
// themes/tokens.ts
interface ThemeTokens {
  colors: {
    bg: string;
    surface: string;
    border: string;
    textPrimary: string;
    textSecondary: string;
    accent: string;
    accentAlt: string;
    tagBg: string;
    success: string;
    codeBg: string;
  };
  fonts: {
    heading: string;
    body: string;
    code: string;
  };
  spacing: {
    borderWidth: string;
    borderRadius: string;  // "0px" for sharp retro, "4px" for softer
    cardPadding: string;
  };
  effects: {
    shadow: string;        // "none" for flat retro, or subtle values
  };
}
```

### Theme Application

```typescript
// themes/retro.ts
export const retroTheme: ThemeTokens = {
  colors: {
    bg: '#FAFAF8',
    surface: '#FFFFFF',
    border: '#2D2D2D',
    textPrimary: '#1A1A1A',
    textSecondary: '#6B6B6B',
    accent: '#E85D3A',
    accentAlt: '#3A7CE8',
    tagBg: '#F0EDE6',
    success: '#2D8A4E',
    codeBg: '#F5F2EB',
  },
  fonts: {
    heading: '"JetBrains Mono", monospace',
    body: '"Inter", system-ui, sans-serif',
    code: '"JetBrains Mono", monospace',
  },
  spacing: {
    borderWidth: '2px',
    borderRadius: '0px',
    cardPadding: '1.5rem',
  },
  effects: {
    shadow: 'none',
  },
};
```

### Tailwind Integration

Tokens map to Tailwind's `theme.extend` in `tailwind.config.ts`, referencing CSS custom properties. This means Tailwind utility classes automatically use the active theme.

```typescript
// tailwind.config.ts
export default {
  theme: {
    extend: {
      colors: {
        bg: 'var(--color-bg)',
        surface: 'var(--color-surface)',
        border: 'var(--color-border)',
        'text-primary': 'var(--color-text-primary)',
        'text-secondary': 'var(--color-text-secondary)',
        accent: 'var(--color-accent)',
        'accent-alt': 'var(--color-accent-alt)',
        'tag-bg': 'var(--color-tag-bg)',
        'code-bg': 'var(--color-code-bg)',
      },
      fontFamily: {
        heading: 'var(--font-heading)',
        body: 'var(--font-body)',
        code: 'var(--font-code)',
      },
      borderWidth: {
        DEFAULT: 'var(--border-width)',
      },
      borderRadius: {
        DEFAULT: 'var(--border-radius)',
      },
    },
  },
};
```

### Swapping Themes

Adding a new theme:

1. Create a new file in `themes/` (e.g., `cyberpunk.ts`)
2. Export a `ThemeTokens` object with the new values
3. Register it in the theme provider
4. It's immediately selectable -- no other code changes needed

```typescript
// themes/cyberpunk.ts
export const cyberpunkTheme: ThemeTokens = {
  colors: {
    bg: '#0D0D1A',
    surface: '#1A1A2E',
    border: '#00FF88',
    textPrimary: '#E0E0E0',
    textSecondary: '#888899',
    accent: '#FF2E88',
    accentAlt: '#00FF88',
    tagBg: '#1A1A2E',
    success: '#00FF88',
    codeBg: '#16162B',
  },
  // ...
};
```

## Component Design Patterns

### Post Card (Retro Style)

```
┌──────────────────────────────────────────┐
│ ● author · 2 hours ago                  │
├──────────────────────────────────────────┤
│                                          │
│  Post Title in Monospace                 │
│                                          │
│  Body text in a clean readable sans-     │
│  serif font with generous line height    │
│  and comfortable measure width...        │
│                                          │
│  [tag-one] [tag-two] [tag-three]         │
│                                          │
├──────────────────────────────────────────┤
│  ♡ 12    ↩ 3 comments    ⊕ share        │
└──────────────────────────────────────────┘
```

- **2px solid borders** (not shadows)
- **0px border-radius** (sharp corners)
- Tags as bordered pills with muted backgrounds
- Interaction bar at the bottom with clear icon + text

### Page Layout

```
┌─────────────────────────────────────────────────────────┐
│  neo                                    [blog] [about]  │
│  ─────────────────────────────────────────────────────  │
├───────────────────────────────────┬─────────────────────┤
│                                   │                     │
│  Main content area                │  Sidebar            │
│  (post list / single post)        │  - Tags cloud       │
│                                   │  - Recent posts     │
│                                   │  - Comment section   │
│                                   │    (on post view)   │
│                                   │                     │
├───────────────────────────────────┴─────────────────────┤
│  footer · built with care · 2026                        │
└─────────────────────────────────────────────────────────┘
```

## Responsive Strategy

- **Desktop (>1024px)**: Two-column layout (content + sidebar)
- **Tablet (768-1024px)**: Single column, sidebar below content
- **Mobile (<768px)**: Full-width cards, collapsible sidebar

Mobile-first Tailwind classes: default styles target mobile, `md:` and `lg:` prefixes add complexity.

## Accessibility

- All colors meet WCAG AA contrast ratios
- Semantic HTML (`<article>`, `<nav>`, `<main>`, `<aside>`)
- Keyboard navigation for all interactive elements
- Reduced motion preferences respected via `prefers-reduced-motion`
