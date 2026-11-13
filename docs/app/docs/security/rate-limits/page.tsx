import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/rate-limits')

export default function RateLimitsPage() {
  return (
    <DocsShell>
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Security</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">Rate limits</h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          One table of labelled rules for the whole API, rather than a limit written at each route.
          A limit written at the route is a limit somebody has to remember, and nothing in a
          handler says &quot;this one needed one and has none&quot;.
        </p>
      </div>

      <div className="prose-grit">
        <h2>How a request is matched</h2>
        <p>
          Every request produces a list of labels, from the most specific to the least. The first
          rule that matches one of them wins.
        </p>
        <CodeBlock
          language="text"
          code={`POST /api/v1/orders
  → "POST /api/v1/orders"   the route and the method
  → "/api/v1/orders"        the route, any method
  → "*:create"              the class of action
  → "/api/v1/"              a prefix
  → "/api/"                 a shorter prefix`}
        />
        <p>
          So one rule can say &quot;every create, twenty in five seconds&quot;, and a route nobody
          thought about is still bounded by the prefix at the end. That last part is the whole
          argument for a table: the set of routes keeps growing, and the set somebody remembered to
          limit does not.
        </p>
        <p>
          Specificity decides, not position in the file. Writing the catch-all first is the obvious
          mistake and it would silently give every route the loosest limit you have; it does not.
        </p>

        <h2>The classes</h2>
        <table>
          <thead>
            <tr>
              <th>Class</th>
              <th>What it covers</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><code>*:auth</code></td>
              <td>
                Anything under <code>/auth/</code>, plus <code>/login</code>,{' '}
                <code>/register</code>, <code>/password</code>, <code>/totp</code>,{' '}
                <code>/passkeys</code>, <code>/magic-link</code>, <code>/sso/</code> and an
                invitation preview. Matched on the path, so a new passkey endpoint is covered the
                day it is written.
              </td>
            </tr>
            <tr>
              <td><code>*:create</code></td>
              <td><code>POST</code>. The expensive one: an id, a row, often a mail or a job.</td>
            </tr>
            <tr>
              <td><code>*:update</code> <code>*:delete</code></td>
              <td><code>PUT</code>/<code>PATCH</code>, and <code>DELETE</code>.</td>
            </tr>
            <tr>
              <td><code>*:list</code> <code>*:read</code></td>
              <td>
                <code>GET</code> without and with a path parameter. A list scans, sorts and counts;
                a row is one index lookup.
              </td>
            </tr>
          </tbody>
        </table>

        <h2>The default table</h2>
        <CodeBlock
          language="go"
          code={`// internal/ratelimit/rules.go
func Default() Rules {
	return Rules{
		{Label: "*:auth", MaxRequests: 10, Window: 10 * time.Second},

		{Label: "*:create", MaxRequests: 20, Window: 5 * time.Second},
		{Label: "*:update", MaxRequests: 30, Window: 5 * time.Second},
		{Label: "*:delete", MaxRequests: 20, Window: 5 * time.Second},

		{Label: "/api/v1/batch", MaxRequests: 3, Window: time.Second},

		{Label: "/api/", MaxRequests: 300, Window: 10 * time.Second, Audience: AudienceGuest},
		{Label: "/api/", MaxRequests: 600, Window: 10 * time.Second, Audience: AudienceAuth},
	}
}`}
        />
        <p>
          Edit it. It is ordinary Go in your own repository, which is the point: a limit that lives
          in a dashboard is a limit nobody reviews, and one that lives in a handler is a limit
          nobody finds. Development multiplies every number by twenty rather than removing the
          rules, because a limit nobody exercises in development is a limit that surprises somebody
          in production.
        </p>

        <Callout type="warning" title="The credential rule has no audience, deliberately">
          The middleware runs before each route&apos;s own authentication, so the only thing it can
          know about a caller is whether the request <em>carries</em> credentials, not whether they
          are real. Split <code>*:auth</code> by audience and a credential-stuffing client sends
          one junk bearer token and buys the signed-in budget. Nothing about a sign-in route needs
          a relaxation for people who are already signed in.
        </Callout>

        <h2>Turning one off</h2>
        <p>
          Set <code>MaxRequests</code> to zero. A rule at zero is <strong>off</strong>, and the
          match stops there rather than falling through to a broader rule, which would silently
          reinstate a limit somebody deliberately removed.
        </p>
        <CodeBlock
          language="go"
          code={`// An inbound webhook you do not control the volume of.
{Label: "POST /api/v1/webhooks/inbound", MaxRequests: 0},`}
        />

        <h2>What the caller sees</h2>
        <CodeBlock
          language="text"
          code={`X-RateLimit-Limit: 300
X-RateLimit-Remaining: 287

HTTP/1.1 429 Too Many Requests
Retry-After: 7
{"error":{"code":"RATE_LIMITED","message":"Too many requests. Try again in 7 seconds."}}`}
        />
        <p>
          The headers go out whether or not the request was allowed, so a client that can see its
          budget can slow down before it is stopped. <code>Retry-After</code> is there because a
          client told only &quot;no&quot; retries immediately, which is the same request again.
        </p>

        <h2>Who is counted</h2>
        <p>
          The user id when the request has one, so a signed-in caller behind a shared office
          address gets their own budget rather than the building&apos;s. The client address
          otherwise, which is all a guest has.
        </p>

        <h2>This is a floor, not a ceiling</h2>
        <p>
          The counters are in process and per instance, so behind three replicas a caller gets
          three times the budget. That is deliberate: it needs no Redis and it cannot fail open.{' '}
          <Link href="/docs/security/defenders-handbook">Sentinel</Link> keeps its counters in
          shared storage and is what bounds the total across replicas. Both are mounted, and where
          both apply the stricter refuses first.
        </p>
        <p>
          Validation runs over the table so a typo is found at start-up:{' '}
          <code>*:craete</code> is not a class, a rule with a limit and no window never resets, and
          two rules with the same label and audience mean the second can never be reached.
        </p>
      </div>

      <div className="flex items-center justify-between mt-12 pt-8 border-t border-border/40">
        <Button asChild variant="ghost">
          <Link href="/docs/security/authorization">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Authorization
          </Link>
        </Button>
        <Button asChild>
          <Link href="/docs/security/defenders-handbook">
            Defenders handbook
            <ArrowRight className="ml-2 h-4 w-4" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
