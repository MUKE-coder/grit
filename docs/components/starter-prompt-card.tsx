'use client'

import Link from 'next/link'
import { Sparkles, ArrowRight, ExternalLink } from 'lucide-react'
import { CopyButton } from '@/components/copy-button'

export interface StarterPromptCardProps {
  /** Human label for the kit, e.g. "Single — Go + embedded SPA". */
  kitLabel: string
  /** Exact CLI command for the kit, e.g. `grit new my-app --single`. */
  command: string
  /** Multi-line prompt body. Will be pasted verbatim into claude.ai. */
  prompt: string
}

/**
 * Drop-in section for each tech-kit page. Renders a copyable "starter
 * prompt" that the user pastes into claude.ai to get back the 4
 * planning files (project-description, project-phases, design-style-guide,
 * prompt.md) which they then feed to Claude Code.
 */
export function StarterPromptCard({ kitLabel, command, prompt }: StarterPromptCardProps) {
  return (
    <div className="mb-16">
      {/* Header */}
      <div className="mb-5 flex items-start justify-between gap-4 flex-wrap">
        <div>
          <div className="inline-flex items-center gap-2 mb-2">
            <Sparkles className="h-4 w-4 text-primary" />
            <span className="tag-mono text-primary/80">Starter prompt for claude.ai</span>
          </div>
          <h2 className="text-2xl font-semibold tracking-tight">
            Plan this kit with an AI
          </h2>
          <p className="text-muted-foreground max-w-2xl mt-1.5 leading-relaxed">
            Paste this into <Link href="https://claude.ai" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">claude.ai</Link> with
            your idea. Claude returns four planning files
            (<code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90">project-description.md</code>,
            <code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90"> project-phases.md</code>,
            <code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90"> design-style-guide.md</code>,
            <code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90"> prompt.md</code>)
            that you feed to Claude Code with the <code className="text-[0.85em] bg-primary/5 px-1 py-0.5 rounded text-primary/90">{command}</code> command.
          </p>
        </div>
        <Link
          href="/docs/ai-integration"
          className="shrink-0 inline-flex items-center gap-1.5 rounded-full border border-primary/30 bg-primary/[0.06] text-primary text-xs font-medium px-3.5 py-2 hover:bg-primary/[0.12] hover:border-primary/50 transition-colors"
        >
          Customize with the wizard
          <ArrowRight className="h-3.5 w-3.5" />
        </Link>
      </div>

      {/* Prompt body */}
      <div className="relative rounded-2xl border border-border/50 bg-card/50 overflow-hidden">
        {/* Top strip */}
        <div className="flex items-center justify-between gap-3 border-b border-border/40 bg-muted/30 px-4 py-2.5">
          <div className="flex items-center gap-2 min-w-0">
            <span className="h-2 w-2 rounded-full bg-rose-400/70" />
            <span className="h-2 w-2 rounded-full bg-amber-400/70" />
            <span className="h-2 w-2 rounded-full bg-emerald-400/70" />
            <span className="ml-2 text-[11px] font-mono text-muted-foreground truncate">
              prompt — {kitLabel}
            </span>
          </div>
          <CopyButton
            text={prompt}
            size="sm"
            variant="ghost"
            className="h-7 px-2.5 text-[11px] text-muted-foreground hover:text-foreground"
          >
            Copy prompt
          </CopyButton>
        </div>

        {/* Body — preserve newlines, soft wrap */}
        <pre className="px-5 py-5 text-[13px] leading-[1.65] font-mono text-foreground/85 whitespace-pre-wrap break-words max-h-[420px] overflow-y-auto">
          {prompt}
        </pre>
      </div>

      {/* How to use strip */}
      <div className="mt-4 grid grid-cols-1 md:grid-cols-3 gap-2.5">
        {[
          { n: '1', t: 'Copy', d: 'Hit Copy prompt above.' },
          { n: '2', t: 'Paste in claude.ai', d: 'Replace [YOUR IDEA] with a paragraph about your project.' },
          { n: '3', t: 'Feed Claude Code', d: 'Save the 4 files Claude returns to your project root, then run prompt.md in Claude Code.' },
        ].map((s) => (
          <div key={s.n} className="rounded-xl border border-border/40 bg-card/30 p-3">
            <div className="flex items-center gap-2 mb-1">
              <span className="h-5 w-5 rounded-full bg-primary/10 text-primary text-[10px] font-mono font-semibold flex items-center justify-center">
                {s.n}
              </span>
              <p className="text-xs font-semibold">{s.t}</p>
            </div>
            <p className="text-xs text-muted-foreground leading-relaxed">{s.d}</p>
          </div>
        ))}
      </div>

      {/* External Claude link */}
      <div className="mt-4 flex justify-end">
        <Link
          href="https://claude.ai"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors"
        >
          Open claude.ai in a new tab
          <ExternalLink className="h-3 w-3" />
        </Link>
      </div>
    </div>
  )
}
