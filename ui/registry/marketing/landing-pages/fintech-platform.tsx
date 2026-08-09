/*
 * A payments-infrastructure page: collage hero, integrations split, a product
 * shot with supporting columns, a dense bento and a testimonial carried by
 * numbers.
 *
 * One file, because the registry installs one file per block. Split it into
 * the sections you want.
 *
 * Three things here that the other landing pages in this category do not do,
 * which is the reason it exists alongside them rather than instead of one:
 *
 * The hero is a collage rather than a single screenshot. A cluster of small
 * artefacts — a card, a request, a fragment of dashboard — claims breadth,
 * which is the claim an infrastructure product has to make. One big screenshot
 * claims depth in a single screen and says nothing about the rest.
 *
 * The testimonial carries two metrics beside it. A quote alone is one person's
 * opinion; a quote with a number attached is a result someone is willing to
 * put their name to. The numbers are in the same figure as the quote, so they
 * are attributed to the same source rather than floating free as marketing.
 *
 * The integrations row uses initials tiles, not brand marks. A template that
 * ships real logos ships someone else's trademark, and the tile is the right
 * placeholder because it makes the shape obvious without pretending.
 *
 * Every mock is markup and aria-hidden. They illustrate the sentence beside
 * them and hold invented figures, and read aloud they are digits with no
 * referent. Nothing here loads an external asset.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

/* Verified on a contact sheet at the size they are shown. Photographs of
   people, markup for product surfaces: a drawn avatar is a worse likeness of a
   person than a photograph, and a stock photo is a worse screenshot of your
   product than a drawing of one. Decorative either way — the caption beside
   each one already names the person. */
const FACES = [
  '1500648767791-00dcc994a43e',
  '1494790108377-be9c29b29330',
  '1531427186611-ecfd6d936c79',
  '1580489944761-15a19d654956',
]

const face = (id: string, size: number) =>
  `https://images.unsplash.com/photo-${id}?w=${size * 2}&h=${size * 2}&fit=crop&crop=faces&q=75`

const NAV = ['Product', 'Solutions', 'Pricing', 'Company']

const LOGOS = ['Northwind', 'Meridian', 'Kestrel', 'Lumen', 'Atlas', 'Verdant']

const TOOLS = [
  { name: 'VS Code', tone: 'bg-sky-100 text-sky-700 dark:bg-sky-500/15 dark:text-sky-300' },
  { name: 'Neovim', tone: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300' },
  { name: 'JetBrains', tone: 'bg-rose-100 text-rose-700 dark:bg-rose-500/15 dark:text-rose-300' },
  { name: 'Zed', tone: 'bg-violet-100 text-violet-700 dark:bg-violet-500/15 dark:text-violet-300' },
]

const INTEGRATION_POINTS = [
  { title: 'Typed clients', body: 'Generated from the same spec the API serves, so the client cannot drift from the server.' },
  { title: 'Idempotent by default', body: 'Every write takes a key. Retries are safe without you designing for them.' },
  { title: 'Sandbox that matches', body: 'The test environment runs the same code path, including the failures you need to handle.' },
]

const DIRECTION_POINTS = [
  { title: 'One timeline', body: 'Charges, refunds, disputes and payouts on a single ordered view rather than four exports.' },
  { title: 'Filters that persist', body: 'A view you set is a view you can link to, so a question asked once can be asked again.' },
  { title: 'Exports that reconcile', body: 'CSV totals match the dashboard totals, because both read the same ledger.' },
]

const BENTO = [
  { title: 'Consistent uptime', body: 'An SLA backed by a status page that shows the incident rather than just the colour.', mock: 'uptime' },
  { title: 'Keyboard shortcuts', body: 'Every destination reachable without the mouse, and discoverable without the docs.', mock: 'keys' },
  { title: 'Currency conversion', body: 'Settle in one, charge in thirty, with the rate recorded against each charge.', mock: 'currency' },
  { title: 'Enterprise security', body: 'Role-based access, audit trails and fine-grained controls that record who saw what.', mock: 'shield' },
  { title: 'Resource insight', body: 'Queue depth, error rate and latency at a glance, alerted on before they bite.', mock: 'memory' },
  { title: 'Blazing performance', body: 'Edge-cached reads and a dashboard that stays responsive at a million rows.', mock: 'speed' },
]

const PROOF = {
  quote:
    'We replaced three vendors and a spreadsheet with this. The part that convinced the board was reconciliation: the number in the bank finally matched the number in the books, on the first close.',
  name: 'Méschac Ngandu',
  role: 'VP Engineering, Kestrel Labs',
  stats: [
    { value: '37%', label: 'increase in checkout completion' },
    { value: '90%', label: 'retention on annual plans' },
  ],
}

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

function CardMock() {
  return (
    <div
      aria-hidden="true"
      className="rounded-xl bg-gradient-to-br from-gray-900 to-gray-700 p-4 text-white shadow-lg select-none dark:from-gray-800 dark:to-gray-700"
    >
      <div className="flex items-start justify-between">
        <span className="size-6 rounded bg-white/20" />
        <span className="text-xs font-semibold tracking-widest">VISA</span>
      </div>
      <p className="mt-6 font-mono text-sm tracking-widest">5387 4987 8981 2345</p>
      <div className="mt-3 flex justify-between text-[10px] text-white/70">
        <span>Member since 2019</span>
        <span>12 / 29</span>
      </div>
    </div>
  )
}

function RequestMock() {
  const lines = [
    { text: 'POST /v1/charges', tone: 'text-emerald-600 dark:text-emerald-400' },
    { text: '  amount: 4200,', tone: 'text-gray-600 dark:text-gray-300' },
    { text: '  currency: "usd",', tone: 'text-gray-600 dark:text-gray-300' },
    { text: '  capture: true', tone: 'text-gray-600 dark:text-gray-300' },
    { text: '201 Created', tone: 'text-sky-600 dark:text-sky-400' },
  ]
  return (
    <Panel className="p-3">
      <div className="space-y-1 font-mono text-[11px]">
        {lines.map((line) => (
          <p key={line.text} className={line.tone}>
            {line.text}
          </p>
        ))}
      </div>
    </Panel>
  )
}

function ConnectedMock() {
  return (
    <Panel className="flex items-center gap-2 px-3 py-2.5">
      <span className="size-2 rounded-full bg-emerald-500" />
      <span className="text-[11px] font-medium text-gray-700 dark:text-gray-200">API connected</span>
    </Panel>
  )
}

function StatMock() {
  return (
    <Panel className="p-3">
      <p className="text-[10px] text-gray-500 dark:text-gray-400">Volume today</p>
      <p className="mt-1 text-lg font-semibold tabular-nums text-gray-900 dark:text-white">
        $128,402
      </p>
      <div className="mt-2 flex items-end gap-0.5">
        {[40, 65, 50, 80, 62, 92].map((h, i) => (
          <span key={i} style={{ height: `${h * 0.28}px` }} className="w-2 rounded-sm bg-sky-500" />
        ))}
      </div>
    </Panel>
  )
}

function DashboardMock() {
  return (
    <Panel>
      <div className="flex items-center gap-1.5 border-b border-gray-200 px-3 py-2.5 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="flex">
        <div className="hidden w-36 shrink-0 space-y-2 border-r border-gray-200 p-3 sm:block dark:border-white/10">
          {Array.from({ length: 8 }, (_, i) => (
            <span
              key={i}
              className={`block h-2 rounded ${
                i === 2 ? 'w-4/5 bg-sky-300 dark:bg-sky-500/40' : 'w-full bg-gray-200 dark:bg-white/10'
              }`}
            />
          ))}
        </div>
        <div className="flex-1 space-y-2 p-3">
          {Array.from({ length: 12 }, (_, i) => (
            <div key={i} className="flex items-center gap-2">
              <span className="size-2 shrink-0 rounded-sm bg-gray-200 dark:bg-white/10" />
              <span className="h-2 flex-1 rounded bg-gray-100 dark:bg-white/5" />
              <span
                className={`h-2 w-12 shrink-0 rounded ${
                  i % 5 === 0
                    ? 'bg-emerald-200 dark:bg-emerald-500/30'
                    : 'bg-gray-100 dark:bg-white/5'
                }`}
              />
            </div>
          ))}
        </div>
      </div>
    </Panel>
  )
}

function UptimeMock() {
  return (
    <div aria-hidden="true" className="flex items-end gap-0.5 select-none">
      {Array.from({ length: 26 }, (_, i) => (
        <span
          key={i}
          style={{ height: `${58 + ((i * 41) % 42)}%` }}
          className={`w-1.5 rounded-sm ${i === 17 ? 'bg-amber-400' : 'bg-sky-500'}`}
        />
      ))}
    </div>
  )
}

function KeysMock() {
  return (
    <div aria-hidden="true" className="flex gap-2 select-none">
      {['⌘', 'K'].map((key) => (
        <span
          key={key}
          className="flex size-10 items-center justify-center rounded-lg border border-gray-300 bg-gray-50 font-mono text-sm text-gray-700 shadow-sm dark:border-white/15 dark:bg-white/5 dark:text-gray-200"
        >
          {key}
        </span>
      ))}
    </div>
  )
}

function CurrencyMock() {
  return (
    <div aria-hidden="true" className="space-y-1.5 select-none">
      {[
        ['USD', '1.00'],
        ['EUR', '0.92'],
        ['GBP', '0.79'],
      ].map(([code, rate]) => (
        <div key={code} className="flex justify-between text-[11px]">
          <span className="font-mono text-gray-500 dark:text-gray-400">{code}</span>
          <span className="tabular-nums text-gray-700 dark:text-gray-200">{rate}</span>
        </div>
      ))}
    </div>
  )
}

function ShieldMock() {
  return (
    <div
      aria-hidden="true"
      className="flex size-14 items-center justify-center rounded-xl bg-gradient-to-br from-sky-500 to-sky-700 select-none"
    >
      <span className="text-xl text-white">🛡</span>
    </div>
  )
}

function MemoryMock() {
  return (
    <div aria-hidden="true" className="select-none">
      <p className="text-[11px] text-gray-500 dark:text-gray-400">Queue depth 41 / 5,000</p>
      <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-gray-200 dark:bg-white/10">
        <span className="block h-full w-[9%] rounded-full bg-emerald-500" />
      </div>
    </div>
  )
}

function SpeedMock() {
  return (
    <div
      aria-hidden="true"
      className="flex size-14 items-center justify-center rounded-xl bg-gray-900 select-none dark:bg-white/10"
    >
      <span className="text-xl text-sky-400">⚡</span>
    </div>
  )
}

const MOCKS: Record<string, ReactNode> = {
  uptime: <UptimeMock />,
  keys: <KeysMock />,
  currency: <CurrencyMock />,
  shield: <ShieldMock />,
  memory: <MemoryMock />,
  speed: <SpeedMock />,
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageFintechPlatform({
  /* Not one of LOGOS below: the default brand appearing in its own
     customer logo cloud reads as a company citing itself as a reference. */
  brand = 'Keystone',
  title = 'The financial layer powering businesses on your platform',
  subtitle = 'Take payments, move money and reconcile it, on one API that behaves the same in test as it does in production.',
}: {
  brand?: string
  title?: string
  subtitle?: string
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
            className="text-lg font-semibold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:text-white"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:text-gray-300 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <a
            href="#"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg border border-gray-300 px-4 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
          >
            Sign in
          </a>
        </nav>
      </header>

      <main>
        {/* Hero with collage */}
        <section className="relative overflow-hidden border-b border-gray-200 dark:border-white/10">
          <div
            aria-hidden="true"
            className="absolute inset-0 bg-gradient-to-b from-sky-50 to-white dark:from-sky-500/10 dark:to-gray-950"
          />
          <div className="relative mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-3xl text-center">
              <h1 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
                {title}
              </h1>
              <p className="mx-auto mt-6 max-w-xl text-lg text-pretty text-gray-600 dark:text-gray-300">
                {subtitle}
              </p>
              <div className="mt-8 flex flex-wrap justify-center gap-3">
                <a
                  href="#"
                  className="inline-flex min-h-12 items-center rounded-lg bg-sky-700 px-6 text-sm font-medium text-white hover:bg-sky-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600"
                >
                  Get started
                </a>
                <a
                  href="#integrations"
                  className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
                >
                  Read the API docs
                </a>
              </div>
            </div>

            {/* A cluster of small artefacts rather than one screenshot: breadth
                is the claim an infrastructure product has to make. Laid out on
                a grid rather than absolutely positioned so it reflows on a
                phone instead of overlapping itself. */}
            <div className="mt-16 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              <div className="lg:pt-10">
                <StatMock />
              </div>
              <div className="lg:col-span-2">
                <CardMock />
                <div className="mt-4">
                  <ConnectedMock />
                </div>
              </div>
              <div className="lg:pt-6">
                <RequestMock />
              </div>
            </div>
          </div>
        </section>

        {/* Logos */}
        <section aria-labelledby="logos" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-12">
            <h2 id="logos" className="text-center text-sm text-gray-500 dark:text-gray-400">
              Trusted by fast-growing companies around the world
            </h2>
            <ul
              role="list"
              className="mt-6 grid grid-cols-2 items-center justify-items-center gap-6 sm:grid-cols-3 lg:grid-cols-6"
            >
              {LOGOS.map((logo) => (
                <li
                  key={logo}
                  /* gray-600 on paper, gray-400 on the dark theme. The
                     gray-400/gray-500 pair this started as measured 2.54:1 and
                     4.16:1 — a logo cloud is quiet by design, but quiet is not
                     the same as unreadable. */
                  className="text-sm font-semibold tracking-tight text-gray-600 dark:text-gray-400"
                >
                  {logo}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Integrations split */}
        <section
          id="integrations"
          aria-labelledby="integrations-heading"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto grid max-w-6xl gap-12 px-6 py-20 lg:grid-cols-2 lg:items-center">
            <div>
              <h2
                id="integrations-heading"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Powerful integrations with our API
              </h2>
              <p className="mt-4 text-gray-600 dark:text-gray-300">
                Generate payment links at scale, automate marketing campaigns, speed up processes
                and ship faster, against an API that does not surprise you in production.
              </p>

              <h3 className="mt-8 text-sm font-semibold text-gray-900 dark:text-white">
                Native editor support
              </h3>
              <ul role="list" className="mt-3 flex flex-wrap gap-2">
                {TOOLS.map((tool) => (
                  <li key={tool.name}>
                    {/* Initials tiles, not brand marks. A template shipping
                        real logos ships someone else's trademark, and the name
                        beside it is what carries the meaning anyway. */}
                    <span className="flex items-center gap-2 rounded-lg border border-gray-200 py-1.5 pr-3 pl-1.5 dark:border-white/10">
                      <span
                        aria-hidden="true"
                        className={`flex size-6 items-center justify-center rounded text-[10px] font-bold ${tool.tone}`}
                      >
                        {tool.name.slice(0, 2).toUpperCase()}
                      </span>
                      <span className="text-xs font-medium text-gray-700 dark:text-gray-200">
                        {tool.name}
                      </span>
                    </span>
                  </li>
                ))}
              </ul>
            </div>

            <DashboardMock />
          </div>

          <div className="mx-auto max-w-6xl px-6 pb-20">
            <ul role="list" className="grid gap-8 sm:grid-cols-3">
              {INTEGRATION_POINTS.map((point) => (
                <li key={point.title}>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {point.title}
                  </h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{point.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Product direction */}
        <section
          aria-labelledby="direction"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase dark:text-sky-400">
                Product direction
              </p>
              <h2
                id="direction"
                className="mt-3 text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                See where the money actually went
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Monitor activity in real time, and resolve the question the same day it is asked.
              </p>
            </div>

            <div className="mt-10">
              <DashboardMock />
            </div>

            <ul role="list" className="mt-10 grid gap-8 sm:grid-cols-3">
              {DIRECTION_POINTS.map((point) => (
                <li key={point.title}>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {point.title}
                  </h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{point.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Bento */}
        <section aria-labelledby="blocks" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase dark:text-sky-400">
                What you get
              </p>
              <h2
                id="blocks"
                className="mt-3 text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Built out of smaller, sturdier blocks
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Everything needed to launch and scale, designed for speed, reliability and a
                developer experience that does not fight you.
              </p>
            </div>

            <ul role="list" className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {BENTO.map((tile) => (
                <li
                  key={tile.title}
                  className="flex flex-col rounded-2xl border border-gray-200 bg-white p-6 dark:border-white/10 dark:bg-gray-900"
                >
                  <h3 className="font-semibold text-gray-900 dark:text-white">{tile.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{tile.body}</p>
                  {/* mt-auto so the mocks line up across a row whatever the
                      copy above them does. */}
                  <div className="mt-auto pt-6">{MOCKS[tile.mock]}</div>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Proof: quote carrying numbers */}
        <section aria-labelledby="proof" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20">
            <h2 id="proof" className="sr-only">
              What customers report
            </h2>

            {/* The numbers live inside the same figure as the quote, so they
                are attributed to the person who said them rather than floating
                free as unsourced marketing. */}
            <figure className="rounded-2xl border border-gray-200 bg-gray-50 p-8 sm:p-10 dark:border-white/10 dark:bg-gray-900/40">
              <blockquote className="text-lg leading-relaxed text-pretty text-gray-900 sm:text-xl dark:text-white">
                {PROOF.quote}
              </blockquote>

              <figcaption className="mt-8 flex flex-col gap-8 sm:flex-row sm:items-center sm:justify-between">
                <span className="flex items-center gap-3">
                  <img
                    src={face(FACES[0], 44)}
                    alt=""
                    aria-hidden="true"
                    loading="lazy"
                    className="size-11 rounded-full object-cover"
                  />
                  <span>
                    <span className="block text-sm font-medium text-gray-900 dark:text-white">
                      {PROOF.name}
                    </span>
                    <span className="block text-xs text-gray-500 dark:text-gray-400">
                      {PROOF.role}
                    </span>
                  </span>
                </span>

                <dl className="flex gap-10">
                  {PROOF.stats.map((stat) => (
                    <div key={stat.label}>
                      <dt className="sr-only">{stat.label}</dt>
                      <dd>
                        <span className="block text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
                          {stat.value}
                        </span>
                        <span className="mt-1 block max-w-[10rem] text-xs text-gray-500 dark:text-gray-400">
                          {stat.label}
                        </span>
                      </dd>
                    </div>
                  ))}
                </dl>
              </figcaption>
            </figure>
          </div>
        </section>

        {/* Closing ask */}
        <section aria-labelledby="cta" className="bg-gray-50 dark:bg-gray-900/40">
          <div className="mx-auto max-w-3xl px-6 py-20 text-center">
            <h2
              id="cta"
              className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              Create, sell and grow
            </h2>
            <p className="mx-auto mt-4 max-w-lg text-gray-600 dark:text-gray-300">
              Join a community of over a thousand companies and developers already building on
              {' '}{brand}.
            </p>
            <a
              href="#"
              className="mt-8 inline-flex min-h-12 items-center rounded-lg bg-sky-700 px-6 text-sm font-medium text-white hover:bg-sky-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600"
            >
              Contact sales
            </a>
          </div>
        </section>
      </main>

      <footer className="border-t border-gray-200 dark:border-white/10">
        <div className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">
              {brand}
            </p>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              The financial layer for platforms.
            </p>
          </div>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
              {[...NAV, 'Privacy', 'Terms'].map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:text-gray-400 dark:hover:text-white"
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
