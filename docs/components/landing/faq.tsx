import Link from 'next/link'
import { ChevronDown } from 'lucide-react'
import { DISCORD_INVITE_URL } from '@/components/community-cta'

/*
 * Plain <details>. Not a scripted accordion.
 *
 * It opens before JavaScript loads, it is findable with the browser's own
 * in-page search, and it prints. An accordion that re-implements all three in
 * React is a larger bundle that does less.
 */

export const FAQS: { q: string; a: string }[] = [
  {
    q: 'Do I need to know Go to use Grit?',
    a: "Basic Go knowledge helps, but Grit generates most of the code for you. The generated code follows clear patterns (handler → service → model) that are easy to extend. If you know any backend language, you'll pick it up fast.",
  },
  {
    q: 'Can I use Grit with an existing project?',
    a: 'Grit is designed for greenfield projects. It scaffolds the full project structure. However, you can use grit generate resource in existing Grit projects to add new features incrementally.',
  },
  {
    q: 'Is Grit production-ready?',
    a: 'Yes. Every scaffolded project includes JWT auth, RBAC, rate limiting (Sentinel), observability (Pulse), error handling, CORS, gzip compression, connection pooling, and graceful shutdown. It is designed for production from day one.',
  },
  {
    q: "What's the difference between Single and Triple architecture?",
    a: 'Single embeds the React SPA into the Go binary via go:embed, leaving one file to deploy. Triple is a Turborepo monorepo with separate web app, admin panel, and API, and is the better fit for teams and complex products.',
  },
  {
    q: 'Can I switch from Next.js to TanStack Router later?',
    a: 'The backend (Go API) is identical regardless of frontend choice. You would need to rebuild the frontend pages, but all hooks, types, and API patterns are the same. The admin panel components are also framework-agnostic React.',
  },
  {
    q: 'How does grit deploy work? Is it like Vercel?',
    a: 'grit deploy is for self-hosted deployments. It SSHs to your server, uploads the binary, configures systemd, and sets up Caddy with auto-TLS. For Vercel or Railway, push to git and the Dockerfile is included.',
  },
  {
    q: 'Is Grit open source?',
    a: 'Yes, Grit is fully open source under the MIT license. The CLI, all plugins, and the documentation are on GitHub.',
  },
]

export function Faq() {
  return (
    <>
      <div className="mt-12 space-y-2">
        {FAQS.map((faq) => (
          <details
            key={faq.q}
            className="group overflow-hidden rounded-xl border border-border bg-surface-1 transition-colors hover:border-border-strong"
          >
            <summary className="flex cursor-pointer list-none items-center justify-between gap-4 px-6 py-4 text-[15px] font-medium text-fg">
              {faq.q}
              <ChevronDown
                className="h-4 w-4 shrink-0 text-fg-subtle transition-transform duration-200 group-open:rotate-180"
                aria-hidden="true"
              />
            </summary>
            <p className="px-6 pb-5 text-[14px] leading-[1.7] text-fg-muted">{faq.a}</p>
          </details>
        ))}
      </div>

      <p className="mt-8 text-center text-[14px] text-fg-muted">
        Anything else?{' '}
        <Link
          href={DISCORD_INVITE_URL}
          target="_blank"
          rel="noreferrer"
          className="g-link"
        >
          Ask in Discord
        </Link>
        . No question is too small.
      </p>
    </>
  )
}
