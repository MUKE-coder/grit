import Link from 'next/link'
import { AlertCircle, ArrowRight, Check, X } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/existing-projects')

export default function ExistingProjectsPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Using Grit with an existing project
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit scaffolds new projects. It does not convert old ones. But several pieces
                of it work standalone in a codebase that has never heard of Grit — and those
                are usually the pieces people actually want.
              </p>
            </div>

            {/* Say the disappointing thing immediately. Someone arriving here has
                a specific question and burying the answer wastes their time. */}
            <div className="rounded-xl border border-amber-500/30 bg-amber-500/[0.05] p-6 mb-12">
              <div className="flex gap-3">
                <AlertCircle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                <div>
                  <h2 className="font-semibold mb-2">There is no `grit init` for an existing app</h2>
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    No command reads your current codebase and Grit-ifies it. The generators
                    write into a known project layout with known marker comments, and neither
                    exists in a codebase that was not scaffolded by Grit. A command that
                    pretended otherwise would produce files that do not compile against your
                    app, which is a worse outcome than not having the command.
                  </p>
                </div>
              </div>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-6">
              What works standalone, today
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-6">
              In rough order of how easy they are to adopt. None of these requires the Grit
              CLI to be anywhere near your production build.
            </p>

            <div className="space-y-4 mb-14">
              {[
                {
                  n: 1,
                  title: 'Grit UI blocks — any React project',
                  body:
                    'The block registry is a plain shadcn registry. Marketing sections, headers, feature sections and the swappable primitives install into any React + Tailwind codebase. Nothing about it is Grit-specific.',
                  code: 'npx shadcn@latest add https://ui.gritframework.dev/r/marketing-headers-simple-with-actions.json',
                },
                {
                  n: 2,
                  title: 'Pulse — any Gin app',
                  body:
                    'Observability as an ordinary Go package. Mount the middleware and you get request tracing, DB query timings and a dashboard, with no other part of Grit involved.',
                  code: `go get github.com/MUKE-coder/pulse

// in your existing router setup
r.Use(pulse.Middleware())
pulse.MountUI(r, "/pulse/ui")`,
                },
                {
                  n: 3,
                  title: 'Sentinel — any Gin app',
                  body:
                    'Rate limiting and WAF rules, same deal. Both of these are the parts of Grit most often wanted on their own, and both were built as standalone packages first.',
                  code: `go get github.com/MUKE-coder/sentinel

r.Use(sentinel.Guard(sentinel.Default()))`,
                },
                {
                  n: 4,
                  title: 'GORM Studio — any GORM app',
                  body:
                    'A visual browser for your existing models. Point it at the same *gorm.DB you already have.',
                  code: `go get github.com/MUKE-coder/gorm-studio

studio.Mount(r, db, "/studio")`,
                },
              ].map((item) => (
                <div key={item.n} className="rounded-xl border border-border/50 bg-card/40 p-5">
                  <div className="flex items-baseline gap-3 mb-2">
                    <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-primary/10 text-xs font-bold text-primary">
                      {item.n}
                    </span>
                    <h3 className="font-semibold text-sm">{item.title}</h3>
                  </div>
                  <p className="text-[13px] text-muted-foreground leading-relaxed mb-3 pl-9">
                    {item.body}
                  </p>
                  <div className="pl-9">
                    <CodeBlock language={item.n === 1 ? 'bash' : 'go'} code={item.code} />
                  </div>
                </div>
              ))}
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">
              Adding a Grit API beside what you have
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              If the goal is a new service rather than a rewrite, generate an API-only project
              into a subdirectory of the existing repo. It gets its own <code className="text-xs">go.mod</code>,
              its own Dockerfile and its own deploy, and it talks to your current system over
              HTTP like any other service.
            </p>
            <CodeBlock
              language="bash"
              code={`# From the root of your existing repo:
mkdir services/reporting && cd services/reporting
grit new . --api --here

# Then generate against it as normal:
grit generate resource Report --fields "title:string,period:string,total:float"`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4 mb-4">
              This is the honest brownfield path: <strong>strangle rather than convert</strong>.
              New surface area gets Grit&apos;s generators, the old system keeps working, and
              nothing has to be migrated on a deadline.
            </p>
            <div className="rounded-xl border border-border/50 bg-card/50 p-5 mb-14">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <code className="text-xs">--here</code> refuses to scaffold into a non-empty
                directory unless you also pass <code className="text-xs">--force</code>. That
                guard exists because the alternative is writing forty files over the top of
                someone&apos;s repository. If you pass <code className="text-xs">--force</code>,
                commit first.
              </p>
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">
              Borrowing the generated code
            </h2>
            <p className="text-muted-foreground leading-relaxed mb-6">
              Underrated and completely legitimate: scaffold a throwaway project, generate the
              resource you were about to hand-write, and read it. The output is ordinary Gin
              and GORM with no framework runtime — handler, service, model, migration, Zod
              schema, React Query hooks. Copy the parts you want.
            </p>
            <CodeBlock
              language="bash"
              code={`cd /tmp
grit new scratch --api
cd scratch
grit generate resource Invoice --fields "customer:belongs_to:Customer,total:float,status:select"

# Now read internal/{models,services,handlers}/invoice.go and take what is useful.`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4 mb-14">
              MIT licensed, and it is your code the moment it is generated. There is no
              attribution requirement and nothing to remove later.
            </p>

            <h2 className="text-2xl font-bold tracking-tight mb-4">What is genuinely not possible</h2>
            <div className="space-y-3 mb-12">
              {[
                'Pointing the admin panel at an existing database schema. Resource definitions are generated from Grit field types, not reflected from tables — you would be hand-writing the definitions anyway.',
                'Running `grit generate resource` inside a project Grit did not scaffold. The injection markers it needs are not there, and it will not invent them.',
                'Migrating an Express, Laravel or Django app. The generators emit Go. There is no path that does not involve rewriting the backend.',
              ].map((t) => (
                <div key={t} className="flex gap-3 rounded-xl border border-border/50 bg-card/40 p-4">
                  <X className="h-4 w-4 shrink-0 text-muted-foreground/60 mt-0.5" />
                  <p className="text-sm text-muted-foreground leading-relaxed">{t}</p>
                </div>
              ))}
            </div>

            <h2 className="text-2xl font-bold tracking-tight mb-4">Deciding</h2>
            <div className="rounded-xl border border-border/50 bg-card/50 p-6 mb-12">
              <ul className="space-y-3 text-sm text-muted-foreground leading-relaxed">
                <li className="flex gap-3">
                  <Check className="h-4 w-4 shrink-0 text-emerald-500 mt-0.5" />
                  <span>
                    <strong className="text-foreground">Want the admin panel and CRUD generation?</strong>{' '}
                    That needs a Grit-scaffolded project. Add one as a new service rather than
                    converting.
                  </span>
                </li>
                <li className="flex gap-3">
                  <Check className="h-4 w-4 shrink-0 text-emerald-500 mt-0.5" />
                  <span>
                    <strong className="text-foreground">Want observability, rate limiting or a DB browser?</strong>{' '}
                    Those are standalone packages. Ten minutes, no restructuring.
                  </span>
                </li>
                <li className="flex gap-3">
                  <Check className="h-4 w-4 shrink-0 text-emerald-500 mt-0.5" />
                  <span>
                    <strong className="text-foreground">Want the UI blocks?</strong> They work in
                    any React project already.
                  </span>
                </li>
              </ul>
            </div>

            <div className="flex flex-wrap items-center gap-6 border-t border-border/40 pt-8">
              <Link
                href="/docs/getting-started/quick-start"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Quick start
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
              <Link
                href="/docs/governance"
                className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:underline"
              >
                Governance &amp; risk
                <ArrowRight className="h-3.5 w-3.5" />
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
