export interface ThemeTokens {
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
    borderRadius: string;
    cardPadding: string;
  };
  effects: {
    shadow: string;
  };
}

export function themeToCSSVars(tokens: ThemeTokens): Record<string, string> {
  return {
    '--color-bg': tokens.colors.bg,
    '--color-surface': tokens.colors.surface,
    '--color-border': tokens.colors.border,
    '--color-text-primary': tokens.colors.textPrimary,
    '--color-text-secondary': tokens.colors.textSecondary,
    '--color-accent': tokens.colors.accent,
    '--color-accent-alt': tokens.colors.accentAlt,
    '--color-tag-bg': tokens.colors.tagBg,
    '--color-success': tokens.colors.success,
    '--color-code-bg': tokens.colors.codeBg,
    '--font-heading': tokens.fonts.heading,
    '--font-body': tokens.fonts.body,
    '--font-code': tokens.fonts.code,
    '--border-width': tokens.spacing.borderWidth,
    '--border-radius': tokens.spacing.borderRadius,
    '--card-padding': tokens.spacing.cardPadding,
    '--shadow': tokens.effects.shadow,
  };
}
