import Link from 'next/link'
import { ArrowUpRight, Github } from 'lucide-react'

/**
 * CornerTicks — the signature frame from DESIGN.md section 4.
 *
 * Four 10px L-marks rather than a heavier border. It reads as a crop mark on
 * a sheet: the card is a region of the page rather than an object floating
 * above it, which is the same reason the card does not take a shadow.
 */
export function CornerTicks() {
  const arm = 'absolute h-2.5 w-2.5 border-line-strong transition-colors duration-150'
  return (
    <span aria-hidden="true" className="pointer-events-none">
      <span className={`${arm} -left-px -top-px border-l border-t`} />
      <span className={`${arm} -right-px -top-px border-r border-t`} />
      <span className={`${arm} -bottom-px -left-px border-b border-l`} />
      <span className={`${arm} -bottom-px -right-px border-b border-r`} />
    </span>
  )
}

export interface PluginEntry {
  name: string
  title: string
  summary: string
  href: string
  /** Lucide icon name, drawn large on the cover. */
  icon: React.ReactNode
  tags: string[]
  /** official | module | community */
  kind: 'official' | 'module' | 'community'
  /** Community only. */
  author?: string
  repo?: string
  install: string
}

const KIND_LABEL: Record<PluginEntry['kind'], string> = {
  official: 'Official',
  module: 'Go module',
  community: 'Community',
}

/**
 * PluginCard.
 *
 * The cover is drawn rather than photographed. Twenty-six plugins would need
 * twenty-six screenshots that all go stale, and a wrong screenshot is worse
 * than none: it promises a screen the plugin does not have. A large icon on
 * the dot grid is honest, consistent, and costs nothing to keep current.
 *
 * Hover moves the border and nothing else. DESIGN.md section 6: a card that
 * lifts or scales is an object; a card that brightens its edge is a region,
 * and the whole page is built from regions.
 */
export function PluginCard({ plugin }: { plugin: PluginEntry }) {
  return (
    <Link
      href={plugin.href}
      className="group relative flex flex-col rounded-lg border border-border bg-card transition-colors duration-150 hover:border-accent/50"
    >
      <CornerTicks />

      {/* Cover */}
      <div className="relative flex h-28 items-center justify-center overflow-hidden rounded-t-lg border-b border-border bg-background/60">
        <div
          aria-hidden="true"
          className="absolute inset-0 opacity-[0.07]"
          style={{
            backgroundImage: 'radial-gradient(currentColor 1px, transparent 1px)',
            backgroundSize: '16px 16px',
          }}
        />
        <span className="relative text-accent [&>svg]:h-8 [&>svg]:w-8">{plugin.icon}</span>
      </div>

      {/* Body */}
      <div className="flex flex-1 flex-col gap-3 p-5">
        <div>
          <h3 className="text-base font-semibold tracking-tight text-foreground">
            {plugin.title}
          </h3>
          <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">{plugin.summary}</p>
        </div>

        {plugin.tags.length > 0 && (
          <ul className="flex flex-wrap gap-1.5">
            {plugin.tags.map((t) => (
              <li
                key={t}
                className="rounded-sm border border-border px-2 py-0.5 font-mono text-[11px] text-muted-foreground"
              >
                {t}
              </li>
            ))}
          </ul>
        )}

        <code className="mt-auto block truncate rounded-sm bg-background px-2.5 py-1.5 font-mono text-[12px] text-muted-foreground">
          {plugin.install}
        </code>
      </div>

      {/* Footer */}
      <div className="flex items-center gap-2 border-t border-border px-5 py-3">
        <span className="font-mono text-[11px] uppercase tracking-[0.12em] text-muted-foreground">
          {KIND_LABEL[plugin.kind]}
        </span>
        {plugin.author && (
          <>
            <span aria-hidden="true" className="text-muted-foreground">
              ·
            </span>
            <span className="truncate text-xs text-muted-foreground">{plugin.author}</span>
          </>
        )}
        <ArrowUpRight className="ml-auto h-4 w-4 shrink-0 text-muted-foreground transition-colors duration-150 group-hover:text-accent" />
      </div>
    </Link>
  )
}

/**
 * ModuleRow — the Go module plugins.
 *
 * A row, not a card, because these install with `go get` rather than
 * `grit plugin add`. Giving them the same card as the others would say they
 * are the same thing, and somebody would run the wrong command.
 */
export function ModuleRow({ name, summary }: { name: string; summary: string }) {
  return (
    <div className="flex items-start gap-4 border-b border-border py-3.5 last:border-0">
      <code className="w-36 shrink-0 font-mono text-sm text-accent">{name}</code>
      <p className="flex-1 text-sm text-muted-foreground">{summary}</p>
    </div>
  )
}

/** Stat — the hero figures. Mono numerals so the row does not jitter. */
export function Stat({ value, label }: { value: string; label: string }) {
  return (
    <div className="relative px-5 py-4">
      <div className="font-mono text-2xl font-bold tabular-nums text-foreground">{value}</div>
      <div className="mt-0.5 font-mono text-[11px] uppercase tracking-[0.12em] text-muted-foreground">
        {label}
      </div>
    </div>
  )
}

/** CommunityEmpty — an empty state that asks for something. */
export function CommunityEmpty() {
  return (
    <div className="relative rounded-lg border border-dashed border-border p-10 text-center">
      <Github className="mx-auto h-5 w-5 text-muted-foreground" aria-hidden="true" />
      <p className="mt-3 font-medium text-foreground">Nothing listed yet</p>
      <p className="mx-auto mt-1.5 max-w-sm text-sm text-muted-foreground">
        If you have written one, yours can be the first. Submission is three commands and a
        GitHub issue.
      </p>
      <Link
        href="/docs/plugins/authoring"
        className="mt-5 inline-flex h-10 items-center gap-2 rounded-full border border-border px-5 text-sm font-medium text-foreground transition-colors duration-150 hover:border-accent/50"
      >
        Write a plugin
        <ArrowUpRight className="h-4 w-4" aria-hidden="true" />
      </Link>
    </div>
  )
}
