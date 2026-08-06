/*
 * Six integrations as bordered cards.
 *
 * On the logos: this ships lettered tiles rather than real brand marks, and
 * takes a `logo` slot so you can drop yours in. That is deliberate. A
 * component library cannot bundle other companies' trademarks on your behalf —
 * most brand guidelines require permission and forbid restyling — and a
 * placeholder that is honestly a placeholder is better than one you forget to
 * replace. Pass an <img> or an inline <svg> and the tile gets out of the way.
 *
 * The whole card is not a link. Only the name is, so a screen reader announces
 * "Gemini, link" rather than reading the name, the description and the tile
 * label as one enormous link name. `after:absolute after:inset-0` on that link
 * stretches its hit area over the card, which keeps the large click target
 * without wrecking the announcement.
 */

import type { ReactNode } from 'react'

export interface Integration {
  name: string
  /** Two-character label for the fallback tile. Defaults to the first two
   *  letters of the name, which collides for names like Redis and Resend. */
  short?: string
  description: string
  href?: string
  logo?: ReactNode
  /** Tailwind classes for the fallback tile, e.g. 'bg-blue-50 text-blue-600'. */
  tone?: string
}

const INTEGRATIONS: Integration[] = [
  {
    name: 'Postgres',
    short: 'PG',
    description: 'Read the schema, generate the models, keep migrations in version control.',
    tone: 'bg-sky-50 text-sky-600 dark:bg-sky-500/10 dark:text-sky-400',
  },
  {
    name: 'Redis',
    short: 'RD',
    description: 'Response caching and the background job queue, wired up out of the box.',
    tone: 'bg-rose-50 text-rose-600 dark:bg-rose-500/10 dark:text-rose-400',
  },
  {
    name: 'S3',
    short: 'S3',
    description: 'Uploads, signed URLs and image processing against any S3-compatible store.',
    tone: 'bg-amber-50 text-amber-600 dark:bg-amber-500/10 dark:text-amber-400',
  },
  {
    name: 'Resend',
    short: 'RS',
    description: 'Transactional email with templates you can preview before you send them.',
    tone: 'bg-violet-50 text-violet-600 dark:bg-violet-500/10 dark:text-violet-400',
  },
  {
    name: 'Stripe',
    short: 'ST',
    description: 'Checkout, subscriptions and webhooks verified against the signing secret.',
    tone: 'bg-indigo-50 text-indigo-600 dark:bg-indigo-500/10 dark:text-indigo-400',
  },
  {
    name: 'Sentry',
    short: 'SN',
    description: 'Traces and errors from the API and both front ends, under one release.',
    tone: 'bg-emerald-50 text-emerald-600 dark:bg-emerald-500/10 dark:text-emerald-400',
  },
]

export function IntegrationLogo({ item }: { item: Integration }) {
  if (item.logo) return <>{item.logo}</>
  return (
    <span
      aria-hidden="true"
      className={`flex size-9 items-center justify-center rounded-lg text-sm font-semibold ${
        item.tone ?? 'bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300'
      }`}
    >
      {item.short ?? item.name.slice(0, 2)}
    </span>
  )
}

export default function IntegrationsCardsGrid({
  title = 'Works with what you already run',
  subtitle = 'Every integration here is a first-class part of the framework, not a community wrapper you have to keep alive.',
  integrations = INTEGRATIONS,
}: {
  title?: string
  subtitle?: string
  integrations?: Integration[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <ul
          role="list"
          className="mx-auto mt-16 grid max-w-5xl grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
        >
          {integrations.map((item) => (
            <li
              key={item.name}
              className="relative rounded-xl border border-gray-200 p-6 hover:bg-gray-50 has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:hover:bg-white/5"
            >
              <IntegrationLogo item={item} />
              <h3 className="mt-4 text-sm font-semibold text-gray-900 dark:text-white">
                {/* Stretched hit area, but only the name is the link. */}
                <a href={item.href ?? '#'} className="after:absolute after:inset-0 focus:outline-none">
                  {item.name}
                </a>
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{item.description}</p>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
