/*
 * A grid for things that have no physical form: templates, kits, icon sets.
 *
 * Digital goods want a wide crop rather than a square one — you are showing a
 * screen, not an object — and they want the file type visible, because "Figma"
 * or "SVG" is the fact that decides whether the buyer can use it at all. That
 * belongs next to the price, not buried in a description.
 *
 * The badges sit on the image, so they need to survive whatever the image is.
 * `bg-white/90` with `backdrop-blur-sm` holds contrast over a dark photograph
 * and a light one; a flat translucent white does not, and a badge you cannot
 * read on half your catalogue is worse than no badge.
 *
 * One link per card, with `after:absolute after:inset-0` stretching the title's
 * hit area over the whole tile. Wrapping the entire card in an anchor instead
 * announces the name, the price and the category as one long link name.
 *
 * Images are plain <img> rather than next/image, which needs every remote host
 * declared in next.config — a block cannot fix that from inside itself, so it
 * would render broken in a project that has not added images.unsplash.com.
 */

export interface DigitalProduct {
  id: string
  name: string
  price: number
  image: string
  /** Shown under the name, e.g. "UI kit, Figma". */
  category?: string
  badges?: string[]
  href?: string
}

const CURRENCY = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

const PRODUCTS: DigitalProduct[] = [
  {
    id: '1',
    name: 'Premium UI kit',
    price: 79,
    image: 'https://images.unsplash.com/photo-1530435460869-d13625c69bbf?w=800&h=600&fit=crop&q=80',
    category: 'UI kit, Figma',
    badges: ['New', 'Trending'],
  },
  {
    id: '2',
    name: 'Landing page templates',
    price: 59,
    image: 'https://images.unsplash.com/photo-1522542550221-31fd19575a2d?w=800&h=600&fit=crop&q=80',
    category: 'Templates, HTML',
  },
  {
    id: '3',
    name: 'Illustration pack',
    price: 39,
    image: 'https://images.unsplash.com/photo-1545670723-196ed0954986?w=800&h=600&fit=crop&q=80',
    category: 'Graphics, SVG',
    badges: ['Featured'],
  },
  {
    id: '4',
    name: 'Icon collection',
    price: 29,
    image: 'https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=800&h=600&fit=crop&q=80',
    category: 'Icons, SVG',
    badges: ['50% off'],
  },
]

export default function DigitalGoodsGrid({
  title = 'Best selling digital items',
  viewAllLabel = 'View all',
  viewAllHref = '#',
  products = PRODUCTS,
}: {
  title?: string
  viewAllLabel?: string
  viewAllHref?: string
  products?: DigitalProduct[]
}) {
  return (
    <section className="bg-white py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="flex items-baseline justify-between gap-6">
          <h2 className="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">
            {title}
          </h2>
          <a
            href={viewAllHref}
            className="inline-flex min-h-11 items-center text-sm font-medium text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
          >
            {viewAllLabel}
            <span aria-hidden="true">&nbsp;&rarr;</span>
          </a>
        </div>

        <ul role="list" className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {products.map((product) => (
            <li
              key={product.id}
              className="group relative has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-4 has-[a:focus-visible]:outline-indigo-600"
            >
              <div className="relative aspect-[4/3] overflow-hidden rounded-xl bg-gray-100 dark:bg-white/5">
                <img
                  src={product.image}
                  alt=""
                  className="size-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
                {product.badges && product.badges.length > 0 && (
                  <ul role="list" className="absolute bottom-3 left-3 flex flex-wrap gap-1.5">
                    {product.badges.map((badge) => (
                      <li
                        key={badge}
                        /* Blur plus 90% white so it stays readable over a dark
                           photograph as well as a light one. */
                        className="rounded-full bg-white/90 px-2 py-1 text-xs font-medium text-gray-900 backdrop-blur-sm"
                      >
                        {badge}
                      </li>
                    ))}
                  </ul>
                )}
              </div>

              <div className="mt-3">
                <h3 className="text-base font-medium text-gray-900 dark:text-white">
                  <a
                    href={product.href ?? '#'}
                    className="after:absolute after:inset-0 focus:outline-none"
                  >
                    {product.name}
                  </a>
                </h3>
                <div className="mt-1 flex items-center justify-between gap-3">
                  <p className="font-semibold text-gray-900 tabular-nums dark:text-white">
                    {CURRENCY.format(product.price)}
                  </p>
                  {product.category && (
                    <p className="text-xs tracking-wider text-gray-500 uppercase dark:text-gray-400">
                      {product.category}
                    </p>
                  )}
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
