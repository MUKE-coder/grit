/*
 * The same ask as simple-centered, set left-aligned and without the paragraph.
 *
 * Use this one lower down a page. By the time someone has scrolled past the
 * features they have already read the pitch, and repeating it in a centred
 * block reads as though the page forgot it had made its case. A heading and a
 * button is enough.
 *
 * Left-aligned also means the heading shares an edge with everything above it,
 * so the eye does not have to re-find the start of the line.
 */

export default function CtaSimpleLeftAligned({
  title = 'Boost your productivity.\nStart using our app today.',
  primaryLabel = 'Get started',
  primaryHref = '#',
  secondaryLabel = 'Learn more',
  secondaryHref = '#',
}: {
  title?: string
  primaryLabel?: string
  primaryHref?: string
  secondaryLabel?: string
  secondaryHref?: string
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <h2 className="max-w-2xl text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>
        <div className="mt-10 flex flex-wrap items-center gap-x-6 gap-y-4">
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
