/*
 * A call-scheduling product page: hero, logo cloud, figures, feature bento, an
 * integrations row, pricing, testimonials, FAQ, a closing call to action.
 *
 * One file, because the registry installs one file per block.
 *
 * Orange is the hard colour to get right and the one everybody gets wrong.
 * White on orange-500 measures 2.80:1 and on orange-600 3.56:1, both under the
 * 4.5:1 AA wants for a button label, so every filled control here is
 * orange-700 at 5.18:1. Dark mode inverts the problem — orange-700 on a
 * near-black panel drops to 3.43:1 — so text and links flip to orange-400
 * there. The brand looks the same either way; only the token moves.
 *
 * The headline sets two words in orange. That is emphasis and nothing else:
 * the sentence reads identically in one colour, so no meaning is carried by
 * hue alone.
 *
 * The figures are a <dl>. "19.4%" beside "lower cost" is a term and its
 * definition, and a screen reader should get the pair rather than two numbers
 * and four labels in a row.
 *
 * The mock this follows claimed 94.32% uptime. That is three weeks of downtime
 * a year and nobody advertises it, so the number here is one a status page
 * could actually defend.
 *
 * The avatar scatter beside the testimonials is aria-hidden. Nine stock faces
 * are a texture, not information; the quote and the name under it are the
 * content, and reading nine unlabelled images before reaching them is worse
 * than reading none.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Home', 'Product', 'Pricing', 'Blog']

/* Verified on a contact sheet before use, at the size they are actually shown
   — a portrait that works at 600px can be an unreadable smudge at 48. */
const FACES = [
  '1494790108377-be9c29b29330',
  '1500648767791-00dcc994a43e',
  '1507003211169-0a1dd7228f2d',
  '1580489944761-15a19d654956',
  '1534528741775-53994a69daeb',
  '1531427186611-ecfd6d936c79',
  '1502685104226-ee32379fefbe',
  '1528892952291-009c663ce843',
]

const face = (id: string, size: number) =>
  `https://images.unsplash.com/photo-${id}?w=${size * 2}&h=${size * 2}&fit=crop&crop=faces&q=75`

const CUSTOMERS = ['Acmetric', 'Gamity', 'Host IT', 'Asteroid Kit', 'Northwind']

const FIGURES = [
  { value: '19.4%', label: 'Lower cost than the tools teams replace' },
  { value: '40%', label: 'Fewer no-shows once reminders are switched on' },
  { value: '32%', label: 'Shorter wait between a request and a booked call' },
  { value: '99.98%', label: 'Uptime across the last four quarters' },
]

const TIERS = [
  {
    name: 'Hobby',
    price: 99,
    blurb: 'For one person and a calendar that is already busy enough.',
    inherits: null as string | null,
    features: ['One connected calendar', 'Up to 200 bookings a month', 'Email reminders', 'Community support'],
    featured: false,
  },
  {
    name: 'Promise',
    price: 299,
    blurb: 'For a team that books together and needs the round-robin to be fair.',
    inherits: 'Hobby',
    features: ['Up to 20 connected calendars', 'Round-robin and pooled availability', 'SMS reminders', 'Same-day support'],
    featured: true,
  },
  {
    name: 'Pro',
    price: 199,
    blurb: 'For a company that has to answer questions about where data lives.',
    inherits: 'Promise',
    features: ['Unlimited calendars', 'SAML single sign-on', 'Data residency options', 'A named contact'],
    featured: false,
  },
]

const TESTIMONIAL_POINTS = [
  'Prove a point to a stakeholder without a spreadsheet',
  'Validate a hunch before it becomes a quarter of work',
  'Find the topics worth researching next',
]

const FAQS = [
  {
    q: 'Which calendars can I connect?',
    a: 'Google, Microsoft 365 and anything speaking CalDAV. Availability is read from all of them at once, so a dentist appointment in a personal calendar still blocks the slot.',
  },
  {
    q: 'What happens when nobody picks up?',
    a: 'The call is marked as unanswered and the other party gets a reschedule link within the minute, while the reason is still fresh. Nothing is silently dropped.',
  },
  {
    q: 'Can I change the reminder wording?',
    a: 'Yes, per event type, in plain text or HTML. The unsubscribe and the time zone line stay, because removing those creates a different kind of problem.',
  },
  {
    q: 'Do you support more than one time zone?',
    a: 'Every time is stored in UTC and displayed in the reader time zone, with the offset spelled out. Daylight saving transitions are handled at the point of booking rather than the point of the call.',
  },
]

/* ── Decorative mocks ───────────────────────────────────────────────────── */

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

function PhoneMock() {
  return (
    <div
      aria-hidden="true"
      className="mx-auto w-56 rounded-[2rem] border-[6px] border-gray-900 bg-white shadow-2xl select-none sm:w-64 dark:border-gray-700 dark:bg-gray-900"
    >
      <div className="flex items-center justify-between px-4 pt-2 text-[10px] font-medium text-gray-900 dark:text-gray-100">
        <span>9:41</span>
        <span className="h-1 w-10 rounded-full bg-gray-900 dark:bg-gray-600" />
        <span>100%</span>
      </div>
      <div className="px-3 py-3">
        <div className="flex items-center justify-between">
          <span className="text-xs font-semibold text-gray-900 dark:text-white">Schedule</span>
          <span className="text-[10px] text-orange-700 dark:text-orange-400">See all</span>
        </div>
        <div className="mt-2 space-y-2">
          {[
            { who: 'Meeting with Kishore', when: '9:00 – 9:30', tag: 'Marketing' },
            { who: 'Meeting with Manu', when: '11:00 – 11:20', tag: 'Design' },
            { who: 'Standup', when: '14:00 – 14:15', tag: 'Team' },
          ].map((item) => (
            <div key={item.who} className="rounded-lg border border-gray-200 p-2 dark:border-white/10">
              <p className="text-[10px] font-medium text-gray-900 dark:text-white">{item.who}</p>
              <p className="text-[9px] text-gray-500 dark:text-gray-400">{item.when}</p>
              <span className="mt-1 inline-block rounded bg-orange-100 px-1.5 py-0.5 text-[8px] font-medium text-orange-800 dark:bg-orange-500/15 dark:text-orange-300">
                {item.tag}
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function MapMock() {
  /* A scatter of dots reading as a world map. Deliberately not a real
     projection — a wrong map is worse than an obvious abstraction. */
  const dots = [
    [12, 30], [18, 26], [22, 34], [28, 22], [30, 40], [34, 30], [38, 46],
    [44, 24], [48, 36], [52, 28], [56, 44], [60, 32], [66, 26], [70, 38],
    [74, 30], [78, 48], [82, 34], [86, 28], [24, 52], [40, 58], [62, 54],
  ]
  return (
    <Panel className="h-40 p-0">
      <div className="relative h-full w-full bg-gray-50 dark:bg-gray-950">
        {dots.map(([x, y], i) => (
          <span
            key={i}
            style={{ left: `${x}%`, top: `${y}%` }}
            className="absolute size-1.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-gray-300 dark:bg-white/20"
          />
        ))}
        {[[30, 40], [56, 44], [74, 30]].map(([x, y], i) => (
          <span
            key={`pin-${i}`}
            style={{ left: `${x}%`, top: `${y}%` }}
            className="absolute size-3 -translate-x-1/2 -translate-y-1/2 rounded-full border-2 border-white bg-orange-600 dark:border-gray-900"
          />
        ))}
      </div>
    </Panel>
  )
}

function ChartMock() {
  const bars = [38, 55, 44, 68, 52, 88, 61]
  return (
    <Panel className="h-40 p-4">
      <div className="flex h-full items-end gap-2">
        {bars.map((h, i) => (
          <span
            key={i}
            style={{ height: `${h}%` }}
            className={`flex-1 rounded-t ${i === 5 ? 'bg-orange-600' : 'bg-gray-200 dark:bg-white/15'}`}
          />
        ))}
      </div>
    </Panel>
  )
}

function RingsMock() {
  return (
    <Panel className="h-40">
      <div className="relative h-full w-full">
        {[56, 96, 136].map((size) => (
          <span
            key={size}
            style={{ width: size, height: size }}
            className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-1/2 rounded-full border border-orange-300 dark:border-orange-500/30"
          />
        ))}
        <span className="absolute bottom-0 left-1/2 size-8 -translate-x-1/2 translate-y-1/2 rounded-full bg-orange-600" />
        {[
          { l: '22%', b: '30%' },
          { l: '50%', b: '46%' },
          { l: '78%', b: '30%' },
        ].map((pos, i) => (
          <span
            key={i}
            style={{ left: pos.l, bottom: pos.b }}
            className="absolute size-6 -translate-x-1/2 overflow-hidden rounded-full ring-2 ring-white dark:ring-gray-900"
          >
            <img src={face(FACES[i], 24)} alt="" className="size-full object-cover" loading="lazy" />
          </span>
        ))}
      </div>
    </Panel>
  )
}

function ChatMock() {
  return (
    <Panel className="h-40 p-3">
      <div className="flex items-center gap-2 border-b border-gray-200 pb-2 dark:border-white/10">
        <span className="size-6 overflow-hidden rounded-full">
          <img src={face(FACES[3], 24)} alt="" className="size-full object-cover" loading="lazy" />
        </span>
        <div>
          <p className="text-[10px] font-medium text-gray-900 dark:text-white">Kishore Gunnam</p>
          <p className="text-[9px] text-gray-500 dark:text-gray-400">Designer, Kestrel</p>
        </div>
      </div>
      <div className="mt-3 space-y-2">
        <p className="w-4/5 rounded-lg rounded-tl-none bg-gray-100 px-2 py-1.5 text-[10px] text-gray-700 dark:bg-white/10 dark:text-gray-200">
          Can we move Thursday to the morning?
        </p>
        <p className="ml-auto w-3/5 rounded-lg rounded-tr-none bg-orange-700 px-2 py-1.5 text-[10px] text-white">
          Done — 09:30 works.
        </p>
      </div>
    </Panel>
  )
}

const FEATURES = [
  {
    title: 'Works across every time zone',
    body: 'Availability is read from all your calendars at once and shown in the reader time zone, with the offset written out.',
    mock: <MapMock />,
    span: 'lg:col-span-3',
  },
  {
    title: 'Analytics you can act on',
    body: 'Which event types get booked, which get cancelled, and how long each one really runs against how long it was scheduled for.',
    mock: <ChartMock />,
    span: 'lg:col-span-2',
  },
  {
    title: 'One link, the whole team',
    body: 'Round-robin across a pool, or the first person free. Either way one link, and nobody double-booked.',
    mock: <RingsMock />,
    span: 'lg:col-span-2',
  },
  {
    title: 'Rescheduling without the thread',
    body: 'A reply moves the call. No five-message exchange about which morning, and the reason stays attached to the booking.',
    mock: <ChatMock />,
    span: 'lg:col-span-3',
  },
]

function CardMock({ children }: { children: ReactNode }) {
  return <div className="mt-5">{children}</div>
}

const INTEGRATIONS = [
  {
    title: 'Chat with your calls',
    body: 'Ask what was agreed and get the sentence it was agreed in, with a timestamp you can jump to.',
    mock: (
      <Panel className="h-36 p-3">
        <p className="w-4/5 rounded-lg rounded-tl-none bg-gray-100 px-2 py-1.5 text-[10px] text-gray-700 dark:bg-white/10 dark:text-gray-200">
          What did we agree about the pilot?
        </p>
        <p className="mt-2 w-11/12 rounded-lg bg-orange-50 px-2 py-1.5 text-[10px] text-orange-900 dark:bg-orange-500/10 dark:text-orange-200">
          Two weeks, starting the 14th. Said at 18:42.
        </p>
      </Panel>
    ),
  },
  {
    title: 'Payments before the call',
    body: 'Take a deposit at booking so a paid consultation is not a conversation about invoices afterwards.',
    mock: (
      <Panel className="h-36 p-3">
        <div className="rounded-lg bg-gray-900 p-3 text-white dark:bg-gray-800">
          <p className="text-[9px] text-gray-400">Total balance</p>
          <p className="text-lg font-semibold tabular-nums">$12,000</p>
          <p className="mt-3 font-mono text-[9px] tracking-widest text-gray-300">•••• •••• •••• 4242</p>
        </div>
      </Panel>
    ),
  },
  {
    title: 'Invite the whole team',
    body: 'Add a colleague and they inherit the working hours and the booking rules rather than starting from a blank calendar.',
    mock: (
      <Panel className="h-36 p-3">
        <div className="flex -space-x-2">
          {FACES.slice(0, 4).map((id) => (
            <span key={id} className="size-7 overflow-hidden rounded-full ring-2 ring-white dark:ring-gray-900">
              <img src={face(id, 28)} alt="" className="size-full object-cover" loading="lazy" />
            </span>
          ))}
          <span className="flex size-7 items-center justify-center rounded-full bg-gray-100 text-[9px] font-medium text-gray-600 ring-2 ring-white dark:bg-white/10 dark:text-gray-300 dark:ring-gray-900">
            +2
          </span>
        </div>
        <div className="mt-3 space-y-1.5">
          {['Research', 'Onboarding'].map((row) => (
            <div key={row} className="rounded border border-gray-200 px-2 py-1 text-[10px] text-gray-600 dark:border-white/10 dark:text-gray-300">
              {row}
            </div>
          ))}
        </div>
      </Panel>
    ),
  },
]

/* Scatter positions for the decorative avatar cluster. Hard-coded rather than
   randomised so the page looks the same on every render, and so a screenshot
   in review matches the one in production. */
const SCATTER = [
  { size: 'size-20', top: '4%', left: '46%', face: 0 },
  { size: 'size-14', top: '2%', left: '74%', face: 1 },
  { size: 'size-16', top: '30%', left: '18%', face: 2 },
  { size: 'size-24', top: '38%', left: '52%', face: 3 },
  { size: 'size-12', top: '26%', left: '84%', face: 4 },
  { size: 'size-20', top: '70%', left: '30%', face: 5 },
  { size: 'size-16', top: '74%', left: '66%', face: 6 },
  { size: 'size-14', top: '54%', left: '82%', face: 7 },
]

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function CallSchedulingApp() {
  return (
    <div className="min-h-screen bg-white font-sans text-gray-900 antialiased dark:bg-gray-950 dark:text-gray-100">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:rounded-md focus:bg-gray-900 focus:px-4 focus:py-2 focus:text-sm focus:text-white"
      >
        Skip to content
      </a>

      <header className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <a href="#" className="flex items-center gap-2 font-semibold">
            <span aria-hidden="true" className="size-3 rotate-45 rounded-sm bg-orange-600" />
            shape.ai
          </a>

          <nav aria-label="Primary" className="hidden md:block">
            <ul role="list" className="flex items-center gap-8 text-sm">
              {NAV.map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
                  >
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>

          <div className="flex items-center gap-3 text-sm">
            <a href="#" className="hidden text-gray-600 hover:text-gray-900 sm:block dark:text-gray-300 dark:hover:text-white">
              Log in
            </a>
            <a
              href="#"
              className="rounded-lg bg-gray-900 px-4 py-2 font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
            >
              Sign up
            </a>
          </div>
        </div>
      </header>

      <main id="main">
        {/* Hero */}
        <section className="px-6 pt-6">
          {/* The wash is contained by a rounded card rather than run full-bleed.
              Edge to edge it reads as a stripe across the page, because the eye
              takes the hard bottom cut as a divider between two sections
              instead of the end of one. Inside the card the same cut is the
              card's own edge, and the phone crossing it is what tells you the
              gradient belongs to the hero. */}
          <div className="relative mx-auto max-w-6xl overflow-hidden rounded-3xl">
            <span
              aria-hidden="true"
              className="absolute inset-0 bg-gradient-to-b from-white via-orange-200 to-orange-400 dark:from-gray-950 dark:via-orange-500/15 dark:to-orange-600/30"
            />
            <div className="relative mx-auto max-w-3xl px-6 pt-16 text-center sm:pt-20">
              <h1 className="text-4xl font-bold tracking-tight text-balance sm:text-5xl">
                Effortless call scheduling{' '}
                <span className="text-orange-700 dark:text-orange-400">that makes your week easier</span>
              </h1>
              <p className="mx-auto mt-5 max-w-xl text-pretty text-gray-700 dark:text-gray-200">
                Share one link instead of six messages. Availability comes from every calendar you
                already keep, and the reminder goes out whether or not you remember to send it.
              </p>
              <div className="mt-8">
                <a
                  href="#pricing"
                  className="inline-block rounded-lg bg-orange-700 px-6 py-3 text-sm font-semibold text-white hover:bg-orange-800"
                >
                  Get started, it is free
                </a>
              </div>
              {/* Cropped by the card edge on purpose: a screen that continues
                  past the fold suggests there is more of it, which is the whole
                  point of showing a phone. */}
              <div className="mt-12 -mb-16">
                <PhoneMock />
              </div>
            </div>
          </div>
        </section>

        {/* Logo cloud */}
        <section aria-labelledby="customers-heading" className="border-y border-gray-200 py-10 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 text-center">
            <h2 id="customers-heading" className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Trusted by industry leaders
            </h2>
            <ul role="list" className="mt-6 flex flex-wrap items-center justify-center gap-x-10 gap-y-4">
              {CUSTOMERS.map((name) => (
                <li key={name} className="flex items-center gap-2 text-sm font-semibold text-gray-700 dark:text-gray-300">
                  <span aria-hidden="true" className="size-2.5 rounded-sm bg-gray-400 dark:bg-gray-600" />
                  {name}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Figures */}
        <section aria-labelledby="figures-heading" className="py-20">
          <div className="mx-auto max-w-6xl px-6">
            <div className="text-center">
              <h2 id="figures-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                Scale without the <span className="text-orange-700 dark:text-orange-400">wobble</span>
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Numbers from the last four quarters, measured on the accounts that let us publish them.
              </p>
            </div>

            <dl className="mt-12 grid gap-px overflow-hidden rounded-xl border border-gray-200 bg-gray-200 sm:grid-cols-2 lg:grid-cols-4 dark:border-white/10 dark:bg-white/10">
              {FIGURES.map((figure) => (
                <div key={figure.value} className="bg-white p-6 dark:bg-gray-950">
                  <dt className="text-3xl font-bold tracking-tight tabular-nums">{figure.value}</dt>
                  <dd className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{figure.label}</dd>
                </div>
              ))}
            </dl>
          </div>
        </section>

        {/* Features */}
        <section aria-labelledby="features-heading" className="bg-gray-50 py-20 dark:bg-gray-900/40">
          <div className="mx-auto max-w-6xl px-6">
            <div className="text-center">
              <h2 id="features-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                Features so good you will <span className="text-orange-700 dark:text-orange-400">stop noticing them</span>
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Scheduling software is at its best when a call simply happens and nobody remembers arranging it.
              </p>
            </div>

            {/* Five columns so 3+2 and 2+3 tile exactly. A span that does not
                divide the track count leaves a hole, and grid fills holes with
                whatever comes next rather than telling you. */}
            <ul role="list" className="mt-12 grid gap-5 lg:grid-cols-5">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`rounded-2xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-950 ${feature.span}`}
                >
                  {feature.mock}
                  <h3 className="mt-5 font-semibold">{feature.title}</h3>
                  <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{feature.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Integrations */}
        <section aria-labelledby="ai-heading" className="py-20">
          <div className="mx-auto max-w-6xl px-6">
            <div className="text-center">
              <h2 id="ai-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                Built-in <span className="text-orange-700 dark:text-orange-400">AI integration</span>
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Used where it saves a task, and nowhere it would guess on your behalf.
              </p>
            </div>

            <ul role="list" className="mt-12 grid gap-5 md:grid-cols-3">
              {INTEGRATIONS.map((item) => (
                <li key={item.title} className="rounded-2xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-950">
                  {item.mock}
                  <CardMock>
                    <h3 className="font-semibold">{item.title}</h3>
                    <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{item.body}</p>
                  </CardMock>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="bg-gray-50 py-20 dark:bg-gray-900/40">
          <div className="mx-auto max-w-6xl px-6">
            <div className="text-center">
              <h2 id="pricing-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                <span className="text-orange-700 dark:text-orange-400">Simple</span> pricing for everyone
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                Every tier includes the whole scheduling engine. What changes is how many calendars
                share it and how quickly we answer.
              </p>
            </div>

            <ul role="list" className="mt-12 grid items-start gap-6 lg:grid-cols-3">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border bg-white p-6 dark:bg-gray-950 ${
                    tier.featured
                      ? 'border-orange-600 ring-1 ring-orange-600 lg:-my-2 lg:py-8'
                      : 'border-gray-200 dark:border-white/10'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="font-semibold">{tier.name}</h3>
                    {tier.featured && (
                      <span className="rounded-full bg-orange-100 px-2.5 py-1 text-xs font-medium text-orange-800 dark:bg-orange-500/15 dark:text-orange-300">
                        Most popular
                      </span>
                    )}
                  </div>

                  <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{tier.blurb}</p>

                  <p className="mt-5">
                    <span className="align-super text-lg text-gray-500 dark:text-gray-400">$</span>
                    <span className="text-4xl font-bold tracking-tight tabular-nums">{tier.price}</span>
                    <span className="text-sm text-gray-500 dark:text-gray-400">/month</span>
                  </p>

                  {/* The link name says which tier it buys. Three links all
                      called "Get started" are three identical rows in a screen
                      reader's link list, and the surrounding card is not
                      available there. */}
                  <a
                    href="#"
                    className={`mt-6 block rounded-lg px-4 py-2.5 text-center text-sm font-semibold ${
                      tier.featured
                        ? 'bg-orange-700 text-white hover:bg-orange-800'
                        : 'bg-gray-900 text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100'
                    }`}
                  >
                    Get {tier.name}
                  </a>

                  <ul role="list" className="mt-6 space-y-2.5 text-sm">
                    {tier.inherits && (
                      <li className="font-medium">Everything in {tier.inherits}, plus:</li>
                    )}
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5">
                        <span aria-hidden="true" className="mt-0.5 text-orange-700 dark:text-orange-400">
                          ✓
                        </span>
                        <span className="text-gray-600 dark:text-gray-300">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Testimonial */}
        <section aria-labelledby="love-heading" className="py-20">
          <div className="mx-auto grid max-w-6xl items-center gap-12 px-6 lg:grid-cols-2">
            <div>
              <h2 id="love-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                <span className="text-orange-700 dark:text-orange-400">People</span> love us
              </h2>
              <p className="mt-4 text-pretty text-gray-600 dark:text-gray-300">
                Sixty thousand calls a week run through this, booked by people who mostly never
                think about it. That is the review we were after.
              </p>

              <ul role="list" className="mt-6 space-y-3 text-sm">
                {TESTIMONIAL_POINTS.map((point) => (
                  <li key={point} className="flex gap-2.5">
                    <span aria-hidden="true" className="mt-0.5 text-orange-700 dark:text-orange-400">
                      ✓
                    </span>
                    <span className="text-gray-700 dark:text-gray-300">{point}</span>
                  </li>
                ))}
              </ul>

              <figure className="mt-8 border-l-2 border-orange-600 pl-5">
                <blockquote className="text-pretty text-gray-700 dark:text-gray-200">
                  We moved forty people onto it in an afternoon. The part I did not expect is that
                  our no-show rate halved, purely because the reminder goes out at a sensible hour in
                  the recipient time zone rather than ours.
                </blockquote>
                <figcaption className="mt-4 flex items-center gap-3">
                  <span aria-hidden="true" className="size-9 shrink-0 overflow-hidden rounded-full">
                    <img src={face(FACES[1], 36)} alt="" className="size-full object-cover" loading="lazy" />
                  </span>
                  <span className="text-sm">
                    <span className="block font-semibold">Daniel Okonkwo</span>
                    <span className="block text-gray-500 dark:text-gray-400">Operations lead, Gamity</span>
                  </span>
                </figcaption>
              </figure>

              <a
                href="#"
                className="mt-8 inline-block rounded-lg bg-gray-900 px-5 py-2.5 text-sm font-semibold text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
              >
                Read the customer stories
              </a>
            </div>

            {/* Texture, not content. Hidden from assistive tech and from small
                screens, where an absolutely positioned scatter is only a way to
                break the layout. */}
            <div aria-hidden="true" className="relative hidden h-96 select-none lg:block">
              {SCATTER.map((item, i) => (
                <span
                  key={i}
                  style={{ top: item.top, left: item.left }}
                  className={`absolute -translate-x-1/2 overflow-hidden rounded-full ring-4 ring-white shadow-lg dark:ring-gray-950 ${item.size}`}
                >
                  <img src={face(FACES[item.face], 96)} alt="" className="size-full object-cover" loading="lazy" />
                </span>
              ))}
              <span className="absolute top-[46%] left-[24%] flex size-12 -translate-x-1/2 items-center justify-center rounded-full bg-orange-600 text-lg text-white shadow-lg">
                ✳
              </span>
            </div>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq-heading" className="bg-gray-50 py-20 dark:bg-gray-900/40">
          <div className="mx-auto max-w-3xl px-6">
            <div className="text-center">
              <h2 id="faq-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
                Frequently <span className="text-orange-700 dark:text-orange-400">asked</span> questions
              </h2>
              <p className="mx-auto mt-3 max-w-xl text-pretty text-gray-600 dark:text-gray-300">
                The four we are asked before every trial. If yours is not here, the answer is a reply away.
              </p>
            </div>

            <div className="mt-10 space-y-3">
              {FAQS.map((item) => (
                <details
                  key={item.q}
                  className="group rounded-xl border border-gray-200 bg-white px-5 dark:border-white/10 dark:bg-gray-950"
                >
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 py-4 text-sm font-medium marker:content-none">
                    {item.q}
                    {/* gray-600, not gray-400. list-none removes the native
                        disclosure marker, which makes this glyph the only
                        visual sign the row opens — functional rather than
                        decorative, so it owes the 3:1 of WCAG 1.4.11. */}
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

        {/* Closing call to action */}
        <section aria-labelledby="cta-heading" className="relative overflow-hidden py-24">
          <div className="relative mx-auto max-w-3xl px-6 text-center">
            {/* Deliberately not a word-for-word repeat of the <h1>. A closing
                call to action restates the promise, but two identical entries
                in the heading outline are two identical rows in a screen
                reader's heading list, with nothing to tell them apart. */}
            <h2 id="cta-heading" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
              Put the back-and-forth{' '}
              <span className="text-orange-700 dark:text-orange-400">behind you</span>
            </h2>
            <p className="mx-auto mt-4 max-w-lg text-pretty text-gray-600 dark:text-gray-300">
              Free for one calendar, for as long as you want it. No card, and no call to get started —
              which would be a strange way to sell this.
            </p>
            <div className="mt-8">
              <a
                href="#"
                className="inline-block rounded-lg bg-orange-700 px-6 py-3 text-sm font-semibold text-white hover:bg-orange-800"
              >
                Get started
              </a>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-gray-200 bg-gray-50 py-14 dark:border-white/10 dark:bg-gray-900/40">
        <div className="mx-auto grid max-w-6xl gap-10 px-6 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <p className="flex items-center gap-2 font-semibold">
              <span aria-hidden="true" className="size-3 rotate-45 rounded-sm bg-orange-600" />
              shape.ai
            </p>
            <p className="mt-3 max-w-xs text-sm text-pretty text-gray-600 dark:text-gray-400">
              One link, every calendar, and a reminder that goes out on time.
            </p>
          </div>

          {[
            { heading: 'Pages', links: ['All products', 'Studio', 'Clients', 'Pricing', 'Blog'] },
            { heading: 'Socials', links: ['Facebook', 'Instagram', 'Twitter', 'LinkedIn'] },
            { heading: 'Legal', links: ['Privacy policy', 'Terms of service', 'Cookie policy'] },
          ].map((column) => (
            <nav key={column.heading} aria-labelledby={`footer-${column.heading.toLowerCase()}`}>
              <h2 id={`footer-${column.heading.toLowerCase()}`} className="text-sm font-semibold">
                {column.heading}
              </h2>
              <ul role="list" className="mt-4 space-y-2.5 text-sm">
                {column.links.map((link) => (
                  <li key={link}>
                    <a href="#" className="text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white">
                      {link}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          ))}
        </div>

        <div className="mx-auto mt-10 max-w-6xl border-t border-gray-200 px-6 pt-6 text-sm text-gray-500 dark:border-white/10 dark:text-gray-400">
          © {'2026'} shape.ai. All rights reserved.
        </div>
      </footer>
    </div>
  )
}
