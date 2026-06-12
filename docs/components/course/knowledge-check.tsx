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
    <div className="not-prose my-8 rounded-2xl border border-primary/30 bg-primary/[0.03]">
      <div className="px-5 py-4 border-b border-primary/20 bg-primary/[0.05]">
        <p className="text-[10px] uppercase tracking-wider text-primary font-mono mb-1.5 inline-flex items-center gap-1.5">
          <HelpCircle className="h-3.5 w-3.5" />
          Quick check
        </p>
        <div className="text-base font-semibold text-foreground leading-relaxed">{question}</div>
      </div>
      <div className="px-5 py-5 space-y-2.5">
        {choices.map((c, i) => {
          const selected = picked === i
          const showState = picked !== null && selected
          return (
            <button
              key={i}
              type="button"
              onClick={() => setPicked(i)}
              className={cn(
                'w-full text-left rounded-xl border-2 px-4 py-3 text-sm transition-all flex items-start gap-3 group',
                !showState && 'border-border bg-card/60 hover:border-primary/50 hover:bg-primary/[0.04]',
                showState && c.correct && 'border-emerald-500 bg-emerald-500/[0.08]',
                showState && !c.correct && 'border-rose-500 bg-rose-500/[0.08]',
              )}
            >
              {/* Radio with clear ring + filled state */}
              <span
                className={cn(
                  'relative h-6 w-6 rounded-full border-2 flex items-center justify-center shrink-0 mt-0.5 transition-colors',
                  !showState && 'border-border bg-background group-hover:border-primary/70',
                  showState && c.correct && 'border-emerald-500 bg-emerald-500',
                  showState && !c.correct && 'border-rose-500 bg-rose-500',
                )}
              >
                {showState && c.correct && <Check className="h-3.5 w-3.5 text-white" strokeWidth={3} />}
                {showState && !c.correct && <X className="h-3.5 w-3.5 text-white" strokeWidth={3} />}
                {!showState && (
                  <span className="h-2.5 w-2.5 rounded-full bg-transparent group-hover:bg-primary/30 transition-colors" />
                )}
              </span>
              <span className="flex-1 text-foreground/95 leading-relaxed">{c.label}</span>
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
