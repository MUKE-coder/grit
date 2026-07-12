import Link from 'next/link'
import { ArrowLeft, ArrowRight, ShieldCheck, AlertTriangle, ExternalLink } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/security')

interface Row {
  id: string
  title: string
  oneLine: string
  defence: React.ReactNode
}

const OWASP_2025: Row[] = [
  {
    id: 'a01',
    title: 'A01 — Broken Access Control',
    oneLine: 'Users acting beyond their permissions; now includes SSRF. Still #1.',
    defence: (
      <>
        <p>
          Grit ships the <code>internal/authz</code> package as the canonical IDOR defence.
          Every object access goes through <code>authz.MustOwn</code>, which loads the row
          by ID, verifies it belongs to the authenticated user, and returns 404 (never 403)
          so existence isn&apos;t leaked through error-message differences.
        </p>
        <CodeBlock filename="internal/handlers/invoice.go" code={`func (h *InvoiceHandler) GetByID(c *gin.Context) {
    var invoice models.Invoice
    if err := authz.MustOwn(c, h.DB, &invoice, c.Param("id")); err != nil {
        return // helper already wrote 404
    }
    c.JSON(http.StatusOK, gin.H{"data": invoice})
}

// The model implements Ownable:
func (i *Invoice) GetOwnerID() string { return i.UserID }`} />
        <p>
          For tenant / team scoping use <code>authz.CheckScope</code>. For admin-only
          routes use <code>authz.RequireRoles(&quot;ADMIN&quot;)</code>. SSRF (absorbed into A01 in
          2025) is handled by the <code>internal/safefetch</code> package — see A05 below.
        </p>
      </>
    ),
  },
  {
    id: 'a02',
    title: 'A02 — Security Misconfiguration',
    oneLine: 'Insecure defaults, open buckets, verbose errors. Up from #5.',
    defence: (
      <>
        <p>
          The <code>SecurityHeaders</code> middleware sets <em>all</em> the OWASP-recommended
          headers by default: HSTS, X-Content-Type-Options, X-Frame-Options,
          Referrer-Policy, Permissions-Policy, Cross-Origin-Opener-Policy,
          Cross-Origin-Resource-Policy, and a strict Content-Security-Policy that
          blocks inline script.
        </p>
        <CodeBlock filename="internal/middleware/security.go" code={`func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()")
        c.Header("Cross-Origin-Opener-Policy", "same-origin")
        c.Header("Cross-Origin-Resource-Policy", "same-origin")
        c.Header("Content-Security-Policy",
            "default-src 'self'; script-src 'self'; ...; frame-ancestors 'none'; "+
            "base-uri 'self'; form-action 'self'; object-src 'none'")
        if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
            c.Header("Strict-Transport-Security",
                "max-age=63072000; includeSubDomains; preload")
        }
        c.Next()
    }
}`} />
        <p>
          Error responses use a generic shape (<code>code</code>, <code>message</code>) that
          never includes a stack trace in production. Sensitive defaults: in production
          mode Sentinel WAF runs in <em>block</em> rather than <em>log</em> mode.
        </p>
      </>
    ),
  },
  {
    id: 'a03',
    title: 'A03 — Software Supply Chain Failures',
    oneLine: 'NEW in 2025 — compromised dependencies, build systems, distribution.',
    defence: (
      <>
        <p>
          Every Grit project ships with <code>.github/dependabot.yml</code> configured for
          Go modules, npm, and GitHub Actions on a weekly schedule, plus
          <code> .github/workflows/security.yml</code> running:
        </p>
        <ul className="list-disc pl-6 my-3">
          <li><strong>govulncheck</strong> — Go vulnerability scan against the Go vulnerability database.</li>
          <li><strong>pnpm audit</strong> — high+ severity gate on the frontend.</li>
          <li><strong>CodeQL</strong> — static analysis for both Go and JavaScript on every PR.</li>
        </ul>
        <p>
          These run on every push and pull-request, plus weekly so newly-disclosed CVEs
          surface even when nothing in your code changed. Go modules are pinned in
          <code> go.sum</code> with checksum verification; pnpm uses
          <code> --frozen-lockfile</code> in CI.
        </p>
      </>
    ),
  },
  {
    id: 'a04',
    title: 'A04 — Cryptographic Failures',
    oneLine: 'Weak/missing encryption exposing sensitive data.',
    defence: (
      <>
        <p>Three layers of defence ship by default:</p>
        <ul className="list-disc pl-6 my-3">
          <li>
            <strong>Passwords</strong> — bcrypt via <code>golang.org/x/crypto/bcrypt</code>
            with the library default cost. Constant-time verify. Never SHA-256, never MD5,
            never reversible encryption.
          </li>
          <li>
            <strong>Transport</strong> — HSTS preload enforced via SecurityHeaders. Any
            HTTPS request locks browsers into HTTPS-only for 2 years.
          </li>
          <li>
            <strong>JWT signing</strong> — HS256 with a 32-character minimum
            <code> JWT_SECRET</code> enforced at startup. The auth middleware verifies the
            algorithm explicitly so the classic <code>alg:none</code> attack is closed.
          </li>
        </ul>
        <p>
          Secrets live in environment variables, never in source. <code>.env</code> is in
          <code> .gitignore</code>; the scaffold ships only <code>.env.example</code>.
        </p>
      </>
    ),
  },
  {
    id: 'a05',
    title: 'A05 — Injection',
    oneLine: 'SQLi, command injection, XSS — hostile input executed.',
    defence: (
      <>
        <p>
          <strong>SQLi</strong> is closed at the protocol level: GORM uses parameterised
          queries for every <code>Where</code>, <code>First</code>, <code>Create</code>,
          <code> Update</code>, etc. The Grit codebase never builds queries by string
          concatenation; if you must use <code>db.Raw</code>, pass arguments as parameters,
          not interpolated strings.
        </p>
        <CodeBlock filename="ok and not ok" code={`// OK — parameterised
db.Where("email = ?", email).First(&user)
db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user)

// NEVER do this
db.Raw("SELECT * FROM users WHERE email = '" + email + "'")`} />
        <p>
          <strong>XSS</strong> — the React SPA escapes by default. The CSP header (A02)
          adds a second layer: even if something slips through, the browser refuses to
          execute inline script. <code>dangerouslySetInnerHTML</code> is intentionally
          banned by ESLint in scaffolded projects.
        </p>
        <p>
          <strong>SSRF</strong> (folded into A01 in 2025) — use the
          <code> internal/safefetch</code> package for any HTTP request whose URL came
          from the caller (webhooks, image-from-URL, PDF render, OEmbed):
        </p>
        <CodeBlock filename="webhook delivery" code={`import "{module}/internal/safefetch"

resp, err := safefetch.Get(ctx, userProvidedURL)
if errors.Is(err, safefetch.ErrBlocked) {
    return fmt.Errorf("URL not allowed: %w", err)
}`} />
        <p>
          <code>safefetch</code> blocks the standard SSRF targets — loopback (127.0.0.1,
          ::1), RFC1918 private ranges, link-local, CGNAT (100.64/10), AWS metadata
          (169.254.169.254 + IPv6 fd00:ec2::/32), and well-known cloud-metadata hostnames.
          It also re-validates the resolved IP at TCP-connect time via a custom dialer,
          closing the DNS-rebinding TOCTOU window.
        </p>
      </>
    ),
  },
  {
    id: 'a06',
    title: 'A06 — Insecure Design',
    oneLine: 'Flaws baked into the architecture, not the code.',
    defence: (
      <>
        <p>
          Grit&apos;s architecture is opinionated specifically to make insecure designs
          hard to express: a single auth middleware that&apos;s applied to a whole route group
          (no per-handler oversight), a single<code> Services</code> struct that owns
          DB access (no rogue connections), and a code generator that wires every new
          resource through the same handler / service / route pattern with auth checks
          pre-applied.
        </p>
        <p>
          When you generate a resource with <code>grit generate resource</code>,
          the produced handlers already call
          <code> authz.MustOwn</code> on every object access — IDOR is closed by the
          generator, not by the developer remembering.
        </p>
      </>
    ),
  },
  {
    id: 'a07',
    title: 'A07 — Authentication Failures',
    oneLine: 'Weak login, sessions, credential handling.',
    defence: (
      <>
        <ul className="list-disc pl-6 my-3">
          <li><strong>Rate limiting</strong> on <code>/api/auth/login</code> — 5 attempts per 15 min per IP via Sentinel.</li>
          <li><strong>Account-level brute-force protection</strong> — Sentinel&apos;s AuthShield locks the account after repeated failures.</li>
          <li><strong>Short access tokens</strong> (15 min) paired with longer-lived refresh tokens (7 days). The frontend ships a <code>SessionExpiryMonitor</code> that refreshes transparently.</li>
          <li><strong>JWT algorithm pinning</strong> — the verifier rejects any token whose <code>alg</code> isn&apos;t the one we signed with. <code>alg:none</code> is impossible.</li>
          <li><strong>TOTP 2FA</strong> shipped out of the box (<code>internal/totp</code>) with trusted-device fingerprinting.</li>
          <li><strong>Password reset</strong> uses single-use, time-bound tokens.</li>
        </ul>
      </>
    ),
  },
  {
    id: 'a08',
    title: 'A08 — Software & Data Integrity Failures',
    oneLine: 'Trusting unverified code/data updates.',
    defence: (
      <>
        <ul className="list-disc pl-6 my-3">
          <li><strong>Webhook signatures</strong> — the <code>internal/webhooks</code> framework verifies HMAC signatures (Stripe, GitHub, Twilio, generic) before any business logic runs.</li>
          <li><strong>Idempotency-Key middleware</strong> — caches the response on the first non-safe request and replays it on retries, preventing duplicate writes from at-least-once delivery.</li>
          <li><strong>Activity-log hash chain</strong> — every mutation is appended to a SHA-256 chain so retroactive tampering with audit logs breaks verification.</li>
          <li><strong>Frontend supply chain</strong> — <code>pnpm install --frozen-lockfile</code> in CI; lockfile committed; Dependabot raises PRs on changes.</li>
        </ul>
      </>
    ),
  },
  {
    id: 'a09',
    title: 'A09 — Security Logging & Alerting Failures',
    oneLine: "Can't detect or respond to attacks. Now includes 'alerting'.",
    defence: (
      <>
        <p>
          Grit ships <code>middleware.LogSecurityEvent(ctx, db, userID, event, ip, ua)</code>
          with a typed enum of every event that matters for SOC 2 / ISO 27001 / GDPR:
          login success/failure, logout, password change, password-reset request, TOTP
          enable/disable, TOTP challenge failure, role change, account lock, AuthZ denial,
          suspicious request.
        </p>
        <p>
          Events flow into the existing <strong>tamper-evident activity log</strong> (each
          row stores a SHA-256 hash of the previous row plus its canonical fields, so
          a retroactive deletion breaks the chain at verification time).
        </p>
        <p>
          <strong>Alerting</strong> is wired through Sentinel&apos;s AuthShield + Anomaly
          modules and visible in the <code>/sentinel/ui</code> dashboard. Hook external
          alerting (PagerDuty, Slack, email) via Sentinel webhooks for spikes in failed
          logins, AuthZ denials, or WAF blocks.
        </p>
      </>
    ),
  },
  {
    id: 'a10',
    title: 'A10 — Mishandling of Exceptional Conditions',
    oneLine: "NEW in 2025 — bad error handling, 'failing open', edge cases.",
    defence: (
      <>
        <p>
          Grit&apos;s middleware fails closed by convention:
          <code> authz.MustOwn</code> returns 404 on <em>any</em> error (DB unavailable,
          not found, wrong owner) — the request never reaches the handler in an
          ambiguous state. <code>gin.Recovery()</code> turns panics into 500s rather than
          leaking stack traces. Production builds set <code>gin.SetMode(gin.ReleaseMode)</code>.
        </p>
        <p>
          The CSRF, rate-limit, and idempotency middleware all return explicit error
          responses on the failure path — there is no &quot;continue without check&quot; branch
          that an attacker can race into.
        </p>
      </>
    ),
  },
]

export default function SecurityGuidePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Security</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Security Guide
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                How Grit defends every category of the OWASP Top 10:2025 by default,
                with the file path and code for each. Use this page as the audit checklist
                when a client asks &quot;are we secure?&quot; — walk it with them line by line.
              </p>
            </div>

            <div className="prose-grit">
              <div className="p-4 rounded-lg border border-primary/20 bg-primary/5 mb-4 flex gap-3">
                <ShieldCheck className="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-foreground/90 mb-1">
                    <strong>Secure by default.</strong> Every fresh
                    <code> grit new</code> project ships with the controls below
                    already wired. You don&apos;t add security — you opt out of pieces you
                    don&apos;t want.
                  </p>
                </div>
              </div>

              <div className="p-4 rounded-lg border border-primary/30 bg-primary/[0.06] mb-6 flex gap-3">
                <ExternalLink className="h-5 w-5 text-primary flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-foreground/90 mb-2">
                    <strong>Prefer the attacker&apos;s perspective?</strong> The companion
                    page maps each attack from JB&apos;s{' '}
                    <Link
                      href="https://jb.desishub.com/blog/defenders-handbook-how-hackers-break-in"
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-primary hover:underline"
                    >
                      Defender&apos;s Handbook
                    </Link>{' '}
                    to Grit&apos;s defence — nmap, Gobuster, Hydra, SQLi, hash cracking,
                    TOTP seed theft, SIM swap, AitM proxies, SSL strip, DDoS — with the
                    actual code path for each, plus a bonus list of defences Grit ships
                    that the handbook doesn&apos;t cover (SSRF, IDOR, CSRF, HMAC webhooks,
                    idempotency, tamper-evident audit log).
                  </p>
                  <Link
                    href="/docs/security/defenders-handbook"
                    className="text-sm font-medium text-primary hover:underline inline-flex items-center gap-1"
                  >
                    Read Defender&apos;s Handbook ↔ Grit
                    <ArrowRight className="h-3.5 w-3.5" />
                  </Link>
                </div>
              </div>

              <h2 id="threat-model">Threat model in one paragraph</h2>
              <p>
                Grit assumes a hostile public internet. The boundary of trust is the
                handler signature: anything coming in over the wire — body, query, header,
                cookie, even a JWT — is untrusted until verified. The defence pattern
                is layered (PHASE 2 §4.1 &quot;defence in depth&quot;) — Sentinel WAF + rate
                limit + auth middleware + AuthZ check + parameterised queries + output
                encoding. A single failed layer never turns into a breach.
              </p>

              <h2 id="owasp-2025">OWASP Top 10:2025 — defence by category</h2>
              <p>
                The 2025 edition is the canonical risk map (8th edition, finalised
                January 2026). For every category below, Grit ships the defence in code
                — not in documentation that asks you to remember.
              </p>

              {OWASP_2025.map((row) => (
                <section key={row.id} id={row.id} className="mt-10">
                  <h3 className="text-xl font-semibold mb-1">{row.title}</h3>
                  <p className="text-sm text-muted-foreground italic mb-4">{row.oneLine}</p>
                  {row.defence}
                </section>
              ))}

              <h2 id="hardening-checklist" className="mt-12">Pre-launch hardening checklist</h2>
              <p>
                The fastest audit pass — work the list, tick the box. Every item is a
                concrete Grit affordance.
              </p>
              <div className="overflow-x-auto my-4">
                <table className="w-full text-sm border-collapse">
                  <thead>
                    <tr className="border-b border-border/40">
                      <th className="text-left px-4 py-2 font-medium">Check</th>
                      <th className="text-left px-4 py-2 font-medium">Where</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">JWT_SECRET ≥ 32 chars, rotated per env</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">APP_ENV=production set</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">SENTINEL_ENABLED=true + WAF in block mode</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">AUTO_SEED=false in production</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">CORS_ORIGINS narrowed to actual frontends</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Every authenticated handler calls authz.MustOwn / RequireRoles</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/handlers/</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Any user-supplied URL fetched via safefetch</td>
                      <td className="px-4 py-2.5 font-mono text-xs">internal/safefetch</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">govulncheck + pnpm audit green</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.github/workflows/security.yml</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">k6 average-load passes p95 SLO</td>
                      <td className="px-4 py-2.5 font-mono text-xs">tests/k6/average-load.js</td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">Sentinel dashboard credentials changed from default</td>
                      <td className="px-4 py-2.5 font-mono text-xs">.env</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <div className="p-4 rounded-lg border border-warning/30 bg-warning/5 my-6 flex gap-3">
                <AlertTriangle className="h-5 w-5 text-warning flex-shrink-0 mt-0.5" />
                <div>
                  <p className="text-sm text-foreground/90 mb-0">
                    <strong>The defences above don&apos;t replace testing.</strong> See
                    the <Link href="/docs/testing" className="text-primary hover:underline">Testing</Link>
                    {' '}page for the k6 + pentest methodology to prove they hold.
                  </p>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/middleware" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Middleware
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/testing" className="gap-1.5">
                  Testing
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
