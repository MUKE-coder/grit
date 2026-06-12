import Link from 'next/link'
import { BookOpen, ArrowRight } from 'lucide-react'

type Prereq = 'golang' | 'nextjs' | 'docker'

const META: Record<Prereq, { label: string; href: string; blurb: string }> = {
  golang: {
    label: 'Go (Golang) primer',
    href: '/docs/prerequisites/golang',
    blurb: 'Variables, structs, pointers, slices, maps, interfaces, goroutines.',
  },
  nextjs: {
    label: 'Next.js primer',
    href: '/docs/prerequisites/nextjs',
    blurb: 'App Router, server vs client components, data fetching, routing.',
  },
  docker: {
    label: 'Docker primer',
    href: '/docs/prerequisites/docker',
    blurb: 'Images, containers, volumes, networks, docker compose.',
  },
}

interface Props {
  /** Which prerequisite pages to nudge the reader toward. */
  prereqs: Prereq[]
  /** Optional intro sentence; otherwise a sensible default. */
  intro?: React.ReactNode
}

/**
 * PrereqLinks — small inline card that links to the prerequisites docs.
 * Drop this near the top of a lesson when the lesson uses concepts a
 * beginner may not yet know (e.g., Go syntax in the API course, Next App
 * Router in the Web course).
 */
export function PrereqLinks({ prereqs, intro }: Props) {
  return (
    <div className="not-prose my-6 rounded-2xl border border-sky-500/25 bg-sky-500/[0.04] p-4">
      <div className="flex items-start gap-3">
        <div className="h-9 w-9 rounded-xl bg-sky-500/15 border border-sky-500/30 flex items-center justify-center shrink-0">
          <BookOpen className="h-4 w-4 text-sky-400" />
        </div>
        <div className="flex-1 min-w-0">
          <p className="text-sm text-foreground/85 leading-relaxed mb-3">
            {intro ?? (
              <>
                New to one of these? Skim the primers first so the rest of
                this lesson sticks.
              </>
            )}
          </p>
          <div className="flex flex-col gap-1.5">
            {prereqs.map((p) => {
              const m = META[p]
              return (
                <Link
                  key={p}
                  href={m.href}
                  className="group flex items-start gap-2 rounded-lg border border-border/50 bg-card/50 hover:border-sky-500/40 hover:bg-sky-500/[0.06] px-3 py-2 transition-colors"
                >
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground group-hover:text-sky-400 transition-colors">
                      {m.label}
                    </p>
                    <p className="text-xs text-muted-foreground leading-relaxed">
                      {m.blurb}
                    </p>
                  </div>
                  <ArrowRight className="h-3.5 w-3.5 text-muted-foreground/70 group-hover:text-sky-400 group-hover:translate-x-0.5 transition-all shrink-0 mt-0.5" />
                </Link>
              )
            })}
          </div>
        </div>
      </div>
    </div>
  )
}
