'use client'

import { useEffect, useRef, useState } from 'react'
import {
  ChevronLeft,
  ChevronRight,
  Globe,
  Headphones,
  Heart,
  LayoutGrid,
  Menu,
  Search,
  ShoppingBag,
  User,
  X,
} from 'lucide-react'

/**
 * A storefront header: promo strip, dark utility bar, a scrolling department
 * row, and a panel built from a category rail and circular tiles.
 *
 * The tiles are the whole idea. Fashion catalogues are browsed by look rather
 * than by name — "Women mid-calf boots" means nothing until you see the shape —
 * so the picture carries the navigation and the label only confirms it. That is
 * the opposite trade from the link-column marketplace header, which has too
 * many destinations for any of them to earn a thumbnail.
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

interface Tile {
  name: string
  /** Real product shot. Omitted, the tile falls back to a tinted placeholder. */
  image?: string
  /** Renders the grid glyph instead of a picture. */
  viewAll?: boolean
}

interface Group {
  title: string
  tiles: Tile[]
}

interface Department {
  label: string
  shopBy: Tile[]
  groups: Group[]
}

const DEPARTMENTS: Department[] = [
  {
    label: 'Shoes',
    shopBy: [
      { name: 'View All', viewAll: true },
      { name: 'Top Rated' },
      { name: 'School Picks' },
      { name: 'Sale' },
      { name: 'Fashion Boots' },
      { name: 'Sneakers' },
      { name: 'Pumps' },
      { name: 'Sandals' },
      { name: 'Flats' },
      { name: 'Slippers' },
      { name: 'Wedges & Flatform' },
      { name: 'Outdoor Shoes' },
    ],
    groups: [
      {
        title: 'Boots',
        tiles: [
          { name: 'View All', viewAll: true },
          { name: 'Ankle Boots & Booties' },
          { name: 'Mid-calf Boots' },
          { name: 'Over-the-knee Boots' },
          { name: 'Knee-high Boots' },
          { name: 'Equestrian & Western' },
          { name: 'Snow Boots' },
        ],
      },
      {
        title: 'Sneakers',
        tiles: [
          { name: 'View All', viewAll: true },
          { name: 'Casual Shoes' },
          { name: 'Sports Shoes' },
          { name: 'Wedge Sneakers' },
          { name: 'Canvas Shoes' },
          { name: 'Running Shoes' },
        ],
      },
    ],
  },
  {
    label: 'Women Clothing',
    shopBy: [
      { name: 'View All', viewAll: true },
      { name: 'New In' },
      { name: 'Top Rated' },
      { name: 'Tops' },
      { name: 'T-shirts' },
      { name: 'Blouses' },
      { name: 'Tank Tops & Camis' },
      { name: 'Dresses' },
      { name: 'Mini Dresses' },
      { name: 'Long Dresses' },
      { name: 'Co-ords' },
      { name: 'Bottoms' },
    ],
    groups: [
      {
        title: 'Dresses',
        tiles: [
          { name: 'View All', viewAll: true },
          { name: 'Long Dresses' },
          { name: 'Short Dresses' },
          { name: 'Midi Dresses' },
          { name: 'Mini Dresses' },
          { name: 'Maxi Dresses' },
        ],
      },
      {
        title: 'Bottoms',
        tiles: [
          { name: 'View All', viewAll: true },
          { name: 'Pants' },
          { name: 'Skirts' },
          { name: 'Shorts' },
          { name: 'Sweatpants' },
          { name: 'Leggings' },
        ],
      },
    ],
  },
]

const RAIL = [
  'New In',
  'Sale',
  'Women Clothing',
  'Beachwear',
  'Kids',
  'Curve',
  'Men Clothing',
  'Shoes',
  'Underwear & Sleepwear',
  'Home & Living',
  'Jewelry & Accessories',
  'Beauty & Health',
  'Baby & Maternity',
  'Bags & Luggage',
]

const TABS = [
  'New In',
  'Sale',
  'Women Clothing',
  'Beachwear',
  'Kids',
  'Curve',
  'Men Clothing',
  'Shoes',
  'Underwear & Sleepwear',
  'Home & Living',
  'Jewelry & Accessories',
  'Beauty & Health',
]

/**
 * Deterministic hue from the label.
 *
 * Math.random() here would produce one colour on the server and another in the
 * browser, and React would report a hydration mismatch on a decorative
 * placeholder. Hashing the name gives every tile its own tint and the same one
 * on both sides.
 */
function hueOf(seed: string): number {
  let h = 0
  for (let i = 0; i < seed.length; i++) h = (h * 31 + seed.charCodeAt(i)) % 360
  return h
}

function CategoryTile({ tile }: { tile: Tile }) {
  const hue = hueOf(tile.name)
  return (
    <a href="#" className="group flex flex-col items-center gap-2 text-center">
      <span
        className="flex size-[4.25rem] shrink-0 items-center justify-center overflow-hidden rounded-full ring-1 ring-gray-200 transition-all duration-200 group-hover:ring-2 group-hover:ring-gray-900 dark:ring-white/10 dark:group-hover:ring-white"
        style={
          tile.viewAll || tile.image
            ? undefined
            : {
                background: `linear-gradient(140deg, hsl(${hue} 42% 88%), hsl(${(hue + 45) % 360} 38% 76%))`,
              }
        }
      >
        {tile.viewAll ? (
          <span className="flex size-full items-center justify-center bg-gray-100 text-gray-500 dark:bg-white/10 dark:text-gray-400">
            <LayoutGrid aria-hidden="true" className="size-6" />
          </span>
        ) : tile.image ? (
          // eslint-disable-next-line @next/next/no-img-element
          <img src={tile.image} alt="" className="size-full object-cover" />
        ) : null}
      </span>
      <span className="max-w-[6.5rem] text-[12.5px] leading-tight text-gray-700 transition-colors group-hover:text-gray-900 group-hover:underline dark:text-gray-300 dark:group-hover:text-white">
        {tile.name}
      </span>
    </a>
  )
}

export default function StorefrontWithCategoryTiles({
  brand = 'GRIT',
  promoLeft = 'Free shipping on orders over $15',
  promoRight = 'Free returns on all orders',
  searchPlaceholder = 'Search products',
  cartCount = 0,
  wishlistCount = 0,
}: {
  brand?: string
  promoLeft?: string
  promoRight?: string
  searchPlaceholder?: string
  cartCount?: number
  wishlistCount?: number
}) {
  const [openTab, setOpenTab] = useState<string | null>(null)
  const [railActive, setRailActive] = useState<string>('Shoes')
  const [mobileOpen, setMobileOpen] = useState(false)
  const rootRef = useRef<HTMLDivElement>(null)
  const tabsRef = useRef<HTMLDivElement>(null)

  // The rail drives the panel once it is open, so hovering "Women Clothing" in
  // the rail swaps the tiles without closing and reopening the menu.
  const department =
    DEPARTMENTS.find((d) => d.label === railActive) ??
    DEPARTMENTS.find((d) => d.label === openTab) ??
    DEPARTMENTS[0]

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

  const scrollTabs = (dir: 1 | -1) =>
    tabsRef.current?.scrollBy({ left: dir * 240, behavior: 'smooth' })

  function openDepartment(tab: string) {
    const known = DEPARTMENTS.some((d) => d.label === tab)
    if (openTab === tab) {
      setOpenTab(null)
      return
    }
    setOpenTab(tab)
    // Point the rail at the department that was clicked when there is a panel
    // for it; otherwise leave the rail where it was rather than showing an
    // empty menu.
    if (known) setRailActive(tab)
  }

  return (
    <div ref={rootRef} className="relative z-50 bg-white dark:bg-gray-950">
      {/* ── Promo strip ──────────────────────────────────────────────────── */}
      <div className="bg-amber-50 dark:bg-amber-500/10">
        <div className="mx-auto flex max-w-[95rem] items-center justify-center gap-8 px-6 py-2.5 text-[12.5px] text-gray-700 dark:text-amber-100/80">
          <span className="font-medium">{promoLeft}</span>
          <span aria-hidden className="hidden h-3 w-px bg-gray-300 sm:block dark:bg-white/20" />
          <span className="hidden font-medium sm:inline">{promoRight}</span>
        </div>
      </div>

      {/* ── Utility bar ──────────────────────────────────────────────────── */}
      <div className="bg-gray-950 text-white dark:bg-black">
        <div className="mx-auto flex max-w-[95rem] items-center gap-6 px-6 py-3.5">
          <a
            href="#"
            className="flex shrink-0 items-center gap-2 transition-transform duration-200 active:scale-[0.97]"
          >
            <GritMark className="size-7" />
            <span className="text-2xl leading-none font-extrabold tracking-[0.12em]">{brand}</span>
          </a>

          <div className="hidden min-w-0 flex-1 justify-center md:flex">
            <div className="flex h-11 w-full max-w-2xl overflow-hidden rounded-sm bg-white">
              <input
                type="search"
                placeholder={searchPlaceholder}
                aria-label="Search products"
                className="min-w-0 flex-1 px-4 text-sm text-gray-900 outline-none placeholder:text-gray-400"
              />
              <button
                type="button"
                aria-label="Search"
                className="flex w-14 shrink-0 items-center justify-center bg-gray-950 text-white transition-colors hover:bg-gray-800"
              >
                <Search aria-hidden="true" className="size-[18px]" />
              </button>
            </div>
          </div>

          <div className="ml-auto flex shrink-0 items-center gap-1">
            <IconLink icon={User} label="Account" />
            <IconLink icon={ShoppingBag} label="Cart" count={cartCount} />
            <IconLink icon={Heart} label="Wishlist" count={wishlistCount} />
            <IconLink icon={Headphones} label="Support" className="hidden lg:inline-flex" />
            <IconLink icon={Globe} label="Language" className="hidden lg:inline-flex" />

            <button
              type="button"
              onClick={() => setMobileOpen(!mobileOpen)}
              aria-expanded={mobileOpen}
              aria-label="Toggle navigation"
              className="-mr-1 inline-flex size-10 items-center justify-center rounded-lg transition-colors hover:bg-white/10 lg:hidden"
            >
              {mobileOpen ? <X className="size-5" /> : <Menu className="size-5" />}
            </button>
          </div>
        </div>

        {/* ── Department tabs ────────────────────────────────────────────── */}
        <div className="hidden border-t border-white/10 lg:block">
          <div className="mx-auto flex max-w-[95rem] items-center px-6">
            <button
              type="button"
              onClick={() => openDepartment('Categories')}
              className="mr-2 shrink-0 px-3 py-3 text-[14px] font-semibold whitespace-nowrap text-white/90 hover:text-white"
            >
              Categories
            </button>

            <div
              ref={tabsRef}
              className="no-scrollbar flex min-w-0 flex-1 items-center overflow-x-auto [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden"
            >
              {TABS.map((tab) => {
                const active = openTab === tab
                return (
                  <button
                    key={tab}
                    type="button"
                    onClick={() => openDepartment(tab)}
                    aria-expanded={active}
                    className={`shrink-0 px-3.5 py-3 text-[14px] whitespace-nowrap transition-colors duration-200 ${
                      active
                        ? 'bg-white/15 font-semibold text-white'
                        : 'text-white/80 hover:text-white'
                    }`}
                  >
                    {tab}
                  </button>
                )
              })}
            </div>

            {/* Arrows sit outside the scroller so they never scroll away with
                the content they control. */}
            <div className="ml-2 flex shrink-0 items-center gap-0.5">
              <button
                type="button"
                onClick={() => scrollTabs(-1)}
                aria-label="Scroll departments left"
                className="flex size-7 items-center justify-center rounded text-white/70 transition-colors hover:bg-white/10 hover:text-white"
              >
                <ChevronLeft aria-hidden="true" className="size-4" />
              </button>
              <button
                type="button"
                onClick={() => scrollTabs(1)}
                aria-label="Scroll departments right"
                className="flex size-7 items-center justify-center rounded text-white/70 transition-colors hover:bg-white/10 hover:text-white"
              >
                <ChevronRight aria-hidden="true" className="size-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* ── Mega panel ───────────────────────────────────────────────────── */}
      {openTab && (
        <div className="absolute inset-x-0 top-full hidden border-b border-gray-200 bg-white shadow-[0_16px_40px_-16px_rgb(15_23_42_/_0.22)] lg:block dark:border-white/10 dark:bg-gray-900">
          <div
            className="mx-auto flex max-h-[32rem] max-w-[95rem] gap-0 px-6 py-5"
            style={{
              maskImage: 'linear-gradient(to bottom, black calc(100% - 1.5rem), transparent)',
              WebkitMaskImage: 'linear-gradient(to bottom, black calc(100% - 1.5rem), transparent)',
            }}
          >
            {/* Rail */}
            <div className="w-56 shrink-0 overflow-y-auto pr-4 pb-2 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
              <nav className="space-y-0.5">
                {RAIL.map((item) => {
                  const active = railActive === item
                  const known = DEPARTMENTS.some((d) => d.label === item)
                  return (
                    <button
                      key={item}
                      type="button"
                      onMouseEnter={() => known && setRailActive(item)}
                      onFocus={() => known && setRailActive(item)}
                      onClick={() => known && setRailActive(item)}
                      className={`flex w-full items-center justify-between gap-2 rounded-lg px-3 py-2.5 text-left text-[13.5px] transition-colors duration-150 ${
                        active
                          ? 'bg-gray-100 font-semibold text-gray-900 dark:bg-white/10 dark:text-white'
                          : 'text-gray-700 hover:bg-gray-50 dark:text-gray-300 dark:hover:bg-white/5'
                      }`}
                    >
                      <span className="truncate">{item}</span>
                      <ChevronRight aria-hidden="true" className="size-4 shrink-0 text-gray-400" />
                    </button>
                  )
                })}
              </nav>
            </div>

            {/* Shop by category */}
            <div className="w-[26rem] shrink-0 overflow-y-auto border-x border-gray-200 px-6 pb-4 dark:border-white/10 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
              <p className="mb-5 flex items-center gap-2 text-[12px] font-bold tracking-wider text-gray-900 uppercase dark:text-white">
                <LayoutGrid aria-hidden="true" className="size-4 text-orange-500" />
                Shop by category
              </p>
              <div className="grid grid-cols-3 gap-x-4 gap-y-6">
                {department.shopBy.map((tile) => (
                  <CategoryTile key={tile.name} tile={tile} />
                ))}
              </div>
            </div>

            {/* Grouped sections */}
            <div className="min-w-0 flex-1 overflow-y-auto pl-6 pb-4 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
              {department.groups.map((group) => (
                <div key={group.title} className="mb-8 last:mb-0">
                  <p className="mb-5 text-[12px] font-bold tracking-wider text-gray-900 uppercase dark:text-white">
                    {group.title}
                  </p>
                  <div className="flex flex-wrap gap-x-5 gap-y-6">
                    {group.tiles.map((tile) => (
                      <div key={tile.name} className="w-[6.5rem]">
                        <CategoryTile tile={tile} />
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* ── Mobile ───────────────────────────────────────────────────────── */}
      {mobileOpen && (
        <div className="border-b border-gray-200 lg:hidden dark:border-white/10">
          <div className="max-h-[70vh] overflow-y-auto px-6 py-4">
            <div className="mb-4 flex h-11 overflow-hidden rounded-sm border border-gray-300 dark:border-white/15">
              <input
                type="search"
                placeholder={searchPlaceholder}
                aria-label="Search products"
                className="min-w-0 flex-1 bg-transparent px-3 text-sm text-gray-900 outline-none placeholder:text-gray-400 dark:text-white"
              />
              <button
                type="button"
                aria-label="Search"
                className="flex w-12 shrink-0 items-center justify-center bg-gray-950 text-white dark:bg-white dark:text-gray-900"
              >
                <Search aria-hidden="true" className="size-[18px]" />
              </button>
            </div>
            {RAIL.map((item) => (
              <a
                key={item}
                href="#"
                className="flex items-center justify-between border-b border-gray-100 py-3 text-[14px] text-gray-800 last:border-0 dark:border-white/5 dark:text-gray-200"
              >
                {item}
                <ChevronRight aria-hidden="true" className="size-4 text-gray-400" />
              </a>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

function IconLink({
  icon: Icon,
  label,
  count,
  className = '',
}: {
  icon: typeof User
  label: string
  count?: number
  className?: string
}) {
  return (
    <a
      href="#"
      aria-label={label}
      className={`relative inline-flex size-10 items-center justify-center rounded-lg text-white/90 transition-colors hover:bg-white/10 hover:text-white ${className}`}
    >
      <Icon aria-hidden="true" className="size-[19px]" />
      {typeof count === 'number' && (
        <span className="absolute top-1 right-0.5 min-w-[15px] rounded-full bg-white px-1 text-center text-[10px] leading-[15px] font-bold text-gray-900 tabular-nums">
          {count > 99 ? '99+' : count}
        </span>
      )}
    </a>
  )
}
