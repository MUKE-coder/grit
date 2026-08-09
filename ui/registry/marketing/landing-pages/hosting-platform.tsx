/*
 * A hosting product page: hero, a live deployments table, a feature bento,
 * cumulative pricing tiers and social proof.
 *
 * One file, because the registry installs one file per block.
 *
 * The deployments table is a real <table> with a status column, not a grid of
 * divs. Each status is a word as well as a colour: "error" in red and "ready"
 * in grey are the same cell to anyone who cannot separate the two, and build
 * state is the one thing this screen exists to report. The dot beside each is
 * aria-hidden, because it repeats the word next to it.
 *
 * The pricing tiers inherit. "Everything in Hobby" is a real list item rather
 * than a footnote, so the chain reads in order for someone stepping through
 * the list, and the tier it refers to is named rather than implied by being
 * the column on the left.
 *
 * The avatar stack is decorative and the sentence beside it carries the claim.
 * Five overlapping photographs of nobody in particular say "people use this";
 * "trusted by 27,000 developers" says it in a form that survives not being
 * seen, so the images take an empty alt and the stack is aria-hidden.
 *
 * The star row has a text equivalent. Five glyphs are five glyphs; the rating
 * they encode has to be written down somewhere.
 *
 * Photography is used for the people. The product surfaces are markup: a stock
 * photo cannot be a screenshot of your product.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Features', 'Pricing', 'Contact']

type Status = 'ready' | 'building' | 'error'

const DEPLOYMENTS: { name: string; status: Status; at: string }[] = [
  { name: 'cache-108', status: 'building', at: '2026-07-26T03:25:00Z' },
  { name: 'web-app-473', status: 'building', at: '2026-05-06T11:15:00Z' },
  { name: 'database-384', status: 'ready', at: '2026-08-20T03:54:00Z' },
  { name: 'backend-832', status: 'ready', at: '2026-05-18T01:15:00Z' },
  { name: 'api-service-19', status: 'ready', at: '2026-08-04T11:46:00Z' },
  { name: 'api-service-576', status: 'error', at: '2026-05-10T13:56:00Z' },
  { name: 'frontend-436', status: 'error', at: '2026-08-07T17:26:00Z' },
  { name: 'api-service-28', status: 'building', at: '2026-08-19T18:08:00Z' },
]

const STATUS_STYLE: Record<Status, { badge: string; dot: string; label: string }> = {
  ready: {
    badge: 'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-200',
    dot: 'bg-gray-400',
    label: 'Ready',
  },
  building: {
    badge: 'bg-gray-900 text-white dark:bg-white dark:text-gray-900',
    dot: 'bg-amber-400',
    label: 'Building',
  },
  error: {
    badge: 'bg-rose-100 text-rose-800 dark:bg-rose-500/20 dark:text-rose-200',
    dot: 'bg-rose-500',
    label: 'Error',
  },
}

const FEATURES = [
  {
    title: 'One-click deploy',
    body: 'Push to a branch and it is live. No pipeline to write and no runner to keep warm.',
    mock: 'git',
    wide: true,
  },
  {
    title: 'Intuitive workflow',
    body: 'Manage the app without learning a second vocabulary for the same operations.',
    mock: 'dashboard',
    wide: false,
  },
  {
    title: 'Hosted at the edge',
    body: 'Served from the city nearest the request, with the origin only for what has to be dynamic.',
    mock: 'globe',
    wide: false,
  },
  {
    title: 'Generated copy',
    body: 'Running out of words for a page. Draft it here and edit rather than starting from a blank box.',
    mock: 'dashboard',
    wide: true,
  },
]

/* Tiers inherit upward. `inherits` names the tier below rather than leaving
   "everything in the previous plan" to depend on column order. */
const TIERS = [
  {
    name: 'Hobby',
    price: 99,
    blurb: 'For a side project that has started getting traffic.',
    features: [
      'Access to basic analytics reports',
      'Up to 10,000 data points per month',
      'Email support',
      'Community forum access',
      'Cancel any time',
    ],
    inherits: null as string | null,
    featured: false,
  },
  {
    name: 'Starter',
    price: 299,
    blurb: 'For a team shipping every week.',
    features: [
      'Advanced analytics dashboard',
      'Customisable reports and charts',
      'Real-time data tracking',
      'Integration with third-party tools',
    ],
    inherits: 'Hobby',
    featured: true,
  },
  {
    name: 'Pro',
    price: 1490,
    blurb: 'For a product the business depends on.',
    features: [
      'Unlimited data storage',
      'Customisable dashboards',
      'Advanced data segmentation',
      'Real-time data processing',
      'Insights and recommendations',
    ],
    inherits: 'Starter',
    featured: false,
  },
]

/* Verified before use: each of these loaded on a contact sheet. */
const AVATARS = [
  '1494790108377-be9c29b29330',
  '1500648767791-00dcc994a43e',
  '1438761681033-6461ffad8d80',
  '1507003211169-0a1dd7228f2d',
  '1544005313-94ddf0286df2',
]

const dates = new Intl.DateTimeFormat('en-GB', {
  day: '2-digit',
  month: 'short',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
})

/* ── Mocks ──────────────────────────────────────────────────────────────── */

function Panel({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div
      aria-hidden="true"
      className={`overflow-hidden rounded-xl border border-gray-200 bg-white select-none dark:border-white/10 dark:bg-gray-900 ${className}`}
    >
      {children}
    </div>
  )
}

function GitMock() {
  return (
    <Panel className="p-4 font-mono text-[11px] text-gray-600 dark:text-gray-300">
      <p>$ git add .</p>
      <p>$ git commit -m &quot;update&quot;</p>
      <p>$ git push</p>
      <p className="mt-2 text-emerald-600 dark:text-emerald-400">deployed in 8s</p>
    </Panel>
  )
}

function DashboardMock() {
  return (
    <Panel className="p-4">
      <div className="flex items-end gap-1.5">
        {[45, 62, 50, 78, 58, 88, 70, 95].map((h, i) => (
          <span
            key={i}
            style={{ height: `${h * 0.4}px` }}
            className="flex-1 rounded-t bg-gradient-to-t from-blue-600 to-blue-400"
          />
        ))}
      </div>
      <div className="mt-3 flex gap-2">
        <span className="h-2 w-20 rounded bg-gray-200 dark:bg-white/10" />
        <span className="h-2 w-12 rounded bg-gray-200 dark:bg-white/10" />
      </div>
    </Panel>
  )
}

function GlobeMock() {
  return (
    <div
      aria-hidden="true"
      className="flex h-32 items-center justify-center rounded-xl bg-gray-950 select-none"
    >
      <span className="size-24 rounded-full bg-gradient-to-br from-blue-500/60 via-blue-700/30 to-transparent" />
    </div>
  )
}

const MOCKS: Record<string, ReactNode> = {
  git: <GitMock />,
  dashboard: <DashboardMock />,
  globe: <GlobeMock />,
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageHostingPlatform({
  brand = 'Unhost',
  title = 'Deploy your website in seconds, not hours',
  subtitle = 'Push to a branch and it is live. No pipeline to write, no runner to keep warm, and a rollback that is one click rather than an incident.',
  rating = 4.9,
  developers = '27,000',
}: {
  brand?: string
  title?: string
  subtitle?: string
  rating?: number
  developers?: string
}) {
  return (
    <div className="bg-white dark:bg-gray-950">
      <header className="border-b border-gray-200 dark:border-white/10">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-6xl items-center justify-between gap-6 px-6 py-4"
        >
          <a
            href="#"
            className="text-lg font-semibold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-white"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-gray-300 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <div className="flex items-center gap-2">
            <a
              href="#"
              className="hidden min-h-11 items-center px-3 text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 sm:inline-flex dark:text-gray-300 dark:hover:text-white"
            >
              Log in
            </a>
            <a
              href="#pricing"
              className="inline-flex min-h-11 items-center rounded-lg border border-gray-300 px-4 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
            >
              Book a call
            </a>
          </div>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-2xl text-center">
              <h1 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
                {title}
              </h1>
              <p className="mx-auto mt-6 text-pretty text-gray-600 dark:text-gray-300">{subtitle}</p>
              <div className="mt-8 flex flex-wrap justify-center gap-3">
                <a
                  href="#pricing"
                  className="inline-flex min-h-12 items-center rounded-lg bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:bg-white dark:text-gray-900"
                >
                  Create account
                </a>
                <a
                  href="#"
                  className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 bg-white px-6 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:border-white/15 dark:bg-transparent dark:text-white dark:hover:bg-white/5"
                >
                  Book a call
                </a>
              </div>
            </div>

            {/* Deployments */}
            <div className="mt-14 overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
              <div className="border-b border-gray-200 px-5 py-4 dark:border-white/10">
                <h2 className="text-sm font-semibold text-gray-900 dark:text-white">
                  Recent deployments
                </h2>
              </div>

              <div className="overflow-x-auto">
                {/* A real table. Build state is the point of this screen, so it
                    is a status column rather than a coloured pill in a div. */}
                <table className="w-full text-left text-sm">
                  <caption className="sr-only">
                    The eight most recent deployments, with build status and time
                  </caption>
                  <thead className="border-b border-gray-200 dark:border-white/10">
                    <tr>
                      <th scope="col" className="px-5 py-3 font-medium text-gray-600 dark:text-gray-300">
                        Name
                      </th>
                      <th scope="col" className="px-5 py-3 font-medium text-gray-600 dark:text-gray-300">
                        Status
                      </th>
                      <th scope="col" className="px-5 py-3 font-medium text-gray-600 dark:text-gray-300">
                        Last deployed
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-gray-100 dark:divide-white/5">
                    {DEPLOYMENTS.map((deployment) => {
                      const style = STATUS_STYLE[deployment.status]
                      return (
                        <tr key={deployment.name}>
                          <th
                            scope="row"
                            className="px-5 py-3 font-mono text-xs font-normal text-gray-900 dark:text-white"
                          >
                            {deployment.name}
                          </th>
                          <td className="px-5 py-3">
                            <span
                              className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${style.badge}`}
                            >
                              {/* Repeats the word beside it, so it is
                                  decoration rather than the only signal. */}
                              <span aria-hidden="true" className={`size-1.5 rounded-full ${style.dot}`} />
                              {style.label}
                            </span>
                          </td>
                          <td className="px-5 py-3 text-gray-600 dark:text-gray-300">
                            <time dateTime={deployment.at}>
                              {dates.format(new Date(deployment.at))}
                            </time>
                          </td>
                        </tr>
                      )
                    })}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </section>

        {/* Features */}
        <section aria-labelledby="features" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="features"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Deployments made easy
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Deploy with ease. Leave the complexity to us.
              </p>
            </div>

            {/* Three columns with spans of 2 + 1 + 1 + 2, which tiles exactly
                over two rows. Spans that do not tile leave CSS grid dropping a
                card into whatever hole it finds. */}
            <ul role="list" className="mt-12 grid gap-6 lg:grid-cols-3">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`flex flex-col rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40 ${
                    feature.wide ? 'lg:col-span-2' : ''
                  }`}
                >
                  <h3 className="font-semibold text-gray-900 dark:text-white">{feature.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{feature.body}</p>
                  <div className="mt-auto pt-6">{MOCKS[feature.mock]}</div>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="pricing-heading"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Simple pricing for advanced people
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Three sizes. Each one includes everything below it.
              </p>
            </div>

            <ul role="list" className="mt-12 grid items-start gap-6 lg:grid-cols-3">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured
                      ? 'border-blue-600 bg-blue-50/50 dark:border-blue-500/40 dark:bg-blue-500/10'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="font-semibold text-gray-900 dark:text-white">{tier.name}</h3>
                    {tier.featured && (
                      /* In words: colour alone reports nothing. */
                      <span className="rounded-full bg-blue-600 px-2.5 py-1 text-xs font-medium text-white">
                        Featured
                      </span>
                    )}
                  </div>

                  <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">{tier.blurb}</p>

                  <p className="mt-5">
                    <span className="align-super text-lg text-gray-500 dark:text-gray-400">$</span>
                    <span className="text-4xl font-bold tracking-tight tabular-nums text-gray-900 dark:text-white">
                      {tier.price}
                    </span>
                    <span className="text-sm text-gray-500 dark:text-gray-400">/month</span>
                  </p>

                  <a
                    href="#"
                    className={`mt-6 inline-flex min-h-12 w-full items-center justify-center rounded-lg text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 ${
                      tier.featured
                        ? 'bg-blue-600 text-white hover:bg-blue-700'
                        : 'bg-gray-900 text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900'
                    }`}
                  >
                    Get {tier.name}
                  </a>

                  <ul role="list" className="mt-6 space-y-2.5">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5 text-sm">
                        <span aria-hidden="true" className="text-blue-600 dark:text-blue-400">
                          ✓
                        </span>
                        <span className="text-gray-700 dark:text-gray-200">{feature}</span>
                      </li>
                    ))}

                    {/* A list item, and it names the tier. "Everything in the
                        previous plan" depends on column order, which is not
                        something a screen reader conveys. */}
                    {tier.inherits && (
                      <li className="flex gap-2.5 border-t border-gray-200 pt-3 text-sm dark:border-white/10">
                        <span aria-hidden="true" className="text-blue-600 dark:text-blue-400">
                          ✓
                        </span>
                        <span className="font-medium text-gray-900 dark:text-white">
                          Everything in {tier.inherits}
                        </span>
                      </li>
                    )}
                  </ul>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Social proof */}
        <section aria-labelledby="proof" className="bg-gray-50 dark:bg-gray-900/40">
          <div className="mx-auto grid max-w-6xl gap-10 px-6 py-20 lg:grid-cols-2 lg:items-center">
            <div>
              <h2
                id="proof"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Host your websites with zero friction today
              </h2>
              <p className="mt-4 text-gray-600 dark:text-gray-300">
                Lightning-fast hosting with unparalleled reliability. Our infrastructure keeps your
                site online around the clock with 99.9% uptime guaranteed.
              </p>

              <div className="mt-8 flex flex-wrap items-center gap-4">
                {/* Decorative: the sentence below carries the claim, and five
                    photographs of nobody in particular carry none of it. */}
                <div aria-hidden="true" className="flex -space-x-3">
                  {AVATARS.map((id) => (
                    <img
                      key={id}
                      src={`https://images.unsplash.com/photo-${id}?w=80&h=80&fit=crop&q=80`}
                      alt=""
                      loading="lazy"
                      className="size-10 rounded-full object-cover ring-2 ring-gray-50 dark:ring-gray-900"
                    />
                  ))}
                </div>

                <p className="flex items-center gap-2">
                  <span aria-hidden="true" className="flex text-amber-400">
                    {'★★★★★'}
                  </span>
                  {/* The glyphs are five glyphs. The rating has to be written
                      down somewhere. */}
                  <span className="text-sm text-gray-600 dark:text-gray-300">
                    {rating} out of 5
                  </span>
                </p>
              </div>

              <p className="mt-3 text-sm text-gray-500 dark:text-gray-400">
                Trusted by {developers} developers
              </p>
            </div>

            <div className="lg:justify-self-end">
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-lg bg-blue-600 px-6 text-sm font-medium text-white hover:bg-blue-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
              >
                Book a call
              </a>
            </div>
          </div>
        </section>
      </main>

      <footer className="relative overflow-hidden border-t border-gray-200 dark:border-white/10">
        <div className="mx-auto grid max-w-6xl gap-8 px-6 py-12 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <p className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">
              {brand}
            </p>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Deploy in seconds, not hours.
            </p>
          </div>

          {[
            { heading: 'Pages', items: ['Home', 'Features', 'Pricing', 'Contact'] },
            { heading: 'Legal', items: ['Privacy policy', 'Terms of service', 'Cookie policy'] },
            { heading: 'Register', items: ['Sign up', 'Log in', 'Book a demo'] },
          ].map((column) => (
            /* Each column is its own labelled nav, so it appears in landmark
               navigation by name rather than as one anonymous block. */
            <nav key={column.heading} aria-label={column.heading}>
              <h2 className="text-sm font-semibold text-gray-900 dark:text-white">
                {column.heading}
              </h2>
              <ul role="list" className="mt-2">
                {column.items.map((item) => (
                  <li key={item}>
                    <a
                      href="#"
                      className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-gray-400 dark:hover:text-white"
                    >
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          ))}
        </div>

        {/* Texture. Read aloud it is the brand said twice. */}
        <span
          aria-hidden="true"
          className="pointer-events-none block text-center text-[5rem] leading-none font-bold tracking-tight text-gray-100 select-none sm:text-[8rem] dark:text-white/5"
        >
          {brand}
        </span>
      </footer>
    </div>
  )
}
