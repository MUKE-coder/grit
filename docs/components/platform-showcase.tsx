'use client'

import { useState } from 'react'
import type React from 'react'
import { Globe, Monitor, Smartphone, Server } from 'lucide-react'

/**
 * One framework, every platform — shown with screenshots, not drawings.
 *
 * Every image here is a capture of a real generated project:
 *   admin.png       the Next.js admin panel (grit new --triple), Atlas theme
 *   desktop.png     the Wails desktop app, including its own window chrome
 *   mobile.png      the Expo app running on an Android emulator
 *   api-scalar.png  the Scalar reference the Go API serves at /docs
 *
 * That distinction is the whole argument of the section. Anyone can draw a
 * dashboard; only a framework that actually ships these can screenshot them.
 * If you retake one, use a real project — never a mockup, never a Figma export.
 *
 * The commands in `steps` are executable as written and were run before this
 * shipped. Treat them as tested code: if a flag changes, re-run them here
 * rather than editing the string to look right.
 */

type Shell = 'browser' | 'window' | 'phone'

interface Platform {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  headline: string
  blurb: string
  image: string
  shell: Shell
  /** Chrome caption: a URL for browser shells, an app name for windows. */
  chromeLabel: string
  steps: { cmd: string; note: string }[]
  footnote?: string
}

const PLATFORMS: Platform[] = [
  {
    key: 'web',
    label: 'Web',
    icon: Globe,
    headline: 'A Next.js app and a real admin panel',
    blurb:
      'Two front-ends against one Go API: a public web app and a resource-driven admin panel with tables, filters, multi-step forms and RBAC already wired up.',
    image: '/images/platforms/admin.png',
    shell: 'browser',
    chromeLabel: 'localhost:3001/dashboard',
    steps: [
      { cmd: 'grit new myapp --triple', note: 'Go API + web app + admin panel' },
      { cmd: 'cd myapp && docker compose up -d', note: 'Postgres, Redis, MinIO, Mailhog' },
      { cmd: 'pnpm install && grit migrate && grit seed', note: 'deps, tables, a demo admin login' },
      { cmd: 'grit start', note: 'all three, one terminal' },
    ],
    footnote: 'Admin on :3001, web on :3000, API on :8080.',
  },
  {
    key: 'desktop',
    label: 'Desktop',
    icon: Monitor,
    headline: 'A native window, offline-capable',
    blurb:
      'A Wails binary with its own title bar and a local SQLite mirror. Reads keep working with no network, writes queue, and the sync indicator tells the user which state they are in.',
    image: '/images/platforms/desktop.png',
    shell: 'window',
    chromeLabel: 'myapp — desktop',
    steps: [
      { cmd: 'grit new myapp --triple --desktop', note: 'adds apps/desktop (Wails + TanStack Router)' },
      { cmd: 'cd myapp && docker compose up -d', note: 'same infrastructure as the web stack' },
      { cmd: 'pnpm install && grit migrate && grit seed', note: 'shared monorepo install' },
      { cmd: 'grit start', note: 'the desktop window opens with the rest' },
    ],
    footnote: 'Wails needs a C toolchain and a WebView runtime. grit package builds the installer.',
  },
  {
    key: 'mobile',
    label: 'Mobile',
    icon: Smartphone,
    headline: 'An Expo app on the same types',
    blurb:
      'React Native through Expo Router, with NativeWind styling and the same generated client the web app uses. Generate a resource and it is on the phone too.',
    image: '/images/platforms/mobile.png',
    shell: 'phone',
    chromeLabel: 'Expo · Android',
    steps: [
      { cmd: 'grit new myapp --triple --expo', note: 'adds apps/expo (Expo Router + NativeWind)' },
      { cmd: 'cd myapp && docker compose up -d', note: 'the API the phone will talk to' },
      { cmd: 'pnpm install && grit migrate && grit seed', note: 'shared monorepo install' },
      { cmd: 'grit start expo', note: 'scan the QR code, or press a for Android' },
    ],
    footnote: 'On a physical device, point EXPO_PUBLIC_API_URL at your machine’s LAN IP.',
  },
  {
    key: 'api',
    label: 'API',
    icon: Server,
    headline: 'Documented the moment it exists',
    blurb:
      'The Go API serves a Scalar reference at /docs, built from your routes and models. Add a resource and its endpoints show up — no annotations to write, no spec to maintain by hand.',
    image: '/images/platforms/api-scalar.png',
    shell: 'browser',
    chromeLabel: 'localhost:8080/docs',
    steps: [
      { cmd: 'grit new myapp --api', note: 'the Go API on its own' },
      { cmd: 'cd myapp && docker compose up -d', note: 'Postgres and Redis' },
      { cmd: 'grit migrate && grit seed', note: 'tables and sample rows' },
      { cmd: 'grit start', note: 'docs at /docs, GORM Studio at /studio' },
    ],
    footnote: 'The OpenAPI document is served at /docs/openapi.json for client codegen.',
  },
]

function Chrome({ shell, label, children }: { shell: Shell; label: string; children: React.ReactNode }) {
  if (shell === 'phone') {
    // A bezel rather than a full device illustration: enough to read as a
    // phone, not so much that it competes with the screenshot.
    return (
      <div className="flex justify-center">
        <div className="relative rounded-[2.25rem] border-[10px] border-[#1c1c22] bg-[#1c1c22] shadow-[0_28px_70px_-20px_rgba(2,6,23,0.75)] max-w-[17rem]">
          <div className="absolute left-1/2 top-2 h-1 w-16 -translate-x-1/2 rounded-full bg-white/20" />
          <div className="overflow-hidden rounded-[1.6rem]">{children}</div>
        </div>
      </div>
    )
  }

  return (
    <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
      <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
        <div className="flex items-center gap-1.5">
          <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </div>
        <span className="mx-auto text-[11px] font-mono text-muted-foreground truncate max-w-[60%]">
          {label}
        </span>
      </div>
      {children}
    </div>
  )
}

export function PlatformShowcase() {
  const [active, setActive] = useState(PLATFORMS[0].key)
  const platform = PLATFORMS.find((p) => p.key === active) ?? PLATFORMS[0]

  return (
    <div>
      <div
        role="tablist"
        aria-label="Platform"
        className="flex flex-wrap items-center gap-2 mb-8"
      >
        {PLATFORMS.map((p) => {
          const selected = p.key === active
          return (
            <button
              key={p.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(p.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <p.icon className="h-4 w-4" />
              {p.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_20rem] gap-8 lg:gap-10 items-start">
        <div>
          <Chrome shell={platform.shell} label={platform.chromeLabel}>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              key={platform.key}
              src={platform.image}
              alt={`${platform.label}: ${platform.headline}`}
              className="w-full h-auto block"
              loading="lazy"
            />
          </Chrome>
          <p className="text-xs text-muted-foreground/70 mt-4">
            A screenshot of a generated project, not a mockup.
          </p>
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-2 leading-snug">{platform.headline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-6">{platform.blurb}</p>

          <ol className="space-y-3">
            {platform.steps.map((step, i) => (
              <li key={step.cmd} className="flex gap-3">
                <span className="mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-primary/10 text-[11px] font-mono font-medium text-primary">
                  {i + 1}
                </span>
                <div className="min-w-0">
                  <code className="block text-[12px] font-mono text-foreground/90 break-words">
                    {step.cmd}
                  </code>
                  <span className="text-[11.5px] text-muted-foreground/80">{step.note}</span>
                </div>
              </li>
            ))}
          </ol>

          {platform.footnote && (
            <p className="text-[11.5px] text-muted-foreground/70 mt-5 pt-4 border-t border-border/40 leading-relaxed">
              {platform.footnote}
            </p>
          )}
        </div>
      </div>
    </div>
  )
}
