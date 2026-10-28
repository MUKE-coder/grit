import Link from 'next/link'
import { ArrowUpRight, Clock } from 'lucide-react'
import { CopyPageButton } from '@/components/docs/copy-page'
import { GRIT_VERSION, siteConfig } from '@/config/site'

/**
 * The top of a docs page: where you are, what it is, and what you can do
 * with it.
 *
 * The meta row is the part worth arguing for. A reading time tells somebody
 * whether to start now or bookmark it. "Updated for v3.415.0" answers the
 * question every reader of framework documentation actually has, which is
 * not when it was written but whether it still applies. Copy page is for
 * pasting into an agent. Edit on GitHub is how a typo gets fixed by the
 * person who found it rather than by the person who has to be told.
 */

export interface Crumb {
  label: string
  href?: string
}

export function DocsPageHeader({
  crumbs,
  title,
  lead,
  /** Minutes. Omit where the page is a reference rather than a read. */
  readingTime,
  /** The repo path of this page, for the Edit link. */
  sourcePath,
  /** Set where a page documents something older than the current release. */
  updatedFor = GRIT_VERSION,
}: {
  crumbs: Crumb[]
  title: string
  lead?: React.ReactNode
  readingTime?: number
  sourcePath?: string
  updatedFor?: string
}) {
  const editUrl = sourcePath
    ? `${siteConfig.github}/edit/main/docs/${sourcePath}`
    : null

  return (
    <header className="mb-10">
      <nav aria-label="Breadcrumb">
        <ol className="flex flex-wrap items-center gap-2 p-0 text-[13px] text-fg-subtle">
          {crumbs.map((crumb, i) => (
            <li key={crumb.label} className="flex items-center gap-2">
              {i > 0 && <span aria-hidden="true">/</span>}
              {crumb.href ? (
                <Link href={crumb.href} className="transition-colors hover:text-fg">
                  {crumb.label}
                </Link>
              ) : (
                <span className="text-fg-muted" aria-current="page">
                  {crumb.label}
                </span>
              )}
            </li>
          ))}
        </ol>
      </nav>

      <h1 className="mt-4 text-[clamp(2.125rem,3.4vw,2.75rem)] font-semibold leading-[1.1] tracking-[-0.035em] text-fg">
        {title}
      </h1>

      {lead && (
        <p className="mt-4 max-w-[var(--measure)] text-[19px] leading-[1.65] text-fg-muted">{lead}</p>
      )}

      <div className="mt-5 flex flex-wrap items-center gap-x-4 gap-y-2 text-[13px] text-fg-subtle">
        {readingTime && (
          <>
            <span className="inline-flex items-center gap-1.5">
              <Clock className="h-3.5 w-3.5" aria-hidden="true" />
              {readingTime} min
            </span>
            <span aria-hidden="true">·</span>
          </>
        )}
        <span>Updated for v{updatedFor}</span>
        <span aria-hidden="true">·</span>
        <CopyPageButton />
        {editUrl && (
          <a
            href={editUrl}
            target="_blank"
            rel="noreferrer"
            className="inline-flex items-center gap-1.5 transition-colors hover:text-fg"
          >
            Edit on GitHub
            <ArrowUpRight className="h-3 w-3" aria-hidden="true" />
          </a>
        )}
      </div>
    </header>
  )
}
