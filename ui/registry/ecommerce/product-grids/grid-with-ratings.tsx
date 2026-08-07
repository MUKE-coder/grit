'use client'

import { ShoppingBag, Star } from 'lucide-react'

/*
 * A product grid with prices, ratings and an add-to-basket button.
 *
 * 'use client' is required and not optional: the add button takes an onClick,
 * and a server component cannot pass a function to one. That combination
 * compiles cleanly and then throws at render, which is the worst kind of
 * mistake to make in a block someone else installs. The sibling
 * digital-goods-grid has no handlers, so it stays a server component.
 *
 * Three things here are the difference between a grid that works and one that
 * only looks like it does.
 *
 * The star rating is not five icons. Five icons is what a sighted person sees;
 * everyone else gets nothing at all, on the element that most influences
 * whether they buy. The stars are `aria-hidden` and the rating is stated in
 * text beside them, so it is announced as "4.8 out of 5, 124 reviews".
 *
 * There is one link per card, not two. Making the image a link and the title
 * another link to the same place doubles every product in the tab order and
 * makes a screen reader read the name twice. The title is the link and
 * `after:absolute after:inset-0` stretches its hit area over the card, which
 * keeps the big click target without the duplicate.
 *
 * The button says "Add" on screen and "Add Urban Hiking Boot to basket" to a
 * screen reader. Four buttons all called "Add" is a list of identical controls
 * with no way to tell which is which.
 *
 * Images are plain <img>, not next/image. next/image needs every remote host
 * declared in next.config, so a block using it renders a broken image in a
 * project that has not added images.unsplash.com — and the block cannot fix
 * that from inside itself. The width, height and crop are in the query string
 * so the browser fetches roughly what it paints.
 */

export interface Product {
  id: string
  name: string
  brand?: string
  price: number
  /** Shown struck through beside the price. */
  originalPrice?: number
  image: string
  rating?: number
  reviews?: number
  href?: string
}

const CURRENCY = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

const PRODUCTS: Product[] = [
  {
    id: '1',
    name: 'Premium leather sneaker',
    brand: 'Luxe Step',
    price: 129.99,
    originalPrice: 179.99,
    image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600&h=600&fit=crop&q=80',
    rating: 4.8,
    reviews: 124,
  },
  {
    id: '2',
    name: 'Urban hiking boot',
    brand: 'TrailMaster',
    price: 149.95,
    image: 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=600&h=600&fit=crop&q=80',
    rating: 4.5,
    reviews: 89,
  },
  {
    id: '3',
    name: 'Performance running shoe',
    brand: 'SpeedRunner',
    price: 119.99,
    originalPrice: 159.99,
    image: 'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?w=600&h=600&fit=crop&q=80',
    rating: 4.7,
    reviews: 213,
  },
  {
    id: '4',
    name: 'Classic oxford dress shoe',
    brand: 'Elegance',
    price: 219.99,
    image: 'https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?w=600&h=600&fit=crop&q=80',
    rating: 4.9,
    reviews: 56,
  },
]

function Rating({ rating, reviews }: { rating: number; reviews?: number }) {
  return (
    <p className="mt-1 flex items-center gap-1.5">
      <span aria-hidden="true" className="flex">
        {[0, 1, 2, 3, 4].map((i) => (
          <Star
            key={i}
            className={`size-4 ${
              i < Math.round(rating)
                ? 'fill-amber-400 text-amber-400'
                : 'fill-gray-200 text-gray-200 dark:fill-white/15 dark:text-white/15'
            }`}
          />
        ))}
      </span>
      {/* The rating in words. The stars above are decoration. */}
      <span className="text-sm text-gray-600 dark:text-gray-400">
        <span className="sr-only">Rated </span>
        {rating.toFixed(1)}
        <span className="sr-only"> out of 5</span>
        {reviews !== undefined && (
          <>
            {' '}
            <span aria-hidden="true">({reviews})</span>
            <span className="sr-only">, {reviews} reviews</span>
          </>
        )}
      </span>
    </p>
  )
}

export default function ProductGridWithRatings({
  title = 'Featured products',
  viewAllLabel = 'View all',
  viewAllHref = '#',
  products = PRODUCTS,
  onAdd,
}: {
  title?: string
  viewAllLabel?: string
  viewAllHref?: string
  products?: Product[]
  onAdd?: (product: Product) => void
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

        <ul
          role="list"
          className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4"
        >
          {products.map((product) => (
            <li
              key={product.id}
              className="group relative flex flex-col overflow-hidden rounded-xl border border-gray-200 bg-white has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:bg-gray-900"
            >
              <div className="aspect-square overflow-hidden bg-gray-50 dark:bg-white/5">
                <img
                  src={product.image}
                  alt=""
                  className="size-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>

              <div className="flex flex-1 flex-col p-4">
                {product.brand && (
                  <p className="text-sm text-gray-500 dark:text-gray-400">{product.brand}</p>
                )}
                <h3 className="mt-0.5 text-base font-medium text-gray-900 dark:text-white">
                  {/* One link per card, with its hit area stretched over the
                      whole tile. */}
                  <a href={product.href ?? '#'} className="after:absolute after:inset-0 focus:outline-none">
                    {product.name}
                  </a>
                </h3>

                {product.rating !== undefined && (
                  <Rating rating={product.rating} reviews={product.reviews} />
                )}

                <div className="mt-4 flex items-center justify-between gap-3">
                  <p className="flex items-baseline gap-2">
                    <span className="text-lg font-semibold text-gray-900 tabular-nums dark:text-white">
                      {CURRENCY.format(product.price)}
                    </span>
                    {product.originalPrice && (
                      <span className="text-sm text-gray-400 line-through tabular-nums dark:text-gray-500">
                        <span className="sr-only">Was </span>
                        {CURRENCY.format(product.originalPrice)}
                      </span>
                    )}
                  </p>

                  {/* relative z-10 keeps this above the stretched link, so the
                      button is clickable rather than covered by it. */}
                  <button
                    type="button"
                    onClick={() => onAdd?.(product)}
                    className="relative z-10 inline-flex min-h-11 items-center gap-2 rounded-full bg-gray-900 px-4 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
                  >
                    <ShoppingBag aria-hidden="true" className="size-4" />
                    Add
                    <span className="sr-only"> {product.name} to basket</span>
                  </button>
                </div>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
