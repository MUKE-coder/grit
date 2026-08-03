import Link from 'next/link'
import { AlertCircle, ArrowRight, Github } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/governance')

export default function GovernancePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Project</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Governance &amp; risk</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Who maintains this, what happens if they stop, and what you are actually
                exposed to if you build on it. Written for the person who has to sign off on
                adopting it.
              </p>
            </div>

            {/* Lead with the objection rather than burying it. A team doing
                vendor review will find this out in five minutes anyway; the only
                question is whether they find it from us or from a stranger. */}
            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-6 mb-12">
              <div className="flex gap-3">
                <AlertCircle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                <div>
                  <h2 className="font-semibold mb-2">The bus factor is one</h2>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Grit, Pulse, Sentinel, gin-docs and GORM Studio are substantially the work
                    of one person. External review of this project has flagged that every
                    time, and it is correct. If you are evaluating Grit for something that has
                    to outlive its maintainer&apos;s interest, that is the real risk — not the
                    architecture.
                  </p>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Pretending otherwise would be worse than useless. What follows is what
                    actually reduces the exposure, and what does not.
                  </p>
                </div>
              </div>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">
              What genuinely limits the damage
            </h2>
            <div className="space-y-4 mb-12">
              {[
                {
                  title: 'Your application does not depend on Grit at runtime',
                  body:
                    'This is the whole answer, and it is structural rather than a promise. Grit generates Go and TypeScript into your repository and then gets out of the way. Your production binary does not import it. If this project stopped today, every app built with it keeps compiling, keeps deploying and keeps running — you would lose future generators, not your software.',
                },
                {
                  title: 'The generated stack is boring on purpose',
                  body:
                    'Gin, GORM, Postgres, Redis, asynq, React, Tailwind. Not one of them is ours. A developer who has never heard of Grit can read a generated handler and know exactly what it does, because it is ordinary code in well-known libraries. There is no bespoke runtime to learn or to inherit.',
                },
                {
                  title: 'MIT, and the whole thing is on GitHub',
                  body:
                    'Including the generators, the templates and this website. A fork is a viable escape hatch rather than a theoretical one.',
                },
              ].map((c) => (
                <div key={c.title} className="rounded-xl border border-border/50 bg-card/40 p-5">
                  <h3 className="font-semibold text-sm mb-2">{c.title}</h3>
                  <p className="text-[13px] text-muted-foreground leading-relaxed">{c.body}</p>
                </div>
              ))}
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">What does not help</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Being honest about the limits of the above:
            </p>
            <ul className="space-y-2.5 mb-12 text-sm text-muted-foreground">
              <li className="leading-relaxed">
                &mdash; <strong className="text-foreground">Pulse and Sentinel are runtime dependencies.</strong>{' '}
                Unlike the generated code, these are Go packages your app imports. If you need
                to remove that exposure, both are replaceable with standard middleware, but it
                is real work rather than a no-op.
              </li>
              <li className="leading-relaxed">
                &mdash; <strong className="text-foreground">A fork is only as good as your team&apos;s Go.</strong>{' '}
                &ldquo;It is open source&rdquo; means little if nobody on the team can maintain
                a code generator.
              </li>
              <li className="leading-relaxed">
                &mdash; <strong className="text-foreground">There is no support SLA.</strong> See{' '}
                <Link href="/docs/versioning" className="text-primary hover:underline">
                  versioning
                </Link>
                . Issues get answered because someone cares, not because someone is contractually
                obliged to.
              </li>
            </ul>

            <h2 className="text-2xl font-bold tracking-tight mb-4">
              How to reduce your exposure today
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              If you are adopting Grit for something that matters, do these three things and the
              maintainer question largely stops being your problem:
            </p>
            <CodeBlock
              language="bash"
              code={`# 1. Pin the CLI in CI. A release then cannot change your build under you.
go install github.com/MUKE-coder/grit/v3/cmd/grit@v3.115.0

# 2. Vendor your Go dependencies, so an upstream disappearing does not stop a build.
go mod vendor

# 3. Confirm the exit is real: your app builds and runs with no Grit installed.
which grit || echo "no grit on PATH"
cd apps/api && go build ./... && echo "builds without the CLI"`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4 mb-12">
              That third one is worth running before you commit to anything. It is the claim
              this whole page rests on, and it takes thirty seconds to verify rather than
              believe.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Where the roadmap lives</h2>
            <p className="text-muted-foreground leading-relaxed mb-6">
              There is no separate roadmap document, because a roadmap nobody updates is worse
              than none. What is true instead:
            </p>
            <div className="grid gap-3 sm:grid-cols-2 mb-12">
              {[
                {
                  href: '/docs/changelog',
                  title: 'Changelog',
                  body: 'Every release, what changed and why. The most honest signal of where effort is going.',
                },
                {
                  href: 'https://github.com/MUKE-coder/grit/issues',
                  title: 'Issues',
                  body: 'What is being worked on and what is being asked for. Filing one is the way to influence priority.',
                  external: true,
                },
                {
                  href: 'https://github.com/MUKE-coder/grit/commits/main',
                  title: 'Commit history',
                  body: 'Release cadence and activity, unfiltered by anything on this website.',
                  external: true,
                },
                {
                  href: '/docs/versioning',
                  title: 'Versioning policy',
                  body: 'What can break, what cannot, and what a version number means.',
                },
              ].map((c) => (
                <Link
                  key={c.href}
                  href={c.href}
                  target={c.external ? '_blank' : undefined}
                  rel={c.external ? 'noreferrer' : undefined}
                  className="group rounded-xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border hover:bg-card/70"
                >
                  <div className="flex items-center gap-1.5 font-semibold text-sm">
                    {c.title}
                    <ArrowRight className="h-3.5 w-3.5 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                  </div>
                  <p className="mt-1.5 text-[13px] text-muted-foreground leading-relaxed">{c.body}</p>
                </Link>
              ))}
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Contributing</h2>
            <p className="text-muted-foreground leading-relaxed mb-6">
              The fastest way to make the bus factor larger than one is for it to be larger than
              one. Generators live in <code className="text-xs">internal/generate</code>,
              scaffold templates in <code className="text-xs">internal/scaffold</code>, and both
              have test suites you can run before touching anything.
            </p>
            <CodeBlock
              language="bash"
              code={`git clone https://github.com/MUKE-coder/grit
cd grit
go test ./internal/...        # should be green before you start
go build -o grit ./cmd/grit   # your build of the CLI`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4 mb-8">
              One thing worth knowing before you start:{' '}
              <code className="text-xs">go build</code> compiles the scaffold templates as{' '}
              <em>strings</em>. It never checks the Go or TypeScript inside them. Generate a
              real project and build that — it is the only thing that catches a broken template.
            </p>

            <Link
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noreferrer"
              className="inline-flex h-11 items-center gap-2 rounded-lg border border-border/60 px-5 text-sm font-semibold transition-colors hover:bg-accent/50"
            >
              <Github className="h-4 w-4" />
              github.com/MUKE-coder/grit
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
