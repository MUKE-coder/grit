import { Layers, Shield, Database, Activity, Bot, Cpu, FileCode2, Container as DockerIcon, Globe } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { BrowserFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'
import { buildKitPrompt } from '@/lib/kit-prompts'

export const metadata = getDocMetadata('/docs/tech-kits/double')

export default function DoubleKitPage() {
  return (
    <TechKitLayout
      name="Double — Web + API monorepo"
      tagline="A public Next.js site talking to a separate Go API. One Turborepo, two deploy targets."
      pitch={
        <>
          <p>
            <code>grit new --double</code> gives you a Turborepo with{' '}
            <code>apps/web</code> (Next.js marketing + product site) and <code>apps/api</code>{' '}
            (Go API). Shared Zod schemas + TS types live in <code>packages/shared</code> so both
            sides stay typed end-to-end. Each app deploys on its own schedule.
          </p>
          <IncludedRow>apps/web (Next.js) + apps/api (Go)</IncludedRow>
          <IncludedRow>Shared Zod schemas + TS types in packages/shared</IncludedRow>
          <IncludedRow>CORS, JWT, cookie-vs-Authorization-header all pre-wired</IncludedRow>
          <IncludedRow>Deploy the web on Vercel, the API on Fly/Railway/your VPS</IncludedRow>
        </>
      }
      command="grit new my-app --double"
      mockup={
        <BrowserFrame url="my-app.com" glow>
          <div className="h-[260px] p-4 flex flex-col">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border/40">
              <div className="h-4 w-4 rounded bg-gradient-to-br from-violet-400 to-fuchsia-500" />
              <span className="text-xs font-semibold">apps/web</span>
              <span className="ml-auto text-[9px] font-mono text-muted-foreground">→ apps/api</span>
            </div>
            <div className="space-y-1.5 text-[10px] font-mono">
              {[
                { src: 'apps/web/app/page.tsx', dst: 'GET /api/posts' },
                { src: 'apps/web/app/blog/[slug]/page.tsx', dst: 'GET /api/posts/:slug' },
                { src: 'apps/web/app/(auth)/login/page.tsx', dst: 'POST /api/auth/login' },
                { src: 'packages/shared/schemas/post.ts', dst: 'shared' },
              ].map((r) => (
                <div key={r.src} className="flex items-center justify-between gap-2 rounded border border-border/40 bg-card/60 px-2 py-1">
                  <span className="text-foreground/80 truncate">{r.src}</span>
                  <span className="text-primary shrink-0">{r.dst}</span>
                </div>
              ))}
            </div>
          </div>
        </BrowserFrame>
      }
      features={[
        { icon: Globe,      title: 'Next.js 14 + App Router', body: 'Server components, ISR, route handlers; everything modern.' },
        { icon: Layers,     title: 'Turborepo',          body: 'Cached builds, single pnpm install, parallel dev.' },
        { icon: Shield,     title: 'Auth + 2FA',         body: 'JWT, OAuth, TOTP. The Go API owns auth; Next.js calls in.' },
        { icon: Database,   title: 'GORM + Postgres',    body: 'Production DB by default; SQLite for tests.' },
        { icon: Cpu,        title: 'Type sharing',       body: 'grit sync regenerates TS from Go structs.' },
        { icon: Activity,   title: 'Background jobs',    body: 'Asynq queue + cron + webhook receiver in the API.' },
        { icon: Bot,        title: 'AI Gateway',         body: '100+ models, streaming, server or client-initiated.' },
        { icon: FileCode2,  title: 'Audit log',          body: 'Tamper-evident SHA-256 chain on every mutation.' },
        { icon: DockerIcon, title: 'Two Dockerfiles',    body: 'One per app; compose ships for local dev.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/double"
      starterPrompt={buildKitPrompt('double')}
      prev={{ label: 'Single + Vite', href: '/docs/tech-kits/single-vite' }}
      next={{ label: 'Triple', href: '/docs/tech-kits/triple' }}
    />
  )
}
