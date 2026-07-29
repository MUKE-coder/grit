'use client'

import Link from 'next/link'
import { ArrowUpRight } from 'lucide-react'
import { CopyButton } from './copy-button'

/**
 * One gallery tile: a scaled-down live preview above the component's name and
 * install command.
 *
 * The preview is a 1280px-wide iframe scaled down to fit, rather than a
 * screenshot. Screenshots go stale the moment a component changes and nobody
 * notices for months; a live frame cannot lie about what you are installing.
 * loading="lazy" keeps a hundred frames from all mounting at once.
 */
export function ComponentCard({
  name,
  title,
  description,
  category,
  installCommand,
}: {
  name: string
  title: string
  description: string
  category: string
  installCommand: string
}) {
  return (
    <div className="group flex flex-col overflow-hidden rounded-2xl border border-border bg-bg-secondary transition-colors hover:border-accent/40">
      {/* Preview */}
      <Link
        href={`/c/${name}`}
        className="relative block h-56 overflow-hidden border-b border-border bg-bg-primary"
      >
        <iframe
          src={`/preview/${name}`}
          title={`${title} preview`}
          loading="lazy"
          tabIndex={-1}
          aria-hidden
          className="pointer-events-none origin-top-left border-0"
          style={{
            width: 1280,
            height: 900,
            transform: 'scale(0.32)',
          }}
        />
        <span className="absolute inset-0 bg-gradient-to-t from-bg-secondary/80 via-transparent to-transparent opacity-0 transition-opacity group-hover:opacity-100" />
      </Link>

      {/* Meta */}
      <div className="flex flex-1 flex-col gap-3 p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0">
            <Link
              href={`/c/${name}`}
              className="flex items-center gap-1 text-sm font-semibold text-foreground hover:text-accent transition-colors"
            >
              {title}
              <ArrowUpRight
                size={13}
                className="opacity-0 transition-opacity group-hover:opacity-100"
              />
            </Link>
            <p className="mt-1 line-clamp-2 text-xs leading-relaxed text-text-muted">
              {description}
            </p>
          </div>
          <span className="shrink-0 rounded-full border border-border bg-bg-elevated px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider text-text-muted">
            {category}
          </span>
        </div>

        <div className="mt-auto flex items-center justify-between gap-2 rounded-lg border border-border bg-bg-primary px-2.5 py-1.5">
          <code className="truncate font-mono text-[11px] text-text-secondary">
            {name}
          </code>
          <CopyButton value={installCommand} label="Install" />
        </div>
      </div>
    </div>
  )
}
