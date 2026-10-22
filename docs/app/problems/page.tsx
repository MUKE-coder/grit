import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CLIENT_PROBLEMS, DEVELOPER_PROBLEMS, type Problem } from '@/config/problems'

export const metadata: Metadata = {
  title: 'The twenty problems — Grit',
  description:
    'Twelve problems that cost the developer the first two weeks, and eight the client asks about before anyone mentions a feature. What each one actually costs, what Grit does about it, and what it does not do.',
  alternates: { canonical: 'https://gritframework.dev/problems' },
}

function ProblemRow({ problem }: { problem: Problem }) {
  return (
    <Link
      href={`/problems/${problem.slug}`}
      className="group flex items-start gap-4 rounded-xl border border-border bg-bg-secondary p-5 transition-colors hover:border-primary/40"
    >
      <span className="mt-0.5 inline-flex shrink-0 items-center rounded-lg bg-accent/15 px-2.5 py-1 font-mono text-xs font-semibold text-primary">
        {problem.id}
      </span>
      <span className="min-w-0 flex-1">
        <span className="block font-medium text-foreground">{problem.label}</span>
        <span className="mt-1 block text-sm leading-relaxed text-muted-foreground">
          {problem.headline}
        </span>
      </span>
      <ArrowRight className="mt-1 h-4 w-4 shrink-0 text-muted-foreground transition-colors group-hover:text-primary" />
    </Link>
  )
}

export default function ProblemsPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-3xl px-6 py-16">
        <h1 className="text-3xl font-bold tracking-tight text-foreground md:text-4xl">
          The twenty problems
        </h1>
        <p className="mt-4 text-lg leading-relaxed text-muted-foreground">
          Every full-stack project runs into the same two lists, and they belong to different
          people. Twelve cost the developer the first two weeks. Eight decide whether the client
          signs at all, and they are asked before anybody mentions a feature.
        </p>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          Each one here says what actually goes wrong, what Grit does about it, the commands, and
          what it does not do. They are the backbone of{' '}
          <Link href="/talk" className="text-primary hover:underline">
            the talk
          </Link>
          .
        </p>

        <section className="mt-12">
          <h2 className="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            What it costs the developer
          </h2>
          <div className="mt-4 flex flex-col gap-3">
            {DEVELOPER_PROBLEMS.map((problem) => (
              <ProblemRow key={problem.slug} problem={problem} />
            ))}
          </div>
        </section>

        <section className="mt-12">
          <h2 className="font-mono text-xs uppercase tracking-widest text-muted-foreground">
            What the client asks
          </h2>
          <div className="mt-4 flex flex-col gap-3">
            {CLIENT_PROBLEMS.map((problem) => (
              <ProblemRow key={problem.slug} problem={problem} />
            ))}
          </div>
        </section>
      </main>
    </div>
  )
}
