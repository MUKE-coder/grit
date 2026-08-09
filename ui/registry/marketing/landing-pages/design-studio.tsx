/*
 * A studio site: process, work, testimonial wall, services, pricing, FAQ.
 *
 * One file, because the registry installs one file per block. Treat it as a
 * page you own rather than a component to configure — the first job after
 * installing is to cut it into the sections you actually want.
 *
 * The section order is the pitch a studio has to make in order: what you do,
 * how you work, who trusted you, what it costs, what you are still worried
 * about. Pricing before proof reads as a quote from a stranger.
 *
 * Product surfaces are markup, never photography. A stock photo cannot be a
 * screenshot of your product, and a studio template full of someone else's
 * screenshots is a template you have to strip before you can use it.
 *
 * People and places are photographed, because the opposite trade applies: a
 * drawn avatar is a worse likeness of a person than a photograph of one, and a
 * studio selling taste cannot open on a grey rectangle. Every photograph here
 * is decorative — alt="" and aria-hidden — since the caption already names the
 * person and the headline already says what the band is for.
 *
 * The highlight behind "earn trust" and "start to finish" is a background on a
 * <mark>, not a coloured <span>. mark carries the meaning; several screen
 * readers announce entering and leaving it, which is the point of highlighting
 * a phrase. Its default yellow is overridden because the default is the one
 * thing that never matches a design.
 *
 * The FAQ is <details>/<summary>. Native disclosure gets keyboard support, the
 * expanded state announced, Ctrl+F opening a closed section, and it works
 * before hydration. A div with useState gets none of that.
 *
 * The recommended pricing tier says "Most popular" in text. Colour alone
 * carries no meaning to anyone not looking at it, and "the blue one" is not a
 * thing a screen reader can report.
 *
 * One <h1>, an <h2> per section, <h3> inside. A long page is navigated by
 * heading, and a page of styled divs cannot be navigated at all.
 */

import type { ReactNode } from 'react'

/* Verified on a contact sheet at the size they are shown before being used
   here. Three candidates across this set turned out to be something other than
   their description — a "lavender field" that was a scoop of soil, a "globe"
   that was a porthole — which is the whole reason the sheet exists. */
const PHOTO = {
  ocean: '1505142468610-359e7d316be0',
  earth: '1451187580459-43490279c0fa',
  consult: '1556740738-b6a63e27c4df',
}

const FACES = [
  '1500648767791-00dcc994a43e',
  '1531427186611-ecfd6d936c79',
  '1494790108377-be9c29b29330',
  '1580489944761-15a19d654956',
  '1528892952291-009c663ce843',
  '1534528741775-53994a69daeb',
]

const photo = (id: string, w: number, h: number) =>
  `https://images.unsplash.com/photo-${id}?w=${w}&h=${h}&fit=crop&q=80`
const face = (id: string, size: number) =>
  `https://images.unsplash.com/photo-${id}?w=${size * 2}&h=${size * 2}&fit=crop&crop=faces&q=75`

const NAV = ['Work', 'Services', 'Pricing', 'About']

const CLIENTS = ['Northwind', 'Meridian', 'Kestrel', 'Lumen', 'Atlas', 'Verdant']

const PROCESS = [
  {
    step: 'Design',
    body: 'Wireframes, then real screens in the browser. You see it working before anyone writes production code.',
    tone: 'bg-blue-600',
  },
  {
    step: 'Build',
    body: 'Typed, tested and deployed behind a preview URL from the first week, so review is continuous rather than a reveal.',
    tone: 'bg-amber-500',
  },
  {
    step: 'Detail',
    body: 'Motion, empty states, focus order and the loading cases. The parts that decide whether it feels finished.',
    tone: 'bg-rose-500',
  },
  {
    step: 'Handover',
    body: 'A repo you own, a written architecture note, and a walkthrough recorded rather than scheduled.',
    tone: 'bg-emerald-600',
  },
]

const SERVICES = [
  {
    title: 'Web design and build',
    body: 'From the first wireframe to production, on a stack your team can actually maintain.',
    mock: 'browser',
  },
  {
    title: 'Product and mobile',
    body: 'Responsive down to the small screens, with the interaction work that makes them feel native.',
    mock: 'phone',
  },
  {
    title: 'Words and search',
    body: 'Copy that says what the thing does, structured so search engines can agree with you.',
    mock: 'copy',
  },
  {
    title: 'Strategy and review',
    body: 'A read on what you have, what is costing you conversions, and what to do in what order.',
    mock: 'consult',
  },
]

const CAPABILITIES = [
  { title: 'Accessible by default', body: 'Keyboard paths, contrast and focus order checked before launch, not after a complaint.' },
  { title: 'Fast on real phones', body: 'Budgets set on mid-range hardware and a throttled connection, because that is the audience.' },
  { title: 'Light and dark', body: 'Both themes designed, not one theme inverted and hoped for.' },
  { title: 'Modern stack', body: 'Typed end to end, deployed on push, with no build step nobody remembers how to run.' },
  { title: 'Maintained after launch', body: 'Dependency updates and a monthly check, so the site does not quietly rot.' },
  { title: 'Measured', body: 'Analytics that answer a question you actually asked, rather than a dashboard nobody opens.' },
]

const TESTIMONIALS = [
  { quote: 'They pushed back on half our brief and were right about most of it. The site is better for the argument.', name: 'Marcus Chen', role: 'CTO, Kestrel Labs' },
  { quote: 'Quick to respond, genuinely professional, and shipped a site within a week. Looking forward to the next one.', name: 'Kenji Sato', role: 'Founder, Lumen' },
  { quote: 'From a loose brief to a polished site in days. Strong taste, strong intuition, zero fluff.', name: 'Priya Nair', role: 'Product lead, Meridian' },
  { quote: 'From rough concept to a shipped product, patient with feedback, good collaboration, and delivery matched our vision.', name: 'Amelia Brooks', role: 'Founder, Verdant' },
  { quote: 'Development was flawless and they acted like partners, not vendors. Very happy we hired them.', name: 'Jonas Keller', role: 'Director, Atlas' },
  { quote: 'Excellent work on our website. They went above and beyond, with a smooth handover we could move forward on confidently.', name: 'Riley Quinn', role: 'Founder, Cedar' },
]

const TIERS = [
  {
    name: 'One pager',
    price: '$2,500',
    cadence: 'one-time payment',
    blurb: 'A high-converting landing page, designed and built end to end.',
    features: [
      'Single high-converting landing page',
      'Custom design, no template',
      'Two revision rounds',
      'Delivered in one to two weeks',
      'SEO, hosting and deployment handled',
      '24-hour support on a private channel',
      'Add on: extra page for $500',
    ],
    cta: 'Start your landing page',
    featured: false,
  },
  {
    name: 'Company site',
    price: '$4,500',
    cadence: 'one-time payment',
    blurb: 'A complete multi-page site that makes the whole business legible.',
    features: [
      'Multiple pages, up to six',
      'CMS integration so you can edit it',
      'Three revision rounds',
      'Delivered in three to four weeks',
      'Everything in One pager',
      'Add on: extra page for $500',
    ],
    cta: 'Start your site',
    featured: true,
  },
]

const FAQS = [
  { q: 'What stack do you use for coding?', a: 'Typed React with a Go or Node backend depending on what you already run. If you have an existing stack we work in it rather than arguing for a rewrite.' },
  { q: 'How do I reach you for support?', a: 'A shared channel on Slack or Discord, plus email. During a project you get a response the same working day.' },
  { q: 'Why should I choose your team for my project?', a: 'Small team, senior people, no handoff to juniors after the pitch. The person who designs it is the person who builds it.' },
  { q: 'Do you do communication progress during a project?', a: 'A preview URL from week one and a written update every Friday. You never have to ask where it is.' },
  { q: 'What happens if I need changes after delivery?', a: 'Revision rounds are included in the price. After those, changes are billed hourly at a rate agreed up front, with an estimate before anything starts.' },
  { q: 'Do you offer custom solutions or only predefined packages?', a: 'The packages cover the common shapes. Anything else is scoped and quoted; get in touch and we will tell you honestly whether we are the right studio for it.' },
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

function BrowserMock() {
  return (
    <Panel>
      <div className="flex items-center gap-1.5 border-b border-gray-200 px-3 py-2.5 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2 rounded-full ${tone}`} />
        ))}
        <span className="ml-2 h-3 w-32 rounded bg-gray-100 dark:bg-white/10" />
      </div>
      <div className="space-y-2.5 p-4">
        <span className="block h-16 rounded-lg bg-gradient-to-br from-blue-100 to-blue-50 dark:from-blue-500/20 dark:to-blue-500/5" />
        <span className="block h-2 w-3/4 rounded bg-gray-200 dark:bg-white/10" />
        <span className="block h-2 w-1/2 rounded bg-gray-200 dark:bg-white/10" />
        <div className="grid grid-cols-3 gap-2 pt-1">
          {[0, 1, 2].map((i) => (
            <span key={i} className="block h-10 rounded bg-gray-100 dark:bg-white/5" />
          ))}
        </div>
      </div>
    </Panel>
  )
}

function PhoneMock() {
  return (
    <div aria-hidden="true" className="flex justify-center py-4 select-none">
      <div className="w-28 rounded-[1.25rem] border-4 border-gray-900 bg-white p-1.5 dark:border-white/20 dark:bg-gray-900">
        <span className="mx-auto mb-1.5 block h-1 w-8 rounded-full bg-gray-300 dark:bg-white/20" />
        <span className="block h-14 rounded-md bg-gradient-to-br from-sky-200 to-sky-50 dark:from-sky-500/25 dark:to-sky-500/5" />
        <span className="mt-1.5 block h-1.5 w-3/4 rounded bg-gray-200 dark:bg-white/10" />
        <span className="mt-1 block h-1.5 w-1/2 rounded bg-gray-200 dark:bg-white/10" />
      </div>
    </div>
  )
}

function CopyMock() {
  return (
    <Panel className="p-4">
      <p className="font-mono text-[11px] text-gray-500 dark:text-gray-400">&lt;h1&gt;</p>
      <span className="mt-1.5 block h-2.5 w-4/5 rounded bg-gray-300 dark:bg-white/20" />
      <p className="mt-3 font-mono text-[11px] text-gray-500 dark:text-gray-400">&lt;p&gt;</p>
      {[0, 1, 2].map((i) => (
        <span
          key={i}
          className={`mt-1.5 block h-2 rounded bg-gray-200 dark:bg-white/10 ${
            i === 2 ? 'w-1/2' : 'w-full'
          }`}
        />
      ))}
    </Panel>
  )
}

function ChartMock() {
  const bars = [40, 65, 45, 80, 60, 95]
  return (
    <Panel className="p-4">
      <div className="flex h-24 items-end gap-2">
        {bars.map((height, i) => (
          <span
            key={i}
            style={{ height: `${height}%` }}
            className="flex-1 rounded-t bg-gradient-to-t from-blue-600 to-blue-400"
          />
        ))}
      </div>
      <div className="mt-3 flex gap-2">
        <span className="h-2 w-16 rounded bg-gray-200 dark:bg-white/10" />
        <span className="h-2 w-10 rounded bg-gray-200 dark:bg-white/10" />
      </div>
    </Panel>
  )
}

function ConsultPhoto() {
  return (
    <img
      src={photo(PHOTO.consult, 800, 500)}
      alt=""
      aria-hidden="true"
      loading="lazy"
      className="h-44 w-full rounded-xl object-cover"
    />
  )
}

const MOCKS: Record<string, ReactNode> = {
  browser: <BrowserMock />,
  phone: <PhoneMock />,
  copy: <CopyMock />,
  chart: <ChartMock />,
  consult: <ConsultPhoto />,
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageDesignStudio({
  brand = 'Mainline',
  email = 'hello@example.com',
}: {
  brand?: string
  email?: string
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
          <a
            href="#contact"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg bg-blue-700 px-4 text-sm font-medium text-white hover:bg-blue-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
          >
            Talk to us
          </a>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-gray-200 dark:border-white/10">
          {/* The photograph is a band above the headline, and no text sits on
              it. The design this follows runs the navigation across the wave,
              which is the one arrangement that cannot be made safe: contrast
              against a photograph is a different number per pixel, so a link
              that reads over the foam disappears over the water two hundred
              pixels along. Below the band the gradient takes over and the type
              sits on a colour that can be measured. */}
          <img
            src={photo(PHOTO.ocean, 1600, 500)}
            alt=""
            aria-hidden="true"
            className="h-44 w-full object-cover sm:h-56"
          />
          <div
            aria-hidden="true"
            className="absolute inset-x-0 top-44 bottom-0 bg-gradient-to-b from-sky-100 via-white to-white sm:top-56 dark:from-sky-500/10 dark:via-gray-950 dark:to-gray-950"
          />
          <div className="relative mx-auto max-w-3xl px-6 py-20 text-center">
            <h1 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              Design-led websites that{' '}
              {/* mark, not a span: the highlight is meaning, and the default
                  yellow is overridden because it matches no design ever. */}
              <mark className="bg-blue-200 px-1 text-gray-900 dark:bg-blue-500/30 dark:text-white">
                earn trust
              </mark>{' '}
              in the first scroll
            </h1>
            <p className="mx-auto mt-6 max-w-xl text-lg text-pretty text-gray-600 dark:text-gray-300">
              We build sites with a clear hierarchy, deliberate motion and the performance to hold
              up when real traffic arrives.
            </p>
            <div className="mt-8 flex flex-wrap justify-center gap-3">
              <a
                href="#contact"
                className="inline-flex min-h-12 items-center rounded-lg bg-blue-700 px-6 text-sm font-medium text-white hover:bg-blue-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
              >
                Talk to us
              </a>
              <a
                href="#work"
                className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
              >
                Explore our work
              </a>
            </div>
          </div>
        </section>

        {/* Clients */}
        <section aria-labelledby="clients" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-12">
            <h2 id="clients" className="text-center text-sm text-gray-500 dark:text-gray-400">
              Some of the companies we have worked with
            </h2>
            <ul
              role="list"
              className="mt-6 grid grid-cols-2 items-center justify-items-center gap-6 sm:grid-cols-3 lg:grid-cols-6"
            >
              {CLIENTS.map((client) => (
                <li
                  key={client}
                  /* gray-600 on paper and gray-400 on the dark theme. The
                      gray-400/gray-500 pair this started as measured 2.54:1 and
                      4.16:1 — a logo cloud is quiet by design, but quiet is not
                      the same as unreadable. */
                  className="text-sm font-semibold tracking-tight text-gray-600 dark:text-gray-400"
                >
                  {/* Wordmarks, not logo files. A template shipping real marks
                      is a template shipping someone else's trademark. */}
                  {client}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Process */}
        <section aria-labelledby="process" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <h2
              id="process"
              className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              A process tuned for momentum
            </h2>
            <p className="mt-3 max-w-xl text-gray-600 dark:text-gray-300">
              Design, build, polish and ship in sync. Nothing gets lost in a handoff.
            </p>

            {/* Ordered, because the steps happen in this order. */}
            <ol className="mt-10 grid gap-6 sm:grid-cols-2">
              {PROCESS.map((phase, index) => (
                <li key={phase.step} className="flex flex-col">
                  {MOCKS[index % 2 === 0 ? 'browser' : 'phone']}
                  <div className="mt-4 flex items-start gap-3">
                    <span
                      aria-hidden="true"
                      className={`flex size-6 shrink-0 items-center justify-center rounded text-xs font-bold text-white ${phase.tone}`}
                    >
                      {index + 1}
                    </span>
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {phase.step}
                      </span>{' '}
                      {phase.body}
                    </p>
                  </div>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* Testimonial wall */}
        <section
          id="work"
          aria-labelledby="testimonials"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <h2
              id="testimonials"
              className="text-center text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              We have worked with some of the best teams around
            </h2>

            {/* CSS columns rather than a masonry grid. Each card is
                independent, so reading down a column loses nothing. */}
            <ul role="list" className="mt-12 gap-5 sm:columns-2 lg:columns-3">
              {TESTIMONIALS.map((item, index) => {
                const featured = index === 1
                return (
                  <li key={item.name} className="mb-5 break-inside-avoid">
                    <figure
                      className={`rounded-xl border p-5 ${
                        featured
                          ? 'border-blue-700 bg-blue-700'
                          : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                      }`}
                    >
                      <blockquote
                        className={`text-sm text-pretty ${
                          featured ? 'text-white' : 'text-gray-700 dark:text-gray-200'
                        }`}
                      >
                        {item.quote}
                      </blockquote>
                      <figcaption className="mt-4 flex items-center gap-3">
                        <img
                          src={face(FACES[index % FACES.length], 36)}
                          alt=""
                          aria-hidden="true"
                          loading="lazy"
                          className="size-9 shrink-0 rounded-full object-cover"
                        />
                        <span>
                          <span
                            className={`block text-sm font-medium ${
                              featured ? 'text-white' : 'text-gray-900 dark:text-white'
                            }`}
                          >
                            {item.name}
                          </span>
                          <span
                            className={`block text-xs ${
                              featured ? 'text-blue-100' : 'text-gray-500 dark:text-gray-400'
                            }`}
                          >
                            {item.role}
                          </span>
                        </span>
                      </figcaption>
                    </figure>
                  </li>
                )
              })}
            </ul>
          </div>
        </section>

        {/* Services */}
        <section
          id="services"
          aria-labelledby="services-heading"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <h2
              id="services-heading"
              className="text-center text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              We take care of things{' '}
              <mark className="bg-blue-200 px-1 text-gray-900 dark:bg-blue-500/30 dark:text-white">
                from start to finish
              </mark>
            </h2>

            <ul role="list" className="mt-12 grid gap-6 md:grid-cols-2">
              {SERVICES.map((service) => (
                <li
                  key={service.title}
                  className="flex flex-col rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40"
                >
                  <div className="mb-5">{MOCKS[service.mock]}</div>
                  <h3 className="mt-auto font-semibold text-gray-900 dark:text-white">
                    {service.title}
                  </h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{service.body}</p>
                </li>
              ))}
            </ul>

            {/* Everything-else band. The earth is scenery, so it carries a
                scrim and the text on it is white — the same treatment every
                photograph on every other block in this set gets, for the same
                reason. */}
            <div className="relative isolate mt-12 overflow-hidden rounded-2xl">
              <img
                src={photo(PHOTO.earth, 1600, 600)}
                alt=""
                aria-hidden="true"
                loading="lazy"
                className="absolute inset-0 -z-10 size-full object-cover"
              />
              <span aria-hidden="true" className="absolute inset-0 -z-10 bg-gray-950/65" />
              <div className="px-6 py-14 text-center">
                <h3 className="text-2xl font-bold tracking-tight text-balance text-white">
                  Anything and everything web
                </h3>
                <p className="mx-auto mt-3 max-w-md text-pretty text-gray-100">
                  Deployments, hosting, maintenance, revisions, design systems — and everything else
                  in between.
                </p>
              </div>
            </div>

            <ul role="list" className="mt-12 grid gap-x-8 gap-y-8 sm:grid-cols-2 lg:grid-cols-3">
              {CAPABILITIES.map((item) => (
                <li key={item.title}>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {item.title}
                  </h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{item.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Pricing */}
        <section
          id="pricing"
          aria-labelledby="pricing-heading"
          className="border-b border-gray-200 dark:border-white/10"
        >
          <div className="mx-auto max-w-5xl px-6 py-20">
            <h2
              id="pricing-heading"
              className="text-center text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              Straightforward, one-time pricing that works for your business
            </h2>

            <ul role="list" className="mt-12 grid items-start gap-6 md:grid-cols-2">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured
                      ? 'border-blue-700 bg-blue-700 text-white'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3
                      className={`text-lg font-semibold ${
                        tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'
                      }`}
                    >
                      {tier.name}
                    </h3>
                    {/* In words, because "the blue one" is not something a
                        screen reader can report. */}
                    {tier.featured && (
                      /* Solid white with blue-800 text. Over the blue tier a
                          white label on bg-white/20 measured 4.42:1: a
                          translucent white lifts the surface just enough to
                          sink the text sitting on it. */
                      <span className="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-blue-800">
                        Most popular
                      </span>
                    )}
                  </div>

                  <p
                    className={`mt-2 text-sm ${
                      tier.featured ? 'text-blue-100' : 'text-gray-600 dark:text-gray-300'
                    }`}
                  >
                    {tier.blurb}
                  </p>

                  <p className="mt-6">
                    <span
                      className={`text-4xl font-bold tracking-tight ${
                        tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'
                      }`}
                    >
                      {tier.price}
                    </span>{' '}
                    <span
                      className={`text-sm ${
                        tier.featured ? 'text-blue-100' : 'text-gray-500 dark:text-gray-400'
                      }`}
                    >
                      {tier.cadence}
                    </span>
                  </p>

                  <ul role="list" className="mt-6 space-y-2.5">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5 text-sm">
                        <span
                          aria-hidden="true"
                          className={tier.featured ? 'text-blue-200' : 'text-blue-700 dark:text-blue-400'}
                        >
                          ✓
                        </span>
                        <span className={tier.featured ? 'text-white' : 'text-gray-700 dark:text-gray-200'}>
                          {feature}
                        </span>
                      </li>
                    ))}
                  </ul>

                  {/* mt-auto on the wrapper, not the anchor: putting mt-auto
                      and mt-8 on one element is a Tailwind conflict where
                      whichever class lands later in the stylesheet wins, which
                      is not something you get to choose. The wrapper absorbs
                      the free space so both buttons sit on the same line
                      however many features each tier lists. */}
                  <div className="mt-auto pt-8">
                    <a
                      href="#contact"
                      className={`inline-flex min-h-12 w-full items-center justify-center rounded-lg text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 ${
                        tier.featured
                          ? 'bg-white text-blue-800 hover:bg-blue-50 focus-visible:outline-white'
                          : 'bg-blue-700 text-white hover:bg-blue-800 focus-visible:outline-blue-600'
                      }`}
                    >
                      {tier.cta}
                    </a>
                  </div>
                </li>
              ))}
            </ul>

            <p className="mt-8 text-center text-sm text-gray-600 dark:text-gray-300">
              Fixed scope, fixed price. No hourly billing, no surprises.{' '}
              <a
                href={`mailto:${email}`}
                className="font-medium text-blue-700 underline-offset-4 hover:underline dark:text-blue-400"
              >
                {email}
              </a>
            </p>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20">
            <h2
              id="faq"
              className="text-center text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              Have questions? Browse the answers
            </h2>

            <div className="mt-10 grid gap-x-8 md:grid-cols-2">
              {FAQS.map((faq) => (
                <details
                  key={faq.q}
                  className="group border-b border-gray-200 py-4 dark:border-white/10"
                >
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-white">
                    {faq.q}
                    <span
                      aria-hidden="true"
                      className="mt-0.5 shrink-0 text-gray-400 transition-transform group-open:rotate-45 motion-reduce:transition-none"
                    >
                      +
                    </span>
                  </summary>
                  <p className="mt-2.5 text-sm text-pretty text-gray-600 dark:text-gray-300">
                    {faq.a}
                  </p>
                </details>
              ))}
            </div>
          </div>
        </section>

        {/* Closing ask */}
        <section id="contact" aria-labelledby="cta" className="bg-gray-50 dark:bg-gray-900/40">
          <div className="mx-auto max-w-3xl px-6 py-20 text-center">
            <h2
              id="cta"
              className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              Let us get your site{' '}
              <mark className="bg-blue-200 px-1 text-gray-900 dark:bg-blue-500/30 dark:text-white">
                built and deployed
              </mark>{' '}
              in no time
            </h2>
            <div className="mt-8 flex flex-wrap justify-center gap-3">
              <a
                href={`mailto:${email}`}
                className="inline-flex min-h-12 items-center rounded-lg bg-blue-700 px-6 text-sm font-medium text-white hover:bg-blue-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
              >
                Talk to us
              </a>
              <a
                href="#pricing"
                className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
              >
                Explore pricing
              </a>
            </div>
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
              Design-led websites that hold up under real traffic.
            </p>
          </div>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
              {[...NAV, 'Privacy', 'Terms'].map((item) => (
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
        </div>
      </footer>
    </div>
  )
}
