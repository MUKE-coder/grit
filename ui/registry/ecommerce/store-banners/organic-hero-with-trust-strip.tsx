import { ArrowRight, ChevronRight, Leaf, Scale, Star, ThumbsUp, Truck } from 'lucide-react'

/*
 * Product hero with a floating price tag, two category cards, and a strip of
 * trust claims along the foot.
 *
 * No mount animation. The source set opacity-0 translate-y-10 and cleared it
 * from a useEffect, so the headline, the price and both buttons were invisible
 * until React hydrated: a fade-in on a fast laptop, a blank hero on a slow
 * phone or a crawler. It also registered a scroll listener with an empty body.
 *
 * The headline is real text rather than bg-clip-text text-transparent.
 * Gradient text disappears under Windows High Contrast and takes the selection
 * highlight with it, and the gradient here ran gray-900 to gray-700.
 *
 * Each category card is one link. The source nested an "Explore" button inside
 * the card, so a keyboard user got two stops for one destination.
 *
 * The trust strip uses lucide icons rather than emoji. Screen readers read an
 * emoji's name, so a truck beside "short supply chain" is the truck twice.
 */

export interface TrustPoint {
  label: string
  Icon: typeof Truck
  tone: string
}

export interface CategoryCard {
  name: string
  href: string
  image: string
  tone: string
}

const TRUST: TrustPoint[] = [
  { label: 'Short supply chain', Icon: Truck, tone: 'bg-green-100 text-green-700' },
  { label: 'Fair prices', Icon: Scale, tone: 'bg-gray-100 text-gray-700' },
  { label: 'High quality', Icon: ThumbsUp, tone: 'bg-amber-100 text-amber-700' },
  { label: 'Full traceability', Icon: Leaf, tone: 'bg-green-100 text-green-700' },
]

const CATEGORIES: CategoryCard[] = [
  {
    name: 'Honey and preserves',
    href: '#',
    image: 'https://images.unsplash.com/photo-1471943311424-646960669fbc?w=400&h=400&fit=crop&q=80',
    tone: 'from-amber-600 to-amber-700',
  },
  {
    name: 'Snack bars',
    href: '#',
    image: 'https://images.unsplash.com/photo-1622484212850-eb596d769edc?w=400&h=400&fit=crop&q=80',
    tone: 'from-orange-700 to-orange-800',
  },
]

export default function OrganicHeroWithTrustStrip({
  eyebrow = 'Premium organic pantry',
  title = 'Organic comes knocking',
  description = 'Nuts bought direct from the growers, shipped the week they were packed.',
  fromPrice = 4900,
  image = 'https://images.unsplash.com/photo-1608797178974-15b35a64ede9?w=900&h=900&fit=crop&q=80',
  categories = CATEGORIES,
  trust = TRUST,
}: {
  eyebrow?: string
  title?: string
  description?: string
  /** Integer cents. */
  fromPrice?: number
  image?: string
  categories?: CategoryCard[]
  trust?: TrustPoint[]
}) {
  const price = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
  }).format(fromPrice / 100)

  return (
    <section className="overflow-hidden rounded-2xl bg-gradient-to-br from-sky-100 to-sky-200 shadow-xl">
      <div className="relative mx-auto max-w-7xl">
        <p className="absolute top-4 right-4 z-20 inline-flex items-center gap-1.5 rounded-full bg-gradient-to-r from-amber-600 to-yellow-500 px-4 py-2 shadow-lg md:top-6 md:right-6">
          <Star aria-hidden="true" className="size-4 fill-white text-white" />
          <span className="text-sm font-semibold text-white">Premium</span>
        </p>

        <div className="grid grid-cols-1 gap-8 px-4 py-12 sm:px-8 md:grid-cols-12 md:py-16">
          <div className="flex flex-col justify-center md:col-span-5">
            <p className="inline-flex w-fit items-center gap-2 rounded-full bg-sky-700/10 px-4 py-1.5 text-sm font-medium text-sky-900">
              {/* Decorative, and it stops under reduced motion. */}
              <span aria-hidden="true" className="relative flex size-2.5">
                <span className="absolute inline-flex size-full animate-ping rounded-full bg-sky-400 opacity-75 motion-reduce:animate-none" />
                <span className="relative inline-flex size-2.5 rounded-full bg-sky-600" />
              </span>
              {eyebrow}
            </p>

            <h2 className="mt-6 text-4xl leading-tight font-bold tracking-tight text-balance text-gray-900 sm:text-5xl md:text-6xl">
              {title}
            </h2>

            <p className="mt-6 max-w-lg text-base leading-relaxed text-pretty text-gray-700 sm:text-lg">
              {description}
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <a
                href="#"
                className="inline-flex min-h-12 items-center gap-2 rounded-md bg-gradient-to-r from-amber-700 to-yellow-500 px-8 text-sm font-medium tracking-wider text-white uppercase shadow-lg transition hover:from-amber-800 hover:to-yellow-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-amber-700"
              >
                Explore collection
                <ArrowRight aria-hidden="true" className="size-4" />
              </a>
              <a
                href="#"
                className="inline-flex min-h-12 items-center rounded-md border border-gray-400 px-8 text-sm font-medium tracking-wider text-gray-700 uppercase transition hover:bg-white/60 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-700"
              >
                About our process
              </a>
            </div>
          </div>

          {/* Clears the Premium badge, which is pinned to the same corner
              this column occupies. Without it the badge sits on the first
              category card. */}
          <div className="md:col-span-7 md:pt-14">
            <div className="flex flex-col items-center gap-6 sm:flex-row sm:items-start sm:justify-center">
              <div className="relative">
                <div
                  aria-hidden="true"
                  className="absolute -inset-4 rounded-full bg-gradient-to-br from-white/30 to-sky-200/40 blur-2xl"
                />
                {/* The heading two elements up already names the product. */}
                <img
                  src={image}
                  alt=""
                  className="relative w-full max-w-sm -rotate-2 rounded-lg object-contain shadow-2xl transition-transform duration-500 hover:rotate-0 motion-reduce:transition-none motion-reduce:hover:rotate-[-2deg]"
                />

                <p className="absolute right-2 -bottom-4 flex size-16 flex-col items-center justify-center rounded-full bg-gradient-to-r from-amber-600 to-yellow-500 text-white ring-4 ring-white">
                  <span className="text-xs font-light">from</span>
                  <span className="text-lg font-bold">{price}</span>
                </p>
              </div>

              <ul role="list" className="flex shrink-0 gap-4 sm:flex-col">
                {categories.map((category) => (
                  <li key={category.name}>
                    <a
                      href={category.href}
                      className={`group block w-32 overflow-hidden rounded-xl bg-gradient-to-br shadow-xl transition-transform duration-300 hover:scale-105 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 motion-reduce:transition-none motion-reduce:hover:scale-100 lg:w-40 ${category.tone}`}
                    >
                      <span className="relative block h-32 lg:h-40">
                        <img
                          src={category.image}
                          alt=""
                          className="size-full object-cover transition-transform duration-700 group-hover:scale-110 motion-reduce:transition-none motion-reduce:group-hover:scale-100"
                        />
                        {/* Bottom-weighted. The source tinted the whole photo
                            at 60% in the card's own colour, which left both
                            products as brown shapes. */}
                        <span
                          aria-hidden="true"
                          className="absolute inset-0 bg-gradient-to-t from-black/60 via-black/10 to-transparent"
                        />
                      </span>
                      <span className="block p-3 text-white">
                        <span className="block text-sm font-medium">{category.name}</span>
                        <span className="mt-1 flex items-center text-sm font-light">
                          Explore
                          <ChevronRight
                            aria-hidden="true"
                            className="ml-1 size-3 transition-all duration-300 group-hover:ml-2 motion-reduce:transition-none"
                          />
                        </span>
                      </span>
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        <ul
          role="list"
          /* Grid below sm. wrap + justify-between pushes the last row's items
             to both edges, and on a narrow screen the right-hand label ends up
             flush against the padding. */
          className="grid grid-cols-2 gap-4 bg-white/80 px-4 py-6 backdrop-blur-md sm:flex sm:justify-between sm:gap-6 sm:px-8"
        >
          {trust.map(({ label, Icon, tone }) => (
            <li key={label} className="flex items-center gap-3">
              <span
                aria-hidden="true"
                className={`flex size-10 items-center justify-center rounded-full shadow-sm ${tone}`}
              >
                <Icon className="size-5" />
              </span>
              <span className="text-xs font-medium text-gray-800 sm:text-sm">{label}</span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
