/*
 * A long-form walkthrough: numbered nodes on a rail, each with real estate for
 * a screenshot, a stat pair or a quote.
 *
 * The rail is a border on the list rather than an absolutely positioned line,
 * so it grows with the content and cannot end early when a step runs long.
 * Each node is pulled back over the rail with a negative left margin, and gets
 * a ring in the page background colour so the line appears to stop behind it
 * instead of running through it.
 *
 * That ring is `ring-white dark:ring-gray-950` — it has to match the section
 * background, not be transparent. If you change the section colour, change the
 * ring with it or the line will show through the node.
 *
 * Reach for this when the steps need explaining rather than listing. If each
 * step is one sentence, use one of the three-across blocks instead; this one
 * will just look empty.
 */

function GanttPanel() {
  const rows = [
    { label: 'Discovery', start: 0, width: 28, tone: 'bg-indigo-500' },
    { label: 'Design', start: 22, width: 34, tone: 'bg-violet-500' },
    { label: 'Build', start: 48, width: 40, tone: 'bg-sky-500' },
    { label: 'Launch', start: 82, width: 16, tone: 'bg-emerald-500' },
  ]
  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
      <div className="flex items-center gap-4 border-b border-gray-200 px-4 py-3 text-xs dark:border-white/10">
        {['Timeline', 'Board', 'Workflow'].map((tab, i) => (
          <span
            key={tab}
            className={
              i === 0
                ? 'font-semibold text-gray-900 dark:text-white'
                : 'text-gray-500 dark:text-gray-400'
            }
          >
            {tab}
          </span>
        ))}
      </div>
      <div className="p-4">
        <div className="flex justify-between text-[10px] text-gray-400 dark:text-gray-600">
          {['Jan', 'Mar', 'May', 'Jul', 'Sep', 'Nov'].map((month) => (
            <span key={month}>{month}</span>
          ))}
        </div>
        <div className="mt-3 space-y-2.5">
          {rows.map((row) => (
            <div key={row.label} className="flex items-center gap-3">
              <span className="w-16 flex-none text-xs text-gray-600 dark:text-gray-400">
                {row.label}
              </span>
              <span className="relative h-4 flex-1 rounded bg-gray-100 dark:bg-white/5">
                <span
                  className={`absolute inset-y-0 rounded ${row.tone}`}
                  style={{ left: `${row.start}%`, width: `${row.width}%` }}
                />
              </span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

function CustomersPanel() {
  const rows = [
    { date: '10/31/2023', status: 'Paid', name: 'Bernard Ng', amount: '$43.99' },
    { date: '10/21/2023', status: 'Refunded', name: 'Méschac Irung', amount: '$19.99' },
    { date: '10/15/2023', status: 'Paid', name: 'Glodie Ng', amount: '$99.99' },
    { date: '10/12/2023', status: 'Cancelled', name: 'Theo Ng', amount: '$19.99' },
  ]
  const tone: Record<string, string> = {
    Paid: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400',
    Refunded: 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400',
    Cancelled: 'bg-rose-50 text-rose-700 dark:bg-rose-500/10 dark:text-rose-400',
  }
  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
      <div className="border-b border-gray-200 px-4 py-3 dark:border-white/10">
        <p className="text-sm font-semibold text-gray-900 dark:text-white">Customers</p>
        <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">
          New users by first primary channel
        </p>
      </div>
      <ul className="divide-y divide-gray-100 dark:divide-white/5">
        {rows.map((row) => (
          <li key={row.name} className="flex items-center gap-3 px-4 py-2.5 text-xs">
            <span className="w-20 flex-none text-gray-500 tabular-nums dark:text-gray-400">
              {row.date}
            </span>
            <span className={`rounded-full px-2 py-0.5 font-medium ${tone[row.status]}`}>
              {row.status}
            </span>
            <span className="min-w-0 flex-1 truncate text-gray-900 dark:text-white">{row.name}</span>
            <span className="font-mono text-gray-900 tabular-nums dark:text-white">
              {row.amount}
            </span>
          </li>
        ))}
      </ul>
    </div>
  )
}

export interface Step {
  title: string
  body: string
  stats?: { value: string; label: string }[]
  Panel?: () => React.JSX.Element
  quote?: { text: string; name: string; handle: string }
}

const STEPS: Step[] = [
  {
    title: 'Turn your data into something you can look at',
    body: 'Point the platform at your database and it builds the charts, filters and drill-downs from the schema it finds. No mapping layer to maintain.',
    stats: [
      { value: '90+', label: 'Integrations' },
      { value: '56%', label: 'Less setup time' },
    ],
  },
  {
    title: 'Manage the work in one place',
    body: 'Plan, track and schedule from a single timeline. Allocate people, watch the critical path, and keep everyone pointed at the same milestone.',
    Panel: GanttPanel,
    quote: {
      text: 'The fusion of simplicity and range here is unusual. We built something that looks considered without a designer on the team.',
      name: 'Glodie Lukose',
      handle: '@glodie',
    },
  },
  {
    title: 'Arrange it the way you work',
    body: 'Move panels, resize them, and save the arrangement per team. The layout is data, so it survives a redeploy.',
    Panel: CustomersPanel,
  },
]

export default function HowItWorksTimelineWithPanels({
  title = 'Set up your pipeline in minutes',
  subtitle = 'The platform reads your data, builds the interface around it, and gets out of the way.',
  steps = STEPS,
}: {
  title?: string
  subtitle?: string
  steps?: Step[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-4xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        {/* The rail is this border, so it is exactly as long as the steps are. */}
        <ol
          role="list"
          className="mt-20 space-y-20 border-l border-gray-200 pl-8 sm:pl-12 dark:border-white/10"
        >
          {steps.map((step, i) => (
            <li key={step.title} className="relative">
              {/* Pulled back over the rail. The ring matches the section
                  background so the line appears to stop behind the node. */}
              <span
                aria-hidden="true"
                className="absolute top-1 -left-8 flex size-7 items-center justify-center rounded-full bg-gray-100 text-xs font-semibold text-gray-700 ring-4 ring-white tabular-nums sm:-left-12 dark:bg-white/10 dark:text-gray-300 dark:ring-gray-950"
              >
                {i + 1}
              </span>

              <h3 className="text-xl font-semibold text-gray-900 dark:text-white">{step.title}</h3>
              <p className="mt-3 max-w-2xl text-base/7 text-gray-600 dark:text-gray-400">
                {step.body}
              </p>

              {step.stats && (
                <dl className="mt-8 flex gap-12">
                  {step.stats.map((stat) => (
                    <div key={stat.label}>
                      <dt className="sr-only">{stat.label}</dt>
                      <dd className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                        {stat.value}
                      </dd>
                      <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">{stat.label}</p>
                    </div>
                  ))}
                </dl>
              )}

              {step.Panel && (
                <div aria-hidden="true" className="mt-8 select-none">
                  <step.Panel />
                </div>
              )}

              {step.quote && (
                <figure className="mt-8 border-l-2 border-indigo-500 pl-4 dark:border-indigo-400">
                  <blockquote className="text-sm/6 text-gray-700 dark:text-gray-300">
                    {step.quote.text}
                  </blockquote>
                  <figcaption className="mt-3 text-sm text-gray-600 dark:text-gray-400">
                    <span className="font-medium text-gray-900 dark:text-white">
                      {step.quote.name}
                    </span>{' '}
                    {step.quote.handle}
                  </figcaption>
                </figure>
              )}
            </li>
          ))}
        </ol>
      </div>
    </section>
  )
}
