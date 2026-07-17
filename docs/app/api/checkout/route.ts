import { NextRequest, NextResponse } from 'next/server'
import Stripe from 'stripe'

/**
 * Sponsorship + donation checkout.
 *
 * Two modes:
 *   - one-time  → Stripe `mode: 'payment'`     (the original /donate flow)
 *   - monthly   → Stripe `mode: 'subscription'` (recurring sponsor tiers)
 *
 * Prices are built inline with `price_data` rather than pre-created Stripe
 * Prices, so changing the ladder in config/sponsors.ts needs no dashboard work.
 */

const MIN_AMOUNT = 100 // $1.00
const MAX_AMOUNT = 99900 // $999.00

export async function POST(req: NextRequest) {
  if (!process.env.STRIPE_SECRET_KEY) {
    return NextResponse.json({ error: 'Stripe is not configured.' }, { status: 500 })
  }

  const stripe = new Stripe(process.env.STRIPE_SECRET_KEY)

  try {
    const body = await req.json()
    const amount: unknown = body.amount
    const interval: unknown = body.interval ?? 'once'
    const tierName: unknown = body.tier

    if (interval !== 'once' && interval !== 'month') {
      return NextResponse.json(
        { error: "Invalid interval. Expected 'once' or 'month'." },
        { status: 400 }
      )
    }

    if (
      !amount ||
      typeof amount !== 'number' ||
      !Number.isInteger(amount) ||
      amount < MIN_AMOUNT
    ) {
      return NextResponse.json(
        { error: 'Invalid amount. Minimum is $1.00.' },
        { status: 400 }
      )
    }

    if (amount > MAX_AMOUNT) {
      return NextResponse.json(
        { error: 'Maximum is $999.00. Get in touch for larger sponsorships.' },
        { status: 400 }
      )
    }

    const recurring = interval === 'month'
    const label =
      typeof tierName === 'string' && tierName.trim()
        ? `Grit ${tierName.trim()} Sponsor`
        : recurring
          ? 'Grit Monthly Sponsor'
          : 'Support Grit Framework'

    const session = await stripe.checkout.sessions.create({
      mode: recurring ? 'subscription' : 'payment',
      line_items: [
        {
          price_data: {
            currency: 'usd',
            product_data: {
              name: label,
              description: recurring
                ? 'Recurring sponsorship — cancel any time.'
                : 'Thank you for supporting open-source development!',
            },
            unit_amount: amount,
            ...(recurring ? { recurring: { interval: 'month' as const } } : {}),
          },
          quantity: 1,
        },
      ],
      // Needed to email the sponsor and to know who to list on /sponsors.
      billing_address_collection: 'auto',
      success_url: `${req.nextUrl.origin}/donate/success`,
      cancel_url: `${req.nextUrl.origin}/sponsor`,
    })

    return NextResponse.json({ url: session.url })
  } catch (err) {
    console.error('Stripe checkout error:', err)
    return NextResponse.json(
      { error: 'Failed to create checkout session.' },
      { status: 500 }
    )
  }
}
