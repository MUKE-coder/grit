'use client'

import { useMemo, useState } from 'react'
import {
  Server,
  Globe,
  LayoutDashboard,
  Smartphone,
  Monitor,
  Check,
  Download,
  ArrowLeft,
  ArrowRight,
  RotateCcw,
  Puzzle,
} from 'lucide-react'
import { CopyButton } from '@/components/copy-button'
import {
  buildFrameworkPrompt,
  deriveCommands,
  type ClientId,
  type Frontend,
  type PluginId,
} from '@/lib/framework-prompt'

const CLIENTS: { id: ClientId; label: string; desc: string; icon: typeof Server }[] = [
  { id: 'api', label: 'API', desc: 'Go backend (Gin + GORM) — the base every client talks to.', icon: Server },
  { id: 'website', label: 'Website', desc: 'Public web frontend — marketing, storefront, SaaS, dashboard.', icon: Globe },
  { id: 'admin', label: 'Admin panel', desc: 'Generated Filament-like dashboard — tables, forms, roles.', icon: LayoutDashboard },
  { id: 'mobile', label: 'Mobile app', desc: 'Expo / React Native — iOS + Android, shared API.', icon: Smartphone },
  { id: 'desktop', label: 'Desktop app', desc: 'Native Wails window — offline-first with local SQLite.', icon: Monitor },
]

const PLUGINS: { id: PluginId; label: string; desc: string }[] = [
  { id: 'multitenant', label: 'Multi-tenancy', desc: 'Organizations, per-org roles, automatic query scoping.' },
  { id: 'impersonate', label: 'Impersonate', desc: 'Admin signs in as another user, with an audit trail.' },
  { id: 'command-palette', label: 'Command palette', desc: '⌘K navigation across the admin. Frontend-only.' },
  { id: 'saved-views', label: 'Saved views', desc: 'Per-user named table views (filters + sort).' },
]

const ALL_PLUGINS: PluginId[] = PLUGINS.map((p) => p.id)
const THEMES = ['atlas', 'aurora', 'pulse']
const STEPS = ['Clients', 'Stack', 'Plugins', 'Prompt'] as const

export function AIIntegrationWizard() {
  const [step, setStep] = useState<0 | 1 | 2 | 3>(0)
  const [clients, setClients] = useState<ClientId[]>(['api'])
  const [frontend, setFrontend] = useState<Frontend>('next')
  const [theme, setTheme] = useState('atlas')
  const [pluginMode, setPluginMode] = useState<'none' | 'all' | 'pick'>('none')
  const [picked, setPicked] = useState<PluginId[]>([])
  const [projectName, setProjectName] = useState('my-app')
  // Which command option the user picked. Held by id (stable across frontend/
  // theme edits); falls back to the first/recommended option when the id no
  // longer exists (e.g. after changing which clients are selected).
  const [chosenId, setChosenId] = useState<string>('')

  const plugins: PluginId[] = pluginMode === 'all' ? ALL_PLUGINS : pluginMode === 'pick' ? picked : []
  const hasWebUI = clients.includes('website') || clients.includes('admin')

  const commands = useMemo(
    () => deriveCommands({ clients, frontend, theme, projectName }),
    [clients, frontend, theme, projectName],
  )

  // Resolve the chosen option, defaulting to the recommended (first) one.
  const chosen = commands.options.find((o) => o.id === chosenId) ?? commands.options[0]

  const prompt = useMemo(
    () => buildFrameworkPrompt({ clients, frontend, theme, plugins, projectName, command: chosen.command }),
    [clients, frontend, theme, plugins, projectName, chosen.command],
  )

  function toggleClient(id: ClientId) {
    setClients((prev) => (prev.includes(id) ? prev.filter((c) => c !== id) : [...prev, id]))
  }
  function togglePlugin(id: PluginId) {
    setPicked((prev) => (prev.includes(id) ? prev.filter((p) => p !== id) : [...prev, id]))
  }

  function download() {
    const blob = new Blob([prompt], { type: 'text/markdown;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'grit-llm-prompt.md'
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
  }

  const canNext = step === 0 ? clients.length > 0 : true

  return (
    <div className="rounded-2xl border border-border bg-card/30">
      {/* Step indicator */}
      <div className="flex flex-wrap items-center gap-2 border-b border-border/60 px-6 py-4">
        {STEPS.map((label, i) => (
          <div key={label} className="flex items-center gap-2">
            <div
              className={
                'flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold ' +
                (i < step
                  ? 'bg-primary text-primary-foreground'
                  : i === step
                    ? 'border border-primary text-primary'
                    : 'border border-border text-muted-foreground')
              }
            >
              {i < step ? <Check className="h-3.5 w-3.5" /> : i + 1}
            </div>
            <span className={'text-sm ' + (i === step ? 'font-medium text-foreground' : 'text-muted-foreground')}>
              {label}
            </span>
            {i < STEPS.length - 1 && <div className="mx-1 h-px w-4 bg-border sm:w-8" />}
          </div>
        ))}
      </div>

      <div className="p-6">
        {/* ── Step 1: Clients ─────────────────────────────── */}
        {step === 0 && (
          <div>
            <h2 className="mb-2 text-xl font-semibold tracking-tight">What are you building?</h2>
            <p className="mb-6 text-sm text-muted-foreground">
              Pick every surface you need — the API is the base, the rest are clients that share it.
              Select as many as apply.
            </p>
            <div className="grid gap-3 sm:grid-cols-2">
              {CLIENTS.map((c) => {
                const on = clients.includes(c.id)
                const Icon = c.icon
                return (
                  <button
                    key={c.id}
                    onClick={() => toggleClient(c.id)}
                    className={
                      'flex items-start gap-3 rounded-xl border p-4 text-left transition-colors ' +
                      (on ? 'border-primary bg-primary/5' : 'border-border hover:border-primary/40')
                    }
                  >
                    <span
                      className={
                        'mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border ' +
                        (on ? 'border-primary bg-primary text-primary-foreground' : 'border-border')
                      }
                    >
                      {on && <Check className="h-3.5 w-3.5" />}
                    </span>
                    <div className="min-w-0">
                      <div className="flex items-center gap-2">
                        <Icon className="h-4 w-4 text-primary" />
                        <span className="font-medium text-foreground">{c.label}</span>
                      </div>
                      <p className="mt-1 text-sm text-muted-foreground">{c.desc}</p>
                    </div>
                  </button>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Step 2: Stack + commands ────────────────────── */}
        {step === 1 && (
          <div>
            <h2 className="mb-2 text-xl font-semibold tracking-tight">Assemble your stack</h2>
            <p className="mb-6 text-sm text-muted-foreground">
              Based on your clients, here&apos;s the command to scaffold it — plus a few variants.
            </p>

            <div className="mb-6 grid gap-4 sm:grid-cols-2">
              <label className="block">
                <span className="mb-1.5 block text-sm font-medium text-foreground">Project name</span>
                <input
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value.replace(/[^a-zA-Z0-9-_]/g, ''))}
                  placeholder="my-app"
                  className="w-full rounded-lg border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
                />
              </label>
              {hasWebUI && (
                <div>
                  <span className="mb-1.5 block text-sm font-medium text-foreground">Frontend framework</span>
                  <div className="flex gap-2">
                    {(['next', 'vite'] as Frontend[]).map((f) => (
                      <button
                        key={f}
                        onClick={() => setFrontend(f)}
                        className={
                          'flex-1 rounded-lg border px-3 py-2 text-sm ' +
                          (frontend === f ? 'border-primary bg-primary/5 text-foreground' : 'border-border text-muted-foreground')
                        }
                      >
                        {f === 'next' ? 'Next.js' : 'Vite / TanStack'}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {hasWebUI && (
              <div className="mb-6">
                <span className="mb-1.5 block text-sm font-medium text-foreground">Theme</span>
                <div className="flex flex-wrap gap-2">
                  {THEMES.map((t) => (
                    <button
                      key={t}
                      onClick={() => setTheme(t)}
                      className={
                        'rounded-lg border px-3 py-1.5 text-sm capitalize ' +
                        (theme === t ? 'border-primary bg-primary/5 text-foreground' : 'border-border text-muted-foreground')
                      }
                    >
                      {t}
                    </button>
                  ))}
                </div>
              </div>
            )}

            <p className="mb-2 text-sm font-medium text-foreground">
              Choose your scaffold command
              {commands.options.length > 1 && (
                <span className="ml-1 font-normal text-muted-foreground">
                  — the selected one is used in your prompt.
                </span>
              )}
            </p>
            <div className="space-y-2">
              {commands.options.map((o, i) => {
                const selected = o.id === chosen.id
                return (
                  <div
                    key={o.id}
                    className={
                      'flex items-start gap-2 rounded-xl border p-2 pr-2 transition-colors ' +
                      (selected ? 'border-primary bg-primary/5' : 'border-border')
                    }
                  >
                    <button
                      type="button"
                      onClick={() => setChosenId(o.id)}
                      aria-pressed={selected}
                      className="flex min-w-0 flex-1 items-start gap-3 rounded-lg p-2 text-left"
                    >
                      <span
                        className={
                          'mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded-full border ' +
                          (selected ? 'border-primary bg-primary text-primary-foreground' : 'border-border')
                        }
                      >
                        {selected && <Check className="h-3 w-3" />}
                      </span>
                      <span className="min-w-0 flex-1">
                        <span className="mb-1 flex flex-wrap items-center gap-2">
                          {i === 0 && (
                            <span className="rounded bg-primary/15 px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-primary">
                              Recommended
                            </span>
                          )}
                          <span className="text-xs text-muted-foreground">{o.label}</span>
                        </span>
                        <code className="block break-all font-mono text-sm text-foreground">{o.command}</code>
                      </span>
                    </button>
                    <CopyButton text={o.command} className="mt-1" />
                  </div>
                )
              })}
            </div>
          </div>
        )}

        {/* ── Step 3: Plugins ─────────────────────────────── */}
        {step === 2 && (
          <div>
            <h2 className="mb-2 text-xl font-semibold tracking-tight">Add plugins?</h2>
            <p className="mb-6 text-sm text-muted-foreground">
              Plugins generate reversible code into the repo. Add none, all, or pick the ones you want.
            </p>

            <div className="mb-4 flex flex-wrap gap-2">
              {(['none', 'all', 'pick'] as const).map((m) => (
                <button
                  key={m}
                  onClick={() => setPluginMode(m)}
                  className={
                    'rounded-lg border px-4 py-2 text-sm ' +
                    (pluginMode === m ? 'border-primary bg-primary/5 text-foreground' : 'border-border text-muted-foreground')
                  }
                >
                  {m === 'none' ? 'None' : m === 'all' ? 'All' : 'Pick each'}
                </button>
              ))}
            </div>

            {(pluginMode === 'pick' || pluginMode === 'all') && (
              <div className="grid gap-3 sm:grid-cols-2">
                {PLUGINS.map((p) => {
                  const on = pluginMode === 'all' || picked.includes(p.id)
                  return (
                    <button
                      key={p.id}
                      disabled={pluginMode === 'all'}
                      onClick={() => togglePlugin(p.id)}
                      className={
                        'flex items-start gap-3 rounded-xl border p-4 text-left transition-colors ' +
                        (on ? 'border-primary bg-primary/5' : 'border-border hover:border-primary/40')
                      }
                    >
                      <span
                        className={
                          'mt-0.5 flex h-5 w-5 shrink-0 items-center justify-center rounded border ' +
                          (on ? 'border-primary bg-primary text-primary-foreground' : 'border-border')
                        }
                      >
                        {on && <Check className="h-3.5 w-3.5" />}
                      </span>
                      <div className="min-w-0">
                        <div className="flex items-center gap-2">
                          <Puzzle className="h-4 w-4 text-primary" />
                          <span className="font-mono text-sm font-medium text-foreground">{p.id}</span>
                        </div>
                        <p className="mt-1 text-sm text-muted-foreground">{p.desc}</p>
                      </div>
                    </button>
                  )
                })}
              </div>
            )}
          </div>
        )}

        {/* ── Step 4: Prompt ──────────────────────────────── */}
        {step === 3 && (
          <div>
            <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
              <div>
                <h2 className="text-xl font-semibold tracking-tight">Your Grit prompt is ready</h2>
                <p className="mt-1 text-sm text-muted-foreground">
                  A full brief that teaches an AI agent Grit from zero — paste it into Claude Code,
                  Cursor, or any coding agent, or download it as <code>grit-llm-prompt.md</code>.
                </p>
              </div>
              <div className="flex items-center gap-2">
                <CopyButton
                  text={prompt}
                  variant="default"
                  className="bg-primary text-primary-foreground hover:bg-primary/90"
                >
                  Copy prompt
                </CopyButton>
                <button
                  onClick={download}
                  className="inline-flex items-center gap-2 rounded-lg border border-border bg-background px-3 py-2 text-sm font-medium text-foreground hover:border-primary/40"
                >
                  <Download className="h-4 w-4" />
                  Download .md
                </button>
              </div>
            </div>

            <div className="max-h-[28rem] overflow-y-auto rounded-xl border border-border bg-background p-4">
              <pre className="whitespace-pre-wrap break-words font-mono text-[12.5px] leading-relaxed text-muted-foreground">
                {prompt}
              </pre>
            </div>
          </div>
        )}
      </div>

      {/* Nav */}
      <div className="flex items-center justify-between border-t border-border/60 px-6 py-4">
        <button
          onClick={() => (step === 0 ? undefined : setStep((s) => (s - 1) as 0 | 1 | 2 | 3))}
          disabled={step === 0}
          className="inline-flex items-center gap-1.5 rounded-lg px-3 py-2 text-sm text-muted-foreground hover:text-foreground disabled:opacity-40"
        >
          <ArrowLeft className="h-4 w-4" />
          Back
        </button>

        {step < 3 ? (
          <button
            onClick={() => canNext && setStep((s) => (s + 1) as 0 | 1 | 2 | 3)}
            disabled={!canNext}
            className="inline-flex items-center gap-1.5 rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-40"
          >
            Next
            <ArrowRight className="h-4 w-4" />
          </button>
        ) : (
          <button
            onClick={() => setStep(0)}
            className="inline-flex items-center gap-1.5 rounded-lg border border-border px-4 py-2 text-sm text-muted-foreground hover:text-foreground"
          >
            <RotateCcw className="h-4 w-4" />
            Start over
          </button>
        )}
      </div>
    </div>
  )
}
