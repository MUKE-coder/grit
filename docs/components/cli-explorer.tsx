'use client'

import { useEffect, useMemo, useRef, useState } from 'react'
import Link from 'next/link'
import {
  Search,
  Terminal,
  FilePlus2,
  FilePen,
  FileX2,
  Play,
  RotateCcw,
  Check,
  Copy,
  ChevronRight,
  X,
} from 'lucide-react'
import {
  CLI_COMMANDS,
  CLI_CATEGORIES,
  type CliCommand,
  type CommandCategory,
  type FileStatus,
} from '@/config/cli-commands'

/**
 * The interactive CLI reference.
 *
 * Two things it is built around:
 *
 * The terminal types the command and then streams the output, because the
 * question people actually have about a CLI is "what happens when I run this",
 * and a static code block answers a different, smaller question.
 *
 * The file list is the other half, and it is the half no CLI reference usually
 * has. Every entry was captured from a real run, so "what does this touch" has
 * an answer you can trust before you run it on a project you care about.
 */

const CATEGORY_TONE: Record<CommandCategory, string> = {
  Scaffold: 'text-violet-400 border-violet-500/30 bg-violet-500/10',
  Generate: 'text-cyan-400 border-cyan-500/30 bg-cyan-500/10',
  Add: 'text-emerald-400 border-emerald-500/30 bg-emerald-500/10',
  Run: 'text-blue-400 border-blue-500/30 bg-blue-500/10',
  Data: 'text-amber-400 border-amber-500/30 bg-amber-500/10',
  Ship: 'text-rose-400 border-rose-500/30 bg-rose-500/10',
  Meta: 'text-slate-400 border-slate-500/30 bg-slate-500/10',
}

const STATUS_META: Record<
  FileStatus,
  { label: string; tone: string; Icon: typeof FilePlus2 }
> = {
  created: {
    label: 'created',
    tone: 'text-emerald-400 bg-emerald-500/10 border-emerald-500/25',
    Icon: FilePlus2,
  },
  modified: {
    label: 'modified',
    tone: 'text-amber-400 bg-amber-500/10 border-amber-500/25',
    Icon: FilePen,
  },
  deleted: {
    label: 'deleted',
    tone: 'text-rose-400 bg-rose-500/10 border-rose-500/25',
    Icon: FileX2,
  },
}

const TYPING_SPEED = 22
const LINE_DELAY = 55

export function CliExplorer() {
  const [query, setQuery] = useState('')
  const [category, setCategory] = useState<CommandCategory | 'All'>('All')
  const [selectedId, setSelectedId] = useState<string>(CLI_COMMANDS[0].id)
  const searchRef = useRef<HTMLInputElement>(null)

  // Deep links: /docs/cli#generate-resource opens that command.
  useEffect(() => {
    const hash = window.location.hash.replace('#', '')
    if (hash && CLI_COMMANDS.some((c) => c.id === hash)) setSelectedId(hash)
  }, [])

  // "/" focuses search, the shortcut every reference page should have.
  useEffect(() => {
    function onKey(e: KeyboardEvent) {
      const target = e.target as HTMLElement | null
      const typing =
        target &&
        (target.tagName === 'INPUT' ||
          target.tagName === 'TEXTAREA' ||
          target.isContentEditable)
      if (e.key === '/' && !typing) {
        e.preventDefault()
        searchRef.current?.focus()
      }
      if (e.key === 'Escape' && document.activeElement === searchRef.current) {
        setQuery('')
        searchRef.current?.blur()
      }
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [])

  const results = useMemo(() => {
    const q = query.trim().toLowerCase()
    return CLI_COMMANDS.filter((c) => {
      if (category !== 'All' && c.category !== category) return false
      if (!q) return true
      // File paths are searchable on purpose: "which command writes
      // routes.go" is a real question and the answer is otherwise a grep.
      const haystack = [
        c.name,
        c.alias ?? '',
        c.summary,
        c.purpose,
        c.example,
        ...(c.keywords ?? []),
        ...c.useCases,
        ...c.files.map((f) => f.path),
      ]
        .join(' ')
        .toLowerCase()
      return q.split(/\s+/).every((term) => haystack.includes(term))
    })
  }, [query, category])

  const selected =
    results.find((c) => c.id === selectedId) ?? results[0] ?? CLI_COMMANDS[0]

  function select(id: string) {
    setSelectedId(id)
    if (typeof window !== 'undefined') {
      window.history.replaceState(null, '', `#${id}`)
    }
  }

  return (
    <div className="not-prose">
      {/* Search + filters */}
      <div className="mb-6 space-y-3">
        <div className="relative">
          <Search className="pointer-events-none absolute left-3.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground/60" />
          <input
            ref={searchRef}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Search commands, use cases, or a file path like routes.go"
            className="w-full rounded-xl border border-border/60 bg-card/40 py-3 pl-10 pr-20 text-sm outline-none transition-colors placeholder:text-muted-foreground/50 focus:border-primary/50"
          />
          {query ? (
            <button
              type="button"
              onClick={() => setQuery('')}
              aria-label="Clear search"
              className="absolute right-3 top-1/2 -translate-y-1/2 rounded-md p-1 text-muted-foreground/60 hover:bg-accent/30 hover:text-foreground"
            >
              <X className="h-4 w-4" />
            </button>
          ) : (
            <kbd className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 rounded border border-border/60 bg-muted/40 px-1.5 py-0.5 font-mono text-[11px] text-muted-foreground/60">
              /
            </kbd>
          )}
        </div>

        <div className="flex flex-wrap gap-2">
          <FilterChip
            label={`All (${CLI_COMMANDS.length})`}
            active={category === 'All'}
            onClick={() => setCategory('All')}
          />
          {CLI_CATEGORIES.map((cat) => {
            const count = CLI_COMMANDS.filter((c) => c.category === cat).length
            return (
              <FilterChip
                key={cat}
                label={`${cat} (${count})`}
                active={category === cat}
                onClick={() => setCategory(cat)}
              />
            )
          })}
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,300px)_minmax(0,1fr)]">
        {/* Command list */}
        <div className="lg:sticky lg:top-24 lg:max-h-[calc(100vh-8rem)] lg:overflow-y-auto lg:pr-1">
          {results.length === 0 ? (
            <p className="rounded-xl border border-border/50 bg-card/30 p-6 text-center text-sm text-muted-foreground">
              Nothing matches <span className="text-foreground">{query}</span>.
            </p>
          ) : (
            <ul className="space-y-1.5">
              {results.map((cmd) => {
                const active = cmd.id === selected.id
                return (
                  <li key={cmd.id}>
                    <button
                      type="button"
                      onClick={() => select(cmd.id)}
                      className={`group w-full rounded-lg border px-3 py-2.5 text-left transition-colors ${
                        active
                          ? 'border-primary/40 bg-primary/[0.07]'
                          : 'border-transparent hover:border-border/60 hover:bg-accent/20'
                      }`}
                    >
                      <div className="flex items-center gap-2">
                        <code
                          className={`truncate font-mono text-[13px] ${
                            active ? 'text-primary' : 'text-foreground'
                          }`}
                        >
                          {cmd.name}
                        </code>
                        <ChevronRight
                          className={`ml-auto h-3.5 w-3.5 shrink-0 transition-transform ${
                            active
                              ? 'translate-x-0 text-primary'
                              : '-translate-x-1 text-muted-foreground/0 group-hover:translate-x-0 group-hover:text-muted-foreground/50'
                          }`}
                        />
                      </div>
                      <p className="mt-0.5 line-clamp-2 text-[12px] leading-snug text-muted-foreground">
                        {cmd.summary}
                      </p>
                    </button>
                  </li>
                )
              })}
            </ul>
          )}
        </div>

        {/* Detail */}
        <CommandDetail key={selected.id} cmd={selected} />
      </div>
    </div>
  )
}

function FilterChip({
  label,
  active,
  onClick,
}: {
  label: string
  active: boolean
  onClick: () => void
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className={`rounded-lg border px-3 py-1.5 text-[12px] font-medium transition-colors ${
        active
          ? 'border-primary/40 bg-primary/10 text-primary'
          : 'border-border/60 text-muted-foreground hover:border-border hover:text-foreground'
      }`}
    >
      {label}
    </button>
  )
}

function CommandDetail({ cmd }: { cmd: CliCommand }) {
  const created = cmd.files.filter((f) => f.status === 'created').length
  const modified = cmd.files.filter((f) => f.status === 'modified').length
  const deleted = cmd.files.filter((f) => f.status === 'deleted').length

  return (
    <div className="min-w-0 space-y-6">
      <header>
        <div className="mb-2 flex flex-wrap items-center gap-2">
          <span
            className={`rounded-md border px-2 py-0.5 text-[11px] font-medium ${CATEGORY_TONE[cmd.category]}`}
          >
            {cmd.category}
          </span>
          {cmd.alias && (
            <span className="rounded-md border border-border/60 px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
              {cmd.alias}
            </span>
          )}
        </div>
        <h2 className="font-mono text-2xl font-semibold tracking-tight">{cmd.name}</h2>
        <p className="mt-2 leading-relaxed text-muted-foreground">{cmd.summary}</p>
      </header>

      <CommandTerminal cmd={cmd} />

      {/* Files */}
      <section>
        <div className="mb-3 flex flex-wrap items-baseline gap-x-3 gap-y-1">
          <h3 className="text-sm font-semibold">What it touches</h3>
          {cmd.files.length > 0 ? (
            <span className="font-mono text-[11px] text-muted-foreground">
              {created > 0 && <span className="text-emerald-400">+{created} created</span>}
              {created > 0 && (modified > 0 || deleted > 0) && ' · '}
              {modified > 0 && <span className="text-amber-400">~{modified} modified</span>}
              {modified > 0 && deleted > 0 && ' · '}
              {deleted > 0 && <span className="text-rose-400">-{deleted} deleted</span>}
            </span>
          ) : null}
        </div>

        {cmd.files.length === 0 ? (
          <p className="rounded-lg border border-border/50 bg-card/30 px-4 py-3 text-sm text-muted-foreground">
            Nothing. This command reads or runs; it does not write to your project.
          </p>
        ) : (
          <ul className="overflow-hidden rounded-xl border border-border/50 bg-card/30">
            {cmd.files.map((file, i) => {
              const meta = STATUS_META[file.status]
              return (
                <li
                  key={`${file.path}-${i}`}
                  className="flex items-start gap-3 border-b border-border/40 px-3 py-2 last:border-0"
                >
                  <meta.Icon className={`mt-0.5 h-3.5 w-3.5 shrink-0 ${meta.tone.split(' ')[0]}`} />
                  <code className="min-w-0 break-all font-mono text-[12.5px] text-foreground/90">
                    {file.path}
                  </code>
                  {file.note && (
                    <span className="ml-auto hidden shrink-0 pl-3 text-[11px] text-muted-foreground/60 sm:block">
                      {file.note}
                    </span>
                  )}
                </li>
              )
            })}
          </ul>
        )}
      </section>

      {/* Purpose + use cases */}
      <section className="grid gap-4 sm:grid-cols-2">
        <div className="rounded-xl border border-border/50 bg-card/30 p-4">
          <h3 className="mb-2 text-sm font-semibold">Why it exists</h3>
          <p className="text-[13px] leading-relaxed text-muted-foreground">{cmd.purpose}</p>
        </div>
        <div className="rounded-xl border border-border/50 bg-card/30 p-4">
          <h3 className="mb-2 text-sm font-semibold">When you reach for it</h3>
          <ul className="space-y-1.5">
            {cmd.useCases.map((use) => (
              <li key={use} className="flex gap-2 text-[13px] leading-relaxed text-muted-foreground">
                <span className="mt-[7px] h-1 w-1 shrink-0 rounded-full bg-primary/60" />
                {use}
              </li>
            ))}
          </ul>
        </div>
      </section>

      {cmd.flags && cmd.flags.length > 0 && (
        <section>
          <h3 className="mb-3 text-sm font-semibold">Flags</h3>
          <div className="overflow-hidden rounded-xl border border-border/50 bg-card/30">
            {cmd.flags.map((f) => (
              <div
                key={f.flag}
                className="grid gap-1 border-b border-border/40 px-4 py-2.5 last:border-0 sm:grid-cols-[minmax(0,220px)_minmax(0,1fr)] sm:gap-4"
              >
                <code className="font-mono text-[12.5px] text-primary/90">{f.flag}</code>
                <span className="text-[13px] text-muted-foreground">{f.desc}</span>
              </div>
            ))}
          </div>
        </section>
      )}

      {cmd.notes && cmd.notes.length > 0 && (
        <section className="rounded-xl border border-amber-500/25 bg-amber-500/[0.06] p-4">
          <h3 className="mb-2 text-sm font-semibold text-amber-300/90">Worth knowing first</h3>
          <ul className="space-y-2">
            {cmd.notes.map((note) => (
              <li key={note} className="flex gap-2 text-[13px] leading-relaxed text-muted-foreground">
                <span className="mt-[7px] h-1 w-1 shrink-0 rounded-full bg-amber-400/70" />
                {note}
              </li>
            ))}
          </ul>
        </section>
      )}

      {cmd.docs && cmd.docs.length > 0 && (
        <section className="flex flex-wrap gap-2">
          {cmd.docs.map((doc) => (
            <Link
              key={doc.href}
              href={doc.href}
              className="inline-flex items-center gap-1.5 rounded-lg border border-border/60 px-3 py-1.5 text-[12px] text-muted-foreground transition-colors hover:border-primary/40 hover:text-primary"
            >
              {doc.label}
              <ChevronRight className="h-3 w-3" />
            </Link>
          ))}
        </section>
      )}
    </div>
  )
}

/**
 * The simulated run.
 *
 * Types the command, then streams the output. Deliberately not instant: the
 * point is to show the shape of what comes back, and a wall of text appearing
 * at once is a code block with extra steps.
 */
function CommandTerminal({ cmd }: { cmd: CliCommand }) {
  const [typed, setTyped] = useState('')
  const [lines, setLines] = useState<string[]>([])
  const [phase, setPhase] = useState<'idle' | 'typing' | 'running' | 'done'>('idle')
  const [copied, setCopied] = useState(false)
  const timers = useRef<ReturnType<typeof setTimeout>[]>([])
  const bodyRef = useRef<HTMLDivElement>(null)

  // Follow the output as it streams. Without this the pane fills past its own
  // height and the reader watches the first ten lines of a thirty-line run,
  // which is the one thing a simulated terminal must not do.
  useEffect(() => {
    const body = bodyRef.current
    if (body) body.scrollTop = body.scrollHeight
  }, [lines])

  // Every timer is tracked so switching command mid-run cannot leave one
  // firing into an unmounted component, or worse, into the next command.
  function clearTimers() {
    timers.current.forEach(clearTimeout)
    timers.current = []
  }
  useEffect(() => clearTimers, [])

  function after(ms: number, fn: () => void) {
    timers.current.push(setTimeout(fn, ms))
  }

  function run() {
    clearTimers()
    setTyped('')
    setLines([])
    setPhase('typing')

    const text = cmd.example
    for (let i = 1; i <= text.length; i++) {
      after(i * TYPING_SPEED, () => setTyped(text.slice(0, i)))
    }

    const typingDone = text.length * TYPING_SPEED + 220
    after(typingDone, () => setPhase('running'))
    cmd.output.forEach((line, i) => {
      after(typingDone + 120 + i * LINE_DELAY, () => setLines((prev) => [...prev, line]))
    })
    after(typingDone + 160 + cmd.output.length * LINE_DELAY, () => setPhase('done'))
  }

  function reset() {
    clearTimers()
    setTyped('')
    setLines([])
    setPhase('idle')
  }

  async function copy() {
    await navigator.clipboard.writeText(cmd.example)
    setCopied(true)
    setTimeout(() => setCopied(false), 1500)
  }

  const running = phase === 'typing' || phase === 'running'

  return (
    <div className="overflow-hidden rounded-xl border border-border/60 bg-[#0b0b12]">
      {/* Chrome */}
      <div className="flex items-center gap-2 border-b border-white/[0.06] bg-white/[0.02] px-4 py-2.5">
        <span className="flex gap-1.5">
          <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
          <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
        </span>
        <span className="ml-1 flex items-center gap-1.5 font-mono text-[11px] text-white/40">
          <Terminal className="h-3 w-3" />
          myapp
        </span>

        <div className="ml-auto flex items-center gap-1.5">
          <button
            type="button"
            onClick={copy}
            className="inline-flex items-center gap-1 rounded-md px-2 py-1 font-mono text-[11px] text-white/45 transition-colors hover:bg-white/[0.06] hover:text-white/80"
          >
            {copied ? <Check className="h-3 w-3 text-emerald-400" /> : <Copy className="h-3 w-3" />}
            {copied ? 'copied' : 'copy'}
          </button>
          {phase === 'idle' ? (
            <button
              type="button"
              onClick={run}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary/90 px-2.5 py-1 font-mono text-[11px] font-medium text-white transition-colors hover:bg-primary"
            >
              <Play className="h-3 w-3" />
              run
            </button>
          ) : (
            <button
              type="button"
              onClick={reset}
              disabled={running}
              className="inline-flex items-center gap-1.5 rounded-md border border-white/10 px-2.5 py-1 font-mono text-[11px] text-white/60 transition-colors hover:bg-white/[0.06] hover:text-white/90 disabled:opacity-40"
            >
              <RotateCcw className="h-3 w-3" />
              {running ? 'running' : 'again'}
            </button>
          )}
        </div>
      </div>

      {/* Body */}
      <div
        ref={bodyRef}
        className="max-h-[420px] overflow-y-auto scroll-smooth px-4 py-3.5 font-mono text-[12.5px] leading-[1.7]"
      >
        <div className="flex gap-2">
          <span className="shrink-0 select-none text-primary/80">$</span>
          <span className="min-w-0 break-all text-white/90">
            {phase === 'idle' ? cmd.example : typed}
            {phase === 'typing' && (
              <span className="ml-0.5 inline-block h-[1.05em] w-[7px] translate-y-[2px] animate-pulse bg-primary/80" />
            )}
          </span>
        </div>

        {phase === 'idle' && (
          <p className="mt-3 select-none text-white/25">
            Press run to see what this does.
          </p>
        )}

        {lines.map((line, i) => (
          <pre key={i} className={`whitespace-pre-wrap break-words ${lineTone(line)}`}>
            {line || ' '}
          </pre>
        ))}

        {phase === 'running' && (
          <span className="inline-block h-[1.05em] w-[7px] translate-y-[2px] animate-pulse bg-white/40" />
        )}
      </div>
    </div>
  )
}

/** Colour a line by what the CLI actually prints at the start of it. */
function lineTone(line: string): string {
  const t = line.trim()
  if (t.startsWith('✓') || t.startsWith('✅')) return 'text-emerald-400/90'
  if (t.startsWith('✗') || t.startsWith('!')) return 'text-rose-400/90'
  if (t.startsWith('⚠')) return 'text-amber-400/90'
  if (t.startsWith('→') || t.startsWith('+') || t.startsWith('~')) return 'text-cyan-400/85'
  if (t.startsWith('#') || t.startsWith('//')) return 'text-white/30'
  return 'text-white/60'
}
