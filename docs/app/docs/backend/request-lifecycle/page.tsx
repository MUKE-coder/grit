import Link from 'next/link'
import { ArrowLeft, ArrowRight, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/request-lifecycle')

export default function RequestLifecyclePage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Backend
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">
            The Request Lifecycle
          </h1>
          <p className="mb-6 text-lg text-muted-foreground">
            Everything a request passes through, in the order it actually happens. Most of the
            confusing bugs in a Go API — a CORS header that never appears, a WAF rule that blocks
            the wrong thing, a CSRF rejection on a mobile client — are ordering problems. This page
            is the map.
          </p>

          <div className="prose-grit">
            <h2 id="order">The order</h2>
            <p>
              From <code>internal/routes/routes.go</code>, top to bottom. Each layer wraps the ones
              below it, so a request goes <em>down</em> the list and the response comes back{' '}
              <em>up</em>.
            </p>

            <CodeBlock language="text" code={`  ┌─ request in
  │
  1  Maintenance          503 for everyone but allowlisted IPs
  2  SecurityHeaders      CSP, HSTS, X-Frame-Options, nosniff
  3  MaxBodySize          10 MB cap — before anything reads the body
  4  RequestID            X-Request-ID, threaded through every log line
  5  Logger               method, path, status, latency, request id
  6  Recovery             turns a panic into a 500 instead of a dead process
  7  CORS                 preflight + Access-Control-* headers
  8  Gzip                 response compression
  9  AutoCSRF             enforces ONLY on cookie-authenticated mutations
 10  Idempotency          replays the cached 2xx when Idempotency-Key repeats
 11  Sentinel             WAF, rate limits, AuthShield, anomaly + geo  (if enabled)
 12  Pulse                tracing, N+1 detection, runtime metrics       (if enabled)
  │
  ├─ route group middleware
  │    protected:  Auth → ActivityLogger
  │    admin:      Auth → RequireRole("ADMIN")
  │
  └─ your handler`} />

            <h2 id="global">The global layers</h2>
            <p>
              Order here is deliberate, and a few positions are load-bearing:
            </p>
            <ul>
              <li>
                <strong>Maintenance is first</strong> so a maintenance window costs nothing — no
                body read, no DB, no auth.
              </li>
              <li>
                <strong>MaxBodySize precedes anything that reads the body.</strong> A cap applied
                after parsing is not a cap.
              </li>
              <li>
                <strong>Recovery sits after Logger, not before.</strong> Gin runs middleware in
                registration order, so Logger wraps Recovery and a panicking request still gets
                logged with its status. Flip them and panics vanish from your logs.
              </li>
              <li>
                <strong>CORS runs before auth</strong>, because a browser preflight
                (<code>OPTIONS</code>) carries no credentials. Put CORS behind auth and every
                cross-origin call fails preflight with a 401 that never reaches your handler — the
                single most common &ldquo;my frontend can&apos;t call my API&rdquo; cause.
              </li>
            </ul>

            <h3>CSRF only applies to cookie sessions</h3>
            <p>
              <code>AutoCSRF</code> enforces the double-submit token only when a request is a
              state-changing method <em>and</em> authenticated by cookie. It deliberately skips:
            </p>
            <ul>
              <li>safe methods (it issues the token there instead);</li>
              <li>
                requests carrying <code>Authorization: Bearer</code> — those authenticate
                explicitly and can&apos;t be forged cross-site;
              </li>
              <li>
                bootstrap endpoints (<code>login</code>, <code>register</code>,{' '}
                <code>refresh</code>, password reset, TOTP verify) which have no token yet;
              </li>
              <li>
                the SAML assertion consumer — the identity provider posts it from its own origin
                and will never have a token.
              </li>
            </ul>
            <p>
              This is why a mobile or desktop client never sends a CSRF header and still works:
              those flows are bearer-authenticated. If you add an endpoint that must be callable
              cross-origin without a session, it needs an explicit exemption — see{' '}
              <code>internal/middleware/csrf.go</code>.
            </p>

            <h2 id="groups">Group middleware</h2>
            <p>
              Below the global stack, routes hang off groups that add their own layers:
            </p>
            <CodeBlock
              language="go"
              filename="internal/routes/routes.go"
              code={`v1 := r.Group("/api/" + APIVersion)

// Public: auth endpoints, SSO, public form submissions.
auth := v1.Group("/auth")

// Authenticated: everything a signed-in user can reach.
protected := v1.Group("")
protected.Use(middleware.Auth(db, authService))
protected.Use(middleware.ActivityLogger(db))

// Admin-only.
admin := v1.Group("")
admin.Use(middleware.Auth(db, authService))
admin.Use(middleware.RequireRole("ADMIN"))`}
            />
            <p>
              <code>Auth</code> populates the gin context with <code>user_id</code>,{' '}
              <code>user_email</code>, <code>user_role</code> and <code>user_grants</code>, so
              everything after it — including your handler — can read the caller without another
              query.
            </p>

            <h2 id="where">Where does my logic go?</h2>
            <table>
              <thead>
                <tr>
                  <th>You want to…</th>
                  <th>Put it…</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Reject a request before any work happens</td>
                  <td>Global middleware, high in the list</td>
                </tr>
                <tr>
                  <td>Require a permission for one route</td>
                  <td>
                    <code>middleware.RequireRole(&quot;perm:invoices.delete&quot;)</code> on that
                    route
                  </td>
                </tr>
                <tr>
                  <td>Require it for a whole section</td>
                  <td>
                    <code>.Use()</code> on a route group
                  </td>
                </tr>
                <tr>
                  <td>Check the caller owns the record</td>
                  <td>
                    In the handler, via <code>authz.MustOwn</code> — it needs the record
                  </td>
                </tr>
                <tr>
                  <td>Business rules</td>
                  <td>
                    The service layer. Handlers stay thin; see{' '}
                    <Link href="/docs/backend/services" className="text-primary hover:underline">
                      Services
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Record that something happened</td>
                  <td>
                    <code>services.LogActivity</code> from the handler, after it succeeds
                  </td>
                </tr>
              </tbody>
            </table>

            <h2 id="adding">Adding your own middleware</h2>
            <p>
              A Grit middleware is an ordinary Gin one. Register it in{' '}
              <code>routes.go</code> at the position its job implies — cheap rejections high,
              anything needing an authenticated user below <code>Auth</code>:
            </p>
            <CodeBlock
              language="go"
              filename="internal/middleware/tenant.go"
              code={`func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := c.GetHeader("X-Tenant")
		if tenant == "" {
			respond.BadRequest(c, "X-Tenant header is required")
			c.Abort() // Abort, not return — return alone continues the chain
			return
		}
		c.Set("tenant", tenant)
		c.Next()
	}
}`}
            />
            <div className="mt-6 rounded-lg border border-amber-500/25 bg-amber-500/5 p-4">
              <p className="!mb-0 text-sm">
                <strong>
                  <code>c.Abort()</code> is not optional.
                </strong>{' '}
                Writing a response and returning does <em>not</em> stop the chain — the handler
                still runs, and you get a rejected request that also did the work. Any middleware
                that denies a request must call <code>c.Abort()</code>.
              </p>
            </div>
          </div>

          <div className="mt-12 flex items-center justify-between border-t border-border/40 pt-6">
            <Button asChild variant="ghost">
              <Link href="/docs/backend/middleware">
                <ArrowLeft className="mr-2 h-4 w-4" />
                Middleware
              </Link>
            </Button>
            <Button asChild variant="ghost">
              <Link href="/docs/backend/response-format">
                API Response Format
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
          </div>
        </div>
      </main>
    </div>
  )
}
