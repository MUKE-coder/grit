import { ChevronRight, Star } from 'lucide-react'

/*
 * A category grid: photographs with the name laid over them.
 *
 * Three things the source of this block did that are worth not copying.
 *
 * It faded itself in on mount — every tile started at `opacity-0` and was
 * animated up by a `useEffect`. That means the entire section is invisible
 * until JavaScript runs, and stays invisible forever if it does not. Content
 * should be visible by default and animate as an enhancement, never the other
 * way round. This block renders complete on the server with no entrance
 * animation at all.
 *
 * It tracked the hovered tile in React state, which re-renders the whole grid
 * every time the pointer crosses a boundary, to do something `group-hover`
 * already does in CSS for free.
 *
 * Its images carried `alt={category.name}`, duplicating the heading directly
 * beneath them, so every tile was announced twice. The photographs here are
 * decorative — the name is right there in text — so their alt is empty.
 *
 * The overlay is a gradient with a solid foot rather than a flat wash, because
 * the text sits at the bottom. A uniform 50% tint either washes out the whole
 * photograph or leaves the text illegible over a bright one; a gradient that
 * is opaque where the words are and clear at the top does both jobs.
 *
 * This is a server component. There is no state left in it once the hover
 * tracking is gone.
 */

export interface Category {
  name: string
  href: string
  image: string
  /** Shown as a badge, e.g. "50% off". */
  discount?: string
  isNew?: boolean
}

const CATEGORIES: Category[] = [
  {
    name: 'Promotions',
    href: '#',
    image: 'https://images.unsplash.com/photo-1607083206968-13611e3d76db?w=500&h=500&fit=crop&q=80',
    discount: '50% off',
  },
  {
    name: 'Home & kitchen',
    href: '#',
    image: 'https://images.unsplash.com/photo-1556909114-f6e7ad7d3136?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Computing',
    href: '#',
    image: 'https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Phones',
    href: '#',
    image: 'https://images.unsplash.com/photo-1556656793-08538906a9f8?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Televisions',
    href: '#',
    image: 'https://images.unsplash.com/photo-1593784991095-a205069470b6?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Appliances',
    href: '#',
    image: 'https://images.unsplash.com/photo-1571175443880-49e1d25b2bc5?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Beauty',
    href: '#',
    image: 'https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9?w=500&h=500&fit=crop&q=80',
    isNew: true,
  },
  {
    name: 'Groceries',
    href: '#',
    image: 'https://images.unsplash.com/photo-1542838132-92c53300491e?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Sport & outdoor',
    href: '#',
    image: 'https://images.unsplash.com/photo-1517649763962-0c623066013b?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Bags & luggage',
    href: '#',
    image: 'https://images.unsplash.com/photo-1547949003-9792a18a2601?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Toys & games',
    href: '#',
    image: 'https://images.unsplash.com/photo-1558060370-d644479cb6f7?w=500&h=500&fit=crop&q=80',
  },
  {
    name: 'Handbags',
    href: '#',
    image: 'https://images.unsplash.com/photo-1590874103328-eac38a683ce7?w=500&h=500&fit=crop&q=80',
    isNew: true,
  },
]

export default function CategoryTileGridWithOverlays({
  title = 'Explore popular categories',
  subtitle = 'Everything we stock, sorted the way people actually shop.',
  viewAllLabel = 'View all categories',
  viewAllHref = '#',
  categories = CATEGORIES,
}: {
  title?: string
  subtitle?: string
  viewAllLabel?: string
  viewAllHref?: string
  categories?: Category[]
}) {
  return (
    <section className="bg-gray-50 py-12 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 md:px-8">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div className="flex items-start gap-3">
            <span
              aria-hidden="true"
              className="mt-1 h-8 w-1 flex-none rounded-full bg-gradient-to-b from-amber-400 to-amber-600"
            />
            <div>
              <h2 className="text-2xl font-bold tracking-tight text-gray-900 md:text-3xl dark:text-white">
                {title}
              </h2>
              <p className="mt-1 text-gray-500 dark:text-gray-400">{subtitle}</p>
            </div>
          </div>
          <a
            href={viewAllHref}
            className="hidden min-h-11 items-center text-sm font-medium text-amber-700 hover:text-amber-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600 md:inline-flex dark:text-amber-400 dark:hover:text-amber-300"
          >
            {viewAllLabel}
            <ChevronRight aria-hidden="true" className="ml-1 size-4" />
          </a>
        </div>

        <ul
          role="list"
          className="mt-8 grid grid-cols-2 gap-4 md:grid-cols-3 md:gap-5 lg:grid-cols-4 xl:grid-cols-6"
        >
          {categories.map((category) => (
            <li key={category.name}>
              <a
                href={category.href}
                className="group relative block overflow-hidden rounded-xl shadow-lg transition-shadow hover:shadow-xl focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600"
              >
                <div className="relative aspect-square overflow-hidden">
                  {/* Decorative: the name is in text directly below. */}
                  <img
                    src={category.image}
                    alt=""
                    className="size-full object-cover transition-transform duration-700 group-hover:scale-110 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                  />

                  {/* Opaque where the words are, clear at the top. A flat tint
                      either washes out the photograph or loses the text. */}
                  <span
                    aria-hidden="true"
                    className="absolute inset-0 bg-gradient-to-t from-gray-950/90 via-gray-950/40 to-transparent"
                  />

                  {category.discount && (
                    <span className="absolute top-3 right-3 rounded-full bg-white/90 px-2.5 py-1 text-xs font-bold text-red-600 shadow-md backdrop-blur-sm">
                      {category.discount}
                    </span>
                  )}
                  {category.isNew && !category.discount && (
                    <span className="absolute top-3 right-3 flex items-center gap-1 rounded-full bg-gradient-to-r from-amber-500 to-amber-600 px-2 py-1 text-xs font-bold text-white shadow-md">
                      <Star aria-hidden="true" className="size-3 fill-white" />
                      New
                    </span>
                  )}

                  <span className="absolute inset-x-0 bottom-0 p-4">
                    <span className="block text-sm leading-tight font-medium text-white md:text-base">
                      {category.name}
                    </span>
                    <span
                      aria-hidden="true"
                      className="mt-2 block h-0.5 w-0 bg-white transition-all duration-500 ease-out group-hover:w-1/2 motion-reduce:transition-none"
                    />
                  </span>
                </div>
              </a>
            </li>
          ))}
        </ul>

        <div className="mt-8 flex justify-center md:hidden">
          <a
            href={viewAllHref}
            className="inline-flex min-h-11 items-center rounded-lg border border-amber-200 bg-white px-6 text-sm font-medium text-amber-700 shadow-sm hover:bg-amber-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600 dark:border-amber-400/30 dark:bg-gray-900 dark:text-amber-400 dark:hover:bg-amber-400/10"
          >
            {viewAllLabel}
            <ChevronRight aria-hidden="true" className="ml-1 size-4" />
          </a>
        </div>
      </div>
    </section>
  )
}
