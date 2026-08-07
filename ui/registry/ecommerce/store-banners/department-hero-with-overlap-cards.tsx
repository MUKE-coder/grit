'use client'

import { useId, useState } from 'react'
import { ChevronLeft, ChevronRight } from 'lucide-react'

/*
 * A wide promo carousel with department cards riding up over its foot.
 *
 * The carousel does not advance on its own. That is a decision, not an
 * omission: the moment something moves by itself for more than five seconds
 * you owe the visitor a pause control, a live region, and a reduced-motion
 * path — see `hero-carousel-with-controls` in this same category, which pays
 * all of that. Here the cards below are the point of the block and the strip
 * is a header, so it waits to be asked. Arrows and dots, nothing on a timer.
 *
 * Off-screen slides are `invisible`, not just translated out of frame. A slide
 * parked at translate-x-full is still a focusable link sitting somewhere off
 * to the right, and tabbing into one scrolls the page sideways to a thing
 * nobody can see. Visibility steps at the end of the transition, so the slide
 * still slides.
 *
 * The four department cards share one component driven by a `layout` field.
 * The source wrote the same card body four times behind `index === 0`,
 * `index === 1`, `index === 2`, `index === 3` — which means four places to fix
 * every time the card design changes, and a fifth department silently renders
 * an empty box.
 *
 * Every tile is a link. In the source the sub-category tiles were plain divs
 * with a single "See all" link underneath, so the twelve things that look most
 * clickable on the block were the twelve that were not.
 */

export interface Slide {
  id: string
  title: string
  href: string
  image: string
}

export interface Tile {
  name: string
  href: string
  image: string
}

export type Department =
  | { layout: 'single'; title: string; href: string; linkText: string; image: string }
  | { layout: 'tiles'; title: string; href: string; linkText: string; tiles: Tile[] }
  | {
      layout: 'feature'
      title: string
      href: string
      linkText: string
      image: string
      caption: string
      tiles: Tile[]
    }

const SLIDES: Slide[] = [
  {
    id: 'gifts',
    title: 'Gifts they will actually use',
    href: '#',
    image: 'https://images.unsplash.com/photo-1512909006721-3d6018887383?w=1600&h=700&fit=crop&q=80',
  },
  {
    id: 'fashion',
    title: 'Fashion finds under $50',
    href: '#',
    image: 'https://images.unsplash.com/photo-1441984904996-e0b6ba687e04?w=1600&h=700&fit=crop&q=80',
  },
  {
    id: 'kitchen',
    title: 'Kitchen essentials',
    href: '#',
    image: 'https://images.unsplash.com/photo-1556909212-d5b604d0c90d?w=1600&h=700&fit=crop&q=80',
  },
  {
    id: 'toys',
    title: 'Toys for little ones',
    href: '#',
    image: 'https://images.unsplash.com/photo-1596461404969-9ae70f2830c1?w=1600&h=700&fit=crop&q=80',
  },
]

const shot = (id: string, size = 300) =>
  `https://images.unsplash.com/photo-${id}?w=${size}&h=${size}&fit=crop&q=80`

const DEPARTMENTS: Department[] = [
  {
    layout: 'single',
    title: 'Get your game on',
    href: '#',
    linkText: 'Shop gaming',
    image: shot('1542751371-adc38448a05e', 600),
  },
  {
    layout: 'tiles',
    title: 'Shop your home essentials',
    href: '#',
    linkText: 'Discover more in Home',
    tiles: [
      { name: 'Cleaning', href: '#', image: shot('1563453392212-326f5e854473') },
      { name: 'Storage', href: '#', image: shot('1595428774223-ef52624120d2') },
      { name: 'Decor', href: '#', image: shot('1513161455079-7dc1de15ef3e') },
      { name: 'Bedding', href: '#', image: shot('1505693416388-ac5ce068fe85') },
    ],
  },
  {
    layout: 'feature',
    title: 'Top categories in kitchen',
    href: '#',
    linkText: 'Explore all of Kitchen',
    image: shot('1590794056226-79ef3a8147e1', 500),
    caption: 'Cookware',
    tiles: [
      { name: 'Coffee', href: '#', image: shot('1509042239860-f550ce710b93') },
      { name: 'Espresso', href: '#', image: shot('1610889556528-9a770e32642f') },
      { name: 'Fridges', href: '#', image: shot('1571175443880-49e1d25b2bc5') },
    ],
  },
  {
    layout: 'tiles',
    title: 'Shop by room',
    href: '#',
    linkText: 'See all rooms',
    tiles: [
      { name: 'Living room', href: '#', image: shot('1567016432779-094069958ea5') },
      { name: 'Bedroom', href: '#', image: shot('1522708323590-d24dbb6b0267') },
      { name: 'Home office', href: '#', image: shot('1600494603989-9650cf6ddd3d') },
      { name: 'Bathroom', href: '#', image: shot('1584622650111-993a426fbf0a') },
    ],
  },
]

function TileLink({ tile, height }: { tile: Tile; height: string }) {
  return (
    <a
      href={tile.href}
      className="group block text-center focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
    >
      <span className={`block overflow-hidden rounded-md ${height}`}>
        <img
          src={tile.image}
          alt=""
          className="size-full object-cover transition-transform duration-500 group-hover:scale-105 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
        />
      </span>
      <span className="mt-1 block text-xs font-medium text-gray-700 group-hover:underline dark:text-gray-300">
        {tile.name}
      </span>
    </a>
  )
}

export default function DepartmentHeroWithOverlapCards({
  slides = SLIDES,
  departments = DEPARTMENTS,
  label = 'Featured promotions',
}: {
  slides?: Slide[]
  departments?: Department[]
  label?: string
}) {
  const [index, setIndex] = useState(0)
  const slidesId = useId()

  const go = (next: number) => setIndex(((next % slides.length) + slides.length) % slides.length)

  return (
    <section className="w-full bg-gray-100 pb-8 dark:bg-gray-950">
      <div
        aria-roledescription="carousel"
        aria-label={label}
        className="relative h-[400px] w-full overflow-hidden"
      >
        {/* Live is polite unconditionally: nothing here moves unless somebody
            moved it, so every change is one they asked to hear about. */}
        <div id={slidesId} aria-live="polite" className="absolute inset-0">
          {slides.map((slide, position) => (
            /* The group wraps the link rather than being the link. Putting
               role="group" on the <a> itself would overwrite its link role,
               and a slide you cannot tell is clickable is worse than a slide
               without a roledescription. */
            <div
              key={slide.id}
              role="group"
              aria-roledescription="slide"
              aria-label={`${position + 1} of ${slides.length}`}
              className={`absolute inset-0 transition-[transform,visibility] duration-500 ease-in-out ${
                position === index
                  ? 'visible translate-x-0'
                  : position < index
                    ? 'invisible -translate-x-full'
                    : 'invisible translate-x-full'
              } motion-reduce:transition-none`}
            >
              <a href={slide.href} className="block size-full">
                {/* Empty alt: the headline below is inside the same link, so
                    it already names the destination. Repeating it here would
                    have the link read its own title twice. */}
                <img src={slide.image} alt="" className="size-full object-cover" />
                <span className="absolute inset-x-0 bottom-0 bg-gradient-to-t from-black/70 to-transparent p-6 pb-28">
                  <span className="mx-auto block max-w-7xl px-4 text-2xl font-bold text-white md:text-3xl">
                    {slide.title}
                  </span>
                </span>
              </a>
            </div>
          ))}
        </div>

        <button
          type="button"
          onClick={() => go(index - 1)}
          aria-controls={slidesId}
          className="absolute top-1/2 left-4 z-10 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-sm bg-white/80 text-gray-800 shadow-md hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
        >
          <ChevronLeft aria-hidden="true" className="size-6" />
          <span className="sr-only">Previous promotion</span>
        </button>
        <button
          type="button"
          onClick={() => go(index + 1)}
          aria-controls={slidesId}
          className="absolute top-1/2 right-4 z-10 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-sm bg-white/80 text-gray-800 shadow-md hover:bg-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
        >
          <ChevronRight aria-hidden="true" className="size-6" />
          <span className="sr-only">Next promotion</span>
        </button>

        {/* Above the overlap so the cards do not bury the dots, and on their
            own scrim: white-on-white is what a bare dot row becomes the first
            time somebody uploads a bright photograph. */}
        <ul
          role="list"
          className="absolute inset-x-0 bottom-24 z-10 mx-auto flex w-fit items-center gap-2 rounded-full bg-black/40 px-3 backdrop-blur-sm"
        >
          {slides.map((slide, position) => (
            <li key={slide.id}>
              <button
                type="button"
                onClick={() => go(position)}
                aria-controls={slidesId}
                aria-current={position === index ? 'true' : undefined}
                className="inline-flex h-11 items-center px-0.5 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
              >
                <span
                  aria-hidden="true"
                  className={`block h-1.5 rounded-full transition-all ${
                    position === index ? 'w-8 bg-white' : 'w-4 bg-white/50'
                  }`}
                />
                <span className="sr-only">{slide.title}</span>
              </button>
            </li>
          ))}
        </ul>
      </div>

      {/* The overlap is the whole visual idea: the cards climb over the foot of
          the promo strip. -mt-20 has to stay smaller than the strip's bottom
          padding above, or the cards eat the slide's own headline. */}
      <div className="relative z-20 -mt-20 px-4">
        <ul
          role="list"
          className="mx-auto grid max-w-7xl grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4"
        >
          {departments.map((department) => (
            <li
              key={department.title}
              className="flex flex-col rounded-lg bg-white p-4 shadow-lg transition-shadow hover:shadow-xl dark:bg-gray-900"
            >
              <h3 className="mb-3 text-lg leading-tight font-semibold text-gray-900 dark:text-white">
                {department.title}
              </h3>

              {department.layout === 'single' && (
                <a
                  href={department.href}
                  className="group block overflow-hidden rounded-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                >
                  <img
                    src={department.image}
                    alt=""
                    className="h-48 w-full object-cover transition-transform duration-500 group-hover:scale-105 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                  />
                </a>
              )}

              {department.layout === 'tiles' && (
                <ul role="list" className="grid grid-cols-2 gap-2">
                  {department.tiles.map((tile) => (
                    <li key={tile.name}>
                      <TileLink tile={tile} height="h-20" />
                    </li>
                  ))}
                </ul>
              )}

              {department.layout === 'feature' && (
                <>
                  <a
                    href={department.href}
                    className="group block overflow-hidden rounded-lg focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                  >
                    <img
                      src={department.image}
                      alt=""
                      className="h-32 w-full object-cover transition-transform duration-500 group-hover:scale-105 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                    />
                    <span className="mt-2 block text-sm font-medium text-gray-800 group-hover:underline dark:text-gray-200">
                      {department.caption}
                    </span>
                  </a>
                  <ul role="list" className="mt-3 grid grid-cols-3 gap-1">
                    {department.tiles.map((tile) => (
                      <li key={tile.name}>
                        <TileLink tile={tile} height="h-16" />
                      </li>
                    ))}
                  </ul>
                </>
              )}

              {/* mt-auto pins every card's link to its own foot. Grid makes
                  the four cards equal height, and without this the link floats
                  wherever its content happened to end — four links at four
                  different heights across one row. */}
              <a
                href={department.href}
                className="mt-auto inline-block pt-3 text-sm font-medium text-blue-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600 dark:text-blue-400"
              >
                {department.linkText}
              </a>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
