/*
 * A sticky title beside numbered steps, each with a panel showing the result.
 *
 * The title column sticks while the steps scroll past it, which is the point:
 * a long sequence loses its heading off the top of the screen at exactly the
 * moment the reader needs to remember what the sequence is for.
 *
 * `lg:self-start` is what makes the sticky work. A grid item stretches to the
 * row height by default, so the element is already as tall as the column and
 * has nothing to stick within. Without it `position: sticky` is silently inert
 * — no error, no warning, it just never sticks, which is why this is the most
 * common way to get this layout wrong.
 *
 * Rows are separated by a border on the list item rather than a divider
 * element, so the last one can drop its rule with `last:border-0`.
 */

function MonitoringPanel() {
  const bars = [38, 62, 45, 78, 55, 88, 70, 96, 64, 82]
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900">
      <p className="text-sm font-semibold text-gray-900 dark:text-white">Monitoring</p>
      <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">January – June 2024</p>
      <div className="mt-6 flex h-28 items-end gap-1.5">
        {bars.map((height, i) => (
          <span
            key={i}
            className="flex-1 rounded-t-sm bg-indigo-500/70 dark:bg-indigo-400/60"
            style={{ height: `${height}%` }}
          />
        ))}
      </div>
    </div>
  )
}

function InvoicePanel() {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900">
      <div className="flex items-start justify-between">
        <div>
          <p className="font-mono text-xs text-gray-500 dark:text-gray-400">INV-456789</p>
          <p className="mt-1 font-mono text-xl font-semibold text-gray-900 tabular-nums dark:text-white">
            $284,342.57
          </p>
        </div>
        <span className="rounded-full bg-emerald-50 px-2 py-1 text-xs font-medium text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400">
          Sent
        </span>
      </div>
      <div className="mt-5 space-y-2">
        {['Design retainer, June', 'Additional revisions', 'Expenses'].map((line, i) => (
          <div
            key={line}
            className="flex items-center justify-between border-t border-gray-100 pt-2 text-xs dark:border-white/5"
          >
            <span className="text-gray-600 dark:text-gray-400">{line}</span>
            <span className="font-mono text-gray-900 tabular-nums dark:text-white">
              ${[240000, 38400, 5942].map((n) => n.toLocaleString())[i]}
            </span>
          </div>
        ))}
      </div>
    </div>
  )
}

function TerminalPanel() {
  const lines = [
    { text: '$ grit deploy --prod', tone: 'text-gray-900 dark:text-white' },
    { text: 'building api ... done in 1.4s', tone: 'text-gray-500 dark:text-gray-400' },
    { text: 'building web ... done in 3.1s', tone: 'text-gray-500 dark:text-gray-400' },
    { text: 'live at https://example.com', tone: 'text-emerald-600 dark:text-emerald-400' },
  ]
  return (
    <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
      <div className="flex gap-1.5 border-b border-gray-200 px-4 py-3 dark:border-white/10">
        {['bg-red-400', 'bg-amber-400', 'bg-emerald-400'].map((tone) => (
          <span key={tone} className={`size-2.5 rounded-full ${tone}`} />
        ))}
      </div>
      <div className="space-y-1.5 p-4 font-mono text-xs">
        {lines.map((line) => (
          <p key={line.text} className={line.tone}>
            {line.text}
          </p>
        ))}
      </div>
    </div>
  )
}

const STEPS = [
  {
    title: 'Connect your data',
    body: 'Point it at your database and it reads the schema. No mapping file, no manual model definitions.',
    Panel: MonitoringPanel,
  },
  {
    title: 'Send the invoice',
    body: 'Generate, send and track from one place, with every state change recorded against the record.',
    Panel: InvoicePanel,
  },
  {
    title: 'Ship it',
    body: 'One command builds the API, the admin panel and the client, and puts them behind your domain.',
    Panel: TerminalPanel,
  },
]

export default function HowItWorksStickyTitleWithSteps({
  title = 'Set up your pipeline in minutes',
  subtitle = 'The platform reads your data, builds the interface around it, and gets out of the way.',
  steps = STEPS,
}: {
  title?: string
  subtitle?: string
  steps?: { title: string; body: string; Panel: () => React.JSX.Element }[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-[minmax(0,20rem)_minmax(0,1fr)]">
          {/* self-start is load-bearing: a stretched grid item has no room to
              stick inside, and sticky fails silently when that happens. */}
          <div className="lg:sticky lg:top-24 lg:self-start">
            <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h2>
            <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
          </div>

          <ol role="list">
            {steps.map((step, i) => (
              <li
                key={step.title}
                className="border-b border-gray-200 pb-12 last:border-0 last:pb-0 [&:not(:first-child)]:pt-12 dark:border-white/10"
              >
                <span className="inline-flex size-8 items-center justify-center rounded-full bg-gray-100 text-sm font-semibold text-gray-700 tabular-nums dark:bg-white/10 dark:text-gray-300">
                  {i + 1}
                </span>
                <h3 className="mt-6 text-xl font-semibold text-gray-900 dark:text-white">
                  {step.title}
                </h3>
                <p className="mt-2 max-w-xl text-base/7 text-gray-600 dark:text-gray-400">
                  {step.body}
                </p>
                <div aria-hidden="true" className="mt-8 select-none">
                  <step.Panel />
                </div>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </section>
  )
}
