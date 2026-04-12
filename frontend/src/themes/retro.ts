import type { ThemeTokens } from './tokens';

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
