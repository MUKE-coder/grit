import { ArrowDownRight, ArrowUpRight, CreditCard, Package, Users, Wallet } from 'lucide-react'

/*
 * A row of metric tiles, each with a value, a trend and a sparkline.
 *
 * It is a <dl>. A metric is a label and a value, which is exactly what a
 * description list is for, and it gets the pairing announced rather than
 * leaving two visually adjacent divs to imply it. The tile is a <div> inside,
 * so the dt/dd relationship survives the styling.
 *
 * The trend is a sentence, not an arrow. "12.5%" beside a green up-arrow tells
 * a screen reader user "12.5 percent" and nothing about direction: the arrow
 * is a glyph and the colour is a colour, and neither is text. Every tile here
 * says "up" or "down" and what the comparison is against, with the arrow
 * marked decorative.
 *
 * The detail links name their metric. The source had four links all reading
 * "View details", which in a screen reader's list of links is four identical
 * entries with no way to tell which is which.
 *
 * Sparklines are aria-hidden. They are a shape, and read aloud a polyline is a
 * string of coordinates. The number and the trend beside them already carry
 * everything the shape suggests.
 *
 * Values are pre-formatted strings rather than numbers formatted in the
 * component, because a metric tile is usually fed something already rounded
 * and suffixed by whatever produced it. Money elsewhere in this registry is
 * integer cents; this deliberately is not, and takes display strings.
 */

export interface Metric {
  label: string
  value: string
  /** Positive is up. The sign decides the wording and the colour. */
  changePercent: number
  comparison: string
  href: string
  Icon: typeof Users
  tone: string
  /** Sparkline points, 0-100, oldest first. Decorative. */
  series: number[]
}

const METRICS: Metric[] = [
  {
    label: 'Revenue',
    value: '$48,120',
    changePercent: 12.5,
    comparison: 'vs last month',
    href: '#revenue',
    Icon: Wallet,
    tone: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300',
    series: [30, 42, 38, 55, 48, 67, 72, 88],
  },
  {
    label: 'Orders',
    value: '1,284',
    changePercent: 4.1,
    comparison: 'vs last month',
    href: '#orders',
    Icon: Package,
    tone: 'bg-sky-100 text-sky-700 dark:bg-sky-500/15 dark:text-sky-300',
    series: [50, 46, 58, 52, 61, 57, 64, 68],
  },
  {
    label: 'Customers',
    value: '892',
    changePercent: -2.4,
    comparison: 'vs last month',
    href: '#customers',
    Icon: Users,
    tone: 'bg-violet-100 text-violet-700 dark:bg-violet-500/15 dark:text-violet-300',
    series: [70, 68, 72, 64, 60, 58, 55, 52],
  },
  {
    label: 'Refunds',
    value: '14',
    changePercent: -18.2,
    comparison: 'vs last month',
    href: '#refunds',
    Icon: CreditCard,
    tone: 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300',
    series: [60, 55, 48, 50, 40, 36, 30, 26],
  },
]

/* A polyline across a 100x32 box. Decorative, so no title and no role. */
function Sparkline({ series, rising }: { series: number[]; rising: boolean }) {
  const points = series
    .map((value, index) => {
      const x = (index / (series.length - 1)) * 100
      const y = 32 - (value / 100) * 32
      return `${x.toFixed(1)},${y.toFixed(1)}`
    })
    .join(' ')

  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 100 32"
      preserveAspectRatio="none"
      className={`h-8 w-full ${rising ? 'text-emerald-500' : 'text-rose-500'}`}
    >
      <polyline
        points={points}
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        vectorEffect="non-scaling-stroke"
      />
    </svg>
  )
}

export default function MetricTilesWithTrend({
  metrics = METRICS,
  heading = 'This month',
}: {
  metrics?: Metric[]
  heading?: string
}) {
  return (
    <section aria-labelledby="metrics" className="bg-gray-50 py-10 dark:bg-gray-950">
      <div className="mx-auto max-w-6xl px-4">
        <h2
          id="metrics"
          className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white"
        >
          {heading}
        </h2>

        {/* dl, not a grid of divs: each tile is a label and a value. */}
        <dl className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          {metrics.map((metric) => {
            const rising = metric.changePercent >= 0
            const magnitude = Math.abs(metric.changePercent)
            return (
              <div
                key={metric.label}
                className="flex flex-col rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900"
              >
                <div className="flex items-start justify-between gap-3">
                  <dt className="text-sm text-gray-600 dark:text-gray-300">{metric.label}</dt>
                  <span
                    aria-hidden="true"
                    className={`flex size-9 shrink-0 items-center justify-center rounded-lg ${metric.tone}`}
                  >
                    <metric.Icon className="size-4.5" />
                  </span>
                </div>

                <dd className="mt-2">
                  <span className="block text-2xl font-semibold tabular-nums text-gray-900 dark:text-white">
                    {metric.value}
                  </span>

                  {/* The direction is a word. An arrow and a colour are not
                      text, so on their own this reads as a bare percentage. */}
                  <span
                    className={`mt-1.5 flex items-center gap-1 text-sm ${
                      rising
                        ? 'text-emerald-700 dark:text-emerald-400'
                        : 'text-rose-700 dark:text-rose-400'
                    }`}
                  >
                    {rising ? (
                      <ArrowUpRight aria-hidden="true" className="size-4" />
                    ) : (
                      <ArrowDownRight aria-hidden="true" className="size-4" />
                    )}
                    {rising ? 'Up' : 'Down'} {magnitude}%{' '}
                    {/* Explicit space: JSX drops the newline between an
                        expression and the next element, so without it the text
                        content reads "12.5%vs last month". The flex gap hides
                        that visually, which is what makes it easy to miss. */}
                    <span className="text-gray-500 dark:text-gray-400">{metric.comparison}</span>
                  </span>
                </dd>

                {/* mt-auto so the sparkline and link sit on one line across the
                    row whatever the label above them wraps to. */}
                <div className="mt-auto pt-4">
                  <Sparkline series={metric.series} rising={rising} />
                  <a
                    href={metric.href}
                    className="mt-3 inline-flex min-h-11 items-center text-sm font-medium text-indigo-700 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
                  >
                    View details
                    {/* Names the metric, so four links are four different
                        entries in a screen reader's link list. */}
                    <span className="sr-only"> for {metric.label}</span>
                  </a>
                </div>
              </div>
            )
          })}
        </dl>
      </div>
    </section>
  )
}
