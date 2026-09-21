import Link from 'next/link'
import type { Metadata } from 'next'
import {
  ArrowRight,
  Check,
  X,
  Sparkles,
  Bot,
  ShieldCheck,
  Layers,
  Database,
  Server,
  FileCode2,
  LayoutDashboard,
  Smartphone,
  Monitor,
  BookOpen,
  Sprout,
} from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { FadeIn, GSAPSection, GlowOrb } from '@/components/motion-primitives'
import { GridFrame } from '@/components/grid-frame'
import { CommunityCTA } from '@/components/community-cta'
import { GRIT_VERSION } from '@/config/site'

export const metadata: Metadata = {
  title: 'The Pitch',
  description:
    'Why Grit exists and what sets it apart: one business product across web, admin, mobile and desktop, generated from one definition, in Go. Where the database fits, how it compares, and when it is the wrong choice.',
  alternates: { canonical: 'https://gritframework.dev/pitch' },
}

/**
 * What one `grit generate resource` writes. Kept to what the generator really
 * emits: a claim here that the code does not back is worse than a shorter list.
 */
const SURFACES = [
  {
    icon: Server,
    title: 'Go API',
    body: 'Model, service and handler with pagination, search, filters, soft delete and optimistic locking, routes registered.',
  },
  {
    icon: Database,
    title: 'Database',
    body: 'The model registered for migration, indexes from the field definitions, and a seeder that writes realistic data in batches.',
  },
  {
    icon: FileCode2,
    title: 'Shared contract',
    body: 'Zod schemas and TypeScript types generated from the same fields, so the frontend and the API validate the same rules.',
  },
  {
    icon: Layers,
    title: 'Data hooks',
    body: 'React Query hooks for list, detail, create, update and delete, typed end to end.',
  },
  {
    icon: LayoutDashboard,
    title: 'Admin panel',
    body: 'A table with sorting, filters and bulk actions, a form with the right input for each field, Excel export and bulk import.',
  },
  {
    icon: Smartphone,
    title: 'Mobile app',
    body: 'Expo screens for the list, the detail and the form, talking to the same API with the same types.',
  },
  {
    icon: Monitor,
    title: 'Desktop app',
    body: 'Screens in the Wails desktop app, with the same relations, the same validation and the same API client.',
  },
  {
    icon: BookOpen,
    title: 'API reference',
    body: 'The endpoints added to the interactive API docs, so the contract is published the moment it exists.',
  },
] as const

/**
 * Every principle is written as a trade, "A over B", because a stated
 * preference that costs nothing is not a principle, it is a feature list. The
 * `gives` field names what you actually give up, so a reader can decide the
 * trade is wrong for them without having to discover it three weeks in.
 */
const PRINCIPLES = [
  {
    title: 'Generated over hidden',
    body: 'Grit writes real files into your repo: the model, service, handler, routes, hooks and screens. Not a runtime dependency that does it invisibly at boot.',
    gives:
      'You give up a small repo. You get code you can read, edit, delete, and debug with a stack trace that points at your own file. There is no fighting the framework, because there is no framework in the way: just Go and TypeScript you own.',
  },
  {
    title: 'Boring over clever',
    body: 'Gin, GORM, Next.js, Postgres, Redis. Every one has years of production behind it and its answers already written down somewhere.',
    gives:
      'You give up novelty. We are not going to invent a router or an ORM. The interesting part of your project should be your product, not the stack underneath it.',
  },
  {
    title: 'Opinionated over configurable',
    body: 'One way to do auth. One queue. One storage layer. Escape hatches where you need them, not a menu where you do not.',
    gives:
      'You give up the pleasure of choosing. You get a codebase where every Grit project looks the same, onboarding takes an afternoon, and an AI agent never has to guess which of six patterns you picked.',
  },
  {
    title: 'One definition over many copies',
    body: 'A resource is not a Go model or a TypeScript type. It is both, plus the handler, the Zod schema, the hooks, the admin page and the mobile and desktop screens, generated together from one definition.',
    gives:
      'You give up picking a different stack per client. You get surfaces that cannot drift apart, because one command writes all of them from the same fields.',
  },
  {
    title: 'Secure on day one over hardened on day ninety',
    body: 'CSRF, strict CSP, rate limiting, SSRF defence, IDOR-safe ownership checks, field-level encryption, server-side sessions and a tamper-evident audit log, in the scaffold.',
    gives:
      'You give up the option to skip it. That is the point: security is the work that never wins a prioritisation meeting against a feature, so it ships before the meeting happens.',
  },
  {
    title: 'Verified over asserted',
    body: 'Every release scaffolds fresh projects in several shapes, including mobile and desktop, then installs, builds, type-checks, lints and tests them, and upgrades a project from the previous release.',
    gives:
      'You give up a faster release cadence. Claims in a changelog are cheap; we would rather find the bug than describe the feature.',
  },
] as const

/** How Grit sits next to what a Go developer already knows. */
const COMPARISON = [
  {
    name: 'Go frameworks',
    examples: 'Gin, Echo, Fiber, Chi',
    gets: 'A router and your choice of ORM.',
    rest: 'Auth, admin, uploads, jobs, email, the frontend and every client are yours to build and keep in step.',
  },
  {
    name: 'Laravel, Rails',
    examples: 'batteries-included web frameworks',
    gets: 'The batteries, conventions and generators.',
    rest: 'For the web app. The mobile and desktop apps are separate projects with their own copy of the model.',
  },
  {
    name: 'Grit',
    examples: 'Go + React, one monorepo',
    gets: 'The batteries, across web, admin, mobile, desktop, jobs, realtime and deployment.',
    rest: 'From one definition, in Go. You build the product; the machinery is already there, on every client.',
  },
] as const

/** Where Grit fits best. The mirror image of NOT_FOR. */
const RIGHT_FOR = [
  'A business product: CRM, ERP, POS, school management, inventory, a marketplace, a booking system',
  'An admin panel the operations team lives in, next to the app customers use',
  'A mobile app or a desktop app that has to agree with the web app on every field',
  'Real data volumes: seeding a million rows to find the breaking point, importing years of records from Excel',
] as const

/** Where Grit is the wrong tool. Written plainly, because a reader who finds
 *  this out on their own three weeks in does not come back. */
const NOT_FOR = [
  {
    title: 'Your team does not want to write Go',
    body: 'Grit is Go on the backend, and that is not configurable. If you are TypeScript end to end, AdonisJS or Nest will make you happier, and we would rather you were happy than converted.',
  },
  {
    title: 'You want to pick your own ORM and router',
    body: 'You can swap them: it is your code. But you would be working against the generators, and that cost compounds with every resource you add.',
  },
  {
    title: 'You are building one small service with no UI',
    body: 'The batteries become weight you never use. `--api` trims a lot of it, but a plain Gin and GORM binary may genuinely be the simpler answer.',
  },
  {
    title: 'You need a decade-old ecosystem',
    body: 'Grit is young. Fewer Stack Overflow answers, fewer third-party plugins, fewer people who have hit your exact bug. What you get instead is a maintainer who ships weekly and answers questions directly.',
  },
] as const

// Small reusable section wrapper that reveals its children on scroll.
function PitchSection({
  children,
  className = '',
}: {
  children: React.ReactNode
  className?: string
}) {
  return (
    <section className={`relative px-6 ${className}`}>
      <div className="max-w-3xl mx-auto">
        <GSAPSection>{children}</GSAPSection>
      </div>
    </section>
  )
}

export default function PitchPage() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      {/* ═══ HERO ═══ */}
      <section className="relative overflow-hidden border-b border-border/40">
        <div
          className="absolute inset-0 -z-10"
          style={{
            background:
              'radial-gradient(ellipse 60% 60% at 50% -10%, hsl(var(--primary) / 0.14), transparent 60%)',
          }}
        />
        <GlowOrb className="-top-32 left-1/3 h-[420px] w-[420px] bg-primary/[0.10]" duration={20} />

        <span className="crosshair absolute top-28 left-[14%] text-foreground/20 hidden md:block" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute top-40 right-[16%] text-primary/30 hidden md:block" style={{ width: 16, height: 16 }} />

        <div className="relative max-w-3xl mx-auto px-6 pt-24 pb-20 md:pt-32 md:pb-24">
          <FadeIn>
            <span className="tag-mono text-primary mb-6 block">The Pitch</span>
          </FadeIn>
          <FadeIn delay={0.08}>
            <h1 className="font-display text-4xl md:text-6xl font-bold tracking-tight text-foreground leading-[1.05] mb-7">
              Build the product once.
              <br />
              <span className="bg-gradient-to-r from-primary via-sky-400 to-primary bg-clip-text text-transparent">
                Ship it to every screen.
              </span>
            </h1>
          </FadeIn>
          <FadeIn delay={0.16}>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed max-w-2xl">
              Grit is the framework for building one business product across{' '}
              <span className="text-foreground font-medium">web, admin, mobile and desktop</span>, with
              the API, background jobs, realtime, security and deployment underneath, without
              rebuilding the machinery for every client. You define a resource once, in Go, and it
              reaches every surface.
            </p>
          </FadeIn>
          <FadeIn delay={0.24}>
            <div className="mt-9 flex flex-col sm:flex-row items-start sm:items-center gap-3">
              <Link
                href="/docs/getting-started/quick-start"
                className="group inline-flex items-center gap-2 h-11 px-6 rounded-full bg-primary text-primary-foreground font-semibold text-sm hover:bg-primary/90 transition-colors glow-primary-sm"
              >
                Start building
                <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
              </Link>
              <Link
                href="#why"
                className="inline-flex items-center h-11 px-6 rounded-full border border-border/60 text-foreground font-medium text-sm hover:bg-accent/30 transition-colors"
              >
                Why Grit exists
              </Link>
            </div>
          </FadeIn>
        </div>
      </section>

      {/* ═══ WHY GRIT EXISTS ═══ */}
      <PitchSection className="py-20 md:py-24 border-b border-border/40">
        <div id="why" className="scroll-mt-24" />
        <p className="tag-mono text-primary mb-4">Why Grit exists</p>
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-6 leading-tight">
          Every product rebuilt the same machinery
        </h2>
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          Every business product we built meant building the same things again: auth and
          sessions, roles, an admin panel, file uploads, background jobs, email, audit logs,
          realtime. Then a mobile app, and sometimes a desktop app, talking to the same API.
        </p>
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          Each client needed the same model defined again: in Go, in TypeScript, in the
          validation, in the forms, in the tables. The copies drifted apart. A field added on the
          backend went missing on mobile; a rule enforced in the form was not enforced by the API.
          Most of the time went into plumbing, not the product.
        </p>
        <p className="text-lg text-foreground leading-relaxed font-medium">
          Grit exists so that work is done once, correctly, and reaches every client from a
          single definition.
        </p>
      </PitchSection>

      {/* ═══ THE PROBLEM ═══ */}
      <PitchSection className="py-20 md:py-24 border-b border-border/40">
        <p className="tag-mono text-muted-foreground mb-4">The machinery tax</p>
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-6 leading-tight">
          The first month is rarely about your product
        </h2>
        <p className="text-lg text-muted-foreground leading-relaxed mb-8">
          Before you write the feature that matters, you re-solve the problems you solved on the
          last project, and the one before that. None of it is hard. All of it is slow, and with
          every extra client you pay for it again.
        </p>
        <div className="grid sm:grid-cols-2 gap-3">
          {[
            'Auth: JWT, refresh, OAuth, 2FA, password reset',
            'An admin panel with tables, forms, and filters',
            'File uploads to S3 / R2 with presigned URLs',
            'Background jobs, queues, and a scheduler',
            'Email templates and a transactional sender',
            'Rate limiting, CORS, security headers, audit logs',
            'The same model again in the mobile app',
            'And again in the desktop app',
          ].map((item) => (
            <div
              key={item}
              className="flex items-start gap-3 rounded-xl border border-border/40 bg-card/40 px-4 py-3 text-sm text-foreground/80"
            >
              <X className="h-4 w-4 text-rose-400/80 mt-0.5 shrink-0" strokeWidth={2.5} />
              {item}
            </div>
          ))}
        </div>
        <p className="text-base text-muted-foreground leading-relaxed mt-8">
          That is months of work that produces <span className="text-foreground font-medium">zero</span>{' '}
          product differentiation. Grit treats all of it as solved: generated, wired, and hardened
          the moment you scaffold.
        </p>
      </PitchSection>

      {/* ═══ ONE DEFINITION, EVERY SURFACE ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40 overflow-hidden">
        <div className="absolute inset-0 -z-20 bg-grit-dots mask-fade-y opacity-60" />
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-10" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">What sets it apart</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight mb-5">
                One definition, every surface
              </h2>
              <p className="text-lg text-muted-foreground leading-relaxed">
                One command writes the resource everywhere it needs to exist, and every piece agrees
                with every other, because they all come from the same fields. Add a field later with{' '}
                <code className="text-primary bg-primary/5 px-1.5 py-0.5 rounded text-base">grit generate field</code>{' '}
                and it reaches the model, the validation, the types and the admin form and table in one
                step.
              </p>
            </div>

            <div data-gsap-reveal className="mb-8">
              <CodeBlock
                terminal
                filename="One command"
                className="!m-0"
                code={`grit generate resource Invoice \\
  --fields "number:string:unique,customer:belongs_to:Customer,total:money,due:date,status:select:draft=Draft|paid=Paid"`}
              />
            </div>

            <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
              {SURFACES.map((s) => (
                <div
                  key={s.title}
                  data-gsap-reveal
                  className="card-grit rounded-2xl border border-border/40 bg-card/40 p-5"
                >
                  <span className="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-primary/10 text-primary mb-3">
                    <s.icon className="h-[18px] w-[18px]" />
                  </span>
                  <h3 className="font-semibold text-foreground text-[15px] mb-1.5">{s.title}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{s.body}</p>
                </div>
              ))}
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ WHERE THE DATABASE FITS ═══ */}
      <PitchSection className="py-20 md:py-24 border-b border-border/40">
        <p className="tag-mono text-primary mb-4">Where the database fits</p>
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-6 leading-tight">
          Not a new ORM. What surrounds it.
        </h2>
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          Underneath, Grit uses GORM, like much of the Go world, and Postgres, MySQL or SQLite. We
          did not reinvent the data layer, and the database is not the headline. What Grit adds is
          everything that usually sits around it by hand:
        </p>
        <div className="space-y-3 mb-8">
          {[
            {
              title: 'Field types that mean something from the column to the form.',
              body: 'money is stored as exact minor units, tel is checked as a real phone number, email, date, country, select and encrypted fields are validated the same way in Go, in Zod and in the input the form renders.',
            },
            {
              title: 'Constraints that answer properly.',
              body: 'A duplicate in a unique field comes back as 409 naming the field, a stale edit as a version conflict, a broken rule as 422 with its message. Not a 500 the form cannot explain.',
            },
            {
              title: 'Seeding built for real volumes.',
              body: 'Most frameworks assume a seeder writes a few rows. Real work means seeding a million to find where your infrastructure breaks, or loading years of records buried in Excel files. Grit generates batched seeders with realistic data, and bulk import in the admin.',
            },
          ].map((item) => (
            <div key={item.title} className="flex items-start gap-3">
              <Sprout className="h-4 w-4 text-primary mt-1 shrink-0" />
              <p className="text-base text-muted-foreground leading-relaxed">
                <span className="text-foreground font-medium">{item.title}</span> {item.body}
              </p>
            </div>
          ))}
        </div>
        <p className="text-base text-muted-foreground leading-relaxed">
          The data layer is one example of the pattern, not the point of it. The point is that none
          of this is yours to rebuild, on any client.
        </p>
      </PitchSection>

      {/* ═══ HOW IT COMPARES ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-10" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">How it compares</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight">
                Batteries included, on every client
              </h2>
            </div>
            <div className="grid md:grid-cols-3 gap-4">
              {COMPARISON.map((c) => {
                const grit = c.name === 'Grit'
                return (
                  <div
                    key={c.name}
                    data-gsap-reveal
                    className={`rounded-2xl border p-6 flex flex-col ${
                      grit ? 'border-primary/30 bg-card/60 glow-primary-sm' : 'border-border/40 bg-card/30'
                    }`}
                  >
                    <h3 className={`font-semibold text-[17px] ${grit ? 'text-primary' : 'text-foreground'}`}>
                      {c.name}
                    </h3>
                    <p className="text-xs font-mono text-muted-foreground mb-4">{c.examples}</p>
                    <p className="text-sm text-foreground/80 leading-relaxed mb-3 flex items-start gap-2">
                      <Check className="h-4 w-4 text-emerald-400 mt-0.5 shrink-0" strokeWidth={2.5} />
                      {c.gets}
                    </p>
                    <p className="text-sm text-muted-foreground leading-relaxed flex items-start gap-2">
                      {grit ? (
                        <Check className="h-4 w-4 text-emerald-400 mt-0.5 shrink-0" strokeWidth={2.5} />
                      ) : (
                        <X className="h-4 w-4 text-rose-400/80 mt-0.5 shrink-0" strokeWidth={2.5} />
                      )}
                      {c.rest}
                    </p>
                  </div>
                )
              })}
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ BEFORE / AFTER ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40 overflow-hidden">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="text-center mb-12" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">Compress the work</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight">
                The same outcome, two timelines
              </h2>
            </div>
            <div className="grid md:grid-cols-2 gap-5 items-stretch">
              <div data-gsap-reveal className="rounded-2xl border border-border/50 bg-card/40 p-6 flex flex-col">
                <div className="flex items-center gap-2 mb-4">
                  <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-rose-500/15 text-rose-400">
                    <X className="h-3.5 w-3.5" strokeWidth={2.5} />
                  </span>
                  <span className="font-semibold text-foreground text-sm">Without a framework</span>
                </div>
                <ul className="space-y-2.5 text-sm text-muted-foreground flex-1">
                  {[
                    'Pick a router, ORM, auth lib, queue, mailer',
                    'Glue them together; debug the seams',
                    'Hand-roll an admin UI per resource',
                    'Wire types between Go and TypeScript by hand',
                    'Define the model again for mobile and desktop',
                    'Stand up Docker, CI, and a deploy story',
                  ].map((l) => (
                    <li key={l} className="flex items-start gap-2.5">
                      <span className="mt-2 h-1 w-1 rounded-full bg-rose-400/60 shrink-0" />
                      {l}
                    </li>
                  ))}
                </ul>
                <div className="mt-5 pt-4 border-t border-border/40 flex items-center justify-between">
                  <span className="text-xs text-muted-foreground font-mono">Time to first feature</span>
                  <span className="text-sm font-semibold text-rose-400">~3–6 weeks</span>
                </div>
              </div>

              <div
                data-gsap-reveal
                className="relative rounded-2xl border border-primary/30 bg-card/60 p-6 flex flex-col glow-primary-sm"
              >
                <div className="flex items-center gap-2 mb-4">
                  <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-emerald-500/15 text-emerald-400">
                    <Check className="h-3.5 w-3.5" strokeWidth={2.5} />
                  </span>
                  <span className="font-semibold text-foreground text-sm">With Grit</span>
                </div>
                <CodeBlock
                  terminal
                  filename="Two commands"
                  className="!m-0 flex-1"
                  code={`grit new my-app --triple --desktop
grit generate resource Product \\
  --fields "name:string,price:money,stock:int"`}
                />
                <div className="mt-5 pt-4 border-t border-border/40 flex items-center justify-between">
                  <span className="text-xs text-muted-foreground font-mono">Time to first feature</span>
                  <span className="text-sm font-semibold text-emerald-400">~5 minutes</span>
                </div>
              </div>
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ WHAT GRIT STANDS FOR ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-12" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">What Grit stands for</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight mb-5">
                Six trades, stated plainly
              </h2>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every one of these costs something. A preference that costs nothing is not a
                principle, it is a feature list, so each one below names what you give up.
              </p>
            </div>

            <div className="grid md:grid-cols-2 gap-4">
              {PRINCIPLES.map((p, i) => (
                <div
                  key={p.title}
                  data-gsap-reveal
                  className="card-grit rounded-2xl border border-border/40 bg-card/40 p-6 flex flex-col"
                >
                  <div className="flex items-baseline gap-2.5 mb-3">
                    <span className="font-mono text-sm font-semibold text-primary">{i + 1}.</span>
                    <h3 className="font-semibold text-foreground text-[17px] leading-snug">{p.title}</h3>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-4">{p.body}</p>
                  <p className="text-sm leading-relaxed text-foreground/70 border-l-2 border-primary/30 pl-4 mt-auto">
                    {p.gives}
                  </p>
                </div>
              ))}
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ THREE PILLARS ═══ */}
      <PitchSection className="py-20 md:py-24 border-b border-border/40 max-w-5xl">
        <div className="max-w-5xl mx-auto">
          <div className="text-center mb-12" data-gsap-reveal>
            <p className="tag-mono text-primary mb-4">What the trades buy you</p>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight">
              The compounding return on a narrow set of opinions
            </h2>
          </div>
          <div className="grid md:grid-cols-3 gap-4">
            {[
              {
                icon: Bot,
                title: 'Agents work, not guess',
                body: 'Because there is one way to do each thing, an AI assistant has nothing to guess. Grit ships a SKILL.md, and grit mcp serve hands an agent the real route table and model definitions over MCP, parsed from your source, read-only.',
              },
              {
                icon: ShieldCheck,
                title: 'Security you did not schedule',
                body: 'The hardening in every scaffold is work no one budgets for: SSRF defence, IDOR-safe ownership checks, GDPR export and erasure, SSO with SAML, field-level encryption. Already there, already tested.',
              },
              {
                icon: Layers,
                title: 'One set of patterns, every shape',
                body: 'Embed a SPA in the Go binary, split web, admin and API into a monorepo, go API-only, or add mobile and desktop. The generators and conventions do not change; only the shape does.',
              },
            ].map((p) => (
              <div key={p.title} className="card-grit rounded-2xl border border-border/40 bg-card/40 p-6">
                <span className="inline-flex h-10 w-10 items-center justify-center rounded-xl bg-primary/10 text-primary mb-4">
                  <p.icon className="h-5 w-5" />
                </span>
                <h3 className="font-semibold text-foreground mb-2">{p.title}</h3>
                <p className="text-sm text-muted-foreground leading-relaxed">{p.body}</p>
              </div>
            ))}
          </div>
        </div>
      </PitchSection>

      {/* ═══ RIGHT CHOICE / WRONG CHOICE ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-10" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">Your next project</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight mb-5">
                When to start it in Grit
              </h2>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit pays off when you are building a real business product, and most of all when it
                has more than one client.
              </p>
            </div>

            <div data-gsap-reveal className="grid sm:grid-cols-2 gap-3 mb-14">
              {RIGHT_FOR.map((r) => (
                <div
                  key={r}
                  className="flex items-start gap-3 rounded-xl border border-primary/20 bg-primary/[0.03] px-4 py-3 text-sm text-foreground/85"
                >
                  <Check className="h-4 w-4 text-emerald-400 mt-0.5 shrink-0" strokeWidth={2.5} />
                  {r}
                </div>
              ))}
            </div>

            <div className="max-w-2xl mb-8" data-gsap-reveal>
              <p className="tag-mono text-muted-foreground mb-4">Honest limits</p>
              <h3 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground leading-tight mb-4">
                When Grit is the wrong choice
              </h3>
              <p className="text-base text-muted-foreground leading-relaxed">
                A framework that fits everything fits nothing. Better to find out on this page than
                three weeks into a rewrite.
              </p>
            </div>

            <div className="grid sm:grid-cols-2 gap-4">
              {NOT_FOR.map((n) => (
                <div key={n.title} data-gsap-reveal className="rounded-2xl border border-border/40 bg-card/30 p-6">
                  <div className="flex items-start gap-2.5 mb-2.5">
                    <X className="h-4 w-4 text-rose-400/80 mt-0.5 shrink-0" strokeWidth={2.5} />
                    <h4 className="font-semibold text-foreground text-[15px] leading-snug">{n.title}</h4>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed pl-[26px]">{n.body}</p>
                </div>
              ))}
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ FRAMING: OPEN SOURCE ═══ */}
      <PitchSection className="py-16 md:py-20 border-b border-border/40">
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          Everything Grit generates is ordinary code in your repo, which means the framework can
          never become the thing standing between you and a fix.
        </p>
        <p className="text-lg text-muted-foreground leading-relaxed">
          It is MIT licensed, developed in the open, and supported by{' '}
          <Link href="/sponsor" className="text-primary hover:underline underline-offset-4">
            sponsors
          </Link>{' '}
          rather than a sales team. Releases are public, versioned, and shipped weekly. For the
          longer story, why Go, why React, and what Grit borrows from Laravel and Rails, read{' '}
          <Link href="/docs/getting-started/philosophy" className="text-primary hover:underline underline-offset-4">
            the philosophy doc
          </Link>
          .
        </p>
      </PitchSection>

      {/* ═══ COMMUNITY ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div data-gsap-reveal>
              <CommunityCTA />
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ CLOSING CTA ═══ */}
      <section className="relative py-24 md:py-32 px-6 overflow-hidden">
        <div className="absolute inset-0 -z-20 bg-grit-grid mask-fade-center opacity-70" />
        <GlowOrb className="top-0 left-1/2 -translate-x-1/2 h-[500px] w-[500px] bg-primary/[0.08]" duration={22} />
        <div className="relative max-w-2xl mx-auto text-center">
          <FadeIn>
            <Sparkles className="h-6 w-6 text-primary mx-auto mb-5" />
            <h2 className="text-3xl md:text-5xl font-bold tracking-tight text-foreground mb-5 leading-tight">
              Your next product is one
              <br />
              command away
            </h2>
            <p className="text-lg text-muted-foreground mb-8">
              Install the CLI and scaffold a production-ready Go + React product, web, admin, mobile
              and desktop, in minutes.
            </p>
          </FadeIn>
          <FadeIn delay={0.1}>
            <div className="max-w-lg mx-auto text-left mb-8">
              <CodeBlock
                terminal
                filename="Install: macOS / Linux"
                code={`curl -fsSL https://gritframework.dev/install.sh | sh`}
              />
            </div>
          </FadeIn>
          <FadeIn delay={0.16}>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
              <Link
                href="/docs/getting-started/quick-start"
                className="group inline-flex items-center gap-2 h-11 px-7 rounded-full bg-primary text-primary-foreground font-semibold text-sm hover:bg-primary/90 transition-colors"
              >
                Get started
                <ArrowRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
              </Link>
              <Link
                href="/courses"
                className="inline-flex items-center h-11 px-7 rounded-full border border-border/60 text-foreground font-medium text-sm hover:bg-accent/30 transition-colors"
              >
                Follow a course
              </Link>
            </div>
            <p className="text-xs text-muted-foreground/60 mt-6 font-mono">
              {`Grit v${GRIT_VERSION} · MIT licensed · Go + React`}
            </p>
          </FadeIn>
        </div>
      </section>
    </div>
  )
}
