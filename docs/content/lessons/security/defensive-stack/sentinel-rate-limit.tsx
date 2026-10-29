import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        A login endpoint with no rate limit is a password guessing service
        you host for free. Rate limiting is the cheapest security control
        there is, and the one most often configured in a way that quietly
        does nothing.
      </p>

      <h2>Two limits, doing different jobs</h2>
      <p>
        Grit mounts Sentinel with a global per-IP limit and a tighter
        per-route limit on the endpoints worth attacking:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`// Everything, per IP.
ipLimit := &sentinel.Limit{Requests: 100, Window: 1 * time.Minute}

// The endpoints an attacker actually wants.
return map[string]sentinel.Limit{
    v + "/login":           {Requests: 30, Window: 5 * time.Minute},
    v + "/register":        {Requests: 10, Window: 15 * time.Minute},
    v + "/forgot-password": {Requests: 10, Window: 15 * time.Minute},
    v + "/reset-password":  {Requests: 10, Window: 15 * time.Minute},
    v + "/totp/verify":     {Requests: 20, Window: 5 * time.Minute},
}`}
      />

      <p>
        The numbers are not arbitrary. 30 logins in five minutes is more
        than any human needs and far less than a guessing attack wants.
        Registration is tighter because the cost of a false negative is an
        inbox full of signups rather than a locked-out user.{' '}
        <code>totp/verify</code> is on the list because a six-digit code is
        a million guesses, and without a limit that is an afternoon.
      </p>

      <h2>The mistake that makes it do nothing</h2>
      <p>
        Read the comment Grit carries at the mount point. It is there
        because it already happened:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`// Built from the prefix the router mounts, because Sentinel matches the
// request path exactly. Written as /api/auth/login these never matched a
// request, so the login limit and AuthShield never fired.
routeLimits := authRouteLimits(false)`}
      />

      <p>
        The limits were configured. The dashboard showed them. The paths were
        wrong by one version segment, so nothing ever matched, and the
        failure mode of a rate limit that does not match is{' '}
        <strong>silence</strong>. No error, no warning, and a login endpoint
        with no limit on it.
      </p>

      <TipBox tone="danger">
        <p>
          A rate limit is one of the few controls whose failure looks exactly
          like success: traffic flows, nothing is blocked, everybody is
          happy. The only way to know it works is to exceed it on purpose
          and watch for the 429.
        </p>
      </TipBox>

      <h2>Development limits, and why they are separate</h2>
      <p>
        The same function returns relaxed numbers when{' '}
        <code>APP_ENV=development</code>: 100 requests a minute instead of
        30 in five, so nobody is locked out while testing a login form.
      </p>

      <p>
        This is worth doing deliberately rather than by disabling the
        middleware in development, because a control that is off in
        development is a control nobody exercises until production. Relaxed
        and present beats absent.
      </p>

      <h2>What rate limiting cannot do</h2>
      <ul>
        <li>
          <strong>A distributed attack.</strong> Per-IP limits assume the
          attacker has few addresses. A botnet has thousands, each staying
          under your limit. That is what Sentinel&apos;s auth shield and
          anomaly detection are for: they look at the account and the
          pattern rather than the address.
        </li>
        <li>
          <strong>A correct password.</strong> Rate limiting buys time
          against guessing; it does nothing against a credential already
          leaked elsewhere. That is what the breach check and 2FA are for.
        </li>
        <li>
          <strong>Expensive single requests.</strong> One report query that
          scans ten million rows is within every limit you have. Body size
          caps and statement timeouts cover that, not request counts.
        </li>
      </ul>

      <h2>Where state lives</h2>
      <p>
        Counters have to be shared or they are per-process, and a limit of
        30 across four replicas is a limit of 120. Grit configures Sentinel
        storage explicitly for that reason, falling back to a local SQLite
        file only when nothing else is set:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/routes/routes.go"
        code={`sentinelStorage := sentinel.StorageConfig{
    Driver: sentinel.SQLite,
    DSN:    config.ResolveSQLitePath("sentinel.db"),
}`}
      />

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> With the API running in production mode,
              post wrong credentials to your login endpoint in a loop.
              Record which attempt first returns 429.
            </p>
            <p>
              <strong>Two.</strong> Change one route limit key to a path that
              does not exist, such as dropping the version segment. Restart
              and repeat the loop.
            </p>
            <p>
              <strong>Three.</strong> Write down what, if anything, told you
              the limit had stopped working.
            </p>
          </>
        }
        hint={
          <p>
            For step three, check the startup output, the Sentinel dashboard
            and the API log, and be precise about which of them mentions it.
          </p>
        }
        solution={
          <>
            <p>
              Step one gives a 429 on attempt 31 within the window. Step two
              runs the loop to completion with no limiting at all.
            </p>
            <p>
              Step three is the point: <strong>nothing</strong> tells you.
              The configuration is still valid, the dashboard still lists the
              rule, and a rule that matches no requests is indistinguishable
              from a rule nobody has hit. This is why the only honest test of
              a rate limit is to exceed it and look for the 429, and why that
              check belongs in your test suite rather than in a runbook.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="You add a rate limit to /api/auth/login, but the router mounts auth under /api/v1/auth. What happens?"
        choices={[
          {
            label: 'Sentinel normalises the path and the limit applies',
            feedback:
              'It matches the request path exactly. A near-miss is a miss.',
          },
          {
            label: 'Startup fails because the route does not exist',
            feedback:
              'A limit is a rule about a path, not a reference to a registered route. There is nothing for startup to validate it against.',
          },
          {
            label: 'Nothing matches, so the endpoint has no limit, and nothing reports it',
            correct: true,
            feedback:
              'Exactly the bug this code carries a comment about. The rule exists, the dashboard shows it, and it never fires. Confirm a limit by exceeding it, never by reading the config.',
          },
          {
            label: 'The global per-IP limit still protects the endpoint fully',
            feedback:
              'It helps, and it is far looser: 100 a minute instead of 30 in five. Enough for a guessing attack to make real progress.',
          },
        ]}
      />
    </>
  )
}
