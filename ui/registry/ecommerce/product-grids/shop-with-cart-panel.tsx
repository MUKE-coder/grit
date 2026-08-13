'use client'

import { useMemo, useState } from 'react'
import {
  ArrowUpRight,
  Bell,
  Heart,
  Minus,
  Plus,
  Search,
  Settings,
  SlidersHorizontal,
  Trash2,
} from 'lucide-react'

/*
 * A storefront with the cart open beside the grid, rather than behind a badge.
 *
 * Self-contained, and the cart is real: adding, removing, changing a quantity
 * and deselecting a line all move the totals. A cart panel with hardcoded
 * figures is the usual version of this block and it cannot answer the only
 * questions that matter, like what the total says when everything is
 * deselected, or whether removing the last line leaves an empty state.
 *
 * ── Quantity steppers ───────────────────────────────────────────────────────
 * A minus, a number and a plus, where each button names its product: "Add one
 * more iPhone 17 Pro", not "+". Six identical plus buttons are six identical
 * entries in a screen reader's control list, and the row a user is on is not
 * recoverable from the label. The count between them is the live value, so it
 * is not a heading or a div: it is the text the buttons change.
 *
 * ── The wishlist toggle ─────────────────────────────────────────────────────
 * aria-pressed, because a heart that fills in is a state change and a filled
 * heart is invisible to a reader. It sits above the stretched card link rather
 * than inside it, so it is clickable and is not swallowed by the link.
 *
 * ── Selection ───────────────────────────────────────────────────────────────
 * Each cart line has a checkbox that decides whether it is paid for. That is a
 * genuine control, so it is a real checkbox with a name, and the total says in
 * text how many lines it covers rather than leaving the user to work out why
 * three items cost less than they expected.
 *
 * ── Prices ──────────────────────────────────────────────────────────────────
 * Integer cents throughout, formatted once with Intl. A struck-through price
 * gets an <s> and a spoken prefix, because strikethrough is a visual convention
 * that carries nothing on its own.
 */

export interface CartProduct {
  id: string
  name: string
  /** Integer cents. */
  price: number
  /** Integer cents, before the discount. */
  was: number
  category: string
  storage: string
  colour: string
  /** Tailwind classes for the placeholder tile. Swap for a real <img>. */
  swatch: string
}

const PRODUCTS: CartProduct[] = [
  { id: 'ip17p', name: 'iPhone 17 Pro', price: 134000, was: 140000, category: 'Mobile', storage: '8/256', colour: 'Silver', swatch: 'from-slate-200 to-slate-400' },
  { id: 'vx300', name: 'Vivo X300 5G', price: 78000, was: 82500, category: 'Mobile', storage: '12/256', colour: 'Sky', swatch: 'from-sky-200 to-sky-400' },
  { id: 'gs25u', name: 'Galaxy S25 Ultra', price: 125000, was: 135000, category: 'Mobile', storage: '12/512', colour: 'Graphite', swatch: 'from-gray-700 to-gray-900' },
  { id: 'ip17', name: 'iPhone 17', price: 78000, was: 82000, category: 'Mobile', storage: '8/128', colour: 'Midnight', swatch: 'from-zinc-500 to-zinc-800' },
  { id: 'ipair', name: 'iPhone Air', price: 89000, was: 93000, category: 'Mobile', storage: '8/256', colour: 'Cloud', swatch: 'from-blue-100 to-blue-300' },
  { id: 'op13', name: 'OnePlus 13', price: 75000, was: 80000, category: 'Mobile', storage: '16/512', colour: 'White', swatch: 'from-neutral-100 to-neutral-300' },
  { id: 'iq15', name: 'iQOO 15', price: 72000, was: 78000, category: 'Mobile', storage: '16/512', colour: 'Silver', swatch: 'from-slate-300 to-slate-500' },
  { id: 'rgt7', name: 'Realme GT 7T', price: 46000, was: 52000, category: 'Mobile', storage: '12/256', colour: 'Blue', swatch: 'from-indigo-300 to-indigo-500' },
  { id: 'px10', name: 'Pixel 10 Pro', price: 99000, was: 105000, category: 'Mobile', storage: '16/256', colour: 'Porcelain', swatch: 'from-stone-200 to-stone-400' },
]

const CATEGORIES = ['Mobile', 'All products', 'Grocery', 'Most purchased', 'Shoes', 'Furniture']
const NAV = ['Dashboard', 'Products', 'My orders', 'Wishlist', 'Cart', 'Settings']

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const PROMO_RATE = 0.12

export default function ShopWithCartPanel({ products = PRODUCTS }: { products?: CartProduct[] }) {
  const [category, setCategory] = useState('Mobile')
  const [query, setQuery] = useState('')
  const [saved, setSaved] = useState<string[]>(['vx300', 'op13'])
  const [lines, setLines] = useState<{ id: string; quantity: number; selected: boolean }[]>([
    { id: 'ip17p', quantity: 1, selected: true },
    { id: 'gs25u', quantity: 1, selected: true },
    { id: 'iq15', quantity: 1, selected: true },
  ])
  const [announcement, setAnnouncement] = useState('')

  const visible = useMemo(
    () =>
      products.filter((product) => {
        if (category !== 'All products' && product.category !== category) return false
        if (query && !product.name.toLowerCase().includes(query.toLowerCase())) return false
        return true
      }),
    [products, category, query],
  )

  const detailed = lines
    .map((line) => ({ line, product: products.find((item) => item.id === line.id) }))
    .filter((entry): entry is { line: (typeof lines)[number]; product: CartProduct } =>
      Boolean(entry.product),
    )

  const selectedLines = detailed.filter((entry) => entry.line.selected)
  const subtotal = selectedLines.reduce(
    (sum, entry) => sum + entry.product.price * entry.line.quantity,
    0,
  )
  const discount = Math.round(subtotal * PROMO_RATE)
  const payable = subtotal - discount

  function announce(message: string) {
    setAnnouncement('')
    requestAnimationFrame(() => setAnnouncement(message))
  }

  function addToCart(product: CartProduct) {
    setLines((current) => {
      const existing = current.find((line) => line.id === product.id)
      if (existing) {
        return current.map((line) =>
          line.id === product.id ? { ...line, quantity: line.quantity + 1 } : line,
        )
      }
      return [...current, { id: product.id, quantity: 1, selected: true }]
    })
    announce(`${product.name} added to the cart.`)
  }

  function setQuantity(product: CartProduct, delta: number) {
    setLines((current) =>
      current.flatMap((line) => {
        if (line.id !== product.id) return [line]
        const quantity = line.quantity + delta
        /* Stepping below one removes the line rather than leaving a zero,
           because a cart row for nothing is not a thing anyone wants. */
        if (quantity < 1) return []
        return [{ ...line, quantity }]
      }),
    )
    announce(
      delta > 0
        ? `Added one more ${product.name}.`
        : `Removed one ${product.name}.`,
    )
  }

  function toggleLine(product: CartProduct) {
    setLines((current) =>
      current.map((line) =>
        line.id === product.id ? { ...line, selected: !line.selected } : line,
      ),
    )
  }

  function removeLine(product: CartProduct) {
    setLines((current) => current.filter((line) => line.id !== product.id))
    announce(`${product.name} removed from the cart.`)
  }

  function toggleSaved(product: CartProduct) {
    setSaved((current) =>
      current.includes(product.id)
        ? current.filter((id) => id !== product.id)
        : [...current, product.id],
    )
  }

  return (
    <div className="bg-gray-100 py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4">
        <p role="status" aria-live="polite" className="sr-only">
          {announcement}
        </p>

        <div className="overflow-hidden rounded-3xl bg-white shadow-sm dark:bg-gray-900">
          {/* ── Top bar ─────────────────────────────────────────────────── */}
          <header className="flex flex-wrap items-center gap-4 border-b border-gray-200 px-5 py-4 dark:border-white/10">
            <a
              href="#"
              className="flex items-center gap-2 text-lg font-semibold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-orange-700 dark:text-white"
            >
              <span
                aria-hidden="true"
                className="grid size-7 place-items-center rounded-full bg-gradient-to-br from-orange-400 to-orange-600 text-xs font-bold text-white"
              >
                S
              </span>
              SellPilot
            </a>

            <nav aria-label="Main" className="hidden min-w-0 flex-1 xl:block">
              <ul className="flex items-center justify-center gap-6 text-sm">
                {NAV.map((item) => {
                  const current = item === 'Cart'
                  return (
                    <li key={item}>
                      <a
                        href="#"
                        aria-current={current ? 'page' : undefined}
                        className={`focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-orange-700 ${
                          current
                            ? 'font-semibold text-orange-700 dark:text-orange-400'
                            : 'text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white'
                        }`}
                      >
                        {item}
                      </a>
                    </li>
                  )
                })}
              </ul>
            </nav>

            <div className="ml-auto flex items-center gap-2">
              <button
                type="button"
                className="grid size-10 place-items-center rounded-full border border-gray-200 text-gray-600 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:text-gray-300 dark:hover:bg-white/5"
              >
                <Bell aria-hidden="true" className="size-4" />
                <span className="sr-only">Notifications</span>
              </button>
              <button
                type="button"
                className="grid size-10 place-items-center rounded-full border border-gray-200 text-gray-600 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:text-gray-300 dark:hover:bg-white/5"
              >
                <Settings aria-hidden="true" className="size-4" />
                <span className="sr-only">Settings</span>
              </button>
              <button
                type="button"
                className="grid size-10 place-items-center rounded-full bg-orange-700 text-xs font-semibold text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
              >
                RS
                <span className="sr-only">Account menu for Rahim Sheikh</span>
              </button>
            </div>
          </header>

          <div className="grid gap-5 p-5 xl:grid-cols-[minmax(0,1fr)_23rem]">
            {/* ── Catalogue ─────────────────────────────────────────────── */}
            <section aria-labelledby="catalogue" className="min-w-0">
              <div className="mb-4 flex flex-wrap items-center gap-3">
                <h2
                  id="catalogue"
                  className="text-xl font-semibold tracking-tight text-gray-900 dark:text-white"
                >
                  Cart products
                </h2>

                <div className="relative ml-auto min-w-0 flex-1 sm:max-w-56">
                  <Search
                    aria-hidden="true"
                    className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-gray-400"
                  />
                  <label htmlFor="shop-search" className="sr-only">
                    Search products by name
                  </label>
                  <input
                    id="shop-search"
                    type="search"
                    value={query}
                    onChange={(event) => setQuery(event.target.value)}
                    placeholder="Search"
                    className="h-10 w-full rounded-full border border-gray-200 bg-gray-50 pl-9 pr-3 text-sm text-gray-900 placeholder:text-gray-500 focus:border-orange-700 focus:outline-none focus:ring-1 focus:ring-orange-700 dark:border-white/15 dark:bg-white/5 dark:text-white dark:placeholder:text-gray-400"
                  />
                </div>

                <button
                  type="button"
                  className="inline-flex min-h-10 items-center gap-2 rounded-full border border-gray-200 px-3.5 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                >
                  <SlidersHorizontal aria-hidden="true" className="size-4" />
                  Filter
                </button>

                <button
                  type="button"
                  className="inline-flex min-h-10 items-center gap-2 rounded-full bg-orange-700 px-4 text-sm font-semibold text-white hover:bg-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                >
                  <Plus aria-hidden="true" className="size-4" />
                  Create order
                </button>
              </div>

              <div className="mb-4 flex flex-wrap gap-2">
                {CATEGORIES.map((item) => (
                  <button
                    key={item}
                    type="button"
                    onClick={() => setCategory(item)}
                    aria-pressed={category === item}
                    className={`min-h-9 rounded-full px-4 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 ${
                      category === item
                        ? 'bg-gray-900 font-medium text-white dark:bg-white dark:text-gray-900'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-white/10 dark:text-gray-200 dark:hover:bg-white/15'
                    }`}
                  >
                    {item}
                  </button>
                ))}
              </div>

              <p role="status" aria-live="polite" className="sr-only">
                {visible.length} products in {category}.
              </p>

              {visible.length === 0 ? (
                <p className="rounded-2xl border border-dashed border-gray-300 px-6 py-16 text-center text-sm text-gray-600 dark:border-white/15 dark:text-gray-300">
                  Nothing in {category} matches that search.
                </p>
              ) : (
                <ul className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                  {visible.map((product) => {
                    const isSaved = saved.includes(product.id)
                    return (
                      <li
                        key={product.id}
                        className="group relative rounded-2xl bg-gray-50 p-3 dark:bg-white/5"
                      >
                        <div className="relative mb-3">
                          <span
                            aria-hidden="true"
                            className={`block h-32 rounded-xl bg-gradient-to-br ${product.swatch}`}
                          />
                          {/* Layered above the stretched link, so it can be
                              clicked, and pressed rather than merely filled. */}
                          <button
                            type="button"
                            onClick={() => toggleSaved(product)}
                            aria-pressed={isSaved}
                            className="absolute right-2 top-2 z-10 grid size-9 place-items-center rounded-full bg-white text-gray-700 shadow-sm hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:bg-gray-900 dark:text-gray-200"
                          >
                            <Heart
                              aria-hidden="true"
                              className={`size-4 ${isSaved ? 'fill-orange-600 text-orange-600' : ''}`}
                            />
                            <span className="sr-only">Save {product.name} to wishlist</span>
                          </button>
                        </div>

                        <h3 className="text-sm font-medium text-gray-900 dark:text-white">
                          <a
                            href="#"
                            className="after:absolute after:inset-0 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                          >
                            {product.name}
                          </a>
                        </h3>

                        <div className="mt-1 flex items-center justify-between gap-2">
                          <p className="flex items-baseline gap-2 text-sm">
                            <span className="font-semibold tabular-nums text-orange-700 dark:text-orange-400">
                              {money.format(product.price / 100)}
                            </span>
                            {/* <s> plus a spoken prefix: strikethrough alone is
                                a visual convention and says nothing. */}
                            <s className="tabular-nums text-gray-500 dark:text-gray-400">
                              <span className="sr-only">Was </span>
                              {money.format(product.was / 100)}
                            </s>
                          </p>

                          <button
                            type="button"
                            onClick={() => addToCart(product)}
                            className="relative z-10 grid size-9 shrink-0 place-items-center rounded-full bg-orange-700 text-white hover:bg-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                          >
                            <ArrowUpRight aria-hidden="true" className="size-4" />
                            <span className="sr-only">Add {product.name} to cart</span>
                          </button>
                        </div>
                      </li>
                    )
                  })}
                </ul>
              )}
            </section>

            {/* ── Cart ──────────────────────────────────────────────────── */}
            <section
              aria-labelledby="cart"
              className="h-fit rounded-2xl border border-gray-200 p-4 dark:border-white/10"
            >
              <div className="mb-4 flex items-center justify-between gap-2">
                <h2
                  id="cart"
                  className="text-lg font-semibold tracking-tight text-gray-900 dark:text-white"
                >
                  Cart
                </h2>
                <button
                  type="button"
                  onClick={() => {
                    setLines([])
                    announce('Cart emptied.')
                  }}
                  disabled={detailed.length === 0}
                  className="min-h-9 rounded-full border border-gray-200 px-3 text-xs font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                >
                  Delete all
                </button>
              </div>

              {detailed.length === 0 ? (
                <p className="rounded-xl border border-dashed border-gray-300 px-4 py-12 text-center text-sm text-gray-600 dark:border-white/15 dark:text-gray-300">
                  Your cart is empty. Add something from the grid.
                </p>
              ) : (
                <ul className="space-y-3">
                  {detailed.map(({ line, product }) => (
                    <li
                      key={product.id}
                      className="flex gap-3 rounded-xl border border-gray-200 p-3 dark:border-white/10"
                    >
                      <label className="flex items-start pt-1">
                        <input
                          type="checkbox"
                          checked={line.selected}
                          onChange={() => toggleLine(product)}
                          className="size-4 accent-orange-600"
                        />
                        <span className="sr-only">Include {product.name} in the total</span>
                      </label>

                      <span
                        aria-hidden="true"
                        className={`h-14 w-11 shrink-0 rounded-lg bg-gradient-to-br ${product.swatch}`}
                      />

                      <div className="min-w-0 flex-1">
                        <div className="flex items-start justify-between gap-2">
                          <h3 className="truncate text-sm font-medium text-gray-900 dark:text-white">
                            {product.name}
                          </h3>
                          <button
                            type="button"
                            onClick={() => removeLine(product)}
                            className="grid size-8 shrink-0 place-items-center rounded-md text-gray-500 hover:bg-rose-50 hover:text-rose-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:text-gray-400 dark:hover:bg-rose-500/10 dark:hover:text-rose-300"
                          >
                            <Trash2 aria-hidden="true" className="size-4" />
                            <span className="sr-only">Remove {product.name} from the cart</span>
                          </button>
                        </div>

                        <p className="mt-0.5 flex flex-wrap items-center gap-2 text-xs text-gray-600 dark:text-gray-300">
                          <span className="rounded-full bg-gray-100 px-2 py-0.5 dark:bg-white/10">
                            {product.storage}
                          </span>
                          <span className="rounded-full bg-gray-100 px-2 py-0.5 dark:bg-white/10">
                            {product.colour}
                          </span>
                        </p>

                        <div className="mt-2 flex items-center justify-between gap-2">
                          <span className="flex items-center gap-1">
                            <button
                              type="button"
                              onClick={() => setQuantity(product, -1)}
                              className="grid size-8 place-items-center rounded-full border border-gray-200 text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                            >
                              <Minus aria-hidden="true" className="size-3.5" />
                              {/* Named per product: six identical minus
                                  buttons are six identical controls. */}
                              <span className="sr-only">
                                {line.quantity === 1
                                  ? `Remove ${product.name} from the cart`
                                  : `Remove one ${product.name}`}
                              </span>
                            </button>
                            <span className="min-w-6 text-center text-sm tabular-nums text-gray-900 dark:text-white">
                              {line.quantity}
                              <span className="sr-only"> in cart</span>
                            </span>
                            <button
                              type="button"
                              onClick={() => setQuantity(product, 1)}
                              className="grid size-8 place-items-center rounded-full bg-orange-700 text-white hover:bg-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                            >
                              <Plus aria-hidden="true" className="size-3.5" />
                              <span className="sr-only">Add one more {product.name}</span>
                            </button>
                          </span>

                          <p className="text-sm font-semibold tabular-nums text-gray-900 dark:text-white">
                            {money.format((product.price * line.quantity) / 100)}
                          </p>
                        </div>
                      </div>
                    </li>
                  ))}
                </ul>
              )}

              <div className="mt-4 flex items-center justify-between gap-2 rounded-xl bg-gray-50 px-3 py-2.5 dark:bg-white/5">
                <p className="text-sm text-gray-700 dark:text-gray-200">
                  Promo applied, {Math.round(PROMO_RATE * 100)}% off
                </p>
                <button
                  type="button"
                  className="min-h-9 rounded-full bg-orange-100 px-3 text-xs font-semibold text-orange-800 hover:bg-orange-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:bg-orange-400/15 dark:text-orange-200"
                >
                  Change promo
                </button>
              </div>

              {/* A description list: each label and its figure are a pair, and
                  should be announced as one. */}
              <dl className="mt-4 space-y-2 border-t border-gray-200 pt-4 text-sm dark:border-white/10">
                <div className="flex items-center justify-between gap-3">
                  <dt className="text-gray-600 dark:text-gray-300">
                    Subtotal
                    <span className="sr-only">
                      , covering {selectedLines.length} of {detailed.length} lines
                    </span>
                  </dt>
                  <dd className="font-medium tabular-nums text-gray-900 dark:text-white">
                    {money.format(subtotal / 100)}
                  </dd>
                </div>
                <div className="flex items-center justify-between gap-3">
                  <dt className="text-gray-600 dark:text-gray-300">Discount</dt>
                  <dd className="font-medium tabular-nums text-emerald-700 dark:text-emerald-400">
                    &minus;{money.format(discount / 100)}
                  </dd>
                </div>
                <div className="flex items-center justify-between gap-3 border-t border-gray-200 pt-2 dark:border-white/10">
                  <dt className="font-semibold text-gray-900 dark:text-white">Total payment</dt>
                  <dd className="text-lg font-semibold tabular-nums text-gray-900 dark:text-white">
                    {money.format(payable / 100)}
                  </dd>
                </div>
              </dl>

              <button
                type="button"
                disabled={selectedLines.length === 0}
                className="mt-4 flex min-h-12 w-full items-center justify-center rounded-full bg-gradient-to-r from-orange-600 to-orange-700 px-4 text-sm font-semibold text-white hover:from-orange-700 hover:to-orange-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 disabled:from-gray-400 disabled:to-gray-400"
              >
                Proceed to payment
                {selectedLines.length > 0 && (
                  <span className="sr-only">
                    , {money.format(payable / 100)} for {selectedLines.length}{' '}
                    {selectedLines.length === 1 ? 'item' : 'items'}
                  </span>
                )}
              </button>
              {selectedLines.length === 0 && detailed.length > 0 && (
                <p className="mt-2 text-center text-xs text-gray-600 dark:text-gray-300">
                  Tick at least one line to check out.
                </p>
              )}
            </section>
          </div>
        </div>
      </div>
    </div>
  )
}
