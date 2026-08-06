/*
 * Figures on the left, someone saying why they matter on the right.
 *
 * The pairing does something neither half does alone: the numbers are the
 * claim and the quote is the reason to believe it. A stats row on its own is
 * assertion, and a testimonial on its own is anecdote.
 *
 * The quote is a <figure> with a <blockquote> and a <figcaption>, which is the
 * markup that ties an attribution to the words it belongs to. A <p> in italics
 * followed by another <p> with a name in it looks the same and says nothing.
 *
 * The logo is text rather than an image, because a customer logo you cannot
 * name yet is a placeholder either way, and a broken <img> is a worse
 * placeholder than a word. Swap in a real mark when you have permission to use
 * one, and give it an alt of the company name.
 */

export interface Stat {
  value: string
  label: string
}

const STATS: Stat[] = [
  { value: '+1200', label: 'Stars on GitHub' },
  { value: '+500', label: 'Apps in production' },
]

export default function StatsWithQuote({
  body = 'From the CLI to the admin panel, everything here exists to get developers and businesses shipping sooner.',
  stats = STATS,
  quote = 'Adopting this was like unlocking a design department. We shipped an admin panel in an afternoon that would have taken us a fortnight.',
  name = 'Adam Wathan',
  role = 'CTO, Tailwind Labs',
  company = 'Tailwind Labs',
}: {
  body?: string
  stats?: Stat[]
  quote?: string
  name?: string
  role?: string
  company?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-2">
          <div>
            <p className="max-w-md text-lg/8 text-pretty text-gray-900 dark:text-white">{body}</p>
            <dl className="mt-10 flex gap-16">
              {stats.map((stat) => (
                <div key={stat.label}>
                  <dd className="text-4xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                    {stat.value}
                  </dd>
                  <dt className="mt-1 text-sm text-gray-600 dark:text-gray-400">{stat.label}</dt>
                </div>
              ))}
            </dl>
          </div>

          <figure className="border-l-2 border-indigo-500 pl-8 dark:border-indigo-400">
            <p className="text-sm font-semibold text-gray-900 dark:text-white">{company}</p>
            <blockquote className="mt-4 text-base/7 text-pretty text-gray-700 dark:text-gray-300">
              &ldquo;{quote}&rdquo;
            </blockquote>
            <figcaption className="mt-6 flex items-center gap-3">
              <span
                aria-hidden="true"
                className="flex size-9 flex-none items-center justify-center rounded-full bg-gray-100 text-xs font-semibold text-gray-600 dark:bg-white/10 dark:text-gray-300"
              >
                {name
                  .split(' ')
                  .map((part) => part[0])
                  .join('')}
              </span>
              <span className="text-sm text-gray-600 dark:text-gray-400">
                <span className="font-medium text-gray-900 dark:text-white">{name}</span>, {role}
              </span>
            </figcaption>
          </figure>
        </div>
      </div>
    </section>
  )
}
