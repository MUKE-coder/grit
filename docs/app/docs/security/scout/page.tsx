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

const categories: { name: string; checks: string }[] = [
  { name: 'Headers', checks: 'CSP, HSTS, clickjacking, nosniff, referrer policy, permissions policy, cookie flags, server version disclosure.' },
  { name: 'TLS', checks: 'Weak protocols and ciphers, certificate chain and expiry, HSTS preload eligibility.' },
  { name: 'DNS', checks: 'SPF, DKIM, DMARC, CAA, dangling CNAMEs and subdomain takeover.' },
  { name: 'Exposure', checks: 'A reachable .git, a .env, directory listings, backup files, an open Swagger or metrics endpoint, a dashboard with no login.' },
  { name: 'Code', checks: 'Secrets committed to the repository, dangerous sinks, insecure randomness, weak crypto, unpinned dependencies, secrets reaching the logs.' },
  { name: 'Dependencies', checks: 'Reachable CVEs in your Go modules, via govulncheck. Reachable, not merely present: a CVE in code nothing calls is not a finding.' },
  { name: 'Config', checks: 'Gin left in debug mode, default credentials, in-memory storage in production, Sentinel not enforcing, Pulse in dev mode.' },
  { name: 'Application', checks: 'Reflected XSS, open redirect, SQL and command injection, path traversal, SSRF, CSRF forms, host-header injection. Most of these are active mode.' },
  { name: 'Auth', checks: 'CORS policy, JWT header and claims, and the route table check below.' },
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
              <span className="tag-mono text-primary/80 mb-3 block">Security</span>
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
              code={`# apps/api/.env
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
              <code>internal/access</code>, which <Link href="/docs/security/doctor" className="text-primary hover:underline"><code>grit doctor</code></Link>{' '}
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
                className="text-primary hover:underline"
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
                      <td className="py-3 pr-4 align-top font-mono text-primary">{m.mode}</td>
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
            <div className="space-y-4">
              {categories.map((c) => (
                <div key={c.name} className="border-l-2 border-border pl-4">
                  <h3 className="font-semibold mb-1">{c.name}</h3>
                  <p className="text-sm text-muted-foreground leading-relaxed">{c.checks}</p>
                </div>
              ))}
            </div>

            <h2 className="text-2xl font-semibold mt-12 mb-4">Findings you have already decided about</h2>
            <p className="text-muted-foreground leading-relaxed">
              A finding can be marked accepted-risk with a date, or a false positive. An accepted
              risk reopens itself when the date passes, which is the difference between deciding
              to live with something and forgetting about it.
            </p>

            <h2 className="text-2xl font-semibold mt-12 mb-4">What it is not</h2>
            <p className="text-muted-foreground leading-relaxed">
              Scout is not <Link href="/docs/batteries/security" className="text-primary hover:underline">Sentinel</Link>.
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
