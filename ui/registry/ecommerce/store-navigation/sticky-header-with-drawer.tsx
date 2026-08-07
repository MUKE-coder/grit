'use client'

import { useEffect, useId, useRef, useState } from 'react'
import { Heart, Menu, Package, Search, ShoppingCart, User } from 'lucide-react'

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
 * Storefront header that shrinks on scroll, with a mobile drawer.
 *
 * The drawer uses the sheet primitive. The source rendered a fixed panel plus
 * a click-catching overlay when a state flag went true, so focus stayed on the
 * trigger, tab went straight into the page behind, and Escape did nothing.
 *
 * Two things the primitive can't decide. Opening focuses the panel, not the
 * search field, so the on-screen keyboard doesn't cover the menu. And the
 * drawer closes when the viewport crosses to desktop, where its trigger is
 * md:hidden.
 *
 * The page offset is measured from the header instead of the source's
 * hardcoded h-[80px] md:h-[130px] spacer, which is only correct until someone
 * adds a promo bar.
 */

export interface NavLink {
  label: string
  href?: string
}

const LINKS: NavLink[] = [
  { label: 'New in' },
  { label: 'Women' },
  { label: 'Men' },
  { label: 'Accessories' },
  { label: 'Sale' },
]

export default function StickyHeaderWithDrawer({
  brand = 'Simple UI',
  links = LINKS,
  wishlistCount = 2,
  basketCount = 3,
  onSearch,
}: {
  brand?: string
  links?: NavLink[]
  wishlistCount?: number
  basketCount?: number
  onSearch?: (query: string) => void
}) {
  const [scrolled, setScrolled] = useState(false)
  const [height, setHeight] = useState(0)
  const [open, setOpen] = useState(false)
  const panel = useRef<HTMLDivElement>(null)
  const searchId = useId()
  const drawerSearchId = useId()

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 10)
    onScroll()
    window.addEventListener('scroll', onScroll, { passive: true })
    return () => window.removeEventListener('scroll', onScroll)
  }, [])

  /* Trigger is md:hidden, so past that breakpoint an open drawer has no
     visible way back to it. matchMedia rather than a resize handler. */
  useEffect(() => {
    const desktop = window.matchMedia('(min-width: 48rem)')
    const sync = () => desktop.matches && setOpen(false)
    sync()
    desktop.addEventListener('change', sync)
    return () => desktop.removeEventListener('change', sync)
  }, [])

  /* Measured. A hardcoded spacer goes silently wrong the moment the header
     changes height. */
  useEffect(() => {
    const el = document.getElementById('storefront-header')
    if (!el) return
    const observer = new ResizeObserver(() => setHeight(el.offsetHeight))
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  function handleSearch(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const data = new FormData(event.currentTarget)
    onSearch?.(String(data.get('q') ?? ''))
  }

  return (
    <>
      <header
        id="storefront-header"
        className={`fixed inset-x-0 top-0 z-40 bg-white transition-shadow dark:bg-gray-950 ${
          scrolled ? 'shadow-md' : 'border-b border-gray-200 dark:border-white/10'
        }`}
      >
        <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-3 sm:px-6">
          <Sheet open={open} onOpenChange={setOpen}>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden">
                <Menu aria-hidden="true" className="size-5" />
                <span className="sr-only">Open menu</span>
              </Button>
            </SheetTrigger>

            {/* Radix focuses the first focusable descendant by default, which
                here is the search box. preventDefault alone would leave focus
                on the now-hidden trigger, so the focus() call is required. */}
            <SheetContent
              ref={panel}
              side="left"
              onOpenAutoFocus={(event) => {
                event.preventDefault()
                panel.current?.focus()
              }}
              className="flex w-72 flex-col overflow-y-auto p-0"
            >
              <SheetHeader className="border-b border-gray-200 p-4 dark:border-white/10">
                <SheetTitle>{brand}</SheetTitle>
              </SheetHeader>

              <form
                role="search"
                onSubmit={handleSearch}
                className="flex border-b border-gray-200 p-4 dark:border-white/10"
              >
                <label htmlFor={drawerSearchId} className="sr-only">
                  Search products
                </label>
                <input
                  id={drawerSearchId}
                  name="q"
                  type="search"
                  placeholder="Search products..."
                  className="min-h-11 w-full min-w-0 flex-1 rounded-l-md border border-gray-300 px-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-transparent dark:text-white"
                />
                <Button type="submit" size="icon" className="rounded-l-none">
                  <Search aria-hidden="true" className="size-4" />
                  <span className="sr-only">Search</span>
                </Button>
              </form>

              <nav aria-label="Main" className="p-2">
                <ul role="list">
                  {links.map((link) => (
                    <li key={link.label}>
                      {/* SheetClose also restores focus. A plain onClick
                          handler doesn't. */}
                      <SheetClose asChild>
                        <a
                          href={link.href ?? '#'}
                          className="flex min-h-11 items-center rounded-md px-4 font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:text-white dark:hover:bg-white/5"
                        >
                          {link.label}
                        </a>
                      </SheetClose>
                    </li>
                  ))}
                </ul>
              </nav>

              <div className="mt-auto border-t border-gray-200 p-2 dark:border-white/10">
                <ul role="list">
                  {[
                    { label: 'My account', Icon: User },
                    { label: 'My orders', Icon: Package },
                  ].map(({ label, Icon }) => (
                    <li key={label}>
                      <SheetClose asChild>
                        <a
                          href="#"
                          className="flex min-h-11 items-center gap-3 rounded-md px-4 text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:text-white dark:hover:bg-white/5"
                        >
                          <Icon aria-hidden="true" className="size-5" />
                          {label}
                        </a>
                      </SheetClose>
                    </li>
                  ))}
                </ul>
              </div>
            </SheetContent>
          </Sheet>

          <a
            href="#"
            className="text-xl font-bold text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-white"
          >
            {brand}
          </a>

          <nav aria-label="Main" className="hidden md:block">
            <ul role="list" className="flex gap-6">
              {links.map((link) => (
                <li key={link.label}>
                  <a
                    href={link.href ?? '#'}
                    className="inline-flex min-h-11 items-center text-sm font-medium text-gray-700 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:text-white"
                  >
                    {link.label}
                  </a>
                </li>
              ))}
            </ul>
          </nav>

          <form
            role="search"
            onSubmit={handleSearch}
            className="ml-auto hidden max-w-xs flex-1 lg:flex"
          >
            <label htmlFor={searchId} className="sr-only">
              Search products
            </label>
            <input
              id={searchId}
              name="q"
              type="search"
              placeholder="Search products..."
              className="min-h-11 w-full min-w-0 flex-1 rounded-l-md border border-gray-300 px-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-transparent dark:text-white"
            />
            <Button type="submit" size="icon" className="rounded-l-none">
              <Search aria-hidden="true" className="size-4" />
              <span className="sr-only">Search</span>
            </Button>
          </form>

          <div className="ml-auto flex items-center gap-1 lg:ml-0">
            <a
              href="#"
              className="hidden size-11 items-center justify-center rounded-md text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:inline-flex dark:text-gray-300 dark:hover:bg-white/10"
            >
              <User aria-hidden="true" className="size-5" />
              <span className="sr-only">Your account</span>
            </a>

            <a
              href="#"
              className="relative inline-flex size-11 items-center justify-center rounded-md text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:bg-white/10"
            >
              <Heart aria-hidden="true" className="size-5" />
              <span
                aria-hidden="true"
                className="absolute top-1 right-1 flex size-4 items-center justify-center rounded-full bg-gray-900 text-[10px] text-white dark:bg-white dark:text-gray-900"
              >
                {wishlistCount}
              </span>
              <span className="sr-only">Wishlist, {wishlistCount} items</span>
            </a>

            <a
              href="#"
              className="relative inline-flex size-11 items-center justify-center rounded-md text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:bg-white/10"
            >
              <ShoppingCart aria-hidden="true" className="size-5" />
              <span
                aria-hidden="true"
                className="absolute top-1 right-1 flex size-4 items-center justify-center rounded-full bg-gray-900 text-[10px] text-white dark:bg-white dark:text-gray-900"
              >
                {basketCount}
              </span>
              <span className="sr-only">Basket, {basketCount} items</span>
            </a>
          </div>
        </div>
      </header>

      {/* Offset by the header's real height. */}
      <div aria-hidden="true" style={{ height }} />
    </>
  )
}
