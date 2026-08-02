import type { Config } from 'tailwindcss'

/**
 * Deliberately close to stock Tailwind.
 *
 * Blocks are authored with default palette classes (bg-white, text-gray-900,
 * bg-indigo-600) so a copied block renders correctly in any Tailwind project
 * with no config to merge and no CSS variables to install. Adding a custom
 * colour scale here would quietly make that untrue — the block would look right
 * on this site and wrong everywhere else.
 *
 * Only fonts are extended, and only for the site chrome.
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
      /* Admin design tokens, mirrored here ONLY so the swappable slot previews
         (Application UI → Elements) render the way they will in a real admin.
         Those blocks are authored against the admin's token names — bg-accent,
         bg-bg-tertiary, text-foreground — because a swapped button has to follow
         whatever theme the project is running. Without these the previews would
         compile fine and render as unstyled boxes, which is the failure mode
         that looks like nothing is wrong. */
      colors: {
        accent: { DEFAULT: 'var(--accent)', hover: 'var(--accent-hover)' },
        foreground: 'var(--text-primary)',
        'text-secondary': 'var(--text-secondary)',
        'text-muted': 'var(--text-muted)',
        'bg-primary': 'var(--bg-primary)',
        'bg-secondary': 'var(--bg-secondary)',
        'bg-tertiary': 'var(--bg-tertiary)',
        'bg-elevated': 'var(--bg-elevated)',
        'bg-hover': 'var(--bg-hover)',
        border: 'var(--border)',
        danger: 'var(--danger)',
        success: 'var(--success)',
        warning: 'var(--warning)',
      },
      fontFamily: {
        sans: ['var(--font-sans)', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['var(--font-mono)', 'ui-monospace', 'SFMono-Regular', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
