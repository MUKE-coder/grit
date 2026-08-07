'use client'

import { useId, useState } from 'react'
import { Check, RefreshCw, ShoppingBag, Star, Truck } from 'lucide-react'

/*
 * A physical product page: gallery, colour, size, and the two facts that
 * decide the purchase — when it arrives and whether you can send it back.
 *
 * Colour and size are two fieldsets of radios, not rows of buttons. Both are
 * single-choice, which is what radios are: arrow keys move between the
 * options, the group announces its legend before them, and the chosen one is
 * announced as checked. The source used buttons whose only accessible name was
 * a `title` attribute — which many screen readers do not announce at all, and
 * which no touch user ever sees.
 *
 * Out-of-stock options stay in the document, disabled, with the reason in
 * their label. Removing them is worse: the reader cannot tell whether size 43
 * is sold out or was never made, and the gap in the sequence is invisible to
 * anyone not looking at it.
 *
 * The gallery thumbnails carry `aria-current`, so the active one is announced
 * rather than only being ringed. Only the main image describes the product;
 * the thumbnails are labelled by position, because five copies of the same
 * sentence is not five pieces of information.
 *
 * The price fields are named for what they are. The source had `price: 49`
 * with `discountPrice: 59.99` — the "discount" being the higher number, which
 * is the sort of thing that survives review and then prints wrong.
 */

export interface Variant {
  name: string
  /** CSS colour for the swatch. */
  value: string
  inStock: boolean
}

export interface SizeOption {
  value: string
  inStock: boolean
}

const CURRENCY = new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' })

const IMAGES = [
  'https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=800&h=800&fit=crop&q=80',
  'https://images.unsplash.com/photo-1607522370275-f14206abe5d3?w=800&h=800&fit=crop&q=80',
  'https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?w=800&h=800&fit=crop&q=80',
  'https://images.unsplash.com/photo-1608231387042-66d1773070a5?w=800&h=800&fit=crop&q=80',
  'https://images.unsplash.com/photo-1605348532760-6753d2c43329?w=800&h=800&fit=crop&q=80',
]

const COLOURS: Variant[] = [
  { name: 'Bone', value: '#F5F1EA', inStock: true },
  { name: 'Black', value: '#111111', inStock: true },
  { name: 'Ember', value: '#E4622A', inStock: true },
  { name: 'Cobalt', value: '#1E40AF', inStock: false },
]

const SIZES: SizeOption[] = [
  { value: '38', inStock: true },
  { value: '39', inStock: true },
  { value: '40', inStock: true },
  { value: '41', inStock: true },
  { value: '42', inStock: true },
  { value: '43', inStock: false },
  { value: '44', inStock: false },
  { value: '45', inStock: true },
]

const HIGHLIGHTS = [
  'Cushioned insole that holds its shape past a year of wear',
  'Mesh upper, so they dry out overnight rather than over a weekend',
  'Rubber outsole with real grip on wet pavement',
  'Reinforced heel counter for runners who overpronate',
]

export default function PhysicalProductWithVariants({
  name = 'Trail running shoe, mens',
  category = 'Running',
  price = 49,
  wasPrice = 59.99,
  rating = 4.8,
  reviewCount = 246,
  description = 'Built for pavement and light trail. The upper is a single piece of mesh with no seams to rub, and the midsole is firm enough to stay honest at the end of a long run.',
  highlights = HIGHLIGHTS,
  images = IMAGES,
  colours = COLOURS,
  sizes = SIZES,
  deliveryEstimate = '2 to 4 working days',
  returnsPolicy = 'Free returns within 30 days',
  onAddToBasket,
}: {
  name?: string
  category?: string
  price?: number
  /** The struck-through original, if there is one. */
  wasPrice?: number
  rating?: number
  reviewCount?: number
  description?: string
  highlights?: string[]
  images?: string[]
  colours?: Variant[]
  sizes?: SizeOption[]
  deliveryEstimate?: string
  returnsPolicy?: string
  onAddToBasket?: (selection: { colour: string; size: string }) => void
}) {
  const [active, setActive] = useState(0)
  const [colour, setColour] = useState(() => colours.find((c) => c.inStock)?.name ?? '')
  const [size, setSize] = useState(() => sizes.find((s) => s.inStock)?.value ?? '')
  const colourGroup = useId()
  const sizeGroup = useId()

  return (
    <section className="bg-white py-8 dark:bg-gray-950">
      <div className="mx-auto max-w-6xl p-4 md:p-8">
        <div className="flex flex-col gap-8 md:flex-row">
          <div className="md:w-1/2">
            <div className="aspect-square overflow-hidden rounded-xl bg-orange-50 dark:bg-white/5">
              <img
                src={images[active]}
                alt={`${name}, image ${active + 1} of ${images.length}`}
                className="size-full object-contain p-4"
              />
            </div>

            <ul role="list" className="mt-4 flex gap-2 overflow-x-auto pb-2">
              {images.map((image, i) => (
                <li key={image}>
                  <button
                    type="button"
                    onClick={() => setActive(i)}
                    aria-current={i === active ? 'true' : undefined}
                    className={`size-16 flex-none overflow-hidden rounded-md focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 ${
                      i === active
                        ? 'ring-2 ring-gray-900 dark:ring-white'
                        : 'ring-1 ring-gray-200 dark:ring-white/15'
                    }`}
                  >
                    <img src={image} alt="" className="size-full object-cover" />
                    <span className="sr-only">Show image {i + 1}</span>
                  </button>
                </li>
              ))}
            </ul>
          </div>

          <div className="md:w-1/2">
            <nav aria-label="Breadcrumb" className="mb-4">
              <ol role="list" className="flex items-center text-sm text-gray-500 dark:text-gray-400">
                <li>
                  <a href="#" className="hover:text-gray-700 dark:hover:text-gray-200">
                    Home
                  </a>
                </li>
                <li aria-hidden="true" className="mx-2">
                  /
                </li>
                <li>
                  <a href="#" className="hover:text-gray-700 dark:hover:text-gray-200">
                    {category}
                  </a>
                </li>
                <li aria-hidden="true" className="mx-2">
                  /
                </li>
                <li aria-current="page" className="text-gray-900 dark:text-white">
                  {name}
                </li>
              </ol>
            </nav>

            <h2 className="text-3xl font-bold tracking-tight text-balance text-gray-900 dark:text-white">
              {name}
            </h2>

            <p className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-1">
              <span className="flex items-center gap-1.5">
                <span aria-hidden="true" className="flex">
                  {[0, 1, 2, 3, 4].map((i) => (
                    <Star
                      key={i}
                      className={`size-4 ${
                        i < Math.round(rating)
                          ? 'fill-amber-400 text-amber-400'
                          : 'fill-gray-200 text-gray-200 dark:fill-white/15 dark:text-white/15'
                      }`}
                    />
                  ))}
                </span>
                <span className="text-sm text-gray-600 dark:text-gray-400">
                  <span className="sr-only">Rated </span>
                  {rating.toFixed(1)}
                  <span className="sr-only"> out of 5</span>
                </span>
              </span>
              <span className="text-sm text-gray-500 dark:text-gray-400">
                {reviewCount} reviews
              </span>
            </p>

            <p className="mt-4 flex items-baseline gap-3">
              <span className="text-3xl font-bold text-gray-900 tabular-nums dark:text-white">
                {CURRENCY.format(price)}
              </span>
              {wasPrice > price && (
                <span className="text-lg text-gray-400 line-through tabular-nums dark:text-gray-500">
                  <span className="sr-only">Was </span>
                  {CURRENCY.format(wasPrice)}
                </span>
              )}
            </p>

            <p className="mt-4 text-base/7 text-gray-600 dark:text-gray-400">{description}</p>

            <h3 className="mt-6 font-medium text-gray-900 dark:text-white">Highlights</h3>
            <ul role="list" className="mt-2 space-y-1">
              {highlights.map((highlight) => (
                <li key={highlight} className="flex items-start gap-2">
                  <Check
                    aria-hidden="true"
                    className="mt-1 size-4 flex-none text-emerald-600 dark:text-emerald-400"
                  />
                  <span className="text-gray-600 dark:text-gray-400">{highlight}</span>
                </li>
              ))}
            </ul>

            {/* A single choice, so radios rather than buttons: arrow keys work
                and the group announces its legend. */}
            <fieldset className="mt-6">
              <div className="flex items-center justify-between">
                <legend className="font-medium text-gray-900 dark:text-white">Colour</legend>
                <span className="text-sm text-gray-500 dark:text-gray-400">{colour}</span>
              </div>
              <div className="mt-2 flex gap-3">
                {colours.map((option) => (
                  <label
                    key={option.name}
                    className={`relative flex size-11 items-center justify-center rounded-full ring-offset-2 has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-indigo-600 ${
                      colour === option.name ? 'ring-2 ring-gray-900 dark:ring-white' : ''
                    } ${option.inStock ? 'cursor-pointer' : 'cursor-not-allowed opacity-40'}`}
                  >
                    <input
                      type="radio"
                      name={colourGroup}
                      className="sr-only"
                      checked={colour === option.name}
                      disabled={!option.inStock}
                      onChange={() => setColour(option.name)}
                    />
                    <span
                      aria-hidden="true"
                      className="absolute inset-1 rounded-full border border-gray-300 dark:border-white/20"
                      style={{ backgroundColor: option.value }}
                    />
                    {/* The reason it cannot be chosen, in the label itself. */}
                    <span className="sr-only">
                      {option.name}
                      {option.inStock ? '' : ', out of stock'}
                    </span>
                  </label>
                ))}
              </div>
            </fieldset>

            <fieldset className="mt-6">
              <div className="flex items-center justify-between">
                <legend className="font-medium text-gray-900 dark:text-white">Size</legend>
                <a
                  href="#"
                  className="text-sm text-indigo-600 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
                >
                  Size guide
                </a>
              </div>
              <div className="mt-2 grid grid-cols-4 gap-2">
                {sizes.map((option) => (
                  <label
                    key={option.value}
                    className={`flex min-h-11 items-center justify-center rounded-md border text-sm has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-indigo-600 ${
                      size === option.value
                        ? 'border-gray-900 bg-gray-900 text-white dark:border-white dark:bg-white dark:text-gray-900'
                        : 'border-gray-300 text-gray-900 dark:border-white/15 dark:text-white'
                    } ${
                      option.inStock
                        ? 'cursor-pointer'
                        : 'cursor-not-allowed text-gray-400 line-through dark:text-gray-600'
                    }`}
                  >
                    <input
                      type="radio"
                      name={sizeGroup}
                      className="sr-only"
                      checked={size === option.value}
                      disabled={!option.inStock}
                      onChange={() => setSize(option.value)}
                    />
                    {option.value}
                    <span className="sr-only">{option.inStock ? '' : ', out of stock'}</span>
                  </label>
                ))}
              </div>
            </fieldset>

            <button
              type="button"
              onClick={() => onAddToBasket?.({ colour, size })}
              className="mt-8 inline-flex min-h-11 w-full items-center justify-center gap-2 rounded-lg bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
            >
              <ShoppingBag aria-hidden="true" className="size-4" />
              Add to basket
              <span className="sr-only">
                : {name}, {colour}, size {size}
              </span>
            </button>

            <dl className="mt-6 space-y-3 border-t border-gray-200 pt-6 text-sm dark:border-white/10">
              <div className="flex gap-3">
                <dt className="flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                  <Truck aria-hidden="true" className="size-4" />
                  Delivery
                </dt>
                <dd className="text-gray-600 dark:text-gray-400">
                  Free, arriving in {deliveryEstimate}
                </dd>
              </div>
              <div className="flex gap-3">
                <dt className="flex items-center gap-2 font-medium text-gray-900 dark:text-white">
                  <RefreshCw aria-hidden="true" className="size-4" />
                  Returns
                </dt>
                <dd className="text-gray-600 dark:text-gray-400">{returnsPolicy}</dd>
              </div>
            </dl>
          </div>
        </div>
      </div>
    </section>
  )
}
