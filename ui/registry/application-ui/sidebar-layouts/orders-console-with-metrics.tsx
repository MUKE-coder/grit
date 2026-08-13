'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import {
  ArrowUpDown,
  BarChart3,
  Bell,
  ChevronDown,
  ChevronRight,
  CreditCard,
  Download,
  FileText,
  Home,
  Megaphone,
  MessageSquare,
  MoreHorizontal,
  Package,
  Plus,
  Search,
  Settings,
  ShoppingBag,
  ShoppingCart,
  SlidersHorizontal,
  Store,
  Users,
  Wallet,
} from 'lucide-react'
import type { LucideIcon } from 'lucide-react'

/*
 * A commerce back-office shell: collapsible sidebar, metric strip, filter tabs
 * and an orders table with selection.
 *
 * Self-contained. The obvious build is a router, a table library and a set of
 * primitives, which is three dependencies for a block whose job is to be a
 * starting point you own.
 *
 * ── The sidebar ─────────────────────────────────────────────────────────────
 * Sections are native <details>, so expanding a group needs no state, works
 * before hydration and gets keyboard support from the browser. The current
 * page carries aria-current="page", which is the only thing that tells a
 * screen reader which of nineteen links you are actually on. Colour alone does
 * not.
 *
 * The collapse control is a real toggle with aria-expanded and aria-controls,
 * not an icon that silently changes a width.
 *
 * ── The metrics ─────────────────────────────────────────────────────────────
 * A <dl>. Each figure's change is written out ("up 25.2% on last week") rather
 * than left to a green triangle, because a triangle carries no direction to
 * anyone not looking at it. The sparklines are aria-hidden: they repeat the
 * number beside them and cannot be read.
 *
 * ── The tabs ────────────────────────────────────────────────────────────────
 * A real tablist with roving focus and arrow-key navigation. A row of buttons
 * that filter a table is the most commonly faked tablist there is, and the
 * fake version leaves a keyboard user tabbing through every filter to reach
 * the table.
 *
 * ── The table ───────────────────────────────────────────────────────────────
 * aria-sort on the sorted column only, sort controls as buttons inside the
 * <th>, per-row checkbox labels naming their row, and one status region for
 * sorting, filtering and selection, none of which moves focus.
 *
 * Money is integer cents, formatted once with Intl.
 */

export interface Order {
  id: string
  /** ISO date. */
  date: string
  customer: string
  payment: 'Pending' | 'Success'
  /** Integer cents. */
  total: number
  delivery: string | null
  items: number
  fulfilment: 'Fulfilled' | 'Unfulfilled'
}

const ORDERS: Order[] = [
  { id: '#1002', date: '2024-02-11', customer: 'Wade Warren', payment: 'Pending', total: 2000, delivery: null, items: 2, fulfilment: 'Unfulfilled' },
  { id: '#1004', date: '2024-02-13', customer: 'Esther Howard', payment: 'Success', total: 2200, delivery: null, items: 3, fulfilment: 'Fulfilled' },
  { id: '#1007', date: '2024-02-15', customer: 'Jenny Wilson', payment: 'Pending', total: 2500, delivery: null, items: 1, fulfilment: 'Unfulfilled' },
  { id: '#1009', date: '2024-02-17', customer: 'Guy Hawkins', payment: 'Success', total: 2700, delivery: null, items: 5, fulfilment: 'Fulfilled' },
  { id: '#1011', date: '2024-02-19', customer: 'Jacob Jones', payment: 'Pending', total: 3200, delivery: null, items: 4, fulfilment: 'Unfulfilled' },
  { id: '#1013', date: '2024-02-21', customer: 'Kristin Watson', payment: 'Success', total: 2500, delivery: null, items: 3, fulfilment: 'Fulfilled' },
  { id: '#1015', date: '2024-02-23', customer: 'Albert Flores', payment: 'Pending', total: 2800, delivery: null, items: 2, fulfilment: 'Unfulfilled' },
  { id: '#1018', date: '2024-02-25', customer: 'Eleanor Pena', payment: 'Success', total: 3500, delivery: null, items: 1, fulfilment: 'Fulfilled' },
  { id: '#1019', date: '2024-02-27', customer: 'Theresa Webb', payment: 'Pending', total: 2000, delivery: null, items: 2, fulfilment: 'Unfulfilled' },
]

const METRICS = [
  { label: 'Total orders', value: '21', change: 25.2, direction: 'up' as const, points: '0,18 12,14 24,16 36,9 48,12 60,5' },
  { label: 'Order items over time', value: '15', change: 18.2, direction: 'up' as const, points: '0,16 12,17 24,11 36,13 48,7 60,4' },
  { label: 'Returns orders', value: '0', change: 1.2, direction: 'down' as const, points: '0,10 12,10 24,10 36,10 48,10 60,10' },
  { label: 'Fulfilled orders over time', value: '12', change: 12.2, direction: 'up' as const, points: '0,17 12,13 24,15 36,8 48,10 60,3' },
]

const TABS = ['All', 'Unfulfilled', 'Unpaid', 'Open', 'Closed'] as const
type Tab = (typeof TABS)[number]

const NAV: { label: string; icon: LucideIcon; children?: string[]; current?: boolean }[] = [
  { label: 'Home', icon: Home },
  { label: 'Orders', icon: ShoppingBag, children: ['Drafts', 'Shipping labels', 'Abandoned checkouts'], current: true },
  { label: 'Products', icon: Package, children: ['Inventory', 'Collections', 'Gift cards'] },
  { label: 'Customers', icon: Users, children: ['Segments'] },
  { label: 'Content', icon: FileText, children: ['Metaobjects', 'Files'] },
  { label: 'Finances', icon: Wallet, children: ['Payouts', 'Billing'] },
  { label: 'Analytics', icon: BarChart3, children: ['Reports', 'Live view'] },
  { label: 'Marketing', icon: Megaphone, children: ['Campaigns', 'Automations'] },
]

const CHANNELS: { label: string; icon: LucideIcon }[] = [
  { label: 'Online store', icon: Store },
  { label: 'Point of sale', icon: CreditCard },
  { label: 'Shop', icon: ShoppingCart },
]

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })
const dates = new Intl.DateTimeFormat('en-US', { day: 'numeric', month: 'short', year: 'numeric' })

const PAYMENT_TONE: Record<Order['payment'], string> = {
  Pending: 'border-amber-300 bg-amber-50 text-amber-800 dark:border-amber-400/30 dark:bg-amber-400/10 dark:text-amber-300',
  Success: 'border-emerald-300 bg-emerald-50 text-emerald-800 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-300',
}

const FULFILMENT_TONE: Record<Order['fulfilment'], string> = {
  Fulfilled: 'border-emerald-300 bg-emerald-50 text-emerald-800 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-300',
  Unfulfilled: 'border-rose-300 bg-rose-50 text-rose-800 dark:border-rose-400/30 dark:bg-rose-400/10 dark:text-rose-300',
}

type SortKey = 'date' | 'total'

export default function OrdersConsoleWithMetrics({ orders = ORDERS }: { orders?: Order[] }) {
  const [collapsed, setCollapsed] = useState(false)
  const [tab, setTab] = useState<Tab>('All')
  const [sortKey, setSortKey] = useState<SortKey>('date')
  const [ascending, setAscending] = useState(true)
  const [selected, setSelected] = useState<string[]>([])
  const [announcement, setAnnouncement] = useState('')

  const selectAllRef = useRef<HTMLInputElement>(null)
  const tabRefs = useRef<(HTMLButtonElement | null)[]>([])

  const visible = useMemo(() => {
    const filtered = orders.filter((order) => {
      if (tab === 'Unfulfilled') return order.fulfilment === 'Unfulfilled'
      if (tab === 'Unpaid') return order.payment === 'Pending'
      if (tab === 'Open') return order.fulfilment === 'Unfulfilled'
      if (tab === 'Closed') return order.fulfilment === 'Fulfilled'
      return true
    })
    /* A copy: sort mutates, and sorting the prop in place reorders the
       caller's array. */
    return [...filtered].sort((a, b) => {
      const factor = ascending ? 1 : -1
      if (sortKey === 'total') return (a.total - b.total) * factor
      return a.date.localeCompare(b.date) * factor
    })
  }, [orders, tab, sortKey, ascending])

  const allSelected = visible.length > 0 && visible.every((order) => selected.includes(order.id))
  const someSelected = selected.length > 0 && !allSelected

  /* indeterminate is a property, not an attribute, so it cannot be written in
     JSX. Without this the box reads as unchecked while three rows are on. */
  useEffect(() => {
    if (selectAllRef.current) selectAllRef.current.indeterminate = someSelected
  }, [someSelected])

  function announce(message: string) {
    setAnnouncement('')
    requestAnimationFrame(() => setAnnouncement(message))
  }

  function chooseTab(next: Tab) {
    setTab(next)
    setSelected([])
    announce(`${next} filter applied.`)
  }

  /* Roving focus: in a tablist, Left and Right move between tabs and Tab
     leaves the group. Without this a keyboard user walks through all five
     filters on the way to the table. */
  function onTabKeyDown(event: React.KeyboardEvent, index: number) {
    const last = TABS.length - 1
    let next = index
    if (event.key === 'ArrowRight') next = index === last ? 0 : index + 1
    else if (event.key === 'ArrowLeft') next = index === 0 ? last : index - 1
    else if (event.key === 'Home') next = 0
    else if (event.key === 'End') next = last
    else return
    event.preventDefault()
    chooseTab(TABS[next])
    tabRefs.current[next]?.focus()
  }

  function toggleSort(key: SortKey) {
    const nextAscending = key === sortKey ? !ascending : true
    setSortKey(key)
    setAscending(nextAscending)
    announce(`Sorted by ${key === 'date' ? 'date' : 'total'}, ${nextAscending ? 'ascending' : 'descending'}.`)
  }

  function toggleRow(order: Order) {
    const next = selected.includes(order.id)
      ? selected.filter((id) => id !== order.id)
      : [...selected, order.id]
    setSelected(next)
    announce(`Order ${order.id} ${next.includes(order.id) ? 'selected' : 'deselected'}. ${next.length} selected.`)
  }

  function toggleAll() {
    const next = allSelected ? [] : visible.map((order) => order.id)
    setSelected(next)
    announce(allSelected ? 'Selection cleared.' : `${next.length} orders selected.`)
  }

  return (
    <div className="flex min-h-[46rem] bg-gray-50 text-gray-900 dark:bg-gray-950 dark:text-white">
      <p role="status" aria-live="polite" className="sr-only">
        {announcement}
      </p>

      {/* ── Sidebar ──────────────────────────────────────────────────────── */}
      <div
        id="console-sidebar"
        className={`hidden shrink-0 flex-col border-r border-gray-200 bg-white transition-[width] lg:flex dark:border-white/10 dark:bg-gray-900 ${
          collapsed ? 'w-16' : 'w-60'
        }`}
      >
        <div className="flex items-center justify-between gap-2 px-3 py-4">
          <span className="flex items-center gap-2 overflow-hidden">
            <span
              aria-hidden="true"
              className="grid size-7 shrink-0 place-items-center rounded-lg bg-violet-700 text-xs font-bold text-white"
            >
              S
            </span>
            {!collapsed && <span className="truncate text-sm font-semibold">ShopZen</span>}
          </span>
          <button
            type="button"
            onClick={() => setCollapsed((value) => !value)}
            aria-expanded={!collapsed}
            aria-controls="console-sidebar"
            className="grid size-8 shrink-0 place-items-center rounded-md text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-400 dark:hover:bg-white/10"
          >
            <ChevronRight
              aria-hidden="true"
              className={`size-4 transition-transform ${collapsed ? '' : 'rotate-180'}`}
            />
            <span className="sr-only">{collapsed ? 'Expand sidebar' : 'Collapse sidebar'}</span>
          </button>
        </div>

        <nav aria-label="Main" className="flex-1 overflow-y-auto px-2 pb-4">
          <ul className="space-y-0.5">
            {NAV.map((item) => {
              const Icon = item.icon
              return (
              <li key={item.label}>
                {item.children && !collapsed ? (
                  <details open={item.current} className="group">
                    <summary className="flex min-h-10 cursor-pointer list-none items-center justify-between gap-2 rounded-lg px-2.5 text-sm text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-200 dark:hover:bg-white/10">
                      <span className="flex items-center gap-2.5">
                        <Icon aria-hidden="true" className="size-4 shrink-0" />
                        {item.label}
                      </span>
                      <ChevronDown
                        aria-hidden="true"
                        className="size-3.5 text-gray-400 transition-transform group-open:rotate-180"
                      />
                    </summary>
                    <ul className="mt-0.5 space-y-0.5 pl-9">
                      {item.children.map((child) => (
                        <li key={child}>
                          <a
                            href="#"
                            className="flex min-h-9 items-center rounded-lg px-2.5 text-sm text-gray-500 hover:bg-gray-100 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-400 dark:hover:bg-white/10 dark:hover:text-white"
                          >
                            {child}
                          </a>
                        </li>
                      ))}
                    </ul>
                  </details>
                ) : (
                  <a
                    href="#"
                    /* The only thing that tells a reader which page this is.
                       A violet pill says it to sighted users alone. */
                    aria-current={item.current ? 'page' : undefined}
                    className={`flex min-h-10 items-center gap-2.5 rounded-lg px-2.5 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 ${
                      item.current
                        ? 'bg-violet-50 font-medium text-violet-800 dark:bg-violet-500/15 dark:text-violet-200'
                        : 'text-gray-700 hover:bg-gray-100 dark:text-gray-200 dark:hover:bg-white/10'
                    }`}
                  >
                    <Icon aria-hidden="true" className="size-4 shrink-0" />
                    {!collapsed && item.label}
                    {collapsed && <span className="sr-only">{item.label}</span>}
                  </a>
                )}
              </li>
              )
            })}
          </ul>

          {!collapsed && (
            <>
              <p className="mt-6 px-2.5 pb-1 text-xs font-semibold uppercase tracking-wide text-gray-400">
                Sales channels
              </p>
              <ul className="space-y-0.5">
                {CHANNELS.map((channel) => {
                  const Icon = channel.icon
                  return (
                    <li key={channel.label}>
                      <a
                        href="#"
                        className="flex min-h-10 items-center gap-2.5 rounded-lg px-2.5 text-sm text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-200 dark:hover:bg-white/10"
                      >
                        <Icon aria-hidden="true" className="size-4 shrink-0" />
                        {channel.label}
                      </a>
                    </li>
                  )
                })}
              </ul>
            </>
          )}
        </nav>

        <div className="border-t border-gray-200 p-2 dark:border-white/10">
          <a
            href="#"
            className="flex min-h-10 items-center gap-2.5 rounded-lg px-2.5 text-sm text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-200 dark:hover:bg-white/10"
          >
            <Settings aria-hidden="true" className="size-4 shrink-0" />
            {!collapsed && 'Settings'}
            {collapsed && <span className="sr-only">Settings</span>}
          </a>
        </div>
      </div>

      {/* ── Main ─────────────────────────────────────────────────────────── */}
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex items-center gap-3 border-b border-gray-200 bg-white px-4 py-3 dark:border-white/10 dark:bg-gray-900">
          <div className="relative min-w-0 flex-1">
            <Search
              aria-hidden="true"
              className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-gray-400"
            />
            <label htmlFor="console-search" className="sr-only">
              Search orders, products and customers
            </label>
            <input
              id="console-search"
              type="search"
              placeholder="Search or type a command"
              className="h-10 w-full rounded-lg border border-gray-300 bg-gray-50 pl-9 pr-16 text-sm text-gray-900 placeholder:text-gray-500 focus:border-violet-700 focus:outline-none focus:ring-1 focus:ring-violet-700 dark:border-white/15 dark:bg-white/5 dark:text-white dark:placeholder:text-gray-400"
            />
            <kbd
              aria-hidden="true"
              className="absolute right-2 top-1/2 hidden -translate-y-1/2 rounded border border-gray-300 bg-white px-1.5 py-0.5 text-[10px] text-gray-600 sm:block dark:border-white/15 dark:bg-white/10 dark:text-gray-300"
            >
              &#8984; /
            </kbd>
          </div>

          <button
            type="button"
            className="relative grid size-10 place-items-center rounded-lg text-gray-600 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-300 dark:hover:bg-white/10"
          >
            <Bell aria-hidden="true" className="size-5" />
            <span className="sr-only">Notifications, 3 unread</span>
            <span
              aria-hidden="true"
              className="absolute right-2 top-2 size-2 rounded-full bg-rose-600 ring-2 ring-white dark:ring-gray-900"
            />
          </button>

          <button
            type="button"
            className="grid size-10 shrink-0 place-items-center rounded-full bg-violet-700 text-xs font-semibold text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700"
          >
            AK
            <span className="sr-only">Account menu for Alex Kim</span>
          </button>
        </header>

        <main className="flex-1 overflow-y-auto p-5">
          <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
            <h1 className="text-2xl font-semibold tracking-tight">Orders</h1>
            <div className="flex flex-wrap items-center gap-2">
              <button
                type="button"
                className="inline-flex min-h-10 items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:border-white/15 dark:bg-transparent dark:text-gray-200 dark:hover:bg-white/5"
              >
                <Download aria-hidden="true" className="size-4" />
                Export
              </button>
              <button
                type="button"
                className="inline-flex min-h-10 items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 text-sm font-medium text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:border-white/15 dark:bg-transparent dark:text-gray-200 dark:hover:bg-white/5"
              >
                More actions
                <ChevronDown aria-hidden="true" className="size-4" />
              </button>
              <button
                type="button"
                className="inline-flex min-h-10 items-center gap-2 rounded-lg bg-violet-700 px-3.5 text-sm font-semibold text-white hover:bg-violet-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700"
              >
                <Plus aria-hidden="true" className="size-4" />
                Create order
              </button>
            </div>
          </div>

          <button
            type="button"
            className="mb-5 inline-flex min-h-10 items-center gap-2 rounded-lg border border-gray-300 bg-white px-3 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:border-white/15 dark:bg-transparent dark:text-gray-200 dark:hover:bg-white/5"
          >
            <FileText aria-hidden="true" className="size-4" />
            1 Jan &ndash; 30 Jan 2024
            <ChevronDown aria-hidden="true" className="size-4" />
          </button>

          {/* ── Metrics ────────────────────────────────────────────────── */}
          <dl className="mb-5 grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
            {METRICS.map((metric) => (
              <div
                key={metric.label}
                className="rounded-xl border border-gray-200 bg-white p-4 dark:border-white/10 dark:bg-gray-900"
              >
                <dt className="text-sm text-gray-600 dark:text-gray-300">{metric.label}</dt>
                <dd className="mt-2">
                  <span className="text-3xl font-semibold tracking-tight">{metric.value}</span>
                  <span className="mt-2 flex items-center justify-between gap-2">
                    {/* Written out. A green triangle carries no direction to
                        anyone who is not looking at it. */}
                    <span
                      className={`text-xs font-medium ${
                        metric.direction === 'up'
                          ? 'text-emerald-700 dark:text-emerald-400'
                          : 'text-rose-700 dark:text-rose-400'
                      }`}
                    >
                      {metric.direction === 'up' ? 'Up' : 'Down'} {metric.change}% on last week
                    </span>
                    <svg
                      aria-hidden="true"
                      viewBox="0 0 60 20"
                      className={`h-5 w-16 shrink-0 ${
                        metric.direction === 'up' ? 'text-emerald-500' : 'text-rose-500'
                      }`}
                      preserveAspectRatio="none"
                    >
                      <polyline
                        points={metric.points}
                        fill="none"
                        stroke="currentColor"
                        strokeWidth="1.5"
                      />
                    </svg>
                  </span>
                </dd>
              </div>
            ))}
          </dl>

          {/* ── Table card ─────────────────────────────────────────────── */}
          <div className="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-white/10 dark:bg-gray-900">
            <div className="flex flex-wrap items-center justify-between gap-3 border-b border-gray-200 px-3 py-2 dark:border-white/10">
              {/* A real tablist: arrow keys move between filters, Tab leaves. */}
              <div role="tablist" aria-label="Filter orders" className="flex flex-wrap gap-1">
                {TABS.map((item, index) => {
                  const active = item === tab
                  return (
                    <button
                      key={item}
                      ref={(node) => {
                        tabRefs.current[index] = node
                      }}
                      type="button"
                      role="tab"
                      id={`orders-tab-${item}`}
                      aria-selected={active}
                      aria-controls="orders-panel"
                      tabIndex={active ? 0 : -1}
                      onClick={() => chooseTab(item)}
                      onKeyDown={(event) => onTabKeyDown(event, index)}
                      className={`min-h-9 rounded-lg px-3 text-sm focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 ${
                        active
                          ? 'bg-gray-100 font-medium text-gray-900 dark:bg-white/10 dark:text-white'
                          : 'text-gray-600 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-white/5'
                      }`}
                    >
                      {item}
                    </button>
                  )
                })}
              </div>

              <div className="flex items-center gap-1">
                {[
                  { icon: Search, label: 'Search within these orders' },
                  { icon: SlidersHorizontal, label: 'Filter orders' },
                  { icon: ArrowUpDown, label: 'Change sort order' },
                  { icon: MoreHorizontal, label: 'More table options' },
                ].map(({ icon: Icon, label }) => (
                  <button
                    key={label}
                    type="button"
                    className="grid size-9 place-items-center rounded-lg text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-400 dark:hover:bg-white/10"
                  >
                    <Icon aria-hidden="true" className="size-4" />
                    {/* Named. Four unlabelled icon buttons are four identical
                        entries in a screen reader's control list. */}
                    <span className="sr-only">{label}</span>
                  </button>
                ))}
              </div>
            </div>

            <div
              id="orders-panel"
              role="tabpanel"
              aria-labelledby={`orders-tab-${tab}`}
              tabIndex={-1}
              className="overflow-x-auto"
            >
              <table className="w-full text-left text-sm">
                <caption className="sr-only">
                  {tab} orders, {visible.length} rows, sortable
                </caption>
                <thead className="border-b border-gray-200 text-gray-600 dark:border-white/10 dark:text-gray-300">
                  <tr>
                    <th scope="col" className="w-12 px-3 py-2.5">
                      <input
                        ref={selectAllRef}
                        type="checkbox"
                        checked={allSelected}
                        onChange={toggleAll}
                        aria-label={`Select all ${visible.length} orders`}
                        className="size-4 accent-violet-700"
                      />
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Order
                    </th>
                    <th
                      scope="col"
                      aria-sort={sortKey === 'date' ? (ascending ? 'ascending' : 'descending') : 'none'}
                      className="px-3 py-2.5 font-medium"
                    >
                      <button
                        type="button"
                        onClick={() => toggleSort('date')}
                        className="inline-flex min-h-9 items-center gap-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700"
                      >
                        Date
                        <ArrowUpDown aria-hidden="true" className="size-3.5 text-gray-400" />
                        <span className="sr-only">
                          {sortKey === 'date'
                            ? ascending
                              ? ', sorted ascending, activate to sort descending'
                              : ', sorted descending, activate to sort ascending'
                            : ', activate to sort ascending'}
                        </span>
                      </button>
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Customer
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Payment
                    </th>
                    <th
                      scope="col"
                      aria-sort={sortKey === 'total' ? (ascending ? 'ascending' : 'descending') : 'none'}
                      className="px-3 py-2.5 font-medium"
                    >
                      <button
                        type="button"
                        onClick={() => toggleSort('total')}
                        className="inline-flex min-h-9 items-center gap-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700"
                      >
                        Total
                        <ArrowUpDown aria-hidden="true" className="size-3.5 text-gray-400" />
                        <span className="sr-only">
                          {sortKey === 'total'
                            ? ascending
                              ? ', sorted ascending, activate to sort descending'
                              : ', sorted descending, activate to sort ascending'
                            : ', activate to sort ascending'}
                        </span>
                      </button>
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Delivery
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Items
                    </th>
                    <th scope="col" className="px-3 py-2.5 font-medium">
                      Fulfilment
                    </th>
                    <th scope="col" className="px-3 py-2.5 text-right font-medium">
                      Actions
                    </th>
                  </tr>
                </thead>

                <tbody className="divide-y divide-gray-100 dark:divide-white/5">
                  {visible.map((order) => {
                    const checked = selected.includes(order.id)
                    return (
                      <tr
                        key={order.id}
                        className={checked ? 'bg-violet-50/70 dark:bg-violet-500/10' : ''}
                      >
                        <td className="px-3 py-2.5">
                          <input
                            type="checkbox"
                            checked={checked}
                            onChange={() => toggleRow(order)}
                            aria-label={`Select order ${order.id} for ${order.customer}`}
                            className="size-4 accent-violet-700"
                          />
                        </td>
                        {/* The cell that identifies the row, so a reader can
                            announce it as context for the cells after it. */}
                        <th scope="row" className="px-3 py-2.5 font-medium">
                          <a
                            href="#"
                            className="text-gray-900 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-white"
                          >
                            {order.id}
                          </a>
                        </th>
                        <td className="whitespace-nowrap px-3 py-2.5 text-gray-600 dark:text-gray-300">
                          <time dateTime={order.date}>{dates.format(new Date(order.date))}</time>
                        </td>
                        <td className="whitespace-nowrap px-3 py-2.5 text-gray-900 dark:text-white">
                          {order.customer}
                        </td>
                        <td className="px-3 py-2.5">
                          <span
                            className={`inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-xs font-medium ${PAYMENT_TONE[order.payment]}`}
                          >
                            <span aria-hidden="true" className="size-1.5 rounded-full bg-current" />
                            {order.payment}
                          </span>
                        </td>
                        <td className="px-3 py-2.5 tabular-nums text-gray-900 dark:text-white">
                          {money.format(order.total / 100)}
                        </td>
                        <td className="px-3 py-2.5 text-gray-500 dark:text-gray-400">
                          {order.delivery ?? 'Not set'}
                        </td>
                        <td className="whitespace-nowrap px-3 py-2.5 text-gray-600 dark:text-gray-300">
                          {order.items} {order.items === 1 ? 'item' : 'items'}
                        </td>
                        <td className="px-3 py-2.5">
                          <span
                            className={`inline-flex items-center gap-1.5 rounded-full border px-2 py-0.5 text-xs font-medium ${FULFILMENT_TONE[order.fulfilment]}`}
                          >
                            <span aria-hidden="true" className="size-1.5 rounded-full bg-current" />
                            {order.fulfilment}
                          </span>
                        </td>
                        <td className="px-3 py-2.5">
                          <span className="flex items-center justify-end gap-1">
                            <button
                              type="button"
                              className="grid size-9 place-items-center rounded-lg text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-400 dark:hover:bg-white/10"
                            >
                              <FileText aria-hidden="true" className="size-4" />
                              <span className="sr-only">Print invoice for {order.id}</span>
                            </button>
                            <button
                              type="button"
                              className="grid size-9 place-items-center rounded-lg text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-violet-700 dark:text-gray-400 dark:hover:bg-white/10"
                            >
                              <MessageSquare aria-hidden="true" className="size-4" />
                              <span className="sr-only">Message the customer on {order.id}</span>
                            </button>
                          </span>
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>

              {visible.length === 0 && (
                <p className="px-4 py-16 text-center text-sm text-gray-600 dark:text-gray-300">
                  No orders match the {tab} filter.
                </p>
              )}
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
