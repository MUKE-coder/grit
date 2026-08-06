/*
 * Three figures in a bordered row, with crosshairs at the corners.
 *
 * The crosshairs are drawn with pseudo-elements on the frame rather than as
 * four positioned <span>s, so they cost nothing in the DOM and cannot be left
 * behind if the frame is restyled. They are decorative and unreachable, which
 * is exactly right: they mark a boundary the border already draws.
 *
 * The dividers between the figures are `divide-x` above the breakpoint and
 * `divide-y` below it, so the line always runs across the direction the
 * figures are laid out in. A vertical rule between stacked rows is a line
 * pointing nowhere.
 *
 * The figures are a <dl>, and the value comes before the label in the markup
 * so the visual order and the DOM order agree.
 */

export interface Stat {
  value: string
  label: string
}

const STATS: Stat[] = [
  { value: '+85%', label: 'Conversion rate' },
  { value: '12K', label: 'Active users' },
  { value: '40%', label: 'Revenue growth' },
]

const CROSSHAIR =
  "before:absolute before:-top-px before:-left-px before:size-2 before:border-t before:border-l before:border-gray-400 before:content-[''] after:absolute after:-right-px after:-bottom-px after:size-2 after:border-r after:border-b after:border-gray-400 after:content-[''] dark:before:border-white/30 dark:after:border-white/30"

export default function StatsBoxedThreeUp({
  title = 'Delivering measurable results',
  subtitle = 'What the platform has done for the companies already running on it.',
  stats = STATS,
}: {
  title?: string
  subtitle?: string
  stats?: Stat[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <dl
          className={`relative mx-auto mt-16 grid max-w-4xl grid-cols-1 divide-y divide-gray-200 border border-gray-200 sm:grid-cols-3 sm:divide-x sm:divide-y-0 dark:divide-white/10 dark:border-white/10 ${CROSSHAIR}`}
        >
          {stats.map((stat) => (
            <div key={stat.label} className="px-8 py-10 text-center">
              <dd className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums sm:text-5xl dark:text-white">
                {stat.value}
              </dd>
              <dt className="mt-2 text-sm text-gray-600 dark:text-gray-400">{stat.label}</dt>
            </div>
          ))}
        </dl>
      </div>
    </section>
  )
}
