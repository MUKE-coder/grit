import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/security')

export default function SecurityPage() {
  return (
    <DocsShell>
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-accent/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Security (Sentinel)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit project ships with Sentinel -- a production-grade security intelligence suite
                that provides a Web Application Firewall, rate limiting, brute-force protection, anomaly detection,
                security headers, and a real-time threat dashboard. Security is not an afterthought.
              </p>
              <LaneFlow
                id="bat-security"
                lanes={['Incoming request', 'Sentinel WAF', 'Outcome']}
                nodes={[
                  { id: 'req', lane: 0, row: 1, title: 'Request', sub: 'client / attacker', tone: 'blue' },
                  { id: 'checks', lane: 1, row: 1, title: 'Sentinel WAF', sub: 'rate · brute · anomaly', tone: 'primary' },
                  { id: 'allow', lane: 2, row: 0, title: 'Allow', sub: '→ handler', tone: 'green' },
                  { id: 'block', lane: 2, row: 1, title: 'Block', sub: '429 / 403', tone: 'rose' },
                  { id: 'dash', lane: 2, row: 2, title: 'Threat dashboard', sub: 'real-time', tone: 'amber' },
                ]}
                edges={[
                  { from: 'req', to: 'checks', label: 'inspect', tone: 'blue' },
                  { from: 'checks', to: 'allow', label: 'clean', tone: 'green' },
                  { from: 'checks', to: 'block', label: 'threat', tone: 'rose' },
                  { from: 'block', to: 'dash', label: 'log', dashed: true, tone: 'amber' },
                ]}
                legend={[
                  { tone: 'green', label: 'Allowed' },
                  { tone: 'rose', label: 'Blocked' },
                  { tone: 'amber', label: 'Streamed to dashboard' },
                ]}
                caption="Every request passes Sentinel's WAF; threats are blocked and streamed to the dashboard"
              />
            </div>

            <div className="prose-grit">
              {/* What's Included */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What&apos;s Included
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sentinel mounts on your Gin router with a single call and provides the following out of the box:
                </p>

                <div className="grid gap-3 sm:grid-cols-2 mb-6">
                  {[
                    { title: 'Web Application Firewall', href: '#waf', desc: 'Fourteen threat classes, matched against the URL, query, headers, cookies and body' },
                    { title: 'Rate Limiting', href: '#rate-limits', desc: 'Per-IP, per-user, per-route and global, with three counting strategies' },
                    { title: 'Auth Shield', href: '#auth-shield-detail', desc: 'Brute force and credential stuffing, locked on the pair rather than the address' },
                    { title: 'Security Headers', href: '#security-headers', desc: 'Six headers, and the CSP reports come back as findings' },
                    { title: 'Anomaly Detection', href: '#anomaly-detection', desc: 'Six behavioural checks: what this account does, not what this request contains' },
                    { title: 'IP Geolocation', href: '#ip-geolocation', desc: 'Country and city per request, which is what makes impossible travel possible' },
                    { title: 'Threat Dashboard', href: '#threat-dashboard', desc: 'Live at /sentinel/ui, backed by the store you point it at' },
                    { title: 'Performance Monitoring', href: '#performance-monitoring', desc: 'Per-route p50, p95 and p99, measured where the WAF already sits' },
                  ].map((item) => (
                    <a
                      key={item.title}
                      href={item.href}
                      className="rounded-lg border border-border/30 bg-card/30 px-4 py-3 block hover:border-accent/40 transition-colors"
                    >
                      <h3 className="text-sm font-semibold text-foreground mb-1">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                    </a>
                  ))}
                </div>
              </div>

              {/* Every feature, in detail */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Every feature, in detail
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-8">
                  What each part actually covers, and where it stops. The lists below are the
                  ones Sentinel compiles in, not a summary of them: if a threat class or a
                  behavioural check is not named here, Sentinel does not look for it.
                </p>

                <h3 id="waf" className="text-xl font-semibold tracking-tight mb-3">
                  Web Application Firewall
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every request is matched against fourteen threat classes before your handler
                  runs. The matching covers the path, the query string, headers, cookies and the
                  body, and payloads are decoded first, so an attack hidden in URL encoding,
                  double encoding, base64 or mixed case is compared in the same form as a plain
                  one.
                </p>
                <div className="grid gap-2 sm:grid-cols-2 mb-4">
                  {[
                    ['SQLi', 'SQL injection, including UNION, boolean and time-based forms'],
                    ['XSS', 'Script injection in markup, attributes and javascript: URLs'],
                    ['PathTraversal', '../ and its encodings, reaching outside the intended directory'],
                    ['CommandInjection', 'Shell metacharacters that would run a second command'],
                    ['SSRF', 'URLs aimed back at your own network or cloud metadata'],
                    ['XXE', 'XML external entities, which read files or make requests for you'],
                    ['LFI', 'Local file inclusion: reading a file the handler never meant to serve'],
                    ['OpenRedirect', 'A redirect target an attacker chose, used to make a phish look like you'],
                    ['PrototypePollution', '__proto__ and constructor keys that rewrite object behaviour'],
                    ['Scanning', 'The shape of an automated probe rather than one bad request'],
                    ['RateLimitExceeded', 'Raised by the limiter below, recorded as a threat like any other'],
                    ['BruteForce', 'Raised by Auth Shield'],
                    ['AnomalyDetected', 'Raised by the behavioural checks'],
                    ['CSPViolation', 'Reported by the browser, back to your own endpoint'],
                  ].map(([name, what]) => (
                    <div key={name} className="rounded-lg border border-border/30 bg-card/20 px-3 py-2">
                      <code className="text-xs text-accent">{name}</code>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed mt-1">{what}</p>
                    </div>
                  ))}
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each rule has a sensitivity of <code>off</code>, <code>low</code>,{' '}
                  <code>medium</code> or <code>strict</code>, so a rule that costs you false
                  positives can be loosened or turned off without turning off the WAF. What
                  happens on a match is the <a href="#waf-modes" className="text-accent hover:underline">WAF mode</a>:
                  log it, block it, or challenge the caller. Your own patterns go in{' '}
                  <a href="#custom-rules" className="text-accent hover:underline">custom rules</a>, and{' '}
                  <a href="#waf-blind-spots" className="text-accent hover:underline">what the WAF does not inspect</a>{' '}
                  is worth reading before you rely on it.
                </p>

                <h3 id="rate-limits" className="text-xl font-semibold tracking-tight mb-3">
                  Rate Limiting
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Limits apply per IP, per authenticated user, per route, and globally, and the
                  four compose: a route limit does not replace the global one. Three counting
                  strategies are available, and the difference between them matters at the edges.
                </p>
                <ul className="text-muted-foreground leading-relaxed mb-4 space-y-2">
                  <li>
                    <code>fixed_window</code>: counts per clock window. Cheapest, and lets through
                    up to twice the limit across a window boundary.
                  </li>
                  <li>
                    <code>sliding_window</code>: counts the last N seconds continuously. No
                    boundary burst, which is why it is the default.
                  </li>
                  <li>
                    <code>token_bucket</code>: refills at a steady rate and allows a burst up to
                    the bucket size. The right one when normal traffic is bursty by nature.
                  </li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A limited response is a 429 carrying <code>Retry-After</code>, and every
                  response carries <code>X-RateLimit-Limit</code> and{' '}
                  <code>X-RateLimit-Remaining</code>, so a client can back off before it is
                  refused rather than after. Exceeding a limit is recorded as a{' '}
                  <code>RateLimitExceeded</code> threat against that actor, which is how a
                  scraper shows up on the dashboard as one caller rather than as noise.
                </p>

                <h3 id="auth-shield-detail" className="text-xl font-semibold tracking-tight mb-3">
                  Auth Shield
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sign-in gets its own protection, because the generic limiter cannot tell a
                  wrong password from a slow typist. Two things are watched: repeated failures
                  against one account, and one password tried across many accounts, which is
                  credential stuffing and looks nothing like brute force from the account&apos;s
                  side.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A lockout bars the <strong>(address, account) pair</strong> rather than the
                  address. Behind an office or a carrier NAT one address is thousands of people,
                  and barring it meant one person&apos;s typing locked out everybody sharing that
                  egress. Addresses that are known to be shared can be exempted from the
                  per-address bar entirely with <code>SENTINEL_LOCKOUT_EXEMPT_IPS</code>, and they
                  are still counted and still subject to the pair lockout.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Between &quot;no friction&quot; and &quot;you are barred&quot; there is a
                  challenge tier, which asks for a CAPTCHA instead of refusing. It needs a
                  provider configured (Turnstile, hCaptcha or reCAPTCHA); with none set the tier
                  exists and can never fire.
                </p>

                <h3 id="auto-block" className="text-xl font-semibold tracking-tight mb-3">
                  Auto-block
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Everything above detects and records. Auto-block is what stops the address
                  doing it again: ten threats of <strong>High</strong> severity or worse from one
                  address inside ten minutes bars that address for a day. On outside development,
                  off in development, where the thing it would block is you.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  It acts on behaviour rather than reputation. <code>IPReputation</code> asks a
                  third party what is already known about an address; this asks what the address
                  has just done to your application, which is the only thing that catches a
                  scanner nobody has reported yet. Counting runs through the same Redis counters
                  as the rate limits, so threats spread over replicas add up to one total and a
                  block made anywhere is enforced everywhere.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The severity floor is what keeps it from turning on your users. High is what the
                  scanner patterns score (path traversal, SQL injection, env-file probes), so a
                  lone medium-severity oddity does not accumulate towards a block, and each threat
                  counts once by its own id. Every automatic block is audited like a manual one,
                  and liftable from the dashboard.
                </p>
                <div className="rounded-lg border border-warning/30 bg-warning/5 p-4 mb-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Loopback and the private ranges are
                    never blocked</strong>, and the reason is worth knowing because it decides
                    whether this feature is safe on your deployment. Auto-block bars an address,
                    and the address is only the attacker&apos;s when Sentinel is reading the real
                    client. Behind a reverse proxy it has not been told to trust, every request
                    carries the proxy&apos;s address instead, so ten probes through that proxy
                    would bar the proxy: every user, for a day, from one scanner. Set{' '}
                    <code>SENTINEL_TRUSTED_PROXIES</code> to the proxy&apos;s address and the real
                    client is read, after which the exemption stops applying to anybody who
                    matters.
                  </p>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <code>SENTINEL_AUTOBLOCK_EXEMPT_IPS</code> adds to that list: your own
                  monitoring, a partner&apos;s integration, an office range. Left unset it follows{' '}
                  <code>SENTINEL_LOCKOUT_EXEMPT_IPS</code>, because both answer the same question
                  and a project that answered it once should not have to answer it twice.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Auto-block only sees threats from the moment it is on. For the addresses already
                  in your threat log, the dashboard&apos;s IP Management page has a{' '}
                  <strong>Block repeat offenders</strong> panel, which takes a minimum threat
                  count and a window and supports a dry run that blocks nobody.
                </p>

                <h3 id="auto-ban-sweep" className="text-xl font-semibold tracking-tight mb-3">
                  The scheduled ban sweep
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Auto-block answers a burst. It cannot answer patience: four
                  high-severity threats an hour, every hour, for a week never trips ten in ten
                  minutes, and that is what most real scanning looks like. One deployment&apos;s
                  threat log held 18,405 rows, nearly all of them env-file probes from a handful
                  of addresses, every one detected and recorded, and not one of them stopped.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  So there is a second rule that reads the window instead of the burst. Every six
                  hours it counts each address&apos;s high-severity threats over the last day and
                  bans the ones over the line. The two are complements: one answers in seconds,
                  the other catches everything slow enough to avoid it.
                </p>
                <CodeBlock language="bash" filename=".env" code={`# More than three high-severity threats in 24 hours, swept every 6 hours.
# MIN_THREATS is exceeded, not met: at 3, an address is banned on its fourth.
SECURITY_AUTO_BAN=true
SECURITY_AUTO_BAN_MIN_THREATS=3
SECURITY_AUTO_BAN_WINDOW_HOURS=24
AUTO_BAN_INTERVAL_HOURS=6
SECURITY_AUTO_BAN_PERMANENT=true`} />
                <p className="text-muted-foreground leading-relaxed mb-4 mt-4">
                  <strong>High severity only</strong>, which is the part worth understanding. An
                  actor&apos;s total threat count includes every low-severity heuristic it ever
                  tripped, so a client that trips one forty times would be banned by a number
                  that says nothing about intent. High is what the scanner patterns score: path
                  traversal, SQL injection, env-file probes.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Loopback and the private ranges are never banned, for the same reason auto-block
                  exempts them: behind a proxy Sentinel has not been told to trust, every request
                  carries the proxy&apos;s address. Addresses already banned are skipped, one
                  sweep is capped, and each ban is written with a reason that says what it
                  counted, which is also how the dashboard tells an automatic ban from a manual
                  one.
                </p>

                <h3 id="security-dashboard" className="text-xl font-semibold tracking-tight mb-3">
                  The Security page
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The admin panel reads all of this through the API rather than embedding
                  Sentinel&apos;s own console, because the two answer different questions. The
                  console is for investigating one address. The Security page is for the question
                  an operator actually arrives with, which is whether anything is happening and
                  whether they need to do something.
                </p>
                <ul className="list-disc pl-6 space-y-2 text-muted-foreground mb-4">
                  <li>The graded score, with what the grades mean. It measures posture, not
                      compromise: heavy scanning against a WAF in log-only mode pulls it down
                      without anything having been breached.</li>
                  <li>Threats, unique addresses, critical count and blocked count over 24 hours,
                      and the hourly shape of them.</li>
                  <li>What is being tried, which routes, and by whom, over seven days.</li>
                  <li>The threat log itself: filtered by severity, searched by address or path,
                      and paged. A log of 18,000 rows cannot be acted on ten at a time.</li>
                  <li>Who is banned, whether automatically or by hand, and the way back out.
                      The unban button matters more than the ban one.</li>
                  <li><strong>Ban high-threat IPs</strong>, which runs the sweep now on the same
                      policy the schedule uses. The schedule only sees what arrives after it is
                      turned on, and the log it is turned on in front of is usually already full.</li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every panel is sourced. When a call to Sentinel fails, the page says which one
                  rather than rendering a zero, because a zero on a security dashboard reads as
                  &quot;you are safe&quot;, and that is the one wrong answer it must never give.
                </p>

                <h3 id="security-headers" className="text-xl font-semibold tracking-tight mb-3">
                  Security Headers
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Six response headers, set on every response:
                </p>
                <ul className="text-muted-foreground leading-relaxed mb-4 space-y-2">
                  <li>
                    <code>Content-Security-Policy</code>: which origins may supply scripts,
                    styles, images and frames. There is a report-only variant, so you can watch
                    what a policy would have blocked before it blocks anything.
                  </li>
                  <li>
                    <code>Strict-Transport-Security</code>: this host is HTTPS only, for the next
                    N seconds, remembered by the browser.
                  </li>
                  <li>
                    <code>X-Frame-Options: DENY</code>: your pages cannot be put in someone
                    else&apos;s iframe, which is what clickjacking needs.
                  </li>
                  <li>
                    <code>X-Content-Type-Options: nosniff</code>: the browser believes your
                    Content-Type instead of guessing, so an uploaded file cannot be guessed into
                    a script.
                  </li>
                  <li>
                    <code>Referrer-Policy</code>: how much of the current URL is sent when a user
                    follows a link away. The default leaks the path to third parties.
                  </li>
                  <li>
                    <code>Permissions-Policy</code>: which browser features (camera, microphone,
                    geolocation) the page and its frames may ask for.
                  </li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  CSP is the one with a feedback loop: browsers post violations back, and
                  Sentinel records them as <code>CSPViolation</code> findings. A sudden run of
                  them usually means a policy that is too tight rather than an attack, and either
                  way you learn it from the dashboard rather than from a user.
                </p>

                <h3 id="anomaly-detection" className="text-xl font-semibold tracking-tight mb-3">
                  Anomaly Detection
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The WAF asks whether a request looks like an attack. These six ask whether this
                  account is behaving like itself, which catches the case where every individual
                  request is perfectly valid, because the credentials are real and stolen.
                </p>
                <div className="grid gap-2 sm:grid-cols-2 mb-4">
                  {[
                    ['impossible_travel', 'Signed in from two places further apart than anyone could travel between in the time elapsed'],
                    ['unusual_access', 'Reaching records or endpoints this account has never touched'],
                    ['data_exfiltration', 'Reading far more rows than this account normally reads'],
                    ['off_hours_access', 'Activity outside this account\u2019s established working pattern'],
                    ['velocity_anomaly', 'Acting faster than a person does, without tripping a rate limit'],
                    ['credential_stuffing', 'One password tried across many accounts, seen across the whole application'],
                  ].map(([name, what]) => (
                    <div key={name} className="rounded-lg border border-border/30 bg-card/20 px-3 py-2">
                      <code className="text-xs text-accent">{name}</code>
                      <p className="text-xs text-muted-foreground/70 leading-relaxed mt-1">{what}</p>
                    </div>
                  ))}
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sensitivity is <code>low</code>, <code>medium</code> or <code>high</code>.
                  These checks compare against a baseline, so they are quiet on a new deployment
                  and get sharper as there is more history to be unlike. A match is recorded as an{' '}
                  <code>AnomalyDetected</code> threat with the check that fired.
                </p>

                <h3 id="ip-geolocation" className="text-xl font-semibold tracking-tight mb-3">
                  IP Geolocation
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every request is resolved to a country and city, from GeoLite2 (free) or GeoIP2
                  (paid). This is less a feature than a prerequisite: impossible travel cannot be
                  computed without it, and a threat actor on the dashboard reads very differently
                  with a country beside it.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The address it resolves is the one Sentinel decided to trust. Behind a proxy
                  that means <code>SENTINEL_TRUSTED_PROXIES</code>: with it unset, forwarded
                  headers are ignored entirely and every request geolocates to your load
                  balancer.
                </p>

                <h3 id="threat-dashboard" className="text-xl font-semibold tracking-tight mb-3">
                  Threat Dashboard
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A dashboard at <code>/sentinel/ui</code>, updating over a WebSocket as requests
                  arrive. It is organised by actor rather than by event: an address or account
                  accumulates findings and carries a status of Active, Blocked or Whitelisted, so
                  the question &quot;what is this one caller doing&quot; has an answer, and
                  whitelisting the pen-tester is one action rather than a config change.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Findings are Low, Medium, High or Critical. They are kept in the store you
                  point Sentinel at: <code>sqlite</code>, <code>postgres</code>, <code>mysql</code>{' '}
                  or <code>memory</code>. Memory is genuinely in-memory, so a restart loses the
                  history, which is a reasonable choice on a laptop and the wrong one in
                  production.
                </p>

                <h3 id="performance-monitoring" className="text-xl font-semibold tracking-tight mb-3">
                  Performance Monitoring
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Per-route latency as p50, p95 and p99. It is here because the WAF is already in
                  the request path and already knows which route this was, so the measurement is
                  free; it is not a replacement for{' '}
                  <Link href="/docs/backend/pulse" className="text-accent hover:underline">Pulse</Link>,
                  which traces a request through its queries and tells you which line issued
                  them. Use these numbers to notice that something got slower, and Pulse to find
                  out why.
                </p>
              </div>

              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sentinel is enabled by default in every Grit project. Configure it through your{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file:
                </p>

                <CodeBlock language="bash" filename=".env" code={`# Security — Sentinel WAF, rate limiting, threat detection
SENTINEL_ENABLED=true
SENTINEL_USERNAME=admin
SENTINEL_PASSWORD=sentinel
SENTINEL_SECRET_KEY=change-me-in-production
SENTINEL_AUDIT_KEY=change-me-in-production`} />

                <div className="mb-6">
                  <table className="w-full border-collapse text-sm">
                    <thead>
                      <tr className="border-b border-border">
                        <th className="text-left font-semibold text-foreground pb-2 pr-4">Variable</th>
                        <th className="text-left font-semibold text-foreground pb-2 pr-4">Default</th>
                        <th className="text-left font-semibold text-foreground pb-2">Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      {[
                        { name: 'SENTINEL_ENABLED', def: 'true', desc: 'Set to "false" to disable Sentinel entirely' },
                        { name: 'SENTINEL_USERNAME', def: 'admin', desc: 'Dashboard login username' },
                        { name: 'SENTINEL_PASSWORD', def: 'sentinel', desc: 'Dashboard login password' },
                        { name: 'SENTINEL_SECRET_KEY', def: 'sentinel-secret...', desc: 'Secret for dashboard JWT sessions' },
                        { name: 'SENTINEL_AUDIT_KEY', def: '(generated)', desc: 'Keys the audit log hash chain, so an entry edited with database access alone fails verification' },
                      ].map((row) => (
                        <tr key={row.name} className="border-b border-border/50">
                          <td className="py-2 pr-4 font-mono text-xs text-accent/60">{row.name}</td>
                          <td className="py-2 pr-4 font-mono text-xs text-muted-foreground">{row.def}</td>
                          <td className="py-2 text-muted-foreground">{row.desc}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Replicas */}
              <div className="mb-12" id="replicas">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Replicas, the audit key and the caller
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A project scaffolded from v3.226.0 counts rate limits and AuthShield lockouts in
                  Redis when Redis is configured, so every replica counts against the same numbers:
                  counted per process, N replicas gave a client N times every limit. It keys the audit
                  log with{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">SENTINEL_AUDIT_KEY</code>,
                  and it tells Sentinel who made each request, which anomaly detection needs to see
                  anything at all.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  An older project gets the new Sentinel from{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit upgrade</code>,
                  but its <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code>{' '}
                  is its own, so these are yours to add to the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">sentinel.MountE</code> call
                  (and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">SentinelAuditKey</code>{' '}
                  to the config, read from{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">SENTINEL_AUDIT_KEY</code>):
                </p>
                <CodeBlock language="go" filename="apps/api/internal/routes/routes.go" code={`import "github.com/MUKE-coder/sentinel/v2/redisstore"

sentinelStorage.AuditKey = cfg.SentinelAuditKey

// Shared by every replica. Without Redis, each process counts alone.
var sentinelCounters sentinel.CounterStore
if svc.Cache != nil {
    sentinelCounters = redisstore.New(svc.Cache.Client())
}

err := sentinel.MountE(r, db, sentinel.Config{
    Storage:  sentinelStorage,
    Counters: sentinelCounters,
    // Read after the handler chain, when the auth middleware has set the caller.
    UserExtractor: func(c *gin.Context) *sentinel.UserContext {
        id := c.GetString("user_id")
        if id == "" {
            return nil
        }
        return &sentinel.UserContext{ID: id, Email: c.GetString("user_email"), Role: c.GetString("user_role")}
    },
    // ...the rest of the config as it was
})`} />
              </div>

              {/* How It Works */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  How It Works
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sentinel is mounted in your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code> Setup
                  function, right after the global middleware (Logger, Recovery, CORS) and before any route registrations.
                  This means every request passes through Sentinel&apos;s security middleware:
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`// Mount Sentinel security suite. MountE returns an error so a bad
// config fails fast at boot instead of silently running unprotected.
if cfg.SentinelEnabled {
    if err := sentinel.MountE(r, db, sentinel.Config{
        Dashboard: sentinel.DashboardConfig{
            Username:  cfg.SentinelUsername,
            Password:  cfg.SentinelPassword,
            SecretKey: cfg.SentinelSecretKey,
            // Only dev may run on default creds; prod (gin.ReleaseMode) refuses them.
            AllowInsecureDefaults: isDev,
        },
        WAF: sentinel.WAFConfig{
            Enabled: true,
            // Log-only in dev, block in production — set automatically.
            Mode: func() sentinel.WAFMode {
                if isDev { return sentinel.ModeLog }
                return sentinel.ModeBlock
            }(),
            MaxBodyBytes:        1 * 1024 * 1024, // 1 MB (richtext payloads)
            RejectOversizedBody: true,
            TrustedProxies:      cfg.SentinelTrustedProxies,
        },
        RateLimit: sentinel.RateLimitConfig{
            Enabled: true,
            ByIP:    &sentinel.Limit{Requests: 100, Window: 1 * time.Minute},
            ByRoute: map[string]sentinel.Limit{
                "/api/auth/login":    {Requests: 5, Window: 15 * time.Minute},
                "/api/auth/register": {Requests: 3, Window: 15 * time.Minute},
            },
        },
        AuthShield: sentinel.AuthShieldConfig{
            Enabled:    true,
            LoginRoute: "/api/auth/login",
        },
        AutoBlock: sentinel.AutoBlockConfig{
            Enabled: true,
            // Sentinel's defaults: 10 threats of High or worse in 10 minutes
            // bars the address for 24 hours.
            //
            // Loopback and the private ranges, because behind a proxy that
            // Sentinel has not been told to trust, every request carries the
            // proxy's address and blocking it blocks everybody.
            ExemptIPs: autoBlockExempt(cfg),
        },
        Anomaly: sentinel.AnomalyConfig{
            Enabled: true,
        },
        Geo: sentinel.GeoConfig{
            Enabled: true,
        },
    }); err != nil {
        log.Fatalf("mounting Sentinel: %v", err)
    }
}`} />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The middleware execution order is: <strong>Auth Shield</strong> (intercepts login attempts) &rarr;{' '}
                  <strong>WAF</strong> (inspects all input vectors) &rarr; <strong>Rate Limiter</strong> (enforces
                  request limits) &rarr; <strong>Security Headers</strong> (injects response headers) &rarr;{' '}
                  <strong>Performance Monitor</strong> (tracks latency).
                </p>
              </div>

              {/* What the WAF does not inspect */}
              <div className="mb-12">
                <h2 id="waf-blind-spots" className="text-2xl font-semibold tracking-tight mb-4">
                  What the WAF does not inspect
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A few routes have to be exempt from body inspection, and{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">wafExcludedRoutes()</code> in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">internal/routes/routes.go</code>{' '}
                  is the list. A photograph is larger than the inspection cap by design. And a
                  richtext body is <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&lt;p&gt;</code>,{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&lt;strong&gt;</code> and{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&lt;img&gt;</code> by
                  definition, so the XSS heuristics flag ordinary markup, and a post whose body
                  quotes a shell command trips the traversal and command-injection rules on the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">../</code> and the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">$(</code> inside a code block.
                </p>
                <CodeBlock language="go" filename="internal/routes/routes.go" code={`paths := []string{
    "/blogs", "/blogs/*",
    "/admin/blogs", "/admin/blogs/*",
    "/uploads", "/uploads/*",
    "/*/import",          // a spreadsheet is over the body cap
    "/public/forms/*",
    // grit generate resource adds a resource here when it emits a richtext field.
    // grit:waf:richtext
}`} />
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every richtext resource is listed twice, at its public path and again under{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/admin</code>. A body is
                  inspected on the way in, so it is the writes that carry the markup, and the admin
                  panel writes through the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/admin</code>{' '}
                  prefix. With only the public path listed, a post could be published once and then
                  never edited: every <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">PUT</code>{' '}
                  from the admin was answered 403 and logged as a critical threat. That was{' '}
                  <a href="https://github.com/MUKE-coder/grit/issues/90" className="underline decoration-dotted underline-offset-4 hover:text-foreground">issue&nbsp;90</a>,
                  fixed in v3.335.0; <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit upgrade</code>{' '}
                  adds the missing entries to a project generated before it.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  <strong className="text-foreground">Exclusion is from body inspection only.</strong>{' '}
                  These routes still pass through authentication, RBAC, binding validation,{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">sanitize:&quot;html&quot;</code>{' '}
                  and the rate limits. The HTML is cleaned on every write whether the WAF looked at
                  it or not, which is what actually protects the reader: a WAF pattern match is a
                  guess about intent, and an allowlist sanitizer is not.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Only <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">richtext</code> fields
                  earn an entry, never <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">text</code>.
                  Nearly every resource has a description or a note, and exempting all of them would
                  leave the WAF inspecting nothing worth inspecting. If a plain text column does
                  carry code, add that resource to the list by hand above the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit:waf:richtext</code> marker.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In development the WAF runs in <strong className="text-foreground">log mode</strong>,
                  so a body it dislikes is recorded and served anyway. It blocks outside development.
                  That is worth knowing when a payload works on your laptop and returns 403 in
                  staging: set <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">APP_ENV</code> to
                  anything but <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">development</code>{' '}
                  to reproduce it locally, and read the evidence in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/sentinel/ui</code>, which
                  names the pattern and the substring that matched.
                </p>
              </div>

              {/* Dashboard */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Security Dashboard
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sentinel ships with an embedded React dashboard at{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">/sentinel/ui</code>. Log in
                  with the username and password from your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file
                  (default: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">admin</code> / <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">sentinel</code>).
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The dashboard provides:
                </p>
                <ul className="text-muted-foreground leading-relaxed mb-4 space-y-1">
                  <li><strong>Threat feed</strong> -- real-time threat events via WebSocket with severity, type, and source IP</li>
                  <li><strong>Threat actors</strong> -- automatic profiling with risk scores, attack types, and geolocation</li>
                  <li><strong>Security score</strong> -- weighted score across 5 dimensions with actionable recommendations</li>
                  <li><strong>Analytics</strong> -- attack trends, geographic distribution, top targeted routes</li>
                  <li><strong>WAF rules</strong> -- view and configure rules, test payloads, add custom patterns</li>
                  <li><strong>Rate limits</strong> -- view active limits and live counter states</li>
                  <li><strong>IP management</strong> -- block/unblock IPs, view blocked list</li>
                  <li><strong>Audit logs</strong> -- GORM-level CREATE/UPDATE/DELETE audit trail with before/after diffs</li>
                  <li><strong>Performance</strong> -- per-route latency metrics (p50/p95/p99) and error rates</li>
                  <li><strong>Compliance reports</strong> -- generate GDPR, PCI-DSS, and SOC2 reports</li>
                </ul>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The dashboard is also embedded in the admin panel under{' '}
                  <strong>System &rarr; Security</strong> for convenience.
                </p>
              </div>

              {/* WAF Modes */}
              <div className="mb-12">
                <h2 id="waf-modes" className="text-2xl font-semibold tracking-tight mb-4">
                  WAF Modes
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Web Application Firewall has two modes:
                </p>
                <div className="grid gap-3 sm:grid-cols-2 mb-4">
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <code className="text-sm font-mono text-accent/70 font-medium">ModeLog</code>
                    <span className="ml-2 text-xs text-warning bg-warning/10 px-1.5 py-0.5 rounded-full">dev default</span>
                    <p className="text-xs text-muted-foreground/60 mt-1 leading-relaxed">
                      Detects and logs threats but does not block requests. Use during development and initial
                      deployment to identify false positives before switching to block mode.
                    </p>
                  </div>
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <code className="text-sm font-mono text-accent/70 font-medium">ModeBlock</code>
                    <span className="ml-2 text-xs text-success bg-success/10 px-1.5 py-0.5 rounded-full">production</span>
                    <p className="text-xs text-muted-foreground/60 mt-1 leading-relaxed">
                      Detects threats and blocks malicious requests with a 403 response. Switch to this
                      mode in production after verifying no false positives in log mode.
                    </p>
                  </div>
                </div>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit picks the mode for you: <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">ModeLog</code> in
                  development and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">ModeBlock</code> in
                  production (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">gin.ReleaseMode</code>), so
                  you log-and-review locally and block for real once deployed. Override the{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">Mode</code> in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code> if you want a fixed
                  mode, or toggle it at runtime from the dashboard.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The WAF detects the following attack types with configurable sensitivity per rule (off, low, medium, strict):
                </p>
                <div className="grid gap-2 sm:grid-cols-2 mb-4">
                  {[
                    'SQL Injection (basic, blind, stacked)',
                    'Cross-Site Scripting (XSS)',
                    'Path Traversal',
                    'Command Injection',
                    'Server-Side Request Forgery (SSRF)',
                    'XML External Entity (XXE)',
                    'Local File Inclusion (LFI)',
                    'Open Redirect',
                  ].map((attack) => (
                    <div key={attack} className="text-sm text-muted-foreground flex items-center gap-2">
                      <span className="text-accent/50">&#10003;</span>
                      {attack}
                    </div>
                  ))}
                </div>
              </div>

              {/* Rate Limiting */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Rate Limiting
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit configures sensible rate limits out of the box. You can customize them in{' '}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">routes.go</code>:
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`RateLimit: sentinel.RateLimitConfig{
    Enabled: true,
    ByIP:    &sentinel.Limit{Requests: 100, Window: 1 * time.Minute},
    ByRoute: map[string]sentinel.Limit{
        "/api/auth/login":    {Requests: 5, Window: 15 * time.Minute},
        "/api/auth/register": {Requests: 3, Window: 15 * time.Minute},
        "/api/uploads":       {Requests: 10, Window: 1 * time.Minute},
    },
    // Optional: per-user and global limits
    // ByUser: &sentinel.Limit{Requests: 200, Window: 1 * time.Minute},
    // Global: &sentinel.Limit{Requests: 1000, Window: 1 * time.Minute},
    Strategy: sentinel.SlidingWindow, // or FixedWindow, TokenBucket
},`} />
              </div>

              {/* Auth Shield */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Auth Shield
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The auth shield protects your login endpoint from brute-force attacks. When enabled, it
                  automatically tracks failed login attempts per IP address and locks out attackers after
                  too many failures.
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`AuthShield: sentinel.AuthShieldConfig{
    Enabled:                     true,
    LoginRoute:                  "/api/auth/login",
    MaxFailedAttempts:           5,               // Lock after 5 failures (default)
    LockoutDuration:             15 * time.Minute, // 15-minute lockout (default)
    CredentialStuffingDetection: true,             // Detect credential lists
    BruteForceDetection:         true,             // Detect brute-force patterns
},`} />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Locked-out users can be unblocked from the dashboard or via the API.
                </p>
              </div>

              {/* Custom WAF Rules */}
              <div className="mb-12">
                <h2 id="custom-rules" className="text-2xl font-semibold tracking-tight mb-4">
                  Custom WAF Rules
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  You can add custom WAF rules to detect application-specific attack patterns:
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`WAF: sentinel.WAFConfig{
    Enabled: true,
    Mode:    sentinel.ModeBlock,
    CustomRules: []sentinel.WAFRule{
        {
            ID:        "block-admin-enum",
            Name:      "Block admin enumeration",
            Pattern:   ` + '`' + `(?i)/(wp-admin|phpmyadmin|administrator)` + '`' + `,
            AppliesTo: []string{"path"},
            Severity:  sentinel.SeverityMedium,
            Action:    "block",
            Enabled:   true,
        },
    },
    ExcludeRoutes: []string{"/api/health"}, // Skip WAF for health checks
},`} />

                <p className="text-muted-foreground leading-relaxed mb-4">
                  Custom rules can also be managed at runtime from the dashboard without redeploying.
                </p>
              </div>

              {/* AI-Powered Analysis */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  AI-Powered Analysis
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If you have an AI API key configured (Claude, OpenAI, or Gemini), Sentinel can use it for
                  intelligent threat analysis:
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`AI: &sentinel.AIConfig{
    Provider: sentinel.Claude,       // or sentinel.OpenAI, sentinel.Gemini
    APIKey:   cfg.AIAPIKey,          // Reuse your existing AI API key
    DailySummary: true,              // Generate daily threat summaries
},`} />

                <ul className="text-muted-foreground leading-relaxed mb-4 space-y-1">
                  <li><strong>Threat analysis</strong> -- AI evaluates individual threats and provides context</li>
                  <li><strong>Actor assessment</strong> -- AI profiles threat actors and predicts intent</li>
                  <li><strong>Daily summaries</strong> -- automated security briefings</li>
                  <li><strong>Natural language queries</strong> -- ask questions about your security data</li>
                  <li><strong>WAF recommendations</strong> -- AI suggests new rules based on attack patterns</li>
                </ul>
              </div>

              {/* Alerts */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Alerts
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Sentinel can send real-time alerts when threats are detected. Configure one or more channels:
                </p>

                <CodeBlock language="go" filename="internal/routes/routes.go" code={`Alerts: sentinel.AlertConfig{
    MinSeverity: sentinel.SeverityHigh,
    Slack: &sentinel.SlackConfig{
        WebhookURL: "https://hooks.slack.com/services/...",
    },
    Email: &sentinel.EmailConfig{
        SMTPHost:   "smtp.example.com",
        SMTPPort:   587,
        Username:   "alerts@example.com",
        Password:   "smtp-password",
        Recipients: []string{"security@example.com"},
    },
    Webhook: &sentinel.WebhookConfig{
        URL:     "https://your-siem.example.com/webhook",
        Headers: map[string]string{"X-Token": "your-token"},
    },
},`} />
              </div>

              {/* Production Recommendations */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Production Checklist
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Before deploying to production, make sure to:
                </p>
                <div className="space-y-2">
                  {[
                    'Change SENTINEL_PASSWORD and SENTINEL_SECRET_KEY to strong, unique values',
                    'Switch WAF mode from ModeLog to ModeBlock after reviewing logs for false positives',
                    'Configure at least one alert channel (Slack, email, or webhook)',
                    'Review rate limits and adjust for your expected traffic',
                    'Set up IP reputation checking with an AbuseIPDB API key for auto-blocking known bad actors',
                    'Consider enabling AI-powered analysis for automated threat intelligence',
                  ].map((item, i) => (
                    <div key={i} className="flex items-start gap-3 text-sm text-muted-foreground">
                      <span className="flex h-5 w-5 items-center justify-center rounded border border-border text-xs text-muted-foreground/50 shrink-0 mt-0.5">
                        {i + 1}
                      </span>
                      <span>{item}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Disabling Sentinel */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Disabling Sentinel
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  If you don&apos;t need Sentinel, set <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">SENTINEL_ENABLED=false</code> in
                  your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file. The app will
                  run normally without any security middleware or dashboard -- no code changes needed.
                </p>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/batteries/ai" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  AI Integration
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/infrastructure/docker" className="gap-1.5">
                  Docker Setup
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
    </DocsShell>
  )
}
