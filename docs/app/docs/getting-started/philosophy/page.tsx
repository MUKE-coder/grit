import Link from 'next/link'
import { ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/getting-started/philosophy')

export default function PhilosophyPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Getting Started</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Philosophy &amp; Inspiration
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Why Grit exists, what inspired it, and the design principles that guide every decision
                in the framework.
              </p>
            </div>

            <div className="prose-grit">
              {/* Why Grit Exists */}
              <h2>Why Grit Exists</h2>
              <p>
                Building a modern full-stack application in 2025/2026 requires stitching together
                15+ tools: Next.js, Prisma, NextAuth, tRPC, Zod, React Query, Tailwind, S3 SDK,
                Resend, Redis client, Bull queues, an admin panel library... Every new project
                starts with 2-3 days of boilerplate. Every developer makes different choices,
                creating inconsistent codebases that are hard to maintain.
              </p>
              <p>
                Meanwhile, Go has incredible performance, simplicity, and deployment story. But the
                Go web ecosystem is fragmented -- Gin, Echo, Fiber, Chi for routing; GORM, sqlc, sqlx
                for databases; dozens of auth libraries. There is no equivalent to Laravel or Rails
                that gives you everything wired together. Go developers spend more time plumbing
                than building.
              </p>
              <p>
                And then there is the admin panel problem. Laravel has <strong>Filament</strong> -- a tool
                that lets you generate beautiful admin dashboards, data tables, forms, and widgets
                from simple PHP definitions. The Go + React world has nothing comparable. Every team
                builds their admin panel from scratch, wasting weeks on repetitive CRUD interfaces.
              </p>
              <p>
                <strong>Grit bridges this gap.</strong> It takes the best ideas from Laravel, Rails,
                and Django -- convention over configuration, batteries included, code generation --
                and applies them to the most performant and modern stack available: Go for the
                backend, React for the frontend.
              </p>

              {/* Inspiration */}
              <h2>Inspiration</h2>
              <p>
                Grit draws inspiration from frameworks that have defined developer experience
                in their respective ecosystems:
              </p>
              <ul>
                <li>
                  <strong>Laravel + Filament (PHP)</strong> -- The gold standard for developer experience.
                  Artisan CLI, Eloquent ORM, and Filament&apos;s resource-based admin panel are the direct
                  inspiration for Grit&apos;s CLI and admin panel generator.
                </li>
                <li>
                  <strong>Ruby on Rails</strong> -- Convention over configuration, code generation with
                  scaffolding, and the idea that frameworks should have opinions. Rails proved that
                  opinionated frameworks make teams more productive.
                </li>
                <li>
                  <strong>Django (Python)</strong> -- The &quot;batteries included&quot; philosophy. Auth, admin,
                  ORM, and email all ship with the framework. Grit applies this same idea to Go + React.
                </li>
                <li>
                  <strong>Next.js</strong> -- File-based routing, server-side rendering, and the best
                  React developer experience. Grit uses Next.js as its frontend runtime.
                </li>
                <li>
                  <strong>GORM Studio</strong> -- Our own visual database browser that evolved into
                  a key feature of the framework. Seeing your data without leaving the browser
                  changes how you develop.
                </li>
              </ul>
              <p>
                The key insight is that no single existing framework combines Go&apos;s backend performance
                with React&apos;s frontend ecosystem AND a Filament-class admin panel in one coherent
                monorepo. Grit exists to fill that gap.
              </p>

              {/* Design Principles */}
              <h2>Design Principles</h2>

              <h3>1. Convention Over Configuration</h3>
              <p>
                There is <strong>one way</strong> to do things in Grit. One auth system. One state management
                approach. One folder structure. One naming convention. Opinions are features, not limitations.
              </p>
              <p>
                This means any developer can jump into any Grit project and immediately understand where
                things live, how data flows, and how to add new features. No more spending the first
                day on a project just figuring out the folder structure.
              </p>

              <h3>2. Code Generation Over Runtime Magic</h3>
              <p>
                Instead of complex runtime coupling between Go and React, Grit uses a CLI code generator.
                When you run <code>grit generate resource Post</code>, it creates all the files -- Go model,
                handler, React hook, Zod schema, admin page. The generated code is readable, editable,
                and debuggable.
              </p>
              <p>
                There are no hidden proxies, no auto-wiring at runtime, no reflection magic. You can open
                any generated file, read it, understand it, and modify it. If you want to change how a
                handler works, just edit the file. The framework gets out of your way.
              </p>

              <h3>3. Own Your Code</h3>
              <p>
                Every piece of code Grit generates belongs to you. It is installed directly into your
                project, not hidden in a <code>node_modules</code> folder or compiled into a binary you
                cannot inspect. No lock-in, no black boxes.
              </p>
              <p>
                If you ever decide to stop using Grit&apos;s CLI, your project still works. It is just
                a Go API, a Next.js app, and some TypeScript schemas. Standard tools, standard deployment.
              </p>

              <h3>4. Batteries Included, Optionally Removable</h3>
              <p>
                Auth, file storage, email, queues, cron, AI, and the admin panel ship with every project.
                But they are modular -- you can remove what you do not need. The framework should be great
                for a freelancer building a simple app AND for a team building a complex SaaS.
              </p>
              <p>
                Every battery is a separate Go package and can be deleted without breaking the rest
                of the application. Nothing is deeply coupled.
              </p>

              <h3>5. Beautiful by Default</h3>
              <p>
                Every UI component -- the admin panel, data tables, forms, login pages -- ships with
                a polished dark theme. First impressions matter. When a developer runs{' '}
                <code>grit new</code> and sees the result, they should think &quot;this looks like a
                real product,&quot; not &quot;this looks like a tutorial.&quot;
              </p>
              <p>
                The design language is inspired by tools like Linear, Vercel Dashboard, and Raycast --
                dark, polished, fast. Not generic Bootstrap. Not Material Design.
              </p>

              <h3>6. AI-Friendly (Vibe Coding Ready)</h3>
              <p>
                Grit&apos;s strict conventions and predictable structure make it ideal for AI-assisted
                development. An AI assistant like Claude Code, Cursor, or GitHub Copilot can
                understand the entire project structure, generate resources, and modify code
                confidently because Grit follows the same patterns everywhere.
              </p>
              <p>
                This is not an afterthought. Every naming convention, every folder structure decision,
                every code pattern in Grit was designed with the question: &quot;Can an AI reliably
                generate this?&quot; If the answer was no, the convention was simplified until the
                answer became yes.
              </p>

              {/* The Vibe Coding Philosophy */}
              <h2>The &quot;Vibe Coding&quot; Philosophy</h2>
              <p>
                We believe the future of software development is a collaboration between human
                developers and AI assistants. Grit is designed from the ground up for this paradigm.
              </p>
              <p>
                When you describe a feature to an AI assistant -- &quot;add a blog with posts,
                categories, and tags&quot; -- the AI should be able to:
              </p>
              <ol>
                <li>Run <code>grit generate resource Post</code>, <code>grit generate resource Category</code>, and <code>grit generate resource Tag</code></li>
                <li>Add relationships between the generated models</li>
                <li>Customize the admin panel columns and form fields</li>
                <li>Add frontend pages that display the blog</li>
              </ol>
              <p>
                This works because Grit&apos;s patterns are <strong>predictable</strong>. The AI knows
                exactly where models live, how handlers are structured, where hooks go, and how
                the admin panel is configured. There is no ambiguity.
              </p>
              <p>
                Frameworks that let developers &quot;choose your own adventure&quot; at every step make
                AI assistance unreliable. Grit&apos;s opinionated nature is what makes vibe coding work.
              </p>

              {/* Why Go + React */}
              <h2>Why Go + React?</h2>
              <p>
                We considered every major backend/frontend pairing. Here is why Go + React won:
              </p>

              <h3>Why Go (and not Node.js, Python, Ruby, or PHP)?</h3>
              <ul>
                <li>
                  <strong>Performance:</strong> Go is 10-50x faster than Node.js/PHP for CPU-bound
                  work. A single Go binary can handle 100k+ concurrent connections with minimal
                  memory usage.
                </li>
                <li>
                  <strong>Simplicity:</strong> Go has 25 keywords. The language is small, explicit,
                  and readable. There is no magic, no metaprogramming, no inheritance hierarchies.
                </li>
                <li>
                  <strong>Deployment:</strong> Go compiles to a single static binary. No runtime
                  dependencies, no version managers, no interpreters. Copy the binary to a server
                  and run it.
                </li>
                <li>
                  <strong>Type Safety:</strong> Go&apos;s static type system catches bugs at compile time.
                  Combined with GORM&apos;s struct tags and Zod on the frontend, you get end-to-end
                  type safety from database to UI.
                </li>
                <li>
                  <strong>Growing Ecosystem:</strong> Go is the fastest-growing backend language.
                  Companies like Google, Uber, Twitch, Cloudflare, and Docker use it in production.
                </li>
              </ul>

              <h3>Why React (and not Vue, Svelte, or HTMX)?</h3>
              <ul>
                <li>
                  <strong>Ecosystem:</strong> React has the largest component ecosystem of any
                  frontend framework. shadcn/ui, Radix, React Query, Recharts -- the best tools
                  are built for React first.
                </li>
                <li>
                  <strong>Hiring:</strong> More developers know React than any other frontend
                  framework. If you are building a team or hiring freelancers, React is the
                  safest bet.
                </li>
                <li>
                  <strong>Next.js:</strong> The App Router provides file-based routing, server-side
                  rendering, static generation, and streaming -- all built on React.
                </li>
                <li>
                  <strong>React Native:</strong> With Grit&apos;s <code>--full</code> flag, you can add
                  an Expo mobile app that shares types and validation with the web frontend. Same
                  language, same patterns, mobile and web.
                </li>
              </ul>

              <h3>Why Not Full-Stack JavaScript?</h3>
              <p>
                Serverless is not always the answer. For many applications -- CRMs, SaaS tools,
                internal dashboards, apps with WebSockets, background jobs, or heavy processing --
                a real backend server is superior. Go gives you a single binary that handles
                massive concurrency, deploys anywhere, and costs a fraction of serverless.
              </p>
              <p>
                And while TypeScript is excellent on the frontend, Go&apos;s simplicity and performance
                make it a better choice for backend services that need to handle real load.
              </p>
            </div>

            {/* Next/Prev navigation */}
            <div className="flex flex-wrap gap-3 mt-12 pt-6 border-t border-border/30">
              <Button variant="outline" asChild className="border-border/60 bg-transparent hover:bg-accent/50">
                <Link href="/docs">
                  Introduction
                </Link>
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
