'use client'

import { useState } from 'react'
import { Database, KeyRound, Send, ShieldCheck } from 'lucide-react'

const TABS = [
  {
    key: 'data',
    icon: Database,
    label: 'Data',
    title: 'Models that generate everything downstream',
    body: 'Write the Go struct. Grit derives the migration, the query layer, the validation schema and the TypeScript type from it — so a field you add in one place appears everywhere it should.',
    points: ['Migrations from struct tags', 'Zod schema generated', 'Typed client, no drift'],
  },
  {
    key: 'auth',
    icon: KeyRound,
    label: 'Auth',
    title: 'Sessions that can actually be revoked',
    body: 'Server-side sessions backed by a row per device, rotation with replay detection, idle and absolute timeouts, and a working sign-out-everywhere. Not a JWT you cannot invalidate.',
    points: ['Per-device revocation', 'Rotation with replay detection', 'SSO via OIDC and SAML'],
  },
  {
    key: 'jobs',
    icon: Send,
    label: 'Jobs',
    title: 'Background work with somewhere to look',
    body: 'A queue, a scheduler and a dashboard that shows what ran, what failed and what is waiting. Retries are configured per job rather than hoped for.',
    points: ['Queue and cron built in', 'Per-job retry policy', 'Failures you can inspect'],
  },
  {
    key: 'security',
    icon: ShieldCheck,
    label: 'Security',
    title: 'The work that never wins a planning meeting',
    body: 'A strict content policy, CSRF, rate limiting, SSRF defence, ownership checks and a tamper-evident audit log — in the scaffold, before the deadline pressure starts.',
    points: ['OWASP patterns by default', 'Field-level encryption', 'GDPR export and erasure'],
  },
] as const

type TabKey = (typeof TABS)[number]['key']

export default function WithFeatureTabs({
  eyebrow = 'Platform',
  title = 'Four things you would otherwise build twice',
}: {
  eyebrow?: string
  title?: string
}) {
  const [active, setActive] = useState<TabKey>('data')
  const current = TABS.find((t) => t.key === active) ?? TABS[0]

  return (
    <section className="bg-gray-50 py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-base/7 font-semibold text-indigo-600 dark:text-indigo-400">
            {eyebrow}
          </p>
          <h2 className="mt-2 text-4xl font-semibold tracking-tight text-pretty text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
        </div>

        {/* Tabs. role=tablist with aria-selected so the control is announced as
            a tab set rather than four unrelated buttons. */}
        <div
          role="tablist"
          aria-label="Platform areas"
          className="mx-auto mt-14 flex max-w-2xl flex-wrap justify-center gap-2"
        >
          {TABS.map(({ key, icon: Icon, label }) => {
            const selected = key === active
            return (
              <button
                key={key}
                role="tab"
                aria-selected={selected}
                onClick={() => setActive(key)}
                className={`inline-flex items-center gap-2 rounded-full px-5 py-2.5 text-sm font-medium transition-colors ${
                  selected
                    ? 'bg-gray-900 text-white dark:bg-white dark:text-gray-900'
                    : 'border border-gray-200 bg-white text-gray-600 hover:text-gray-900 dark:border-white/10 dark:bg-white/5 dark:text-gray-400 dark:hover:text-white'
                }`}
              >
                <Icon aria-hidden="true" className="size-4" />
                {label}
              </button>
            )
          })}
        </div>

        {/* Panel */}
        <div className="mx-auto mt-12 max-w-5xl">
          <div className="grid grid-cols-1 items-center gap-10 rounded-3xl border border-gray-200 bg-white p-8 shadow-sm lg:grid-cols-2 lg:p-12 dark:border-white/10 dark:bg-gray-900">
            <div>
              <h3 className="text-2xl font-semibold tracking-tight text-pretty text-gray-900 dark:text-white">
                {current.title}
              </h3>
              <p className="mt-4 text-base/7 text-gray-600 dark:text-gray-400">
                {current.body}
              </p>
              <ul role="list" className="mt-6 space-y-2.5">
                {current.points.map((p) => (
                  <li
                    key={p}
                    className="flex items-center gap-2.5 text-sm text-gray-700 dark:text-gray-300"
                  >
                    <span className="size-1.5 rounded-full bg-indigo-600 dark:bg-indigo-400" />
                    {p}
                  </li>
                ))}
              </ul>
            </div>

            {/* A different illustration per tab, so switching visibly changes
                something other than the paragraph. */}
            <div className="rounded-2xl bg-gray-50 p-6 dark:bg-white/5">
              {active === 'data' && <DataArt />}
              {active === 'auth' && <AuthArt />}
              {active === 'jobs' && <JobsArt />}
              {active === 'security' && <SecurityArt />}
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

function Line({ w, tone = 'gray' }: { w: string; tone?: 'gray' | 'indigo' }) {
  return (
    <div
      className={`h-2 rounded ${
        tone === 'indigo' ? 'bg-indigo-500/70' : 'bg-gray-200 dark:bg-white/10'
      }`}
      style={{ width: w }}
    />
  )
}

function DataArt() {
  return (
    <div className="space-y-3 font-mono text-xs">
      <p className="text-gray-500 dark:text-gray-400">type Invoice struct</p>
      {[
        ['Number', 'string'],
        ['Total', 'float64'],
        ['Paid', 'bool'],
      ].map(([f, t]) => (
        <div key={f} className="flex items-center gap-3">
          <span className="w-16 text-gray-900 dark:text-white">{f}</span>
          <span className="text-indigo-600 dark:text-indigo-400">{t}</span>
        </div>
      ))}
      <div className="!mt-5 border-t border-gray-200 pt-4 dark:border-white/10">
        <p className="mb-2 text-gray-500 dark:text-gray-400">invoice.ts</p>
        <div className="space-y-2">
          <Line w="80%" tone="indigo" />
          <Line w="60%" />
          <Line w="70%" />
        </div>
      </div>
    </div>
  )
}

function AuthArt() {
  return (
    <div className="space-y-3">
      {[
        ['MacBook Pro', 'This device', true],
        ['iPhone 15', 'Nairobi · 2h ago', false],
        ['Chrome, Windows', 'Kampala · 3d ago', false],
      ].map(([device, meta, current]) => (
        <div
          key={device as string}
          className="flex items-center gap-3 rounded-lg border border-gray-200 bg-white p-3 dark:border-white/10 dark:bg-gray-900"
        >
          <span
            className={`size-2 rounded-full ${current ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'}`}
          />
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm text-gray-900 dark:text-white">
              {device as string}
            </p>
            <p className="truncate text-xs text-gray-500 dark:text-gray-400">
              {meta as string}
            </p>
          </div>
          {!current && (
            <span className="shrink-0 text-xs font-medium text-rose-600 dark:text-rose-400">
              Revoke
            </span>
          )}
        </div>
      ))}
    </div>
  )
}

function JobsArt() {
  return (
    <div className="space-y-2.5">
      {[
        ['send-welcome-email', 'succeeded', 'green'],
        ['generate-invoice-pdf', 'running', 'indigo'],
        ['nightly-backup', 'queued', 'gray'],
        ['sync-search-index', 'failed · 2 retries', 'rose'],
      ].map(([job, status, tone]) => (
        <div
          key={job as string}
          className="flex items-center justify-between gap-3 rounded-lg border border-gray-200 bg-white px-3 py-2.5 dark:border-white/10 dark:bg-gray-900"
        >
          <span className="truncate font-mono text-xs text-gray-900 dark:text-white">
            {job as string}
          </span>
          <span
            className={`shrink-0 rounded-full px-2 py-0.5 text-[11px] font-medium ${
              tone === 'green'
                ? 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-400'
                : tone === 'indigo'
                  ? 'bg-indigo-100 text-indigo-700 dark:bg-indigo-500/15 dark:text-indigo-400'
                  : tone === 'rose'
                    ? 'bg-rose-100 text-rose-700 dark:bg-rose-500/15 dark:text-rose-400'
                    : 'bg-gray-100 text-gray-600 dark:bg-white/10 dark:text-gray-400'
            }`}
          >
            {status as string}
          </span>
        </div>
      ))}
    </div>
  )
}

function SecurityArt() {
  return (
    <div className="space-y-2.5">
      {[
        'Content-Security-Policy',
        'Strict-Transport-Security',
        'X-Content-Type-Options',
        'Referrer-Policy',
        'Cross-Origin-Opener-Policy',
      ].map((header) => (
        <div key={header} className="flex items-center gap-2.5">
          <ShieldCheck
            aria-hidden="true"
            className="size-4 shrink-0 text-green-600 dark:text-green-400"
          />
          <span className="truncate font-mono text-xs text-gray-700 dark:text-gray-300">
            {header}
          </span>
        </div>
      ))}
      <p className="!mt-5 border-t border-gray-200 pt-4 text-xs text-gray-500 dark:border-white/10 dark:text-gray-400">
        Set on every response by default.
      </p>
    </div>
  )
}
