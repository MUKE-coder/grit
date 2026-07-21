import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata: Metadata = {
  title: 'Grit vs — honest framework comparisons',
  description:
    'How Grit compares to Rails, Laravel, Django, Next.js, NestJS, Express, Encore, and the MERN and MEAN stacks. When to pick each, without the marketing spin.',
  alternates: { canonical: 'https://gritframework.dev/compare' },
}

const comparisons: { slug: string; name: string; blurb: string }[] = [
  { slug: 'nextjs', name: 'Next.js', blurb: 'A React framework with API routes vs a real Go backend + generated admin.' },
  { slug: 'rails', name: 'Ruby on Rails', blurb: 'The batteries-included original — on Ruby vs on Go, with a React admin.' },
  { slug: 'laravel', name: 'Laravel', blurb: 'PHP&apos;s productivity king vs Go + React, and Filament vs a generated admin.' },
  { slug: 'django', name: 'Django', blurb: 'Python&apos;s "batteries included" vs Go + a typed React frontend.' },
  { slug: 'nestjs', name: 'NestJS', blurb: 'A structured Node backend vs Go, with codegen and clients included.' },
  { slug: 'express', name: 'Express', blurb: 'A minimal Node router you assemble vs an opinionated full stack.' },
  { slug: 'encore', name: 'Encore', blurb: 'Two Go backend frameworks — infra-from-code vs batteries + a React admin.' },
  { slug: 'mern', name: 'MERN stack', blurb: 'Mongo + Express + React + Node glued by hand vs one generated stack.' },
  { slug: 'mean', name: 'MEAN stack', blurb: 'Mongo + Express + Angular + Node vs Go + React, generated.' },
]

export default function ComparePage() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <span className="font-mono text-xs uppercase tracking-wider text-primary">Compare</span>
        <h1 className="mb-4 mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
          How Grit compares
        </h1>
        <p className="mb-12 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Honest, side-by-side comparisons — where Grit wins, where the alternative wins, and how to
          choose. No trashing other tools; several of them inspired Grit.
        </p>

        <div className="grid gap-3 sm:grid-cols-2">
          {comparisons.map((c) => (
            <Link
              key={c.slug}
              href={`/compare/${c.slug}`}
              className="group flex items-start justify-between gap-3 rounded-2xl border border-border bg-card/40 p-5 transition-colors hover:border-primary/40"
            >
              <div>
                <div className="mb-1 font-semibold text-foreground group-hover:text-primary">
                  Grit vs {c.name}
                </div>
                <p className="text-sm leading-relaxed text-muted-foreground" dangerouslySetInnerHTML={{ __html: c.blurb }} />
              </div>
              <ArrowRight className="mt-1 h-4 w-4 shrink-0 text-muted-foreground/50 transition-colors group-hover:text-primary" />
            </Link>
          ))}
        </div>
      </main>
    </div>
  )
}
