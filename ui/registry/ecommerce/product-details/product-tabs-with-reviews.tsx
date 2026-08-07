'use client'

import { useId, useRef, useState } from 'react'
import { Star } from 'lucide-react'

/*
 * The panel below a product: description, reviews, shipping and returns.
 *
 * A real tablist. The source rendered three plain buttons and three
 * conditional divs, so nothing announced that the buttons were tabs, which one
 * was current, or what any of them controlled, and the arrow keys did nothing.
 * Here the tabs carry role, aria-selected and aria-controls, arrow keys move
 * between them with Home and End at the ends, and one roving tabindex means
 * Tab steps over the group rather than through it.
 *
 * Tab labels are stored, not derived. The source mapped over the strings
 * 'description', 'reviews', 'support' and made them look right with a
 * capitalize class. CSS text-transform changes the pixels and not the text, so
 * the accessible name stayed lowercase.
 *
 * The rating distribution is passed in, not computed from the reviews on this
 * page. Six loaded reviews cannot describe 157, and a bar chart derived from
 * whatever happened to render is a chart that changes when you press "show
 * all".
 *
 * Reviewers are initials rather than photographs. Putting a real person's face
 * next to review copy nobody wrote is the wrong default for a template, and
 * the source pulled its faces from randomuser.me, so every avatar was a live
 * dependency on a third-party service.
 *
 * Dates are <time datetime>. "March 14, 2021" as a bare string is unparseable
 * by anything, and ambiguous between locales the moment it becomes 03/14/21.
 */

export interface Review {
  id: string
  author: string
  rating: number
  /** ISO date. Rendered through toLocaleDateString. */
  date: string
  body: string
  verified?: boolean
}

export interface Tab {
  id: string
  label: string
}

const REVIEWS: Review[] = [
  {
    id: 'r1',
    author: 'Kristin Watson',
    rating: 5,
    date: '2026-03-14',
    body: 'Runs about half a size small, so size up. Past that it is the most comfortable thing I own and the stitching has not moved after four months of daily wear.',
    verified: true,
  },
  {
    id: 'r2',
    author: 'Jenny Wilson',
    rating: 5,
    date: '2026-01-28',
    body: 'Arrived two days early. The colour is closer to the second photo than the first, which is what I wanted, but worth knowing before you order.',
    verified: true,
  },
  {
    id: 'r3',
    author: 'Bessie Cooper',
    rating: 4,
    date: '2026-01-11',
    body: 'Good quality for the price. Docking a star because the packaging was enormous for what is inside.',
  },
  {
    id: 'r4',
    author: 'Devon Lane',
    rating: 5,
    date: '2025-12-02',
    body: 'Second one I have bought. The first is three years old and still going, which is why I did not think twice.',
    verified: true,
  },
  {
    id: 'r5',
    author: 'Arlene McCoy',
    rating: 3,
    date: '2025-11-19',
    body: 'Fine, but not remarkable. Does the job and I have no complaints, I just would not rave about it.',
  },
  {
    id: 'r6',
    author: 'Guy Hawkins',
    rating: 5,
    date: '2025-10-30',
    body: 'Customer service replaced a faulty one without asking me to send the original back first. That is why the rating is what it is.',
    verified: true,
  },
]

const TABS: Tab[] = [
  { id: 'description', label: 'Description' },
  { id: 'reviews', label: 'Reviews' },
  { id: 'shipping', label: 'Shipping and returns' },
]

/** Counts per star, 5 down to 1, over every review rather than the loaded page. */
const DISTRIBUTION = [104, 31, 14, 5, 3]

function Stars({ rating }: { rating: number }) {
  return (
    <span aria-hidden="true" className="flex">
      {[0, 1, 2, 3, 4].map((index) => (
        <Star
          key={index}
          className={`size-4 ${
            index < Math.round(rating)
              ? 'fill-amber-400 text-amber-400'
              : 'text-gray-300 dark:text-gray-600'
          }`}
        />
      ))}
    </span>
  )
}

function initials(name: string) {
  return name
    .split(' ')
    .map((part) => part[0])
    .slice(0, 2)
    .join('')
}

export default function ProductTabsWithReviews({
  description = 'Cut from a mid-weight cotton twill that softens without losing its shape, with a half-lining through the body so it sits flat over a jumper. Every seam is felled rather than overlocked, which is slower to make and does not fray. Machine washable at 30 degrees.',
  tabs = TABS,
  reviews = REVIEWS,
  totalReviews = 157,
  averageRating = 4.5,
  distribution = DISTRIBUTION,
  supportEmail = 'support@example.com',
  supportPhone = '+1 (800) 123-4567',
}: {
  description?: string
  tabs?: Tab[]
  reviews?: Review[]
  totalReviews?: number
  averageRating?: number
  distribution?: number[]
  supportEmail?: string
  supportPhone?: string
}) {
  const [active, setActive] = useState(0)
  const [expanded, setExpanded] = useState(false)
  const baseId = useId()
  const tabRefs = useRef<(HTMLButtonElement | null)[]>([])

  const tabId = (index: number) => `${baseId}-tab-${index}`
  const panelId = (index: number) => `${baseId}-panel-${index}`

  const visible = expanded ? reviews : reviews.slice(0, 3)

  function onKeyDown(event: React.KeyboardEvent) {
    const last = tabs.length - 1
    let next: number | null = null
    if (event.key === 'ArrowRight') next = active === last ? 0 : active + 1
    else if (event.key === 'ArrowLeft') next = active === 0 ? last : active - 1
    else if (event.key === 'Home') next = 0
    else if (event.key === 'End') next = last
    if (next === null) return
    event.preventDefault()
    setActive(next)
    tabRefs.current[next]?.focus()
  }

  const current = tabs[active]

  return (
    <section className="border-t border-gray-200 bg-white py-12 dark:border-white/10 dark:bg-gray-950">
      <div className="mx-auto max-w-4xl px-4">
        <div
          role="tablist"
          aria-label="Product information"
          onKeyDown={onKeyDown}
          className="flex gap-8 overflow-x-auto border-b border-gray-200 dark:border-white/10"
        >
          {tabs.map((tab, index) => (
            <button
              key={tab.id}
              ref={(node) => {
                tabRefs.current[index] = node
              }}
              type="button"
              role="tab"
              id={tabId(index)}
              aria-selected={index === active}
              aria-controls={panelId(index)}
              tabIndex={index === active ? 0 : -1}
              onClick={() => setActive(index)}
              className={`-mb-px shrink-0 border-b-2 px-1 py-4 text-sm font-medium whitespace-nowrap focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-indigo-600 ${
                index === active
                  ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white'
                  : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200'
              }`}
            >
              {tab.label}
              {tab.id === 'reviews' && ` (${totalReviews})`}
            </button>
          ))}
        </div>

        <div
          role="tabpanel"
          id={panelId(active)}
          aria-labelledby={tabId(active)}
          tabIndex={0}
          className="mt-8 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-indigo-600"
        >
          {current.id === 'description' && (
            <p className="text-pretty text-gray-600 dark:text-gray-300">{description}</p>
          )}

          {current.id === 'reviews' && (
            <div>
              <div className="grid gap-8 border-b border-gray-200 pb-8 sm:grid-cols-[auto_1fr] dark:border-white/10">
                <div className="text-center sm:text-left">
                  <p className="text-4xl font-bold text-gray-900 dark:text-white">
                    {averageRating}
                    <span className="text-lg font-normal text-gray-500 dark:text-gray-400">
                      {' '}
                      / 5
                    </span>
                  </p>
                  <div className="mt-1 flex justify-center sm:justify-start">
                    <Stars rating={averageRating} />
                  </div>
                  <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    {totalReviews} reviews
                  </p>
                </div>

                {/* Each bar is a labelled row, not a bare div with a width.
                    A chart nobody can read is decoration. */}
                <ul role="list" className="flex flex-col justify-center gap-1.5">
                  {distribution.map((count, index) => {
                    const stars = 5 - index
                    const percent = totalReviews ? Math.round((count / totalReviews) * 100) : 0
                    return (
                      <li key={stars} className="flex items-center gap-3 text-sm">
                        <span className="w-12 shrink-0 text-gray-600 dark:text-gray-400">
                          {stars} star
                        </span>
                        <span
                          aria-hidden="true"
                          className="h-2 flex-1 overflow-hidden rounded-full bg-gray-200 dark:bg-white/10"
                        >
                          <span
                            className="block h-full rounded-full bg-amber-400"
                            style={{ width: `${percent}%` }}
                          />
                        </span>
                        <span className="w-20 shrink-0 text-right whitespace-nowrap text-gray-500 tabular-nums dark:text-gray-400">
                          {count} <span className="sr-only">reviews, </span>
                          <span aria-hidden="true">({percent}%)</span>
                          <span className="sr-only">{percent} percent</span>
                        </span>
                      </li>
                    )
                  })}
                </ul>
              </div>

              <ul role="list" className="mt-8 space-y-8">
                {visible.map((review) => (
                  <li
                    key={review.id}
                    className="flex gap-4 border-b border-gray-200 pb-8 last:border-0 dark:border-white/10"
                  >
                    <span
                      aria-hidden="true"
                      className="flex size-12 shrink-0 items-center justify-center rounded-full bg-gray-100 text-sm font-semibold text-gray-600 dark:bg-white/10 dark:text-gray-300"
                    >
                      {initials(review.author)}
                    </span>

                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        <Stars rating={review.rating} />
                        <span className="sr-only">{review.rating} out of 5</span>
                      </div>

                      <p className="mt-1 text-sm font-medium text-gray-900 dark:text-white">
                        {review.author}
                        {review.verified && (
                          <span className="ml-2 rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-normal text-emerald-700 dark:bg-emerald-500/10 dark:text-emerald-300">
                            Verified purchase
                          </span>
                        )}
                      </p>

                      <p className="mt-0.5 text-sm text-gray-500 dark:text-gray-400">
                        <time dateTime={review.date}>
                          {new Date(review.date).toLocaleDateString('en-US', {
                            year: 'numeric',
                            month: 'long',
                            day: 'numeric',
                          })}
                        </time>
                      </p>

                      <p className="mt-2 text-sm text-pretty text-gray-600 dark:text-gray-300">
                        {review.body}
                      </p>
                    </div>
                  </li>
                ))}
              </ul>

              {reviews.length > 3 && (
                <>
                  <button
                    type="button"
                    onClick={() => setExpanded((current) => !current)}
                    /* aria-expanded and aria-controls, because this button
                       changes the length of the list above it rather than
                       navigating anywhere. */
                    aria-expanded={expanded}
                    aria-controls={panelId(active)}
                    className="mt-6 inline-flex min-h-11 items-center text-sm font-medium text-indigo-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
                  >
                    {expanded ? 'Show fewer reviews' : `Show all ${reviews.length} reviews`}
                  </button>

                  {/* The list grows below the button, so nothing about the
                      change is visible from where focus is sitting. */}
                  <p role="status" aria-live="polite" className="sr-only">
                    Showing {visible.length} of {reviews.length} loaded reviews.
                  </p>
                </>
              )}
            </div>
          )}

          {current.id === 'shipping' && (
            <div className="space-y-6">
              <div>
                <h3 className="font-medium text-gray-900 dark:text-white">Shipping</h3>
                <p className="mt-2 text-gray-600 dark:text-gray-300">
                  Free worldwide over $100, otherwise flat rate at checkout. Orders placed before
                  2pm ship the same working day and take three to five days to arrive.
                </p>
              </div>
              <div>
                <h3 className="font-medium text-gray-900 dark:text-white">Returns and exchanges</h3>
                <p className="mt-2 text-gray-600 dark:text-gray-300">
                  Thirty days from delivery, in original condition with tags attached. Return
                  postage is on us if the item is faulty or we sent the wrong thing.
                </p>
              </div>
              <div>
                <h3 className="font-medium text-gray-900 dark:text-white">Contact</h3>
                <p className="mt-2 text-gray-600 dark:text-gray-300">
                  <a href={`mailto:${supportEmail}`} className="underline">
                    {supportEmail}
                  </a>{' '}
                  or{' '}
                  <a href={`tel:${supportPhone.replace(/[^+\d]/g, '')}`} className="underline">
                    {supportPhone}
                  </a>
                  .
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
