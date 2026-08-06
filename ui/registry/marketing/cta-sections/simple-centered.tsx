/*
 * The plainest CTA in the set: heading, one line of support, two actions.
 *
 * Worth saying out loud, because it is the one people skip past for something
 * with a gradient in it: this converts. There is exactly one thing to read and
 * exactly one thing to do, and nothing on the page competes with the button.
 *
 * The secondary action is a link rather than a second button. Two buttons of
 * equal weight is not a choice, it is a decision handed back to the visitor.
 */

export default function CtaSimpleCentered({
  title = 'Boost your productivity.\nStart using our app today.',
  subtitle = 'Ship the API, the admin panel and the typed client from one definition, then deploy the whole thing as a single binary.',
  primaryLabel = 'Get started',
  primaryHref = '#',
  secondaryLabel = 'Learn more',
  secondaryHref = '#',
}: {
  title?: string
  subtitle?: string
  primaryLabel?: string
  primaryHref?: string
  secondaryLabel?: string
  secondaryHref?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-2xl px-6 text-center lg:px-8">
        {/* whitespace-pre-line honours the \n in the default title, so the
            author controls the line break instead of hoping the measure does. */}
        <h2 className="text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>
        <p className="mx-auto mt-6 max-w-xl text-lg/8 text-pretty text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>
        <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-4">
          <a
            href={primaryHref}
            className="inline-flex min-h-11 items-center rounded-md bg-indigo-600 px-4 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            {primaryLabel}
          </a>
          <a
            href={secondaryHref}
            className="inline-flex min-h-11 items-center rounded-md px-1 text-sm font-semibold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-white"
          >
            {secondaryLabel}<span aria-hidden="true">&nbsp;&rarr;</span>
          </a>
        </div>
      </div>
    </section>
  )
}
