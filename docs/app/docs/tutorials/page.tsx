import Link from 'next/link'
import {
  ArrowLeft,
  ArrowRight,
  Clock,
  BookOpen,
  Globe,
  Monitor,
  Smartphone,
  Wifi,
  ShieldCheck,
  Flag,
  Cable,
  FileSpreadsheet,
  Layers,
  Zap,
  Server,
  Dumbbell,
  FileText,
  Shield,
  Activity,
  Database,
  Code2,
  Rocket,
  Receipt,
  TestTube2,
  HardDrive,
  MessageSquare,
  CreditCard,
  GitBranch,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tutorials')

type Icon = React.ComponentType<{ className?: string }>

/* -- Guided tracks (multi-lesson) ----------------------------------------- */

const TRACKS: {
  title: string
  outline: string
  href: string
  icon: Icon
  cmd: string
  courses: number
  time: string
}[] = [
  {
    title: 'Grit Web',
    outline: 'Building Web Applications',
    href: '/courses/grit-web',
    icon: Globe,
    cmd: '$ grit new app --triple',
    courses: 9,
    time: '~4.5 hours',
  },
  {
    title: 'Grit Desktop',
    outline: 'Building Desktop Applications',
    href: '/courses/grit-desktop',
    icon: Monitor,
    cmd: '$ grit new app --desktop',
    courses: 6,
    time: '~3 hours',
  },
  {
    title: 'Grit Mobile',
    outline: 'Building Mobile Applications',
    href: '/courses/grit-mobile',
    icon: Smartphone,
    cmd: '$ grit new app --mobile',
    courses: 5,
    time: '~2.5 hours',
  },
]

/* -- Focused tutorials (standalone) --------------------------------------- */

const TUTORIALS: {
  title: string
  outline: string
  href: string
  icon: Icon
  cmd: string
  duration: string
}[] = [
  { title: 'Offline-First Desktop', outline: 'Local SQLite + outbox + Git-style sync', href: '/docs/desktop/offline', icon: Wifi, cmd: '$ grit new app --desktop', duration: '30 min' },
  { title: 'Audit Log + Hash Chain', outline: 'Tamper-evident activity tracking for SOC2', href: '/courses/audit-log', icon: ShieldCheck, cmd: 'GET /admin/audit-log', duration: '30 min' },
  { title: 'Feature Flags & A/B Testing', outline: 'Sticky bucketing, percentage rollouts, realtime push', href: '/courses/feature-flags', icon: Flag, cmd: 'flags.IsEnabled("beta")', duration: '30 min' },
  { title: 'Webhook Receiver', outline: 'Stripe / GitHub / HMAC verifiers + replay', href: '/courses/webhook-receiver', icon: Cable, cmd: 'POST /webhooks/:provider', duration: '30 min' },
  { title: 'CSV / Excel Export', outline: 'Auto-generated per resource via grit generate', href: '/courses/export', icon: FileSpreadsheet, cmd: '$ grit generate resource', duration: '30 min' },
  { title: 'Realtime + WebSocket Hub', outline: 'SendToUser, Broadcast, useRealtimeEvent', href: '/courses/realtime-chat', icon: Layers, cmd: 'hub.Broadcast(event)', duration: '30 min' },
  { title: 'Batteries Included', outline: 'Every feature that ships with Grit', href: '/courses/batteries', icon: Zap, cmd: '$ grit new app --triple', duration: '30 min' },
  { title: 'API-Only Masterclass', outline: 'Build & deploy a REST API with Go', href: '/courses/api-masterclass', icon: Server, cmd: '$ grit new api --api', duration: '30 min' },
  { title: 'Build a Fitness App', outline: 'Go API + Expo React Native', href: '/courses/mobile-fitness-app', icon: Dumbbell, cmd: '$ grit new fit --mobile', duration: '30 min' },
  { title: 'API Docs: Scalar & Swagger', outline: 'Auto-generated API documentation', href: '/courses/api-docs-scalar', icon: FileText, cmd: 'GET /docs', duration: '30 min' },
  { title: 'Security Deep Dive', outline: 'Auth, 2FA, WAF & rate limiting', href: '/courses/security-deep-dive', icon: Shield, cmd: 'sentinel.Protect()', duration: '30 min' },
  { title: 'Pulse Analytics', outline: 'Tracing, metrics & monitoring', href: '/courses/pulse-analytics', icon: Activity, cmd: 'GET /pulse', duration: '30 min' },
  { title: 'GORM Studio', outline: 'The visual database browser', href: '/courses/gorm-studio', icon: Database, cmd: '$ grit studio', duration: '30 min' },
  { title: 'React + Vite + Go', outline: 'Building with TanStack Router', href: '/courses/react-vite-go', icon: Code2, cmd: '$ grit new app --vite', duration: '30 min' },
  { title: 'Deployment Guide', outline: 'Dokploy, Orbita, VPS & Vercel', href: '/courses/deployment-guide', icon: Rocket, cmd: '$ dokploy deploy', duration: '30 min' },
  { title: 'Invoice Generator', outline: 'Desktop app with Wails + PDF export', href: '/courses/invoice-desktop', icon: Receipt, cmd: '$ wails build', duration: '30 min' },
  { title: 'Testing Your Grit App', outline: 'Go, Vitest & Playwright', href: '/courses/testing', icon: TestTube2, cmd: '$ go test ./...', duration: '30 min' },
  { title: 'Database Mastery', outline: 'GORM models, migrations & queries', href: '/courses/gorm-mastery', icon: HardDrive, cmd: 'db.AutoMigrate(&User{})', duration: '30 min' },
  { title: 'Real-Time Chat', outline: 'WebSockets with grit-websockets', href: '/courses/realtime-chat', icon: MessageSquare, cmd: 'useRealtimeEvent()', duration: '30 min' },
  { title: 'Stripe Payments', outline: 'Subscriptions & billing for SaaS', href: '/courses/stripe-payments', icon: CreditCard, cmd: 'stripe.Subscribe(plan)', duration: '30 min' },
  { title: 'CI/CD with GitHub Actions', outline: 'Automated testing & deployment', href: '/courses/cicd-github', icon: GitBranch, cmd: '$ gh workflow run ci', duration: '30 min' },
]

/* -- Reusable mockup thumbnail -------------------------------------------- */

function MockupThumb({ tag, cmd }: { tag: string; cmd: string }) {
  return (
    <div className="relative border-b border-border/40 bg-[#0d1117] px-4 py-3 overflow-hidden">
      <div className="flex items-center gap-1.5 mb-2">
        <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
        <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
        <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        <span className="ml-2 text-[10px] font-mono uppercase tracking-wider text-primary/50">{tag}</span>
      </div>
      <pre className="text-[11px] leading-5 font-mono text-slate-300/80 whitespace-pre-wrap">{cmd}</pre>
    </div>
  )
}

export default function TutorialsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-5xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Learn Grit</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Tutorials</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Learn Grit by building. Follow a{' '}
                <span className="text-foreground font-medium">guided track</span> end-to-end for web,
                desktop, or mobile &mdash; or grab a{' '}
                <span className="text-foreground font-medium">focused tutorial</span> and finish
                something concrete in about 30 minutes. All free and self-paced.
              </p>
            </div>

            {/* Guided tracks */}
            <div className="mb-12">
              <div className="flex items-center gap-2 mb-5">
                <BookOpen className="h-4 w-4 text-primary" />
                <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                  Guided tracks
                </h2>
              </div>
              <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
                {TRACKS.map((t) => (
                  <Link
                    key={t.href}
                    href={t.href}
                    className="group rounded-xl border border-border/50 bg-card/40 overflow-hidden hover:border-primary/30 transition-colors flex flex-col"
                  >
                    <MockupThumb tag="Track" cmd={t.cmd} />
                    <div className="p-5 flex flex-col flex-1">
                      <div className="flex items-center gap-2.5 mb-1.5">
                        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15">
                          <t.icon className="h-4 w-4 text-primary" />
                        </div>
                        <h3 className="text-base font-semibold flex items-center gap-1.5 group-hover:text-primary transition-colors">
                          {t.title}
                          <ArrowRight className="h-3.5 w-3.5 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                        </h3>
                      </div>
                      <p className="text-sm text-muted-foreground/70 leading-relaxed mb-4">{t.outline}</p>
                      <div className="mt-auto flex items-center gap-4 text-[11px] font-mono text-muted-foreground/70">
                        <span className="inline-flex items-center gap-1">
                          <BookOpen className="h-3 w-3 text-primary/60" />
                          {t.courses} lessons
                        </span>
                        <span className="inline-flex items-center gap-1">
                          <Clock className="h-3 w-3 text-primary/60" />
                          {t.time}
                        </span>
                      </div>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            {/* Focused tutorials */}
            <div>
              <div className="flex items-center gap-2 mb-5">
                <Zap className="h-4 w-4 text-primary" />
                <h2 className="text-sm font-semibold uppercase tracking-wider text-muted-foreground">
                  Focused tutorials
                </h2>
              </div>
              <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
                {TUTORIALS.map((t) => (
                  <Link
                    key={t.title}
                    href={t.href}
                    className="group rounded-xl border border-border/50 bg-card/40 overflow-hidden hover:border-primary/30 transition-colors flex flex-col"
                  >
                    <MockupThumb tag={t.duration} cmd={t.cmd} />
                    <div className="p-5 flex flex-col flex-1">
                      <div className="flex items-center gap-2.5 mb-1.5">
                        <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary/20 transition-colors">
                          <t.icon className="h-4 w-4" />
                        </div>
                        <h3 className="text-sm font-semibold flex items-center gap-1.5 group-hover:text-primary transition-colors">
                          {t.title}
                          <ArrowRight className="h-3.5 w-3.5 opacity-0 -translate-x-1 group-hover:opacity-100 group-hover:translate-x-0 transition-all" />
                        </h3>
                      </div>
                      <p className="text-sm text-muted-foreground/70 leading-relaxed mb-3">{t.outline}</p>
                      <span className="mt-auto inline-flex items-center gap-1 text-[11px] font-mono text-muted-foreground/60">
                        <Clock className="h-3 w-3" />
                        {t.duration}
                      </span>
                    </div>
                  </Link>
                ))}
              </div>
            </div>

            <div className="flex items-center justify-between border-t border-border pt-8 mt-12 max-w-3xl">
              <Button variant="ghost" asChild>
                <Link href="/courses" className="gap-2">
                  <ArrowLeft className="h-4 w-4" />
                  All courses
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/getting-started/create-a-project" className="gap-2">
                  Create a project
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
