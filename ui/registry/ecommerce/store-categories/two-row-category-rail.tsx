'use client'

import { useEffect, useRef, useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'

/*
 * Two rows of circular departments that scroll together, one column at a time.
 *
 * The sibling `circular-category-rail` is one row and lets CSS decide how many
 * fit. This one is for a catalogue with forty departments, where one row means
 * either a rail nobody scrolls to the end of or a second component underneath.
 *
 * Two rows is a layout problem, not a data problem. The tempting version splits
 * the array in half and renders two independent scrollers, which drift out of
 * step the moment either one is touched. Here the list is a single grid with
 * `grid-rows-2` and `grid-flow-col`: the browser fills down then across, so one
 * scroll container moves both rows and they cannot disagree.
 *
 * An odd number of categories leaves a gap in the bottom row rather than
 * rebalancing the columns. That is deliberate: the alternative reorders the
 * departments, and the order was somebody's merchandising decision.
 *
 * Scrolling is by `clientWidth`, so the buttons move by whatever is visible.
 * A page count computed from a column width has to be kept in sync with the CSS
 * that decides that width, and it never is.
 *
 * The arrow buttons are hidden from assistive technology. They are a
 * convenience over a scroll container that is already reachable with Tab and
 * arrow keys, and announcing them adds two controls that do nothing new.
 */

export interface RailCategory {
  name: string
  href: string
  image: string
  /** Backdrop behind the photo while it loads. Any Tailwind background class. */
  tone?: string
}

const CATEGORIES: RailCategory[] = [
  { name: 'Women', href: '#', image: 'https://images.unsplash.com/photo-1483985988355-763728e1935b?w=300&h=300&fit=crop&q=80', tone: 'bg-rose-50' },
  { name: 'Curve', href: '#', image: 'https://images.unsplash.com/photo-1445205170230-053b83016050?w=300&h=300&fit=crop&q=80', tone: 'bg-amber-50' },
  { name: 'Kids', href: '#', image: 'https://images.unsplash.com/photo-1503944583220-79d8926ad5e2?w=300&h=300&fit=crop&q=80', tone: 'bg-sky-50' },
  { name: 'Men', href: '#', image: 'https://images.unsplash.com/photo-1516257984-b1b4d707412e?w=300&h=300&fit=crop&q=80', tone: 'bg-slate-100' },
  { name: 'Baby & Maternity', href: '#', image: 'https://images.unsplash.com/photo-1519689680058-324335c77eba?w=300&h=300&fit=crop&q=80', tone: 'bg-blue-50' },
  { name: 'Tops', href: '#', image: 'https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=300&h=300&fit=crop&q=80', tone: 'bg-neutral-100' },
  { name: 'Dresses', href: '#', image: 'https://images.unsplash.com/photo-1595777457583-95e059d581b8?w=300&h=300&fit=crop&q=80', tone: 'bg-pink-50' },
  { name: 'Cell Phones', href: '#', image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=300&h=300&fit=crop&q=80', tone: 'bg-violet-50' },
  { name: 'Shoes', href: '#', image: 'https://images.unsplash.com/photo-1549298916-b41d501d3772?w=300&h=300&fit=crop&q=80', tone: 'bg-orange-50' },
  { name: 'Jewelry', href: '#', image: 'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=300&h=300&fit=crop&q=80', tone: 'bg-yellow-50' },
  { name: 'Home & Living', href: '#', image: 'https://images.unsplash.com/photo-1522771739844-6a9f6d5f14af?w=300&h=300&fit=crop&q=80', tone: 'bg-emerald-50' },
  { name: 'Beauty & Health', href: '#', image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=300&h=300&fit=crop&q=80', tone: 'bg-fuchsia-50' },
  { name: 'Sports & Outdoor', href: '#', image: 'https://images.unsplash.com/photo-1517836357463-d25dfeac3438?w=300&h=300&fit=crop&q=80', tone: 'bg-lime-50' },
  { name: 'Beachwear', href: '#', image: 'https://images.unsplash.com/photo-1507525428034-b723cf961d3e?w=300&h=300&fit=crop&q=80', tone: 'bg-cyan-50' },
  { name: 'Underwear', href: '#', image: 'https://images.unsplash.com/photo-1618354691373-d851c5c3a990?w=300&h=300&fit=crop&q=80', tone: 'bg-stone-100' },
  { name: 'Toys & Games', href: '#', image: 'https://images.unsplash.com/photo-1558060370-d644479cb6f7?w=300&h=300&fit=crop&q=80', tone: 'bg-indigo-50' },
]

export default function TwoRowCategoryRail({
  title,
  categories = CATEGORIES,
}: {
  /** Optional. The rail reads fine with no heading above it. */
  title?: string
  categories?: RailCategory[]
}) {
  const rail = useRef<HTMLUListElement>(null)
  const [atStart, setAtStart] = useState(true)
  const [atEnd, setAtEnd] = useState(false)

  /* Derived from scrollLeft, so a swipe or a scrollbar drag updates the buttons
     too. A page counter only knows about button presses. */
  function sync() {
    const el = rail.current
    if (!el) return
    setAtStart(el.scrollLeft <= 1)
    setAtEnd(el.scrollLeft + el.clientWidth >= el.scrollWidth - 1)
  }

  useEffect(() => {
    sync()
    const el = rail.current
    if (!el) return
    // The content can change width without the window resizing: a font loads,
    // an image arrives, the container is animated open. Watching the element
    // catches all of those; a window resize listener catches none of them.
    const observer = new ResizeObserver(sync)
    observer.observe(el)
    return () => observer.disconnect()
  }, [categories])

  function scrollBy(direction: -1 | 1) {
    const el = rail.current
    if (!el) return
    el.scrollBy({ left: direction * el.clientWidth, behavior: 'smooth' })
  }

  return (
    <section className="bg-gray-50 py-10 dark:bg-gray-950">
      <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {title && (
          <h2 className="mb-6 text-xl font-semibold tracking-tight text-gray-900 dark:text-white">
            {title}
          </h2>
        )}

        <button
          type="button"
          onClick={() => scrollBy(-1)}
          disabled={atStart}
          aria-hidden="true"
          tabIndex={-1}
          className="absolute left-0 top-1/2 z-10 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white shadow-md transition-opacity disabled:pointer-events-none disabled:opacity-0 md:flex dark:bg-gray-800"
        >
          <ChevronLeft className="h-5 w-5 text-gray-700 dark:text-gray-200" />
        </button>

        {/*
          grid-rows-2 with grid-flow-col is the whole trick: items fill down the
          first column, then the second, so one scroller drives both rows.
          Splitting the array into two scrollers puts them out of step forever.
        */}
        <ul
          ref={rail}
          onScroll={sync}
          className="grid snap-x snap-mandatory grid-flow-col grid-rows-2 gap-x-4 gap-y-6 overflow-x-auto scroll-smooth px-1 py-2 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
        >
          {categories.map((category) => (
            <li key={category.name} className="w-24 shrink-0 snap-start md:w-28">
              <a
                href={category.href}
                className="group flex flex-col items-center gap-2 rounded-lg p-1 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900"
              >
                <span
                  className={`block aspect-square w-full overflow-hidden rounded-full transition-shadow group-hover:shadow-lg ${category.tone ?? 'bg-gray-100'}`}
                >
                  {/* Decorative: the name is directly beneath it in text. */}
                  <img
                    src={category.image}
                    alt=""
                    loading="lazy"
                    className="h-full w-full object-cover"
                  />
                </span>
                <span className="text-center text-sm leading-tight text-gray-800 group-hover:underline dark:text-gray-200">
                  {category.name}
                </span>
              </a>
            </li>
          ))}
        </ul>

        <button
          type="button"
          onClick={() => scrollBy(1)}
          disabled={atEnd}
          aria-hidden="true"
          tabIndex={-1}
          className="absolute right-0 top-1/2 z-10 hidden h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full bg-white shadow-md transition-opacity disabled:pointer-events-none disabled:opacity-0 md:flex dark:bg-gray-800"
        >
          <ChevronRight className="h-5 w-5 text-gray-700 dark:text-gray-200" />
        </button>
      </div>
    </section>
  )
}
