'use client'

import { useState } from 'react'
import Link from 'next/link'
import { ArrowRight, AlertTriangle, Check, X } from 'lucide-react'
import { DEPLOYMENT_PROVIDERS } from '@/config/deployment-providers'

/**
 * Where to deploy, with the steps for each.
 *
 * Reads DEPLOYMENT_PROVIDERS rather than restating it. That config already
 * drives the /docs/deployment/[provider] pages, and a second hand-written copy
 * on the homepage is how the two drift until one is wrong — usually the one
 * nobody re-reads.
 *
 * Only the first few steps are shown per provider; the full walkthrough lives
 * on the docs page and is linked. That is a deliberate cap, and it is stated
 * in the UI rather than silently truncating.
 */

const STEP_PREVIEW = 3

const EFFORT_TONE: Record<string, string> = {
  Lowest: 'text-emerald-400',
  Low: 'text-emerald-400',
  Medium: 'text-amber-400',
  High: 'text-rose-400',
}

function Capability({ ok, label }: { ok: boolean; label: string }) {
  return (
    <div className="flex items-center gap-1.5">
      {ok ? (
        <Check className="h-3.5 w-3.5 text-emerald-400" strokeWidth={2.5} />
      ) : (
        <X className="h-3.5 w-3.5 text-muted-foreground/50" strokeWidth={2.5} />
      )}
      <span className={ok ? 'text-foreground/80' : 'text-muted-foreground/60'}>{label}</span>
    </div>
  )
}

export function DeployShowcase() {
  const [active, setActive] = useState(DEPLOYMENT_PROVIDERS[0].slug)
  const provider =
    DEPLOYMENT_PROVIDERS.find((p) => p.slug === active) ?? DEPLOYMENT_PROVIDERS[0]
  const shown = provider.steps.slice(0, STEP_PREVIEW)
  const hidden = provider.steps.length - shown.length

  return (
    <div>
      <div role="tablist" aria-label="Hosting provider" className="flex flex-wrap gap-2 mb-8">
        {DEPLOYMENT_PROVIDERS.map((p) => {
          const selected = p.slug === active
          return (
            <button
              key={p.slug}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(p.slug)}
              className={`rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              {p.name}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_19rem] gap-8 lg:gap-10 items-start">
        <div>
          <h3 className="text-lg font-semibold mb-2 leading-snug">{provider.tagline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-6">
            <span className="text-foreground/80">Good fit:</span> {provider.bestFor}
          </p>

          <ol className="space-y-4">
            {shown.map((step, i) => (
              <li key={step.title} className="flex gap-3">
                <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-medium text-primary">
                  {i + 1}
                </span>
                <div className="min-w-0 flex-1">
                  <p className="text-[13.5px] font-medium text-foreground mb-1">{step.title}</p>
                  <p className="text-[12.5px] text-muted-foreground leading-relaxed mb-2">
                    {step.body}
                  </p>
                  {step.code && (
                    <div className="overflow-x-auto rounded-lg border border-border/50 bg-background/60 p-3">
                      <pre className="text-[11.5px] font-mono text-foreground/85 leading-relaxed whitespace-pre">
                        {step.code.code}
                      </pre>
                    </div>
                  )}
                </div>
              </li>
            ))}
          </ol>

          <div className="mt-6 flex flex-wrap items-center gap-4">
            <Link
              href={`/docs/deployment/${provider.slug}`}
              className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
            >
              {hidden > 0
                ? `The remaining ${hidden} step${hidden === 1 ? '' : 's'}`
                : 'Full walkthrough'}
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>
            <a
              href={provider.docsUrl}
              target="_blank"
              rel="noreferrer"
              className="text-sm text-muted-foreground hover:text-foreground"
            >
              {provider.name} docs
            </a>
          </div>
        </div>

        <div className="space-y-4">
          <div className="rounded-xl border border-border/60 bg-card/40 p-4">
            <div className="grid grid-cols-2 gap-y-3 text-[12.5px] mb-4">
              <span className="text-muted-foreground">From</span>
              <span className="text-foreground/90 text-right">{provider.costFrom}</span>
              <span className="text-muted-foreground">Ops effort</span>
              <span className={`text-right ${EFFORT_TONE[provider.effort] ?? 'text-foreground/90'}`}>
                {provider.effort}
              </span>
            </div>
            <div className="space-y-2 border-t border-border/40 pt-3 text-[12.5px]">
              <Capability ok={provider.managedPostgres} label="Managed Postgres" />
              <Capability ok={provider.managedRedis} label="Managed Redis" />
              <Capability ok={provider.persistentDisk} label="Persistent disk" />
            </div>
          </div>

          {provider.gotchas[0] && (
            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.07] p-4">
              <div className="flex items-center gap-2 mb-2">
                <AlertTriangle className="h-3.5 w-3.5 text-amber-500 shrink-0" />
                <span className="text-[11.5px] font-semibold text-amber-200/90">
                  What catches people out
                </span>
              </div>
              <p className="text-[12px] text-muted-foreground leading-relaxed">
                {provider.gotchas[0]}
              </p>
            </div>
          )}

          <div className="rounded-xl border border-border/60 bg-card/40 p-4">
            <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-2">
              Not for
            </div>
            <p className="text-[12px] text-muted-foreground leading-relaxed">{provider.notFor}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
