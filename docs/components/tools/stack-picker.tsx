'use client'

import { useMemo, useState } from 'react'
import Link from 'next/link'
import { Check, Copy } from 'lucide-react'

/**
 * Answer five questions, get the grit new command for what you described.
 *
 * Every flag here is one `grit new` actually accepts: --arch, --frontend,
 * --db, --storage, --expo, --desktop, --i18n, --multitenant, --admin, --full.
 * Taken from cmd/grit/main.go rather than from memory, because a picker that
 * emits a flag the CLI rejects is worse than a comparison table.
 *
 * The recommendation is deliberately explained rather than just printed: the
 * architecture is the one decision that is awkward to change later, so a
 * reader should leave understanding why, not just holding a command.
 */

type Audience = 'customers' | 'internal' | 'both'
type Clients = 'web' | 'web-mobile' | 'web-desktop' | 'everything' | 'api'
type Data = 'postgres' | 'sqlite'
type Router = 'next' | 'vite'

interface Answers {
  audience: Audience
  clients: Clients
  data: Data
  router: Router
  tenants: boolean
}

const DEFAULTS: Answers = {
  audience: 'both',
  clients: 'web',
  data: 'postgres',
  router: 'next',
  tenants: false,
}

interface Choice<T> {
  value: T
  label: string
  blurb: string
}

const AUDIENCE: Choice<Audience>[] = [
  { value: 'both', label: 'Customers and a team', blurb: 'A product people sign up to, plus an operator console.' },
  { value: 'customers', label: 'Customers only', blurb: 'No separate back office yet.' },
  { value: 'internal', label: 'A team only', blurb: 'An internal tool. Nobody outside uses it.' },
]

const CLIENTS: Choice<Clients>[] = [
  { value: 'web', label: 'A website', blurb: 'The browser, and that is all for now.' },
  { value: 'web-mobile', label: 'Website and phone app', blurb: 'Adds an Expo app against the same API.' },
  { value: 'web-desktop', label: 'Website and desktop app', blurb: 'Adds a Wails app that can work offline.' },
  { value: 'everything', label: 'Web, phone and desktop', blurb: 'Every target, plus a docs site.' },
  { value: 'api', label: 'No interface at all', blurb: 'An API other things call.' },
]

const DATA: Choice<Data>[] = [
  { value: 'postgres', label: 'Postgres', blurb: 'The production answer. Needs Docker, or a hosted database.' },
  { value: 'sqlite', label: 'SQLite', blurb: 'A file. No Docker, no database server, nothing to install.' },
]

const ROUTER: Choice<Router>[] = [
  { value: 'next', label: 'Next.js', blurb: 'Server components, and the frontend most people know.' },
  { value: 'vite', label: 'Vite and TanStack Router', blurb: 'A single-page app, typed routes, very fast reloads.' },
]

/** The architecture that fits the answers, and why. */
function decide(a: Answers): { arch: string; reason: string } {
  if (a.clients === 'api') {
    return {
      arch: 'api',
      reason:
        'Nothing to render, so there is no frontend to scaffold. You get the Go API, its tests and its Docker setup.',
    }
  }
  if (a.audience === 'internal') {
    return {
      arch: 'api --admin',
      reason:
        'An internal tool is the admin. This puts the API and the admin in one binary with the admin served at the root, and leaves out the customer-facing site entirely.',
    }
  }
  if (a.audience === 'customers') {
    return {
      arch: 'single',
      reason:
        'One app, which is also the project root, with the Go API beside it. Least to run and least to deploy. The admin is a section of the same router if you turn out to need one.',
    }
  }
  return {
    arch: 'triple',
    reason:
      'Two audiences, two apps: a customer-facing site and a separate admin, over one API. They want different navigation and different density, and keeping them apart means a role check is not the only thing standing between a customer and the dangerous buttons.',
  }
}

function buildCommand(a: Answers, name: string): string {
  const { arch } = decide(a)
  const flags: string[] = []

  if (a.clients === 'everything') {
    // --full is every target at once, so the architecture and the per-client
    // flags are already implied by it.
    flags.push('--full')
  } else {
    flags.push(arch === 'api --admin' ? '--api --admin' : `--arch ${arch}`)
    if (a.clients === 'web-mobile') flags.push('--expo')
    if (a.clients === 'web-desktop') flags.push('--desktop')
  }

  if (a.clients !== 'api' && a.router === 'vite') flags.push('--vite')
  if (a.data === 'sqlite') {
    // Paired on purpose: SQLite removes the database container, and local
    // storage removes the MinIO one. Together they are a project that runs
    // with nothing but Go and Node.
    flags.push('--db sqlite', '--storage local')
  }
  if (a.tenants) flags.push('--multitenant')

  return `grit new ${name.trim() || 'myapp'} ${flags.join(' ')}`
}

function Question<T extends string>({
  legend,
  help,
  choices,
  value,
  onChange,
  name,
}: {
  legend: string
  help?: string
  choices: Choice<T>[]
  value: T
  onChange: (v: T) => void
  name: string
}) {
  return (
    <fieldset className="mb-9">
      <legend className="mb-1 font-display text-base font-semibold text-foreground">
        {legend}
      </legend>
      {help && <p className="mb-3 text-sm text-muted-foreground">{help}</p>}
      <div className="grid gap-2 sm:grid-cols-2">
        {choices.map((c) => {
          const id = `${name}-${c.value}`
          const active = value === c.value
          return (
            <label
              key={c.value}
              htmlFor={id}
              className={`cursor-pointer rounded-xl border p-3.5 transition-colors ${
                active
                  ? 'border-accent/60 bg-accent/[0.07]'
                  : 'border-border/60 bg-card/20 hover:border-accent/40'
              }`}
            >
              <span className="flex items-start gap-2.5">
                <input
                  id={id}
                  type="radio"
                  name={name}
                  value={c.value}
                  checked={active}
                  onChange={() => onChange(c.value)}
                  className="mt-1 h-4 w-4 shrink-0 accent-[var(--accent)]"
                />
                <span>
                  <span className="block text-sm font-medium text-foreground">{c.label}</span>
                  <span className="mt-0.5 block text-xs leading-relaxed text-muted-foreground">
                    {c.blurb}
                  </span>
                </span>
              </span>
            </label>
          )
        })}
      </div>
    </fieldset>
  )
}

export function StackPicker() {
  const [a, setA] = useState<Answers>(DEFAULTS)
  const [name, setName] = useState('myapp')
  const [copied, setCopied] = useState(false)

  const decision = useMemo(() => decide(a), [a])
  const command = useMemo(() => buildCommand(a, name), [a, name])

  async function copy() {
    try {
      await navigator.clipboard.writeText(command)
      setCopied(true)
      setTimeout(() => setCopied(false), 1800)
    } catch {
      setCopied(false)
    }
  }

  return (
    <div className="grid gap-10 lg:grid-cols-[minmax(0,1fr)_380px]">
      <div>
        <div className="mb-9">
          <label htmlFor="project" className="mb-1.5 block font-display text-base font-semibold text-foreground">
            What is it called?
          </label>
          <input
            id="project"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="myapp"
            className="w-full max-w-xs rounded-lg border border-border bg-background px-3 py-2 font-mono text-sm outline-none focus-visible:border-accent focus-visible:ring-2 focus-visible:ring-accent/30"
          />
        </div>

        <Question
          name="audience"
          legend="Who uses it?"
          help="This is the question that decides the architecture, and it is the hardest one to change later."
          choices={AUDIENCE}
          value={a.audience}
          onChange={(v) => setA({ ...a, audience: v })}
        />

        <Question
          name="clients"
          legend="Where do they use it?"
          choices={CLIENTS}
          value={a.clients}
          onChange={(v) => setA({ ...a, clients: v })}
        />

        <Question
          name="data"
          legend="Where does the data go?"
          help="SQLite is the one to pick while learning: there is no container to start and nothing to install."
          choices={DATA}
          value={a.data}
          onChange={(v) => setA({ ...a, data: v })}
        />

        {a.clients !== 'api' && (
          <Question
            name="router"
            legend="Which frontend?"
            choices={ROUTER}
            value={a.router}
            onChange={(v) => setA({ ...a, router: v })}
          />
        )}

        <fieldset className="mb-4">
          <legend className="mb-1 font-display text-base font-semibold text-foreground">
            Will one deployment serve several customers?
          </legend>
          <p className="mb-3 text-sm text-muted-foreground">
            If each customer is an organisation with its own users and its own data, say yes now.
            Adding tenancy later means touching every query, every permission check and every
            billing record.
          </p>
          <label
            htmlFor="tenants"
            className="inline-flex cursor-pointer items-center gap-2.5 rounded-xl border border-border/60 bg-card/20 px-3.5 py-3 text-sm"
          >
            <input
              id="tenants"
              type="checkbox"
              checked={a.tenants}
              onChange={(e) => setA({ ...a, tenants: e.target.checked })}
              className="h-4 w-4 accent-[var(--accent)]"
            />
            Yes, it is multi-tenant
          </label>
        </fieldset>
      </div>

      {/* The answer */}
      <aside className="lg:sticky lg:top-24 lg:self-start">
        <div className="rounded-2xl border border-accent/30 bg-card/60 p-6">
          <span className="font-mono text-xs uppercase tracking-wider text-accent">
            Your stack
          </span>
          <h2 className="mb-2 mt-2 font-display text-2xl font-semibold tracking-tight">
            {decision.arch === 'api --admin' ? 'API with admin' : decision.arch}
          </h2>
          <p className="mb-5 text-sm leading-relaxed text-muted-foreground">{decision.reason}</p>

          <div className="mb-3 flex items-center justify-between gap-3">
            <span className="text-xs font-medium text-muted-foreground">Run this</span>
            <button
              type="button"
              onClick={copy}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border/60 px-2.5 py-1 text-xs transition-colors hover:border-accent/60"
            >
              {copied ? (
                <>
                  <Check className="h-3.5 w-3.5 text-emerald-400" aria-hidden="true" />
                  Copied
                </>
              ) : (
                <>
                  <Copy className="h-3.5 w-3.5" aria-hidden="true" />
                  Copy
                </>
              )}
            </button>
          </div>
          <pre className="whitespace-pre-wrap break-words rounded-lg bg-background/80 p-3.5 font-mono text-[12.5px] leading-relaxed text-foreground">
            {command}
          </pre>

          {a.data === 'sqlite' && (
            <p className="mt-4 text-xs leading-relaxed text-muted-foreground">
              With SQLite and local storage there is no Docker to start. Run{' '}
              <code className="font-mono">grit start</code> and the project is up.
            </p>
          )}

          <p className="mt-5 border-t border-border/50 pt-4 text-xs leading-relaxed text-muted-foreground">
            Running <code className="font-mono">grit new</code> with no flags asks these same
            questions in the terminal, and pressing Enter on every prompt gives you the defaults.
          </p>
        </div>

        <p className="mt-4 px-1 text-xs leading-relaxed text-muted-foreground">
          Not sure about the architecture?{' '}
          <Link href="/docs/concepts/architecture" className="text-accent hover:underline">
            The architecture overview
          </Link>{' '}
          compares the five side by side, and{' '}
          <Link href="/docs/stack-selector" className="text-accent hover:underline">
            the stack selector
          </Link>{' '}
          walks through the trade-offs in prose.
        </p>
      </aside>
    </div>
  )
}
