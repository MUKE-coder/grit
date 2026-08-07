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
 * A storefront header that shrinks on scroll, with a mobile drawer.
 *
 * The drawer is the reason this block declares a primitive. The source
 * rendered a fixed panel and a click-catching overlay when a state flag went
 * true, which *looks* modal and is not: focus stayed on the trigger, tabbing
 * walked straight into the page behind, Escape did nothing, the content behind
 * was still announced, and closing left focus nowhere. shadcn's sheet is Radix
 * Dialog underneath, which moves focus in, traps it, closes on Escape, marks
 * the rest of the page `aria-hidden`, and puts focus back on the trigger
 * afterwards.
 *
 * Two things the primitive cannot decide for you, both handled below. Opening
 * focuses the panel rather than the search field, because the first focusable
 * thing in a mobile drawer is usually an input and focusing it throws up the
 * on-screen keyboard over the menu the person just asked to see. And the
 * drawer closes when the viewport crosses to desktop, because its trigger is
 * `md:hidden` — resize with it open and you have a panel whose only visible
 * exit is the close button.
 *
 * The header measures its own height and offsets the page by that, rather than
 * the source's hardcoded `h-[80px] md:h-[130px]` spacer. A fixed header with a
 * guessed spacer is correct only while the two numbers agree, and they stop
 * agreeing the first time someone adds a promo bar.
 *
 * Both search fields are real forms with real labels. A placeholder is not a
 * label: it disappears the moment someone types, and it is not reliably
 * announced.
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

  /* The trigger is md:hidden, so past that breakpoint an open drawer has no
     visible way back to it. Match the breakpoint here, not a resize handler
     firing on every pixel. */
  useEffect(() => {
    const desktop = window.matchMedia('(min-width: 48rem)')
    const sync = () => desktop.matches && setOpen(false)
    sync()
    desktop.addEventListener('change', sync)
    return () => desktop.removeEventListener('change', sync)
  }, [])

  /* Measured, not guessed. A hardcoded spacer is right until the header
     changes height, and then it is silently wrong. */
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

            {/* Radix handles focus in, focus trapped, Escape, hidden
                background and focus restored. None of that is free in a
                hand-rolled drawer, and all of it is missing from most of them.

                The override redirects the initial focus to the panel itself.
                Left alone Radix focuses the first focusable descendant, which
                here is the search box — and a keyboard sliding up over a menu
                is not what "open the menu" asked for. preventDefault alone
                would leave focus on the now-hidden trigger, so the focus() call
                is not optional. */}
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
                      {/* SheetClose closes the drawer AND restores focus, which
                          a plain onClick handler does not. */}
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

      {/* Offset by the header's real height, tracked as it changes. */}
      <div aria-hidden="true" style={{ height }} />
    </>
  )
}
