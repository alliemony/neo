# Design: Design System & Theming

## Technical Approach

### Theme Token Flow

```
ThemeTokens (TS object)
    ↓ ThemeProvider
CSS Custom Properties (on <html>)
    ↓ tailwind.config.ts
Tailwind Utility Classes (theme-aware)
    ↓ Components
Rendered UI
```

### Key Files

```
frontend/src/
├── themes/
│   ├── tokens.ts           # ThemeTokens interface
│   ├── retro.ts            # Default retro theme values
│   ├── dark.ts             # Dark variant values
│   └── ThemeProvider.tsx    # React context + CSS variable injection
├── components/
│   └── layout/
│       ├── Layout.tsx       # Two-column shell with responsive breakpoints
│       ├── Header.tsx       # Site name + nav links
│       ├── Footer.tsx       # Footer text
│       └── Sidebar.tsx      # Right sidebar container
```

### ThemeProvider Implementation

```tsx
// ThemeProvider.tsx
const ThemeContext = createContext<{
  theme: string;
  setTheme: (name: string) => void;
}>();

const themes: Record<string, ThemeTokens> = { retro, dark };

function ThemeProvider({ children, defaultTheme = 'retro' }) {
  const [theme, setTheme] = useState(defaultTheme);

  useEffect(() => {
    const tokens = themes[theme];
    const root = document.documentElement;
    // Map each token to a CSS variable
    root.style.setProperty('--color-bg', tokens.colors.bg);
    root.style.setProperty('--color-surface', tokens.colors.surface);
    // ... all tokens
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}
```

### Layout Component

```tsx
// Layout.tsx
function Layout({ children, sidebar }: { children: ReactNode; sidebar?: ReactNode }) {
  return (
    <div className="min-h-screen bg-bg text-text-primary font-body">
      <Header />
      <div className="max-w-6xl mx-auto px-4 py-8 lg:flex lg:gap-8">
        <main className="flex-1">{children}</main>
        {sidebar && (
          <aside className="w-full lg:w-80 mt-8 lg:mt-0">{sidebar}</aside>
        )}
      </div>
      <Footer />
    </div>
  );
}
```

### Font Loading

Use `@fontsource` npm packages for self-hosted fonts (no external CDN dependency):

```bash
npm install @fontsource/jetbrains-mono @fontsource/inter
```

Import in `main.tsx`:
```tsx
import '@fontsource/jetbrains-mono/400.css';
import '@fontsource/jetbrains-mono/700.css';
import '@fontsource/inter/400.css';
import '@fontsource/inter/600.css';
```

### Retro Aesthetic Details

- `border-radius: 0px` globally (sharp corners)
- `border-width: 2px` for all bordered elements
- `box-shadow: none` (flat, no elevation)
- Monospace headings create the terminal/technical feel
- High line-height (1.7) for body text readability
