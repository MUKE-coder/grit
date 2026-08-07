'use client'

import { useEffect, useId, useRef, useState } from 'react'
import { Menu, Minus, Plus, Search, ShoppingBag, Trash2, User } from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Sheet,
  SheetClose,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from '@/components/ui/sheet'

/*
 * A storefront header whose cart is a real drawer: quantity steppers, a live
 * subtotal, and an empty state.
 *
 * Money is integer cents. Floats are the wrong type for money and `0.1 + 0.2`
 * is the usual demonstration of why; a cart is exactly where that surfaces,
 * because it multiplies and sums. Formatting happens once, at the edge, with
 * Intl — which also gets the currency symbol and separators right for the
 * locale instead of hardcoding a dollar sign.
 *
 * Three things a cart drawer gets wrong more often than not, fixed here:
 *
 * 1. The steppers are named per item. Six buttons all called "Increase
 *    quantity" are six identical entries in a screen reader's control list.
 *    Each one here says which product it belongs to.
 *
 * 2. Changing a quantity or removing a line updates the subtotal silently. A
 *    sighted user sees the number move; nobody else is told anything. The
 *    `role="status"` region below announces what changed and what the cart now
 *    totals.
 *
 * 3. Removing a line destroys the button that had focus, which drops focus to
 *    the body — and in a dialog that means the next Tab starts from the top of
 *    the document. Focus moves deliberately to the next line's remove button,
 *    or to the heading when the cart empties.
 *
 * The source's mobile nav links each carried a ChevronDown, which says "this
 * expands" to anyone who reads icons, and then navigated instead. They are
 * links, so they look like links.
 */

export interface CartLine {
  id: string
  name: string
  /** Integer cents. See the note above. */
  price: number
  quantity: number
  image: string
  variant?: string
}

const LINES: CartLine[] = [
  {
    id: 'headphones',
    name: 'Noise-cancelling headphones',
    price: 29999,
    quantity: 1,
    variant: 'Midnight',
    image: 'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=200&h=200&fit=crop&q=80',
  },
  {
    id: 'watch',
    name: 'Titanium sport watch',
    price: 42999,
    quantity: 1,
    variant: '45mm',
    image: 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=200&h=200&fit=crop&q=80',
  },
]

const NAV = ['New arrivals', 'Women', 'Men', 'Accessories', 'Collections']

const money = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

function format(cents: number) {
  return money.format(cents / 100)
}

export default function HeaderWithCartDrawer({
  brand = 'Simple UI',
  announcement = 'Free worldwide shipping on orders over $100',
  initialLines = LINES,
  navigation = NAV,
}: {
  brand?: string
  announcement?: string
  initialLines?: CartLine[]
  navigation?: string[]
}) {
  const [lines, setLines] = useState(initialLines)
  const [announcementText, setAnnouncementText] = useState('')
  const searchId = useId()
  const menuSearchId = useId()
  const cartPanel = useRef<HTMLDivElement>(null)
  const cartTitle = useRef<HTMLHeadingElement>(null)
  const removeButtons = useRef(new Map<string, HTMLButtonElement>())

  const subtotal = lines.reduce((sum, line) => sum + line.price * line.quantity, 0)
  const count = lines.reduce((sum, line) => sum + line.quantity, 0)

  /* The status region is emptied between messages. Setting it to the same
     string twice — two clicks on the same stepper — is not a change, so
     assistive tech has nothing to announce the second time. */
  function announce(message: string) {
    setAnnouncementText('')
    requestAnimationFrame(() => setAnnouncementText(message))
  }

  /* Both the new state and the announcement are derived from `lines` up here
     rather than inside the updater. Calling setState from within another
     updater is a rendering side effect, and React runs updaters twice in
     development precisely to make that kind of thing show up. */
  function setQuantity(id: string, delta: number) {
    const line = lines.find((item) => item.id === id)
    if (!line) return
    const quantity = Math.max(1, line.quantity + delta)
    if (quantity === line.quantity) return

    setLines(lines.map((item) => (item.id === id ? { ...item, quantity } : item)))
    announce(
      `${line.name}, quantity ${quantity}. Subtotal ${format(subtotal + (quantity - line.quantity) * line.price)}.`,
    )
  }

  function remove(id: string) {
    const index = lines.findIndex((line) => line.id === id)
    const removed = lines[index]
    const next = lines.filter((line) => line.id !== id)
    setLines(next)
    announce(
      next.length
        ? `${removed.name} removed. Subtotal ${format(subtotal - removed.price * removed.quantity)}.`
        : `${removed.name} removed. Your bag is empty.`,
    )

    /* The button that had focus is about to unmount. Hand focus to the line
       that takes its place, or to the heading if there is none. */
    const successor = next[index] ?? next[index - 1]
    requestAnimationFrame(() => {
      if (successor) removeButtons.current.get(successor.id)?.focus()
      else cartTitle.current?.focus()
    })
  }

  useEffect(() => () => removeButtons.current.clear(), [])

  return (
    <header className="sticky top-0 z-40 w-full border-b border-gray-200 bg-white dark:border-white/10 dark:bg-gray-950">
      <p className="bg-orange-600 px-4 py-2 text-center text-sm font-medium text-white">
        {announcement}
      </p>

      <div className="mx-auto flex max-w-7xl items-center gap-4 px-4 py-3 md:px-6">
        <a
          href="#"
          className="flex shrink-0 items-center gap-1 text-2xl font-bold tracking-tight text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-600 dark:text-white"
        >
          {brand}
          <span
            aria-hidden="true"
            className="flex size-5 items-center justify-center rounded-full bg-orange-600 text-xs text-white"
          >
            ★
          </span>
        </a>

        <form role="search" className="mx-auto hidden max-w-xl flex-1 md:flex">
          <label htmlFor={searchId} className="sr-only">
            Search products, brands and categories
          </label>
          <input
            id={searchId}
            name="q"
            type="search"
            placeholder="Search products, brands and categories"
            className="min-h-11 w-full min-w-0 flex-1 rounded-l-md border border-gray-300 px-3 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:border-white/15 dark:bg-transparent dark:text-white"
          />
          <Button type="submit" className="rounded-l-none bg-orange-600 hover:bg-orange-700">
            Search
          </Button>
        </form>

        <div className="ml-auto flex items-center gap-1 md:ml-0">
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="hidden md:inline-flex">
                <User aria-hidden="true" className="size-5" />
                <span className="sr-only">Your account</span>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuLabel>My account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Orders</DropdownMenuItem>
              <DropdownMenuItem>Wishlist</DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>Sign out</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>

          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="relative">
                <ShoppingBag aria-hidden="true" className="size-5" />
                {count > 0 && (
                  <span
                    aria-hidden="true"
                    className="absolute top-0 right-0 flex size-5 items-center justify-center rounded-full bg-orange-600 text-xs font-medium text-white"
                  >
                    {count}
                  </span>
                )}
                <span className="sr-only">
                  Shopping bag, {count} {count === 1 ? 'item' : 'items'}
                </span>
              </Button>
            </SheetTrigger>

            {/* Focus the panel, not its first focusable child. That child is
                the first line's remove button, so the default behaviour hands
                a keyboard user an armed delete the instant the cart opens. */}
            <SheetContent
              ref={cartPanel}
              side="right"
              onOpenAutoFocus={(event) => {
                event.preventDefault()
                cartPanel.current?.focus()
              }}
              className="flex w-full flex-col p-0 sm:max-w-md"
            >
              <SheetHeader className="border-b border-gray-200 p-6 dark:border-white/10">
                {/* tabIndex so focus can land here when the last line goes. */}
                <SheetTitle ref={cartTitle} tabIndex={-1} className="text-xl">
                  Your shopping bag
                </SheetTitle>
              </SheetHeader>

              {/* Outside the conditional so it survives the cart emptying —
                  a live region mounted at the same moment as its message is
                  frequently missed. */}
              <p role="status" aria-live="polite" className="sr-only">
                {announcementText}
              </p>

              {lines.length === 0 ? (
                <div className="flex flex-1 flex-col items-center justify-center p-6 text-center">
                  <ShoppingBag aria-hidden="true" className="size-16 text-gray-300 dark:text-gray-700" />
                  <p className="mt-6 text-lg font-medium text-gray-900 dark:text-white">
                    Your shopping bag is empty
                  </p>
                  <p className="mt-2 max-w-xs text-sm text-gray-600 dark:text-gray-400">
                    Nothing in here yet. Have a look at what is new.
                  </p>
                  <SheetClose asChild>
                    <Button className="mt-8 bg-orange-600 hover:bg-orange-700">
                      Continue shopping
                    </Button>
                  </SheetClose>
                </div>
              ) : (
                <>
                  <ul role="list" className="flex-1 space-y-6 overflow-y-auto p-6">
                    {lines.map((line) => (
                      <li key={line.id} className="grid grid-cols-[80px_1fr] gap-4">
                        <img
                          src={line.image}
                          alt=""
                          className="size-20 rounded-md bg-gray-50 object-cover dark:bg-white/5"
                        />
                        <div>
                          <div className="flex items-start justify-between gap-2">
                            <div>
                              <h3 className="font-medium leading-tight text-gray-900 dark:text-white">
                                {line.name}
                              </h3>
                              {line.variant && (
                                <p className="mt-0.5 text-sm text-gray-500 dark:text-gray-400">
                                  {line.variant}
                                </p>
                              )}
                            </div>
                            <button
                              type="button"
                              ref={(node) => {
                                if (node) removeButtons.current.set(line.id, node)
                                else removeButtons.current.delete(line.id)
                              }}
                              onClick={() => remove(line.id)}
                              className="-m-2 inline-flex size-11 shrink-0 items-center justify-center rounded-full text-gray-500 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-orange-600 dark:text-gray-400 dark:hover:bg-white/10"
                            >
                              {/* A bin, not an X. The panel's own close is an
                                  X in the corner directly above this one, and
                                  two adjacent X buttons where one dismisses
                                  and the other destroys is a trap. */}
                              <Trash2 aria-hidden="true" className="size-4" />
                              <span className="sr-only">Remove {line.name}</span>
                            </button>
                          </div>

                          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                            {format(line.price)}
                          </p>

                          <div className="mt-3 flex items-center gap-3">
                            <div className="flex items-center rounded-full border border-gray-300 dark:border-white/15">
                              <button
                                type="button"
                                onClick={() => setQuantity(line.id, -1)}
                                disabled={line.quantity === 1}
                                className="inline-flex size-11 items-center justify-center rounded-full text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 disabled:opacity-40 dark:text-gray-300 dark:hover:bg-white/10"
                              >
                                <Minus aria-hidden="true" className="size-3" />
                                <span className="sr-only">Decrease quantity of {line.name}</span>
                              </button>
                              {/* The number is decoration for assistive tech —
                                  the stepper announces the new value, and this
                                  would otherwise be read as a bare digit. */}
                              <span
                                aria-hidden="true"
                                className="w-8 text-center text-sm text-gray-900 dark:text-white"
                              >
                                {line.quantity}
                              </span>
                              <button
                                type="button"
                                onClick={() => setQuantity(line.id, 1)}
                                className="inline-flex size-11 items-center justify-center rounded-full text-gray-700 hover:bg-gray-100 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:text-gray-300 dark:hover:bg-white/10"
                              >
                                <Plus aria-hidden="true" className="size-3" />
                                <span className="sr-only">Increase quantity of {line.name}</span>
                              </button>
                            </div>
                            <p className="ml-auto font-medium text-gray-900 dark:text-white">
                              <span className="sr-only">Line total </span>
                              {format(line.price * line.quantity)}
                            </p>
                          </div>
                        </div>
                      </li>
                    ))}
                  </ul>

                  <div className="space-y-4 border-t border-gray-200 p-6 dark:border-white/10">
                    <dl className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <dt className="text-gray-500 dark:text-gray-400">Subtotal</dt>
                        <dd className="font-medium text-gray-900 dark:text-white">
                          {format(subtotal)}
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500 dark:text-gray-400">Shipping</dt>
                        <dd className="font-medium text-gray-900 dark:text-white">
                          Calculated at checkout
                        </dd>
                      </div>
                      <div className="flex justify-between border-t border-gray-200 pt-3 text-base dark:border-white/10">
                        <dt className="font-semibold text-gray-900 dark:text-white">Total</dt>
                        <dd className="font-semibold text-gray-900 dark:text-white">
                          {format(subtotal)}
                        </dd>
                      </div>
                    </dl>

                    <Button className="min-h-12 w-full bg-orange-600 hover:bg-orange-700">
                      Proceed to checkout
                    </Button>
                    <SheetClose asChild>
                      <Button variant="outline" className="w-full">
                        Continue shopping
                      </Button>
                    </SheetClose>
                    <p className="text-center text-xs text-gray-500 dark:text-gray-400">
                      Shipping and taxes calculated at checkout.
                    </p>
                  </div>
                </>
              )}
            </SheetContent>
          </Sheet>

          <Sheet>
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="md:hidden">
                <Menu aria-hidden="true" className="size-5" />
                <span className="sr-only">Open menu</span>
              </Button>
            </SheetTrigger>
            <SheetContent side="left" className="flex flex-col overflow-y-auto p-0">
              <SheetHeader className="border-b border-gray-200 p-4 dark:border-white/10">
                <SheetTitle>Menu</SheetTitle>
              </SheetHeader>

              <form role="search" className="relative p-4">
                <label htmlFor={menuSearchId} className="sr-only">
                  Search products
                </label>
                <Search
                  aria-hidden="true"
                  className="pointer-events-none absolute top-1/2 left-7 size-4 -translate-y-1/2 text-gray-400"
                />
                <input
                  id={menuSearchId}
                  name="q"
                  type="search"
                  placeholder="Search"
                  className="min-h-11 w-full rounded-md border border-gray-300 pr-3 pl-9 text-sm text-gray-900 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:border-white/15 dark:bg-transparent dark:text-white"
                />
              </form>

              <nav aria-label="Main" className="px-2">
                <ul role="list">
                  {navigation.map((item) => (
                    <li key={item}>
                      <SheetClose asChild>
                        <a
                          href="#"
                          className="flex min-h-11 items-center rounded-md px-4 font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:text-white dark:hover:bg-white/5"
                        >
                          {item}
                        </a>
                      </SheetClose>
                    </li>
                  ))}
                </ul>
              </nav>

              <div className="mt-auto border-t border-gray-200 px-2 py-2 dark:border-white/10">
                <ul role="list">
                  {['Account', 'Wishlist', 'Order tracking', 'Help & contact'].map((item) => (
                    <li key={item}>
                      <SheetClose asChild>
                        <a
                          href="#"
                          className="flex min-h-11 items-center rounded-md px-4 text-sm text-gray-700 hover:bg-gray-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-orange-600 dark:text-gray-300 dark:hover:bg-white/5"
                        >
                          {item}
                        </a>
                      </SheetClose>
                    </li>
                  ))}
                </ul>
              </div>
            </SheetContent>
          </Sheet>
        </div>
      </div>
    </header>
  )
}
