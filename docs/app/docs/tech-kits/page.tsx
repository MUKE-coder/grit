import Link from 'next/link'
import { ArrowRight, Zap, Layers, Server, Database, Smartphone, Monitor, Cpu } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tech-kits')

interface KitCard {
  slug: string
  name: string
  tagline: string
  command: string
  icon: React.ComponentType<{ className?: string }>
  accentClass: string
  bullets: string[]
}

const KITS: KitCard[] = [
  {
    slug: 'single',
    name: 'Single — Go + embedded SPA',
    tagline: 'One binary. Embed your React frontend in the Go server. Smallest possible deploy.',
    command: 'grit new my-app --single',
    icon: Zap,
    accentClass: 'from-sky-500/20 to-cyan-500/10 text-sky-400',
    bullets: ['Single Go binary', 'React (Next or Vite) embedded via go:embed', 'No CORS, no reverse proxy', 'Perfect for solo / indie ship'],
  },
  {
    slug: 'single-vite',
    name: 'Single — Vite + TanStack Router',
    tagline: 'Same single-binary shape but with a Vite SPA + TanStack Router instead of Next.js.',
    command: 'grit new my-app --single --vite',
    icon: Zap,
    accentClass: 'from-cyan-500/20 to-blue-500/10 text-cyan-400',
    bullets: ['Vite + TanStack Router', 'Faster cold dev start', 'Server-rendered nothing — pure SPA', 'Drop-in for Inertia-style apps'],
  },
  {
    slug: 'double',
    name: 'Double — Web + API monorepo',
    tagline: 'Two apps in a Turborepo. Public Next.js frontend talking to a separate Go API.',
    command: 'grit new my-app --double',
    icon: Layers,
    accentClass: 'from-violet-500/20 to-fuchsia-500/10 text-violet-400',
    bullets: ['apps/web (Next.js) + apps/api (Go)', 'Shared TS types in packages/shared', 'CORS pre-wired', 'Deploy each on its own schedule'],
  },
  {
    slug: 'triple',
    name: 'Triple — Web + Admin + API',
    tagline: 'The full SaaS shape. Marketing site, Filament-style admin, Go API — all in one monorepo.',
    command: 'grit new my-app --triple',
    icon: Server,
    accentClass: 'from-emerald-500/20 to-green-500/10 text-emerald-400',
    bullets: ['apps/web + apps/admin + apps/api', 'Filament-style admin panel', 'RBAC + invitation flow shipped', 'The recommended SaaS starter'],
  },
  {
    slug: 'api',
    name: 'API — Go backend only',
    tagline: 'Headless Go API. Bring your own frontend — mobile, native desktop, Discord bot.',
    command: 'grit new my-app --api',
    icon: Database,
    accentClass: 'from-amber-500/20 to-orange-500/10 text-amber-400',
    bullets: ['No frontend; pure Gin + GORM', 'OpenAPI docs at /docs', 'JWT + OAuth + 2FA', 'Smallest scaffold; fastest to deploy'],
  },
  {
    slug: 'mobile',
    name: 'Mobile — Expo + API',
    tagline: 'React Native (Expo) frontend talking to a Grit API. Same generated types front + back.',
    command: 'grit new my-app --mobile',
    icon: Smartphone,
    accentClass: 'from-rose-500/20 to-pink-500/10 text-rose-400',
    bullets: ['Expo (React Native) frontend', 'Shared Zod schemas + TS types', 'Mobile-friendly auth flow', 'EAS Build / OTA ready'],
  },
  {
    slug: 'desktop',
    name: 'Desktop — Wails + GORM',
    tagline: 'Native desktop app via Wails. Offline-first SQLite, frameless window, draggable panels.',
    command: 'grit new-desktop my-app',
    icon: Monitor,
    accentClass: 'from-indigo-500/20 to-violet-500/10 text-indigo-400',
    bullets: ['Wails v2 + React + Tailwind', 'Local SQLite / Postgres', 'Local auth (bcrypt)', 'Export PDF / Excel built-in'],
  },
]

export default function TechKitsIndexPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-5xl mx-auto">

            {/* Header */}
            <div className="mb-12">
              <span className="tag-mono text-primary/80 mb-3 block">Tech Kits</span>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4 leading-tight">
                The fastest way to start your<br className="hidden md:block" /> next application
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-3xl">
                With authentication, registration, settings, and best practices built in. Our tech
                kits let you skip the setup and start building what really matters. Pick the kit
                that matches your shape — single binary, monorepo, mobile, desktop, or pure API —
                and ship in minutes.
              </p>
            </div>

            {/* Featured kits — first 4 in a 2x2 grid, with bigger framing */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              {KITS.slice(0, 4).map((kit) => (
                <KitCardLarge key={kit.slug} kit={kit} />
              ))}
            </div>

            {/* Remaining kits — 3 col strip */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-16">
              {KITS.slice(4).map((kit) => (
                <KitCardSmall key={kit.slug} kit={kit} />
              ))}
            </div>

            {/* Auth options section — Laravel cribs this from its Built-in vs WorkOS layout */}
            <div className="mb-16">
              <h2 className="text-2xl font-semibold tracking-tight mb-2 flex items-center gap-2">
                <Cpu className="h-5 w-5 text-primary" />
                An authenticated application out of the box
              </h2>
              <p className="text-muted-foreground mb-6">
                Authentication options to suit your application&apos;s needs — flip a flag in
                <code className="text-primary text-sm bg-primary/5 px-1.5 py-0.5 rounded mx-1">.env</code>
                or compose your own.
              </p>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="rounded-xl border border-border/50 bg-card/50 p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <div className="h-1.5 w-1.5 rounded-full bg-primary" />
                    <p className="text-xs font-mono text-primary uppercase tracking-wider">Default</p>
                  </div>
                  <h3 className="text-lg font-semibold mb-1.5">Built-in Grit auth</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Robust and secure authentication: registration, login, password reset, email
                    verification, JWT access + refresh tokens, OAuth (Google / GitHub), and TOTP
                    two-factor with trusted-device cookies.
                  </p>
                </div>

                <div className="rounded-xl border border-border/50 bg-card/50 p-5">
                  <div className="flex items-center gap-2 mb-2">
                    <div className="h-1.5 w-1.5 rounded-full bg-amber-400" />
                    <p className="text-xs font-mono text-amber-500 uppercase tracking-wider">Optional</p>
                  </div>
                  <h3 className="text-lg font-semibold mb-1.5">WorkOS AuthKit + SSO</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Drop in WorkOS / Auth0 / Clerk for enterprise SSO, SAML, SCIM provisioning,
                    and magic links. The Grit auth middleware is interface-driven so the swap is
                    one config line.
                  </p>
                </div>
              </div>
            </div>

            {/* Ready-to-build strip */}
            <div className="rounded-2xl border border-primary/20 bg-primary/[0.04] p-8 mb-12 text-center">
              <h3 className="text-2xl font-semibold mb-2">Ready to build?</h3>
              <p className="text-muted-foreground mb-6">
                Create the incredible, with Grit.
              </p>
              <div className="flex flex-wrap justify-center gap-2.5">
                <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
                  <Link href="/docs/getting-started/quick-start">
                    Quick start <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
                  </Link>
                </Button>
                <Button asChild variant="outline" className="border-border/60 rounded-full">
                  <Link href="/docs/stack-selector">Use the stack selector</Link>
                </Button>
              </div>
            </div>

          </div>
        </div>
      </main>
    </div>
  )
}

function KitCardLarge({ kit }: { kit: KitCard }) {
  const Icon = kit.icon
  return (
    <Link
      href={`/docs/tech-kits/${kit.slug}`}
      className="group relative rounded-2xl border border-border/50 bg-card/50 p-6 hover:border-primary/40 hover:bg-card/80 transition-all overflow-hidden"
    >
      {/* Decorative gradient blob */}
      <div className={`absolute -top-12 -right-12 h-32 w-32 rounded-full bg-gradient-to-br ${kit.accentClass} blur-2xl opacity-60 pointer-events-none`} />

      <div className="relative">
        <div className="flex items-center gap-3 mb-3">
          <div className={`h-10 w-10 rounded-lg bg-gradient-to-br ${kit.accentClass} flex items-center justify-center`}>
            <Icon className="h-5 w-5" />
          </div>
          <h3 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors">
            {kit.name}
          </h3>
        </div>
        <p className="text-sm text-muted-foreground leading-relaxed mb-4">{kit.tagline}</p>
        <code className="block text-[11px] font-mono text-primary/80 bg-primary/5 border border-primary/10 rounded-md px-2.5 py-1.5 mb-4">
          {kit.command}
        </code>
        <div className="inline-flex items-center gap-1 text-xs font-mono text-primary/80 group-hover:text-primary">
          View kit <ArrowRight className="h-3 w-3 transition-transform group-hover:translate-x-0.5" />
        </div>
      </div>
    </Link>
  )
}

function KitCardSmall({ kit }: { kit: KitCard }) {
  const Icon = kit.icon
  return (
    <Link
      href={`/docs/tech-kits/${kit.slug}`}
      className="group rounded-2xl border border-border/50 bg-card/40 p-5 hover:border-primary/40 hover:bg-card/80 transition-all"
    >
      <div className="flex items-center gap-2.5 mb-2">
        <div className={`h-8 w-8 rounded-md bg-gradient-to-br ${kit.accentClass} flex items-center justify-center`}>
          <Icon className="h-4 w-4" />
        </div>
        <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors leading-tight">
          {kit.name}
        </h3>
      </div>
      <p className="text-[13px] text-muted-foreground leading-relaxed mb-3 line-clamp-2">{kit.tagline}</p>
      <code className="block text-[10px] font-mono text-primary/80 bg-primary/5 rounded px-2 py-1">
        {kit.command}
      </code>
    </Link>
  )
}
