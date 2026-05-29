import { Monitor, Shield, Database, FileCode2, FileSpreadsheet, Container as DockerIcon, Cpu, Wifi, Activity } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { DesktopFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/tech-kits/desktop')

export default function DesktopKitPage() {
  return (
    <TechKitLayout
      name="Desktop — Wails + GORM"
      tagline="Native desktop binary. Local SQLite, frameless window, draggable panels — all offline."
      pitch={
        <>
          <p>
            <code>grit new-desktop my-app</code> scaffolds a standalone Wails v2 desktop app
            with React + TypeScript + Tailwind on the front and Go + GORM (SQLite or Postgres)
            on the back. Local auth (bcrypt), blog + contact CRUD, export to PDF / Excel, and a
            polished frameless window are all there from the start.
          </p>
          <IncludedRow>Wails v2 + React + TypeScript + Tailwind</IncludedRow>
          <IncludedRow>Local SQLite by default; Postgres optional</IncludedRow>
          <IncludedRow>Local auth (bcrypt) — no server required</IncludedRow>
          <IncludedRow>Export to PDF + Excel out of the box</IncludedRow>
          <IncludedRow>Triple-mode: <code>--desktop --triple</code> ships web + admin + desktop together</IncludedRow>
        </>
      }
      command="grit new-desktop my-app"
      mockup={
        <DesktopFrame title="My App — Desktop" glow>
          <div className="h-[260px] flex">
            <div className="w-28 border-r border-border/40 p-2 bg-card/40">
              <div className="flex items-center gap-1.5 mb-3">
                <div className="h-3.5 w-3.5 rounded bg-gradient-to-br from-indigo-400 to-violet-500" />
                <span className="text-[10px] font-semibold">My App</span>
              </div>
              {[
                { l: 'Dashboard', sel: true },
                { l: 'Blogs', sel: false },
                { l: 'Contacts', sel: false },
                { l: 'Reports', sel: false },
                { l: 'Settings', sel: false },
              ].map((row) => (
                <div key={row.l} className={`rounded px-1.5 py-0.5 mb-0.5 text-[9px] ${row.sel ? 'bg-primary/10 text-primary' : 'text-muted-foreground'}`}>
                  {row.l}
                </div>
              ))}
            </div>
            <div className="flex-1 p-3">
              <div className="flex items-center justify-between mb-2">
                <p className="text-sm font-semibold">Dashboard</p>
                <div className="flex items-center gap-1 text-[9px] text-emerald-500">
                  <Wifi className="h-2.5 w-2.5" />
                  offline · synced 12s ago
                </div>
              </div>
              <div className="grid grid-cols-2 gap-1.5 mb-3">
                {[
                  { l: 'BLOGS', v: '32' },
                  { l: 'CONTACTS', v: '128' },
                ].map((s) => (
                  <div key={s.l} className="rounded border border-border/40 bg-card/60 p-1.5">
                    <p className="text-[8px] uppercase tracking-wider text-muted-foreground">{s.l}</p>
                    <p className="text-sm font-semibold">{s.v}</p>
                  </div>
                ))}
              </div>
              <div className="rounded border border-border/40 bg-card/40 p-2 h-16">
                <div className="flex items-end gap-[2px] h-full">
                  {[40, 55, 48, 62, 70, 65, 80, 72, 88, 84, 75, 68].map((h, i) => (
                    <div key={i} className="flex-1 rounded-sm bg-indigo-400/50" style={{ height: `${h}%` }} />
                  ))}
                </div>
              </div>
            </div>
          </div>
        </DesktopFrame>
      }
      features={[
        { icon: Monitor,         title: 'Wails v2',            body: 'Bind Go methods to React — call backend functions like normal JS.' },
        { icon: Database,        title: 'Local SQLite',        body: 'Ships with GORM + SQLite for true offline work.' },
        { icon: Shield,          title: 'Local auth',          body: 'bcrypt password hashing; no server, no JWT.' },
        { icon: Activity,        title: 'Offline-first sync',  body: 'Optional triple-mode adds outbox-based server sync.' },
        { icon: FileCode2,       title: 'PDF export',          body: 'Reports + receipts as PDF via go-pdf.' },
        { icon: FileSpreadsheet, title: 'Excel export',        body: 'Sheets via excelize; streaming-friendly.' },
        { icon: Cpu,             title: 'Native binary',       body: 'macOS / Windows / Linux — single binary per platform.' },
        { icon: DockerIcon,      title: 'No Docker needed',    body: 'Ship a real desktop binary; no container runtime.' },
        { icon: Wifi,             title: 'Frameless window',    body: 'Draggable panels, custom titlebar, dark/light theme.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes"
      prev={{ label: 'Mobile', href: '/docs/tech-kits/mobile' }}
    />
  )
}
