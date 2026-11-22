/**
 * The maturity of a subsystem, on the page that documents it.
 *
 * A reader of /docs/backend/outbox could not tell whether durable events were
 * a core guarantee with live tests behind them or something that shipped last
 * month. /docs/stability knew, and nothing on the feature page said so
 * (grit#247).
 *
 * It reads config/stability.ts, which is where the answer already lives. A
 * second file saying the same thing in different words is the bug this project
 * keeps finding, so there is no maturity.ts: there is one table, and this is a
 * second view of it.
 *
 * Server component: static text and a link, with nothing to hydrate.
 */

import Link from 'next/link'
import { STABILITY_AREAS, statusStyle, type Status } from '@/config/stability'

/** What each badge means in one line, so the word is not doing it alone. */
const MEANING: Record<Status, string> = {
  stable:
    'Proven: named tests, and the live suite drives it over HTTP on every push.',
  beta:
    'Tested and in use, with less live coverage or a shorter history than the stable subsystems.',
  new: 'Recent. It works and it is tested; it has had fewer projects through it.',
}

interface MaturityBadgeProps {
  /**
   * One or more ids from config/stability.ts. Several because one page can
   * document two subsystems: optimistic locking and encryption at rest are
   * both on the database page.
   */
  areas: string[]
}

export function MaturityBadge({ areas }: MaturityBadgeProps) {
  const found = areas
    .map((id) => STABILITY_AREAS.find((a) => a.id === id))
    .filter((a): a is (typeof STABILITY_AREAS)[number] => Boolean(a))

  // An id that names no area is a typo, and a badge that silently renders
  // nothing is how it would survive. scripts/check-maturity-badges.py fails
  // the build on one, so this branch is the belt to that pair of braces.
  if (found.length === 0) return null

  return (
    <div className="mb-8 flex flex-wrap items-center gap-x-3 gap-y-2 rounded-xl border border-border/50 bg-card/30 px-4 py-3">
      {found.map((area) => (
        <span key={area.id} className="inline-flex items-center gap-2">
          <span
            className={`inline-flex items-center rounded-full border px-2.5 py-0.5 text-[11px] font-medium uppercase tracking-wide ${statusStyle[area.status]}`}
          >
            {area.status}
          </span>
          <span className="text-[12.5px] text-muted-foreground">
            {found.length > 1 ? `${area.area}. ` : ''}
            {MEANING[area.status]}
          </span>
        </span>
      ))}
      <Link
        href={`/docs/stability#${found[0].id}`}
        className="text-[12.5px] text-accent hover:underline whitespace-nowrap"
      >
        What is tested
      </Link>
    </div>
  )
}
