'use client'

import { useId } from 'react'
import { ChevronDown, Heart, Menu, Search, ShoppingCart, User } from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

/*
 * The three-tier storefront header: utility strip, search and account row,
 * then departments.
 *
 * Four things here that the obvious version of this header gets wrong, and
 * that are worth reading before you copy it.
 *
 * The search is a <form>. In the source it was a bare <input> next to a
 * <button>, so pressing Enter did nothing and the button submitted nothing —
 * the single most-used control on a storefront, wired to nothing. It also had
 * only a placeholder, which is not a label: it disappears the moment someone
 * types, and it is not reliably announced. There is a real <label>, visually
 * hidden.
 *
 * The icon links have names. Three links containing nothing but an SVG
 * announce as "link, link, link", and the account, wishlist and basket are
 * exactly the controls someone needs to find in a hurry.
 *
 * The counts are announced with their units. A badge reading "3" beside a
 * heart is "3" to a screen reader; here it reads "wishlist, 3 items".
 *
 * The departments menu is Radix rather than hand-rolled. The source toggled a
 * div with `useState` and no `aria-expanded`, no Escape handling, no outside
 * click, and list items that were plain <li> — a menu a keyboard user could
 * open and then never leave or activate. That is what `registryDependencies`
 * is for: this pulls in shadcn's dropdown-menu, which does focus management,
 * type-ahead and Escape properly.
 */

export interface NavItem {
  label: string
  href?: string
}

const UTILITY: NavItem[] = [
  { label: 'About us' },
  { label: 'Order tracking' },
  { label: 'Contact us' },
  { label: 'FAQs' },
]

const DEPARTMENTS: NavItem[] = [
  { label: 'Men' },
  { label: 'Women' },
  { label: 'Kids' },
  { label: 'Electronics' },
  { label: 'Kitchen' },
  { label: 'News & blog' },
  { label: 'Contact' },
]

const LANGUAGES = ['English', 'Español', 'Français']
const CURRENCIES = ['USD', 'EUR', 'GBP']

export default function ThreeTierStorefrontHeader({
  brand = 'Simple UI',
  supportLabel = 'Need help?',
  supportNumber = '+001 123 456 789',
  wishlistCount = 3,
  basketCount = 1,
  utility = UTILITY,
  departments = DEPARTMENTS,
  onSearch,
}: {
  brand?: string
  supportLabel?: string
  supportNumber?: string
  wishlistCount?: number
  basketCount?: number
  utility?: NavItem[]
  departments?: NavItem[]
  onSearch?: (query: string) => void
}) {
  const searchId = useId()

  return (
    <header className="w-full bg-white dark:bg-gray-950">
      <div className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-1 text-sm">
          <nav aria-label="Secondary">
            <ul role="list" className="hidden gap-6 sm:flex">
              {utility.map((item) => (
                <li key={item.label}>
                  <a
                    href={item.href ?? '#'}
                    className="inline-flex min-h-11 items-center text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
                  >
                    {item.label}
                  </a>
                </li>
              ))}
            </ul>
          </nav>

          <div className="ml-auto flex items-center gap-2">
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="gap-1 text-gray-600 dark:text-gray-400">
                  English
                  <span className="sr-only">, change language</span>
                  <ChevronDown aria-hidden="true" className="size-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {LANGUAGES.map((language) => (
                  <DropdownMenuItem key={language}>{language}</DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="gap-1 text-gray-600 dark:text-gray-400">
                  USD
                  <span className="sr-only">, change currency</span>
                  <ChevronDown aria-hidden="true" className="size-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {CURRENCIES.map((currency) => (
                  <DropdownMenuItem key={currency}>{currency}</DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        </div>
      </div>

      <div className="border-b border-gray-200 py-4 dark:border-white/10">
        <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-4 px-6">
          <a
            href="#"
            className="flex items-center gap-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
          >
            <span
              aria-hidden="true"
              className="flex size-10 items-center justify-center rounded-full bg-indigo-600 text-white"
            >
              {brand[0]}
            </span>
            <span className="text-2xl font-bold text-indigo-600 dark:text-indigo-400">{brand}</span>
          </a>

          {/* A real form: Enter submits, and the button is a submit button. */}
          <form
            role="search"
            onSubmit={(event) => {
              event.preventDefault()
              const data = new FormData(event.currentTarget)
              onSearch?.(String(data.get('q') ?? ''))
            }}
            className="mx-4 flex max-w-xl flex-1 items-center"
          >
            <label htmlFor={searchId} className="sr-only">
              Search products
            </label>

            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  type="button"
                  variant="outline"
                  className="rounded-r-none border-r-0 bg-gray-100 text-gray-700 hover:bg-gray-200 dark:bg-white/5 dark:text-gray-300"
                >
                  <span className="hidden sm:inline">All categories</span>
                  <span className="sr-only">Filter search by category</span>
                  <ChevronDown aria-hidden="true" className="size-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="start">
                {departments.slice(0, 5).map((item) => (
                  <DropdownMenuItem key={item.label}>{item.label}</DropdownMenuItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>

            <input
              id={searchId}
              name="q"
              type="search"
              placeholder="I'm shopping for..."
              className="min-h-11 w-full min-w-0 flex-1 border border-gray-300 px-4 text-sm text-gray-900 placeholder:text-gray-400 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:bg-transparent dark:text-white"
            />

            <Button type="submit" className="rounded-l-none">
              <Search aria-hidden="true" className="size-5" />
              <span className="sr-only">Search</span>
            </Button>
          </form>

          <div className="flex items-center gap-5">
            <p className="hidden text-right text-sm sm:block">
              <span className="block font-medium text-gray-900 dark:text-white">
                {supportLabel}
              </span>
              <a
                href={`tel:${supportNumber.replace(/\s/g, '')}`}
                className="text-gray-600 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-400 dark:hover:text-white"
              >
                {supportNumber}
              </a>
            </p>

            <a
              href="#"
              className="hidden size-11 items-center justify-center text-gray-700 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 sm:inline-flex dark:text-gray-300 dark:hover:text-white"
            >
              <User aria-hidden="true" className="size-6" />
              <span className="sr-only">Your account</span>
            </a>

            {/* The count is announced with its unit. A badge reading "3" on its
                own is just "3". */}
            <a
              href="#"
              className="relative inline-flex size-11 items-center justify-center text-gray-700 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:text-white"
            >
              <Heart aria-hidden="true" className="size-6" />
              <span
                aria-hidden="true"
                className="absolute top-0 right-0 flex size-5 items-center justify-center rounded-full bg-indigo-600 text-xs text-white"
              >
                {wishlistCount}
              </span>
              <span className="sr-only">Wishlist, {wishlistCount} items</span>
            </a>

            <a
              href="#"
              className="relative inline-flex size-11 items-center justify-center text-gray-700 hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:text-white"
            >
              <ShoppingCart aria-hidden="true" className="size-6" />
              <span
                aria-hidden="true"
                className="absolute top-0 right-0 flex size-5 items-center justify-center rounded-full bg-indigo-600 text-xs text-white"
              >
                {basketCount}
              </span>
              <span className="sr-only">Basket, {basketCount} items</span>
            </a>
          </div>
        </div>
      </div>

      <div className="border-b border-gray-200 dark:border-white/10">
        <div className="mx-auto flex max-w-7xl items-center px-6">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button className="gap-2 rounded-none md:min-w-56">
                <Menu aria-hidden="true" className="size-5" />
                <span className="hidden md:inline">All categories</span>
                <span className="sr-only md:hidden">All categories</span>
                <ChevronDown aria-hidden="true" className="size-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="start" className="min-w-56">
              {departments.map((item) => (
                <DropdownMenuItem key={item.label} asChild>
                  <a href={item.href ?? '#'}>{item.label}</a>
                </DropdownMenuItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>

          <nav aria-label="Departments" className="hidden md:block">
            <ul role="list" className="flex gap-8 px-6">
              {departments.map((item) => (
                <li key={item.label}>
                  <a
                    href={item.href ?? '#'}
                    className="inline-flex min-h-11 items-center text-sm text-gray-700 hover:text-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400"
                  >
                    {item.label}
                  </a>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>
    </header>
  )
}
