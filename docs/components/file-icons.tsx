// VSCode-Material-style file-type icons, rendered as inline SVG (no dependency).
// Brand-coloured letter badges for most types + a real React atom for .tsx/.jsx.
// Add or retint a type by editing the ICONS map below — it drives <Files> and <FileTree>.

import type { CSSProperties, ReactNode } from 'react'

type Badge = { label: string; bg: string; fg?: string }

// ext → coloured badge. Colours track each tool's brand.
const ICONS: Record<string, Badge> = {
  go: { label: 'GO', bg: '#00ADD8' },
  mod: { label: 'GO', bg: '#00758D' },
  sum: { label: 'GO', bg: '#00758D' },
  js: { label: 'JS', bg: '#F7DF1E', fg: '#1a1a1a' },
  mjs: { label: 'JS', bg: '#F7DF1E', fg: '#1a1a1a' },
  cjs: { label: 'JS', bg: '#F7DF1E', fg: '#1a1a1a' },
  ts: { label: 'TS', bg: '#3178C6' },
  json: { label: '{ }', bg: '#F59E0B', fg: '#1a1a1a' },
  md: { label: 'MD', bg: '#4B5563' },
  mdx: { label: 'MDX', bg: '#4B5563' },
  css: { label: 'CSS', bg: '#38BDF8', fg: '#0a2a3a' },
  scss: { label: 'SASS', bg: '#EC4899' },
  html: { label: '</>', bg: '#E34F26' },
  env: { label: 'ENV', bg: '#10B981', fg: '#062b1e' },
  yml: { label: 'YML', bg: '#8B5CF6' },
  yaml: { label: 'YML', bg: '#8B5CF6' },
  toml: { label: 'TOML', bg: '#FB923C', fg: '#3a1c04' },
  sql: { label: 'SQL', bg: '#FB923C', fg: '#3a1c04' },
  sh: { label: '>_', bg: '#10B981', fg: '#062b1e' },
  bash: { label: '>_', bg: '#10B981', fg: '#062b1e' },
  prisma: { label: 'PSL', bg: '#14B8A6', fg: '#052e29' },
  dockerfile: { label: '🐳', bg: '#2496ED' },
  lock: { label: '🔒', bg: '#374151' },
  png: { label: 'IMG', bg: '#D946EF' },
  jpg: { label: 'IMG', bg: '#D946EF' },
  jpeg: { label: 'IMG', bg: '#D946EF' },
  ico: { label: 'IMG', bg: '#D946EF' },
}

function extOf(name: string): string {
  const clean = name.replace(/\/+$/, '')
  if (clean.toLowerCase() === 'dockerfile') return 'dockerfile'
  return clean.includes('.')
    ? clean.split('.').pop()!.toLowerCase()
    : clean.replace(/^\./, '').toLowerCase()
}

/** A crisp 24×24 rounded badge with a short label — VSCode-Material feel. */
function BadgeIcon({ badge, style }: { badge: Badge; style?: CSSProperties }) {
  const len = [...badge.label].length
  const fontSize = len <= 2 ? 10 : len === 3 ? 8 : 6.5
  return (
    <svg viewBox="0 0 24 24" style={style} className="shrink-0" aria-hidden>
      <rect x="1.5" y="1.5" width="21" height="21" rx="5" fill={badge.bg} />
      <text
        x="12"
        y="12.5"
        textAnchor="middle"
        dominantBaseline="central"
        fontFamily="var(--font-mono, ui-monospace, monospace)"
        fontWeight="700"
        fontSize={fontSize}
        fill={badge.fg ?? '#ffffff'}
      >
        {badge.label}
      </text>
    </svg>
  )
}

/** The React atom, for .tsx / .jsx. */
function ReactIcon({ style }: { style?: CSSProperties }) {
  return (
    <svg viewBox="0 0 24 24" style={style} className="shrink-0" aria-hidden>
      <g transform="translate(12 12)" fill="none" stroke="#61DAFB" strokeWidth="1.1">
        <ellipse rx="10.5" ry="4" />
        <ellipse rx="10.5" ry="4" transform="rotate(60)" />
        <ellipse rx="10.5" ry="4" transform="rotate(120)" />
      </g>
      <circle cx="12" cy="12" r="2" fill="#61DAFB" />
    </svg>
  )
}

/** Generic document, for unknown types. */
function DocIcon({ style }: { style?: CSSProperties }) {
  return (
    <svg viewBox="0 0 24 24" style={style} className="shrink-0" aria-hidden fill="none">
      <path
        d="M6 2.5h7l5 5V21a.5.5 0 0 1-.5.5h-11A.5.5 0 0 1 6 21V3a.5.5 0 0 1 .5-.5Z"
        fill="#475569"
      />
      <path d="M13 2.5v5h5" fill="#64748b" />
    </svg>
  )
}

/**
 * File-type glyph for `name`. `className` controls the box size (e.g. "h-4 w-4").
 * Returns a React atom for .tsx/.jsx, a brand badge for known types, else a doc.
 */
export function FileIcon({ name, className }: { name: string; className?: string }) {
  const ext = extOf(name)
  const style = { width: '1em', height: '1em' } as CSSProperties
  const wrap = (node: ReactNode) => (
    <span className={className} style={{ display: 'inline-flex', lineHeight: 0 }}>
      {node}
    </span>
  )
  if (ext === 'tsx' || ext === 'jsx') return wrap(<ReactIcon style={style} />)
  const badge = ICONS[ext]
  if (badge) return wrap(<BadgeIcon badge={badge} style={style} />)
  return wrap(<DocIcon style={style} />)
}
