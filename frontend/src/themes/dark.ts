import type { ThemeTokens } from './tokens';

export const darkTheme: ThemeTokens = {
  colors: {
    bg: '#1A1A1A',
    surface: '#242424',
    border: '#404040',
    textPrimary: '#E8E6E3',
    textSecondary: '#999999',
    accent: '#FF7B5C',
    accentAlt: '#5C9BFF',
    tagBg: '#2D2A26',
    success: '#3DBA6A',
    codeBg: '#2A2725',
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
