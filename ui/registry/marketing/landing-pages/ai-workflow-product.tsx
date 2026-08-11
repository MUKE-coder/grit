'use client'

import { useId, useState } from 'react'
import type { ReactNode } from 'react'

/*
 * A product page for an automation tool: hero, integrations, feature bento,
 * industries, pricing with a billing toggle, compliance and an FAQ.
 *
 * One file, because the registry installs one file per block.
 *
 * The billing period is a radio group in a fieldset, not a switch. This is the
 * thing pages like this almost always get wrong. A switch means on or off, and
 * announces as such: "yearly billing, switch, off" tells you nothing about what
 * is on instead. Monthly and yearly are two named options and exactly one is
 * chosen, which is what radios are for — arrow keys move between them, the
 * legend is announced before them, and the chosen one is announced as checked.
 *
 * Changing the period rewrites every price on the page without moving focus,
 * so the change goes through a status region. Otherwise the control reports
 * "yearly, selected" and stays silent about the four numbers that just moved.
 *
 * The inputs are sr-only rather than hidden. `hidden` removes them from the tab
 * order and the accessibility tree, taking the keyboard support with them; the
 * visible pill is the <label>, so clicking works with no handler.
 *
 * Compliance badges carry their names in text. A row of shield glyphs is a row
 * of pictures, and "we are certified" is exactly the claim that has to survive
 * not being seen.
 *
 * The product surfaces are markup and aria-hidden; the industry cards use
 * photographs. That split is deliberate: a stock photo cannot be a screenshot
 * of your product, so presenting one as if it were is a claim the page has not
 * earned, whereas a warehouse really is a warehouse.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

const NAV = ['Pricing', 'About', 'Careers', 'Blog']

const LOGOS = ['Bridge', 'Arch', 'Bill', 'Paperwork', 'Foldera', 'Granola', 'Modal', 'Replicate']

const STEPS = [
  {
    title: 'Design your workflow',
    body: 'Drag and drop nodes to order and configure agents into logical workflows.',
  },
  {
    title: 'Connect your tools',
    body: 'Agents operate independently and coordinate tasks to complete complex goals together.',
  },
  {
    title: 'Deploy and scale',
    body: 'Push to production in a sandbox that runs the same code path, including the failures.',
  },
]

const FEATURES = [
  { title: 'Model selector', body: 'Track real-time activity of agents with a powerful visual interface designed for technical teams.', mock: 'models', wide: false },
  { title: 'Text to workflow builder', body: 'Describe what you want in a sentence and get a first draft of the graph to edit.', mock: 'chat', wide: false },
  { title: 'Native tool integration', body: 'Meeting summaries, code review and customer support wired in without a connector to write.', mock: 'tools', wide: true },
]

const SMALL_FEATURES = [
  { title: 'One-click auth', body: 'A drag-and-drop interface to create, connect and configure agents into logical workflows.' },
  { title: 'Realtime sync', body: 'Agents operate independently and coordinate tasks to complete complex goals.' },
  { title: 'Custom connector SDK', body: 'Run agent workflows in a sandbox to preview behaviour, debug logic and time interventions.' },
]

/* Photographs, not markup. These are scenes rather than product surfaces, so a
   real image says more than a grey rectangle. The dashboard mocks further up
   stay as markup on purpose: a stock photo cannot be a screenshot of your
   product, and using one as if it were is a claim the page has not earned. */
const shot = (id: string) =>
  `https://images.unsplash.com/photo-${id}?w=600&h=400&fit=crop&q=80`

const INDUSTRIES = [
  { name: 'DevOps', body: 'Ship, roll back and triage without a human copying values between dashboards.', image: shot('1558494949-ef010cbdcc31') },
  { name: 'SalesOps', body: 'Enrich, route and follow up on leads from the systems you already pay for.', image: shot('1600880292203-757bb62b4baf') },
  { name: 'Supply chain', body: 'Reconcile stock across warehouses and flag the discrepancy before the customer does.', image: shot('1553413077-190dd305871c') },
  { name: 'Customer support', body: 'Draft, categorise and escalate, with the model kept out of the send decision.', image: shot('1552664730-d307ca884978') },
  { name: 'DataOps', body: 'Pipelines that retry sensibly and tell you which upstream actually broke.', image: shot('1551288049-bebda4e38f71') },
  { name: 'FinOps', body: 'Spot the spend anomaly on the day it happens rather than at the month-end close.', image: shot('1554224155-6726b3ff858f') },
]

/* Prices in integer cents, per seat. The yearly figure is the effective
   monthly rate when billed annually, so the two are comparable at a glance
   rather than one being an annual total the reader has to divide. */
const TIERS = [
  {
    name: 'Growth',
    blurb: 'Everything to start.',
    monthly: 0,
    yearly: 0,
    cta: 'Start building',
    featured: false,
    features: ['Up to 5 active agents', '50 simulation runs', 'Visual builder access', 'GitHub and Zapier integration', 'Basic support', 'Community Slack access'],
  },
  {
    name: 'Scale',
    blurb: 'For growing teams.',
    monthly: 1500,
    yearly: 1200,
    cta: 'Start for free',
    featured: true,
    features: ['Up to 50 active agents', '500 simulation runs', 'Visual builder access', 'GitHub and Zapier integration', 'Priority support', 'Priority Slack access'],
  },
  {
    name: 'Enterprise',
    blurb: 'For the whole company.',
    monthly: 3200,
    yearly: 2500,
    cta: 'Contact sales',
    featured: false,
    features: ['Unlimited active agents', 'Unlimited simulation runs', 'Visual builder access', 'GitHub and Zapier integration', 'Dedicated support', 'Access to Flight Club'],
  },
]

const COMPLIANCE = [
  { name: 'SOC 2 Type II', body: 'Audited annually' },
  { name: 'GDPR', body: 'EU data residency' },
  { name: 'ISO 27001', body: 'Certified' },
]

const FAQS = [
  { q: 'What exactly does this platform do?', a: 'It runs and monitors agent workflows: you describe the steps, connect the tools, and it executes them with a record of what happened.' },
  { q: 'How do I get started with creating my first workflow?', a: 'Start from a template or describe the workflow in a sentence and edit the draft it produces. Nothing needs deploying to try it.' },
  { q: 'What tools and services can I integrate?', a: 'The common developer and sales tools out of the box, plus anything with an HTTP API through the connector SDK.' },
  { q: 'Is my data secure when using agents?', a: 'Data stays in the region you choose, access is role-based and audited, and model providers are configurable per workspace.' },
  { q: 'Can I test workflows before they go live?', a: 'Yes. Runs execute in a sandbox against the same code path, so the failures you need to handle actually appear.' },
  { q: 'What is the difference between automated and manual steps?', a: 'A manual step pauses the run and waits for a person. Use it wherever a mistake would be expensive to undo.' },
]

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', minimumFractionDigits: 0 })

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

function AppMock() {
  return (
    <Panel>
      <div className="flex items-center gap-1.5 border-b border-gray-200 px-3 py-2.5 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="flex">
        <div className="hidden w-36 shrink-0 space-y-2 border-r border-gray-200 p-3 sm:block dark:border-white/10">
          {Array.from({ length: 7 }, (_, i) => (
            <span
              key={i}
              className={`block h-2 rounded ${i === 0 ? 'w-4/5 bg-rose-300 dark:bg-rose-500/40' : 'w-full bg-gray-200 dark:bg-white/10'}`}
            />
          ))}
        </div>
        <div className="flex-1 p-4">
          <div className="grid grid-cols-4 gap-2">
            {Array.from({ length: 8 }, (_, i) => (
              <span key={i} className="block h-8 rounded bg-gray-100 dark:bg-white/5" />
            ))}
          </div>
          <div className="mt-4 flex items-end gap-1.5">
            {[40, 62, 48, 78, 56, 88, 70, 94].map((h, i) => (
              <span
                key={i}
                style={{ height: `${h * 0.5}px` }}
                className="flex-1 rounded-t bg-gradient-to-t from-rose-500 to-rose-300"
              />
            ))}
          </div>
        </div>
      </div>
    </Panel>
  )
}

function ModelsMock() {
  return (
    <Panel className="p-4">
      {['Claude', 'GPT', 'Llama'].map((model, i) => (
        <div key={model} className="flex items-center gap-2 py-1.5 text-[11px]">
          <span className={`size-2 rounded-full ${i === 0 ? 'bg-emerald-500' : 'bg-gray-300 dark:bg-white/20'}`} />
          <span className="text-gray-700 dark:text-gray-200">{model}</span>
        </div>
      ))}
    </Panel>
  )
}

function ChatMock() {
  return (
    <Panel className="space-y-2 p-4">
      <span className="ml-auto block h-8 w-3/5 rounded-lg bg-sky-500" />
      <span className="block h-6 w-2/3 rounded-lg bg-gray-100 dark:bg-white/10" />
    </Panel>
  )
}

function ToolsMock() {
  return (
    <Panel className="p-4">
      <div className="space-y-2">
        {['Meeting summariser', 'Code reviewer', 'Customer support'].map((tool) => (
          <div
            key={tool}
            className="flex items-center justify-between rounded-lg border border-gray-200 px-3 py-2 text-[11px] text-gray-700 dark:border-white/10 dark:text-gray-200"
          >
            {tool}
            <span className="rounded bg-emerald-100 px-1.5 py-0.5 text-[10px] text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300">
              Connected
            </span>
          </div>
        ))}
      </div>
    </Panel>
  )
}

const MOCKS: Record<string, ReactNode> = {
  models: <ModelsMock />,
  chat: <ChatMock />,
  tools: <ToolsMock />,
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function LandingPageAiWorkflowProduct({
  brand = 'Notus',
  title = 'Manage and simulate agentic workflows',
  subtitle = 'Empower developers and technical teams to model, simulate and manage AI-driven workflows visually.',
}: {
  brand?: string
  title?: string
  subtitle?: string
}) {
  const [yearly, setYearly] = useState(false)
  const [announcement, setAnnouncement] = useState('')
  const billingId = useId()
  const newsletterId = useId()

  function changePeriod(next: boolean) {
    if (next === yearly) return
    setYearly(next)
    setAnnouncement(
      next
        ? 'Showing yearly prices, billed annually. Two months free.'
        : 'Showing monthly prices, billed monthly.',
    )
  }

  return (
    <div className="bg-white dark:bg-gray-950">
      <header className="border-b border-gray-200 dark:border-white/10">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-6xl items-center justify-between gap-6 px-6 py-4"
        >
          <a
            href="#"
            className="text-lg font-semibold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:text-white"
          >
            {brand}
          </a>
          <ul role="list" className="hidden gap-7 md:flex">
            {NAV.map((item) => (
              <li key={item}>
                <a
                  href="#"
                  className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:text-gray-300 dark:hover:text-white"
                >
                  {item}
                </a>
              </li>
            ))}
          </ul>
          <a
            href="#pricing"
            className="inline-flex min-h-11 shrink-0 items-center rounded-lg bg-gray-900 px-4 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:bg-white dark:text-gray-900"
          >
            Start building
          </a>
        </nav>
      </header>

      <main>
        {/* Hero */}
        <section className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-2xl text-center">
              <h1 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
                Manage and simulate{' '}
                <span className="text-rose-600 dark:text-rose-400">agentic workflows</span>
              </h1>
              <p className="mx-auto mt-6 text-lg text-pretty text-gray-600 dark:text-gray-300">
                {subtitle}
              </p>
              <div className="mt-8 flex flex-wrap justify-center gap-3">
                <a
                  href="#pricing"
                  className="inline-flex min-h-12 items-center rounded-lg bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:bg-white dark:text-gray-900"
                >
                  Start building
                </a>
                <a
                  href="#pricing"
                  className="inline-flex min-h-12 items-center rounded-lg border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
                >
                  View pricing
                </a>
              </div>
            </div>

            <div className="mt-14">
              <AppMock />
            </div>
          </div>
        </section>

        {/* Logos */}
        <section aria-labelledby="logos" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-12">
            <h2 id="logos" className="sr-only">
              Teams building with {brand}
            </h2>
            <ul
              role="list"
              className="grid grid-cols-2 items-center justify-items-center gap-6 sm:grid-cols-4"
            >
              {LOGOS.map((logo) => (
                <li key={logo} className="text-sm font-semibold tracking-tight text-gray-600 dark:text-gray-400">
                  {logo}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Integrations */}
        <section aria-labelledby="integrates" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-5xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="integrates"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Integrates easily
              </h2>
              <p className="mt-3 text-gray-600 dark:text-gray-300">
                Three steps from an empty canvas to a workflow doing real work.
              </p>
            </div>

            {/* Ordered, because the steps happen in order. */}
            <ol className="mt-12 grid gap-6 sm:grid-cols-3">
              {STEPS.map((step, index) => (
                <li
                  key={step.title}
                  className="rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40"
                >
                  <span
                    aria-hidden="true"
                    className="flex size-8 items-center justify-center rounded-lg bg-rose-600 text-sm font-bold text-white"
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
        <section aria-labelledby="features" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="features"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Built for agentic intelligence
              </h2>
            </div>

            {/* Two columns: two narrow tiles, then one spanning both. Spans
                tile exactly at 1 + 1 + 2 = 4 cells over two rows. */}
            <ul role="list" className="mt-12 grid gap-6 md:grid-cols-2">
              {FEATURES.map((feature) => (
                <li
                  key={feature.title}
                  className={`flex flex-col rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/40 ${
                    feature.wide ? 'md:col-span-2' : ''
                  }`}
                >
                  <h3 className="font-semibold text-gray-900 dark:text-white">{feature.title}</h3>
                  <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{feature.body}</p>
                  <div className="mt-auto pt-6">{MOCKS[feature.mock]}</div>
                </li>
              ))}
            </ul>

            <ul role="list" className="mt-10 grid gap-8 sm:grid-cols-3">
              {SMALL_FEATURES.map((feature) => (
                <li key={feature.title}>
                  <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    {feature.title}
                  </h3>
                  <p className="mt-1.5 text-sm text-gray-600 dark:text-gray-300">{feature.body}</p>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Industries */}
        <section
          aria-labelledby="industries"
          className="border-b border-gray-200 bg-gray-50 dark:border-white/10 dark:bg-gray-900/40"
        >
          <div className="mx-auto max-w-6xl px-6 py-20">
            <div className="mx-auto max-w-xl text-center">
              <h2
                id="industries"
                className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Across various industries
              </h2>
            </div>
            <ul role="list" className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
              {INDUSTRIES.map((industry) => (
                <li
                  key={industry.name}
                  className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900"
                >
                  {/* Empty alt: the heading directly below names the industry,
                      and describing the photograph as well would say it
                      twice. */}
                  <img
                    src={industry.image}
                    alt=""
                    loading="lazy"
                    className="h-36 w-full object-cover"
                  />
                  <div className="p-6">
                    <h3 className="font-semibold text-rose-700 dark:text-rose-400">
                      {industry.name}
                    </h3>
                    <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{industry.body}</p>
                  </div>
                </li>
              ))}
            </ul>
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
                Simple and feasible pricing
              </h2>
            </div>

            {/* A radio group, not a switch. A switch means on or off and
                announces as such; this is a choice between two named options
                where exactly one applies. */}
            <fieldset className="mt-8 flex justify-center">
              <legend className="sr-only">Billing period</legend>
              <div className="inline-flex rounded-lg border border-gray-300 p-1 dark:border-white/15">
                {[
                  { value: false, label: 'Monthly' },
                  { value: true, label: 'Yearly' },
                ].map((option) => {
                  const id = `${billingId}-${option.label}`
                  const checked = option.value === yearly
                  return (
                    <div key={option.label}>
                      {/* sr-only, not hidden: hidden removes it from the tab
                          order and the a11y tree, losing arrow-key support. */}
                      <input
                        type="radio"
                        id={id}
                        name={billingId}
                        checked={checked}
                        onChange={() => changePeriod(option.value)}
                        className="peer sr-only"
                      />
                      <label
                        htmlFor={id}
                        className={`flex min-h-11 cursor-pointer items-center gap-2 rounded-md px-4 text-sm font-medium peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-rose-600 ${
                          checked
                            ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'
                            : 'text-gray-700 dark:text-gray-300'
                        }`}
                      >
                        {option.label}
                        {option.value && (
                          <span
                            className={`rounded-full px-2 py-0.5 text-xs ${
                              checked
                                ? 'bg-white/20 text-white dark:bg-gray-900/10 dark:text-gray-900'
                                : 'bg-emerald-100 text-emerald-800 dark:bg-emerald-500/15 dark:text-emerald-300'
                            }`}
                          >
                            2 months free
                          </span>
                        )}
                      </label>
                    </div>
                  )
                })}
              </div>
            </fieldset>

            {/* Every price on the page just changed and focus did not move. */}
            <p role="status" aria-live="polite" className="sr-only">
              {announcement}
            </p>

            <ul role="list" className="mt-10 grid items-start gap-6 lg:grid-cols-3">
              {TIERS.map((tier) => {
                const price = yearly ? tier.yearly : tier.monthly
                return (
                  <li
                    key={tier.name}
                    className={`flex h-full flex-col rounded-2xl border p-8 ${
                      tier.featured
                        ? 'border-rose-600 bg-rose-50/50 dark:border-rose-500/40 dark:bg-rose-500/10'
                        : 'border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900'
                    }`}
                  >
                    <div className="flex items-center justify-between gap-3">
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                        {tier.name}
                      </h3>
                      {tier.featured && (
                        /* In words: colour reports nothing. */
                        <span className="rounded-full bg-rose-600 px-2.5 py-1 text-xs font-medium text-white">
                          Most popular
                        </span>
                      )}
                    </div>
                    <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">{tier.blurb}</p>

                    <p className="mt-6">
                      <span className="text-4xl font-bold tracking-tight tabular-nums text-gray-900 dark:text-white">
                        {money.format(price / 100)}
                      </span>{' '}
                      <span className="text-sm text-gray-500 dark:text-gray-400">
                        per seat, per month
                      </span>
                    </p>
                    {/* Says what the number means under this period, rather
                        than leaving the toggle to imply it. */}
                    <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                      {yearly ? 'Billed annually' : 'Billed monthly'}
                    </p>

                    <ul role="list" className="mt-6 space-y-2.5">
                      {tier.features.map((feature) => (
                        <li key={feature} className="flex gap-2.5 text-sm">
                          <span aria-hidden="true" className="text-rose-600 dark:text-rose-400">
                            ✓
                          </span>
                          <span className="text-gray-700 dark:text-gray-200">{feature}</span>
                        </li>
                      ))}
                    </ul>

                    {/* Wrapper owns the spacing: mt-auto and mt-8 on one
                        element is a conflict resolved by stylesheet order. */}
                    <div className="mt-auto pt-8">
                      <a
                        href="#"
                        className={`inline-flex min-h-12 w-full items-center justify-center rounded-lg text-sm font-medium focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 ${
                          tier.featured
                            ? 'bg-rose-600 text-white hover:bg-rose-700'
                            : 'border border-gray-300 text-gray-900 hover:bg-gray-50 dark:border-white/15 dark:text-white dark:hover:bg-white/5'
                        }`}
                      >
                        {tier.cta}
                        <span className="sr-only"> with the {tier.name} plan</span>
                      </a>
                    </div>
                  </li>
                )
              })}
            </ul>
          </div>
        </section>

        {/* Compliance */}
        <section aria-labelledby="compliance" className="border-b border-gray-200 dark:border-white/10">
          <div className="mx-auto max-w-5xl px-6 py-16">
            <div className="rounded-2xl border border-gray-200 bg-gray-50 p-8 dark:border-white/10 dark:bg-gray-900/40">
              <h2
                id="compliance"
                className="text-2xl font-bold tracking-tight text-balance text-gray-900 dark:text-white"
              >
                Scale securely with confidence
              </h2>
              <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
                Built to satisfy the enterprise-grade security questionnaire before it arrives.
              </p>

              {/* Names in text. A row of shield glyphs is a row of pictures,
                  and this is exactly the claim that has to survive not being
                  seen. */}
              <dl className="mt-6 grid gap-4 sm:grid-cols-3">
                {COMPLIANCE.map((item) => (
                  <div
                    key={item.name}
                    className="rounded-xl border border-gray-200 bg-white p-4 dark:border-white/10 dark:bg-gray-900"
                  >
                    <dt className="text-sm font-semibold text-gray-900 dark:text-white">
                      {item.name}
                    </dt>
                    <dd className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">{item.body}</dd>
                  </div>
                ))}
              </dl>
            </div>
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
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 text-sm font-medium text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:text-white">
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
        <section aria-labelledby="cta" className="bg-gray-50 dark:bg-gray-900/40">
          <div className="mx-auto max-w-3xl px-6 py-20 text-center">
            <h2
              id="cta"
              className="text-3xl font-bold tracking-tight text-balance text-gray-900 sm:text-4xl dark:text-white"
            >
              Connect your current stack and start automating
            </h2>
            <a
              href="#pricing"
              className="mt-8 inline-flex min-h-12 items-center rounded-lg bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:bg-white dark:text-gray-900"
            >
              Start building for free
            </a>
          </div>
        </section>
      </main>

      <footer className="border-t border-gray-200 dark:border-white/10">
        <div className="mx-auto grid max-w-6xl gap-8 px-6 py-12 sm:grid-cols-2">
          <div>
            <p className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white">
              {brand}
            </p>
            <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
              Manage and simulate agentic workflows.
            </p>
          </div>

          <form onSubmit={(event) => event.preventDefault()} className="sm:justify-self-end">
            {/* A real label, not a placeholder. A placeholder disappears the
                moment someone types and is not reliably announced. */}
            <label
              htmlFor={newsletterId}
              className="block text-sm font-medium text-gray-900 dark:text-white"
            >
              Get product updates
            </label>
            <div className="mt-2 flex gap-2">
              <input
                id={newsletterId}
                type="email"
                autoComplete="email"
                placeholder="you@example.com"
                className="min-h-11 w-full min-w-0 rounded-lg border border-gray-300 px-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-rose-600 sm:w-64 dark:border-white/15 dark:bg-transparent dark:text-white"
              />
              <button
                type="submit"
                className="inline-flex min-h-11 shrink-0 items-center rounded-lg bg-gray-900 px-4 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:bg-white dark:text-gray-900"
              >
                Subscribe
              </button>
            </div>
          </form>
        </div>

        <div className="border-t border-gray-200 dark:border-white/10">
          <nav aria-label="Footer" className="mx-auto max-w-6xl px-6 py-6">
            <ul role="list" className="flex flex-wrap gap-x-6 gap-y-2">
              {[...NAV, 'Privacy', 'Terms', 'Security'].map((item) => (
                <li key={item}>
                  <a
                    href="#"
                    className="inline-flex min-h-11 items-center text-sm text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-600 dark:text-gray-400 dark:hover:text-white"
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
