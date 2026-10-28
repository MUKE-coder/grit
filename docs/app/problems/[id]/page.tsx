import Link from 'next/link'
import { notFound } from 'next/navigation'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { PROBLEMS, getProblem } from '@/config/problems'

/**
 * One page per problem, linked from the slide that states it.
 *
 * The talk gives each problem six words and moves on, which is right on a stage
 * and useless to the person who photographs the slide and reads it on the bus.
 * These pages are where that person lands: the same problem, the detail a slide
 * has no room for, the commands, and what it does not do.
 *
 * Rendered from one template rather than twenty files, because the useful thing
 * about them is that they answer the same questions in the same order. Twenty
 * hand-written pages diverge, and by the sixth nobody can tell whether a missing
 * "what it does not do" means there isn't one or nobody wrote it.
 */

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const problem = getProblem(id)
  if (!problem) return { title: 'Not found' }
  return {
    title: `${problem.id}: ${problem.label} — Grit`,
    description: `${problem.symptom} What Grit does about it, the commands, and what it does not do.`,
    alternates: { canonical: `https://gritframework.dev/problems/${problem.slug}` },
  }
}

export function generateStaticParams() {
  return PROBLEMS.map((p) => ({ id: p.slug }))
}

export const dynamicParams = false

export default async function ProblemPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params
  const problem = getProblem(id)
  if (!problem) notFound()

  const index = PROBLEMS.findIndex((p) => p.slug === problem.slug)
  const previous = PROBLEMS[index - 1]
  const next = PROBLEMS[index + 1]
  const groupLabel =
    problem.group === 'developer' ? 'What it costs the developer' : 'What the client asks'

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-3xl px-6 py-16">
        <Link
          href="/problems"
          className="inline-flex items-center gap-2 text-sm text-muted-foreground transition-colors hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          All twenty problems
        </Link>

        <div className="mt-8 flex items-center gap-3">
          <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 font-mono text-sm font-semibold text-accent">
            {problem.id}
          </span>
          <span className="text-sm text-muted-foreground">{groupLabel}</span>
        </div>

        <h1 className="mt-4 text-3xl font-bold tracking-tight text-foreground md:text-4xl">
          {problem.headline}
        </h1>
        <p className="mt-4 text-lg leading-relaxed text-muted-foreground">{problem.symptom}</p>

        <section className="mt-12">
          <h2 className="text-xl font-semibold text-foreground">What actually goes wrong</h2>
          <div className="prose-grit mt-4">
            {problem.detail.map((paragraph) => (
              <p key={paragraph}>{paragraph}</p>
            ))}
          </div>
        </section>

        <section className="mt-12">
          <h2 className="text-xl font-semibold text-foreground">What Grit does about it</h2>
          <ul className="mt-4 space-y-3">
            {problem.answer.map((line) => (
              <li key={line} className="flex gap-3 text-muted-foreground">
                <span aria-hidden="true" className="mt-2 h-1.5 w-1.5 shrink-0 rounded-full bg-primary" />
                <span className="leading-relaxed">{line}</span>
              </li>
            ))}
          </ul>
        </section>

        <section className="mt-12">
          <h2 className="text-xl font-semibold text-foreground">The commands</h2>
          <div className="mt-4 space-y-4">
            {problem.commands.map((command) => (
              <CodeBlock
                key={command.code}
                language="bash"
                filename={command.label}
                code={command.code}
              />
            ))}
          </div>
        </section>

        <section className="mt-12 rounded-xl border border-border bg-bg-secondary p-6">
          <h2 className="text-base font-semibold text-foreground">What it does not do</h2>
          <p className="mt-2 leading-relaxed text-muted-foreground">{problem.limit}</p>
        </section>

        {problem.docs.length > 0 && (
          <section className="mt-12">
            <h2 className="text-xl font-semibold text-foreground">Read next</h2>
            <div className="mt-4 flex flex-col gap-2">
              {problem.docs.map((doc) => (
                <Link
                  key={doc.href}
                  href={doc.href}
                  className="inline-flex items-center gap-2 text-accent transition-opacity hover:opacity-80"
                >
                  {doc.label}
                  <ArrowRight className="h-4 w-4" />
                </Link>
              ))}
            </div>
          </section>
        )}

        <nav className="mt-16 flex items-center justify-between gap-4 border-t border-border pt-8 text-sm">
          {previous ? (
            <Link
              href={`/problems/${previous.slug}`}
              className="inline-flex items-center gap-2 text-muted-foreground transition-colors hover:text-foreground"
            >
              <ArrowLeft className="h-4 w-4" />
              {previous.id} {previous.label}
            </Link>
          ) : (
            <span />
          )}
          {next && (
            <Link
              href={`/problems/${next.slug}`}
              className="inline-flex items-center gap-2 text-right text-muted-foreground transition-colors hover:text-foreground"
            >
              {next.id} {next.label}
              <ArrowRight className="h-4 w-4" />
            </Link>
          )}
        </nav>
      </main>
    </div>
  )
}
