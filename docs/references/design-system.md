# Design System

## Aesthetic Direction

**Warm minimal**: Clean, friendly, readable — like a well-tended personal space. System fonts, warm cream tones, soft borders, and colorful tag pills. Single-column 700px layout matching allieg.dev.

### Core Visual Characteristics

- **System UI typography** for both headings and body — no monospace headings
- **Soft 1px borders** in a warm beige tone — not heavy outlines
- **Flat color blocks** — no gradients, no skeuomorphism
- **Single-column layout** at max-width 700px
- **Colorful hash-based tag pills** — each tag consistently maps to a color from an 8-color palette
- **Warm cream background** with parchment surface tones

## Color Palette

### Default "Retro" Theme (allieg.dev warm palette)

```
Background:     #fefcf8  (warm cream)
Surface:        #f4f0e8  (parchment, for cards/surfaces)
Border:         #e2dbd0  (soft warm border)
Border Dark:    #cec6b8  (slightly stronger border)
Text Primary:   #302b26  (warm dark brown)
Text Secondary: #5a5550  (medium warm gray)
Text Dim:       #908880  (muted text)
Accent:         #c04830  (brick red, for links and highlights)
Accent Alt:     #a33c28  (darker brick red, hover state)
Tag Background: #f4f0e8  (same as surface)
Success:        #287848  (green for positive states)
Code Background:#eae5db  (warm light beige for code blocks)
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
Headings:    system-ui, -apple-system, 'Segoe UI', sans-serif
Body:        system-ui, -apple-system, 'Segoe UI', sans-serif
Code:        'Courier Prime', 'Courier New', monospace

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

### Post Card

```
April 10, 2026            ← date in text-dim
Post Title                ← font-semibold, hover → accent color
Excerpt text showing...   ← text-secondary, 180 char truncation

[tag-one] [tag-two]       ← colorful hash-based pill tags
```

- **1px border-top** separator between cards (no box borders)
- **6px border-radius** on pill tags
- Tags use deterministic hash-based colors from 8-color palette

### Tag Color System

Tags use `getTagColor(tag)` from `utils/tagColor.ts` — a hash function maps each tag string to one of 8 color pairs (bg + text). This makes tag colors consistent across all pages.

### Page Layout

```
┌────────────────────────────────────────────┐
│  allieg.dev    home blog widgets about recs │
│  ────────────────────────────────────────  │
│                                            │
│         Main content (max 700px)           │
│         hero / post list / page content    │
│                                            │
│  ──────────────────────────────────────── │
│  allieg.dev © 2026 — made with too much    │
│  coffee                                    │
└────────────────────────────────────────────┘
```

Single-column layout. No sidebar. Comments appear below post content.

## Responsive Strategy

- Single-column at all breakpoints (max-width 700px, centered)
- Mobile: same layout with `px-6` padding

## Accessibility

- All colors meet WCAG AA contrast ratios
- Semantic HTML (`<article>`, `<nav>`, `<main>`, `<aside>`)
- Keyboard navigation for all interactive elements
- Reduced motion preferences respected via `prefers-reduced-motion`
