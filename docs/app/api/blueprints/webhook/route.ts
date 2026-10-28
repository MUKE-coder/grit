import { NextRequest, NextResponse } from 'next/server'
import Stripe from 'stripe'

import { getBlueprint } from '@/config/blueprints'
import { isPlanId } from '@/config/blueprint-plans'
import {
  cleanUsername,
  membershipTargets,
  parseRepo,
  worthRetrying,
  type Target,
} from '@/lib/blueprint-access'

/**
 * What happens after the money: the GitHub invitation.
 *
 * Stripe is the only record. There is no database here, and there does not
 * need to be one, because every question this code asks ("did they pay?",
 * "is the membership still live?") is a question Stripe already answers, and
 * both GitHub calls are idempotent: adding a collaborator who is already one
 * changes nothing, and removing one who is already gone returns 404, which is
 * the state we wanted. So a redelivered event is harmless, which is the
 * property that lets the whole thing be stateless.
 *
 * Two events matter:
 *
 *   checkout.session.completed    invite to one repository, or to all of them
 *   customer.subscription.deleted a membership ended, so take the access back
 *
 * A one-off purchase is never revoked. That is what "forever" means on the
 * page, and it is why the single and the membership are different events.
 */

export const dynamic = 'force-dynamic'

const GITHUB_API = 'https://api.github.com'

function usernameFrom(session: Stripe.Checkout.Session): string | null {
  const field = session.custom_fields?.find((f) => f.key === 'github')
  return cleanUsername(field?.text?.value)
}

async function invite(target: Target, username: string, token: string) {
  const res = await fetch(
    `${GITHUB_API}/repos/${target.owner}/${target.repo}/collaborators/${username}`,
    {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: 'application/vnd.github+json',
        'X-GitHub-Api-Version': '2022-11-28',
      },
      body: JSON.stringify({ permission: 'pull' }),
    }
  )
  // 201 invitation created, 204 already a collaborator. Both are success.
  return { ok: res.ok, status: res.status, body: res.ok ? '' : await res.text() }
}

async function revoke(target: Target, username: string, token: string) {
  const res = await fetch(
    `${GITHUB_API}/repos/${target.owner}/${target.repo}/collaborators/${username}`,
    {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${token}`,
        Accept: 'application/vnd.github+json',
        'X-GitHub-Api-Version': '2022-11-28',
      },
    }
  )
  // 204 removed, 404 was not there. Both leave them without access.
  const ok = res.ok || res.status === 404
  return { ok, status: res.status, body: ok ? '' : await res.text() }
}

export async function POST(req: NextRequest) {
  const secret = process.env.STRIPE_WEBHOOK_SECRET
  const stripeKey = process.env.STRIPE_SECRET_KEY
  const githubToken = process.env.BLUEPRINT_GITHUB_TOKEN

  if (!secret || !stripeKey) {
    console.error('blueprint webhook: Stripe is not configured on this host.')
    return NextResponse.json({ error: 'Not configured.' }, { status: 500 })
  }

  const signature = req.headers.get('stripe-signature')
  if (!signature) {
    return NextResponse.json({ error: 'No signature.' }, { status: 400 })
  }

  const stripe = new Stripe(stripeKey)

  // The RAW body, before any JSON parsing: the signature is over these exact
  // bytes, and re-serialising a parsed object does not reproduce them.
  const payload = await req.text()

  let event: Stripe.Event
  try {
    event = stripe.webhooks.constructEvent(payload, signature, secret)
  } catch (err) {
    console.error('blueprint webhook: bad signature', err)
    return NextResponse.json({ error: 'Bad signature.' }, { status: 400 })
  }

  if (!githubToken) {
    // Loudly, and with a 500, so Stripe keeps the event and retries once the
    // token is there. Swallowing this would lose a paid invitation.
    console.error('blueprint webhook: BLUEPRINT_GITHUB_TOKEN is not set.')
    return NextResponse.json({ error: 'Not configured.' }, { status: 500 })
  }

  try {
    switch (event.type) {
      case 'checkout.session.completed': {
        const session = event.data.object as Stripe.Checkout.Session
        const plan = session.metadata?.plan
        if (!isPlanId(plan)) {
          console.error('blueprint webhook: session has no known plan', session.id)
          return NextResponse.json({ received: true })
        }

        const username = usernameFrom(session)
        if (!username) {
          // Nothing to retry: the value is wrong and will stay wrong. Logged
          // so it can be fixed by hand, and acknowledged so Stripe stops.
          console.error(
            'blueprint webhook: no usable GitHub username on session',
            session.id,
            session.customer_details?.email
          )
          return NextResponse.json({ received: true })
        }

        let targets: Target[] = []
        if (plan === 'single') {
          const blueprint = getBlueprint(session.metadata?.blueprint ?? '')
          const target = parseRepo(blueprint?.repo)
          if (!target) {
            console.error(
              'blueprint webhook: no repository for',
              session.metadata?.blueprint
            )
            return NextResponse.json({ received: true })
          }
          targets = [target]
        } else {
          targets = membershipTargets()
        }

        // The subscription carries the username from here on, because
        // cancellation arrives as a subscription event with no session
        // attached and no other way to know whose access to take back.
        if (typeof session.subscription === 'string') {
          await stripe.subscriptions.update(session.subscription, {
            metadata: { plan, github: username },
          })
        }

        const retry: string[] = []
        for (const target of targets) {
          const result = await invite(target, username, githubToken)
          if (result.ok) {
            console.log(
              `blueprint webhook: invited ${username} to ${target.owner}/${target.repo}`
            )
            continue
          }
          console.error(
            `blueprint webhook: could not invite ${username} to ${target.repo}`,
            result.status,
            result.body
          )
          if (worthRetrying(result.status)) retry.push(target.repo)
        }

        if (retry.length > 0) {
          return NextResponse.json(
            { error: `GitHub was unavailable for ${retry.join(', ')}` },
            { status: 500 }
          )
        }
        return NextResponse.json({ received: true })
      }

      case 'customer.subscription.deleted': {
        const subscription = event.data.object as Stripe.Subscription
        const username = cleanUsername(subscription.metadata?.github)
        if (!username) {
          console.error(
            'blueprint webhook: cancelled subscription has no username',
            subscription.id
          )
          return NextResponse.json({ received: true })
        }

        const retry: string[] = []
        for (const target of membershipTargets()) {
          const result = await revoke(target, username, githubToken)
          if (result.ok) {
            console.log(
              `blueprint webhook: removed ${username} from ${target.owner}/${target.repo}`
            )
            continue
          }
          console.error(
            `blueprint webhook: could not remove ${username} from ${target.repo}`,
            result.status,
            result.body
          )
          if (worthRetrying(result.status)) retry.push(target.repo)
        }

        if (retry.length > 0) {
          return NextResponse.json(
            { error: `GitHub was unavailable for ${retry.join(', ')}` },
            { status: 500 }
          )
        }
        return NextResponse.json({ received: true })
      }

      default:
        return NextResponse.json({ received: true })
    }
  } catch (err) {
    console.error('blueprint webhook: failed', event.type, err)
    // A 500 keeps the event in Stripe's queue, which is what we want for
    // anything unexpected: a dropped event is a paid invitation that never
    // arrives, and nobody finds out until the buyer emails.
    return NextResponse.json({ error: 'Handler failed.' }, { status: 500 })
  }
}
