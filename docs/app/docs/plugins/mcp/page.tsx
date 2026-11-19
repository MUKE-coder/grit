import Link from 'next/link'
import { ArrowLeft, ArrowRight, Bot, KeyRound, ShieldCheck, ScrollText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { PluginSteps, RelatedPlugins } from '@/components/plugin-detail'
import type { PluginStep, RelatedPlugin } from '@/components/plugin-detail'

export const metadata = {
  title: 'MCP plugin — Grit',
  description:
    'An MCP server for your app, so Claude or Cursor can use it through the same permissions, tenancy and audit trail a person does.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/mcp' },
}

const STEPS: PluginStep[] = [
  {
    title: 'Install it',
    run: 'grit plugin add mcp',
    then: 'Writes internal/mcp (the transport, the tool registry, the in-process caller), the connection model and service, and mounts POST /mcp at the root.',
  },
  {
    title: 'Run the migration',
    run: 'grit migrate',
    then: 'Two tables: the connections, and what every agent did with them.',
  },
  {
    title: 'Make a connection',
    run: 'POST /api/v1/mcp-connections',
    then: 'A name and the tools it may call. It answers with a key, once. Any signed-in person can make their own: it acts as them and can do no more than they can.',
  },
  {
    title: 'Point a client at it',
    then: 'The mcpServers entry below. The key goes in the Authorization header.',
  },
]

const RELATED: RelatedPlugin[] = [
  {
    name: 'Multitenant',
    href: '/docs/plugins/multitenant',
    why: 'A connection carries its organization, and every tool call is scoped by the same callbacks a browser request is. Nothing extra to do.',
  },
  {
    name: 'Impersonate',
    href: '/docs/plugins/impersonate',
    why: 'The other way a session exists without the account holder typing a password. Both are audited, and for the same reason.',
  },
]

export default function MCPPluginPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Plugins</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">MCP</h1>
        <p className="text-lg text-muted-foreground">
          An MCP server for your app, so an agent works through the same rules a person does.
        </p>
      </div>

      <div className="prose-grit">
        <div className="not-prose mb-10 rounded-xl border border-border/60 bg-muted/20 p-5">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Not the same thing as{' '}
            <Link href="/docs/ai-workflows/mcp" className="text-accent hover:underline">
              the MCP server under AI Workflows
            </Link>.</strong>{' '}
            That one runs on your machine while you build, and answers questions about your
            code: which routes exist, what fields a model has, which files you have edited.
            This one runs inside your deployed app and exposes its <em>data</em>, to an agent
            acting as one of your users. One helps you write the app; the other lets somebody
            use it.
          </p>
        </div>

        <h2>The decision the rest follows from</h2>
        <p>
          A tool call does not touch the database. It builds an HTTP request and runs it in
          process against the same gin engine that serves the browser.
        </p>
        <p>
          That is not a shortcut taken for want of a cheaper path. A Grit API decides what a
          caller may do in middleware and in services: authentication, the organization scope,
          the role, the per-resource policy, the module switch, the record-level scope, the
          validation tags, the audit log. A tool that queried GORM directly would bypass every
          one of them, and would have to reimplement them. Two implementations that start
          identical and drift, and the one that drifts is the one nobody is looking at.
        </p>
        <p>
          Running through the router costs a function call and some JSON. What it buys is that
          an agent is <strong>structurally incapable</strong> of doing something the same person
          could not do in the browser, and that every tool call lands in the audit log the same
          way a click does.
        </p>

        <h2>One connection is one key</h2>
        <p>
          A person makes a connection in the app: a name, and the tools it may use. That mints a
          secret API key that acts as them, and the key is narrowed twice, both times from the
          tools rather than by choice:
        </p>
        <ul>
          <li>
            <strong>Scopes</strong>, to the permissions those tools declare, so the key cannot
            reach a permission no chosen tool needs.
          </li>
          <li>
            <strong>Endpoints</strong>, to the routes those tools call. This is the one that
            makes a read-only connection structurally read-only: it holds no POST route, so a
            write is refused by the API key middleware before any handler runs, whatever the
            agent asks for.
          </li>
        </ul>
        <p>
          Derived and not chosen, so nobody widens a connection by accident, and adding a tool
          to a resource later does not quietly widen the connections that already exist. A tool
          whose own role does not grant its permission is refused when the connection is made,
          because an agent must not be able to do more than the person who set it up.
        </p>
        <p>
          Revoking the connection revokes the key, in that order and immediately. The
          alternative, one <code>MCP_API_KEY</code> for the whole app, cannot say who did what,
          cannot be narrowed to the twelve tools one person needs, and cannot be revoked without
          breaking everybody.
        </p>

        <PluginSteps steps={STEPS} />

        <h2>Pointing a client at it</h2>
        <CodeBlock
          language="json"
          filename=".mcp.json"
          code={`{
  "mcpServers": {
    "my-app": {
      "type": "http",
      "url": "https://api.example.com/mcp",
      "headers": {
        "Authorization": "Bearer grit_sk_..."
      }
    }
  }
}`}
        />
        <p>
          Streamable HTTP, and deliberately the plain half of it: every JSON-RPC request gets a
          JSON answer rather than an event stream. The transport allows that, and no tool
          reports progress, so a stream would be a second code path with nothing to carry. A{' '}
          <code>GET</code> on the endpoint, which is how a client opens the stream half, answers
          405 with <code>Allow: POST</code> rather than leaving it waiting.
        </p>
        <p>
          Stateless: every request carries its own key, so there is no session to create, lose
          or pin to one process. It can sit behind a load balancer and scale to zero.
        </p>
        <p>
          It implements <code>initialize</code> (protocol 2025-11-25 down to 2024-11-05),{' '}
          <code>ping</code>, <code>tools/list</code> and <code>tools/call</code>, and it handles
          batches. A version it has never heard of gets the oldest it speaks rather than a
          refusal, because a client from the future can still talk over 2024-11-05 and a refusal
          looks to the person like the server being down.
        </p>

        <h2>Writing a tool</h2>
        <p>
          The generated CRUD tools are the base, not the ceiling: the interesting tools in any
          app are the ones that are not create-read-update-delete. Yours go in{' '}
          <code>internal/mcp/custom_tools.go</code>, which Grit never rewrites.
        </p>
        <CodeBlock
          language="go"
          filename="internal/mcp/custom_tools.go"
          code={`func init() {
	Register(Tool{
		Name:        "send_invoice",
		Title:       "Send an invoice",
		Description: "Emails an invoice to its client and marks it sent. " +
			"The invoice must be a draft. Takes the invoice number or its id.",
		Group:      "Invoices",
		Permission: "invoices.edit",
		Write:      true,
		Routes: []string{
			"GET /api/v1/invoices",
			"POST /api/v1/invoices/:id/send",
		},
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"invoice": map[string]any{
					"type":        "string",
					"description": "The invoice number, or its id.",
				},
			},
			"required": []string{"invoice"},
		},
		Exec: func(ctx context.Context, call *Call, args map[string]any) (any, error) {
			term, err := String(args, "invoice")
			if err != nil {
				return nil, err
			}
			id, err := Resolve(ctx, call, "invoice", "/api/v1/invoices", term)
			if err != nil {
				return nil, err
			}
			result, err := call.Post(ctx, "/api/v1/invoices/"+id+"/send", nil)
			if err != nil {
				return nil, err
			}
			if !result.OK() {
				return nil, fmt.Errorf("%s", result.Message())
			}
			return map[string]any{"sent": true, "invoice_id": id}, nil
		},
	})
}`}
        />
        <p>
          Four things in there are load-bearing. <code>Routes</code> is a security boundary and
          not documentation: the connection&apos;s key is restricted to the union of its tools&apos;
          routes, so a route left out is refused at runtime by the app itself.{' '}
          <code>Permission</code> is what keeps the tool off the list for somebody whose role
          lacks it. <code>Write</code> and <code>Destructive</code> drive the hints a client uses
          to decide whether to confirm with the person first. And <code>Exec</code> goes through{' '}
          <code>call.Post</code> rather than GORM, for the reason at the top of this page.
        </p>

        <h2>Arguments, as a model actually sends them</h2>
        <p>
          The schema says a field is a number and the model sends <code>&quot;5&quot;</code>,
          because that is what the person typed. It says the field is an id and the model sends
          the customer&apos;s name. Refusing both is correct by the schema and useless in
          practice: the call fails, the model apologises, and the person is no further forward.
        </p>
        <p>
          So the helpers are forgiving about shape and strict about meaning.{' '}
          <code>Int</code> reads <code>&quot;7&quot;</code>, <code>Float</code> reads{' '}
          <code>&quot;$1,250.50&quot;</code>, <code>Strings</code> reads a bare item as a list of
          one. <code>Resolve</code> turns a name into an id by asking the app, so the lookup is
          scoped to what this connection may see.
        </p>
        <p>
          What they never do is guess which record was meant. A name matching more than one
          returns an error naming the candidates, which the model can put back to the person.
          Picking the first match is the behaviour that quietly invoices the wrong customer.
        </p>
        <p>
          Paging is capped at 100 whatever the model asks for. &ldquo;List every
          invoice&rdquo; with no page size is a model requesting a hundred thousand rows it will
          then put in its own context.
        </p>

        <h2>What an agent did</h2>
        <p>
          Every call is recorded: the connection, the person, the tool, the arguments, whether
          it worked, the error if not, and how long it took. Refusals as well as successes, which
          are often the more interesting half.
        </p>
        <CodeBlock
          language="text"
          code={`GET    /api/v1/mcp-connections            your own connections
POST   /api/v1/mcp-connections            make one; answers with the key, once
GET    /api/v1/mcp-connections/tools      the tools your role lets you choose
DELETE /api/v1/mcp-connections/:id        revoke it, and its key
GET    /api/v1/mcp-connections/:id/calls  what that agent did`}
        />
        <p>
          <code>call_count</code> counts tool calls and nothing else: the handshake and every{' '}
          <code>tools/list</code> are excluded, because otherwise an idle client looks busy.{' '}
          <code>last_used_at</code> is every request, since &ldquo;when did this agent last
          connect&rdquo; is the question that one answers.
        </p>

        <h2>What is not in this release</h2>
        <p>
          Generated CRUD tools per resource, so{' '}
          <code>grit generate resource Invoice</code> writes{' '}
          <code>list_invoices</code>, <code>get_invoice</code> and the rest with schemas built
          from the field types. That is the next one, and it is where the leverage is. This
          release is the server, the connections, and the registry those tools will register
          into, which is also everything a hand-written tool needs today.
        </p>
        <p>
          The connection screens in the admin come with it. Until then a connection is made over
          the API, which is also how a headless project would do it.
        </p>

        <RelatedPlugins items={RELATED} />
      </div>

      <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
        <Link href="/docs/plugins/marketplace">
          <Button variant="ghost" className="gap-2">
            <ArrowLeft className="h-4 w-4" />
            Marketplace
          </Button>
        </Link>
        <Link href="/docs/plugins/multitenant">
          <Button variant="ghost" className="gap-2">
            Multitenant
            <ArrowRight className="h-4 w-4" />
          </Button>
        </Link>
      </div>
    </DocsShell>
  )
}
