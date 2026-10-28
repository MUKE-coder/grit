import Link from 'next/link'
import { ArrowRight, Puzzle } from 'lucide-react'
import { CodeBlock } from '@/components/code-block'

/**
 * One step in a worked example.
 *
 * `run` is the command; `then` is what it does, in a sentence. A step whose
 * explanation is "runs the command" is not a step, it is a line of shell, and
 * it should be folded into the one above it.
 */
export interface PluginStep {
  run?: string
  language?: string
  title: string
  then: string
}

/**
 * A related plugin, with the reason it is related.
 *
 * `why` is required on purpose. A row of names under "see also" is decoration;
 * what makes it useful is knowing whether the next plugin solves the problem
 * you actually have, and that is the sentence.
 */
export interface RelatedPlugin {
  name: string
  href: string
  why: string
}

/**
 * PluginSteps renders a worked example on a real project.
 *
 * The steps are passed in rather than read from a shared config because the
 * CLI's docs indexer only extracts const arrays a page maps over in its own
 * file. Content kept in a config is invisible to `grit docs search`, silently,
 * which is the worst way for documentation to be missing.
 */
export function PluginSteps({
  steps,
  intro,
}: {
  steps: PluginStep[]
  intro?: string
}) {
  return (
    <>
      <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Using it, step by step</h2>
      {intro ? <p className="mb-6 leading-relaxed text-muted-foreground">{intro}</p> : null}

      <ol className="mb-8 space-y-6">
        {steps.map((step, i) => (
          <li key={step.title} className="relative pl-10">
            <span
              className="absolute left-0 top-0 flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 font-mono text-xs font-semibold text-primary"
              aria-hidden="true"
            >
              {i + 1}
            </span>
            <p className="mb-2 font-medium leading-snug">{step.title}</p>
            {step.run ? (
              <div className="mb-2">
                <CodeBlock language={step.language ?? 'bash'} code={step.run} />
              </div>
            ) : null}
            <p className="text-sm leading-relaxed text-muted-foreground">{step.then}</p>
          </li>
        ))}
      </ol>
    </>
  )
}

/** RelatedPlugins is the "what else is near this" block at the foot of a page. */
export function RelatedPlugins({ items }: { items: RelatedPlugin[] }) {
  if (items.length === 0) return null

  return (
    <>
      <h2 className="mb-2 mt-12 text-2xl font-semibold tracking-tight">Related plugins</h2>
      <p className="mb-5 leading-relaxed text-muted-foreground">
        Near neighbours, and what each one is for instead.
      </p>
      <div className="mb-6 space-y-2">
        {items.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className="group flex items-start gap-3 rounded-xl border border-border p-4 transition-colors hover:border-primary"
          >
            <Puzzle className="mt-0.5 h-4 w-4 shrink-0 text-primary" aria-hidden="true" />
            <span className="flex-1">
              <span className="block font-medium">{item.name}</span>
              <span className="text-sm text-muted-foreground">{item.why}</span>
            </span>
            <ArrowRight className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
          </Link>
        ))}
      </div>
      <p className="text-sm text-muted-foreground">
        Everything else is on the{' '}
        <Link href="/docs/plugins/marketplace" className="text-primary hover:underline">
          marketplace
        </Link>
        , and{' '}
        <Link href="/docs/extending" className="text-primary hover:underline">
          four ways to extend
        </Link>{' '}
        says when a plugin is the wrong tool for the job.
      </p>
    </>
  )
}
