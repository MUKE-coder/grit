import { Zap, Shield, Database, Layers, Bot, Activity, Cpu, FileCode2, Container as DockerIcon } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { BrowserFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'
import { buildKitPrompt } from '@/lib/kit-prompts'

export const metadata = getDocMetadata('/docs/tech-kits/single')

export default function SingleKitPage() {
  return (
    <TechKitLayout
      name="Single — Go + embedded SPA"
      tagline="One binary. Your React frontend is embedded into the Go server via go:embed."
      pitch={
        <>
          <p>
            The smallest possible Grit deploy. <code>grit new --single</code> scaffolds a single
            Go module with the React frontend bundled into the same binary. No CORS, no reverse
            proxy, no separate front-end server — drop the binary on a box, run it, done.
          </p>
          <IncludedRow>Single Go binary; production deploy = one file</IncludedRow>
          <IncludedRow>React (Next.js) frontend embedded via <code>//go:embed</code></IncludedRow>
          <IncludedRow>SPA fallback with the redirect-loop-safe pattern wired</IncludedRow>
          <IncludedRow>Auth, jobs, storage, AI, Pulse + Sentinel — all the same batteries</IncludedRow>
        </>
      }
      command="grit new my-app --single"
      mockup={
        <BrowserFrame url="my-app.com" glow>
          <div className="h-[260px] p-4 flex flex-col">
            <div className="flex items-center gap-2 mb-3 pb-2 border-b border-border/40">
              <div className="h-4 w-4 rounded bg-gradient-to-br from-sky-400 to-blue-500" />
              <span className="text-xs font-semibold">My App</span>
              <span className="ml-auto text-[10px] font-mono text-muted-foreground">localhost:8080</span>
            </div>
            <div className="grid grid-cols-2 gap-2 mb-2">
              <div className="rounded border border-border/40 bg-card/60 p-2">
                <p className="text-[9px] text-muted-foreground uppercase tracking-wider">USERS</p>
                <p className="text-base font-semibold tabular-nums">1,284</p>
              </div>
              <div className="rounded border border-border/40 bg-card/60 p-2">
                <p className="text-[9px] text-muted-foreground uppercase tracking-wider">RPS</p>
                <p className="text-base font-semibold tabular-nums">42.1</p>
              </div>
            </div>
            <div className="rounded border border-border/40 bg-card/40 p-2 mt-auto">
              <div className="flex items-end gap-[2px] h-8">
                {[40, 55, 48, 62, 70, 65, 80, 72, 88, 84].map((h, i) => (
                  <div key={i} className="flex-1 rounded-sm bg-primary/40" style={{ height: `${h}%` }} />
                ))}
              </div>
            </div>
          </div>
        </BrowserFrame>
      }
      features={[
        { icon: Zap,        title: 'One binary',        body: 'Go server + embedded React bundle. Single file deploy.' },
        { icon: Shield,     title: 'Auth + 2FA',         body: 'JWT, OAuth (Google, GitHub), TOTP, trusted devices.' },
        { icon: Database,   title: 'GORM + Postgres',    body: 'Production DB by default; SQLite for local dev.' },
        { icon: Layers,     title: 'Sentinel + Pulse',   body: 'In-app Security + Observability admin pages.' },
        { icon: Bot,        title: 'AI Gateway',         body: '100+ models, streaming, one API key.' },
        { icon: Activity,   title: 'Background jobs',    body: 'Asynq queue, cron, webhook receiver, idempotency.' },
        { icon: Cpu,        title: 'Code generator',     body: 'grit generate resource → model + service + handler + admin page.' },
        { icon: FileCode2,  title: 'Audit log',          body: 'Tamper-evident SHA-256 hash chain on every mutation.' },
        { icon: DockerIcon, title: 'Docker + deploy',    body: 'Multi-stage Dockerfile + grit deploy --domain.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/single"
      starterPrompt={buildKitPrompt('single')}
      next={{ label: 'Single + Vite', href: '/docs/tech-kits/single-vite' }}
    />
  )
}
