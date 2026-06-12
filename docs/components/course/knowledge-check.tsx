'use client'

import { useState } from 'react'
import { Check, X, HelpCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Choice {
  label: string
  /** True for the correct answer(s). */
  correct?: boolean
  /** Feedback shown when this choice is selected. */
  feedback: string
}

interface Props {
  question: React.ReactNode
  choices: Choice[]
}

/**
 * KnowledgeCheck — a single multiple-choice question. The learner
 * picks a choice; immediate feedback explains why it's right or wrong.
 * Optimised for "learned something just now, did it stick?" — not
 * for graded testing.
 */
export function KnowledgeCheck({ question, choices }: Props) {
  const [picked, setPicked] = useState<number | null>(null)

  return (
    <div className="not-prose my-8 rounded-2xl border border-border/50 bg-card/30">
      <div className="px-5 py-4 border-b border-border/40">
        <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-mono mb-1.5 inline-flex items-center gap-1.5">
          <HelpCircle className="h-3 w-3" />
          Quick check
        </p>
        <div className="text-sm font-medium text-foreground/90 leading-relaxed">{question}</div>
      </div>
      <div className="px-5 py-4 space-y-2">
        {choices.map((c, i) => {
          const selected = picked === i
          const showState = picked !== null && selected
          return (
            <button
              key={i}
              type="button"
              onClick={() => setPicked(i)}
              className={cn(
                'w-full text-left rounded-lg border px-3 py-2.5 text-sm transition-colors flex items-start gap-2.5',
                !showState && 'border-border/40 hover:border-primary/40 hover:bg-card/60',
                showState && c.correct && 'border-emerald-500/40 bg-emerald-500/10',
                showState && !c.correct && 'border-rose-500/40 bg-rose-500/10',
              )}
            >
              <span
                className={cn(
                  'h-5 w-5 rounded-full border flex items-center justify-center shrink-0 mt-0.5',
                  !showState && 'border-border/60',
                  showState && c.correct && 'border-emerald-500 bg-emerald-500/20',
                  showState && !c.correct && 'border-rose-500 bg-rose-500/20',
                )}
              >
                {showState && c.correct && <Check className="h-3 w-3 text-emerald-500" />}
                {showState && !c.correct && <X className="h-3 w-3 text-rose-500" />}
              </span>
              <span className="flex-1">{c.label}</span>
            </button>
          )
        })}
        {picked !== null && (
          <div
            className={cn(
              'mt-3 p-3 rounded-md text-sm leading-relaxed',
              choices[picked].correct
                ? 'bg-emerald-500/10 text-emerald-200 border border-emerald-500/20'
                : 'bg-rose-500/10 text-rose-200 border border-rose-500/20',
            )}
          >
            <p className="text-xs font-mono uppercase tracking-wider mb-1 opacity-80">
              {choices[picked].correct ? 'Correct' : 'Not quite'}
            </p>
            <p className="text-foreground/90">{choices[picked].feedback}</p>
          </div>
        )}
      </div>
    </div>
  )
}
