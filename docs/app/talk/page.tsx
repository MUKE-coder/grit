import Link from 'next/link'
import type { Metadata } from 'next'
import { Geist_Mono } from 'next/font/google'
import { ArrowRight, Presentation } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DeckViewer } from '@/components/talk/deck-viewer'
import { getTalkDeck } from '@/lib/talk-deck'
import { GRIT_VERSION } from '@/config/site'

// Every slide is set in Geist and Geist Mono. The site already loads Geist for
// its own body text; the mono is only needed here, so it is scoped to this page
// rather than added to the root layout and shipped to every route.
const geistMono = Geist_Mono({
  subsets: ['latin'],
  variable: '--font-geist-mono',
  weight: ['400', '500', '600'],
})

export const metadata: Metadata = {
  title: 'Grit at GopherCon Africa',
  description:
    'The GopherCon Africa talk, with the speaker notes: the twenty problems every project runs into, twelve the developer pays for and eight the client asks about, what it takes to generate the shape of an application you still own afterwards, and a client billing system built and deployed from an empty folder.',
  alternates: { canonical: 'https://gritframework.dev/talk' },
  openGraph: {
    title: 'Grit at GopherCon Africa',
    description:
      'Twenty problems every full-stack project runs into, and what a framework can honestly do about them. The full deck, with speaker notes.',
    url: 'https://gritframework.dev/talk',
  },
}

export default function TalkPage() {
  const deck = getTalkDeck('gophercon-africa')

  return (
    <div className={`${geistMono.variable} min-h-screen bg-background`}>
      <SiteHeader />

      <main className="mx-auto w-full max-w-6xl px-4 py-12 sm:px-6 lg:py-16">
        <div className="mb-10 flex flex-col gap-5">
          <p className="inline-flex items-center gap-2 font-mono text-[11px] uppercase tracking-[0.2em] text-muted-foreground">
            <Presentation className="h-3.5 w-3.5" />
            GopherCon Africa 2026
          </p>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            Grit, as presented
          </h1>
          <p className="max-w-3xl text-lg leading-relaxed text-muted-foreground">
            The whole talk, {deck.slides.length} slides with the speaker notes attached. It builds
            one argument: every project runs into the same two sets of problems, one that costs the
            developer the first two weeks and one that decides whether there is a project at all,
            and the question worth answering is whether the shape of an application can be
            generated, owned outright by the team that generated it, and still improve after they
            own it.
          </p>
        </div>

        <DeckViewer deck={deck} />

        <section className="mt-16 border-t border-border pt-10">
          <h2 className="mb-6 text-xl font-semibold tracking-tight">What it covers</h2>
          <ol className="flex flex-col divide-y divide-border">
            {deck.sections.map((section) => (
              <li key={section.start} className="flex gap-5 py-3">
                <Link
                  href={`#${section.start}`}
                  className="font-mono text-xs text-muted-foreground transition-colors hover:text-accent"
                  aria-label={`Jump to slide ${section.number}`}
                >
                  {String(section.number).padStart(2, '0')}
                </Link>
                <span className="text-sm leading-relaxed text-foreground/90">
                  {section.description}
                </span>
              </li>
            ))}
          </ol>
        </section>

        <section className="mt-14 rounded-xl border border-border bg-card/40 p-7">
          <h2 className="mb-3 text-xl font-semibold tracking-tight">
            Everything in the deck is a command you can run
          </h2>
          <p className="mb-6 max-w-2xl text-sm leading-relaxed text-muted-foreground">
            The application built on stage is the one in the quick start, and the numbers quoted
            are measured in the posts they come from rather than estimated for the slide.
          </p>
          <div className="flex flex-col gap-3 sm:flex-row">
            <Link
              href="/docs/getting-started/quick-start"
              className="group inline-flex h-11 items-center gap-2 rounded-full bg-primary px-7 text-sm font-semibold text-primary-foreground transition-colors hover:bg-fg/85"
            >
              Build one yourself
              <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="/blog"
              className="inline-flex h-11 items-center rounded-full border border-border/60 px-7 text-sm font-medium text-foreground transition-colors hover:bg-surface-2"
            >
              Read the write-ups
            </Link>
          </div>
          <p className="mt-6 font-mono text-xs text-muted-foreground/60">
            {`Grit v${GRIT_VERSION} · MIT licensed · Go + React`}
          </p>
        </section>
      </main>
    </div>
  )
}
