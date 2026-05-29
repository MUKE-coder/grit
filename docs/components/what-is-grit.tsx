'use client'

// "What is Grit — in 60 seconds." A clear, low-jargon explainer that
// lives between the hero and the deep technical sections, so a
// first-time visitor never has to guess what Grit is or what it
// produces. Three numbered steps + a tiny terminal showing the exact
// sequence of commands you run to get a working app.

import { Sparkles, Wand2, Rocket, ArrowRight, Check } from 'lucide-react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'

const STEPS = [
  {
    n: '01',
    icon: Sparkles,
    title: 'Scaffold a real project',
    body: 'Pick your architecture (single binary, monorepo, mobile, desktop) and one command scaffolds the entire repo — Go API, React frontend, admin panel, Docker, CI, all wired together.',
    code: 'grit new my-app --triple',
  },
  {
    n: '02',
    icon: Wand2,
    title: 'Generate full-stack features',
    body: "Describe a resource (Product · Order · Invoice). Grit emits the Go model, service, handler, routes, Zod schema, React hooks, and admin page — typed end-to-end, no glue code.",
    code: 'grit generate resource Product',
  },
  {
    n: '03',
    icon: Rocket,
    title: 'Ship to a real server',
    body: "grit deploy cross-compiles, uploads, sets up systemd, and configures Caddy with auto-TLS. Or push to git and let your platform of choice run the bundled Dockerfile.",
    code: 'grit deploy --domain acme.app',
  },
]

const HIGHLIGHTS = [
  'No glue code between backend and frontend — generated together, always in sync',
  'Production batteries: auth + 2FA, OAuth, storage, jobs, AI, observability',
  'Secure-by-default headers; OWASP Top 10:2025 hardened out of the box',
]

export function WhatIsGrit() {
  return (
    <section className="relative py-24 px-6 overflow-hidden border-t border-border/40">
      {/* Layered backdrop */}
      <div className="absolute inset-0 -z-30 bg-gradient-to-b from-card/30 via-background to-card/20" />
      <div
        className="absolute inset-0 -z-20 opacity-[0.35]"
        style={{
          backgroundImage:
            'radial-gradient(circle at 1px 1px, hsl(var(--foreground) / 0.06) 1px, transparent 0)',
          backgroundSize: '32px 32px',
        }}
      />

      <div className="max-w-6xl mx-auto">

        {/* Heading */}
        <div className="text-center mb-14">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-5">
            <span className="h-1.5 w-1.5 rounded-full bg-primary" /> 60-second tour
          </span>
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5 leading-[1.1]">
            What is Grit?
          </h2>
          <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            A full-stack meta-framework that turns one CLI command into a working Go API,
            React frontend, and admin panel — with everything you usually wire up by hand
            already wired up.
          </p>
        </div>

        {/* Three-step flow */}
        <div className="relative">
          {/* Animated horizontal connector between steps (desktop only) */}
          <div
            aria-hidden
            className="hidden md:block absolute top-[88px] left-[16%] right-[16%] h-px"
            style={{
              background:
                'linear-gradient(to right, transparent 0%, hsl(var(--primary)/0.3) 20%, hsl(var(--primary)/0.6) 50%, hsl(var(--primary)/0.3) 80%, transparent 100%)',
            }}
          />

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {STEPS.map((s) => {
              const Icon = s.icon
              return (
                <div
                  key={s.n}
                  className="group relative rounded-2xl border-2 border-border/50 bg-card/80 backdrop-blur p-6 hover:border-primary/40 hover:-translate-y-0.5 transition-all"
                >
                  {/* Number ribbon */}
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-2.5">
                      <div className="h-9 w-9 rounded-xl bg-primary/10 text-primary flex items-center justify-center ring-1 ring-primary/20">
                        <Icon className="h-4 w-4" />
                      </div>
                      <span className="font-mono font-bold text-[11px] tracking-widest text-muted-foreground">
                        STEP / {s.n}
                      </span>
                    </div>
                  </div>

                  <h3 className="text-lg font-semibold text-foreground mb-2 leading-snug">{s.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-4">{s.body}</p>

                  {/* Inline command */}
                  <div className="flex items-center gap-2 rounded-lg border border-border/50 bg-background/60 px-3 py-2 font-mono text-[12px]">
                    <span className="text-primary/70 select-none">$</span>
                    <span className="text-foreground/90 truncate">{s.code}</span>
                  </div>
                </div>
              )
            })}
          </div>
        </div>

        {/* Highlights row */}
        <div className="mt-14 rounded-2xl border-2 border-border/50 bg-card/60 p-6 md:p-8">
          <div className="grid grid-cols-1 md:grid-cols-[1fr_auto] gap-6 items-center">
            <div>
              <p className="text-xs font-mono font-medium text-primary uppercase tracking-wider mb-3">
                What you get out of the box
              </p>
              <ul className="space-y-2.5">
                {HIGHLIGHTS.map((line) => (
                  <li key={line} className="flex items-start gap-2.5 text-sm md:text-base text-foreground/85">
                    <Check className="h-4 w-4 text-primary mt-0.5 shrink-0" strokeWidth={2.5} />
                    <span className="leading-relaxed">{line}</span>
                  </li>
                ))}
              </ul>
            </div>
            <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground self-start md:self-center px-7 h-11">
              <Link href="/docs/getting-started/quick-start">
                Read quick start <ArrowRight className="ml-2 h-3.5 w-3.5" />
              </Link>
            </Button>
          </div>
        </div>

      </div>
    </section>
  )
}
