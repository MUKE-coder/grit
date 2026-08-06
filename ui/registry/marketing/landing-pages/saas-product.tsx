/*
 * A complete landing page: nav, hero, logo cloud, features, bento, product
 * shot, testimonial, CTA, footer.
 *
 * This is one file because the registry installs one file per block. That is a
 * constraint worth knowing before you install it: you get a page you own
 * outright with nothing to import, and the first thing you should do is split
 * it into the sections you actually want. It is a starting point, not a
 * component to configure.
 *
 * The section order is the argument the page makes, and it is deliberate:
 * claim, proof that other people believed it, what it does, how it looks, what
 * it costs you to try. Moving the testimonial above the features breaks that,
 * because a stranger vouching for something the reader cannot yet picture is
 * noise.
 *
 * One <h1>, then <h2> per section, then <h3> inside them. Screen reader users
 * navigate long pages by heading, and a page of <div class="text-4xl"> is a
 * page they cannot navigate at all.
 *
 * Every product mock is markup and `aria-hidden`. They illustrate the sentence
 * next to them and hold invented data.
 */

import type { ReactNode } from 'react'

const NAV = ['Product', 'Features', 'Pricing', 'Docs']

const LOGOS = ['Northwind', 'Acme', 'Globex', 'Initech', 'Umbrella', 'Soylent']

const FEATURES = [
  {
    title: 'Generated, not scaffolded',
    body: 'One definition produces the model, the migration, the handler, the validation schema and the typed client. Change a field and every layer follows.',
  },
  {
    title: 'Sessions you can revoke',
    body: 'A row per device, rotation with replay detection, and a sign-out-everywhere that actually signs you out everywhere.',
  },
  {
    title: 'One binary to deploy',
    body: 'The API, the admin panel and the assets ship as a single artifact. There is no runtime to install on the box.',
  },
]

const BENTO = [
  {
    title: 'Background jobs',
    body: 'Queues, retries and a dashboard that shows what failed and why.',
    span: '',
  },
  {
    title: 'File storage',
    body: 'Uploads, signed URLs and image processing against any S3-compatible store.',
    span: '',
  },
  {
    title: 'Observability',
    body: 'Traces from the browser through the API to the query, under one release.',
    span: '',
  },
  {
    title: 'An admin panel you did not write',
    body: 'Every resource you generate gets a table with sorting, filtering and pagination, plus a form that validates against the same schema the API does.',
    span: 'lg:col-span-2',
  },
  {
    title: 'Typed end to end',
    body: 'A renamed Go field fails the TypeScript build.',
    span: '',
  },
]

const STATS = [
  { value: '1', label: 'Command to scaffold' },
  { value: '90+', label: 'Integrations' },
  { value: '12ms', label: 'Median response' },
]

function Card({ children }: { children: ReactNode }) {
  return (
    <div className="rounded-lg border border-gray-200 bg-white p-3 shadow-sm dark:border-white/10 dark:bg-gray-900">
      {children}
    </div>
  )
}

/* The dashboard under the hero. Markup, not a screenshot: it stays sharp,
   weighs nothing, and does not go stale when a nav item is renamed. */
function DashboardMock() {
  const rows = [
    { team: 'Planetaria', name: 'ios-app', live: false, meta: 'Initiated 1m 32s ago' },
    { team: 'Planetaria', name: 'mobile-api', live: true, meta: 'Deployed 3m ago · 23s' },
    { team: 'Grit Labs', name: 'gritframework.dev', live: true, meta: 'Deployed 5m ago · 3m 4s' },
    { team: 'Protocol', name: 'relay-service', live: true, meta: 'Deployed 3h ago · 8s' },
  ]
  return (
    <div
      aria-hidden="true"
      className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl select-none dark:border-white/10 dark:bg-gray-900"
    >
      <div className="flex items-center gap-1.5 border-b border-gray-200 px-4 py-3 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2.5 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="grid grid-cols-[minmax(0,12rem)_minmax(0,1fr)]">
        <div className="hidden border-r border-gray-200 p-4 sm:block dark:border-white/10">
          <p className="text-[10px] font-medium tracking-wide text-gray-400 uppercase">
            Navigation
          </p>
          <ul className="mt-3 space-y-2.5">
            {['Projects', 'Deployments', 'Activity', 'Domains', 'Settings'].map((item) => (
              <li key={item} className="flex items-center gap-2.5">
                <span className="size-3.5 rounded-sm bg-gray-100 dark:bg-white/10" />
                <span className="text-xs text-gray-600 dark:text-gray-400">{item}</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="p-4">
          <p className="text-sm font-semibold text-gray-900 dark:text-white">All projects</p>
          <ul className="mt-3 divide-y divide-gray-100 dark:divide-white/5">
            {rows.map((row) => (
              <li key={row.name} className="py-3">
                <div className="flex items-center gap-2">
                  <span
                    className={`size-2 flex-none rounded-full ${
                      row.live ? 'bg-emerald-400' : 'bg-gray-300 dark:bg-gray-600'
                    }`}
                  />
                  <span className="truncate text-xs font-semibold text-gray-900 dark:text-white">
                    {row.team}
                  </span>
                  <span className="text-xs text-gray-300 dark:text-gray-600">/</span>
                  <span className="truncate text-xs font-semibold text-gray-900 dark:text-white">
                    {row.name}
                  </span>
                </div>
                <p className="mt-1 pl-4 text-[11px] text-gray-500 dark:text-gray-400">{row.meta}</p>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  )
}

function TerminalMock() {
  const lines = [
    { text: '$ grit new storefront --api', tone: 'text-gray-900 dark:text-white' },
    { text: '$ grit generate resource Post title:string body:text', tone: 'text-gray-900 dark:text-white' },
    { text: '  model, migration, handler, schema, client ... done', tone: 'text-gray-500 dark:text-gray-400' },
    { text: '$ grit dev', tone: 'text-gray-900 dark:text-white' },
    { text: '  api    http://localhost:8080', tone: 'text-gray-500 dark:text-gray-400' },
    { text: '  admin  http://localhost:3001', tone: 'text-emerald-600 dark:text-emerald-400' },
  ]
  return (
    <div
      aria-hidden="true"
      className="overflow-hidden rounded-xl border border-gray-200 bg-white select-none dark:border-white/10 dark:bg-gray-900"
    >
      <div className="flex gap-1.5 border-b border-gray-200 px-4 py-3 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2.5 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="space-y-1.5 p-4 font-mono text-xs">
        {lines.map((line) => (
          <p key={line.text} className={line.tone}>
            {line.text}
          </p>
        ))}
      </div>
    </div>
  )
}

export default function LandingPageSaasProduct({
  brand = 'Grit',
  title = 'The backend, the admin panel and the client, from one definition',
  subtitle = 'Describe a resource once. Get the Go API, the migration, the validation schema, the typed client and a working admin screen, and deploy the whole thing as a single binary.',
  primaryLabel = 'Get started',
  secondaryLabel = 'Read the docs',
}: {
  brand?: string
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
}) {
  return (
    <div className="bg-white dark:bg-gray-950">
      <header className="border-b border-gray-200 dark:border-white/10">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4 lg:px-8"
        >
          <a href="#" className="flex items-center gap-2">
            <span
              aria-hidden="true"
              className="flex size-8 items-center justify-center rounded-lg bg-indigo-600 text-sm font-bold text-white"
            >
              {brand[0]}
            </span>
            <span className="text-sm font-semibold text-gray-900 dark:text-white">{brand}</span>
          </a>
          <ul role="list" className="hidden gap-8 lg:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="text-sm font-medium text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <a
            href="#"
            className="inline-flex min-h-11 items-center rounded-lg bg-gray-900 px-4 text-sm font-semibold text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
          >
            Sign in
          </a>
        </nav>
      </header>

      <main>
        {/* Hero — the only h1 on the page. */}
        <section className="relative isolate overflow-hidden px-6 pt-20 pb-16 sm:pt-28 lg:px-8">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-x-0 top-0 -z-10 h-[40rem]"
            style={{
              background:
                'radial-gradient(48rem 26rem at 50% -4rem, rgba(99,102,241,0.16), transparent 68%)',
            }}
          />
          <div className="mx-auto max-w-3xl text-center">
            <h1 className="text-5xl font-semibold tracking-tight text-balance text-gray-900 sm:text-6xl dark:text-white">
              {title}
            </h1>
            <p className="mx-auto mt-6 max-w-2xl text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
              {subtitle}
            </p>
            <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-4">
              <a
                href="#"
                className="inline-flex min-h-11 items-center rounded-lg bg-indigo-600 px-5 text-sm font-semibold text-white hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                {primaryLabel}
              </a>
              <a
                href="#"
                className="inline-flex min-h-11 items-center rounded-lg px-1 text-sm font-semibold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-white"
              >
                {secondaryLabel}
                <span aria-hidden="true">&nbsp;&rarr;</span>
              </a>
            </div>
          </div>
          <div className="mx-auto mt-16 max-w-5xl">
            <DashboardMock />
          </div>
        </section>

        {/* Proof, before the explanation. */}
        <section aria-labelledby="customers" className="px-6 py-16 lg:px-8">
          <div className="mx-auto max-w-7xl">
            <h2 id="customers" className="text-center text-sm text-gray-600 dark:text-gray-400">
              Trusted by teams shipping every week
            </h2>
            <ul
              role="list"
              className="mt-8 flex flex-wrap items-center justify-center gap-x-10 gap-y-6"
            >
              {LOGOS.map((logo) => (
                <li
                  key={logo}
                  className="text-lg font-semibold text-gray-400 dark:text-gray-600"
                >
                  {logo}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* What it does. */}
        <section aria-labelledby="features" className="px-6 py-24 sm:py-32 lg:px-8">
          <div className="mx-auto max-w-7xl">
            <div className="max-w-2xl">
              <p className="text-base font-semibold text-indigo-600 dark:text-indigo-400">
                Everything included
              </p>
              <h2
                id="features"
                className="mt-2 text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white"
              >
                The parts you always end up building anyway
              </h2>
            </div>

            <div className="mt-16 grid grid-cols-1 items-center gap-x-16 gap-y-12 lg:grid-cols-2">
              <dl className="space-y-10">
                {FEATURES.map((feature) => (
                  <div key={feature.title} className="border-l-2 border-indigo-500 pl-6 dark:border-indigo-400">
                    <dt className="text-lg font-semibold text-gray-900 dark:text-white">
                      {feature.title}
                    </dt>
                    <dd className="mt-2 text-base/7 text-gray-600 dark:text-gray-400">
                      {feature.body}
                    </dd>
                  </div>
                ))}
              </dl>
              <TerminalMock />
            </div>
          </div>
        </section>

        {/* The wider surface area, as a bento. Spans tile the grid exactly:
            1+1+1 then 2+1 across three columns. */}
        <section aria-labelledby="platform" className="px-6 pb-24 sm:pb-32 lg:px-8">
          <div className="mx-auto max-w-7xl">
            <h2
              id="platform"
              className="max-w-2xl text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white"
            >
              And the parts you would rather not
            </h2>
            <div className="mt-12 grid grid-cols-1 gap-4 lg:grid-cols-3">
              {BENTO.map((tile) => (
                <div
                  key={tile.title}
                  className={`rounded-2xl border border-gray-200 p-6 dark:border-white/10 ${tile.span}`}
                >
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {tile.title}
                  </h3>
                  <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{tile.body}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* Numbers and a voice, together: one is the claim, one is the reason
            to believe it. */}
        <section
          aria-labelledby="results"
          className="border-y border-gray-200 px-6 py-24 sm:py-32 lg:px-8 dark:border-white/10"
        >
          <div className="mx-auto grid max-w-7xl grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-2">
            <div>
              <h2
                id="results"
                className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white"
              >
                Less of your week spent on plumbing
              </h2>
              <dl className="mt-10 flex flex-wrap gap-x-12 gap-y-8">
                {STATS.map((stat) => (
                  <div key={stat.label}>
                    <dd className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                      {stat.value}
                    </dd>
                    <dt className="mt-1 text-sm text-gray-600 dark:text-gray-400">{stat.label}</dt>
                  </div>
                ))}
              </dl>
            </div>

            <figure className="border-l-2 border-indigo-500 pl-8 dark:border-indigo-400">
              <blockquote className="text-lg/8 text-pretty text-gray-700 dark:text-gray-300">
                &ldquo;We replaced four repositories and a fragile codegen script with one
                definition file. The admin panel we had been putting off for a year existed by
                the end of the afternoon.&rdquo;
              </blockquote>
              <figcaption className="mt-6 flex items-center gap-3">
                <span
                  aria-hidden="true"
                  className="flex size-10 items-center justify-center rounded-full bg-gray-100 text-sm font-semibold text-gray-600 dark:bg-white/10 dark:text-gray-300"
                >
                  ML
                </span>
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  <span className="font-medium text-gray-900 dark:text-white">Maria Lukose</span>,
                  engineering lead
                </span>
              </figcaption>
            </figure>
          </div>
        </section>

        {/* The ask, last, when the case has been made. */}
        <section aria-labelledby="cta" className="px-6 py-24 sm:py-32 lg:px-8">
          <div className="mx-auto max-w-2xl text-center">
            <h2
              id="cta"
              className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white"
            >
              Start with one command
            </h2>
            <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
              Scaffold a project, generate a resource, and have something running locally in
              about a minute.
            </p>
            <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-4">
              <a
                href="#"
                className="inline-flex min-h-11 items-center rounded-lg bg-indigo-600 px-5 text-sm font-semibold text-white hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                {primaryLabel}
              </a>
              <a
                href="#"
                className="inline-flex min-h-11 items-center rounded-lg px-1 text-sm font-semibold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-white"
              >
                {secondaryLabel}
                <span aria-hidden="true">&nbsp;&rarr;</span>
              </a>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-gray-200 px-6 py-12 lg:px-8 dark:border-white/10">
        <div className="mx-auto flex max-w-7xl flex-col gap-8 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex items-center gap-2">
            <span
              aria-hidden="true"
              className="flex size-8 items-center justify-center rounded-lg bg-indigo-600 text-sm font-bold text-white"
            >
              {brand[0]}
            </span>
            <span className="text-sm text-gray-600 dark:text-gray-400">
              &copy; {brand}. All rights reserved.
            </span>
          </div>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-8 gap-y-2">
              {['Docs', 'Changelog', 'GitHub', 'Privacy'].map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
                  >
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </footer>
    </div>
  )
}
