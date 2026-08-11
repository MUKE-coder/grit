'use client'

import { useId, useMemo, useState } from 'react'
import { Heart, LayoutGrid, List, Search, SearchX, ShoppingCart, Star } from 'lucide-react'

/*
 * A product listing: search, category, sort, and a grid or list view.
 *
 * The result count is a live region. Filtering is the one interaction where
 * all the feedback lands somewhere other than the control you touched, so
 * typing in the search box and having nothing announced leaves a screen reader
 * user with no idea whether they narrowed it to forty results or none. The
 * count says what changed, and the empty state is inside the same region.
 *
 * The view toggle carries aria-pressed. The source gave both buttons an
 * aria-label and nothing else, so the pair announced identically whichever one
 * was active.
 *
 * Every control has a real label. Two selects and a search box relied on a
 * placeholder and their own first option, neither of which is a label.
 *
 * The "newest" sort is a date comparison. The source wrote
 * `a.isNew ? -1 : b.isNew ? 1 : 0`, which is not a valid comparator: it is not
 * antisymmetric, so the order it produces depends on the engine's sort
 * implementation and the array's starting order.
 *
 * Prices are integer cents. The source's floats meant every subtotal a
 * consumer of this data computed would drift.
 *
 * Dropped: a "Filters" button with no handler, and a hover-only image cycle.
 * A control that advances a gallery on mouseenter does not exist for touch or
 * keyboard, so the extra photographs were unreachable for most people.
 */

export interface Product {
  id: string
  name: string
  description: string
  category: string
  image: string
  /** Integer cents. */
  price: number
  /** Integer cents. Undefined when not discounted. */
  salePrice?: number
  rating: number
  reviews: number
  /** ISO date. Sorted on, not guessed from a boolean. */
  addedAt: string
  isNew?: boolean
}

const PRODUCTS: Product[] = [
  {
    id: 'weekender',
    name: 'Leather weekender bag',
    description: 'Full-grain leather, brass hardware, cotton lining.',
    category: 'Bags',
    image: 'https://images.unsplash.com/photo-1590874103328-eac38a683ce7?w=600&h=600&fit=crop&q=80',
    price: 24999,
    salePrice: 19999,
    rating: 4.8,
    reviews: 214,
    addedAt: '2026-05-02',
  },
  {
    id: 'headphones',
    name: 'Over-ear headphones',
    description: 'Active noise cancelling, 40 hours on a charge.',
    category: 'Audio',
    image: 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=600&h=600&fit=crop&q=80',
    price: 44900,
    salePrice: 37900,
    rating: 4.6,
    reviews: 1032,
    addedAt: '2026-06-18',
    isNew: true,
  },
  {
    id: 'watch',
    name: 'Automatic wristwatch',
    description: 'Sapphire crystal, exhibition case back, 100m.',
    category: 'Watches',
    image: 'https://images.unsplash.com/photo-1522312346375-d1a52e2b99b3?w=600&h=600&fit=crop&q=80',
    price: 89900,
    rating: 4.9,
    reviews: 87,
    addedAt: '2026-03-11',
  },
  {
    id: 'sneakers',
    name: 'Everyday sneakers',
    description: 'Leather upper, foam midsole, replaceable insole.',
    category: 'Footwear',
    image: 'https://images.unsplash.com/photo-1560769629-975ec94e6a86?w=600&h=600&fit=crop&q=80',
    price: 12900,
    salePrice: 9900,
    rating: 4.3,
    reviews: 640,
    addedAt: '2026-06-30',
    isNew: true,
  },
  {
    id: 'camera',
    name: 'Instant film camera',
    description: 'Fixed 35mm lens, built-in flash, takes pack film.',
    category: 'Cameras',
    image: 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=600&h=600&fit=crop&q=80',
    price: 15900,
    rating: 4.1,
    reviews: 298,
    addedAt: '2026-01-22',
  },
  {
    id: 'laptop',
    name: 'Studio laptop',
    description: '32GB memory, 2TB storage, colour-accurate display.',
    category: 'Computers',
    image: 'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=600&h=600&fit=crop&q=80',
    price: 349900,
    salePrice: 319900,
    rating: 4.7,
    reviews: 156,
    addedAt: '2026-04-09',
  },
  {
    id: 'skincare',
    name: 'Daily skincare set',
    description: 'Cleanser, serum and moisturiser, unscented.',
    category: 'Beauty',
    image: 'https://images.unsplash.com/photo-1583209814683-c023dd293cc6?w=600&h=600&fit=crop&q=80',
    price: 8400,
    rating: 4.4,
    reviews: 512,
    addedAt: '2026-06-05',
  },
  {
    id: 'dutch-oven',
    name: 'Enamelled dutch oven',
    description: 'Cast iron, 5.5 litres, oven safe to 260C.',
    category: 'Kitchen',
    image: 'https://images.unsplash.com/photo-1590794056226-79ef3a8147e1?w=600&h=600&fit=crop&q=80',
    price: 19900,
    salePrice: 14900,
    rating: 4.8,
    reviews: 421,
    addedAt: '2026-02-14',
  },
]

const SORTS = [
  { value: 'featured', label: 'Featured' },
  { value: 'price-low', label: 'Price: low to high' },
  { value: 'price-high', label: 'Price: high to low' },
  { value: 'rating', label: 'Top rated' },
  { value: 'newest', label: 'Newest' },
] as const

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const format = (cents: number) => money.format(cents / 100)
const effective = (product: Product) => product.salePrice ?? product.price

function Stars({ rating, reviews }: { rating: number; reviews: number }) {
  return (
    <p className="flex items-center gap-1">
      <span aria-hidden="true" className="flex">
        {[0, 1, 2, 3, 4].map((index) => (
          <Star
            key={index}
            className={`size-3.5 ${
              index < Math.round(rating)
                ? 'fill-amber-400 text-amber-400'
                : 'text-gray-300 dark:text-gray-600'
            }`}
          />
        ))}
      </span>
      {/* The stars are decoration; this is the actual rating, with its scale.
          "4.6" alone does not say out of what. */}
      <span className="text-xs text-gray-500 dark:text-gray-400">
        {rating} out of 5, {reviews} reviews
      </span>
    </p>
  )
}

function Price({ product }: { product: Product }) {
  return (
    <p className="flex items-baseline gap-2">
      <span className="font-semibold text-gray-900 dark:text-white">
        {format(effective(product))}
      </span>
      {product.salePrice !== undefined && (
        <span className="text-sm text-gray-500 dark:text-gray-400">
          <span className="sr-only">was </span>
          <s>{format(product.price)}</s>
        </span>
      )}
    </p>
  )
}

export default function ListingWithSortAndViewToggle({
  title = 'All products',
  products = PRODUCTS,
}: {
  title?: string
  products?: Product[]
}) {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState('All')
  const [sort, setSort] = useState<(typeof SORTS)[number]['value']>('featured')
  const [view, setView] = useState<'grid' | 'list'>('grid')
  const [favourites, setFavourites] = useState<string[]>([])

  const searchId = useId()
  const categoryId = useId()
  const sortId = useId()

  const categories = useMemo(
    () => ['All', ...Array.from(new Set(products.map((p) => p.category))).sort()],
    [products],
  )

  const results = useMemo(() => {
    const needle = query.trim().toLowerCase()
    const filtered = products.filter((product) => {
      const inCategory = category === 'All' || product.category === category
      const matches =
        !needle ||
        product.name.toLowerCase().includes(needle) ||
        product.description.toLowerCase().includes(needle) ||
        product.category.toLowerCase().includes(needle)
      return inCategory && matches
    })

    /* Sorted on a copy. Array.prototype.sort mutates, and sorting the source
       array here would reorder the caller's prop. */
    return [...filtered].sort((a, b) => {
      switch (sort) {
        case 'price-low':
          return effective(a) - effective(b)
        case 'price-high':
          return effective(b) - effective(a)
        case 'rating':
          return b.rating - a.rating
        case 'newest':
          return b.addedAt.localeCompare(a.addedAt)
        default:
          return 0
      }
    })
  }, [products, query, category, sort])

  function toggleFavourite(id: string) {
    setFavourites((current) =>
      current.includes(id) ? current.filter((item) => item !== id) : [...current, id],
    )
  }

  const controlClass =
    'min-h-11 rounded-lg border border-gray-300 bg-white px-3 text-sm text-gray-700 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-blue-600 dark:border-white/15 dark:bg-gray-900 dark:text-gray-200'

  return (
    <section className="bg-white py-10 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4">
        <div className="mb-8 rounded-2xl bg-gradient-to-r from-gray-900 to-gray-800 p-8">
          <h2 className="text-3xl font-bold text-white">{title}</h2>
          <form role="search" onSubmit={(event) => event.preventDefault()} className="mt-6 max-w-xl">
            <label htmlFor={searchId} className="sr-only">
              Search products
            </label>
            <div className="relative">
              <input
                id={searchId}
                type="search"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Search for products..."
                className="min-h-12 w-full rounded-xl bg-white/90 pr-12 pl-4 text-gray-800 placeholder-gray-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-500"
              />
              <Search
                aria-hidden="true"
                className="pointer-events-none absolute top-1/2 right-4 size-5 -translate-y-1/2 text-gray-500"
              />
            </div>
          </form>
        </div>

        <div className="mb-6 flex flex-col items-start justify-between gap-4 lg:flex-row lg:items-center">
          <div className="flex flex-wrap items-end gap-3">
            <div className="flex flex-col gap-1">
              <label htmlFor={categoryId} className="text-xs font-medium text-gray-600 dark:text-gray-400">
                Category
              </label>
              <select
                id={categoryId}
                value={category}
                onChange={(event) => setCategory(event.target.value)}
                className={`${controlClass} w-48`}
              >
                {categories.map((name) => (
                  <option key={name} value={name}>
                    {name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex flex-col gap-1">
              <label htmlFor={sortId} className="text-xs font-medium text-gray-600 dark:text-gray-400">
                Sort by
              </label>
              <select
                id={sortId}
                value={sort}
                onChange={(event) => setSort(event.target.value as typeof sort)}
                className={`${controlClass} w-48`}
              >
                {SORTS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="flex items-center gap-4">
            {/* Live, because this is where the result of every control above
                actually shows up. */}
            <p
              role="status"
              aria-live="polite"
              className="text-sm text-gray-500 dark:text-gray-400"
            >
              Showing{' '}
              <span className="font-medium text-gray-700 dark:text-gray-200">{results.length}</span>{' '}
              {results.length === 1 ? 'product' : 'products'}
            </p>

            <div className="flex overflow-hidden rounded-lg border border-gray-300 dark:border-white/15">
              {(
                [
                  { mode: 'grid', Icon: LayoutGrid, label: 'Grid view' },
                  { mode: 'list', Icon: List, label: 'List view' },
                ] as const
              ).map(({ mode, Icon, label }) => (
                <button
                  key={mode}
                  type="button"
                  onClick={() => setView(mode)}
                  /* aria-pressed, not just a label. Without it both buttons
                     announce the same whichever view is showing. */
                  aria-pressed={view === mode}
                  className={`inline-flex size-11 items-center justify-center focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-blue-600 ${
                    view === mode
                      ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'
                      : 'bg-white text-gray-700 hover:bg-gray-100 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-white/5'
                  }`}
                >
                  <Icon aria-hidden="true" className="size-4.5" />
                  <span className="sr-only">{label}</span>
                </button>
              ))}
            </div>
          </div>
        </div>

        {results.length === 0 ? (
          <div className="rounded-xl border border-dashed border-gray-300 py-16 text-center dark:border-white/15">
            <SearchX aria-hidden="true" className="mx-auto size-12 text-gray-400" />
            <p className="mt-4 text-lg font-medium text-gray-700 dark:text-gray-200">
              No products match that
            </p>
            <p className="mt-2 text-gray-500 dark:text-gray-400">
              Try a different search, or pick another category.
            </p>
          </div>
        ) : (
          <ul
            role="list"
            className={
              view === 'grid'
                ? 'grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4'
                : 'flex flex-col gap-6'
            }
          >
            {results.map((product) => (
              <li
                key={product.id}
                /* flex column in both views so the price row can be pinned
                   to the foot. Descriptions wrap to one or two lines, and the
                   grid stretches every card to the tallest, so without this
                   the Add buttons sit at different heights across a row. */
                className={`group relative flex rounded-xl border border-gray-200 bg-white transition-shadow hover:shadow-lg dark:border-white/10 dark:bg-gray-900 ${
                  view === 'list' ? 'gap-6 p-4' : 'flex-col overflow-hidden'
                }`}
              >
                <img
                  src={product.image}
                  alt=""
                  loading="lazy"
                  className={
                    view === 'list'
                      ? 'size-40 shrink-0 rounded-lg object-cover'
                      : 'aspect-square w-full object-cover'
                  }
                />

                <div className={`flex min-w-0 flex-1 flex-col ${view === 'list' ? '' : 'p-4'}`}>
                  <p className="text-xs tracking-wide text-gray-500 uppercase dark:text-gray-400">
                    {product.category}
                  </p>

                  <h3 className="mt-1 font-medium text-gray-900 dark:text-white">
                    {/* Stretched link: the whole card is the target, and the
                        heart above it stays clickable because of its z-index. */}
                    <a href="#" className="after:absolute after:inset-0 focus-visible:outline-none">
                      {product.name}
                    </a>
                  </h3>

                  <p
                    className={`mt-1 text-sm text-gray-600 dark:text-gray-400 ${
                      view === 'grid' ? 'line-clamp-2' : ''
                    }`}
                  >
                    {product.description}
                  </p>

                  <div className="mt-2">
                    <Stars rating={product.rating} reviews={product.reviews} />
                  </div>

                  <div className="mt-auto flex items-center justify-between gap-3 pt-3">
                    <Price product={product} />
                    <button
                      type="button"
                      className="relative z-10 inline-flex min-h-11 items-center gap-2 rounded-lg bg-gray-900 px-4 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
                    >
                      <ShoppingCart aria-hidden="true" className="size-4" />
                      <span className="sr-only">Add {product.name} to cart</span>
                      <span aria-hidden="true">Add</span>
                    </button>
                  </div>
                </div>

                {product.isNew && (
                  <p className="absolute top-3 left-3 z-10 rounded-full bg-emerald-700 px-2.5 py-1 text-xs font-semibold text-white">
                    New
                  </p>
                )}

                {/* z-10 keeps it above the stretched link, which otherwise
                    covers the whole card and swallows this. */}
                <button
                  type="button"
                  onClick={() => toggleFavourite(product.id)}
                  aria-pressed={favourites.includes(product.id)}
                  className="absolute top-3 right-3 z-10 inline-flex size-11 items-center justify-center rounded-full bg-white/90 text-gray-600 shadow-sm hover:text-rose-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:bg-gray-900/90 dark:text-gray-300"
                >
                  <Heart
                    aria-hidden="true"
                    className={`size-5 ${
                      favourites.includes(product.id) ? 'fill-rose-600 text-rose-600' : ''
                    }`}
                  />
                  <span className="sr-only">Save {product.name}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </section>
  )
}
