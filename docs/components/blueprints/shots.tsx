import type { Blueprint } from '@/config/blueprints'

/*
 * Three screenshots per blueprint, arranged so a card shows the application
 * rather than describing it.
 *
 * The pictures are real: every blueprint repository ships its own mockups for
 * its README, and scripts/blueprint-shots.py crops and resizes them into
 * public/blueprints. Nothing here is drawn in CSS to look like an app, because
 * a drawing of an app is the one screenshot that can never go out of date,
 * which is exactly the problem with it.
 *
 * Two arrangements, because the sources are two shapes. A phone blueprint
 * fans three portrait screens; a web one stacks three browser windows. Both
 * sit in the same 16:10 frame with the same wash behind them, so a row of
 * cards reads as one row.
 */

function Shot({
  src,
  alt,
  className,
  priority,
}: {
  src: string
  alt: string
  className: string
  priority?: boolean
}) {
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={src}
      alt={alt}
      {...(alt ? {} : { 'aria-hidden': true as const, 'data-decorative': '' })}
      loading={priority ? 'eager' : 'lazy'}
      decoding="async"
      className={className}
    />
  )
}

export function BlueprintShots({
  blueprint,
  className = '',
}: {
  blueprint: Blueprint
  className?: string
}) {
  const shots = blueprint.shots
  if (!shots || shots.length === 0) return null

  // The data is in reading order, because the row below renders it that way
  // and a caption list that starts in the middle reads oddly. The card then
  // picks its own front: the middle screen for a fan of phones, the first for
  // a stack of windows.
  const phone = blueprint.shotKind === 'phone'
  const front = phone ? shots[1] : shots[0]
  const left = phone ? shots[0] : shots[1]
  const right = shots[2]

  return (
    <div
      className={`relative isolate aspect-[16/10] overflow-hidden rounded-[12px] border border-border bg-surface-2 ${className}`.trim()}
    >
      {/* The wash: the same accent glow the landing sections use, kept behind
          everything so the screenshots stay the brightest thing in the frame. */}
      <div
        aria-hidden="true"
        data-decorative=""
        className="absolute inset-0 -z-10"
        style={{
          background:
            'radial-gradient(120% 90% at 50% 0%, color-mix(in oklab, var(--accent) 22%, transparent) 0%, transparent 62%)',
        }}
      />

      {phone ? (
        <>
          <Shot
            src={left.src}
            alt=""
            className="absolute bottom-[-12%] left-[6%] w-[31%] -rotate-[9deg] rounded-[10px] border border-border-strong shadow-[0_18px_40px_-18px_rgba(0,0,0,0.75)]"
          />
          <Shot
            src={right.src}
            alt=""
            className="absolute bottom-[-12%] right-[6%] w-[31%] rotate-[9deg] rounded-[10px] border border-border-strong shadow-[0_18px_40px_-18px_rgba(0,0,0,0.75)]"
          />
          <Shot
            src={front.src}
            alt={front.caption}
            priority
            className="absolute bottom-[-10%] left-1/2 w-[35%] -translate-x-1/2 rounded-[12px] border border-border-strong shadow-[0_24px_60px_-20px_rgba(0,0,0,0.85)]"
          />
        </>
      ) : (
        <>
          <Shot
            src={right.src}
            alt=""
            className="absolute left-[16%] right-[16%] top-[8%] rounded-[8px] border border-border opacity-70"
          />
          <Shot
            src={left.src}
            alt=""
            className="absolute left-[9%] right-[9%] top-[17%] rounded-[8px] border border-border-strong opacity-85"
          />
          <Shot
            src={front.src}
            alt={front.caption}
            priority
            className="absolute left-[2%] right-[2%] top-[28%] rounded-[10px] border border-border-strong shadow-[0_24px_60px_-20px_rgba(0,0,0,0.85)]"
          />
        </>
      )}
    </div>
  )
}

/**
 * The same pictures, laid out flat for a detail page: three in a row, each
 * captioned, nothing overlapping. A card wants a composition; a page that
 * somebody is reading before paying $10 wants to show the screens clearly.
 */
export function BlueprintShotRow({ blueprint }: { blueprint: Blueprint }) {
  const shots = blueprint.shots
  if (!shots || shots.length === 0) return null

  const phone = blueprint.shotKind === 'phone'

  return (
    <ul
      className={`m-0 grid list-none gap-5 p-0 ${
        phone
          ? 'mx-auto max-w-3xl grid-cols-2 sm:grid-cols-3'
          : 'grid-cols-1 sm:grid-cols-3'
      }`}
    >
      {shots.map((shot, i) => (
        <li key={shot.src} className="flex flex-col gap-3">
          <div className="overflow-hidden rounded-[10px] border border-border bg-surface-2">
            <Shot
              src={shot.src}
              alt={shot.caption}
              priority={i === 0}
              className="block w-full"
            />
          </div>
          <p className="m-0 text-body-sm text-fg-muted">{shot.caption}</p>
        </li>
      ))}
    </ul>
  )
}
