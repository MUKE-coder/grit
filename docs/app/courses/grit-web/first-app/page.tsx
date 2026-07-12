import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Your First Grit App — Grit Web Course',
  description: 'Learn to install Grit, scaffold your first full-stack project, understand the project structure, and run all development servers.',
}

export default function FirstAppCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

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
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">18 challenges</span>
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

        {/* ═══ What is Grit? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Grit?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit is a <strong className="text-foreground">full-stack framework</strong> that combines Go (for the backend) with React (for the frontend).
          </p>

          <Definition term="Framework">
            A framework is a pre-built foundation that provides structure, tools, and conventions so you
            don{"'"}t have to build everything from scratch. Think of it like a house blueprint — it gives you
            the walls, plumbing, and wiring, and you add the furniture.
          </Definition>

          <Definition term="Full-stack">
            Full-stack means both the <strong className="text-foreground">backend</strong> (server, database, API — the part users don{"'"}t see)
            and the <strong className="text-foreground">frontend</strong> (the website or app the user interacts with).
            Grit handles both.
          </Definition>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Grit uses <strong className="text-foreground">Go</strong> with the <a href="https://gin-gonic.com" target="_blank" rel="noreferrer" className="text-primary hover:underline">Gin</a> web framework and <a href="https://gorm.io" target="_blank" rel="noreferrer" className="text-primary hover:underline">GORM</a> ORM for the API</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit uses <strong className="text-foreground">React</strong> with <a href="https://nextjs.org" target="_blank" rel="noreferrer" className="text-primary hover:underline">Next.js</a> or <a href="https://tanstack.com/router" target="_blank" rel="noreferrer" className="text-primary hover:underline">TanStack Router</a> for the frontend</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit generates a complete project with authentication, admin panel, and more</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit has a <strong className="text-foreground">CLI</strong> (Command Line Interface) that helps you generate code, run migrations, and deploy</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Grit is <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer" className="text-primary hover:underline">open source</a> and free to use (MIT license)</li>
          </ul>

          <Definition term="CLI (Command Line Interface)">
            A program you run in your terminal (like Terminal on Mac or Command Prompt on Windows).
            Instead of clicking buttons in a graphical interface, you type commands like <Code>grit new myapp</Code>.
          </Definition>
        </section>

        {/* ═══ Prerequisites ═══ */}
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
                  <th className="text-left px-4 py-3 font-semibold text-foreground border-b border-border/40">Check</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><a href="https://go.dev/dl/" target="_blank" rel="noreferrer" className="text-primary hover:underline">Go</a></td>
                  <td className="px-4 py-3"><Code>1.21+</Code></td>
                  <td className="px-4 py-3">Runs the backend API server</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">go version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><a href="https://nodejs.org" target="_blank" rel="noreferrer" className="text-primary hover:underline">Node.js</a></td>
                  <td className="px-4 py-3"><Code>18+</Code></td>
                  <td className="px-4 py-3">Runs the frontend apps</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">node --version</code></td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground"><a href="https://pnpm.io/installation" target="_blank" rel="noreferrer" className="text-primary hover:underline">pnpm</a></td>
                  <td className="px-4 py-3"><Code>8+</Code></td>
                  <td className="px-4 py-3">Installs JavaScript packages</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">pnpm --version</code></td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground"><a href="https://docker.com/get-started" target="_blank" rel="noreferrer" className="text-primary hover:underline">Docker</a></td>
                  <td className="px-4 py-3"><Code>20+</Code></td>
                  <td className="px-4 py-3">Runs PostgreSQL, Redis, MinIO</td>
                  <td className="px-4 py-3"><code className="text-xs bg-muted/30 px-1.5 py-0.5 rounded">docker --version</code></td>
                </tr>
              </tbody>
            </table>
          </div>

          <Definition term="pnpm">
            A fast, disk-efficient package manager for JavaScript. It{"'"}s like npm but faster and uses less disk space.
            Grit uses pnpm instead of npm because it handles monorepos better. Install it with: <Code>npm install -g pnpm</Code>
          </Definition>

          <Definition term="Docker">
            Docker runs applications in isolated containers. Instead of installing PostgreSQL, Redis, and MinIO
            directly on your computer, Docker runs them in lightweight virtual environments. This keeps your
            system clean and makes setup identical across all operating systems.
          </Definition>

          <Challenge number={1} title="Check Your Tools">
            <p>Open your terminal and run all four check commands from the table above. Write down the version of each tool. All four must return a version number before you continue.</p>
          </Challenge>
        </section>

        {/* ═══ Install Grit ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Install Grit</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The easiest way: a one-line install script. It detects an existing
            install and runs <Code>grit update</Code>; otherwise it downloads
            the release binary for your OS / arch. If you have Go installed,
            <Code>go install</Code> works too (shown below the script).
          </p>

          <CodeBlock filename="Terminal">
{`# macOS / Linux
curl -fsSL https://gritframework.dev/install.sh | sh

# Windows (PowerShell)
iwr -useb https://gritframework.dev/install.ps1 | iex

# Power users with Go installed (alternative)
go install github.com/MUKE-coder/grit/v3/cmd/grit@latest`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Command Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>go install</Code> tells Go to download and compile a package</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <code className="text-xs bg-muted/30 px-1 rounded">github.com/MUKE-coder/grit/v3/cmd/grit</code> is the path to the Grit CLI on <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noreferrer" className="text-primary hover:underline">GitHub</a></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>@latest</Code> means get the newest version</li>
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
            Add <Code>export PATH=$PATH:$(go env GOPATH)/bin</Code> to your shell profile (~/.bashrc or ~/.zshrc), then restart your terminal.
          </Note>

          <Challenge number={2} title="Install Grit">
            <p>Run the install command, then run <Code>grit version</Code> to verify. You should see {'"'}grit version 3.5.0{'"'} (or newer).</p>
          </Challenge>
        </section>

        {/* ═══ Create a Project ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Create Your First Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>grit new</Code> command creates a new project.
            When you run it without flags, it enters <strong className="text-foreground">interactive mode</strong> and asks you questions.
          </p>

          <CodeBlock filename="Terminal">
{`grit new myapp`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Question 1: Choose your architecture</h3>

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

          <Definition term="Architecture">
            The structure of your project — how many apps it contains and how they{"'"}re organized.
            Triple means three apps (web + admin + API). Double means two (web + API). Single means one Go binary
            that serves everything.
          </Definition>

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

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Question 2: Choose your frontend</h3>

          <CodeBlock filename="Terminal">
{`? Select frontend:
  > Next.js — SSR, SEO, App Router
    TanStack Router — Vite, fast builds, small bundle (SPA)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Select <strong className="text-foreground">Next.js</strong> for this course.
          </p>

          <Definition term="SSR (Server-Side Rendering)">
            The server generates the HTML before sending it to the browser. This is better for SEO because
            search engines can read the content immediately. Next.js does this by default.
          </Definition>

          <Definition term="SPA (Single Page Application)">
            The browser downloads one HTML file and JavaScript handles all navigation. Faster page
            transitions but harder for SEO. TanStack Router with Vite creates SPAs.
          </Definition>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Question 3: Choose your admin style</h3>

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
            <code className="block mt-2 text-xs bg-muted/30 px-2 py-1 rounded font-mono">grit new myapp --triple --next --style default</code>
          </Tip>

          <Challenge number={3} title="Create a Project">
            <p>Run <Code>grit new myapp</Code> and select Triple, Next.js, and Default style. Wait for it to finish.</p>
          </Challenge>

          <Challenge number={4} title="Create with Flags">
            <p>Delete the myapp folder (<Code>rm -rf myapp</Code>). Now create the same project using flags instead of interactive mode: <Code>grit new myapp --triple --next</Code>. Same result, no prompts.</p>
          </Challenge>
        </section>

        {/* ═══ Project Structure ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Understanding the Project Structure</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit created a folder called <Code>myapp/</Code>. Let{"'"}s look inside:
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
├── pnpm-workspace.yaml ← Workspace definition
└── .env               ← Environment variables`}
          </CodeBlock>

          <Definition term="Monorepo">
            A single Git repository that contains multiple projects (apps and libraries). Instead of having
            separate repos for your API and frontend, everything lives together. This makes it easy to share
            code and keep things in sync. Grit uses <a href="https://turbo.build" target="_blank" rel="noreferrer" className="text-primary hover:underline">Turborepo</a> to
            manage the monorepo.
          </Definition>

          <Definition term="ORM (Object-Relational Mapping)">
            A tool that lets you interact with your database using Go structs instead of writing raw SQL.
            <a href="https://gorm.io" target="_blank" rel="noreferrer" className="text-primary hover:underline"> GORM</a> is Go{"'"}s most popular ORM.
            You define a <Code>User</Code> struct and GORM creates the <Code>users</Code> table automatically.
          </Definition>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Structure Explained</h3>

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
                  <td className="px-4 py-3">Shared <a href="https://zod.dev" target="_blank" rel="noreferrer" className="text-primary hover:underline">Zod</a> validation schemas and TypeScript types</td>
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

          <Challenge number={5} title="Explore the Folders">
            <p>Open the <Code>myapp</Code> folder in VS Code (<Code>code myapp</Code>). Look at each folder in the table above. Can you find the Go API entry point at <Code>apps/api/cmd/server/main.go</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Inside the Go API ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Inside the Go API</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Go API is the brain of your application. It lives in <Code>apps/api/</Code>:
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
    ├── ai/                ← AI service (Vercel AI Gateway)
    └── totp/              ← Two-factor authentication`}
          </CodeBlock>

          <Definition term="API (Application Programming Interface)">
            A set of endpoints (URLs) that your frontend calls to get or send data. For example,
            <Code>POST /api/auth/login</Code> logs a user in and returns a token.
            The frontend sends JSON, the API processes it, and returns JSON back.
          </Definition>

          <Definition term="REST API">
            A style of API that uses HTTP methods (GET, POST, PUT, DELETE) to perform actions.
            GET reads data, POST creates data, PUT updates data, DELETE removes data.
            Grit generates REST APIs automatically.
          </Definition>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">How a Request Flows Through the API</h3>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user does something (like clicking {'"'}Login{'"'}), here{"'"}s what happens:
          </p>

          <CodeBlock filename="Request Flow">
{`Browser → HTTP Request → Middleware → Handler → Service → Database
                            ↓                                 ↓
                      (Auth check)                      (GORM query)
                            ↓                                 ↓
                   Handler ← Service ← Database Response
                            ↓
                   JSON Response → Browser`}
          </CodeBlock>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Middleware</strong> runs first — checks if the user is logged in, logs the request, applies rate limits</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Handler</strong> receives the request — validates input data, then calls the service</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Service</strong> contains business logic — talks to the database, processes data, applies rules</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Model</strong> defines the database table — what columns exist, their types, and relationships</li>
          </ul>

          <Definition term="Middleware">
            Code that runs <em>before</em> your handler. Like a security guard at a building entrance — it checks
            your credentials before letting you in. Grit{"'"}s auth middleware checks the JWT token on every protected request.
          </Definition>

          <Challenge number={6} title="Read the Entry Point">
            <p>Open <Code>apps/api/cmd/server/main.go</Code> and read through it. Can you find where: (1) the database connects, (2) the router is set up, and (3) the server starts listening?</p>
          </Challenge>

          <Challenge number={7} title="Find the User Model">
            <p>Open <Code>apps/api/internal/models/user.go</Code>. List all the fields the User model has. Can you identify the GORM tags (<Code>gorm:&quot;...&quot;</Code>) and JSON tags (<Code>json:&quot;...&quot;</Code>)?</p>
          </Challenge>
        </section>

        {/* ═══ Start Docker ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start Docker Services</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your app needs a database (PostgreSQL), cache (Redis), file storage (MinIO), and a mail catcher (Mailhog).
            All of these run inside Docker containers.
          </p>

          <Definition term="Container">
            A lightweight, isolated environment that runs an application. Think of it like a tiny virtual computer
            dedicated to running one thing (like PostgreSQL). Docker manages these containers for you.
          </Definition>

          <CodeBlock filename="Terminal">
{`cd myapp
docker compose up -d`}
          </CodeBlock>

          <h3 className="text-lg font-semibold text-foreground mb-3 mt-6">Command Explained</h3>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>cd myapp</Code> — moves into your project folder</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>docker compose up</Code> — reads <Code>docker-compose.yml</Code> and starts all services</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>-d</Code> means {'"'}detached{'"'} — they run in the background so you get your terminal back</li>
          </ul>

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
                  <td className="px-4 py-3">5434</td>
                  <td className="px-4 py-3">Main database — stores users, posts, everything</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Redis</td>
                  <td className="px-4 py-3">6380</td>
                  <td className="px-4 py-3">Cache (fast temporary storage) and job queue</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">MinIO</td>
                  <td className="px-4 py-3">9002 / 9003</td>
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
            Check if everything is running:
          </p>

          <CodeBlock filename="Terminal">
{`docker compose ps
# All 4 services should show "running"`}
          </CodeBlock>

          <Challenge number={8} title="Start Docker Services">
            <p>Run <Code>docker compose up -d</Code> inside your project. Then run <Code>docker compose ps</Code> to verify all 4 services are running. If any service failed, run <Code>docker compose logs</Code> to see what went wrong.</p>
          </Challenge>
        </section>

        {/* ═══ Start the App ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start the Development Servers</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Install JavaScript dependencies first:
          </p>

          <CodeBlock filename="Terminal">
{`pnpm install`}
          </CodeBlock>

          <Definition term="Dependencies">
            External libraries your project uses. For example, React, Tailwind CSS, and Zod are dependencies.
            <Code>pnpm install</Code> reads the <Code>package.json</Code> file and downloads everything your project needs.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Then start all apps at once:
          </p>

          <CodeBlock filename="Terminal">
{`grit start`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The bare <Code>grit start</Code> boots all three servers at once. You can also start them
            individually: <Code>grit start server</Code> for just the Go API, or <Code>grit start client</Code> for
            just the frontend apps. Under the hood it uses <a href="https://turbo.build" target="_blank" rel="noreferrer" className="text-primary hover:underline">Turborepo</a>:
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
                  <td className="px-4 py-3"><Code>http://localhost:3000</Code></td>
                  <td className="px-4 py-3">Main frontend with login/register/dashboard</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Admin Panel</td>
                  <td className="px-4 py-3"><Code>http://localhost:3001</Code></td>
                  <td className="px-4 py-3">Admin dashboard with user management</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Go API</td>
                  <td className="px-4 py-3"><Code>http://localhost:8080</Code></td>
                  <td className="px-4 py-3">REST API (JSON responses)</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Note>
            If the Go API doesn{"'"}t come up with <Code>grit start</Code>, open a separate terminal and run just the API with <Code>grit start server</Code>.
          </Note>

          <Challenge number={9} title="Start Everything">
            <p>Run <Code>pnpm install</Code> then <Code>grit start</Code>. Open all three URLs in your browser.</p>
          </Challenge>

          <Challenge number={10} title="Register an Account">
            <p>Go to <Code>http://localhost:3000</Code> and click {'"'}Register{'"'}. Create an account with your email and password. You should be redirected to the dashboard after registration.</p>
          </Challenge>

          <Challenge number={11} title="Log Into the Admin Panel">
            <p>Go to <Code>http://localhost:3001</Code> and log in with the same credentials. Explore the sidebar — click on Users, Dashboard, and System pages.</p>
          </Challenge>
        </section>

        {/* ═══ Built-in Tools ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Built-in Tools</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Your API comes with five built-in tools, each with its own web interface:
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
                  <td className="px-4 py-3"><Code>localhost:8080/studio</Code></td>
                  <td className="px-4 py-3">Browse database tables, view/edit records, run SQL</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">API Docs</td>
                  <td className="px-4 py-3"><Code>localhost:8080/docs</Code></td>
                  <td className="px-4 py-3">Interactive API documentation — test every endpoint</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Pulse</td>
                  <td className="px-4 py-3"><Code>localhost:8080/pulse/ui</Code></td>
                  <td className="px-4 py-3">Request tracing, performance metrics, database monitoring</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-medium text-foreground">Sentinel</td>
                  <td className="px-4 py-3"><Code>localhost:8080/sentinel/ui</Code></td>
                  <td className="px-4 py-3">Security dashboard — rate limits, blocked IPs, threats</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-medium text-foreground">Mailhog</td>
                  <td className="px-4 py-3"><Code>localhost:8025</Code></td>
                  <td className="px-4 py-3">Email inbox — catches all emails sent during development</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={12} title="Visit GORM Studio">
            <p>Open <Code>http://localhost:8080/studio</Code>. Log in (default: admin/studio). Find the {'"'}users{'"'} table and look for the account you just registered.</p>
          </Challenge>

          <Challenge number={13} title="Test the API Docs">
            <p>Open <Code>http://localhost:8080/docs</Code>. Find the {'"'}POST /api/auth/login{'"'} endpoint. Click it, fill in your email and password, and click {'"'}Send{'"'}. You should get a JSON response with a token.</p>
          </Challenge>

          <Challenge number={14} title="Check Pulse">
            <p>Open <Code>http://localhost:8080/pulse/ui</Code>. Refresh your web app a few times, then check Pulse again. Can you see the requests being logged with their response times?</p>
          </Challenge>

          <Challenge number={15} title="Check Mailhog">
            <p>Open <Code>http://localhost:8025</Code>. Is there a welcome email from when you registered? Open it and look at the HTML template.</p>
          </Challenge>
        </section>

        {/* ═══ The .env File ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The .env File</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project has a <Code>.env</Code> file at the root. This file stores all your configuration
            as key-value pairs:
          </p>

          <CodeBlock filename=".env (partial)">
{`# Core
APP_NAME=myapp
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://postgres:postgres@localhost:5434/myapp?sslmode=disable
JWT_SECRET=your-secret-key
REDIS_URL=redis://localhost:6380

# Storage (MinIO for local development)
STORAGE_DRIVER=minio
MINIO_ENDPOINT=http://localhost:9002

# Email
RESEND_API_KEY=
MAIL_FROM=noreply@localhost

# AI (Vercel AI Gateway)
AI_GATEWAY_API_KEY=
AI_GATEWAY_MODEL=anthropic/claude-sonnet-4-6`}
          </CodeBlock>

          <Definition term="Environment Variables">
            Configuration values that change between environments (development vs production). You use
            different database URLs locally vs on your server. The <Code>.env</Code> file keeps these values
            outside your code so you can change them without modifying any Go or TypeScript files.
          </Definition>

          <Note>
            Never commit your <Code>.env</Code> file to Git — it contains secrets like API keys and database passwords.
            Grit adds it to <Code>.gitignore</Code> automatically. Use <Code>.env.example</Code> as a template for other developers.
          </Note>

          <Challenge number={16} title="Read the .env File">
            <p>Open the <Code>.env</Code> file. Find: (1) What database name is it using? (2) What port does the API run on? (3) What is the JWT_SECRET set to?</p>
          </Challenge>
        </section>

        {/* ═══ Essential CLI Commands ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Essential CLI Commands</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are commands you{"'"}ll use every day:
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
                  <td className="px-4 py-3">Lists all API endpoints in a formatted table</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit migrate</td>
                  <td className="px-4 py-3">Creates or updates database tables from your Go models</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit seed</td>
                  <td className="px-4 py-3">Fills the database with demo data (admin user, sample posts)</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit studio</td>
                  <td className="px-4 py-3">Opens the database browser in your browser</td>
                </tr>
                <tr>
                  <td className="px-4 py-3 font-mono text-foreground text-xs">grit version</td>
                  <td className="px-4 py-3">Shows the installed Grit version</td>
                </tr>
              </tbody>
            </table>
          </div>

          <Challenge number={17} title="List Your Routes">
            <p>Run <Code>grit routes</Code> in your project folder. How many routes does your app have? Can you find the login endpoint? The register endpoint? Which routes are {'"'}protected{'"'} (require authentication)?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> What Grit is — a full-stack Go + React framework</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to install Grit with <Code>go install</Code></li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to scaffold a project with <Code>grit new</Code> (interactive and with flags)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The monorepo project structure (apps/api, apps/web, apps/admin, packages/shared)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How the Go API is organized (middleware → handler → service → database)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How to start Docker services and development servers</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> The 5 built-in tools (GORM Studio, API Docs, Pulse, Sentinel, Mailhog)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How the .env file works and why you never commit it to Git</li>
          </ul>

          <Challenge number={18} title="Final Challenge: Start Fresh">
            <p>Delete the <Code>myapp</Code> folder completely. Now create a new project called <Code>bookstore</Code> using Triple architecture but with <strong>TanStack Router (Vite)</strong> instead of Next.js: <Code>grit new bookstore --triple --vite</Code>. Start everything and verify it works. Notice any differences from the Next.js version?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-web/introduction', label: 'Prev: Introduction to Grit' }}
            next={{ href: '/courses/grit-web/code-generator', label: 'Next: Code Generator Mastery' }}
          />
        </div>
      </main>
    </div>
  )
}
