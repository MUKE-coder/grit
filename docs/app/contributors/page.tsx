import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowUpRight, GitPullRequest, MessageSquare, Bug } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { getContributors } from '@/lib/contributors'

export const metadata: Metadata = {
  title: 'Contributors',
  description:
    'The people who have shipped code into Grit, the Go + React full-stack framework, and what their work became.',
  alternates: { canonical: 'https://gritframework.dev/contributors' },
}

// An hour, matching lib/contributors.ts. This page is a list of people: it does
// not need to be fresher than that, and ISR means a GitHub outage costs a stale
// page rather than a broken one.
export const revalidate = 3600

/**
 * What a contribution became, for the ones worth pointing at.
 *
 * Deliberately not a changelog. The GitHub API counts commits, so somebody who
 * argued a design decision into existence and then shipped it reads as "1", and
 * a page that printed only the number would be telling the opposite of the
 * truth. These are the sentences that say what actually happened.
 */
const SHIPPED: Record<string, { text: string; href: string }[]> = {
  daljo25: [
    {
      text: 'argued for grit plugin new over a template repository, and was right about why',
      href: 'https://github.com/MUKE-coder/grit/issues/143',
    },
    {
      text: 'built grit setup, the guided installer',
      href: 'https://github.com/MUKE-coder/grit/pull/170',
    },
  ],
}

function ContributorCard({
  login,
  avatarUrl,
  profileUrl,
  commits,
}: {
  login: string
  avatarUrl: string
  profileUrl: string
  commits: number
}) {
  const shipped = SHIPPED[login]
  return (
    <a
      href={profileUrl}
      target="_blank"
      rel="noopener noreferrer"
      className="group flex gap-4 rounded-xl border border-border/40 bg-card/40 p-5 transition-colors hover:border-accent/40 hover:bg-card/70"
    >
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={avatarUrl}
        alt=""
        width={48}
        height={48}
        loading="lazy"
        className="h-12 w-12 shrink-0 rounded-full border border-border/40"
      />
      <div className="min-w-0 flex-1">
        <div className="flex items-start justify-between gap-2">
          {/* No truncation. A name cut to "jdjchelp-..." is a name the page
              failed to show, and the fix is a column that fits rather than an
              ellipsis. break-words lets a long handle wrap instead. */}
          <span className="font-medium break-words leading-snug text-foreground">{login}</span>
          <ArrowUpRight className="mt-0.5 h-4 w-4 shrink-0 text-muted-foreground/60 transition-colors group-hover:text-accent" />
        </div>
        <p className="mt-0.5 text-[13px] text-muted-foreground">
          {commits.toLocaleString()} {commits === 1 ? 'commit' : 'commits'}
        </p>
        {shipped && (
          <ul className="mt-3 space-y-1.5">
            {shipped.map((item) => (
              <li key={item.href} className="text-[13px] leading-relaxed text-muted-foreground">
                {item.text}
              </li>
            ))}
          </ul>
        )}
      </div>
    </a>
  )
}

export default async function ContributorsPage() {
  const { contributors, live } = await getContributors()

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <main className="mx-auto w-full max-w-4xl px-6 pb-24 pt-16">
        <h1 className="text-4xl font-semibold tracking-tight">Contributors</h1>
        <p className="mt-4 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          Grit is one maintainer and a short list of people who have shipped code into it. The
          list is short on purpose: it is read off GitHub rather than written by hand, so nobody
          is on it who is not.
        </p>

        <section className="mt-12">
          <div className="grid gap-4 sm:grid-cols-2">
            {contributors.map((c) => (
              <ContributorCard key={c.login} {...c} />
            ))}
          </div>
          {!live && (
            <p className="mt-4 text-xs text-muted-foreground/70">
              GitHub could not be reached, so this is the last list we have. The live one is on{' '}
              <a
                href="https://github.com/MUKE-coder/grit/graphs/contributors"
                className="underline hover:text-accent"
                rel="noopener noreferrer"
                target="_blank"
              >
                the repository
              </a>
              .
            </p>
          )}
        </section>

        {/* The honest paragraph, and the reason there is no "total
            contributions" number at the top of this page. One person has
            written about 99% of the commits, so that number would be a
            statement about the maintainer dressed up as a statement about the
            project. */}
        <section className="mt-16 rounded-2xl border border-border/40 bg-card/30 p-7">
          <h2 className="text-lg font-semibold tracking-tight">Commits are not the whole list</h2>
          <p className="mt-3 leading-relaxed text-muted-foreground">
            GitHub counts commits, so the numbers above undersell almost everybody but the
            maintainer. The clearest example is on this page: <code>daljo25</code> shows one
            commit. That commit is <code>grit setup</code>, and before it they filed the issue
            that became <code>grit plugin new</code> and made the argument for building it as a
            command rather than a template repository, which is the reason it works the way it
            does.
          </p>
          <div className="mt-6 grid gap-4 sm:grid-cols-3">
            {[
              {
                Icon: GitPullRequest,
                title: 'Pull requests',
                body: 'The part that counts here. Read CONTRIBUTING.md first: it is short.',
              },
              {
                Icon: MessageSquare,
                title: 'Issues and reviews',
                body: 'A well-argued issue has changed the shape of this framework more than once.',
              },
              {
                Icon: Bug,
                title: 'Bug reports',
                body: 'Especially from a real project. Those are the ones that find what tests do not.',
              },
            ].map(({ Icon, title, body }) => (
              <div key={title}>
                <Icon className="h-4 w-4 text-accent" aria-hidden />
                <h3 className="mt-2 text-sm font-medium text-foreground">{title}</h3>
                <p className="mt-1 text-[13px] leading-relaxed text-muted-foreground">{body}</p>
              </div>
            ))}
          </div>
        </section>

        <section className="mt-10 flex flex-wrap items-center gap-3">
          <a
            href="https://github.com/MUKE-coder/grit/blob/main/CONTRIBUTING.md"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-lg bg-foreground px-5 py-2.5 text-sm font-medium text-background transition-opacity hover:opacity-90"
          >
            How to contribute
          </a>
          <a
            href="https://github.com/MUKE-coder/grit/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22"
            target="_blank"
            rel="noopener noreferrer"
            className="rounded-lg border border-border/60 px-5 py-2.5 text-sm font-medium transition-colors hover:border-accent/50 hover:text-accent"
          >
            Good first issues
          </a>
          <Link
            href="/sponsors"
            className="px-2 py-2.5 text-sm text-muted-foreground transition-colors hover:text-accent"
          >
            Or sponsor the work
          </Link>
        </section>
      </main>
    </div>
  )
}
