'use client'

import { useEffect, useRef, useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'

/*
 * Categories as circular photo tiles on a scrolling rail.
 *
 * There is no page arithmetic. The source kept a visibleItems number in state
 * and rendered categories.slice(page * visibleItems, ...), while the number of
 * columns actually on screen came from grid-cols-3 sm:4 md:5 lg:6 xl:8. Those
 * two never agreed: at xl the grid had room for eight and the slice handed it
 * six, so the last page was ragged at every width except one. Here CSS decides
 * how many fit and the buttons scroll by whatever is currently visible, so
 * there is nothing to keep in sync.
 *
 * Scrolling is native. That gives touch swipe, trackpad, momentum, and
 * scrollbar dragging at no cost. The source hand-rolled touchstart, touchmove
 * and touchend with a 50px threshold and got touch only, plus a `sliding` flag
 * that ignored input for 500ms after every move.
 *
 * The rail is a focusable labelled region. A scroll container that is not
 * focusable cannot be scrolled with the keyboard at all, which is easy to miss
 * because the mouse and the arrows both work.
 *
 * The buttons disable at the ends, computed from scroll position rather than
 * from a page counter, so they are still right after a swipe or a scrollbar
 * drag. Nothing wraps: wrapping a rail that shows its position is disorienting.
 *
 * Four of the source's ten photographs were 404s, including one whose Unsplash
 * id was malformed, and two categories shared a single image.
 */

export interface Category {
  name: string
  href: string
  image: string
  /** Backdrop behind the photo while it loads. */
  tone: string
}

const CATEGORIES: Category[] = [
  {
    name: 'Beauty and fragrance',
    href: '#',
    image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=300&h=300&fit=crop&q=80',
    tone: 'bg-amber-50',
  },
  {
    name: 'Skincare',
    href: '#',
    image: 'https://images.unsplash.com/photo-1583209814683-c023dd293cc6?w=300&h=300&fit=crop&q=80',
    tone: 'bg-rose-50',
  },
  {
    name: "Men's fashion",
    href: '#',
    image: 'https://images.unsplash.com/photo-1487222477894-8943e31ef7b2?w=300&h=300&fit=crop&q=80',
    tone: 'bg-blue-50',
  },
  {
    name: "Women's fashion",
    href: '#',
    image: 'https://images.unsplash.com/photo-1525507119028-ed4c629a60a3?w=300&h=300&fit=crop&q=80',
    tone: 'bg-indigo-50',
  },
  {
    name: 'Shoes',
    href: '#',
    image: 'https://images.unsplash.com/photo-1560769629-975ec94e6a86?w=300&h=300&fit=crop&q=80',
    tone: 'bg-slate-50',
  },
  {
    name: 'Watches',
    href: '#',
    image: 'https://images.unsplash.com/photo-1522312346375-d1a52e2b99b3?w=300&h=300&fit=crop&q=80',
    tone: 'bg-teal-50',
  },
  {
    name: 'Electronics',
    href: '#',
    image: 'https://images.unsplash.com/photo-1498049794561-7780e7231661?w=300&h=300&fit=crop&q=80',
    tone: 'bg-gray-100',
  },
  {
    name: 'Cameras',
    href: '#',
    image: 'https://images.unsplash.com/photo-1526170375885-4d8ecf77b99f?w=300&h=300&fit=crop&q=80',
    tone: 'bg-zinc-100',
  },
  {
    name: 'Sports and outdoors',
    href: '#',
    image: 'https://images.unsplash.com/photo-1599058917212-d750089bc07e?w=300&h=300&fit=crop&q=80',
    tone: 'bg-purple-50',
  },
  {
    name: 'Automotive',
    href: '#',
    image: 'https://images.unsplash.com/photo-1504215680853-026ed2a45def?w=300&h=300&fit=crop&q=80',
    tone: 'bg-orange-50',
  },
  {
    name: 'Home and living',
    href: '#',
    image: 'https://images.unsplash.com/photo-1441986300917-64674bd600d8?w=300&h=300&fit=crop&q=80',
    tone: 'bg-emerald-50',
  },
  {
    name: 'Fabrics and craft',
    href: '#',
    image: 'https://images.unsplash.com/photo-1523381210434-271e8be1f52b?w=300&h=300&fit=crop&q=80',
    tone: 'bg-lime-50',
  },
]

export default function CircularCategoryRail({
  title = 'Shop by category',
  subtitle = 'Every department, one row.',
  categories = CATEGORIES,
}: {
  title?: string
  subtitle?: string
  categories?: Category[]
}) {
  const rail = useRef<HTMLUListElement>(null)
  const [atStart, setAtStart] = useState(true)
  const [atEnd, setAtEnd] = useState(false)

  /* Derived from scrollLeft, so a swipe or a scrollbar drag updates the
     buttons too. A page counter only knows about button presses. */
  function sync() {
    const el = rail.current
    if (!el) return
    setAtStart(el.scrollLeft <= 1)
    setAtEnd(el.scrollLeft >= el.scrollWidth - el.clientWidth - 1)
  }

  useEffect(() => {
    sync()
    const el = rail.current
    if (!el) return
    /* ResizeObserver as well as scroll: rotating a phone changes how many
       tiles fit, which can put a rail that was scrollable at the end. */
    const observer = new ResizeObserver(sync)
    observer.observe(el)
    return () => observer.disconnect()
  }, [])

  function scrollByPage(direction: 1 | -1) {
    const el = rail.current
    if (!el) return
    el.scrollBy({ left: direction * el.clientWidth, behavior: 'smooth' })
  }

  return (
    <section className="border-y border-amber-100/60 bg-gradient-to-b from-amber-50/70 to-amber-50/30 py-8 dark:border-white/10 dark:from-gray-900 dark:to-gray-950">
      <div className="mx-auto max-w-7xl px-4 md:px-8">
        <div className="mb-6 flex items-end justify-between gap-4">
          <div>
            <div className="flex items-center gap-3">
              <span aria-hidden="true" className="h-8 w-1.5 rounded-full bg-amber-500" />
              <h2 className="text-xl font-bold text-gray-800 md:text-2xl dark:text-white">
                {title}
              </h2>
            </div>
            <p className="mt-1 ml-[1.125rem] text-sm text-gray-500 dark:text-gray-400">{subtitle}</p>
          </div>

          <div className="flex shrink-0 gap-3">
            <button
              type="button"
              onClick={() => scrollByPage(-1)}
              disabled={atStart}
              className="inline-flex size-11 items-center justify-center rounded-full border border-amber-200 bg-white text-gray-600 shadow-sm transition hover:bg-amber-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600 disabled:opacity-40 dark:border-white/15 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-white/5"
            >
              <ChevronLeft aria-hidden="true" className="size-5" />
              <span className="sr-only">Scroll categories left</span>
            </button>
            <button
              type="button"
              onClick={() => scrollByPage(1)}
              disabled={atEnd}
              className="inline-flex size-11 items-center justify-center rounded-full border border-amber-200 bg-white text-gray-600 shadow-sm transition hover:bg-amber-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600 disabled:opacity-40 dark:border-white/15 dark:bg-gray-900 dark:text-gray-300 dark:hover:bg-white/5"
            >
              <ChevronRight aria-hidden="true" className="size-5" />
              <span className="sr-only">Scroll categories right</span>
            </button>
          </div>
        </div>

        {/* tabIndex and a label make this a scroll region a keyboard user can
            reach and move with the arrow keys. Without them the rail is
            mouse-and-touch only, which is invisible in testing because both of
            those work fine. */}
        <ul
          ref={rail}
          role="group"
          aria-label={title}
          tabIndex={0}
          onScroll={sync}
          className="flex snap-x snap-mandatory gap-4 overflow-x-auto pb-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600 md:gap-6"
        >
          {categories.map((category) => (
            <li key={category.name} className="w-24 shrink-0 snap-start md:w-28">
              <a
                href={category.href}
                className="group flex flex-col items-center gap-2 rounded-lg p-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-600"
              >
                <span
                  className={`block aspect-square w-full overflow-hidden rounded-full ring-2 ring-white transition-shadow group-hover:shadow-lg ${category.tone}`}
                >
                  {/* Decorative: the name is directly beneath it in text. */}
                  <img
                    src={category.image}
                    alt=""
                    loading="lazy"
                    className="size-full object-cover transition-transform duration-500 group-hover:scale-110 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                  />
                </span>
                <span className="text-center text-xs leading-tight font-medium text-gray-700 group-hover:text-gray-900 md:text-sm dark:text-gray-300 dark:group-hover:text-white">
                  {category.name}
                </span>
              </a>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
