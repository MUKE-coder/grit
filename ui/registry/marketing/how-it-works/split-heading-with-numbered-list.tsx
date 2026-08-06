/*
 * Heading on the left, the steps as a plain numbered list on the right.
 *
 * The least decorated block in this set, and the one to use when the steps are
 * genuinely simple. Everything else here spends real estate on previews and
 * rails; if your three steps are one sentence each, that scaffolding is
 * packaging around nothing and readers can tell.
 *
 * The numbers come from CSS counters rather than being written into the
 * markup, so reordering the array renumbers the list and there is no chance of
 * a hand-written "3." sitting in position two. The counter output is
 * `aria-hidden` because the <ol> already conveys position — otherwise the
 * number is announced twice.
 */

export interface Step {
  title: string
  body: string
}

const STEPS: Step[] = [
  {
    title: 'Describe the resource',
    body: 'One definition names the fields and their types. It is the only place the shape of the thing is written down.',
  },
  {
    title: 'Generate the stack',
    body: 'The model, migration, service, handler, validation schema, typed client and admin screen all come from that definition.',
  },
  {
    title: 'Deploy the binary',
    body: 'One artifact carries the API, the admin panel and the assets. There is no runtime to install on the box.',
  },
]

export default function HowItWorksSplitHeadingWithNumberedList({
  eyebrow = 'How it works',
  title = 'Three steps from an idea to something deployed',
  subtitle = 'No scaffolding to wire together afterwards, and nothing generated that you cannot read.',
  ctaLabel = 'Read the quickstart',
  ctaHref = '#',
  steps = STEPS,
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  ctaLabel?: string
  ctaHref?: string
  steps?: Step[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="grid grid-cols-1 gap-x-16 gap-y-12 lg:grid-cols-2">
          <div>
            <p className="text-base font-semibold text-indigo-600 dark:text-indigo-400">
              {eyebrow}
            </p>
            <h2 className="mt-2 text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {title}
            </h2>
            <p className="mt-6 max-w-md text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
              {subtitle}
            </p>
            <a
              href={ctaHref}
              className="mt-8 inline-flex min-h-11 items-center rounded-lg px-1 text-sm font-semibold text-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
            >
              {ctaLabel}
              <span aria-hidden="true">&nbsp;&rarr;</span>
            </a>
          </div>

          {/* Counter-generated numbers: reorder the array and they follow. */}
          <ol role="list" className="space-y-10 [counter-reset:step]">
            {steps.map((step) => (
              <li
                key={step.title}
                className="relative pl-14 [counter-increment:step] before:absolute before:top-0 before:left-0 before:flex before:size-10 before:items-center before:justify-center before:rounded-full before:bg-gray-100 before:text-sm before:font-semibold before:text-gray-700 before:tabular-nums before:content-[counter(step)] dark:before:bg-white/10 dark:before:text-gray-300"
              >
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                  {step.title}
                </h3>
                <p className="mt-2 text-base/7 text-gray-600 dark:text-gray-400">{step.body}</p>
              </li>
            ))}
          </ol>
        </div>
      </div>
    </section>
  )
}
