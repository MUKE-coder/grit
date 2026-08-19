/*
 * The card wall every large retailer puts under its hero: a heading, four
 * labelled pictures, and a link out to the department.
 *
 * It exists because a flat list of category names sells nothing. "Kitchen" is a
 * word; a photo of a pressure cooker labelled Cooker is a thing somebody wants.
 * Each card is a small merchandising slot, and the point of the component is
 * that the slot is data rather than markup.
 *
 * Two card shapes from one type. A card with four tiles renders a 2x2; a card
 * with one renders it large. That is not a variant prop, it is `tiles.length`,
 * because the alternative is a `layout: 'quad' | 'single'` that can disagree
 * with the number of tiles actually supplied.
 *
 * Every tile is a link and so is the footer, which is a real difference from
 * the other category blocks here. On those, the whole card goes to one place
 * and a second link would be a duplicate tab stop. Here each tile goes
 * somewhere different, so they are genuinely separate destinations.
 *
 * The card heading is an h3 inside a section labelled by the h2 above the wall.
 * Eight h2s in a row flattens the page outline into a list of unrelated
 * headings with nothing above them.
 *
 * Images are decorative and the label is text beneath, so alt is empty. A photo
 * of a saucepan captioned "Pots and Pans" does not need to be announced twice.
 */

export interface CategoryTile {
  label: string
  href: string
  image: string
}

export interface MerchandisedCard {
  title: string
  /** One tile renders large, four render as a 2x2. */
  tiles: CategoryTile[]
  linkLabel: string
  linkHref: string
}

const CARDS: MerchandisedCard[] = [
  {
    title: 'New home arrivals under $50',
    linkLabel: 'Shop the latest from Home',
    linkHref: '#',
    tiles: [
      { label: 'Kitchen & Dining', href: '#', image: 'https://images.unsplash.com/photo-1556910103-1c02745aae4d?w=300&h=300&fit=crop&q=80' },
      { label: 'Home Improvement', href: '#', image: 'https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?w=300&h=300&fit=crop&q=80' },
      { label: 'Décor', href: '#', image: 'https://images.unsplash.com/photo-1513519245088-0e12902e5a38?w=300&h=300&fit=crop&q=80' },
      { label: 'Bedding & Bath', href: '#', image: 'https://images.unsplash.com/photo-1522771739844-6a9f6d5f14af?w=300&h=300&fit=crop&q=80' },
    ],
  },
  {
    title: 'Top categories in Kitchen appliances',
    linkLabel: 'Explore all products in Kitchen',
    linkHref: '#',
    tiles: [
      { label: 'Cooker', href: '#', image: 'https://images.unsplash.com/photo-1585515320310-259814833e62?w=600&h=400&fit=crop&q=80' },
    ],
  },
  {
    title: 'Fashion trends you like',
    linkLabel: 'Explore more',
    linkHref: '#',
    tiles: [
      { label: 'Dresses', href: '#', image: 'https://images.unsplash.com/photo-1595777457583-95e059d581b8?w=300&h=300&fit=crop&q=80' },
      { label: 'Knits', href: '#', image: 'https://images.unsplash.com/photo-1576871337622-98d48d1cf531?w=300&h=300&fit=crop&q=80' },
      { label: 'Jackets', href: '#', image: 'https://images.unsplash.com/photo-1551028719-00167b16eac5?w=300&h=300&fit=crop&q=80' },
      { label: 'Jewelry', href: '#', image: 'https://images.unsplash.com/photo-1515562141207-7a88fb7ce338?w=300&h=300&fit=crop&q=80' },
    ],
  },
  {
    title: 'Easy updates for elevated spaces',
    linkLabel: 'Shop home products',
    linkHref: '#',
    tiles: [
      { label: 'Baskets & hampers', href: '#', image: 'https://images.unsplash.com/photo-1595428774223-ef52624120d2?w=300&h=300&fit=crop&q=80' },
      { label: 'Hardware', href: '#', image: 'https://images.unsplash.com/photo-1581092160562-40aa08e78837?w=300&h=300&fit=crop&q=80' },
      { label: 'Accent furniture', href: '#', image: 'https://images.unsplash.com/photo-1555041469-a586c61ea9bc?w=300&h=300&fit=crop&q=80' },
      { label: 'Wallpaper & paint', href: '#', image: 'https://images.unsplash.com/photo-1562113530-57ba467cea38?w=300&h=300&fit=crop&q=80' },
    ],
  },
]

function Tile({ tile, large }: { tile: CategoryTile; large?: boolean }) {
  return (
    <a href={tile.href} className="group block focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
      <span className={`block overflow-hidden bg-gray-100 dark:bg-gray-800 ${large ? 'aspect-[3/2]' : 'aspect-square'}`}>
        {/* Decorative: the label is directly beneath it in text. */}
        <img
          src={tile.image}
          alt=""
          loading="lazy"
          className="h-full w-full object-cover transition-transform duration-200 group-hover:scale-105"
        />
      </span>
      <span className="mt-1.5 block text-xs text-gray-700 group-hover:underline dark:text-gray-300">
        {tile.label}
      </span>
    </a>
  )
}

export default function MerchandisedCategoryCards({
  heading = 'Shop by department',
  cards = CARDS,
}: {
  /* Visually hidden by default: this wall usually sits under a hero that has
     already said what the page is, and a second heading there reads as noise.
     It is still in the outline, which is what a screen reader navigates by. */
  heading?: string
  cards?: MerchandisedCard[]
}) {
  return (
    <section className="bg-gray-100 py-8 dark:bg-gray-950" aria-labelledby="merchandised-heading">
      <h2 id="merchandised-heading" className="sr-only">
        {heading}
      </h2>

      <div className="mx-auto grid max-w-7xl gap-5 px-4 sm:grid-cols-2 sm:px-6 lg:grid-cols-4 lg:px-8">
        {cards.map((card) => (
          <div key={card.title} className="flex flex-col bg-white p-5 dark:bg-gray-900">
            <h3 className="mb-4 text-xl font-bold leading-tight text-gray-900 dark:text-white">
              {card.title}
            </h3>

            {/* One tile is a feature, four are a quad. Driven by the data
                rather than a prop that can contradict it. */}
            {card.tiles.length === 1 ? (
              <Tile tile={card.tiles[0]} large />
            ) : (
              <div className="grid grid-cols-2 gap-3">
                {card.tiles.slice(0, 4).map((tile) => (
                  <Tile key={tile.label} tile={tile} />
                ))}
              </div>
            )}

            {/* mt-auto so the link sits on the baseline of every card in the
                row, whatever the heading wrapped to. */}
            <a
              href={card.linkHref}
              className="mt-auto pt-4 text-sm text-indigo-700 hover:underline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600 dark:text-indigo-400"
            >
              {card.linkLabel}
            </a>
          </div>
        ))}
      </div>
    </section>
  )
}
