'use client'

import * as React from 'react'
import {
  Layers, Boxes, Globe, Server, Smartphone, Monitor,
  Check, Terminal, Sparkles, Palette, Zap, RotateCcw, FolderTree,
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
  { id: 'full', name: 'Full', desc: 'Triple plus a docs site, an Expo mobile app, and a desktop app, everything.', icon: Boxes, hasFrontend: true, hasAdmin: true },
  { id: 'double', name: 'Double', desc: 'Web + API. A public site with no admin panel.', icon: Globe, hasFrontend: true, hasAdmin: false },
  { id: 'single', name: 'Single', desc: 'Go API with an embedded React SPA, one deployable binary.', icon: Server, hasFrontend: true, hasAdmin: false },
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

// Preview tree — what `grit new` actually writes for the current selection.
// The shapes here are ground-truthed against real scaffolds: single is a
// single-binary Go app with an embedded frontend (no monorepo); api is just
// apps/api; the rest are Turborepo monorepos whose apps/ set follows the arch
// and add-ons. The extra mobile app scaffolds to apps/expo, not apps/mobile.
interface TreeNode {
  name: string
  note?: string
  children?: TreeNode[]
}

function buildTree(s: BuilderState): TreeNode {
  const name = (s.name.trim() || 'my-app') + '/'

  if (s.arch === 'single') {
    return {
      name,
      children: [
        { name: 'cmd/', note: 'server entrypoint' },
        { name: 'internal/', note: 'handlers, services, models' },
        { name: 'frontend/', note: s.frontend === 'vite' ? 'embedded TanStack SPA' : 'embedded React SPA' },
        { name: 'main.go' },
        { name: 'go.mod' },
        { name: 'Dockerfile' },
        { name: 'docker-compose.yml' },
        { name: 'README.md' },
      ],
    }
  }

  const feNote = s.frontend === 'vite' ? 'TanStack Router (Vite)' : 'Next.js (App Router)'
  const apps: TreeNode[] = [{ name: 'api/', note: 'Go + Gin + GORM' }]
  if (s.arch === 'double' || s.arch === 'triple' || s.arch === 'full') apps.push({ name: 'web/', note: feNote })
  if (s.arch === 'triple' || s.arch === 'full') apps.push({ name: 'admin/', note: 'admin panel' })
  const withExpo = s.arch === 'mobile' || s.arch === 'full' || (expoAvailable(s.arch) && s.expo)
  if (withExpo) apps.push({ name: 'expo/', note: 'Expo mobile app' })
  const withDesktop = s.arch === 'full' || (desktopAvailable(s.arch) && s.desktop)
  if (withDesktop) apps.push({ name: 'desktop/', note: 'Wails desktop app' })
  if (s.arch === 'full') apps.push({ name: 'docs/', note: 'documentation site' })

  const children: TreeNode[] = [{ name: 'apps/', children: apps }]
  // api-only skips the monorepo tooling and the shared package.
  if (s.arch !== 'api') {
    children.push({ name: 'packages/', children: [{ name: 'shared/', note: 'Zod schemas + TS types' }] })
  }
  children.push({ name: 'docker-compose.yml' })
  if (s.arch !== 'api') {
    children.push({ name: 'turbo.json' })
    children.push({ name: 'pnpm-workspace.yaml' })
    children.push({ name: 'grit.config.ts' })
  }
  children.push({ name: 'README.md' })
  return { name, children }
}

function flattenTree(node: TreeNode, prefix = '', isRoot = true): { line: string; note?: string }[] {
  const rows: { line: string; note?: string }[] = []
  if (isRoot) rows.push({ line: node.name })
  const kids = node.children ?? []
  kids.forEach((k, i) => {
    const last = i === kids.length - 1
    rows.push({ line: prefix + (last ? '└─ ' : '├─ ') + k.name, note: k.note })
    if (k.children) rows.push(...flattenTree(k, prefix + (last ? '   ' : '│  '), false))
  })
  return rows
}

// Presets are common starting points — one click sets the whole stack. The
// name is kept out of them so a preset never clobbers what the user typed.
const PRESETS: { name: string; desc: string; state: Omit<BuilderState, 'name'> }[] = [
  { name: 'Full-stack SaaS', desc: 'Web + admin + API, glass admin', state: { arch: 'triple', frontend: 'next', desktop: false, expo: false, style: 'glass', theme: 'atlas' } },
  { name: 'API service', desc: 'Go API only', state: { arch: 'api', frontend: 'next', desktop: false, expo: false, style: 'default', theme: 'atlas' } },
  { name: 'Web + mobile', desc: 'Web, API and an Expo app', state: { arch: 'triple', frontend: 'next', desktop: false, expo: true, style: 'default', theme: 'atlas' } },
  { name: 'Everything', desc: 'Full: web, admin, docs, mobile, desktop', state: { arch: 'full', frontend: 'next', desktop: false, expo: false, style: 'default', theme: 'atlas' } },
]

// URL state — the selection round-trips through the query string so a stack is
// shareable by link. Only non-default values are written, keeping URLs short;
// everything is validated on the way back in so a hand-edited URL can never put
// the builder into a state that emits an invalid command.
const ARCH_IDS = ARCHITECTURES.map((a) => a.id)
const FE_IDS = FRONTENDS.map((f) => f.id)
const STYLE_IDS = STYLES.map((s) => s.id)
const THEME_IDS = THEMES.map((t) => t.id)

function encodeState(s: BuilderState): string {
  const p = new URLSearchParams()
  if (s.name.trim() && s.name.trim() !== 'my-app') p.set('name', s.name.trim())
  p.set('arch', s.arch)
  if (archById(s.arch).hasFrontend && s.frontend !== 'next') p.set('fe', s.frontend)
  if (s.arch !== 'full' && s.desktop && desktopAvailable(s.arch)) p.set('desktop', '1')
  if (s.arch !== 'full' && s.expo && expoAvailable(s.arch)) p.set('expo', '1')
  if (archById(s.arch).hasAdmin && s.style !== 'default') p.set('style', s.style)
  if (archById(s.arch).hasFrontend && s.theme !== 'atlas') p.set('theme', s.theme)
  return p.toString()
}

function decodeState(qs: string): BuilderState {
  const p = new URLSearchParams(qs)
  const pick = <T extends string>(v: string | null, allowed: readonly T[], fallback: T): T =>
    v && (allowed as readonly string[]).includes(v) ? (v as T) : fallback
  const arch = pick(p.get('arch'), ARCH_IDS, DEFAULT_STATE.arch)
  return {
    name: p.get('name')?.slice(0, 60) || DEFAULT_STATE.name,
    arch,
    frontend: pick(p.get('fe'), FE_IDS, DEFAULT_STATE.frontend),
    desktop: p.get('desktop') === '1',
    expo: p.get('expo') === '1',
    style: pick(p.get('style'), STYLE_IDS, DEFAULT_STATE.style),
    theme: pick(p.get('theme'), THEME_IDS, DEFAULT_STATE.theme),
  }
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
  const [view, setView] = React.useState<'configure' | 'preview'>('configure')
  const set = <K extends keyof BuilderState>(key: K, value: BuilderState[K]) =>
    setState((s) => ({ ...s, [key]: value }))

  const arch = archById(state.arch)
  const command = React.useMemo(() => buildCommand(state), [state])

  // Hydrate from the URL once on mount. This runs client-only, so the server
  // always renders DEFAULT_STATE and there's no hydration mismatch.
  const hydrated = React.useRef(false)
  React.useEffect(() => {
    const qs = window.location.search.replace(/^\?/, '')
    if (qs) setState(decodeState(qs))
    hydrated.current = true
  }, [])

  // Mirror state back into the URL so any stack is shareable by link. Only
  // non-defaults are written (see encodeState), keeping the URL short.
  const [shareUrl, setShareUrl] = React.useState('')
  React.useEffect(() => {
    const qs = encodeState(state)
    const path = window.location.pathname
    if (hydrated.current) window.history.replaceState(null, '', qs ? `${path}?${qs}` : path)
    setShareUrl(`${window.location.origin}${path}${qs ? `?${qs}` : ''}`)
  }, [state])

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

  const treeRows = React.useMemo(() => flattenTree(buildTree(state)), [state])

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
            <div className="flex items-center gap-1">
              <CopyButton text={shareUrl} variant="ghost" className="h-7 px-2 text-xs">Share link</CopyButton>
              <CopyButton text={command} variant="ghost" className="h-7 px-2 text-xs">Copy</CopyButton>
            </div>
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
          Run the command, then pick your database in <code className="text-foreground">.env</code>:
          Postgres by default, or SQLite for a zero-setup start. Optional modules (AI, jobs,
          storage, 2FA…) toggle via <code className="text-foreground">MODULE_*</code> flags.
        </p>
      </aside>

      {/* Right: Configure / Preview */}
      <div className="space-y-10">
        <div className="inline-flex rounded-lg border border-border bg-card/40 p-1">
          {(['configure', 'preview'] as const).map((v) => (
            <button
              key={v}
              type="button"
              onClick={() => setView(v)}
              className={cn(
                'inline-flex items-center gap-1.5 rounded-md px-3 py-1.5 font-mono text-xs uppercase tracking-wider transition-colors',
                view === v ? 'bg-primary/15 text-primary' : 'text-muted-foreground hover:text-foreground',
              )}
            >
              {v === 'configure' ? <Terminal className="h-3.5 w-3.5" /> : <FolderTree className="h-3.5 w-3.5" />}
              {v}
            </button>
          ))}
        </div>

        {view === 'preview' && (
          <Section label="Project structure" hint="What grit new writes">
            <div className="overflow-x-auto rounded-xl border border-border bg-background/60 p-4">
              <pre className="font-mono text-[13px] leading-relaxed">
                {treeRows.map((r, i) => (
                  <div key={i} className="flex">
                    <span className="whitespace-pre text-foreground">{r.line}</span>
                    {r.note && <span className="ml-3 shrink-0 text-muted-foreground"># {r.note}</span>}
                  </div>
                ))}
              </pre>
            </div>
            <p className="mt-3 text-[12px] leading-relaxed text-muted-foreground">
              Top-level layout only. Each app ships with its own tests, Dockerfile, and config;
              run the command to get the full tree.
            </p>
          </Section>
        )}

        {view === 'configure' && (
        <>
        <Section label="Presets" hint="Start from a common stack">
          <div className="flex flex-wrap gap-2">
            {PRESETS.map((p) => (
              <button
                key={p.name}
                type="button"
                onClick={() => setState((s) => ({ ...p.state, name: s.name }))}
                title={p.desc}
                className="rounded-lg border border-border bg-card/40 px-3 py-2 text-left transition-all hover:border-primary/40 hover:bg-card/60"
              >
                <span className="block text-[13px] font-medium text-foreground">{p.name}</span>
                <span className="block text-[11px] text-muted-foreground">{p.desc}</span>
              </button>
            ))}
          </div>
        </Section>

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
              site, an Expo mobile app, and a Wails desktop app: no extra flags needed.
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
        </>
        )}
      </div>
    </div>
  )
}
