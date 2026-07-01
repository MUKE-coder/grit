import { Globe } from 'lucide-react'

// Publisher / author card shown at the end of blog posts and the blog index.
const AVATAR = 'https://14j7oh8kso.ufs.sh/f/HLxTbDBCDLwfAUUBxSZezIN7vwylkF1PXSCqAuseUG0gx8mh'
const LINKEDIN = 'https://www.linkedin.com/newsletters/the-daily-grit-7477930448136912896/'
const PORTFOLIO = 'https://jb.desishub.com/'
const X = 'https://x.com/MJohnbaptist'

const iconLink =
  'h-8 w-8 rounded-md border border-border/40 flex items-center justify-center text-muted-foreground hover:text-foreground hover:border-border/60 transition-colors'

export function AuthorCard() {
  return (
    <div className="rounded-2xl border border-border/50 bg-card/40 p-6 flex flex-col sm:flex-row items-start sm:items-center gap-5">
      {/* eslint-disable-next-line @next/next/no-img-element */}
      <img
        src={AVATAR}
        alt="JB — Creator of Grit"
        className="h-16 w-16 rounded-full object-cover ring-2 ring-primary/20 shrink-0"
      />

      <div className="flex-1 min-w-0">
        <div className="text-[11px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-1">
          Written &amp; published by
        </div>
        <div className="font-semibold text-foreground">JB — Creator of Grit</div>
        <p className="text-sm text-muted-foreground mt-1 leading-relaxed">
          Founder of Grit and author of <span className="text-foreground">The Daily Grit</span> — a
          5-minute morning read on building full-stack apps with Go + React.
        </p>
        <div className="mt-3 flex items-center gap-2">
          <a href={LINKEDIN} target="_blank" rel="noreferrer" aria-label="The Daily Grit on LinkedIn" className={iconLink}>
            <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
              <path d="M19 3a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h14zM8.339 18.337v-8.59H5.667v8.59h2.672zM7.003 8.574a1.548 1.548 0 1 0 0-3.096 1.548 1.548 0 0 0 0 3.096zm11.335 9.763V13.64c0-2.465-1.338-3.612-3.123-3.612-1.44 0-2.087.791-2.448 1.348v-1.157h-2.671c.034.751 0 8.59 0 8.59h2.671v-4.79c0-.243.018-.487.09-.66.196-.485.642-.989 1.39-.989.982 0 1.376.752 1.376 1.852v4.587h2.715z" />
            </svg>
          </a>
          <a href={PORTFOLIO} target="_blank" rel="noreferrer" aria-label="Portfolio" className={iconLink}>
            <Globe className="h-4 w-4" />
          </a>
          <a href={X} target="_blank" rel="noreferrer" aria-label="X (Twitter)" className={iconLink}>
            <svg className="h-3.5 w-3.5" viewBox="0 0 24 24" fill="currentColor">
              <path d="M18.244 2.25h3.308l-7.227 8.26 8.502 11.24H16.17l-5.214-6.817L4.99 21.75H1.68l7.73-8.835L1.254 2.25H8.08l4.713 6.231zm-1.161 17.52h1.833L7.084 4.126H5.117z" />
            </svg>
          </a>
        </div>
      </div>

      <a
        href={LINKEDIN}
        target="_blank"
        rel="noreferrer"
        className="shrink-0 inline-flex items-center justify-center h-10 px-5 rounded-full bg-primary text-primary-foreground font-semibold text-sm hover:bg-primary/90 transition-colors"
      >
        Subscribe
      </a>
    </div>
  )
}
