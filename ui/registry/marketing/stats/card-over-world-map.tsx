/*
 * A stat card sitting on a dotted world map.
 *
 * The map is drawn from a coarse land mask rather than loaded as an image or
 * pasted in as a 200KB SVG path. Each row lists the column ranges that are
 * land, and the component turns that into circles. It weighs almost nothing,
 * takes its colour from the theme, and stays crisp at any size — none of which
 * is true of a PNG.
 *
 * It is deliberately low resolution. A dot map is a texture that says
 * "everywhere", not a reference you look countries up in, and detail past this
 * point only makes it noisier at the size it is actually rendered.
 *
 * The whole map is `aria-hidden`. It carries no information the stats beside
 * it do not already state, and "graphic" announced before every figure is
 * noise.
 *
 * The figures are a <dl>: each one is a value with a label, which is what a
 * description list is for, and it pairs them in the accessibility tree instead
 * of leaving a screen reader to guess which caption belongs to which number.
 */

/* Column ranges that are land, one entry per row of the mask.
   The grid is 54 columns of longitude by 26 rows of latitude. */
const LAND: [number, number][][] = [
  [[8, 14], [17, 19], [33, 48]],
  [[3, 4], [7, 15], [17, 19], [26, 28], [31, 49]],
  [[2, 5], [7, 16], [17, 19], [25, 28], [30, 50]],
  [[2, 6], [7, 16], [18, 19], [25, 28], [30, 50]],
  [[3, 6], [8, 16], [25, 31], [32, 50]],
  [[4, 7], [8, 16], [26, 30], [32, 49]],
  [[5, 8], [9, 16], [26, 31], [33, 48]],
  [[6, 9], [10, 16], [26, 32], [34, 47]],
  [[7, 10], [11, 15], [26, 33], [35, 46]],
  [[8, 11], [12, 15], [26, 34], [36, 45]],
  [[9, 12], [13, 15], [27, 34], [37, 44]],
  [[11, 14], [26, 35], [38, 43]],
  [[12, 15], [26, 36], [38, 42]],
  [[14, 16], [27, 36], [39, 42], [46, 48]],
  [[16, 21], [28, 36], [40, 41], [45, 49]],
  [[18, 22], [29, 35], [45, 50]],
  [[18, 23], [29, 35], [46, 50]],
  [[18, 23], [29, 35], [47, 50]],
  [[18, 23], [29, 35], [45, 51]],
  [[18, 23], [29, 34], [45, 51]],
  [[19, 23], [30, 34], [45, 51]],
  [[19, 22], [30, 33], [46, 50]],
  [[19, 22], [31, 33], [47, 49]],
  [[19, 21], [48, 49]],
  [[19, 20]],
  [[19, 19]],
]

const COLS = 54
const GAP = 8
const DOT = 1.6

function WorldMap({ className = '' }: { className?: string }) {
  const dots: React.JSX.Element[] = []
  LAND.forEach((ranges, row) => {
    ranges.forEach(([from, to]) => {
      for (let col = from; col <= to; col++) {
        dots.push(
          <circle key={`${row}-${col}`} cx={col * GAP + GAP / 2} cy={row * GAP + GAP / 2} r={DOT} />,
        )
      }
    })
  })

  return (
    <svg
      aria-hidden="true"
      viewBox={`0 0 ${COLS * GAP} ${LAND.length * GAP}`}
      className={`fill-gray-300 dark:fill-white/15 ${className}`}
    >
      {dots}
    </svg>
  )
}

export interface Stat {
  value: string
  unit?: string
  label: string
  body: string
}

const STATS: Stat[] = [
  { value: '99.9', unit: '%', label: 'Uptime guarantee', body: 'across every region we run in.' },
  { value: '24/7', label: 'Support', body: 'from engineers who work on the product.' },
  { value: '12', unit: '×', label: 'Faster processing', body: 'than the previous generation.' },
  { value: '90', unit: '+', label: 'Integrations', body: 'with the tools you already use.' },
]

export default function StatsCardOverWorldMap({
  title = 'Building the next generation of developer tooling',
  subtitle = 'The platform runs in fourteen regions and has never lost a customer record.',
  stats = STATS,
}: {
  title?: string
  subtitle?: string
  stats?: Stat[]
}) {
  return (
    <section className="relative isolate overflow-hidden bg-white py-24 sm:py-32 dark:bg-gray-950">
      {/* Width is not cosmetic here: the mask is 54 columns, so the on-screen
          gap between dots is the rendered width divided by 54. Much wider than
          this and it stops reading as a texture and starts reading as a grid
          of circles. */}
      <WorldMap className="pointer-events-none absolute top-1/2 right-0 -z-10 w-[46rem] -translate-y-1/2 opacity-80 lg:right-8" />

      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-xl rounded-3xl bg-white p-8 shadow-xl ring-1 ring-gray-900/5 sm:p-10 dark:bg-gray-900 dark:ring-white/10">
          <h2 className="text-3xl font-semibold tracking-tight text-balance text-gray-900 sm:text-4xl dark:text-white">
            {title}
          </h2>
          <p className="mt-4 text-base/7 text-pretty text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>

          <dl className="mt-10 grid grid-cols-1 gap-3 sm:grid-cols-2">
            {stats.map((stat) => (
              <div
                key={stat.label}
                className="rounded-xl bg-gray-50 p-5 dark:bg-white/5"
              >
                <dd className="text-3xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                  {stat.value}
                  {stat.unit && (
                    <span className="ml-1 text-lg font-normal text-gray-500 dark:text-gray-400">
                      {stat.unit}
                    </span>
                  )}
                </dd>
                <dt className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">
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
