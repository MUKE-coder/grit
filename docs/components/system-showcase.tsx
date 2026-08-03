'use client'

import { useState } from 'react'
import type React from 'react'
import { KeyRound, ShieldAlert, Gauge, ScrollText, DatabaseBackup, FolderOpen, Layers3, CalendarClock } from 'lucide-react'

/**
 * The System Hub, shown page by page.
 *
 * Every image is a screenshot of a running generated admin — the same project,
 * the same session, captured at 1440px. Several needed real infrastructure to
 * produce: the jobs page shows image:process jobs that five actual uploads
 * enqueued through Redis, and the file page shows thumbnails MinIO served back.
 * That is the point. A drawn "background jobs" panel proves nothing; a queue
 * with five completed jobs in it proves the queue exists and works.
 *
 * If you retake these, run a real project with docker compose up, do the thing
 * the page reports on, then screenshot. Never mock the numbers.
 */

interface Page {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  route: string
  headline: string
  body: string
  /** The specific thing to notice in the screenshot. */
  detail: string
  image: string
}

const PAGES: Page[] = [
  {
    key: 'roles',
    label: 'Roles & permissions',
    icon: KeyRound,
    route: '/system/roles',
    headline: 'Permissions you can actually reason about',
    body:
      'Every resource you generate registers its own create / view / edit / delete permissions automatically. Roles are rows in your database, not constants in a file, so the people running the app can change them without a deploy.',
    detail:
      'Grant a whole resource and it keeps any action added to it later — so generating a new resource never silently widens or narrows an existing role.',
    image: '/images/system/roles.png',
  },
  {
    key: 'security',
    label: 'Security',
    icon: ShieldAlert,
    route: '/system/security',
    headline: 'Brute force gets expensive, automatically',
    body:
      'Sentinel watches auth endpoints and bans IPs that trip the rate limit. Re-offenders escalate — 5 hours, then 8, then 24, then a week — so a bot cannot simply wait out a fixed cooldown.',
    detail:
      'Active bans, IPs currently hitting limits, and recent threats are all on one page, with the escalation policy stated rather than buried in config.',
    image: '/images/system/security.png',
  },
  {
    key: 'performance',
    label: 'Performance',
    icon: Gauge,
    route: '/system/performance',
    headline: 'The four golden signals, already wired',
    body:
      'Latency, traffic, errors and saturation — the SRE signals — measured by Pulse from the moment your API starts. No agent to install, no dashboard to build.',
    detail:
      'p50/p95/p99 rather than an average, because the average hides the tail. The slowest-routes table names the actual endpoint and its error rate.',
    image: '/images/system/performance.png',
  },
  {
    key: 'activity',
    label: 'User activity',
    icon: ScrollText,
    route: '/system/activity',
    headline: 'Who did what, with an IP attached',
    body:
      'Sign-ins, writes and operator actions land in one timeline, grouped by day and tagged by severity. Filter to flagged, critical, or after-hours events, then export the range.',
    detail:
      'Each row carries the event name (auth.login), the host and the IP — the three things you actually want when someone asks what happened on Tuesday.',
    image: '/images/system/activity.png',
  },
  {
    key: 'backups',
    label: 'Data & backup',
    icon: DatabaseBackup,
    route: '/system/backups',
    headline: 'Backups that are not your problem yet',
    body:
      'A full backup runs on a schedule you pick — daily, weekly, monthly, yearly — or on demand. The four most recent are kept and downloadable from the panel.',
    detail:
      'Each archive is a ZIP: one CSV per table, a dump.sql of INSERTs, and a metadata.json manifest. Restore with grit restore backup.',
    image: '/images/system/backups.png',
  },
  {
    key: 'files',
    label: 'File storage',
    icon: FolderOpen,
    route: '/system/files',
    headline: 'S3 storage with the plumbing done',
    body:
      'Uploads go browser-to-storage through a presigned URL, so files never travel through your API. MinIO locally, S3 or R2 or B2 in production — same code, one env var.',
    detail:
      'Totals, a per-type breakdown and thumbnails served back from storage. These five files were uploaded through the panel while this screenshot was taken.',
    image: '/images/system/files.png',
  },
  {
    key: 'jobs',
    label: 'Background jobs',
    icon: Layers3,
    route: '/system/jobs',
    headline: 'A queue you can see into',
    body:
      'Redis-backed jobs via asynq, with active / pending / completed / failed / retry all visible, and a dead-letter queue for the ones that never made it.',
    detail:
      'The five image:process jobs here were enqueued by uploading five images — the queue is wired into the framework, not something you bolt on.',
    image: '/images/system/jobs.png',
  },
  {
    key: 'cron',
    label: 'Cron schedules',
    icon: CalendarClock,
    route: '/system/cron',
    headline: 'Recurring work, declared and visible',
    body:
      'The scheduler ships with the jobs a real app needs — expiring tokens, orphaned uploads, database backups — and shows every task with its cron expression and next run.',
    detail:
      'Adding your own is a registration call, and it appears here with the rest instead of living only in a crontab someone forgot about.',
    image: '/images/system/cron.png',
  },
]

export function SystemShowcase() {
  const [active, setActive] = useState(PAGES[0].key)
  const page = PAGES.find((p) => p.key === active) ?? PAGES[0]

  return (
    <div>
      <div
        role="tablist"
        aria-label="System page"
        className="flex flex-wrap gap-2 mb-8"
      >
        {PAGES.map((p) => {
          const selected = p.key === active
          return (
            <button
              key={p.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(p.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-3.5 py-2 text-[13px] font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <p.icon className="h-3.5 w-3.5" />
              {p.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_19rem] gap-8 lg:gap-10 items-start">
        <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
          <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
            <div className="flex items-center gap-1.5">
              <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
            </div>
            <span className="mx-auto text-[11px] font-mono text-muted-foreground">
              localhost:3001{page.route}
            </span>
          </div>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            key={page.key}
            src={page.image}
            alt={`${page.label}: ${page.headline}`}
            className="w-full h-auto block"
            loading="lazy"
          />
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-3 leading-snug">{page.headline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5">{page.body}</p>
          <div className="rounded-xl border border-border/60 bg-card/40 p-4">
            <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-2">
              What to notice
            </div>
            <p className="text-[12.5px] text-foreground/80 leading-relaxed">{page.detail}</p>
          </div>
          <p className="text-[11.5px] text-muted-foreground/70 mt-4 leading-relaxed">
            Every one of these ships with{' '}
            <code className="text-foreground/70">grit new</code>. There is nothing to install
            and nothing to wire up.
          </p>
        </div>
      </div>
    </div>
  )
}
