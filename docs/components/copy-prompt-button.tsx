'use client'

import * as React from 'react'
import Link from 'next/link'
import { Check, Copy, Sparkles } from 'lucide-react'
import { cn } from '@/lib/utils'
import { BUILD_WITH_AI_PROMPT } from '@/config/build-prompt'

type Variant = 'hero' | 'inline' | 'panel'

interface CopyPromptButtonProps {
  variant?: Variant
  className?: string
  /** Override the label. Keep it a verb phrase. */
  label?: string
}

/**
 * "Copy prompt to build with AI".
 *
 * One button, three sizes, one prompt. The prompt itself lives in
 * config/build-prompt.ts so that every placement copies the same text: a page
 * that quietly carried its own variant would be the version that goes stale.
 *
 * Clipboard writes fail in a few real situations: an insecure origin, a
 * browser that withholds the permission, an iframe without the policy. A
 * button that silently does nothing is worse than one that admits it, so the
 * failure path says what happened and offers the file instead.
 */
export function CopyPromptButton({ variant = 'inline', className, label }: CopyPromptButtonProps) {
  const [state, setState] = React.useState<'idle' | 'copied' | 'failed'>('idle')

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(BUILD_WITH_AI_PROMPT)
      setState('copied')
    } catch {
      setState('failed')
    }
    setTimeout(() => setState('idle'), 4000)
  }

  const text =
    state === 'copied'
      ? 'Prompt copied, paste it into your AI'
      : state === 'failed'
        ? 'Could not copy'
        : (label ?? 'Copy prompt to build with AI')

  const Icon = state === 'copied' ? Check : state === 'failed' ? Copy : Sparkles

  if (variant === 'panel') {
    return (
      <div
        className={cn(
          'rounded-xl border border-primary/20 bg-primary/[0.04] p-5 sm:p-6',
          className,
        )}
      >
        <div className="flex flex-wrap items-start justify-between gap-4">
          <div className="min-w-0">
            <div className="mb-1.5 inline-flex items-center gap-2">
              <Sparkles className="h-4 w-4 text-primary" />
              <span className="tag-mono text-primary/80">Build with AI</span>
            </div>
            <h3 className="text-lg font-semibold tracking-tight">
              Let your AI build it
            </h3>
            <p className="mt-1 max-w-xl text-sm leading-relaxed text-muted-foreground">
              One prompt that installs the Grit skill, teaches your agent the framework,
              and links every concept it needs. Works in Claude Code, Cursor, Codex,
              Windsurf, Copilot or any chat.
            </p>
          </div>
          <button type="button" onClick={copy} className={buttonClass('hero')}>
            <Icon className="h-4 w-4" aria-hidden="true" />
            {text}
          </button>
        </div>
        <Failed show={state === 'failed'} />
      </div>
    )
  }

  return (
    <div className={cn('inline-flex flex-col items-start gap-1.5', className)}>
      <button type="button" onClick={copy} className={buttonClass(variant)}>
        <Icon className="h-4 w-4" aria-hidden="true" />
        {text}
      </button>
      <Failed show={state === 'failed'} />
    </div>
  )
}

function buttonClass(variant: Variant): string {
  const base =
    'inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-background'
  if (variant === 'hero') {
    return cn(base, 'h-12 px-6 text-[15px] bg-primary text-primary-foreground hover:bg-primary/90')
  }
  return cn(
    base,
    'h-10 px-4 text-sm border border-primary/30 bg-primary/10 text-primary hover:bg-primary/15',
  )
}

function Failed({ show }: { show: boolean }) {
  if (!show) return null
  return (
    <p role="status" className="text-xs text-muted-foreground">
      Your browser blocked the clipboard.{' '}
      <Link href="/prompt" className="text-primary underline" target="_blank" rel="noopener noreferrer">
        Open the prompt
      </Link>{' '}
      and copy it by hand.
    </p>
  )
}
