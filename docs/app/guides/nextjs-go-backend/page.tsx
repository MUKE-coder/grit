import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  H3,
  Outcome,
  Note,
  Warn,
  Mistake,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Connect a Next.js frontend to a Go backend',
  description:
    'CORS, credentials and SameSite explained: why the browser throws away a response your Go API sent correctly, and why the login cookie disappears.',
  alternates: { canonical: 'https://gritframework.dev/guides/nextjs-go-backend' },
  openGraph: {
    title: 'Connect a Next.js frontend to a Go backend',
    description: 'Why the request works in Postman and fails in the browser, and how to fix it.',
    url: 'https://gritframework.dev/guides/nextjs-go-backend',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="nextjs-go-backend"
      lead={
        <>
          <p className="mb-4">
            Your Go API answers <code>curl</code> and Postman perfectly, and the browser throws its
            response away. Then you fix that, the login returns 200, and the very next request is
            still unauthenticated. Both are the same subject, and neither is a bug in your code.
          </p>
          <p>
            A browser enforces rules that no other HTTP client does, because a browser carries your
            users&apos; credentials and is pointed at untrusted pages all day. This explains the
            three rules that matter, in the order you hit them, and ends with a configuration you
            can deploy.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A Go API a browser will read from, with the allowed origin coming from configuration.</p>
        <p>
          A login whose cookie is actually sent back on the next request, in development across two
          ports and in production across two subdomains.
        </p>
        <p>
          Enough understanding of <code>SameSite</code> to predict which of your deployments will
          break before it does.
        </p>
      </Outcome>

      <Section id="origins" title="An origin is scheme, host and port">
        <p>
          Two URLs share an origin when all three match. These are all different origins from{' '}
          <code>http://localhost:3000</code>:
        </p>
        <Checklist
          items={[
            <>
              <code>http://localhost:8080</code>, because the port differs. This is the one that
              catches you in development.
            </>,
            <>
              <code>https://localhost:3000</code>, because the scheme differs.
            </>,
            <>
              <code>http://127.0.0.1:3000</code>, because the host string differs, even though it
              resolves to the same machine.
            </>,
            <>
              <code>https://api.example.com</code> against <code>https://example.com</code>, because
              the host differs. Subdomains are separate origins.
            </>,
          ]}
        />
        <p>
          That last one matters more than people expect: the common production layout, frontend on
          the apex and API on a subdomain, is cross-origin, so everything on this page applies in
          production and not only locally.
        </p>
      </Section>

      <Section id="cors" title="Rule one: the response is withheld unless the server allows it">
        <p>Here is the error, as it appears in the console:</p>
        <CodeBlock
          language="text"
          code={`Access to fetch at 'http://localhost:8080/api/users' from origin
'http://localhost:3000' has been blocked by CORS policy: No
'Access-Control-Allow-Origin' header is present on the requested resource.`}
        />
        <p>
          Read that carefully, because the wording hides what happened. The request{' '}
          <strong>was sent</strong>. Your server <strong>did run the handler</strong>. If it was a{' '}
          <code>POST</code> that created a row, the row exists. The browser received the response
          and refused to hand it to your JavaScript, because the response did not say which origins
          are allowed to read it.
        </p>
        <p>
          So CORS is not a firewall and does not protect your API. It protects the user&apos;s other
          tabs: it is what stops a page on a random site reading your bank&apos;s API using your
          logged-in session. Authorisation is still entirely your job.
        </p>

        <H3>The preflight</H3>
        <p>
          For anything beyond a simple <code>GET</code> or form post, the browser asks permission
          first, with an <code>OPTIONS</code> request:
        </p>
        <CodeBlock
          language="text"
          code={`OPTIONS /api/users HTTP/1.1
Origin: http://localhost:3000
Access-Control-Request-Method: POST
Access-Control-Request-Headers: content-type, authorization`}
        />
        <p>
          Your server must answer that with the permissions and no body. If it instead runs your
          route handler, or returns a 404 because you only registered <code>POST</code>, the real
          request never happens. A <code>Content-Type: application/json</code> header is enough on
          its own to trigger a preflight, which is why it appears the moment you start sending
          JSON.
        </p>

        <H3>The middleware, written out</H3>
        <CodeBlock
          filename="internal/middleware/cors.go"
          code={`package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS allows exactly the origins given, which come from configuration
// rather than from this file.
//
// There is no wildcard branch on purpose. A browser refuses to send
// credentials to a wildcard origin, so the moment cookies are involved the
// wildcard silently stops working, and the failure looks like an auth bug
// rather than a CORS one.
func CORS(allowed []string) gin.HandlerFunc {
	set := make(map[string]bool, len(allowed))
	for _, o := range allowed {
		set[strings.TrimSpace(o)] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && set[origin] {
			// Echo the caller's origin back. Not a list: the header takes
			// exactly one value, and browsers reject a comma-separated one.
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			// Caches key on this, so a response cached for one origin is not
			// served to another. Leaving it out is a real bug behind a CDN.
			c.Header("Vary", "Origin")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		// How long the browser may remember this answer, so it stops asking
		// before every request.
		c.Header("Access-Control-Max-Age", "600")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}`}
        />
        <p>Register it before your routes, with the origins read from the environment:</p>
        <CodeBlock
          filename="cmd/server/main.go"
          code={`r := gin.Default()
r.Use(middleware.CORS(strings.Split(os.Getenv("CORS_ORIGINS"), ",")))
// routes after this line`}
        />
        <CodeBlock
          language="bash"
          filename=".env"
          code={`CORS_ORIGINS=http://localhost:3000
# production:
# CORS_ORIGINS=https://app.example.com,https://admin.example.com`}
        />
        <Warn title="Middleware registered after a route does not run for it">
          <p>
            Gin builds the handler chain when the route is registered. A <code>r.Use(...)</code>{' '}
            below your <code>r.GET(...)</code> applies to nothing above it, and this is quiet: no
            error, no warning, just a route that behaves differently from its neighbours.
          </p>
        </Warn>
      </Section>

      <Section id="credentials" title="Rule two: cookies are not sent unless you ask">
        <p>
          CORS is fixed, the login endpoint returns 200 and sets a cookie, and the next request is
          still unauthenticated. Three separate things have to be true, and missing any one of them
          produces exactly this.
        </p>

        <H3>The browser must be told to send credentials</H3>
        <p>
          By default <code>fetch</code> sends no cookies cross-origin, and attaches none from the
          response either. You opt in per request:
        </p>
        <CodeBlock
          language="ts"
          code={`const res = await fetch(\`\${API}/auth/login\`, {
  method: 'POST',
  credentials: 'include',          // <- without this, nothing else matters
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password }),
})`}
        />
        <p>
          Set it once, in your API client, rather than on each call. A single <code>fetch</code>{' '}
          somewhere that forgets it is a bug that only shows up on one screen.
        </p>
        <CodeBlock
          language="ts"
          filename="lib/api.ts"
          code={`const API = process.env.NEXT_PUBLIC_API_URL!

export async function api<T>(path: string, init: RequestInit = {}): Promise<T> {
  const res = await fetch(\`\${API}\${path}\`, {
    ...init,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...init.headers },
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error(body?.error?.message ?? \`request failed: \${res.status}\`)
  }
  return res.json()
}`}
        />

        <H3>The server must allow credentials, and name the origin</H3>
        <p>
          <code>Access-Control-Allow-Credentials: true</code>, which the middleware above sends, and{' '}
          <strong>not</strong> a wildcard origin. A browser refuses the combination of{' '}
          <code>credentials: &apos;include&apos;</code> and{' '}
          <code>Access-Control-Allow-Origin: *</code> outright, and the console message talks about
          the wildcard rather than about your cookie, so the connection is easy to miss.
        </p>

        <H3>The cookie&apos;s own attributes must permit it</H3>
        <p>This is where most of the remaining time goes.</p>
        <CodeBlock
          filename="internal/handlers/auth.go"
          code={`http.SetCookie(c.Writer, &http.Cookie{
	Name:     "refresh_token",
	Value:    token,
	Path:     "/api/auth",
	HttpOnly: true,                    // JavaScript cannot read it, which is the point
	Secure:   cfg.Env == "production", // required whenever SameSite is None
	SameSite: sameSite,                // see below
	MaxAge:   int(30 * 24 * time.Hour / time.Second),
})`}
        />
        <p>
          <code>SameSite</code> decides whether the browser attaches the cookie to a request
          triggered by a different site. The three values, and what each means for you:
        </p>
        <Checklist
          items={[
            <>
              <code>Lax</code>, the default in modern browsers. Sent on top-level navigation, not on
              a cross-site <code>fetch</code>. This is why a cookie that exists in devtools is still
              not sent.
            </>,
            <>
              <code>Strict</code>. Never sent cross-site at all, including when the user clicks a
              link from an email into your app, which logs them out at the worst moment.
            </>,
            <>
              <code>None</code>. Always sent, and <strong>requires</strong> <code>Secure</code>,
              which requires HTTPS. A browser silently drops a <code>SameSite=None</code> cookie
              that is not <code>Secure</code>.
            </>,
          ]}
        />
        <p>
          There is one more rule nobody tells you:{' '}
          <strong>
            <code>localhost:3000</code> and <code>localhost:8080</code> are the <em>same</em> site
          </strong>
          , even though they are different origins. SameSite is judged on the registrable domain,
          which ignores the port. So a <code>Lax</code> cookie works fine across your two
          development ports, and then fails in production when the frontend and API are on
          different domains.
        </p>
        <Note title="Pick by deployment, not by environment">
          <p>
            If your frontend and API share a registrable domain in production,{' '}
            <code>app.example.com</code> and <code>api.example.com</code>, then <code>Lax</code>{' '}
            works and is the safer choice. If they are on genuinely different domains, you need{' '}
            <code>None</code> with <code>Secure</code>, and therefore HTTPS everywhere.
          </p>
          <p>
            Deciding this from <code>NODE_ENV</code> rather than from the actual deployment is how
            a configuration works locally and in staging and fails in production.
          </p>
        </Note>
        <p>
          Also set <code>Domain=.example.com</code> when the cookie must be readable across
          subdomains. Without it, a cookie set by <code>api.example.com</code> is not sent to{' '}
          <code>app.example.com</code> at all.
        </p>
      </Section>

      <Section id="server-side" title="Rule three: server components are not in the browser">
        <p>
          A Next.js server component runs on the server, where there is no browser, no ambient
          cookie jar and no CORS. Three consequences people trip over.
        </p>
        <Checklist
          items={[
            <>
              <strong>CORS does not apply.</strong> A fetch from a server component to your Go API
              is server to server. If that one works and the client one does not, you have a CORS
              problem, and you have just proved it.
            </>,
            <>
              <strong>Cookies are not attached automatically.</strong> You have to read them and
              forward them yourself.
            </>,
            <>
              <strong>
                <code>localhost</code> may not mean the same thing.
              </strong>{' '}
              In Docker, the server component resolves <code>localhost</code> inside its own
              container. It needs the service name, while the browser needs the public URL. Two
              variables, not one.
            </>,
          ]}
        />
        <CodeBlock
          language="ts"
          filename="app/dashboard/page.tsx"
          code={`import { cookies } from 'next/headers'

export default async function Dashboard() {
  const jar = await cookies()

  const res = await fetch(\`\${process.env.API_URL_INTERNAL}/api/me\`, {
    headers: { cookie: jar.toString() },   // forward it by hand
    cache: 'no-store',                     // never cache a per-user response
  })

  const { data } = await res.json()
  return <h1>Hello {data.name}</h1>
}`}
        />
        <CodeBlock
          language="bash"
          filename=".env"
          code={`# What the browser calls. Must be public and reachable from outside.
NEXT_PUBLIC_API_URL=http://localhost:8080

# What the Next.js server calls. Inside Docker this is the service name.
API_URL_INTERNAL=http://api:8080`}
        />
        <Warn title="NEXT_PUBLIC_ is compiled into the bundle">
          <p>
            Any variable with that prefix is substituted at build time and shipped to the browser,
            where anyone can read it. It is right for an API URL and wrong for anything else. A
            secret with that prefix is a published secret, and the build will not warn you.
          </p>
        </Warn>
      </Section>

      <Section id="proxy" title="The alternative: make it same-origin">
        <p>
          Everything above exists because the two are on different origins. You can simply not do
          that: have Next.js proxy <code>/api</code> through to Go, so the browser only ever talks
          to one origin.
        </p>
        <CodeBlock
          language="ts"
          filename="next.config.ts"
          code={`const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: \`\${process.env.API_URL_INTERNAL}/api/:path*\`,
      },
    ]
  },
}

export default nextConfig`}
        />
        <p>
          Now the browser fetches <code>/api/users</code>, same origin, no CORS, no{' '}
          <code>credentials: &apos;include&apos;</code>, no SameSite question. Cookies are
          first-party.
        </p>
        <p>
          What it costs: every API request now goes through the Node process, which adds a hop and
          makes Node a bottleneck under load. Streaming responses and file uploads need care. And
          mobile or desktop clients cannot use the proxy, so the direct path has to work anyway.
        </p>
        <p>
          A reasonable rule: proxy if the web app is the only client, go direct with proper CORS if
          there will be others.
        </p>
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="No 'Access-Control-Allow-Origin' header is present on the requested resource">
          <p>
            The server did not send the header. Three things to check in order: that the middleware
            is registered <strong>before</strong> the routes; that the origin in{' '}
            <code>CORS_ORIGINS</code> matches <strong>exactly</strong>, including scheme and port
            and with no trailing slash; and that you restarted the Go server, which does not reload
            itself.
          </p>
        </Mistake>

        <Mistake error="Response to preflight request doesn't pass access control check">
          <p>
            The <code>OPTIONS</code> request got an answer the browser did not accept. Usually your
            router returned 404 or 405 because only <code>POST</code> is registered for that path,
            so the middleware never ran. Handle <code>OPTIONS</code> in middleware and abort there,
            as above.
          </p>
          <CodeBlock
            terminal
            code={`curl -i -X OPTIONS http://localhost:8080/api/users \\
  -H "Origin: http://localhost:3000" \\
  -H "Access-Control-Request-Method: POST"
# want: 204, with Access-Control-Allow-Origin echoing your origin`}
          />
        </Mistake>

        <Mistake error="The value of the 'Access-Control-Allow-Origin' header must not be the wildcard '*' when the request's credentials mode is 'include'">
          <p>
            Exactly what it says, and the fix is to echo the caller&apos;s origin rather than send a
            wildcard. This is the error that tells you your wildcard CORS was always going to break
            the day you added login.
          </p>
        </Mistake>

        <Mistake error="Request header field authorization is not allowed by Access-Control-Allow-Headers">
          <p>
            The preflight listed a header your allow-list does not name. Add it. The list is
            literal: a header you forgot is a request the browser refuses to send, and it names the
            missing one for you.
          </p>
        </Mistake>

        <Mistake error="The login returns 200 and the next request is still unauthenticated">
          <p>
            Work through it in this order, because each step makes the next observable.
          </p>
          <Checklist
            items={[
              <>
                Devtools, Network, the login response. Is there a <code>Set-Cookie</code> header at
                all? If not, the server never set it.
              </>,
              <>
                Devtools, Application, Cookies. Is the cookie stored? If the response had it and the
                browser did not keep it, the attributes are wrong:{' '}
                <code>SameSite=None</code> without <code>Secure</code> is dropped in silence.
              </>,
              <>
                The next request, Request Headers. Is the cookie sent? If it is stored and not sent,
                it is <code>SameSite</code>, or a <code>Path</code> narrower than the path you are
                calling.
              </>,
              <>
                Is <code>credentials: &apos;include&apos;</code> on that particular call? One
                forgotten <code>fetch</code> is enough.
              </>,
            ]}
          />
        </Mistake>

        <Mistake error="It works on localhost and fails in production">
          <p>
            The usual cause is <code>SameSite=Lax</code>. Your two development ports are the same
            site, so <code>Lax</code> is fine there, and your production frontend and API are on
            different domains, where it is not. The same mechanism also explains the reverse: a
            config built for production, with <code>SameSite=None; Secure</code>, drops the cookie
            on plain-HTTP localhost.
          </p>
        </Mistake>
      </Section>

      <Section id="with-grit" title="How a generated project does it">
        <p>
          A Grit project ships this configured. The API reads its allowed origins from{' '}
          <code>CORS_ORIGINS</code>, the client sets <code>credentials: &apos;include&apos;</code>{' '}
          in one place, the refresh cookie is <code>HttpOnly</code> and scoped to{' '}
          <code>/api/auth</code>, and the access token is held in memory rather than in{' '}
          <code>localStorage</code>, where any script on the page can read it.
        </p>
        <p>
          The scoped <code>Path</code> is worth copying even if you take nothing else: a refresh
          token that is only sent to the endpoints that refresh it is not exposed on the other
          several hundred requests your app makes.
        </p>
        <p>
          <Link href="/docs/backend/authentication" className="text-accent hover:underline">
            The auth reference
          </Link>{' '}
          has the whole scheme, including rotation and replay detection.
        </p>
      </Section>

      <Faq
        slug="nextjs-go-backend"
        items={[
          {
            q: 'Why does it work in Postman but not in the browser?',
            a: (
              <>
                Because CORS is a browser rule and nothing else implements it. Postman, curl, your
                Go tests and a native mobile app all ignore it entirely. A tool that bypasses the
                rule cannot tell you whether you have satisfied it, which is why &quot;it works in
                Postman&quot; is not evidence of anything here.
              </>
            ),
            text:
              'Because CORS is enforced by browsers and nothing else. Postman, curl, Go tests and native mobile apps ignore it. A tool that bypasses the rule cannot tell you whether you satisfy it.',
          },
          {
            q: 'Should I store the JWT in localStorage?',
            a: (
              <>
                No. Any script running on your page can read <code>localStorage</code>, including
                one that arrived through a compromised dependency, so a single cross-site scripting
                hole becomes every user&apos;s token. Keep the refresh token in an{' '}
                <code>HttpOnly</code> cookie, which JavaScript cannot read, and hold the short-lived
                access token in memory.
              </>
            ),
            text:
              'No. Any script on the page can read localStorage, including one from a compromised dependency, so one XSS hole exposes every token. Keep the refresh token in an HttpOnly cookie that JavaScript cannot read, and hold the short-lived access token in memory.',
          },
          {
            q: 'Is a wildcard CORS origin ever acceptable?',
            a: (
              <>
                Only for an API that is genuinely public and takes no credentials: a read-only
                endpoint serving data anyone may have. The moment a cookie or an Authorization
                header is involved, the browser refuses the wildcard anyway, so the choice is made
                for you.
              </>
            ),
            text:
              'Only for a genuinely public API that takes no credentials, such as a read-only endpoint serving public data. Once cookies or an Authorization header are involved the browser refuses a wildcard, so the decision is made for you.',
          },
          {
            q: 'Should I proxy through Next.js or call the Go API directly?',
            a: (
              <>
                Proxy if the web app is your only client: it removes CORS, cookies become
                first-party, and there is less to get wrong. Call directly if you will also have a
                mobile app, a desktop app or third-party consumers, because those cannot use the
                proxy and the direct path has to be correct regardless.
              </>
            ),
            text:
              'Proxy if the web app is the only client: CORS disappears and cookies become first-party. Call directly if there will also be mobile, desktop or third-party clients, since they cannot use the proxy and the direct path must be correct anyway.',
          },
          {
            q: 'Why can my server component reach the API when my browser code cannot?',
            a: (
              <>
                Because the server component runs on the server, where CORS does not exist. That
                asymmetry is a useful diagnostic: if the server fetch works and the client fetch
                does not, the problem is CORS and not your API. Note that the two may also need
                different URLs, since inside Docker the server reaches the API by service name.
              </>
            ),
            text:
              'Because a server component runs on the server, where CORS does not apply. That is a useful diagnostic: if the server fetch works and the client fetch does not, the problem is CORS rather than the API. The two may also need different URLs, since inside Docker the server reaches the API by service name.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'Where this problem first appears, with the smallest fix that gets you moving.',
          },
          {
            href: '/guides/common-go-docker-errors',
            label: 'Fix the errors you hit starting out',
            blurb: 'The exact messages, with the cause and the fix per operating system.',
          },
          {
            href: '/docs/backend/authentication',
            label: 'Authentication reference',
            blurb: 'Sessions, refresh rotation and replay detection as a generated project does them.',
          },
          {
            href: '/docs/frontend/fetching-data',
            label: 'Fetching data in the frontend',
            blurb: 'The generated client and hooks, where credentials and the base URL are decided.',
          },
        ]}
      />
    </GuidePage>
  )
}
