'use client'

import { useId, useState } from 'react'
import { Heart, ShoppingCart } from 'lucide-react'

/*
 * A single-product banner: the pitch, a variant picker, and the numbers that
 * make the case.
 *
 * The variant picker actually changes the product. That sounds like a low bar
 * and it is the thing this pattern most often gets wrong — the source of this
 * block had a `getImagePath()` that returned the same photograph for every
 * colour, with a comment admitting it. A swatch that recolours the page
 * background but not the chair teaches people the picker is decorative, and
 * they stop using it.
 *
 * It is a fieldset of radios rather than a row of buttons. This is a
 * single-choice selection, which is exactly what radios are: arrow keys move
 * between them, the group is announced with its legend, and the chosen one is
 * announced as checked. A row of buttons with aria-pressed says "three
 * toggles, one of which happens to be on", which is a different thing.
 *
 * The inputs are `sr-only`, not hidden. Hidden removes them from the tab order
 * and the accessibility tree, taking the keyboard support with them. The
 * visible swatch is the <label>, so clicking works with no handler.
 *
 * The tint is applied as an rgb() with an alpha channel rather than by
 * appending "30" to a hex string. The hex trick assumes every colour you will
 * ever pass is six digits, and silently produces nothing when someone passes a
 * named colour or a three-digit hex.
 */

export interface Variant {
  id: string
  name: string
  /** The swatch colour, and the tint behind the banner. */
  rgb: [number, number, number]
  image: string
  /** Shown on the annotation pin over the photo. */
  material: string
}

const VARIANTS: Variant[] = [
  {
    id: 'slate',
    name: 'Slate',
    rgb: [148, 180, 184],
    image: 'https://images.unsplash.com/photo-1592078615290-033ee584e267?w=900&h=900&fit=crop&q=80',
    material: 'Moulded shell, oak legs',
  },
  {
    id: 'cream',
    name: 'Cream',
    rgb: [222, 214, 200],
    image: 'https://images.unsplash.com/photo-1567538096630-e0c55bd6374c?w=900&h=900&fit=crop&q=80',
    material: 'Buttoned velvet, turned legs',
  },
  {
    id: 'ochre',
    name: 'Ochre',
    rgb: [226, 178, 84],
    image: 'https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=900&h=900&fit=crop&q=80',
    material: 'Woven wool, tapered legs',
  },
]

const STATS = [
  { value: '2.5k+', label: 'Chairs sold' },
  { value: '98%', label: 'Satisfaction' },
  { value: '1.2k', label: 'Five-star reviews' },
  { value: '10 yr', label: 'Warranty' },
]

export default function ProductSpotlightWithVariants({
  name = 'Lento lounge chair',
  description = 'Ergonomic curves with enough give to sink into, and a frame that holds its shape after a decade of sitting in it.',
  price = 'USD $798',
  variants = VARIANTS,
  stats = STATS,
  onAdd,
}: {
  name?: string
  description?: string
  price?: string
  variants?: Variant[]
  stats?: { value: string; label: string }[]
  onAdd?: (variant: Variant) => void
}) {
  const [selectedId, setSelectedId] = useState(variants[0]?.id)
  const groupName = useId()
  const selected = variants.find((v) => v.id === selectedId) ?? variants[0]

  /* rgb() with an alpha channel. Appending "30" to a hex string assumes every
     colour is six digits and fails silently when one is not. */
  const tint = (alpha: number) =>
    `rgb(${selected.rgb[0]} ${selected.rgb[1]} ${selected.rgb[2]} / ${alpha})`

  return (
    <section
      className="w-full transition-colors duration-300 dark:bg-gray-950"
      style={{ backgroundColor: tint(0.18) }}
    >
      <div className="mx-auto max-w-7xl">
        <div className="grid grid-cols-1 gap-8 px-6 py-12 md:px-10 lg:grid-cols-2">
          <div className="flex flex-col justify-center">
            <h2 className="text-4xl font-bold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
              {name}
            </h2>
            <p className="mt-6 max-w-lg text-lg/8 text-pretty text-gray-700 dark:text-gray-300">
              {description}
            </p>

            <div className="mt-8 flex flex-wrap gap-3">
              <button
                type="button"
                onClick={() => onAdd?.(selected)}
                className="inline-flex min-h-11 items-center gap-2 rounded-md bg-gray-900 px-6 text-sm font-medium text-white hover:bg-gray-800 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-100"
              >
                <ShoppingCart aria-hidden="true" className="size-5" />
                Add to cart
                <span className="sr-only">
                  : {name}, {selected.name}
                </span>
              </button>
              <a
                href="#"
                className="inline-flex min-h-11 items-center rounded-md border border-gray-300 px-6 text-sm font-medium text-gray-900 hover:bg-white/60 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:border-white/20 dark:text-white dark:hover:bg-white/10"
              >
                Learn more
              </a>
              <button
                type="button"
                className="inline-flex size-11 items-center justify-center rounded-md border border-gray-300 text-gray-900 hover:bg-white/60 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900 dark:border-white/20 dark:text-white dark:hover:bg-white/10"
              >
                <Heart aria-hidden="true" className="size-5" />
                <span className="sr-only">Save {name} to wishlist</span>
              </button>
            </div>

            <fieldset className="mt-8">
              <legend className="sr-only">Choose a finish</legend>
              <div className="flex flex-wrap items-center gap-4">
                {variants.map((variant) => {
                  const active = variant.id === selected.id
                  return (
                    <label
                      key={variant.id}
                      className={`flex size-11 cursor-pointer items-center justify-center rounded-full ring-offset-2 transition-shadow has-[:focus-visible]:outline-2 has-[:focus-visible]:outline-offset-2 has-[:focus-visible]:outline-gray-900 ${
                        active ? 'ring-2 ring-gray-500' : ''
                      }`}
                      style={{
                        backgroundColor: `rgb(${variant.rgb[0]} ${variant.rgb[1]} ${variant.rgb[2]})`,
                      }}
                    >
                      {/* sr-only, not hidden: hidden takes the keyboard with it. */}
                      <input
                        type="radio"
                        name={groupName}
                        className="sr-only"
                        checked={active}
                        onChange={() => setSelectedId(variant.id)}
                      />
                      {/* The swatch shows the colour; the name is for anyone
                          not seeing it. Colour alone is never a label. */}
                      <span className="sr-only">{variant.name} finish</span>
                    </label>
                  )
                })}
                <span className="text-sm text-gray-700 dark:text-gray-300">
                  {selected.name}
                </span>
              </div>
            </fieldset>
          </div>

          <div className="relative">
            <div className="relative aspect-square w-full overflow-hidden rounded-lg">
              {/* The picture changes with the variant. A swatch that only
                  recolours the background teaches people to ignore it. */}
              <img
                key={selected.id}
                src={selected.image}
                alt={`${name} in ${selected.name}`}
                className="size-full object-cover"
              />
            </div>

            <p className="absolute top-1/3 right-4 flex items-center gap-2">
              <span
                aria-hidden="true"
                className="flex size-4 items-center justify-center rounded-full border-2 border-gray-300 bg-white"
              >
                <span className="size-1.5 rounded-full bg-gray-500" />
              </span>
              <span className="rounded-md border border-gray-200 bg-white px-3 py-1 text-sm text-gray-900 shadow-md">
                {selected.material}
              </span>
            </p>
          </div>
        </div>

        <dl className="grid grid-cols-2 border-t border-gray-900/10 md:grid-cols-4 dark:border-white/10">
          {stats.map((stat, i) => (
            <div
              key={stat.label}
              className={`px-6 py-8 text-center ${
                i < stats.length - 1
                  ? 'border-r border-gray-900/10 dark:border-white/10'
                  : ''
              }`}
            >
              <dd className="text-2xl font-bold text-gray-900 tabular-nums md:text-3xl dark:text-white">
                {stat.value}
              </dd>
              <dt className="mt-1 text-sm text-gray-600 dark:text-gray-400">{stat.label}</dt>
            </div>
          ))}
        </dl>

        <div
          className="flex flex-col items-center justify-between gap-4 px-6 py-8 transition-colors duration-300 md:flex-row md:px-10"
          style={{ backgroundColor: tint(0.4) }}
        >
          <div>
            <p className="text-2xl font-bold text-gray-900 dark:text-white">Buy now</p>
            <p className="text-lg text-gray-800 tabular-nums dark:text-gray-200">{price}</p>
          </div>
          <button
            type="button"
            onClick={() => onAdd?.(selected)}
            className="inline-flex min-h-11 items-center rounded-md bg-white px-8 text-sm font-medium text-gray-900 hover:bg-gray-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-gray-900"
          >
            Purchase
            <span className="sr-only">
              {' '}
              {name} in {selected.name}
            </span>
          </button>
        </div>
      </div>
    </section>
  )
}
