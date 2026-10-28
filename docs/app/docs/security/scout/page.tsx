import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security/scout')

const modes: { mode: string; sends: string; use: string }[] = [
  {
    mode: 'passive',
    sends: 'Nothing. Reads response headers, the TLS handshake, DNS records, your Dockerfile, your repository and the route table.',
    use: 'A scan you can run against production without telling anyone first.',
  },
  {
    mode: 'safe',
    sends: 'Harmless probes: a GET for /.git/config, a version banner, a canary string. The default.',
    use: 'Routine. Nothing it sends would be logged as an attack.',
  },
  {
    mode: 'active',
    sends: 'Real injection payloads: SQL, command, path traversal, SSRF, auth bypass.',
    use: 'A target you own, with permission, knowing it will appear in your logs as an attack.',
  },
  {
    mode: 'intrusive',
    sends: 'Brute force, load, state-changing requests.',
    use: 'A throwaway environment. Not production, not staging that anyone depends on.',
  },
]

const categories: { name: string; href: string; checks: string }[] = [
  { name: 'Headers', href: '#cat-headers', checks: 'Eight checks on what every response carries.' },
  { name: 'TLS', href: '#cat-tls', checks: 'Five checks on the handshake and the certificate.' },
  { name: 'DNS', href: '#cat-dns', checks: 'Six checks on mail authentication and dangling records.' },
  { name: 'Exposure', href: '#cat-exposure', checks: 'Seven checks for files and endpoints that should not answer.' },
  { name: 'Code', href: '#cat-code', checks: 'Seven checks over the repository itself, plus govulncheck.' },
  { name: 'Config', href: '#cat-config', checks: 'Six checks on how this deployment is configured.' },
  { name: 'Application', href: '#cat-app', checks: 'Nine that send something and read the answer, plus two on API surface.' },
  { name: 'Auth', href: '#cat-auth', checks: 'CORS, JWT claims, and the route table.' },
  { name: 'Network', href: '#cat-network', checks: 'Ports answering that nobody meant to publish.' },
]

export default function ScoutPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Security</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Scout: scanning your own app</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Scout is a security scanner that runs inside your API, the way GORM Studio and
                Pulse do. It checks the things that are true of a deployment rather than of a
                line of code: which headers you send, what your TLS offers, whether a .env is
                reachable, whether a dependency has a CVE anything actually calls.
              </p>
            </div>

            <Callout type="note" title="Opt-in, unlike Sentinel and Pulse">
              Scout is not in a new project. Run <code>grit add scout</code> when you want it.
              A scanner that appears because you ran <code>grit new</code> is one nobody decided
              to run, and an undecided scanner gets either ignored or trusted without being read.
            </Callout>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Adding it</h2>
            <CodeBlock
              language="bash"
              code={`grit add scout`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              That writes <code>apps/api/internal/routes/scout.go</code>, adds four settings to
              your config and <code>.env.example</code>, and mounts the dashboard at{' '}
              <code>/scout</code>. It installs disabled. Turn it on when you have set a password:
            </p>
            <CodeBlock
              language="bash"
              code={`# .env, beside grit.json
SCOUT_ENABLED=true
SCOUT_USERNAME=admin
SCOUT_PASSWORD=something-long
SCOUT_DB=scout.db`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              With no password it logs a line saying so and does not mount. The dashboard lists
              everywhere your application is weakest, which is the one page you least want to
              serve to anybody who finds the URL.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              <Link href="/docs/security/doctor" className="text-accent hover:underline">
                <code>grit doctor</code>
              </Link>{' '}
              reports both ways of getting this wrong: enabled with no password, where you think
              a scanner is running and none is, and enabled with a guessable one, which is an
              error in production. A project that has not added Scout reports nothing.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Running a scan</h2>
            <p className="text-muted-foreground leading-relaxed">
              There are two scans, they look at different things, and neither does the other
              one&apos;s job. Knowing which one you are running is most of knowing how to read
              the result.
            </p>
            <div className="overflow-x-auto mt-4">
              <table className="w-full text-sm table-fixed">
                <colgroup>
                  <col className="w-[22%]" />
                  <col className="w-[39%]" />
                  <col className="w-[39%]" />
                </colgroup>
                <thead>
                  <tr className="border-b border-border text-left">
                    <th className="py-2 pr-3 font-semibold"></th>
                    <th className="py-2 pr-3 font-semibold">The mounted dashboard</th>
                    <th className="py-2 font-semibold">The CLI, with <code>--repo</code></th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top font-medium">Looks at</td>
                    <td className="py-2 pr-3 align-top">
                      Your app as it is running, over HTTP, at <code>APP_URL</code>
                    </td>
                    <td className="py-2 align-top">Your source tree on disk</td>
                  </tr>
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top font-medium">Finds</td>
                    <td className="py-2 pr-3 align-top">
                      Headers, cookies, TLS, DNS, served files that should not be served, default
                      credentials, and the route-table checks that use the access table
                    </td>
                    <td className="py-2 align-top">
                      Hard-coded secrets, dangerous sinks, unpinned and vulnerable dependencies:
                      the <code>code.*</code> checks
                    </td>
                  </tr>
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top font-medium">Needs</td>
                    <td className="py-2 pr-3 align-top">
                      The app to be running. In development it also reads the source, so it
                      covers both columns
                    </td>
                    <td className="py-2 align-top">
                      Nothing running, and the <code>scout</code> binary
                    </td>
                  </tr>
                  <tr>
                    <td className="py-2 pr-3 align-top font-medium">Use it</td>
                    <td className="py-2 pr-3 align-top">
                      To see what the internet can see, and to keep a history you can diff
                    </td>
                    <td className="py-2 align-top">
                      In CI, and before a release, where there is no deployment to point at yet
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p className="text-muted-foreground leading-relaxed mt-4">
              Since v3.411.6 the dashboard runs the <code>code.*</code> checks too, in
              development. <code>grit add scout</code> hands Scout the project root, so the
              source findings and the live findings land in one scan series you can diff, rather
              than in two. Scout honours that path <strong>only</strong> when the environment is
              dev and logs a line otherwise, because a finding quotes the source line it found
              and a dashboard reachable in staging or production has no business reading the
              tree. So the CLI stays the way to get code findings anywhere but your own machine,
              and it is still the one for CI.
            </p>

            <h3 className="text-xl font-semibold mt-8 mb-3">The dashboard</h3>
            <p className="text-muted-foreground leading-relaxed">
              Start the API and open it. Sign in with the username and password you set above:
              it is not your application&apos;s login, and the two have no accounts in common.
            </p>
            <CodeBlock
              language="bash"
              code={`grit start                       # or: cd apps/api && go run cmd/server/main.go

# then open, and sign in with SCOUT_USERNAME / SCOUT_PASSWORD
http://localhost:8080/scout/ui/`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              Press <strong>Scan</strong>. Nothing runs until you do: Grit mounts Scout with no
              schedule, so it scans when asked and never on its own. The run streams as it goes,
              and when it finishes you get the findings, the attack surface, a compliance view
              and a diff against the previous scan of the same target. Every scan is kept in{' '}
              <code>SCOUT_DB</code>, which is what makes the diff possible: history is the
              feature, because the useful question is not &quot;what is wrong&quot; but
              &quot;what is new since Friday&quot;.
            </p>

            <h3 className="text-xl font-semibold mt-8 mb-3">The CLI</h3>
            <p className="text-muted-foreground leading-relaxed">
              Install it once. It is the same engine the dashboard runs, driven from a terminal.
            </p>
            <CodeBlock
              language="bash"
              code={`go install github.com/MUKE-coder/scout/cmd/scout@v0.4.9

# the code checks, from the project root. No server needed
scout scan --repo .

# the live checks, against a running app
scout scan --target http://localhost:8080

# both at once, written out as files you can read later
scout scan --repo . --target http://localhost:8080 --export md,html

# what changed since the last scan of this target
scout diff

# fail a build. Exit 1 on a finding at or above the gate, 2 on an error
scout ci --repo . --fail-on high

# the same gate, counting only what is new since the last scan. This is the
# one that works on a project that already has findings
scout ci --repo . --fail-on high --new-only

# the dashboard, without mounting anything in your app
scout serve`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              <code>scout scan --json</code> prints the findings as JSON, which is the form to
              reach for when something else has to read them.{' '}
              <code>--sarif &lt;file&gt;</code> writes SARIF 2.1.0, which GitHub code scanning
              accepts directly. <code>scout init</code> writes a <code>scout.yaml</code> holding
              every default, for when flags stop being enough, and{' '}
              <code>scout prune</code> drops scans past the retention period.
            </p>

            <h3 className="text-xl font-semibold mt-8 mb-3">Choosing a mode</h3>
            <p className="text-muted-foreground leading-relaxed">
              <code>--mode</code> decides how hard Scout pushes, and the environment caps it, so
              a mode you ask for is not always the mode you get. Safe is the default and the one
              to leave alone unless you have a reason.
            </p>
            <CodeBlock
              language="bash"
              code={`scout scan --target https://staging.example.com --mode safe       # the default
scout scan --target https://staging.example.com --mode active     # sends real payloads`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              Active against production and intrusive against staging both need two things: the
              target marked <code>allow_active</code> in your <code>scout.yaml</code>, and proof
              you own the host. Scout prints what to publish:
            </p>
            <CodeBlock
              language="bash"
              code={`scout ownership --target https://example.com`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              Publish the token it gives you as a DNS TXT record on{' '}
              <code>_scout-verify.&lt;host&gt;</code> with the value{' '}
              <code>scout-verify=&lt;token&gt;</code>, or serve it at{' '}
              <code>/.well-known/scout-verify.txt</code>, and the scan will run. A local address
              needs none of this, which is why development never asks. The point of the
              ceremony is that an active scan sends real attack payloads, and the one thing a
              scanner must never do is send those at somebody else&apos;s server.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              If you want to try the aggressive modes against something built to be attacked,{' '}
              <code>scout lab</code> clones a deliberately vulnerable app and scans it.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">
              The part Grit can give it that a scanner usually cannot
            </h2>
            <p className="text-muted-foreground leading-relaxed">
              A scanner looking at your application from outside sees URLs. It cannot see which
              of them are protected, so when it wants to report &quot;this endpoint has no
              authentication&quot; it has to infer that from the shape of the path, and the shape
              of the path does not say. <code>/api/v1/auth/login</code> and{' '}
              <code>/api/v1/auth/me</code> sit in the same place and want opposite things.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              Nor can it ask the router. Gin&apos;s route table carries each route&apos;s method,
              path and last handler, and nothing about the middleware in front of it. The running
              process genuinely does not know.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              Grit already solved this for its own Routes screen: <code>grit</code> parses your
              route files, resolves each group&apos;s middleware chain, and writes{' '}
              <code>internal/access</code>, which <Link href="/docs/security/doctor" className="text-accent hover:underline"><code>grit doctor</code></Link>{' '}
              fails on when it goes stale. <code>grit add scout</code> hands that table over, so
              Scout reads what each route requires instead of guessing:
            </p>
            <CodeBlock
              language="json"
              code={`POST /api/v1/auth/login     middleware: []
GET  /api/v1/auth/me        middleware: ["AuthRequired"]
GET  /api/v1/admin/blogs    middleware: ["AuthRequired",
                                         "RequireRole(ADMIN)",
                                         "RequirePermission(blogs.view)"]`}
            />
            <p className="text-muted-foreground leading-relaxed mt-4">
              A route that should be behind a role and is not is then a finding Scout can state
              rather than suspect, and it finds it without sending a single request.
            </p>

            <Callout type="warning" title="One rough edge, in active mode only">
              Scout currently treats any path containing <code>/api</code> as sensitive, so its
              missing-auth check reports every deliberately public route, your sign-in included,
              as lacking authentication. The check is active mode, so a default scan never runs
              it. Tracked at{' '}
              <a
                href="https://github.com/MUKE-coder/scout/issues/2"
                className="text-accent hover:underline"
                target="_blank"
                rel="noopener noreferrer"
              >
                scout#2
              </a>
              .
            </Callout>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Modes</h2>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Every check declares the most aggressive thing it does, and a scan runs only the
              checks at or below the mode you asked for. The default is safe.
            </p>
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-border">
                    <th className="text-left py-2 pr-4 font-semibold">Mode</th>
                    <th className="text-left py-2 pr-4 font-semibold">What it sends</th>
                    <th className="text-left py-2 font-semibold">When</th>
                  </tr>
                </thead>
                <tbody>
                  {modes.map((m) => (
                    <tr key={m.mode} className="border-b border-border/50">
                      <td className="py-3 pr-4 align-top font-mono text-accent">{m.mode}</td>
                      <td className="py-3 pr-4 align-top text-muted-foreground">{m.sends}</td>
                      <td className="py-3 align-top text-muted-foreground">{m.use}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <Callout type="warning" title="Active and intrusive are attacks">
              They are real payloads against a real application. Run them against something you
              own, on purpose, knowing your logs will record an attack because one happened.
            </Callout>

            <h2 className="text-2xl font-semibold mt-12 mb-4">What it checks</h2>
            <div className="grid gap-3 sm:grid-cols-2">
              {categories.map((c) => (
                <a
                  key={c.name}
                  href={c.href}
                  className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 block hover:border-accent/40 transition-colors"
                >
                  <h3 className="text-sm font-semibold text-foreground mb-1">{c.name}</h3>
                  <p className="text-xs text-muted-foreground/60 leading-relaxed">{c.checks}</p>
                </a>
              ))}
            </div>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Every category, in detail</h2>
            <p className="text-muted-foreground leading-relaxed mb-8">
              Fifty-four checks. The names below are the ids Scout reports findings under, so a
              finding on the dashboard can be traced back to what produced it. The mode in
              brackets is the lowest mode that runs the check.
            </p>

            <h3 id="cat-headers" className="text-xl font-semibold mb-3">
              Headers <span className="text-sm font-normal text-muted-foreground">(passive)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Reads what your responses already carry; sends nothing unusual.{' '}
              <code>headers.csp</code>, <code>headers.hsts</code>,{' '}
              <code>headers.clickjacking</code>, <code>headers.nosniff</code>,{' '}
              <code>headers.referrer-policy</code>, <code>headers.permissions-policy</code>,{' '}
              <code>headers.cookie-flags</code> and <code>headers.server-disclosure</code>.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Cookie flags is the one that finds real bugs most often: a session cookie without{' '}
              <code>HttpOnly</code> is readable by any script on the page, and one without{' '}
              <code>Secure</code> goes out over plain HTTP. Expect one finding you will keep:
              Grit&apos;s <code>grit_csrf</code> cookie is deliberately not <code>HttpOnly</code>,
              because the double-submit pattern needs JavaScript to read it and echo it back.
              Mark it accepted and it stays marked.
            </p>

            <h3 id="cat-tls" className="text-xl font-semibold mb-3">
              TLS <span className="text-sm font-normal text-muted-foreground">(passive)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>tls.weak-protocols</code> and <code>tls.weak-ciphers</code> report anything
              old enough to be a problem. <code>tls.cert-chain</code> catches a chain that your
              browser completes from its own cache but a Go client or a phone will not, which is
              the classic &quot;works for me&quot; outage. <code>tls.cert-expiry</code> gives you
              the renewal date before your users do, and <code>tls.hsts-preload</code> says
              whether the domain qualifies for the browser preload list.
            </p>

            <h3 id="cat-dns" className="text-xl font-semibold mb-3">
              DNS <span className="text-sm font-normal text-muted-foreground">(passive)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>dns.spf</code>, <code>dns.dkim</code> and <code>dns.dmarc</code> are about
              whether somebody else can send mail as you, which is phishing with your name on
              it. <code>dns.caa</code> is which authorities may issue certificates for the
              domain.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>dns.dangling-cname</code> and <code>dns.subdomain-takeover</code> are the
              pair worth understanding: a CNAME pointing at a service you stopped paying for can
              be claimed by whoever registers that name next, and then they are serving content
              on your subdomain, with your cookies in scope.
            </p>

            <h3 id="cat-exposure" className="text-xl font-semibold mb-3">
              Exposure <span className="text-sm font-normal text-muted-foreground">(safe)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Asks for things that should not answer: <code>exposure.git</code> (a served{' '}
              <code>.git</code> directory is your whole source history),{' '}
              <code>exposure.dotenv</code>, <code>exposure.backup-files</code>,{' '}
              <code>exposure.dir-listing</code>, <code>exposure.swagger</code>,{' '}
              <code>exposure.metrics</code> and <code>exposure.dashboards</code>.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              These are plain GETs for well-known paths. Nothing here is an attack; a 404 is the
              answer it hopes for.
            </p>

            <h3 id="cat-code" className="text-xl font-semibold mb-3">
              Code and dependencies{' '}
              <span className="text-sm font-normal text-muted-foreground">(passive)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Reads the repository rather than the running app.{' '}
              <code>code.secrets</code> finds credentials committed to the tree,{' '}
              <code>code.secret-logging</code> finds them being written to logs, which is the
              same leak with a longer fuse. <code>code.dangerous-sink</code>,{' '}
              <code>code.insecure-random</code> and <code>code.weak-crypto</code> are the static
              checks: <code>math/rand</code> where a token is being generated, a hash that should
              not still be in use. <code>code.unpinned-deps</code> covers the supply chain, and{' '}
              <code>code.dockerfile</code> covers image practice (filed under config, since it
              is about how this is built and run rather than what the code says).
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>code.go-vulncheck</code> runs govulncheck over the module and reports{' '}
              <strong>reachable</strong> CVEs: a vulnerability in a function nothing in your code
              path calls is not a finding, which is what makes the output short enough to read.
              It degrades quietly, so a repo that is not a Go module, or a machine without
              govulncheck installed, reports nothing rather than failing the scan.
            </p>

            <h3 id="cat-config" className="text-xl font-semibold mb-3">
              Config <span className="text-sm font-normal text-muted-foreground">(passive)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The checks that know what framework this is.{' '}
              <code>config.gin-debug</code> catches Gin left in debug mode, which prints routes
              and stack traces. <code>config.default-creds</code> catches a dashboard still on
              what it shipped with. <code>config.inmemory-storage-prod</code> catches a store
              that loses everything on restart being used where that matters.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Two are specifically about the tools beside it.{' '}
              <code>config.sentinel-misconfig</code> reads{' '}
              <code>/sentinel/api/health</code> and reports a WAF left in log-only mode, which
              detects attacks and allows them through: protection that exists on the dashboard
              and not in the request path. <code>config.pulse-misconfig</code> reports Pulse
              running in dev mode against a production target.
            </p>

            <h3 id="cat-app" className="text-xl font-semibold mb-3">
              Application{' '}
              <span className="text-sm font-normal text-muted-foreground">(safe and active)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              The ones that send something and read the answer. Safe mode covers{' '}
              <code>app.reflected-xss</code>, <code>app.open-redirect</code>,{' '}
              <code>app.csrf-form</code>, <code>app.info-disclosure</code>,{' '}
              <code>app.host-header</code>, plus <code>api.graphql-introspection</code> and{' '}
              <code>api.improper-asset-mgmt</code>, which is an old API version still answering
              after everyone moved to the new one.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              Active mode adds <code>app.sqli</code>, <code>app.cmd-injection</code>,{' '}
              <code>app.path-traversal</code> and <code>app.ssrf</code>. These send real
              payloads. Your logs will record an attack, because one happened, and Sentinel will
              block some of them, which is the correct outcome and makes the results harder to
              read. Run them against something you own, deliberately.
            </p>

            <h3 id="cat-auth" className="text-xl font-semibold mb-3">
              Auth{' '}
              <span className="text-sm font-normal text-muted-foreground">(passive and active)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>auth.cors</code> reads your CORS policy for the combination that undoes it:
              reflecting any origin while allowing credentials means any site can make
              authenticated requests as your signed-in user. <code>auth.jwt</code> inspects a
              token you supply, passively, for <code>alg: none</code>, a missing signature and a
              missing expiry. It does not try to crack anything.
            </p>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>auth.missing</code> is the route-table check described above, and the one
              this integration exists for. Two checks are deliberately not registered at all,{' '}
              <code>auth.idor</code> and <code>auth.mass-assignment</code>, because both need two
              configured test identities to compare what each can reach, and shipping them
              without that would mean a check that always passes.
            </p>

            <h3 id="cat-network" className="text-xl font-semibold mb-3">
              Network <span className="text-sm font-normal text-muted-foreground">(safe)</span>
            </h3>
            <p className="text-muted-foreground leading-relaxed mb-4">
              <code>network.exposed-services</code> checks the common ports for something
              answering that nobody meant to publish: a database, a cache, an admin port that was
              bound to <code>0.0.0.0</code> instead of <code>127.0.0.1</code>. It is a short list
              of well-known ports, not a full port sweep.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Findings you have already decided about</h2>
            <p className="text-muted-foreground leading-relaxed">
              A finding can be marked accepted-risk with a date, or a false positive. An accepted
              risk reopens itself when the date passes, which is the difference between deciding
              to live with something and forgetting about it.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">What a clean project scores</h2>
            <p className="text-muted-foreground leading-relaxed">
              Worth knowing before your first scan, because a number with nothing to compare it
              to tells you very little. A freshly generated monorepo, scanned with{' '}
              <code>scout scan --repo . --mode safe</code> on Scout v0.4.10, reports 66
              findings. Not one of them is a way in. They break down like this.
            </p>
            <div className="overflow-x-auto mt-4">
              <table className="w-full text-sm table-fixed">
                <colgroup>
                  <col className="w-[12%]" />
                  <col className="w-[26%]" />
                  <col className="w-[62%]" />
                </colgroup>
                <thead>
                  <tr className="border-b border-border text-left">
                    <th className="py-2 pr-3 font-semibold">Count</th>
                    <th className="py-2 pr-3 font-semibold">Check</th>
                    <th className="py-2 font-semibold">What it is</th>
                  </tr>
                </thead>
                <tbody className="text-muted-foreground">
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top">62</td>
                    <td className="py-2 pr-3 align-top"><code>code.unpinned-deps</code></td>
                    <td className="py-2 align-top">
                      A caret range in a <code>package.json</code>. A real decision rather than a
                      bug: a range means a fresh install can pick up a version nobody tested.
                    </td>
                  </tr>
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top">1</td>
                    <td className="py-2 pr-3 align-top"><code>code.secrets</code></td>
                    <td className="py-2 align-top">
                      The demo users&apos; shared <code>admin123</code> in{' '}
                      <code>users_seeder.go</code>. The literal is real and the seeder refuses to
                      run outside <code>APP_ENV=development</code>, which a pattern match cannot
                      see. Accept it.
                    </td>
                  </tr>
                  <tr className="border-b border-border/40">
                    <td className="py-2 pr-3 align-top">2</td>
                    <td className="py-2 pr-3 align-top"><code>code.vulnerable-deps</code></td>
                    <td className="py-2 align-top">
                      <code>golang.org/x/crypto</code> (the advisory that{' '}
                      <code>openpgp</code> is unmaintained, and nothing imports{' '}
                      <code>openpgp</code>) and <code>disintegration/imaging</code> (its TIFF
                      decoder, which is never registered). Both are in the module graph and
                      neither is in the build. Neither has a fixed version to move to.
                    </td>
                  </tr>
                  <tr>
                    <td className="py-2 pr-3 align-top">1</td>
                    <td className="py-2 pr-3 align-top"><code>code.secret-logging</code></td>
                    <td className="py-2 align-top">
                      The API key seeder printing the secret it just minted. That line is inside
                      a development-only branch as of v3.411.0, and the check reads the call
                      rather than the branch around it. Accept it.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            <p className="text-muted-foreground leading-relaxed mt-4">
              The useful reading of that table is the ratio. A scanner is worth running when a
              new finding stands out, and it stops being worth running when there are so many
              that nobody looks. If your own scan is much longer than this, the extra findings
              are yours and worth reading one at a time.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              It was 73 on v0.4.9. Seven of those were signing keys in test files, which is not
              a finding anyone can act on: a test that signs a token has to contain a key, so
              there is no version of that file without one. Scout skips test fixtures as of
              v0.4.10.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              That is a repo scan. Point Scout at a running app with{' '}
              <code>--target</code> and you also get the <code>grit_csrf</code> cookie finding
              described above, which is the one deliberate exception in the cookie checks.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">What it is not</h2>
            <p className="text-muted-foreground leading-relaxed">
              Scout is not <Link href="/docs/batteries/security" className="text-accent hover:underline">Sentinel</Link>.
              Sentinel is a WAF: it sits in the request path and blocks things as they happen.
              Scout runs on a schedule and tells you what it found. One defends, the other looks.
              They answer different questions and the answers do not substitute for each other.
            </p>
            <p className="text-muted-foreground leading-relaxed mt-4">
              It is also not a penetration test. It finds the classes of problem a tool can find.
              Business logic, a broken permission model, an authorization rule that is right for
              every case but one: those need somebody thinking about your application, and
              nothing here replaces that.
            </p>

            <div className="flex items-center justify-between mt-16 pt-8 border-t border-border">
              <Button variant="ghost" asChild>
                <Link href="/docs/security/doctor">
                  <ArrowLeft className="mr-2 h-4 w-4" />
                  grit doctor
                </Link>
              </Button>
              <Button variant="ghost" asChild>
                <Link href="/docs/security/compliance">
                  Privacy &amp; Compliance
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
