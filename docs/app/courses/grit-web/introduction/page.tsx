import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Introduction to Grit — Grit Web Course',
  description: 'Learn what Grit is, why it exists, what problems it solves, how it compares to other frameworks, and what you will learn in the Grit Web course series.',
}

export default function IntroductionToGritCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Introduction to Grit</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 0 of 9</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Introduction to Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            This course introduces Grit Framework — what it is, why it was created, how it compares to other
            tools, and what you{"'"}ll learn in this course series. No coding in this course — just concepts and
            understanding. By the end, you{"'"}ll know exactly what Grit gives you and whether it{"'"}s right for
            your project.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: Welcome to Grit ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Welcome to Grit</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Welcome! You{"'"}re about to learn a framework that will change how you build web applications.
            Whether you{"'"}re a Go developer looking for a batteries-included solution, a JavaScript developer
            curious about Go, or a solo founder who needs to ship fast — Grit was built for you.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is Course 0 — the conceptual foundation. There{"'"}s no coding here, no terminal commands, no
            setup. Just ideas. You{"'"}ll understand what Grit is, why it exists, what problems it solves, how
            it compares to frameworks you already know, and what you{"'"}ll build throughout this course series.
          </p>

          <Tip>
            Take your time with this course. The concepts here will make every future course click faster. Understanding
            the {'"'}why{'"'} behind Grit is just as important as learning the {'"'}how{'"'}.
          </Tip>
        </section>

        {/* ═══ Section 2: What is Grit? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Grit?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit is a <strong className="text-foreground">full-stack meta-framework</strong> that combines
            Go (for the backend) with React (for the frontend). It gives you a complete, production-ready
            project in one command — with authentication, admin panel, file uploads, background jobs, email,
            and more already wired up and working.
          </p>

          <Definition term="Meta-framework">
            A framework that builds on top of other frameworks. Grit combines Gin (Go web framework), GORM
            (Go ORM), Next.js or TanStack Router (React frameworks), Tailwind CSS, and shadcn/ui into one
            cohesive system. You don{"'"}t configure each one separately — Grit does it for you.
          </Definition>

          <Definition term="Full-stack">
            Both the backend (server, database, API) and the frontend (the UI users interact with). Some
            frameworks only handle one side — Grit handles both.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What Grit is <strong className="text-foreground">NOT</strong>:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Not a boilerplate</strong> — it{"'"}s a living framework with a CLI that keeps helping you after the initial scaffold. Generate resources, sync types, deploy — all from the same tool.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Not a template</strong> — it generates real, customizable code you own. Every file is yours to modify, extend, or delete.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Not locked in</strong> — you can stop using Grit tomorrow and your code still works. It{"'"}s standard Go, standard React, standard PostgreSQL. No proprietary runtime.</li>
          </ul>

          <Challenge number={1} title="Boilerplate vs Framework">
            <p>In your own words, explain the difference between a boilerplate/template and a framework. Why is
            Grit a framework and not just a starter template? Think about what happens <em>after</em> the initial
            project is created.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: The Problem Grit Solves ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Problem Grit Solves</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every time you start a new project, you pay what we call the <strong className="text-foreground">{'"'}Setup Tax{'"'}</strong> — weeks
            of work building infrastructure before you write a single line of business logic:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Task</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Typical Time</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Authentication (JWT, password hashing, tokens)</td>
                  <td className="px-4 py-3">2–3 days</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Admin panel (data tables, forms, dashboard)</td>
                  <td className="px-4 py-3">3–5 days</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">File uploads (S3, presigned URLs, image processing)</td>
                  <td className="px-4 py-3">1–2 days</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Email (templates, sending, testing)</td>
                  <td className="px-4 py-3">1 day</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Background jobs (queue, workers, retry logic)</td>
                  <td className="px-4 py-3">1–2 days</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">API documentation (OpenAPI spec, interactive UI)</td>
                  <td className="px-4 py-3">1 day</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Security (rate limiting, WAF, headers)</td>
                  <td className="px-4 py-3">1 day</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Docker setup (PostgreSQL, Redis, MinIO)</td>
                  <td className="px-4 py-3">half day</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3">Deployment (systemd, reverse proxy, TLS)</td>
                  <td className="px-4 py-3">1 day</td>
                </tr>
                <tr>
                  <td className="px-4 py-3">Type safety between backend and frontend</td>
                  <td className="px-4 py-3">ongoing pain</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            <strong className="text-foreground">Total: 2–4 weeks</strong> before you write a single line of
            business logic. And this happens every single time you start a new project.
          </p>

          <Note>
            Grit gives you ALL of this in one command, in under 30 seconds. Run <Code>grit new myapp</Code> and
            every item in that table is already done — configured, tested, and ready to use.
          </Note>

          <Challenge number={2} title="Your Setup Tax">
            <p>Think about the last project you started. How long did you spend on setup and boilerplate before
            writing actual features? List the things you set up manually. How many of them are in the table above?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: The Motivation — Why Grit Was Created ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Motivation — Why Grit Was Created</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit was created by <strong className="text-foreground">Muke JohnBaptist</strong>, a full-stack
            developer who was tired of rebuilding the same infrastructure for every project. Auth, admin panels,
            file uploads, background jobs — they should just work. Now they do.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The key insight: most web applications need the <strong className="text-foreground">same 80% of
            infrastructure</strong>. Only 20% is unique business logic. Grit gives you the 80% so you can
            focus on the 20% that makes your app special.
          </p>

          <Definition term="Convention over Configuration">
            One of Grit{"'"}s core principles. Instead of making you choose and configure every tool, Grit makes
            sensible decisions for you. File structure, naming conventions, API response format — all decided.
            But everything is overridable if you disagree.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Think about it: how many times have you debated whether to put your handlers in <Code>/controllers</Code> or <Code>/handlers</Code>?
            Whether to use <Code>camelCase</Code> or <Code>snake_case</Code> for API fields? Whether to structure
            your project by feature or by layer? Grit makes those decisions for you — good decisions, based on
            years of building production apps — so you can focus on building your product.
          </p>

          <Challenge number={3} title="The 80/20 Rule">
            <p>Name 3 things that are the same in almost every web app (hint: authentication is one). For each,
            think about how much time you typically spend setting them up from scratch.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: What Makes Grit Different ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What Makes Grit Different</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            There are plenty of frameworks out there. Here{"'"}s what sets Grit apart:
          </p>

          <div className="space-y-6 mb-6">
            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">1. Choose Your Architecture</h3>
              <p className="text-muted-foreground leading-relaxed">
                No other framework lets you pick between 5 architectures (single, double, triple, API, mobile) from
                the same CLI. Laravel is always single. Next.js is always single. Rails is always single. Grit
                adapts to <em>your</em> project.
              </p>
            </div>

            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">2. Choose Your Frontend</h3>
              <p className="text-muted-foreground leading-relaxed">
                Next.js for SEO, TanStack Router for speed. Same backend, same components, different routing.
                Pick what fits your project without changing anything else.
              </p>
            </div>

            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">3. Code Generation That Works Both Ways</h3>
              <p className="text-muted-foreground leading-relaxed">
                <Code>grit generate</Code> creates files. <Code>grit remove</Code> cleanly deletes them. No
                leftover code, no broken imports. Most generators only go one direction — Grit goes both.
              </p>
            </div>

            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">4. Type Safety Across the Stack</h3>
              <p className="text-muted-foreground leading-relaxed">
                Go types → TypeScript types → Zod schemas. Change a Go model, run <Code>grit sync</Code>, and
                your frontend types update automatically. No API contract drift. No runtime surprises.
              </p>
            </div>

            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">5. One-Command Deployment</h3>
              <p className="text-muted-foreground leading-relaxed">
                <Code>grit deploy</Code> builds, uploads, and configures your production server. No CI/CD needed
                to get started. From localhost to live in one command.
              </p>
            </div>

            <div className="bg-muted/10 border border-border/40 rounded-lg p-5">
              <h3 className="text-lg font-semibold text-foreground mb-2">6. Built-in Security</h3>
              <p className="text-muted-foreground leading-relaxed">
                Sentinel WAF, rate limiting, 2FA with TOTP — not plugins, not add-ons. Security by default,
                not security as an afterthought.
              </p>
            </div>
          </div>

          <Challenge number={4} title="Framework Comparison">
            <p>Pick another framework you{"'"}ve used (Laravel, Rails, Next.js, Django, Express). List 3 things
            Grit does that your framework doesn{"'"}t include out of the box.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: The Tech Stack ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Tech Stack</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit is opinionated about its tech stack. Every technology was chosen for a reason — performance,
            developer experience, community support, or all three. Here{"'"}s the complete picture:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Layer</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Technology</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Why</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Backend</td>
                  <td className="px-4 py-3">Go 1.21+</td>
                  <td className="px-4 py-3">Fast, compiled, type-safe</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Web Framework</td>
                  <td className="px-4 py-3">Gin</td>
                  <td className="px-4 py-3">Most popular Go web framework</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">ORM</td>
                  <td className="px-4 py-3">GORM</td>
                  <td className="px-4 py-3">Auto-migration, relationships, hooks</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Database</td>
                  <td className="px-4 py-3">PostgreSQL / SQLite</td>
                  <td className="px-4 py-3">Production / Desktop</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Cache</td>
                  <td className="px-4 py-3">Redis</td>
                  <td className="px-4 py-3">In-memory speed, pub/sub, queues</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Job Queue</td>
                  <td className="px-4 py-3">asynq</td>
                  <td className="px-4 py-3">Redis-backed, reliable, Go-native</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">File Storage</td>
                  <td className="px-4 py-3">S3-compatible</td>
                  <td className="px-4 py-3">AWS S3, Cloudflare R2, MinIO</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Email</td>
                  <td className="px-4 py-3">Resend</td>
                  <td className="px-4 py-3">Modern email API, great DX</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">AI</td>
                  <td className="px-4 py-3">Vercel AI Gateway</td>
                  <td className="px-4 py-3">Multi-provider, streaming</td>
                </tr>
                <tr className="border-b border-border/20 bg-muted/10">
                  <td className="px-4 py-3 font-medium text-foreground" colSpan={3}>Frontend</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Router</td>
                  <td className="px-4 py-3">Next.js 14+ / TanStack Router</td>
                  <td className="px-4 py-3">SEO (Next) or SPA speed (TanStack)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">UI Library</td>
                  <td className="px-4 py-3">React</td>
                  <td className="px-4 py-3">Largest ecosystem, most talent</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Styling</td>
                  <td className="px-4 py-3">Tailwind CSS + shadcn/ui</td>
                  <td className="px-4 py-3">Utility-first, copy-paste components</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Data Fetching</td>
                  <td className="px-4 py-3">TanStack Query</td>
                  <td className="px-4 py-3">Caching, refetching, optimistic updates</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Validation</td>
                  <td className="px-4 py-3">Zod</td>
                  <td className="px-4 py-3">TypeScript-first schema validation</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Forms</td>
                  <td className="px-4 py-3">React Hook Form</td>
                  <td className="px-4 py-3">Performant, minimal re-renders</td>
                </tr>
                <tr className="border-b border-border/20 bg-muted/10">
                  <td className="px-4 py-3 font-medium text-foreground" colSpan={3}>Infrastructure</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Monorepo</td>
                  <td className="px-4 py-3">Turborepo + pnpm</td>
                  <td className="px-4 py-3">Fast builds, shared packages</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Security</td>
                  <td className="px-4 py-3">Sentinel</td>
                  <td className="px-4 py-3">WAF + rate limiting</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Observability</td>
                  <td className="px-4 py-3">Pulse</td>
                  <td className="px-4 py-3">Tracing + metrics</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">DB Browser</td>
                  <td className="px-4 py-3">GORM Studio</td>
                  <td className="px-4 py-3">Visual database management</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Desktop</td>
                  <td className="px-4 py-3">Wails v2</td>
                  <td className="px-4 py-3">Native desktop apps with Go + React</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Deployment</td>
                  <td className="px-4 py-3">SSH + systemd + Caddy</td>
                  <td className="px-4 py-3">Simple, reliable, auto-HTTPS</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Definition term="Go">
            A programming language created by Google. Known for fast compilation, excellent concurrency, and
            producing single binary executables. Used by Docker, Kubernetes, and many high-performance systems.
          </Definition>

          <Definition term="React">
            A JavaScript library for building user interfaces, created by Meta (Facebook). Uses components and
            a virtual DOM for efficient UI updates. The most popular frontend library in the world.
          </Definition>

          <Definition term="Tailwind CSS">
            A utility-first CSS framework. Instead of writing custom CSS, you apply small utility classes
            directly in your HTML: <Code>{'"'}flex items-center gap-4 bg-blue-500 text-white{'"'}</Code>.
          </Definition>

          <Definition term="shadcn/ui">
            A collection of beautifully-designed, accessible React components built on Tailwind CSS. Not a
            component library you install — you copy the source code into your project and own it. This means
            you can customize every pixel.
          </Definition>

          <Challenge number={5} title="Tech Stack Reasoning">
            <p>For each layer in the tech stack (backend, ORM, database, frontend, styling, data fetching),
            explain in one sentence why that technology was chosen. Think about popularity, performance, and
            developer experience.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Who is Grit For? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Who is Grit For?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit was designed for seven types of developers:
          </p>

          <div className="space-y-4 mb-6">
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">1.</span>
              <div>
                <strong className="text-foreground">Go Developers</strong>
                <p className="text-muted-foreground">Who want a batteries-included framework like Laravel or Rails, but in Go. No more gluing together 15 libraries for every project.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">2.</span>
              <div>
                <strong className="text-foreground">Laravel Developers</strong>
                <p className="text-muted-foreground">Migrating to Go and wanting familiar patterns — single-app mode, artisan-like CLI, convention-driven structure.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">3.</span>
              <div>
                <strong className="text-foreground">Next.js Developers</strong>
                <p className="text-muted-foreground">Who want a proper backend instead of serverless functions. Real database connections, real background jobs, real performance.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">4.</span>
              <div>
                <strong className="text-foreground">MERN Stack Developers</strong>
                <p className="text-muted-foreground">Who want a structured monorepo with type safety instead of loosely connected services.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">5.</span>
              <div>
                <strong className="text-foreground">Solo Developers</strong>
                <p className="text-muted-foreground">Building SaaS products who need to ship fast without cutting corners on quality or security.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">6.</span>
              <div>
                <strong className="text-foreground">Agencies</strong>
                <p className="text-muted-foreground">Building client projects who need consistent, production-ready scaffolding for every engagement.</p>
              </div>
            </div>
            <div className="flex gap-3">
              <span className="text-primary font-bold shrink-0">7.</span>
              <div>
                <strong className="text-foreground">Students and Beginners</strong>
                <p className="text-muted-foreground">Learning full-stack development with a guided framework that teaches best practices by example.</p>
              </div>
            </div>
          </div>

          <Challenge number={6} title="Find Your Persona">
            <p>Which persona fits you best? Write down your background and what you hope to build with Grit.
            If none of these fit perfectly, describe your own persona — who are you and what do you need?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: The 5 Architecture Modes ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The 5 Architecture Modes</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            One of Grit{"'"}s most powerful features is architecture flexibility. When you run <Code>grit new</Code>,
            you choose the architecture that fits your project. You{"'"}ll learn each in detail in Course 1, but
            here{"'"}s a quick overview:
          </p>

          <div className="space-y-4 mb-6">
            <div className="bg-muted/10 border border-border/40 rounded-lg p-4">
              <h3 className="font-semibold text-foreground mb-1">Triple: Web + Admin + API</h3>
              <p className="text-sm text-muted-foreground">A monorepo with three apps — a public website, an admin dashboard, and a Go API. Best for SaaS apps, platforms, and anything that needs a separate admin experience.</p>
            </div>
            <div className="bg-muted/10 border border-border/40 rounded-lg p-4">
              <h3 className="font-semibold text-foreground mb-1">Double: Web + API</h3>
              <p className="text-sm text-muted-foreground">A monorepo with two apps — a public website and a Go API. Best for simpler apps, blogs, and landing pages with a backend.</p>
            </div>
            <div className="bg-muted/10 border border-border/40 rounded-lg p-4">
              <h3 className="font-semibold text-foreground mb-1">Single: One Go Binary</h3>
              <p className="text-sm text-muted-foreground">One Go binary with the frontend embedded. Like Laravel — everything in one place. Best for microservices and internal tools.</p>
            </div>
            <div className="bg-muted/10 border border-border/40 rounded-lg p-4">
              <h3 className="font-semibold text-foreground mb-1">API Only: Go Backend</h3>
              <p className="text-sm text-muted-foreground">Just the Go API — no frontend. Best for mobile backends, headless APIs, and services consumed by other apps.</p>
            </div>
            <div className="bg-muted/10 border border-border/40 rounded-lg p-4">
              <h3 className="font-semibold text-foreground mb-1">Mobile: API + Expo</h3>
              <p className="text-sm text-muted-foreground">Go API plus a React Native app built with Expo. Best for mobile-first products.</p>
            </div>
          </div>

          <Tip>
            There{"'"}s also <Code>grit new-desktop</Code> for standalone Wails desktop apps with SQLite. It{"'"}s
            a separate command because desktop apps have different requirements than web apps.
          </Tip>

          <Challenge number={7} title="Choose Your Architecture">
            <p>For each project idea below, which architecture would you choose and why?</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>(a) A blog with an admin panel</li>
              <li>(b) A mobile fitness app</li>
              <li>(c) An internal company tool</li>
              <li>(d) A SaaS with billing and admin</li>
              <li>(e) A REST API for a mobile game</li>
            </ul>
          </Challenge>
        </section>

        {/* ═══ Section 9: What Ships With Every Project ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What Ships With Every Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project comes with 18 features out of the box. You{"'"}ll learn each one in detail
            throughout this course series, but here{"'"}s what you get on day one:
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mb-6">
            {[
              'Auth (JWT + 2FA + OAuth)',
              'Admin Panel',
              'Code Generator',
              'File Storage (S3)',
              'Email (Resend)',
              'Background Jobs',
              'Cron Scheduler',
              'AI Integration',
              'Redis Cache',
              'GORM Studio',
              'API Docs (OpenAPI)',
              'Security (Sentinel)',
              'Observability (Pulse)',
              '100 UI Components',
              'Docker Setup',
              'Test Suite',
              'One-Cmd Deployment',
              'Type Sync (Go → TS)',
            ].map((feature) => (
              <div key={feature} className="flex items-center gap-2 text-sm text-muted-foreground bg-muted/10 border border-border/30 rounded-lg px-3 py-2">
                <CheckCircle2 className="w-4 h-4 text-primary shrink-0" />
                <span>{feature}</span>
              </div>
            ))}
          </div>

          <Challenge number={8} title="Your Top 3">
            <p>Which 3 features are you most excited about? Which ones would save you the most time compared
            to building them yourself? Write down your top 3 and estimate how many hours each would take
            to build from scratch.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: The Grit Philosophy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Grit Philosophy</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every framework has a philosophy — a set of beliefs that guide every decision. Here are the six
            principles behind Grit:
          </p>

          <div className="space-y-5 mb-6">
            <div>
              <h3 className="font-semibold text-foreground mb-1">1. Convention over Configuration</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Sensible defaults for everything. File structure, naming, API format — all decided. But every
                convention is overridable. You get speed by default and flexibility when you need it.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-foreground mb-1">2. Code Generation over Runtime Magic</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                You own the code. When Grit generates a file, it{"'"}s real code you can read, modify, and debug.
                No hidden runtime behavior, no magic annotations, no black boxes. If something breaks, you can
                see exactly why.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-foreground mb-1">3. Batteries Included, Optionally Removable</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Everything is included by default. Don{"'"}t need email? Delete the email service file. Don{"'"}t
                need AI? Remove the AI handler. Features are modular — removing one doesn{"'"}t break the others.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-foreground mb-1">4. Beautiful by Default</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Every generated UI uses a premium dark theme with 100 pre-built components. Your app looks polished
                from minute one — not like a default Bootstrap template.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-foreground mb-1">5. Monorepo Native</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Shared types between Go and TypeScript. Change a model on the backend, run <Code>grit sync</Code>,
                and your frontend types update. One repo, one source of truth.
              </p>
            </div>
            <div>
              <h3 className="font-semibold text-foreground mb-1">6. Vibe Coding Ready</h3>
              <p className="text-muted-foreground text-sm leading-relaxed">
                Grit generates an AI skill file that teaches AI assistants (like Claude, Cursor, or Copilot) your
                project{"'"}s structure, conventions, and patterns. Your AI pair programmer understands your project
                from day one.
              </p>
            </div>
          </div>

          <Challenge number={9} title="Philosophy Reflection">
            <p>Which philosophy resonates most with you? Have you been burned by {'"'}runtime magic{'"'} or
            {'"'}configuration hell{'"'} before? Share a specific experience where one of these principles
            would have saved you time or frustration.</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Comparison with Other Frameworks ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Comparison with Other Frameworks</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            How does Grit stack up against the most popular frameworks? Here{"'"}s an honest comparison:
          </p>

          <div className="overflow-x-auto mb-6">
            <table className="w-full text-xs border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Feature</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Grit</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Laravel</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Rails</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Next.js</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">T3 Stack</th>
                  <th className="text-center px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Django</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">Go Performance</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">React Frontend</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">Admin Panel</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-yellow-400">Nova ($)</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">Code Generator</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">Built-in 2FA</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2.5 font-medium text-foreground">One-Cmd Deploy</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-yellow-400">Forge ($)</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-yellow-400">Vercel</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
                <tr>
                  <td className="px-3 py-2.5 font-medium text-foreground">Type-Safe API</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-green-400">Yes</td>
                  <td className="px-3 py-2.5 text-center text-red-400">No</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            Grit is not {'"'}better{'"'} than these frameworks. They each have strengths. Laravel has the largest
            ecosystem. Rails has the most mature conventions. Next.js has the best React SSR. Django has the best
            admin. Grit fills a specific gap: the <strong className="text-foreground">Go + React full-stack
            framework</strong> that doesn{"'"}t exist elsewhere.
          </Note>

          <Challenge number={10} title="Closest Competitor">
            <p>Look at the comparison table. Which framework is closest to Grit? What does Grit offer that
            it doesn{"'"}t? And what does that framework offer that Grit doesn{"'"}t?</p>
          </Challenge>
        </section>

        {/* ═══ Section 12: What You'll Learn in This Course Series ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You{"'"}ll Learn in This Course Series</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This course series has 9 courses that take you from zero to deploying a production-ready application.
            Here{"'"}s the roadmap:
          </p>

          <div className="space-y-3 mb-6">
            {[
              { num: 0, title: 'Introduction to Grit', desc: 'What Grit is, why it exists, and what you will learn (this course)' },
              { num: 1, title: 'Your First Grit App', desc: 'Install Grit, scaffold a project, explore the folder structure, run dev servers' },
              { num: 2, title: 'Code Generator Mastery', desc: 'Generate full-stack CRUD resources with one command' },
              { num: 3, title: 'Authentication & Authorization', desc: 'JWT, 2FA with TOTP, OAuth providers, role-based access' },
              { num: 4, title: 'Admin Panel Customization', desc: 'DataTable, FormBuilder, dashboard widgets, custom styles' },
              { num: 5, title: 'File Storage & Uploads', desc: 'S3-compatible storage, presigned URLs, image processing, MinIO' },
              { num: 6, title: 'Background Jobs & Email', desc: 'asynq job queue, cron scheduling, Resend email integration' },
              { num: 7, title: 'AI-Powered Features', desc: 'Vercel AI Gateway, streaming responses, multi-provider support' },
              { num: 8, title: 'Deploy to Production', desc: 'grit deploy, systemd, Caddy reverse proxy, Docker, auto-HTTPS' },
            ].map((course) => (
              <div key={course.num} className={`flex gap-3 items-start border border-border/30 rounded-lg px-4 py-3 ${course.num === 0 ? 'bg-primary/10 border-primary/30' : 'bg-muted/10'}`}>
                <span className={`text-xs font-mono font-bold shrink-0 mt-0.5 ${course.num === 0 ? 'text-primary' : 'text-muted-foreground'}`}>{course.num}.</span>
                <div>
                  <span className="font-semibold text-foreground text-sm">{course.title}</span>
                  <p className="text-muted-foreground text-xs mt-0.5">{course.desc}</p>
                </div>
              </div>
            ))}
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            By the end of this series, you{"'"}ll be able to build and deploy a production-ready full-stack
            application with Go and React — including authentication, admin panel, file uploads, background
            jobs, AI features, and automated deployment.
          </p>

          <Challenge number={11} title="Your Project Idea">
            <p>Write down one project idea you want to build with Grit. What architecture would you use
            (triple, double, single, API, or mobile)? What resources would you need to generate (e.g.,
            users, products, orders)? Which of the 18 built-in features would your project need?</p>
          </Challenge>
        </section>

        {/* ═══ Section 13: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Grit is a full-stack meta-framework combining Go (backend) with React (frontend)',
              'It solves the "Setup Tax" — giving you 2-4 weeks of infrastructure in 30 seconds',
              'Grit was created to let developers focus on the 20% that makes their app unique',
              'It offers 5 architecture modes: triple, double, single, API, and mobile',
              'You can choose between Next.js and TanStack Router for the frontend',
              'Every project ships with 18 production-ready features out of the box',
              'The philosophy centers on convention over configuration and code you own',
              'Type safety flows from Go to TypeScript to Zod with grit sync',
              'Grit fills the gap of a Go + React full-stack framework',
              'The course series covers 9 courses from introduction to production deployment',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={12} title="Prepare for Course 1">
            <p>Before starting Course 1, prepare your development environment. Visit these links and bookmark them:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Grit Docs: <a href="https://gritframework.dev/" target="_blank" rel="noreferrer" className="text-primary hover:underline">https://gritframework.dev/</a></li>
              <li>GitHub: <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer" className="text-primary hover:underline">https://github.com/MUKE-coder/grit</a></li>
              <li>YouTube: <a href="https://www.youtube.com/@GritFramework" target="_blank" rel="noreferrer" className="text-primary hover:underline">https://www.youtube.com/@GritFramework</a></li>
              <li>LinkedIn: <a href="https://www.linkedin.com/company/grit-framework" target="_blank" rel="noreferrer" className="text-primary hover:underline">https://www.linkedin.com/company/grit-framework</a></li>
            </ul>
            <p className="mt-2">Star the GitHub repo, subscribe to YouTube, follow on LinkedIn. You{"'"}re ready to start building.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web', label: 'Back to Grit Web' }}
            next={{ href: '/courses/grit-web/first-app', label: 'Next: Your First Grit App' }}
          />
        </div>
      </main>
    </div>
  )
}
