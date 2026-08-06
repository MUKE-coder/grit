/*
 * A dotted world map above three figures.
 *
 * The map is drawn from a coarse land mask rather than loaded as an image:
 * each row lists the column ranges that are land, and the component turns that
 * into circles. It weighs almost nothing, takes its colour from the theme, and
 * stays crisp at any size.
 *
 * The mask is repeated in this file rather than imported from the sibling
 * block that also uses it. The registry installs exactly one file per block,
 * so a shared import would land in a project with nothing to resolve.
 *
 * The map is `aria-hidden`. It says "everywhere", which the figures below
 * already say in words, and announcing "graphic" before every number is noise.
 * The figures are a <dl> with the value before the label, so the DOM order and
 * the visual order agree.
 */

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

function WorldMap() {
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
      /* The mask is 54 columns, so the on-screen dot spacing is this width
         divided by 54. Wider than about 40rem and it reads as a grid of
         circles rather than as a texture. */
      className="mx-auto w-full max-w-2xl fill-gray-300 dark:fill-white/15"
    >
      {dots}
    </svg>
  )
}

export interface Stat {
  value: string
  label: string
}

const STATS: Stat[] = [
  { value: '+85%', label: 'Conversion rate' },
  { value: '12K', label: 'Active users' },
  { value: '40%', label: 'Revenue growth' },
]

export default function StatsMapAboveThreeUp({
  stats = STATS,
}: {
  stats?: Stat[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <WorldMap />

        <dl className="mx-auto -mt-6 grid max-w-4xl grid-cols-1 divide-y divide-gray-200 rounded-2xl bg-white shadow-lg ring-1 ring-gray-900/5 sm:grid-cols-3 sm:divide-x sm:divide-y-0 dark:divide-white/10 dark:bg-gray-900 dark:ring-white/10">
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
