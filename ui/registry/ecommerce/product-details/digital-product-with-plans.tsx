'use client'

import { useId, useState } from 'react'
import { Star } from 'lucide-react'

/*
 * A detail page for something you download: screenshots and prose on the left,
 * a sticky package picker on the right.
 *
 * The breadcrumb is an ordered list with `aria-current="page"` on the last
 * item, and its separators are `aria-hidden`. In the source they were bare "/"
 * characters in the text, which a screen reader reads out: "Home slash
 * Products slash Clarity UI". Separators are punctuation drawn for the eye;
 * the list structure is what carries the meaning.
 *
 * The package picker is a fieldset with a legend rather than a heading above
 * loose radios. Radios need a group name announced before the options, or the
 * first one is read as though it were the only question on the page.
 *
 * The rating states itself in text. Five filled stars and a review count told
 * you how many people reviewed it and never what they said.
 *
 * The price lives in one field. The source carried both `name: "$19"` and
 * `price: 19` and rendered them in different places, which is two sources of
 * truth for the same number and a currency symbol hardcoded in the data. There
 * is one number, formatted once.
 *
 * `lg:self-start` is what makes the sidebar stick. A grid item stretches to the
 * row height by default, leaving nothing to stick within, and `position:
 * sticky` then fails silently.
 */

export interface Plan {
  id: string
  price: number
  /** What the price buys, e.g. "6 months of updates". */
  period: string
  isPopular?: boolean
}

export interface DigitalProduct {
  name: string
  rating: number
  reviewCount: number
  updatedOn: string
  description: string
  features: string[]
  highlights: string[]
  images: { src: string; caption: string }[]
  plans: Plan[]
}

const CURRENCY = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD',
  maximumFractionDigits: 0,
})

const PRODUCT: DigitalProduct = {
  name: 'Clarity landing UI kit',
  rating: 4.8,
  reviewCount: 39,
  updatedOn: '18 April 2026',
  description:
    'Ninety coded blocks for landing pages, built with Tailwind and nothing else. No component library to learn, no runtime to keep in step: you copy the block you want and own the markup from that moment on.',
  features: [
    '90+ coded blocks',
    'Built with Tailwind CSS, no other dependencies',
    'Responsive from 320px up',
    'Light and dark themes throughout',
  ],
  highlights: [
    'Personal and commercial licence',
    'Free updates for the term you buy',
    'Figma source files included',
  ],
  images: [
    {
      src: 'https://images.unsplash.com/photo-1559028012-481c04fa702d?w=1000&h=563&fit=crop&q=80',
      caption: 'Landing page blocks in the design file',
    },
    {
      src: 'https://images.unsplash.com/photo-1587440871875-191322ee64b0?w=1000&h=563&fit=crop&q=80',
      caption: 'Wireframes for the layouts the kit covers',
    },
  ],
  plans: [
    { id: 'monthly', price: 19, period: '1 month of updates' },
    { id: 'biannual', price: 29, period: '6 months of updates', isPopular: true },
    { id: 'annual', price: 49, period: '12 months of updates' },
  ],
}

const BODY = [
  'Every block is a single file with no imports beyond the icons it draws. Install one and it renders; there is no provider to mount, no theme to merge and no configuration that has to match ours.',
  'The layouts cover the sections a landing page actually needs: heroes, feature grids, pricing tables, FAQs and the closing ask. Each one ships in light and dark, and each one has been driven with a keyboard before it went in.',
]

export default function DigitalProductWithPlans({
  product = PRODUCT,
  body = BODY,
  onPurchase,
}: {
  product?: DigitalProduct
  body?: string[]
  onPurchase?: (plan: Plan) => void
}) {
  const [selectedId, setSelectedId] = useState(
    product.plans.find((p) => p.isPopular)?.id ?? product.plans[0]?.id,
  )
  const groupName = useId()
  const selected = product.plans.find((p) => p.id === selectedId) ?? product.plans[0]

  return (
    <section className="bg-white py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Separators are drawn, not spoken. The list carries the structure. */}
        <nav aria-label="Breadcrumb" className="mb-8">
          <ol role="list" className="flex items-center text-sm text-gray-500 dark:text-gray-400">
            <li>
              <a
                href="#"
                className="hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:hover:text-white"
              >
                Home
              </a>
            </li>
            <li aria-hidden="true" className="mx-2">
              /
            </li>
            <li>
              <a
                href="#"
                className="hover:text-gray-900 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:hover:text-white"
              >
                Products
              </a>
            </li>
            <li aria-hidden="true" className="mx-2">
              /
            </li>
            <li aria-current="page" className="text-gray-900 dark:text-white">
              {product.name}
            </li>
          </ol>
        </nav>

        <div className="grid grid-cols-1 gap-12 lg:grid-cols-12">
          <div className="lg:col-span-8">
            <ul role="list" className="mb-10 grid grid-cols-1 gap-4 md:grid-cols-2">
              {product.images.map((image) => (
                <li
                  key={image.src}
                  className="overflow-hidden rounded-lg border border-gray-200 dark:border-white/10"
                >
                  <img
                    src={image.src}
                    alt={image.caption}
                    className="aspect-video w-full object-cover"
                  />
                </li>
              ))}
            </ul>

            <p className="flex items-center gap-2">
              <span aria-hidden="true" className="flex">
                {[0, 1, 2, 3, 4].map((i) => (
                  <Star
                    key={i}
                    className={`size-5 ${
                      i < Math.round(product.rating)
                        ? 'fill-amber-400 text-amber-400'
                        : 'fill-gray-200 text-gray-200 dark:fill-white/15 dark:text-white/15'
                    }`}
                  />
                ))}
              </span>
              <span className="text-gray-600 dark:text-gray-400">
                <span className="sr-only">Rated </span>
                {product.rating.toFixed(1)}
                <span className="sr-only"> out of 5</span> from {product.reviewCount} reviews
              </span>
            </p>

            <h2 className="mt-3 text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
              {product.name}
            </h2>
            <p className="mt-1 text-gray-500 dark:text-gray-400">
              Updated on {product.updatedOn}
            </p>

            <p className="mt-6 text-base/7 text-gray-700 dark:text-gray-300">
              {product.description}
            </p>

            <h3 className="mt-10 text-xl font-bold text-gray-900 dark:text-white">
              What is included
            </h3>
            <ul
              role="list"
              className="mt-4 list-disc space-y-2 pl-5 text-gray-700 dark:text-gray-300"
            >
              {product.features.map((feature) => (
                <li key={feature}>{feature}</li>
              ))}
            </ul>

            <h3 className="mt-10 text-xl font-bold text-gray-900 dark:text-white">
              Built to be pulled apart
            </h3>
            {body.map((paragraph) => (
              <p key={paragraph} className="mt-4 text-base/7 text-gray-700 dark:text-gray-300">
                {paragraph}
              </p>
            ))}
          </div>

          {/* self-start is load-bearing: a stretched grid item has no room to
              stick inside, and sticky then fails silently. */}
          <div className="lg:col-span-4 lg:self-start">
            <div className="space-y-8 lg:sticky lg:top-8">
              <div className="rounded-lg border border-gray-200 p-6 dark:border-white/10">
                <h3 className="text-lg font-bold text-gray-900 dark:text-white">Highlights</h3>
                <ul role="list" className="mt-4 space-y-3">
                  {product.highlights.map((highlight) => (
                    <li key={highlight} className="flex items-start gap-2">
                      <span
                        aria-hidden="true"
                        className="mt-2 size-1.5 flex-none rounded-full bg-gray-400 dark:bg-gray-600"
                      />
                      <span className="text-gray-700 dark:text-gray-300">{highlight}</span>
                    </li>
                  ))}
                </ul>
              </div>

              <div className="overflow-hidden rounded-lg border border-gray-200 dark:border-white/10">
                <fieldset>
                  <legend className="w-full border-b border-gray-200 p-6 text-lg font-bold text-gray-900 dark:border-white/10 dark:text-white">
                    Select a package
                  </legend>

                  <div className="divide-y divide-gray-200 dark:divide-white/10">
                    {product.plans.map((plan) => {
                      const active = plan.id === selected.id
                      return (
                        <label
                          key={plan.id}
                          className={`flex cursor-pointer items-center gap-4 p-6 hover:bg-gray-50 has-[:focus-visible]:outline-2 has-[:focus-visible]:-outline-offset-2 has-[:focus-visible]:outline-indigo-600 dark:hover:bg-white/5 ${
                            active ? 'bg-gray-50 dark:bg-white/5' : ''
                          }`}
                        >
                          <input
                            type="radio"
                            name={groupName}
                            value={plan.id}
                            checked={active}
                            onChange={() => setSelectedId(plan.id)}
                            className="size-4 accent-indigo-600"
                          />
                          <span className="flex w-full items-center justify-between gap-3">
                            <span>
                              <span className="block text-lg font-bold text-gray-900 tabular-nums dark:text-white">
                                {CURRENCY.format(plan.price)}
                              </span>
                              <span className="block text-sm text-gray-600 dark:text-gray-400">
                                {plan.period}
                              </span>
                            </span>
                            {plan.isPopular && (
                              <span className="rounded-full bg-indigo-100 px-2.5 py-0.5 text-xs font-medium text-indigo-800 dark:bg-indigo-500/15 dark:text-indigo-300">
                                Popular
                              </span>
                            )}
                          </span>
                        </label>
                      )
                    })}
                  </div>
                </fieldset>

                <div className="space-y-4 p-6">
                  <button
                    type="button"
                    onClick={() => onPurchase?.(selected)}
                    className="inline-flex min-h-11 w-full items-center justify-center rounded-md bg-indigo-600 px-4 text-sm font-medium text-white hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                  >
                    Purchase for {CURRENCY.format(selected.price)}
                    <span className="sr-only">, {selected.period}</span>
                  </button>
                  <a
                    href="#"
                    className="inline-flex min-h-11 w-full items-center justify-center rounded-md border border-gray-300 px-4 text-sm font-medium text-gray-900 hover:bg-gray-50 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:border-white/15 dark:text-white dark:hover:bg-white/5"
                  >
                    Live preview
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}
