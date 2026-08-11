'use client'

import { useState } from 'react'
import { Eye, Flame, Heart, ShoppingBag, Timer, Zap } from 'lucide-react'

/*
 * A deals grid: discount, what you save, how many are left, and a countdown.
 *
 * The hover actions are the part worth reading before you copy this. They fade
 * in on hover, which is fine, but they also fade in on `group-focus-within` —
 * without that they stay at `opacity-0` while remaining focusable, so a
 * keyboard user tabs into a button they cannot see, presses it, and something
 * happens for no visible reason. Invisible-but-focusable is the most common
 * hover-reveal bug and it is one class to fix.
 *
 * The stock meter is a real `progressbar` with `aria-valuenow`, `aria-valuemax`
 * and a text label, rather than a coloured div whose width encodes the number.
 * A bar with no value announces nothing, which on an "only 4 left" signal is
 * the whole message lost. The width comes from sold against a declared total,
 * so the proportion means something — the original divided by a hardcoded 150.
 *
 * Urgency copy is data, not decoration: `timeLeft` and `stockLeft` are props,
 * because a countdown that is hardcoded in markup is a lie the moment it ships.
 *
 * Images are plain <img> with verified remote URLs. The source of this block
 * pointed at local files in its own public folder, which render as broken
 * images in any project that does not happen to have them.
 */

export interface Deal {
  id: string
  name: string
  image: string
  currentPrice: number
  originalPrice: number
  /** Units sold and the total in the drop, which is what the meter shows. */
  sold: number
  total: number
  timeLeft?: string
  isNew?: boolean
  href?: string
}

const CURRENCY = new Intl.NumberFormat('en-NG', {
  style: 'currency',
  currency: 'NGN',
  maximumFractionDigits: 0,
})

const DEALS: Deal[] = [
  {
    id: '1',
    name: 'Wireless over-ear headphones, active noise cancelling',
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=600&h=450&fit=crop&q=80',
    currentPrice: 84000,
    originalPrice: 140000,
    sold: 128,
    total: 150,
    timeLeft: '4h 12m',
    isNew: true,
  },
  {
    id: '2',
    name: 'Smart watch with heart rate and sleep tracking',
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=600&h=450&fit=crop&q=80',
    currentPrice: 115000,
    originalPrice: 189000,
    sold: 41,
    total: 120,
    timeLeft: '6h 40m',
  },
  {
    id: '3',
    name: 'Wireless keyboard, low-profile keys',
    image: 'https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=600&h=450&fit=crop&q=80',
    currentPrice: 62000,
    originalPrice: 95000,
    sold: 64,
    total: 200,
    timeLeft: '1d 3h',
  },
  {
    id: '4',
    name: '13-inch ultrabook, 16GB RAM',
    image: 'https://images.unsplash.com/photo-1593642632823-8f785ba67e45?w=600&h=450&fit=crop&q=80',
    currentPrice: 940000,
    originalPrice: 1250000,
    sold: 187,
    total: 200,
    timeLeft: '2h 05m',
    isNew: true,
  },
  {
    id: '5',
    name: 'Instant camera with 20-shot film pack',
    image: 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=600&h=450&fit=crop&q=80',
    currentPrice: 96000,
    originalPrice: 145000,
    sold: 12,
    total: 40,
  },
  {
    id: '6',
    name: 'Mid-tower gaming PC, RGB cooling',
    image: 'https://images.unsplash.com/photo-1587202372775-e229f172b9d7?w=600&h=450&fit=crop&q=80',
    currentPrice: 1450000,
    originalPrice: 1980000,
    sold: 96,
    total: 100,
    timeLeft: '9h 22m',
  },
]

export default function DealsGridWithStockMeter({
  title = "Today's deals",
  subtitle = 'Limited time offers',
  viewAllLabel = 'See all',
  viewAllHref = '#',
  deals = DEALS,
  onAdd,
}: {
  title?: string
  subtitle?: string
  viewAllLabel?: string
  viewAllHref?: string
  deals?: Deal[]
  onAdd?: (deal: Deal) => void
}) {
  const [saved, setSaved] = useState<Record<string, boolean>>({})

  return (
    <section className="bg-gray-50 py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="overflow-hidden rounded-2xl border border-gray-200 dark:border-white/10">
          <div className="flex flex-wrap items-center justify-between gap-4 bg-gray-900 px-6 py-5">
            <div className="flex items-center gap-3">
              <Flame aria-hidden="true" className="size-7 text-amber-400" />
              <div>
                <h2 className="text-xl font-bold text-white">{title}</h2>
                <p className="text-sm text-gray-400">{subtitle}</p>
              </div>
            </div>
            <a
              href={viewAllHref}
              className="inline-flex min-h-11 items-center rounded-lg border border-white/20 bg-white/10 px-4 text-sm font-medium text-white hover:bg-white/15 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
            >
              {viewAllLabel}
            </a>
          </div>

          <ul
            role="list"
            className="grid grid-cols-1 gap-6 bg-white p-6 md:grid-cols-2 lg:grid-cols-3 dark:bg-gray-900"
          >
            {deals.map((deal) => {
              const discount = Math.round(
                ((deal.originalPrice - deal.currentPrice) / deal.originalPrice) * 100,
              )
              const remaining = Math.max(deal.total - deal.sold, 0)
              const soldPct = Math.min(Math.round((deal.sold / deal.total) * 100), 100)

              return (
                <li
                  key={deal.id}
                  className="group relative flex flex-col overflow-hidden rounded-lg border border-gray-200 bg-white transition-shadow hover:shadow-lg has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:bg-gray-950"
                >
                  <div className="absolute top-3 left-3 z-10 flex flex-col items-start gap-1">
                    {deal.isNew && (
                      <span className="flex items-center gap-1 rounded bg-emerald-700 px-2 py-1 text-xs font-bold text-white">
                        <Zap aria-hidden="true" className="size-3" />
                        New
                      </span>
                    )}
                    <span className="flex items-center gap-1 rounded bg-red-600 px-2 py-1 text-xs font-bold text-white">
                      <Flame aria-hidden="true" className="size-3" />
                      <span className="sr-only">Save </span>
                      {discount}%<span className="sr-only"> off</span>
                    </span>
                  </div>

                  {/* group-focus-within is load-bearing: without it these stay
                      invisible while still being focusable. */}
                  <div className="absolute top-3 right-3 z-10 flex flex-col gap-2 opacity-0 transition-opacity duration-200 group-hover:opacity-100 group-focus-within:opacity-100">
                    <button
                      type="button"
                      onClick={() => setSaved((s) => ({ ...s, [deal.id]: !s[deal.id] }))}
                      aria-pressed={Boolean(saved[deal.id])}
                      className="inline-flex size-11 items-center justify-center rounded-full bg-white/90 shadow-md backdrop-blur-sm hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                    >
                      <Heart
                        aria-hidden="true"
                        className={`size-4 ${
                          saved[deal.id] ? 'fill-red-500 text-red-500' : 'text-gray-700'
                        }`}
                      />
                      <span className="sr-only">Save {deal.name} to wishlist</span>
                    </button>
                    <a
                      href={deal.href ?? '#'}
                      className="inline-flex size-11 items-center justify-center rounded-full bg-white/90 shadow-md backdrop-blur-sm hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                    >
                      <Eye aria-hidden="true" className="size-4 text-gray-700" />
                      <span className="sr-only">Quick view of {deal.name}</span>
                    </a>
                  </div>

                  <div className="h-48 overflow-hidden bg-gray-50 dark:bg-white/5">
                    <img
                      src={deal.image}
                      alt=""
                      className="size-full object-cover transition-transform duration-500 group-hover:scale-110"
                    />
                    {deal.timeLeft && (
                      <span className="absolute top-40 right-3 flex items-center gap-1 rounded bg-orange-700 px-2 py-1 text-xs font-bold text-white">
                        <Timer aria-hidden="true" className="size-3" />
                        <span className="sr-only">Ends in </span>
                        {deal.timeLeft}
                      </span>
                    )}
                  </div>

                  <div className="flex flex-1 flex-col gap-3 p-4">
                    <h3 className="line-clamp-2 text-sm font-semibold text-gray-900 dark:text-white">
                      <a
                        href={deal.href ?? '#'}
                        className="after:absolute after:inset-0 focus:outline-none"
                      >
                        {deal.name}
                      </a>
                    </h3>

                    <div>
                      <p className="flex items-baseline gap-2">
                        <span className="text-xl font-bold text-gray-900 tabular-nums dark:text-white">
                          {CURRENCY.format(deal.currentPrice)}
                        </span>
                        <span className="text-sm text-gray-600 line-through tabular-nums dark:text-gray-400">
                          <span className="sr-only">Was </span>
                          {CURRENCY.format(deal.originalPrice)}
                        </span>
                      </p>
                      <p className="mt-0.5 text-sm font-medium text-emerald-700 dark:text-emerald-400">
                        You save {CURRENCY.format(deal.originalPrice - deal.currentPrice)}
                      </p>
                    </div>

                    <div className="mt-auto">
                      <div className="flex items-center justify-between text-xs text-gray-600 dark:text-gray-400">
                        <span>{deal.sold} sold</span>
                        {remaining <= 20 && remaining > 0 && (
                          <span className="font-medium text-orange-700 dark:text-orange-400">
                            Only {remaining} left
                          </span>
                        )}
                      </div>
                      {/* A real progressbar. A bare coloured div announces
                          nothing, which loses the whole "almost gone" signal. */}
                      <div
                        role="progressbar"
                        aria-valuenow={deal.sold}
                        aria-valuemin={0}
                        aria-valuemax={deal.total}
                        aria-label={`${deal.sold} of ${deal.total} sold`}
                        className="mt-1.5 h-1.5 w-full overflow-hidden rounded-full bg-gray-200 dark:bg-white/10"
                      >
                        <span
                          className="block h-full rounded-full bg-gradient-to-r from-orange-500 to-red-500"
                          style={{ width: `${soldPct}%` }}
                        />
                      </div>
                    </div>

                    <button
                      type="button"
                      onClick={() => onAdd?.(deal)}
                      className="relative z-10 inline-flex min-h-11 w-full items-center justify-center gap-2 rounded bg-red-600 text-sm font-medium text-white hover:bg-red-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-600"
                    >
                      <ShoppingBag aria-hidden="true" className="size-4" />
                      Add to cart
                      <span className="sr-only">: {deal.name}</span>
                    </button>
                  </div>
                </li>
              )
            })}
          </ul>
        </div>
      </div>
    </section>
  )
}
