'use client'

import { useState } from 'react'
import { Heart, Star } from 'lucide-react'

/*
 * A product grid where each card can be saved to a wishlist.
 *
 * The wishlist control is the interesting part, and it is the one most
 * implementations get wrong. It is a toggle, so it needs three things a plain
 * icon button does not have:
 *
 *   - `aria-pressed`, which is what makes a screen reader say "pressed" rather
 *     than leaving the state visible only as a filled heart;
 *   - a name that says which product it belongs to, because "Save" repeated
 *     six times is six identical controls;
 *   - a label that does not change with the state. "Save Hydrating Serum" with
 *     aria-pressed is correct; flipping the label to "Remove" as well tells the
 *     user the state changed twice.
 *
 * It is also `relative z-10`, so it sits above the stretched card link rather
 * than under it. Miss that and the heart is unclickable — the link covers it —
 * which looks like a broken button and is really a stacking bug.
 *
 * State is local here so the block works on its own. `onToggle` reports the
 * change, so wiring it to a real wishlist is one prop.
 *
 * Images are plain <img>, since next/image needs every remote host declared in
 * next.config and a block cannot do that from inside itself.
 */

export interface Product {
  id: string
  name: string
  brand?: string
  price: number
  image: string
  rating?: number
  reviews?: number
  href?: string
  /** Saved when the block first renders. */
  saved?: boolean
}

const CURRENCY = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

const PRODUCTS: Product[] = [
  {
    id: '1',
    name: 'Hydrating serum with hyaluronic acid',
    brand: 'PureGlow',
    price: 54.99,
    image: 'https://images.unsplash.com/photo-1571781926291-c477ebfd024b?w=600&h=600&fit=crop&q=80',
    rating: 4.8,
    reviews: 312,
    saved: true,
  },
  {
    id: '2',
    name: 'Purifying clay mask',
    brand: 'EarthTone',
    price: 34.99,
    image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=600&h=600&fit=crop&q=80',
    rating: 4.6,
    reviews: 189,
  },
  {
    id: '3',
    name: 'Daily moisturising cream',
    brand: 'PureGlow',
    price: 42.0,
    image: 'https://images.unsplash.com/photo-1620916566398-39f1143ab7be?w=600&h=600&fit=crop&q=80',
    rating: 4.4,
    reviews: 97,
  },
  {
    id: '4',
    name: 'Restorative night oil',
    brand: 'Lumen',
    price: 68.0,
    image: 'https://images.unsplash.com/photo-1608248543803-ba4f8c70ae0b?w=600&h=600&fit=crop&q=80',
    rating: 4.9,
    reviews: 143,
  },
]

function Rating({ rating, reviews }: { rating: number; reviews?: number }) {
  return (
    <p className="mt-1 flex items-center gap-1.5">
      <span aria-hidden="true" className="flex">
        {[0, 1, 2, 3, 4].map((i) => (
          <Star
            key={i}
            className={`size-3.5 ${
              i < Math.round(rating)
                ? 'fill-amber-400 text-amber-400'
                : 'fill-gray-200 text-gray-200 dark:fill-white/15 dark:text-white/15'
            }`}
          />
        ))}
      </span>
      <span className="text-xs text-gray-600 dark:text-gray-400">
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

export default function ProductGridWithWishlist({
  title = 'Skincare',
  viewAllLabel = 'View all',
  viewAllHref = '#',
  products = PRODUCTS,
  onToggle,
}: {
  title?: string
  viewAllLabel?: string
  viewAllHref?: string
  products?: Product[]
  onToggle?: (product: Product, saved: boolean) => void
}) {
  const [saved, setSaved] = useState<Record<string, boolean>>(() =>
    Object.fromEntries(products.map((p) => [p.id, Boolean(p.saved)])),
  )

  function toggle(product: Product) {
    const next = !saved[product.id]
    setSaved((current) => ({ ...current, [product.id]: next }))
    onToggle?.(product, next)
  }

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
              className="group relative flex flex-col rounded-xl border border-gray-200 bg-white p-4 has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:bg-gray-900"
            >
              {/* z-10 lifts the toggle above the stretched card link. Without
                  it the link covers the heart and it cannot be clicked. */}
              <button
                type="button"
                onClick={() => toggle(product)}
                aria-pressed={saved[product.id]}
                className="absolute top-2 right-2 z-10 inline-flex size-11 items-center justify-center rounded-full text-gray-400 hover:text-rose-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-500"
              >
                <Heart
                  aria-hidden="true"
                  className={`size-5 ${saved[product.id] ? 'fill-rose-500 text-rose-500' : ''}`}
                />
                {/* The label stays the same in both states; aria-pressed is
                    what reports the change. */}
                <span className="sr-only">Save {product.name} to wishlist</span>
              </button>

              <div className="aspect-square overflow-hidden rounded-lg bg-gray-50 dark:bg-white/5">
                <img
                  src={product.image}
                  alt=""
                  className="size-full object-cover transition-transform duration-300 group-hover:scale-105"
                />
              </div>

              <div className="mt-4 flex flex-1 flex-col">
                {product.brand && (
                  <p className="text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400">
                    {product.brand}
                  </p>
                )}
                <h3 className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                  <a
                    href={product.href ?? '#'}
                    className="after:absolute after:inset-0 focus:outline-none"
                  >
                    {product.name}
                  </a>
                </h3>
                {product.rating !== undefined && (
                  <Rating rating={product.rating} reviews={product.reviews} />
                )}
                <p className="mt-3 text-base font-semibold text-gray-900 tabular-nums dark:text-white">
                  {CURRENCY.format(product.price)}
                </p>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
