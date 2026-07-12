import { Smartphone, Shield, Database, Bot, Activity, Cpu, FileCode2, Container as DockerIcon, Bell } from 'lucide-react'
import { TechKitLayout, IncludedRow } from '@/components/tech-kit-layout'
import { MobileFrame } from '@/components/ui/device-frames'
import { getDocMetadata } from '@/config/docs-metadata'
import { buildKitPrompt } from '@/lib/kit-prompts'

export const metadata = getDocMetadata('/docs/tech-kits/mobile')

export default function MobileKitPage() {
  return (
    <TechKitLayout
      name="Mobile — Expo + API"
      tagline="React Native (Expo) frontend on top of a Grit API. Shared types both ways."
      pitch={
        <>
          <p>
            <code>grit new --mobile</code> scaffolds a Turborepo with{' '}
            <code>apps/expo</code> (Expo + React Native) and <code>apps/api</code> (Go). Same
            shared types in <code>packages/shared</code> as the web kits — Zod schemas validate
            on both ends. EAS Build, OTA updates, and a mobile-friendly auth flow are wired in.
          </p>
          <IncludedRow>Expo SDK 50+ with React Native</IncludedRow>
          <IncludedRow>Shared Zod schemas + generated TS types</IncludedRow>
          <IncludedRow>Mobile-friendly auth (refresh tokens, expo-secure-store (SecureStore))</IncludedRow>
          <IncludedRow>EAS Build configuration shipped; OTA-ready</IncludedRow>
        </>
      }
      command="grit new my-app --mobile"
      mockup={
        <div className="flex justify-center">
          <MobileFrame glow className="scale-90 origin-top">
            <div className="h-[300px] px-4">
              <div className="flex items-center justify-between mb-3">
                <div>
                  <p className="text-[9px] text-muted-foreground uppercase tracking-wider">Good morning</p>
                  <p className="text-sm font-semibold">My App</p>
                </div>
                <div className="h-7 w-7 rounded-full bg-primary/15 flex items-center justify-center">
                  <Bell className="h-3 w-3 text-primary" />
                </div>
              </div>
              <div className="space-y-2 mb-3">
                {[
                  { l: 'Tasks today', v: '4', t: 'text-emerald-500' },
                  { l: 'Unread', v: '7', t: 'text-sky-500' },
                ].map((s) => (
                  <div key={s.l} className="rounded-lg border border-border/50 bg-card/60 p-2.5">
                    <p className="text-[9px] uppercase tracking-wider text-muted-foreground">{s.l}</p>
                    <p className={`text-base font-semibold ${s.t}`}>{s.v}</p>
                  </div>
                ))}
              </div>
              <p className="text-[9px] uppercase tracking-wider text-muted-foreground mb-1.5">Recent</p>
              <div className="space-y-1.5">
                {['Maya signed in', 'New comment on Post #12', 'Invoice #1284 paid'].map((row, i) => (
                  <div key={i} className="flex items-center gap-2 text-[10px] text-foreground/80">
                    <span className={`h-1 w-1 rounded-full ${i === 0 ? 'bg-emerald-500' : i === 1 ? 'bg-sky-500' : 'bg-amber-500'}`} />
                    <span className="flex-1 truncate">{row}</span>
                  </div>
                ))}
              </div>
            </div>
          </MobileFrame>
        </div>
      }
      features={[
        { icon: Smartphone, title: 'Expo SDK 50+',       body: 'React Native with Expo Router; works on iOS, Android, web.' },
        { icon: Shield,     title: 'Auth + 2FA',         body: 'JWT refresh in expo-secure-store (SecureStore); biometric unlock optional.' },
        { icon: Bell,       title: 'Push notifications', body: 'Expo Push wired through the notification bell.' },
        { icon: Database,   title: 'GORM API',           body: 'Same Go API; offline-first sync available.' },
        { icon: Bot,        title: 'AI Gateway',         body: 'Stream completions to mobile, structured tool calls.' },
        { icon: Cpu,        title: 'Code generator',     body: 'grit generate resource → React Native screens + hooks.' },
        { icon: Activity,   title: 'Background jobs',    body: 'Asynq queue + cron in the Go API.' },
        { icon: FileCode2,  title: 'Audit log',          body: 'Tamper-evident SHA-256 chain on the API side.' },
        { icon: DockerIcon, title: 'EAS Build',          body: 'eas.json shipped; OTA updates through Expo.' },
      ]}
      architectureDeepLink="/docs/concepts/architecture-modes/mobile"
      starterPrompt={buildKitPrompt('mobile')}
      prev={{ label: 'API only', href: '/docs/tech-kits/api' }}
      next={{ label: 'Desktop', href: '/docs/tech-kits/desktop' }}
    />
  )
}
