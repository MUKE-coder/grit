import {
  ArrowRight,
  Award,
  Blocks,
  Box,
  ChevronRight,
  Code2,
  Layers,
  PenTool,
  Shapes,
  Smile,
} from 'lucide-react'

/*
 * Category grid for a catalogue with nothing to photograph: icon, name, one
 * line of description, item count.
 *
 * Renders without JavaScript. The source faded the section in from a useEffect
 * and staggered each card with an inline transitionDelay, so the whole grid
 * started at opacity 0 and stayed there until React hydrated.
 *
 * It also tracked a hoveredCategory in state, setting it on every mouseenter
 * and clearing it on every mouseleave, and then never read the value. That is
 * a re-render of the whole grid per pointer move for nothing; the hover styling
 * it looks like it exists for is done by group-hover in CSS.
 *
 * The descriptions come back. They were in the data and their render was
 * commented out, along with the isNew and discount badges, so the block was
 * carrying three fields it never showed.
 *
 * Icons are component references rather than JSX stored in the data array. As
 * elements they cannot be serialised, cannot come from an API, and bake their
 * own sizing into the content.
 */

export interface Category {
  name: string
  href: string
  description: string
  itemCount: number
  Icon: typeof Box
  /** Tailwind text colour for the icon. */
  tone: string
}

const CATEGORIES: Category[] = [
  {
    name: '3D illustrations',
    href: '#',
    description: 'Ready-to-render scenes and objects',
    itemCount: 1240,
    Icon: Box,
    tone: 'text-violet-600',
  },
  {
    name: 'Illustration constructors',
    href: '#',
    description: 'Mix-and-match character and scene kits',
    itemCount: 486,
    Icon: Blocks,
    tone: 'text-sky-600',
  },
  {
    name: 'Vector illustrations',
    href: '#',
    description: 'Editable SVG sets at any size',
    itemCount: 3105,
    Icon: PenTool,
    tone: 'text-emerald-600',
  },
  {
    name: 'Mockups',
    href: '#',
    description: 'Device, print and packaging templates',
    itemCount: 872,
    Icon: Layers,
    tone: 'text-amber-600',
  },
  {
    name: 'Coded templates',
    href: '#',
    description: 'Production-ready pages and layouts',
    itemCount: 314,
    Icon: Code2,
    tone: 'text-rose-600',
  },
  {
    name: 'Web UI kits',
    href: '#',
    description: 'Component libraries with design files',
    itemCount: 529,
    Icon: Shapes,
    tone: 'text-indigo-600',
  },
  {
    name: 'Icons',
    href: '#',
    description: 'Line, solid and duotone families',
    itemCount: 8460,
    Icon: Smile,
    tone: 'text-teal-600',
  },
  {
    name: 'Fonts',
    href: '#',
    description: 'Typefaces and typography collections',
    itemCount: 967,
    Icon: Award,
    tone: 'text-orange-600',
  },
]

const counts = new Intl.NumberFormat('en-US')

export default function IconCardGridWithCounts({
  title = 'Top categories',
  subtitle = 'Browse the collection by what you are looking for.',
  seeAllHref = '#',
  categories = CATEGORIES,
}: {
  title?: string
  subtitle?: string
  seeAllHref?: string
  categories?: Category[]
}) {
  return (
    <section className="bg-white py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 md:px-8">
        <div className="mb-10 flex items-end justify-between gap-6">
          <div>
            <h2 className="text-2xl font-bold tracking-tight text-gray-900 md:text-3xl dark:text-white">
              {title}
            </h2>
            <p className="mt-2 text-gray-500 dark:text-gray-400">{subtitle}</p>
          </div>

          {/* Hidden rather than merely invisible below md, so the duplicate at
              the foot is not a second announcement of the same link. */}
          <a
            href={seeAllHref}
            className="group hidden shrink-0 items-center gap-1 text-sm font-medium text-gray-700 hover:text-black focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 md:flex dark:text-gray-300 dark:hover:text-white"
          >
            See all
            <ChevronRight
              aria-hidden="true"
              className="size-4 transition-transform group-hover:translate-x-1 motion-reduce:transition-none"
            />
          </a>
        </div>

        <ul role="list" className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
          {categories.map((category) => (
            <li key={category.name}>
              <a
                href={category.href}
                className="group flex h-full gap-5 rounded-xl border border-gray-100 bg-white p-6 transition-shadow hover:border-gray-200 hover:shadow-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/10 dark:bg-gray-900 dark:hover:border-white/20"
              >
                <span
                  aria-hidden="true"
                  className="flex size-16 shrink-0 self-start items-center justify-center rounded-xl bg-gray-50 transition-colors group-hover:bg-gray-100 dark:bg-white/5 dark:group-hover:bg-white/10"
                >
                  <category.Icon className={`size-8 ${category.tone}`} />
                </span>

                {/* Flex column so the count row can be pinned to the foot.
                    Descriptions wrap to one, two or three lines, and without
                    this the counts sit at three different heights across a
                    row of equal-height cards. */}
                <span className="flex min-w-0 flex-1 flex-col">
                  <span className="block font-semibold text-gray-900 dark:text-white">
                    {category.name}
                  </span>
                  <span className="mt-1 block text-sm text-gray-500 dark:text-gray-400">
                    {category.description}
                  </span>

                  <span className="mt-auto flex items-center justify-between gap-2 pt-3">
                    <span className="text-xs text-gray-400 dark:text-gray-500">
                      {counts.format(category.itemCount)} items
                    </span>
                    <ArrowRight
                      aria-hidden="true"
                      className="size-4 text-gray-400 transition-transform group-hover:translate-x-1 group-hover:text-gray-900 motion-reduce:transition-none dark:group-hover:text-white"
                    />
                  </span>
                </span>
              </a>
            </li>
          ))}
        </ul>

        <div className="mt-10 flex justify-center md:hidden">
          <a
            href={seeAllHref}
            className="inline-flex min-h-11 items-center gap-1 rounded-lg bg-gray-100 px-6 text-sm font-medium text-gray-800 hover:bg-gray-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:bg-white/10 dark:text-white dark:hover:bg-white/15"
          >
            See all categories
            <ChevronRight aria-hidden="true" className="size-4" />
          </a>
        </div>
      </div>
    </section>
  )
}
