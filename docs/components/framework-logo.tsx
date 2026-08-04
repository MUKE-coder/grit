import Image from 'next/image'

/*
 * A framework's official mark, rendered legibly on the dark theme.
 *
 * Marks that are near-black on transparent (Next.js) disappear against
 * #0d0d12. Inverting them turns a black logo white, but it also mangles any
 * light detail the mark already has — the Next.js sticker outline comes out as
 * a dark smear. So those sit on a white chip instead: the logo is shown exactly
 * as its owner drew it, on the background it was drawn for. Nobody's trademark
 * gets recoloured by us.
 */

export function FrameworkLogo({
  src,
  alt,
  onLight,
  className = 'h-7',
}: {
  src: string
  alt: string
  /** the mark is dark and needs a light backing to be visible */
  onLight?: boolean
  /** height utility plus any extra classes; width is always auto */
  className?: string
}) {
  const img = (
    <Image
      src={src}
      alt={alt}
      width={200}
      height={64}
      className={`${className} w-auto object-contain`}
    />
  )

  if (!onLight) return img

  return (
    <span className="inline-flex items-center rounded-md bg-white px-1.5 py-1">{img}</span>
  )
}
