'use client'

import { useCallback, useEffect, useRef, useState } from 'react'
import { ChevronLeft, ChevronRight, Heart, ShoppingBag, Sparkles, Star } from 'lucide-react'

/*
 * A horizontally scrolling row of products.
 *
 * The source combined two mistakes that are individually bad and together
 * fatal: the scroll container had no `tabIndex`, and its scrollbar was hidden
 * with `scrollbar-width: none`. A scroll region that cannot be focused cannot
 * be scrolled from the keyboard, and with the scrollbar gone there is no
 * visual affordance either — so everything past the third card was unreachable
 * without a mouse and invisible as a possibility. Both are fixed here: the
 * region is a labelled, focusable `region`, and the scrollbar is left alone.
 *
 * The arrows scroll by the measured width of a card rather than a hardcoded
 * 340px. A magic pixel number is wrong the moment the card is restyled or the
 * breakpoint changes, and it lands mid-card, which looks like a rendering bug.
 *
 * They also disable at each end. An arrow that stays lit and does nothing
 * teaches people the control is broken; `aria-disabled` plus a real `disabled`
 * keeps it out of the tab order once it is useless.
 *
 * Smooth scrolling is `motion-safe`. An unrequested animated scroll is exactly
 * what someone with vestibular sensitivity turned that preference on to avoid.
 */

export interface Product {
  id: string
  name: string
  category: string
  price: number
  originalPrice?: number
  image: string
  rating: number
  reviews: number
  isNew?: boolean
  isBestseller?: boolean
  href?: string
}

const CURRENCY = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
})

const PRODUCTS: Product[] = [
  {
    id: '1',
    name: 'Rainbow runner',
    category: 'Sustainable footwear',
    price: 749,
    originalPrice: 1249,
    image: 'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600&h=600&fit=crop&q=80',
    rating: 4.8,
    reviews: 234,
    isNew: true,
  },
  {
    id: '2',
    name: 'Court classic',
    category: 'Athletic wear',
    price: 899,
    originalPrice: 1199,
    image: 'https://images.unsplash.com/photo-1607522370275-f14206abe5d3?w=600&h=600&fit=crop&q=80',
    rating: 4.9,
    reviews: 456,
    isBestseller: true,
  },
  {
    id: '3',
    name: 'Everyday casual',
    category: 'Casual wear',
    price: 999,
    originalPrice: 1249,
    image: 'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?w=600&h=600&fit=crop&q=80',
    rating: 4.7,
    reviews: 189,
  },
  {
    id: '4',
    name: 'Performance runner',
    category: 'Running',
    price: 649,
    originalPrice: 849,
    image: 'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=600&h=600&fit=crop&q=80',
    rating: 4.6,
    reviews: 312,
  },
  {
    id: '5',
    name: 'Urban sneaker',
    category: 'Urban style',
    price: 799,
    originalPrice: 999,
    image: 'https://images.unsplash.com/photo-1605348532760-6753d2c43329?w=600&h=600&fit=crop&q=80',
    rating: 4.8,
    reviews: 278,
    isNew: true,
  },
  {
    id: '6',
    name: 'Oxford dress shoe',
    category: 'Formal',
    price: 1099,
    originalPrice: 1399,
    image: 'https://images.unsplash.com/photo-1614252235316-8c857d38b5f4?w=600&h=600&fit=crop&q=80',
    rating: 4.9,
    reviews: 167,
    isBestseller: true,
  },
]

export default function ScrollingProductCarousel({
  eyebrow = 'Premium collection',
  title = 'Trending products',
  subtitle = 'Recycled plastic, turned into shoes people actually want to wear.',
  products = PRODUCTS,
  onAdd,
}: {
  eyebrow?: string
  title?: string
  subtitle?: string
  products?: Product[]
  onAdd?: (product: Product) => void
}) {
  const track = useRef<HTMLUListElement>(null)
  const [saved, setSaved] = useState<Record<string, boolean>>({})
  const [atStart, setAtStart] = useState(true)
  const [atEnd, setAtEnd] = useState(false)

  const sync = useCallback(() => {
    const el = track.current
    if (!el) return
    setAtStart(el.scrollLeft <= 1)
    // 1px of slack: sub-pixel layout means scrollLeft rarely lands exactly on
    // the maximum, and a strict check leaves the arrow permanently enabled.
    setAtEnd(el.scrollLeft + el.clientWidth >= el.scrollWidth - 1)
  }, [])

  useEffect(() => {
    sync()
    const el = track.current
    if (!el) return
    const observer = new ResizeObserver(sync)
    observer.observe(el)
    return () => observer.disconnect()
  }, [sync])

  /* Scroll by a real card, measured now, rather than a number typed once. */
  function scrollByCard(direction: 1 | -1) {
    const el = track.current
    if (!el) return
    const card = el.querySelector('li')
    const step = card ? card.getBoundingClientRect().width + 32 : el.clientWidth * 0.8
    el.scrollBy({ left: direction * step, behavior: 'smooth' })
  }

  return (
    <section className="bg-gray-50 py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6">
        <div className="mx-auto max-w-3xl text-center">
          <p className="inline-flex items-center gap-2 rounded-full border border-indigo-200 bg-indigo-50 px-4 py-1.5 text-sm font-semibold text-indigo-700 dark:border-indigo-400/30 dark:bg-indigo-500/10 dark:text-indigo-300">
            <Sparkles aria-hidden="true" className="size-4" />
            {eyebrow}
          </p>
          <h2 className="mt-6 text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-4 text-lg/8 text-pretty text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <div className="relative mt-12">
          <button
            type="button"
            onClick={() => scrollByCard(-1)}
            disabled={atStart}
            className="absolute top-1/2 left-0 z-10 inline-flex size-12 -translate-x-2 -translate-y-1/2 items-center justify-center rounded-2xl border border-gray-200 bg-white/90 text-gray-700 shadow-lg backdrop-blur-sm transition hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:pointer-events-none disabled:opacity-0 dark:border-white/10 dark:bg-gray-900/90 dark:text-gray-200"
          >
            <ChevronLeft aria-hidden="true" className="size-6" />
            <span className="sr-only">Scroll products left</span>
          </button>

          {/* Focusable and labelled. Without tabIndex this content simply
              cannot be reached with a keyboard. The scrollbar is left visible
              on purpose: it is the only affordance a pointer user gets. */}
          <ul
            ref={track}
            onScroll={sync}
            role="region"
            aria-label={`${title}, scrollable`}
            tabIndex={0}
            className="flex snap-x snap-mandatory gap-8 overflow-x-auto py-8 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 motion-safe:scroll-smooth"
          >
            {products.map((product) => (
              <li
                key={product.id}
                className="group relative w-72 flex-none snap-start rounded-3xl border border-gray-200 bg-white p-4 shadow-sm transition-shadow hover:shadow-xl has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:bg-gray-900"
              >
                <div className="relative aspect-square overflow-hidden rounded-2xl bg-gray-50 dark:bg-white/5">
                  <img
                    src={product.image}
                    alt=""
                    className="size-full object-cover transition-transform duration-500 group-hover:scale-105 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                  />

                  <span className="absolute top-3 left-3 flex flex-col items-start gap-1">
                    {product.isNew && (
                      <span className="rounded-full bg-emerald-600 px-2.5 py-1 text-xs font-bold text-white">
                        New
                      </span>
                    )}
                    {product.isBestseller && (
                      <span className="rounded-full bg-amber-500 px-2.5 py-1 text-xs font-bold text-white">
                        Bestseller
                      </span>
                    )}
                  </span>

                  <button
                    type="button"
                    onClick={() => setSaved((s) => ({ ...s, [product.id]: !s[product.id] }))}
                    aria-pressed={Boolean(saved[product.id])}
                    className="absolute top-2 right-2 z-10 inline-flex size-11 items-center justify-center rounded-full bg-white/90 shadow backdrop-blur-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                  >
                    <Heart
                      aria-hidden="true"
                      className={`size-4 ${
                        saved[product.id] ? 'fill-rose-500 text-rose-500' : 'text-gray-700'
                      }`}
                    />
                    <span className="sr-only">Save {product.name} to wishlist</span>
                  </button>
                </div>

                <div className="mt-4">
                  <p className="text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400">
                    {product.category}
                  </p>
                  <h3 className="mt-1 text-base font-semibold text-gray-900 dark:text-white">
                    <a
                      href={product.href ?? '#'}
                      className="after:absolute after:inset-0 focus:outline-none"
                    >
                      {product.name}
                    </a>
                  </h3>

                  <p className="mt-1 flex items-center gap-1.5">
                    <span aria-hidden="true" className="flex">
                      {[0, 1, 2, 3, 4].map((i) => (
                        <Star
                          key={i}
                          className={`size-4 ${
                            i < Math.round(product.rating)
                              ? 'fill-amber-400 text-amber-400'
                              : 'fill-gray-200 text-gray-200 dark:fill-white/15 dark:text-white/15'
                          }`}
                        />
                      ))}
                    </span>
                    <span className="text-sm text-gray-600 dark:text-gray-400">
                      <span className="sr-only">Rated </span>
                      {product.rating.toFixed(1)}
                      <span className="sr-only"> out of 5, {product.reviews} reviews</span>
                      <span aria-hidden="true"> ({product.reviews})</span>
                    </span>
                  </p>

                  <button
                    type="button"
                    onClick={() => onAdd?.(product)}
                    className="relative z-10 mt-4 flex min-h-11 w-full items-center justify-between gap-3 rounded-2xl bg-gray-900 px-4 text-sm font-semibold text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
                  >
                    <span className="flex items-center gap-2">
                      <ShoppingBag aria-hidden="true" className="size-4" />
                      Add
                      <span className="sr-only">{product.name} to basket</span>
                    </span>
                    <span className="flex items-baseline gap-2">
                      <span className="text-lg font-bold tabular-nums">
                        {CURRENCY.format(product.price)}
                      </span>
                      {product.originalPrice && (
                        <span className="text-xs line-through opacity-60 tabular-nums">
                          <span className="sr-only">was </span>
                          {CURRENCY.format(product.originalPrice)}
                        </span>
                      )}
                    </span>
                  </button>
                </div>
              </li>
            ))}
          </ul>

          <button
            type="button"
            onClick={() => scrollByCard(1)}
            disabled={atEnd}
            className="absolute top-1/2 right-0 z-10 inline-flex size-12 translate-x-2 -translate-y-1/2 items-center justify-center rounded-2xl border border-gray-200 bg-white/90 text-gray-700 shadow-lg backdrop-blur-sm transition hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 disabled:pointer-events-none disabled:opacity-0 dark:border-white/10 dark:bg-gray-900/90 dark:text-gray-200"
          >
            <ChevronRight aria-hidden="true" className="size-6" />
            <span className="sr-only">Scroll products right</span>
          </button>
        </div>
      </div>
    </section>
  )
}
