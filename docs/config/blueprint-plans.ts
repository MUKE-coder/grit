/*
 * What a blueprint costs, in one place.
 *
 * Read by three things that must agree: the pricing UI, the Stripe checkout
 * route (which takes the amount from HERE and never from the request body),
 * and internal/blueprint/registry.go, which prints the same numbers in
 * `grit blueprint info`. TestBlueprintPricesMatchTheWebsite compares the Go
 * constants with this file, because two writers with nothing comparing them
 * drift, and a price that drifts is a price somebody disputes.
 *
 * Deliberately NOT in blueprints.ts: a Go test scans that file for every
 * `name: '...'` and fails on any value that is not a blueprint, so a plan
 * object with a `name` field there would break it in a confusing way.
 */

export type PlanId = 'single' | 'monthly' | 'yearly' | 'lifetime'

export interface Plan {
  id: PlanId
  /** What it is called on the page. */
  label: string
  priceUsd: number
  /** Rendered after the price: 'once', 'a month', 'a year'. */
  cadence: string
  /** One line: what the money buys. */
  blurb: string
  /** Stripe checkout mode. */
  mode: 'payment' | 'subscription'
  /** Only for mode 'subscription'. */
  interval?: 'month' | 'year'
}

export const SINGLE_USD = 10
export const MEMBERSHIP_MONTHLY_USD = 20
export const MEMBERSHIP_YEARLY_USD = 200
export const LIFETIME_USD = 400

export const PLANS: Record<PlanId, Plan> = {
  single: {
    id: 'single',
    label: 'One blueprint',
    priceUsd: SINGLE_USD,
    cadence: 'once',
    blurb: 'Permanent access to that one repository.',
    mode: 'payment',
  },
  monthly: {
    id: 'monthly',
    label: 'Membership',
    priceUsd: MEMBERSHIP_MONTHLY_USD,
    cadence: 'a month',
    blurb: 'Every blueprint, including the ones added while you are a member.',
    mode: 'subscription',
    interval: 'month',
  },
  yearly: {
    id: 'yearly',
    label: 'Membership, yearly',
    priceUsd: MEMBERSHIP_YEARLY_USD,
    cadence: 'a year',
    blurb: 'The same thing for the price of ten months.',
    mode: 'subscription',
    interval: 'year',
  },
  lifetime: {
    id: 'lifetime',
    label: 'Lifetime',
    priceUsd: LIFETIME_USD,
    cadence: 'once',
    blurb: 'Every blueprint there is, and every one there will be. Nothing to renew.',
    mode: 'payment',
  },
}

/** The three ways to get all of them, cheapest first. */
export const MEMBERSHIP_PLANS: Plan[] = [PLANS.monthly, PLANS.yearly, PLANS.lifetime]

/** Stripe wants cents, and wants them as an integer. */
export function priceCents(id: PlanId): number {
  return PLANS[id].priceUsd * 100
}

/** $10, $200. No trailing .00, because every price here is whole dollars. */
export function usd(amount: number): string {
  return '$' + amount.toLocaleString('en-US')
}

export function isPlanId(value: unknown): value is PlanId {
  return typeof value === 'string' && value in PLANS
}
