import Link from 'next/link'
import type { Metadata } from 'next'
import Stripe from 'stripe'
import { ArrowRight, Check, Github, Mail } from 'lucide-react'

import { SiteHeader } from '@/components/site-header'
import { Button } from '@/components/grit/button'
import { Section, Highlight } from '@/components/grit/section'
import { BLUEPRINTS, getBlueprint, PAID_BLUEPRINTS } from '@/config/blueprints'
import { isPlanId, PLANS } from '@/config/blueprint-plans'

/*
 * Where Stripe sends somebody after they pay.
 *
 * It reads the checkout session to say exactly what was bought and which
 * GitHub account the invitation went to, because "thank you for your
 * purchase" answers neither of the two questions somebody has at this moment:
 * did it work, and where do I go now.
 *
 * The invitation itself is the webhook's job, not this page's. This page may
 * well render before that event arrives, which is why it says "within a
 * minute" rather than claiming the invitation has already been sent.
 */

export const dynamic = 'force-dynamic'

export const metadata: Metadata = {
  title: 'Your blueprint access',
  description: 'What you bought, and where the repository invitation went.',
  robots: { index: false, follow: false },
}

interface Purchase {
  planLabel: string
  repos: { title: string; url: string }[]
  username: string | null
  email: string | null
}

async function readSession(id: string | undefined): Promise<Purchase | null> {
  if (!id || !process.env.STRIPE_SECRET_KEY) return null
  try {
    const stripe = new Stripe(process.env.STRIPE_SECRET_KEY)
    const session = await stripe.checkout.sessions.retrieve(id)
    const plan = session.metadata?.plan
    if (!isPlanId(plan)) return null

    const single = plan === 'single' ? getBlueprint(session.metadata?.blueprint ?? '') : null
    const covered = single ? [single] : PAID_BLUEPRINTS

    return {
      planLabel: single ? single.title : PLANS[plan].label,
      repos: covered
        .filter((b) => b.repo)
        .map((b) => ({ title: b.title, url: b.repo as string })),
      username:
        session.custom_fields?.find((f) => f.key === 'github')?.text?.value ?? null,
      email: session.customer_details?.email ?? null,
    }
  } catch {
    // A stale or forged session id lands here. The page still renders, with
    // the generic copy: it is a confirmation, not an entitlement check.
    return null
  }
}

export default async function BlueprintAccessPage({
  searchParams,
}: {
  searchParams: Promise<{ session_id?: string }>
}) {
  const { session_id } = await searchParams
  const purchase = await readSession(session_id)

  const repos =
    purchase?.repos ??
    BLUEPRINTS.filter((b) => b.access === 'paid' && b.repo).map((b) => ({
      title: b.title,
      url: b.repo as string,
    }))

  return (
    <div className="min-h-screen bg-bg">
      <SiteHeader />

      <main>
        <Section first grid wash="right">
          <div className="mx-auto max-w-2xl">
            <span className="g-badge g-badge--accent font-mono text-[11px]">
              <Check className="h-3.5 w-3.5" aria-hidden="true" />
              Payment received
            </span>

            <h1 className="mt-6 text-[clamp(2.25rem,4.5vw,3rem)] font-semibold leading-[1.07] tracking-[-0.04em] text-fg [text-wrap:balance]">
              {purchase ? (
                <>
                  {purchase.planLabel} is <Highlight>yours.</Highlight>
                </>
              ) : (
                <>
                  Thank you. <Highlight>Check GitHub.</Highlight>
                </>
              )}
            </h1>

            <p className="mt-5 text-lead text-fg-muted">
              {purchase?.username ? (
                <>
                  The repository invitation is on its way to{' '}
                  <span className="font-mono text-fg">{purchase.username}</span>, usually
                  within a minute. Accept it and the code is yours to clone.
                </>
              ) : (
                <>
                  The repository invitation is on its way to the GitHub account you gave
                  at checkout, usually within a minute. Accept it and the code is yours
                  to clone.
                </>
              )}
            </p>

            <div className="mt-10 space-y-3">
              <p className="g-eyebrow">What to do now</p>
              {[
                {
                  icon: Github,
                  title: 'Accept the invitation',
                  body: 'It appears at github.com/notifications and in an email from GitHub. An unaccepted invitation expires after seven days.',
                  href: 'https://github.com/notifications',
                  cta: 'Open GitHub notifications',
                },
                {
                  icon: ArrowRight,
                  title: 'Run it',
                  body: 'Each blueprint page has the exact commands, what each one does, and the account to sign in as.',
                  href: '/blueprints',
                  cta: 'The blueprint pages',
                },
              ].map((step) => (
                <div key={step.title} className="g-card gap-2">
                  <div className="flex items-center gap-3">
                    <span className="g-card__icon">
                      <step.icon aria-hidden="true" />
                    </span>
                    <h2 className="g-card__title m-0">{step.title}</h2>
                  </div>
                  <p className="g-card__body">{step.body}</p>
                  <div className="mt-2">
                    <Button href={step.href} variant="secondary" size="sm">
                      {step.cta}
                    </Button>
                  </div>
                </div>
              ))}
            </div>

            {repos.length > 0 && (
              <div className="mt-10">
                <p className="g-eyebrow mb-3">
                  {purchase && purchase.repos.length === 1
                    ? 'The repository'
                    : 'The repositories'}
                </p>
                <ul className="m-0 list-none space-y-2 p-0">
                  {repos.map((repo) => (
                    <li key={repo.url}>
                      <a
                        href={repo.url}
                        target="_blank"
                        rel="noreferrer"
                        className="g-link font-mono text-body-sm"
                      >
                        {repo.url.replace('https://github.com/', '')}
                      </a>
                      <span className="ml-2 text-body-sm text-fg-subtle">
                        {repo.title}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            )}

            <div className="mt-12 border-t border-border pt-8">
              <p className="m-0 flex items-start gap-2.5 text-body-sm text-fg-muted">
                <Mail className="mt-0.5 h-4 w-4 flex-none text-fg-subtle" aria-hidden="true" />
                <span>
                  Nothing after five minutes, or a typo in the username? Send the
                  address you paid with
                  {purchase?.email ? (
                    <>
                      {' '}
                      (<span className="font-mono text-fg">{purchase.email}</span>)
                    </>
                  ) : null}{' '}
                  to{' '}
                  <a href="mailto:desishub.info@gmail.com" className="g-link">
                    desishub.info@gmail.com
                  </a>{' '}
                  and it gets fixed by hand.{' '}
                  <Link href="/blueprints" className="g-link">
                    Back to the blueprints
                  </Link>
                  .
                </span>
              </p>
            </div>
          </div>
        </Section>
      </main>
    </div>
  )
}
