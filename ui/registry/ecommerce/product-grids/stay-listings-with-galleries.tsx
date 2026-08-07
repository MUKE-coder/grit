'use client'

import { useState } from 'react'
import { ChevronLeft, ChevronRight, Heart, Star } from 'lucide-react'

/*
 * Stay listings, each card with its own small gallery.
 *
 * The gallery is where this pattern usually breaks, and the source broke it in
 * the most complete way available: the previous and next buttons were rendered
 * only while `showControls` was true, and `showControls` was set by
 * onMouseEnter. Not hidden — absent from the DOM. A keyboard user could never
 * reach them, and a touch user only got them if their browser synthesised a
 * hover. The controls here are always in the document and fade in on hover OR
 * focus-within, so they can be reached and are visible once they are.
 *
 * Only the current image is exposed. In the source every image in every
 * gallery sat in the DOM at `opacity-0`, with a real alt, so a screen reader
 * read four near-identical descriptions per card and twenty-four for the grid.
 * The off-screen ones are `aria-hidden` and `inert`, which takes them out of
 * both the accessibility tree and the tab order.
 *
 * Changing image announces the position through a live region. Without it,
 * pressing "next" moves a picture a screen reader user cannot see and reports
 * nothing at all.
 *
 * The dots use `aria-current`, so the active one is announced rather than only
 * being wider than the others.
 *
 * The title is a link with its hit area stretched over the card. In the source
 * nothing on the card opened the listing — the card was a div and the only
 * interactive parts were the gallery controls.
 */

export interface Listing {
  id: string
  title: string
  location: string
  images: string[]
  /** Per night. */
  price: number
  rating: number
  reviewCount: number
  dates: string
  badges?: string[]
  saved?: boolean
  href?: string
}

const CURRENCY = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
})

const LISTINGS: Listing[] = [
  {
    id: '1',
    title: 'Beachfront villa with infinity pool',
    location: 'Malibu, California',
    images: [
      'https://images.unsplash.com/photo-1512917774080-9991f1c4c750?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1600585152220-90363fe7e115?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1600566753376-12c8ab7fb75b?w=700&h=700&fit=crop&q=80',
    ],
    price: 899,
    rating: 4.98,
    reviewCount: 156,
    dates: '1–6 Jun',
    badges: ['Beachfront'],
    saved: true,
  },
  {
    id: '2',
    title: 'Modern farmhouse with mountain views',
    location: 'Asheville, North Carolina',
    images: [
      'https://images.unsplash.com/photo-1623298317883-6b70254edf31?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1505692952047-1a78307da8f2?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1618221118493-9cfa1a1c00da?w=700&h=700&fit=crop&q=80',
    ],
    price: 345,
    rating: 4.92,
    reviewCount: 86,
    dates: '10–15 Jul',
    badges: ['Countryside'],
  },
  {
    id: '3',
    title: 'Penthouse with city skyline views',
    location: 'Chicago, Illinois',
    images: [
      'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1554995207-c18c203602cb?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1567767292278-a4f21aa2d36e?w=700&h=700&fit=crop&q=80',
    ],
    price: 429,
    rating: 4.88,
    reviewCount: 104,
    dates: '20–25 Aug',
    badges: ['New listing'],
  },
  {
    id: '4',
    title: 'Rustic cabin in the woods',
    location: 'Big Sur, California',
    images: [
      'https://images.unsplash.com/photo-1449158743715-0a90ebb6d2d8?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1464288550599-43d5a73451b8?w=700&h=700&fit=crop&q=80',
      'https://images.unsplash.com/photo-1470770841072-f978cf4d019e?w=700&h=700&fit=crop&q=80',
    ],
    price: 275,
    rating: 4.96,
    reviewCount: 132,
    dates: '5–10 Sep',
    badges: ['Nature'],
  },
]

function ListingCard({
  listing,
  saved,
  onToggleSave,
}: {
  listing: Listing
  saved: boolean
  onToggleSave: () => void
}) {
  const [index, setIndex] = useState(0)
  const count = listing.images.length

  const go = (next: number) => setIndex((next + count) % count)

  return (
    <li className="group/card relative flex flex-col">
      <div className="group/gallery relative mb-3 aspect-square overflow-hidden rounded-xl bg-gray-100 dark:bg-white/5">
        {listing.images.map((image, i) => {
          const current = i === index
          return (
            <img
              key={image}
              src={image}
              alt={current ? `${listing.title}, image ${i + 1} of ${count}` : ''}
              /* Off-screen slides leave both the accessibility tree and the
                 tab order, rather than sitting there at opacity 0 being read
                 out. */
              aria-hidden={!current}
              inert={!current ? true : undefined}
              className={`absolute inset-0 size-full object-cover transition-opacity duration-300 ${
                current ? 'opacity-100' : 'opacity-0'
              }`}
            />
          )
        })}

        {count > 1 && (
          <>
            {/* Always in the DOM. The source rendered these only on hover, so
                they did not exist for a keyboard user at all. */}
            <button
              type="button"
              onClick={() => go(index - 1)}
              className="absolute top-1/2 left-2 z-10 inline-flex size-11 -translate-y-1/2 items-center justify-center rounded-full bg-white/90 text-gray-800 opacity-0 shadow-md transition-opacity hover:bg-white focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 group-hover/gallery:opacity-100"
            >
              <ChevronLeft aria-hidden="true" className="size-4" />
              <span className="sr-only">Previous image of {listing.title}</span>
            </button>
            <button
              type="button"
              onClick={() => go(index + 1)}
              className="absolute top-1/2 right-2 z-10 inline-flex size-11 -translate-y-1/2 items-center justify-center rounded-full bg-white/90 text-gray-800 opacity-0 shadow-md transition-opacity hover:bg-white focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 group-hover/gallery:opacity-100"
            >
              <ChevronRight aria-hidden="true" className="size-4" />
              <span className="sr-only">Next image of {listing.title}</span>
            </button>

            <div className="absolute bottom-3 left-1/2 z-10 flex -translate-x-1/2 gap-1.5">
              {listing.images.map((image, i) => (
                <button
                  key={image}
                  type="button"
                  onClick={() => setIndex(i)}
                  aria-current={i === index ? 'true' : undefined}
                  className={`h-1.5 rounded-full transition-all focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white ${
                    i === index ? 'w-6 bg-white' : 'w-1.5 bg-white/60 hover:bg-white/80'
                  }`}
                >
                  <span className="sr-only">Image {i + 1}</span>
                </button>
              ))}
            </div>
          </>
        )}

        {/* Announced on change. Without this, pressing next moves a picture
            and reports nothing. */}
        <p aria-live="polite" className="sr-only">
          Image {index + 1} of {count}
        </p>

        <button
          type="button"
          onClick={onToggleSave}
          aria-pressed={saved}
          className="absolute top-2 right-2 z-10 inline-flex size-11 items-center justify-center rounded-full focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
        >
          <Heart
            aria-hidden="true"
            className={`size-6 ${
              saved ? 'fill-rose-500 text-white' : 'fill-black/20 text-white'
            }`}
          />
          <span className="sr-only">Save {listing.title}</span>
        </button>

        {listing.badges && listing.badges.length > 0 && (
          <ul role="list" className="absolute top-3 left-3 z-10 flex max-w-[70%] flex-wrap gap-1.5">
            {listing.badges.map((badge) => (
              <li
                key={badge}
                className="rounded-md bg-white/90 px-2 py-1 text-xs font-medium text-gray-900 backdrop-blur-sm"
              >
                {badge}
              </li>
            ))}
          </ul>
        )}
      </div>

      <div className="flex items-start justify-between gap-3">
        <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
          <a
            href={listing.href ?? '#'}
            className="after:absolute after:inset-0 focus:outline-none"
          >
            {listing.title}
          </a>
        </h3>
        <p className="flex flex-none items-center gap-1 text-sm text-gray-900 dark:text-white">
          <Star aria-hidden="true" className="size-4 fill-current" />
          <span className="sr-only">Rated </span>
          {listing.rating}
          <span className="sr-only"> out of 5 from {listing.reviewCount} reviews</span>
        </p>
      </div>
      <p className="text-sm text-gray-500 dark:text-gray-400">{listing.location}</p>
      <p className="text-sm text-gray-500 dark:text-gray-400">{listing.dates}</p>
      <p className="mt-1 text-sm text-gray-900 dark:text-white">
        <span className="font-semibold tabular-nums">{CURRENCY.format(listing.price)}</span> night
      </p>
    </li>
  )
}

export default function StayListingsWithGalleries({
  title = 'Places to stay',
  listings = LISTINGS,
  onToggleSave,
}: {
  title?: string
  listings?: Listing[]
  onToggleSave?: (listing: Listing, saved: boolean) => void
}) {
  const [saved, setSaved] = useState<Record<string, boolean>>(() =>
    Object.fromEntries(listings.map((l) => [l.id, Boolean(l.saved)])),
  )

  return (
    <section className="bg-white py-16 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <h2 className="text-2xl font-semibold tracking-tight text-gray-900 dark:text-white">
          {title}
        </h2>

        <ul
          role="list"
          className="mt-8 grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-4"
        >
          {listings.map((listing) => (
            <ListingCard
              key={listing.id}
              listing={listing}
              saved={Boolean(saved[listing.id])}
              onToggleSave={() => {
                const next = !saved[listing.id]
                setSaved((current) => ({ ...current, [listing.id]: next }))
                onToggleSave?.(listing, next)
              }}
            />
          ))}
        </ul>
      </div>
    </section>
  )
}
