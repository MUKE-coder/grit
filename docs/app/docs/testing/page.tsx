import Link from 'next/link'
import { ArrowLeft, ArrowRight, Activity, Bug, FileCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/testing')

export default function TestingPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Testing</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Performance &amp; Security Testing
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                How to run load tests with k6 and a methodology-driven penetration
                test against a Grit app. Both produce evidence — a before/after latency
                graph, a list of closed vulnerabilities — that&apos;s exactly what serious
                clients pay a premium for.
              </p>
              <LaneFlow
                id="testing-svc"
                lanes={['Test a Grit app', 'Evidence clients pay for']}
                nodes={[
                  { id: 'load', lane: 0, row: 0, title: 'Load test', sub: 'latency under load', tone: 'cyan' },
                  { id: 'sec', lane: 0, row: 1, title: 'Security scan', sub: 'vulnerabilities', tone: 'rose' },
                  { id: 'graph', lane: 1, row: 0, title: 'Latency graph', sub: 'before / after', tone: 'blue' },
                  { id: 'vulns', lane: 1, row: 1, title: 'Closed vulns', sub: 'report', tone: 'green' },
                ]}
                edges={[
                  { from: 'load', to: 'graph', label: 'produces', tone: 'blue' },
                  { from: 'sec', to: 'vulns', label: 'produces', tone: 'green' },
                ]}
                legend={[{ tone: 'cyan', label: 'Performance' }, { tone: 'rose', label: 'Security' }]}
                caption="Load and security tests both produce hard evidence — the deliverables clients pay a premium for"
              />
            </div>

            <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-8">
              <p className="text-sm text-muted-foreground leading-relaxed">
                <strong className="text-foreground">Go deeper:</strong> the{' '}
                <Link href="/courses/testing" className="text-primary hover:underline">Testing Your Grit App</Link>{' '}
                course walks through Go, Vitest and Playwright suites end to end.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── 1. PERFORMANCE TESTING ─────────────────────────────────── */}
              <h2 id="performance-testing" className="flex items-center gap-2">
                <Activity className="h-5 w-5 text-primary" />
                Performance testing with k6
              </h2>
              <p>
                Every fresh Grit project ships a complete k6 suite in
                <code> tests/k6/</code> covering the six load-test types: smoke,
                average-load, stress, spike, soak, breakpoint. They share a single user
                journey via <code>tests/k6/lib/common.js</code> — edit the journey once,
                reshape the load profile per test.
              </p>

              <h3 id="install-k6" className="mt-6">Install k6</h3>
              <CodeBlock language="bash" code={`# macOS
brew install k6

# Linux (Debian/Ubuntu)
sudo gpg -k && sudo gpg --no-default-keyring \\
  --keyring /usr/share/keyrings/k6-archive-keyring.gpg \\
  --keyserver hkp://keyserver.ubuntu.com:80 \\
  --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" \\
  | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update && sudo apt-get install k6

# Windows
winget install k6`} />

              <h3 id="six-test-types" className="mt-6">The six test types — when to run each</h3>
              <div className="overflow-x-auto my-4">
                <table className="w-full text-sm border-collapse">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left px-4 py-2 font-medium">Type</th>
                      <th className="text-left px-4 py-2 font-medium">Question it answers</th>
                      <th className="text-left px-4 py-2 font-medium">When</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">smoke.js</td>
                      <td className="px-4 py-2.5">Script + system handle minimal load?</td>
                      <td className="px-4 py-2.5">Every PR (CI gate)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">average-load.js</td>
                      <td className="px-4 py-2.5">Normal expected traffic behaviour?</td>
                      <td className="px-4 py-2.5">Every PR (CI gate)</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">stress.js</td>
                      <td className="px-4 py-2.5">Failure mode at the limit?</td>
                      <td className="px-4 py-2.5">Before launch</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">spike.js</td>
                      <td className="px-4 py-2.5">Survive a sudden surge + recover?</td>
                      <td className="px-4 py-2.5">Before launch</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">soak.js</td>
                      <td className="px-4 py-2.5">Memory leaks / resource creep over hours?</td>
                      <td className="px-4 py-2.5">Before major releases</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">breakpoint.js</td>
                      <td className="px-4 py-2.5">Exact capacity — VU count at failure?</td>
                      <td className="px-4 py-2.5">Capacity planning</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <h3 id="run-tests" className="mt-6">Run a test</h3>
              <CodeBlock language="bash" code={`# 1. Start the API (in another terminal)
grit start server

# 2. Run any test
export BASE_URL=http://localhost:8080
k6 run tests/k6/smoke.js              # 30s — does it work?
k6 run tests/k6/average-load.js       # 9m — baseline at 100 VUs
k6 run tests/k6/stress.js             # 9m — 4× load
k6 run tests/k6/spike.js              # 4m — 50 → 1000 → 50
k6 run tests/k6/breakpoint.js         # 1h — slow ramp to 5000
k6 run tests/k6/soak.js               # 4h — overnight`} />

              <h3 id="ci-gate" className="mt-6">Wire smoke + average-load into CI</h3>
              <p>
                k6 exits non-zero when a threshold breaches — that&apos;s how it gates the
                pipeline. Add a job to your workflow:
              </p>
              <CodeBlock filename=".github/workflows/perf.yml" code={`name: perf
on: [pull_request]
jobs:
  k6:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build and start API
        run: |
          go build -o /tmp/server .
          /tmp/server &
          sleep 5
      - uses: grafana/setup-k6-action@v1
      - name: Smoke + average-load
        env:
          BASE_URL: http://localhost:8080
        run: |
          k6 run tests/k6/smoke.js
          k6 run tests/k6/average-load.js`} />

              <h3 id="reading-results" className="mt-6">Reading the result</h3>
              <p>
                Open <code>/pulse/ui/</code> in another tab during the run — Pulse shows
                you live request traces and DB query timings so you can see the bottleneck
                appear in real time. Look for:
              </p>
              <ul className="list-disc pl-6 my-3">
                <li><strong>Smoke</strong>: must always pass. Failure = the script is broken, not the system.</li>
                <li><strong>Average-load</strong>: p95 ≤ 500 ms and error-rate ≤ 1% is the default SLO. Loosen / tighten in <code>lib/common.js</code>.</li>
                <li><strong>Stress</strong>: graceful degradation (latency climbs, errors stay low) vs collapse (error spike). The former is fine; the latter needs work.</li>
                <li><strong>Spike</strong>: the recovery curve is the answer. p95 must return to baseline within ~1 min after the surge.</li>
                <li><strong>Soak</strong>: a slowly tilting latency line over 4h = memory leak or unclosed connections. Check Pulse&apos;s runtime metrics for steady memory growth.</li>
                <li><strong>Breakpoint</strong>: the VU count at which thresholds first breach is your real capacity. Plan launches at 50% of this.</li>
              </ul>

              {/* ── 2. SECURITY / PENTEST ──────────────────────────────────── */}
              <h2 id="security-testing" className="flex items-center gap-2 mt-12">
                <Bug className="h-5 w-5 text-primary" />
                Security testing — the pentest methodology
              </h2>

              <div className="p-4 rounded-lg border border-danger/30 bg-danger/5 my-4">
                <p className="text-sm text-foreground/90 mb-0">
                  <strong className="text-danger">Only test what you&apos;re authorised to test.</strong>
                  Techniques are identical for attack and defence — the line between
                  &quot;security professional&quot; and &quot;criminal&quot; is the signed scope. Practise
                  on your own apps, OWASP Juice Shop, PortSwigger labs, HackTheBox.
                  Anything else is a crime in most countries, regardless of intent.
                </p>
              </div>

              <h3 id="five-phases" className="mt-6">The five-phase methodology</h3>
              <p>
                Follow these in order, every time. The methodology (not the tools) is
                what makes a pentest thorough and repeatable.
              </p>

              <ol className="list-decimal pl-6 my-3 space-y-1">
                <li><strong>Scope &amp; authorisation</strong> — signed scope + rules of engagement before anything technical.</li>
                <li><strong>Recon &amp; mapping</strong> — passive (OSINT, Shodan) + active (Nmap, Burp + ffuf).</li>
                <li><strong>Vulnerability discovery</strong> — automated scan + manual testing of logic &amp; access control.</li>
                <li><strong>Exploitation</strong> — prove impact with the minimum necessary; chain low-sev findings into bigger ones.</li>
                <li><strong>Reporting &amp; remediation</strong> — CVSS-scored report with reproduction steps + fixes.</li>
              </ol>

              <h3 id="toolkit" className="mt-6">The toolkit</h3>
              <div className="overflow-x-auto my-4">
                <table className="w-full text-sm border-collapse">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left px-4 py-2 font-medium">Tool</th>
                      <th className="text-left px-4 py-2 font-medium">For</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Burp Suite (Community)</td>
                      <td className="px-4 py-2.5">Intercepting proxy — see / replay / modify every request. The center of any web pentest.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Nmap</td>
                      <td className="px-4 py-2.5">Port + service scanning — what&apos;s exposed.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">ffuf / Gobuster</td>
                      <td className="px-4 py-2.5">Content discovery — brute-force hidden directories &amp; endpoints.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">sqlmap</td>
                      <td className="px-4 py-2.5">Automate SQL-injection detection + exploitation safely.</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono text-xs">Nuclei</td>
                      <td className="px-4 py-2.5">Template-based vuln scanning against a huge community library.</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono text-xs">govulncheck + pnpm audit</td>
                      <td className="px-4 py-2.5">Supply chain — already wired into <code>.github/workflows/security.yml</code>.</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <h3 id="pentest-grit-app" className="mt-6">Running a pentest against a Grit app — what to test</h3>
              <p>
                Map each OWASP Top 10:2025 category to a concrete test against the
                generated API. Cross-reference defences on the
                {' '}<Link href="/docs/security" className="text-primary hover:underline">Security Guide</Link>{' '}
                page.
              </p>

              <h4 className="mt-4">Broken Access Control / IDOR (A01)</h4>
              <CodeBlock language="bash" code={`# Login as user A
TOKEN_A=$(curl -s -X POST http://localhost:8080/api/auth/login \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"a@example.com","password":"password"}' | jq -r '.data.tokens.access_token')

# Create an invoice
INV_ID=$(curl -s -X POST http://localhost:8080/api/invoices \\
  -H "Authorization: Bearer $TOKEN_A" \\
  -d '{"amount":100}' | jq -r '.data.id')

# Now login as user B and try to read user A's invoice
TOKEN_B=$(curl -s -X POST http://localhost:8080/api/auth/login \\
  -H 'Content-Type: application/json' \\
  -d '{"email":"b@example.com","password":"password"}' | jq -r '.data.tokens.access_token')

curl -s http://localhost:8080/api/invoices/$INV_ID \\
  -H "Authorization: Bearer $TOKEN_B"
# Expected: 404 (NOT 403, NOT 200). authz.MustOwn returns 404 so existence
# of A's invoice doesn't leak to B.`} />

              <h4 className="mt-4">SQL Injection (A05)</h4>
              <CodeBlock language="bash" code={`# Try classic payloads on any parameter that hits the DB
curl "http://localhost:8080/api/users?email=' OR '1'='1"
curl "http://localhost:8080/api/users?email=admin' --"
# Expected: 200 with a normal (empty) response. GORM parameterises so the
# input is interpreted as data, never SQL.

# Time-based blind probe
curl "http://localhost:8080/api/users?email=' OR pg_sleep(5)--"
# Expected: response in <100ms (no delay). The DB never sees the payload as SQL.`} />

              <h4 className="mt-4">XSS (A05)</h4>
              <CodeBlock language="bash" code={`# Try stored XSS via a field that's later rendered
curl -X POST http://localhost:8080/api/blogs \\
  -H "Authorization: Bearer $TOKEN_ADMIN" \\
  -d '{"title":"<script>alert(1)</script>","content":"x"}'

# View the blog in the SPA. React escapes by default → the script tag
# renders as text. The CSP header blocks inline script as a 2nd layer.
# Check the response headers:
curl -I http://localhost:8080/ | grep -i content-security
# Expected: Content-Security-Policy: default-src 'self'; script-src 'self'; ...`} />

              <h4 className="mt-4">SSRF (A01 — 2025)</h4>
              <CodeBlock language="bash" code={`# Try to make the server fetch the AWS metadata endpoint via any feature
# that fetches user-provided URLs (webhook delivery, image-from-URL, etc).
curl -X POST http://localhost:8080/api/webhooks/dispatch \\
  -H "Authorization: Bearer $TOKEN" \\
  -d '{"url":"http://169.254.169.254/latest/meta-data/iam/security-credentials/"}'

# Expected: 400 with "URL not allowed". internal/safefetch blocks the
# request at validation AND re-blocks at TCP-connect time if DNS rebinds.`} />

              <h4 className="mt-4">Authentication brute force (A07)</h4>
              <CodeBlock language="bash" code={`# Hammer the login endpoint
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/api/auth/login \\
    -H 'Content-Type: application/json' \\
    -d '{"email":"victim@example.com","password":"guess'$i'"}'
done
# Expected after ~5 attempts: 429 "Rate limited" (Sentinel's per-route limit).
# After more attempts at the same email: the account locks (AuthShield).`} />

              <h4 className="mt-4">Misconfiguration / verbose errors (A02)</h4>
              <CodeBlock language="bash" code={`# Probe for verbose error pages
curl http://localhost:8080/api/this-route-does-not-exist
curl http://localhost:8080/api/users/invalid-uuid

# Expected: generic error JSON, no stack trace, no DB driver names.
# Also check the security headers are set:
curl -I http://localhost:8080/api/health | grep -iE 'x-frame|x-content|content-security|strict-transport|referrer'`} />

              {/* ── 3. AUDIT REPORT ─────────────────────────────────────────── */}
              <h2 id="audit-report" className="flex items-center gap-2 mt-12">
                <FileCheck className="h-5 w-5 text-primary" />
                The audit report — what to deliver
              </h2>
              <p>
                A polished, CVSS-scored, evidence-backed report is what justifies the
                fee. Structure it for two readers — an executive who needs the bottom
                line on page 1, and an engineer who needs enough detail to reproduce
                each finding.
              </p>

              <ol className="list-decimal pl-6 my-3 space-y-1">
                <li><strong>Executive summary</strong> (1 page) — overall risk posture, finding counts by severity, top 3 business risks.</li>
                <li><strong>Scope &amp; methodology</strong> — what was tested, what wasn&apos;t, dates, approach (e.g. &quot;authorised black-box web pentest per OWASP WSTG&quot;).</li>
                <li><strong>Findings</strong> — one entry per vuln, sorted critical-first. Each finding needs: title, CVSS score + severity, affected component, plain-English risk description, reproduction steps, evidence (screenshots, request/response), and a specific fix.</li>
                <li><strong>Remediation roadmap</strong> — prioritised list with SLAs (Critical: 7 days, High: 14 days, Medium: 30 days, Low: 90 days).</li>
                <li><strong>Appendices</strong> — raw tool output, full logs.</li>
              </ol>

              <h3 id="cvss" className="mt-6">CVSS scoring (FIRST.org)</h3>
              <div className="overflow-x-auto my-4">
                <table className="w-full text-sm border-collapse">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left px-4 py-2 font-medium">Score</th>
                      <th className="text-left px-4 py-2 font-medium">Severity</th>
                      <th className="text-left px-4 py-2 font-medium">SLA</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono">9.0–10.0</td>
                      <td className="px-4 py-2.5 text-danger font-medium">Critical</td>
                      <td className="px-4 py-2.5">24h – 7 days</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono">7.0–8.9</td>
                      <td className="px-4 py-2.5 text-warning font-medium">High</td>
                      <td className="px-4 py-2.5">within 7 days</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5 font-mono">4.0–6.9</td>
                      <td className="px-4 py-2.5">Medium</td>
                      <td className="px-4 py-2.5">14–30 days</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5 font-mono">0.1–3.9</td>
                      <td className="px-4 py-2.5 text-muted-foreground">Low</td>
                      <td className="px-4 py-2.5">60–90 days</td>
                    </tr>
                  </tbody>
                </table>
              </div>
              <p>
                Adjust the raw CVSS by business context — a &quot;Critical&quot; on an air-gapped
                internal tool may be a real-world Low; a &quot;Medium&quot; on a public payment
                endpoint may be a real-world Critical. Document the adjustment and the
                reasoning. That documented judgment is the senior-level deliverable.
              </p>

              <h2 id="continuous-evidence" className="mt-10">Continuous evidence — between pentests</h2>
              <p>
                A pentest is a snapshot. Between tests, three things keep the system
                defensible and prove it:
              </p>
              <ul className="list-disc pl-6 my-3">
                <li><strong>Audit trails</strong> — Grit&apos;s <code>middleware.LogSecurityEvent</code> + the activity-log hash chain provide tamper-evident records of every authN/authZ event.</li>
                <li><strong>Continuous scanning</strong> — <code>.github/workflows/security.yml</code> runs govulncheck + pnpm audit + CodeQL on every PR and weekly. Dependabot raises PRs the moment a CVE drops.</li>
                <li><strong>Remediation tracking</strong> — every finding flows from discovery → ticket (severity + owner + SLA) → fix → re-test. That documented loop is what SOC 2 / ISO 27001 auditors ask to see.</li>
              </ul>

              <h2 id="resources" className="mt-10">Resources</h2>
              <ul className="list-disc pl-6 my-3">
                <li><a className="text-primary hover:underline" href="https://portswigger.net/web-security" target="_blank" rel="noreferrer">PortSwigger Web Security Academy</a> — the best free pentest labs anywhere.</li>
                <li><a className="text-primary hover:underline" href="https://owasp.org/www-project-juice-shop/" target="_blank" rel="noreferrer">OWASP Juice Shop</a> — the deliberately vulnerable practice app.</li>
                <li><a className="text-primary hover:underline" href="https://owasp.org/Top10/" target="_blank" rel="noreferrer">OWASP Top 10:2025</a> — the canonical risk map.</li>
                <li><a className="text-primary hover:underline" href="https://owasp.org/www-project-web-security-testing-guide/" target="_blank" rel="noreferrer">OWASP WSTG</a> — 90+ web-app test cases mapped to the Top 10.</li>
                <li><a className="text-primary hover:underline" href="https://grafana.com/docs/k6/latest/" target="_blank" rel="noreferrer">k6 docs</a> — load-testing reference.</li>
                <li><a className="text-primary hover:underline" href="https://first.org/cvss/" target="_blank" rel="noreferrer">FIRST CVSS</a> — official scoring spec + calculator.</li>
              </ul>
            </div>

            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/security" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Security Guide
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/pulse" className="gap-1.5">
                  Pulse (Observability)
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
