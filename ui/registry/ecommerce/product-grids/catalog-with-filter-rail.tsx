'use client'

import { useId, useMemo, useState } from 'react'
import { Plus, RefreshCw, Search, SlidersHorizontal, X } from 'lucide-react'

/*
 * A catalogue manager: a filter rail beside a product grid, with tabs, search
 * and pagination.
 *
 * Self-contained, and everything filters for real. A filter rail that renders
 * but does not filter is the most common way this block ships, and it hides
 * every question worth answering: what the empty state says, whether the count
 * is announced, what happens when the last filter is removed.
 *
 * ── The filter groups ───────────────────────────────────────────────────────
 * Each group is a <fieldset> with a <legend>. A heading above a pile of
 * checkboxes looks identical and means nothing: the legend is what makes a
 * screen reader say "Brands, Gucci, checkbox, not checked" instead of leaving
 * eleven unrelated boxes in a row.
 *
 * The colour swatches are checkboxes with real names. A grid of coloured
 * circles is unusable to anyone who cannot see them and ambiguous to anyone
 * who can: "the blue one" is three of these. The name is visible on focus and
 * always present for a reader.
 *
 * ── The price range ─────────────────────────────────────────────────────────
 * Two labelled number inputs, plus a histogram that is decoration. A dual
 * thumb slider is the design everybody draws and nobody can operate with a
 * keyboard without real work; two inputs are operable by everyone, and the
 * histogram behind them still communicates the distribution.
 *
 * ── Counts ──────────────────────────────────────────────────────────────────
 * The result count is a live region. Ticking a checkbox changes the grid
 * without moving focus, so without it the page silently becomes something
 * else.
 *
 * ── The cards ───────────────────────────────────────────────────────────────
 * One link per card, stretched over the whole tile, with the row of actions
 * layered above it rather than nested inside it. Two links to the same place
 * is two entries in a link list. The rating is written out beside the star,
 * because a star is a picture of a number.
 */

export interface Product {
  id: string
  name: string
  /** Integer cents. */
  price: number
  rating: number
  type: 'Clothing' | 'Shoes' | 'Bags' | 'Accessories' | 'Jewelry'
  brand: string
  colour: string
  active: boolean
  /** Tailwind classes for the placeholder tile. Swap for a real <img>. */
  swatch: string
}

const PRODUCTS: Product[] = [
  { id: 'p1', name: 'Studio Display 27"', price: 159900, rating: 4.6, type: 'Accessories', brand: 'Aurio', colour: 'Violet', active: true, swatch: 'from-violet-300 to-violet-500' },
  { id: 'p2', name: 'Halo Phone 16', price: 99900, rating: 4.4, type: 'Accessories', brand: 'Nimbus', colour: 'Black', active: true, swatch: 'from-gray-700 to-gray-900' },
  { id: 'p3', name: 'Air Laptop 13"', price: 129900, rating: 4.8, type: 'Accessories', brand: 'Aurio', colour: 'Blue', active: true, swatch: 'from-sky-300 to-blue-500' },
  { id: 'p4', name: 'Fold Handset 6', price: 179900, rating: 4.2, type: 'Accessories', brand: 'Kestrel', colour: 'Green', active: false, swatch: 'from-emerald-300 to-emerald-600' },
  { id: 'p5', name: 'Quiet Headphones', price: 34900, rating: 4.7, type: 'Accessories', brand: 'Sonora', colour: 'White', active: true, swatch: 'from-gray-100 to-gray-300' },
  { id: 'p6', name: 'Orb Speaker Mini', price: 9900, rating: 4.1, type: 'Accessories', brand: 'Sonora', colour: 'Orange', active: true, swatch: 'from-orange-300 to-orange-500' },
  { id: 'p7', name: 'Over-ear Max', price: 54900, rating: 4.5, type: 'Accessories', brand: 'Sonora', colour: 'Green', active: false, swatch: 'from-teal-200 to-emerald-400' },
  { id: 'p8', name: 'Buds Pro 2nd gen', price: 23000, rating: 4.6, type: 'Accessories', brand: 'Sonora', colour: 'White', active: true, swatch: 'from-gray-50 to-gray-200' },
  { id: 'p9', name: 'Tablet Pro 12.9"', price: 109900, rating: 4.3, type: 'Accessories', brand: 'Aurio', colour: 'Violet', active: true, swatch: 'from-fuchsia-400 to-purple-600' },
  { id: 'p10', name: 'Canvas Weekender', price: 18900, rating: 4.0, type: 'Bags', brand: 'Kestrel', colour: 'Yellow', active: true, swatch: 'from-amber-200 to-amber-400' },
  { id: 'p11', name: 'Runner Low Top', price: 12900, rating: 4.4, type: 'Shoes', brand: 'Kestrel', colour: 'Blue', active: true, swatch: 'from-indigo-300 to-indigo-500' },
  { id: 'p12', name: 'Field Watch 42mm', price: 44900, rating: 4.9, type: 'Jewelry', brand: 'Aurio', colour: 'Red', active: true, swatch: 'from-rose-200 to-rose-400' },
]

const TYPES: Product['type'][] = ['Clothing', 'Shoes', 'Bags', 'Accessories', 'Jewelry']
const BRANDS = ['Aurio', 'Nimbus', 'Kestrel', 'Sonora']
const COLOURS = ['Red', 'Orange', 'Yellow', 'Green', 'Blue', 'Violet', 'Black', 'White']

const SWATCH: Record<string, string> = {
  Red: 'bg-red-500',
  Orange: 'bg-orange-500',
  Yellow: 'bg-yellow-400',
  Green: 'bg-emerald-500',
  Blue: 'bg-blue-500',
  Violet: 'bg-violet-500',
  Black: 'bg-gray-900',
  White: 'bg-white',
}

/* Bar heights for the price histogram, as percentages. Decoration behind the
   two inputs, not a control. */
const HISTOGRAM = [22, 34, 28, 46, 62, 40, 55, 78, 48, 66, 90, 58, 44, 70, 36, 26]

const TABS = ['All', 'Active', 'Non active'] as const
type Tab = (typeof TABS)[number]

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const PAGE_SIZE = 8

/*
 * The legend IS the flex row, rather than a legend sitting inside one. A
 * <legend> only names its group when it is a direct child of the <fieldset>;
 * wrap it in a div for layout and the association is silently gone, which
 * leaves eleven unrelated checkboxes in a row. A button is phrasing content,
 * so it is allowed inside the legend.
 */
function Legend({ children, onClear }: { children: string; onClear?: () => void }) {
  return (
    <legend className="mb-3 flex w-full items-center justify-between gap-2 text-sm font-semibold text-gray-900 dark:text-white">
      <span>{children}</span>
      {onClear && (
        <button
          type="button"
          onClick={onClear}
          className="grid size-7 place-items-center rounded-md text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:text-gray-400 dark:hover:bg-white/10"
        >
          <X aria-hidden="true" className="size-3.5" />
          {/* Named per group. Five buttons called "Clear" are five identical
              entries in a control list. */}
          <span className="sr-only">Clear the {children.toLowerCase()} filter</span>
        </button>
      )}
    </legend>
  )
}

export default function CatalogWithFilterRail({ products = PRODUCTS }: { products?: Product[] }) {
  const baseId = useId()
  const [tab, setTab] = useState<Tab>('All')
  const [types, setTypes] = useState<string[]>([])
  const [brands, setBrands] = useState<string[]>([])
  const [colours, setColours] = useState<string[]>([])
  const [minPrice, setMinPrice] = useState('')
  const [maxPrice, setMaxPrice] = useState('')
  const [query, setQuery] = useState('')
  const [page, setPage] = useState(1)

  const results = useMemo(() => {
    const min = minPrice === '' ? 0 : Number(minPrice) * 100
    const max = maxPrice === '' ? Infinity : Number(maxPrice) * 100
    return products.filter((product) => {
      if (tab === 'Active' && !product.active) return false
      if (tab === 'Non active' && product.active) return false
      if (types.length > 0 && !types.includes(product.type)) return false
      if (brands.length > 0 && !brands.includes(product.brand)) return false
      if (colours.length > 0 && !colours.includes(product.colour)) return false
      if (product.price < min || product.price > max) return false
      if (query && !product.name.toLowerCase().includes(query.toLowerCase())) return false
      return true
    })
  }, [products, tab, types, brands, colours, minPrice, maxPrice, query])

  const pageCount = Math.max(1, Math.ceil(results.length / PAGE_SIZE))
  const safePage = Math.min(page, pageCount)
  const visible = results.slice((safePage - 1) * PAGE_SIZE, safePage * PAGE_SIZE)

  const activeFilters = types.length + brands.length + colours.length

  function toggle(list: string[], setList: (next: string[]) => void, value: string) {
    setList(list.includes(value) ? list.filter((item) => item !== value) : [...list, value])
    setPage(1)
  }

  function clearAll() {
    setTypes([])
    setBrands([])
    setColours([])
    setMinPrice('')
    setMaxPrice('')
    setQuery('')
    setPage(1)
  }

  return (
    <div className="bg-gray-50 py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4">
        <div className="overflow-hidden rounded-2xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-5 py-4 dark:border-white/10">
            <h2 className="text-xl font-semibold tracking-tight text-gray-900 dark:text-white">
              Product list
            </h2>
            <p className="inline-flex items-center gap-2 rounded-lg border border-gray-300 px-3 py-1.5 text-xs text-gray-600 dark:border-white/15 dark:text-gray-300">
              <RefreshCw aria-hidden="true" className="size-3.5" />
              Last updated <time dateTime="2024-02-28">28 Feb 2024</time>
            </p>
          </div>

          <div className="grid gap-5 p-5 lg:grid-cols-[17rem_minmax(0,1fr)]">
            {/* ── Filter rail ─────────────────────────────────────────── */}
            <form
              aria-label="Filter products"
              onSubmit={(event) => event.preventDefault()}
              className="h-fit rounded-xl border border-gray-200 p-4 dark:border-white/10"
            >
              <div className="mb-4 flex items-center justify-between gap-2">
                <h3 className="text-base font-semibold text-gray-900 dark:text-white">Filters</h3>
                <button
                  type="button"
                  onClick={clearAll}
                  className="text-xs font-medium text-orange-700 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:text-orange-400"
                >
                  Clear all
                </button>
              </div>

              <fieldset className="mb-5 border-t border-gray-200 pt-4 dark:border-white/10">
                <Legend onClear={() => setTypes([])}>Type</Legend>
                <div className="flex flex-wrap gap-x-4 gap-y-2">
                  {TYPES.map((type) => (
                    <label
                      key={type}
                      className="flex min-h-9 cursor-pointer items-center gap-2 text-sm text-gray-700 dark:text-gray-200"
                    >
                      <input
                        type="checkbox"
                        checked={types.includes(type)}
                        onChange={() => toggle(types, setTypes, type)}
                        className="size-4 accent-orange-600"
                      />
                      {type}
                    </label>
                  ))}
                </div>
              </fieldset>

              <fieldset className="mb-5 border-t border-gray-200 pt-4 dark:border-white/10">
                <Legend onClear={() => setBrands([])}>Brands</Legend>
                {/* A scroll region with a keyboard-reachable container, so the
                    brands past the fold are not unreachable without a mouse. */}
                <div
                  tabIndex={0}
                  className="max-h-40 space-y-1 overflow-y-auto pr-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                >
                  {BRANDS.map((brand) => {
                    const count = products.filter((product) => product.brand === brand).length
                    return (
                      <label
                        key={brand}
                        className="flex min-h-9 cursor-pointer items-center gap-2 text-sm text-gray-700 dark:text-gray-200"
                      >
                        <input
                          type="checkbox"
                          checked={brands.includes(brand)}
                          onChange={() => toggle(brands, setBrands, brand)}
                          className="size-4 accent-orange-600"
                        />
                        {brand}
                        <span className="text-gray-500 dark:text-gray-400">({count})</span>
                      </label>
                    )
                  })}
                </div>
              </fieldset>

              <fieldset className="mb-5 border-t border-gray-200 pt-4 dark:border-white/10">
                <Legend
                  onClear={() => {
                    setMinPrice('')
                    setMaxPrice('')
                  }}
                >
                  Price
                </Legend>
                <div className="flex items-end gap-3">
                  <span className="flex-1">
                    <label
                      htmlFor={`${baseId}-min`}
                      className="mb-1 block text-xs text-gray-600 dark:text-gray-300"
                    >
                      From
                    </label>
                    <input
                      id={`${baseId}-min`}
                      type="number"
                      min={0}
                      inputMode="numeric"
                      value={minPrice}
                      onChange={(event) => {
                        setMinPrice(event.target.value)
                        setPage(1)
                      }}
                      placeholder="0"
                      className="h-10 w-full rounded-lg border border-gray-300 px-2.5 text-sm text-gray-900 focus:border-orange-700 focus:outline-none focus:ring-1 focus:ring-orange-700 dark:border-white/15 dark:bg-transparent dark:text-white"
                    />
                  </span>
                  <span className="flex-1">
                    <label
                      htmlFor={`${baseId}-max`}
                      className="mb-1 block text-xs text-gray-600 dark:text-gray-300"
                    >
                      To
                    </label>
                    <input
                      id={`${baseId}-max`}
                      type="number"
                      min={0}
                      inputMode="numeric"
                      value={maxPrice}
                      onChange={(event) => {
                        setMaxPrice(event.target.value)
                        setPage(1)
                      }}
                      placeholder="2000"
                      className="h-10 w-full rounded-lg border border-gray-300 px-2.5 text-sm text-gray-900 focus:border-orange-700 focus:outline-none focus:ring-1 focus:ring-orange-700 dark:border-white/15 dark:bg-transparent dark:text-white"
                    />
                  </span>
                </div>
                {/* Decoration: the shape of the distribution, behind controls
                    that are operable. */}
                <div aria-hidden="true" className="mt-3 flex h-10 items-end gap-0.5">
                  {HISTOGRAM.map((height, index) => (
                    <span
                      key={index}
                      style={{ height: `${height}%` }}
                      className="flex-1 rounded-sm bg-gray-200 dark:bg-white/10"
                    />
                  ))}
                </div>
              </fieldset>

              <fieldset className="mb-5 border-t border-gray-200 pt-4 dark:border-white/10">
                <Legend onClear={() => setColours([])}>Colour</Legend>
                <div className="flex flex-wrap gap-2">
                  {COLOURS.map((colour) => {
                    const checked = colours.includes(colour)
                    return (
                      <label
                        key={colour}
                        className="group relative cursor-pointer"
                        title={colour}
                      >
                        {/* The input is visually hidden but focusable, so the
                            ring below tracks real focus rather than hover. */}
                        <input
                          type="checkbox"
                          checked={checked}
                          onChange={() => toggle(colours, setColours, colour)}
                          className="peer sr-only"
                        />
                        <span
                          aria-hidden="true"
                          className={`block size-7 rounded-full border border-gray-300 ring-offset-2 peer-focus-visible:ring-2 peer-focus-visible:ring-orange-700 dark:border-white/20 dark:ring-offset-gray-900 ${SWATCH[colour]} ${
                            checked ? 'ring-2 ring-orange-700' : ''
                          }`}
                        />
                        {/* The name a reader gets. "The blue one" is three of
                            these swatches. */}
                        <span className="sr-only">{colour}</span>
                      </label>
                    )
                  })}
                </div>
              </fieldset>

              <fieldset className="border-t border-gray-200 pt-4 dark:border-white/10">
                <Legend>Availability</Legend>
                <div className="space-y-1">
                  {TABS.map((item) => (
                    <label
                      key={item}
                      className="flex min-h-9 cursor-pointer items-center gap-2 text-sm text-gray-700 dark:text-gray-200"
                    >
                      <input
                        type="radio"
                        name={`${baseId}-availability`}
                        checked={tab === item}
                        onChange={() => {
                          setTab(item)
                          setPage(1)
                        }}
                        className="size-4 accent-orange-600"
                      />
                      {item === 'All' ? 'Any availability' : item}
                    </label>
                  ))}
                </div>
              </fieldset>
            </form>

            {/* ── Grid ────────────────────────────────────────────────── */}
            <div>
              <div className="mb-4 flex flex-wrap items-center gap-3">
                <div className="flex rounded-lg border border-gray-300 p-1 dark:border-white/15">
                  {TABS.map((item) => (
                    <button
                      key={item}
                      type="button"
                      onClick={() => {
                        setTab(item)
                        setPage(1)
                      }}
                      aria-pressed={tab === item}
                      className={`min-h-9 rounded-md px-3 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 ${
                        tab === item
                          ? 'bg-gray-900 font-medium text-white dark:bg-white dark:text-gray-900'
                          : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-white/10'
                      }`}
                    >
                      {item}
                    </button>
                  ))}
                </div>

                <div className="relative min-w-0 flex-1">
                  <Search
                    aria-hidden="true"
                    className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-gray-400"
                  />
                  <label htmlFor={`${baseId}-search`} className="sr-only">
                    Search products by name
                  </label>
                  <input
                    id={`${baseId}-search`}
                    type="search"
                    value={query}
                    onChange={(event) => {
                      setQuery(event.target.value)
                      setPage(1)
                    }}
                    placeholder="Search product"
                    className="h-10 w-full rounded-lg border border-gray-300 pl-9 pr-3 text-sm text-gray-900 placeholder:text-gray-500 focus:border-orange-700 focus:outline-none focus:ring-1 focus:ring-orange-700 dark:border-white/15 dark:bg-transparent dark:text-white dark:placeholder:text-gray-400"
                  />
                </div>

                <p className="inline-flex min-h-10 items-center gap-2 rounded-lg border border-gray-300 px-3 text-sm text-gray-700 dark:border-white/15 dark:text-gray-200">
                  <SlidersHorizontal aria-hidden="true" className="size-4" />
                  {activeFilters} active
                </p>

                <button
                  type="button"
                  className="inline-flex min-h-10 items-center gap-2 rounded-lg bg-gray-900 px-3.5 text-sm font-semibold text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
                >
                  <Plus aria-hidden="true" className="size-4" />
                  New product
                </button>
              </div>

              {/* Ticking a checkbox changes the grid without moving focus. */}
              <p role="status" aria-live="polite" className="sr-only">
                {results.length} products match. Page {safePage} of {pageCount}.
              </p>

              {visible.length === 0 ? (
                <div className="rounded-xl border border-dashed border-gray-300 px-6 py-20 text-center dark:border-white/15">
                  <p className="font-medium text-gray-900 dark:text-white">No products match</p>
                  <p className="mt-1 text-sm text-gray-600 dark:text-gray-300">
                    Try removing a filter, or widening the price range.
                  </p>
                  <button
                    type="button"
                    onClick={clearAll}
                    className="mt-4 inline-flex min-h-10 items-center rounded-lg border border-gray-300 px-4 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                  >
                    Clear all filters
                  </button>
                </div>
              ) : (
                <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
                  {visible.map((product) => (
                    <li
                      key={product.id}
                      className="group relative rounded-xl border border-gray-200 p-3 hover:border-gray-300 dark:border-white/10 dark:hover:border-white/20"
                    >
                      <div className="mb-3 flex items-start justify-between">
                        <span
                          className={`rounded-full px-2 py-0.5 text-[11px] font-medium ${
                            product.active
                              ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-400/15 dark:text-emerald-300'
                              : 'bg-gray-200 text-gray-700 dark:bg-white/10 dark:text-gray-300'
                          }`}
                        >
                          {product.active ? 'Active' : 'Draft'}
                        </span>
                        {/* Above the stretched link, so it is clickable. */}
                        <button
                          type="button"
                          className="relative z-10 grid size-8 place-items-center rounded-md text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 dark:text-gray-400 dark:hover:bg-white/10"
                        >
                          <span aria-hidden="true" className="text-lg leading-none">
                            &hellip;
                          </span>
                          <span className="sr-only">More actions for {product.name}</span>
                        </button>
                      </div>

                      {/* Placeholder art. Swap for an <img> with real alt text
                          describing the product, not the word "product". */}
                      <span
                        aria-hidden="true"
                        className={`mb-3 block h-28 rounded-lg bg-gradient-to-br ${product.swatch}`}
                      />

                      <h3 className="text-sm font-medium text-gray-900 dark:text-white">
                        {/* One link, stretched over the tile. Two links to the
                            same product is two entries in a link list. */}
                        <a
                          href="#"
                          className="after:absolute after:inset-0 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700"
                        >
                          {product.name}
                        </a>
                      </h3>

                      <p className="mt-1 flex items-center justify-between text-sm">
                        <span className="font-semibold text-gray-900 tabular-nums dark:text-white">
                          {money.format(product.price / 100)}
                        </span>
                        <span className="flex items-center gap-1 text-gray-700 dark:text-gray-200">
                          <span aria-hidden="true" className="text-amber-500">
                            &#9733;
                          </span>
                          {/* The number in text. A star is a picture of it. */}
                          <span className="tabular-nums">{product.rating.toFixed(1)}</span>
                          <span className="sr-only">out of 5</span>
                        </span>
                      </p>
                    </li>
                  ))}
                </ul>
              )}

              <div className="mt-5 flex flex-wrap items-center justify-between gap-3">
                <p className="text-sm text-gray-600 dark:text-gray-300">
                  Showing {visible.length} of {results.length}
                </p>
                <nav aria-label="Product pages">
                  <ul className="flex items-center gap-1">
                    <li>
                      <button
                        type="button"
                        onClick={() => setPage(Math.max(1, safePage - 1))}
                        disabled={safePage === 1}
                        className="inline-flex min-h-9 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                      >
                        Previous
                      </button>
                    </li>
                    {Array.from({ length: pageCount }).map((_, index) => {
                      const number = index + 1
                      const current = number === safePage
                      return (
                        <li key={number}>
                          <button
                            type="button"
                            onClick={() => setPage(number)}
                            aria-current={current ? 'page' : undefined}
                            className={`grid size-9 place-items-center rounded-lg text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 ${
                              current
                                ? 'bg-gray-900 font-semibold text-white dark:bg-white dark:text-gray-900'
                                : 'border border-gray-300 text-gray-700 hover:bg-gray-50 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5'
                            }`}
                          >
                            {number}
                            <span className="sr-only">
                              {current ? ', current page' : `, go to page ${number}`}
                            </span>
                          </button>
                        </li>
                      )
                    })}
                    <li>
                      <button
                        type="button"
                        onClick={() => setPage(Math.min(pageCount, safePage + 1))}
                        disabled={safePage === pageCount}
                        className="inline-flex min-h-9 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-700 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                      >
                        Next
                      </button>
                    </li>
                  </ul>
                </nav>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
