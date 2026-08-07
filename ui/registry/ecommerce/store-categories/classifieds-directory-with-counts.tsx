'use client'

import { useId, useRef, useState } from 'react'
import {
  Armchair,
  Bike,
  Building2,
  Bus,
  Car,
  ChevronRight,
  Cpu,
  Dog,
  Home,
  Laptop,
  Monitor,
  Shirt,
  Ship,
  Smartphone,
  Sofa,
  Sparkles,
  Truck,
  Tv,
  Warehouse,
  Watch,
  Wrench,
} from 'lucide-react'

/*
 * Classifieds directory: top-level categories on the left, that category's
 * subcategories and listing counts on the right.
 *
 * The source opened its submenus on hover and nothing else, so the
 * subcategories were unreachable by keyboard and unreachable on a phone. It
 * also threaded an isExpanded flag and an onToggle handler down to every row,
 * where the toggle set state that nothing read, so the rows looked interactive
 * and were not. The rows were divs with onClick, no role and no tabindex.
 *
 * This is a tablist instead. Arrow keys move between categories, Home and End
 * jump to the ends, and the panel is reachable with Tab. Selection follows
 * focus, which APG allows when switching panels is cheap.
 *
 * aria-orientation is vertical because that is the wide layout, but both axes
 * of arrow keys are handled: below lg the same tablist renders as a scrolling
 * row, and someone pressing ArrowRight there should not find nothing happens.
 *
 * Counts are exact and grouped by Intl rather than run through a hand-written
 * "if (n > 1000)" helper. Compact notation was the first attempt and it prints
 * 1,006 as "1K" directly under a literal "398", which reads as a rounding bug.
 */

export interface Subcategory {
  name: string
  href: string
  count: number
  Icon: typeof Car
}

export interface Category {
  id: string
  name: string
  href: string
  count: number
  Icon: typeof Car
  subcategories: Subcategory[]
}

const CATEGORIES: Category[] = [
  {
    id: 'vehicles',
    name: 'Vehicles',
    href: '#',
    count: 56143,
    Icon: Car,
    subcategories: [
      { name: 'Cars', href: '#', count: 11876, Icon: Car },
      { name: 'Buses and microbuses', href: '#', count: 398, Icon: Bus },
      { name: 'Motorcycles and scooters', href: '#', count: 7704, Icon: Bike },
      { name: 'Trucks and trailers', href: '#', count: 1006, Icon: Truck },
      { name: 'Parts and accessories', href: '#', count: 34867, Icon: Wrench },
      { name: 'Boats and watercraft', href: '#', count: 67, Icon: Ship },
    ],
  },
  {
    id: 'property',
    name: 'Property',
    href: '#',
    count: 21430,
    Icon: Home,
    subcategories: [
      { name: 'Houses for sale', href: '#', count: 8450, Icon: Home },
      { name: 'Apartments for rent', href: '#', count: 5200, Icon: Building2 },
      { name: 'Land and plots', href: '#', count: 4120, Icon: Warehouse },
      { name: 'Commercial property', href: '#', count: 2380, Icon: Building2 },
      { name: 'Short lets', href: '#', count: 1280, Icon: Home },
    ],
  },
  {
    id: 'electronics',
    name: 'Mobile and electronics',
    href: '#',
    count: 48902,
    Icon: Smartphone,
    subcategories: [
      { name: 'Mobile phones', href: '#', count: 22140, Icon: Smartphone },
      { name: 'Laptops', href: '#', count: 9860, Icon: Laptop },
      { name: 'Televisions', href: '#', count: 6410, Icon: Tv },
      { name: 'Monitors', href: '#', count: 3990, Icon: Monitor },
      { name: 'Components', href: '#', count: 4310, Icon: Cpu },
      { name: 'Watches and wearables', href: '#', count: 2192, Icon: Watch },
    ],
  },
  {
    id: 'home',
    name: 'Home and garden',
    href: '#',
    count: 33517,
    Icon: Sofa,
    subcategories: [
      { name: 'Furniture', href: '#', count: 14200, Icon: Armchair },
      { name: 'Home appliances', href: '#', count: 9870, Icon: Warehouse },
      { name: 'Garden and outdoor', href: '#', count: 5240, Icon: Home },
      { name: 'Kitchen and dining', href: '#', count: 4207, Icon: Sofa },
    ],
  },
  {
    id: 'fashion',
    name: 'Fashion',
    href: '#',
    count: 27884,
    Icon: Shirt,
    subcategories: [
      { name: 'Clothing', href: '#', count: 15400, Icon: Shirt },
      { name: 'Shoes', href: '#', count: 7120, Icon: Shirt },
      { name: 'Bags and luggage', href: '#', count: 3140, Icon: Shirt },
      { name: 'Jewellery and watches', href: '#', count: 2224, Icon: Watch },
    ],
  },
  {
    id: 'pets',
    name: 'Animals and pets',
    href: '#',
    count: 6218,
    Icon: Dog,
    subcategories: [
      { name: 'Dogs and puppies', href: '#', count: 2840, Icon: Dog },
      { name: 'Cats and kittens', href: '#', count: 1620, Icon: Dog },
      { name: 'Birds', href: '#', count: 980, Icon: Dog },
      { name: 'Pet accessories', href: '#', count: 778, Icon: Sparkles },
    ],
  },
]

const counts = new Intl.NumberFormat('en-US')

export default function ClassifiedsDirectoryWithCounts({
  title = 'Browse categories',
  categories = CATEGORIES,
}: {
  title?: string
  categories?: Category[]
}) {
  const [active, setActive] = useState(0)
  const baseId = useId()
  const tabs = useRef<(HTMLButtonElement | null)[]>([])

  const tabId = (index: number) => `${baseId}-tab-${index}`
  const panelId = (index: number) => `${baseId}-panel-${index}`

  /* Selection follows focus, so moving is one keypress rather than move-then-
     confirm. Both axes because the tablist is a column above lg and a row
     below it. */
  function onKeyDown(event: React.KeyboardEvent) {
    const last = categories.length - 1
    let next: number | null = null

    if (event.key === 'ArrowDown' || event.key === 'ArrowRight') next = active === last ? 0 : active + 1
    else if (event.key === 'ArrowUp' || event.key === 'ArrowLeft') next = active === 0 ? last : active - 1
    else if (event.key === 'Home') next = 0
    else if (event.key === 'End') next = last

    if (next === null) return
    event.preventDefault()
    setActive(next)
    tabs.current[next]?.focus()
  }

  const current = categories[active]

  return (
    <section className="bg-gray-50 py-12 dark:bg-gray-950">
      <div className="mx-auto max-w-5xl px-4">
        <h2 className="mb-8 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
          {title}
        </h2>

        <div className="grid grid-cols-1 gap-4 lg:grid-cols-[20rem_1fr]">
          <div
            role="tablist"
            aria-orientation="vertical"
            aria-label={title}
            onKeyDown={onKeyDown}
            className="flex gap-2 overflow-x-auto rounded-lg border border-gray-200 bg-white p-2 lg:flex-col lg:gap-0 lg:overflow-x-visible lg:divide-y lg:divide-gray-100 lg:p-0 dark:border-white/10 dark:bg-gray-900 lg:dark:divide-white/5"
          >
            {categories.map((category, index) => (
              <button
                key={category.id}
                ref={(node) => {
                  tabs.current[index] = node
                }}
                type="button"
                role="tab"
                id={tabId(index)}
                aria-selected={index === active}
                aria-controls={panelId(index)}
                /* Roving tabindex: one stop for the whole list, then arrows
                   inside it. Nine tabbable rows would be nine Tab presses to
                   get past a navigation aid. */
                tabIndex={index === active ? 0 : -1}
                onClick={() => setActive(index)}
                className={`flex shrink-0 items-center gap-3 rounded-lg p-4 text-left whitespace-nowrap focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 lg:w-full lg:rounded-none lg:whitespace-normal ${
                  index === active
                    ? 'bg-indigo-50 text-indigo-900 dark:bg-indigo-500/10 dark:text-indigo-200'
                    : 'text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-white/5'
                }`}
              >
                <category.Icon aria-hidden="true" className="size-6 shrink-0" />
                <span className="flex-1">
                  <span className="block font-medium">{category.name}</span>
                  <span className="block text-sm text-gray-500 dark:text-gray-400">
                    {counts.format(category.count)} ads
                  </span>
                </span>
                <ChevronRight
                  aria-hidden="true"
                  className="hidden size-4 shrink-0 text-gray-400 lg:block"
                />
              </button>
            ))}
          </div>

          {/* tabIndex so the panel itself is reachable with Tab. Without it a
              keyboard user leaving the tablist lands on the first subcategory
              link with no announcement of what changed. */}
          <div
            role="tabpanel"
            id={panelId(active)}
            aria-labelledby={tabId(active)}
            tabIndex={0}
            className="rounded-lg border border-gray-200 bg-white p-4 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/10 dark:bg-gray-900"
          >
            <div className="mb-3 flex items-baseline justify-between gap-4 px-2">
              <h3 className="font-semibold text-gray-900 dark:text-white">{current.name}</h3>
              <a
                href={current.href}
                className="text-sm font-medium text-indigo-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
              >
                All {counts.format(current.count)} ads
              </a>
            </div>

            <ul role="list" className="grid grid-cols-1 gap-1 sm:grid-cols-2">
              {current.subcategories.map((sub) => (
                <li key={sub.name}>
                  <a
                    href={sub.href}
                    className="flex items-center gap-3 rounded-lg p-3 hover:bg-gray-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 dark:hover:bg-white/5"
                  >
                    <sub.Icon
                      aria-hidden="true"
                      className="size-5 shrink-0 text-gray-500 dark:text-gray-400"
                    />
                    <span>
                      <span className="block text-gray-900 dark:text-white">{sub.name}</span>
                      <span className="block text-sm text-gray-500 dark:text-gray-400">
                        {counts.format(sub.count)} ads
                      </span>
                    </span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  )
}
