# Spec: Design System & Theming

## Purpose

Establish a swappable theme system with CSS custom properties, Tailwind integration, and base layout components styled with a retro aesthetic inspired by radicle.xyz.

## Requirements

### Requirement: Theme tokens SHALL be defined as a TypeScript interface covering colors, fonts, spacing, and effects

Every visual property that varies between themes must be captured in the `ThemeTokens` type.

#### Scenario: Token interface includes all required categories

- **GIVEN:** The `ThemeTokens` interface is defined in `themes/tokens.ts`
- **WHEN:** A developer creates a new theme object
- **THEN:** TypeScript enforces that all required properties (bg, surface, border, textPrimary, textSecondary, accent, accentAlt, tagBg, codeBg, heading font, body font, code font, borderWidth, borderRadius) are provided

### Requirement: A default retro theme SHALL be provided with monospace headings, stark borders, and a warm off-white palette

#### Scenario: Retro theme applies correct visual properties

- **GIVEN:** The retro theme is active
- **WHEN:** A component renders with `className="bg-bg text-text-primary border"`
- **THEN:** The background is warm off-white (#FAFAF8), text is near-black (#1A1A1A), and borders are 2px solid dark (#2D2D2D)

### Requirement: A dark theme variant SHOULD be provided alongside the default retro theme

#### Scenario: Dark theme inverts the palette appropriately

- **GIVEN:** The dark theme is active
- **WHEN:** A component renders with `className="bg-bg text-text-primary"`
- **THEN:** The background is dark (#1A1A1A) and text is light (#E8E6E3)

### Requirement: ThemeProvider SHALL inject theme tokens as CSS custom properties on the document root

#### Scenario: CSS variables are set when theme changes

- **GIVEN:** The app is wrapped in `<ThemeProvider>`
- **WHEN:** The active theme changes from retro to dark
- **THEN:** CSS variable `--color-bg` on `<html>` updates from `#FAFAF8` to `#1A1A1A`
- **AND:** All components using Tailwind theme classes re-render with new colors

### Requirement: Tailwind config SHALL reference CSS custom properties so utility classes are theme-aware

#### Scenario: Tailwind utility classes respond to theme tokens

- **GIVEN:** `tailwind.config.ts` maps `colors.accent` to `var(--color-accent)`
- **WHEN:** A component uses `className="text-accent"`
- **THEN:** The rendered text color matches the active theme's accent value

### Requirement: Web fonts SHALL be loaded for heading (monospace) and body (sans-serif) typography

#### Scenario: Fonts load and apply correctly

- **GIVEN:** JetBrains Mono and Inter fonts are available (via CDN or local files)
- **WHEN:** The page renders
- **THEN:** Headings use JetBrains Mono and body text uses Inter
- **AND:** Fallback system fonts are specified in case loading fails

### Requirement: Base layout components SHALL provide a responsive page shell

The layout must include Header, Footer, Sidebar, and a Layout wrapper.

#### Scenario: Desktop layout renders two columns

- **GIVEN:** The viewport width is >= 1024px
- **WHEN:** The Layout component renders with children and sidebar content
- **THEN:** Content appears in the left column and sidebar in the right column

#### Scenario: Mobile layout collapses to single column

- **GIVEN:** The viewport width is < 768px
- **WHEN:** The Layout component renders
- **THEN:** Content and sidebar stack vertically

### Requirement: Header SHALL display the site name and navigation links

#### Scenario: Header renders site title and nav

- **GIVEN:** The Header component is rendered
- **WHEN:** The page loads
- **THEN:** The site name "neo" is visible
- **AND:** Navigation links are displayed

### Requirement: Adding a new theme SHALL require only creating a new token file and registering it

#### Scenario: New theme is added without modifying existing components

- **GIVEN:** A developer creates `themes/cyberpunk.ts` exporting a `ThemeTokens` object
- **WHEN:** They register it in the ThemeProvider's theme map
- **THEN:** The cyberpunk theme is selectable and all components render correctly with it
- **AND:** No component code was modified
