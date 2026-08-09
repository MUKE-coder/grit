/*
 * An AI automation agency page in an editorial, almost print register: numbered
 * section markers, hairline rules, monospaced labels, two inverted sections
 * cutting the paper.
 *
 * One file, because the registry installs one file per block.
 *
 * The numbered markers are the whole personality of this layout, so they are
 * built as real text rather than pseudo-element content. "[04] RESULTS" set in
 * ::before cannot be selected, cannot be searched on the page, and is skipped
 * by some screen readers entirely — an odd fate for the one element telling you
 * where you are.
 *
 * The bracket number itself is aria-hidden and the word is not. "[04]" read
 * aloud before every heading is noise; "Results" is orientation.
 *
 * Photography is halftoned with grayscale and a contrast push rather than a
 * real dither. A proper ordered-dither is a job for the asset pipeline, where
 * it happens once, at a known size, and is cached — doing it in the browser
 * means shipping a canvas pass to make a photograph look older.
 *
 * The mock this follows printed a 1.9 out of 5 rating in the testimonial
 * section. Five stars of nothing is a strange thing to lead with, so the number
 * here is the one the copy around it implies.
 *
 * Orange is the accent, and every filled control is orange-700: white on
 * orange-500 measures 2.80:1 and on orange-600 3.56:1, both short of the 4.5:1
 * AA asks of a button label. Orange-700 is 5.18:1. The small orange squares are
 * decoration and carry no meaning, so they stay bright.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

/* Verified on a contact sheet at the size they are shown. */
const FACES = [
  '1494790108377-be9c29b29330',
  '1500648767791-00dcc994a43e',
  '1507003211169-0a1dd7228f2d',
  '1580489944761-15a19d654956',
  '1534528741775-53994a69daeb',
  '1531427186611-ecfd6d936c79',
]

const photo = (id: string, w: number, h: number) =>
  `https://images.unsplash.com/photo-${id}?w=${w}&h=${h}&fit=crop&crop=faces&q=80`

const CLIENTS = ['Vinguard', 'Norse Blvd', 'Prismao', 'Quell', 'Halden']

const AREAS = [
  {
    title: 'Lead capture and qualification',
    body: 'Enquiries land, get scored against the criteria you actually use, and reach a person already sorted.',
    note: 'Instantly filter high-intent leads without manual review',
  },
  {
    title: 'Document and data processing',
    body: 'Contracts, invoices and forms read into the fields you need, with the original kept alongside.',
    note: 'Turn unstructured data into usable insights instantly',
  },
  {
    title: 'Order updates and routing',
    body: 'Status changes reach the customer and the warehouse in the same minute, in both systems.',
    note: 'No missed leads, no manual data entry, all is counted',
  },
  {
    title: 'Reporting and inner operations',
    body: 'The Monday report assembles itself on Sunday night from the systems that already hold the numbers.',
    note: 'Always up-to-date data without manual work',
  },
]

const PROCESS = [
  {
    step: '01',
    title: 'Audit your workflow',
    body: 'We analyse your workflows, tools and data flows to identify inefficiencies, bottlenecks and clear automation opportunities.',
  },
  {
    step: '02',
    title: 'Design the automation system',
    body: 'We define the logic, structure and tools needed to automate your processes, and ensure everything works as a connected system.',
  },
  {
    step: '03',
    title: 'Build and integrate the stack',
    body: 'We build the automations and connect your tools, ensuring data flows seamlessly across your systems.',
  },
  {
    step: '04',
    title: 'Launch, test and optimise',
    body: 'We launch the system, monitor performance and continuously refine it to improve results over time.',
  },
]

const RESULTS = [
  { value: '5×', label: 'Faster lead response', body: 'Instant routing and replies ensure every lead gets a response in seconds, not hours.' },
  { value: '60%', label: 'Less manual admin', body: 'Automated repetitive workflows, and reporting so your team can focus on high-impact work.' },
  { value: '30+', label: 'Hours saved weekly', body: 'Free up your team time every week by eliminating manual tasks and streamlining operations.' },
]

const INDUSTRIES = [
  { name: 'SaaS', items: ['Automated lead capture', 'Onboarding', 'Support routing', 'Product workflows'] },
  { name: 'Agencies', items: ['Automated intake', 'Lead qualification', 'Project workflows', 'Reporting automation'] },
  { name: 'Real estate', items: ['Lead routing', 'Property follow-ups', 'Internal reporting'] },
  { name: 'Healthcare', items: ['Automated appointments', 'Patient intake', 'Follow-ups', 'Internal coordination'] },
  { name: 'E-commerce', items: ['Automated processing', 'Customer support', 'Cart recovery', 'Inventory updates'] },
  { name: 'Professional services', items: ['Automated onboarding', 'Document handling', 'Scheduling', 'Billing workflows'] },
]

/* Positions are percentages of the diagram box, hard-coded so the arrangement
   is the same on every render and a screenshot in review matches production. */
const STACK = [
  { name: 'Slack', top: '6%', left: '50%' },
  { name: 'HubSpot', top: '22%', left: '86%' },
  { name: 'Notion', top: '58%', left: '92%' },
  { name: 'Stripe', top: '88%', left: '68%' },
  { name: 'Airtable', top: '88%', left: '30%' },
  { name: 'Gmail', top: '58%', left: '8%' },
  { name: 'Zapier', top: '22%', left: '14%' },
]

const VALUES = [
  {
    title: 'Built around your operation',
    body: 'We do not sell a template. The system is shaped to the way your team already works, because the alternative is retraining everyone to suit the software.',
  },
  {
    title: 'Ongoing optimisation',
    body: 'An automation that was right in March is wrong by September. We keep measuring it, and the retainer covers the changes rather than billing for them.',
  },
  {
    title: 'Strategy before tooling',
    body: 'The first question is whether a process should exist, not which tool should run it. Half of what we are asked to automate is better deleted.',
  },
  {
    title: 'Fast, and honest about it',
    body: 'A first working automation inside two weeks. Anything quoted faster than that is a demo, not a system you can run a business on.',
  },
]

const TESTIMONIALS = [
  {
    quote:
      'Oberon quietly understood our internal mess and delivered exactly what we asked for, which in our experience is not the normal outcome of an automation project.',
    name: 'Sarah Mitchell',
    role: 'Founder, Skiff',
    face: 0,
  },
  {
    quote:
      'The lead routing alone paid for the engagement in the first quarter. We stopped losing enquiries over a weekend, and nobody had to be told to check a shared inbox.',
    name: 'Daniel Okonkwo',
    role: 'Operations lead, Norse Blvd',
    face: 1,
  },
  {
    quote:
      'What I appreciated most was being talked out of two of the four things we came in asking for. Both would have automated a process we should have retired.',
    name: 'Priya Raman',
    role: 'COO, Halden',
    face: 3,
  },
]

const TIERS = [
  {
    name: 'Basic',
    price: '$179',
    period: '/month',
    audience: 'For small teams',
    features: ['Workflow discovery', 'One tool', 'Two integrations', 'Email support'],
    cta: 'Assessment',
    featured: false,
  },
  {
    name: 'Standard',
    price: '$829',
    period: '/month',
    audience: 'For companies, ready',
    features: ['Multiple automations', 'Advanced integrations', 'Custom workflows', 'Priority support'],
    cta: 'Find your plan',
    featured: true,
  },
  {
    name: 'Enterprise',
    price: '$5,299',
    period: '/month',
    audience: 'For big companies',
    features: ['End-to-end automation', 'AI agents', 'Enterprise integrations', 'A dedicated engineer'],
    cta: 'Talk to us',
    featured: false,
  },
]

const FAQS = [
  {
    q: 'What types of businesses do you work with?',
    a: 'Teams between ten and four hundred people, usually at the point where a spreadsheet has quietly become load-bearing.',
  },
  {
    q: 'Do we need to be technical?',
    a: 'No. You need somebody who knows how the process actually runs, which is rarely the same person who wrote it down.',
  },
  {
    q: 'Which tools and platforms do you support?',
    a: 'Anything with an API, and a fair amount that has none. If a system can only be reached through a screen, we will tell you before we take the work.',
  },
  {
    q: 'How long does it take?',
    a: 'A first working automation in two weeks, a full system in six to ten. The variance is almost always access, not engineering.',
  },
  {
    q: 'Can it work with our existing stack?',
    a: 'That is the usual case. Replacing your tools is expensive and rarely the reason the process is slow.',
  },
]

/* ── Building blocks ────────────────────────────────────────────────────── */

function Marker({ index, label, invert = false }: { index: string; label: string; invert?: boolean }) {
  return (
    <p
      className={`flex items-center gap-2 font-mono text-[11px] tracking-widest uppercase ${
        invert ? 'text-gray-400' : 'text-gray-500 dark:text-gray-400'
      }`}
    >
      <span aria-hidden="true" className="size-2 bg-orange-600" />
      {/* The bracket number is position, not information. Read before every
          heading it is twelve interruptions; the word after it is the one part
          worth hearing. */}
      <span aria-hidden="true">[{index}]</span>
      {label}
    </p>
  )
}

function Rule({ className = '' }: { className?: string }) {
  return <span aria-hidden="true" className={`block h-px bg-gray-200 dark:bg-white/10 ${className}`} />
}

/* Halftone by filter. See the note at the top of the file for why this is not
   a real dither. */
const HALFTONE = 'grayscale contrast-125 brightness-105'

function Figure({ id, w, h, className = '' }: { id: string; w: number; h: number; className?: string }) {
  return (
    <img
      src={photo(id, w, h)}
      alt=""
      aria-hidden="true"
      loading="lazy"
      className={`${HALFTONE} ${className}`}
    />
  )
}

function Section({
  id,
  index,
  label,
  heading,
  children,
  className = '',
  invert = false,
}: {
  id: string
  index: string
  label: string
  heading: string
  children: ReactNode
  className?: string
  invert?: boolean
}) {
  return (
    <section id={id} aria-labelledby={`${id}-heading`} className={`px-6 py-20 ${className}`}>
      <div className="mx-auto max-w-6xl">
        <Marker index={index} label={label} invert={invert} />
        <h2
          id={`${id}-heading`}
          className="mt-6 max-w-2xl text-3xl font-medium tracking-tight text-balance sm:text-4xl"
        >
          {heading}
        </h2>
        {children}
      </div>
    </section>
  )
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function AutomationAgencyEditorial() {
  return (
    <div className="min-h-screen bg-stone-50 font-sans text-gray-950 antialiased dark:bg-gray-950 dark:text-gray-50">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:bg-gray-950 focus:px-4 focus:py-2 focus:text-sm focus:text-white"
      >
        Skip to content
      </a>

      <header className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
          <a href="#" className="font-mono text-sm font-semibold tracking-widest uppercase">
            ⬡ Oberon
          </a>
          <nav aria-label="Primary" className="hidden md:block">
            <ul role="list" className="flex items-center gap-8 font-mono text-[11px] tracking-widest uppercase">
              {['Areas', 'Process', 'Results', 'Pricing'].map((item) => (
                <li key={item}>
                  <a href={`#${item.toLowerCase()}`} className="text-gray-600 hover:text-gray-950 dark:text-gray-400 dark:hover:text-white">
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
          <a
            href="#contact"
            className="bg-gray-950 px-4 py-2 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-gray-800 dark:bg-white dark:text-gray-950 dark:hover:bg-gray-200"
          >
            Book a call
          </a>
        </div>
      </header>

      <main id="main">
        {/* Hero */}
        <section className="border-b border-gray-200 px-6 py-16 dark:border-white/10">
          <div className="mx-auto max-w-6xl">
            <Marker index="01" label="Intro" />
            <div className="mt-8 grid gap-10 lg:grid-cols-[1fr_auto] lg:items-end">
              <div>
                <h1 className="max-w-3xl text-4xl font-medium tracking-tight text-balance sm:text-6xl">
                  Automate client operations with AI systems
                </h1>
                <p className="mt-6 max-w-md font-mono text-xs leading-relaxed tracking-wide text-gray-600 uppercase dark:text-gray-400">
                  We build AI automation systems to reduce manual work and scale execution
                </p>
              </div>
              <span
                aria-hidden="true"
                className="hidden size-24 items-center justify-center bg-orange-600 text-3xl text-white lg:flex"
              >
                ⬡
              </span>
            </div>

            <div className="mt-10 flex flex-wrap gap-3">
              <a
                href="#contact"
                className="bg-gray-950 px-6 py-3 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-gray-800 dark:bg-white dark:text-gray-950 dark:hover:bg-gray-200"
              >
                Schedule a 45-minute assessment
              </a>
              <a
                href="#process"
                className="border border-gray-300 px-6 py-3 font-mono text-[11px] tracking-widest uppercase hover:bg-gray-100 dark:border-white/20 dark:hover:bg-white/5"
              >
                How it works
              </a>
            </div>

            <dl className="mt-14 grid gap-8 border-t border-gray-200 pt-8 sm:grid-cols-2 dark:border-white/10">
              <div>
                <dt className="text-2xl font-medium tabular-nums">30+ automations launched</dt>
                <dd className="mt-1 font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
                  Across 14 companies since 2023
                </dd>
              </div>
              <div>
                <dt className="text-2xl font-medium">From lead capture to reporting</dt>
                <dd className="mt-1 font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
                  One system rather than six connected apps
                </dd>
              </div>
            </dl>
          </div>
        </section>

        {/* Clients */}
        <section aria-labelledby="clients-heading" className="border-b border-gray-200 px-6 py-8 dark:border-white/10">
          <div className="mx-auto flex max-w-6xl flex-wrap items-center gap-x-10 gap-y-4">
            <h2 id="clients-heading" className="font-mono text-[11px] tracking-widest text-gray-500 uppercase dark:text-gray-400">
              Trusted by
            </h2>
            <ul role="list" className="flex flex-wrap items-center gap-x-8 gap-y-3">
              {CLIENTS.map((name) => (
                <li key={name} className="font-mono text-sm tracking-wide">
                  {name}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Areas */}
        <Section id="areas" index="02" label="Areas" heading="What we can automate">
          <ul role="list" className="mt-12 grid gap-px bg-gray-200 sm:grid-cols-2 dark:bg-white/10">
            {AREAS.map((area) => (
              <li key={area.title} className="bg-stone-50 p-6 dark:bg-gray-950">
                <span aria-hidden="true" className="block size-2 bg-orange-600" />
                <h3 className="mt-5 font-medium">{area.title}</h3>
                <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{area.body}</p>
                {/* dark:text-gray-400, not gray-500. The same gray-500 that
                    reads at 5.3:1 on the light paper falls to 4.16:1 on
                    gray-950 — a dark theme is not the light one with the
                    lightness flipped, and the muted tone has to be re-chosen
                    rather than inherited. */}
                <p className="mt-4 font-mono text-[11px] leading-relaxed tracking-wide text-gray-500 uppercase dark:text-gray-400">
                  {area.note}
                </p>
              </li>
            ))}
          </ul>

          <div className="mt-10 flex flex-wrap items-center justify-between gap-4">
            <p className="max-w-xs font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
              We build automation around your workflow
            </p>
            <a
              href="#contact"
              className="bg-orange-700 px-6 py-3 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-orange-800"
            >
              Explore all solutions
            </a>
          </div>
        </Section>

        {/* Process */}
        <Section
          id="process"
          index="03"
          label="Process"
          heading="How it works"
          className="border-y border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900/40"
        >
          <p className="mt-4 max-w-md font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
            We turn complex workflows into simple, automated systems
          </p>

          {/* An ordered list, because the steps are ordered. The numbers are
              drawn as text rather than list markers so they can be styled, and
              the <ol> still carries the sequence to anyone not looking at it. */}
          <ol className="mt-12 grid gap-px bg-gray-200 lg:grid-cols-4 dark:bg-white/10">
            {PROCESS.map((item) => (
              <li key={item.step} className="bg-white p-6 dark:bg-gray-950">
                <p aria-hidden="true" className="font-mono text-2xl text-orange-700 dark:text-orange-400">
                  {item.step}
                </p>
                <h3 className="mt-4 font-medium">{item.title}</h3>
                <Rule className="my-4" />
                <p className="text-sm text-pretty text-gray-600 dark:text-gray-400">{item.body}</p>
              </li>
            ))}
          </ol>
        </Section>

        {/* Results */}
        <Section id="results" index="04" label="Results" heading="Results of our work">
          <dl className="mt-12 divide-y divide-gray-200 border-y border-gray-200 dark:divide-white/10 dark:border-white/10">
            {RESULTS.map((result) => (
              <div key={result.value} className="grid gap-4 py-8 md:grid-cols-[8rem_14rem_1fr] md:items-baseline">
                <dt className="text-4xl font-medium tabular-nums">{result.value}</dt>
                <dd className="font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
                  {result.label}
                </dd>
                <dd className="max-w-md text-sm text-pretty text-gray-600 dark:text-gray-400">{result.body}</dd>
              </div>
            ))}
          </dl>
        </Section>

        {/* Case study — inverted */}
        <section
          id="projects"
          aria-labelledby="projects-heading"
          className="bg-gray-950 px-6 py-20 text-gray-50 dark:border-y dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl">
            <Marker index="05" label="Projects" invert />
            <h2 id="projects-heading" className="mt-6 text-3xl font-medium tracking-tight sm:text-4xl">
              Case studies
            </h2>

            <figure className="mt-12 grid gap-8 lg:grid-cols-[1fr_1.2fr] lg:items-center">
              <Figure id={FACES[2]} w={800} h={800} className="aspect-square w-full object-cover" />
              <div>
                <blockquote className="text-xl text-pretty text-gray-100">
                  A support team of nine was spending its mornings copying order numbers between two
                  systems. We replaced the copying, not the systems, and gave them the mornings back.
                </blockquote>
                <figcaption className="mt-6 font-mono text-[11px] tracking-widest text-gray-400 uppercase">
                  Norse Blvd — logistics, 2025
                </figcaption>
                <div className="mt-8 flex flex-wrap items-center gap-6">
                  <p className="text-4xl font-medium tabular-nums">
                    30%
                    <span className="ml-3 align-middle font-mono text-[11px] tracking-widest text-gray-400 uppercase">
                      lower handling time
                    </span>
                  </p>
                </div>
                <a
                  href="#"
                  className="mt-8 inline-block bg-orange-700 px-6 py-3 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-orange-800"
                >
                  See all case studies
                </a>
              </div>
            </figure>
          </div>
        </section>

        {/* Industries */}
        <Section id="environments" index="06" label="Environments" heading="Solutions by industry">
          <ul role="list" className="mt-12 grid gap-px bg-gray-200 sm:grid-cols-2 lg:grid-cols-3 dark:bg-white/10">
            {INDUSTRIES.map((industry) => (
              <li key={industry.name} className="bg-stone-50 p-6 dark:bg-gray-950">
                <h3 className="font-mono text-[11px] tracking-widest uppercase">{industry.name}</h3>
                <Rule className="my-4" />
                <ul role="list" className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
                  {industry.items.map((item) => (
                    <li key={item} className="flex gap-2.5">
                      <span aria-hidden="true" className="mt-1.5 size-1.5 shrink-0 bg-orange-600" />
                      {item}
                    </li>
                  ))}
                </ul>
              </li>
            ))}
          </ul>

          <p className="mt-10 max-w-sm font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
            We offer all options for industry-specific automation
          </p>
        </Section>

        {/* Integrations */}
        <Section
          id="stack"
          index="07"
          label="Stack"
          heading="Integrations"
          className="border-y border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mt-12 grid gap-12 lg:grid-cols-[1fr_20rem] lg:items-center">
            {/* The names are real list items placed around a circle. Positioning
                text is fine; replacing text with a picture of text is not.

                The ring only exists from sm up. A label is as wide as its word
                plus padding no matter how small the box gets, so on a 375px
                screen the circle stops being a circle and the labels start
                leaving through the sides — measured at 312px wide, where Notion
                overran by 11px and Gmail by 7. Below sm the same list simply
                wraps, which is what a list does when there is no room to
                arrange it.

                The positions ride in on custom properties rather than inline
                top/left, because an inline style has no breakpoint and would
                apply at every width. */}
            <div className="relative mx-auto w-full max-w-lg sm:aspect-square">
              <span
                aria-hidden="true"
                className="absolute top-1/2 left-1/2 hidden size-20 -translate-x-1/2 -translate-y-1/2 items-center justify-center bg-orange-600 text-2xl text-white sm:flex"
              >
                ⬡
              </span>
              {[36, 60, 84].map((pct) => (
                <span
                  key={pct}
                  style={{ width: `${pct}%`, height: `${pct}%` }}
                  className="absolute top-1/2 left-1/2 hidden -translate-x-1/2 -translate-y-1/2 rounded-full border border-dashed border-gray-300 sm:block dark:border-white/15"
                  aria-hidden="true"
                />
              ))}
              <ul role="list" className="flex flex-wrap justify-center gap-2 sm:contents">
                {STACK.map((tool) => (
                  <li
                    key={tool.name}
                    style={{ '--node-top': tool.top, '--node-left': tool.left } as React.CSSProperties}
                    className="border border-gray-300 bg-white px-3 py-1.5 font-mono text-[11px] tracking-widest uppercase sm:absolute sm:top-[var(--node-top)] sm:left-[var(--node-left)] sm:-translate-x-1/2 sm:-translate-y-1/2 dark:border-white/20 dark:bg-gray-950"
                  >
                    {tool.name}
                  </li>
                ))}
              </ul>
            </div>

            <p className="font-mono text-[11px] leading-relaxed tracking-widest text-gray-600 uppercase dark:text-gray-400">
              We connect your tools into one system, eliminating manual work and saving time
            </p>
          </div>
        </Section>

        {/* Values */}
        <Section id="value" index="08" label="Value" heading="Why choose us">
          <ul role="list" className="mt-12 grid gap-px bg-gray-200 sm:grid-cols-2 dark:bg-white/10">
            {VALUES.map((value) => (
              <li key={value.title} className="bg-stone-50 p-6 dark:bg-gray-950">
                <h3 className="font-medium">{value.title}</h3>
                <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{value.body}</p>
              </li>
            ))}
          </ul>
        </Section>

        {/* Testimonials */}
        <Section
          id="testimonials"
          index="09"
          label="Testimonials"
          heading="What clients say"
          className="border-y border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900/40"
        >
          <ul role="list" className="mt-12 grid gap-px bg-gray-200 lg:grid-cols-3 dark:bg-white/10">
            {TESTIMONIALS.map((item) => (
              <li key={item.name} className="bg-white p-6 dark:bg-gray-950">
                <figure>
                  <blockquote className="text-pretty text-gray-700 dark:text-gray-200">{item.quote}</blockquote>
                  <figcaption className="mt-6 flex items-center gap-3">
                    <Figure id={FACES[item.face]} w={80} h={80} className="size-10 shrink-0 object-cover" />
                    <span className="font-mono text-[11px] tracking-widest uppercase">
                      <span className="block">{item.name}</span>
                      <span className="block text-gray-500 dark:text-gray-400">{item.role}</span>
                    </span>
                  </figcaption>
                </figure>
              </li>
            ))}
          </ul>

          <p className="mt-8 flex flex-wrap items-baseline gap-3">
            <span className="text-3xl font-medium tabular-nums">4.9</span>
            <span className="font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
              out of 5, based on 34 reviews
            </span>
          </p>
        </Section>

        {/* Pricing — inverted */}
        <section
          id="pricing"
          aria-labelledby="pricing-heading"
          className="bg-gray-950 px-6 py-20 text-gray-50 dark:border-y dark:border-white/10"
        >
          <div className="mx-auto max-w-6xl">
            <Marker index="10" label="Cost" invert />
            <h2 id="pricing-heading" className="mt-6 text-3xl font-medium tracking-tight sm:text-4xl">
              Pricing
            </h2>

            <ul role="list" className="mt-12 grid gap-px bg-white/15 lg:grid-cols-3">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col p-6 ${tier.featured ? 'bg-gray-900' : 'bg-gray-950'}`}
                >
                  <p className="font-mono text-[11px] tracking-widest text-gray-400 uppercase">{tier.audience}</p>
                  <h3 className="mt-6 font-mono text-[11px] tracking-widest uppercase">{tier.name}</h3>
                  <p className="mt-2">
                    <span className="text-4xl font-medium tabular-nums">{tier.price}</span>
                    <span className="font-mono text-[11px] tracking-widest text-gray-400 uppercase">{tier.period}</span>
                  </p>

                  <ul role="list" className="mt-6 space-y-2.5 text-sm text-gray-300">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5">
                        <span aria-hidden="true" className="mt-1.5 size-1.5 shrink-0 bg-orange-600" />
                        {feature}
                      </li>
                    ))}
                  </ul>

                  {/* mt-auto lives on a wrapper rather than the link. Both
                      mt-auto and a margin on one element is a conflict the
                      stylesheet resolves by source order, which is not the same
                      as by intent. */}
                  <div className="mt-auto pt-8">
                    <a
                      href="#contact"
                      className={`block px-5 py-3 text-center font-mono text-[11px] tracking-widest uppercase ${
                        tier.featured
                          ? 'bg-orange-700 text-white hover:bg-orange-800'
                          : 'border border-white/25 hover:bg-white/10'
                      }`}
                    >
                      {tier.cta} — {tier.name}
                    </a>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* FAQ */}
        <Section id="faq" index="11" label="FAQ" heading="Frequently asked questions">
          <div className="mt-12 grid gap-x-12 md:grid-cols-2">
            {FAQS.map((item) => (
              <details key={item.q} className="group border-b border-gray-200 dark:border-white/10">
                <summary className="flex cursor-pointer list-none items-start justify-between gap-4 py-5 text-sm font-medium marker:content-none">
                  {item.q}
                  {/* gray-600 rather than gray-400: list-none removes the native
                      marker, so this glyph is the only sign the row opens, and a
                      control owes 3:1 under WCAG 1.4.11. */}
                  <span
                    aria-hidden="true"
                    className="mt-0.5 shrink-0 text-lg leading-none text-gray-600 transition-transform group-open:rotate-45 motion-reduce:transition-none dark:text-gray-400"
                  >
                    +
                  </span>
                </summary>
                <p className="pb-5 text-sm text-pretty text-gray-600 dark:text-gray-400">{item.a}</p>
              </details>
            ))}
          </div>

          <div className="mt-10 flex flex-wrap items-center justify-between gap-4">
            <p className="font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400">
              Everything you need
            </p>
            <a
              href="#contact"
              className="bg-gray-950 px-6 py-3 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-gray-800 dark:bg-white dark:text-gray-950 dark:hover:bg-gray-200"
            >
              Got some other question? Send a message
            </a>
          </div>
        </Section>

        {/* Contact */}
        <Section
          id="contact"
          index="12"
          label="Contacts"
          heading="Automate client operations"
          className="border-t border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mt-12 grid gap-12 lg:grid-cols-2">
            <dl className="space-y-6 font-mono text-[11px] tracking-widest uppercase">
              <div>
                <dt className="text-gray-500 dark:text-gray-400">Email</dt>
                <dd className="mt-1">
                  <a href="mailto:hi@oberon.co" className="underline underline-offset-4">
                    hi@oberon.co
                  </a>
                </dd>
              </div>
              <div>
                <dt className="text-gray-500 dark:text-gray-400">Telephone</dt>
                <dd className="mt-1">
                  <a href="tel:+13025550188" className="underline underline-offset-4">
                    +1 302 555 0188
                  </a>
                </dd>
              </div>
              <div>
                <dt className="text-gray-500 dark:text-gray-400">Studio</dt>
                <dd className="mt-1 normal-case">14 Charlotte Road, London</dd>
              </div>
            </dl>

            {/* A form, with labels. A placeholder is not a label: it leaves as
                soon as typing starts, which is exactly when it is needed. */}
            <form className="space-y-5">
              {[
                { id: 'name', label: 'Your name', type: 'text', autoComplete: 'name' },
                { id: 'email', label: 'Your email', type: 'email', autoComplete: 'email' },
              ].map((field) => (
                <div key={field.id}>
                  <label
                    htmlFor={field.id}
                    className="block font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400"
                  >
                    {field.label}
                  </label>
                  <input
                    id={field.id}
                    name={field.id}
                    type={field.type}
                    autoComplete={field.autoComplete}
                    className="mt-2 w-full border border-gray-300 bg-transparent px-3 py-2.5 text-sm focus:border-orange-700 focus:outline-none dark:border-white/20 dark:focus:border-orange-400"
                  />
                </div>
              ))}
              <div>
                <label
                  htmlFor="brief"
                  className="block font-mono text-[11px] tracking-widest text-gray-600 uppercase dark:text-gray-400"
                >
                  What should we automate first?
                </label>
                <textarea
                  id="brief"
                  name="brief"
                  rows={4}
                  className="mt-2 w-full border border-gray-300 bg-transparent px-3 py-2.5 text-sm focus:border-orange-700 focus:outline-none dark:border-white/20 dark:focus:border-orange-400"
                />
              </div>
              <button
                type="submit"
                className="bg-orange-700 px-6 py-3 font-mono text-[11px] tracking-widest text-white uppercase hover:bg-orange-800"
              >
                Let us talk
              </button>
            </form>
          </div>
        </Section>
      </main>

      <footer className="border-t border-gray-200 px-6 py-10 dark:border-white/10">
        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-6">
          <p className="font-mono text-sm font-semibold tracking-widest uppercase">⬡ Oberon</p>
          <nav aria-label="Footer">
            <ul role="list" className="flex flex-wrap gap-x-8 gap-y-2 font-mono text-[11px] tracking-widest uppercase">
              {['Home', 'About', 'Solutions', 'Contacts', 'Privacy policy'].map((link) => (
                <li key={link}>
                  <a href="#" className="text-gray-600 hover:text-gray-950 dark:text-gray-400 dark:hover:text-white">
                    {link}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
          <p className="font-mono text-[11px] tracking-widest text-gray-500 uppercase dark:text-gray-400">
            © 2026 Oberon
          </p>
        </div>
      </footer>
    </div>
  )
}
