/*
 * A task-management SaaS page: hero, logo wall, a three-step setup, a feature
 * browser, integrations, pricing, testimonials and an FAQ.
 *
 * One file, because the registry installs one file per block.
 *
 * Two controls in the design this follows are interactive, and this block is
 * static markup with no state. That is a fork worth being explicit about,
 * because the usual resolution — draw the control, wire nothing — ships a
 * button that lies.
 *
 *   The billing switch is a fieldset of two radios rather than a switch. A
 *   switch announces "on" and "off", but nobody is turning annual billing on;
 *   they are choosing between two named periods, and the price on screen has to
 *   agree with the one that is chosen.
 *
 *   The feature browser is a list and a panel, not ARIA tabs. Tabs are a
 *   keyboard contract — arrow keys move between them, Home and End jump to the
 *   ends, only the selected tab is tabbable — and role="tab" with none of that
 *   behind it is worse than the plain list, because it promises the contract to
 *   exactly the people relying on it.
 *
 * The carousel dots under the testimonials are gone. All three quotes are on
 * screen at every width, so the dots controlled nothing and counted to three
 * for no reason.
 *
 * One <h1>, an <h2> per section, <h3> inside.
 */

import type { ReactNode } from 'react'

const FACES = ['1494790108377-be9c29b29330', '1500648767791-00dcc994a43e', '1507003211169-0a1dd7228f2d']
const face = (id: string, size = 40) =>
  `https://images.unsplash.com/photo-${id}?w=${size * 2}&h=${size * 2}&fit=crop&crop=faces&q=75`

const NAV = ['Features', 'Solution', 'Integration', 'Pricing']

const TRUST = ['No credit card required', '14-day free trial', 'Cancel any time']

const BRANDS = [
  'Quixotic', 'Goodwill', 'Watchtower', 'Foresight', 'Constellation',
  'Codecraft', 'Northwind', 'Portals', 'Kestrel', 'Visionwork',
]

const STEPS = [
  {
    title: 'Create an account',
    body: 'Push through it in seconds. Set up your workspace and you are ready to go.',
    mock: 'signup',
  },
  {
    title: 'Invite the team',
    body: 'Bring your teammates on board instantly. Manage roles, track projects, and stay aligned without the chaos.',
    mock: 'invite',
  },
  {
    title: 'Assign and track',
    body: 'Delegate tasks, set due dates, and track progress without losing the thread.',
    mock: 'table',
    wide: true,
  },
]

const CAPABILITIES = [
  { name: 'Dashboard', body: 'See your work clearly — real-time progress, blockers and deadlines in one customisable, distraction-free dashboard that adapts to your role.' },
  { name: 'Global search', body: 'One field across every project, task, comment and file, with results ranked by what you touched most recently.' },
  { name: 'Multiple workspaces', body: 'Keep client work, internal projects and personal tasks apart without a second account.' },
  { name: 'Messages', body: 'Discussion attached to the task it is about, so the decision and the work do not live in different apps.' },
  { name: 'Multiple task views', body: 'Board, list, calendar and timeline over the same data. Changing the view never changes the truth.' },
  { name: 'Private notepad', body: 'A scratch space that is yours, until you decide a note should become a task.' },
  { name: 'Two-factor security', body: 'TOTP and passkeys, enforced per workspace, with recovery codes you can print.' },
]

const INTEGRATIONS = ['Slack', 'GitHub', 'Figma', 'Drive', 'Zoom', 'Notion', 'Linear', 'Calendar']

const TIERS = [
  {
    name: 'Starter',
    price: 'Free',
    period: null as string | null,
    blurb: 'Perfect for individuals and small teams getting started.',
    cta: 'Get started',
    featured: false,
    inherits: null as string | null,
    features: ['Unlimited projects', 'Task creation and due dates', 'Team collaboration, up to 3 members', 'Basic analytics'],
  },
  {
    name: 'Professional',
    price: '$24',
    period: 'per member, billed yearly',
    blurb: 'For growing teams who need advanced features.',
    cta: 'Start Professional',
    featured: true,
    inherits: 'Starter',
    badge: '25% off',
    features: ['Unlimited members', 'Automations and recurring tasks', 'Advanced integrations', 'Priority support'],
  },
  {
    name: 'Enterprise',
    price: 'Custom',
    period: null,
    blurb: 'Tailored for large organisations with complex workflows.',
    cta: 'Talk to sales',
    featured: false,
    inherits: 'Professional',
    features: ['Dedicated account manager', 'Advanced reporting and security', 'Custom integrations and onboarding', 'A signed DPA'],
  },
]

const TESTIMONIALS = [
  {
    quote:
      'Tasker has completely streamlined how our design team collaborates. It keeps everyone on track without the meetings overload.',
    name: 'Maya Ortiz',
    role: 'Design lead',
    company: 'Visionwork',
    face: 0,
  },
  {
    quote:
      'Before Tasker, keeping projects aligned felt impossible. Now everything runs smoothly, and our clients notice the difference.',
    name: 'Daniel Okonkwo',
    role: 'Operations',
    company: 'Constellation',
    face: 1,
  },
  {
    quote:
      'We deliver projects faster and with less stress — all thanks to Tasker. It is our team’s secret weapon for scaling.',
    name: 'Samir Haddad',
    role: 'Founder',
    company: 'Goodwill',
    face: 2,
  },
]

const FAQS = [
  {
    q: 'Do you offer fixed pricing or packages?',
    a: 'We offer custom plans based on your goals, services needed and project scope. Pricing starts at $600 a month and scales with the team, so a five-person workspace is not paying for a fifty-person one.',
  },
  {
    q: 'Is there a minimum commitment?',
    a: 'No. Monthly plans cancel from the billing page, and annual plans are refunded pro rata if you leave in the first quarter.',
  },
  {
    q: 'What kind of teams do you typically work with?',
    a: 'Product, design and operations teams between three and four hundred people, usually at the point where a spreadsheet has quietly become load-bearing.',
  },
  {
    q: 'How quickly can we get started?',
    a: 'Minutes. Import from CSV, Trello, Asana or Linear, and the structure comes across with the tasks rather than after them.',
  },
  {
    q: 'Can I start small and scale up later?',
    a: 'That is the usual path. Upgrades apply mid-cycle and are prorated, and nothing you created on the free plan is held hostage.',
  },
  {
    q: 'Will I have access to your team directly?',
    a: 'On Professional and above, yes — a shared channel with the people who build the product, not a queue in front of them.',
  },
]

/* ── Mocks ──────────────────────────────────────────────────────────────── */

function Frame({ children, className = '' }: { children: ReactNode; className?: string }) {
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
    <Frame className="shadow-2xl">
      <div className="flex items-center gap-2 border-b border-gray-200 px-4 py-2.5 dark:border-white/10">
        <span className="flex items-center gap-1.5 text-xs font-semibold">
          <span className="size-4 rounded bg-blue-600" />
          tasklyn
        </span>
        <span className="ml-4 h-6 w-48 rounded-md bg-gray-100 dark:bg-white/10" />
      </div>
      <div className="grid gap-4 p-4 sm:grid-cols-[8rem_1fr]">
        <div className="hidden space-y-1.5 sm:block">
          {['Home', 'My tasks', 'Inbox', 'Projects', 'Analytics'].map((row, i) => (
            <div
              key={row}
              className={`rounded-md px-2.5 py-1.5 text-[11px] ${
                i === 0 ? 'bg-blue-600 text-white' : 'text-gray-500 dark:text-gray-400'
              }`}
            >
              {row}
            </div>
          ))}
        </div>
        <div>
          <div className="grid grid-cols-3 gap-2 sm:grid-cols-5">
            {[['7', 'Total'], ['49', 'Active'], ['12', 'Assigned'], ['6', 'Completed'], ['3', 'Overdue']].map(([v, k]) => (
              <div key={k} className="rounded-lg border border-gray-200 p-2 dark:border-white/10">
                <p className="text-sm font-semibold tabular-nums">{v}</p>
                <p className="text-[10px] text-gray-500 dark:text-gray-400">{k}</p>
              </div>
            ))}
          </div>
          <div className="mt-3 grid gap-3 sm:grid-cols-2">
            <div className="space-y-1.5">
              <p className="text-[10px] font-medium text-gray-500 dark:text-gray-400">Assigned tasks</p>
              {['Web mockup', 'Landing page', 'Cart flow'].map((row) => (
                <div key={row} className="flex items-center justify-between rounded-md border border-gray-200 px-2 py-1.5 text-[10px] dark:border-white/10">
                  <span>{row}</span>
                  <span className="rounded bg-blue-100 px-1.5 py-0.5 text-blue-800 dark:bg-blue-500/20 dark:text-blue-200">Doing</span>
                </div>
              ))}
            </div>
            <div className="space-y-1.5">
              <p className="text-[10px] font-medium text-gray-500 dark:text-gray-400">Projects</p>
              {['New project', 'Video branding', 'Futurework'].map((row) => (
                <div key={row} className="flex items-center gap-2 rounded-md border border-gray-200 px-2 py-1.5 text-[10px] dark:border-white/10">
                  <span className="size-1.5 rounded-full bg-emerald-500" />
                  {row}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </Frame>
  )
}

function StepMock({ kind }: { kind: string }) {
  if (kind === 'signup') {
    return (
      <Frame className="p-4">
        <p className="text-center text-xs font-semibold">Sign up</p>
        <div className="mt-3 space-y-2">
          <div className="rounded-md border border-gray-200 px-2 py-2 text-[10px] text-gray-500 dark:border-white/10 dark:text-gray-400">
            user@example.com
          </div>
          <div className="rounded-md bg-blue-600 py-2 text-center text-[10px] font-medium text-white">Sign up free</div>
        </div>
      </Frame>
    )
  }
  if (kind === 'invite') {
    return (
      <Frame className="p-4">
        <p className="text-center text-xs font-semibold">Invite members</p>
        <div className="mt-3 space-y-2">
          {[['Mira Lydon', 'Project manager'], ['Sam Reyes', 'Team lead']].map(([n, r]) => (
            <div key={n} className="flex items-center gap-2 rounded-md border border-gray-200 p-2 dark:border-white/10">
              <span className="size-6 rounded-full bg-gray-200 dark:bg-white/10" />
              <span className="text-[10px]">
                <span className="block font-medium">{n}</span>
                <span className="block text-gray-500 dark:text-gray-400">{r}</span>
              </span>
            </div>
          ))}
        </div>
      </Frame>
    )
  }
  return (
    <Frame className="p-4">
      <div className="flex items-center justify-between text-[10px]">
        <span className="font-semibold">Tasks</span>
        <span className="rounded bg-blue-600 px-2 py-1 text-white">New</span>
      </div>
      <div className="mt-3 space-y-1.5">
        {[
          ['Design review', 'In progress', 'bg-amber-100 text-amber-900'],
          ['API contract', 'Blocked', 'bg-rose-100 text-rose-900'],
          ['Onboarding copy', 'Done', 'bg-emerald-100 text-emerald-900'],
          ['Billing migration', 'In progress', 'bg-amber-100 text-amber-900'],
        ].map(([task, status, tone]) => (
          <div key={task} className="flex items-center justify-between rounded-md border border-gray-200 px-2.5 py-2 text-[10px] dark:border-white/10">
            <span>{task}</span>
            <span className={`rounded px-1.5 py-0.5 ${tone}`}>{status}</span>
          </div>
        ))}
      </div>
    </Frame>
  )
}

/* ── Page ───────────────────────────────────────────────────────────────── */

export default function TaskManagementSaas() {
  return (
    <div className="min-h-screen bg-white font-sans text-gray-950 antialiased dark:bg-gray-950 dark:text-gray-50">
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:rounded-lg focus:bg-gray-950 focus:px-4 focus:py-2 focus:text-sm focus:text-white"
      >
        Skip to content
      </a>

      <header className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
          <a href="#" className="flex items-center gap-2 font-semibold">
            <span aria-hidden="true" className="size-5 rounded-md bg-blue-600" />
            tasklyn
          </a>
          <nav aria-label="Primary" className="hidden md:block">
            <ul role="list" className="flex items-center gap-8 text-sm">
              {NAV.map((item) => (
                <li key={item}>
                  <a href={`#${item.toLowerCase()}`} className="text-gray-600 hover:text-gray-950 dark:text-gray-300 dark:hover:text-white">
                    {item}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
          <a href="#" className="rounded-lg bg-gray-950 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
            View demo
          </a>
        </div>
      </header>

      <main id="main">
        {/* Hero */}
        <section className="relative isolate overflow-hidden px-6 pt-16 pb-0">
          <span
            aria-hidden="true"
            className="absolute inset-x-0 top-0 -z-10 h-[34rem] bg-gradient-to-b from-indigo-100 via-amber-50 to-white dark:from-indigo-500/10 dark:via-amber-500/5 dark:to-gray-950"
          />
          <div className="mx-auto max-w-4xl text-center">
            <p className="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-medium shadow-sm ring-1 ring-gray-200 dark:bg-gray-900 dark:ring-white/10">
              <span aria-hidden="true" className="size-1.5 rounded-full bg-blue-600" />
              Tasker raises $20M to build the future of work
            </p>
            <h1 className="mt-6 text-4xl font-semibold tracking-tight text-balance sm:text-6xl">
              Empower your team to achieve more, faster
              {/* The caret is a flourish from the design, not content. It is
                  hidden, and it stops blinking when the reader has asked for
                  less motion. */}
              <span
                aria-hidden="true"
                className="ml-1 inline-block h-[0.9em] w-[3px] translate-y-[0.08em] bg-blue-600 motion-safe:animate-pulse"
              />
            </h1>
            <p className="mx-auto mt-5 max-w-lg text-pretty text-gray-600 dark:text-gray-300">
              Whether you are just starting out or scaling fast, Tasker helps your team work smarter
              on what actually matters.
            </p>
            <div className="mt-8">
              <a href="#pricing" className="inline-block rounded-lg bg-blue-600 px-6 py-3 text-sm font-semibold text-white hover:bg-blue-700">
                Get started for free
              </a>
            </div>
            <ul role="list" className="mt-6 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-sm text-gray-600 dark:text-gray-300">
              {TRUST.map((item) => (
                <li key={item} className="flex items-center gap-1.5">
                  <span aria-hidden="true" className="text-emerald-600 dark:text-emerald-400">✓</span>
                  {item}
                </li>
              ))}
            </ul>
          </div>

          <div className="mx-auto mt-14 max-w-5xl">
            <AppMock />
          </div>
        </section>

        {/* Brands */}
        <section aria-labelledby="brands-heading" className="border-y border-gray-200 px-6 py-12 dark:border-white/10">
          <div className="mx-auto max-w-6xl text-center">
            <h2 id="brands-heading" className="text-sm text-gray-500 dark:text-gray-400">
              Brands that put trust in us
            </h2>
            <ul role="list" className="mt-8 grid grid-cols-2 gap-x-8 gap-y-5 sm:grid-cols-3 lg:grid-cols-5">
              {BRANDS.map((brand) => (
                <li key={brand} className="flex items-center justify-center gap-2 text-sm font-semibold text-gray-500 dark:text-gray-400">
                  <span aria-hidden="true" className="size-3 rounded-sm bg-gray-300 dark:bg-gray-600" />
                  {brand}
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Setup */}
        <section id="solution" aria-labelledby="setup-heading" className="px-6 py-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="setup-heading" className="text-center text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Simplify a complicated process
            </h2>

            {/* Ordered, so an <ol>. Six columns so 3+3 and 6 both tile exactly. */}
            <ol className="mt-12 grid gap-6 lg:grid-cols-6">
              {STEPS.map((step) => (
                <li
                  key={step.title}
                  className={`rounded-2xl border border-gray-200 bg-gray-50 p-6 dark:border-white/10 dark:bg-gray-900/50 ${
                    step.wide ? 'lg:col-span-6' : 'lg:col-span-3'
                  }`}
                >
                  <StepMock kind={step.mock} />
                  <h3 className="mt-6 font-semibold">{step.title}</h3>
                  <p className="mt-2 max-w-md text-sm text-pretty text-gray-600 dark:text-gray-400">{step.body}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>

        {/* Capabilities */}
        <section id="features" aria-labelledby="capabilities-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <h2 id="capabilities-heading" className="text-center text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Everything you need to succeed. Fast.
            </h2>

            {/* A list beside a panel, not role="tab". See the note at the top of
                the file: tabs are a keyboard contract and this block has no
                script to honour it. */}
            <div className="mt-12 grid gap-8 overflow-hidden rounded-2xl border border-gray-200 lg:grid-cols-[16rem_1fr] lg:gap-0 dark:border-white/10">
              <ul role="list" className="divide-y divide-gray-200 border-b border-gray-200 lg:border-r lg:border-b-0 dark:divide-white/10 dark:border-white/10">
                {CAPABILITIES.map((item, i) => (
                  <li key={item.name}>
                    <a
                      href="#"
                      className={`flex items-center gap-3 px-5 py-3.5 text-sm hover:bg-gray-50 dark:hover:bg-white/5 ${
                        i === 0 ? 'bg-blue-50 font-medium text-blue-800 dark:bg-blue-500/10 dark:text-blue-200' : ''
                      }`}
                    >
                      <span aria-hidden="true" className="size-1.5 rounded-full bg-current opacity-40" />
                      {item.name}
                    </a>
                  </li>
                ))}
              </ul>

              <div className="p-6 sm:p-8">
                <h3 className="font-semibold">{CAPABILITIES[0].name}</h3>
                <p className="mt-2 max-w-xl text-sm text-pretty text-gray-600 dark:text-gray-400">{CAPABILITIES[0].body}</p>
                <div className="mt-6">
                  <Frame className="p-4">
                    <div className="flex h-40 items-end gap-2">
                      {[35, 52, 44, 68, 55, 78, 62, 88, 72, 95, 80, 66].map((h, i) => (
                        <span
                          key={i}
                          style={{ height: `${h}%` }}
                          className="flex-1 rounded-t bg-gradient-to-t from-blue-600/25 to-blue-600"
                        />
                      ))}
                    </div>
                  </Frame>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Integrations */}
        <section id="integration" aria-labelledby="integrations-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-3xl rounded-2xl border border-gray-200 p-8 text-center sm:p-12 dark:border-white/10">
            <p className="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1 text-xs font-medium dark:bg-white/10">
              <span aria-hidden="true" className="size-1.5 rounded-full bg-blue-600" />
              With a global community
            </p>
            <h2 id="integrations-heading" className="mt-5 text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
              Seamless integrations, zero friction
            </h2>

            {/* The hub occupies a cell of its own rather than floating over the
                middle of the grid. Absolutely centred, it landed on the corners
                of four tiles and covered their labels — a hub that hides the
                things it connects to.

                It is placed explicitly at column 2, row 2. Grid resolves
                explicitly placed items first, so the eight tiles auto-place
                into the remaining cells and flow around it. The list keeps
                role="list" because display:contents drops list semantics in
                some browsers. */}
            <div className="mx-auto mt-10 grid max-w-sm grid-cols-3 gap-3">
              <span
                aria-hidden="true"
                className="col-start-2 row-start-2 mx-auto flex size-14 items-center justify-center rounded-full bg-blue-600 text-xl text-white"
              >
                ✓
              </span>
              <ul role="list" className="contents">
                {INTEGRATIONS.map((tool) => (
                  <li
                    key={tool}
                    className="flex items-center justify-center rounded-xl border border-gray-200 px-2 py-3 text-xs font-medium dark:border-white/10"
                  >
                    {tool}
                  </li>
                ))}
              </ul>
            </div>

            <p className="mx-auto mt-8 max-w-sm text-pretty text-gray-600 dark:text-gray-300">
              Plug Tasker into your workflow in a click. No messy setup, no dev time needed, just
              instant productivity.
            </p>
            <div className="mt-6">
              <a href="#pricing" className="inline-block rounded-lg bg-gray-950 px-6 py-3 text-sm font-semibold text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950">
                Try Tasker for free
              </a>
            </div>
          </div>
        </section>

        {/* Pricing */}
        <section id="pricing" aria-labelledby="pricing-heading" className="px-6 pb-20">
          <div className="mx-auto max-w-6xl">
            <div className="text-center">
              <p className="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1 text-xs font-medium dark:bg-white/10">
                <span aria-hidden="true" className="size-1.5 rounded-full bg-blue-600" />
                Pricing
              </p>
              <h2 id="pricing-heading" className="mt-5 text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Start free, upgrade as you scale
              </h2>

              {/* Two named periods, so two radios in a fieldset with a legend.
                  A switch would announce "annual pricing, on", which is not the
                  question being asked. */}
              <fieldset className="mt-6 inline-block">
                <legend className="sr-only">Billing period</legend>
                <div className="inline-flex rounded-full border border-gray-200 p-1 dark:border-white/10">
                  {[
                    { id: 'billing-monthly', label: 'Monthly', checked: false },
                    { id: 'billing-yearly', label: 'Yearly, save 25%', checked: true },
                  ].map((option) => (
                    <div key={option.id}>
                      <input
                        type="radio"
                        id={option.id}
                        name="billing-period"
                        defaultChecked={option.checked}
                        className="peer sr-only"
                      />
                      <label
                        htmlFor={option.id}
                        className="block cursor-pointer rounded-full px-4 py-1.5 text-sm text-gray-600 peer-checked:bg-blue-600 peer-checked:font-medium peer-checked:text-white peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-blue-600 dark:text-gray-300"
                      >
                        {option.label}
                      </label>
                    </div>
                  ))}
                </div>
              </fieldset>
            </div>

            <ul role="list" className="mt-12 grid items-start gap-6 lg:grid-cols-3">
              {TIERS.map((tier) => (
                <li
                  key={tier.name}
                  className={`flex h-full flex-col rounded-2xl border p-8 ${
                    tier.featured
                      ? 'border-blue-600 ring-1 ring-blue-600'
                      : 'border-gray-200 dark:border-white/10'
                  }`}
                >
                  <div className="flex items-center justify-between gap-3">
                    <h3 className="font-semibold">{tier.name}</h3>
                    {tier.badge && (
                      <span className="rounded-full bg-blue-600 px-2.5 py-1 text-xs font-medium text-white">
                        {tier.badge}
                      </span>
                    )}
                  </div>
                  <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-400">{tier.blurb}</p>

                  <p className="mt-8">
                    <span className="text-4xl font-semibold tracking-tight tabular-nums">{tier.price}</span>
                  </p>
                  {tier.period && <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">{tier.period}</p>}

                  {/* The link name says which plan it buys. Three links called
                      "Get started" are three identical rows in a link list, and
                      the card around them is not available there. */}
                  <a
                    href="#"
                    className={`mt-6 block rounded-lg px-5 py-3 text-center text-sm font-semibold ${
                      tier.featured
                        ? 'bg-blue-600 text-white hover:bg-blue-700'
                        : 'bg-gray-950 text-white hover:bg-gray-800 dark:bg-white dark:text-gray-950'
                    }`}
                  >
                    {tier.cta}
                  </a>

                  <ul role="list" className="mt-8 space-y-2.5 text-sm">
                    {tier.inherits && <li className="font-medium">Everything in {tier.inherits}, plus:</li>}
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex gap-2.5">
                        <span aria-hidden="true" className="mt-0.5 text-blue-700 dark:text-blue-400">✓</span>
                        <span className="text-gray-600 dark:text-gray-300">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* Testimonials */}
        <section aria-labelledby="love-heading" className="bg-gray-50 px-6 py-20 dark:bg-gray-900/40">
          <div className="mx-auto max-w-6xl">
            <div className="text-center">
              <p className="inline-flex items-center gap-2 rounded-full bg-white px-3 py-1 text-xs font-medium ring-1 ring-gray-200 dark:bg-gray-900 dark:ring-white/10">
                <span aria-hidden="true" className="size-1.5 rounded-full bg-blue-600" />
                Wall of love
              </p>
              <h2 id="love-heading" className="mt-5 text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Why teams love Tasker
              </h2>
              <p className="mx-auto mt-3 max-w-lg text-pretty text-gray-600 dark:text-gray-300">
                Thousands of teams trust Tasker to deliver projects with clarity and speed.
              </p>
            </div>

            <ul role="list" className="mt-12 grid gap-6 lg:grid-cols-3">
              {TESTIMONIALS.map((item) => (
                <li key={item.name}>
                  <figure className="flex h-full flex-col rounded-2xl border border-gray-200 bg-white p-6 dark:border-white/10 dark:bg-gray-950">
                    <p className="text-sm font-semibold text-gray-500 dark:text-gray-400">{item.company}</p>
                    <blockquote className="mt-4 text-sm text-pretty text-gray-700 dark:text-gray-200">
                      {item.quote}
                    </blockquote>
                    <figcaption className="mt-auto flex items-center gap-3 pt-6">
                      <img
                        src={face(FACES[item.face])}
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
                </li>
              ))}
            </ul>
          </div>
        </section>

        {/* FAQ */}
        <section aria-labelledby="faq-heading" className="px-6 py-20">
          <div className="mx-auto max-w-3xl">
            <div className="text-center">
              <p className="inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1 text-xs font-medium dark:bg-white/10">
                <span aria-hidden="true" className="size-1.5 rounded-full bg-blue-600" />
                FAQs
              </p>
              <h2 id="faq-heading" className="mt-5 text-3xl font-semibold tracking-tight text-balance sm:text-4xl">
                Frequent queries explained
              </h2>
              <p className="mx-auto mt-3 max-w-md text-pretty text-gray-600 dark:text-gray-300">
                A global community, countless trusted teams, and top-notch ratings.
              </p>
            </div>

            <div className="mt-10">
              {FAQS.map((item) => (
                <details key={item.q} className="group border-b border-gray-200 dark:border-white/10">
                  <summary className="flex cursor-pointer list-none items-start justify-between gap-4 py-4 text-sm font-medium marker:content-none">
                    {item.q}
                    {/* gray-600 rather than gray-400: list-none removes the
                        native marker, so this glyph is the only sign the row
                        opens and it owes the 3:1 of WCAG 1.4.11. */}
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
      </main>

      <footer className="border-t border-gray-200 px-6 pt-14 dark:border-white/10">
        <div className="mx-auto grid max-w-6xl gap-10 sm:grid-cols-2 lg:grid-cols-4">
          <div>
            <p className="flex items-center gap-2 font-semibold">
              <span aria-hidden="true" className="size-5 rounded-md bg-blue-600" />
              tasklyn
            </p>
            <p className="mt-3 max-w-xs text-sm text-pretty text-gray-600 dark:text-gray-400">
              Work that stays in one place, from the first task to the last invoice.
            </p>
          </div>
          {[
            { heading: 'Product', links: ['Features', 'Integration', 'Pricing', 'Changelog'] },
            { heading: 'Company', links: ['Solution', 'Testimonials', 'About', 'Careers'] },
            { heading: 'Social', links: ['Dribbble', 'LinkedIn', 'YouTube'] },
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

        {/* The oversized wordmark repeats the name above it, so it is a graphic
            rather than a second announcement of the same word. */}
        <p
          aria-hidden="true"
          className="mx-auto mt-14 max-w-6xl text-6xl font-semibold tracking-tight text-gray-100 select-none sm:text-8xl dark:text-white/5"
        >
          tasklyn
        </p>

        <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-4 border-t border-gray-200 py-6 text-sm text-gray-500 dark:border-white/10 dark:text-gray-400">
          <p>© 2026 Tasker Limited. All rights reserved.</p>
          <ul role="list" className="flex gap-6">
            <li>
              <a href="#" className="hover:text-gray-950 dark:hover:text-white">
                Privacy policy
              </a>
            </li>
            <li>
              <a href="#" className="hover:text-gray-950 dark:hover:text-white">
                Terms and conditions
              </a>
            </li>
          </ul>
        </div>
      </footer>
    </div>
  )
}
