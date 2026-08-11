import { Tag, TrendingUp, Zap } from 'lucide-react'

/*
 * Categories as a bento: a hero tile, two tall ones, and a run of squares.
 *
 * The spans tile the grid exactly, and that is arithmetic you have to redo
 * whenever the list changes. Across four columns this is 2×2 + 1×2 + 1×2 for
 * the first two rows (eight cells), four squares for the third, and two
 * double-width tiles for the fourth. Get it wrong and CSS grid drops tiles
 * into whatever hole it finds, which reads as a layout bug rather than a
 * missing span.
 *
 * Row height is set with `auto-rows` rather than left to the content. Tall
 * tiles span two rows, so without a defined row height the row sizes itself to
 * whatever happens to be in it and the "large" tile stops being large.
 *
 * Three things from the source deliberately dropped, all of them the same
 * mistakes as the sibling tile grid: the mount fade-in that leaves the section
 * invisible until JavaScript runs, the hovered-tile React state that
 * re-renders the grid to do what `group-hover` already does, and
 * `alt={category.name}` sitting directly above a heading with the same words.
 *
 * "Explore" appears on hover but is inside the link, so it is part of the
 * link's name either way — it is progressive decoration, not hidden content.
 */

export interface Category {
  name: string
  href: string
  image: string
  /** Grid span above lg. Must tile the grid — see the note above. */
  span: string
  caption?: string
  discount?: string
  isNew?: boolean
  isTrending?: boolean
  /** Overrides the default scrim. Keep the foot opaque or the label vanishes. */
  scrim?: string
}

const CATEGORIES: Category[] = [
  {
    name: 'Electronics',
    href: '#',
    image: 'https://images.unsplash.com/photo-1526738549149-8e07eca6c147?w=900&h=900&fit=crop&q=80',
    span: 'lg:col-span-2 lg:row-span-2',
    caption: 'The latest of everything',
    isTrending: true,
    scrim: 'from-blue-950/90 via-blue-900/40 to-transparent',
  },
  {
    name: 'Smartphones',
    href: '#',
    image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=600&h=900&fit=crop&q=80',
    span: 'lg:row-span-2',
    scrim: 'from-indigo-950/90 via-indigo-900/40 to-transparent',
  },
  {
    name: 'Wearables',
    href: '#',
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600&h=900&fit=crop&q=80',
    span: 'lg:row-span-2',
    caption: 'Watches and trackers',
    scrim: 'from-teal-950/90 via-teal-900/40 to-transparent',
  },
  {
    name: 'Laptops',
    href: '#',
    image: 'https://images.unsplash.com/photo-1593642632823-8f785ba67e45?w=600&h=600&fit=crop&q=80',
    span: '',
    discount: 'Up to 30% off',
  },
  {
    name: 'Headphones',
    href: '#',
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600&h=600&fit=crop&q=80',
    span: '',
    isNew: true,
    scrim: 'from-amber-950/90 via-amber-900/40 to-transparent',
  },
  {
    name: 'Cameras',
    href: '#',
    image: 'https://images.unsplash.com/photo-1516035069371-29a1b244cc32?w=600&h=600&fit=crop&q=80',
    span: '',
  },
  {
    name: 'Gaming',
    href: '#',
    image: 'https://images.unsplash.com/photo-1552820728-8b83bb6b773f?w=600&h=600&fit=crop&q=80',
    span: '',
    scrim: 'from-purple-950/90 via-purple-900/40 to-transparent',
  },
  {
    name: 'Home theatre',
    href: '#',
    image: 'https://images.unsplash.com/photo-1593784991095-a205069470b6?w=900&h=600&fit=crop&q=80',
    span: 'lg:col-span-2',
    caption: 'Screens, speakers and everything between',
  },
  {
    name: 'Desk & workspace',
    href: '#',
    image: 'https://images.unsplash.com/photo-1601524909162-ae8725290836?w=900&h=600&fit=crop&q=80',
    span: 'lg:col-span-2',
    caption: 'Monitors, stands and cable you will actually use',
    scrim: 'from-emerald-950/90 via-emerald-900/40 to-transparent',
  },
]

const DEFAULT_SCRIM = 'from-gray-950/90 via-gray-900/40 to-transparent'

export default function BentoCategoryGrid({
  title = 'Shop by category',
  subtitle = 'Everything we stock, grouped the way people actually look for it.',
  categories = CATEGORIES,
}: {
  title?: string
  subtitle?: string
  categories?: Category[]
}) {
  return (
    <section className="bg-gray-50 py-12 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 md:px-8">
        <div className="mx-auto mb-10 max-w-2xl text-center">
          <h2 className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white">
            {title}
          </h2>
          <p className="mt-3 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        {/* auto-rows is what makes a row-span-2 tile actually twice as tall.
            Without it the row sizes to its content and "large" means nothing. */}
        <ul
          role="list"
          className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6 lg:grid-cols-4 lg:auto-rows-[11rem]"
        >
          {categories.map((category) => (
            <li key={category.name} className={category.span}>
              <a
                href={category.href}
                className="group relative block size-full overflow-hidden rounded-2xl shadow-lg transition-shadow hover:shadow-xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
              >
                <div className="aspect-[4/3] lg:aspect-auto lg:size-full">
                  <img
                    src={category.image}
                    alt=""
                    className="size-full object-cover transition-transform duration-700 group-hover:scale-110 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                  />
                </div>

                <span
                  aria-hidden="true"
                  className={`absolute inset-0 bg-gradient-to-t ${category.scrim ?? DEFAULT_SCRIM}`}
                />

                <span className="absolute inset-0 flex flex-col justify-between p-5 md:p-6">
                  <span className="flex flex-wrap gap-2">
                    {category.isNew && (
                      <span className="inline-flex items-center gap-1 rounded-full bg-emerald-700 px-2.5 py-1 text-xs font-bold text-white shadow-md">
                        <Zap aria-hidden="true" className="size-3" />
                        New
                      </span>
                    )}
                    {category.isTrending && (
                      <span className="inline-flex items-center gap-1 rounded-full bg-gradient-to-r from-red-500 to-orange-500 px-2.5 py-1 text-xs font-bold text-white shadow-md">
                        <TrendingUp aria-hidden="true" className="size-3" />
                        Trending
                      </span>
                    )}
                  </span>

                  <span className="transition-transform duration-300 group-hover:-translate-y-1 motion-reduce:transition-none motion-reduce:group-hover:translate-y-0">
                    <span
                      className={`block font-bold text-white ${
                        category.span.includes('col-span-2') && category.span.includes('row-span-2')
                          ? 'text-3xl'
                          : category.span
                            ? 'text-2xl'
                            : 'text-xl'
                      }`}
                    >
                      {category.name}
                    </span>

                    {category.caption && (
                      <span className="mt-1 block text-sm text-white/80">{category.caption}</span>
                    )}

                    {category.discount && (
                      <span className="mt-2 inline-flex items-center gap-1.5 rounded-lg bg-white/20 px-3 py-1 text-sm text-white backdrop-blur-sm">
                        <Tag aria-hidden="true" className="size-3.5" />
                        {category.discount}
                      </span>
                    )}

                    <span
                      aria-hidden="true"
                      className="mt-3 block h-0.5 w-0 bg-white transition-all duration-500 ease-out group-hover:w-1/3 motion-reduce:transition-none"
                    />
                  </span>
                </span>
              </a>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
