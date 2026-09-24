'use client'

import * as React from 'react'
import Image from 'next/image'
import { Pause, Play, Terminal } from 'lucide-react'
import { cn } from '@/lib/utils'

/**
 * Grit in action: the commands on the left, what they produced on the right.
 *
 * The rest of this page argues. This section demonstrates. Somebody who has
 * never heard of Grit should be able to watch it once and know what the tool
 * does, without reading a word of prose, which is why the screenshot is tied to
 * the command rather than sitting in a carousel of its own: the pairing is the
 * whole point.
 *
 * The screenshots are real admin panels from generated projects, not mockups.
 *
 * Accessibility: the run advances on a timer, so there is a pause control and
 * it stops on hover or keyboard focus; the step buttons are real buttons, in
 * tab order, and each one jumps straight to its step. Anyone who has asked for
 * reduced motion gets the finished state of every step with no typing and no
 * auto-advance.
 */

interface Step {
  /** Typed into the terminal, character by character. */
  command: string
  /** Printed under it once the command finishes typing. */
  output: string[]
  shot: string
  alt: string
  /** Shown under the screenshot. */
  caption: string
  label: string
}

const STEPS: Step[] = [
  {
    command: 'grit new myapp --triple --next',
    output: [
      'Scaffolding Go API...',
      'Adding batteries (cache, storage, mail, jobs, cron, AI)...',
      'Scaffolding Next.js web app...',
      'Scaffolding admin panel...',
      'Project created successfully!',
    ],
    shot: '/images/auth/login.png',
    alt: 'The generated sign-in page, with email and password fields and social login buttons',
    caption: 'Sign in, register, reset, verify. Generated, themed, working.',
    label: 'Scaffold',
  },
  {
    command: 'cd myapp && grit start',
    output: [
      'API on :8080   web on :3000   admin on :3001',
      'GORM Studio on /studio, Scalar docs on /docs',
      'Ready.',
    ],
    shot: '/images/platforms/admin.png',
    alt: 'The generated admin dashboard, with statistic cards, a chart and a recent activity list',
    caption: 'The admin panel, already populated. Nothing assembled by hand.',
    label: 'Run',
  },
  {
    command:
      'grit generate resource Product --fields "name:string,price:money,category:belongs_to:Category"',
    output: [
      'apps/api/internal/models, services, handlers',
      'packages/shared/schemas + types',
      'apps/web/hooks/use-products.ts',
      'apps/admin/resources + pages',
      'Injecting into AutoMigrate, routes, /docs, permissions...',
      '12 files written, 13 injections applied',
    ],
    shot: '/images/forms/relationship.png',
    alt: 'A generated admin form showing a searchable relationship picker for selecting a related record',
    caption: 'One command. The form, the picker and the API behind it.',
    label: 'Generate',
  },
  {
    command: 'grit generate resource Invoice --fields "..." --seed --faker --count 500',
    output: [
      'Resource generated.',
      'apps/api/internal/database/invoices_seeder.go',
      'seeded 500 rows in 312ms (1600 rows/s)',
    ],
    shot: '/images/forms/line-items.png',
    alt: 'A generated invoice form with repeatable line items and a running total',
    caption: 'Line items, running totals, money that is not a float.',
    label: 'Seed',
  },
  {
    command: 'grit add role EDITOR',
    output: [
      'Role added and wired into every guard.',
      'Permissions screen updated.',
    ],
    shot: '/images/system/roles.png',
    alt: 'The roles and permissions screen, listing roles with a matrix of permission checkboxes',
    caption: 'Roles and per-resource permissions, from one command.',
    label: 'Secure',
  },
  {
    command: 'grit deploy --host you@server.com --domain myapp.com',
    output: [
      'Building... uploading... systemd + Caddy...',
      'Live at https://myapp.com with TLS',
    ],
    shot: '/images/platforms/api-scalar.png',
    alt: 'The generated API reference, rendered with Scalar, listing endpoints and schemas',
    caption: 'Shipped, with the API reference generated from the routes.',
    label: 'Deploy',
  },
]

const TYPING_MS = 26
const OUTPUT_MS = 260
const HOLD_MS = 2200

export function GritInAction() {
  const [step, setStep] = React.useState(0)
  const [typed, setTyped] = React.useState('')
  const [shownOutput, setShownOutput] = React.useState(0)
  const [paused, setPaused] = React.useState(false)
  const [reduced, setReduced] = React.useState(false)

  React.useEffect(() => {
    const query = window.matchMedia('(prefers-reduced-motion: reduce)')
    const apply = () => setReduced(query.matches)
    apply()
    query.addEventListener('change', apply)
    return () => query.removeEventListener('change', apply)
  }, [])

  const current = STEPS[step]

  // Typing, then output lines, then hold, then the next step. One effect owns
  // the whole cycle so the screenshot can never get out of step with the
  // command that produced it.
  React.useEffect(() => {
    if (reduced) {
      setTyped(current.command)
      setShownOutput(current.output.length)
      return
    }
    if (paused) return

    if (typed.length < current.command.length) {
      const t = setTimeout(() => setTyped(current.command.slice(0, typed.length + 1)), TYPING_MS)
      return () => clearTimeout(t)
    }
    if (shownOutput < current.output.length) {
      const t = setTimeout(() => setShownOutput((n) => n + 1), OUTPUT_MS)
      return () => clearTimeout(t)
    }
    const t = setTimeout(() => goTo((step + 1) % STEPS.length), HOLD_MS)
    return () => clearTimeout(t)
  }, [typed, shownOutput, step, paused, reduced, current])

  function goTo(next: number) {
    setStep(next)
    setTyped('')
    setShownOutput(0)
  }

  return (
    <section className="border-t border-border/40 py-20 sm:py-28">
      <div className="container max-w-screen-xl px-6">
        <div className="mx-auto mb-12 max-w-2xl text-center">
          <span className="tag-mono mb-3 block text-primary/80">Grit in action</span>
          <h2 className="mb-4 text-3xl font-bold leading-tight tracking-tight sm:text-4xl">
            Six commands.
            <br className="hidden sm:block" /> A working application.
          </h2>
          <p className="text-base leading-relaxed text-muted-foreground sm:text-lg">
            Left: the commands, in the order you actually run them. Right: what each one
            produced. These are screenshots of real generated projects.
          </p>
        </div>

        <div
          className="grid items-start gap-6 lg:grid-cols-2 lg:gap-8"
          onMouseEnter={() => setPaused(true)}
          onMouseLeave={() => setPaused(false)}
          onFocusCapture={() => setPaused(true)}
          onBlurCapture={() => setPaused(false)}
        >
          {/* ── Left: the terminal ─────────────────────────────── */}
          <div className="overflow-hidden rounded-xl border border-border bg-[#0b0b12] shadow-2xl">
            <div className="flex items-center gap-2 border-b border-white/5 px-4 py-3">
              <span className="h-3 w-3 rounded-full bg-[#ff5f57]" />
              <span className="h-3 w-3 rounded-full bg-[#febc2e]" />
              <span className="h-3 w-3 rounded-full bg-[#28c840]" />
              <span className="ml-2 inline-flex items-center gap-1.5 font-mono text-xs text-white/40">
                <Terminal className="h-3.5 w-3.5" aria-hidden="true" />
                myapp
              </span>
              <button
                type="button"
                onClick={() => setPaused((p) => !p)}
                className="ml-auto inline-flex items-center gap-1.5 rounded px-2 py-1 font-mono text-[11px] text-white/40 transition-colors hover:bg-white/5 hover:text-white/70 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white/40"
              >
                {paused ? (
                  <>
                    <Play className="h-3 w-3" aria-hidden="true" /> Play
                  </>
                ) : (
                  <>
                    <Pause className="h-3 w-3" aria-hidden="true" /> Pause
                  </>
                )}
              </button>
            </div>

            <div className="min-h-[300px] p-5 font-mono text-[13px] leading-relaxed sm:min-h-[340px]">
              <div className="flex gap-2">
                <span className="shrink-0 text-[#6c5ce7]">$</span>
                <span className="break-all text-white/90">
                  {typed}
                  {!reduced && typed.length < current.command.length && (
                    <span className="ml-0.5 inline-block h-4 w-2 translate-y-0.5 animate-pulse bg-white/70" />
                  )}
                </span>
              </div>
              <div className="mt-3 space-y-1.5" aria-live="polite">
                {current.output.slice(0, shownOutput).map((line, i) => (
                  <div
                    key={line}
                    className={cn(
                      'pl-4',
                      i === current.output.length - 1 ? 'text-[#00b894]' : 'text-white/45',
                    )}
                  >
                    {i === current.output.length - 1 ? '✓ ' : '  '}
                    {line}
                  </div>
                ))}
              </div>
            </div>

            {/* Step buttons: also the way to skip straight to a step. */}
            <div className="flex flex-wrap gap-1.5 border-t border-white/5 p-3">
              {STEPS.map((s, i) => (
                <button
                  key={s.label}
                  type="button"
                  onClick={() => goTo(i)}
                  aria-current={i === step ? 'step' : undefined}
                  className={cn(
                    'rounded-md px-2.5 py-1.5 font-mono text-[11px] transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-white/40',
                    i === step
                      ? 'bg-[#6c5ce7] text-white'
                      : 'text-white/40 hover:bg-white/5 hover:text-white/70',
                  )}
                >
                  {i + 1}. {s.label}
                </button>
              ))}
            </div>
          </div>

          {/* ── Right: what the command produced ───────────────── */}
          <div>
            <div className="relative overflow-hidden rounded-xl border border-border bg-muted/20 shadow-2xl">
              {/* Every shot is mounted and cross-faded, so switching steps does
                  not show a blank frame while the next image loads. */}
              <div className="relative aspect-[16/10]">
                {STEPS.map((s, i) => (
                  <Image
                    key={s.shot}
                    src={s.shot}
                    alt={i === step ? s.alt : ''}
                    fill
                    sizes="(max-width: 1024px) 100vw, 640px"
                    priority={i === 0}
                    aria-hidden={i !== step}
                    className={cn(
                      'object-cover object-top transition-opacity duration-500',
                      i === step ? 'opacity-100' : 'opacity-0',
                    )}
                  />
                ))}
              </div>
            </div>
            <p className="mt-4 text-center text-sm leading-relaxed text-muted-foreground">
              {current.caption}
            </p>
          </div>
        </div>
      </div>
    </section>
  )
}
