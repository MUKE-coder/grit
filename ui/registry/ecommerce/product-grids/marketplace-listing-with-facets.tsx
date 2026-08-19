'use client'

import { useMemo, useState } from 'react'
import { ChevronRight, Heart, Star } from 'lucide-react'

/*
 * A whole marketplace listing page: breadcrumb, facets, results, cards.
 *
 * The sibling `catalog-with-filter-rail` is a filter rail beside a grid. This
 * is the page a marketplace actually ships, and the differences are all things
 * that only show up once real inventory is in it.
 *
 * Price is bands, not a slider. "Under 2000" and "2000 to 5000" are how
 * shoppers think and how merchandising teams report, and a radio group of six
 * bands is operable with a keyboard by anybody, which a two-handle slider is
 * not. The custom range underneath is for the minority who want it, and it is
 * two labelled number inputs with an explicit Go: applying on every keystroke
 * refilters the page while somebody is halfway through typing 20000.
 *
 * Brands carry counts and a search box. A marketplace has three hundred brands,
 * so the list is filtered by typing rather than scrolled, and each count comes
 * from the data rather than being decoration: `(5)` next to a brand that
 * returns nothing is worse than no number at all.
 *
 * Every facet group is a fieldset with a legend. A heading followed by loose
 * checkboxes reads to a screen reader as eleven unrelated controls; the legend
 * is what makes it say "Brand, Binatone, checkbox, not checked".
 *
 * The result count is a live region, so filtering announces "12 of 40 results"
 * instead of silently changing a grid the user cannot see.
 *
 * Card badges are the merchandising layer and they are deliberately typed
 * separately: a discount is computed from the two prices, a fulfilment badge
 * and an official-store flag are facts about the seller. Letting a caller pass
 * "-13%" as free text is how a card ends up claiming a discount that the prices
 * beside it contradict.
 *
 * "No reviews" is stated rather than shown as five empty stars alone. An unrated
 * product and a badly rated one look identical at a glance otherwise.
 */

export interface ListingProduct {
  id: string
  name: string
  price: number
  /** Struck through beside the price. The discount badge is computed from it. */
  wasPrice?: number
  image: string
  href?: string
  seller?: string
  /** e.g. "Same Day Delivery, Lagos". Rendered as-is. */
  fulfilment?: string
  officialStore?: boolean
  rating?: number
  reviews?: number
  brand?: string
}

export interface PriceBand {
  label: string
  min?: number
  max?: number
}

export interface FacetCategory {
  name: string
  href: string
}

const CURRENCY = new Intl.NumberFormat('en-NG', {
  style: 'currency',
  currency: 'NGN',
  maximumFractionDigits: 0,
})

const CATEGORIES: FacetCategory[] = [
  { name: 'Home and Kitchen', href: '#' },
  { name: 'Bulk Products', href: '#' },
]

const BANDS: PriceBand[] = [
  { label: 'Under 2,000', max: 2000 },
  { label: '2,000 - 5,000', min: 2000, max: 5000 },
  { label: '5,000 - 10,000', min: 5000, max: 10000 },
  { label: '10,000 - 20,000', min: 10000, max: 20000 },
  { label: '20,000 - 40,000', min: 20000, max: 40000 },
  { label: 'Above 40,000', min: 40000 },
]

const PRODUCTS: ListingProduct[] = [
  { id: '1', name: 'Centrepiece oval serving dish, stainless steel', price: 39400, image: 'https://images.unsplash.com/photo-1591261730799-ee4e6c2d16d7?w=400&h=400&fit=crop&q=80', seller: 'ORCA OFFICIAL', brand: 'ORCA' },
  { id: '2', name: 'Binatone 2 litres blender with grinder', price: 122099, wasPrice: 140999, image: 'https://images.unsplash.com/photo-1570222094114-d054a817e56b?w=400&h=400&fit=crop&q=80', seller: 'Binatone E-Store', fulfilment: 'Same Day Delivery, Lagos', brand: 'Binatone' },
  { id: '3', name: 'Bosch Series 2 built-in oven, 66 litres', price: 665640, image: 'https://images.unsplash.com/photo-1585659722983-3a675dabf23d?w=400&h=400&fit=crop&q=80', seller: 'Konga Plus', fulfilment: 'Same Day Delivery, Lagos', officialStore: true, brand: 'Bosch' },
  { id: '4', name: 'Koolboks 200L solar chest freezer', price: 2589000, image: 'https://images.unsplash.com/photo-1571175443880-49e1d25b2bc5?w=400&h=400&fit=crop&q=80', seller: 'KOOLBOKS LIMITED', fulfilment: 'Same Day Delivery, Lagos', brand: 'Koolboks' },
  { id: '5', name: 'Haier Thermocool automatic hand dryer', price: 20000, wasPrice: 39000, image: 'https://images.unsplash.com/photo-1584622650111-993a426fbf0a?w=400&h=400&fit=crop&q=80', seller: 'Konga', fulfilment: 'Same Day Delivery, Lagos', officialStore: true, brand: 'Haier Thermocool' },
  { id: '6', name: 'Binatone blender and coffee grinder set', price: 37499, wasPrice: 56299, image: 'https://images.unsplash.com/photo-1522336572468-97b06e8ef143?w=400&h=400&fit=crop&q=80', seller: 'Binatone E-Store', fulfilment: 'Same Day Delivery, Lagos', brand: 'Binatone' },
  { id: '7', name: 'Orosioe curtain and blind cleaning set', price: 21467, wasPrice: 25000, image: 'https://images.unsplash.com/photo-1513694203232-719a280e022f?w=400&h=400&fit=crop&q=80', seller: 'ORCA OFFICIAL', brand: 'ORCA' },
  { id: '8', name: 'Jewellery box with mirror and drawer', price: 55900, wasPrice: 58000, image: 'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=400&h=400&fit=crop&q=80', seller: 'Oliver Davids LTD', rating: 4.4, reviews: 18, brand: 'AEON' },
]

/** Discount is computed, never passed in, so a badge cannot contradict a price. */
function discountOf(product: ListingProduct): number | null {
  if (!product.wasPrice || product.wasPrice <= product.price) return null
  return Math.round(((product.wasPrice - product.price) / product.wasPrice) * 100)
}

function Rating({ rating, reviews }: { rating?: number; reviews?: number }) {
  const filled = Math.round(rating ?? 0)
  return (
    <p className="flex items-center gap-1.5 text-xs">
      <span className="flex" aria-hidden="true">
        {[1, 2, 3, 4, 5].map((n) => (
          <Star
            key={n}
            className={`h-3.5 w-3.5 ${n <= filled ? 'fill-amber-400 text-amber-400' : 'fill-gray-200 text-gray-200 dark:fill-gray-700 dark:text-gray-700'}`}
          />
        ))}
      </span>
      <span className="text-gray-500 dark:text-gray-400">
        {rating ? (
          <>
            {rating.toFixed(1)} out of 5, {reviews} review{reviews === 1 ? '' : 's'}
          </>
        ) : (
          '(No reviews)'
        )}
      </span>
    </p>
  )
}

export default function MarketplaceListingWithFacets({
  breadcrumb = ['Home', 'Listing'],
  categories = CATEGORIES,
  bands = BANDS,
  products = PRODUCTS,
  onAddToCart,
}: {
  breadcrumb?: string[]
  categories?: FacetCategory[]
  bands?: PriceBand[]
  products?: ListingProduct[]
  onAddToCart?: (product: ListingProduct) => void
}) {
  const [band, setBand] = useState<string>('')
  const [brands, setBrands] = useState<string[]>([])
  const [brandQuery, setBrandQuery] = useState('')
  const [minPrice, setMinPrice] = useState('')
  const [maxPrice, setMaxPrice] = useState('')
  /* Held separately from the inputs so typing does not refilter the page. The
     Go button copies one into the other. */
  const [appliedRange, setAppliedRange] = useState<{ min?: number; max?: number }>({})
  const [sort, setSort] = useState('relevance')

  /* Counts come from the data. A hardcoded (5) beside a brand that returns
     nothing is worse than showing no number. */
  const brandCounts = useMemo(() => {
    const counts = new Map<string, number>()
    for (const product of products) {
      if (!product.brand) continue
      counts.set(product.brand, (counts.get(product.brand) ?? 0) + 1)
    }
    return [...counts.entries()].sort((a, b) => b[1] - a[1])
  }, [products])

  const visibleBrands = brandCounts.filter(([name]) =>
    name.toLowerCase().includes(brandQuery.toLowerCase()),
  )

  const results = useMemo(() => {
    const selected = bands.find((b) => b.label === band)
    const min = appliedRange.min ?? selected?.min
    const max = appliedRange.max ?? selected?.max

    const filtered = products.filter((product) => {
      if (min !== undefined && product.price < min) return false
      if (max !== undefined && product.price > max) return false
      if (brands.length && (!product.brand || !brands.includes(product.brand))) return false
      return true
    })

    if (sort === 'low') return [...filtered].sort((a, b) => a.price - b.price)
    if (sort === 'high') return [...filtered].sort((a, b) => b.price - a.price)
    return filtered
  }, [products, bands, band, brands, appliedRange, sort])

  function toggleBrand(name: string) {
    setBrands((current) =>
      current.includes(name) ? current.filter((b) => b !== name) : [...current, name],
    )
  }

  return (
    <div className="bg-gray-50 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <nav aria-label="Breadcrumb" className="mb-4">
          <ol className="flex items-center gap-1 text-sm">
            {breadcrumb.map((crumb, i) => {
              const last = i === breadcrumb.length - 1
              return (
                <li key={crumb} className="flex items-center gap-1">
                  {last ? (
                    <span aria-current="page" className="font-medium text-pink-600 dark:text-pink-400">
                      {crumb}
                    </span>
                  ) : (
                    <>
                      <a href="#" className="text-gray-600 hover:underline dark:text-gray-400">
                        {crumb}
                      </a>
                      <ChevronRight className="h-4 w-4 text-gray-400" aria-hidden="true" />
                    </>
                  )}
                </li>
              )
            })}
          </ol>
        </nav>

        <div className="flex flex-col gap-6 lg:flex-row">
          <aside className="w-full shrink-0 lg:w-64">
            <div className="rounded-lg bg-white p-4 dark:bg-gray-900">
              <h2 className="text-base font-semibold text-gray-900 dark:text-white">Filters</h2>

              <div className="mt-4 border-t border-gray-200 pt-4 dark:border-gray-800">
                <h3 className="mb-2 text-sm font-semibold text-gray-900 dark:text-white">
                  Browse categories
                </h3>
                <ul className="space-y-1">
                  {categories.map((category) => (
                    <li key={category.name}>
                      <a
                        href={category.href}
                        className="flex min-h-9 items-center justify-between rounded px-1 text-sm text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
                      >
                        {category.name}
                        <ChevronRight className="h-4 w-4 text-gray-400" aria-hidden="true" />
                      </a>
                    </li>
                  ))}
                </ul>
              </div>

              <fieldset className="mt-4 border-t border-gray-200 pt-4 dark:border-gray-800">
                <legend className="mb-2 text-sm font-semibold text-gray-900 dark:text-white">
                  Price
                </legend>
                <div className="space-y-1">
                  {bands.map((option) => (
                    <label
                      key={option.label}
                      className="flex min-h-9 cursor-pointer items-center gap-2 text-sm text-gray-700 dark:text-gray-300"
                    >
                      <input
                        type="radio"
                        name="price-band"
                        checked={band === option.label}
                        onChange={() => {
                          setBand(option.label)
                          // A band and a custom range are two answers to one
                          // question. Choosing a band clears the range rather
                          // than silently intersecting with it.
                          setAppliedRange({})
                          setMinPrice('')
                          setMaxPrice('')
                        }}
                        className="h-4 w-4 border-gray-300 text-pink-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
                      />
                      {option.label}
                    </label>
                  ))}
                </div>

                <div className="mt-3">
                  <p className="mb-1.5 text-xs font-medium text-gray-700 dark:text-gray-300">
                    Custom price range
                  </p>
                  <div className="flex items-center gap-2">
                    <label className="sr-only" htmlFor="min-price">
                      Minimum price
                    </label>
                    <input
                      id="min-price"
                      inputMode="numeric"
                      value={minPrice}
                      onChange={(e) => setMinPrice(e.target.value)}
                      placeholder="Min"
                      className="min-h-9 w-full rounded border border-gray-300 px-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                    />
                    <label className="sr-only" htmlFor="max-price">
                      Maximum price
                    </label>
                    <input
                      id="max-price"
                      inputMode="numeric"
                      value={maxPrice}
                      onChange={(e) => setMaxPrice(e.target.value)}
                      placeholder="Max"
                      className="min-h-9 w-full rounded border border-gray-300 px-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                    />
                    {/* Explicit, because applying on keystroke refilters the
                        page while somebody is halfway through typing 20000. */}
                    <button
                      type="button"
                      onClick={() => {
                        setBand('')
                        setAppliedRange({
                          min: minPrice ? Number(minPrice) : undefined,
                          max: maxPrice ? Number(maxPrice) : undefined,
                        })
                      }}
                      className="min-h-9 shrink-0 rounded bg-pink-600 px-3 text-sm font-medium text-white hover:bg-pink-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
                    >
                      Go
                    </button>
                  </div>
                </div>
              </fieldset>

              <fieldset className="mt-4 border-t border-gray-200 pt-4 dark:border-gray-800">
                <legend className="mb-2 text-sm font-semibold text-gray-900 dark:text-white">
                  Brand
                </legend>
                <label className="sr-only" htmlFor="brand-search">
                  Search brands
                </label>
                <input
                  id="brand-search"
                  type="search"
                  value={brandQuery}
                  onChange={(e) => setBrandQuery(e.target.value)}
                  placeholder="Search brands..."
                  className="mb-2 min-h-9 w-full rounded border border-gray-300 px-2 text-sm dark:border-gray-700 dark:bg-gray-800 dark:text-white"
                />
                <div className="max-h-56 space-y-1 overflow-y-auto">
                  {visibleBrands.map(([name, count]) => (
                    <label
                      key={name}
                      className="flex min-h-9 cursor-pointer items-center gap-2 text-sm text-gray-700 dark:text-gray-300"
                    >
                      <input
                        type="checkbox"
                        checked={brands.includes(name)}
                        onChange={() => toggleBrand(name)}
                        className="h-4 w-4 rounded border-gray-300 text-pink-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
                      />
                      <span className="flex-1">{name}</span>
                      <span className="text-gray-500 dark:text-gray-400">({count})</span>
                    </label>
                  ))}
                  {visibleBrands.length === 0 && (
                    <p className="py-2 text-sm text-gray-500 dark:text-gray-400">
                      No brands match that.
                    </p>
                  )}
                </div>
              </fieldset>
            </div>
          </aside>

          <div className="min-w-0 flex-1">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3 rounded-lg bg-white px-4 py-3 dark:bg-gray-900">
              {/* Live, so filtering announces the new count to a screen reader
                  rather than silently changing a grid it cannot see. */}
              <p aria-live="polite" className="text-sm text-gray-700 dark:text-gray-300">
                {results.length === products.length
                  ? `1 - ${results.length} of ${results.length} results`
                  : `${results.length} of ${products.length} results`}
              </p>
              <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                Sort by:
                <select
                  value={sort}
                  onChange={(e) => setSort(e.target.value)}
                  className="min-h-9 rounded border border-gray-300 bg-white px-2 text-sm font-medium text-pink-600 dark:border-gray-700 dark:bg-gray-800"
                >
                  <option value="relevance">Relevance</option>
                  <option value="low">Price: low to high</option>
                  <option value="high">Price: high to low</option>
                </select>
              </label>
            </div>

            <ul className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {results.map((product) => {
                const discount = discountOf(product)
                return (
                  <li
                    key={product.id}
                    className="flex flex-col overflow-hidden rounded-lg bg-white dark:bg-gray-900"
                  >
                    <div className="relative bg-gray-50 dark:bg-gray-800">
                      <img
                        src={product.image}
                        alt=""
                        loading="lazy"
                        className="aspect-square w-full object-cover"
                      />
                      {discount !== null && (
                        <span className="absolute right-2 top-2 rounded bg-pink-100 px-1.5 py-0.5 text-xs font-semibold text-pink-700">
                          -{discount}%
                        </span>
                      )}
                      {product.officialStore && (
                        <span className="absolute bottom-2 left-2 rounded bg-amber-400 px-1.5 py-0.5 text-[11px] font-semibold text-gray-900">
                          Official Store
                        </span>
                      )}
                    </div>

                    <div className="flex flex-1 flex-col gap-1.5 p-3">
                      <h3 className="text-sm text-gray-900 dark:text-white">
                        <a
                          href={product.href ?? '#'}
                          className="line-clamp-2 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
                        >
                          {product.name}
                        </a>
                      </h3>

                      <p className="flex items-baseline gap-2">
                        <span className="text-base font-bold text-gray-900 dark:text-white">
                          {CURRENCY.format(product.price)}
                        </span>
                        {product.wasPrice && (
                          <s className="text-xs text-gray-500 dark:text-gray-400">
                            {CURRENCY.format(product.wasPrice)}
                          </s>
                        )}
                      </p>

                      {product.fulfilment && (
                        <p className="text-xs text-pink-600 dark:text-pink-400">{product.fulfilment}</p>
                      )}
                      {product.seller && (
                        <p className="text-xs text-gray-600 dark:text-gray-400">
                          Sold by:{' '}
                          <span className="text-pink-600 dark:text-pink-400">{product.seller}</span>
                        </p>
                      )}

                      <Rating rating={product.rating} reviews={product.reviews} />

                      <div className="mt-auto flex items-center gap-2 pt-2">
                        {/* Named for a screen reader. Eight buttons all called
                            "Add to Cart" is a list of identical controls with
                            no way to tell which is which. */}
                        <button
                          type="button"
                          onClick={() => onAddToCart?.(product)}
                          className="min-h-11 flex-1 rounded bg-pink-600 text-sm font-semibold text-white hover:bg-pink-700 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600"
                        >
                          Add to Cart
                          <span className="sr-only">: {product.name}</span>
                        </button>
                        <button
                          type="button"
                          className="flex h-11 w-11 shrink-0 items-center justify-center rounded border border-gray-200 text-gray-500 hover:text-pink-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-pink-600 dark:border-gray-700"
                        >
                          <Heart className="h-4 w-4" aria-hidden="true" />
                          <span className="sr-only">Save {product.name} for later</span>
                        </button>
                      </div>
                    </div>
                  </li>
                )
              })}
            </ul>

            {results.length === 0 && (
              <p className="rounded-lg bg-white p-8 text-center text-sm text-gray-600 dark:bg-gray-900 dark:text-gray-400">
                Nothing matches those filters. Widen the price range or clear a brand.
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
