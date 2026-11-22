import Link from 'next/link'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Errors',
  description:
    'Every error the Grit CLI explains, each with its cause and its fix: port in use, Docker not running, go or node not found, a Go or Node version below the floor, pnpm missing.',
  alternates: { canonical: 'https://gritframework.dev/e' },
}

/**
 * The index of the explanation pages.
 *
 * Hand-listed rather than generated from internal/clierr, because a config
 * -driven list is invisible to both of the site's indexers. The Go catalogue
 * and the pages are paired by a test, so this list cannot quietly fall behind
 * either: TestEveryProblemHasAPage fails when a slug has no page.
 */
const ERRORS = [
  {
    slug: 'port-in-use',
    title: 'A port the project needs is already taken',
    cause: 'Usually a server from an earlier run whose terminal was closed.',
  },
  {
    slug: 'docker-not-running',
    title: 'Docker is not running',
    cause: 'The daemon did not answer, which usually means Docker Desktop is closed.',
  },
  {
    slug: 'go-not-found',
    title: 'go: command not found',
    cause: 'Go is not installed, or is installed and not on PATH.',
  },
  {
    slug: 'go-too-old',
    title: 'Your Go is older than Grit needs',
    cause: 'Older patches carry standard-library vulnerabilities a project reaches.',
  },
  {
    slug: 'node-not-found',
    title: 'node: command not found',
    cause: 'The frontends need Node. An API-only project does not.',
  },
  {
    slug: 'node-too-old',
    title: 'Your Node is older than the frontend needs',
    cause: 'Next declares a floor, and an old Node fails at build rather than at install.',
  },
  {
    slug: 'pnpm-not-found',
    title: 'pnpm: command not found',
    cause: 'Generated projects are a pnpm workspace, so npm produces a broken tree.',
  },
  {
    slug: 'not-a-grit-project',
    title: 'This is not a Grit project',
    cause: 'No grit.json here or in any directory above.',
  },
]

export default function Page() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-3xl px-6 py-16">
        <h1 className="mb-4 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Errors
        </h1>
        <p className="mb-8 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          When the CLI stops, it prints the address of the page that explains why. These are those
          pages: one problem each, the fix first, the reasoning underneath.
        </p>

        <div className="mb-12">
          <p className="mb-3 text-sm text-muted-foreground">
            To check for all of them at once, before you start:
          </p>
          <CodeBlock terminal code={`grit doctor --env`} />
        </div>

        <ul className="space-y-2.5">
          {ERRORS.map((e) => (
            <li key={e.slug}>
              <Link
                href={`/e/${e.slug}`}
                className="block rounded-xl border border-border/60 bg-card/30 p-4 transition-colors hover:border-accent/50"
              >
                <span className="block font-display text-base font-semibold text-foreground">
                  {e.title}
                </span>
                <span className="mt-1 block text-sm text-muted-foreground">{e.cause}</span>
                <span className="mt-2 block font-mono text-[11px] text-muted-foreground/50">
                  /e/{e.slug}
                </span>
              </Link>
            </li>
          ))}
        </ul>

        <section className="mt-14 border-t border-border/60 pt-8">
          <p className="text-sm leading-relaxed text-muted-foreground">
            A message not listed here is worth reading about in{' '}
            <Link href="/guides/common-go-docker-errors" className="text-accent hover:underline">
              the errors guide
            </Link>
            , which covers every message a Go, Docker and Node project produces in its first
            hour, or in{' '}
            <Link
              href="/docs/getting-started/troubleshooting"
              className="text-accent hover:underline"
            >
              the troubleshooting page
            </Link>
            , which is about Grit specifically.
          </p>
        </section>
      </main>
    </div>
  )
}
