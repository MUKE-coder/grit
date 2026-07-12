import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/pulse')

export default function PulsePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Pulse (Observability)
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Every Grit project ships with{' '}
                <a href="https://github.com/MUKE-coder/pulse" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Pulse</a>,
                a self-hosted observability and performance monitoring SDK. It provides request tracing,
                database monitoring, runtime metrics, error tracking, health checks, alerting, and a
                real-time dashboard &mdash; all with a single <code>Mount()</code> call.
              </p>
            </div>

            <div className="prose-grit">
              {/* What You Get */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  What You Get
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  After starting your API, the Pulse dashboard is available at{' '}
                  <code>/pulse/ui/</code> (default credentials: <code>admin</code> / <code>pulse</code>).
                </p>

                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Feature</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Request Tracing</td>
                        <td className="py-2 text-muted-foreground">Automatic trace IDs, slow request detection, latency tracking</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Database Monitoring</td>
                        <td className="py-2 text-muted-foreground">GORM query capture, N+1 detection, connection pool stats</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Runtime Metrics</td>
                        <td className="py-2 text-muted-foreground">Heap memory, goroutines, GC pauses, leak detection</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Error Tracking</td>
                        <td className="py-2 text-muted-foreground">Panic recovery, stack traces, request body capture, error fingerprinting</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Health Checks</td>
                        <td className="py-2 text-muted-foreground">Kubernetes-compatible <code>/live</code> and <code>/ready</code> endpoints</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Alerting</td>
                        <td className="py-2 text-muted-foreground">Threshold-based rules with Slack, Discord, email, and webhook notifications</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4 font-medium text-foreground/80">Prometheus</td>
                        <td className="py-2 text-muted-foreground">18+ metrics at <code>/pulse/metrics</code></td>
                      </tr>
                      <tr>
                        <td className="py-2 pr-4 font-medium text-foreground/80">Dashboard</td>
                        <td className="py-2 text-muted-foreground">React 19 UI with 8 pages: Overview, Routes, Database, Errors, Runtime, Health, Alerts, Settings</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Pulse is mounted in <code>routes.go</code> and configured via environment variables:
                </p>

                <CodeBlock
                  language="bash"
                  filename=".env"
                  code={`# Observability — Pulse
PULSE_ENABLED=true
PULSE_USERNAME=admin
PULSE_PASSWORD=pulse
PULSE_STORAGE=memory                 # "memory" (default) | "sqlite"
PULSE_STORAGE_DSN=pulse.db           # path used when PULSE_STORAGE=sqlite`}
                />

                <p className="text-muted-foreground leading-relaxed mt-6 mb-4">
                  Pulse v1.0 is mounted with <strong>functional options</strong> and a{' '}
                  <code>context.Context</code> (the context drives clean shutdown of the
                  dashboard WebSocket and background samplers). The scaffolded mount call:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`// Mount Pulse observability
if cfg.PulseEnabled {
    pulseOpts := []pulse.Option{
        pulse.WithAppName(cfg.AppName),
        pulse.WithCredentials(cfg.PulseUsername, cfg.PulsePassword),
        pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),
        pulse.WithPrometheus(),
        // Pulse's error middleware snapshots a request-body snippet for context,
        // but restores ONLY the first 4096 bytes — silently truncating larger
        // uploads/JSON. Disable capture so the full body reaches the handler.
        pulse.WithRequestBodyCaptureDisabled(),
    }
    if cfg.IsDevelopment() {
        pulseOpts = append(pulseOpts, pulse.WithDevMode())
    }
    if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
        pulseOpts = append(pulseOpts, pulse.WithSQLite(cfg.PulseStorageDSN))
    }

    p := pulse.Mount(context.Background(), r, db, pulseOpts...)

    // Register health checks for connected services
    if svc.Cache != nil {
        p.AddHealthCheck(pulse.HealthCheck{
            Name:     "redis",
            Type:     "redis",
            Critical: false,
            CheckFunc: func(ctx context.Context) error {
                return svc.Cache.Client().Ping(ctx).Err()
            },
        })
    }
}`}
                />

                <div className="overflow-x-auto my-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Option</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Effect</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[13px]">
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithAppName(name)</td><td className="py-2 text-muted-foreground font-sans">Labels the dashboard</td></tr>
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithCredentials(user, pass)</td><td className="py-2 text-muted-foreground font-sans">Dashboard basic-auth login</td></tr>
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithExcludePaths(globs…)</td><td className="py-2 text-muted-foreground font-sans">Skip tracing for these routes</td></tr>
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithPrometheus()</td><td className="py-2 text-muted-foreground font-sans">Expose /pulse/metrics</td></tr>
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithDevMode()</td><td className="py-2 text-muted-foreground font-sans">Verbose, dev-only behaviour</td></tr>
                      <tr className="border-b border-border/20"><td className="py-2 pr-4">WithSQLite(dsn)</td><td className="py-2 text-muted-foreground font-sans">Persist data across restarts</td></tr>
                      <tr><td className="py-2 pr-4">WithRequestBodyCaptureDisabled()</td><td className="py-2 text-muted-foreground font-sans">Stop the 4096-byte body truncation</td></tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Endpoints */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Endpoints
                </h2>

                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Endpoint</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Description</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/ui/</code></td>
                        <td className="py-2 text-muted-foreground">Dashboard UI</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/health</code></td>
                        <td className="py-2 text-muted-foreground">Full health status (public)</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/health/live</code></td>
                        <td className="py-2 text-muted-foreground">Kubernetes liveness probe</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/health/ready</code></td>
                        <td className="py-2 text-muted-foreground">Kubernetes readiness probe</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/metrics</code></td>
                        <td className="py-2 text-muted-foreground">Prometheus metrics</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>POST /pulse/api/auth/login</code></td>
                        <td className="py-2 text-muted-foreground">Dashboard authentication</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4"><code>GET /pulse/api/overview</code></td>
                        <td className="py-2 text-muted-foreground">Dashboard summary (auth required)</td>
                      </tr>
                      <tr>
                        <td className="py-2 pr-4"><code>WS /pulse/ws/live</code></td>
                        <td className="py-2 text-muted-foreground">WebSocket live updates</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Health Checks */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Custom Health Checks
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The <code>Mount()</code> function returns a Pulse instance. Use it to register
                  health checks for any dependency:
                </p>

                <CodeBlock
                  language="go"
                  code={`p := pulse.Mount(context.Background(), r, db, pulseOpts...)

// Redis health check
p.AddHealthCheck(pulse.HealthCheck{
    Name:     "redis",
    Type:     "redis",
    Critical: false,  // non-critical = won't fail /ready
    CheckFunc: func(ctx context.Context) error {
        return redisClient.Ping(ctx).Err()
    },
})

// External API health check
p.AddHealthCheck(pulse.HealthCheck{
    Name:     "payment-api",
    Type:     "http",
    Critical: true,  // critical = will fail /ready if down
    CheckFunc: func(ctx context.Context) error {
        resp, err := http.Get("https://api.stripe.com/v1/health")
        if err != nil {
            return err
        }
        defer resp.Body.Close()
        if resp.StatusCode != 200 {
            return fmt.Errorf("status %d", resp.StatusCode)
        }
        return nil
    },
})`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  The database health check is registered automatically when you pass a <code>*gorm.DB</code> to <code>Mount()</code>.
                </p>
              </div>

              {/* Alerting */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Alerting
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Pulse ships with 5 default alert rules (high latency, high error rate, high memory,
                  goroutine leak, health check failure) that run out of the box &mdash; no
                  configuration needed. They fire against the metrics Pulse already collects and
                  surface in the dashboard.
                </p>

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mb-6">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Custom rules &amp; channels.</strong> The
                    scaffold enables the built-in rules and doesn&apos;t wire any notification
                    channel by default. To add your own rules or a Slack/Discord/webhook channel,
                    pass the matching Pulse option when you build <code>pulseOpts</code> in
                    routes.go &mdash; see the{' '}
                    <a href="https://github.com/MUKE-coder" className="text-primary hover:underline">Pulse package docs</a>{' '}
                    for the current option set, since these evolve with the library.
                  </p>
                </div>

                <h3 className="text-lg font-semibold mt-6 mb-3">Built-in Default Rules</h3>

                <div className="overflow-x-auto mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30">
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Rule</th>
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Metric</th>
                        <th className="text-left py-2 pr-4 text-foreground/80 font-medium">Condition</th>
                        <th className="text-left py-2 text-foreground/80 font-medium">Severity</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4">high_latency</td>
                        <td className="py-2 pr-4 text-muted-foreground">P95 latency</td>
                        <td className="py-2 pr-4 text-muted-foreground">&gt; 2000ms for 5min</td>
                        <td className="py-2 text-muted-foreground">critical</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4">high_error_rate</td>
                        <td className="py-2 pr-4 text-muted-foreground">Error rate</td>
                        <td className="py-2 pr-4 text-muted-foreground">&gt; 10% for 3min</td>
                        <td className="py-2 text-muted-foreground">critical</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4">high_memory</td>
                        <td className="py-2 pr-4 text-muted-foreground">Heap alloc</td>
                        <td className="py-2 pr-4 text-muted-foreground">&gt; 500MB for 5min</td>
                        <td className="py-2 text-muted-foreground">warning</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="py-2 pr-4">goroutine_leak</td>
                        <td className="py-2 pr-4 text-muted-foreground">Goroutine growth</td>
                        <td className="py-2 pr-4 text-muted-foreground">&gt; 100/hr for 10min</td>
                        <td className="py-2 text-muted-foreground">warning</td>
                      </tr>
                      <tr>
                        <td className="py-2 pr-4">health_check_failure</td>
                        <td className="py-2 pr-4 text-muted-foreground">Health status</td>
                        <td className="py-2 pr-4 text-muted-foreground">unhealthy for 2min</td>
                        <td className="py-2 text-muted-foreground">critical</td>
                      </tr>
                    </tbody>
                  </table>
                </div>
              </div>

              {/* Storage */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Data Storage
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  By default, Pulse uses in-memory ring buffers &mdash; fastest, but data resets on
                  restart. For persistence across restarts, set <code>PULSE_STORAGE=sqlite</code> in
                  your <code>.env</code>; the mount call then adds the{' '}
                  <code>pulse.WithSQLite()</code> option automatically:
                </p>

                <CodeBlock
                  language="bash"
                  filename=".env"
                  code={`PULSE_STORAGE=sqlite       # switch from in-memory to persistent
PULSE_STORAGE_DSN=pulse.db # SQLite file (or ":memory:")`}
                />

                <CodeBlock
                  language="go"
                  filename="internal/routes/routes.go"
                  code={`if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
    pulseOpts = append(pulseOpts, pulse.WithSQLite(cfg.PulseStorageDSN))
}`}
                />
              </div>

              {/* Prometheus */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Prometheus Integration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Pulse exposes 18+ metrics in Prometheus exposition format at <code>/pulse/metrics</code>.
                  Connect it to Grafana or any Prometheus-compatible tool.
                </p>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Key metrics include:
                </p>
                <ul className="space-y-1 text-muted-foreground mb-6">
                  <li><code>pulse_http_requests_total</code> &mdash; Total HTTP requests by method, path, status</li>
                  <li><code>pulse_http_request_duration_seconds</code> &mdash; Request duration histogram</li>
                  <li><code>pulse_http_error_rate</code> &mdash; Current error rate percentage</li>
                  <li><code>pulse_runtime_goroutines</code> &mdash; Active goroutine count</li>
                  <li><code>pulse_runtime_heap_bytes</code> &mdash; Heap allocation in bytes</li>
                  <li><code>pulse_db_query_duration_seconds</code> &mdash; Database query duration</li>
                  <li><code>pulse_health_check_status</code> &mdash; Health check status (1=healthy, 0=unhealthy)</li>
                </ul>
              </div>

              {/* Disabling */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Disabling Pulse
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Set <code>PULSE_ENABLED=false</code> in your <code>.env</code> to disable Pulse entirely.
                  This skips the mount and adds zero overhead to your application.
                </p>
              </div>
            </div>

            {/* Navigation */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild>
                <Link href="/docs/backend/api-docs">
                  <ArrowLeft className="h-3.5 w-3.5 mr-1" />
                  API Documentation
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild>
                <Link href="/docs/admin/overview">
                  Admin Overview
                  <ArrowRight className="h-3.5 w-3.5 ml-1" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
