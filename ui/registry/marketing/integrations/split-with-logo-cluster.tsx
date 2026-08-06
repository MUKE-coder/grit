/*
 * Copy on the left, a loose cluster of tiles on the right.
 *
 * The offsets are per-item classes rather than random values, because a
 * component that scatters its own tiles with Math.random renders differently
 * on the server and on the client, and React will tell you about it in the
 * console before your users tell you about it in a bug report. Deliberate
 * offsets also let you keep the arrangement balanced, which randomness does
 * not.
 *
 * On a phone the offsets are dropped and the tiles sit on a plain grid. A
 * scattered arrangement needs room around it to read as scattered; in a narrow
 * column it just reads as misaligned.
 *
 * Logos are lettered tiles with a `logo` slot — a component library should not
 * bundle other companies' trademarks on your behalf, and an honest placeholder
 * beats one you forget to replace.
 */

import type { ReactNode } from 'react'

export interface Integration {
  name: string
  /** Two-character label for the fallback tile. Defaults to the first two
   *  letters of the name, which collides for names like Redis and Resend. */
  short?: string
  logo?: ReactNode
  tone?: string
  /** Vertical nudge above sm, e.g. 'sm:mt-8'. Purely visual. */
  offset?: string
}

const INTEGRATIONS: Integration[] = [
  { name: 'Postgres', short: 'PG', tone: 'bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400', offset: 'sm:mt-10' },
  { name: 'Redis', short: 'RD', tone: 'bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400' },
  { name: 'S3', short: 'S3', tone: 'bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400', offset: 'sm:mt-6' },
  { name: 'Resend', short: 'RS', tone: 'bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400', offset: 'sm:-mt-4' },
  { name: 'Stripe', short: 'ST', tone: 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400', offset: 'sm:mt-4' },
  { name: 'Sentry', short: 'SN', tone: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400' },
  { name: 'Docker', short: 'DK', tone: 'bg-cyan-50 text-cyan-600 dark:bg-cyan-500/10 dark:text-cyan-400', offset: 'sm:mt-8' },
  { name: 'GitHub', short: 'GH', tone: 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300', offset: 'sm:-mt-2' },
  { name: 'Vercel', short: 'VC', tone: 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300', offset: 'sm:mt-6' },
]

export default function IntegrationsSplitWithLogoCluster({
  title = 'Seamlessly integrate with your favourite tools',
  subtitle = 'Everything the framework ships with is a real integration rather than a wrapper: the config lives in one place, the credentials in your environment, and the failure modes are documented.',
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
        <div className="grid grid-cols-1 items-center gap-x-16 gap-y-12 lg:grid-cols-2">
          <div>
            <h2 className="max-w-md text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h2>
            <p className="mt-6 max-w-lg text-base/7 text-pretty text-gray-600 dark:text-gray-400">
              {subtitle}
            </p>
            <a
              href={ctaHref}
              className="mt-8 inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-5 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
            >
              {ctaLabel}
            </a>
          </div>

          {/* Offsets are declared per item, not generated. A random scatter
              renders differently on the server and the client. */}
          <ul
            role="list"
            className="grid grid-cols-3 justify-items-center gap-4 sm:gap-6"
          >
            {integrations.map((item) => (
              <li key={item.name} className={item.offset ?? ''}>
                {item.logo ?? (
                  <span
                    aria-hidden="true"
                    className={`flex size-16 items-center justify-center rounded-2xl border border-gray-200 text-base font-semibold shadow-sm dark:border-white/10 ${
                      item.tone ?? 'bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300'
                    }`}
                  >
                    {item.short ?? item.name.slice(0, 2)}
                  </span>
                )}
                <span className="sr-only">{item.name}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </section>
  )
}
