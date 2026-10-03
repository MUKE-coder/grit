import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/ai-workflows/llms-txt')

export default function LLMSTxtPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">AI agent workflows</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">llms.txt</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                An agent pointed at Grit, or at an API you built with it, should not have to
                crawl a site or guess at conventions. Both serve an{' '}
                <a href="https://llmstxt.org" target="_blank" rel="noreferrer">
                  llmstxt.org
                </a>{' '}
                index at the path models look for it.
              </p>
            </div>

            <div className="prose-grit">
              <h2>For working with the framework</h2>
              <p>
                <a href="/llms.txt">gritframework.dev/llms.txt</a> is the index: what Grit is,
                the one thing to understand about it, and every documentation page with the
                description it carries, grouped by subject.{' '}
                <a href="/llms-full.txt">gritframework.dev/llms-full.txt</a> is that index, plus
                the working guide an agent should follow, plus the complete CLI reference, as one
                file of about 80 KB.
              </p>
              <CodeBlock
                language="bash"
                code={`curl https://gritframework.dev/llms-full.txt`}
              />
              <p>
                Both are built from the files the site itself is built from: the page list is the
                same metadata the pages use for their titles, the CLI reference is the same
                catalogue <code>/docs/cli</code> renders, and the working guide is the same file{' '}
                <code>grit init</code> writes into a project. None of it is written twice, so
                none of it can drift from the docs.
              </p>

              <h2>For calling an API you built</h2>
              <p>
                A scaffolded API serves its own pair. <code>/llms.txt</code> is the orientation:
                how versioning works, which header carries the token, what a response and an error
                look like, and the parameters every list endpoint takes.
              </p>
              <CodeBlock
                language="text"
                filename="GET /llms.txt"
                code={`# acme API

> A REST API built with Grit (Go, Gin, GORM). This file is the orientation;
> the contract is the OpenAPI spec linked below.

## Start here

- [OpenAPI 3.1 spec](https://api.acme.com/docs/openapi.json): every endpoint, its
  parameters and its schemas. Generate a client from this.
- [Every route this process serves](https://api.acme.com/llms-full.txt)
- [Dependency health](https://api.acme.com/api/health)

## Versioning
...
## Authentication
...`}
              />
              <p>
                <code>/llms-full.txt</code> adds every route the router holds, grouped by what it
                is about. That is the one thing the OpenAPI spec cannot answer: the spec documents
                the routes somebody wrote an override for, and this is the router&apos;s own
                table, so a route you add appears without anybody maintaining a list.
              </p>

              <Callout type="note" title="The spec is still the contract">
                For generating a client, read <code>/docs/openapi.json</code>. llms.txt exists
                because 300 KB of JSON says nothing about the conventions every endpoint shares,
                and because an agent that has just been handed a URL needs somewhere to start.
              </Callout>

              <h2>Both are gated with the API reference</h2>
              <p>
                They describe the whole surface, admin routes included, so they are mounted under
                the same condition as <code>/docs</code>: off in production unless{' '}
                <code>API_DOCS_PUBLIC=true</code>. A route list is not a secret, and it is also not
                something to hand out by default.
              </p>

              <h2>Adding to what they say</h2>
              <p>
                The text is <code>internal/llms</code> in your project, which is yours to edit
                like any generated file: the conventions section is a string, and a service with
                its own idempotency header or its own rate limit policy should say so there. The
                route listing is read from the engine on each request, so nothing has to be kept
                in step by hand.
              </p>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/ai-workflows/mcp">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    MCP Server
                  </Button>
                </Link>
                <Link href="/docs/backend/api-docs">
                  <Button variant="outline" size="sm" className="gap-2">
                    API Documentation
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
