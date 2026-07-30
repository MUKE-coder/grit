/**
 * The Grit UI mark.
 *
 * A rounded tile with the G knocked out of it — one path, one fill, evenodd.
 * That construction is deliberate: the glyph is transparent rather than
 * painted, so the mark is a single `currentColor` shape that takes the colour
 * of whatever text context it sits in and shows the surface behind through the
 * letter. It cannot turn into a solid blob the way a two-colour version does
 * when the tile and the glyph resolve to the same colour.
 *
 * Legible down to 16px, which is the size that actually decides whether a mark
 * works.
 */
export function GritUIMark({
  size = 24,
  className = '',
}: {
  size?: number
  className?: string
}) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 32 32"
      fill="none"
      aria-hidden="true"
      className={className}
    >
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

/** Mark plus wordmark, as used in the site header. */
export function GritUILogo({
  size = 24,
  className = '',
}: {
  size?: number
  className?: string
}) {
  return (
    <span className={`inline-flex items-center gap-2 ${className}`}>
      <GritUIMark size={size} />
      <span className="text-sm font-semibold tracking-tight">
        grit<span className="text-indigo-600 dark:text-indigo-400">UI</span>
      </span>
    </span>
  )
}
