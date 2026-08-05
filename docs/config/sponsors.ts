import sponsorsData from '@/data/sponsors.json'

/**
 * Single source of truth for sponsorship — the tier ladder, the roster, and the
 * derived totals. Every surface reads from here: /sponsor, /sponsors, the home
 * strip, and the docs CTA. Never hardcode a tier price or perk in a page.
 *
 * Roster lives in data/sponsors.json (see sponsors.schema.json) so it can be
 * hand-edited, or written by a sync script later, without touching components.
 */

export type SponsorTierId = 'supporter' | 'backer' | 'special' | 'partner'
export type SponsorInterval = 'month' | 'once'

export interface SponsorTier {
  id: SponsorTierId
  name: string
  /** USD per month. Also the Stripe unit_amount (x100) for the monthly plan. */
  price: number
  /** One-line pitch shown under the tier name. */
  blurb: string
  perks: string[]
  /** The tier we visually lead with on /sponsor. */
  featured?: boolean
}

/**
 * The ladder. Perks are promises — only add one we will actually honour.
 * Ordered cheapest → dearest; the UI relies on this order.
 */
export const SPONSOR_TIERS: SponsorTier[] = [
  {
    id: 'supporter',
    name: 'Supporter',
    price: 5,
    blurb: 'Keep the lights on and get your name on the project.',
    perks: [
      'Your name, photo and links on the Sponsors page',
      'Listed in the Grit README',
      'Supporter badge on your GitHub profile',
    ],
  },
  {
    id: 'backer',
    name: 'Backer',
    price: 25,
    blurb: 'For developers and small teams shipping with Grit.',
    perks: [
      'Everything in Supporter',
      'Your logo on the Sponsors page',
      'Logo in the docs footer',
      'Early access to new releases',
    ],
  },
  {
    id: 'special',
    name: 'Special Sponsor',
    price: 100,
    blurb: 'Top billing across the project. The sweet spot for companies.',
    perks: [
      'Everything in Backer',
      'Highlighted above all other sponsors',
      'Logo on the Grit home page',
      'Logo inside the CLI (grit new)',
      'Link from the README header',
    ],
    featured: true,
  },
  {
    id: 'partner',
    name: 'Partner',
    price: 500,
    blurb: 'Grit is core to your stack and you want everyone to know it.',
    perks: [
      'Everything in Special Sponsor',
      'Hero placement on the home page',
      '“Sponsored by” credit on the docs site',
      'Roadmap input: a direct line to the maintainer',
      'Priority support',
    ],
  },
]

/** One-time amounts offered on /sponsor. `custom` is handled separately. */
export const ONE_TIME_AMOUNTS = [25, 100, 500] as const

export const SPONSOR_MIN_USD = 1
/** Mirrors the cap enforced in app/api/checkout/route.ts. */
export const SPONSOR_MAX_USD = 999

export interface Sponsor {
  name: string
  tier: SponsorTierId
  /** YYYY-MM */
  since: string
  interval?: SponsorInterval
  amount?: number
  lifetime?: number
  avatar?: string
  github?: string
  website?: string
  past?: boolean
}

interface SponsorsFile {
  lastSync: string
  sponsors: Sponsor[]
}

const data = sponsorsData as unknown as SponsorsFile

export const SPONSORS_LAST_SYNC: string = data.lastSync
export const ALL_SPONSORS: Sponsor[] = data.sponsors ?? []

export function getTier(id: SponsorTierId): SponsorTier {
  const tier = SPONSOR_TIERS.find((t) => t.id === id)
  // A roster entry naming an unknown tier is a data bug; fail loudly in dev
  // rather than rendering a blank card.
  if (!tier) throw new Error(`Unknown sponsor tier: ${id}`)
  return tier
}

/** Rank for sorting — higher tier first, then bigger lifetime. */
function tierRank(id: SponsorTierId): number {
  return SPONSOR_TIERS.findIndex((t) => t.id === id)
}

export function sponsorAmount(sponsor: Sponsor): number {
  return sponsor.amount ?? getTier(sponsor.tier).price
}

/** Active sponsors, best billing first — the order every grid renders in. */
export function activeSponsors(): Sponsor[] {
  return ALL_SPONSORS.filter((s) => !s.past).sort(
    (a, b) =>
      tierRank(b.tier) - tierRank(a.tier) ||
      (b.lifetime ?? 0) - (a.lifetime ?? 0) ||
      a.name.localeCompare(b.name)
  )
}

export function specialSponsors(): Sponsor[] {
  return activeSponsors().filter(
    (s) => s.tier === 'special' || s.tier === 'partner'
  )
}

export interface SponsorStats {
  lifetimeFunding: number
  monthlyRecurring: number
  activeCount: number
  specialCount: number
  allTimeCount: number
  hasAny: boolean
}

/**
 * Derived, never stored — so the numbers can't drift from the roster.
 * `hasAny` lets surfaces render a founding-sponsor state instead of a wall of
 * zeros while the roster is empty.
 */
export function sponsorStats(): SponsorStats {
  const active = activeSponsors()
  return {
    lifetimeFunding: ALL_SPONSORS.reduce((sum, s) => sum + (s.lifetime ?? 0), 0),
    monthlyRecurring: active
      .filter((s) => (s.interval ?? 'month') === 'month')
      .reduce((sum, s) => sum + sponsorAmount(s), 0),
    activeCount: active.length,
    specialCount: specialSponsors().length,
    allTimeCount: ALL_SPONSORS.length,
    hasAny: ALL_SPONSORS.length > 0,
  }
}

/** "2026-07" → "JUL 2026" (the SINCE label on a sponsor card). */
export function formatSince(since: string): string {
  const [year, month] = since.split('-')
  const date = new Date(Number(year), Number(month) - 1, 1)
  return date
    .toLocaleDateString('en-US', { month: 'short', year: 'numeric' })
    .toUpperCase()
}

/** "$100 a month" / "$500 one time" */
export function formatContribution(sponsor: Sponsor): string {
  const amount = sponsorAmount(sponsor)
  return (sponsor.interval ?? 'month') === 'once'
    ? `$${amount.toLocaleString()} one time`
    : `$${amount.toLocaleString()} a month`
}
