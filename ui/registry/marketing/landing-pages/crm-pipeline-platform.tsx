/*
 * A B2B CRM product page: hero with a product shot, logo cloud, headline
 * figures, feature bento, an accordion beside a live panel, and pricing.
 *
 * One file, because the registry installs one file per block.
 *
 * ── The orange ──────────────────────────────────────────────────────────────
 * Orange is the hardest accent to use honestly, because the shade everyone
 * reaches for does not pass. White on orange-500 measures 2.80:1 and on
 * orange-600 3.56:1, both under the 4.5:1 AA wants for a label. orange-700 is
 * 5.18:1. So every filled button and every solid badge here is orange-700,
 * and the bright orange-500 survives only where it is decoration: a chart
 * bar, a ring, a dot, a gradient. The hero's accented word is orange-700 in
 * light mode and orange-400 in dark, since orange-700 on near-black is 3.1:1.
 *
 * ── The product shot ────────────────────────────────────────────────────────
 * Markup, not a screenshot, so it stays sharp and themable. It is also a
 * drawing of an app that does not exist, which makes it decoration: the whole
 * subtree is aria-hidden. Otherwise a screen reader reads out "Pipeline
 * Activity, 274, Sunday, Monday" as though those were facts on the page, and
 * the deliberately quiet styling becomes a contrast failure precisely because
 * it is announced as text. Nothing inside it is focusable, which is what makes
 * hiding it safe rather than a new bug.
 *
 * The floating name tags around the headline are the same: decoration that
 * says "people collaborate here", not information.
 *
 * ── The figures ─────────────────────────────────────────────────────────────
 * A <dl>. "90%" beside "Productivity growth" is a term and its definition, and
 * that pairing should be announced rather than inferred from two boxes sitting
 * near each other.
 *
 * ── The accordion ───────────────────────────────────────────────────────────
 * Native <details>/<summary>, so it works before hydration and keyboard
 * support is the browser's problem rather than ours. One is open by default
 * because an accordion where everything is shut looks broken.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Features', 'Pricing', 'Blogs', 'Resources', 'Contact']

/* Eight, so the four-column grid fills exactly. Nine leaves a dead cell that
   reads as a missing logo rather than a deliberate gap. */
const LOGOS = [
  'Powersurge',
  'Foresight',
  'Stacked Lab',
  'CoreOS',
  'Wildcrafted',
  'Sonorous',
  'Quixotic',
  'Luminary',
]

const FIGURES = [
  { value: '90%', label: 'Productivity growth' },
  { value: '270+', label: 'Teams trust Flowen' },
  { value: '80%', label: 'Boost in revenue' },
  { value: '12K+', label: 'Deals managed weekly' },
]

const FEATURES = [
  {
    title: 'Know your best-performing regions',
    body: 'Analyse sales performance across countries to identify high-performing regions and new opportunities.',
    mock: 'regions' as const,
  },
  {
    title: 'Monitor sales growth in real time',
    body: 'View customer orders, activity, and growth metrics in one simple and intuitive dashboard.',
    mock: 'growth' as const,
  },
  {
    title: 'Real-time sales analytics',
    body: 'Stay on top of revenue trends with powerful analytics that show your team a real-time performance.',
    mock: 'profit' as const,
  },
  {
    title: 'Simplify your sales operations',
    body: 'Track deals, manage tasks, and collaborate with your team, all from a single organised CRM dashboard.',
    mock: 'orders' as const,
  },
]

const CAPABILITIES = [
  {
    title: 'Proven sales performance',
    body: 'Every number on this page comes out of the same pipeline your reps already work in, so the report and the reality cannot drift apart.',
  },
  {
    title: 'Intelligent sales automation',
    body: 'Our automation engine handles the repetitive work, assigns leads on your rules, and keeps the pipeline moving without manual effort.',
    open: true,
  },
  {
    title: 'Seamless CRM integrations',
    body: 'Connect the mailbox, the calendar and the billing system you already pay for, in minutes, without a data migration project.',
  },
  {
    title: 'Real-time pipeline insights',
    body: 'Stage-by-stage value, ageing deals and win rate, recalculated as the deal moves rather than overnight.',
  },
  {
    title: 'Collaborative sales workspace',
    body: 'Notes, tasks and files live on the deal, so the person covering next week is not starting from an empty inbox.',
  },
]

const PLANS = [
  {
    name: 'Free Plan',
    price: '$0',
    blurb: 'Perfect for individuals or small teams just getting started.',
    cta: 'Start free trial',
    featured: false,
    includes: 'Includes:',
    features: [
      'Email support',
      '1 team member',
      'Standard reports',
      'Basic analytics dashboard',
      'Product and customer tracking',
    ],
  },
  {
    name: 'Growth Plan',
    price: '$49',
    blurb: 'Best for growing businesses that need deeper insights.',
    cta: 'Get started',
    featured: true,
    includes: 'Everything in Free, plus:',
    features: [
      'Priority support',
      'Messaging inbox',
      'Team collaboration tools',
      'Advanced filters and segmentation',
      'Automated workflow triggers',
    ],
  },
  {
    name: 'Business Plan',
    price: '$99',
    blurb: 'Ideal for scaling companies that need full control and automation.',
    cta: 'Get started',
    featured: false,
    includes: 'Everything in Growth, plus:',
    features: [
      'Custom dashboards',
      'Automated reporting',
      'Unlimited team members',
      'Dedicated account manager',
      'Role-based access controls',
    ],
  },
]

/* Weekday bar heights for the hero chart, as percentages. Fixed values rather
   than random ones: a mock that reshuffles on every render is a mock that
   cannot be reviewed. */
const ACTIVITY = [58, 74, 96, 62, 88, 47]
const DAYS = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri']

function Tag({ name, className, tone }: { name: string; className: string; tone: string }) {
  return (
    <span className={`absolute hidden items-center gap-1.5 lg:inline-flex ${className}`}>
      <span className={`rounded-md px-2 py-1 text-xs font-medium text-white ${tone}`}>{name}</span>
    </span>
  )
}

function Eyebrow({ children }: { children: ReactNode }) {
  return (
    <span className="inline-flex items-center gap-2 rounded-full border border-orange-200 bg-orange-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-orange-800 dark:border-orange-400/25 dark:bg-orange-400/10 dark:text-orange-300">
      <span aria-hidden="true" className="size-1.5 rounded-full bg-orange-600" />
      {children}
    </span>
  )
}

function Card({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div
      className={`rounded-xl border border-gray-200 bg-white p-4 dark:border-white/10 dark:bg-gray-900 ${className}`}
    >
      {children}
    </div>
  )
}

/* Every mock below is a drawing. The section that renders them is aria-hidden
   as a whole, so none of this text is announced. */
function FeatureMock({ kind }: { kind: (typeof FEATURES)[number]['mock'] }) {
  if (kind === 'regions') {
    return (
      <div className="space-y-3">
        <div className="flex items-baseline justify-between">
          <p className="text-xs font-medium text-gray-900 dark:text-white">Sales by countries</p>
          <p className="text-[10px] text-gray-500 dark:text-gray-400">All products</p>
        </div>
        <div className="grid grid-cols-3 gap-3">
          <div className="col-span-1 space-y-2">
            <div>
              <p className="text-[10px] text-gray-500 dark:text-gray-400">Top performing</p>
              <p className="text-sm font-semibold text-gray-900 dark:text-white">$120,000</p>
            </div>
            <div>
              <p className="text-[10px] text-gray-500 dark:text-gray-400">Revenue growth</p>
              <p className="text-sm font-semibold text-orange-700 dark:text-orange-400">+34%</p>
            </div>
          </div>
          {/* A stylised map: dots on a grid, not a projection anybody should
              read a country off. */}
          <div className="col-span-2 grid grid-cols-12 gap-1 rounded-lg bg-gray-100 p-2 dark:bg-white/5">
            {Array.from({ length: 48 }).map((_, index) => (
              <span
                key={index}
                className={`size-1 rounded-full ${
                  [4, 6, 15, 17, 18, 26, 29, 33, 38, 41].includes(index)
                    ? 'bg-orange-500'
                    : 'bg-gray-300 dark:bg-white/15'
                }`}
              />
            ))}
          </div>
        </div>
      </div>
    )
  }

  if (kind === 'growth') {
    return (
      <div className="space-y-3">
        <p className="text-xs font-medium text-gray-900 dark:text-white">Customer orders</p>
        <div className="flex items-end gap-2">
          <p className="text-2xl font-semibold text-gray-900 dark:text-white">45,637</p>
          <span className="rounded-md bg-emerald-100 px-1.5 py-0.5 text-[10px] font-semibold text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300">
            +9.4%
          </span>
        </div>
        <svg viewBox="0 0 200 48" className="h-14 w-full" preserveAspectRatio="none">
          <path
            d="M0 38 L28 30 L56 34 L84 22 L112 26 L140 12 L168 18 L200 6"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            className="text-orange-500"
          />
        </svg>
        <div className="flex justify-between text-[10px] text-gray-500 dark:text-gray-400">
          <span>May</span>
          <span>Jul</span>
          <span>Sep</span>
          <span>Nov</span>
        </div>
      </div>
    )
  }

  if (kind === 'profit') {
    return (
      <div className="space-y-3">
        <p className="text-xs font-medium text-gray-900 dark:text-white">Total profit overview</p>
        <p className="text-xl font-semibold text-gray-900 dark:text-white">$98,643.24</p>
        <div className="flex h-16 items-end gap-1.5">
          {[45, 62, 38, 74, 52, 90, 58, 41, 66, 49].map((height, index) => (
            <span
              key={index}
              style={{ height: `${height}%` }}
              className={`flex-1 rounded-sm ${index === 5 ? 'bg-orange-500' : 'bg-orange-200 dark:bg-orange-400/25'}`}
            />
          ))}
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-2">
      <p className="text-xs font-medium text-gray-900 dark:text-white">Recent orders</p>
      {[
        ['#1248', 'Smartwatch Pro X', 'Shakil Ahmed'],
        ['#1246', 'Office Chair Deluxe', 'Arafat Hossain'],
        ['#1247', 'Laptop Stand Ergonomic', 'Nazrul Jahan'],
      ].map(([id, product, customer]) => (
        <div
          key={id}
          className="flex items-center justify-between rounded-lg bg-gray-50 px-2 py-1.5 text-[10px] dark:bg-white/5"
        >
          <span className="font-medium text-gray-900 dark:text-white">{id}</span>
          <span className="flex-1 truncate px-2 text-gray-600 dark:text-gray-300">{product}</span>
          <span className="text-gray-500 dark:text-gray-400">{customer}</span>
        </div>
      ))}
    </div>
  )
}

export default function CrmPipelinePlatform() {
  return (
    <div className="bg-white text-gray-900 dark:bg-gray-950 dark:text-white">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-50 focus:rounded-lg focus:bg-orange-700 focus:px-4 focus:py-2 focus:text-sm focus:font-semibold focus:text-white"
      >
        Skip to content
      </a>

      {/* ── Header ───────────────────────────────────────────────────────── */}
      <header className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex max-w-6xl items-center justify-between gap-6 px-6 py-4">
          <a href="#" className="flex items-center gap-2 text-lg font-semibold">
            <span
              aria-hidden="true"
              className="grid size-8 place-items-center rounded-lg bg-orange-700 text-sm font-bold text-white"
            >
              F
            </span>
            Flowen
          </a>

          <nav aria-label="Main" className="hidden lg:block">
            <ul className="flex items-center gap-7 text-sm">
              {NAV.map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-orange-700 dark:text-gray-300 dark:hover:text-white"
                  >
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>

          <div className="flex items-center gap-2">
            <a
              href="#"
              className="hidden min-h-11 items-center rounded-lg px-3 text-sm font-medium text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 sm:inline-flex dark:text-gray-200 dark:hover:bg-white/10"
            >
              Sign in
            </a>
            <a
              href="#"
              className="inline-flex min-h-11 items-center rounded-lg bg-orange-700 px-4 text-sm font-semibold text-white hover:bg-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
            >
              Start for free
            </a>
          </div>
        </div>
      </header>

      <main id="main">
        {/* ── Hero ───────────────────────────────────────────────────────── */}
        <section className="relative overflow-hidden px-6 pb-16 pt-14 text-center">
          <div className="relative mx-auto max-w-3xl">
            <p className="mb-6">
              <a
                href="#"
                className="inline-flex items-center gap-2 rounded-full border border-gray-200 bg-white py-1 pl-1 pr-3 text-sm text-gray-700 hover:border-gray-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:bg-white/5 dark:text-gray-200"
              >
                <span className="rounded-full bg-orange-700 px-2 py-0.5 text-xs font-semibold text-white">
                  New
                </span>
                Trusted by 999+ growing B2B teams
                <span aria-hidden="true">&rsaquo;</span>
              </a>
            </p>

            {/* Decoration. Four names floating near a headline are a mood, not
                a fact, so they are hidden rather than read out mid-sentence. */}
            <div aria-hidden="true">
              <Tag name="Jenny" className="-left-4 top-24 xl:-left-16" tone="bg-fuchsia-600" />
              <Tag name="Emma" className="-left-2 top-52 xl:-left-10" tone="bg-violet-600" />
              <Tag name="Conner" className="-right-4 top-24 xl:-right-16" tone="bg-emerald-700" />
              <Tag name="Maria" className="-right-2 top-52 xl:-right-10" tone="bg-orange-600" />
            </div>

            <h1 className="text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
              The CRM built to turn pipeline into{' '}
              <span className="text-orange-700 dark:text-orange-400">revenue</span>
            </h1>

            <p className="mx-auto mt-5 max-w-xl text-lg text-gray-600 dark:text-gray-300">
              Bring every deal, contact and rep into one aligned workspace, so nothing slips and
              every quarter closes strong.
            </p>

            <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-xl bg-orange-700 px-6 text-sm font-semibold text-white hover:bg-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
              >
                Start for free
              </a>
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-xl border border-gray-300 bg-white px-6 text-sm font-semibold text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:bg-transparent dark:text-white dark:hover:bg-white/5"
              >
                Book a demo
              </a>
            </div>
          </div>

          {/* The product shot. A drawing of an app, so the whole subtree is
              hidden: nothing in it is focusable and none of it is a fact. */}
          <div
            aria-hidden="true"
            className="mx-auto mt-14 max-w-5xl overflow-hidden rounded-t-2xl border border-b-0 border-gray-200 bg-white shadow-2xl shadow-gray-900/10 dark:border-white/10 dark:bg-gray-900"
          >
            <div className="flex">
              <div className="hidden w-48 shrink-0 border-r border-gray-200 p-3 sm:block dark:border-white/10">
                <p className="mb-3 flex items-center gap-2 text-sm font-semibold">
                  <span className="grid size-6 place-items-center rounded-md bg-orange-700 text-[10px] font-bold text-white">
                    F
                  </span>
                  Flowen
                </p>
                <p className="px-2 py-1 text-[10px] uppercase tracking-wide text-gray-400">
                  Daily operation
                </p>
                {['Dashboard', 'My deals', 'Contacts', 'Opportunities', 'Activities', 'Automations'].map(
                  (item, index) => (
                    <p
                      key={item}
                      className={`rounded-md px-2 py-1.5 text-xs ${
                        index === 0
                          ? 'bg-orange-50 font-medium text-orange-800 dark:bg-orange-400/10 dark:text-orange-300'
                          : 'text-gray-600 dark:text-gray-300'
                      }`}
                    >
                      {item}
                    </p>
                  ),
                )}
              </div>

              <div className="min-w-0 flex-1 p-4">
                <div className="mb-4 flex items-center justify-between">
                  <p className="text-sm font-semibold">Dashboard</p>
                  <span className="rounded-lg bg-orange-700 px-2.5 py-1 text-[10px] font-semibold text-white">
                    New workflow
                  </span>
                </div>

                <div className="grid gap-3 sm:grid-cols-3">
                  <div className="rounded-xl border border-gray-200 p-3 sm:col-span-2 dark:border-white/10">
                    <div className="mb-2 flex items-baseline gap-2">
                      <p className="text-xs text-gray-500 dark:text-gray-400">Pipeline activity</p>
                      <p className="text-lg font-semibold">274</p>
                    </div>
                    {/* h-full on the column, not just the row: a percentage
                        height resolves against the parent's height, and a
                        content-sized flex column has none, so the bars
                        collapse to nothing. */}
                    <div className="flex h-24 items-stretch gap-3">
                      {ACTIVITY.map((height, index) => (
                        <span
                          key={index}
                          className="flex h-full flex-1 flex-col justify-end gap-1"
                        >
                          <span
                            style={{ height: `${height}%` }}
                            className="w-full rounded-t bg-orange-300 dark:bg-orange-400/40"
                          />
                          <span className="text-center text-[9px] text-gray-400">
                            {DAYS[index]}
                          </span>
                        </span>
                      ))}
                    </div>
                  </div>

                  <div className="rounded-xl border border-gray-200 p-3 dark:border-white/10">
                    <p className="mb-2 text-xs text-gray-500 dark:text-gray-400">Deal stage</p>
                    <div className="relative mx-auto grid size-24 place-items-center">
                      <svg viewBox="0 0 36 36" className="size-24 -rotate-90">
                        <circle
                          cx="18"
                          cy="18"
                          r="15.9"
                          fill="none"
                          strokeWidth="4"
                          className="stroke-gray-200 dark:stroke-white/10"
                        />
                        <circle
                          cx="18"
                          cy="18"
                          r="15.9"
                          fill="none"
                          strokeWidth="4"
                          strokeDasharray="78 100"
                          className="stroke-orange-500"
                        />
                      </svg>
                      <span className="absolute text-sm font-semibold">150</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* ── Logo cloud ─────────────────────────────────────────────────── */}
        <section aria-labelledby="customers" className="px-6 py-14">
          <div className="mx-auto max-w-5xl">
            <h2
              id="customers"
              className="text-center text-lg font-semibold text-gray-900 dark:text-white"
            >
              Trusted by 104+ businesses
            </h2>
            <ul className="mt-8 grid grid-cols-2 gap-px overflow-hidden rounded-xl border border-gray-200 bg-gray-200 sm:grid-cols-3 lg:grid-cols-4 dark:border-white/10 dark:bg-white/10">
              {LOGOS.map((logo) => (
                <li
                  key={logo}
                  className="flex items-center justify-center gap-2 bg-white px-4 py-6 dark:bg-gray-950"
                >
                  <span aria-hidden="true" className="size-4 rounded-sm bg-gray-400" />
                  {/* gray-600 rather than gray-400: a wordmark is text, and a
                      logo cloud nobody can read is a row of grey smudges. */}
                  <span className="text-sm font-semibold text-gray-600 dark:text-gray-300">
                    {logo}
                  </span>
                </li>
              ))}
            </ul>
            <p className="mt-4 text-center text-sm text-gray-500 dark:text-gray-400">
              + more companies
            </p>
          </div>
        </section>

        {/* ── Figures ────────────────────────────────────────────────────── */}
        <section aria-labelledby="results" className="px-6 pb-16">
          <h2 id="results" className="sr-only">
            Results our customers see
          </h2>
          <dl className="mx-auto grid max-w-5xl gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {FIGURES.map((figure) => (
              <div
                key={figure.label}
                className="rounded-xl border border-gray-200 p-6 dark:border-white/10"
              >
                <span
                  aria-hidden="true"
                  className="mb-4 block size-9 rounded-lg bg-orange-700 opacity-90"
                />
                <dt className="text-3xl font-semibold tracking-tight">{figure.value}</dt>
                <dd className="mt-1 text-sm text-gray-600 dark:text-gray-300">{figure.label}</dd>
              </div>
            ))}
          </dl>
        </section>

        {/* ── Feature bento ──────────────────────────────────────────────── */}
        <section
          aria-labelledby="features"
          className="border-t border-gray-200 px-6 py-16 dark:border-white/10"
        >
          <div className="mx-auto max-w-5xl">
            <div className="mb-10 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
              <div>
                <Eyebrow>Features</Eyebrow>
                <h2
                  id="features"
                  className="mt-4 max-w-md text-3xl font-semibold tracking-tight sm:text-4xl"
                >
                  Smart features for modern sales teams
                </h2>
              </div>
              <p className="max-w-sm text-sm text-gray-600 dark:text-gray-300">
                Automate tasks, track performance and manage customers effortlessly in one powerful
                CRM platform.
              </p>
            </div>

            <ul className="grid gap-4 md:grid-cols-2">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className="rounded-2xl border border-gray-200 p-4 dark:border-white/10"
                >
                  {/* The panel is a drawing of the product, so it is hidden.
                      The heading and copy underneath carry the meaning. */}
                  <div aria-hidden="true">
                    <Card className="mb-4">
                      <FeatureMock kind={feature.mock} />
                    </Card>
                  </div>
                  <h3 className="text-base font-semibold">{feature.title}</h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{feature.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* ── Capabilities accordion ─────────────────────────────────────── */}
        <section
          aria-labelledby="growth"
          className="border-t border-gray-200 px-6 py-16 dark:border-white/10"
        >
          <div className="mx-auto max-w-5xl">
            <div className="mb-10 flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
              <div>
                <Eyebrow>Why Flowen</Eyebrow>
                <h2
                  id="growth"
                  className="mt-4 max-w-md text-3xl font-semibold tracking-tight sm:text-4xl"
                >
                  Power your growth with smart, effortless CRM
                </h2>
              </div>
              <p className="max-w-sm text-sm text-gray-600 dark:text-gray-300">
                All the tools you need to scale your revenue team, beautifully designed and easy to
                use.
              </p>
            </div>

            <div className="grid gap-6 lg:grid-cols-2">
              {/* Native disclosure widgets: keyboard behaviour and the
                  expanded state are the browser's job, and they work with
                  JavaScript switched off. */}
              <ul className="space-y-3">
                {CAPABILITIES.map((item) => (
                  <li key={item.title}>
                    <details
                      open={item.open}
                      className="group rounded-xl border border-gray-200 open:border-orange-300 open:bg-orange-50/60 dark:border-white/10 dark:open:border-orange-400/30 dark:open:bg-orange-400/5"
                    >
                      <summary className="flex min-h-12 cursor-pointer list-none items-center justify-between gap-4 px-4 py-3 text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700">
                        <span className="flex items-center gap-3">
                          <span
                            aria-hidden="true"
                            className="grid size-7 shrink-0 place-items-center rounded-md bg-orange-700 text-xs font-bold text-white"
                          >
                            {item.title.charAt(0)}
                          </span>
                          {item.title}
                        </span>
                        <span
                          aria-hidden="true"
                          className="text-lg leading-none text-gray-400 group-open:hidden"
                        >
                          +
                        </span>
                        <span
                          aria-hidden="true"
                          className="hidden text-lg leading-none text-gray-400 group-open:inline"
                        >
                          &minus;
                        </span>
                      </summary>
                      <p className="px-4 pb-4 pl-14 text-sm text-gray-600 dark:text-gray-300">
                        {item.body}
                      </p>
                    </details>
                  </li>
                ))}
              </ul>

              <div
                aria-hidden="true"
                className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-orange-100 to-amber-50 p-6 dark:from-orange-500/15 dark:to-amber-500/5"
              >
                <Card>
                  <p className="text-xs text-gray-500 dark:text-gray-400">Success rate</p>
                  <div className="mt-4 flex items-end justify-center gap-1">
                    {[30, 44, 58, 72, 86, 100, 86, 72, 58, 44, 30].map((height, index) => (
                      <span
                        key={index}
                        style={{ height: `${height * 0.6}px` }}
                        className={`w-2 rounded-full ${index === 5 ? 'bg-orange-600' : 'bg-orange-300 dark:bg-orange-400/40'}`}
                      />
                    ))}
                  </div>
                  <p className="mt-3 text-center text-2xl font-semibold">72.5%</p>
                </Card>
                <Card className="mt-4 w-2/3">
                  <p className="text-xs text-gray-500 dark:text-gray-400">Total revenue</p>
                  <p className="text-lg font-semibold">$68,837</p>
                  <p className="text-[10px] font-medium text-emerald-700 dark:text-emerald-400">
                    +4.1% vs last month
                  </p>
                </Card>
              </div>
            </div>
          </div>
        </section>

        {/* ── Pricing ────────────────────────────────────────────────────── */}
        <section
          aria-labelledby="pricing"
          className="border-t border-gray-200 px-6 py-16 dark:border-white/10"
        >
          <div className="mx-auto max-w-5xl">
            <div className="mb-10 text-center">
              <Eyebrow>Pricing plans</Eyebrow>
              <h2 id="pricing" className="mt-4 text-3xl font-semibold tracking-tight sm:text-4xl">
                Select the plan that fits your needs
              </h2>
              <p className="mx-auto mt-3 max-w-lg text-sm text-gray-600 dark:text-gray-300">
                Choose the plan that fits your workflow and start growing your business.
              </p>
            </div>

            <ul className="grid items-start gap-6 lg:grid-cols-3">
              {PLANS.map((plan) => (
                <li
                  key={plan.name}
                  className={`rounded-2xl border p-6 ${
                    plan.featured
                      ? 'border-orange-700 bg-orange-700 text-white lg:-mt-4 lg:pb-10'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="text-base font-semibold">{plan.name}</h3>
                    {/* Said in words. Elevation and colour say "popular" only
                        to someone looking at the page. */}
                    {plan.featured && (
                      <span className="rounded-full bg-white px-2.5 py-0.5 text-xs font-semibold text-orange-800">
                        Most popular
                      </span>
                    )}
                  </div>

                  <p
                    className={`mt-2 text-sm ${plan.featured ? 'text-orange-50' : 'text-gray-600 dark:text-gray-300'}`}
                  >
                    {plan.blurb}
                  </p>

                  <p className="mt-6 flex items-baseline gap-1">
                    <span className="text-4xl font-semibold tracking-tight">{plan.price}</span>
                    <span
                      className={`text-sm ${plan.featured ? 'text-orange-100' : 'text-gray-500 dark:text-gray-400'}`}
                    >
                      /month
                    </span>
                  </p>

                  <a
                    href="#"
                    className={`mt-6 flex min-h-12 items-center justify-center rounded-xl px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 ${
                      plan.featured
                        ? 'bg-white text-orange-800 hover:bg-orange-50 focus-visible:outline-white'
                        : 'bg-gray-900 text-white hover:bg-gray-800 focus-visible:outline-orange-700 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100'
                    }`}
                  >
                    {/* Names the plan, so three "Get started" links are not
                        three identical entries in a list of links. */}
                    {plan.cta}
                    <span className="sr-only"> with the {plan.name}</span>
                  </a>

                  <p
                    className={`mt-6 text-xs font-semibold uppercase tracking-wide ${plan.featured ? 'text-orange-100' : 'text-gray-500 dark:text-gray-400'}`}
                  >
                    {plan.includes}
                  </p>
                  <ul className="mt-3 space-y-2 text-sm">
                    {plan.features.map((feature) => (
                      <li key={feature} className="flex items-start gap-2">
                        <span
                          aria-hidden="true"
                          className={plan.featured ? 'text-orange-100' : 'text-orange-700 dark:text-orange-400'}
                        >
                          &#10003;
                        </span>
                        <span className={plan.featured ? '' : 'text-gray-700 dark:text-gray-200'}>
                          {feature}
                        </span>
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
            </ul>
          </div>
        </section>
      </main>
    </div>
  )
}
