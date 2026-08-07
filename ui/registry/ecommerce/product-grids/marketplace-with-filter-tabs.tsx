'use client'

import { useId, useState } from 'react'
import { MapPin, SearchX } from 'lucide-react'

/*
 * A marketplace listing filtered by tabs: price first, then condition, then
 * where the thing is.
 *
 * The ordering is the design. On a classifieds listing the price decides
 * whether the rest of the card gets read at all, and the location decides
 * whether the sale can happen — so those two frame everything else rather than
 * being footnotes under a product name.
 *
 * Two things here are worth copying, and both are things the obvious version
 * gets wrong.
 *
 * The tabs are a real tablist: `role="tab"`, `aria-selected`, `aria-controls`,
 * and arrow keys that move selection and focus together. Two styled buttons
 * with onClick look identical and announce as two unrelated buttons, with no
 * indication that one of them is currently active.
 *
 * The cards are links, not divs with an onClick. A div that responds to clicks
 * cannot be reached by keyboard at all — no tab stop, no Enter, nothing — and
 * a screen reader has no reason to mention it. It is the single most common
 * way a card grid becomes unusable, and it looks perfect while doing it.
 *
 * The result is announced through `aria-live` when the filter changes.
 * Otherwise a keyboard user presses a tab, the grid silently swaps underneath,
 * and nothing tells them how many results they now have.
 */

export interface Listing {
  id: string
  title: string
  /** Pre-formatted, because a marketplace is rarely single-currency. */
  price: string
  description: string
  location: string
  condition: string
  category: string
  image: string
  href?: string
}

export interface Tab {
  id: string
  label: string
}

const TABS: Tab[] = [
  { id: 'retail', label: 'Retail' },
  { id: 'warehouse', label: 'Warehouse' },
]

const LISTINGS: Listing[] = [
  {
    id: '1',
    title: 'Samsung Galaxy S21+ 5G 128GB, silver',
    price: 'USh 450,000',
    description: 'Boxed with the charger. Small mark on the frame, screen is clean.',
    location: 'Kampala, Central Division',
    condition: 'Used',
    category: 'retail',
    image: 'https://images.unsplash.com/photo-1610945265064-0e34e5519bbf?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '2',
    title: 'Apple iPhone 11 Pro 256GB, silver',
    price: 'USh 900,000',
    description: 'Battery health 89%. Triple camera all working.',
    location: 'Kampala, Central Division',
    condition: 'Used',
    category: 'retail',
    image: 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '3',
    title: 'Xiaomi Mi 11 Lite 128GB, yellow or blue',
    price: 'USh 520,000',
    description: 'Two units available. Both sealed, both with the local warranty.',
    location: 'Kampala, Nakawa',
    condition: 'New',
    category: 'retail',
    image: 'https://images.unsplash.com/photo-1621330396173-e41b1cafd17f?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '4',
    title: 'Apple iPhone 12 128GB, black',
    price: 'USh 750,000',
    description: 'Face ID working, no scratches. Refurbished and tested.',
    location: 'Kampala, Central Division',
    condition: 'Refurbished',
    category: 'retail',
    image: 'https://images.unsplash.com/photo-1616348436168-de43ad0db179?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '5',
    title: 'Apple iPhone 11 64GB, blue',
    price: 'USh 480,000',
    description: 'Bought abroad, used for eight months. Charger included.',
    location: 'Kampala, Makindye',
    condition: 'Used',
    category: 'warehouse',
    image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '6',
    title: 'Apple iPhone 12 with AirPods Pro',
    price: 'USh 950,000',
    description: 'Bundle. Phone and buds both boxed, both in good order.',
    location: 'Wakiso, Kira',
    condition: 'Used',
    category: 'warehouse',
    image: 'https://images.unsplash.com/photo-1592899677977-9c10ca588bbd?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '7',
    title: 'Apple iPhone 11 64GB, yellow',
    price: 'USh 180,000',
    description: 'Selling for parts. Screen intact, board is not.',
    location: 'Kampala, Makindye',
    condition: 'For parts',
    category: 'warehouse',
    image: 'https://images.unsplash.com/photo-1574755393849-623942496936?w=600&h=450&fit=crop&q=80',
  },
  {
    id: '8',
    title: 'Apple iPhone X 64GB, silver',
    price: 'USh 320,000',
    description: 'Battery replaced last month. Slight burn-in at the top of the screen.',
    location: 'Kampala, Central Division',
    condition: 'Refurbished',
    category: 'warehouse',
    image: 'https://images.unsplash.com/photo-1580910051074-3eb694886505?w=600&h=450&fit=crop&q=80',
  },
]

/* Any condition not listed here falls back to neutral rather than being
   mislabelled. A ternary would paint "For parts" the same as "Refurbished". */
const CONDITION_TONE: Record<string, string> = {
  New: 'bg-indigo-50 text-indigo-700 dark:bg-indigo-500/10 dark:text-indigo-400',
  Used: 'bg-emerald-50 text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-400',
  Refurbished: 'bg-sky-50 text-sky-700 dark:bg-sky-500/10 dark:text-sky-400',
  'For parts': 'bg-amber-50 text-amber-700 dark:bg-amber-500/10 dark:text-amber-400',
}

export default function MarketplaceWithFilterTabs({
  title = 'Mobile phone marketplace',
  subtitle = 'Quality handsets from sellers across the city, listed with what is actually wrong with them.',
  tabs = TABS,
  listings = LISTINGS,
}: {
  title?: string
  subtitle?: string
  tabs?: Tab[]
  listings?: Listing[]
}) {
  const [active, setActive] = useState(0)
  const id = useId()
  const shown = listings.filter((l) => l.category === tabs[active]?.id)

  /* Arrow keys move selection and focus together. This is required for the
     pattern, and it is the part a pair of onClick buttons never has. */
  function onKeyDown(event: React.KeyboardEvent) {
    if (event.key !== 'ArrowRight' && event.key !== 'ArrowLeft') return
    event.preventDefault()
    const next =
      event.key === 'ArrowRight'
        ? (active + 1) % tabs.length
        : (active - 1 + tabs.length) % tabs.length
    setActive(next)
    document.getElementById(`${id}-tab-${next}`)?.focus()
  }

  return (
    <section className="bg-white py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-6xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-3xl font-semibold tracking-tight text-balance text-gray-900 sm:text-4xl dark:text-white">
            {title}
          </h2>
          <p className="mt-4 text-base/7 text-pretty text-gray-600 dark:text-gray-400">
            {subtitle}
          </p>
        </div>

        <div className="mt-10 flex justify-center">
          <div
            role="tablist"
            aria-label="Listing type"
            className="flex border-b border-gray-200 dark:border-white/10"
          >
            {tabs.map((tab, i) => (
              <button
                key={tab.id}
                id={`${id}-tab-${i}`}
                type="button"
                role="tab"
                aria-selected={i === active}
                aria-controls={`${id}-panel`}
                tabIndex={i === active ? 0 : -1}
                onClick={() => setActive(i)}
                onKeyDown={onKeyDown}
                className={`-mb-px min-h-11 border-b-2 px-6 text-base font-medium focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 ${
                  i === active
                    ? 'border-indigo-600 text-indigo-600 dark:border-indigo-400 dark:text-indigo-400'
                    : 'border-transparent text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>
        </div>

        {/* Announced when the filter changes, so a keyboard user is told how
            many results they now have instead of the grid silently swapping. */}
        <p aria-live="polite" className="sr-only">
          {shown.length} {shown.length === 1 ? 'listing' : 'listings'} in{' '}
          {tabs[active]?.label}
        </p>

        <div id={`${id}-panel`} role="tabpanel" aria-labelledby={`${id}-tab-${active}`}>
          {shown.length > 0 ? (
            <ul role="list" className="mt-10 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
              {shown.map((listing) => (
                <li
                  key={listing.id}
                  className="group relative flex flex-col overflow-hidden rounded-xl border border-gray-200 bg-white transition-shadow hover:shadow-lg has-[a:focus-visible]:outline-2 has-[a:focus-visible]:outline-offset-2 has-[a:focus-visible]:outline-indigo-600 dark:border-white/10 dark:bg-gray-900"
                >
                  <div className="h-48 overflow-hidden bg-gray-50 dark:bg-white/5">
                    <img
                      src={listing.image}
                      alt=""
                      className="size-full object-cover transition-transform duration-300 group-hover:scale-105"
                    />
                  </div>

                  <div className="flex flex-1 flex-col p-6">
                    <div className="flex items-start justify-between gap-3">
                      <p className="text-2xl font-bold text-gray-900 tabular-nums dark:text-white">
                        {listing.price}
                      </p>
                      <span
                        className={`rounded px-2 py-1 text-xs font-semibold ${
                          CONDITION_TONE[listing.condition] ??
                          'bg-gray-100 text-gray-700 dark:bg-white/10 dark:text-gray-300'
                        }`}
                      >
                        {listing.condition}
                      </span>
                    </div>

                    <h3 className="mt-3 line-clamp-2 text-lg font-semibold text-gray-900 dark:text-white">
                      {/* A link, not a div with an onClick. */}
                      <a
                        href={listing.href ?? '#'}
                        className="after:absolute after:inset-0 focus:outline-none"
                      >
                        {listing.title}
                      </a>
                    </h3>
                    <p className="mt-2 line-clamp-2 text-sm/6 text-gray-600 dark:text-gray-400">
                      {listing.description}
                    </p>

                    <p className="mt-auto flex items-center gap-2 pt-4 text-sm text-gray-500 dark:text-gray-400">
                      <MapPin aria-hidden="true" className="size-4 flex-none" />
                      <span className="sr-only">Location: </span>
                      {listing.location}
                    </p>
                  </div>
                </li>
              ))}
            </ul>
          ) : (
            <div className="mt-16 text-center">
              <SearchX aria-hidden="true" className="mx-auto size-12 text-gray-400" />
              <h3 className="mt-4 text-lg font-medium text-gray-900 dark:text-white">
                Nothing listed here yet
              </h3>
              <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
                There are no {tabs[active]?.label.toLowerCase()} listings at the moment.
              </p>
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
