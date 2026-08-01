'use client'

import { useEffect, useRef, useState } from 'react'
import {
  ChevronDown,
  ChevronRight,
  Heart,
  Languages,
  MapPin,
  Menu,
  Package,
  Search,
  ShoppingCart,
  User,
  X,
} from 'lucide-react'

/**
 * A marketplace header: utility bar, category tabs, and a full-bleed panel of
 * plain link columns with a promo tile and a brand rail.
 *
 * This is the retail shape, not the SaaS one. The menus elsewhere in this
 * subcategory sell a handful of features and can afford an icon tile and a
 * sentence per entry. A marketplace has three hundred destinations, and the
 * only layout that survives that is dense text columns under bold headings —
 * anything richer turns one department into a full screen of scrolling.
 */

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

interface Column {
  title: string
  links: string[]
}

interface Department {
  label: string
  columns: Column[]
  promo: { eyebrow: string; title: string }
  brands: string[]
}

/* Brand names are invented. A component library that ships real wordmarks hands
   every consumer someone else's trademark, and the layout is the point here —
   swap in your own logos via the `brands` data. */
const DEPARTMENTS: Department[] = [
  {
    label: 'Electronics',
    columns: [
      {
        title: 'Mobiles & Accessories',
        links: [
          'Flagship phones',
          'Premium Androids',
          'Tablets',
          'Headsets & speakers',
          'Wearables',
          'Power banks',
        ],
      },
      {
        title: 'Laptops & Accessories',
        links: [
          'Ultrabooks',
          'Creator laptops',
          'Gaming laptops',
          'Budget laptops',
          'Monitors',
          'Printers',
          'Storage devices',
          'Input devices',
        ],
      },
      {
        title: 'Gaming Essentials',
        links: [
          'Consoles',
          'Gaming accessories',
          'Video games',
          'Gaming monitors',
          'Digital cards',
        ],
      },
      {
        title: 'TVs & Home Entertainment',
        links: ['LED', 'QLED', 'OLED', '4K', '8K', 'Projectors', 'Soundbars', 'Streaming devices'],
      },
      {
        title: 'Cameras',
        links: [
          'Action cameras',
          'DSLR cameras',
          'Surveillance cameras',
          'Instant cameras',
          'Camera accessories',
        ],
      },
    ],
    promo: { eyebrow: 'This week', title: 'Gaming & accessories' },
    brands: ['Northwind', 'Kestrel', 'Lumon', 'Vertex', 'Anvil', 'Halcyon', 'Meridian', 'Corvus'],
  },
  {
    label: 'Home & Kitchen',
    columns: [
      {
        title: 'Kitchen & Dining',
        links: [
          'Cookware',
          'Storage & organisation',
          'Dinnerware',
          'Cutlery',
          'Coffee & tea',
          'Bakeware',
          'Drinkware',
        ],
      },
      {
        title: 'Furniture',
        links: [
          'Sofas & couches',
          'Coffee tables',
          'Gaming chairs',
          'Bean bags',
          'Desks & desk chairs',
          'TV & media units',
          'Storage & cabinets',
        ],
      },
      {
        title: 'Tools & Home Improvement',
        links: [
          'Power tools',
          'Hand tools',
          'Cleaning supplies',
          'Home organisation',
          'Laundry care',
          'Safety & security',
          'Paints & wall supplies',
        ],
      },
      {
        title: 'Home Decor',
        links: [
          'Lighting',
          'Home fragrance',
          'Mats & carpets',
          'Mirrors',
          'Window treatments',
          'Decorative pillows',
        ],
      },
      {
        title: 'Large Appliances',
        links: [
          'Refrigerators',
          'Washing machines',
          'Air conditioners',
          'Cooking ranges',
          'Dishwashers',
          'Dryers',
        ],
      },
    ],
    promo: { eyebrow: 'New season', title: 'Transform your kitchen' },
    brands: ['Prestwick', 'Rowan', 'Copperline', 'Bellhaus', 'Fernwood', 'Ridgely', 'Tallow', 'Oaklet'],
  },
]

const TABS = [
  'Electronics',
  'Beauty & Fragrance',
  'Home & Kitchen',
  'Grocery',
  "Men's Fashion",
  "Women's Fashion",
  'Mom & Baby',
  'Toys',
  "Kids' Fashion",
  'Sports & Outdoors',
  'Health',
]

export default function RetailWithLinkColumns({
  brand = 'noon',
  city = 'Dubai',
  cartCount = 1,
  searchPlaceholder = 'What are you looking for?',
  membershipLabel = 'Get free delivery with membership',
}: {
  brand?: string
  city?: string
  cartCount?: number
  searchPlaceholder?: string
  membershipLabel?: string
}) {
  const [openTab, setOpenTab] = useState<string | null>(null)
  const [mobileOpen, setMobileOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)

  const department = DEPARTMENTS.find((d) => d.label === openTab)

  useEffect(() => {
    function onPointerDown(e: PointerEvent) {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpenTab(null)
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        setOpenTab(null)
        setMobileOpen(false)
      }
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  return (
    <div ref={rootRef} className="relative z-50 bg-white dark:bg-gray-950">
      {/* ── Utility bar ───────────────────────────────────────────────────
          Stays yellow in dark mode. A brand bar is a brand asset, not a
          surface — inverting it is the same mistake as inverting a logo. */}
      <div className="bg-yellow-300 text-gray-900">
        <div className="mx-auto flex max-w-[95rem] items-center gap-4 px-6 py-3.5">
          <a href="#" className="flex shrink-0 items-center gap-2 transition-transform duration-200 active:scale-[0.97]">
            <GritMark className="size-7" />
            <span className="text-2xl leading-none font-extrabold tracking-tight lowercase">
              {brand}
            </span>
          </a>

          <button
            type="button"
            className="hidden shrink-0 items-center gap-1.5 rounded-lg px-2 py-1.5 text-sm font-medium transition-colors hover:bg-black/10 lg:inline-flex"
          >
            <MapPin aria-hidden="true" className="size-4" />
            Deliver to <span className="font-bold">{city}</span>
            <ChevronDown aria-hidden="true" className="size-4" />
          </button>

          <div className="relative min-w-0 flex-1">
            <Search
              aria-hidden="true"
              className="pointer-events-none absolute top-1/2 left-3.5 size-[18px] -translate-y-1/2 text-gray-400"
            />
            <input
              type="search"
              placeholder={searchPlaceholder}
              aria-label="Search products"
              className="h-11 w-full rounded-lg bg-white pr-4 pl-11 text-sm text-gray-900 outline-none placeholder:text-gray-400 focus:ring-2 focus:ring-gray-900/20"
            />
          </div>

          <div className="flex shrink-0 items-center gap-1">
            <UtilityLink icon={Languages} label="EN" className="hidden xl:inline-flex" />
            <UtilityLink icon={User} label="Log in" />
            <UtilityLink icon={Package} label="Orders" className="hidden lg:inline-flex" />
            <UtilityLink icon={Heart} label="Wishlist" className="hidden lg:inline-flex" />
            <a
              href="#"
              className="relative inline-flex items-center gap-2 rounded-lg px-2.5 py-2 text-sm font-semibold transition-colors hover:bg-black/10"
            >
              <span className="relative">
                <ShoppingCart aria-hidden="true" className="size-[18px]" />
                {cartCount > 0 && (
                  <span className="absolute -top-2 -right-2 flex size-[17px] items-center justify-center rounded-full bg-red-600 text-[10px] font-bold text-white tabular-nums">
                    {cartCount > 9 ? '9+' : cartCount}
                  </span>
                )}
              </span>
              <span className="hidden sm:inline">Cart</span>
            </a>

            <button
              type="button"
              onClick={() => setMobileOpen(!mobileOpen)}
              aria-expanded={mobileOpen}
              aria-label="Toggle navigation"
              className="-mr-1 inline-flex size-10 items-center justify-center rounded-lg transition-colors hover:bg-black/10 lg:hidden"
            >
              {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* ── Department tabs ──────────────────────────────────────────────── */}
      <div className="hidden border-b border-gray-200 lg:block dark:border-white/10">
        <div className="mx-auto flex max-w-[95rem] items-center gap-1 px-6">
          {/* The tab strip scrolls rather than wrapping. A marketplace adds
              departments seasonally, and a wrapping row silently becomes two
              rows tall on the day someone adds one. */}
          {/* Masked right edge. Without it the last tab is sliced mid-word right
              against the membership pill, which reads as a broken layout rather
              than as "this row continues". */}
          <div
            className="no-scrollbar flex min-w-0 flex-1 items-center gap-1 overflow-x-auto [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
            style={{
              maskImage: 'linear-gradient(to right, black calc(100% - 2.5rem), transparent)',
              WebkitMaskImage: 'linear-gradient(to right, black calc(100% - 2.5rem), transparent)',
            }}
          >
            {TABS.map((tab) => {
              const hasPanel = DEPARTMENTS.some((d) => d.label === tab)
              const active = openTab === tab
              return (
                <button
                  key={tab}
                  type="button"
                  onClick={() => setOpenTab(active ? null : hasPanel ? tab : null)}
                  aria-expanded={active}
                  className={`relative shrink-0 border-b-2 px-3 py-3.5 text-[14.5px] font-medium whitespace-nowrap transition-colors duration-200 ${
                    active
                      ? 'border-gray-900 text-gray-900 dark:border-white dark:text-white'
                      : 'border-transparent text-gray-600 hover:text-gray-900 dark:text-gray-400 dark:hover:text-white'
                  }`}
                >
                  {tab}
                </button>
              )
            })}
          </div>

          <a
            href="#"
            className="ml-4 hidden shrink-0 items-center gap-1.5 rounded-full border border-yellow-400 bg-yellow-50 px-4 py-1.5 text-[13.5px] font-semibold text-gray-900 transition-colors hover:bg-yellow-100 xl:inline-flex dark:border-yellow-400/40 dark:bg-yellow-400/10 dark:text-yellow-200"
          >
            {membershipLabel}
            <ChevronRight aria-hidden="true" className="size-4" />
          </a>
        </div>
      </div>

      {/* ── Mega panel ───────────────────────────────────────────────────── */}
      {department && (
        <div className="absolute inset-x-0 top-full hidden border-b border-gray-200 bg-white shadow-[0_16px_40px_-16px_rgb(15_23_42_/_0.22)] lg:block dark:border-white/10 dark:bg-gray-900">
          <div className="mx-auto max-w-[95rem] px-6 py-8">
            <div className="flex gap-8">
              {/* Columns take the space they need and no more, so a department
                  with four columns does not stretch them across the viewport. */}
              <div className="grid min-w-0 flex-1 grid-cols-2 gap-x-8 gap-y-6 xl:grid-cols-5">
                {department.columns.map((col) => (
                  <div key={col.title} className="min-w-0">
                    <p className="mb-3 text-[14.5px] leading-snug font-bold text-gray-900 dark:text-white">
                      {col.title}
                    </p>
                    <ul className="space-y-2">
                      {col.links.map((link) => (
                        <li key={link}>
                          <a
                            href="#"
                            className="block text-[13.5px] text-gray-600 transition-colors hover:text-gray-900 hover:underline dark:text-gray-400 dark:hover:text-white"
                          >
                            {link}
                          </a>
                        </li>
                      ))}
                    </ul>
                  </div>
                ))}
              </div>

              <a
                href="#"
                className="group hidden w-72 shrink-0 overflow-hidden rounded-xl 2xl:block"
              >
                <div className="relative flex h-full min-h-[15rem] flex-col justify-end bg-[linear-gradient(150deg,#c7d2e4_0%,#9aa8c4_45%,#2b3446_100%)] p-5 transition-transform duration-300 group-hover:scale-[1.02] dark:bg-[linear-gradient(150deg,#2a3244_0%,#1d2432_45%,#0e131c_100%)]">
                  <p className="text-[11px] font-semibold tracking-wider text-white/70 uppercase">
                    {department.promo.eyebrow}
                  </p>
                  <p className="mt-1 text-lg leading-tight font-bold text-white">
                    {department.promo.title}
                  </p>
                </div>
              </a>
            </div>

            {/* ── Brand rail ─────────────────────────────────────────────── */}
            <div className="mt-8 border-t border-gray-200 pt-6 dark:border-white/10">
              <p className="mb-4 text-[11px] font-bold tracking-wider text-gray-500 uppercase">
                Top brands
              </p>
              <div className="flex flex-wrap gap-3">
                {/* One name per tile. The reference pairs a LOGO with a caption;
                    with a wordmark standing in for the logo, printing the name
                    above and below reads as a duplication bug. Swap the span for
                    an <img> and the caption becomes worth restoring. */}
                {department.brands.map((name) => (
                  <a
                    key={name}
                    href="#"
                    aria-label={name}
                    className="group flex h-14 w-[7.5rem] items-center justify-center rounded-lg border border-gray-200 bg-white transition-colors hover:border-gray-400 dark:border-white/10 dark:bg-white/5 dark:hover:border-white/30"
                  >
                    <span className="text-[13px] font-bold tracking-tight text-gray-800 dark:text-gray-200">
                      {name}
                    </span>
                  </a>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* ── Mobile ───────────────────────────────────────────────────────── */}
      {mobileOpen && (
        <div className="border-b border-gray-200 lg:hidden dark:border-white/10">
          <div className="max-h-[70vh] overflow-y-auto px-6 py-4">
            {DEPARTMENTS.map((dept) => (
              <div key={dept.label} className="border-b border-gray-100 py-3 last:border-0 dark:border-white/5">
                <p className="mb-2 text-sm font-bold text-gray-900 dark:text-white">{dept.label}</p>
                <div className="grid grid-cols-2 gap-x-4 gap-y-1.5">
                  {dept.columns.map((col) => (
                    <a
                      key={col.title}
                      href="#"
                      className="text-[13px] text-gray-600 dark:text-gray-400"
                    >
                      {col.title}
                    </a>
                  ))}
                </div>
              </div>
            ))}
            {/* Departments without a panel still need to be reachable on a
                phone — they are real destinations, just flat ones. */}
            <div className="grid grid-cols-2 gap-x-4 gap-y-1.5 pt-3">
              {TABS.filter((t) => !DEPARTMENTS.some((d) => d.label === t)).map((tab) => (
                <a key={tab} href="#" className="text-[13px] font-medium text-gray-800 dark:text-gray-200">
                  {tab}
                </a>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

function UtilityLink({
  icon: Icon,
  label,
  className = '',
}: {
  icon: typeof User
  label: string
  className?: string
}) {
  return (
    <a
      href="#"
      className={`inline-flex items-center gap-2 rounded-lg px-2.5 py-2 text-sm font-semibold transition-colors hover:bg-black/10 ${className}`}
    >
      <Icon aria-hidden="true" className="size-[18px]" />
      <span className="hidden sm:inline">{label}</span>
    </a>
  )
}
