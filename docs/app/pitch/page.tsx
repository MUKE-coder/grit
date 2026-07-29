import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Terminal, Sparkles, Bot, ShieldCheck, Layers } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { FadeIn, GSAPSection, GlowOrb } from '@/components/motion-primitives'
import { GridFrame } from '@/components/grid-frame'
import { CommunityCTA } from '@/components/community-cta'
import { GRIT_VERSION } from '@/config/site'

export const metadata: Metadata = {
  title: 'The Pitch',
  description:
    'The thinking behind Grit: what we optimise for, what we deliberately trade away, and when Grit is the wrong choice. The case for Go + React, batteries included.',
  alternates: { canonical: 'https://gritframework.dev/pitch' },
}

/**
 * Every principle is written as a trade — "A over B" — because a stated
 * preference that costs nothing is not a principle, it is a feature list. The
 * `gives` field names what you actually give up, so a reader can decide the
 * trade is wrong for them without having to discover it three weeks in.
 */
const PRINCIPLES = [
  {
    title: 'Generated over hidden',
    body: 'Grit writes real files into your repo — the model, service, handler, routes, hooks and admin page. Not a runtime dependency that does it invisibly at boot.',
    gives:
      'You give up a small repo. You get code you can read, edit, delete, and debug with a stack trace that points at your own file. There is no "fighting the framework" because there is no framework in the way — just Go and TypeScript you own.',
  },
  {
    title: 'Boring over clever',
    body: 'Gin, GORM, Next.js, Postgres, Redis. Every one has years of production behind it and its answers already written down somewhere.',
    gives:
      'You give up novelty. We are not going to invent a router. The interesting part of your project should be your product, not the stack underneath it.',
  },
  {
    title: 'Opinionated over configurable',
    body: 'One way to do auth. One queue. One storage layer. Escape hatches where you need them, not a menu where you do not.',
    gives:
      'You give up the pleasure of choosing. You get a codebase where every Grit project looks the same, onboarding takes an afternoon, and an AI agent never has to guess which of six patterns you picked.',
  },
  {
    title: 'Whole-stack over layer-by-layer',
    body: 'A resource is not a Go model or a TypeScript type. It is both — plus the handler, the Zod schema, the React Query hook and the admin page — generated together from one definition.',
    gives:
      'You give up picking your own frontend and backend languages. You get two halves that cannot drift apart, because one command writes both.',
  },
  {
    title: 'Secure on day one over hardened on day ninety',
    body: 'CSRF, strict CSP, rate limiting, SSRF defence, IDOR-safe ownership checks, field-level encryption, server-side sessions and a tamper-evident audit log — in the scaffold.',
    gives:
      'You give up the option to skip it. That is the point: security is the work that never wins a prioritisation meeting against a feature, so it ships before the meeting happens.',
  },
  {
    title: 'Verified over asserted',
    body: 'Every release runs a 73-check matrix across 13 project shapes — scaffold, install, build, typecheck, test — then gets re-verified by downloading the published binary and building a fresh project with it.',
    gives:
      'You give up a faster release cadence. Claims in a changelog are cheap; we would rather find the bug than describe the feature.',
  },
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
    body: 'You can swap them — it is your code. But you would be working against the generators, and that cost compounds with every resource you add.',
  },
  {
    title: 'You are building one small service with no UI',
    body: 'The batteries become weight you never use. `--api` trims a lot of it, but a plain Gin binary may genuinely be the simpler answer.',
  },
  {
    title: 'You need a decade-old ecosystem',
    body: 'Grit is young. Fewer Stack Overflow answers, fewer third-party plugins, fewer people who have hit your exact bug. What you get instead is a maintainer who ships weekly and answers questions directly.',
  },
] as const

// Small reusable section wrapper that pins a faint crosshair to a corner —
// the Firecrawl signature — and reveals its children on scroll.
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
        {/* soft glow only — grid + rails come from the shared GridFrame */}
        <div
          className="absolute inset-0 -z-10"
          style={{
            background:
              'radial-gradient(ellipse 60% 60% at 50% -10%, hsl(var(--primary) / 0.14), transparent 60%)',
          }}
        />
        <GlowOrb className="-top-32 left-1/3 h-[420px] w-[420px] bg-primary/[0.10]" duration={20} />

        {/* crosshair accents */}
        <span className="crosshair absolute top-28 left-[14%] text-foreground/20 hidden md:block" style={{ width: 16, height: 16 }} />
        <span className="crosshair absolute top-40 right-[16%] text-primary/30 hidden md:block" style={{ width: 16, height: 16 }} />

        <div className="relative max-w-3xl mx-auto px-6 pt-24 pb-20 md:pt-32 md:pb-24">
          <FadeIn>
            <span className="tag-mono text-primary mb-6 block">The Pitch</span>
          </FadeIn>
          <FadeIn delay={0.08}>
            <h1 className="font-display text-4xl md:text-6xl font-bold tracking-tight text-foreground leading-[1.05] mb-7">
              What if your backend
              <br />
              shipped with the frontend
              <br />
              <span className="bg-gradient-to-r from-primary via-sky-400 to-primary bg-clip-text text-transparent">
                already wired?
              </span>
            </h1>
          </FadeIn>
          <FadeIn delay={0.16}>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed max-w-2xl">
              Every full-stack app starts the same way: weeks of plumbing before a single
              feature ships. Grit asks a different question — what is the smallest set of
              commands that takes you from <span className="text-foreground font-medium">idea</span> to a{' '}
              <span className="text-foreground font-medium">production-ready</span> Go + React app?
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
                href="/docs"
                className="inline-flex items-center h-11 px-6 rounded-full border border-border/60 text-foreground font-medium text-sm hover:bg-accent/30 transition-colors"
              >
                Read the docs
              </Link>
            </div>
          </FadeIn>
        </div>
      </section>

      {/* ═══ FRAMING — this page is about the why ═══ */}
      <PitchSection className="py-16 md:py-20 border-b border-border/40">
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          If you are evaluating Grit, or just curious how we think about building a
          full-stack framework, this page explains the{' '}
          <span className="text-foreground font-medium">why</span> behind the decisions —
          including the ones we know are not for everyone.
        </p>
        <p className="text-lg text-muted-foreground leading-relaxed mb-5">
          Grit brings coherence and speed to Go + React development without hiding the
          platform underneath it. Everything it generates is ordinary code in your repo,
          which means the framework can never become the thing standing between you and a
          fix.
        </p>
        <p className="text-lg text-muted-foreground leading-relaxed">
          It is MIT licensed, developed in the open, and supported by{' '}
          <Link href="/sponsor" className="text-primary hover:underline underline-offset-4">
            sponsors
          </Link>{' '}
          rather than a sales team. Releases are public, versioned, and shipped weekly.
        </p>
        <p className="text-base text-muted-foreground leading-relaxed mt-6">
          For the longer story — why Go, why React, and what Grit borrows from Laravel and
          Rails — read{' '}
          <Link
            href="/docs/getting-started/philosophy"
            className="text-primary hover:underline underline-offset-4"
          >
            the philosophy doc
          </Link>
          .
        </p>
      </PitchSection>

      {/* ═══ THE PROBLEM ═══ */}
      <PitchSection className="py-20 md:py-24 border-b border-border/40">
        <p className="tag-mono text-muted-foreground mb-4">The boilerplate tax</p>
        <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground mb-6 leading-tight">
          The first month is rarely about your product
        </h2>
        <p className="text-lg text-muted-foreground leading-relaxed mb-8">
          Before you write the feature that matters, you re-solve the same problems you
          solved on the last project — and the one before that. None of it is hard. All
          of it is slow.
        </p>
        <div className="grid sm:grid-cols-2 gap-3">
          {[
            'Auth: JWT, refresh, OAuth, 2FA, password reset',
            'An admin panel with tables, forms, and filters',
            'File uploads to S3 / R2 with presigned URLs',
            'Background jobs, queues, and a scheduler',
            'Email templates and a transactional sender',
            'Rate limiting, CORS, security headers, audit logs',
            'Type-safe API clients and React Query hooks',
            'Docker, migrations, seeders, CI, deploy scripts',
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
          That is months of work that produces <span className="text-foreground font-medium">zero</span> product
          differentiation. Grit treats all of it as solved — generated, wired, and hardened the
          moment you scaffold.
        </p>
      </PitchSection>

      {/* ═══ BEFORE / AFTER ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40 overflow-hidden">
        <div className="absolute inset-0 -z-20 bg-grit-dots mask-fade-y opacity-60" />
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="text-center mb-12" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">Compress the work</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight">
                The same outcome, two timelines
              </h2>
            </div>
            <div className="grid md:grid-cols-2 gap-5 items-stretch">
              {/* Without */}
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
                    'Stand up Docker, CI, and a deploy story',
                  ].map((l) => (
                    <li key={l} className="flex items-start gap-2.5">
                      <span className="text-rose-400/60 font-mono text-xs mt-0.5">—</span>
                      {l}
                    </li>
                  ))}
                </ul>
                <div className="mt-5 pt-4 border-t border-border/40 flex items-center justify-between">
                  <span className="text-xs text-muted-foreground font-mono">Time to first feature</span>
                  <span className="text-sm font-semibold text-rose-400">~3–6 weeks</span>
                </div>
              </div>

              {/* With */}
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
                  code={`grit new my-app --triple
grit generate resource Product \\
  --fields "name:string,price:float,stock:int"`}
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
                principle, it is a feature list — so each one below names what you give up.
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
                    <span className="font-mono text-sm font-semibold text-primary">
                      {i + 1}.
                    </span>
                    <h3 className="font-semibold text-foreground text-[17px] leading-snug">
                      {p.title}
                    </h3>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-4">
                    {p.body}
                  </p>
                  <p className="text-sm leading-relaxed text-foreground/70 border-l-2 border-primary/30 pl-4 mt-auto">
                    {p.gives}
                  </p>
                </div>
              ))}
            </div>
          </GSAPSection>
        </div>
      </section>

      {/* ═══ ONE COMMAND, FULL STACK ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-10" data-gsap-reveal>
              <p className="tag-mono text-primary mb-4">One command, full stack</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight mb-5">
                Generate the backend and the frontend in the same breath
              </h2>
              <p className="text-lg text-muted-foreground leading-relaxed">
                <code className="text-primary bg-primary/5 px-1.5 py-0.5 rounded text-base">grit generate resource</code>{' '}
                emits the Go model, service, handler, routes, Zod schema, TypeScript types,
                React Query hooks, and an admin page — all consistent, all type-safe.
              </p>
            </div>

            <div
              data-gsap-reveal
              className="rounded-2xl overflow-hidden border border-border/50 bg-white dark:bg-[#0d1117] shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]"
            >
              <div className="flex items-center gap-3 px-4 py-2.5 bg-[#f6f8fa] dark:bg-[#161b22] border-b border-[#d0d7de] dark:border-white/[0.08]">
                <div className="flex items-center gap-1.5">
                  <div className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
                  <div className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
                </div>
                <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-400">generated · one command</span>
              </div>
              <div className="grid md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-[#d0d7de] dark:divide-white/[0.06]">
                <CodeBlock
                  language="go"
                  code={`func (h *ProductHandler) List(c *gin.Context) {
    var products []models.Product
    h.DB.
        Where("user_id = ?", c.GetString("user_id")).
        Order("created_at desc").
        Find(&products)

    c.JSON(http.StatusOK, gin.H{"data": products})
}`}
                  className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
                />
                <CodeBlock
                  language="tsx"
                  code={`export function useProducts() {
  return useQuery({
    queryKey: ['products'],
    queryFn: async () => {
      const res = await api.get('/api/products')
      return res.data.data as Product[]
    },
  })
}`}
                  className="!border-0 !rounded-none !shadow-none !bg-transparent dark:!bg-transparent !m-0"
                />
              </div>
              <div className="flex items-center justify-center gap-2 px-4 py-2.5 bg-[#f6f8fa] dark:bg-[#161b22] border-t border-[#d0d7de] dark:border-white/[0.08]">
                <Terminal className="h-3 w-3 text-[#57606a] dark:text-slate-500" />
                <span className="text-[11px] font-mono text-[#57606a] dark:text-slate-400">
                  Both files written by{' '}
                  <span className="text-primary font-semibold">grit generate resource Product</span>
                </span>
              </div>
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
                body: 'Because there is one way to do each thing, an AI assistant has nothing to guess. Grit ships a SKILL.md, and grit mcp serve hands an agent the real route table and model definitions over MCP — parsed from your source, read-only.',
              },
              {
                icon: ShieldCheck,
                title: 'Security you did not schedule',
                body: 'The hardening in every scaffold is work no one budgets for: SSRF defence, IDOR-safe ownership checks, GDPR export and erasure, SSO with SAML, field-level encryption. Already there, already tested.',
              },
              {
                icon: Layers,
                title: 'One set of patterns, five shapes',
                body: 'Embed a SPA in the Go binary, split web / admin / API into a monorepo, go API-only, or add mobile and desktop. The generators and conventions do not change — only the shape does.',
              },
            ].map((p) => (
              <div
                key={p.title}
                className="card-grit rounded-2xl border border-border/40 bg-card/40 p-6"
              >
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

      {/* ═══ WHEN GRIT IS THE WRONG CHOICE ═══ */}
      <section className="relative py-20 md:py-24 px-6 border-b border-border/40">
        <div className="max-w-5xl mx-auto">
          <GSAPSection>
            <div className="max-w-2xl mb-10" data-gsap-reveal>
              <p className="tag-mono text-muted-foreground mb-4">Honest limits</p>
              <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight mb-5">
                When Grit is the wrong choice
              </h2>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A framework that fits everything fits nothing. Here is where Grit does not
                fit — better to find out on this page than three weeks into a rewrite.
              </p>
            </div>

            <div className="grid sm:grid-cols-2 gap-4">
              {NOT_FOR.map((n) => (
                <div
                  key={n.title}
                  data-gsap-reveal
                  className="rounded-2xl border border-border/40 bg-card/30 p-6"
                >
                  <div className="flex items-start gap-2.5 mb-2.5">
                    <X
                      className="h-4 w-4 text-rose-400/80 mt-0.5 shrink-0"
                      strokeWidth={2.5}
                    />
                    <h3 className="font-semibold text-foreground text-[15px] leading-snug">
                      {n.title}
                    </h3>
                  </div>
                  <p className="text-sm text-muted-foreground leading-relaxed pl-[26px]">
                    {n.body}
                  </p>
                </div>
              ))}
            </div>

            <p className="text-base text-muted-foreground leading-relaxed mt-8 max-w-2xl">
              If none of those describe you, the rest of the decisions on this page were
              made with your project in mind.
            </p>
          </GSAPSection>
        </div>
      </section>

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
              Your next project is one
              <br />
              command away
            </h2>
            <p className="text-lg text-muted-foreground mb-8">
              Install the CLI and scaffold a production-ready Go + React app in minutes.
            </p>
          </FadeIn>
          <FadeIn delay={0.1}>
            <div className="max-w-lg mx-auto text-left mb-8">
              <CodeBlock
                terminal
                filename="Install — macOS / Linux"
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
