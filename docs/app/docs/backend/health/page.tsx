import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { Callout } from '@/components/callout'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/health')

export default function HealthPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Health checks</h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                <code>GET /api/health</code> reports every component this deployment depends on, in
                four states rather than a boolean, so &quot;we have no Redis&quot; and &quot;Redis
                is down&quot; stop being the same answer.
              </p>
            </div>

            <div className="prose-grit">
              <h2>The four states</h2>
              <p>
                A boolean cannot distinguish a component that is broken from one nobody asked for,
                and a dashboard that cannot tell them apart either cries wolf about a dependency you
                never configured or stays quiet about one that just died.
              </p>
              <CodeBlock
                language="text"
                code={`ok        wired up, and nothing known to be wrong
degraded  wired up, and not working, or working worse than it should
off       deliberately not wired up; this is not a problem
unknown   wired up, but this probe could not find out`}
              />
              <p>
                Only <code>degraded</code> lowers the deployment&apos;s overall status. A deployment
                with no Redis is healthy. One whose queue depth has not been sampled since the last
                restart is healthy until something says otherwise. That rule lives in one place,{' '}
                <code>health.Overall</code>, so a new component cannot quietly invent its own.
              </p>

              <h2>What a response looks like</h2>
              <p>
                Each component carries its state, a <code>detail</code> an operator can act on, and
                whatever that component measures, flattened into the same object.
              </p>
              <CodeBlock
                language="json"
                filename="GET /api/health"
                code={`{
  "status": "ok",
  "version": "0.1.0",
  "api":      { "state": "ok", "ok": true, "detail": "answering requests" },
  "database": { "state": "ok", "ok": true, "latency_ms": 1, "tables": 42 },
  "redis":    { "state": "off", "ok": false,
                "detail": "REDIS_URL is empty, so caching, jobs and cron are off on purpose" },
  "jobs":     { "state": "off", "ok": false, "detail": "no queue: asynq needs Redis" },
  "storage":  { "state": "ok", "ok": true, "configured": true, "driver": "local" },
  "email":    { "state": "ok", "ok": true, "configured": true, "driver": "log" },
  "events":   { "state": "ok", "ok": true, "subscribers": 2, "queued": 0,
                "capacity": 1024, "dropped": 0 },
  "realtime": { "connections": 0, "users": 0, "channels": 0 }
}`}
              />
              <p>
                <code>ok</code> is still there, and is exactly <code>state === &quot;ok&quot;</code>
                . Load balancers, uptime probes and the desktop client&apos;s heartbeat read it, and
                it means what it always meant.
              </p>

              <Callout type="note" title="Redis is the one probe that asks about the environment">
                A missing cache has two causes. <code>REDIS_URL=</code> turned it off on purpose, and
                that is <code>off</code>. Anything else means the API dialled Redis at boot, was
                refused, and carried on with caching, jobs and cron disabled. On a laptop that is
                still <code>off</code>, with a detail saying so. In production it is{' '}
                <code>degraded</code>, because there it is an incident.
              </Callout>

              <h2>Adding your own component</h2>
              <p>
                Anything the route file does not already hold a handle to reports itself through{' '}
                <code>health.Register</code>: a plugin, a client you wired in <code>main.go</code>, a
                dependency only your code knows about. Register once while the application is being
                built, never per request.
              </p>
              <CodeBlock
                language="go"
                filename="cmd/server/main.go"
                code={`health.Register("billing", func() health.Report {
    if !gateway.Configured() {
        return health.Off("no payment gateway; set BILLING_KEY")
    }
    if err := gateway.LastPing(); err != nil {
        return health.Degraded("the gateway refused our last call", nil)
    }
    return health.OK("", map[string]any{"charges_settled": settled.Load()})
})`}
              />
              <p>
                It appears beside the framework&apos;s own components, counts towards the overall
                status under the same rule, and shows up on the admin&apos;s System Health page
                without touching the handler. Registering the same name twice replaces the first, so
                a reload cannot show one component twice.
              </p>

              <h2>What a probe may do</h2>
              <p>
                Nothing slow. A probe reports state its subsystem already holds, or makes one bounded
                call. It must not block on a dependency that is already unwell:{' '}
                <code>/api/health</code> is exactly what somebody reaches for during an outage, and a
                probe that waits on a hung database turns the one endpoint that could explain the
                outage into another symptom of it. The framework&apos;s own database and Redis probes
                are bounded at 500ms; the queue counts come from a snapshot refreshed in the
                background, because counting asynq&apos;s keys inside Redis on every probe stalled
                every other Redis client while it ran.
              </p>
              <p>
                A probe that panics is reported as <code>unknown</code> rather than taking the
                endpoint down. Somebody else&apos;s broken probe is not a reason to stop answering
                &quot;is the database up&quot;.
              </p>

              <Callout type="warning" title="No probe reports a secret">
                <code>/api/health</code> is reachable without authentication so a load balancer can
                use it. Report a credential by presence, never by value:{' '}
                <code>{'{"signing_key": true}'}</code>, not the key. The framework&apos;s own probes
                follow the same rule, which is why storage reports <code>local</code> or{' '}
                <code>s3</code> and not the bucket, and why a failed database ping says &quot;the
                database did not answer a ping&quot; while the driver error, which names the host and
                often the user, goes to the log.
              </Callout>

              <h2>Where to see it</h2>
              <p>
                The admin panel renders it at <code>/system/health</code>, one card per component,
                with <code>off</code> and <code>unknown</code> drawn in neutral rather than in red.
                The endpoint is registered at both <code>/api/health</code> and{' '}
                <code>/api/v1/health</code>: the first for probes, load balancers and the desktop
                client&apos;s heartbeat, which are configured outside your repo, and the second for
                the frontends, whose client rewrites every call to the versioned path.
              </p>

              <div className="mt-16 flex items-center justify-between border-t border-border/30 pt-8">
                <Link href="/docs/backend/pulse">
                  <Button variant="outline" size="sm" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Pulse (Observability)
                  </Button>
                </Link>
                <Link href="/docs/backend/feature-flags">
                  <Button variant="outline" size="sm" className="gap-2">
                    Feature Flags
                    <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
