/*
 * The same six integrations in one frame, divided rather than boxed.
 *
 * Six separate cards read as six things; one divided frame reads as one set.
 * Use this where the integrations are a capability of the product, and the
 * card version where each one is a destination you want clicked.
 *
 * The dividers are on the grid, so neighbours share a single line instead of
 * two borders meeting at slightly different opacities. Below `sm` only the
 * horizontal rules survive, because a vertical rule between stacked rows
 * points nowhere.
 *
 * Logos are lettered tiles with a `logo` slot, for the reasons in the sibling
 * cards block: a component library cannot bundle other companies' trademarks
 * on your behalf, and an honest placeholder beats one you forget to replace.
 */

import type { ReactNode } from 'react'

export interface Integration {
  name: string
  /** Two-character label for the fallback tile. Defaults to the first two
   *  letters of the name, which collides for names like Redis and Resend. */
  short?: string
  description: string
  logo?: ReactNode
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

export default function IntegrationsFramedGrid({
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

        {/* One line between neighbours, and only where there are neighbours. */}
        <ul
          role="list"
          className="mx-auto mt-16 grid max-w-5xl grid-cols-1 divide-y divide-gray-200 rounded-2xl border border-gray-200 sm:grid-cols-2 sm:divide-x lg:grid-cols-3 dark:divide-white/10 dark:border-white/10"
        >
          {integrations.map((item) => (
            <li key={item.name} className="p-8">
              {item.logo ?? (
                <span
                  aria-hidden="true"
                  className={`flex size-9 items-center justify-center rounded-lg text-sm font-semibold ${
                    item.tone ?? 'bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-300'
                  }`}
                >
                  {item.short ?? item.name.slice(0, 2)}
                </span>
              )}
              <h3 className="mt-4 text-sm font-semibold text-gray-900 dark:text-white">
                {item.name}
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{item.description}</p>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
