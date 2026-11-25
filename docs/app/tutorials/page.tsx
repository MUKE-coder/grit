import Link from 'next/link'
import { ArrowUpRight, Clock, Youtube } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { VideoFacade } from '@/components/video-facade'
import { VIDEOS } from '@/config/videos'

export const metadata = {
  // The root layout appends the site name, so the title is the page's own.
  title: 'Tutorials',
  description:
    'Grit on video: install the CLI and scaffold a project, generate a resource and its relationships, the admin panel you get without writing it, roles and permissions, deploying to a VPS, and building with an AI coding agent. Every video says which Grit version it was recorded on.',
  alternates: { canonical: 'https://gritframework.dev/tutorials' },
}

/**
 * The video tutorials library.
 *
 * ── Why this is not under /docs ──────────────────────────────────────────
 *
 * It is the same decision /guides and /plugins made. A video library is a way
 * in rather than a reference: somebody arrives from YouTube or a search, and
 * the docs shell around them, a sidebar of 230 reference pages, is furniture
 * for a different job. These pages get the marketing chrome instead, which is
 * a header and the page.
 *
 * The cost is explicit and worth knowing. `scripts/docs-index.py` builds the
 * index the binary embeds for `grit docs search`, and it walks app/docs and
 * nothing else, so this page is not in it, exactly as /guides and /plugins
 * are not. The website's own search does cover it, through a curated entry in
 * config/search-keywords.json, which is how /plugins is handled and for the
 * same reason: one page is a list entry, thirty would need a walker.
 *
 * ── Why the prose is specific rather than welcoming ─────────────────────
 *
 * What somebody can search for is what is written here literally, and
 * "watch and learn" indexes nothing. Naming the subjects is the page.
 *
 * ── Why there is no list of planned videos ──────────────────────────────
 *
 * A page that offers a video nobody can watch wastes the trip. The series is
 * tracked on an issue; an entry appears here when the video does.
 */
export default function TutorialsPage() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      {/* ------------------------------------------------------------ hero */}
      <section className="relative overflow-hidden border-b border-border">
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-0 opacity-[0.05]"
          style={{
            backgroundImage: 'radial-gradient(currentColor 1px, transparent 1px)',
            backgroundSize: '32px 32px',
          }}
        />
        <div className="relative mx-auto max-w-5xl px-6 py-16">
          <p className="font-mono text-[11px] uppercase tracking-[0.12em] text-muted-foreground">
            Tutorials
          </p>
          <h1 className="mt-3 font-display text-4xl font-bold tracking-tight md:text-5xl">
            Learn Grit on video
          </h1>
          <p className="mt-4 max-w-2xl text-lg leading-relaxed text-muted-foreground">
            Installing the CLI and scaffolding a project, generating a resource and its
            relationships, the admin panel you get without writing it, roles and
            permissions, deploying to a VPS, and building with an AI coding agent.
            Watch first if you would rather see the framework than read about it.
          </p>

          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              href="/docs/tutorials"
              className="inline-flex h-11 items-center gap-2 rounded-full bg-primary px-6 text-sm font-medium text-primary-foreground transition-opacity duration-150 hover:opacity-90"
            >
              Written tutorials
              <ArrowUpRight className="h-4 w-4" aria-hidden="true" />
            </Link>
            <Link
              href="/docs/start"
              className="inline-flex h-11 items-center gap-2 rounded-full border border-border px-6 text-sm font-medium text-foreground transition-colors duration-150 hover:border-accent/50"
            >
              Install Grit and start
            </Link>
          </div>
        </div>
      </section>

      <div className="mx-auto max-w-5xl px-6 py-16">
        <p className="mb-12 max-w-2xl text-[15px] leading-relaxed text-muted-foreground">
          Every video says which Grit version it was recorded against. Grit ships several
          releases a week, so a command in a video is a command as it was: where one
          differs from the written documentation, the documentation is what moves with
          the framework. Nothing is requested from YouTube until you press play, and the
          thumbnails are served from this site, so opening this page sets no third-party
          cookies and downloads no player.
        </p>

        <div className="space-y-16">
          {VIDEOS.map((video, i) => (
            <section key={video.videoId}>
              <div className="mb-4 flex flex-wrap items-center gap-2">
                <span className="inline-flex items-center rounded-md bg-accent/12 px-2 py-0.5 text-xs font-semibold text-accent">
                  {video.level}
                </span>
                <span className="inline-flex items-center gap-1.5 text-xs text-muted-foreground">
                  <Clock className="h-3.5 w-3.5" aria-hidden="true" />
                  {video.length}
                </span>
                <span className="inline-flex items-center rounded-md border border-border px-2 py-0.5 font-mono text-[11px] text-muted-foreground">
                  recorded on {video.recordedOn}
                </span>
              </div>

              <h2 className="mb-2 font-display text-2xl font-semibold tracking-tight">
                {video.title}
              </h2>
              <p className="mb-6 max-w-[62ch] text-[15px] leading-relaxed text-muted-foreground">
                {video.blurb}
              </p>

              <VideoFacade
                videoId={video.videoId}
                title={video.title}
                poster={video.poster}
                priority={i === 0}
              />
            </section>
          ))}
        </div>

        <div className="mt-16 rounded-xl border border-border bg-surface p-6">
          <h2 className="mb-2 flex items-center gap-2 font-display text-lg font-semibold">
            <Youtube className="h-5 w-5 text-accent" aria-hidden="true" />
            More of these are coming
          </h2>
          <p className="mb-5 max-w-[62ch] text-[15px] leading-relaxed text-muted-foreground">
            The series in progress covers your first project, resources and
            relationships, roles and permissions, deploying with{' '}
            <code className="rounded bg-code-bg px-1.5 py-0.5 font-mono text-[13px]">
              grit deploy
            </code>
            , and using Grit with an AI coding agent. If there is one you want next, or
            one you would record yourself, say so on the issue that tracks them.
          </p>
          <a
            href="https://github.com/MUKE-coder/grit/issues/309"
            target="_blank"
            rel="noreferrer"
            className="inline-flex h-10 items-center gap-2 rounded-full border border-border px-5 text-sm font-medium text-foreground transition-colors duration-150 hover:border-accent/50"
          >
            The video series issue
            <ArrowUpRight className="h-4 w-4" aria-hidden="true" />
          </a>
        </div>
      </div>
    </div>
  )
}
