import Link from 'next/link'
import { ArrowRight, ListChecks, Wand2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'

export const metadata = {
  title: 'Tools',
  description:
    'Small interactive tools for working with Grit: build a --fields string visually, and pick the architecture and flags for what you are making.',
  alternates: { canonical: 'https://gritframework.dev/tools' },
}

const TOOLS = [
  {
    href: '/tools/stack-picker',
    icon: ListChecks,
    title: 'Stack picker',
    blurb:
      'Five questions about what you are building, and the grit new command for it, with the architecture explained rather than just named.',
  },
  {
    href: '/tools/field-builder',
    icon: Wand2,
    title: 'Field builder',
    blurb:
      'Pick from all 32 field types, set relationships, choices and modifiers, and copy a --fields string that the parser accepts.',
  },
]

export default function Page() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-4xl px-6 py-16">
        <h1 className="mb-4 font-display text-4xl font-bold tracking-tight md:text-5xl">Tools</h1>
        <p className="mb-12 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Two things the CLI asks you to remember, done in a browser instead. Both produce a
          command you copy and run; neither needs an account and neither sends anything anywhere.
        </p>

        <ul className="grid gap-4 sm:grid-cols-2">
          {TOOLS.map((t) => (
            <li key={t.href}>
              <Link
                href={t.href}
                className="group flex h-full flex-col rounded-2xl border border-border/60 bg-card/30 p-6 transition-colors hover:border-accent/50"
              >
                <t.icon className="mb-4 h-5 w-5 text-accent" aria-hidden="true" />
                <h2 className="font-display text-lg font-semibold text-foreground">{t.title}</h2>
                <p className="mt-2 flex-1 text-sm leading-relaxed text-muted-foreground">
                  {t.blurb}
                </p>
                <span className="mt-4 inline-flex items-center gap-1 text-sm text-accent opacity-0 transition-opacity group-hover:opacity-100">
                  Open
                  <ArrowRight className="h-3.5 w-3.5" aria-hidden="true" />
                </span>
              </Link>
            </li>
          ))}
        </ul>

        <section className="mt-16 border-t border-border/60 pt-10">
          <h2 className="mb-3 font-display text-xl font-semibold tracking-tight">
            Elsewhere on the site
          </h2>
          <p className="max-w-2xl text-sm leading-relaxed text-muted-foreground">
            <Link href="/playground" className="text-accent hover:underline">
              The playground
            </Link>{' '}
            runs a generated project in the browser.{' '}
            <Link href="/docs/glossary" className="text-accent hover:underline">
              The glossary
            </Link>{' '}
            defines every term in one sentence. And{' '}
            <Link href="/guides" className="text-accent hover:underline">
              the guides
            </Link>{' '}
            answer one problem each, in plain Go and React.
          </p>
        </section>
      </main>
    </div>
  )
}
