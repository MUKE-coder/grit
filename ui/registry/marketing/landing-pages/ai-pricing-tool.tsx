/*
 * A product page for a pricing tool: hero with a wide app shot, a feature
 * bento, one pull quote, a three-step workflow and a testimonial wall.
 *
 * One file, because the registry installs one file per block. Split it into
 * the sections you want; it is a starting point you own.
 *
 * The distinct part here is the workflow. Three steps, numbered, as an <ol>,
 * each with its own illustration — because the claim being made is sequence,
 * and a grid of three cards makes no claim about order at all. The connector
 * between them is drawn with a pseudo-element on the list item rather than
 * with separate arrow elements, so there is nothing decorative in the DOM to
 * mark hidden and nothing that lands between two steps in the reading order.
 *
 * The single quote sits between the features and the workflow on purpose. One
 * person saying one specific thing, at the point the reader has seen enough to
 * be sceptical, does more than a wall of five-star cards. The wall comes later,
 * for whoever wants volume.
 *
 * Every mock is markup and aria-hidden. They illustrate the sentence beside
 * them and hold invented figures, which read aloud as digits with no referent.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Product', 'Solutions', 'Pricing', 'Company']

const LOGOS = ['Northwind', 'Meridian', 'Kestrel', 'Lumen', 'Atlas', 'Verdant']

const FEATURES = [
  {
    title: 'Pricing performance analytics',
    body: 'Monitor conversion, ARPU and churn against cohort charts and device breakdowns, so a price change has a before and an after.',
    mock: 'chart',
    wide: true,
  },
  {
    title: 'AI pricing recommendations',
    body: 'Generated summaries of what is working, plus suggested experiments, copy and layout tweaks for each plan.',
    mock: 'suggest',
    wide: false,
  },
  {
    title: 'Localised pricing at scale',
    body: 'Auto-translate plans and regional messaging into a hundred and eighty languages with currency and tax-inclusive formatting.',
    mock: 'locales',
    wide: false,
  },
  {
    title: 'One-click checkout links',
    body: 'Share a secure payment link that remembers saved details and verifies the code before it charges anyone.',
    mock: 'link',
    wide: true,
  },
]

const WORKFLOW = [
  {
    title: 'Collect',
    body: 'Import data from the sources and formats you already have. No schema to design first.',
    mock: 'collect',
  },
  {
    title: 'Analyse',
    body: 'The model reads the set, flags the patterns and proposes the experiments worth running.',
    mock: 'analyse',
  },
  {
    title: 'Act',
    body: 'Turn the findings into live price changes, with a rollback that is one click and not a migration.',
    mock: 'act',
  },
]

const BENEFITS = [
  { title: 'Team collaboration', body: 'Discuss in context, mention teammates and resolve threads without leaving the plan you are editing.' },
  { title: 'Workflow automation', body: 'Trigger agents from events and chain tools with conditions, so routine changes stop needing a person.' },
  { title: 'Omnichannel messaging', body: 'Send campaigns across email, chat and more, from one unified inbox rather than four tabs.' },
  { title: 'Enterprise-grade security', body: 'Role-based access, audit trails and fine-grained controls, with data kept in the region you choose.' },
  { title: 'Collaborative analysis', body: 'Turn scattered signals into shared understanding with insights the whole team can explore.' },
]

const TESTIMONIALS = [
  { quote: 'The component library has been a game-changer for our development team. We build consistent interfaces across our payment platform with minimal effort.', name: 'Adam Wathan', role: 'CEO, Kestrel Labs' },
  { quote: 'We needed something that could handle our compliance requirements while keeping performance. It was delivered exactly that, and integrated with our existing architecture.', name: 'Ghalie Lukose', role: 'Frontend engineer, Meridian' },
  { quote: 'Implementing this helped us create a more engaging consumer interface. The responsive system works flawlessly across devices.', name: 'Shadcn', role: 'Design engineer, Lumen' },
  { quote: 'The documentation is excellent and the customisation options are exactly what we needed, which is a rarer pairing than it sounds.', name: 'Priya Nair', role: 'Staff engineer, Atlas' },
]

const PULL_QUOTE = {
  quote:
    'Using this has been like unlocking a secret design superpower. It is the perfect fusion of simplicity and versatility, and it lets us create interfaces that are as stunning as they are user-friendly.',
  name: 'Méschac Ngandu',
  role: 'UI engineer',
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

function AppShotMock() {
  return (
    <Panel>
      <div className="flex items-center gap-1.5 border-b border-gray-200 px-3 py-2.5 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="flex">
        <div className="hidden w-40 shrink-0 space-y-2 border-r border-gray-200 p-3 sm:block dark:border-white/10">
          {Array.from({ length: 7 }, (_, i) => (
            <span
              key={i}
              className={`block h-2 rounded ${
                i === 1 ? 'w-4/5 bg-sky-300 dark:bg-sky-500/40' : 'w-full bg-gray-200 dark:bg-white/10'
              }`}
            />
          ))}
        </div>
        <div className="flex-1 space-y-2 p-3">
          {Array.from({ length: 11 }, (_, i) => (
            <div key={i} className="flex items-center gap-2">
              <span className="size-2 shrink-0 rounded-sm bg-gray-200 dark:bg-white/10" />
              <span className="h-2 flex-1 rounded bg-gray-100 dark:bg-white/5" />
              <span
                className={`h-2 w-10 shrink-0 rounded ${
                  i % 4 === 0 ? 'bg-emerald-200 dark:bg-emerald-500/30' : 'bg-gray-100 dark:bg-white/5'
                }`}
              />
            </div>
          ))}
        </div>
      </div>
    </Panel>
  )
}

function ChartMock() {
  return (
    <Panel className="p-4">
      <div className="flex h-28 items-end gap-1.5">
        {[35, 52, 44, 68, 58, 76, 64, 88, 72, 95].map((height, i) => (
          <span
            key={i}
            style={{ height: `${height}%` }}
            className="flex-1 rounded-t bg-gradient-to-t from-sky-600 to-sky-400"
          />
        ))}
      </div>
      <div className="mt-3 flex gap-3 text-[11px] text-gray-500 dark:text-gray-400">
        <span>Conversion</span>
        <span className="text-emerald-600 dark:text-emerald-400">+18%</span>
      </div>
    </Panel>
  )
}

function SuggestMock() {
  return (
    <Panel className="p-4">
      <p className="text-[11px] text-gray-500 dark:text-gray-400">Suggested</p>
      <div className="mt-2 space-y-2">
        {['Raise Pro to $29', 'Add an annual toggle', 'Trim the free tier'].map((row) => (
          <div
            key={row}
            className="rounded-lg border border-gray-200 px-2.5 py-1.5 text-[11px] text-gray-700 dark:border-white/10 dark:text-gray-200"
          >
            {row}
          </div>
        ))}
      </div>
    </Panel>
  )
}

function LocalesMock() {
  return (
    <Panel className="p-4">
      <div className="space-y-1.5">
        {[
          ['English', '$29'],
          ['Français', '27 €'],
          ['日本語', '¥4,400'],
        ].map(([label, price]) => (
          <div key={label} className="flex justify-between text-[11px]">
            <span className="text-gray-500 dark:text-gray-400">{label}</span>
            <span className="tabular-nums text-gray-700 dark:text-gray-200">{price}</span>
          </div>
        ))}
      </div>
    </Panel>
  )
}

function LinkMock() {
  return (
    <Panel className="p-4">
      <p className="text-[11px] font-medium text-gray-900 dark:text-white">Checkout with link</p>
      <div className="mt-3 flex items-center gap-2 rounded-lg border border-gray-200 px-2.5 py-2 dark:border-white/10">
        <span className="h-2 flex-1 rounded bg-gray-200 dark:bg-white/10" />
        <span className="rounded bg-sky-600 px-2 py-1 text-[10px] text-white">Copy</span>
      </div>
      <p className="mt-2 text-[10px] text-gray-500 dark:text-gray-400">Powered by secure checkout</p>
    </Panel>
  )
}

function CollectMock() {
  return (
    <Panel className="p-4">
      <div className="grid grid-cols-3 gap-1.5">
        {Array.from({ length: 9 }, (_, i) => (
          <span key={i} className="block h-5 rounded bg-gray-100 dark:bg-white/5" />
        ))}
      </div>
    </Panel>
  )
}

function AnalyseMock() {
  return (
    <Panel className="p-4">
      <div className="flex items-end gap-1">
        {[30, 60, 45, 80, 55, 70].map((h, i) => (
          <span
            key={i}
            style={{ height: `${h * 0.5}px` }}
            className={`flex-1 rounded-t ${i === 3 ? 'bg-sky-600' : 'bg-sky-200 dark:bg-sky-500/30'}`}
          />
        ))}
      </div>
    </Panel>
  )
}

function ActMock() {
  return (
    <Panel className="p-4">
      <div className="space-y-2">
        <div className="flex items-center justify-between rounded-lg bg-emerald-50 px-2.5 py-2 text-[11px] text-emerald-800 dark:bg-emerald-500/10 dark:text-emerald-300">
          Price live
          <span>✓</span>
        </div>
        <span className="block h-2 w-2/3 rounded bg-gray-200 dark:bg-white/10" />
      </div>
    </Panel>
  )
}

const MOCKS: Record<string, ReactNode> = {
  chart: <ChartMock />,
  suggest: <SuggestMock />,
  locales: <LocalesMock />,
  link: <LinkMock />,
  collect: <CollectMock />,
  analyse: <AnalyseMock />,
  act: <ActMock />,
}

function initials(name: string) {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageAiPricingTool({
  brand = 'Tenor',
  title = 'Cutting-edge tools to build AI-based pricing',
  subtitle = 'Manage your sales team with tools that let you test, price and ship changes without waiting on an engineer.',
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
        {/* Hero */}
        <section className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
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
                  href="#workflow"
                  className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
                >
                  How it works
                </a>
              </div>
            </div>

            <div className="mt-14">
              <AppShotMock />
            </div>
          </div>
        </section>

        {/* Logos */}
        <section aria-labelledby="logos" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-12">
            <h2 id="logos" className="sr-only">
              Companies using {brand}
            </h2>
            <ul
              role="list"
              className="grid grid-cols-2 items-center justify-items-center gap-6 sm:grid-cols-3 lg:grid-cols-6"
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

        {/* Features */}
        <section aria-labelledby="features" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <h2
                id="features"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Built to price in the AI era
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Fast, considered tooling that helps you set a number and defend it.
              </p>
            </div>

            {/* Six columns so a wide tile can take four and a narrow one two,
                which tiles exactly. Spans that do not tile leave CSS grid
                dropping cards into whatever hole it finds. */}
            <ul role="list" className="mt-12 grid gap-6 lg:grid-cols-6">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`flex flex-col rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40 ${
                    feature.wide ? 'lg:col-span-4' : 'lg:col-span-2'
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

        {/* Benefits */}
        <section
          aria-labelledby="benefits"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase dark:text-sky-400">
                Benefits
              </p>
              <h2
                id="benefits"
                className="mt-3 text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Tools that cut the busywork out of subscriptions
              </h2>
            </div>

            <ul role="list" className="mt-10 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {BENEFITS.map((benefit) => (
                <li
                  key={benefit.title}
                  className="rounded-2xl border border-gray-200 bg-white p-6 dark:border-white/10 dark:bg-gray-900"
                >
                  <h3 className="font-semibold text-gray-900 dark:text-white">{benefit.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{benefit.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Pull quote */}
        <section aria-labelledby="quote" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-3xl px-6 py-20 text-center">
            <h2 id="quote" className="sr-only">
              What one customer says
            </h2>
            <figure>
              <blockquote className="text-xl leading-relaxed text-balance text-gray-900 sm:text-2xl dark:text-white">
                {PULL_QUOTE.quote}
              </blockquote>
              <figcaption className="mt-6 flex items-center justify-center gap-3">
                <span
                  aria-hidden="true"
                  className="flex size-10 items-center justify-center rounded-full bg-gray-100 text-xs font-semibold text-gray-600 dark:bg-white/10 dark:text-gray-300"
                >
                  {initials(PULL_QUOTE.name)}
                </span>
                <span className="text-left">
                  <span className="block text-sm font-medium text-gray-900 dark:text-white">
                    {PULL_QUOTE.name}
                  </span>
                  <span className="block text-xs text-gray-500 dark:text-gray-400">
                    {PULL_QUOTE.role}
                  </span>
                </span>
              </figcaption>
            </figure>
          </div>
        </section>

        {/* Workflow */}
        <section
          id="workflow"
          aria-labelledby="workflow-heading"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-5xl px-6 py-20">
            <div className="mx-auto max-w-2xl text-center">
              <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase dark:text-sky-400">
                Our process
              </p>
              <h2
                id="workflow-heading"
                className="mt-3 text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                A simple three-step workflow
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                From raw data to a live price change, without a quarter of engineering in between.
              </p>
            </div>

            {/* An ordered list, because the claim is sequence. Three cards in a
                grid say nothing about order. The connector is a pseudo-element
                on the item, so no decorative node sits between two steps in
                the reading order. */}
            <ol className="mt-14 grid gap-10 sm:grid-cols-3">
              {WORKFLOW.map((step, index) => (
                <li
                  key={step.title}
                  className="relative text-center sm:before:absolute sm:before:top-5 sm:before:left-[calc(50%+2rem)] sm:before:h-px sm:before:w-[calc(100%-4rem)] sm:before:bg-gray-200 sm:last:before:hidden sm:dark:before:bg-white/10"
                >
                  <span
                    aria-hidden="true"
                    className="relative z-10 mx-auto flex size-10 items-center justify-center rounded-full bg-sky-700 text-sm font-bold text-white"
                  >
                    {index + 1}
                  </span>
                  {/* Fixed height, contents centred. The three mocks have
                      different natural heights (106, 74 and 82px), so without
                      a common well the step headings land at three different
                      vertical positions and the row reads as misaligned. */}
                  {/* grid, not flex: grid items stretch on the inline axis by
                      default, so the mock keeps its full width while
                      items-center handles the vertical centring. flex would
                      shrink each panel to its content width. */}
                  <div className="mt-6 grid h-28 items-center">{MOCKS[step.mock]}</div>
                  <h3 className="mt-5 font-semibold text-gray-900 dark:text-white">{step.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{step.body}</p>
                </li>
              ))}
            </ol>

            <div className="mt-12 text-center">
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-lg bg-sky-700 px-6 text-sm font-medium text-white hover:bg-sky-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-sky-600"
              >
                Get started
              </a>
            </div>
          </div>
        </section>

        {/* Testimonial wall */}
        <section
          aria-labelledby="testimonials"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="max-w-2xl">
              <p className="text-sm font-semibold tracking-wide text-sky-700 uppercase dark:text-sky-400">
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
              Join a community of over a thousand companies and developers already pricing with
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
              Pricing tooling for teams that ship.
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
