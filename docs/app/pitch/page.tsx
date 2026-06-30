import Link from 'next/link'
import type { Metadata } from 'next'
import { ArrowRight, Check, X, Terminal, Sparkles, Bot, ShieldCheck, Layers } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { CodeBlock } from '@/components/code-block'
import { FadeIn, GSAPSection, GlowOrb } from '@/components/motion-primitives'
import { GridFrame } from '@/components/grid-frame'
import { GRIT_VERSION } from '@/config/site'

export const metadata: Metadata = {
  title: 'The Pitch',
  description:
    'What is actually painful about building full-stack apps — and how Grit compresses weeks of plumbing into a single command. The case for Go + React, batteries included.',
  alternates: { canonical: 'https://gritframework.dev/pitch' },
}

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
    <div className="min-h-screen bg-background">
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
            <h1 className="text-4xl md:text-6xl font-bold tracking-tight text-foreground leading-[1.05] mb-7">
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
            <p className="tag-mono text-primary mb-4">Why it holds up</p>
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground leading-tight">
              Opinions where it counts, escape hatches everywhere else
            </h2>
          </div>
          <div className="grid md:grid-cols-3 gap-4">
            {[
              {
                icon: Bot,
                title: 'Built for agents too',
                body: 'Grit has an opinion on routing, queues, auth, storage, and AI — thousands of decisions an AI assistant never has to guess. It ships a SKILL.md so agents extend it correctly.',
              },
              {
                icon: ShieldCheck,
                title: 'Hardened by default',
                body: 'OWASP-2025 patterns baked in: SSRF defence, IDOR-safe ownership checks, CSRF, strict CSP, rate limiting, and a tamper-evident audit log — on day one, not day ninety.',
              },
              {
                icon: Layers,
                title: 'Five architectures',
                body: 'Embed a SPA in the Go binary, split web / admin / API into a monorepo, go API-only, or add mobile and desktop. Same generators, same patterns, your call.',
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
