import type { Config } from 'tailwindcss'

/**
 * The Grit UI palette.
 *
 * The component sources use two families of class: shadcn-standard ones like
 * `bg-background` and `text-foreground`, and Grit-specific ones like
 * `text-text-muted`, `bg-bg-elevated` and `bg-accent`. The nested `text` and
 * `bg` scales below are what make the second family resolve — `text-text-muted`
 * is the `text-` utility applied to the `text.muted` colour.
 *
 * Every value is a CSS variable rather than a literal so the exact same palette
 * can be shipped to a consumer project through the registry item's `cssVars`
 * block. A component that silently loses its colours in someone else's app is
 * worse than one that fails to compile, because it looks like a design choice.
 */
const config: Config = {
  darkMode: 'class',
  content: [
    './app/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './registry/**/*.{ts,tsx}',
  ],
  theme: {
    extend: {
      colors: {
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
      },
      fontFamily: {
        sans: ['Onest', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['JetBrains Mono', 'ui-monospace', 'SFMono-Regular', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
