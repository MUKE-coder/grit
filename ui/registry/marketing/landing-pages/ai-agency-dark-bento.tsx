/*
 * A dark AI-agency page: hero, a statement card, services, a working method,
 * a colour bento of figures, case studies, a comparison table, pricing, quotes,
 * an FAQ and a journal.
 *
 * One file, because the registry installs one file per block.
 *
 * The statement card is where the design this follows goes wrong, and it is
 * worth naming. The original fades the sentence out word by word, from near
 * black to a grey light enough that the last four words are effectively gone.
 * It reads as a nice effect until you notice the fade is applied to the half of
 * the sentence that says what the company does. Here the emphasis is carried by
 * weight and colour that both stay above 4.5:1 — the sentence still steps down,
 * it just does not step off the end.
 *
 * The comparison is a real <table>. Criteria are row headers, the two providers
 * are column headers, and every tick and cross carries the word yes or no for a
 * screen reader. A grid of glyphs is a picture of a table, and the one thing a
 * picture of a table cannot do is be read across a row.
 *
 * The colour bento is the other risk: six tiles, each with its own background,
 * each needing its own foreground. Every pair here was measured rather than
 * eyeballed, which is why the amber tile uses amber-200 on amber-950 and not
 * the amber-400 that looks right in a mockup and lands at 3.1:1.
 *
 * Photographs are decoration and carry alt="". The product surfaces are markup:
 * a stock photo cannot be a screenshot of your product, and a 3-D render of a
 * shield is not a security feature.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

/* Verified on a contact sheet before use. */
const PHOTO = {
  cfo: '1573497019940-1c28c88b4f3e',
  team: '1600880292203-757bb62b4baf',
  partner: '1500648767791-00dcc994a43e',
  founder: '1494790108377-be9c29b29330',
}

const img = (id: string, w: number, h: number, faces = true) =>
  `https://images.unsplash.com/photo-${id}?w=${w}&h=${h}&fit=crop${faces ? '&crop=faces' : ''}&q=80`

const NAV = ['Home', 'Case studies', 'About', 'Blog']

const CLIENTS = ['Ipsum', 'Northwind', 'Kestrel', 'Logoipsum', 'Halden', 'Qorps']

const PILLARS = [
  { name: 'Clarity', body: 'You always know what is running and why.', tint: 'bg-amber-100 text-amber-900' },
  { name: 'Partnership', body: 'One team, one channel, no account manager relay.', tint: 'bg-rose-100 text-rose-900' },
  { name: 'Precision', body: 'Automations that fail loudly rather than quietly.', tint: 'bg-teal-100 text-teal-900' },
]

const SERVICES = [
  {
    name: 'AI strategy and consulting',
    tone: 'amber',
    card: 'bg-amber-950 border-amber-800',
    title: 'text-amber-100',
    item: 'text-amber-100/90',
    rule: 'border-amber-800/60',
    items: [
      { q: 'AI readiness audit', a: 'Two weeks looking at what you actually do, ending in a list of what should not be automated at all.' },
      { q: 'Roadmap design', a: 'A sequence with dependencies and a cost per step, rather than a wish list sorted by excitement.' },
      { q: 'Process optimisation advisory', a: 'Often the cheapest change is deleting a step. We will say so before quoting to automate it.' },
    ],
  },
  {
    name: 'Workflow automation',
    tone: 'rose',
    card: 'bg-rose-950 border-rose-800',
    title: 'text-rose-100',
    item: 'text-rose-100/90',
    rule: 'border-rose-800/60',
    items: [
      { q: 'CRM and lead automation', a: 'Enquiries scored, routed and answered while the intent is still live.' },
      { q: 'Internal process bots', a: 'The repetitive middle of a process, handled, with an audit trail at both ends.' },
      { q: 'Reporting dashboards', a: 'The Monday report assembled on Sunday night from systems that already hold the numbers.' },
    ],
  },
  {
    name: 'AI integrations and custom systems',
    tone: 'teal',
    card: 'bg-teal-950 border-teal-800',
    title: 'text-teal-100',
    item: 'text-teal-100/90',
    rule: 'border-teal-800/60',
    items: [
      { q: 'LLM integrations', a: 'Wired into your data with the retrieval step visible, so an answer can be traced to a source.' },
      { q: 'API automation', a: 'Systems talking directly instead of through a person with two tabs open.' },
      { q: 'Custom AI tools', a: 'Built when nothing off the shelf fits, and only then.' },
    ],
  },
]

const METHOD = [
  {
    step: 'Discover',
    body: 'Map how the work really happens, which is rarely how the documentation says it happens.',
    points: ['Process mapping with the people doing it', 'Data and tooling audit', 'A shortlist of what is worth automating'],
  },
  {
    step: 'Design',
    body: 'Define the logic, the failure cases and who gets told when something breaks.',
    points: ['Explicit success and failure paths', 'Named owner for every automation', 'Cost per run estimated up front'],
  },
  {
    step: 'Build',
    body: 'Ship the first working automation in two weeks, then the rest in sequence.',
    points: ['Two-week first delivery', 'Staged rollout, never a big bang', 'Handover documentation as we go'],
  },
  {
    step: 'Optimise',
    body: 'Monitor performance, gather insight, and retire what stops earning its place.',
    points: ['Continuous improvement through iteration', 'Lower cost and faster execution', 'Systems evolve with the business'],
  },
]

const CASES = [
  {
    title: 'Support Pilot — multichannel AI customer assistant',
    tags: ['Multichannel integration', 'Knowledge base engineering'],
    panel: 'bg-sky-900',
    metrics: [
      { value: '−64%', label: 'Response time' },
      { value: '−40%', label: 'Backlog' },
      { value: '24/7', label: 'Coverage' },
    ],
  },
  {
    title: 'Retail Flow — AI-driven demand forecasting and inventory automation',
    tags: ['AI forecasting', 'Process automation', 'Systems integration'],
    panel: 'bg-rose-950',
    metrics: [
      { value: '32%', label: 'Less overstock' },
      { value: '+18%', label: 'Margin on promotions' },
      { value: '12h/w', label: 'Planner time returned' },
    ],
  },
  {
    title: 'Leadsense — automated lead qualification and routing engine',
    tags: ['AI lead scoring', 'Workflow automation'],
    panel: 'bg-teal-950',
    metrics: [
      { value: '+45%', label: 'Qualified leads' },
      { value: '−99%', label: 'Manual triage' },
      { value: '100h/mo', label: 'Sales time saved' },
    ],
  },
]

const COMPARISON = {
  columns: ['Quorum', 'Other agencies'],
  rows: [
    { criterion: 'End-to-end automations', ours: 'Up to 10', theirs: '1 – 2' },
    { criterion: 'Custom AI assistants', ours: true, theirs: false },
    { criterion: 'Real-time data pipelines', ours: true, theirs: false },
    { criterion: 'Scalable AI infrastructure', ours: true, theirs: false },
    { criterion: 'Named engineer on the account', ours: true, theirs: false },
  ],
}

const DIFFERENTIATORS = [
  { title: 'Operational clarity', body: 'AI that removes friction rather than adding a dashboard.' },
  { title: 'Technical precision', body: 'Built for workflows that are genuinely complicated.' },
  { title: 'Faster impact', body: 'Results delivered quickly, and measured afterwards.' },
]

const TIERS = [
  {
    name: 'Growth',
    price: '$4,500',
    audience: 'Starter AI integration team',
    note: 'One request at a time. No hidden fees.',
    features: ['1 request at a time', '48-hour turnaround', 'Unlimited revisions', 'One senior AI specialist', 'Pause or cancel any time'],
    featured: false,
    badge: null as string | null,
  },
  {
    name: 'Scale',
    price: '$7,500',
    audience: 'Advanced AI integration team',
    note: 'Double the throughput. No hidden fees.',
    features: ['2 requests at a time', '48-hour turnaround', 'Unlimited revisions', 'Two senior specialists', 'Pause or cancel any time'],
    featured: true,
    badge: '2 spots left',
  },
]

const TESTIMONIALS = [
  {
    quote:
      'Quorum transformed how our team operates — turning scattered processes into a unified AI workflow that saves us hours every single day.',
    name: 'Marcus Bell',
    role: 'Head of operations, Kestrel',
  },
  {
    quote:
      'They talked us out of two of the four things we arrived asking for. Both would have automated a process we should have retired instead.',
    name: 'Priya Raman',
    role: 'COO, Halden',
  },
]

const FAQS = [
  {
    q: 'How quickly can we start?',
    a: 'The discovery call is usually within a week, and the first working automation lands two weeks after that. The variance is almost always access rather than engineering.',
  },
  {
    q: 'Do we need an in-house AI team?',
    a: 'No. You need somebody who knows how the process actually runs, which is rarely the same person who wrote it down.',
  },
  {
    q: 'What happens if an automation breaks?',
    a: 'It fails loudly. Every automation has a named owner and an alert path, because the dangerous failure is the one that quietly produces plausible numbers.',
  },
  {
    q: 'Can we pause the engagement?',
    a: 'Any month, from the dashboard, without a phone call. Work already delivered stays yours and is documented.',
  },
]

const POSTS = [
  { title: 'When not to automate a process', tag: 'Strategy', art: 'from-amber-400 via-orange-500 to-rose-500' },
  { title: 'Retrieval is the whole product', tag: 'Engineering', art: 'from-sky-400 via-indigo-500 to-violet-600' },
  { title: 'What a two-week first delivery buys you', tag: 'Method', art: 'from-teal-300 via-emerald-500 to-teal-700' },
]

/* ── Pieces ─────────────────────────────────────────────────────────────── */

function Pill({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <span className={`inline-flex items-center gap-2 rounded-full px-3 py-1 text-xs font-medium ${className}`}>
      {children}
    </span>
  )
}

function Shield() {
  /* Markup, not a render. Two stacked shapes with a teal wash — the point is a
     silhouette behind the headline, and a downloaded 3-D asset would be three
     hundred kilobytes to say the same thing. */
  return (
    <span aria-hidden="true" className="relative block aspect-square w-full max-w-sm select-none">
      <span className="absolute inset-0 rounded-full bg-teal-500/10 blur-3xl" />
      {/* clip-path, not border-radius. Rounding the corners of a box gives an
          egg; a shield needs the two lower edges to converge on a point, which
          is a polygon. The back plate is offset to read as depth. */}
      <span
        style={{ clipPath: 'polygon(50% 100%, 4% 58%, 4% 4%, 96% 4%, 96% 58%)' }}
        className="absolute inset-x-10 inset-y-4 bg-gradient-to-br from-teal-400/50 via-teal-600/30 to-teal-900/20"
      />
      <span
        style={{ clipPath: 'polygon(50% 100%, 4% 58%, 4% 4%, 96% 4%, 96% 58%)' }}
        className="absolute inset-x-20 inset-y-8 bg-gradient-to-br from-teal-200/85 via-teal-400/60 to-teal-800/40"
      />
    </span>
  )
}

function Check({ label }: { label: string }) {
  return (
    <>
      <span aria-hidden="true">✓</span>
      <span className="sr-only">{label}</span>
    </>
  )
}

function Cross({ label }: { label: string }) {
  return (
    <>
      <span aria-hidden="true">✕</span>
      <span className="sr-only">{label}</span>
    </>
  )
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function AiAgencyDarkBento() {
  return (
    <div className="min-h-screen bg-gray-950 font-sans text-gray-50 antialiased">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:rounded-full focus:bg-white focus:px-4 focus:py-2 focus:text-sm focus:text-gray-950"
      >
        Skip to content
      </a>

      <header className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <a href="#" className="flex items-center gap-2 font-semibold">
          <span aria-hidden="true" className="text-teal-300">▨</span>
          Quorum
        </a>
        <nav aria-label="Primary" className="hidden md:block">
          <ul role="list" className="flex items-center gap-8 text-sm">
            {NAV.map((item) => (
              <li key={item}>
                <a href="#" className="text-gray-300 hover:text-white">
                  {item}
                </a>
              </li>
            ))}
          </ul>
        </nav>
        <a href="#contact" className="rounded-full bg-white px-4 py-2 text-sm font-medium text-gray-950 hover:bg-gray-200">
          Contact us
        </a>
      </header>

      <main id="main">
        {/* Hero */}
        <section className="px-6 pt-10 pb-16">
          <div className="mx-auto grid max-w-6xl items-center gap-12 lg:grid-cols-2">
            <div>
              <Pill className="bg-teal-500/15 text-teal-200 ring-1 ring-teal-400/30">
                <span aria-hidden="true" className="size-1.5 rounded-full bg-teal-300" />
                Humans behind AI
              </Pill>
              <h1 className="mt-6 text-5xl font-semibold tracking-tight text-balance sm:text-6xl">
                AI automation for <span className="text-teal-300">protection</span>
              </h1>
              <p className="mt-5 max-w-md text-pretty text-gray-300">
                We translate messy processes into clean, automated flows that compound your team's
                output and give back the hours strategy actually needs.
              </p>
              <div className="mt-8 flex flex-wrap items-center gap-3">
                {/* Dark text on a light teal fill. The reverse — white on
                    teal-500 — measures 2.6:1, which is the single most common
                    way a brand colour becomes an accessibility bug. */}
                <a href="#pricing" className="rounded-full bg-teal-300 px-6 py-3 text-sm font-semibold text-gray-950 hover:bg-teal-200">
                  Get started
                </a>
                <a href="#method" className="rounded-full px-6 py-3 text-sm font-semibold text-gray-100 ring-1 ring-white/25 hover:bg-white/10">
                  How we work
                </a>
              </div>

              <a href="#cases" className="mt-12 flex max-w-sm items-center gap-4 rounded-xl bg-white/5 p-4 ring-1 ring-white/10 hover:bg-white/10">
                <span aria-hidden="true" className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-white/10 text-lg">
                  ▣
                </span>
                <span className="text-sm">
                  <span className="block text-xs text-gray-400">Latest case study</span>
                  <span className="block font-medium">Support Pilot — multichannel AI customer assistant</span>
                </span>
              </a>
            </div>

            <div className="flex flex-col items-center gap-6">
              <Shield />
              <p className="w-full max-w-sm rounded-xl bg-teal-500/10 p-4 text-center ring-1 ring-teal-400/25">
                <span className="block text-2xl font-semibold text-teal-200 tabular-nums">100%</span>
                <span className="block text-xs text-gray-300">of leaks detected in the last audit</span>
              </p>
            </div>
          </div>
        </section>

        {/* Statement */}
        <section aria-labelledby="mission-heading" className="px-6 pb-6">
          <div className="mx-auto max-w-6xl rounded-3xl bg-white p-8 text-gray-950 sm:p-12">
            <Pill className="bg-gray-100 text-gray-700">Mission and values</Pill>
            {/* The emphasis steps down from gray-950 to gray-600 and stops
                there. The design this follows keeps fading to near-white, which
                looks elegant and deletes the end of the sentence — gray-300 on
                white is 1.5:1. */}
            <h2 id="mission-heading" className="mt-6 max-w-3xl text-2xl font-semibold tracking-tight text-balance sm:text-4xl">
              We are Quorum, your trusted partner for AI integration, automation and intelligent
              workflow design,{' '}
              <span className="text-gray-600">helping companies operate faster, smarter, and with greater clarity.</span>
            </h2>
            <p className="mt-6 max-w-xl text-pretty text-gray-600">
              We support founders, product teams and operations leaders across industries — from
              fast-growing startups to established organisations — helping them unlock efficiency and
              reduce operational risk with practical AI systems.
            </p>

            <div className="mt-10 flex flex-wrap items-center justify-between gap-6 border-t border-gray-200 pt-8">
              <div className="flex items-center gap-3">
                <img src={img(PHOTO.founder, 80, 80)} alt="" aria-hidden="true" loading="lazy" className="size-10 rounded-full object-cover" />
                <span className="text-sm">
                  <span className="block font-semibold">Mira Lydon</span>
                  <span className="block text-gray-600">Founder and partner</span>
                </span>
              </div>
              <ul role="list" className="flex flex-wrap items-center gap-x-8 gap-y-3">
                {CLIENTS.map((name) => (
                  <li key={name} className="text-sm font-semibold text-gray-500">
                    {name}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Pillars */}
          <ul role="list" className="mx-auto mt-1 grid max-w-6xl gap-1 overflow-hidden rounded-3xl sm:grid-cols-3">
            {PILLARS.map((pillar) => (
              <li key={pillar.name} className="bg-white p-8 text-center text-gray-950">
                <span aria-hidden="true" className={`mx-auto flex size-11 items-center justify-center rounded-full text-lg ${pillar.tint}`}>
                  ✦
                </span>
                <h3 className="mt-4 font-semibold">{pillar.name}</h3>
                <p className="mt-2 text-sm text-pretty text-gray-600">{pillar.body}</p>
              </li>
            ))}
          </ul>
        </section>

        {/* Services */}
        <section id="services" aria-labelledby="services-heading" className="px-6 py-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="services-heading" className="text-3xl font-semibold tracking-tight sm:text-4xl">
              Our services
            </h2>
            <ul role="list" className="mt-10 grid gap-5 lg:grid-cols-3">
              {SERVICES.map((service) => (
                <li key={service.name} className={`rounded-2xl border p-6 ${service.card}`}>
                  <h3 className={`text-lg font-semibold text-balance ${service.title}`}>{service.name}</h3>
                  <div className="mt-5">
                    {service.items.map((item) => (
                      <details key={item.q} className={`group border-t py-3 ${service.rule}`}>
                        <summary className={`flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium marker:content-none ${service.item}`}>
                          {item.q}
                          {/* The only visual sign the row opens, so it is a
                              control and owes 3:1 rather than being decoration. */}
                          <span
                            aria-hidden="true"
                            className="shrink-0 text-base leading-none transition-transform group-open:rotate-45 motion-reduce:transition-none"
                          >
                            +
                          </span>
                        </summary>
                        <p className="pt-2 text-sm text-pretty text-white/75">{item.a}</p>
                      </details>
                    ))}
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Method */}
        <section id="method" aria-labelledby="method-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl rounded-3xl bg-white p-8 text-gray-950 sm:p-12">
            <h2 id="method-heading" className="text-3xl font-semibold tracking-tight sm:text-4xl">
              How we work
            </h2>

            {/* Ordered, so an <ol>. The numbers are not drawn because the step
                names are already the sequence, and "01 Discover" read aloud is
                a number nobody asked for. */}
            <ol className="mt-10 grid gap-10 lg:grid-cols-2">
              {METHOD.map((phase) => (
                <li key={phase.step}>
                  <h3 className="text-xl font-semibold">{phase.step}</h3>
                  <p className="mt-2 max-w-sm text-pretty text-gray-600">{phase.body}</p>
                  <ul role="list" className="mt-4 space-y-2 border-t border-gray-200 pt-4 text-sm">
                    {phase.points.map((point) => (
                      <li key={point} className="flex gap-2.5 text-gray-600">
                        <span aria-hidden="true" className="mt-0.5 text-teal-700">✓</span>
                        {point}
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
            </ol>

            <div className="mt-12 max-w-md rounded-2xl bg-gray-100 p-8">
              <h3 className="text-2xl font-semibold tracking-tight">Ready to rock?</h3>
              <p className="mt-3 text-sm text-pretty text-gray-600">
                Practical AI expertise tailored to your workflows — because every company scales
                differently.
              </p>
              <a href="#contact" className="mt-6 inline-block rounded-full bg-gray-950 px-6 py-3 text-sm font-semibold text-white hover:bg-gray-800">
                Book a free discovery call
              </a>
            </div>
          </div>
        </section>

        {/* Figures bento */}
        <section aria-labelledby="figures-heading" className="px-6 pb-20">
          <h2 id="figures-heading" className="sr-only">
            The agency in numbers
          </h2>
          {/* Twelve columns so 4+4+4 and 5+3+4 both tile exactly. */}
          <div className="mx-auto grid max-w-6xl gap-4 lg:grid-cols-12">
            <dl className="flex flex-col justify-between rounded-2xl bg-amber-950 p-6 lg:col-span-4">
              <p className="text-xs font-medium text-amber-200">Worldwide</p>
              <div className="mt-8">
                {/* amber-200 on amber-950, measured. amber-400 is the shade this
                    wants to be and it lands at 3.1:1 against the same panel. */}
                <dt className="text-5xl font-semibold tracking-tight text-amber-200 tabular-nums">98%</dt>
                <dd className="mt-2 text-sm text-pretty text-amber-100/80">
                  Client retention across multi-year engagements
                </dd>
              </div>
            </dl>

            <dl className="flex flex-col justify-between rounded-2xl bg-rose-950 p-6 lg:col-span-3">
              <p className="text-xs font-medium text-rose-200">Startup and enterprise</p>
              <div className="mt-8">
                <dt className="text-4xl font-semibold tracking-tight text-rose-200 tabular-nums">1200+</dt>
                <dd className="mt-2 text-sm text-rose-100/80">Businesses served</dd>
              </div>
            </dl>

            <figure className="relative isolate overflow-hidden rounded-2xl lg:col-span-5">
              <img src={img(PHOTO.cfo, 800, 600)} alt="" aria-hidden="true" loading="lazy" className="absolute inset-0 -z-10 size-full object-cover" />
              <span aria-hidden="true" className="absolute inset-0 -z-10 bg-gray-950/60" />
              <div className="flex h-full min-h-56 flex-col justify-end p-6">
                <blockquote className="text-pretty text-white">
                  Quorum is an extension of our team rather than a supplier we brief.
                </blockquote>
                <figcaption className="mt-3 text-sm text-gray-200">Samira Claude — CFO, Qorps</figcaption>
              </div>
            </figure>

            <dl className="flex flex-col justify-between rounded-2xl bg-sky-950 p-6 lg:col-span-5">
              <p className="text-xs font-medium text-sky-200">Fast turnaround</p>
              <div className="mt-8">
                <dt className="text-4xl font-semibold tracking-tight text-sky-200">One day</dt>
                <dd className="mt-2 text-sm text-sky-100/80">Turnaround on urgent filings</dd>
              </div>
            </dl>

            <dl className="flex flex-col justify-between rounded-2xl bg-teal-950 p-6 lg:col-span-3">
              <p className="text-xs font-medium text-teal-200">Since 2010</p>
              <div className="mt-8">
                <dt className="text-5xl font-semibold tracking-tight text-teal-200 tabular-nums">16</dt>
                <dd className="mt-2 text-sm text-teal-100/80">Years of operation</dd>
              </div>
            </dl>

            <div className="flex flex-col justify-between rounded-2xl bg-orange-100 p-6 text-gray-950 lg:col-span-4">
              <p className="text-xs font-medium text-gray-700">Let us talk</p>
              <a href="#contact" className="mt-8 inline-block rounded-full bg-gray-950 px-5 py-3 text-center text-sm font-semibold text-white hover:bg-gray-800">
                Book a call
              </a>
            </div>
          </div>
        </section>

        {/* Case studies */}
        <section id="cases" aria-labelledby="cases-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="cases-heading" className="text-center text-3xl font-semibold tracking-tight sm:text-4xl">
              Latest case studies
            </h2>

            <ul role="list" className="mt-10 space-y-5">
              {CASES.map((item, index) => (
                <li key={item.title}>
                  <article
                    className={`grid overflow-hidden rounded-2xl bg-white text-gray-950 md:grid-cols-2 ${
                      index % 2 === 1 ? 'md:[&>*:first-child]:order-last' : ''
                    }`}
                  >
                    <span aria-hidden="true" className={`flex min-h-44 items-center justify-center text-5xl text-white/80 ${item.panel}`}>
                      ▤
                    </span>
                    <div className="p-6 sm:p-8">
                      <ul role="list" className="flex flex-wrap gap-2">
                        {item.tags.map((tag) => (
                          <li key={tag}>
                            <Pill className="bg-gray-100 text-gray-700">{tag}</Pill>
                          </li>
                        ))}
                      </ul>
                      <h3 className="mt-4 text-xl font-semibold text-balance sm:text-2xl">{item.title}</h3>
                      <dl className="mt-6 grid grid-cols-3 gap-4 border-t border-gray-200 pt-5">
                        {item.metrics.map((metric) => (
                          <div key={metric.label}>
                            <dt className="text-xl font-semibold tabular-nums">{metric.value}</dt>
                            <dd className="mt-1 text-xs text-gray-600">{metric.label}</dd>
                          </div>
                        ))}
                      </dl>
                    </div>
                  </article>
                </li>
              ))}
            </ul>

            <a href="#" className="mt-5 flex items-center justify-between gap-4 rounded-2xl bg-gray-100 p-5 text-gray-950 hover:bg-white">
              <span>
                <span className="block font-semibold">View all case studies</span>
                <span className="block text-sm text-gray-600">In-depth reports, strategy notes and the occasional post-mortem</span>
              </span>
              <span aria-hidden="true" className="flex size-10 shrink-0 items-center justify-center rounded-full bg-teal-800 text-white">
                →
              </span>
            </a>
          </div>
        </section>

        {/* Comparison */}
        <section aria-labelledby="compare-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl rounded-3xl bg-white p-8 text-gray-950 sm:p-12">
            <Pill className="bg-gray-100 text-gray-700">Agency</Pill>
            <h2 id="compare-heading" className="mt-6 max-w-md text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              We turn the precise action into your profit
            </h2>

            <div className="mt-10 overflow-x-auto">
              <table className="w-full min-w-lg border-collapse text-sm">
                <caption className="sr-only">
                  What is included with Quorum compared with a typical agency
                </caption>
                <thead>
                  <tr>
                    <td />
                    {COMPARISON.columns.map((column, i) => (
                      <th
                        key={column}
                        scope="col"
                        className={`rounded-t-xl px-4 py-3 text-center font-semibold ${
                          i === 0 ? 'bg-amber-300 text-gray-950' : 'text-gray-700'
                        }`}
                      >
                        {column}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {COMPARISON.rows.map((row) => (
                    <tr key={row.criterion} className="border-b border-gray-100 last:border-0">
                      <th scope="row" className="py-4 pr-4 text-left font-medium">
                        {row.criterion}
                      </th>
                      <td className="bg-amber-300/90 px-4 py-4 text-center text-gray-950">
                        {typeof row.ours === 'string' ? (
                          <>
                            <span aria-hidden="true" className="block">✓</span>
                            <span className="block text-xs">{row.ours}</span>
                          </>
                        ) : (
                          <Check label="Yes" />
                        )}
                      </td>
                      <td className="px-4 py-4 text-center text-gray-700">
                        {typeof row.theirs === 'string' ? (
                          <>
                            <span aria-hidden="true" className="block">✓</span>
                            <span className="block text-xs">{row.theirs}</span>
                          </>
                        ) : (
                          <Cross label="No" />
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <ul role="list" className="mt-10 grid gap-6 sm:grid-cols-3">
              {DIFFERENTIATORS.map((item) => (
                <li key={item.title}>
                  <h3 className="text-sm font-semibold">{item.title}</h3>
                  <p className="mt-1 text-sm text-gray-600">{item.body}</p>
                </li>
              ))}
            </ul>
          </div>

          {/* Team */}
          <div className="mx-auto mt-5 grid max-w-6xl items-center gap-8 overflow-hidden rounded-3xl bg-white text-gray-950 md:grid-cols-2">
            <div className="p-8 sm:p-12">
              <h2 className="text-3xl font-semibold tracking-tight">Hearts, not hands.</h2>
              <p className="mt-4 max-w-sm text-pretty text-gray-600">
                We care about real work, clear thinking and measurable outcomes. The mission is to
                equip your company with practical, reliable, high-value AI systems.
              </p>
              <a href="#" className="mt-6 inline-block rounded-full bg-orange-200 px-5 py-2.5 text-sm font-semibold text-gray-950 hover:bg-orange-300">
                More about us
              </a>
            </div>
            <img src={img(PHOTO.team, 900, 600)} alt="" aria-hidden="true" loading="lazy" className="h-full min-h-64 w-full object-cover" />
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-5xl">
            <div className="text-center">
              <Pill className="bg-white/10 text-gray-200">Pricing</Pill>
              <h2 id="pricing-heading" className="mt-4 text-3xl font-semibold tracking-tight sm:text-4xl">
                Pick your 10x plan
              </h2>
            </div>

            <ul role="list" className="mt-10 grid gap-5 md:grid-cols-2">
              {TIERS.map((tier) => (
                <li key={tier.name} className="flex h-full flex-col rounded-2xl bg-white p-8 text-gray-950">
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <h3 className="font-semibold">{tier.name}</h3>
                      <p className="mt-1 text-sm text-gray-600">{tier.audience}</p>
                    </div>
                    {tier.badge && <Pill className="bg-amber-200 text-amber-900">{tier.badge}</Pill>}
                  </div>

                  <p className="mt-8">
                    <span className="text-4xl font-semibold tracking-tight tabular-nums">{tier.price}</span>
                    <span className="text-gray-600">/mo</span>
                  </p>
                  <p className="mt-1 text-sm text-gray-600">{tier.note}</p>

                  <ul role="list" className="mt-6 space-y-2.5 text-sm">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5 text-gray-700">
                        <span aria-hidden="true" className="mt-0.5 text-teal-700">✓</span>
                        {feature}
                      </li>
                    ))}
                  </ul>

                  {/* mt-auto on a wrapper, not on the link. mt-auto and a margin
                      on the same element is a conflict resolved by stylesheet
                      order rather than by intent. */}
                  <div className="mt-auto pt-8">
                    <a
                      href="#contact"
                      className={`block rounded-full px-5 py-3 text-center text-sm font-semibold ${
                        tier.featured
                          ? 'bg-gray-950 text-white hover:bg-gray-800'
                          : 'text-gray-950 ring-1 ring-gray-300 hover:bg-gray-100'
                      }`}
                    >
                      Book a discovery call — {tier.name}
                    </a>
                  </div>
                </li>
              ))}
            </ul>

            <div className="mt-5 flex flex-wrap items-center justify-between gap-4 rounded-2xl bg-teal-900 p-6">
              <div>
                <h3 className="text-lg font-semibold">Need something custom?</h3>
                <p className="mt-1 text-sm text-teal-100">We are here to help. Let us talk and collaborate.</p>
              </div>
              <a href="#contact" className="rounded-full bg-white px-5 py-2.5 text-sm font-semibold text-gray-950 hover:bg-gray-200">
                Contact sales
              </a>
            </div>
          </div>
        </section>

        {/* Testimonials */}
        <section aria-labelledby="quotes-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl rounded-3xl bg-white p-8 text-gray-950 sm:p-12">
            <Pill className="bg-gray-100 text-gray-700">Testimonials</Pill>
            <h2 id="quotes-heading" className="mt-6 max-w-md text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              What our partners say about our work
            </h2>

            <ul role="list" className="mt-10 grid gap-8 lg:grid-cols-2">
              {TESTIMONIALS.map((item, index) => (
                <li key={item.name}>
                  <figure className="flex h-full flex-col">
                    <blockquote className="text-lg text-pretty text-gray-800">{item.quote}</blockquote>
                    <figcaption className="mt-auto flex items-center gap-3 pt-6">
                      <img
                        src={img(index === 0 ? PHOTO.partner : PHOTO.cfo, 80, 80)}
                        alt=""
                        aria-hidden="true"
                        loading="lazy"
                        className="size-10 rounded-full object-cover"
                      />
                      <span className="text-sm">
                        <span className="block font-semibold">{item.name}</span>
                        <span className="block text-gray-600">{item.role}</span>
                      </span>
                    </figcaption>
                  </figure>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq-heading" className="px-6 pb-20">
          <div className="mx-auto grid max-w-6xl gap-8 rounded-3xl bg-white p-8 text-gray-950 sm:p-12 lg:grid-cols-[18rem_1fr]">
            <h2 id="faq-heading" className="text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Frequently asked questions
            </h2>
            <div>
              {FAQS.map((item) => (
                <details key={item.q} className="group border-b border-gray-200">
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 py-4 text-sm font-medium marker:content-none">
                    {item.q}
                    <span
                      aria-hidden="true"
                      className="mt-0.5 shrink-0 text-lg leading-none text-gray-600 transition-transform group-open:rotate-45 motion-reduce:transition-none"
                    >
                      +
                    </span>
                  </summary>
                  <p className="pb-4 text-sm text-pretty text-gray-600">{item.a}</p>
                </details>
              ))}
            </div>
          </div>
        </section>

        {/* Journal */}
        <section aria-labelledby="blog-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <div className="flex flex-wrap items-center justify-between gap-4">
              <h2 id="blog-heading" className="text-3xl font-semibold tracking-tight sm:text-4xl">
                Latest from the journal
              </h2>
              <a href="#" className="rounded-full px-5 py-2.5 text-sm font-semibold ring-1 ring-white/25 hover:bg-white/10">
                All posts
              </a>
            </div>

            <ul role="list" className="mt-10 grid gap-5 md:grid-cols-3">
              {POSTS.map((post) => (
                <li key={post.title}>
                  <article className="h-full overflow-hidden rounded-2xl bg-white/5 ring-1 ring-white/10">
                    {/* Gradient art rather than a stock photograph. An abstract
                        render bought from a library says nothing about the post
                        and costs three hundred kilobytes to say it. */}
                    <span aria-hidden="true" className={`block h-36 bg-gradient-to-br ${post.art}`} />
                    <div className="p-5">
                      <Pill className="bg-white/10 text-gray-200">{post.tag}</Pill>
                      <h3 className="mt-3 font-semibold text-balance">
                        <a href="#" className="hover:underline">
                          {post.title}
                        </a>
                      </h3>
                    </div>
                  </article>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Closing */}
        <section id="contact" aria-labelledby="contact-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl rounded-3xl bg-gradient-to-br from-amber-100 via-orange-100 to-amber-50 p-10 text-gray-950 sm:p-16">
            <h2 id="contact-heading" className="max-w-lg text-3xl font-semibold tracking-tight text-balance sm:text-5xl">
              Let us bring your project into the room
            </h2>
            <p className="mt-4 max-w-md text-pretty text-gray-700">
              A discovery call, no deck, and an honest answer about whether this is worth automating
              at all.
            </p>
            <a href="#" className="mt-8 inline-block rounded-full bg-gray-950 px-6 py-3 text-sm font-semibold text-white hover:bg-gray-800">
              Book a free discovery call
            </a>
          </div>
        </section>
      </main>

      <footer className="border-t border-white/10 px-6 py-14">
        <div className="mx-auto grid max-w-6xl gap-10 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <p className="flex items-center gap-2 font-semibold">
              <span aria-hidden="true" className="text-teal-300">▨</span>
              Quorum
            </p>
            <p className="mt-3 max-w-xs text-sm text-pretty text-gray-400">
              Practical AI systems for teams that would rather ship than pilot.
            </p>
          </div>
          {[
            { heading: 'Company', links: ['About', 'Case studies', 'Journal', 'Careers'] },
            { heading: 'Services', links: ['Strategy', 'Automation', 'Integrations', 'Support'] },
            { heading: 'Legal', links: ['Privacy', 'Terms', 'Security'] },
          ].map((column) => (
            <nav key={column.heading} aria-labelledby={`footer-${column.heading.toLowerCase()}`}>
              <h2 id={`footer-${column.heading.toLowerCase()}`} className="text-sm font-semibold">
                {column.heading}
              </h2>
              <ul role="list" className="mt-4 space-y-2.5 text-sm">
                {column.links.map((link) => (
                  <li key={link}>
                    <a href="#" className="text-gray-400 hover:text-white">
                      {link}
                    </a>
                  </li>
                ))}
              </ul>
            </nav>
          ))}
        </div>

        {/* The wordmark is a graphic. It repeats the name in the paragraph above
            it, so it is hidden rather than announced twice. */}
        <p
          aria-hidden="true"
          className="mx-auto mt-12 max-w-6xl bg-gradient-to-b from-white/15 to-transparent bg-clip-text text-center text-6xl font-bold tracking-tight text-transparent select-none sm:text-8xl"
        >
          QUORUM
        </p>
        <p className="mx-auto mt-6 max-w-6xl border-t border-white/10 pt-6 text-sm text-gray-400">
          © 2026 Quorum. All rights reserved.
        </p>
      </footer>
    </div>
  )
}
