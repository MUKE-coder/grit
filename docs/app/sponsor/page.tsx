'use client'

import { useState } from 'react'
import Link from 'next/link'
import { Heart, Check, Loader2, Star, Sparkles, ArrowRight } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { SiteHeader } from '@/components/site-header'
import {
  SPONSOR_TIERS,
  ONE_TIME_AMOUNTS,
  SPONSOR_MAX_USD,
  type SponsorTierId,
  type SponsorInterval,
} from '@/config/sponsors'

export default function SponsorPage() {
  const [interval, setInterval] = useState<SponsorInterval>('month')
  const [tierId, setTierId] = useState<SponsorTierId>('special')
  const [oneTime, setOneTime] = useState<number>(100)
  const [customAmount, setCustomAmount] = useState('')
  const [isCustom, setIsCustom] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const selectedTier = SPONSOR_TIERS.find((t) => t.id === tierId)!

  const amountUsd =
    interval === 'month'
      ? selectedTier.price
      : isCustom
        ? parseFloat(customAmount)
        : oneTime

  async function checkout() {
    setError('')

    if (interval === 'once' && isCustom) {
      const parsed = parseFloat(customAmount)
      if (!customAmount || isNaN(parsed) || parsed < 1) {
        setError('Please enter a valid amount ($1 minimum).')
        return
      }
      if (parsed > SPONSOR_MAX_USD) {
        setError(`Maximum is $${SPONSOR_MAX_USD}. Get in touch for larger sponsorships.`)
        return
      }
    }

    setLoading(true)
    try {
      const res = await fetch('/api/checkout', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          amount: Math.round(amountUsd * 100),
          interval,
          tier: interval === 'month' ? selectedTier.name : undefined,
        }),
      })
      const data = await res.json()
      if (!res.ok) {
        setError(data.error || 'Something went wrong.')
        return
      }
      window.location.href = data.url
    } catch {
      setError('Failed to connect. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />

      {/* Hero */}
      <section className="relative overflow-hidden border-b border-border/30">
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute left-1/2 top-0 h-[420px] w-[720px] -translate-x-1/2 rounded-full bg-primary/[0.07] blur-[120px]" />
          <div className="absolute inset-0 bg-grit-grid-sm mask-fade-y opacity-[0.35]" />
        </div>
        <div className="container max-w-screen-xl px-6 py-20 md:py-28">
          <div className="mx-auto max-w-3xl text-center">
            <div className="mx-auto mb-6 flex w-fit items-center gap-2 rounded-full border border-border/60 bg-primary/10 px-4 py-1.5">
              <Heart className="h-3.5 w-3.5 text-primary" />
              <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
                Become a sponsor
              </span>
            </div>
            <h1 className="mb-6 font-display text-4xl font-bold tracking-tight md:text-5xl lg:text-6xl">
              <span className="gradient-text">Sponsor Grit</span>
            </h1>
            <p className="mx-auto max-w-2xl text-lg leading-relaxed text-muted-foreground md:text-xl">
              Grit is free, open source, and built in the open. Sponsorship pays for
              the features, docs and releases — and puts your name in front of every
              developer who uses it.
            </p>
          </div>
        </div>
      </section>

      {/* Plans */}
      <section>
        <div className="container max-w-screen-xl px-6 py-16 md:py-20">
          {/* Interval toggle */}
          <div className="mb-10 flex justify-center">
            <div
              role="tablist"
              aria-label="Sponsorship frequency"
              className="inline-flex rounded-xl border border-border/50 bg-card/50 p-1"
            >
              {(
                [
                  { id: 'month', label: 'Monthly' },
                  { id: 'once', label: 'One-time' },
                ] as const
              ).map((opt) => (
                <button
                  key={opt.id}
                  role="tab"
                  aria-selected={interval === opt.id}
                  onClick={() => {
                    setInterval(opt.id)
                    setError('')
                  }}
                  className={cn(
                    'rounded-lg px-6 py-2 text-sm font-semibold transition-all',
                    interval === opt.id
                      ? 'bg-primary text-primary-foreground'
                      : 'text-muted-foreground hover:text-foreground'
                  )}
                >
                  {opt.label}
                </button>
              ))}
            </div>
          </div>

          {interval === 'month' ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
              {SPONSOR_TIERS.map((tier) => {
                const selected = tierId === tier.id
                return (
                  <button
                    key={tier.id}
                    onClick={() => {
                      setTierId(tier.id)
                      setError('')
                    }}
                    className={cn(
                      'card-grit relative flex flex-col rounded-xl border p-5 text-left transition-all',
                      selected
                        ? 'border-primary/50 bg-primary/[0.06] ring-1 ring-primary/25'
                        : 'border-border/40 bg-card/50 hover:border-primary/25'
                    )}
                  >
                    {tier.featured && (
                      <span className="absolute -top-2.5 right-4 flex items-center gap-1 rounded-full border border-amber-500/30 bg-amber-500/10 px-2 py-0.5 font-mono text-[10px] uppercase tracking-wider text-amber-400">
                        <Star className="h-2.5 w-2.5 fill-current" />
                        Most popular
                      </span>
                    )}
                    <h3 className="font-display text-lg font-bold text-foreground">
                      {tier.name}
                    </h3>
                    <div className="mt-2 flex items-baseline gap-1">
                      <span className="font-display text-3xl font-bold text-foreground">
                        ${tier.price}
                      </span>
                      <span className="text-sm text-muted-foreground">/month</span>
                    </div>
                    <p className="mt-2 text-xs leading-relaxed text-muted-foreground">
                      {tier.blurb}
                    </p>
                    <ul className="mt-4 flex-1 space-y-2 border-t border-border/30 pt-4">
                      {tier.perks.map((perk) => (
                        <li key={perk} className="flex items-start gap-2">
                          <Check className="mt-0.5 h-3.5 w-3.5 shrink-0 text-primary" />
                          <span className="text-xs text-foreground/80">{perk}</span>
                        </li>
                      ))}
                    </ul>
                    <span
                      className={cn(
                        'mt-5 block rounded-lg py-2 text-center text-sm font-semibold transition-colors',
                        selected
                          ? 'bg-primary text-primary-foreground'
                          : 'border border-border/50 text-foreground'
                      )}
                    >
                      {selected ? 'Selected' : 'Select'}
                    </span>
                  </button>
                )
              })}
            </div>
          ) : (
            <div className="mx-auto max-w-xl">
              <div className="mb-4 grid grid-cols-3 gap-3">
                {ONE_TIME_AMOUNTS.map((amount) => (
                  <button
                    key={amount}
                    onClick={() => {
                      setOneTime(amount)
                      setIsCustom(false)
                      setError('')
                    }}
                    className={cn(
                      'rounded-xl border p-5 text-center transition-all',
                      !isCustom && oneTime === amount
                        ? 'border-primary/50 bg-primary/[0.06] ring-1 ring-primary/25'
                        : 'border-border/40 bg-card/50 hover:border-primary/25'
                    )}
                  >
                    <span className="font-display text-2xl font-bold text-foreground">
                      ${amount}
                    </span>
                  </button>
                ))}
              </div>
              <button
                onClick={() => {
                  setIsCustom(true)
                  setError('')
                }}
                className={cn(
                  'w-full rounded-xl border p-4 text-left transition-all',
                  isCustom
                    ? 'border-primary/50 bg-primary/[0.06] ring-1 ring-primary/25'
                    : 'border-border/40 bg-card/50 hover:border-primary/25'
                )}
              >
                <span className="text-sm font-semibold text-foreground">Custom amount</span>
                {isCustom && (
                  <div className="mt-3 flex items-center gap-2">
                    <span className="text-lg font-bold text-muted-foreground">$</span>
                    <Input
                      type="number"
                      min="1"
                      max={SPONSOR_MAX_USD}
                      step="1"
                      placeholder="Enter amount"
                      value={customAmount}
                      onChange={(e) => setCustomAmount(e.target.value)}
                      onClick={(e) => e.stopPropagation()}
                      className="border-border/40 bg-background/50"
                      autoFocus
                    />
                  </div>
                )}
              </button>
              <p className="mt-4 text-center text-xs text-muted-foreground">
                One-time sponsors are listed on the{' '}
                <Link href="/sponsors" className="text-primary hover:underline">
                  Sponsors page
                </Link>{' '}
                too.
              </p>
            </div>
          )}

          {/* Checkout */}
          <div className="mx-auto mt-10 max-w-xl">
            {error && (
              <p className="mb-4 text-center text-sm text-destructive">{error}</p>
            )}
            <Button
              size="lg"
              className="glow-primary-sm h-12 w-full text-base"
              onClick={checkout}
              disabled={loading}
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Redirecting to Stripe…
                </>
              ) : (
                <>
                  <Heart className="mr-2 h-4 w-4" />
                  {interval === 'month'
                    ? `Sponsor $${selectedTier.price}/month`
                    : `Sponsor ${
                        isCustom && customAmount ? `$${customAmount}` : `$${oneTime}`
                      }`}
                </>
              )}
            </Button>
            <p className="mt-4 text-center text-xs text-muted-foreground/50">
              {interval === 'month'
                ? 'Cancel any time. Payments processed securely by Stripe — no card details touch our servers.'
                : 'Payments processed securely by Stripe — no card details touch our servers.'}
            </p>
          </div>
        </div>
      </section>

      {/* Where support goes */}
      <section className="border-t border-border/30 bg-accent/20">
        <div className="container max-w-screen-xl px-6 py-20">
          <div className="mb-12 text-center">
            <h2 className="mb-4 font-display text-3xl font-bold tracking-tight md:text-4xl">
              Where your support goes
            </h2>
            <p className="mx-auto max-w-2xl text-lg text-muted-foreground">
              Every contribution funds the work that makes Grit better for everyone.
            </p>
          </div>
          <div className="mx-auto grid max-w-3xl gap-4 sm:grid-cols-2">
            {[
              'New features and framework improvements',
              'Documentation, tutorials, and video courses',
              'Bug fixes and security patches',
              'Community support and issue triage',
              'Infrastructure and hosting costs',
              'Open-source sustainability',
            ].map((point) => (
              <div
                key={point}
                className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-3"
              >
                <Sparkles className="mt-0.5 h-4 w-4 shrink-0 text-primary" />
                <span className="text-sm text-foreground/80">{point}</span>
              </div>
            ))}
          </div>
          <div className="mt-10 text-center">
            <Link
              href="/sponsors"
              className="inline-flex items-center gap-2 text-sm font-semibold text-primary hover:underline"
            >
              See everyone who sponsors Grit
              <ArrowRight className="h-4 w-4" />
            </Link>
          </div>
        </div>
      </section>

      <footer className="border-t border-border/30 px-6 py-8">
        <div className="container flex max-w-screen-xl flex-col items-center justify-between gap-4 text-sm text-muted-foreground/50 sm:flex-row">
          <span>Grit Framework — Go + React. Built with Grit.</span>
          <div className="flex items-center gap-4">
            <Link href="/docs" className="transition-colors hover:text-foreground">
              Docs
            </Link>
            <Link href="/sponsors" className="transition-colors hover:text-foreground">
              Sponsors
            </Link>
            <Link href="/hire" className="transition-colors hover:text-foreground">
              Hire Us
            </Link>
            <Link
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noreferrer"
              className="transition-colors hover:text-foreground"
            >
              GitHub
            </Link>
          </div>
        </div>
      </footer>
    </div>
  )
}
