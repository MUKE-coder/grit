/*
 * A platform page for a startup audience: hero, capability tiles, a bento of
 * small proofs, testimonials and a grouped FAQ.
 *
 * One file, because the registry installs one file per block. It is a page you
 * own, not a component to configure — split it into the sections you want.
 *
 * What makes this one different from the other landing pages here is the FAQ.
 * The questions are grouped under real headings rather than dumped into one
 * list of fourteen. Someone looking for a shipping answer should be able to
 * jump to Shipping by heading and read three items, instead of reading past
 * six billing questions to find it. Each group is an <h3> with its own list,
 * so heading navigation lands where the labels promise.
 *
 * Every disclosure is <details>/<summary>. That gets keyboard support, the
 * expanded state announced, Ctrl+F opening a closed section, and it works
 * before hydration. A div with useState gets none of those.
 *
 * The bento tiles hold small mocks rather than photographs, so there is no
 * external asset to 404 and nothing to strip before you use the page. Each is
 * aria-hidden: they illustrate the sentence beside them and hold invented
 * numbers, and read aloud they are a string of digits with no referent.
 *
 * The icon row uses text labels beside every glyph. An icon on its own is a
 * picture of a concept, and which concept is a guess that changes per reader.
 *
 * One <h1>, an <h2> per section, <h3> inside. Long pages are navigated by
 * heading.
 */

import type { ReactNode } from 'react'

const NAV = ['Product', 'Solutions', 'Pricing', 'Company']

const LOGOS = ['Northwind', 'Meridian', 'Kestrel', 'Lumen', 'Atlas', 'Verdant']

const CAPABILITIES = [
  {
    title: 'Payment links in a line',
    body: 'Generate a hosted checkout from an API call, or from the dashboard when you have not built one yet.',
    mock: 'checkout',
  },
  {
    title: 'Reconciliation that matches',
    body: 'Every payout maps back to the charges inside it, so the number in your bank matches the number in your books.',
    mock: 'ledger',
  },
  {
    title: 'Spend controls',
    body: 'Per-card limits, category rules and a running total that updates as the month goes rather than after it.',
    mock: 'limit',
  },
  {
    title: 'Scheduled everything',
    body: 'Retries, dunning and payout timing on a schedule you set, with the next run shown before it happens.',
    mock: 'schedule',
  },
]

const CHECKOUT_FEATURES = [
  { title: 'Fraud scoring', body: 'Risk assessed per charge, with the reasons attached rather than a bare score.' },
  { title: 'Automated retries', body: 'Failed charges retried on a schedule tuned to why they failed.' },
  { title: 'Predictive targeting', body: 'Identify high-value customers from behaviour you already record.' },
  { title: 'Smart personalisation', body: 'Checkout adapts to returning customers without a separate flow to maintain.' },
  { title: 'Real-time insight', body: 'Charges, refunds and disputes on one timeline as they land.' },
  { title: 'Cross-channel', body: 'The same balance whether the sale came from web, app or an invoice.' },
]

const BENTO = [
  { title: 'Consistent uptime', body: 'Backed by an SLA and a status page that shows the incident, not just the colour.', mock: 'uptime' },
  { title: 'Keyboard shortcuts', body: 'Every destination reachable without the mouse, with the shortcuts discoverable.', mock: 'keys' },
  { title: 'Currency conversion', body: 'Settle in one currency, charge in thirty, with the rate recorded per charge.', mock: 'currency' },
  { title: 'Audited access', body: 'Role-based access, audit trails and fine-grained controls that log who saw what.', mock: 'shield' },
  { title: 'Resource insight', body: 'Memory, queue depth and error rate at a glance, alerted on before they bite.', mock: 'memory' },
  { title: 'Fast by default', body: 'Edge-cached reads and hardware-accelerated rendering on the dashboard.', mock: 'speed' },
]

const TESTIMONIALS = [
  { quote: 'We moved off a homegrown billing service in a fortnight. The reconciliation alone paid for the migration.', name: 'Adam Wathan', role: 'CEO, Kestrel Labs' },
  { quote: 'We needed something that could handle a complex compliance surface without a dedicated team. It delivered exactly that.', name: 'Ghalie Lukose', role: 'Finance lead, Meridian' },
  { quote: 'Implementing this helped us create a more engaging checkout. Conversion moved four points in the first month.', name: 'Shadcn', role: 'Design engineer, Lumen' },
  { quote: 'The documentation is excellent and the defaults are sensible, which is a rarer combination than it should be.', name: 'Priya Nair', role: 'Staff engineer, Atlas' },
]

const FAQ_GROUPS = [
  {
    group: 'General',
    items: [
      { q: 'How long does onboarding take?', a: 'Most teams are taking live charges within a week. The long pole is your own compliance review, not the integration.' },
      { q: 'What payment methods do you accept?', a: 'Cards, direct debit, and the common regional wallets. New methods are enabled per account rather than per integration, so adding one is a settings change.' },
      { q: 'Can I change or cancel my plan?', a: 'Any time, from the dashboard. Changes take effect at the start of the next billing period and we prorate the difference.' },
    ],
  },
  {
    group: 'Payouts',
    items: [
      { q: 'Do you pay out internationally?', a: 'To most countries, in local currency. The fee and the expected arrival date are shown before you confirm.' },
      { q: 'What is your refund policy?', a: 'Refunds are issued back to the original method and reverse the associated fees, provided they are within ninety days of the charge.' },
      { q: 'How quickly do payouts arrive?', a: 'Two working days on the standard schedule, same day on the accelerated one. Both are shown per payout rather than as an average.' },
    ],
  },
]

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

function CheckoutMock() {
  return (
    <Panel className="p-4">
      <p className="text-xs font-medium text-gray-900 dark:text-white">Checkout setup</p>
      <div className="mt-3 space-y-2">
        {['Business details', 'Bank account', 'Card payments'].map((row, i) => (
          <div key={row} className="flex items-center gap-2">
            <span
              className={`flex size-4 items-center justify-center rounded-full text-[9px] text-white ${
                i < 2 ? 'bg-emerald-600' : 'bg-gray-300 dark:bg-white/20'
              }`}
            >
              {i < 2 ? '✓' : ''}
            </span>
            <span className="h-2 flex-1 rounded bg-gray-200 dark:bg-white/10" />
          </div>
        ))}
      </div>
      <span className="mt-4 block h-7 w-24 rounded-md bg-violet-600" />
    </Panel>
  )
}

function LedgerMock() {
  const rows = [
    ['Subscription', '$3,847.16'],
    ['Marketplace', '$1,120.40'],
    ['Refunds', '-$92.55'],
  ]
  return (
    <Panel className="p-4">
      <div className="flex items-baseline justify-between">
        <p className="text-xs text-gray-500 dark:text-gray-400">Payout</p>
        <p className="text-sm font-semibold text-gray-900 dark:text-white">$5,167.56</p>
      </div>
      <div className="mt-3 space-y-2">
        {rows.map(([label, value]) => (
          <div key={label} className="flex justify-between text-[11px]">
            <span className="text-gray-500 dark:text-gray-400">{label}</span>
            <span className="text-gray-700 tabular-nums dark:text-gray-200">{value}</span>
          </div>
        ))}
      </div>
    </Panel>
  )
}

function LimitMock() {
  return (
    <Panel className="p-4">
      <p className="text-xs font-medium text-amber-700 dark:text-amber-400">Spending limit</p>
      <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-gray-200 dark:bg-white/10">
        <span className="block h-full w-[62%] rounded-full bg-violet-600" />
      </div>
      <p className="mt-2 text-[11px] text-gray-500 dark:text-gray-400">$6,200 of $10,000 this month</p>
    </Panel>
  )
}

function ScheduleMock() {
  return (
    <Panel className="p-4">
      <p className="text-xs font-medium text-gray-900 dark:text-white">Next runs</p>
      <div className="mt-3 space-y-2">
        {['Payout · in 2 days', 'Dunning · tomorrow', 'Invoices · Friday'].map((row) => (
          <div key={row} className="flex items-center gap-2 text-[11px] text-gray-600 dark:text-gray-300">
            <span className="size-1.5 rounded-full bg-violet-600" />
            {row}
          </div>
        ))}
      </div>
    </Panel>
  )
}

function UptimeMock() {
  return (
    <div aria-hidden="true" className="flex items-end gap-0.5 select-none">
      {Array.from({ length: 28 }, (_, i) => (
        <span
          key={i}
          style={{ height: `${60 + ((i * 37) % 40)}%` }}
          className={`w-1.5 rounded-sm ${i === 19 ? 'bg-amber-400' : 'bg-violet-500'}`}
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
      className="flex size-14 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-violet-700 select-none"
    >
      <span className="text-xl text-white">🛡</span>
    </div>
  )
}

function MemoryMock() {
  return (
    <div aria-hidden="true" className="select-none">
      <p className="text-[11px] text-gray-500 dark:text-gray-400">Memory 14 GB / 128 GB</p>
      <div className="mt-2 h-1.5 overflow-hidden rounded-full bg-gray-200 dark:bg-white/10">
        <span className="block h-full w-[11%] rounded-full bg-emerald-500" />
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
      <span className="text-xl text-violet-400">⚡</span>
    </div>
  )
}

const MOCKS: Record<string, ReactNode> = {
  checkout: <CheckoutMock />,
  ledger: <LedgerMock />,
  limit: <LimitMock />,
  schedule: <ScheduleMock />,
  uptime: <UptimeMock />,
  keys: <KeysMock />,
  currency: <CurrencyMock />,
  shield: <ShieldMock />,
  memory: <MemoryMock />,
  speed: <SpeedMock />,
}

function initials(name: string) {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageStartupPlatform({
  brand = 'Ledgerline',
  eyebrow = 'For startups',
  title = 'Simple cash flow for startups',
  subtitle = 'Take payments, reconcile them and pay yourself, without wiring four services together and hoping the totals agree.',
}: {
  brand?: string
  eyebrow?: string
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
            className="text-lg font-semibold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:text-white"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:text-gray-300 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <a
            href="#"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg border border-gray-300 px-4 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
          >
            Sign in
          </a>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-gray-200 dark:border-white/10">
          <div
            aria-hidden="true"
            className="absolute inset-0 bg-gradient-to-b from-violet-50 to-white dark:from-violet-500/10 dark:to-gray-950"
          />
          <div className="relative mx-auto max-w-3xl px-6 py-20 text-center">
            <p className="text-sm font-semibold tracking-wide text-violet-700 uppercase dark:text-violet-400">
              {eyebrow}
            </p>
            <h1 className="mt-4 text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h1>
            <p className="mx-auto mt-6 max-w-xl text-lg text-pretty text-gray-600 dark:text-gray-300">
              {subtitle}
            </p>
            <div className="mt-8 flex flex-wrap justify-center gap-3">
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-lg bg-violet-700 px-6 text-sm font-medium text-white hover:bg-violet-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600"
              >
                Start for free
              </a>
              <a
                href="#pricing"
                className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
              >
                Talk to sales
              </a>
            </div>

            <div className="mx-auto mt-12 max-w-sm">{MOCKS.checkout}</div>
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
                  className="text-sm font-semibold tracking-tight text-gray-400 dark:text-gray-500"
                >
                  {logo}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Capabilities */}
        <section
          aria-labelledby="capabilities"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-2xl text-center">
              <h2
                id="capabilities"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Ship with confidence on one platform
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                The pieces that usually take a quarter to wire together, already wired together.
              </p>
            </div>

            <ul role="list" className="mt-12 grid gap-6 md:grid-cols-2">
              {CAPABILITIES.map((item) => (
                <li
                  key={item.title}
                  className="flex flex-col rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40"
                >
                  <div className="mb-5">{MOCKS[item.mock]}</div>
                  <h3 className="mt-auto font-semibold text-gray-900 dark:text-white">
                    {item.title}
                  </h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{item.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Checkout features */}
        <section aria-labelledby="checkout" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-2xl text-center">
              <h2
                id="checkout"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Checkout, enterprise grade
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                The same checkout whether you take ten payments a day or ten thousand.
              </p>
            </div>

            <ul role="list" className="mt-12 grid gap-x-8 gap-y-10 sm:grid-cols-2 lg:grid-cols-3">
              {CHECKOUT_FEATURES.map((feature) => (
                <li key={feature.title}>
                  {/* The dot is decoration; the label beside it is the content.
                      An icon alone is a picture of a concept and which concept
                      is a guess that changes per reader. */}
                  <span
                    aria-hidden="true"
                    className="mb-3 flex size-9 items-center justify-center rounded-lg bg-violet-100 dark:bg-violet-500/15"
                  >
                    <span className="size-2 rounded-full bg-violet-700 dark:bg-violet-400" />
                  </span>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {feature.title}
                  </h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{feature.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Bento */}
        <section
          aria-labelledby="blocks"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <h2
                id="blocks"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Built out of smaller, sturdier blocks
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Everything you need to launch and scale, designed for speed, reliability and a
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

        {/* Testimonials */}
        <section
          aria-labelledby="testimonials"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <p className="text-sm font-semibold tracking-wide text-violet-700 uppercase dark:text-violet-400">
                Testimonials
              </p>
              <h2
                id="testimonials"
                className="mt-3 text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                You are in good company
              </h2>
            </div>

            <ul role="list" className="mt-12 gap-5 sm:columns-2">
              {TESTIMONIALS.map((item) => (
                <li key={item.name} className="mb-5 break-inside-avoid">
                  <figure className="rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900">
                    <blockquote className="text-sm text-pretty text-gray-700 dark:text-gray-200">
                      {item.quote}
                    </blockquote>
                    <figcaption className="mt-4 flex items-center gap-3">
                      <span
                        aria-hidden="true"
                        className="flex size-9 shrink-0 items-center justify-center rounded-full bg-gray-100 text-xs font-semibold text-gray-600 dark:bg-white/10 dark:text-gray-300"
                      >
                        {initials(item.name)}
                      </span>
                      <span>
                        <span className="block text-sm font-medium text-gray-900 dark:text-white">
                          {item.name}
                        </span>
                        <span className="block text-xs text-gray-500 dark:text-gray-400">
                          {item.role}
                        </span>
                      </span>
                    </figcaption>
                  </figure>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Grouped FAQ */}
        <section aria-labelledby="faq" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-3xl px-6 py-20">
            <div className="text-center">
              <h2
                id="faq"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Frequently asked questions
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Answers to what teams ask before they switch.
              </p>
            </div>

            {/* Grouped under real headings. Someone after a payout answer can
                jump to Payouts and read three items rather than reading past
                the billing questions to find it. */}
            <div className="mt-10 space-y-10">
              {FAQ_GROUPS.map((group) => (
                <section key={group.group} aria-labelledby={`faq-${group.group}`}>
                  <h3
                    id={`faq-${group.group}`}
                    className="text-sm font-semibold tracking-wide text-gray-500 uppercase dark:text-gray-400"
                  >
                    {group.group}
                  </h3>
                  <div className="mt-3">
                    {group.items.map((item) => (
                      <details
                        key={item.q}
                        className="group border-b border-gray-200 py-4 dark:border-white/10"
                      >
                        <summary className="flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:text-white">
                          {item.q}
                          <span
                            aria-hidden="true"
                            className="mt-0.5 shrink-0 text-gray-400 transition-transform group-open:rotate-45 motion-reduce:transition-none"
                          >
                            +
                          </span>
                        </summary>
                        <p className="mt-2.5 text-sm text-pretty text-gray-600 dark:text-gray-300">
                          {item.a}
                        </p>
                      </details>
                    ))}
                  </div>
                </section>
              ))}
            </div>

            <p className="mt-10 text-center text-sm text-gray-600 dark:text-gray-300">
              Cannot find what you are looking for?{' '}
              <a
                href="#"
                className="font-medium text-violet-700 underline-offset-4 hover:underline dark:text-violet-400"
              >
                Contact support
              </a>
              .
            </p>
          </div>
        </section>

        {/* Closing ask */}
        <section id="pricing" aria-labelledby="cta" className="bg-gray-50 dark:bg-gray-900/40">
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
              className="mt-8 inline-flex min-h-12 items-center rounded-lg bg-violet-700 px-6 text-sm font-medium text-white hover:bg-violet-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600"
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
              Payments, reconciliation and payouts in one place.
            </p>
          </div>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
              {[...NAV, 'Privacy', 'Terms'].map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-600 dark:text-gray-400 dark:hover:text-white"
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
