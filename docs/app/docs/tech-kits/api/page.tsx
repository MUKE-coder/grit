import { Database, Shield, Bot, Activity, Cpu, FileCode2, Container as DockerIcon, FileJson, Terminal as TerminalIcon, Layers } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { BrowserFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'
import { buildKitPrompt } from '@/lib/kit-prompts'

export const metadata = getDocMetadata('/docs/tech-kits/api')

export default function ApiKitPage() {
  return (
    <TechKitLayout
      name="API — Go backend only"
      tagline="Pure Go API. Bring your own frontend — mobile, native desktop, Discord bot, anything."
      pitch={
        <>
          <p>
            <code>grit new --api</code> ships zero frontend code. Pure Gin + GORM + the
            framework&apos;s batteries. You get JWT + OAuth + 2FA, file storage, jobs, AI, audit
            log, and an OpenAPI spec auto-served at <code>/docs</code>. The fastest scaffold to
            deploy when you already have a frontend or want one in a separate repo.
          </p>
          <IncludedRow>Pure Gin + GORM; no front-end</IncludedRow>
          <IncludedRow>OpenAPI 3.0 auto-generated from routes + models</IncludedRow>
          <IncludedRow>JWT + OAuth (Google, GitHub) + TOTP 2FA</IncludedRow>
          <IncludedRow>Scalar / Swagger / Postman / Insomnia exports at /docs</IncludedRow>
        </>
      }
      command="grit new my-app --api"
      mockup={
        <BrowserFrame url="api.my-app.com/docs" glow>
          <div className="h-[260px] flex">
            <div className="w-28 border-r border-border/40 p-2 bg-card/40">
              <p className="text-[8px] uppercase font-mono text-muted-foreground tracking-wider mb-1.5">/api</p>
              {[
                { l: 'POST /auth/login', m: 'text-emerald-500' },
                { l: 'GET /users', m: 'text-sky-500' },
                { l: 'GET /posts', m: 'text-sky-500' },
                { l: 'POST /posts', m: 'text-emerald-500' },
                { l: 'PUT /posts/:id', m: 'text-amber-500' },
                { l: 'DELETE /posts/:id', m: 'text-rose-500' },
              ].map((r, i) => (
                <div key={i} className={`text-[8.5px] font-mono ${r.m} mb-0.5 truncate`}>{r.l}</div>
              ))}
            </div>
            <div className="flex-1 p-3 font-mono text-[10px]">
              <p className="text-foreground/80 mb-1">POST /api/posts</p>
              <pre className="rounded bg-card/60 border border-border/40 p-2 text-foreground/70 leading-relaxed">{`{
  "title": "string",
  "body": "string",
  "tags": ["string"]
}`}</pre>
              <p className="text-[8.5px] text-muted-foreground mt-2">200 OK · 401 Unauthorized · 422 Validation</p>
            </div>
          </div>
        </BrowserFrame>
      }
      features={[
        { icon: TerminalIcon, title: 'Gin + GORM',           body: 'The reference Go stack; fast, idiomatic, batteries included.' },
        { icon: FileJson,     title: 'OpenAPI auto-spec',    body: 'gin-docs reflects routes + GORM models — no annotations needed.' },
        { icon: Shield,       title: 'Auth + 2FA',           body: 'JWT, OAuth, TOTP — exactly the same auth as the full kits.' },
        { icon: Database,     title: 'Postgres + GORM',      body: 'AutoMigrate + GORM Studio at /studio.' },
        { icon: Layers,       title: 'Sentinel + Pulse',     body: 'Mount /sentinel and /pulse for security + observability.' },
        { icon: Activity,     title: 'Background jobs',      body: 'Asynq queue + cron + webhook receiver.' },
        { icon: Bot,          title: 'AI Gateway',           body: 'Stream from 100+ models with one API key.' },
        { icon: Cpu,          title: 'Code generator',       body: 'grit generate resource → model + service + handler + routes.' },
        { icon: DockerIcon,   title: 'Single Dockerfile',    body: 'Multi-stage build → ~20MB final image.' },
        { icon: FileCode2,    title: 'Audit log',            body: 'Tamper-evident SHA-256 chain on every mutation.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/api-only"
      starterPrompt={buildKitPrompt('api')}
      prev={{ label: 'Triple', href: '/docs/tech-kits/triple' }}
      next={{ label: 'Mobile', href: '/docs/tech-kits/mobile' }}
    />
  )
}
