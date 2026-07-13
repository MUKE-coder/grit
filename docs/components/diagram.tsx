import { ArrowDown, ArrowRight } from 'lucide-react'
import { FileIcon } from '@/components/file-icons'
import type { ReactNode } from 'react'
import { folderEmoji } from '@/lib/file-meta'

/* ─────────────────────────────────────────────────────────────
   Visual diagram toolkit for the docs — styled boxes, flow
   arrows, and file trees, replacing monospace ASCII art.
   All presentational (server components).
   ───────────────────────────────────────────────────────────── */

type Tone = 'default' | 'primary' | 'blue' | 'green' | 'amber' | 'rose' | 'cyan' | 'violet'

const TONES: Record<Tone, string> = {
  default: 'border-border/60 bg-card/60 text-foreground',
  primary: 'border-primary/30 bg-primary/[0.07] text-foreground',
  blue: 'border-blue-500/30 bg-blue-500/[0.07] text-foreground',
  green: 'border-emerald-500/30 bg-emerald-500/[0.07] text-foreground',
  amber: 'border-amber-500/30 bg-amber-500/[0.07] text-foreground',
  rose: 'border-rose-500/30 bg-rose-500/[0.07] text-foreground',
  cyan: 'border-cyan-500/30 bg-cyan-500/[0.07] text-foreground',
  violet: 'border-violet-500/30 bg-violet-500/[0.07] text-foreground',
}

const ACCENT: Record<Tone, string> = {
  default: 'text-muted-foreground',
  primary: 'text-primary',
  blue: 'text-blue-400',
  green: 'text-emerald-400',
  amber: 'text-amber-400',
  rose: 'text-rose-400',
  cyan: 'text-cyan-400',
  violet: 'text-violet-400',
}

/** Outer container for a diagram. */
export function Diagram({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`my-6 rounded-xl border border-border/50 bg-gradient-to-b from-card/40 to-background/40 p-5 sm:p-6 ${className}`}>
      {children}
    </div>
  )
}

/** A labelled node/box in a flow diagram. */
export function DiagramBox({
  title,
  sub,
  tone = 'default',
  icon,
  className = '',
}: {
  title: ReactNode
  sub?: ReactNode
  tone?: Tone
  icon?: ReactNode
  className?: string
}) {
  return (
    <div className={`flex-1 rounded-lg border px-3.5 py-2.5 text-center ${TONES[tone]} ${className}`}>
      <div className="flex items-center justify-center gap-1.5 text-[13px] font-semibold">
        {icon && <span className={ACCENT[tone]}>{icon}</span>}
        {title}
      </div>
      {sub && <div className="mt-0.5 text-[11px] font-mono text-muted-foreground/70">{sub}</div>}
    </div>
  )
}

/** A horizontal row of boxes (a layer). */
export function DiagramRow({ children, className = '' }: { children: ReactNode; className?: string }) {
  return <div className={`flex flex-col sm:flex-row gap-3 ${className}`}>{children}</div>
}

/** A vertical connector between layers, optionally labelled. */
export function DiagramArrow({ label }: { label?: ReactNode }) {
  return (
    <div className="flex items-center justify-center gap-2 py-1.5 text-muted-foreground/60">
      <ArrowDown className="h-4 w-4" />
      {label && <span className="text-[11px] font-mono">{label}</span>}
    </div>
  )
}

/** An inline right-arrow for horizontal flows (e.g. Handler → Service). */
export function FlowArrow() {
  return <ArrowRight className="h-4 w-4 shrink-0 self-center text-muted-foreground/50" />
}

/* ── File tree ──────────────────────────────────────────────── */

export interface TreeNode {
  name: string
  /** indentation level (0 = root child) */
  depth?: number
  type?: 'folder' | 'file'
  /** muted note shown to the right */
  comment?: string
  /** highlight this row (e.g. a generated file) */
  highlight?: boolean
}

/** A styled file/folder tree — replaces monospace directory listings. */
export function FileTree({ title, nodes }: { title?: string; nodes: TreeNode[] }) {
  return (
    <div className="my-6 overflow-hidden rounded-xl border border-border/50 bg-card/40">
      {title && (
        <div className="border-b border-border/50 bg-muted/20 px-4 py-2 text-[12px] font-mono text-muted-foreground">
          {title}
        </div>
      )}
      <div className="p-2">
        {nodes.map((n, i) => {
          const isFolder = n.type !== 'file'
          return (
            <div
              key={i}
              className={`flex items-center gap-2 rounded-md px-2 py-1 text-[13px] ${
                n.highlight ? 'bg-primary/[0.07] border-l-2 border-primary' : 'border-l-2 border-transparent'
              }`}
              style={{ paddingLeft: `${(n.depth ?? 0) * 18 + 8}px` }}
            >
              {isFolder ? (
                <span className="shrink-0 text-[13px] leading-none">{folderEmoji(true)}</span>
              ) : (
                <FileIcon name={n.name} className="h-4 w-4 shrink-0" />
              )}
              <span className={`font-mono ${isFolder ? 'font-medium text-foreground/90' : 'text-muted-foreground'}`}>
                {n.name}
              </span>
              {n.comment && (
                <span className="ml-auto pl-3 text-[11px] text-muted-foreground/50">{n.comment}</span>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

/** A small legend row for coloured diagrams. */
export function DiagramLegend({ items }: { items: { tone: Tone; label: string }[] }) {
  return (
    <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1.5">
      {items.map((it) => (
        <div key={it.label} className="flex items-center gap-1.5 text-[11px] text-muted-foreground">
          <span className={`h-2.5 w-2.5 rounded-sm border ${TONES[it.tone]}`} />
          {it.label}
        </div>
      ))}
    </div>
  )
}

/**
 * A tinted, labelled grouping box — for wrapping related nodes in a diagram
 * (e.g. a "Session Management" or "Middleware stack" cluster). The label
 * floats on the top border, Next.js-style.
 */
export function HighlightBox({
  label,
  tone = 'primary',
  children,
  className = '',
}: {
  label?: ReactNode
  tone?: Tone
  children: ReactNode
  className?: string
}) {
  return (
    <div className={`relative rounded-lg border-2 ${TONES[tone]} px-3 pb-3 pt-4 ${className}`}>
      {label && (
        <span className={`absolute -top-2.5 left-3 rounded bg-background px-1.5 text-[11px] font-semibold ${ACCENT[tone]}`}>
          {label}
        </span>
      )}
      <div className="flex flex-wrap gap-2">{children}</div>
    </div>
  )
}
