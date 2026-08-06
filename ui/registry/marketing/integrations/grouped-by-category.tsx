/*
 * Integrations grouped into the categories people actually search by.
 *
 * The grouping is the whole value. A flat wall of forty logos answers "do you
 * have a lot of integrations?", which nobody is asking. Grouping answers "do
 * you work with my database, my model provider, my host?", which is the
 * question that decides the sale.
 *
 * Each group is a <section> with its heading referenced by `aria-labelledby`,
 * so a screen reader announces "Hosting, list, 3 items" rather than reading
 * twelve tile labels in a row with no idea where one group ended.
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

export interface Group {
  title: string
  items: Integration[]
}

const GROUPS: Group[] = [
  {
    title: 'Data',
    items: [
      { name: 'Postgres', short: 'PG', tone: 'bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400' },
      { name: 'Redis', short: 'RD', tone: 'bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400' },
      { name: 'SQLite', short: 'SQ', tone: 'bg-slate-100 text-slate-600 dark:bg-white/10 dark:text-slate-300' },
    ],
  },
  {
    title: 'Models',
    items: [
      { name: 'Claude', short: 'CL', tone: 'bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400' },
      { name: 'OpenAI', short: 'AI', tone: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400' },
      { name: 'Ollama', short: 'OL', tone: 'bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400' },
    ],
  },
  {
    title: 'Hosting',
    items: [
      { name: 'Dokploy', short: 'DK', tone: 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400' },
      { name: 'Railway', short: 'RW', tone: 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300' },
      { name: 'Fly.io', short: 'FL', tone: 'bg-cyan-50 text-cyan-600 dark:bg-cyan-500/10 dark:text-cyan-400' },
    ],
  },
]

export default function IntegrationsGroupedByCategory({
  title = 'Seamless integration',
  subtitle = 'Connect the services you already pay for. Configuration lives in one file and credentials stay in your environment.',
  ctaLabel = 'Get started',
  ctaHref = '#',
  groups = GROUPS,
}: {
  title?: string
  subtitle?: string
  ctaLabel?: string
  ctaHref?: string
  groups?: Group[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 text-center lg:px-8">
        <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>
        <p className="mx-auto mt-6 max-w-2xl text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>
        <a
          href={ctaHref}
          className="mt-8 inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-5 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
        >
          {ctaLabel}
        </a>

        <div className="mx-auto mt-16 grid max-w-4xl grid-cols-1 gap-6 sm:grid-cols-3">
          {groups.map((group, groupIndex) => (
            /* Labelled group: "Hosting, list, 3 items", not twelve loose tiles. */
            <section
              key={group.title}
              aria-labelledby={`integration-group-${groupIndex}`}
              className="rounded-2xl border border-gray-200 p-6 dark:border-white/10"
            >
              <h3
                id={`integration-group-${groupIndex}`}
                className="text-sm font-medium text-gray-600 dark:text-gray-400"
              >
                {group.title}
              </h3>
              <ul role="list" className="mt-4 flex justify-center gap-3">
                {group.items.map((item) => (
                  <li key={item.name}>
                    {item.logo ?? (
                      <span
                        aria-hidden="true"
                        className={`flex size-12 items-center justify-center rounded-xl border border-gray-200 text-sm font-semibold shadow-sm dark:border-white/10 ${
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
            </section>
          ))}
        </div>
      </div>
    </section>
  )
}
