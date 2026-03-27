import Link from 'next/link'
import { ArrowLeft, ArrowRight, Terminal, CheckCircle2, AlertCircle, BookOpen } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Your First Grit App — Grit Web Course',
  description: 'Learn to install Grit, scaffold your first full-stack project, understand the project structure, and run all development servers.',
}

function CodeBlock({ filename, children }: { filename?: string; children: string }) {
  return (
    <div className="my-4 rounded-lg border border-border/40 overflow-hidden">
      {filename && (
        <div className="px-4 py-2 bg-muted/30 border-b border-border/40 text-xs font-mono text-muted-foreground">
          {filename}
        </div>
      )}
      <pre className="p-4 overflow-x-auto bg-[#0d1117] text-sm font-mono leading-relaxed">
        <code className="text-gray-300">{children}</code>
      </pre>
    </div>
  )
}

function Challenge({ number, title, children }: { number: number; title: string; children: React.ReactNode }) {
  return (
    <div className="my-6 rounded-lg border border-primary/30 bg-primary/5 p-5">
      <div className="flex items-center gap-2 mb-3">
        <span className="flex items-center justify-center h-6 w-6 rounded-full bg-primary/20 text-xs font-bold text-primary">
          {number}
        </span>
        <h4 className="font-semibold text-foreground text-sm">Challenge: {title}</h4>
      </div>
      <div className="text-sm text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

function Note({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-4 flex gap-3 rounded-lg border border-yellow-500/20 bg-yellow-500/5 p-4 text-sm">
      <AlertCircle className="h-5 w-5 text-yellow-500 shrink-0 mt-0.5" />
      <div className="text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

function Tip({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-4 flex gap-3 rounded-lg border border-primary/20 bg-primary/5 p-4 text-sm">
      <CheckCircle2 className="h-5 w-5 text-primary shrink-0 mt-0.5" />
      <div className="text-muted-foreground leading-relaxed">{children}</div>
    </div>
  )
}

export default function FirstAppCourse() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-web" className="hover:text-foreground transition-colors">Grit Web</Link>
          <span>/</span>
          <span className="text-foreground">Your First Grit App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 1 of 8</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Your First Grit App
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will install Grit, create your first full-stack project,
            understand every folder and file it generates, and run all the development servers.
            By the end, you will have a working app with authentication, admin panel, and database.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 1: What is Grit?                                  */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Grit?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit is a <strong className="text-foreground">full-stack framework</strong> that combines Go (for the backend) with React (for the frontend).
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Grit uses <strong className="text-foreground">Go</strong> with the Gin web framework and GORM ORM for the API</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit uses <strong className="text-foreground">React</strong> with Next.js or TanStack Router for the frontend</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit generates a complete project with authentication, admin panel, and more</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit has a CLI that helps you generate code, run migrations, and deploy</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit is open source and free to use (MIT license)</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With Grit, you write one command and get a production-ready app. No more spending weeks setting up boilerplate.
          </p>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 2: Prerequisites                                  */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Prerequisites</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before installing Grit, you need four tools on your computer:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Tool</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Version</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Check command</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Go</td>
                  <td className="px-4 py-3"><code className="text-primary">1.21+</code></td>
                  <td className="px-4 py-3">Runs the backend API server</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">go version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Node.js</td>
                  <td className="px-4 py-3"><code className="text-primary">18+</code></td>
                  <td className="px-4 py-3">Runs the frontend apps</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">node --version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">pnpm</td>
                  <td className="px-4 py-3"><code className="text-primary">8+</code></td>
                  <td className="px-4 py-3">Installs JavaScript packages</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">pnpm --version</code></td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Docker</td>
                  <td className="px-4 py-3"><code className="text-primary">20+</code></td>
                  <td className="px-4 py-3">Runs PostgreSQL, Redis, MinIO</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">docker --version</code></td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Open your terminal and run each check command. If any tool is missing, install it from its official website.
          </p>

          <Challenge number={1} title="Check Your Tools">
            <p>Open your terminal and run all four check commands above. Write down the version of each tool you have installed. All four must return a version number.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 3: Install Grit                                   */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Install Grit</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit is installed using Go{"'"}s <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">go install</code> command. This downloads the Grit CLI binary and puts it in your Go bin directory.
          </p>

          <CodeBlock filename="Terminal">
{`go install github.com/MUKE-coder/grit/v2/cmd/grit@latest`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Command Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">go install</code> tells Go to download and compile a package</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">github.com/MUKE-coder/grit/v2/cmd/grit</code> is the path to the Grit CLI</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">@latest</code> means get the newest version</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After installation, verify it works:
          </p>

          <CodeBlock filename="Terminal">
{`grit version
# Output: grit version 3.5.0`}
          </CodeBlock>

          <Note>
            If you get {'"'}command not found{'"'}, your Go bin directory is not in your PATH.
            Add <code className="text-xs bg-muted/30 px-1 rounded">export PATH=$PATH:$(go env GOPATH)/bin</code> to your shell profile (~/.bashrc or ~/.zshrc).
          </Note>

          <Challenge number={2} title="Install Grit">
            <p>Run the install command, then run <code className="text-primary bg-muted/30 px-1 rounded">grit version</code> to verify. You should see {"\""}grit version 3.5.0{"\""} (or newer).</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 4: Create a Project                               */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Create Your First Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">grit new</code> command creates a new project.
            When you run it without flags, it enters <strong className="text-foreground">interactive mode</strong> and asks you questions.
          </p>

          <CodeBlock filename="Terminal">
{`grit new myapp`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">The Interactive Prompts</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit will ask you three questions:
          </p>

          <h4 className="text-base font-semibold text-foreground mb-2">Question 1: Choose your architecture</h4>

          <CodeBlock filename="Terminal">
{`? Select architecture:
  > Triple — Web + Admin + API (Turborepo)
    Double — Web + API (Turborepo)
    Single — Go API + embedded React SPA (one binary)
    API Only — Go API (no frontend)
    Mobile — API + Expo (React Native)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For this course, select <strong className="text-foreground">Triple</strong>. This gives you the full experience:
            a web app, an admin panel, and a Go API — all in one project.
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Architecture</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What you get</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Best for</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Triple</td>
                  <td className="px-4 py-3">Web + Admin + API</td>
                  <td className="px-4 py-3">SaaS apps, platforms</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Double</td>
                  <td className="px-4 py-3">Web + API</td>
                  <td className="px-4 py-3">Simple apps, blogs</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Single</td>
                  <td className="px-4 py-3">One Go binary</td>
                  <td className="px-4 py-3">Microservices, internal tools</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">API Only</td>
                  <td className="px-4 py-3">Go backend only</td>
                  <td className="px-4 py-3">Mobile backends, headless APIs</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Mobile</td>
                  <td className="px-4 py-3">API + Expo</td>
                  <td className="px-4 py-3">React Native apps</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h4 className="text-base font-semibold text-foreground mb-2">Question 2: Choose your frontend</h4>

          <CodeBlock filename="Terminal">
{`? Select frontend:
  > Next.js — SSR, SEO, App Router
    TanStack Router — Vite, fast builds, small bundle (SPA)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Select <strong className="text-foreground">Next.js</strong> for this course.
            Next.js gives you server-side rendering and is better for SEO.
            TanStack Router is faster but is a client-side SPA.
          </p>

          <h4 className="text-base font-semibold text-foreground mb-2">Question 3: Choose your admin style</h4>

          <CodeBlock filename="Terminal">
{`? Select admin panel style:
  > Default — Clean dark theme
    Modern — Gradient accents
    Minimal — Ultra clean
    Glass — Glassmorphism`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Select <strong className="text-foreground">Default</strong>. You can always change this later.
          </p>

          <Tip>
            You can skip the prompts by passing flags directly:
            <code className="block mt-2 text-xs bg-muted/30 px-2 py-1 rounded">grit new myapp --triple --next --style default</code>
          </Tip>

          <Challenge number={3} title="Create a Project">
            <p>Run <code className="text-primary bg-muted/30 px-1 rounded">grit new myapp</code> and select Triple, Next.js, and Default style. Wait for it to finish.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 5: Project Structure                              */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Understanding the Project Structure</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit created a folder called <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">myapp/</code>. Let{"'"}s look inside:
          </p>

          <CodeBlock filename="Project Structure">
{`myapp/
├── apps/
│   ├── api/           ← Go backend (Gin + GORM)
│   ├── web/           ← Next.js frontend
│   └── admin/         ← Next.js admin panel
├── packages/
│   └── shared/        ← Zod schemas + TypeScript types
├── docker-compose.yml ← PostgreSQL, Redis, MinIO, Mailhog
├── turbo.json         ← Monorepo task runner config
├── package.json       ← Root package.json
└── .env               ← Environment variables`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Structure Explained</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This is a <strong className="text-foreground">monorepo</strong> — multiple apps in one repository.
            Turborepo manages them so they can share code and run together.
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Folder</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What{"'"}s inside</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Language</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">apps/api/</td>
                  <td className="px-4 py-3">Go REST API — models, handlers, services, middleware, routes</td>
                  <td className="px-4 py-3">Go</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">apps/web/</td>
                  <td className="px-4 py-3">Main frontend — pages, components, hooks</td>
                  <td className="px-4 py-3">TypeScript + React</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">apps/admin/</td>
                  <td className="px-4 py-3">Admin dashboard — resource management, data tables</td>
                  <td className="px-4 py-3">TypeScript + React</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">packages/shared/</td>
                  <td className="px-4 py-3">Shared Zod validation schemas and TypeScript types</td>
                  <td className="px-4 py-3">TypeScript</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">.env</td>
                  <td className="px-4 py-3">All configuration — database URL, API keys, secrets</td>
                  <td className="px-4 py-3">Key=Value</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={4} title="Explore the Folders">
            <p>Open the <code className="text-primary bg-muted/30 px-1 rounded">myapp</code> folder in your code editor (VS Code recommended). Look at each folder listed above. Can you find the Go API entry point at <code className="text-primary bg-muted/30 px-1 rounded">apps/api/cmd/server/main.go</code>?</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 6: The Go API                                     */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Inside the Go API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Go API lives in <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">apps/api/</code>. Here{"'"}s what{"'"}s inside:
          </p>

          <CodeBlock filename="apps/api/">
{`apps/api/
├── cmd/
│   ├── server/main.go    ← Entry point (starts the API)
│   ├── migrate/main.go   ← Database migration runner
│   └── seed/main.go      ← Database seeder
└── internal/
    ├── config/            ← Reads .env variables
    ├── database/          ← Connects to PostgreSQL
    ├── models/            ← Database tables (User, Upload, Blog...)
    ├── handlers/          ← HTTP request handlers
    ├── services/          ← Business logic
    ├── middleware/         ← Auth, CORS, logging, rate limiting
    ├── routes/            ← Route definitions
    ├── cache/             ← Redis caching
    ├── storage/           ← File uploads (S3/MinIO)
    ├── mail/              ← Email service (Resend)
    ├── jobs/              ← Background jobs (asynq)
    ├── cron/              ← Scheduled tasks
    ├── ai/                ← AI service (Vercel AI Gateway)
    └── totp/              ← Two-factor authentication`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">How the API Works</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user makes a request (like logging in), this is what happens:
          </p>

          <CodeBlock filename="Request Flow">
{`Browser → Request → Middleware → Handler → Service → Database
                      ↓                                  ↓
                 (Auth check)                      (GORM query)
                      ↓                                  ↓
              Handler ← Service ← Database Response
                      ↓
              JSON Response → Browser`}
          </CodeBlock>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Middleware</strong> runs first — checks authentication, logs the request, checks rate limits</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Handler</strong> receives the request — validates input, calls the service</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Service</strong> contains business logic — talks to the database, processes data</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Database</strong> stores everything — users, posts, files, sessions</li>
          </ul>

          <Challenge number={5} title="Read the Entry Point">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">apps/api/cmd/server/main.go</code> and read through it. Can you find where the database connects? Where the router is set up? Where the server starts listening?</p>
          </Challenge>

          <Challenge number={6} title="Find the User Model">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">apps/api/internal/models/user.go</code>. What fields does the User model have? Can you identify the GORM tags (<code className="text-primary bg-muted/30 px-1 rounded">gorm:&quot;...&quot;</code>) and JSON tags (<code className="text-primary bg-muted/30 px-1 rounded">json:&quot;...&quot;</code>)?</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 7: Start Docker Services                          */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start Docker Services</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your app needs a database (PostgreSQL), cache (Redis), file storage (MinIO), and a mail catcher (Mailhog).
            All of these run in Docker containers.
          </p>

          <CodeBlock filename="Terminal">
{`cd myapp
docker compose up -d`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Command Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">cd myapp</code> — moves into your project folder</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">docker compose up -d</code> — starts all services in the background</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">-d</code> means {"\""} detached{"\""} — they run in the background so you get your terminal back</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Docker starts four services:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Service</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Port</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">PostgreSQL</td>
                  <td className="px-4 py-3">5432</td>
                  <td className="px-4 py-3">Main database — stores users, posts, everything</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Redis</td>
                  <td className="px-4 py-3">6379</td>
                  <td className="px-4 py-3">Cache and job queue</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">MinIO</td>
                  <td className="px-4 py-3">9000 / 9001</td>
                  <td className="px-4 py-3">File storage (S3-compatible, for uploads)</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Mailhog</td>
                  <td className="px-4 py-3">8025</td>
                  <td className="px-4 py-3">Catches all emails locally (for testing)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You can check if all services are running:
          </p>

          <CodeBlock filename="Terminal">
{`docker compose ps
# All 4 services should show "running"`}
          </CodeBlock>

          <Challenge number={7} title="Start Docker">
            <p>Run <code className="text-primary bg-muted/30 px-1 rounded">docker compose up -d</code> inside your project. Then run <code className="text-primary bg-muted/30 px-1 rounded">docker compose ps</code> to verify all 4 services are running.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 8: Start the App                                  */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start the Development Servers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now start the frontend apps. First, install JavaScript dependencies:
          </p>

          <CodeBlock filename="Terminal">
{`pnpm install`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Then start all apps at once:
          </p>

          <CodeBlock filename="Terminal">
{`pnpm dev`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This starts three servers simultaneously:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">App</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">URL</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What you see</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Web App</td>
                  <td className="px-4 py-3"><code className="text-primary">http://localhost:3000</code></td>
                  <td className="px-4 py-3">Main frontend with login/register/dashboard</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Admin Panel</td>
                  <td className="px-4 py-3"><code className="text-primary">http://localhost:3001</code></td>
                  <td className="px-4 py-3">Admin dashboard with user management</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Go API</td>
                  <td className="px-4 py-3"><code className="text-primary">http://localhost:8080</code></td>
                  <td className="px-4 py-3">REST API (JSON responses)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            The Go API starts automatically with <code className="text-xs bg-muted/30 px-1 rounded">pnpm dev</code> if you are running the Turborepo setup. If it does not start, open a separate terminal and run: <code className="text-xs bg-muted/30 px-1 rounded">cd apps/api && go run cmd/server/main.go</code>
          </Note>

          <Challenge number={8} title="Start Everything">
            <p>Run <code className="text-primary bg-muted/30 px-1 rounded">pnpm install</code> then <code className="text-primary bg-muted/30 px-1 rounded">pnpm dev</code>. Open all three URLs in your browser. You should see the web app, admin panel, and API welcome message.</p>
          </Challenge>

          <Challenge number={9} title="Register an Account">
            <p>Go to <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:3000</code> and click {"\""}Register{"\""}.
            Create an account with your email and password. After registering, you should be redirected to the dashboard.</p>
          </Challenge>

          <Challenge number={10} title="Log Into Admin">
            <p>Go to <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:3001</code> and log in with the same account.
            You should see the admin dashboard with stats cards and a sidebar menu.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 9: Built-in Tools                                 */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Built-in Tools</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your API comes with four built-in tools. Each one has a web interface:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Tool</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">URL</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">GORM Studio</td>
                  <td className="px-4 py-3"><code className="text-primary">localhost:8080/studio</code></td>
                  <td className="px-4 py-3">Browse database tables, view and edit records, run SQL</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">API Docs</td>
                  <td className="px-4 py-3"><code className="text-primary">localhost:8080/docs</code></td>
                  <td className="px-4 py-3">Interactive API documentation — test every endpoint</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Pulse</td>
                  <td className="px-4 py-3"><code className="text-primary">localhost:8080/pulse/ui</code></td>
                  <td className="px-4 py-3">Request tracing, metrics, database monitoring</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Sentinel</td>
                  <td className="px-4 py-3"><code className="text-primary">localhost:8080/sentinel/ui</code></td>
                  <td className="px-4 py-3">Security dashboard — rate limits, blocked IPs, threats</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Mailhog</td>
                  <td className="px-4 py-3"><code className="text-primary">localhost:8025</code></td>
                  <td className="px-4 py-3">Email inbox — catches all emails sent during development</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={11} title="Visit GORM Studio">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:8080/studio</code> in your browser. Log in (default: admin/studio). Find the {"\""}users{"\""} table and look for the account you just registered.</p>
          </Challenge>

          <Challenge number={12} title="Test the API Docs">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:8080/docs</code>. Find the {"\""}POST /api/auth/login{"\""} endpoint. Try logging in with your email and password using the interactive form. You should get a JSON response with a token.</p>
          </Challenge>

          <Challenge number={13} title="Check Pulse">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:8080/pulse/ui</code>. Refresh your web app a few times, then go back to Pulse. Can you see the requests being logged?</p>
          </Challenge>

          <Challenge number={14} title="Check Mailhog">
            <p>Open <code className="text-primary bg-muted/30 px-1 rounded">http://localhost:8025</code>. Is there a welcome email from when you registered? If yes, open it and see the HTML template.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 10: The .env File                                  */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The .env File</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project has a <code className="text-primary text-sm bg-muted/30 px-1.5 py-0.5 rounded">.env</code> file at the root. This file contains all your configuration:
          </p>

          <CodeBlock filename=".env (partial)">
{`# Core
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable
JWT_SECRET=your-secret-key
REDIS_URL=redis://localhost:6379

# Storage (MinIO for local development)
STORAGE_DRIVER=minio
STORAGE_ENDPOINT=localhost:9000

# Email
RESEND_API_KEY=
MAIL_FROM=noreply@localhost

# AI (Vercel AI Gateway)
AI_GATEWAY_API_KEY=
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Key Variables Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">DATABASE_URL</code> — connection string for PostgreSQL</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">JWT_SECRET</code> — secret key for signing authentication tokens (change this in production!)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">REDIS_URL</code> — connection to Redis for caching and job queues</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">AI_GATEWAY_API_KEY</code> — optional, for AI features (get from vercel.com/ai-gateway)</li>
          </ul>

          <Note>
            Never commit your <code className="text-xs bg-muted/30 px-1 rounded">.env</code> file to Git. It contains secrets. Grit adds it to <code className="text-xs bg-muted/30 px-1 rounded">.gitignore</code> automatically.
          </Note>

          <Challenge number={15} title="Read the .env File">
            <p>Open the <code className="text-primary bg-muted/30 px-1 rounded">.env</code> file in your editor. Find the <code className="text-primary bg-muted/30 px-1 rounded">DATABASE_URL</code> — what database name is it using? Find the <code className="text-primary bg-muted/30 px-1 rounded">APP_PORT</code> — what port does the API run on?</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 11: Useful CLI Commands                           */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Useful CLI Commands</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are commands you will use every day while developing with Grit:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Command</th>
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">What it does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit routes</td>
                  <td className="px-4 py-3">Lists all your API endpoints in a table</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit migrate</td>
                  <td className="px-4 py-3">Runs database migrations (creates/updates tables)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit seed</td>
                  <td className="px-4 py-3">Fills the database with demo data</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit studio</td>
                  <td className="px-4 py-3">Opens the database browser</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit version</td>
                  <td className="px-4 py-3">Shows the installed Grit version</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={16} title="List Your Routes">
            <p>Run <code className="text-primary bg-muted/30 px-1 rounded">grit routes</code> in your project folder. How many routes does your app have? Can you find the login endpoint? The register endpoint?</p>
          </Challenge>

          <Challenge number={17} title="Stop and Restart">
            <p>Stop your dev servers (Ctrl+C). Stop Docker with <code className="text-primary bg-muted/30 px-1 rounded">docker compose down</code>. Then start everything again: <code className="text-primary bg-muted/30 px-1 rounded">docker compose up -d</code> and <code className="text-primary bg-muted/30 px-1 rounded">pnpm dev</code>. Verify everything still works.</p>
          </Challenge>
        </section>

        {/* ═══════════════════════════════════════════════════════════ */}
        {/* SECTION 12: Summary                                       */}
        {/* ═══════════════════════════════════════════════════════════ */}

        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            In this course you learned:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to install Grit with <code className="text-xs bg-muted/30 px-1 rounded">go install</code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to scaffold a project with <code className="text-xs bg-muted/30 px-1 rounded">grit new</code> (interactive and with flags)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The monorepo project structure (apps/api, apps/web, apps/admin, packages/shared)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How the Go API is organized (models → services → handlers → routes)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to start Docker services and development servers</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The 5 built-in tools (GORM Studio, API Docs, Pulse, Sentinel, Mailhog)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How the .env file works</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Essential CLI commands (grit routes, grit migrate, grit studio)</li>
          </ul>

          <Challenge number={18} title="Final Challenge: Start From Scratch">
            <p>Delete the <code className="text-primary bg-muted/30 px-1 rounded">myapp</code> folder completely. Now create a new project called <code className="text-primary bg-muted/30 px-1 rounded">bookstore</code> using Triple architecture with TanStack Router (Vite) instead of Next.js. Start everything and verify it works. Notice any differences from the Next.js version?</p>
          </Challenge>
        </section>

        {/* Navigation */}
        <div className="flex items-center justify-between pt-8 border-t border-border/40">
          <Link
            href="/courses/grit-web"
            className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to Grit Web
          </Link>
          <Link
            href="/courses/grit-web/code-generator"
            className="flex items-center gap-2 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
          >
            Next: Code Generator Mastery
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </main>
    </div>
  )
}
