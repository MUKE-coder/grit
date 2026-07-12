import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/security')

export default function SecurityPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Batteries</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Security (Sentinel)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit project ships with Sentinel -- a production-grade security intelligence suite
                that provides a Web Application Firewall, rate limiting, brute-force protection, anomaly detection,
                security headers, and a real-time threat dashboard. Security is not an afterthought.
              </p>
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
                    { title: 'Web Application Firewall', desc: 'Detects SQL injection, XSS, path traversal, command injection, SSRF, XXE, and more' },
                    { title: 'Rate Limiting', desc: 'Per-IP, per-user, per-route, and global limits with sliding window strategy' },
                    { title: 'Auth Shield', desc: 'Brute-force protection with per-IP lockouts and credential stuffing detection' },
                    { title: 'Security Headers', desc: 'CSP, HSTS, X-Frame-Options, X-Content-Type-Options, Referrer-Policy' },
                    { title: 'Anomaly Detection', desc: 'Off-hours access, velocity anomalies, impossible travel, data exfiltration' },
                    { title: 'IP Geolocation', desc: 'Automatic geo-lookup for every request with country and city data' },
                    { title: 'Threat Dashboard', desc: 'Embedded React dashboard with real-time WebSocket updates at /sentinel/ui' },
                    { title: 'Performance Monitoring', desc: 'Per-route latency tracking with p50/p95/p99 metrics' },
                  ].map((item) => (
                    <div key={item.title} className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                      <h3 className="text-sm font-semibold text-foreground mb-1">{item.title}</h3>
                      <p className="text-xs text-muted-foreground/60 leading-relaxed">{item.desc}</p>
                    </div>
                  ))}
                </div>
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
SENTINEL_SECRET_KEY=change-me-in-production`} />

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
                      ].map((row) => (
                        <tr key={row.name} className="border-b border-border/50">
                          <td className="py-2 pr-4 font-mono text-xs text-primary/60">{row.name}</td>
                          <td className="py-2 pr-4 font-mono text-xs text-muted-foreground">{row.def}</td>
                          <td className="py-2 text-muted-foreground">{row.desc}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
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
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  WAF Modes
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Web Application Firewall has two modes:
                </p>
                <div className="grid gap-3 sm:grid-cols-2 mb-4">
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <code className="text-sm font-mono text-primary/70 font-medium">ModeLog</code>
                    <span className="ml-2 text-xs text-warning bg-warning/10 px-1.5 py-0.5 rounded-full">dev default</span>
                    <p className="text-xs text-muted-foreground/60 mt-1 leading-relaxed">
                      Detects and logs threats but does not block requests. Use during development and initial
                      deployment to identify false positives before switching to block mode.
                    </p>
                  </div>
                  <div className="rounded-lg border border-border/30 bg-card/30 px-4 py-3">
                    <code className="text-sm font-mono text-primary/70 font-medium">ModeBlock</code>
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
                      <span className="text-primary/50">&#10003;</span>
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
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
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
          </div>
        </div>
      </main>
    </div>
  )
}
