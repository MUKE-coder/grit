import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-workflows/mcp')

export default function MCPServerPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">AI Workflows</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">MCP Server</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Let your AI coding agent ask Grit about your project instead of guessing from a
                grep. <code>grit mcp serve</code> speaks the Model Context Protocol and answers with
                the real route table, the real model definitions, and the real layout.
              </p>
            </div>

            <div className="prose-grit">
              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Why</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  An agent working in a Grit project has to answer the same questions over and over.
                  What is the URL for listing users? Does that endpoint need an admin token? What
                  fields does <code>Invoice</code> actually have? Without a way to ask, it greps,
                  infers, and gets it subtly wrong &mdash; usually by dropping the{' '}
                  <code>/api/v1</code> prefix or inventing a field that isn&apos;t there.
                </p>
                <p className="text-muted-foreground leading-relaxed">
                  The MCP server answers those three questions from your source files, so the agent
                  reads facts instead of guessing.
                </p>
              </div>

              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Setup</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  With Claude Code, one command:
                </p>
                <CodeBlock
                  language="bash"
                  code={`claude mcp add grit -- grit mcp serve --project /path/to/your-project`}
                />
                <p className="text-muted-foreground leading-relaxed mt-6 mb-4">
                  For any other MCP client, add it to the client&apos;s config:
                </p>
                <CodeBlock
                  language="json"
                  code={`{
  "mcpServers": {
    "grit": {
      "command": "grit",
      "args": ["mcp", "serve", "--project", "/path/to/your-project"]
    }
  }
}`}
                />
                <p className="text-muted-foreground leading-relaxed mt-6">
                  Omit <code>--project</code> and the server searches upward from its working
                  directory for <code>grit.json</code>, the same way every other Grit command finds
                  your project.
                </p>
              </div>

              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">The tools</h2>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_project_info</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Architecture, frontend framework, Go module path, the CLI version that scaffolded
                  the project, and which apps exist. Worth calling first &mdash; it tells the agent
                  whether Go code lives at the root or under <code>apps/api</code>, which is the
                  thing most often assumed wrongly.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_list_routes</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every registered route with its method, full path including the{' '}
                  <code>/api/v1</code> prefix, handler, and access level (
                  <code>public</code>, <code>protected</code>, <code>admin</code>). Takes optional{' '}
                  <code>method</code> and <code>contains</code> filters so the agent can ask a narrow
                  question instead of pulling 140 routes into its context.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_describe_models</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every GORM model with its fields, Go types, JSON names, and GORM tags &mdash; the
                  exact shape of a request or response body, and the column constraints behind it.
                  Pass <code>model</code> to fetch just one.
                </p>
              </div>

              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Read-only, and static on purpose
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every tool answers by parsing your source. None of them connects to a running
                  server or a database, and none of them writes anything. That is a design decision,
                  not a first-draft limitation, and it buys four things:
                </p>
                <ul className="text-muted-foreground leading-relaxed space-y-2 mb-4">
                  <li>
                    It works on a checkout that has never been started &mdash; no{' '}
                    <code>docker compose up</code> first.
                  </li>
                  <li>It needs no credentials, so there is no secret to leak into an agent&apos;s context.</li>
                  <li>It cannot mutate your repo.</li>
                  <li>
                    It cannot be talked into running a migration by instructions hidden in a README
                    or an issue comment.
                  </li>
                </ul>
                <p className="text-muted-foreground leading-relaxed">
                  An agent that wants to <em>change</em> your project still has to call the CLI, where
                  the change lands in your diff and you review it like any other.
                </p>
              </div>

              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">Scope</h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The server targets Grit&apos;s web and API architectures &mdash;{' '}
                  <code>single</code>, <code>double</code>, <code>triple</code>, <code>api</code>,
                  and <code>mobile</code>. Standalone desktop projects created with{' '}
                  <code>grit new-desktop</code> have a different shape (Wails bindings rather than
                  HTTP routes) and are not covered yet.
                </p>
                <p className="text-muted-foreground leading-relaxed">
                  Two further tools are planned but deliberately not shipped yet:{' '}
                  <code>openapi</code> and <code>recent_errors</code>. Both need a running server and
                  a database connection, which means a credential and connection story the read-only
                  tools above do not require &mdash; a meaningfully different surface, worth doing
                  separately rather than bolting on.
                </p>
              </div>
            </div>

            {/* Prev / Next */}
            <div className="flex items-center justify-between border-t border-border/40 pt-6 mt-12">
              <Button variant="ghost" asChild>
                <Link href="/docs/ai-workflows/antigravity">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  Using Grit with Antigravity
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/ai-skill">
                  LLM Skill Guide
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
