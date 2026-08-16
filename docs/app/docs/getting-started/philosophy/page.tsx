import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CommunityCTA } from '@/components/community-cta'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/philosophy')

export default function PhilosophyPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Philosophy</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Why Grit exists, what it borrows from the frameworks that came before it, and
                what each of its decisions costs you. Every choice below has a downside; they
                are named rather than hidden.
              </p>
            </div>

            <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-8">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Related reading.</strong> The{' '}
                <Link href="/pitch" className="text-primary hover:underline">
                  pitch
                </Link>{' '}
                states the six trade-offs in short form. The long-form founder story is on the
                blog:{' '}
                <Link href="/blog/why-i-built-grit" className="text-primary hover:underline">
                  Why I built Grit &rarr;
                </Link>
              </p>
            </div>

            {/* The one-line answer to "what exactly is Grit?". Everything on
                this page and in the framework fits under one of these three. */}
            <div className="mb-10">
              <h2 className="text-2xl font-semibold tracking-tight mb-3">What Grit is</h2>
              <p className="text-muted-foreground leading-relaxed mb-6">
                Grit does enough things that the honest answer needs three parts rather
                than one. Everything in the framework fits under one of them, and anything
                that does not fit under any of them probably should not be in the core.
              </p>
              <div className="grid gap-4 sm:grid-cols-3">
                <div className="rounded-lg border border-border/40 bg-card/40 p-5">
                  <div className="text-sm font-semibold text-foreground mb-2">
                    Application Infrastructure
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Auth, RBAC, sessions, jobs, cron, cache, storage, mail, realtime,
                    observability, backups, security. The three weeks you spend before you
                    have written any of your own product.
                  </p>
                </div>
                <div className="rounded-lg border border-border/40 bg-card/40 p-5">
                  <div className="text-sm font-semibold text-foreground mb-2">
                    Application UI
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    An admin panel with tables, forms, filters, bulk actions and widgets,
                    plus a component library. Not a theme: components that understand what a
                    resource is.
                  </p>
                </div>
                <div className="rounded-lg border border-border/40 bg-card/40 p-5">
                  <div className="text-sm font-semibold text-foreground mb-2">
                    Application Generation
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    Describe a resource once and get the Go model, service, API, migrations,
                    OpenAPI, Zod schemas, TypeScript types, React hooks and admin screens.
                    Code you own, not runtime magic.
                  </p>
                </div>
              </div>
              <p className="text-muted-foreground leading-relaxed mt-6">
                The third is what connects the first two. A framework with good
                infrastructure and a good admin panel still leaves you wiring one to the
                other by hand for every entity in your domain. Grit&apos;s bet is that the
                wiring is the boilerplate worth eliminating, and that eliminating it is only
                worth anything if the result is code you can read, edit and take over.
              </p>
            </div>

            <div className="prose-grit">
              {/* ============================================================ */}
              <h2>Why Grit exists</h2>
              <p>
                Building a full-stack application today means stitching together fifteen or more
                tools: a framework, an ORM, an auth library, a validation library, a data-fetching
                library, an S3 client, a mailer, a queue, an admin panel. Every project starts with
                days of wiring. Every developer wires it differently, and the result is a codebase
                where the interesting decisions are buried under the boring ones.
              </p>
              <p>
                Go has an excellent performance and deployment story, but its web ecosystem is
                deliberately unopinionated &mdash; Gin, Echo, Fiber and Chi for routing; GORM, sqlc
                and sqlx for data; dozens of auth libraries. That is a genuine strength for
                infrastructure work and a genuine cost for product work, where you would rather
                inherit good defaults than assemble them.
              </p>
              <p>
                And then there is the admin panel. Laravel has <strong>Filament</strong>: define a
                resource, get a real dashboard with tables, filters, forms and widgets. The Go +
                React world has no equivalent, so every team rebuilds the same CRUD screens.
              </p>
              <p>
                <strong>Grit exists to close that gap</strong> &mdash; convention over
                configuration, batteries included, and code generation, applied to Go on the
                backend and React on the front.
              </p>

              {/* ============================================================ */}
              <h2>What Grit borrows</h2>
              <p>
                Almost nothing here is a new idea. Grit&apos;s contribution is the combination, not
                the invention:
              </p>
              <ul>
                <li>
                  <strong>Laravel + Filament (PHP)</strong> &mdash; the Artisan CLI and
                  resource-driven admin panel are the direct model for{' '}
                  <code>grit generate</code> and Grit&apos;s admin.
                </li>
                <li>
                  <strong>Ruby on Rails</strong> &mdash; convention over configuration, and the
                  argument that frameworks <em>should</em> have opinions.
                </li>
                <li>
                  <strong>Django (Python)</strong> &mdash; batteries included. Auth, admin, ORM and
                  email belong in the box.
                </li>
                <li>
                  <strong>Next.js</strong> &mdash; file-based routing and the React developer
                  experience Grit builds its frontends on.
                </li>
                <li>
                  <strong>GORM Studio</strong> &mdash; our own database browser, folded into the
                  framework. Seeing your data without leaving the browser changes how you work.
                </li>
              </ul>

              {/* ============================================================ */}
              <h2>Design principles, and what each costs</h2>
              <p>
                A principle that costs nothing is a slogan. Each of these buys something real and
                charges for it.
              </p>

              <h3>1. Convention over configuration</h3>
              <p>
                One auth system. One folder structure. One naming convention. Any developer can
                open any Grit project and know where things live, and an AI agent does not have to
                infer which of six patterns you chose.
              </p>
              <p>
                <strong>The cost:</strong> if your application genuinely does not fit the
                convention, you are swimming upstream, and the framework will not help you. Grit
                is a good fit for CRUD-shaped products with an admin surface. It is a poor fit for
                software whose core is an unusual data model or an unusual request lifecycle.
              </p>

              <h3>2. Code generation over runtime magic</h3>
              <p>
                <code>grit generate resource Post</code> writes real files &mdash; Go model,
                service, handler, routes, Zod schema, TypeScript types, React Query hook, admin
                page. No hidden proxies, no auto-wiring, no reflection. Open any file, read it,
                change it.
              </p>
              <p>
                <strong>The cost:</strong> generated code is yours to maintain from the moment it
                lands. Upgrading Grit does not retroactively improve resources you generated six
                months ago &mdash; new templates apply to new code. That is the honest trade for
                never being blocked by a framework internal you cannot see.
              </p>

              <h3>3. Own your code</h3>
              <p>
                Everything Grit produces is in your repository, not in <code>node_modules</code>{' '}
                and not compiled into a binary you cannot inspect. Stop using the CLI tomorrow and
                the project still builds: it is a Go API, a Next.js app, and some TypeScript.
              </p>
              <p>
                <strong>The cost:</strong> a bigger repository and more code with your name on it.
                Some teams would rather depend on a library than own the source. That is a
                reasonable preference, and it is the opposite of this one.
              </p>

              <h3>4. Batteries included, individually removable</h3>
              <p>
                Auth, storage, email, queues, cron, AI and the admin panel ship in every project.
                Each is a separate Go package that can be switched off with{' '}
                <code>MODULE_&lt;NAME&gt;=false</code> in <code>.env</code>, or deleted outright
                without breaking the rest of the app.
              </p>
              <p>
                <strong>The cost:</strong> more surface area to read on day one. A project that
                needs three of the batteries still starts larger than a hand-rolled service that
                needs three things.
              </p>

              <h3>5. Secure before it is convenient</h3>
              <p>
                CSRF, a strict CSP, rate limiting, SSRF defence, IDOR-safe ownership checks,
                field-level encryption, server-side sessions with revocation, GDPR export and
                erasure, and a tamper-evident audit log are in the scaffold, not in a checklist.
              </p>
              <p>
                <strong>The cost:</strong> some of it will be in your way. A strict CSP means you
                cannot paste an inline script and move on. That friction is the feature &mdash;
                security work is what loses every argument against a deadline, so it ships before
                the argument happens.
              </p>

              <h3>6. Predictable enough for machines</h3>
              <p>
                Because there is one way to do each thing, an AI assistant has nothing to guess.
                Grit ships a <code>SKILL.md</code>, and{' '}
                <Link href="/docs/ai-workflows/mcp">
                  <code>grit mcp serve</code>
                </Link>{' '}
                hands an agent the real route table and model definitions over the Model Context
                Protocol &mdash; parsed from your source, read-only.
              </p>
              <p>
                <strong>The cost:</strong> conventions were sometimes simplified past what an
                experienced developer would have chosen, because the question asked of every
                pattern was &ldquo;can this be generated and read back reliably?&rdquo; Where the
                answer was no, the clever version lost.
              </p>

              {/* ============================================================ */}
              <h2>Why Go &mdash; and when it is the wrong answer</h2>
              <p>
                Go is small, explicit, and boring in the way infrastructure should be. It compiles
                to a single static binary with no runtime to install, which collapses most of what
                deployment usually costs. Its concurrency model suits the work a product backend
                actually does &mdash; many concurrent requests, each mostly waiting on I/O. For
                CPU-bound work it is substantially faster than interpreted runtimes, though for the
                request/response work most applications do, the deployment story and the predictable
                memory profile matter more than raw benchmark numbers.
              </p>
              <p>
                <strong>Go is the wrong answer when:</strong> your core problem is data science or
                machine learning, where Python&apos;s ecosystem is not close to matched; your team
                is genuinely all-TypeScript and adding a second language costs more than the
                runtime saves; or you want an expressive, metaprogramming-heavy style, which Go
                will actively resist.
              </p>

              <h3>Why React &mdash; and when it is not</h3>
              <p>
                React has the largest component ecosystem of any frontend framework, the tools Grit
                depends on are built for it first, and more developers know it than any
                alternative, which matters when hiring. It also reaches mobile through Expo, so a
                Grit project can share types and validation between web and native.
              </p>
              <p>
                <strong>Pick something else when:</strong> you want the smallest possible bundle and
                the simplest mental model, where Svelte is a better answer; you are building
                server-rendered pages with sprinkles of interactivity, where HTMX plus Go templates
                is genuinely lighter; or your team already knows Vue well, in which case that
                knowledge outweighs the ecosystem gap.
              </p>

              <h3>Why not full-stack JavaScript</h3>
              <p>
                For CRMs, SaaS tools, internal dashboards, and anything with WebSockets, background
                jobs, or sustained processing, a long-running server beats a serverless function,
                and Go gives you one binary that handles high concurrency and deploys anywhere.
              </p>
              <p>
                But if your team is small, entirely TypeScript, and shipping speed dominates every
                other concern, a single-language stack is a real advantage.{' '}
                <Link href="https://adonisjs.com" target="_blank" rel="noreferrer">
                  AdonisJS
                </Link>{' '}
                and Nest solve a similar problem well inside that constraint. Grit is the answer
                when you want the Laravel experience and a Go backend, not when you want one
                language everywhere.
              </p>

              {/* ============================================================ */}
              <h2>What Grit deliberately does not do</h2>
              <ul>
                <li>
                  <strong>Support every database.</strong> Postgres in production, SQLite for
                  development and tests. Adding more would dilute the migration and Studio work
                  that makes those two good.
                </li>
                <li>
                  <strong>Abstract the cloud.</strong> Grit generates Docker, Compose, and
                  deployment config you can read. It does not wrap your provider in an interface
                  that hides what is actually happening.
                </li>
                <li>
                  <strong>Ship a plugin for everything.</strong> Plugins are added when someone
                  needs one, not in anticipation. An unused abstraction is a maintenance cost with
                  no user.
                </li>
                <li>
                  <strong>Chase releases.</strong> Every release runs a 73-check matrix across 13
                  project shapes and is re-verified by downloading the published binary and
                  building a fresh project with it. That is slower than shipping on green unit
                  tests, and it is the point.
                </li>
              </ul>

              <p>
                If those trades sound right for what you are building, the{' '}
                <Link href="/docs/getting-started/quick-start">quick start</Link> takes about five
                minutes. If they do not, that is genuinely useful to know now rather than later.
              </p>
            </div>

            <CommunityCTA className="mt-10" />

            {/* Next/Prev navigation */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs">Introduction</Link>
              </Button>
              <Button asChild className="glow-purple-sm ml-auto">
                <Link href="/docs/getting-started/quick-start">
                  Quick Start
                  <ArrowRight className="ml-1.5 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
