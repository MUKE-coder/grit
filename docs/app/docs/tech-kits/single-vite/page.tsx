import { Zap, Shield, Database, Layers, Bot, Activity, Cpu, FileCode2, Container as DockerIcon, Sparkles } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { BrowserFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tech-kits/single-vite')

export default function SingleViteKitPage() {
  return (
    <TechKitLayout
      name="Single — Vite + TanStack Router"
      tagline="Same single-binary shape, but with a Vite SPA and TanStack Router instead of Next.js."
      pitch={
        <>
          <p>
            For teams who prefer a pure SPA over Next.js: <code>grit new --single --vite</code>
            ships the same single-Go-binary deploy with a Vite-built TanStack Router frontend
            embedded. Cold dev starts are faster, the bundle is smaller, and you skip Next.js
            altogether.
          </p>
          <IncludedRow>Vite 6 + TanStack Router file-based routing</IncludedRow>
          <IncludedRow>Same Go API, same generated React Query hooks</IncludedRow>
          <IncludedRow>Sub-second dev server cold starts</IncludedRow>
          <IncludedRow>Pre-built <code>lib/auth.ts</code> with refresh-on-401 baked in</IncludedRow>
        </>
      }
      command="grit new my-app --single --vite"
      mockup={
        <BrowserFrame url="my-app.com" glow>
          <div className="h-[260px] p-4">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border/40">
              <Sparkles className="h-4 w-4 text-primary" />
              <span className="text-xs font-semibold">SaaS Dashboard · Vite</span>
            </div>
            <div className="grid grid-cols-4 gap-1.5 mb-3">
              {[
                { l: 'REV', v: '$84.2K', t: 'text-emerald-500' },
                { l: 'USERS', v: '1,284', t: 'text-sky-500' },
                { l: 'CONV', v: '3.2%', t: 'text-violet-500' },
                { l: 'NPS', v: '72', t: 'text-amber-500' },
              ].map((s) => (
                <div key={s.l} className="rounded border border-border/40 bg-card/60 p-1.5">
                  <p className="text-[7px] uppercase text-muted-foreground tracking-wider">{s.l}</p>
                  <p className={`text-xs font-semibold ${s.t}`}>{s.v}</p>
                </div>
              ))}
            </div>
            <div className="rounded border border-border/40 bg-card/40 p-2 h-24">
              <svg viewBox="0 0 200 50" className="w-full h-full">
                <defs>
                  <linearGradient id="va" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="0.4" />
                    <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0" />
                  </linearGradient>
                </defs>
                <path d="M 0 40 L 20 32 L 40 35 L 60 28 L 80 22 L 100 25 L 120 18 L 140 20 L 160 12 L 180 14 L 200 8 L 200 50 L 0 50 Z" fill="url(#va)" />
                <path d="M 0 40 L 20 32 L 40 35 L 60 28 L 80 22 L 100 25 L 120 18 L 140 20 L 160 12 L 180 14 L 200 8" fill="none" stroke="hsl(var(--primary))" strokeWidth="1.5" />
              </svg>
            </div>
          </div>
        </BrowserFrame>
      }
      features={[
        { icon: Zap,        title: 'Vite 6',             body: 'Sub-second cold dev starts; HMR you can feel.' },
        { icon: Sparkles,   title: 'TanStack Router',    body: 'File-based routing, type-safe links, search params as data.' },
        { icon: Shield,     title: 'Auth + refresh',     body: 'lib/auth.ts handles JWT refresh on 401 transparently.' },
        { icon: Database,   title: 'GORM + Postgres',    body: 'Same Go stack, no front-end framework lock-in.' },
        { icon: Layers,     title: 'Sentinel + Pulse',   body: 'In-app Security + Observability admin pages.' },
        { icon: Bot,        title: 'AI Gateway',         body: 'Stream from 100+ models with one API key.' },
        { icon: Activity,   title: 'Background jobs',    body: 'Asynq queue + cron + webhook receiver.' },
        { icon: Cpu,        title: 'Code generator',     body: 'grit generate resource targets TanStack hooks.' },
        { icon: DockerIcon, title: 'Docker + deploy',    body: 'Single binary, single Docker image, one TLS hop.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/single"
      prev={{ label: 'Single (Next)', href: '/docs/tech-kits/single' }}
      next={{ label: 'Double', href: '/docs/tech-kits/double' }}
    />
  )
}
