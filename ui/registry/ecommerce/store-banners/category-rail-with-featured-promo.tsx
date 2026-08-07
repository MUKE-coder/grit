'use client'

import { useId, useState } from 'react'
import {
  ChevronLeft,
  ChevronRight,
  Gamepad2,
  Headphones,
  Home,
  Laptop,
  Camera,
  Shirt,
  ShoppingBag,
  Smartphone,
  Sparkles,
  Star,
  Watch,
} from 'lucide-react'

/*
 * Marketplace homepage layout: a department rail down the side, one featured
 * promotion beside it.
 *
 * The promo does not run on a timer. The source advanced every five seconds
 * and treated a click on a dot or an arrow as a ten-second reprieve, so
 * choosing a slide came with a countdown attached. Moving on its own needs a
 * real stop, not a snooze; hero-carousel-with-controls shows what that costs.
 * Here the rail is the navigation and the promo is a poster.
 *
 * One lucide component per department. The source carried a getCategoryIcon
 * map from emoji string to hand-written inline SVG, roughly two hundred lines
 * for twelve glyphs, so adding a department meant writing a path.
 *
 * The rail is a nav of links rather than divs with hover styling. It is the
 * primary navigation of whatever page this sits on.
 *
 * Prices are integer cents and the discount is derived. The source stored
 * "1,599", "1,799" and discount: 11 as three independent strings. 1599 off
 * 1799 is 11.1%, so it had already rounded in a direction nobody chose.
 */

export interface Department {
  label: string
  href: string
  Icon: typeof Smartphone
  featured?: boolean
}

export interface Promo {
  id: string
  eyebrow: string
  title: string
  subtitle: string
  product: string
  specs: string
  /** Integer cents. */
  price: number
  /** Integer cents, struck through. */
  wasPrice: number
  image: string
  badge: string
  /** Tailwind gradient stops for the panel. */
  tone: string
}

const DEPARTMENTS: Department[] = [
  { label: 'Smartphones', href: '#', Icon: Smartphone, featured: true },
  { label: 'Laptops', href: '#', Icon: Laptop },
  { label: 'Watches', href: '#', Icon: Watch },
  { label: 'Audio', href: '#', Icon: Headphones },
  { label: 'Cameras', href: '#', Icon: Camera },
  { label: 'Gaming', href: '#', Icon: Gamepad2 },
  { label: 'Home', href: '#', Icon: Home },
  { label: 'Fashion', href: '#', Icon: Shirt },
  { label: 'Limited editions', href: '#', Icon: Sparkles, featured: true },
]

/*
 * Named by category, not by brand. The source named three specific flagship
 * devices and paired every one with the wrong photograph: the "iPhone 15 Pro
 * Max" slide showed an Android handset, the "Galaxy Z Fold" showed a slab that
 * does not fold, and the "MacBook Pro" showed a phone on a yellow background.
 */
const PROMOS: Promo[] = [
  {
    id: 'phone',
    eyebrow: 'April collection',
    title: 'Exclusive edition',
    subtitle: 'Limited availability',
    product: 'Flagship smartphone',
    specs: '1TB · Titanium · 120Hz',
    price: 159900,
    wasPrice: 179900,
    image: 'https://images.unsplash.com/photo-1592750475338-74b7b21085ab?w=700&h=700&fit=crop&q=80',
    badge: 'Limited',
    tone: 'from-indigo-900 via-purple-800 to-violet-900',
  },
  {
    id: 'laptop',
    eyebrow: 'Featured product',
    title: 'Elite collection',
    subtitle: 'Built for long sessions',
    product: 'Studio laptop',
    specs: '32GB · 2TB · XDR display',
    price: 349900,
    wasPrice: 379900,
    image: 'https://images.unsplash.com/photo-1541807084-5c52b6b3adef?w=700&h=700&fit=crop&q=80',
    badge: 'Elite',
    tone: 'from-emerald-900 via-teal-800 to-emerald-800',
  },
  {
    id: 'audio',
    eyebrow: 'New arrival',
    title: 'Signature series',
    subtitle: 'Member exclusive',
    product: 'Over-ear headphones',
    specs: 'Active noise cancelling · 40h',
    price: 44900,
    wasPrice: 54900,
    image: 'https://images.unsplash.com/photo-1546435770-a3e426bf472b?w=700&h=700&fit=crop&q=80',
    badge: 'Premium',
    tone: 'from-slate-900 via-blue-900 to-slate-800',
  },
]

const money = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
})

export default function CategoryRailWithFeaturedPromo({
  departments = DEPARTMENTS,
  promos = PROMOS,
  railTitle = 'Premium categories',
}: {
  departments?: Department[]
  promos?: Promo[]
  railTitle?: string
}) {
  const [index, setIndex] = useState(0)
  const panelId = useId()
  const promo = promos[index]
  const saved = Math.round((1 - promo.price / promo.wasPrice) * 100)

  const go = (next: number) => setIndex(((next % promos.length) + promos.length) % promos.length)

  return (
    <section className="mx-auto max-w-7xl overflow-hidden rounded-2xl border border-gray-200 shadow-xl dark:border-white/10">
      <div className="flex flex-col md:flex-row">
        {/* A scrolling chip row below md. The source hid the rail outright at
            that breakpoint, which deletes the block's navigation on the screen
            size most people shop from. */}
        <nav
          aria-label={railTitle}
          className="flex w-full shrink-0 flex-col bg-gradient-to-b from-gray-50 to-white md:w-1/4 dark:from-gray-900 dark:to-gray-950"
        >
          <h2 className="border-b border-gray-100 px-6 py-4 font-bold text-gray-800 dark:border-white/10 dark:text-white">
            {railTitle}
          </h2>

          <ul
            role="list"
            className="flex gap-2 overflow-x-auto p-3 md:flex-1 md:flex-col md:gap-0 md:overflow-x-visible md:overflow-y-auto"
          >
            {departments.map(({ label, href, Icon, featured }) => (
              <li key={label} className="shrink-0 md:shrink">
                <a
                  href={href}
                  className="group flex min-h-11 items-center gap-3 rounded-full border border-gray-200 px-3 text-sm whitespace-nowrap text-gray-700 hover:bg-amber-50 md:rounded-lg md:border-0 dark:border-white/10 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-amber-600 dark:text-gray-300 dark:hover:bg-white/5"
                >
                  <span
                    aria-hidden="true"
                    className={`flex size-8 shrink-0 items-center justify-center rounded-lg ${
                      featured
                        ? 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-400'
                        : 'bg-gray-100 text-gray-600 dark:bg-white/5 dark:text-gray-400'
                    }`}
                  >
                    <Icon className="size-4" />
                  </span>
                  <span className="group-hover:underline">{label}</span>
                  {featured && (
                    /* The word, not just the amber tint. Colour alone is not
                       a label. */
                    <span className="ml-auto hidden rounded-full bg-amber-100 px-2 py-0.5 text-[10px] font-semibold tracking-wide text-amber-800 uppercase md:inline dark:bg-amber-500/15 dark:text-amber-300">
                      Featured
                    </span>
                  )}
                </a>
              </li>
            ))}
          </ul>

          <a
            href="#"
            className="group flex min-h-12 items-center justify-between border-t border-gray-100 bg-gray-50 px-6 text-sm font-medium text-gray-600 hover:text-amber-700 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-amber-600 dark:border-white/10 dark:bg-white/5 dark:text-gray-400"
          >
            View all categories
            <ChevronRight aria-hidden="true" className="size-4 transition-transform group-hover:translate-x-1 motion-reduce:transition-none" />
          </a>
        </nav>

        <div className="relative w-full md:w-3/4">
          <div
            id={panelId}
            /* px-16 is the arrows' gutter. At p-6 the previous arrow lands on
               the product name and the first character of the specs. */
            className={`relative flex min-h-[400px] flex-col items-center gap-8 overflow-hidden bg-gradient-to-br px-16 py-6 md:min-h-[550px] md:flex-row md:px-24 md:py-12 ${promo.tone}`}
          >
            <div aria-hidden="true" className="absolute inset-0 opacity-10">
              <span className="absolute top-0 left-0 size-1/3 rounded-full bg-white blur-3xl" />
              <span className="absolute right-0 bottom-0 size-1/2 rounded-full bg-white blur-3xl" />
            </div>

            <p className="absolute top-6 right-6 z-20 inline-flex items-center gap-1 rounded-full bg-amber-500 px-4 py-1 shadow-lg">
              <Star aria-hidden="true" className="size-3 fill-white text-white" />
              <span className="text-xs font-bold tracking-widest text-white uppercase">
                {promo.badge}
              </span>
            </p>

            {/* Announced as a whole. The title, specs and price all change
                together, and reading only the heading would leave someone with
                the previous slide's price. */}
            <div aria-live="polite" aria-atomic="true" className="z-10 w-full md:w-1/2">
              <p className="inline-block rounded-full bg-white/20 px-3 py-1.5 text-xs font-medium text-white backdrop-blur-sm">
                {promo.eyebrow}
              </p>

              <h2 className="mt-3 text-4xl leading-tight font-bold tracking-wide text-balance text-white drop-shadow-md sm:text-5xl md:text-6xl">
                {promo.title}
              </h2>
              <p className="mt-2 text-lg font-light text-white drop-shadow-sm">{promo.subtitle}</p>

              <p className="mt-8 text-3xl font-bold text-white drop-shadow-md md:text-4xl">
                {promo.product}
              </p>
              <p className="mt-2 text-xl font-medium text-white/90 md:text-2xl">{promo.specs}</p>

              <p className="mt-6 inline-flex flex-col rounded-xl border border-white/20 bg-white/10 px-5 py-3 shadow-lg backdrop-blur-sm">
                <span className="text-2xl font-bold text-white md:text-3xl">
                  {money.format(promo.price / 100)}
                </span>
                <span className="mt-1 flex items-center gap-2">
                  <span className="text-sm text-white/80">
                    <span className="sr-only">was </span>
                    <s>{money.format(promo.wasPrice / 100)}</s>
                  </span>
                  <span className="rounded-full bg-amber-500 px-2 py-0.5 text-xs font-bold text-white">
                    Save {saved}%
                  </span>
                </span>
              </p>

              <p className="mt-8">
                <a
                  href="#"
                  className="inline-flex min-h-12 items-center gap-2 rounded-lg bg-white px-8 font-bold text-gray-900 shadow-lg transition hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                >
                  <ShoppingBag aria-hidden="true" className="size-4" />
                  Shop now
                </a>
              </p>

              <p className="mt-6 text-xs text-white/80">
                Premium delivery included. Terms and conditions apply.
              </p>
            </div>

            <div className="z-10 flex w-full justify-center md:w-1/2">
              <div className="relative">
                <span
                  aria-hidden="true"
                  className="absolute -inset-8 rounded-full bg-amber-500/20 blur-3xl"
                />
                {/* Product name and specs are two elements away. */}
                <img
                  src={promo.image}
                  alt=""
                  className="relative max-h-64 rounded-2xl object-contain shadow-2xl sm:max-h-80 md:max-h-96"
                />
              </div>
            </div>

            <button
              type="button"
              onClick={() => go(index - 1)}
              aria-controls={panelId}
              className="absolute top-1/2 left-3 z-20 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-full border border-white/10 bg-black/30 text-white shadow-lg backdrop-blur-md hover:bg-black/50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white md:left-6"
            >
              <ChevronLeft aria-hidden="true" className="size-5" />
              <span className="sr-only">Previous promotion</span>
            </button>
            <button
              type="button"
              onClick={() => go(index + 1)}
              aria-controls={panelId}
              className="absolute top-1/2 right-3 z-20 inline-flex size-12 -translate-y-1/2 items-center justify-center rounded-full border border-white/10 bg-black/30 text-white shadow-lg backdrop-blur-md hover:bg-black/50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white md:right-6"
            >
              <ChevronRight aria-hidden="true" className="size-5" />
              <span className="sr-only">Next promotion</span>
            </button>
          </div>

          <ul
            role="list"
            className="absolute inset-x-0 bottom-2 z-30 mx-auto flex w-fit items-center gap-2 rounded-full bg-black/30 px-3 backdrop-blur-sm"
          >
            {promos.map((item, position) => (
              <li key={item.id}>
                <button
                  type="button"
                  onClick={() => go(position)}
                  aria-controls={panelId}
                  aria-current={position === index ? 'true' : undefined}
                  className="inline-flex h-11 items-center px-0.5 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-white"
                >
                  <span
                    aria-hidden="true"
                    className={`block h-2 rounded-full transition-all ${
                      position === index ? 'w-8 bg-white' : 'w-2 bg-white/50'
                    }`}
                  />
                  <span className="sr-only">{item.product}</span>
                </button>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </section>
  )
}
