/*
 * The heading in the middle, integrations arranged on a ring around it.
 *
 * The positions are computed with trigonometry from the number of items, so
 * adding a seventh integration redistributes all seven evenly instead of
 * leaving a gap where hand-written percentages ran out. Angles start at
 * -90 degrees so the first item sits at the top rather than at three o'clock,
 * which is where 0 radians actually points.
 *
 * Each tile is placed with `top`/`left` percentages and then pulled back by
 * half its own size with a translate, because a percentage position places the
 * element's corner, not its centre. Skip the translate and every tile hangs
 * down and to the right of where you meant.
 *
 * The ring is `aspect-square` with a max width rather than a fixed pixel size,
 * so it scales down with the viewport instead of pushing the page sideways.
 * Below `sm` the ring is dropped entirely and the tiles become a plain wrapped
 * row: a circle of logos inside a 320px column is a pile, not a ring.
 *
 * Logos are lettered tiles with a `logo` slot — a component library should not
 * bundle other companies' trademarks on your behalf.
 */

import type { ReactNode } from 'react'

export interface Integration {
  name: string
  /** Two-character label for the fallback tile. Defaults to the first two
   *  letters of the name, which collides for names like Redis and Resend. */
  short?: string
  logo?: ReactNode
  tone?: string
}

const INTEGRATIONS: Integration[] = [
  { name: 'Postgres', short: 'PG', tone: 'bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400' },
  { name: 'Redis', short: 'RD', tone: 'bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400' },
  { name: 'S3', short: 'S3', tone: 'bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400' },
  { name: 'Resend', short: 'RS', tone: 'bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400' },
  { name: 'Stripe', short: 'ST', tone: 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400' },
  { name: 'Sentry', short: 'SN', tone: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400' },
  { name: 'Docker', short: 'DK', tone: 'bg-cyan-50 text-cyan-600 dark:bg-cyan-500/10 dark:text-cyan-400' },
  { name: 'GitHub', short: 'GH', tone: 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300' },
]

function Tile({ item }: { item: Integration }) {
  if (item.logo) return <>{item.logo}</>
  return (
    /* aria-hidden: the sr-only name beside it is the accessible label, and
       without this a screen reader reads the monogram too — "PG, Postgres". */
    <span
      aria-hidden="true"
      className={`flex size-12 items-center justify-center rounded-xl border border-gray-200 text-sm font-semibold shadow-sm dark:border-white/10 ${
        item.tone ?? 'bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300'
      }`}
    >
      {item.short ?? item.name.slice(0, 2)}
    </span>
  )
}

export default function IntegrationsOrbitRing({
  title = 'Seamlessly integrate with the tools you already use',
  subtitle = 'Connect the services you already pay for and keep your workflow where it is.',
  ctaLabel = 'Browse integrations',
  ctaHref = '#',
  integrations = INTEGRATIONS,
}: {
  title?: string
  subtitle?: string
  ctaLabel?: string
  ctaHref?: string
  integrations?: Integration[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        {/* Below sm the ring is a wrapped row. A circle in a 320px column is a
            pile of overlapping tiles, not a ring. */}
        <ul
          role="list"
          className="mx-auto flex max-w-xs flex-wrap justify-center gap-4 sm:hidden"
        >
          {integrations.map((item) => (
            <li key={item.name}>
              <Tile item={item} />
              <span className="sr-only">{item.name}</span>
            </li>
          ))}
        </ul>

        <div className="mt-16 sm:mt-0">
          <div className="relative mx-auto hidden aspect-square w-full max-w-2xl sm:block">
            <ul role="list">
              {integrations.map((item, i) => {
                /* -90deg so the first tile sits at the top, not at 3 o'clock. */
                const angle = (i / integrations.length) * 2 * Math.PI - Math.PI / 2
                const radius = 44 // percent of half the container
                const left = 50 + radius * Math.cos(angle)
                const top = 50 + radius * Math.sin(angle)
                return (
                  <li
                    key={item.name}
                    /* -translate-1/2 recentres the tile: a percentage position
                       places its corner, not its middle. */
                    className="absolute -translate-x-1/2 -translate-y-1/2"
                    style={{ left: `${left}%`, top: `${top}%` }}
                  >
                    <Tile item={item} />
                    <span className="sr-only">{item.name}</span>
                  </li>
                )
              })}
            </ul>

            <div
              aria-hidden="true"
              className="absolute inset-[6%] rounded-full border border-dashed border-gray-200 dark:border-white/10"
            />

            <div className="absolute inset-0 flex flex-col items-center justify-center px-16 text-center">
              <h2 className="max-w-md text-3xl font-semibold tracking-tight text-balance text-gray-900 sm:text-4xl dark:text-white">
                {title}
              </h2>
              <p className="mt-4 max-w-sm text-base/7 text-pretty text-gray-600 dark:text-gray-400">
                {subtitle}
              </p>
              <a
                href={ctaHref}
                className="mt-8 inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-5 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
              >
                {ctaLabel}
              </a>
            </div>
          </div>

          {/* The heading again for the stacked layout, where the ring is gone. */}
          <div className="mt-10 text-center sm:hidden">
            <h2 className="text-3xl font-semibold tracking-tight text-balance text-gray-900 dark:text-white">
              {title}
            </h2>
            <p className="mt-4 text-base/7 text-pretty text-gray-600 dark:text-gray-400">
              {subtitle}
            </p>
            <a
              href={ctaHref}
              className="mt-8 inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-5 text-sm font-semibold text-gray-900 dark:border-white/15 dark:text-white"
            >
              {ctaLabel}
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
