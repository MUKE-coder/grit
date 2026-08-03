import Link from 'next/link'
import { ArrowRight, Check, Clock } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/start')

/**
 * The guided path — a numbered spine over docs that already exist.
 *
 * Grit's problem was never missing documentation. Four independent reviews of
 * the site failed to find the MCP server, generated tests, the audit log,
 * webhooks, multi-tenancy and the deployment guides — all shipped, all
 * documented, all reachable only if you already knew to look.
 *
 * A reference organised by subsystem answers "where is X" for someone who knows
 * X exists. It answers nothing for someone who does not know what to ask. This
 * page is the other thing: one ordered route from nothing to deployed, where
 * each step says what you will have when it is done.
 */

interface Step {
  n: number
  title: string
  href: string
  time: string
  body: string
  outcome: string
}

const STEPS: Step[] = [
  {
    n: 1,
    title: 'Create a project',
    href: '/docs/getting-started/create-a-project',
    time: '5 min',
    body:
      'One command scaffolds the API, the frontend and the admin panel, already wired together. Pick the architecture when it asks — if you are unsure, take the recommendation below.',
    outcome: 'A running app at localhost, with auth already working.',
  },
  {
    n: 2,
    title: 'Understand the layout',
    href: '/docs/getting-started/project-structure',
    time: '10 min',
    body:
      'Where models, services, handlers and resource definitions live, and — more usefully — which files the generators own and which are yours permanently.',
    outcome: 'You can find anything, and you know what is safe to edit.',
  },
  {
    n: 3,
    title: 'Generate your first resource',
    href: '/docs/concepts/code-generation',
    time: '10 min',
    body:
      'The command that justifies the framework. One line produces the model, migration, service, handler, routes, Zod schema, TypeScript types, React Query hooks and a working admin page.',
    outcome: 'Full CRUD, backend to admin UI, from one command.',
  },
  {
    n: 4,
    title: 'Model the relationships',
    href: '/docs/concepts/field-types',
    time: '15 min',
    body:
      'Every field type, and the two that matter most: belongs_to and many_to_many. This is where a toy CRUD app becomes something with a real schema behind it.',
    outcome: 'Related records, with the admin forms and filters to match.',
  },
  {
    n: 5,
    title: 'Lock it down',
    href: '/docs/backend/authentication',
    time: '15 min',
    body:
      'Auth ships working, so this is about the parts you configure: roles, permissions, who can reach which route, and what the frontend is allowed to show.',
    outcome: 'Routes enforced server-side, with the UI mirroring the rules.',
  },
  {
    n: 6,
    title: 'Add the batteries you need',
    href: '/docs/batteries',
    time: 'as needed',
    body:
      'File uploads, background jobs, email, realtime, AI, webhooks, multi-tenancy, feature flags. All present, all off until you switch them on. Read this page once so you know what is there before you build it yourself.',
    outcome: 'You stop writing things the framework already ships.',
  },
  {
    n: 7,
    title: 'Deploy it',
    href: '/docs/deployment',
    time: '30 min',
    body:
      'Pick a host, set the environment, run the migrations, work the checklist. Railway or Render if you have no opinion; a $5 VPS if cost matters.',
    outcome: 'A URL other people can use.',
  },
]

export default function StartPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Start here</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                The path from nothing to deployed
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Seven steps, about ninety minutes end to end. Everything else in these docs is
                reference you can reach for later — this is the part to read in order.
              </p>
            </div>

            {/* The golden path. Grit supports five architectures and three
                frontends; that flexibility is a feature and a liability, because
                a newcomer's first decision should not be a matrix. */}
            <div className="rounded-xl border border-primary/25 bg-primary/[0.04] p-6 mb-12">
              <h2 className="font-semibold mb-3">If you do not have a preference, use this</h2>
              <CodeBlock language="bash" code="grit new my-app --triple --next" />
              <p className="text-sm text-muted-foreground leading-relaxed mt-4">
                Go API, Next.js web app, admin panel, Postgres, Redis. It is the combination
                that gets the most use and the most testing, and every guide here assumes it
                unless it says otherwise.
              </p>
              <p className="text-sm text-muted-foreground leading-relaxed mt-3">
                The other four architectures are real and supported —{' '}
                <Link href="/docs/getting-started/create-a-project" className="text-primary hover:underline">
                  single binary, double, API-only and mobile
                </Link>{' '}
                — but choose one because you have a reason, not because you were made to
                choose on day one.
              </p>
            </div>

            <div className="space-y-4 mb-14">
              {STEPS.map((step, i) => (
                <div key={step.n} className="relative">
                  <Link
                    href={step.href}
                    className="group block rounded-xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border hover:bg-card/70"
                  >
                    <div className="flex items-start gap-4">
                      <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-primary/10 text-sm font-bold text-primary">
                        {step.n}
                      </span>
                      <div className="min-w-0 flex-1">
                        <div className="flex flex-wrap items-baseline gap-x-3 gap-y-1">
                          <h3 className="font-semibold group-hover:text-primary transition-colors">
                            {step.title}
                          </h3>
                          <span className="inline-flex items-center gap-1 text-xs text-muted-foreground">
                            <Clock className="h-3 w-3" />
                            {step.time}
                          </span>
                        </div>
                        <p className="mt-2 text-sm text-muted-foreground leading-relaxed">
                          {step.body}
                        </p>
                        {/* Naming the outcome is the difference between a list of
                            links and a path. It tells you whether you are done. */}
                        <p className="mt-3 flex items-start gap-2 text-[13px] text-emerald-600 dark:text-emerald-400">
                          <Check className="h-3.5 w-3.5 shrink-0 mt-0.5" />
                          {step.outcome}
                        </p>
                      </div>
                      <ArrowRight className="h-4 w-4 shrink-0 text-muted-foreground transition-transform group-hover:translate-x-0.5 mt-1" />
                    </div>
                  </Link>
                  {i < STEPS.length - 1 && (
                    <div aria-hidden className="ml-9 h-4 w-px bg-border/60" />
                  )}
                </div>
              ))}
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Then build something real</h2>
            <p className="text-muted-foreground leading-relaxed mb-6">
              Once the path above makes sense, the tutorials build complete applications rather
              than explaining features — a blog, a product catalogue, an ecommerce store, a SaaS
              with billing and tenants.
            </p>
            <div className="grid gap-3 sm:grid-cols-2 mb-12">
              {[
                { href: '/docs/tutorials', title: 'All tutorials', body: 'Complete applications, start to finish.' },
                {
                  href: '/docs/getting-started/existing-projects',
                  title: 'Already have an app?',
                  body: 'What works standalone in a codebase Grit did not scaffold.',
                },
                {
                  href: '/docs/getting-started/coming-from',
                  title: 'Coming from Laravel, Django or Next?',
                  body: 'The translation table for concepts you already know.',
                },
                { href: '/docs/ai-integration', title: 'Building with an AI agent?', body: 'The generated prompt and the MCP server.' },
              ].map((c) => (
                <Link
                  key={c.href}
                  href={c.href}
                  className="group rounded-xl border border-border/50 bg-card/40 p-5 transition-colors hover:border-border hover:bg-card/70"
                >
                  <div className="flex items-center gap-1.5 font-semibold text-sm">
                    {c.title}
                    <ArrowRight className="h-3.5 w-3.5 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
                  </div>
                  <p className="mt-1.5 text-[13px] text-muted-foreground leading-relaxed">{c.body}</p>
                </Link>
              ))}
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
