import { SCENES, scene } from './photos'

/*
 * Culture section: a claim, the numbers behind it, and photographs of the team
 * actually doing something.
 *
 * The photographs here are of moments rather than of people, which is why this
 * block takes `caption` and not `name`. A room, a whiteboard, a launch. Those
 * are much easier to source honestly than eight matched headshots, and they say
 * more about what it is like to work somewhere.
 *
 * The tilt is applied per item from a fixed list rather than at random. A random
 * rotation re-rolls on every render in development and moves under you while you
 * are trying to place something next to it.
 */

export interface Moment {
  caption: string
  photo?: string
}

export interface Stat {
  value: string
  label: string
}

const STATS: Stat[] = [
  { value: '42', label: 'Team members' },
  { value: '12', label: 'Countries' },
  { value: '4', label: 'Offices' },
]

const MOMENTS: Moment[] = [
  { caption: 'Strategy offsite, Austin 2024', photo: scene(SCENES[0]) },
  { caption: 'Product launch day, NYC 2024', photo: scene(SCENES[1]) },
  { caption: 'Annual retreat, Tokyo 2024', photo: scene(SCENES[3]) },
]

/* Fixed, deliberately uneven, and small. Past about three degrees a polaroid
   stops reading as casually placed and starts reading as broken CSS. */
const TILTS = ['-rotate-2', 'rotate-1', '-rotate-1', 'rotate-2']

export default function TeamStatsWithPolaroids({
  title = 'We build together',
  subtitle = 'Our team thrives on collaboration, creativity, and a shared commitment to excellence. Every day brings new challenges and opportunities to grow.',
  stats = STATS,
  moments = MOMENTS,
}: {
  title?: string
  subtitle?: string
  stats?: Stat[]
  moments?: Moment[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>

          <dl className="mt-12 grid grid-cols-2 gap-8 sm:grid-cols-3">
            {stats.map((stat) => (
              <div key={stat.label}>
                <dt className="sr-only">{stat.label}</dt>
                <dd>
                  <span className="block text-3xl font-semibold tracking-tight text-gray-900 tabular-nums dark:text-white">
                    {stat.value}
                  </span>
                  <span className="mt-1 block text-sm text-gray-500 dark:text-gray-400">
                    {stat.label}
                  </span>
                </dd>
              </div>
            ))}
          </dl>
        </div>

        <ul
          role="list"
          className="mt-20 grid grid-cols-1 gap-10 sm:grid-cols-2 lg:grid-cols-3"
        >
          {moments.map((moment, i) => (
            <li
              key={moment.caption}
              className={`${TILTS[i % TILTS.length]} bg-white p-3 pb-10 shadow-lg ring-1 ring-gray-900/5 dark:bg-gray-900 dark:ring-white/10`}
            >
              {moment.photo ? (
                <img
                  src={moment.photo}
                  alt=""
                  className="aspect-[4/3] w-full object-cover"
                />
              ) : (
                <div
                  aria-hidden="true"
                  className="flex aspect-[4/3] w-full items-center justify-center bg-gray-100 dark:bg-white/5"
                >
                  {/* Says what belongs here rather than pretending to be it. */}
                  <span className="px-4 text-center text-xs text-gray-400 dark:text-gray-500">
                    Your photo
                  </span>
                </div>
              )}
              <figcaption className="mt-4 text-sm text-gray-500 dark:text-gray-400">
                {moment.caption}
              </figcaption>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
