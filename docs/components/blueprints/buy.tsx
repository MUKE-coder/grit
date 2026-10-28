'use client'

import { useState } from 'react'
import { Check, Loader2 } from 'lucide-react'

import { Button, type ButtonVariant } from '@/components/grit/button'
import {
  MEMBERSHIP_PLANS,
  PLANS,
  usd,
  type PlanId,
} from '@/config/blueprint-plans'

/*
 * The buy buttons.
 *
 * They send a plan id and, for a single purchase, a blueprint name. They do
 * not send a price: the server looks that up. See the checkout route for why
 * that is not an oversight.
 */

function useCheckout() {
  const [pending, setPending] = useState<PlanId | null>(null)
  const [error, setError] = useState('')

  async function start(plan: PlanId, blueprint?: string) {
    setError('')
    setPending(plan)
    try {
      const res = await fetch('/api/blueprints/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ plan, blueprint }),
      })
      const data = await res.json()
      if (!res.ok || !data.url) {
        setError(data.error || 'Checkout could not be started.')
        setPending(null)
        return
      }
      window.location.href = data.url
    } catch {
      setError('Failed to reach the server. Please try again.')
      setPending(null)
    }
  }

  return { pending, error, start }
}

export function BuyButton({
  plan,
  blueprint,
  children,
  variant = 'primary',
  size = 'lg',
}: {
  plan: PlanId
  blueprint?: string
  children: React.ReactNode
  variant?: ButtonVariant
  size?: 'sm' | 'md' | 'lg'
}) {
  const { pending, error, start } = useCheckout()
  const busy = pending === plan

  return (
    <div className="flex flex-col gap-2">
      <Button
        variant={variant}
        size={size}
        className="disabled:cursor-not-allowed disabled:opacity-70"
        onClick={() => start(plan, blueprint)}
        disabled={busy}
        aria-busy={busy}
      >
        {busy && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
        {children}
      </Button>
      {error && (
        <p role="alert" className="m-0 text-body-sm text-danger">
          {error}
        </p>
      )}
    </div>
  )
}

/**
 * The three ways to get all of them.
 *
 * Yearly is marked rather than monthly, because it is the one that is better
 * value and saying so is more useful than making the cheapest number look
 * like the recommendation.
 */
export function MembershipPlans({ className = '' }: { className?: string }) {
  const { pending, error, start } = useCheckout()

  return (
    <div className={className}>
      <ul className="m-0 grid list-none grid-cols-1 gap-4 p-0 lg:grid-cols-3">
        {MEMBERSHIP_PLANS.map((plan) => {
          const best = plan.id === 'yearly'
          return (
            <li key={plan.id}>
              <div
                className={`g-card h-full ${best ? 'border-accent-line' : ''}`}
                data-recommended={best || undefined}
              >
                <div className="flex items-center justify-between gap-3">
                  <h3 className="g-card__title">{plan.label}</h3>
                  {best && (
                    <span className="g-badge g-badge--accent shrink-0 font-mono text-[11px]">
                      Two months free
                    </span>
                  )}
                </div>

                <p className="m-0 flex items-baseline gap-1.5">
                  <span className="text-[32px] font-semibold leading-none tracking-[-0.03em] text-fg">
                    {usd(plan.priceUsd)}
                  </span>
                  <span className="text-body-sm text-fg-subtle">{plan.cadence}</span>
                </p>

                <p className="g-card__body">{plan.blurb}</p>

                <div className="mt-auto pt-5">
                  <Button
                    variant={best ? 'primary' : 'secondary'}
                    className="w-full justify-center disabled:cursor-not-allowed disabled:opacity-70"
                    onClick={() => start(plan.id)}
                    disabled={pending === plan.id}
                    aria-busy={pending === plan.id}
                  >
                    {pending === plan.id && (
                      <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
                    )}
                    {plan.mode === 'subscription' ? 'Join' : 'Buy it once'}
                  </Button>
                </div>
              </div>
            </li>
          )
        })}
      </ul>

      {error && (
        <p role="alert" className="mt-4 text-body-sm text-danger">
          {error}
        </p>
      )}

      <ul className="mt-8 grid list-none grid-cols-1 gap-x-8 gap-y-3 p-0 text-body-sm text-fg-muted sm:grid-cols-2">
        {[
          `Or ${usd(PLANS.single.priceUsd)} for a single blueprint, kept permanently.`,
          'The invitation reaches your GitHub account within a minute.',
          'A membership covers every blueprint added while it runs.',
          'Cancel whenever. What you have already cloned stays yours.',
        ].map((line) => (
          <li key={line} className="flex gap-2.5">
            <Check
              className="mt-0.5 h-4 w-4 flex-none text-accent"
              aria-hidden="true"
            />
            {line}
          </li>
        ))}
      </ul>
    </div>
  )
}
