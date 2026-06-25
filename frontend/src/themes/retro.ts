import type { ThemeTokens } from './tokens';

export const retroTheme: ThemeTokens = {
  colors: {
    bg: '#fefcf8',
    surface: '#f4f0e8',
    border: '#e2dbd0',
    textPrimary: '#302b26',
    textSecondary: '#5a5550',
    accent: '#c04830',
    accentAlt: '#a33c28',
    tagBg: '#f4f0e8',
    success: '#287848',
    codeBg: '#eae5db',
  },
  fonts: {
    heading: "system-ui, -apple-system, 'Segoe UI', sans-serif",
    body: "system-ui, -apple-system, 'Segoe UI', sans-serif",
    code: "'Courier Prime', 'Courier New', monospace",
  },
  spacing: {
    borderWidth: '1px',
    borderRadius: '6px',
    cardPadding: '1.25rem',
  },
  effects: {
    shadow: '0 2px 12px rgba(30,26,22,0.06)',
  },
};
