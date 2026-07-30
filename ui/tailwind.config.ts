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
      fontFamily: {
        sans: ['var(--font-sans)', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['var(--font-mono)', 'ui-monospace', 'SFMono-Regular', 'monospace'],
      },
    },
  },
  plugins: [],
}

export default config
