'use client'

import { useMemo, useState } from 'react'
import Link from 'next/link'
import {
  ArrowLeft,
  ArrowRight,
  Check,
  Globe,
  Smartphone,
  Monitor,
  Server,
  Layers,
  Sparkles,
  Building,
  ShoppingCart,
  FileText,
  Wrench,
  Bot,
  Box,
  ExternalLink,
} from 'lucide-react'
import { Textarea } from '@/components/ui/textarea'
import { CopyButton } from '@/components/copy-button'
import { buildKitPrompt, KITS, type KitSlug } from '@/lib/kit-prompts'
import { cn } from '@/lib/utils'

type Platform = 'web' | 'mobile' | 'desktop' | 'api' | 'multi'

interface PlatformOption {
  id: Platform
  title: string
  body: string
  icon: React.ComponentType<{ className?: string }>
  /** Which kit slugs this platform unlocks. */
  kits: KitSlug[]
}

const PLATFORMS: PlatformOption[] = [
  {
    id: 'web',
    title: 'Web app',
    body: 'A browser-first product — marketing site, dashboard, SaaS, or internal tool.',
    icon: Globe,
    kits: ['single', 'single-vite', 'double', 'triple'],
  },
  {
    id: 'mobile',
    title: 'Mobile app',
    body: 'React Native via Expo, with a shared Go API. iOS / Android.',
    icon: Smartphone,
    kits: ['mobile'],
  },
  {
    id: 'desktop',
    title: 'Desktop app',
    body: 'Native Wails desktop binary — offline-first with local SQLite.',
    icon: Monitor,
    kits: ['desktop'],
  },
  {
    id: 'api',
    title: 'API only',
    body: 'Headless Go API. Bring your own frontend (mobile, bot, native).',
    icon: Server,
    kits: ['api'],
  },
  {
    id: 'multi',
    title: 'Multi-platform',
    body: 'Web + admin + mobile (and optionally desktop) sharing one API.',
    icon: Layers,
    kits: ['triple', 'mobile', 'desktop'],
  },
]

interface KitOption {
  slug: KitSlug
  title: string
  body: string
  badge?: string
}

const KIT_BLURBS: Record<KitSlug, KitOption> = {
  single: {
    slug: 'single',
    title: 'Single — Next.js + Go',
    body: 'One binary, Next.js embedded. Smallest deploy; best for marketing + product in one app.',
  },
  'single-vite': {
    slug: 'single-vite',
    title: 'Single — Vite + TanStack Router',
    body: 'One binary, Vite SPA embedded. Sub-second cold starts; best for pure dashboard apps.',
  },
  double: {
    slug: 'double',
    title: 'Double — Web + API',
    body: 'Turborepo with apps/web (Next.js) + apps/api (Go). Two deploy targets.',
  },
  triple: {
    slug: 'triple',
    title: 'Triple — Web + Admin + API',
    body: 'The full SaaS shape. Public site + Filament-style admin + Go API.',
    badge: 'Recommended',
  },
  api: {
    slug: 'api',
    title: 'API only',
    body: 'Pure Gin + GORM. OpenAPI auto-served at /docs.',
  },
  mobile: {
    slug: 'mobile',
    title: 'Mobile — Expo + API',
    body: 'React Native + Go API, shared Zod schemas, EAS Build configured.',
  },
  desktop: {
    slug: 'desktop',
    title: 'Desktop — Wails',
    body: 'Standalone Wails v2 app with local SQLite + bcrypt auth.',
  },
}

interface UseCaseOption {
  id: string
  title: string
  body: string
  icon: React.ComponentType<{ className?: string }>
  /** One-line prompt-side hint. */
  hint: string
}

const USE_CASES: UseCaseOption[] = [
  {
    id: 'saas',
    title: 'SaaS product',
    body: 'Multi-tenant, billing, admin, end-user dashboards.',
    icon: Building,
    hint: 'A multi-tenant SaaS — each tenant gets its own scoped data and roles; billing via Stripe; an admin panel for staff plus a dashboard for end users.',
  },
  {
    id: 'ecommerce',
    title: 'E-commerce / marketplace',
    body: 'Products, orders, checkout, fulfilment.',
    icon: ShoppingCart,
    hint: 'An online store — product catalogue, cart, checkout with payments, order management, inventory, and a back-office for admin.',
  },
  {
    id: 'internal',
    title: 'Internal tool',
    body: 'Operations dashboard, ops admin, CRM-style.',
    icon: Wrench,
    hint: 'An internal operations tool — focused on a small team, RBAC matters, dashboards summarize KPIs, audit log on every mutation.',
  },
  {
    id: 'content',
    title: 'Content / blog / CMS',
    body: 'Articles, authors, drafts, public reading surface.',
    icon: FileText,
    hint: 'A content publishing platform — drafts and published states, multi-author, SEO-first public pages, an admin for editorial.',
  },
  {
    id: 'ai',
    title: 'AI-native product',
    body: 'LLM workflows, chat, agentic flows.',
    icon: Bot,
    hint: 'An AI-native product that leans on the Grit AI Gateway — chat history, streaming responses, model selection, and rate-limited usage tracking per user.',
  },
  {
    id: 'other',
    title: 'Something else',
    body: 'Anything not on this list — describe it below.',
    icon: Sparkles,
    hint: '',
  },
]

const STEPS = ['Platform', 'Tech Kit', 'Use case', 'Your prompt'] as const

export function AIIntegrationWizard() {
  const [step, setStep] = useState<0 | 1 | 2 | 3>(0)
  const [platform, setPlatform] = useState<Platform | null>(null)
  const [kit, setKit] = useState<KitSlug | null>(null)
  const [useCase, setUseCase] = useState<string | null>(null)
  const [idea, setIdea] = useState('')

  const availableKits = useMemo<KitOption[]>(() => {
    if (!platform) return []
    const slugs = PLATFORMS.find((p) => p.id === platform)?.kits ?? []
    return slugs.map((s) => KIT_BLURBS[s])
  }, [platform])

  const generatedPrompt = useMemo(() => {
    if (!kit) return ''
    const useCaseLabel = USE_CASES.find((u) => u.id === useCase)
    return buildKitPrompt(kit, {
      useCase: useCaseLabel?.hint,
      idea: idea.trim(),
    })
  }, [kit, useCase, idea])

  const canAdvance =
    (step === 0 && platform) ||
    (step === 1 && kit) ||
    (step === 2 && useCase) ||
    step === 3

  const handleNext = () => {
    if (!canAdvance) return
    if (step < 3) setStep((step + 1) as 0 | 1 | 2 | 3)
  }
  const handleBack = () => {
    if (step > 0) setStep((step - 1) as 0 | 1 | 2 | 3)
  }
  const handleRestart = () => {
    setStep(0)
    setPlatform(null)
    setKit(null)
    setUseCase(null)
    setIdea('')
  }

  return (
    <div>
      {/* Progress dots */}
      <div className="mb-10 flex items-center gap-3 flex-wrap">
        {STEPS.map((label, i) => {
          const isActive = i === step
          const isDone = i < step
          return (
            <div key={label} className="flex items-center gap-3">
              <div className="flex items-center gap-2.5">
                <div
                  className={cn(
                    'h-7 w-7 rounded-full text-[11px] font-semibold flex items-center justify-center transition-colors border',
                    isActive && 'bg-primary text-primary-foreground border-primary',
                    isDone && 'bg-primary/15 text-primary border-primary/30',
                    !isActive && !isDone && 'bg-card/40 text-muted-foreground border-border/40',
                  )}
                >
                  {isDone ? <Check className="h-3.5 w-3.5" /> : i + 1}
                </div>
                <span
                  className={cn(
                    'text-[13px] font-medium hidden sm:block',
                    isActive && 'text-foreground',
                    isDone && 'text-muted-foreground',
                    !isActive && !isDone && 'text-muted-foreground/60',
                  )}
                >
                  {label}
                </span>
              </div>
              {i < STEPS.length - 1 && (
                <div className={cn('h-px w-8 sm:w-14', isDone ? 'bg-primary/40' : 'bg-border/40')} />
              )}
            </div>
          )
        })}
      </div>

      {/* Card */}
      <div className="rounded-2xl border border-border/50 bg-card/40 p-6 sm:p-8 mb-6 min-h-[480px]">
        {step === 0 && (
          <StepPlatform
            value={platform}
            onChange={(v) => {
              setPlatform(v)
              // Reset downstream choices when platform changes
              setKit(null)
            }}
          />
        )}
        {step === 1 && (
          <StepKit
            options={availableKits}
            value={kit}
            onChange={setKit}
          />
        )}
        {step === 2 && (
          <StepUseCase
            value={useCase}
            onChange={setUseCase}
            idea={idea}
            onIdeaChange={setIdea}
          />
        )}
        {step === 3 && kit && <StepPrompt prompt={generatedPrompt} kitSlug={kit} />}
      </div>

      {/*
        Nav buttons. The Nexora chat widget mounts a fixed bottom-right
        button (~96px square + label). Bottom padding here ensures the Next
        button never sits beneath it; the right-side spacer on `Next` keeps
        them side-by-side on narrow screens instead of stacked.
      */}
      <div className="flex items-center justify-between gap-3 pb-28">
        <button
          type="button"
          onClick={handleBack}
          disabled={step === 0}
          className={cn(
            'inline-flex items-center gap-1.5 rounded-full border px-4 py-2 text-sm font-medium transition-colors',
            step === 0
              ? 'border-border/30 text-muted-foreground/40 cursor-not-allowed'
              : 'border-border/60 text-foreground hover:bg-accent/40',
          )}
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Back
        </button>

        <span className="text-xs text-muted-foreground hidden sm:block">
          Step {step + 1} of {STEPS.length}
        </span>

        {step < 3 ? (
          <button
            type="button"
            onClick={handleNext}
            disabled={!canAdvance}
            className={cn(
              'inline-flex items-center gap-1.5 rounded-full px-5 py-2 text-sm font-medium transition-colors sm:mr-32',
              canAdvance
                ? 'bg-primary text-primary-foreground hover:bg-primary/90'
                : 'bg-primary/30 text-primary-foreground/70 cursor-not-allowed',
            )}
          >
            Next
            <ArrowRight className="h-3.5 w-3.5" />
          </button>
        ) : (
          <button
            type="button"
            onClick={handleRestart}
            className="inline-flex items-center gap-1.5 rounded-full border border-border/60 px-5 py-2 text-sm font-medium text-foreground hover:bg-accent/40 transition-colors sm:mr-32"
          >
            Start over
          </button>
        )}
      </div>
    </div>
  )
}

/* ─── Step 1: Platform ─────────────────────────────────────────────── */

function StepPlatform({
  value,
  onChange,
}: {
  value: Platform | null
  onChange: (p: Platform) => void
}) {
  return (
    <div>
      <h2 className="text-xl font-semibold tracking-tight mb-2">What are you building?</h2>
      <p className="text-sm text-muted-foreground mb-6">
        The prompt adapts to your stack. Mobile picks up Expo + React Query patterns; web tunes
        for SSR vs SPA.
      </p>
      <div className="space-y-2.5">
        {PLATFORMS.map((p) => {
          const Icon = p.icon
          const selected = value === p.id
          return (
            <button
              key={p.id}
              type="button"
              onClick={() => onChange(p.id)}
              className={cn(
                'w-full text-left rounded-xl border p-4 transition-all flex items-start gap-4',
                selected
                  ? 'border-primary/60 bg-primary/[0.06] ring-1 ring-primary/30'
                  : 'border-border/40 bg-card/40 hover:border-primary/30 hover:bg-card/60',
              )}
            >
              <div
                className={cn(
                  'h-10 w-10 rounded-lg flex items-center justify-center shrink-0',
                  selected
                    ? 'bg-primary/15 text-primary'
                    : 'bg-muted/40 text-muted-foreground',
                )}
              >
                <Icon className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-0.5">
                  <p className="font-semibold text-foreground">{p.title}</p>
                  {selected && (
                    <Check className="h-4 w-4 text-primary" strokeWidth={3} />
                  )}
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed">{p.body}</p>
              </div>
            </button>
          )
        })}
      </div>
    </div>
  )
}

/* ─── Step 2: Kit ──────────────────────────────────────────────────── */

function StepKit({
  options,
  value,
  onChange,
}: {
  options: KitOption[]
  value: KitSlug | null
  onChange: (k: KitSlug) => void
}) {
  return (
    <div>
      <h2 className="text-xl font-semibold tracking-tight mb-2">Pick a tech kit</h2>
      <p className="text-sm text-muted-foreground mb-6">
        Each kit is a complete scaffold — pre-wired auth, jobs, mail, storage, AI, security, and
        observability.
      </p>
      <div className="space-y-2.5">
        {options.map((o) => {
          const selected = value === o.slug
          return (
            <button
              key={o.slug}
              type="button"
              onClick={() => onChange(o.slug)}
              className={cn(
                'w-full text-left rounded-xl border p-4 transition-all flex items-start gap-4',
                selected
                  ? 'border-primary/60 bg-primary/[0.06] ring-1 ring-primary/30'
                  : 'border-border/40 bg-card/40 hover:border-primary/30 hover:bg-card/60',
              )}
            >
              <div
                className={cn(
                  'h-10 w-10 rounded-lg flex items-center justify-center shrink-0',
                  selected
                    ? 'bg-primary/15 text-primary'
                    : 'bg-muted/40 text-muted-foreground',
                )}
              >
                <Box className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-0.5 flex-wrap">
                  <p className="font-semibold text-foreground">{o.title}</p>
                  {o.badge && (
                    <span className="text-[9px] font-mono uppercase tracking-wider bg-primary/10 text-primary border border-primary/20 px-1.5 py-0.5 rounded-full">
                      {o.badge}
                    </span>
                  )}
                  {selected && (
                    <Check className="h-4 w-4 text-primary ml-auto" strokeWidth={3} />
                  )}
                </div>
                <p className="text-sm text-muted-foreground leading-relaxed mb-1.5">{o.body}</p>
                <code className="text-[10px] font-mono text-primary/70 bg-primary/5 px-1.5 py-0.5 rounded">
                  {KITS[o.slug].command}
                </code>
              </div>
            </button>
          )
        })}
      </div>
    </div>
  )
}

/* ─── Step 3: Use case ─────────────────────────────────────────────── */

function StepUseCase({
  value,
  onChange,
  idea,
  onIdeaChange,
}: {
  value: string | null
  onChange: (id: string) => void
  idea: string
  onIdeaChange: (v: string) => void
}) {
  return (
    <div>
      <h2 className="text-xl font-semibold tracking-tight mb-2">What kind of project is it?</h2>
      <p className="text-sm text-muted-foreground mb-6">
        This tunes the prompt so Claude understands the project shape and recommends the right
        Grit primitives.
      </p>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5 mb-7">
        {USE_CASES.map((u) => {
          const Icon = u.icon
          const selected = value === u.id
          return (
            <button
              key={u.id}
              type="button"
              onClick={() => onChange(u.id)}
              className={cn(
                'text-left rounded-xl border p-4 transition-all flex items-start gap-3',
                selected
                  ? 'border-primary/60 bg-primary/[0.06] ring-1 ring-primary/30'
                  : 'border-border/40 bg-card/40 hover:border-primary/30 hover:bg-card/60',
              )}
            >
              <div
                className={cn(
                  'h-9 w-9 rounded-lg flex items-center justify-center shrink-0',
                  selected
                    ? 'bg-primary/15 text-primary'
                    : 'bg-muted/40 text-muted-foreground',
                )}
              >
                <Icon className="h-4 w-4" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-0.5">
                  <p className="font-semibold text-foreground text-sm">{u.title}</p>
                  {selected && (
                    <Check className="h-3.5 w-3.5 text-primary ml-auto" strokeWidth={3} />
                  )}
                </div>
                <p className="text-xs text-muted-foreground leading-relaxed">{u.body}</p>
              </div>
            </button>
          )
        })}
      </div>

      <div>
        <label htmlFor="idea" className="block text-sm font-medium text-foreground mb-2">
          Describe your idea
          <span className="text-muted-foreground font-normal"> — optional, but recommended</span>
        </label>
        <Textarea
          id="idea"
          value={idea}
          onChange={(e) => onIdeaChange(e.target.value)}
          rows={6}
          placeholder="E.g. A platform for music teachers to schedule lessons, send reminders, and accept payments. Single tenant per teacher — they sign up, set their availability, and share a public booking link with students. Inspired by Calendly but tuned for recurring lessons."
          className="resize-y bg-background/40 border-border/40 text-sm leading-relaxed"
        />
        <p className="mt-2 text-xs text-muted-foreground">
          A paragraph here gets pasted directly into the prompt. Leave blank and you&apos;ll see
          a <code className="text-[0.85em]">[YOUR IDEA]</code> placeholder in the result that you
          fill in later.
        </p>
      </div>
    </div>
  )
}

/* ─── Step 4: Generated prompt ─────────────────────────────────────── */

function StepPrompt({ prompt, kitSlug }: { prompt: string; kitSlug: KitSlug }) {
  const kit = KITS[kitSlug]
  return (
    <div>
      <div className="flex items-start justify-between gap-4 mb-5 flex-wrap">
        <div>
          <h2 className="text-xl font-semibold tracking-tight mb-1.5 flex items-center gap-2">
            <Sparkles className="h-5 w-5 text-primary" />
            Your prompt is ready
          </h2>
          <p className="text-sm text-muted-foreground max-w-2xl">
            Paste this into{' '}
            <Link
              href="https://claude.ai"
              target="_blank"
              rel="noopener noreferrer"
              className="text-primary hover:underline"
            >
              claude.ai
            </Link>
            . Claude will return four planning files. Save them next to your scaffold and run{' '}
            <code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90">
              {kit.command}
            </code>
            .
          </p>
        </div>
        <div className="flex items-center gap-2 shrink-0">
          <CopyButton text={prompt} variant="default" className="bg-primary text-primary-foreground hover:bg-primary/90">
            Copy prompt
          </CopyButton>
        </div>
      </div>

      <div className="rounded-xl border border-border/50 bg-background/40 overflow-hidden">
        <div className="flex items-center justify-between gap-3 border-b border-border/40 bg-muted/30 px-4 py-2.5">
          <div className="flex items-center gap-2 min-w-0">
            <span className="h-2 w-2 rounded-full bg-rose-400/70" />
            <span className="h-2 w-2 rounded-full bg-amber-400/70" />
            <span className="h-2 w-2 rounded-full bg-emerald-400/70" />
            <span className="ml-2 text-[11px] font-mono text-muted-foreground truncate">
              prompt — {kit.label}
            </span>
          </div>
          <CopyButton
            text={prompt}
            size="sm"
            variant="ghost"
            className="h-7 px-2.5 text-[11px] text-muted-foreground hover:text-foreground"
          >
            Copy
          </CopyButton>
        </div>
        <pre className="px-5 py-5 text-[12.5px] leading-[1.65] font-mono text-foreground/85 whitespace-pre-wrap break-words max-h-[460px] overflow-y-auto">
          {prompt}
        </pre>
      </div>

      {/* Next steps */}
      <div className="mt-6 grid grid-cols-1 md:grid-cols-3 gap-3">
        {[
          {
            n: '1',
            t: 'Paste into claude.ai',
            d: 'Open a fresh chat and paste the full prompt above.',
            href: 'https://claude.ai',
            external: true,
          },
          {
            n: '2',
            t: 'Save the 4 files',
            d: 'Claude returns project-description.md, project-phases.md, design-style-guide.md, prompt.md.',
          },
          {
            n: '3',
            t: 'Run prompt.md in Claude Code',
            d: `Scaffold with \`${kit.command}\`, then paste prompt.md to begin building.`,
          },
        ].map((s) => (
          <div key={s.n} className="rounded-xl border border-border/40 bg-card/30 p-4">
            <div className="flex items-center gap-2 mb-1.5">
              <span className="h-6 w-6 rounded-full bg-primary/10 text-primary text-[11px] font-mono font-semibold flex items-center justify-center">
                {s.n}
              </span>
              <p className="text-sm font-semibold">{s.t}</p>
            </div>
            <p className="text-xs text-muted-foreground leading-relaxed">{s.d}</p>
            {s.href && (
              <Link
                href={s.href}
                target="_blank"
                rel="noopener noreferrer"
                className="mt-2 inline-flex items-center gap-1 text-[11px] text-primary hover:underline"
              >
                Open claude.ai
                <ExternalLink className="h-3 w-3" />
              </Link>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
