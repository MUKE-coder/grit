import { NextRequest, NextResponse } from 'next/server'
import Stripe from 'stripe'

const stripe = new Stripe(process.env.STRIPE_SECRET_KEY!)

export async function POST(req: NextRequest) {
  if (!process.env.STRIPE_SECRET_KEY) {
    return NextResponse.json(
      { error: 'Stripe is not configured.' },
      { status: 500 }
    )
  }

  try {
    const { amount } = await req.json()

    if (!amount || typeof amount !== 'number' || amount < 100 || !Number.isInteger(amount)) {
      return NextResponse.json(
        { error: 'Invalid amount. Minimum donation is $1.00.' },
        { status: 400 }
      )
    }

    if (amount > 99900) {
      return NextResponse.json(
        { error: 'Maximum donation is $999.00.' },
        { status: 400 }
      )
    }

    const session = await stripe.checkout.sessions.create({
      mode: 'payment',
      line_items: [
        {
          price_data: {
            currency: 'usd',
            product_data: {
              name: 'Support Grit Framework',
              description: 'Thank you for supporting open-source development!',
            },
            unit_amount: amount,
          },
          quantity: 1,
        },
      ],
      success_url: `${req.nextUrl.origin}/donate/success`,
      cancel_url: `${req.nextUrl.origin}/donate`,
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
