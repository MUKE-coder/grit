/*
 * Two figures and a sentence, in one row.
 *
 * The smallest stats block here, and the one that fits between two larger
 * sections without announcing itself. Use it when the numbers support a claim
 * made elsewhere on the page rather than being the claim.
 *
 * It is a <dl>. Each figure is a value with a label, which is exactly what a
 * description list is for, and it pairs them in the accessibility tree instead
 * of leaving a screen reader to work out which caption goes with which number
 * from their position on screen.
 *
 * The <dd> comes before the <dt> in the markup, which looks backwards and is
 * fine: a description list has no required order within a group, and putting
 * the number first means the visual order and the DOM order agree. Reversing
 * it with CSS would make them disagree, and the DOM order is the one a screen
 * reader follows.
 */

export interface Stat {
  value: string
  label: string
}

const STATS: Stat[] = [
  { value: '90+', label: 'Integrations' },
  { value: '56%', label: 'Less setup time' },
]

export default function StatsInlineWithCopy({
  body = 'The platform keeps growing with the developers and businesses building on it.',
  stats = STATS,
}: {
  body?: string
  stats?: Stat[]
}) {
  return (
    <section className="bg-white py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="flex flex-col items-start gap-8 sm:flex-row sm:items-center sm:gap-12">
          <dl className="flex gap-12">
            {stats.map((stat) => (
              <div key={stat.label}>
                <dd className="text-3xl font-semibold tracking-tight text-indigo-600 tabular-nums dark:text-indigo-400">
                  {stat.value}
                </dd>
                <dt className="mt-1 text-sm text-gray-600 dark:text-gray-400">{stat.label}</dt>
              </div>
            ))}
          </dl>

          <p className="max-w-md border-gray-200 text-base/7 text-pretty text-gray-600 sm:border-l sm:pl-12 dark:border-white/10 dark:text-gray-400">
            {body}
          </p>
        </div>
      </div>
    </section>
  )
}
