import { SCENES, scene } from './photos'

/*
 * A statement about the team beside a bento collage of them at work.
 *
 * The grid is deliberately uneven: one tall frame, one wide one, two small.
 * Four equal squares reads as a contact sheet, and the eye has nowhere to land.
 * Giving one frame more area makes it the subject and the rest supporting.
 *
 * The tall frame keeps its aspect ratio through `row-span-2` rather than a
 * fixed height, so the block survives being dropped into a narrower container
 * than the one it was designed in.
 */

export interface CollageItem {
  /** Described for a screen reader. Say what is happening, not "team photo". */
  alt: string
  photo?: string
}

const ITEMS: CollageItem[] = [
  { alt: 'Two engineers pairing at a laptop in the studio', photo: scene(SCENES[5], 600) },
  { alt: 'The team on the annual hike, taking a group photo', photo: scene(SCENES[3], 1200) },
  { alt: 'A product review in the meeting room', photo: scene(SCENES[1], 600) },
  { alt: 'The open-plan office on a working afternoon', photo: scene(SCENES[6], 600) },
]

export default function TeamBentoCollage({
  title = 'Built by developers. Backed by experience.',
  subtitle = "Our team has built products you already use. Now we are changing how businesses connect with their customers.",
  items = ITEMS,
}: {
  title?: string
  subtitle?: string
  items?: CollageItem[]
}) {
  const [tall, wide, ...rest] = items

  const frame =
    'overflow-hidden rounded-2xl bg-gray-100 ring-1 ring-gray-900/5 dark:bg-white/5 dark:ring-white/10'

  const Placeholder = () => (
    <div
      aria-hidden="true"
      className="flex size-full min-h-32 items-center justify-center"
    >
      <span className="text-xs text-gray-400 dark:text-gray-500">Your photo</span>
    </div>
  )

  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="mt-16 grid grid-cols-1 gap-5 sm:grid-cols-3 sm:grid-rows-2">
          {tall && (
            <figure className={`${frame} sm:row-span-2`}>
              {tall.photo ? (
                <img src={tall.photo} alt={tall.alt} className="size-full object-cover" />
              ) : (
                <Placeholder />
              )}
            </figure>
          )}

          {wide && (
            <figure className={`${frame} sm:col-span-2`}>
              {wide.photo ? (
                <img src={wide.photo} alt={wide.alt} className="size-full object-cover" />
              ) : (
                <Placeholder />
              )}
            </figure>
          )}

          {rest.map((item) => (
            <figure key={item.alt} className={frame}>
              {item.photo ? (
                <img src={item.photo} alt={item.alt} className="size-full object-cover" />
              ) : (
                <Placeholder />
              )}
            </figure>
          ))}
        </div>
      </div>
    </section>
  )
}
