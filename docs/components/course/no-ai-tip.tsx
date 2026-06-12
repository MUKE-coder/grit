import { Brain, Keyboard, X } from 'lucide-react'

interface Props {
  /** Override the default body. */
  children?: React.ReactNode
}

/**
 * NoAITip — sits at the top of every course / first lesson, reminding
 * students that the GOAL of this course is to *learn*, not ship fastest.
 * That means hand-typing the code and disabling AI suggestions while
 * working through lessons. The point isn't to be anti-AI — it's that
 * AI mid-completion robs you of the small mistakes that make concepts
 * stick.
 */
export function NoAITip({ children }: Props) {
  return (
    <div className="not-prose my-6 rounded-2xl border-2 border-amber-500/40 bg-gradient-to-br from-amber-500/[0.08] via-orange-500/[0.04] to-amber-500/[0.02] p-5">
      <div className="flex items-start gap-3">
        <div className="h-10 w-10 rounded-xl bg-amber-500/15 border border-amber-500/30 flex items-center justify-center shrink-0">
          <Brain className="h-5 w-5 text-amber-500" />
        </div>
        <div className="flex-1">
          <p className="text-[10px] uppercase tracking-wider text-amber-500 font-mono font-semibold mb-1.5 inline-flex items-center gap-1.5">
            <Keyboard className="h-3 w-3" />
            Disable AI suggestions while you learn
          </p>
          {children ? (
            <div className="text-sm text-foreground/90 leading-relaxed">{children}</div>
          ) : (
            <div className="text-sm text-foreground/90 leading-relaxed space-y-2">
              <p>
                This course teaches you to <strong>hand-write every line</strong> of
                code. Open VS Code (or your editor of choice) and{' '}
                <strong>turn off Copilot, Cursor Tab, Tabnine, Codeium</strong>,
                and any inline AI autocomplete <em>before you start a lesson</em>.
              </p>
              <p>
                AI mid-completion robs you of the small mistakes that make
                concepts stick. You&apos;ll be a faster, more independent
                developer at the end of the course if you type every
                character yourself. Re-enable AI for your real work after
                — never during a lesson.
              </p>
            </div>
          )}
          <p className="mt-3 text-xs text-muted-foreground inline-flex items-center gap-1.5">
            <X className="h-3 w-3 text-rose-400" />
            <span>Goal of this course: <em>learn</em>, not <em>ship fastest</em>.</span>
          </p>
        </div>
      </div>
    </div>
  )
}
