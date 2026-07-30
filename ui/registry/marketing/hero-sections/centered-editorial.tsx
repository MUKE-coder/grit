import { Download, GitBranch } from 'lucide-react'

/** The Grit UI mark, inlined so the block stays self-contained. */
function GritMark({ className = 'size-7' }: { className?: string }) {
  return (
    <svg viewBox="0 0 32 32" fill="none" aria-hidden="true" className={className}>
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M8 0H24A8 8 0 0 1 32 8V24A8 8 0 0 1 24 32H8A8 8 0 0 1 0 24V8A8 8 0 0 1 8 0ZM16 9.4A6.6 6.6 0 1 0 21.4 19.9V17.4H17.2A1.7 1.7 0 0 1 17.2 14H23.1A1.7 1.7 0 0 1 24.8 15.7V20.6A1.7 1.7 0 0 1 24.4 21.7A10 10 0 1 1 22.6 8.2A1.7 1.7 0 0 1 20.4 10.8A6.6 6.6 0 0 0 16 9.4Z"
        fill="currentColor"
      />
    </svg>
  )
}

const NAV = ['Product', 'Docs', 'Extensions', 'Pricing']

const PILLARS = [
  {
    title: 'Fast',
    body: 'A Go binary that starts in milliseconds and holds thousands of concurrent connections without breaking a sweat.',
  },
  {
    title: 'Generated',
    body: 'One resource definition writes the model, the API, the typed client and the admin screen — all of it yours to edit.',
  },
  {
    title: 'Complete',
    body: 'Auth, storage, queues, email and an audit log ship in the scaffold, not in a checklist for later.',
  },
]

export default function CenteredEditorial({
  announcement = 'Early access',
  announcementLink = 'Grit UI — 100 blocks for Go + React apps',
  title = 'Your last full-stack framework',
  subtitle = 'Grit is a Go and React framework crafted for speed and for building alongside AI.',
  primaryLabel = 'Download now',
  secondaryLabel = 'Clone source',
  availability = 'Available for macOS, Linux, and Windows',
}: {
  announcement?: string
  announcementLink?: string
  title?: string
  subtitle?: string
  primaryLabel?: string
  secondaryLabel?: string
  availability?: string
}) {
  return (
    <div className="bg-[#faf9f7] dark:bg-gray-950">
      {/* Nav */}
      <header className="border-b border-gray-900/[0.07] dark:border-white/10">
        <nav
          aria-label="Global"
          className="mx-auto flex max-w-6xl items-center gap-8 px-6 py-4"
        >
          <a href="#" className="flex items-center gap-2 text-gray-900 dark:text-white">
            <GritMark className="size-7" />
            <span className="text-base font-semibold tracking-tight">Grit</span>
          </a>
          <div className="hidden gap-6 md:flex">
            {NAV.map((item) => (
              <a
                key={item}
                href="#"
                className="text-sm text-gray-600 transition-colors hover:text-gray-900 dark:text-gray-400 dark:hover:text-white"
              >
                {item}
              </a>
            ))}
          </div>
          <div className="ml-auto flex items-center gap-3">
            <a
              href="#"
              className="hidden text-sm text-gray-600 transition-colors hover:text-gray-900 sm:block dark:text-gray-400 dark:hover:text-white"
            >
              Sign in
            </a>
            <a
              href="#"
              className="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-500"
            >
              Download
            </a>
          </div>
        </nav>
      </header>

      {/* Announcement strip */}
      <div className="border-b border-gray-900/[0.07] py-3 text-center text-sm dark:border-white/10">
        <span className="text-blue-700 dark:text-blue-400">{announcement}:</span>{' '}
        <a
          href="#"
          className="text-gray-700 underline-offset-4 hover:underline dark:text-gray-300"
        >
          {announcementLink} <span aria-hidden="true">&rarr;</span>
        </a>
      </div>

      {/* Hero */}
      <div className="mx-auto max-w-3xl px-6 py-24 text-center sm:py-32">
        {/* The serif italic display face is the whole personality of this
            layout — it is what stops a centred hero reading as generic. */}
        <h1 className="font-serif text-5xl font-normal italic tracking-tight text-blue-700 sm:text-6xl dark:text-blue-400">
          {title}
        </h1>
        <p className="mx-auto mt-7 max-w-lg text-lg text-gray-600 dark:text-gray-400">
          {subtitle}
        </p>

        <div className="mt-10 flex flex-wrap items-center justify-center gap-3">
          <a
            href="#"
            className="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-blue-500"
          >
            <Download aria-hidden="true" className="size-4" />
            {primaryLabel}
            <kbd className="ml-1 rounded bg-white/20 px-1.5 py-0.5 font-mono text-[11px]">
              D
            </kbd>
          </a>
          <a
            href="#"
            className="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-5 py-2.5 text-sm font-semibold text-gray-900 transition-colors hover:bg-gray-50 dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10"
          >
            <GitBranch aria-hidden="true" className="size-4" />
            {secondaryLabel}
            <kbd className="ml-1 rounded bg-gray-900/10 px-1.5 py-0.5 font-mono text-[11px] dark:bg-white/10">
              C
            </kbd>
          </a>
        </div>

        <p className="mt-6 font-mono text-xs text-gray-500 dark:text-gray-500">
          {availability}
        </p>
      </div>

      {/* Pillars */}
      <div className="border-t border-gray-900/[0.07] dark:border-white/10">
        <div className="mx-auto grid max-w-6xl grid-cols-1 divide-y divide-gray-900/[0.07] md:grid-cols-3 md:divide-x md:divide-y-0 dark:divide-white/10">
          {PILLARS.map((pillar) => (
            <div key={pillar.title} className="px-6 py-8">
              <h2 className="text-base font-semibold text-gray-900 dark:text-white">
                {pillar.title}
              </h2>
              <p className="mt-2 text-sm/6 text-gray-600 dark:text-gray-400">
                {pillar.body}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
