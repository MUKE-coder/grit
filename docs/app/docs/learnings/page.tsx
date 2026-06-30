import Link from 'next/link'
import { ArrowRight, Lightbulb, Activity, Calendar } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/learnings')

interface Entry {
  slug: string
  title: string
  summary: string
  date: string
  tags: string[]
  milestone: string
  icon: React.ComponentType<{ className?: string }>
}

const ENTRIES: Entry[] = [
  {
    slug: 'stateless-service-load-test',
    title: 'Stateless service + k6 load test',
    summary:
      'Scaffold a stateless Go API with grit new --api, load-test the /api/health endpoint with k6, and capture a latency chart (p50 / p95 / p99) of the run.',
    date: '2026-05-31',
    tags: ['Go', 'Gin', 'k6', 'Performance'],
    milestone: 'A committed latency chart of the service under load',
    icon: Activity,
  },
]

export default function LearningsIndexPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-4xl mx-auto">
            {/* Header */}
            <div className="mb-12">
              <span className="tag-mono text-primary/80 mb-3 inline-flex items-center gap-1.5">
                <Lightbulb className="h-3.5 w-3.5" />
                Learnings
              </span>
              <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-4 leading-tight">
                Engineering journal
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-3xl">
                A running log of hands-on challenges I work through on top of the Grit framework.
                Each entry sticks to the same shape — the problem, the build, the numbers, and
                what I took away. Every script is copy-pasteable and every milestone is
                reproducible.
              </p>
            </div>

            {/* Entries */}
            <div className="space-y-3">
              {ENTRIES.map((e) => {
                const Icon = e.icon
                return (
                  <Link
                    key={e.slug}
                    href={`/docs/learnings/${e.slug}`}
                    className="group block rounded-2xl border border-border/50 bg-card/40 p-5 hover:border-primary/40 hover:bg-card/70 transition-all"
                  >
                    <div className="flex items-start gap-4">
                      <div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center shrink-0">
                        <Icon className="h-5 w-5" />
                      </div>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2 mb-1 flex-wrap">
                          <h2 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors">
                            {e.title}
                          </h2>
                          <span className="inline-flex items-center gap-1 text-[10px] font-mono text-muted-foreground">
                            <Calendar className="h-3 w-3" />
                            {e.date}
                          </span>
                        </div>
                        <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                          {e.summary}
                        </p>
                        <div className="flex items-center gap-2 flex-wrap mb-3">
                          {e.tags.map((t) => (
                            <span
                              key={t}
                              className="text-[10px] font-mono text-foreground/70 bg-card/80 border border-border/40 rounded-full px-2 py-0.5"
                            >
                              {t}
                            </span>
                          ))}
                        </div>
                        <p className="text-xs text-muted-foreground/80">
                          <span className="text-primary font-semibold mr-1.5">MILESTONE</span>
                          {e.milestone}
                        </p>
                      </div>
                      <ArrowRight className="h-4 w-4 text-muted-foreground/60 group-hover:text-primary group-hover:translate-x-0.5 transition-all shrink-0 mt-1" />
                    </div>
                  </Link>
                )
              })}
            </div>

            {/* More to come */}
            <div className="mt-10 rounded-xl border border-dashed border-border/40 bg-card/20 p-6 text-center">
              <p className="text-sm text-muted-foreground">
                More entries coming as I work through them. Each will land here with the same
                shape: problem · build · numbers · takeaways.
              </p>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
