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
                  The MCP server answers those questions from your source files, so the agent reads
                  facts instead of guessing. Some of them nothing else can answer: which files you
                  have edited since Grit wrote them, and which permission keys actually exist.
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
                  Every GORM model with its fields, Go types, JSON names, and GORM tags: the exact
                  shape of a request or response body, and the column constraints behind it. Pass{' '}
                  <code>model</code> to fetch just one.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_list_resources</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  What <code>grit generate resource</code> has made, and every file each one owns:
                  the model, service, handler and routes, the admin resource definition, its overlay,
                  the hooks, the Zod schema and the TypeScript types. Found by the bulk-request type
                  the generator declares, so the framework&apos;s own endpoints are not reported as
                  things somebody generated. Worth calling before generating, to see whether the
                  thing already exists.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_file_ownership</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Which generated files are still exactly as Grit wrote them, which you have edited,
                  and which have been deleted. Nothing else can answer this: it comes from the
                  manifest Grit records when it writes. It matters twice over, because editing a
                  pristine file is what stops <code>grit upgrade</code> updating it, and a modified
                  file is where an upgrade will report a conflict.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_list_permissions</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every permission key the application understands, grouped as the admin&apos;s
                  permission tree groups them. A guard written against a key that is not in the
                  catalogue matches nothing and fails silently, which is exactly the kind of mistake
                  an agent makes and nobody notices until a role does not work.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_doctor</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The same audit as the <code>grit doctor</code> command, as structured findings:
                  the check that found each one, what it is about, and the fix. Worth running after
                  generating a resource that holds personal data, and before claiming a project is
                  ready to deploy.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_env_keys</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The environment variables the project reads, with the comment documenting each and
                  whether <code>.env</code> sets it. Values are never returned, from either file:
                  this is where the database password and the signing keys live, and a tool that
                  returned one would put it in a transcript, a log and a model&apos;s context in a
                  single call.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_cli_reference</code>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The CLI&apos;s commands and flags, read from the binary that is running the server.
                  An agent proposing a flag from a different version costs you a confusing error, and
                  the binary is the only thing that knows what this version accepts.
                </p>

                <h3 className="text-lg font-semibold tracking-tight mt-8 mb-2">
                  <code>grit_generate_resource</code> <span className="text-sm font-normal text-muted-foreground">(write mode)</span>
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Generates the model, service, handler and routes, the Zod schema and TypeScript
                  types, the hooks and the admin page, with the wiring injected. This writes files,
                  so it exists only on a server started with <code>--mode write</code>.
                </p>
              </div>

              {/* ============================================================ */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Two modes, and the writing tools are absent from one
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The default server is read-only. Not because each handler checks a flag: the tools
                  that write files are <strong>not registered</strong>, so there is no code path that
                  reaches one. No token, no scope, no misconfigured client and no instruction hidden
                  in a README can call what is not in the map.
                </p>
                <CodeBlock
                  language="bash"
                  code={`grit mcp serve                 # read-only: nine tools that answer questions
grit mcp serve --mode write    # also the generators, which write files`}
                />
                <p className="text-muted-foreground leading-relaxed mt-6 mb-4">
                  The alternative, which is what most servers do, is to register everything and check
                  a flag inside each handler. That works until somebody adds a tool and forgets the
                  check, and the failure is silent and total: the tool simply works for everyone.
                  Letting the mode decide what gets built removes the class of mistake rather than
                  asking every future contributor to remember.
                </p>
                <p className="text-muted-foreground leading-relaxed">
                  Even in write mode, every change lands in your diff and you review it like any
                  other. Reading is still static: the answers come from parsing your source, so the
                  server works on a checkout that has never been started and needs no credentials.
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
                  One further kind of tool is deliberately not shipped: anything that asks a{' '}
                  <em>running</em> service about itself, such as its recent errors or its live
                  traces. Those need a connection and a credential, which is a meaningfully
                  different surface from parsing a checkout, and worth doing separately rather than
                  bolting on. A running Grit service already answers{' '}
                  <a href="/docs/ai-workflows/llms-txt">/llms.txt</a> and{' '}
                  <code>/docs/openapi.json</code> over HTTP, which covers most of what an agent
                  wants from one.
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
