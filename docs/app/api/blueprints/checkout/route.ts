import { NextRequest, NextResponse } from 'next/server'
import Stripe from 'stripe'

import { getBlueprint } from '@/config/blueprints'
import { PLANS, isPlanId, priceCents, type PlanId } from '@/config/blueprint-plans'

/**
 * Blueprint checkout.
 *
 * Four things can be bought: one blueprint ($10, forever), or a membership to
 * all of them monthly, yearly or for life. The plan id arrives from the
 * browser; the AMOUNT NEVER DOES. It is looked up in config/blueprint-plans.ts
 * on the server, because the sponsorship route next door takes the amount from
 * the request body (correct for a donation, where the payer picks the number)
 * and copying that here would let anybody buy a lifetime membership for $1.
 *
 * Access is a GitHub repository invitation, so checkout asks for the buyer's
 * GitHub username as a required custom field. The webhook route reads it and
 * sends the invitation. Nothing is stored here: Stripe holds the record of who
 * paid for what, which is why there is no database behind any of this.
 */

// GitHub's limit: 39 characters, alphanumeric and single hyphens.
const GITHUB_USERNAME_MAX = 39

export const GITHUB_FIELD_KEY = 'github'

export async function POST(req: NextRequest) {
  if (!process.env.STRIPE_SECRET_KEY) {
    return NextResponse.json(
      { error: 'Checkout is not configured on this host.' },
      { status: 500 }
    )
  }

  const stripe = new Stripe(process.env.STRIPE_SECRET_KEY)

  try {
    const body = await req.json()
    const planId: unknown = body.plan
    const blueprintName: unknown = body.blueprint

    if (!isPlanId(planId)) {
      return NextResponse.json(
        { error: 'Unknown plan. Expected single, monthly, yearly or lifetime.' },
        { status: 400 }
      )
    }

    const plan = PLANS[planId as PlanId]

    // A single purchase has to say WHICH blueprint, and it has to be one that
    // is actually for sale: an open blueprint is free and a planned one does
    // not exist, so neither can be bought.
    let label = 'Grit blueprints, every one'
    let description = plan.blurb
    let cancelPath = '/blueprints#pricing'

    if (plan.id === 'single') {
      if (typeof blueprintName !== 'string' || !blueprintName) {
        return NextResponse.json(
          { error: 'Which blueprint? A single purchase needs a name.' },
          { status: 400 }
        )
      }
      const blueprint = getBlueprint(blueprintName)
      if (!blueprint) {
        return NextResponse.json(
          { error: `There is no blueprint called ${blueprintName}.` },
          { status: 400 }
        )
      }
      if (blueprint.access === 'open') {
        return NextResponse.json(
          { error: `${blueprint.title} is open source. Clone it instead.` },
          { status: 400 }
        )
      }
      if (blueprint.access === 'planned') {
        return NextResponse.json(
          { error: `${blueprint.title} is not built yet.` },
          { status: 400 }
        )
      }
      label = `${blueprint.title}, a Grit blueprint`
      description = 'Permanent access to the private repository.'
      cancelPath = `/blueprints/${blueprint.name}`
    }

    const recurring = plan.mode === 'subscription'
    const metadata: Record<string, string> = { plan: plan.id }
    if (plan.id === 'single' && typeof blueprintName === 'string') {
      metadata.blueprint = blueprintName
    }

    const session = await stripe.checkout.sessions.create({
      mode: plan.mode,
      line_items: [
        {
          price_data: {
            currency: 'usd',
            product_data: { name: label, description },
            unit_amount: priceCents(plan.id),
            ...(recurring && plan.interval
              ? { recurring: { interval: plan.interval } }
              : {}),
          },
          quantity: 1,
        },
      ],
      // The invitation goes to a GitHub account, so the username is required
      // rather than collected afterwards by email. Getting it here means the
      // webhook can send the invitation without anybody being in the loop.
      custom_fields: [
        {
          key: GITHUB_FIELD_KEY,
          label: { type: 'custom', custom: 'Your GitHub username' },
          type: 'text',
          optional: false,
          text: { minimum_length: 1, maximum_length: GITHUB_USERNAME_MAX },
        },
      ],
      custom_text: {
        submit: {
          message:
            'The repository invitation goes to this GitHub account, usually within a minute.',
        },
      },
      // Needed to reach the buyer if the username turns out to be wrong.
      billing_address_collection: 'auto',
      ...(recurring
        ? { subscription_data: { metadata } }
        : { customer_creation: 'always' as const }),
      metadata,
      success_url: `${req.nextUrl.origin}/blueprints/access?session_id={CHECKOUT_SESSION_ID}`,
      cancel_url: `${req.nextUrl.origin}${cancelPath}`,
    })

    return NextResponse.json({ url: session.url })
  } catch (err) {
    console.error('Blueprint checkout error:', err)
    return NextResponse.json(
      { error: 'Failed to start checkout. Please try again.' },
      { status: 500 }
    )
  }
}
