/*
 * Three steps in a framed grid, each with a preview above the copy.
 *
 * The frame is drawn with dividers on the grid rather than borders on the
 * cards, so there is one line between neighbours instead of two sitting on top
 * of each other at slightly different opacities. `divide-x` only applies above
 * the breakpoint where the columns exist; below it `divide-y` takes over,
 * because a vertical rule between stacked rows is a line pointing nowhere.
 *
 * The preview sits above the title here, which is the opposite of the arrow
 * layout. It works when the previews are interesting enough to be the reason
 * someone stops scrolling, and badly when they are generic icons.
 */

function CampaignPreview() {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-900">
      <p className="text-sm font-semibold text-gray-900 dark:text-white">Campaign</p>
      <p className="mt-0.5 text-xs text-gray-500 dark:text-gray-400">Loyalty programme</p>
      <div className="mt-4 grid grid-cols-2 gap-2">
        {['Start date', 'End date'].map((label) => (
          <div key={label} className="border-l-2 border-indigo-500 pl-2 dark:border-indigo-400">
            <p className="text-xs font-medium text-gray-900 dark:text-white">{label}</p>
            <p className="text-[11px] text-gray-500 tabular-nums dark:text-gray-400">
              Feb 6, 2024
            </p>
          </div>
        ))}
      </div>
      <p className="mt-4 text-xs text-gray-500 dark:text-gray-400">
        Connected to 12{' '}
        <span className="text-indigo-600 dark:text-indigo-400">marketing campaigns</span>.
      </p>
    </div>
  )
}

function PollPreview() {
  const events = [
    { time: '06 AM', label: 'Poll created', highlight: false },
    { time: '12 PM', label: '+50 users voted', highlight: true },
    { time: '01 PM', label: 'Poll closed', highlight: false },
  ]
  return (
    <div className="space-y-2">
      {events.map((event) => (
        <div
          key={event.label}
          className={`flex items-center gap-3 rounded-lg px-3 py-2 ${
            event.highlight
              ? 'border border-gray-200 bg-white shadow-sm dark:border-white/10 dark:bg-gray-900'
              : ''
          }`}
        >
          <span className="size-1.5 flex-none rounded-full border border-gray-400 dark:border-gray-600" />
          <span className="text-[11px] text-gray-500 tabular-nums dark:text-gray-400">
            {event.time}
          </span>
          <span className="text-xs font-semibold text-gray-900 dark:text-white">{event.label}</span>
        </div>
      ))}
    </div>
  )
}

function MemoryPreview() {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-gray-900">
      <p className="text-sm font-semibold text-gray-900 dark:text-white">Memory usage</p>
      <div className="mt-3 flex items-baseline justify-between">
        <p className="text-lg font-semibold text-gray-900 tabular-nums dark:text-white">
          56 GB <span className="text-sm font-normal text-gray-500">/ 128 GB</span>
        </p>
        <p className="text-sm text-indigo-600 tabular-nums dark:text-indigo-400">45%</p>
      </div>
      <div className="mt-3 h-2 overflow-hidden rounded-full bg-gray-100 dark:bg-white/10">
        <span className="block h-full w-[45%] rounded-full bg-indigo-500 dark:bg-indigo-400" />
      </div>
    </div>
  )
}

const STEPS = [
  {
    title: 'Secure messaging',
    body: 'End-to-end encrypted communication, with keys your provider never holds.',
    Preview: CampaignPreview,
  },
  {
    title: 'Analytics dashboard',
    body: 'Visualisation tools that turn a pile of metrics into something you can act on.',
    Preview: PollPreview,
  },
  {
    title: 'Resource monitoring',
    body: 'Real-time tracking of system performance, so you find the ceiling before your users do.',
    Preview: MemoryPreview,
  },
]

export default function HowItWorksFramedColumnsWithPreviews({
  title = 'Simple three-step workflow',
  subtitle = 'A streamlined approach that lets your team make informed decisions quickly.',
  steps = STEPS,
}: {
  title?: string
  subtitle?: string
  steps?: { title: string; body: string; Preview: () => React.JSX.Element }[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        {/* One line between neighbours, not two borders meeting. */}
        <ol
          role="list"
          className="mt-16 grid grid-cols-1 divide-y divide-gray-200 rounded-2xl border border-gray-200 lg:grid-cols-3 lg:divide-x lg:divide-y-0 dark:divide-white/10 dark:border-white/10"
        >
          {steps.map((step, i) => (
            <li key={step.title} className="flex flex-col p-8">
              <span className="text-sm font-semibold text-gray-400 tabular-nums dark:text-gray-600">
                {String(i + 1).padStart(2, '0')}
              </span>
              <div aria-hidden="true" className="mt-6 flex min-h-40 items-center select-none">
                <div className="w-full">
                  <step.Preview />
                </div>
              </div>
              <h3 className="mt-8 text-lg font-semibold text-gray-900 dark:text-white">
                {step.title}
              </h3>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">{step.body}</p>
            </li>
          ))}
        </ol>
      </div>
    </section>
  )
}
