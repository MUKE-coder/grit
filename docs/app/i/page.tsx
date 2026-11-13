import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Guides: choosing a Go stack',
  description:
    'Long-form, honest guides to the decisions a Go team actually faces: full-stack frameworks, routers, admin panels, React frontends, moving from Laravel, and code generators.',
  alternates: { canonical: 'https://gritframework.dev/i' },
}

const guides = [
  {
    slug: 'go-full-stack-framework',
    title: 'Full-stack frameworks for Go',
    blurb:
      'Go has no Rails, and the phrase means three different things. Buffalo, Beego, Bud, templ with HTMX, Encore, PocketBase and Grit, compared.',
    reading: '12 min',
  },
  {
    slug: 'go-web-framework',
    title: 'Go web frameworks: net/http, Chi, Gin, Echo and Fiber',
    blurb:
      'Four of the five are routers, not frameworks. Where that line falls, and the list of things none of them give you.',
    reading: '10 min',
  },
  {
    slug: 'go-admin-panel',
    title: 'Admin panels for Go',
    blurb:
      'Why Go has no Filament, the four ways teams solve it, and the question that decides between them: where your schema lives.',
    reading: '9 min',
  },
  {
    slug: 'go-react-stack',
    title: 'A Go backend with a React frontend',
    blurb:
      'The four problems the tutorials skip: type drift, auth across two origins, the dev proxy, and one artefact or two.',
    reading: '9 min',
  },
  {
    slug: 'laravel-alternative-go',
    title: 'Coming to Go from Laravel',
    blurb:
      'Every part of Laravel mapped to its nearest Go equivalent, including the four gaps with no good answer.',
    reading: '8 min',
  },
  {
    slug: 'go-crud-generator',
    title: 'CRUD and scaffolding generators for Go',
    blurb:
      'sqlc, ent, GORM gen, goa and full-stack scaffolders. What each writes, and what it leaves you.',
    reading: '8 min',
  },
]

export default function GuidesPage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-accent">Guides</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Choosing a Go stack
        </h1>
        <p className="mb-12 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Long-form answers to the decisions a Go team actually faces, written to be useful whether
          or not you end up using Grit. Every project mentioned is checked against its own source,
          and the ones that have gone quiet are named as such.
        </p>

        <ul className="space-y-4">
          {guides.map((g) => (
            <li key={g.slug}>
              <Link
                href={`/i/${g.slug}`}
                className="group block rounded-xl border border-border/60 bg-card/30 p-6 transition-colors hover:border-accent/50"
              >
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <h2 className="font-display text-xl font-semibold text-foreground group-hover:text-accent">
                      {g.title}
                    </h2>
                    <p className="mt-2 leading-relaxed text-muted-foreground">{g.blurb}</p>
                  </div>
                  <ArrowRight className="mt-1 h-5 w-5 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-1 group-hover:text-accent" />
                </div>
                <p className="mt-4 font-mono text-xs text-muted-foreground/70">{g.reading} read</p>
              </Link>
            </li>
          ))}
        </ul>

        <div className="mt-12 rounded-xl border border-border/60 p-6">
          <h2 className="mb-2 font-display text-lg font-semibold">Looking for a direct comparison?</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            These guides are about the decision. If you already know the two things you are choosing
            between, the side-by-side pages are more direct.
          </p>
          <Link
            href="/compare"
            className="inline-flex items-center gap-2 text-sm font-semibold text-accent hover:underline"
          >
            Grit compared to nine alternatives
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </main>
    </div>
  )
}
