'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import {
  Activity,
  AlertTriangle,
  ArrowUpDown,
  Award,
  Columns3,
  Download,
  HelpCircle,
  Megaphone,
  MessageSquare,
  MoreHorizontal,
  Package,
  Receipt,
  Search,
  Settings,
  SlidersHorizontal,
  Upload,
  XCircle,
} from 'lucide-react'

/*
 * A dense inventory table with a bulk-action bar that appears on selection.
 *
 * Self-contained. Sorting, selection, filtering and paging all work, because
 * a bulk-action bar that never appears cannot be reviewed and a stock badge
 * with no data behind it cannot be checked for contrast against a real value.
 *
 * ── The bulk bar ────────────────────────────────────────────────────────────
 * The hard part of this pattern is not the styling, it is that the bar appears
 * out of nowhere, usually over the content, and a keyboard user never learns
 * it exists. Three things fix that here. It is a labelled region, so it turns
 * up in a landmark list. Its arrival is announced through the status region,
 * because selecting a row does not move focus. And it sits in the flow at the
 * bottom of the card rather than floating over rows, so it cannot cover the
 * row you are working on.
 *
 * Destructive and reversible actions are not the same colour. Delete is the
 * only red thing in the bar.
 *
 * ── Stock ───────────────────────────────────────────────────────────────────
 * Three states, and each says what it is in words next to its icon and colour.
 * "Low stock" beside an amber dot is legible to everyone; an amber dot alone
 * is legible to people who can see it and can remember what amber means here.
 * The number of units is in the same cell, so the badge is a summary rather
 * than the only information.
 *
 * ── Sorting ─────────────────────────────────────────────────────────────────
 * aria-sort on the <th>, a <button> inside it, and the attribute set on
 * exactly one column. A <th> with an onClick is a cell nobody can press.
 *
 * Money is minor units, formatted once with Intl, in the currency the store
 * actually sells in rather than a hardcoded dollar sign.
 */

export interface InventoryProduct {
  id: string
  name: string
  sku: string
  category: string
  brand: string
  /** Minor units of the store currency. */
  price: number
  stock: number
  /** Tailwind classes for the placeholder thumbnail. Swap for a real <img>. */
  swatch: string
}

const PRODUCTS: InventoryProduct[] = [
  { id: 'g1', name: 'Gourmia GCM3100 Coffee Maker', sku: '2348678', category: 'Coffee machines', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-red-300 to-red-500' },
  { id: 'g2', name: 'Gourmia GAF7020G Panorama 7-quart Air Fryer', sku: '2348679', category: 'Air fryers', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-slate-300 to-slate-500' },
  { id: 'g3', name: 'Gourmia GSI1020 2-quart Auto Ice Cream Maker', sku: '2348680', category: 'Ice cream makers', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-stone-300 to-stone-500' },
  { id: 'g4', name: 'Gourmia GCM4225 15-bar Espresso Machine', sku: '2348681', category: 'Coffee machines', brand: 'Gourmia', price: 300000, stock: 0, swatch: 'from-zinc-400 to-zinc-600' },
  { id: 'g5', name: 'Gourmia GAF823W 8-quart Digital Air Fryer', sku: '2348682', category: 'Air fryers', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-amber-200 to-amber-400' },
  { id: 'g6', name: 'Gourmia GPM1270 All-in-one Pizza Oven', sku: '2348683', category: 'Pizza ovens', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-orange-300 to-orange-500' },
  { id: 'g7', name: 'Gourmia GCM6500 Espresso and Cappuccino Maker', sku: '2348684', category: 'Coffee machines', brand: 'Gourmia', price: 300000, stock: 50, swatch: 'from-neutral-400 to-neutral-600' },
  { id: 'g8', name: 'Gourmia GK220 1.8-quart Cordless Electric Kettle', sku: '2348685', category: 'Kettles and tea brewers', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-violet-200 to-violet-400' },
  { id: 'g9', name: 'Gourmia GDK368 Digital 3-function Oven', sku: '2348686', category: 'Rice and grain', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-gray-500 to-gray-700' },
  { id: 'g10', name: 'Gourmia GAF680 Digital Free Fry Air Fryer', sku: '2348687', category: 'Air fryers', brand: 'Gourmia', price: 300000, stock: 500, swatch: 'from-sky-200 to-sky-400' },
]

const NAV = [
  { label: 'Analytics', icon: Activity, badge: null },
  { label: 'Transactions', icon: Receipt, badge: '99+' },
  { label: 'Products', icon: Package, badge: null, current: true },
  { label: 'Promotions', icon: Megaphone, badge: null },
  { label: 'Loyalty', icon: Award, badge: null },
  { label: 'Ticketing', icon: MessageSquare, badge: '2' },
]

type SortKey = 'name' | 'category' | 'brand' | 'price' | 'stock'
type Status = 'In stock' | 'Low stock' | 'Out of stock'

const COLUMNS: { key: SortKey; label: string; numeric?: boolean }[] = [
  { key: 'name', label: 'Product name' },
  { key: 'category', label: 'Category' },
  { key: 'brand', label: 'Brand' },
  { key: 'price', label: 'Price', numeric: true },
  { key: 'stock', label: 'Stock', numeric: true },
]

/* Minor units, and the currency the store sells in. A hardcoded "$" in front
   of a number formatted for another locale is how prices end up wrong. */
const money = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  maximumFractionDigits: 0,
})

function statusOf(stock: number): Status {
  if (stock === 0) return 'Out of stock'
  if (stock <= 50) return 'Low stock'
  return 'In stock'
}

const STATUS_TONE: Record<Status, string> = {
  'In stock': 'bg-emerald-100 text-emerald-800 dark:bg-emerald-400/15 dark:text-emerald-300',
  'Low stock': 'bg-amber-100 text-amber-800 dark:bg-amber-400/15 dark:text-amber-300',
  'Out of stock': 'bg-rose-100 text-rose-800 dark:bg-rose-400/15 dark:text-rose-300',
}

const PAGE_SIZE = 10

export default function ProductTableWithBulkActions({
  products = PRODUCTS,
}: {
  products?: InventoryProduct[]
}) {
  const [archived, setArchived] = useState(false)
  const [query, setQuery] = useState('')
  const [sortKey, setSortKey] = useState<SortKey>('name')
  const [ascending, setAscending] = useState(true)
  const [selected, setSelected] = useState<string[]>([])
  const [page, setPage] = useState(1)
  const [announcement, setAnnouncement] = useState('')

  const selectAllRef = useRef<HTMLInputElement>(null)

  const sorted = useMemo(() => {
    const filtered = products.filter((product) =>
      query
        ? product.name.toLowerCase().includes(query.toLowerCase()) || product.sku.includes(query)
        : true,
    )
    return [...filtered].sort((a, b) => {
      const factor = ascending ? 1 : -1
      if (sortKey === 'price' || sortKey === 'stock') return (a[sortKey] - b[sortKey]) * factor
      return String(a[sortKey]).localeCompare(String(b[sortKey])) * factor
    })
  }, [products, query, sortKey, ascending])

  const pageCount = Math.max(1, Math.ceil(sorted.length / PAGE_SIZE))
  const safePage = Math.min(page, pageCount)
  const visible = sorted.slice((safePage - 1) * PAGE_SIZE, safePage * PAGE_SIZE)

  const visibleIds = visible.map((product) => product.id)
  const allSelected = visibleIds.length > 0 && visibleIds.every((id) => selected.includes(id))
  const someSelected = selected.length > 0 && !allSelected

  useEffect(() => {
    if (selectAllRef.current) selectAllRef.current.indeterminate = someSelected
  }, [someSelected])

  function announce(message: string) {
    setAnnouncement('')
    requestAnimationFrame(() => setAnnouncement(message))
  }

  function toggleSort(key: SortKey) {
    const nextAscending = key === sortKey ? !ascending : true
    setSortKey(key)
    setAscending(nextAscending)
    setPage(1)
    const label = COLUMNS.find((column) => column.key === key)?.label ?? key
    announce(`Sorted by ${label}, ${nextAscending ? 'ascending' : 'descending'}.`)
  }

  function toggleRow(product: InventoryProduct) {
    const next = selected.includes(product.id)
      ? selected.filter((id) => id !== product.id)
      : [...selected, product.id]
    setSelected(next)
    /* Selecting a row does not move focus, and the bulk bar appears out of
       nowhere, so its arrival has to be said. */
    announce(
      next.length === 0
        ? 'Selection cleared.'
        : `${next.length} ${next.length === 1 ? 'product' : 'products'} selected. Bulk actions available.`,
    )
  }

  function toggleAll() {
    const next = allSelected
      ? selected.filter((id) => !visibleIds.includes(id))
      : [...new Set([...selected, ...visibleIds])]
    setSelected(next)
    announce(
      next.length === 0
        ? 'Selection cleared.'
        : `${next.length} products selected. Bulk actions available.`,
    )
  }

  return (
    <div className="flex min-h-[44rem] bg-indigo-50/60 text-gray-900 dark:bg-gray-950 dark:text-white">
      <p role="status" aria-live="polite" className="sr-only">
        {announcement}
      </p>

      {/* ── Sidebar ──────────────────────────────────────────────────────── */}
      <div className="hidden w-56 shrink-0 flex-col border-r border-gray-200 bg-white lg:flex dark:border-white/10 dark:bg-gray-900">
        <p className="flex items-center gap-2 px-4 py-5 text-base font-semibold">
          <span
            aria-hidden="true"
            className="grid size-7 place-items-center rounded-lg bg-indigo-700 text-xs font-bold text-white"
          >
            A
          </span>
          Northwind
        </p>

        <nav aria-label="Main" className="flex-1 px-2">
          <ul className="space-y-0.5">
            {NAV.map((item) => {
              const Icon = item.icon
              return (
              <li key={item.label}>
                <a
                  href="#"
                  aria-current={item.current ? 'page' : undefined}
                  className={`flex min-h-10 items-center justify-between gap-2 rounded-lg px-2.5 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 ${
                    item.current
                      ? 'bg-indigo-50 font-medium text-indigo-800 dark:bg-indigo-500/15 dark:text-indigo-200'
                      : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-white/10'
                  }`}
                >
                  <span className="flex items-center gap-2.5">
                    <Icon aria-hidden="true" className="size-4 shrink-0" />
                    {item.label}
                  </span>
                  {item.badge && (
                    <span className="rounded-full bg-rose-600 px-1.5 py-0.5 text-[10px] font-semibold text-white">
                      {item.badge}
                      {/* The number alone is ambiguous: 2 what? */}
                      <span className="sr-only"> unread</span>
                    </span>
                  )}
                </a>
              </li>
              )
            })}
          </ul>
        </nav>

        <div className="space-y-0.5 border-t border-gray-200 p-2 dark:border-white/10">
          {[
            { icon: HelpCircle, label: 'Help' },
            { icon: Settings, label: 'Settings' },
          ].map(({ icon: Icon, label }) => (
            <a
              key={label}
              href="#"
              className="flex min-h-10 items-center gap-2.5 rounded-lg px-2.5 text-sm text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:text-gray-200 dark:hover:bg-white/10"
            >
              <Icon aria-hidden="true" className="size-4" />
              {label}
            </a>
          ))}
          <p className="mt-2 flex items-center gap-2 rounded-lg bg-gray-50 p-2 dark:bg-white/5">
            <span
              aria-hidden="true"
              className="grid size-8 shrink-0 place-items-center rounded-full bg-amber-600 text-xs font-semibold text-white"
            >
              CW
            </span>
            <span className="min-w-0">
              <span className="block truncate text-sm font-medium">Craig Westervelt</span>
              <span className="block truncate text-xs text-gray-500 dark:text-gray-400">
                craig@northwind.example
              </span>
            </span>
          </p>
        </div>
      </div>

      {/* ── Main ─────────────────────────────────────────────────────────── */}
      <div className="flex min-w-0 flex-1 flex-col p-5">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
          <h1 className="text-2xl font-semibold tracking-tight">Products</h1>
          <div className="flex flex-wrap items-center gap-2">
            <button
              type="button"
              className="inline-flex min-h-10 items-center gap-2 rounded-lg border border-indigo-200 bg-white px-3 text-sm font-medium text-indigo-800 hover:bg-indigo-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:border-indigo-400/30 dark:bg-transparent dark:text-indigo-200 dark:hover:bg-indigo-500/10"
            >
              <Upload aria-hidden="true" className="size-4" />
              Import
            </button>
            <button
              type="button"
              className="inline-flex min-h-10 items-center gap-2 rounded-lg border border-indigo-200 bg-white px-3 text-sm font-medium text-indigo-800 hover:bg-indigo-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:border-indigo-400/30 dark:bg-transparent dark:text-indigo-200 dark:hover:bg-indigo-500/10"
            >
              <Download aria-hidden="true" className="size-4" />
              Export
            </button>
            <button
              type="button"
              className="inline-flex min-h-10 items-center gap-2 rounded-lg bg-indigo-700 px-3.5 text-sm font-semibold text-white hover:bg-indigo-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700"
            >
              <span aria-hidden="true" className="text-base leading-none">
                +
              </span>
              Create product
            </button>
          </div>
        </div>

        <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
          <div className="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 p-3 dark:border-white/10">
            <div className="flex rounded-lg bg-gray-100 p-1 dark:bg-white/5">
              {[
                { label: 'Published', value: false },
                { label: 'Archived', value: true },
              ].map((option) => (
                <button
                  key={option.label}
                  type="button"
                  onClick={() => {
                    setArchived(option.value)
                    setSelected([])
                    setPage(1)
                  }}
                  aria-pressed={archived === option.value}
                  className={`min-h-9 rounded-md px-3 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 ${
                    archived === option.value
                      ? 'bg-indigo-100 font-medium text-indigo-800 dark:bg-indigo-500/20 dark:text-indigo-200'
                      : 'text-gray-600 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white'
                  }`}
                >
                  {option.label}
                </button>
              ))}
            </div>

            <div className="flex flex-1 items-center justify-end gap-2">
              <div className="relative min-w-0 max-w-64 flex-1">
                <Search
                  aria-hidden="true"
                  className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-gray-400"
                />
                <label htmlFor="inventory-search" className="sr-only">
                  Search products by name or SKU
                </label>
                <input
                  id="inventory-search"
                  type="search"
                  value={query}
                  onChange={(event) => {
                    setQuery(event.target.value)
                    setPage(1)
                  }}
                  placeholder="Search"
                  className="h-10 w-full rounded-lg border border-gray-300 pl-9 pr-3 text-sm text-gray-900 placeholder:text-gray-500 focus:border-indigo-700 focus:outline-none focus:ring-1 focus:ring-indigo-700 dark:border-white/15 dark:bg-transparent dark:text-white dark:placeholder:text-gray-400"
                />
              </div>
              {[
                { icon: Columns3, label: 'Choose visible columns' },
                { icon: SlidersHorizontal, label: 'Filter products' },
              ].map(({ icon: Icon, label }) => (
                <button
                  key={label}
                  type="button"
                  className="grid size-10 shrink-0 place-items-center rounded-lg border border-gray-300 text-gray-600 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:border-white/15 dark:text-gray-300 dark:hover:bg-white/5"
                >
                  <Icon aria-hidden="true" className="size-4" />
                  <span className="sr-only">{label}</span>
                </button>
              ))}
            </div>
          </div>

          <div className="min-h-0 flex-1 overflow-auto">
            <table className="w-full text-left text-sm">
              <caption className="sr-only">
                {archived ? 'Archived' : 'Published'} products, {sorted.length} rows, sortable
              </caption>
              <thead className="sticky top-0 z-10 border-b border-gray-200 bg-white text-gray-600 dark:border-white/10 dark:bg-gray-900 dark:text-gray-300">
                <tr>
                  <th scope="col" className="w-12 px-3 py-2.5">
                    <input
                      ref={selectAllRef}
                      type="checkbox"
                      checked={allSelected}
                      onChange={toggleAll}
                      aria-label={`Select all ${visible.length} products on this page`}
                      className="size-4 accent-indigo-700"
                    />
                  </th>
                  {COLUMNS.map((column) => {
                    const active = column.key === sortKey
                    return (
                      <th
                        key={column.key}
                        scope="col"
                        aria-sort={active ? (ascending ? 'ascending' : 'descending') : 'none'}
                        className={`px-3 py-2.5 text-xs font-medium uppercase tracking-wide ${column.numeric ? 'text-right' : ''}`}
                      >
                        <button
                          type="button"
                          onClick={() => toggleSort(column.key)}
                          className={`inline-flex min-h-9 items-center gap-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 ${
                            column.numeric ? 'flex-row-reverse' : ''
                          } ${active ? 'text-gray-900 dark:text-white' : ''}`}
                        >
                          {column.label}
                          <ArrowUpDown aria-hidden="true" className="size-3 text-gray-400" />
                          <span className="sr-only">
                            {active
                              ? ascending
                                ? ', sorted ascending, activate to sort descending'
                                : ', sorted descending, activate to sort ascending'
                              : ', activate to sort ascending'}
                          </span>
                        </button>
                      </th>
                    )
                  })}
                  <th scope="col" className="px-3 py-2.5 text-xs font-medium uppercase tracking-wide">
                    Status
                  </th>
                  <th scope="col" className="w-12 px-3 py-2.5">
                    <span className="sr-only">Row actions</span>
                  </th>
                </tr>
              </thead>

              <tbody className="divide-y divide-gray-100 dark:divide-white/5">
                {visible.map((product) => {
                  const checked = selected.includes(product.id)
                  const status = statusOf(product.stock)
                  return (
                    <tr
                      key={product.id}
                      className={checked ? 'bg-indigo-50/70 dark:bg-indigo-500/10' : ''}
                    >
                      <td className="px-3 py-2.5">
                        <input
                          type="checkbox"
                          checked={checked}
                          onChange={() => toggleRow(product)}
                          aria-label={`Select ${product.name}`}
                          className="size-4 accent-indigo-700"
                        />
                      </td>

                      <th scope="row" className="px-3 py-2.5 font-normal">
                        <span className="flex items-center gap-3">
                          <span
                            aria-hidden="true"
                            className={`size-9 shrink-0 rounded-md bg-gradient-to-br ${product.swatch}`}
                          />
                          <span className="min-w-0">
                            <span className="block max-w-72 truncate font-medium text-gray-900 dark:text-white">
                              {product.name}
                            </span>
                            <span className="block text-xs text-gray-500 dark:text-gray-400">
                              SN: {product.sku}
                            </span>
                          </span>
                        </span>
                      </th>

                      <td className="whitespace-nowrap px-3 py-2.5 text-gray-600 dark:text-gray-300">
                        {product.category}
                      </td>
                      <td className="whitespace-nowrap px-3 py-2.5 text-gray-600 dark:text-gray-300">
                        {product.brand}
                      </td>
                      <td className="whitespace-nowrap px-3 py-2.5 text-right tabular-nums text-gray-900 dark:text-white">
                        {money.format(product.price)}
                      </td>
                      <td className="whitespace-nowrap px-3 py-2.5 text-right">
                        <span
                          className={`inline-flex items-center gap-1.5 tabular-nums ${
                            status === 'Out of stock'
                              ? 'text-rose-700 dark:text-rose-300'
                              : status === 'Low stock'
                                ? 'text-amber-700 dark:text-amber-300'
                                : 'text-gray-900 dark:text-white'
                          }`}
                        >
                          {status === 'Out of stock' && (
                            <XCircle aria-hidden="true" className="size-3.5" />
                          )}
                          {status === 'Low stock' && (
                            <AlertTriangle aria-hidden="true" className="size-3.5" />
                          )}
                          {product.stock}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-3 py-2.5">
                        {/* The words are the message. The colour agrees with
                            them; it does not replace them. */}
                        <span
                          className={`inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${STATUS_TONE[status]}`}
                        >
                          {status}
                        </span>
                      </td>
                      <td className="px-3 py-2.5 text-right">
                        <button
                          type="button"
                          className="grid size-9 place-items-center rounded-lg text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:text-gray-400 dark:hover:bg-white/10"
                        >
                          <MoreHorizontal aria-hidden="true" className="size-4" />
                          <span className="sr-only">More actions for {product.name}</span>
                        </button>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>

            {visible.length === 0 && (
              <p className="px-4 py-16 text-center text-sm text-gray-600 dark:text-gray-300">
                No products match that search.
              </p>
            )}
          </div>

          {/* In the flow at the foot of the card, not floating over the rows:
              a bar that covers the row you are editing is a bar in the way. */}
          {selected.length > 0 && (
            <section
              aria-label="Bulk actions"
              className="flex flex-wrap items-center gap-3 border-t border-gray-200 bg-gray-50 px-3 py-2.5 dark:border-white/10 dark:bg-white/5"
            >
              <p className="text-sm font-medium">
                {selected.length} {selected.length === 1 ? 'product' : 'products'} selected
              </p>
              <div className="flex flex-wrap items-center gap-2">
                <button
                  type="button"
                  className="min-h-9 rounded-lg border border-gray-300 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:border-white/15 dark:bg-transparent dark:text-gray-200 dark:hover:bg-white/10"
                >
                  Edit
                  <span className="sr-only"> the {selected.length} selected products</span>
                </button>
                <button
                  type="button"
                  className="min-h-9 rounded-lg border border-gray-300 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:border-white/15 dark:bg-transparent dark:text-gray-200 dark:hover:bg-white/10"
                >
                  Archive
                  <span className="sr-only"> the {selected.length} selected products</span>
                </button>
                {/* The only red control in the bar. If everything is red then
                    nothing is dangerous. */}
                <button
                  type="button"
                  className="min-h-9 rounded-lg border border-rose-300 bg-white px-3 text-sm font-medium text-rose-700 hover:bg-rose-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-rose-700 dark:border-rose-400/30 dark:bg-transparent dark:text-rose-300 dark:hover:bg-rose-500/10"
                >
                  Delete
                  <span className="sr-only"> the {selected.length} selected products</span>
                </button>
              </div>
              <button
                type="button"
                onClick={() => {
                  setSelected([])
                  announce('Selection cleared.')
                }}
                className="ml-auto min-h-9 rounded-lg px-3 text-sm font-medium text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 dark:text-gray-200 dark:hover:bg-white/10"
              >
                Clear selection
              </button>
            </section>
          )}

          <nav
            aria-label="Product pages"
            className="flex flex-wrap items-center justify-between gap-3 border-t border-gray-200 px-3 py-2.5 dark:border-white/10"
          >
            <p className="text-sm text-gray-600 dark:text-gray-300">
              Showing {visible.length} of {sorted.length} products
            </p>
            <ul className="flex items-center gap-1">
              <li>
                <button
                  type="button"
                  onClick={() => setPage(Math.max(1, safePage - 1))}
                  disabled={safePage === 1}
                  className="inline-flex min-h-9 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
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
                      className={`grid size-9 place-items-center rounded-lg text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 ${
                        current
                          ? 'bg-indigo-700 font-semibold text-white'
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
                  className="inline-flex min-h-9 items-center rounded-lg border border-gray-300 px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-700 disabled:opacity-40 dark:border-white/15 dark:text-gray-200 dark:hover:bg-white/5"
                >
                  Next
                </button>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </div>
  )
}
