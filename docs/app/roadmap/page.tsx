import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowUpRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import {
  ROADMAP,
  ROADMAP_REVIEWED,
  bucketLabel,
  bucketStyle,
  bucketBlurb,
  type Bucket,
} from '@/config/roadmap'

export const metadata: Metadata = {
  title: 'Roadmap',
  description:
    'What is being built next in Grit, what is agreed but unscheduled, what is waiting on a decision, and what Grit has decided not to do.',
  alternates: { canonical: 'https://gritframework.dev/roadmap' },
}

const ORDER: Bucket[] = ['next', 'later', 'question']

/**
 * Things already decided against, with the reason.
 *
 * Short on purpose. An issue nobody has got to is not "not planned", and
 * filing one here would be using the roadmap to close the backlog. Each of
 * these is a decision the project actually made and can be argued with.
 */
const NOT_PLANNED: { title: string; why: string }[] = [
  {
    title: 'A second web framework, ORM or CSS library',
    why: 'Gin, GORM and Tailwind are the stack, and the value of a convention-led framework is that there is one answer. Supporting Echo as well as Gin would double the generator and halve the conventions.',
  },
  {
    title: 'A Pages Router option for the Next apps',
    why: 'App Router only. Every generated page, layout and server action assumes it, and the second shape would be the one nobody tests.',
  },
  {
    title: 'Radix as the primitive library',
    why: 'Where a new primitive genuinely needs behaviour rather than appearance, it uses Base UI, which is the successor the shadcn maintainers point at. Existing hand-written primitives stay: they work and they are tested, and rewriting them would trade known behaviour for unknown.',
  },
  {
    title: 'A component registry inside every scaffolded project',
    why: 'It was there and was removed in v3.31.78. Every project carried a UIComponent model, a registry handler, a 91-component seeder and a browser page, for a feature almost nobody used. Grit UI is a standalone library instead.',
  },
  {
    title: 'An access column on the routes screen',
    why: 'gin exposes a route’s last handler and nothing in front of it, so the middleware chain cannot be recovered at runtime. The static table in internal/access is how anything learns which routes are protected, and inventing the column from path shapes would be a guess presented as a fact.',
  },
]

export default function RoadmapPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <header className="mb-12">
          <span className="tag-mono text-accent/80 mb-3 block">Roadmap</span>
          <h1 className="text-4xl font-bold tracking-tight mb-4">What is next</h1>
          <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl">
            Grit ships every merged change as a release, so the{' '}
            <Link href="/docs/changelog" className="text-accent hover:underline">
              changelog
            </Link>{' '}
            says what happened and says it in a lot of detail. This says what has not
            happened yet, which is the other half of deciding whether to build on
            something.
          </p>
          <p className="text-sm text-muted-foreground/80 leading-relaxed mt-4 max-w-2xl">
            Reviewed {ROADMAP_REVIEWED}. Every line links to its issue, which is where the
            real state is: this page is a statement of intent, and intent is the part that
            only a person can refresh. Nothing here is a date, because a one-maintainer
            project that promised dates would be lying on a schedule.
          </p>
        </header>

        {ORDER.map((bucket) => {
          const items = ROADMAP.filter((item) => item.bucket === bucket)
          if (items.length === 0) return null
          return (
            <section key={bucket} className="mb-14">
              <div className="flex flex-wrap items-center gap-3 mb-2">
                <h2 className="text-2xl font-semibold tracking-tight">
                  {bucketLabel[bucket]}
                </h2>
                <span
                  className={`inline-flex items-center rounded border px-2 py-0.5 text-xs font-medium ${bucketStyle[bucket]}`}
                >
                  {items.length}
                </span>
              </div>
              <p className="text-sm text-muted-foreground leading-relaxed mb-5 max-w-2xl">
                {bucketBlurb[bucket]}
              </p>

              <ul className="space-y-3">
                {items.map((item) => (
                  <li
                    key={item.issue}
                    className="rounded-xl border border-border/40 bg-card/30 p-5"
                  >
                    <a
                      href={`https://github.com/MUKE-coder/grit/issues/${item.issue}`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="group inline-flex items-start gap-1.5 font-medium hover:text-accent"
                    >
                      {item.title}
                      <ArrowUpRight className="h-3.5 w-3.5 shrink-0 mt-1 opacity-50 group-hover:opacity-100" />
                    </a>
                    <p className="text-sm text-muted-foreground leading-relaxed mt-1.5">
                      {item.note}
                    </p>
                    <p className="text-xs font-mono text-muted-foreground/60 mt-2">
                      #{item.issue}
                    </p>
                  </li>
                ))}
              </ul>
            </section>
          )
        })}

        <section className="mb-14">
          <h2 className="text-2xl font-semibold tracking-tight mb-2">Not planned</h2>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5 max-w-2xl">
            Decisions already made, with the reason, so you can disagree with the reason
            rather than guess at one. A request nobody has got to yet is in the list above,
            not here.
          </p>
          <ul className="space-y-3">
            {NOT_PLANNED.map((item) => (
              <li
                key={item.title}
                className="rounded-xl border border-border/30 bg-card/20 p-5"
              >
                <p className="font-medium">{item.title}</p>
                <p className="text-sm text-muted-foreground leading-relaxed mt-1.5">
                  {item.why}
                </p>
              </li>
            ))}
          </ul>
        </section>

        <section className="rounded-xl border border-accent/25 bg-accent/[0.04] p-6">
          <h2 className="text-lg font-semibold mb-2">Want one of these sooner</h2>
          <p className="text-sm text-muted-foreground leading-relaxed">
            Say so on its issue, with what you are building and what the absence costs you.
            That is not a formality: several items moved up because one person explained
            the case, and three of the things on{' '}
            <Link href="/docs/changelog" className="text-accent hover:underline">
              this month&apos;s changelog
            </Link>{' '}
            exist because somebody filed the issue rather than working around it.
          </p>
          <div className="mt-4 flex flex-wrap gap-4 text-sm">
            <a
              href="https://github.com/MUKE-coder/grit/issues"
              target="_blank"
              rel="noopener noreferrer"
              className="text-accent hover:underline inline-flex items-center gap-1"
            >
              The live issue list
              <ArrowUpRight className="h-3.5 w-3.5" />
            </a>
            <Link href="/docs/governance" className="text-accent hover:underline">
              Who maintains Grit, and how decisions are made
            </Link>
            <Link href="/contributors" className="text-accent hover:underline">
              Who has shipped into it
            </Link>
          </div>
        </section>
      </main>
    </div>
  )
}
