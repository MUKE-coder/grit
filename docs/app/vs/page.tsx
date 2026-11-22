import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata = {
  title: 'Comparisons',
  description:
    'Grit next to PocketBase, Encore and Laravel: how each is built, what it costs to leave, and when the other one is the better answer.',
  alternates: { canonical: 'https://gritframework.dev/vs' },
}

const PAGES = [
  {
    href: '/vs/pocketbase',
    title: 'Grit and PocketBase',
    blurb:
      'A server you write against, or a generator that writes a codebase you own. The fastest start against the lowest ceiling.',
  },
  {
    href: '/vs/encore',
    title: 'Grit and Encore',
    blurb:
      'Infrastructure declared in code, or the application plumbing generated. Two halves of the same week.',
  },
  {
    // The Laravel comparison has a page already and a longer one. Sending
    // somebody to a second, thinner version of it would split the subject in
    // two and rank neither.
    href: '/i/laravel-alternative-go',
    title: 'Grit and Laravel',
    blurb:
      'What transfers from Laravel, the four things that do not, and whether moving is worth it.',
  },
]

export default function Page() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-3xl px-6 py-16">
        <h1 className="mb-4 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Comparisons
        </h1>
        <p className="mb-12 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Short, and written by one side, which is worth remembering. Each sticks to how the two
          are built rather than to numbers, and each has a section saying plainly when the other
          one is the better choice.
        </p>

        <ul className="space-y-3">
          {PAGES.map((p) => (
            <li key={p.href}>
              <Link
                href={p.href}
                className="group flex items-start justify-between gap-6 rounded-xl border border-border/60 bg-card/30 p-5 transition-colors hover:border-accent/50"
              >
                <span>
                  <span className="block font-display text-lg font-semibold text-foreground">
                    {p.title}
                  </span>
                  <span className="mt-1.5 block text-sm leading-relaxed text-muted-foreground">
                    {p.blurb}
                  </span>
                </span>
                <ArrowRight
                  className="mt-1 h-4 w-4 shrink-0 text-accent opacity-0 transition-opacity group-hover:opacity-100"
                  aria-hidden="true"
                />
              </Link>
            </li>
          ))}
        </ul>

        <section className="mt-14 border-t border-border/60 pt-10">
          <h2 className="mb-3 font-display text-xl font-semibold tracking-tight">
            Looking at the field rather than one rival
          </h2>
          <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
            <Link
              href="/guides/go-full-stack-frameworks-compared"
              className="text-accent hover:underline"
            >
              Choosing a Go full-stack approach
            </Link>{' '}
            groups the options by the kind of bargain each one makes and measures them by exit
            cost, including six cases where Grit is the wrong choice.{' '}
            <Link href="/i/go-full-stack-framework" className="text-accent hover:underline">
              The survey page
            </Link>{' '}
            covers eight tools individually.
          </p>
        </section>
      </main>
    </div>
  )
}
