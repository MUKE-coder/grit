'use client'

import { useState } from 'react'
import { ChevronDown, ChevronRight, Target, Lightbulb } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Props {
  /** Imperative prompt — what the learner should try. */
  prompt: React.ReactNode
  /** Optional hint shown before they reveal the solution. */
  hint?: React.ReactNode
  /** The correct answer / sample solution. */
  solution: React.ReactNode
}

/**
 * Exercise — "Now you try". Renders an exercise prompt with an
 * optional hint and a collapsible "Reveal solution" panel. The
 * learner is encouraged to attempt first; the solution stays hidden
 * until they ask for it.
 */
export function Exercise({ prompt, hint, solution }: Props) {
  const [revealed, setRevealed] = useState(false)
  const [hintShown, setHintShown] = useState(false)

  return (
    <div className="not-prose my-8 rounded-2xl border border-primary/30 bg-primary/[0.04] overflow-hidden">
      <div className="px-5 py-4 border-b border-primary/20 bg-primary/[0.06]">
        <p className="text-[10px] uppercase tracking-wider text-primary font-mono mb-1 inline-flex items-center gap-1.5">
          <Target className="h-3 w-3" />
          Try it
        </p>
        <div className="text-sm text-foreground/90 leading-relaxed">{prompt}</div>
      </div>

      {hint && (
        <div className="px-5 py-3 border-b border-primary/10">
          <button
            type="button"
            onClick={() => setHintShown((v) => !v)}
            className="text-xs text-amber-500/90 hover:text-amber-500 inline-flex items-center gap-1.5 font-medium"
          >
            <Lightbulb className="h-3.5 w-3.5" />
            {hintShown ? 'Hide hint' : 'Need a hint?'}
          </button>
          {hintShown && (
            <div className="mt-2 text-sm text-foreground/80 leading-relaxed">{hint}</div>
          )}
        </div>
      )}

      <div className="px-5 py-3">
        <button
          type="button"
          onClick={() => setRevealed((v) => !v)}
          className={cn(
            'inline-flex items-center gap-1.5 text-sm font-medium rounded-md px-3 py-1.5 transition-colors',
            revealed
              ? 'text-muted-foreground hover:text-foreground'
              : 'bg-primary text-primary-foreground hover:bg-primary/90',
          )}
        >
          {revealed ? (
            <>
              <ChevronDown className="h-3.5 w-3.5" />
              Hide solution
            </>
          ) : (
            <>
              <ChevronRight className="h-3.5 w-3.5" />
              Reveal solution
            </>
          )}
        </button>
        {revealed && (
          <div className="mt-3 text-sm text-foreground/90 prose-grit max-w-none">
            {solution}
          </div>
        )}
      </div>
    </div>
  )
}
