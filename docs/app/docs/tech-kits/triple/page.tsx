import { Server, Shield, Database, Layers, Bot, Activity, Cpu, FileCode2, Container as DockerIcon, Globe, Users } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { BrowserFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tech-kits/triple')

export default function TripleKitPage() {
  return (
    <TechKitLayout
      name="Triple — Web + Admin + API"
      tagline="The full SaaS shape. Marketing site, Filament-style admin panel, Go API — one monorepo."
      pitch={
        <>
          <p>
            <code>grit new --triple</code> is the recommended starter for serious SaaS products.
            Three apps in one Turborepo: <code>apps/web</code> (public marketing + product UI),
            <code> apps/admin</code> (Filament-style admin panel), and <code>apps/api</code>{' '}
            (Go API). RBAC, invitation flow, audit log, and the in-app Security + Observability
            dashboards all come pre-wired.
          </p>
          <IncludedRow>apps/web + apps/admin + apps/api in one Turborepo</IncludedRow>
          <IncludedRow>Filament-style admin panel with DataTable, FormBuilder, dashboard widgets</IncludedRow>
          <IncludedRow>RBAC (admin / manager / staff) + workspace switcher</IncludedRow>
          <IncludedRow>Invitation flow + audit log + in-app Security/Observability admin pages</IncludedRow>
        </>
      }
      command="grit new my-app --triple"
      mockup={
        <BrowserFrame url="admin.my-app.com" glow>
          <div className="h-[260px] flex">
            <div className="w-24 border-r border-border/40 p-2 bg-card/40">
              <p className="text-[8px] uppercase font-mono text-muted-foreground tracking-wider mb-1.5">Admin</p>
              {[
                { l: 'Dashboard', sel: true },
                { l: 'Users', sel: false },
                { l: 'Orders', sel: false },
                { l: 'Sales', sel: false },
                { l: 'Settings', sel: false },
              ].map((r) => (
                <div key={r.l} className={`text-[9px] rounded px-1.5 py-0.5 mb-0.5 ${r.sel ? 'bg-primary/15 text-primary' : 'text-muted-foreground'}`}>
                  {r.l}
                </div>
              ))}
            </div>
            <div className="flex-1 p-3">
              <div className="flex items-center gap-1.5 mb-2">
                <Users className="h-3 w-3 text-primary" />
                <span className="text-xs font-semibold">Users · 1,284</span>
              </div>
              <div className="space-y-1">
                {['Maya · admin', 'Carlos · manager', 'Stella · staff', 'Lawrence · staff'].map((r, i) => (
                  <div key={i} className="flex items-center justify-between gap-1 text-[10px] rounded border border-border/40 bg-card/60 px-2 py-1">
                    <span className="truncate">{r}</span>
                    <span className="text-[8px] font-mono text-emerald-500">active</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </BrowserFrame>
      }
      features={[
        { icon: Globe,      title: 'Public web app',     body: 'apps/web — Next.js marketing + product UI.' },
        { icon: Server,     title: 'Admin panel',        body: 'apps/admin — Filament-style DataTable + FormBuilder + widgets.' },
        { icon: Shield,     title: 'RBAC + invitations', body: 'Roles, workspace switcher, accept-by-token invitation flow.' },
        { icon: Layers,     title: 'In-app Security',    body: 'Sentinel summary inside admin: score, threats, AuthShield, CSP.' },
        { icon: Activity,   title: 'In-app Observability', body: 'Pulse summary inside admin: p95/p99, SLOs, USE grid, top N+1.' },
        { icon: Database,   title: 'Postgres + GORM',    body: 'Production DB, GORM Studio at /studio.' },
        { icon: Cpu,        title: 'Code generator',     body: 'grit generate resource emits to all three apps.' },
        { icon: Bot,        title: 'AI Gateway',         body: 'Streaming completions, agent endpoints, structured outputs.' },
        { icon: FileCode2,  title: 'Audit + notifications', body: 'Tamper-evident log + bell notifications for high-sev events.' },
        { icon: DockerIcon, title: 'Three Dockerfiles',  body: 'One per app; docker-compose for local + grit deploy for prod.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/triple"
      prev={{ label: 'Double', href: '/docs/tech-kits/double' }}
      next={{ label: 'API only', href: '/docs/tech-kits/api' }}
    />
  )
}
