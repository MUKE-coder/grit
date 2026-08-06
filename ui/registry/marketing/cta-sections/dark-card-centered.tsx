/*
 * A dark card dropped into a light page, with a violet glow rising from the
 * lower edge.
 *
 * Two notes on the glow. It is a radial gradient on a positioned, aria-hidden
 * element rather than a background image, so there is nothing to download and
 * nothing to go blurry on a retina screen. And it sits behind the content on
 * its own layer via `isolate` plus `-z-10`, which keeps the stacking context
 * local: drop this card next to a sticky header and the glow will not climb
 * over it.
 *
 * The card is dark in both themes on purpose. A dark panel that turns light in
 * dark mode stops being the accent it was chosen to be.
 */

export default function CtaDarkCardCentered({
  title = 'Boost your productivity today',
  subtitle = 'Generate the API, the admin panel and the typed client from one definition. Deploy the whole thing as a single binary.',
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
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="relative isolate overflow-hidden rounded-3xl bg-gray-950 px-6 py-24 text-center shadow-2xl sm:px-16 sm:py-32">
          <div
            aria-hidden="true"
            className="pointer-events-none absolute inset-0 -z-10"
            style={{
              background:
                'radial-gradient(42rem 26rem at 50% 118%, rgba(139,92,246,0.45), transparent 68%)',
            }}
          />

          <h2 className="mx-auto max-w-2xl text-4xl font-semibold tracking-tight text-balance text-white sm:text-5xl">
            {title}
          </h2>
          <p className="mx-auto mt-6 max-w-xl text-lg/8 text-pretty text-gray-300">{subtitle}</p>
          <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-4">
            <a
              href={primaryHref}
              className="inline-flex min-h-11 items-center rounded-md bg-white px-4 text-sm font-semibold text-gray-900 shadow-sm hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
            >
              {primaryLabel}
            </a>
            <a
              href={secondaryHref}
              className="inline-flex min-h-11 items-center rounded-md px-1 text-sm font-semibold text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
            >
              {secondaryLabel}<span aria-hidden="true">&nbsp;&rarr;</span>
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
