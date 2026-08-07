'use client'

import { useEffect, useId, useRef, useState } from 'react'
import { Minus, Plus, ShoppingBag } from 'lucide-react'

/*
 * The bottom half of a product page: a frequently-bought-together bundle with
 * a running total, and a buy bar that appears once the real one scrolls away.
 *
 * The bar is driven by an IntersectionObserver watching the main add-to-cart
 * button, so it appears exactly when that button leaves the viewport, which is
 * the only moment it is useful. The source read getBoundingClientRect on every
 * scroll event, which forces a synchronous layout on a handler that fires
 * dozens of times a second, and pinned the threshold to an element top rather
 * than to the control being replaced.
 *
 * When the bar is away it is `invisible`, not merely translated off-screen. A
 * bar at translate-y-full still holds a quantity stepper and an add-to-cart
 * button, so tabbing from the top of the page walks a keyboard user through
 * controls sitting above the viewport. visibility takes them out of the tab
 * order and still transitions, stepping at the end of the slide-out.
 *
 * Both steppers name their product, and the quantity is announced when it
 * changes. A stepper whose entire accessible name is "-" tells you nothing,
 * and a number that only updates in a span tells you nothing either.
 *
 * The bundle is a fieldset of checkboxes. Multi-select is what checkboxes are;
 * the source used clickable divs holding a styled square, so the state was
 * visible and nothing else.
 *
 * Money is integer cents, and the total is recomputed from the selection
 * rather than accumulated, so it cannot drift out of step with the boxes.
 */

export interface Addon {
  id: string
  name: string
  detail: string
  /** Integer cents. */
  price: number
  image: string
}

const ADDONS: Addon[] = [
  {
    id: 'watch',
    name: 'Automatic wristwatch',
    detail: 'Sapphire crystal, 100m',
    price: 89900,
    image: 'https://images.unsplash.com/photo-1522312346375-d1a52e2b99b3?w=200&h=200&fit=crop&q=80',
  },
  {
    id: 'headphones',
    name: 'Over-ear headphones',
    detail: 'Noise cancelling, 40h',
    price: 37900,
    image: 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=200&h=200&fit=crop&q=80',
  },
  {
    id: 'skincare',
    name: 'Travel skincare set',
    detail: 'Under 100ml, cabin safe',
    price: 8400,
    image: 'https://images.unsplash.com/photo-1583209814683-c023dd293cc6?w=200&h=200&fit=crop&q=80',
  },
]

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const format = (cents: number) => money.format(cents / 100)

export default function BundlePickerWithStickyBuyBar({
  name = 'Leather weekender bag',
  price = 19999,
  image = 'https://images.unsplash.com/photo-1590874103328-eac38a683ce7?w=400&h=400&fit=crop&q=80',
  addons = ADDONS,
}: {
  name?: string
  /** Integer cents. */
  price?: number
  image?: string
  addons?: Addon[]
}) {
  const [quantity, setQuantity] = useState(1)
  const [selected, setSelected] = useState<string[]>([addons[0]?.id].filter(Boolean) as string[])
  const [barVisible, setBarVisible] = useState(false)
  const [announcement, setAnnouncement] = useState('')

  const mainButton = useRef<HTMLButtonElement>(null)
  const bundleId = useId()

  /* Watches the button the bar stands in for. A scroll threshold would be a
     guess about where that button happens to sit. */
  useEffect(() => {
    const target = mainButton.current
    if (!target) return
    const observer = new IntersectionObserver(
      ([entry]) => setBarVisible(!entry.isIntersecting),
      { threshold: 0 },
    )
    observer.observe(target)
    return () => observer.disconnect()
  }, [])

  /* Recomputed from the selection every render. Adding and subtracting from a
     stored total as boxes are ticked is how a bundle price ends up disagreeing
     with the boxes. */
  const addonsTotal = addons
    .filter((addon) => selected.includes(addon.id))
    .reduce((sum, addon) => sum + addon.price, 0)
  const total = price * quantity + addonsTotal

  function announce(message: string) {
    setAnnouncement('')
    requestAnimationFrame(() => setAnnouncement(message))
  }

  function changeQuantity(delta: number) {
    const next = Math.max(1, quantity + delta)
    if (next === quantity) return
    setQuantity(next)
    announce(`Quantity ${next}. Total ${format(price * next + addonsTotal)}.`)
  }

  function toggleAddon(addon: Addon) {
    const isOn = selected.includes(addon.id)
    const nextSelected = isOn
      ? selected.filter((id) => id !== addon.id)
      : [...selected, addon.id]
    setSelected(nextSelected)

    const nextTotal =
      price * quantity +
      addons.filter((a) => nextSelected.includes(a.id)).reduce((sum, a) => sum + a.price, 0)
    announce(`${addon.name} ${isOn ? 'removed' : 'added'}. Total ${format(nextTotal)}.`)
  }

  function Stepper({ compact }: { compact?: boolean }) {
    return (
      <div className="flex items-center rounded-lg border border-gray-300 dark:border-white/15">
        <button
          type="button"
          onClick={() => changeQuantity(-1)}
          disabled={quantity === 1}
          className={`inline-flex items-center justify-center rounded-l-lg text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 disabled:opacity-40 dark:text-gray-300 dark:hover:bg-white/10 ${
            compact ? 'size-11' : 'size-12'
          }`}
        >
          <Minus aria-hidden="true" className="size-4" />
          <span className="sr-only">Decrease quantity of {name}</span>
        </button>
        <span
          aria-hidden="true"
          className="w-10 text-center text-sm tabular-nums text-gray-900 dark:text-white"
        >
          {quantity}
        </span>
        <button
          type="button"
          onClick={() => changeQuantity(1)}
          className={`inline-flex items-center justify-center rounded-r-lg text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:bg-white/10 ${
            compact ? 'size-11' : 'size-12'
          }`}
        >
          <Plus aria-hidden="true" className="size-4" />
          <span className="sr-only">Increase quantity of {name}</span>
        </button>
      </div>
    )
  }

  return (
    <>
      {/* One status region for the whole block. Quantity and bundle changes
          both land here, and neither is visible from the control you used. */}
      <p role="status" aria-live="polite" className="sr-only">
        {announcement}
      </p>

      <section className="bg-white py-12 dark:bg-gray-950">
        <div className="mx-auto max-w-3xl px-4">
          <div className="flex flex-wrap items-center gap-6 rounded-xl border border-gray-200 p-6 dark:border-white/10">
            <img src={image} alt="" className="size-24 rounded-lg object-cover" />
            <div className="min-w-0 flex-1">
              <h2 className="font-semibold text-gray-900 dark:text-white">{name}</h2>
              <p className="mt-1 text-lg font-bold text-gray-900 dark:text-white">
                {format(price)}
              </p>
            </div>
            <Stepper />
            <button
              ref={mainButton}
              type="button"
              className="inline-flex min-h-12 items-center gap-2 rounded-lg bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:bg-white dark:text-gray-900"
            >
              <ShoppingBag aria-hidden="true" className="size-4" />
              Add to cart
            </button>
          </div>

          <fieldset className="mt-10">
            <legend className="text-lg font-semibold text-gray-900 dark:text-white">
              Frequently bought together
            </legend>
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Add any of these to the same order.
            </p>

            <ul role="list" className="mt-4 space-y-3">
              {addons.map((addon) => {
                const inputId = `${bundleId}-${addon.id}`
                const checked = selected.includes(addon.id)
                return (
                  <li key={addon.id}>
                    {/* The whole row is the label, so the hit area is the row
                        and no click handler is needed to make it work. */}
                    <label
                      htmlFor={inputId}
                      className={`flex cursor-pointer items-center gap-4 rounded-lg border p-4 transition-colors has-[:focus-visible]:outline has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-indigo-600 ${
                        checked
                          ? 'border-gray-900 bg-gray-50 dark:border-white dark:bg-white/5'
                          : 'border-gray-200 hover:bg-gray-50 dark:border-white/10 dark:hover:bg-white/5'
                      }`}
                    >
                      <input
                        id={inputId}
                        type="checkbox"
                        checked={checked}
                        onChange={() => toggleAddon(addon)}
                        className="size-4 shrink-0 accent-gray-900 dark:accent-white"
                      />
                      <img src={addon.image} alt="" className="size-14 rounded object-cover" />
                      <span className="min-w-0 flex-1">
                        <span className="block font-medium text-gray-900 dark:text-white">
                          {addon.name}
                        </span>
                        <span className="block text-sm text-gray-500 dark:text-gray-400">
                          {addon.detail}
                        </span>
                      </span>
                      <span className="font-semibold text-gray-900 dark:text-white">
                        {format(addon.price)}
                      </span>
                    </label>
                  </li>
                )
              })}
            </ul>

            <p className="mt-6 flex items-baseline justify-between border-t border-gray-200 pt-4 dark:border-white/10">
              <span className="text-gray-600 dark:text-gray-300">
                Total for {quantity + selected.length}{' '}
                {quantity + selected.length === 1 ? 'item' : 'items'}
              </span>
              <span className="text-2xl font-bold text-gray-900 dark:text-white">
                {format(total)}
              </span>
            </p>
          </fieldset>

          <div className="h-[60vh]" aria-hidden="true" />
        </div>
      </section>

      {/* invisible, not just translated away, or the stepper and the button
          below stay in the tab order while the bar is off screen. */}
      <div
        className={`fixed inset-x-0 top-0 z-50 border-b border-gray-200 bg-white transition-[transform,visibility] duration-300 dark:border-white/10 dark:bg-gray-950 ${
          barVisible ? 'visible translate-y-0 shadow-md' : 'invisible -translate-y-full'
        }`}
      >
        <div className="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
          <div className="flex min-w-0 items-center gap-3">
            <img src={image} alt="" className="size-12 shrink-0 rounded object-cover" />
            <div className="min-w-0">
              <p className="truncate text-sm font-medium text-gray-900 dark:text-white">{name}</p>
              <p className="text-sm text-gray-500 dark:text-gray-400">{format(total)}</p>
            </div>
          </div>
          <div className="flex shrink-0 items-center gap-3">
            <div className="hidden sm:block">
              <Stepper compact />
            </div>
            <button
              type="button"
              className="inline-flex min-h-11 items-center gap-2 rounded-lg bg-gray-900 px-5 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:bg-white dark:text-gray-900"
            >
              <ShoppingBag aria-hidden="true" className="size-4" />
              Add to cart
            </button>
          </div>
        </div>
      </div>
    </>
  )
}
