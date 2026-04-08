import { createContext, useContext, useEffect, useState, type ReactNode } from 'react';
import type { ThemeTokens } from './tokens';
import { themeToCSSVars } from './tokens';
import { retroTheme } from './retro';
import { darkTheme } from './dark';

const themes: Record<string, ThemeTokens> = {
  retro: retroTheme,
  dark: darkTheme,
};

interface ThemeContextValue {
  theme: string;
  setTheme: (name: string) => void;
  availableThemes: string[];
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

interface ThemeProviderProps {
  children: ReactNode;
  defaultTheme?: string;
}

export function ThemeProvider({ children, defaultTheme = 'retro' }: ThemeProviderProps) {
  const [theme, setTheme] = useState(defaultTheme);

  useEffect(() => {
    const tokens = themes[theme];
    if (!tokens) return;

    const vars = themeToCSSVars(tokens);
    const root = document.documentElement;
    for (const [key, value] of Object.entries(vars)) {
      root.style.setProperty(key, value);
    }
  }, [theme]);

  const value: ThemeContextValue = {
    theme,
    setTheme,
    availableThemes: Object.keys(themes),
  };

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return ctx;
}

export { themes };
