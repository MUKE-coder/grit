'use client'

import * as React from 'react'
import {
  Layers, Boxes, Globe, Server, Smartphone, Monitor,
  Check, Terminal, Sparkles, Palette, Zap, RotateCcw,
} from 'lucide-react'
import { CopyButton } from '@/components/copy-button'
import { cn } from '@/lib/utils'

// The builder mirrors the real `grit new` flag surface. Every command it emits
// must be one the CLI actually accepts — so the option data and the constraints
// below track cmd/grit/main.go, not a wishlist. When a flag changes there, it
// changes here.

type ArchId = 'triple' | 'full' | 'double' | 'single' | 'api' | 'mobile'
type FrontendId = 'next' | 'vite'
type StyleId = 'default' | 'modern' | 'minimal' | 'glass'
type ThemeId = 'atlas' | 'aurora' | 'pulse'

interface ArchOption {
  id: ArchId
  name: string
  desc: string
  icon: React.ComponentType<{ className?: string }>
  hasFrontend: boolean
  hasAdmin: boolean
  badge?: string
}

const ARCHITECTURES: ArchOption[] = [
  { id: 'triple', name: 'Triple', desc: 'Web + Admin + API in a Turborepo. The full-stack default.', icon: Layers, hasFrontend: true, hasAdmin: true, badge: 'Popular' },
  { id: 'full', name: 'Full', desc: 'Triple plus a docs site, an Expo mobile app, and a desktop app — everything.', icon: Boxes, hasFrontend: true, hasAdmin: true },
  { id: 'double', name: 'Double', desc: 'Web + API. A public site with no admin panel.', icon: Globe, hasFrontend: true, hasAdmin: false },
  { id: 'single', name: 'Single', desc: 'Go API with an embedded React SPA — one deployable binary.', icon: Server, hasFrontend: true, hasAdmin: false },
  { id: 'api', name: 'API only', desc: 'Just the Go API. No frontend.', icon: Terminal, hasFrontend: false, hasAdmin: false },
  { id: 'mobile', name: 'Mobile', desc: 'Go API paired with an Expo mobile app.', icon: Smartphone, hasFrontend: false, hasAdmin: false },
]

const FRONTENDS: { id: FrontendId; name: string; desc: string }[] = [
  { id: 'next', name: 'Next.js', desc: 'App Router, SSR. The default.' },
  { id: 'vite', name: 'TanStack Router', desc: 'Vite SPA with type-safe routing.' },
]

const STYLES: { id: StyleId; name: string; desc: string }[] = [
  { id: 'default', name: 'Default', desc: 'The standard premium dark admin.' },
  { id: 'modern', name: 'Modern', desc: 'Softer surfaces, rounded.' },
  { id: 'minimal', name: 'Minimal', desc: 'Flat, restrained, high-contrast.' },
  { id: 'glass', name: 'Glass', desc: 'Frosted, translucent panels.' },
]

const THEMES: { id: ThemeId; name: string; desc: string }[] = [
  { id: 'atlas', name: 'Atlas', desc: 'The default. Purple on deep navy.' },
  { id: 'aurora', name: 'Aurora', desc: 'Cool teal / green accents.' },
  { id: 'pulse', name: 'Pulse', desc: 'Warm, high-energy accents.' },
]

interface BuilderState {
  name: string
  arch: ArchId
  frontend: FrontendId
  desktop: boolean
  expo: boolean
  style: StyleId
  theme: ThemeId
}

const DEFAULT_STATE: BuilderState = {
  name: 'my-app',
  arch: 'triple',
  frontend: 'next',
  desktop: false,
  expo: false,
  style: 'default',
  theme: 'atlas',
}

const archById = (id: ArchId) => ARCHITECTURES.find((a) => a.id === id)!

// desktop needs a monorepo API to talk to, so the CLI rejects it with --single;
// it's meaningless without a frontend (api/mobile). Offer it only where valid.
// full already bundles it.
function desktopAvailable(arch: ArchId): boolean {
  return arch === 'double' || arch === 'triple'
}

// an extra Expo app is additive on the web archs; mobile IS expo, and full
// already bundles it, so neither needs the toggle.
function expoAvailable(arch: ArchId): boolean {
  return arch === 'double' || arch === 'triple'
}

// buildCommand assembles a `grit new` invocation from the current selection,
// omitting anything that's already the CLI default (Next.js frontend, atlas
// theme, default style) so the command stays as short as what you'd actually
// type — the same way the CLI's own docs show it.
function buildCommand(s: BuilderState): string {
  const name = s.name.trim() || 'my-app'
  const parts = ['grit new', name]
  const arch = archById(s.arch)

  if (s.arch === 'api') {
    parts.push('--api')
  } else if (s.arch === 'mobile') {
    parts.push('--mobile')
  } else if (s.arch === 'full') {
    parts.push('--full')
    if (s.frontend === 'vite') parts.push('--vite') // full defaults to next
  } else {
    parts.push('--' + s.arch, '--' + s.frontend)
  }

  // Add-ons — never for full (it bundles them) or where they're unavailable.
  if (s.arch !== 'full') {
    if (s.desktop && desktopAvailable(s.arch)) parts.push('--desktop')
    if (s.expo && expoAvailable(s.arch)) parts.push('--expo')
  }

  // Admin style only applies where there's an admin panel, and only when
  // it's not the default.
  if (arch.hasAdmin && s.style !== 'default') parts.push('--style', s.style)

  // Theme only applies where there's a frontend, and only when not the default.
  if (arch.hasFrontend && s.theme !== 'atlas') parts.push('--theme', s.theme)

  return parts.join(' ')
}

function Card({
  active, onClick, disabled, icon: Icon, name, desc, badge,
}: {
  active: boolean
  onClick: () => void
  disabled?: boolean
  icon?: React.ComponentType<{ className?: string }>
  name: string
  desc: string
  badge?: string
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={disabled}
      className={cn(
        'group relative flex w-full flex-col rounded-xl border p-4 text-left transition-all',
        active
          ? 'border-primary bg-primary/5 ring-1 ring-primary/40'
          : 'border-border bg-card/40 hover:border-primary/40 hover:bg-card/60',
        disabled && 'cursor-not-allowed opacity-40 hover:border-border hover:bg-card/40',
      )}
    >
      <div className="flex items-center gap-2">
        {Icon && <Icon className={cn('h-4 w-4', active ? 'text-primary' : 'text-muted-foreground')} />}
        <span className="font-medium text-foreground">{name}</span>
        {badge && (
          <span className="ml-auto rounded-full bg-primary/15 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-primary">
            {badge}
          </span>
        )}
        {active && !badge && <Check className="ml-auto h-4 w-4 text-primary" />}
      </div>
      <p className="mt-1.5 text-[13px] leading-relaxed text-muted-foreground">{desc}</p>
    </button>
  )
}

function Section({ label, hint, children }: { label: string; hint?: string; children: React.ReactNode }) {
  return (
    <section className="scroll-mt-24">
      <div className="mb-3 flex items-baseline gap-3">
        <h2 className="font-mono text-sm font-semibold uppercase tracking-wider text-foreground">{label}</h2>
        {hint && <span className="text-xs text-muted-foreground">{hint}</span>}
      </div>
      {children}
    </section>
  )
}

export function StackBuilder() {
  const [state, setState] = React.useState<BuilderState>(DEFAULT_STATE)
  const set = <K extends keyof BuilderState>(key: K, value: BuilderState[K]) =>
    setState((s) => ({ ...s, [key]: value }))

  const arch = archById(state.arch)
  const command = React.useMemo(() => buildCommand(state), [state])

  // Selected-stack chips — only the things that actually shape the command.
  const chips: string[] = [state.arch]
  if (arch.hasFrontend) chips.push(state.frontend === 'vite' ? 'tanstack' : 'next')
  if (state.arch !== 'full') {
    if (state.desktop && desktopAvailable(state.arch)) chips.push('desktop')
    if (state.expo && expoAvailable(state.arch)) chips.push('expo')
  } else {
    chips.push('docs', 'expo', 'desktop')
  }
  if (arch.hasAdmin && state.style !== 'default') chips.push(state.style)
  if (arch.hasFrontend && state.theme !== 'atlas') chips.push(state.theme)

  return (
    <div className="grid grid-cols-1 gap-8 lg:grid-cols-[340px_1fr]">
      {/* Left: sticky summary */}
      <aside className="lg:sticky lg:top-24 lg:h-fit space-y-5">
        <div>
          <label htmlFor="project-name" className="mb-2 block font-mono text-xs uppercase tracking-wider text-muted-foreground">
            Project name
          </label>
          <input
            id="project-name"
            value={state.name}
            onChange={(e) => set('name', e.target.value)}
            spellCheck={false}
            placeholder="my-app"
            className="w-full rounded-lg border border-border bg-background px-3 py-2.5 font-mono text-sm text-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary/40"
          />
        </div>

        <div>
          <div className="mb-2 flex items-center justify-between">
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Command</span>
            <CopyButton text={command} variant="ghost" className="h-7 px-2 text-xs">Copy</CopyButton>
          </div>
          <div className="rounded-lg border border-border bg-background/60 p-3">
            <code className="block whitespace-pre-wrap break-words font-mono text-[13px] leading-relaxed text-foreground">
              <span className="text-primary">$</span> {command}
            </code>
          </div>
        </div>

        <div>
          <div className="mb-2 flex items-center justify-between">
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Selected stack</span>
            <span className="text-xs text-muted-foreground">{chips.length} picks</span>
          </div>
          <div className="flex flex-wrap gap-1.5">
            {chips.map((c) => (
              <span key={c} className="rounded-md border border-border bg-card/60 px-2 py-1 font-mono text-[11px] text-foreground">
                {c}
              </span>
            ))}
          </div>
        </div>

        <button
          type="button"
          onClick={() => setState(DEFAULT_STATE)}
          className="inline-flex items-center gap-1.5 text-xs text-muted-foreground transition-colors hover:text-foreground"
        >
          <RotateCcw className="h-3.5 w-3.5" />
          Reset to defaults
        </button>

        <p className="rounded-lg border border-border bg-muted/30 p-3 text-[12px] leading-relaxed text-muted-foreground">
          Run the command, then pick your database in <code className="text-foreground">.env</code> —
          Postgres by default, or SQLite for a zero-setup start. Optional modules (AI, jobs,
          storage, 2FA…) toggle via <code className="text-foreground">MODULE_*</code> flags.
        </p>
      </aside>

      {/* Right: options */}
      <div className="space-y-10">
        <Section label="Architecture" hint="What gets scaffolded">
          <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
            {ARCHITECTURES.map((a) => (
              <Card
                key={a.id}
                active={state.arch === a.id}
                onClick={() => set('arch', a.id)}
                icon={a.icon}
                name={a.name}
                desc={a.desc}
                badge={a.badge}
              />
            ))}
          </div>
        </Section>

        {arch.hasFrontend && (
          <Section label="Web Frontend" hint="Router & rendering">
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              {FRONTENDS.map((f) => (
                <Card
                  key={f.id}
                  active={state.frontend === f.id}
                  onClick={() => set('frontend', f.id)}
                  name={f.name}
                  desc={f.desc}
                />
              ))}
            </div>
          </Section>
        )}

        {(desktopAvailable(state.arch) || expoAvailable(state.arch)) && (
          <Section label="Add-ons" hint="Extra apps in the monorepo">
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              {desktopAvailable(state.arch) && (
                <Card
                  active={state.desktop}
                  onClick={() => set('desktop', !state.desktop)}
                  icon={Monitor}
                  name="Desktop (Wails)"
                  desc="A native desktop app sharing the monorepo API."
                />
              )}
              {expoAvailable(state.arch) && (
                <Card
                  active={state.expo}
                  onClick={() => set('expo', !state.expo)}
                  icon={Smartphone}
                  name="Expo mobile"
                  desc="A React Native app alongside the web apps."
                />
              )}
            </div>
          </Section>
        )}

        {state.arch === 'full' && (
          <Section label="Add-ons" hint="Included with Full">
            <div className="rounded-xl border border-border bg-card/40 p-4 text-[13px] text-muted-foreground">
              <span className="font-medium text-foreground">Full</span> already bundles the docs
              site, an Expo mobile app, and a Wails desktop app — no extra flags needed.
            </div>
          </Section>
        )}

        {arch.hasAdmin && (
          <Section label="Admin Style" hint="Look of the admin panel">
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              {STYLES.map((st) => (
                <Card
                  key={st.id}
                  active={state.style === st.id}
                  onClick={() => set('style', st.id)}
                  icon={Palette}
                  name={st.name}
                  desc={st.desc}
                />
              ))}
            </div>
          </Section>
        )}

        {arch.hasFrontend && (
          <Section label="Theme" hint="Colors, fonts, auth pages">
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
              {THEMES.map((t) => (
                <Card
                  key={t.id}
                  active={state.theme === t.id}
                  onClick={() => set('theme', t.id)}
                  icon={t.id === 'aurora' ? Zap : t.id === 'pulse' ? Sparkles : Palette}
                  name={t.name}
                  desc={t.desc}
                />
              ))}
            </div>
          </Section>
        )}
      </div>
    </div>
  )
}
