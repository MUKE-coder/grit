/*
 * A recruiting product page: hero, process, feature bento, reach, pricing,
 * testimonials, FAQ and a newsletter.
 *
 * One file, because the registry installs one file per block.
 *
 * The green is emerald-700 wherever white text sits on it, and emerald-500 for
 * link text in dark mode. Measured rather than picked: white on emerald-500 is
 * 2.54:1 and on emerald-600 is 3.77:1, both under the 4.5:1 WCAG AA wants for
 * a label, while emerald-700 is 5.48:1. In dark mode it inverts — emerald-700
 * on a near-black panel falls to 3.23:1. A bright green button with white text
 * is the same failure as a bright orange one, and both are everywhere.
 *
 * The featured pricing tier is raised visually and says "Most popular" in
 * text. Elevation, scale and colour are three ways of saying the same thing to
 * the same audience; none of them reaches anyone reading the page rather than
 * looking at it, so the words do the work and the styling decorates it.
 *
 * The reach section's figures are a <dl>. "23,000+" beside "happy customers"
 * is a term and its definition, and the pairing should be announced rather
 * than inferred from two boxes being near each other.
 *
 * Photography is used for people. Product surfaces stay as markup, because a
 * stock photo cannot be a screenshot of your product.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const NAV = ['Home', 'Pricing', 'Blog']

const PROCESS = [
  {
    title: 'Centralise the search',
    body: 'One place for every candidate, with easy access to the whole database rather than four inboxes.',
  },
  {
    title: 'Accelerate the speed',
    body: 'Shortlist candidates, make notes as you scan, and decide later with the reasons still attached.',
  },
  {
    title: 'Track every stage',
    body: 'See what is backed up and where, so you never lose a week to a stage nobody owned.',
  },
]

const FEATURES = [
  { title: 'Resume upload', body: 'Drop a batch of CVs and get a parsed, searchable record with the original still attached.', mock: 'upload', wide: false },
  { title: 'Interview feedback', body: 'Structured scoring against the questions you set, so two interviewers can actually be compared.', mock: 'feedback', wide: false },
  { title: 'Ranking that holds up', body: 'Rank on the criteria you defined at the start rather than the impression left by the last call.', mock: 'chart', wide: true },
]

const REACH = [
  { value: '23,000+', label: 'Happy customers worldwide' },
  { value: '120', label: 'Countries supported' },
  { value: '4.9 / 5', label: 'Average customer rating' },
]

const TIERS = [
  {
    name: 'Hobby',
    price: 99,
    features: ['Access to basic analytics reports', 'Up to 10,000 data points per month', 'Email support', 'Community forum access', 'Cancel any time'],
    inherits: null as string | null,
    featured: false,
  },
  {
    name: 'Starter',
    price: 299,
    features: ['Advanced analytics dashboard', 'Customisable reports and charts', 'Real-time data tracking', 'Integration with third-party tools'],
    inherits: 'Hobby',
    featured: true,
  },
  {
    name: 'Pro',
    price: 1490,
    features: ['Unlimited data storage', 'Customisable dashboards', 'Advanced data segmentation', 'AI-powered insights'],
    inherits: 'Starter',
    featured: false,
  },
]

/* Verified on a contact sheet before use. */
const PORTRAIT = '1507003211169-0a1dd7228f2d'

const TESTIMONIALS = [
  { quote: 'We hired our first candidate through it inside a fortnight, which was fast enough that I assumed we had cut a corner somewhere. We had not. What changed is that the notes from the first screen were still attached at the offer stage, so nobody had to reconstruct why we liked her from memory. Three hires later that is still the part I would not give up.', name: 'Marcus Kirby', role: 'Head of talent, Northwind' },
  { quote: 'Working with this platform has transformed our recruitment process. The AI-powered matching system saved us countless hours in finding the perfect candidates.', name: 'Sarah Chen', role: 'Head of talent, Northwind' },
  { quote: 'The level of customisation and flexibility in the platform is outstanding. We have been able to adapt it perfectly to our unique hiring workflows.', name: 'Emily Nakamura', role: 'Talent acquisition, Kestrel' },
  { quote: 'The platform’s intuitive interface and powerful analytics have completely revolutionised how we approach talent acquisition.', name: 'James Rodriguez', role: 'VP of HR, Meridian' },
]

const FAQS = [
  { q: 'What is the purpose of this platform?', a: 'It records and centralises interviews so feedback is comparable across interviewers, rather than living in six sets of private notes.' },
  { q: 'How do I contact support?', a: 'A shared channel plus email. During onboarding you get a response the same working day.' },
  { q: 'How do I find the best candidates?', a: 'Rank against the criteria you set before the first interview. The tool will not invent a shortlist for you, and should not.' },
  { q: 'Can I export my data?', a: 'Yes, in full, at any time. Your candidate records are yours and leaving is not an obstacle course.' },
  { q: 'Do you support international hiring?', a: 'We support over a hundred countries, including the data-residency rules that come with several of them.' },
  { q: 'How do I track my pipeline?', a: 'Every stage shows what is waiting and for how long, so a bottleneck is visible before it costs you a candidate.' },
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

function UploadMock() {
  return (
    <Panel className="p-4">
      <div className="rounded-lg border border-dashed border-gray-300 p-5 text-center text-[11px] text-gray-500 dark:border-white/15 dark:text-gray-400">
        Drop CVs here
      </div>
      <div className="mt-3 space-y-1.5">
        {['ada-lovelace.pdf', 'grace-hopper.pdf'].map((file) => (
          <div key={file} className="flex items-center justify-between text-[11px]">
            <span className="text-gray-600 dark:text-gray-300">{file}</span>
            <span className="text-emerald-700 dark:text-emerald-400">Parsed</span>
          </div>
        ))}
      </div>
    </Panel>
  )
}

function FeedbackMock() {
  return (
    <Panel className="p-4">
      {['Communication', 'Systems design', 'Ownership'].map((row, i) => (
        <div key={row} className="py-1.5">
          <p className="text-[11px] text-gray-600 dark:text-gray-300">{row}</p>
          <div className="mt-1 h-1.5 overflow-hidden rounded-full bg-gray-200 dark:bg-white/10">
            <span
              style={{ width: `${[85, 60, 72][i]}%` }}
              className="block h-full rounded-full bg-emerald-600"
            />
          </div>
        </div>
      ))}
    </Panel>
  )
}

function ChartMock() {
  return (
    <Panel className="p-4">
      <div className="flex h-24 items-end gap-2">
        {[45, 70, 55, 85, 62, 92, 74].map((h, i) => (
          <span
            key={i}
            style={{ height: `${h}%` }}
            className="flex-1 rounded-t bg-gradient-to-t from-emerald-700 to-emerald-500"
          />
        ))}
      </div>
    </Panel>
  )
}

const MOCKS: Record<string, ReactNode> = {
  upload: <UploadMock />,
  feedback: <FeedbackMock />,
  chart: <ChartMock />,
}

function initials(name: string) {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageRecruitingPlatform({
  brand = 'Playful',
  title = 'Record interviews. Centralise feedback automatically.',
  subtitle = 'Record and organise every interview, so hiring decisions rest on what candidates actually said rather than on who remembers it best.',
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
            className="text-lg font-semibold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700 dark:text-white"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700 dark:text-gray-300 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          {/* emerald-700: white on emerald-600 is 3.77:1 and fails AA. */}
          <a
            href="#pricing"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg bg-emerald-700 px-4 text-sm font-medium text-white hover:bg-emerald-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700"
          >
            Get started
          </a>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="border-b border-gray-200 bg-emerald-50/40 dark:border-white/10 dark:bg-gray-900/40">
          <div className="mx-auto max-w-3xl px-6 py-20 text-center">
            <p className="inline-flex items-center rounded-full border border-emerald-700/20 bg-white px-3 py-1 text-xs font-medium text-emerald-800 dark:border-emerald-500/30 dark:bg-transparent dark:text-emerald-300">
              Interviews recorded without a meeting bot
            </p>
            <h1 className="mt-6 text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h1>
            <p className="mx-auto mt-6 text-pretty text-gray-600 dark:text-gray-300">{subtitle}</p>
            <a
              href="#pricing"
              className="mt-8 inline-flex min-h-12 items-center rounded-lg bg-emerald-700 px-6 text-sm font-medium text-white hover:bg-emerald-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700"
            >
              Get started, it is free
            </a>
          </div>
        </section>

        {/* Process */}
        <section aria-labelledby="process" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-5xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="process"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Go from question to hired
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Everything needed to evaluate candidates, track pipelines and end an engineer&rsquo;s
                week knowing where things stand.
              </p>
            </div>

            <ol className="mt-12 grid gap-8 sm:grid-cols-3">
              {PROCESS.map((step, index) => (
                <li key={step.title}>
                  <span
                    aria-hidden="true"
                    className="flex size-9 items-center justify-center rounded-lg bg-emerald-700 text-sm font-bold text-white"
                  >
                    {index + 1}
                  </span>
                  <h3 className="mt-4 font-semibold text-gray-900 dark:text-white">{step.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{step.body}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* Features */}
        <section
          aria-labelledby="features"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="features"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Features that make the week easier
              </h2>
            </div>

            {/* 1 + 1 + 2 over two columns, which tiles exactly. */}
            <ul role="list" className="mt-12 grid gap-6 md:grid-cols-2">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`flex flex-col rounded-2xl border border-gray-200 bg-white p-6 dark:border-white/10 dark:bg-gray-900 ${
                    feature.wide ? 'md:col-span-2' : ''
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

        {/* Reach */}
        <section aria-labelledby="reach" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-4xl px-6 py-20 text-center">
            <h2
              id="reach"
              className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              We are available everywhere
            </h2>
            <p className="mt-3 text-gray-600 dark:text-gray-300">
              Available in every country we can legally support, which is most of them.
            </p>

            {/* dl: a figure and what it counts are a term and its definition. */}
            <dl className="mt-12 grid gap-8 sm:grid-cols-3">
              {REACH.map((item) => (
                <div key={item.label}>
                  <dt className="sr-only">{item.label}</dt>
                  <dd>
                    <span className="block text-3xl font-bold tracking-tight tabular-nums text-emerald-700 dark:text-emerald-400">
                      {item.value}
                    </span>
                    <span className="mt-1 block text-sm text-gray-600 dark:text-gray-300">
                      {item.label}
                    </span>
                  </dd>
                </div>
              ))}
            </dl>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="pricing-heading"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Pricing so simple you would buy instantly
              </h2>
            </div>

            <ul role="list" className="mt-12 grid items-center gap-6 lg:grid-cols-3">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  /* Raised and scaled, and it also says so in words. Elevation
                     and colour are two ways of telling the same audience the
                     same thing; neither reaches anyone reading rather than
                     looking. */
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured
                      ? 'border-emerald-700 bg-emerald-700 text-white shadow-xl lg:scale-105'
                      : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3
                      className={`font-semibold ${tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'}`}
                    >
                      {tier.name}
                    </h3>
                    {tier.featured && (
                      /* Solid white with emerald-800 text. On bg-white/20 over
                         emerald-700 the white label measured 3.75:1, because a
                         translucent white lifts the surface just enough to
                         sink the text. */
                      <span className="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-emerald-800">
                        Most popular
                      </span>
                    )}
                  </div>

                  <p className="mt-5">
                    <span
                      className={`align-super text-lg ${tier.featured ? 'text-emerald-100' : 'text-gray-500 dark:text-gray-400'}`}
                    >
                      $
                    </span>
                    <span
                      className={`text-4xl font-bold tracking-tight tabular-nums ${tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'}`}
                    >
                      {tier.price}
                    </span>
                    <span
                      className={`text-sm ${tier.featured ? 'text-emerald-100' : 'text-gray-500 dark:text-gray-400'}`}
                    >
                      /month
                    </span>
                  </p>

                  <a
                    href="#"
                    className={`mt-6 inline-flex min-h-12 w-full items-center justify-center rounded-lg text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 ${
                      tier.featured
                        ? 'bg-white text-emerald-800 hover:bg-emerald-50 focus-visible:outline-white'
                        : 'bg-emerald-700 text-white hover:bg-emerald-800 focus-visible:outline-emerald-700'
                    }`}
                  >
                    Get {tier.name}
                  </a>

                  <ul role="list" className="mt-6 space-y-2.5">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5 text-sm">
                        <span
                          aria-hidden="true"
                          className={tier.featured ? 'text-emerald-200' : 'text-emerald-700 dark:text-emerald-400'}
                        >
                          ✓
                        </span>
                        <span className={tier.featured ? 'text-white' : 'text-gray-700 dark:text-gray-200'}>
                          {feature}
                        </span>
                      </li>
                    ))}
                    {tier.inherits && (
                      /* Names the tier: "the previous plan" depends on column
                         order, which is not conveyed to a screen reader. */
                      <li
                        className={`flex gap-2.5 border-t pt-3 text-sm ${
                          tier.featured ? 'border-white/20' : 'border-gray-200 dark:border-white/10'
                        }`}
                      >
                        <span
                          aria-hidden="true"
                          className={tier.featured ? 'text-emerald-200' : 'text-emerald-700 dark:text-emerald-400'}
                        >
                          ✓
                        </span>
                        <span
                          className={`font-medium ${tier.featured ? 'text-white' : 'text-gray-900 dark:text-white'}`}
                        >
                          Everything in {tier.inherits}
                        </span>
                      </li>
                    )}
                  </ul>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Testimonials */}
        <section
          aria-labelledby="testimonials"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="testimonials"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Recruiters love us
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                People have chosen us from all over the world to help them with their hiring.
              </p>
            </div>

            <ul role="list" className="mt-12 gap-5 sm:columns-2 lg:columns-3">
              {TESTIMONIALS.map((item, index) => (
                <li key={item.name} className="mb-5 break-inside-avoid">
                  <figure className="rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900">
                    {/* One card carries a photograph, which is where a
                        photograph belongs on this page: a person. */}
                    {index === 3 && (
                      <img
                        src={`https://images.unsplash.com/photo-${PORTRAIT}?w=600&h=400&fit=crop&q=80`}
                        alt=""
                        loading="lazy"
                        className="mb-4 h-40 w-full rounded-lg object-cover"
                      />
                    )}
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

        {/* FAQ */}
        <section aria-labelledby="faq" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-3xl px-6 py-20">
            <h2
              id="faq"
              className="text-center text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
            >
              Frequently asked questions
            </h2>
            <div className="mt-10">
              {FAQS.map((faq) => (
                <details
                  key={faq.q}
                  className="group border-b border-gray-200 py-4 dark:border-white/10"
                >
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700 dark:text-white">
                    {faq.q}
                    <span
                      aria-hidden="true"
                      /* gray-500, not gray-400. list-none hides the native
                          disclosure marker, so this glyph is the only visual
                          affordance for expand and collapse. That makes it
                          functional rather than decorative, and it owes the
                          3:1 of WCAG 1.4.11 -- gray-400 on white measured
                          2.54:1. */
                      className="mt-0.5 shrink-0 text-gray-500 transition-transform group-open:rotate-45 motion-reduce:transition-none dark:text-gray-400"
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
      </main>

      <footer className="border-t border-gray-200 dark:border-white/10">
        <div className="mx-auto grid max-w-6xl gap-8 px-6 py-12 sm:grid-cols-2">
          <div>
            <p className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">
              {brand}
            </p>
            <p className="mt-2 max-w-sm text-2xl font-bold tracking-tight text-balance text-gray-900 dark:text-white">
              {title}
            </p>
          </div>

          <div className="grid grid-cols-2 gap-8 sm:justify-items-end">
            {[
              { heading: 'Pages', items: ['Home', 'Pricing', 'Blog', 'Contact'] },
              { heading: 'Legal', items: ['Privacy policy', 'Terms of service', 'Cookie policy'] },
            ].map((column) => (
              <nav key={column.heading} aria-label={column.heading}>
                <h2 className="text-sm font-semibold text-gray-900 dark:text-white">
                  {column.heading}
                </h2>
                <ul role="list" className="mt-1">
                  {column.items.map((item) => (
                    <li key={item}>
                      <a
                        href="#"
                        className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-emerald-700 dark:text-gray-400 dark:hover:text-white"
                      >
                        {item}
                      </a>
                    </li>
                  ))}
                </ul>
              </nav>
            ))}
          </div>
        </div>
      </footer>
    </div>
  )
}
