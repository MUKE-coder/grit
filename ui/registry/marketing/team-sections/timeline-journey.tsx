import { SCENES, scene } from './photos'

/*
 * The company story as an alternating timeline.
 *
 * Two things worth knowing about the markup.
 *
 * It is an ordered list, because it is a sequence and the order carries meaning.
 * A screen reader announces "3 of 4" and a keyboard user can move through it;
 * the same thing built from divs is a pile of unrelated paragraphs.
 *
 * The alternating left/right arrangement only exists above the `lg` breakpoint.
 * On a phone every entry stacks in one column in reading order, because a
 * zig-zag on a 375px screen is not a design, it is two very narrow columns.
 */

export interface Milestone {
  year: string
  title: string
  body: string
  caption?: string
  photo?: string
}

const MILESTONES: Milestone[] = [
  {
    year: '2021',
    title: 'The beginning',
    body: 'Three founders, one vision, and a cramped garage office.',
    caption: 'Where it all began, 2021',
    photo: scene(SCENES[0]),
  },
  {
    year: '2022',
    title: 'Growing pains',
    body: 'Moved into our first office. The team grew to fifteen.',
    caption: 'First real office, 2022',
    photo: scene(SCENES[4]),
  },
  {
    year: '2024',
    title: 'Today',
    body: 'Fifty or so people across twelve countries, still shipping every week.',
    caption: 'Global team, 2024',
    photo: scene(SCENES[2]),
  },
]

export default function TeamTimelineJourney({
  title = 'Our journey so far',
  subtitle = 'From a garage startup to a global team. Every moment has shaped who we are.',
  milestones = MILESTONES,
}: {
  title?: string
  subtitle?: string
  milestones?: Milestone[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="relative mx-auto mt-20 max-w-5xl">
          {/* The spine. Decorative, and only drawn where the alternating
              layout exists for it to run through. */}
          <div
            aria-hidden="true"
            className="absolute inset-y-0 left-1/2 hidden w-px -translate-x-1/2 bg-gray-200 lg:block dark:bg-white/10"
          />

          <ol role="list" className="space-y-16 lg:space-y-24">
            {milestones.map((milestone, i) => {
              const photoFirst = i % 2 === 0
              return (
                <li
                  key={milestone.year}
                  className="relative grid grid-cols-1 items-center gap-8 lg:grid-cols-2 lg:gap-16"
                >
                  {/* The node on the spine. */}
                  <span
                    aria-hidden="true"
                    className="absolute left-1/2 hidden size-4 -translate-x-1/2 rounded-full border-2 border-gray-300 bg-white lg:block dark:border-white/20 dark:bg-gray-950"
                  />

                  <figure
                    className={`${photoFirst ? 'lg:order-1' : 'lg:order-2'} ${
                      photoFirst ? '-rotate-1' : 'rotate-1'
                    } bg-white p-3 pb-9 shadow-lg ring-1 ring-gray-900/5 dark:bg-gray-900 dark:ring-white/10`}
                  >
                    {milestone.photo ? (
                      <img
                        src={milestone.photo}
                        alt=""
                        className="aspect-[4/3] w-full object-cover"
                      />
                    ) : (
                      <div
                        aria-hidden="true"
                        className="flex aspect-[4/3] w-full items-center justify-center bg-gray-100 dark:bg-white/5"
                      >
                        <span className="text-xs text-gray-400 dark:text-gray-500">
                          Your photo
                        </span>
                      </div>
                    )}
                    {milestone.caption && (
                      <figcaption className="mt-3 text-sm text-gray-500 dark:text-gray-400">
                        {milestone.caption}
                      </figcaption>
                    )}
                  </figure>

                  <div
                    className={`${photoFirst ? 'lg:order-2 lg:pl-8' : 'lg:order-1 lg:pr-8 lg:text-right'}`}
                  >
                    <p className="text-sm font-semibold text-indigo-600 tabular-nums dark:text-indigo-400">
                      {milestone.year}
                    </p>
                    <h3 className="mt-2 text-xl font-semibold text-gray-900 dark:text-white">
                      {milestone.title}
                    </h3>
                    <p className="mt-2 text-base/7 text-gray-600 dark:text-gray-400">
                      {milestone.body}
                    </p>
                  </div>
                </li>
              )
            })}
          </ol>
        </div>
      </div>
    </section>
  )
}
