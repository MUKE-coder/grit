/**
 * The single source of truth for the Grit UI palette.
 *
 * Used in three places that must never disagree:
 *   1. app/globals.css — so previews on this site render correctly
 *   2. every registry item's `cssVars` — so a component dropped into someone
 *      else's project renders correctly too
 *   3. the install docs
 *
 * Keeping one exported object means adding a colour cannot leave the registry
 * shipping a palette the site does not actually use.
 */
export const GRIT_TOKENS: Record<string, string> = {
  '--bg-primary': '#0a0a0f',
  '--bg-secondary': '#111118',
  '--bg-tertiary': '#1a1a24',
  '--bg-elevated': '#22222e',
  '--bg-hover': '#2a2a38',
  '--border': '#2a2a3a',
  '--text-primary': '#e8e8f0',
  '--text-secondary': '#9090a8',
  '--text-muted': '#606078',
  '--accent': '#6c5ce7',
  '--accent-hover': '#7c6cf7',
  '--accent-foreground': '#ffffff',
  '--success': '#00b894',
  '--danger': '#ff6b6b',
  '--warning': '#fdcb6e',
  '--info': '#74b9ff',
}

/** The Tailwind colour scale a consumer needs so `text-text-muted` and
 *  `bg-bg-elevated` resolve. Shipped in the registry item's `tailwind` block. */
export const GRIT_TAILWIND_COLORS = {
  background: 'var(--bg-primary)',
  foreground: 'var(--text-primary)',
  border: 'var(--border)',
  bg: {
    primary: 'var(--bg-primary)',
    secondary: 'var(--bg-secondary)',
    tertiary: 'var(--bg-tertiary)',
    elevated: 'var(--bg-elevated)',
    hover: 'var(--bg-hover)',
  },
  text: {
    primary: 'var(--text-primary)',
    secondary: 'var(--text-secondary)',
    muted: 'var(--text-muted)',
  },
  accent: {
    DEFAULT: 'var(--accent)',
    hover: 'var(--accent-hover)',
    foreground: 'var(--accent-foreground)',
  },
  success: 'var(--success)',
  danger: 'var(--danger)',
  warning: 'var(--warning)',
  info: 'var(--info)',
}

/** Emits the token block as CSS, for globals.css and the docs snippet. */
export function tokensAsCSS(indent = '  '): string {
  return Object.entries(GRIT_TOKENS)
    .map(([k, v]) => `${indent}${k}: ${v};`)
    .join('\n')
}
