/**
 * Section — the landing page's structural shell.
 *
 * The design's crispness is structure that is drawn, not implied: a
 * full-bleed 1px rule above each section, a 1200px frame with 1px rails down
 * both sides running the whole page, and a 7x7 accent square where a rail
 * meets a rule. Those squares are the brand's signature detail, and the only
 * thing that makes them read as a signature rather than as decoration is that
 * every section has them.
 */

export function Section({
  children,
  id,
  /** A wash is one radial violet glow, anchored off a corner. At most one per
   *  section, and only on hero, batteries, benchmarks, deploy and the CTA. */
  wash,
  /** The 56px line grid. Hero, stats band and the closing CTA only. */
  grid = false,
  /** First section on the page: no rule above it, the header owns that edge. */
  first = false,
  className = '',
}: {
  children: React.ReactNode
  id?: string
  wash?: 'left' | 'right'
  grid?: boolean
  first?: boolean
  className?: string
}) {
  return (
    <section
      id={id}
      className={`relative ${first ? '' : 'border-t border-[#1a1d22]'} ${className}`.trim()}
    >
      {grid && (
        <div aria-hidden="true" className="g-gridbg pointer-events-none absolute inset-0" />
      )}
      {wash && (
        <div
          aria-hidden="true"
          className={`g-wash absolute h-[520px] w-[640px] ${
            wash === 'left' ? '-left-40' : '-right-40'
          } -top-40`}
        />
      )}

      <div className="g-frame px-6 py-16 md:py-24 lg:py-32">
        {/* The corner marks sit on the frame's own edges, so they land exactly
            where the rail crosses the rule above. */}
        <span className="g-cnr g-cnr--tl" aria-hidden="true" />
        <span className="g-cnr g-cnr--tr" aria-hidden="true" />
        {children}
      </div>
    </section>
  )
}

/**
 * SectionHeader — eyebrow, title, lead.
 *
 * The eyebrow is mono, uppercase and wrapped in braces: `{ BENCHMARKS }`.
 * It is the only uppercase on the site, which is what makes it read as a
 * label rather than as shouting.
 */
export function SectionHeader({
  eyebrow,
  title,
  lead,
  align = 'center',
  className = '',
}: {
  eyebrow?: string
  title: React.ReactNode
  lead?: React.ReactNode
  align?: 'center' | 'left'
  className?: string
}) {
  const centred = align === 'center'
  return (
    <header className={`${centred ? 'mx-auto max-w-3xl text-center' : 'max-w-3xl'} ${className}`.trim()}>
      {eyebrow && (
        <p className="g-eyebrow mb-4">
          {'{ '}
          {eyebrow}
          {' }'}
        </p>
      )}
      <h2 className="text-[clamp(2.25rem,5vw,3.5rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg">
        {title}
      </h2>
      {lead && (
        <p
          className={`mt-5 text-lead text-fg-muted ${centred ? 'mx-auto' : ''}`}
          style={{ maxWidth: '60ch' }}
        >
          {lead}
        </p>
      )}
    </header>
  )
}

/**
 * Highlight — one phrase per heading, in the violet gradient.
 *
 * Never a whole heading and never body copy: the gradient is an emphasis, and
 * emphasis applied evenly is no emphasis.
 */
export function Highlight({ children }: { children: React.ReactNode }) {
  return <span className="g-hl">{children}</span>
}
