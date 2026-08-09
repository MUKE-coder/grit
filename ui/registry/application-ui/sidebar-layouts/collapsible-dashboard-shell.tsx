'use client'

import { useEffect, useId, useState } from 'react'
import {
  Bell,
  ChevronsLeft,
  CreditCard,
  FileText,
  LayoutDashboard,
  LifeBuoy,
  Menu,
  Package,
  Search,
  Settings,
  ShoppingCart,
  Users,
} from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'

/*
 * An application shell: persistent sidebar, collapsible to icons, with a
 * drawer on small screens.
 *
 * The active item carries aria-current="page". The sources this is based on
 * marked it with a background colour and a text colour and nothing else, so
 * the one piece of information the navigation exists to convey — where you
 * are — was available only to people looking at it. This is the single most
 * common defect in a hand-built app shell.
 *
 * There is a skip link, and it is the first thing in the tab order. A shell
 * with twenty sidebar links means every page begins with twenty links to tab
 * past before reaching the content. The link is positioned off-screen and
 * pulled back on focus rather than hidden with `display: none`, which would
 * make it unfocusable and therefore useless.
 *
 * Links live in lists, and each group is a <nav> with its own accessible name.
 * Screen readers announce "list, six items" and offer landmark navigation, so
 * the grouping becomes a way to move rather than just a visual rhythm. The
 * sources put a bare <h2> above a pile of anchors, which looks the same and
 * navigates like a wall.
 *
 * Collapsing hides the labels with sr-only, not by removing them. Deleting the
 * text turns every item into an icon-only link with no accessible name, which
 * is how a tidy-looking rail becomes unusable. The width is the only thing
 * that actually changes.
 *
 * No magic pixel heights. The sources scrolled the nav inside a hardcoded
 * h-[445px], which is correct at exactly one window height. This is a flex
 * column with the nav taking the remaining space and scrolling inside it.
 */

export interface NavItem {
  label: string
  href: string
  Icon: typeof LayoutDashboard
  badge?: string
}

export interface NavGroup {
  label: string
  items: NavItem[]
}

const GROUPS: NavGroup[] = [
  {
    label: 'Overview',
    items: [
      { label: 'Dashboard', href: '#dashboard', Icon: LayoutDashboard },
      { label: 'Orders', href: '#orders', Icon: ShoppingCart, badge: '6' },
      { label: 'Products', href: '#products', Icon: Package },
      { label: 'Customers', href: '#customers', Icon: Users },
    ],
  },
  {
    label: 'Finance',
    items: [
      { label: 'Invoices', href: '#invoices', Icon: FileText },
      { label: 'Payments', href: '#payments', Icon: CreditCard },
    ],
  },
  {
    label: 'Workspace',
    items: [
      { label: 'Settings', href: '#settings', Icon: Settings },
      { label: 'Support', href: '#support', Icon: LifeBuoy },
    ],
  },
]

function NavList({
  groups,
  current,
  collapsed,
  onNavigate,
}: {
  groups: NavGroup[]
  current: string
  collapsed: boolean
  onNavigate?: (href: string) => void
}) {
  const base = useId()

  return (
    <div className="flex flex-col gap-6">
      {groups.map((group, index) => {
        const labelId = `${base}-group-${index}`
        return (
          /* Each group is a labelled nav, so it turns up in landmark
             navigation by name instead of being one anonymous region. */
          <nav key={group.label} aria-labelledby={labelId}>
            <h2
              id={labelId}
              className={`px-3 text-xs font-semibold tracking-wide text-gray-500 uppercase dark:text-gray-400 ${
                collapsed ? 'sr-only' : ''
              }`}
            >
              {group.label}
            </h2>

            <ul role="list" className="mt-2 space-y-1">
              {group.items.map((item) => {
                const active = item.href === current
                return (
                  <li key={item.href}>
                    <a
                      href={item.href}
                      onClick={() => onNavigate?.(item.href)}
                      /* The whole point of the navigation, and the thing the
                         sources left to colour alone. */
                      aria-current={active ? 'page' : undefined}
                      className={`flex min-h-11 items-center gap-3 rounded-lg px-3 text-sm transition-colors focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 ${
                        active
                          ? 'bg-indigo-50 font-medium text-indigo-900 dark:bg-indigo-500/15 dark:text-indigo-100'
                          : 'text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-white/5'
                      } ${collapsed ? 'justify-center px-0' : ''}`}
                    >
                      <item.Icon aria-hidden="true" className="size-5 shrink-0" />
                      {/* sr-only when collapsed, never removed: an icon-only
                          link with no text has no accessible name. */}
                      <span className={collapsed ? 'sr-only' : 'flex-1'}>{item.label}</span>
                      {item.badge && (
                        <span
                          className={`rounded-full bg-gray-200 px-2 py-0.5 text-xs tabular-nums text-gray-700 dark:bg-white/10 dark:text-gray-200 ${
                            collapsed ? 'sr-only' : ''
                          }`}
                        >
                          {item.badge}
                          <span className="sr-only"> pending</span>
                        </span>
                      )}
                    </a>
                  </li>
                )
              })}
            </ul>
          </nav>
        )
      })}
    </div>
  )
}

export default function CollapsibleDashboardShell({
  brand = 'Northwind',
  groups = GROUPS,
  initialCurrent = '#dashboard',
}: {
  brand?: string
  groups?: NavGroup[]
  initialCurrent?: string
}) {
  const [collapsed, setCollapsed] = useState(false)
  const [current, setCurrent] = useState(initialCurrent)
  const [drawerOpen, setDrawerOpen] = useState(false)
  const contentId = useId()
  const sidebarId = useId()
  const searchId = useId()

  /* The drawer trigger only exists below lg, so an open drawer past that
     breakpoint has no visible way back to it. */
  useEffect(() => {
    const desktop = window.matchMedia('(min-width: 64rem)')
    const sync = () => desktop.matches && setDrawerOpen(false)
    sync()
    desktop.addEventListener('change', sync)
    return () => desktop.removeEventListener('change', sync)
  }, [])

  const activeLabel =
    groups.flatMap((group) => group.items).find((item) => item.href === current)?.label ?? 'Dashboard'

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      {/* First in the tab order, and pulled back on focus rather than hidden
          with display:none, which would make it unfocusable. */}
      <a
        href={`#${contentId}`}
        className="sr-only focus:not-sr-only focus:absolute focus:top-3 focus:left-3 focus:z-50 focus:rounded-lg focus:bg-indigo-700 focus:px-4 focus:py-2.5 focus:text-sm focus:font-medium focus:text-white"
      >
        Skip to content
      </a>

      {/* Desktop sidebar */}
      <div
        id={sidebarId}
        className={`fixed inset-y-0 left-0 z-30 hidden flex-col border-r border-gray-200 bg-white transition-[width] duration-200 lg:flex dark:border-white/10 dark:bg-gray-900 ${
          collapsed ? 'w-16' : 'w-64'
        } motion-reduce:transition-none`}
      >
        <div
          className={`flex h-16 shrink-0 items-center border-b border-gray-200 dark:border-white/10 ${
            collapsed ? 'justify-center px-0' : 'gap-2.5 px-4'
          }`}
        >
          <span
            aria-hidden="true"
            className="flex size-8 shrink-0 items-center justify-center rounded-lg bg-indigo-700 text-sm font-bold text-white"
          >
            {brand.slice(0, 1)}
          </span>
          <span
            className={`text-sm font-semibold text-gray-900 dark:text-white ${
              collapsed ? 'sr-only' : ''
            }`}
          >
            {brand}
          </span>
        </div>

        {/* flex-1 with its own scroll: no hardcoded height that is right at
            exactly one window size. */}
        <div className="flex-1 overflow-y-auto p-3">
          <NavList
            groups={groups}
            current={current}
            collapsed={collapsed}
            onNavigate={setCurrent}
          />
        </div>

        <div className="border-t border-gray-200 p-3 dark:border-white/10">
          <button
            type="button"
            onClick={() => setCollapsed((value) => !value)}
            /* aria-controls points at the sidebar this collapses, not at
               main. Pointing it at the content would claim the button
               controls the page body. aria-expanded describes the sidebar,
               and the label says what the press will do rather than what the
               state currently is. */
            aria-expanded={!collapsed}
            aria-controls={sidebarId}
            className={`flex min-h-11 w-full items-center gap-3 rounded-lg px-3 text-sm text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:bg-white/5 ${
              collapsed ? 'justify-center px-0' : ''
            }`}
          >
            <ChevronsLeft
              aria-hidden="true"
              className={`size-5 shrink-0 transition-transform duration-200 motion-reduce:transition-none ${
                collapsed ? 'rotate-180' : ''
              }`}
            />
            <span className={collapsed ? 'sr-only' : ''}>
              {collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
            </span>
          </button>
        </div>
      </div>

      {/* Top bar + content */}
      <div className={`transition-[padding] duration-200 motion-reduce:transition-none ${collapsed ? 'lg:pl-16' : 'lg:pl-64'}`}>
        <header className="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-gray-200 bg-white px-4 dark:border-white/10 dark:bg-gray-900">
          <Sheet open={drawerOpen} onOpenChange={setDrawerOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="lg:hidden">
                <Menu aria-hidden="true" className="size-5" />
                <span className="sr-only">Open navigation</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex w-72 flex-col p-0">
              <SheetHeader className="h-16 shrink-0 justify-center border-b border-gray-200 px-4 dark:border-white/10">
                <SheetTitle>{brand}</SheetTitle>
              </SheetHeader>
              <div className="flex-1 overflow-y-auto p-3">
                {/* Not collapsed in the drawer: there is no rail to collapse
                    to, and the labels are the whole point on a phone. */}
                <NavList
                  groups={groups}
                  current={current}
                  collapsed={false}
                  onNavigate={(href) => {
                    setCurrent(href)
                    setDrawerOpen(false)
                  }}
                />
              </div>
            </SheetContent>
          </Sheet>

          <form role="search" className="relative hidden flex-1 sm:block" onSubmit={(e) => e.preventDefault()}>
            <label htmlFor={searchId} className="sr-only">
              Search
            </label>
            <Search
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-gray-400"
            />
            <input
              id={searchId}
              type="search"
              placeholder="Search orders, products, customers"
              className="min-h-11 w-full max-w-md rounded-lg border border-gray-300 pr-3 pl-9 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-transparent dark:text-white"
            />
          </form>

          <div className="ml-auto flex items-center gap-1">
            <Button variant="ghost" size="icon" className="relative">
              <Bell aria-hidden="true" className="size-5" />
              <span
                aria-hidden="true"
                className="absolute top-1.5 right-1.5 size-2 rounded-full bg-rose-500"
              />
              <span className="sr-only">Notifications, 3 unread</span>
            </Button>
            <span
              aria-hidden="true"
              className="flex size-9 items-center justify-center rounded-full bg-gray-200 text-xs font-semibold text-gray-700 dark:bg-white/10 dark:text-gray-200"
            >
              AR
            </span>
          </div>
        </header>

        {/* tabIndex so the skip link has somewhere to land that then continues
            the tab order into the content. */}
        <main id={contentId} tabIndex={-1} className="p-6 focus-visible:outline-none">
          <h1 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            {activeLabel}
          </h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
            Replace this with the page body. The shell handles navigation, the
            drawer and the collapse; nothing below is load-bearing.
          </p>

          <div className="mt-6 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {[
              { label: 'Revenue', value: '$48,120' },
              { label: 'Orders', value: '1,284' },
              { label: 'Customers', value: '892' },
              { label: 'Refunds', value: '14' },
            ].map((stat) => (
              <div
                key={stat.label}
                className="rounded-xl border border-gray-200 bg-white p-5 dark:border-white/10 dark:bg-gray-900"
              >
                <p className="text-sm text-gray-500 dark:text-gray-400">{stat.label}</p>
                <p className="mt-1 text-2xl font-semibold tabular-nums text-gray-900 dark:text-white">
                  {stat.value}
                </p>
              </div>
            ))}
          </div>

          <div
            aria-hidden="true"
            className="mt-6 h-64 rounded-xl border border-dashed border-gray-300 dark:border-white/15"
          />
        </main>
      </div>
    </div>
  )
}
