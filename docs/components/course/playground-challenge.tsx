import Link from 'next/link'
import { Play, ArrowRight } from 'lucide-react'

interface Props {
  /** Short challenge title — what the reader will attempt. */
  title: string
  /** A 1-2 sentence brief describing the challenge. */
  children: React.ReactNode
  /** Optional href into the playground. Defaults to /playground. */
  href?: string
}

/**
 * PlaygroundChallenge — drops a "try this in the playground" card into
 * a lesson. Used to give learners a tiny, immediate code-run challenge
 * tied to what they just read, before they go and modify their actual
 * project.
 */
export function PlaygroundChallenge({ title, children, href = '/playground' }: Props) {
  return (
    <div className="not-prose my-6 rounded-2xl border border-emerald-500/25 bg-emerald-500/[0.04] p-5">
      <div className="flex items-start gap-3">
        <div className="h-9 w-9 rounded-xl bg-emerald-500/15 border border-emerald-500/30 flex items-center justify-center shrink-0">
          <Play className="h-4 w-4 text-emerald-400" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-[10px] uppercase tracking-wider text-emerald-400 font-mono mb-1">
            Playground challenge
          </p>
          <h4 className="text-sm font-semibold mb-1.5 text-foreground">{title}</h4>
          <div className="text-sm text-foreground/85 leading-relaxed mb-3 prose-grit max-w-none">
            {children}
          </div>
          <Link
            href={href}
            target="_blank"
            className="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/15 hover:bg-emerald-500/25 border border-emerald-500/30 text-emerald-300 text-xs font-medium px-3 py-1.5 transition-colors"
          >
            Open the playground
            <ArrowRight className="h-3 w-3" />
          </Link>
        </div>
      </div>
    </div>
  )
}
