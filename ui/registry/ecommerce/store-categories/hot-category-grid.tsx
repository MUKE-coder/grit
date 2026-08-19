/*
 * Departments as a dense grid of cut-out product photos.
 *
 * The pattern every large marketplace lands on, and it is worth being explicit
 * about why it looks the way it does, because the obvious version is worse.
 *
 * The photos are not cropped to a square. A kettle, a bed and a lipstick
 * palette are different shapes, and forcing them into identical tiles crops the
 * one part that identifies them. Each image sits inside a fixed box with
 * `object-contain`, so the box is uniform and the product is not.
 *
 * There is no card, no border and no shadow. At this density the chrome outweighs
 * the content: eighteen bordered boxes read as a table of empty containers, and
 * the same eighteen photos on white read as products.
 *
 * The whole tile is one link with a real label. A grid where the image is a link
 * and the caption is a second link to the same place doubles the tab stops and
 * makes a screen reader announce every department twice.
 *
 * Two rows of nine on a wide screen, collapsing to three across on a phone.
 * Nine is not a magic number: it is what fits before the labels start wrapping
 * at the width this grid is usually given.
 */

export interface HotCategory {
  name: string
  href: string
  image: string
}

const CATEGORIES: HotCategory[] = [
  { name: 'TVs', href: '#', image: 'https://images.unsplash.com/photo-1593784991095-a205069470b6?w=300&h=300&fit=crop&q=80' },
  { name: 'Appliances', href: '#', image: 'https://images.unsplash.com/photo-1585659722983-3a675dabf23d?w=300&h=300&fit=crop&q=80' },
  { name: 'Kitchen', href: '#', image: 'https://images.unsplash.com/photo-1556909212-d5b604d0c90d?w=300&h=300&fit=crop&q=80' },
  { name: 'Home', href: '#', image: 'https://images.unsplash.com/photo-1522771739844-6a9f6d5f14af?w=300&h=300&fit=crop&q=80' },
  { name: 'Phones', href: '#', image: 'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=300&h=300&fit=crop&q=80' },
  { name: 'Earphones', href: '#', image: 'https://images.unsplash.com/photo-1590658268037-6bf12165a8df?w=300&h=300&fit=crop&q=80' },
  { name: 'Smartwatches', href: '#', image: 'https://images.unsplash.com/photo-1546868871-7041f2a55e12?w=300&h=300&fit=crop&q=80' },
  { name: 'Personal Care', href: '#', image: 'https://images.unsplash.com/photo-1556228720-195a672e8a03?w=300&h=300&fit=crop&q=80' },
  { name: 'Beauty', href: '#', image: 'https://images.unsplash.com/photo-1596462502278-27bfdc403348?w=300&h=300&fit=crop&q=80' },
  { name: 'Health Care', href: '#', image: 'https://images.unsplash.com/photo-1631549916768-4119b2e5f926?w=300&h=300&fit=crop&q=80' },
  { name: 'Men Shoes', href: '#', image: 'https://images.unsplash.com/photo-1549298916-b41d501d3772?w=300&h=300&fit=crop&q=80' },
  { name: 'Women Shoes', href: '#', image: 'https://images.unsplash.com/photo-1543163521-1bf539c55dd2?w=300&h=300&fit=crop&q=80' },
  { name: 'Kids Shoes', href: '#', image: 'https://images.unsplash.com/photo-1514989940723-e8e51635b782?w=300&h=300&fit=crop&q=80' },
  { name: 'Women Bags', href: '#', image: 'https://images.unsplash.com/photo-1584917865442-de89df76afd3?w=300&h=300&fit=crop&q=80' },
  { name: 'Men Bags', href: '#', image: 'https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=300&h=300&fit=crop&q=80' },
  { name: 'Fridge', href: '#', image: 'https://images.unsplash.com/photo-1571175443880-49e1d25b2bc5?w=300&h=300&fit=crop&q=80' },
  { name: 'Storage Devices', href: '#', image: 'https://images.unsplash.com/photo-1531492746076-161ca9bcad58?w=300&h=300&fit=crop&q=80' },
  { name: 'Laptops', href: '#', image: 'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?w=300&h=300&fit=crop&q=80' },
]

export default function HotCategoryGrid({
  title = 'Hot Category',
  categories = CATEGORIES,
}: {
  title?: string
  categories?: HotCategory[]
}) {
  return (
    <section className="bg-gray-50 py-10 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="rounded-lg bg-white p-6 shadow-sm dark:bg-gray-900">
          <h2 className="text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>

          {/* The rule is on the heading rather than a bare <hr>, which would be
              a separator a screen reader announces for no reason. */}
          <div className="mt-4 border-t border-gray-200 pt-6 dark:border-gray-800">
            <ul className="grid grid-cols-3 gap-x-4 gap-y-8 sm:grid-cols-5 lg:grid-cols-9">
              {categories.map((category) => (
                <li key={category.name}>
                  <a
                    href={category.href}
                    className="group flex flex-col items-center gap-3 rounded-lg p-2 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                  >
                    {/* Fixed box, unfixed product. object-contain keeps the
                        kettle a kettle instead of cropping it to a square. */}
                    <span className="flex h-24 w-24 items-center justify-center">
                      {/* Decorative: the name is directly beneath it in text. */}
                      <img
                        src={category.image}
                        alt=""
                        loading="lazy"
                        className="max-h-24 max-w-24 object-contain transition-transform duration-200 group-hover:scale-105"
                      />
                    </span>
                    <span className="text-center text-sm text-gray-700 group-hover:text-gray-900 dark:text-gray-300 dark:group-hover:text-white">
                      {category.name}
                    </span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  )
}
