import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Sentinel is Grit&apos;s in-process Web Application Firewall (WAF). It
        does rate limiting, AuthShield (brute-force lockouts), Anomaly
        detection, and Geo blocking — all in middleware that runs in
        sub-millisecond time. This lesson covers what it ships and how to
        configure it.
      </p>

      <h2>The Sentinel dashboard</h2>
      <p>
        Open <code>http://localhost:8080/sentinel/ui</code>. Log in with{' '}
        <code>SENTINEL_USERNAME</code> / <code>SENTINEL_PASSWORD</code> from
        your <code>.env</code>. You&apos;ll see:
      </p>
      <ul>
        <li>Live request rate + p95 latency</li>
        <li>Blocked requests (by reason: WAF rule, rate limit, AuthShield)</li>
        <li>Recent suspicious IPs</li>
        <li>Per-route rate-limit counters</li>
      </ul>

      <h2>What it intercepts</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go (excerpt)"
        code={`sentinel.MountE(r, db, sentinel.Config{
    Dashboard: sentinel.DashboardConfig{
        Username:              cfg.SentinelUsername,
        Password:              cfg.SentinelPassword,
        SecretKey:             cfg.SentinelSecretKey,
        AllowInsecureDefaults: cfg.AppEnv != "production",
    },
    WAF: sentinel.WAFConfig{
        Enabled: true,
        Mode:    sentinel.ModeBlock,  // or ModeLog
    },
    RateLimit: sentinel.RateLimitConfig{
        Enabled: true,
        ByIP:    &sentinel.Limit{Requests: 100, Window: time.Minute},
        ByRoute: map[string]sentinel.Limit{
            "/api/auth/login": {Requests: 5, Window: 15 * time.Minute},
        },
    },
    AuthShield: sentinel.AuthShieldConfig{
        Enabled:    true,
        LoginRoute: "/api/auth/login",
    },
    Anomaly: sentinel.AnomalyConfig{Enabled: true},
    Geo:     sentinel.GeoConfig{Enabled: true},
})`}
      />

      <h2>WAF — basic injection / scanner blocks</h2>
      <p>
        OWASP CRS-lite ruleset embedded. Blocks:
      </p>
      <ul>
        <li>SQLi probes (<code>&apos; OR 1=1--</code>, UNION SELECT, time-based)</li>
        <li>XSS payloads in query / body</li>
        <li>Path traversal (<code>../../etc/passwd</code>)</li>
        <li>Dotfile probes (<code>.git</code>, <code>.env</code>, <code>.bak</code>)</li>
      </ul>
      <p>
        Two modes: <code>ModeBlock</code> (return 403 + record) or{' '}
        <code>ModeLog</code> (record but pass through — useful during
        rollout).
      </p>

      <h2>Rate limit — per-IP + per-route</h2>
      <p>The default config:</p>
      <ul>
        <li>100 req/min per IP across all routes</li>
        <li>5 req/15min per IP on <code>/api/auth/login</code></li>
        <li>Custom rates per route via the <code>ByRoute</code> map</li>
      </ul>
      <p>
        Exceeded → 429 Too Many Requests. Counters live in Redis (shared
        across API replicas).
      </p>

      <h2>AuthShield — account-level brute-force</h2>
      <p>
        After N failed logins for an email, the account enters a cool-down
        regardless of source IP. This defends against distributed
        credential-stuffing where every attempt comes from a different IP.
      </p>
      <CodeBlock
        language="text"
        code={`6 failed login attempts for alex@example.com
→ AuthShield locks the account for 15 minutes
→ Even valid password retries return "invalid credentials"
→ Audit log records the lock event
→ Sentinel dashboard shows it`}
      />

      <h2>Anomaly + Geo</h2>
      <ul>
        <li>
          <strong>Anomaly:</strong> behavioural rules — rapid 4xx bursts,
          suspicious user-agent strings, scanner signatures. Bots get
          throttled before they exhaust your rate-limit budget.
        </li>
        <li>
          <strong>Geo:</strong> country-level allow/block list. &quot;Block
          everything outside the US&quot; or &quot;flag high-fraud-rate countries
          for manual review&quot; — both are config flags.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Trust your proxy header chain.</strong> If you&apos;re behind a
        reverse proxy (Traefik, Cloudflare), set{' '}
        <code>SENTINEL_TRUSTED_PROXIES</code> to your LB&apos;s CIDR. Without
        it, the WAF sees every request as coming from the LB&apos;s IP — your
        rate limit immediately exhausts.
      </TipBox>

      <h2>The dashboard credential gate</h2>
      <p>
        Sentinel refuses to mount in production with the default password{' '}
        <code>sentinel</code>. v3.25.1 of Grit fixes this by generating a
        random password at scaffold time, but if you somehow kept the
        default, the API won&apos;t start in <code>APP_ENV=production</code>.
        Set strong credentials or set <code>DEMO_MODE=true</code> for the
        public-demo carve-out.
      </p>

      <KnowledgeCheck
        question="Six failed logins in a row from your IP. Sentinel takes action. What's the response?"
        choices={[
          {
            label: 'API returns 200 OK but silently doesn\'t actually authenticate you',
            feedback:
              "Wrong — that'd be confusing and give attackers signal. Failed login is always 401.",
          },
          {
            label: 'API returns 429 Too Many Requests for the next 15 minutes',
            correct: true,
            feedback:
              "Right — rate limit on /api/auth/login is 5 per 15 minutes by default. After 5, you get 429 until the window resets. The account-level AuthShield also activates if multiple IPs are at it.",
          },
          {
            label: 'API blocks your IP forever',
            feedback:
              "Wrong — rate limits are time-windowed by design. Forever-blocks are a manual action via the Sentinel dashboard.",
          },
          {
            label: 'Sentinel emails you to confirm',
            feedback:
              "Not built-in. Could be wired via the audit log + a job, but that's not the default behaviour.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Trip the rate limit:</p>
            <ol>
              <li>Run 6 login attempts in quick succession with curl:</li>
            </ol>
            <CodeBlock
              terminal
              code={`for i in 1 2 3 4 5 6; do
  curl -s -o /dev/null -w "%{http_code}\\n" -X POST http://localhost:8080/api/auth/login \\
    -H 'Content-Type: application/json' \\
    -d '{"email":"alex@example.com","password":"wrong"}'
done`}
            />
            <ol start={2}>
              <li>You should see five 401s, then a 429.</li>
              <li>Open the Sentinel dashboard and look at the rate-limit page.</li>
              <li>Paste the curl output and a dashboard screenshot in <code>notes.md</code>.</li>
            </ol>
          </>
        }
        hint={
          <>
            If you see 6 × 401 (no 429), check that{' '}
            <code>SENTINEL_ENABLED=true</code> in <code>.env</code> and the
            API is running in production mode (<code>APP_ENV=production</code>).
            Rate limiting is off by default in dev.
          </>
        }
        solution={
          <>
            <p>You should see:</p>
            <CodeBlock
              language="text"
              code={`401
401
401
401
401
429`}
            />
            <p>
              Sentinel did its job. The 6th request hit the per-route rate
              limit. The audit log records the event; the dashboard shows
              the IP.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Security on. Now to see <em>what&apos;s actually happening</em> —
        Pulse&apos;s p50/p95/p99 dashboards, per-route metrics, slow-query
        visibility.
      </p>
    </>
  )
}
