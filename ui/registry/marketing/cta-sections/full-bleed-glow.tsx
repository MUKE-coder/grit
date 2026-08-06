/*
 * The same dark treatment as dark-card-centered, run edge to edge.
 *
 * Reach for this when the CTA is the last thing on the page and sits directly
 * above the footer: a rounded card floating above a dark footer leaves a strip
 * of page colour between them that reads as a mistake. Full bleed just ends.
 *
 * The primary button here is a translucent white rather than solid. On a band
 * this large a solid white button is a hole punched in the design; at 10%
 * opacity over the glow it still clears contrast for its text while sitting in
 * the surface rather than on top of it. The visible border is what keeps it
 * legible as a button, so do not drop it.
 */

export default function CtaFullBleedGlow({
  title = 'Boost your productivity.\nStart using our app today.',
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
    <section className="relative isolate overflow-hidden bg-gray-950 px-6 py-28 text-center sm:py-36 lg:px-8">
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(56rem 30rem at 50% 120%, rgba(139,92,246,0.4), transparent 66%)',
        }}
      />

      <h2 className="mx-auto max-w-3xl text-4xl font-semibold tracking-tight text-balance whitespace-pre-line text-white sm:text-5xl">
        {title}
      </h2>
      <p className="mx-auto mt-6 max-w-xl text-lg/8 text-pretty text-gray-300">{subtitle}</p>
      <div className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-4">
        <a
          href={primaryHref}
          className="inline-flex min-h-11 items-center rounded-md border border-white/15 bg-white/10 px-4 text-sm font-semibold text-white hover:bg-white/15 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
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
    </section>
  )
}
