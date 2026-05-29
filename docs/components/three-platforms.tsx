'use client'

// "Three Platforms. One Framework." — Grit's killer superpower made
// visible. Same Go API, same generated code, three different targets:
// a web browser (Next.js / Vite SPA), a desktop window (Wails), and a
// phone (Expo). Each device frame renders a realistic mockup of the
// same product so the parity message is undeniable.

import {
  Globe,
  Monitor,
  Smartphone,
  Activity,
  Users,
  ArrowRight,
  Search,
  Bell,
  ChevronRight,
} from 'lucide-react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { BrowserFrame, DesktopFrame, MobileFrame } from '@/components/ui/device-frames'

// Mini product mockup that lives inside each frame — the "Acme Dashboard"
// app. Renders three flavours of the same view depending on the device.
function DashboardMockup({ variant }: { variant: 'web' | 'desktop' | 'mobile' }) {
  if (variant === 'mobile') {
    return (
      <div className="px-4">
        <div className="flex items-center justify-between mb-3">
          <div>
            <p className="text-[9px] text-muted-foreground uppercase tracking-wider">Good morning</p>
            <p className="text-sm font-semibold text-foreground">Acme — Sales</p>
          </div>
          <div className="h-7 w-7 rounded-full bg-primary/15 flex items-center justify-center">
            <Bell className="h-3 w-3 text-primary" />
          </div>
        </div>
        {/* Stat cards stacked vertically */}
        <div className="space-y-2 mb-3">
          {[
            { label: 'Revenue', value: '$84.2K', tone: 'text-emerald-500' },
            { label: 'New users', value: '1,284', tone: 'text-sky-500' },
          ].map((s) => (
            <div key={s.label} className="rounded-lg border border-border/50 bg-card/60 p-2.5">
              <p className="text-[9px] uppercase tracking-wider text-muted-foreground">{s.label}</p>
              <p className={`text-base font-semibold ${s.tone} tabular-nums`}>{s.value}</p>
            </div>
          ))}
        </div>
        {/* Activity rows */}
        <p className="text-[9px] text-muted-foreground uppercase tracking-wider mb-1.5">Activity</p>
        <div className="space-y-1.5">
          {['Maya signed in', 'Invoice #1284 paid', 'New review (4.8★)'].map((row, i) => (
            <div key={i} className="flex items-center gap-2 py-1 text-[10px] text-foreground/80">
              <span className={`h-1 w-1 rounded-full ${i === 0 ? 'bg-emerald-500' : i === 1 ? 'bg-sky-500' : 'bg-amber-500'}`} />
              <span className="flex-1 truncate">{row}</span>
            </div>
          ))}
        </div>
      </div>
    )
  }

  if (variant === 'desktop') {
    return (
      <div className="flex h-full">
        {/* Sidebar */}
        <div className="w-28 border-r border-border/60 p-2.5 bg-card/40">
          <div className="flex items-center gap-1.5 mb-3">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src="/grit_logo.png" alt="" className="h-3.5 w-3.5 rounded" />
            <span className="text-[10px] font-semibold text-foreground">Acme</span>
          </div>
          <div className="space-y-0.5">
            {[
              { l: 'Dashboard', sel: true },
              { l: 'Customers', sel: false },
              { l: 'Sales', sel: false },
              { l: 'Reports', sel: false },
              { l: 'Settings', sel: false },
            ].map((row) => (
              <div
                key={row.l}
                className={`rounded px-1.5 py-0.5 text-[9px] ${
                  row.sel ? 'bg-primary/10 text-primary' : 'text-muted-foreground'
                }`}
              >
                {row.l}
              </div>
            ))}
          </div>
        </div>
        {/* Body */}
        <div className="flex-1 p-3">
          <div className="flex items-center justify-between mb-2.5">
            <p className="text-sm font-semibold text-foreground">Dashboard</p>
            <div className="text-[9px] font-mono text-muted-foreground">offline · synced</div>
          </div>
          <div className="grid grid-cols-3 gap-2 mb-3">
            {[
              { label: 'Revenue', value: '$84.2K', dot: 'bg-emerald-500' },
              { label: 'Users', value: '1,284', dot: 'bg-sky-500' },
              { label: 'NPS', value: '72', dot: 'bg-violet-500' },
            ].map((s) => (
              <div key={s.label} className="rounded border border-border/50 bg-card/60 p-2">
                <div className="flex items-center gap-1 mb-0.5">
                  <span className={`h-1 w-1 rounded-full ${s.dot}`} />
                  <p className="text-[8px] uppercase tracking-wider text-muted-foreground">{s.label}</p>
                </div>
                <p className="text-xs font-semibold text-foreground tabular-nums">{s.value}</p>
              </div>
            ))}
          </div>
          {/* Mini chart */}
          <div className="rounded border border-border/50 bg-card/40 p-2">
            <div className="flex items-end gap-[2px] h-10">
              {[40, 55, 48, 62, 58, 70, 75, 62, 80, 72, 88, 84].map((h, i) => (
                <div
                  key={i}
                  className="flex-1 rounded-sm bg-gradient-to-t from-primary/60 to-primary/30"
                  style={{ height: `${h}%` }}
                />
              ))}
            </div>
            <p className="text-[8px] font-mono text-muted-foreground mt-1">Last 12 days · trending +18%</p>
          </div>
        </div>
      </div>
    )
  }

  // web variant
  return (
    <div className="p-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between mb-3 pb-2 border-b border-border/40">
        <div className="flex items-center gap-2">
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img src="/grit_logo.png" alt="" className="h-4 w-4 rounded" />
          <span className="text-xs font-semibold text-foreground">Acme · Dashboard</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex items-center gap-1 rounded-md border border-border/40 px-1.5 py-0.5">
            <Search className="h-2.5 w-2.5 text-muted-foreground" />
            <span className="text-[9px] text-muted-foreground">Search…</span>
          </div>
          <div className="h-5 w-5 rounded-full bg-gradient-to-br from-violet-400 to-rose-400" />
        </div>
      </div>
      {/* Stats row */}
      <div className="grid grid-cols-4 gap-2 mb-3">
        {[
          { label: 'REVENUE', value: '$84.2K', sub: '+12%', tone: 'text-emerald-500' },
          { label: 'USERS', value: '1,284', sub: '+8%', tone: 'text-sky-500' },
          { label: 'CONVERSION', value: '3.2%', sub: '+0.4', tone: 'text-violet-500' },
          { label: 'NPS', value: '72', sub: '+5', tone: 'text-amber-500' },
        ].map((s) => (
          <div key={s.label} className="rounded-md border border-border/40 bg-card/60 p-1.5">
            <p className="text-[8px] uppercase tracking-wider text-muted-foreground">{s.label}</p>
            <p className="text-xs font-semibold text-foreground tabular-nums">{s.value}</p>
            <p className={`text-[8px] font-mono ${s.tone}`}>▲ {s.sub}</p>
          </div>
        ))}
      </div>
      {/* Chart + table */}
      <div className="grid grid-cols-3 gap-2">
        <div className="col-span-2 rounded-md border border-border/40 bg-card/40 p-2">
          <p className="text-[9px] uppercase tracking-wider text-muted-foreground mb-1">Revenue · 30 days</p>
          <svg viewBox="0 0 200 50" className="w-full h-10">
            <defs>
              <linearGradient id="rev-area" x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="hsl(var(--primary))" stopOpacity="0.4" />
                <stop offset="100%" stopColor="hsl(var(--primary))" stopOpacity="0" />
              </linearGradient>
            </defs>
            <path
              d="M 0 40 L 20 32 L 40 35 L 60 28 L 80 22 L 100 25 L 120 18 L 140 20 L 160 12 L 180 14 L 200 8 L 200 50 L 0 50 Z"
              fill="url(#rev-area)"
            />
            <path
              d="M 0 40 L 20 32 L 40 35 L 60 28 L 80 22 L 100 25 L 120 18 L 140 20 L 160 12 L 180 14 L 200 8"
              fill="none"
              stroke="hsl(var(--primary))"
              strokeWidth="1.5"
            />
          </svg>
        </div>
        <div className="rounded-md border border-border/40 bg-card/40 p-2">
          <p className="text-[9px] uppercase tracking-wider text-muted-foreground mb-1">Activity</p>
          <div className="space-y-1">
            {['Maya', 'Carlos', 'Inez'].map((n, i) => (
              <div key={n} className="flex items-center gap-1.5 text-[9px] text-foreground/80">
                <span className={`h-3 w-3 rounded-full ${
                  i === 0 ? 'bg-emerald-400' : i === 1 ? 'bg-sky-400' : 'bg-violet-400'
                } flex items-center justify-center text-[7px] text-white font-bold`}>
                  {n[0]}
                </span>
                <span className="truncate">{n}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

export function ThreePlatforms() {
  return (
    <section className="relative py-24 px-6 overflow-hidden border-t border-border/40">
      {/* Layered backdrop */}
      <div className="absolute inset-0 -z-30 bg-gradient-to-b from-background via-card/20 to-background" />
      <div
        className="absolute inset-0 -z-20 opacity-[0.4]"
        style={{
          backgroundImage:
            'radial-gradient(circle at 1px 1px, hsl(var(--foreground) / 0.06) 1px, transparent 0)',
          backgroundSize: '32px 32px',
        }}
      />
      <div className="absolute left-1/2 top-1/3 -translate-x-1/2 -z-10 h-[500px] w-[700px] rounded-full bg-primary/[0.05] blur-[120px]" />

      <div className="max-w-6xl mx-auto">

        {/* Heading */}
        <div className="text-center mb-12">
          <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-5">
            <span className="h-1.5 w-1.5 rounded-full bg-primary" /> Grit&apos;s Superpower
          </span>
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5 leading-[1.1]">
            One framework.<br className="md:hidden" />{' '}
            <span className="bg-gradient-to-r from-primary via-sky-400 to-cyan-400 bg-clip-text text-transparent">
              Three platforms.
            </span>
          </h2>
          <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            The same Go API and the same generated code ship as a web app, a desktop binary,
            and a mobile app. Pick one. Or ship all three from the same repo.
          </p>
        </div>

        {/* Platform stat row */}
        <div className="grid grid-cols-3 gap-4 max-w-3xl mx-auto mb-14">
          {[
            { icon: Globe,      label: 'Web',     stack: 'Next.js · Vite · TanStack', flag: '--single / --triple' },
            { icon: Monitor,    label: 'Desktop', stack: 'Wails · Frameless · GORM',  flag: '--desktop' },
            { icon: Smartphone, label: 'Mobile',  stack: 'Expo · React Native',       flag: '--mobile' },
          ].map((p) => {
            const Icon = p.icon
            return (
              <div
                key={p.label}
                className="rounded-xl border-2 border-border/50 bg-card/60 backdrop-blur p-4 text-center hover:border-primary/40 transition-colors"
              >
                <div className="inline-flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary mb-2">
                  <Icon className="h-5 w-5" />
                </div>
                <p className="font-semibold text-foreground text-sm">{p.label}</p>
                <p className="text-[11px] text-muted-foreground mt-0.5">{p.stack}</p>
                <code className="text-[9px] font-mono text-primary/80 bg-primary/5 px-1.5 py-0.5 rounded mt-1.5 inline-block">{p.flag}</code>
              </div>
            )
          })}
        </div>

        {/* Three device frames showing the SAME product */}
        <div className="grid grid-cols-1 lg:grid-cols-[1.4fr_1fr_0.5fr] gap-6 items-end">
          {/* Browser (largest) */}
          <div>
            <p className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground mb-2 ml-1">
              Web — built via <span className="text-foreground">grit new --triple --vite</span>
            </p>
            <BrowserFrame url="acme.app" glow>
              <div className="h-[260px]">
                <DashboardMockup variant="web" />
              </div>
            </BrowserFrame>
          </div>

          {/* Desktop (medium) */}
          <div>
            <p className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground mb-2 ml-1">
              Desktop — <span className="text-foreground">grit new-desktop</span>
            </p>
            <DesktopFrame title="Acme — Desktop" glow>
              <div className="h-[260px]">
                <DashboardMockup variant="desktop" />
              </div>
            </DesktopFrame>
          </div>

          {/* Mobile (smallest, scaled down) */}
          <div className="flex flex-col items-center lg:items-start">
            <p className="text-[10px] font-mono uppercase tracking-wider text-muted-foreground mb-2 self-start ml-1">
              Mobile — <span className="text-foreground">--mobile</span>
            </p>
            <MobileFrame glow className="scale-[0.85] origin-top-left">
              <div className="h-[280px]">
                <DashboardMockup variant="mobile" />
              </div>
            </MobileFrame>
          </div>
        </div>

        {/* Bottom CTA strip */}
        <div className="mt-14 flex flex-col sm:flex-row items-center justify-between gap-4 rounded-2xl border-2 border-border/50 bg-card/60 p-6">
          <div className="flex items-center gap-4">
            <div className="hidden sm:flex h-10 w-10 rounded-lg bg-primary/10 text-primary items-center justify-center">
              <Users className="h-5 w-5" />
            </div>
            <div>
              <p className="font-semibold text-foreground">Same Go backend. Same generated React. Same auth.</p>
              <p className="text-sm text-muted-foreground">Pick the targets you need — the API doesn&apos;t change.</p>
            </div>
          </div>
          <Button asChild variant="outline" className="border-border/60 rounded-full">
            <Link href="/docs/concepts/architecture-modes">
              Compare modes <ArrowRight className="ml-2 h-3.5 w-3.5" />
            </Link>
          </Button>
        </div>

      </div>
    </section>
  )
}

// Tiny re-export for chevrons used in some icon rows (kept here so the
// section is self-contained when imported elsewhere).
export { ChevronRight as _ChevronRight, Activity as _Activity }
