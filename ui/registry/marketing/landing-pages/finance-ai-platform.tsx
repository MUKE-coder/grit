/*
 * A finance platform page built on photography: landscape bands behind most
 * sections, product surfaces floating over them on cards.
 *
 * One file, because the registry installs one file per block.
 *
 * Text over a photograph is the whole risk of this design. Contrast against a
 * photograph is not one number — it is a different number per pixel, and a
 * headline that reads over a bright sky disappears over the dark ridge two
 * hundred pixels along. So no text here sits directly on an image. Either it
 * sits on a card, or it sits on a scrim of known opacity over a known base,
 * which makes the effective background a colour that can be measured instead
 * of hoped about.
 *
 * Every photograph is decoration: alt="" and aria-hidden. A hillside behind a
 * pricing table is mood, and describing it to somebody who cannot see it wastes
 * their time in the middle of choosing a plan. The persona photographs are the
 * same — the heading already says who the section is for, and the picture only
 * repeats it less precisely.
 *
 * The frosted chips over the photographs are tinted with gray-950/40 rather
 * than the white/15 this look is usually built from. A translucent white under
 * white text raises the floor and sinks the contrast: measured against the
 * brightest pixel of the hero photograph, the white version of the pill and the
 * ghost button came out at 4.1:1. Darkening the chip instead of lightening it
 * keeps the same glassy effect and takes them past 10:1.
 *
 * The figures are description lists. "99.9%" beside "uptime" is a term and its
 * definition, not two paragraphs that happen to be stacked.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

/* Verified on a contact sheet before use. One candidate in the batch was
   labelled lavender and turned out to be a scoop of soil, which is the reason
   the sheet exists. */
const LAND = {
  hills: '1501854140801-50d01698950b',
  valley: '1472214103451-9374bd1c798e',
  poppies: '1465146344425-f00d5f5c8f07',
  ridge: '1500534314209-a25ddb2bd429',
  clouds: '1506905925346-21bda4d32df4',
}

const PEOPLE = {
  team: '1552664730-d307ca884978',
  advisor: '1573497019940-1c28c88b4f3e',
  desk: '1522071820081-009f0129c71c',
}

const FACES = ['1494790108377-be9c29b29330', '1500648767791-00dcc994a43e', '1507003211169-0a1dd7228f2d']

const land = (id: string, w = 1600, h = 900) =>
  `https://images.unsplash.com/photo-${id}?w=${w}&h=${h}&fit=crop&q=80`
const face = (id: string, size = 80) =>
  `https://images.unsplash.com/photo-${id}?w=${size * 2}&h=${size * 2}&fit=crop&crop=faces&q=75`

const CUSTOMERS = ['Lattice', 'Ampry', 'Ardent', 'Phoenix', 'Vercel', 'Linden']

const CHALLENGES = [
  'Financial data is spread across four dashboards and a spreadsheet',
  'Reporting arrives late enough that the decision is already made',
  'Nobody can say which assumption a number came from',
]

const CHALLENGE_FIGURES = [
  { value: '58%', label: 'of teams still reconcile by hand each month' },
  { value: '51%', label: 'cannot trace a figure back to its source' },
]

const FEATURES = [
  {
    title: 'Advanced risk analysis',
    body: 'Exposure by position, sector and currency, recalculated as the market moves rather than overnight.',
    tone: 'plain',
    span: 'lg:col-span-2',
  },
  {
    title: 'Market insights',
    body: 'The moves that touch your holdings, with the reason attached and the source named.',
    tone: 'chart',
    span: 'lg:col-span-2',
  },
  {
    title: 'AI-powered insights',
    body: 'A recommendation is always shown with the reasoning that produced it. A rating with no argument behind it is a horoscope.',
    tone: 'dark',
    span: 'lg:col-span-2',
  },
  {
    title: 'Portfolio tracking',
    body: 'Every account in one balance, updated through the day.',
    tone: 'donut',
    span: 'lg:col-span-3',
  },
  {
    title: 'Smart alerts',
    body: 'Told once, at the threshold you set, in the channel you already read.',
    tone: 'plain',
    span: 'lg:col-span-3',
  },
]

const START_FIGURES = [
  { value: '100%', label: 'Bank-level encryption in transit and at rest' },
  { value: '2 minutes', label: 'Median time from sign-up to a connected account' },
]

const COMPLIANCE = ['SOC 2 Type II', 'ISO 27001', 'GDPR', 'PCI DSS']

const AUDIENCES = [
  {
    name: 'Financial teams',
    body: 'Close the month without three people rebuilding the same figure from three exports.',
    photo: PEOPLE.team,
    points: ['Shared workspace', 'Audit trail on every figure'],
  },
  {
    name: 'Wealth managers',
    body: 'Client reporting that comes out of the system rather than out of a Sunday evening.',
    photo: PEOPLE.advisor,
    points: ['Per-client views', 'Scheduled statements'],
  },
  {
    name: 'Individual investors',
    body: 'One balance across every account you hold, without a spreadsheet in the middle.',
    photo: PEOPLE.desk,
    points: ['All accounts in one place', 'Alerts you choose'],
  },
]

const INTEGRATIONS = ['Plaid', 'Xero', 'QuickBooks', 'Stripe', 'Slack', 'Notion', 'Sheets', 'Snowflake']

const FIGURES = [
  { value: '10,000+', label: 'Active investors' },
  { value: '$1M+', label: 'Assets tracked daily' },
  { value: '99.9%', label: 'Uptime over twelve months' },
  { value: '120+', label: 'Institutions connected' },
]

const TESTIMONIALS = [
  {
    quote:
      'The reporting that used to take our team a full week now takes an afternoon, and I can see where every number came from without asking anybody.',
    name: 'Claire Devereaux',
    role: 'Finance director, Ardent',
    face: 0,
  },
  {
    quote:
      'I moved eleven client portfolios across in a morning. The part that sold me was that it refused to guess at two transactions and asked me instead.',
    name: 'Sam Thompson',
    role: 'Wealth manager, Linden',
    face: 1,
  },
  {
    quote:
      'It is the first tool of this kind I have used that shows its working. When it flags a position it tells me which signal fired and when.',
    name: 'Michael Chen',
    role: 'Private investor',
    face: 2,
  },
]

const TIERS = [
  {
    name: 'Starter',
    price: '$19',
    blurb: 'For one person keeping track of their own accounts.',
    features: ['Up to 5 connected accounts', 'Daily balance updates', 'Basic alerts', 'Email support'],
    featured: false,
    inherits: null as string | null,
  },
  {
    name: 'Pro',
    price: '$39',
    blurb: 'For an adviser or a small team reporting to somebody else.',
    features: ['Unlimited accounts', 'Client reporting', 'Custom alerts and rules', 'Same-day support'],
    featured: true,
    inherits: 'Starter',
  },
]

const FAQS = [
  {
    q: 'How do the AI insights work?',
    a: 'Signals are computed from your positions and public market data, and every insight is shown with the signals that produced it. Nothing is generated from a prompt and presented as analysis.',
  },
  {
    q: 'Is my financial data safe?',
    a: 'Connections are read-only by default and credentials never reach our servers — the aggregator holds them. Data is encrypted in transit and at rest, and you can delete an account and its history at any time.',
  },
  {
    q: 'Do you offer plans for organisations?',
    a: 'Yes. Above roughly twenty seats it stops being a per-seat conversation and becomes one about single sign-on, data residency and who signs the DPA.',
  },
  {
    q: 'Can I export everything?',
    a: 'CSV and a documented API, including history. Leaving should be as easy as arriving, and a product that makes it hard is telling you something.',
  },
]

/* ── Pieces ─────────────────────────────────────────────────────────────── */

/* A photograph plus a scrim of known opacity. The scrim is what makes the text
   on top measurable: without it the effective background is whatever the
   photograph happens to be doing at that pixel. */
function PhotoBand({
  id,
  scrim = 'bg-gray-950/55',
  className = '',
  children,
}: {
  id: string
  scrim?: string
  className?: string
  children: ReactNode
}) {
  return (
    <div className={`relative isolate overflow-hidden ${className}`}>
      <img src={land(id)} alt="" aria-hidden="true" className="absolute inset-0 -z-10 size-full object-cover" />
      <span aria-hidden="true" className={`absolute inset-0 -z-10 ${scrim}`} />
      {children}
    </div>
  )
}

function Card({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div className={`rounded-2xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900 ${className}`}>
      {children}
    </div>
  )
}

function DashboardMock({ compact = false }: { compact?: boolean }) {
  return (
    <div
      aria-hidden="true"
      className="overflow-hidden rounded-2xl border border-gray-200 bg-white shadow-2xl select-none dark:border-white/10 dark:bg-gray-900"
    >
      <div className="flex items-center gap-2 border-b border-gray-200 px-4 py-2.5 dark:border-white/10">
        <span className="size-2.5 rounded-full bg-gray-300 dark:bg-white/20" />
        <span className="size-2.5 rounded-full bg-gray-300 dark:bg-white/20" />
        <span className="size-2.5 rounded-full bg-gray-300 dark:bg-white/20" />
        <span className="ml-3 h-5 w-40 rounded bg-gray-100 dark:bg-white/10" />
      </div>
      <div className="grid gap-4 p-4 sm:grid-cols-[9rem_1fr]">
        <div className="hidden space-y-2 sm:block">
          {['Overview', 'Portfolio', 'Insights', 'Reports', 'Settings'].map((row, i) => (
            <div
              key={row}
              className={`rounded-lg px-3 py-2 text-[11px] ${
                i === 0
                  ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'
                  : 'text-gray-500 dark:text-gray-400'
              }`}
            >
              {row}
            </div>
          ))}
        </div>
        <div className="space-y-4">
          <div className="grid gap-3 sm:grid-cols-3">
            {[
              { k: 'Balance', v: '$482,910' },
              { k: 'Today', v: '+1.4%' },
              { k: 'Cash', v: '$21,004' },
            ].map((tile) => (
              <div key={tile.k} className="rounded-xl border border-gray-200 p-3 dark:border-white/10">
                <p className="text-[10px] text-gray-500 dark:text-gray-400">{tile.k}</p>
                <p className="mt-1 text-sm font-semibold tabular-nums">{tile.v}</p>
              </div>
            ))}
          </div>
          {!compact && (
            <div className="rounded-xl border border-gray-200 p-3 dark:border-white/10">
              <div className="flex h-28 items-end gap-1.5">
                {[38, 52, 44, 61, 49, 72, 58, 80, 66, 88, 74, 95].map((h, i) => (
                  <span
                    key={i}
                    style={{ height: `${h}%` }}
                    className="flex-1 rounded-t bg-gradient-to-t from-emerald-600/30 to-emerald-600"
                  />
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function Donut() {
  /* conic-gradient rather than an SVG: one element, no viewBox arithmetic, and
     it is decoration so it never needs to be read. */
  return (
    <span
      aria-hidden="true"
      className="block size-24 rounded-full"
      style={{
        background:
          'conic-gradient(#059669 0 45%, #0d9488 45% 70%, #6366f1 70% 88%, #e5e7eb 88% 100%)',
        mask: 'radial-gradient(circle, transparent 55%, black 56%)',
        WebkitMask: 'radial-gradient(circle, transparent 55%, black 56%)',
      }}
    />
  )
}

function Sparkline() {
  return (
    <span aria-hidden="true" className="flex h-16 items-end gap-1">
      {[30, 48, 36, 58, 44, 70, 52, 78, 64, 86].map((h, i) => (
        <span key={i} style={{ height: `${h}%` }} className="flex-1 rounded-t bg-emerald-600/80" />
      ))}
    </span>
  )
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function FinanceAiPlatform() {
  return (
    <div className="min-h-screen bg-gray-50 font-sans text-gray-950 antialiased dark:bg-gray-950 dark:text-gray-50">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:rounded-full focus:bg-gray-950 focus:px-4 focus:py-2 focus:text-sm focus:text-white"
      >
        Skip to content
      </a>

      <main id="main">
        {/* Hero over photography */}
        <PhotoBand id={LAND.hills} scrim="bg-gray-950/55" className="px-6 pt-6 pb-40">
          <div className="mx-auto max-w-6xl">
            <header className="flex h-14 items-center justify-between rounded-full bg-white/95 px-5 shadow-sm backdrop-blur dark:bg-gray-900/95">
              <a href="#" className="flex items-center gap-2 text-sm font-semibold">
                <span aria-hidden="true" className="size-5 rounded-md bg-emerald-700" />
                Financeai
              </a>
              <nav aria-label="Primary" className="hidden md:block">
                <ul role="list" className="flex items-center gap-7 text-sm">
                  {['Product', 'Features', 'Pricing', 'Blog'].map((item) => (
                    <li key={item}>
                      <a href="#" className="text-gray-600 hover:text-gray-950 dark:text-gray-300 dark:hover:text-white">
                        {item}
                      </a>
                    </li>
                  ))}
                </ul>
              </nav>
              <div className="flex items-center gap-3 text-sm">
                <a href="#" className="hidden text-gray-600 hover:text-gray-950 sm:block dark:text-gray-300 dark:hover:text-white">
                  Sign in
                </a>
                <a href="#pricing" className="rounded-full bg-gray-950 px-4 py-2 font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
                  Get started
                </a>
              </div>
            </header>

            <div className="mx-auto mt-16 max-w-3xl text-center">
              {/* White on a 55% scrim over the photograph. The scrim is the
                  point: it turns "white on a hillside" into a value that can be
                  measured, and it is why the pill below is not white-on-white
                  when the sun happens to be in frame. */}
              <p className="inline-flex items-center gap-2 rounded-full bg-gray-950/40 px-3 py-1 text-xs font-medium text-white ring-1 ring-white/30">
                <span aria-hidden="true" className="size-1.5 rounded-full bg-emerald-300" />
                Now with automated reconciliation
              </p>
              <h1 className="mt-5 text-4xl font-semibold tracking-tight text-balance text-white sm:text-6xl">
                Finance AI platform
              </h1>
              <p className="mx-auto mt-5 max-w-xl text-pretty text-gray-100">
                Optimise your investments with AI-driven analysis, real-time tracking and intelligent
                recommendations that show their working.
              </p>
              <div className="mt-8 flex flex-wrap justify-center gap-3">
                <a href="#pricing" className="rounded-full bg-white px-6 py-3 text-sm font-semibold text-gray-950 hover:bg-gray-100">
                  Start for free
                </a>
                <a href="#features" className="rounded-full bg-gray-950/40 px-6 py-3 text-sm font-semibold text-white ring-1 ring-white/40 hover:bg-gray-950/55">
                  See how it works
                </a>
              </div>
            </div>
          </div>
        </PhotoBand>

        {/* Dashboard lifted over the seam between the photograph and the page */}
        <div className="mx-auto -mt-32 max-w-5xl px-6">
          <DashboardMock />
        </div>

        {/* Customers */}
        <section aria-labelledby="customers-heading" className="px-6 py-14">
          <div className="mx-auto max-w-6xl text-center">
            <h2 id="customers-heading" className="text-sm text-gray-500 dark:text-gray-400">
              Trusted by teams managing real money
            </h2>
            <ul role="list" className="mt-6 flex flex-wrap items-center justify-center gap-x-10 gap-y-4">
              {CUSTOMERS.map((name) => (
                <li key={name} className="text-sm font-semibold text-gray-500 dark:text-gray-400">
                  {name}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Problem */}
        <section aria-labelledby="clarity-heading" className="px-6 py-16">
          <div className="mx-auto max-w-4xl text-center">
            <h2 id="clarity-heading" className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Smarter decisions start with clear data
            </h2>
            <p className="mx-auto mt-4 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
              Most teams do not lack data. They lack one version of it that everybody agrees on.
            </p>
          </div>

          <Card className="mx-auto mt-12 max-w-4xl overflow-hidden">
            <div className="grid gap-8 p-8 md:grid-cols-2">
              <div>
                <h3 className="text-sm font-semibold">Challenges of managing investments today</h3>
                <ul role="list" className="mt-4 space-y-3 text-sm">
                  {CHALLENGES.map((item) => (
                    <li key={item} className="flex gap-2.5 text-gray-600 dark:text-gray-300">
                      <span aria-hidden="true" className="mt-1.5 size-1.5 shrink-0 rounded-full bg-rose-500" />
                      {item}
                    </li>
                  ))}
                </ul>
              </div>
              <dl className="space-y-6">
                {CHALLENGE_FIGURES.map((figure) => (
                  <div key={figure.value} className="rounded-xl bg-gray-50 p-5 dark:bg-white/5">
                    <dt className="text-3xl font-semibold tracking-tight tabular-nums text-rose-700 dark:text-rose-400">
                      {figure.value}
                    </dt>
                    <dd className="mt-1 text-sm text-pretty text-gray-600 dark:text-gray-400">{figure.label}</dd>
                  </div>
                ))}
              </dl>
            </div>
          </Card>
        </section>

        {/* Features */}
        <section id="features" aria-labelledby="features-heading" className="px-6 py-16">
          <div className="mx-auto max-w-6xl">
            <div className="flex flex-wrap items-end justify-between gap-6">
              <div>
                <p className="text-sm font-medium text-emerald-700 dark:text-emerald-400">Capabilities</p>
                <h2 id="features-heading" className="mt-2 max-w-lg text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                  Everything you need to invest confidently
                </h2>
              </div>
              <a href="#" className="rounded-full bg-gray-950 px-5 py-2.5 text-sm font-semibold text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
                See all features
              </a>
            </div>

            {/* Six columns so 2+2+2 and 3+3 both tile exactly. A span that does
                not divide the track count leaves a hole, and grid quietly fills
                it with the next tile instead of telling you. */}
            <ul role="list" className="mt-12 grid gap-5 lg:grid-cols-6">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`rounded-2xl border p-6 ${feature.span} ${
                    feature.tone === 'dark'
                      ? 'border-gray-900 bg-gray-950 text-gray-50 dark:border-white/10'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  {feature.tone === 'chart' && <Sparkline />}
                  {feature.tone === 'donut' && <Donut />}
                  {feature.tone === 'dark' && (
                    <p aria-hidden="true" className="font-mono text-3xl font-bold tracking-[0.3em] text-emerald-400">
                      HOLD
                    </p>
                  )}
                  <h3 className={`font-semibold ${feature.tone === 'plain' ? '' : 'mt-5'}`}>{feature.title}</h3>
                  <p
                    className={`mt-2 text-sm text-pretty ${
                      feature.tone === 'dark' ? 'text-gray-300' : 'text-gray-600 dark:text-gray-400'
                    }`}
                  >
                    {feature.body}
                  </p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Product over photography */}
        <PhotoBand id={LAND.valley} scrim="bg-gray-950/60" className="px-6 py-20">
          <div className="mx-auto max-w-5xl text-center">
            <h2 className="text-3xl font-semibold tracking-tight text-balance text-white sm:text-4xl">
              See your financial intelligence in action
            </h2>
            <p className="mx-auto mt-4 max-w-xl text-pretty text-gray-100">
              One workspace for balances, exposure and the reasoning behind every flag.
            </p>
            <div className="mt-10 text-left">
              <DashboardMock />
            </div>
          </div>
        </PhotoBand>

        {/* Onboarding */}
        <section aria-labelledby="start-heading" className="px-6 py-20">
          <div className="mx-auto grid max-w-6xl gap-10 lg:grid-cols-2">
            <div>
              <h2 id="start-heading" className="max-w-md text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Start investing in minutes
              </h2>
              <p className="mt-4 max-w-md text-pretty text-gray-600 dark:text-gray-300">
                Connect an account, confirm what we read, and the rest fills itself in. There is no
                implementation call, because there is nothing to implement.
              </p>
              <dl className="mt-8 grid gap-6 sm:grid-cols-2">
                {START_FIGURES.map((figure) => (
                  <div key={figure.value}>
                    <dt className="text-3xl font-semibold tracking-tight tabular-nums">{figure.value}</dt>
                    <dd className="mt-1 text-sm text-pretty text-gray-600 dark:text-gray-400">{figure.label}</dd>
                  </div>
                ))}
              </dl>
            </div>

            <Card className="p-8">
              <h3 className="font-semibold">Connect your accounts</h3>
              <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">
                Read-only by default. Credentials go to the aggregator, never to us, and you can
                revoke a connection from either side.
              </p>
              <ul role="list" className="mt-6 space-y-2.5 text-sm">
                {['Securely link your bank, brokerage and investment accounts', 'Confirm what we read before anything is stored', 'Revoke access at any time, from either end'].map((item) => (
                  <li key={item} className="flex gap-2.5">
                    <span aria-hidden="true" className="mt-0.5 text-emerald-700 dark:text-emerald-400">✓</span>
                    <span className="text-gray-600 dark:text-gray-300">{item}</span>
                  </li>
                ))}
              </ul>
            </Card>
          </div>
        </section>

        {/* Security */}
        <section aria-labelledby="security-heading" className="px-6 pb-20">
          <Card className="mx-auto max-w-6xl p-8 sm:p-12">
            <div className="text-center">
              <h2 id="security-heading" className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Your data is protected at every level
              </h2>
              <p className="mx-auto mt-4 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Audited annually, and the reports are available before you ask rather than after you
                sign.
              </p>
            </div>
            <ul role="list" className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {COMPLIANCE.map((badge) => (
                <li
                  key={badge}
                  className="rounded-xl border border-gray-200 px-4 py-6 text-center text-sm font-semibold dark:border-white/10"
                >
                  {badge}
                </li>
              ))}
            </ul>
          </Card>
        </section>

        {/* Audiences */}
        <section aria-labelledby="audiences-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="audiences-heading" className="text-center text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Who this platform is built for
            </h2>
            <ul role="list" className="mt-12 grid gap-6 lg:grid-cols-3">
              {AUDIENCES.map((audience) => (
                <li key={audience.name}>
                  <Card className="h-full overflow-hidden">
                    <img
                      src={land(audience.photo, 800, 500)}
                      alt=""
                      aria-hidden="true"
                      loading="lazy"
                      className="h-44 w-full object-cover"
                    />
                    <div className="p-6">
                      <h3 className="font-semibold">{audience.name}</h3>
                      <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{audience.body}</p>
                      <ul role="list" className="mt-4 space-y-2 text-sm">
                        {audience.points.map((point) => (
                          <li key={point} className="flex gap-2.5 text-gray-600 dark:text-gray-300">
                            <span aria-hidden="true" className="mt-0.5 text-emerald-700 dark:text-emerald-400">✓</span>
                            {point}
                          </li>
                        ))}
                      </ul>
                    </div>
                  </Card>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Integrations over photography */}
        <PhotoBand id={LAND.clouds} scrim="bg-gray-950/65" className="px-6 py-20">
          <div className="mx-auto max-w-4xl text-center">
            <h2 className="text-3xl font-semibold tracking-tight text-balance text-white sm:text-4xl">
              Connect with the tools you already use
            </h2>
            <p className="mx-auto mt-4 max-w-xl text-pretty text-gray-100">
              Over a hundred integrations, and a documented API for the ones we have not built yet.
            </p>
            <ul role="list" className="mx-auto mt-10 flex max-w-2xl flex-wrap justify-center gap-3">
              {INTEGRATIONS.map((tool) => (
                <li
                  key={tool}
                  className="rounded-full bg-gray-950/40 px-4 py-2 text-sm font-medium text-white ring-1 ring-white/30"
                >
                  {tool}
                </li>
              ))}
            </ul>
          </div>
        </PhotoBand>

        {/* Figures */}
        <section aria-labelledby="figures-heading" className="px-6 py-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="figures-heading" className="sr-only">
              Platform in numbers
            </h2>
            <dl className="grid gap-px overflow-hidden rounded-2xl border border-gray-200 bg-gray-200 sm:grid-cols-2 lg:grid-cols-4 dark:border-white/10 dark:bg-white/10">
              {FIGURES.map((figure) => (
                <div key={figure.value} className="bg-gray-50 p-6 dark:bg-gray-950">
                  <dt className="text-3xl font-semibold tracking-tight tabular-nums">{figure.value}</dt>
                  <dd className="mt-2 text-sm text-gray-600 dark:text-gray-400">{figure.label}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>

        {/* Testimonials over photography */}
        {/* A lighter scrim than the other bands. This photograph is much
            brighter than the rest (mean luminance 0.71 against 0.12 for the
            hero), so the 55% used elsewhere flattened it to grey and the point
            of putting a photograph there was lost. 45% still leaves the heading
            at 3.4:1 against the brightest pixel in the frame, and the quotes do
            not depend on it because they sit on cards. */}
        <PhotoBand id={LAND.poppies} scrim="bg-gray-950/45" className="px-6 py-20">
          <div className="mx-auto max-w-6xl">
            <h2 className="text-center text-3xl font-semibold tracking-tight text-balance text-white sm:text-4xl">
              What investors say about the platform
            </h2>
            {/* The quotes sit on cards rather than on the flowers. White text on
                a poppy field is legible in the sky and gone in the petals. */}
            <ul role="list" className="mt-12 grid gap-6 lg:grid-cols-3">
              {TESTIMONIALS.map((item) => (
                <li key={item.name}>
                  <Card className="h-full p-6">
                    <figure>
                      <p aria-hidden="true" className="text-sm tracking-widest text-amber-500">
                        ★★★★★
                      </p>
                      <blockquote className="mt-4 text-sm text-pretty text-gray-700 dark:text-gray-200">
                        {item.quote}
                      </blockquote>
                      <figcaption className="mt-5 flex items-center gap-3">
                        <img
                          src={face(FACES[item.face], 36)}
                          alt=""
                          aria-hidden="true"
                          loading="lazy"
                          className="size-9 rounded-full object-cover"
                        />
                        <span className="text-sm">
                          <span className="block font-semibold">{item.name}</span>
                          <span className="block text-gray-500 dark:text-gray-400">{item.role}</span>
                        </span>
                      </figcaption>
                    </figure>
                  </Card>
                </li>
              ))}
            </ul>
          </div>
        </PhotoBand>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="px-6 py-20">
          <div className="mx-auto max-w-5xl">
            <div className="text-center">
              <h2 id="pricing-heading" className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Transparent pricing without hidden fees
              </h2>
              <p className="mx-auto mt-4 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Billed monthly, cancel whenever. The price you see is the price on the invoice.
              </p>
            </div>

            <ul role="list" className="mt-12 grid items-start gap-6 md:grid-cols-2">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured
                      ? 'border-gray-950 bg-gray-950 text-gray-50 dark:border-white/20'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="font-semibold">{tier.name}</h3>
                    {tier.featured && (
                      <span className="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-gray-950">
                        Most popular
                      </span>
                    )}
                  </div>
                  <p className={`mt-2 text-sm text-pretty ${tier.featured ? 'text-gray-300' : 'text-gray-600 dark:text-gray-400'}`}>
                    {tier.blurb}
                  </p>
                  <p className="mt-6">
                    <span className="text-4xl font-semibold tracking-tight tabular-nums">{tier.price}</span>
                    <span className={`text-sm ${tier.featured ? 'text-gray-400' : 'text-gray-500 dark:text-gray-400'}`}>
                      /month
                    </span>
                  </p>

                  {/* The link says which plan it buys. Two links both called
                      "Get started" are two identical rows in a link list, and
                      the card around them is not available there. */}
                  <a
                    href="#"
                    className={`mt-6 block rounded-full px-5 py-3 text-center text-sm font-semibold ${
                      tier.featured
                        ? 'bg-white text-gray-950 hover:bg-gray-100'
                        : 'bg-gray-950 text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950'
                    }`}
                  >
                    Choose {tier.name}
                  </a>

                  <ul role="list" className="mt-6 space-y-2.5 text-sm">
                    {tier.inherits && <li className="font-medium">Everything in {tier.inherits}, plus:</li>}
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5">
                        <span aria-hidden="true" className={`mt-0.5 ${tier.featured ? 'text-emerald-400' : 'text-emerald-700 dark:text-emerald-400'}`}>
                          ✓
                        </span>
                        <span className={tier.featured ? 'text-gray-300' : 'text-gray-600 dark:text-gray-300'}>
                          {feature}
                        </span>
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
            </ul>

            <Card className="mt-6 flex flex-wrap items-center justify-between gap-4 p-6">
              <div>
                <h3 className="font-semibold">Enterprise</h3>
                <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                  Single sign-on, data residency, and somebody who answers the DPA.
                </p>
              </div>
              <a href="#" className="rounded-full bg-gray-950 px-5 py-2.5 text-sm font-semibold text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
                Contact sales
              </a>
            </Card>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq-heading" className="px-6 pb-20">
          <div className="mx-auto grid max-w-5xl gap-10 lg:grid-cols-[20rem_1fr]">
            <div>
              <h2 id="faq-heading" className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Frequently asked questions
              </h2>
              <p className="mt-4 text-pretty text-gray-600 dark:text-gray-300">
                Still have a question? The answer is usually a reply away.
              </p>
              <a href="#" className="mt-6 inline-block rounded-full bg-gray-950 px-5 py-2.5 text-sm font-semibold text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
                Talk to us
              </a>
            </div>

            <div className="space-y-3">
              {FAQS.map((item) => (
                <details key={item.q} className="group rounded-2xl border border-gray-200 bg-white px-5 dark:border-white/10 dark:bg-gray-900">
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 py-4 text-sm font-medium marker:content-none">
                    {item.q}
                    {/* gray-600, not gray-400: list-none removes the native
                        marker, so this glyph is the only sign the row opens and
                        it owes the 3:1 of WCAG 1.4.11. */}
                    <span
                      aria-hidden="true"
                      className="mt-0.5 shrink-0 text-lg leading-none text-gray-600 transition-transform group-open:rotate-45 motion-reduce:transition-none dark:text-gray-400"
                    >
                      +
                    </span>
                  </summary>
                  <p className="pb-4 text-sm text-pretty text-gray-600 dark:text-gray-400">{item.a}</p>
                </details>
              ))}
            </div>
          </div>
        </section>

        {/* Closing */}
        <PhotoBand id={LAND.ridge} scrim="bg-gray-950/60" className="px-6 py-24">
          <div className="mx-auto max-w-3xl text-center">
            <h2 className="text-3xl font-semibold tracking-tight text-balance text-white sm:text-4xl">
              Ready to invest with more clarity?
            </h2>
            <p className="mx-auto mt-4 max-w-lg text-pretty text-gray-100">
              Free for one account, for as long as you like. No card, and no call before you can see
              the product.
            </p>
            <div className="mt-8 flex flex-wrap justify-center gap-3">
              <a href="#" className="rounded-full bg-white px-6 py-3 text-sm font-semibold text-gray-950 hover:bg-gray-100">
                Start for free
              </a>
              <a href="#" className="rounded-full bg-gray-950/40 px-6 py-3 text-sm font-semibold text-white ring-1 ring-white/40 hover:bg-gray-950/55">
                Book a walkthrough
              </a>
            </div>
          </div>
        </PhotoBand>
      </main>

      <footer className="border-t border-gray-200 px-6 py-14 dark:border-white/10">
        <div className="mx-auto grid max-w-6xl gap-10 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <p className="flex items-center gap-2 font-semibold">
              <span aria-hidden="true" className="size-5 rounded-md bg-emerald-700" />
              Financeai
            </p>
            <p className="mt-3 max-w-xs text-sm text-pretty text-gray-600 dark:text-gray-400">
              One version of your numbers, with the working shown.
            </p>
          </div>
          {[
            { heading: 'Product', links: ['Features', 'Pricing', 'Integrations', 'Changelog'] },
            { heading: 'Company', links: ['About', 'Careers', 'Blog', 'Contact'] },
            { heading: 'Legal', links: ['Privacy', 'Terms', 'Security', 'Status'] },
          ].map((column) => (
            <nav key={column.heading} aria-labelledby={`footer-${column.heading.toLowerCase()}`}>
              <h2 id={`footer-${column.heading.toLowerCase()}`} className="text-sm font-semibold">
                {column.heading}
              </h2>
              <ul role="list" className="mt-4 space-y-2.5 text-sm">
                {column.links.map((link) => (
                  <li key={link}>
                    <a href="#" className="text-gray-600 hover:text-gray-950 dark:text-gray-400 dark:hover:text-white">
                      {link}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          ))}
        </div>
        <p className="mx-auto mt-10 max-w-6xl border-t border-gray-200 pt-6 text-sm text-gray-500 dark:border-white/10 dark:text-gray-400">
          © 2026 Financeai. Nothing here is financial advice.
        </p>
      </footer>
    </div>
  )
}
