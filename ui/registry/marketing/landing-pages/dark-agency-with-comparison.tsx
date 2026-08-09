/*
 * A dark studio page: hero, work, comparison table, pricing, founder's note
 * and an FAQ.
 *
 * One file, because the registry installs one file per block.
 *
 * The section that makes this page worth having is the comparison. "Us versus
 * them" is tabular data — criteria down the side, two providers across the top
 * — and almost every site builds it as two columns of divs, which loses the
 * relationship between a row's label and its two answers. Here it is a real
 * <table>: each criterion is a <th scope="row">, each provider a
 * <th scope="col">, so a screen reader announces "Design and engineering in
 * sync, us, yes" rather than reading two disconnected lists of ticks.
 *
 * The ticks and crosses carry text. A green check is a picture; on its own a
 * cell reads as empty. Each one has an sr-only "yes" or "no", which is what
 * makes the table answerable rather than merely visible.
 *
 * This block commits to dark rather than being theme-aware, and says so in one
 * place: the wrapper sets the colours and everything inside inherits. Blocks
 * that support both themes are the norm in this registry; a page whose whole
 * design is a glow on near-black is not improved by a light mode bolted on.
 *
 * The watermark words behind two sections are aria-hidden. They are texture,
 * and read aloud "Projects" twice in a row is a stutter with no meaning.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Work', 'Services', 'Pricing', 'Blog']

const CLIENTS = [
  'Cursor', 'Rogue', 'Loopback', 'Verbatim', 'Appsmith', 'Contentbox',
  'Kearney', 'Hexbind', 'Halcyon', 'Junipero', 'Northbeam', 'Thrust',
]

const CAPABILITIES = [
  {
    title: 'Design and development',
    body: 'One team from wireframe to deploy, so nothing is lost explaining the design to the people building it.',
    mock: 'browser',
  },
  {
    title: 'Rapid iteration',
    body: 'A preview URL in week one and a written update every Friday. You never have to ask where it is.',
    mock: 'progress',
  },
  {
    title: 'Hosting and maintenance',
    body: 'Deployment, dependency updates and a monthly check, so the site does not quietly rot after launch.',
    mock: 'stack',
  },
]

const PROJECTS = [
  { name: 'Kitchen commerce', tone: 'from-orange-500/30 to-amber-500/10', span: 'lg:col-span-2 lg:row-span-2' },
  { name: 'Transit app', tone: 'from-yellow-400/30 to-yellow-500/5', span: '' },
  { name: 'Analytics console', tone: 'from-emerald-500/25 to-teal-500/5', span: '' },
  /* col-span-3, not 2. Across three columns the first tile takes two of them
     for two rows, so rows one and two are full. A two-column tile on row three
     leaves the third cell empty, and CSS grid renders that as a visible hole
     rather than stretching anything to cover it. The spans have to tile
     exactly: 4 + 1 + 1 + 3 = 9 cells over three rows. */
  { name: 'Field reports', tone: 'from-sky-500/25 to-indigo-500/5', span: 'lg:col-span-3' },
]

const TESTIMONIALS = [
  { quote: 'Excellent communication and professionalism, open to ideas, and humble when it counted. We will engage again.', name: 'Marcus Chen', role: 'CTO, Cursor' },
  { quote: 'Quick to respond, very professional, and shipped a site within a week. Looking forward to the next collaboration.', name: 'Rachel Han', role: 'Founder, Loopback' },
  { quote: 'They took the brief and dropped a product our team is proud of. The requirements landed exactly as specified.', name: 'Jonas Keller', role: 'Director, Verbatim' },
]

/* Criteria down the side, providers across the top. The shape of the data is
   the shape of the markup. */
const COMPARISON = {
  providers: ['This studio', 'Traditional agency'],
  rows: [
    { criterion: 'Approach', values: [true, false], detail: ['Design and engineering in sync', 'Handoff between siloed teams'] },
    { criterion: 'Process', values: [true, false], detail: ['Streamlined and transparent', 'Endless calls, vague timelines'] },
    { criterion: 'Design philosophy', values: [true, false], detail: ['Modern, minimal, purposeful', 'Template-based and dated'] },
    { criterion: 'Development stack', values: [true, false], detail: ['Built on current frameworks', 'Outdated stacks'] },
    { criterion: 'Communication', values: [true, false], detail: ['Clear updates on a schedule', 'Multiple middlemen'] },
    { criterion: 'Deliverables', values: [true, false], detail: ['Production-ready design systems', 'Static mockups'] },
    { criterion: 'Support', values: [true, false], detail: ['Long-term partnership', 'One-and-done projects'] },
  ],
}

const TIERS = [
  {
    name: 'Components',
    price: '$4,995',
    cadence: 'per month',
    blurb: 'Tailored components for fast-moving brands.',
    features: [
      'Custom strategy and wireframes',
      'Development in React or Next',
      'Unlimited revisions within scope',
      'Direct channel with the team',
      'Deployment handled',
    ],
    featured: false,
  },
  {
    name: 'Website',
    price: '$6,995',
    cadence: 'per month',
    blurb: 'A complete site, designed and built end to end.',
    features: [
      'Everything in Components',
      'Full multi-page site',
      'CMS integration',
      'Performance and SEO budget',
      'Ongoing maintenance',
    ],
    featured: true,
  },
]

const FAQS = [
  { q: 'What exactly does this platform do?', a: 'We design and build websites and product interfaces, from the first wireframe through to a deployed, maintained site.' },
  { q: 'What is a typical use case?', a: 'A funded startup that needs a marketing site and a product surface that agree with each other, without hiring a design team first.' },
  { q: 'Can I connect this with my existing stack?', a: 'Yes. We work in what you already run rather than arguing for a rewrite, unless the rewrite is genuinely the cheaper option.' },
  { q: 'How does the model selection work?', a: 'Scope is fixed up front and priced monthly. You can pause or stop at the end of any month.' },
  { q: 'Do I need to be technical to use this?', a: 'No. The handover includes a written architecture note and a recorded walkthrough rather than a scheduled call you have to attend.' },
  { q: 'How does onboarding work?', a: 'A kickoff call, a shared channel, and a preview URL inside the first week. Nothing waits on a contract being countersigned.' },
]

/* ── Mocks ──────────────────────────────────────────────────────────────── */

function Panel({ children, className = '' }: { children: ReactNode; className?: string }) {
  return (
    <div
      aria-hidden="true"
      className={`overflow-hidden rounded-xl border border-white/10 bg-white/5 select-none ${className}`}
    >
      {children}
    </div>
  )
}

function BrowserMock() {
  return (
    <Panel>
      <div className="flex gap-1.5 border-b border-white/10 px-3 py-2.5">
        {['bg-red-400/70', 'bg-amber-400/70', 'bg-emerald-400/70'].map((tone) => (
          <span key={tone} className={`size-2 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="space-y-2 p-4">
        <span className="block h-14 rounded-lg bg-gradient-to-br from-white/15 to-white/5" />
        <span className="block h-2 w-3/4 rounded bg-white/10" />
        <span className="block h-2 w-1/2 rounded bg-white/10" />
      </div>
    </Panel>
  )
}

function ProgressMock() {
  return (
    <Panel className="p-4">
      <p className="text-[11px] text-white/60">Sprint progress</p>
      <div className="mt-3 space-y-2.5">
        {[80, 55, 30].map((value) => (
          <div key={value} className="h-1.5 overflow-hidden rounded-full bg-white/10">
            <span style={{ width: `${value}%` }} className="block h-full rounded-full bg-amber-400" />
          </div>
        ))}
      </div>
    </Panel>
  )
}

function StackMock() {
  return (
    <Panel className="p-4">
      <div className="grid grid-cols-4 gap-2">
        {Array.from({ length: 8 }, (_, i) => (
          <span key={i} className="block h-8 rounded-lg bg-white/10" />
        ))}
      </div>
    </Panel>
  )
}

const MOCKS: Record<string, ReactNode> = {
  browser: <BrowserMock />,
  progress: <ProgressMock />,
  stack: <StackMock />,
}

function initials(name: string) {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageDarkAgencyWithComparison({
  brand = 'Aurelia',
  title = 'The design and engineering team you were about to hire',
  subtitle = 'We design and build websites that drive results and keep your business growing. Cal.com, No BS, real results.',
}: {
  brand?: string
  title?: string
  subtitle?: string
}) {
  return (
    /* Dark by commitment, set once here and inherited. A page whose design is
       a glow on near-black is not improved by a light mode bolted on. */
    <div className="bg-gray-950 text-white">
      <header className="border-b border-white/10">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-6xl items-center justify-between gap-6 px-6 py-4"
        >
          <a
            href="#"
            className="text-lg font-semibold tracking-tight focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-white/70 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <a
            href="#contact"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg bg-amber-400 px-4 text-sm font-medium text-gray-950 hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
          >
            Chat with us
          </a>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-white/10">
          <div
            aria-hidden="true"
            className="absolute -top-40 left-1/2 size-[36rem] -translate-x-1/2 rounded-full bg-amber-500/20 blur-3xl"
          />
          <div className="relative mx-auto grid max-w-6xl gap-10 px-6 py-24 lg:grid-cols-[1.4fr_1fr] lg:items-end">
            <h1 className="text-4xl font-bold tracking-tight text-balance sm:text-5xl">{title}</h1>
            <div>
              <p className="text-pretty text-white/70">{subtitle}</p>
              <a
                href="#contact"
                className="mt-6 inline-flex min-h-12 items-center rounded-lg bg-amber-400 px-6 text-sm font-medium text-gray-950 hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
              >
                Chat with us
              </a>
            </div>
          </div>
        </section>

        {/* Clients */}
        <section aria-labelledby="clients" className="border-b border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-14">
            <h2 id="clients" className="sr-only">
              Companies we have worked with
            </h2>
            <ul
              role="list"
              className="grid grid-cols-2 items-center justify-items-center gap-6 sm:grid-cols-4 lg:grid-cols-6"
            >
              {CLIENTS.map((client) => (
                /* Wordmarks, not logo files: a template shipping real marks
                   ships someone else's trademark. */
                <li key={client} className="text-sm font-semibold tracking-tight text-white/40">
                  {client}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Capabilities */}
        <section aria-labelledby="capabilities" className="border-b border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <h2 id="capabilities" className="text-3xl font-bold tracking-tight text-balance">
              Replace your engineering team
            </h2>
            <ul role="list" className="mt-10 grid gap-6 md:grid-cols-3">
              {CAPABILITIES.map((item) => (
                <li
                  key={item.title}
                  className="flex flex-col rounded-2xl border border-white/10 bg-white/[0.03] p-6"
                >
                  <div className="mb-5">{MOCKS[item.mock]}</div>
                  <h3 className="mt-auto font-semibold">{item.title}</h3>
                  <p className="mt-2 text-sm text-white/60">{item.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Projects */}
        <section aria-labelledby="projects" className="relative border-b border-white/10">
          {/* Texture. Read aloud it is the heading said twice. */}
          <span
            aria-hidden="true"
            className="pointer-events-none absolute inset-x-0 top-8 text-center text-[6rem] font-bold text-white/[0.04] select-none sm:text-[9rem]"
          >
            Projects
          </span>
          <div className="relative mx-auto max-w-6xl px-6 py-20">
            <h2 id="projects" className="sr-only">
              Selected projects
            </h2>
            {/* auto-rows so a row-span-2 tile is genuinely twice as tall. */}
            <ul role="list" className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 lg:auto-rows-[11rem]">
              {PROJECTS.map((project) => (
                <li key={project.name} className={project.span}>
                  <a
                    href="#"
                    className={`flex size-full min-h-44 items-end rounded-2xl border border-white/10 bg-gradient-to-br p-5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400 ${project.tone}`}
                  >
                    <span className="text-sm font-medium">{project.name}</span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Testimonials */}
        <section aria-labelledby="testimonials" className="border-b border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <h2 id="testimonials" className="text-3xl font-bold tracking-tight text-balance">
              Insights straight from our users
            </h2>
            <ul role="list" className="mt-10 grid gap-6 md:grid-cols-3">
              {TESTIMONIALS.map((item) => (
                <li key={item.name}>
                  <figure className="h-full rounded-2xl border border-white/10 bg-white/[0.03] p-6">
                    <blockquote className="text-sm text-pretty text-white/80">
                      {item.quote}
                    </blockquote>
                    <figcaption className="mt-5 flex items-center gap-3">
                      <span
                        aria-hidden="true"
                        className="flex size-9 shrink-0 items-center justify-center rounded-full bg-white/10 text-xs font-semibold"
                      >
                        {initials(item.name)}
                      </span>
                      <span>
                        <span className="block text-sm font-medium">{item.name}</span>
                        <span className="block text-xs text-white/50">{item.role}</span>
                      </span>
                    </figcaption>
                  </figure>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Comparison */}
        <section aria-labelledby="comparison" className="border-b border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20">
            <h2 id="comparison" className="text-3xl font-bold tracking-tight text-balance">
              {brand} versus a traditional agency
            </h2>
            <p className="mt-3 text-white/60">
              The same seven questions, answered honestly for both.
            </p>

            {/* A real table. Criteria are row headers and providers are column
                headers, so each cell is announced with both, instead of two
                lists of ticks with nothing tying a row's answers together. */}
            <div className="mt-10 overflow-x-auto">
              <table className="w-full border-collapse text-left text-sm">
                <caption className="sr-only">
                  How {brand} compares with a traditional agency across seven criteria
                </caption>
                <thead>
                  <tr className="border-b border-white/10">
                    <th scope="col" className="py-3 pr-4 font-medium text-white/60">
                      Criterion
                    </th>
                    {COMPARISON.providers.map((provider, index) => (
                      <th
                        key={provider}
                        scope="col"
                        className={`py-3 pl-4 font-medium ${index === 0 ? 'text-amber-300' : 'text-white/60'}`}
                      >
                        {provider}
                      </th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {COMPARISON.rows.map((row) => (
                    <tr key={row.criterion} className="border-b border-white/5">
                      <th scope="row" className="py-4 pr-4 font-medium">
                        {row.criterion}
                      </th>
                      {row.values.map((value, index) => (
                        <td key={index} className="py-4 pl-4 align-top">
                          <span className="flex items-start gap-2">
                            <span
                              aria-hidden="true"
                              className={`mt-0.5 flex size-4 shrink-0 items-center justify-center rounded-full text-[10px] ${
                                value ? 'bg-emerald-500/20 text-emerald-300' : 'bg-rose-500/20 text-rose-300'
                              }`}
                            >
                              {value ? '✓' : '✕'}
                            </span>
                            {/* The tick is a picture. Without this the cell
                                reads as just the detail, or as nothing. */}
                            <span className="sr-only">{value ? 'Yes. ' : 'No. '}</span>
                            <span className={value ? 'text-white/80' : 'text-white/50'}>
                              {row.detail[index]}
                            </span>
                          </span>
                        </td>
                      ))}
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="border-b border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20">
            <h2 id="pricing-heading" className="text-3xl font-bold tracking-tight text-balance">
              Extensive pricing plans
            </h2>

            <ul role="list" className="mt-10 grid items-start gap-6 md:grid-cols-2">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured ? 'border-amber-400/40 bg-amber-400/[0.06]' : 'border-white/10 bg-white/[0.03]'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="text-lg font-semibold">{tier.name}</h3>
                    {/* In words: colour alone reports nothing. */}
                    {tier.featured && (
                      <span className="rounded-full bg-amber-400 px-2.5 py-1 text-xs font-medium text-gray-950">
                        Most popular
                      </span>
                    )}
                  </div>
                  <p className="mt-2 text-sm text-white/60">{tier.blurb}</p>
                  <p className="mt-6">
                    <span className="text-4xl font-bold tracking-tight">{tier.price}</span>{' '}
                    <span className="text-sm text-white/50">{tier.cadence}</span>
                  </p>
                  <ul role="list" className="mt-6 space-y-2.5">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5 text-sm text-white/80">
                        <span aria-hidden="true" className="text-amber-300">
                          ✓
                        </span>
                        {feature}
                      </li>
                    ))}
                  </ul>
                  {/* Wrapper owns the spacing: mt-auto and mt-8 on one element
                      is a Tailwind conflict resolved by stylesheet order. */}
                  <div className="mt-auto pt-8">
                    <a
                      href="#contact"
                      className={`inline-flex min-h-12 w-full items-center justify-center rounded-lg text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 ${
                        tier.featured
                          ? 'bg-amber-400 text-gray-950 hover:bg-amber-300 focus-visible:outline-amber-400'
                          : 'border border-white/15 hover:bg-white/5 focus-visible:outline-white'
                      }`}
                    >
                      Get started
                    </a>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Founder's note */}
        <section aria-labelledby="founder" className="border-b border-white/10">
          <div className="mx-auto grid max-w-5xl gap-10 px-6 py-20 lg:grid-cols-[1fr_1.4fr] lg:items-center">
            <div
              aria-hidden="true"
              className="aspect-4/5 rounded-2xl bg-gradient-to-br from-white/15 to-white/5"
            />
            <div>
              <h2 id="founder" className="text-3xl font-bold tracking-tight text-balance">
                The founder&rsquo;s desk
              </h2>
              <p className="mt-5 text-pretty text-white/70">
                I started this studio because agencies kept selling process where a product was
                wanted. We are small on purpose: the person who designs your site is the person who
                builds it, and the person who answers when something breaks.
              </p>
              <p className="mt-4 text-pretty text-white/70">
                If that sounds like the arrangement you want, the fastest way to find out whether we
                suit each other is a conversation rather than a proposal.
              </p>
              <p className="mt-6 text-sm font-medium">Ada Whitfield, founder</p>
            </div>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq" className="border-b border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20">
            <h2 id="faq" className="text-3xl font-bold tracking-tight text-balance">
              Frequently asked questions
            </h2>
            <div className="mt-8">
              {FAQS.map((faq) => (
                <details key={faq.q} className="group border-b border-white/10 py-4">
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400">
                    {faq.q}
                    <span
                      aria-hidden="true"
                      className="mt-0.5 shrink-0 text-white/40 transition-transform group-open:rotate-45 motion-reduce:transition-none"
                    >
                      +
                    </span>
                  </summary>
                  <p className="mt-2.5 text-sm text-pretty text-white/60">{faq.a}</p>
                </details>
              ))}
            </div>
          </div>
        </section>

        {/* Closing ask */}
        <section id="contact" aria-labelledby="cta" className="relative overflow-hidden">
          <span
            aria-hidden="true"
            className="pointer-events-none absolute inset-x-0 bottom-0 text-center text-[6rem] font-bold text-white/[0.04] select-none sm:text-[9rem]"
          >
            {brand}
          </span>
          <div className="relative mx-auto max-w-3xl px-6 py-24 text-center">
            <h2 id="cta" className="text-3xl font-bold tracking-tight text-balance sm:text-4xl">
              Make your website a sales machine
            </h2>
            <a
              href="#"
              className="mt-8 inline-flex min-h-12 items-center rounded-lg bg-amber-400 px-6 text-sm font-medium text-gray-950 hover:bg-amber-300 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
            >
              Chat with us
            </a>
          </div>
        </section>
      </main>

      <footer className="border-t border-white/10">
        <div className="mx-auto flex max-w-6xl flex-col gap-6 px-6 py-10 sm:flex-row sm:items-center sm:justify-between">
          <p className="text-lg font-semibold tracking-tight">{brand}</p>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
              {[...NAV, 'Privacy', 'Terms'].map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="inline-flex min-h-11 items-center text-sm text-white/60 hover:text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-400"
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
