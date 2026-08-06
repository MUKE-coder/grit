/*
 * A claim on the left, the figures that back it on the right.
 *
 * The rule beside each figure is a left border on the group rather than a
 * separate element, so it is exactly as tall as the text next to it and cannot
 * be left behind when the copy wraps to another line at a narrower width.
 *
 * The heading takes a `highlight` prop rather than accepting HTML, because a
 * block that renders arbitrary markup from a prop is an injection hole in
 * every project that installs it. Two plain strings and a <strong> between
 * them does the same job with nothing to escape.
 */

export interface Stat {
  value: string
  label: string
  body: string
}

const STATS: Stat[] = [
  {
    value: '99.9%',
    label: 'Uptime guarantee',
    body: 'for every service, measured from outside our own network.',
  },
  {
    value: '24/7',
    label: 'Support available',
    body: 'from the engineers who work on the product, not a ticket queue.',
  },
]

export default function StatsHeadingWithRuledStats({
  titleStart = 'Building the next generation of ',
  highlight = 'AI-powered marketing tools',
  titleEnd = '',
  body = 'The pipeline reads, interprets and acts on data that used to need a person in the loop.',
  stats = STATS,
}: {
  titleStart?: string
  highlight?: string
  titleEnd?: string
  body?: string
  stats?: Stat[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-2">
          <div>
            <h2 className="max-w-md text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {titleStart}
              <strong className="font-semibold">{highlight}</strong>
              {titleEnd}
            </h2>
            <p className="mt-6 max-w-md text-base/7 text-pretty text-gray-600 dark:text-gray-400">
              {body}
            </p>
          </div>

          <dl className="space-y-10">
            {stats.map((stat) => (
              /* The rule is this border, so it matches the height of the text. */
              <div
                key={stat.label}
                className="border-l-2 border-gray-200 pl-6 dark:border-white/15"
              >
                <dd className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                  {stat.value}
                </dd>
                <dt className="mt-2 max-w-sm text-sm/6 text-gray-600 dark:text-gray-400">
                  <span className="font-medium text-gray-900 dark:text-white">{stat.label}</span>{' '}
                  {stat.body}
                </dt>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </section>
  )
}
