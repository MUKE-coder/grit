/**
 * WindowDots — the three macOS title-bar dots.
 *
 * This existed as five hand-written lines in 66 places, with the three colours
 * spelled as hex literals every time. That is 216 of the 284 hard-coded hex
 * values on the site, and it meant the dots could be a different size or a
 * different gap on any page without anybody noticing.
 *
 * The colours stay literal rather than becoming tokens on purpose: they are
 * not Grit's palette, they are macOS's, and a theme switch must not recolour
 * them. Keeping them in one file is the point, not theming them.
 */
/**
 * Two sizes, because the hand-written copies used two: 10px on a docs code
 * frame and 12px on the hero terminal. Anything else is a third size nobody
 * decided on.
 */
export function WindowDots({
  className = '',
  size = 'sm',
}: {
  className?: string
  size?: 'sm' | 'md'
}) {
  const dot = size === 'md' ? 'h-3 w-3' : 'h-2.5 w-2.5'
  return (
    <div className={`flex items-center gap-1.5 ${className}`} aria-hidden="true">
      <span className={`${dot} rounded-full bg-[#ff5f57]`} />
      <span className={`${dot} rounded-full bg-[#febc2e]`} />
      <span className={`${dot} rounded-full bg-[#28c840]`} />
    </div>
  )
}
