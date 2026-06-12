import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to write a real k6 test. We&apos;ll hit your local Grit
        API&apos;s <code>/healthz</code> endpoint with 10 virtual users
        for 30 seconds, then read the output.
      </p>

      <h2>The script</h2>
      <CodeBlock
        language="js"
        filename="tests/k6/smoke.js"
        code={`import http from 'k6/http'
import { check, sleep } from 'k6'

// Configuration
export const options = {
  vus: 10,             // virtual users
  duration: '30s',     // how long
  thresholds: {
    http_req_duration: ['p(95)<200'],   // 95% of requests under 200ms
    http_req_failed:   ['rate<0.01'],   // <1% errors
  },
}

// The test function runs in a loop for each VU
export default function () {
  const res = http.get('http://localhost:8080/healthz')

  check(res, {
    'status is 200':   (r) => r.status === 200,
    'body has "ok"':   (r) => r.body.includes('ok'),
  })

  sleep(1)   // pause 1s between iterations per VU
}`}
      />
      <p>Three sections, each non-negotiable:</p>
      <ul>
        <li>
          <strong><code>options</code></strong> — how many VUs, for how
          long, with what thresholds.
        </li>
        <li>
          <strong>The default function</strong> — what each VU does
          repeatedly until time runs out.
        </li>
        <li>
          <strong><code>check</code></strong> — assertions on the
          response. Failed checks count toward <code>checks</code>{' '}
          metric.
        </li>
      </ul>

      <h2>Run it</h2>
      <CodeBlock
        language="bash"
        code={`# Make sure your API is running
go run apps/api/cmd/server/main.go

# In another terminal
k6 run tests/k6/smoke.js`}
      />

      <h2>Reading the output</h2>
      <CodeBlock
        language="text"
        code={`     checks.........................: 100.00% ✓ 600  ✗ 0
     data_received..................: 240 kB  8.0 kB/s
     data_sent......................: 60 kB   2.0 kB/s
     http_req_blocked...............: avg=15µs    p(95)=24µs
     http_req_connecting............: avg=4µs     p(95)=0s
     http_req_duration..............: avg=12ms    p(95)=19ms    ← key metric
       { expected_response:true }...: avg=12ms    p(95)=19ms
     http_req_failed................: 0.00%   ✓ 0    ✗ 300
     http_req_receiving.............: avg=80µs    p(95)=120µs
     http_req_sending...............: avg=22µs    p(95)=39µs
     http_req_tls_handshaking.......: avg=0s      p(95)=0s
     http_req_waiting...............: avg=11ms    p(95)=18ms
     http_reqs......................: 300     10/s              ← throughput
     iteration_duration.............: avg=1.01s   p(95)=1.02s
     iterations.....................: 300     10/s
     vus............................: 10
     vus_max........................: 10`}
      />
      <p>Four metrics matter on a smoke test:</p>
      <ul>
        <li>
          <strong><code>checks</code></strong> — 100% means every
          response was healthy.
        </li>
        <li>
          <strong><code>http_req_duration</code></strong> — latency.
          The p(95) tells you the 95th-percentile latency — the bar
          most teams set SLOs against.
        </li>
        <li>
          <strong><code>http_req_failed</code></strong> — error rate.
          Anything &gt; 0% on a healthy smoke test is a red flag.
        </li>
        <li>
          <strong><code>http_reqs</code></strong> — total + rate.
          Useful for capacity planning.
        </li>
      </ul>

      <h2>Threshold pass / fail</h2>
      <p>
        At the bottom of the output:
      </p>
      <CodeBlock
        language="text"
        code={`✓ http_req_duration..p(95)<200
✓ http_req_failed....rate<0.01`}
      />
      <p>
        ✓ means the threshold held. ✗ means a breach — and k6 exits
        with non-zero status. That&apos;s how CI knows your PR
        regressed the API.
      </p>

      <h2>VUs vs. iterations — the mental model</h2>
      <p>
        VUs are virtual users — independent loops running your default
        function. With 10 VUs, ten copies of the function are running
        concurrently. Each VU loops until the duration ends.
      </p>
      <p>
        An iteration is one trip through the default function. With{' '}
        <code>sleep(1)</code> and a 30s duration, each VU does about
        30 iterations, so 10 VUs × 30 iterations = ~300 total. That
        lines up with what we saw above.
      </p>

      <TipBox tone="warning">
        <strong>VUs are NOT real concurrent users in a browser
        sense.</strong> 100 VUs in k6 sends about as many requests as
        100 humans actively clicking. If your real app has 10,000
        casual users browsing, that&apos;s probably more like 200-500
        VUs in load-testing terms (because most users are idle most of
        the time).
      </TipBox>

      <h2>Adding auth — a more realistic test</h2>
      <CodeBlock
        language="js"
        filename="tests/k6/load.js"
        code={`import http from 'k6/http'
import { check, sleep } from 'k6'

const BASE = 'http://localhost:8080'

export const options = {
  vus: 50,
  duration: '2m',
  thresholds: {
    http_req_duration: ['p(95)<300'],
    http_req_failed:   ['rate<0.01'],
  },
}

// Runs ONCE per VU at start
export function setup() {
  const r = http.post(\`\${BASE}/api/auth/login\`,
    JSON.stringify({ email: 'load-test@example.com', password: 'secret123' }),
    { headers: { 'Content-Type': 'application/json' } },
  )
  return { token: r.json('data.token') }
}

// Runs in a loop per VU; receives setup data as argument
export default function (data) {
  const res = http.get(\`\${BASE}/api/me\`, {
    headers: { Authorization: \`Bearer \${data.token}\` },
  })
  check(res, { 'me 200': (r) => r.status === 200 })
  sleep(1)
}`}
      />
      <p>
        <code>setup()</code> runs once before VUs start; its return
        value is passed to every iteration as <code>data</code>. Use
        it for: log in once, share the token; create a fixture user,
        share the ID.
      </p>

      <KnowledgeCheck
        question="Your smoke test shows p(95) http_req_duration of 47ms. Then you push a refactor and re-run: p(95) is now 230ms. Threshold is `p(95)<200`. What does k6 do?"
        choices={[
          {
            label: "Prints a warning but exits 0",
            feedback:
              "No — a breached threshold makes k6 exit non-zero, by design. CI catches this.",
          },
          {
            label: "Exits with a non-zero status code, failing CI",
            correct: true,
            feedback:
              "Right — that's the whole point of thresholds. A failed threshold breaks the build, so a regression PR can't merge without explicit acknowledgement.",
          },
          {
            label: 'Auto-rolls back the deploy',
            feedback:
              "k6 doesn't touch your infra. It tells the truth; humans/CI decide what to do.",
          },
          {
            label: 'Logs the regression to k6 cloud',
            feedback: "Only if you opted in to k6 cloud. Locally, it just exits non-zero.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Smoke-test your real API:</p>
            <ol>
              <li>Save the smoke.js script from this lesson.</li>
              <li>
                Boot your local API:{' '}
                <code>cd apps/api && go run cmd/server/main.go</code>.
              </li>
              <li>Run k6: <code>k6 run tests/k6/smoke.js</code>.</li>
              <li>
                Capture the p(95) for http_req_duration. Paste in{' '}
                <code>notes.md</code>.
              </li>
              <li>
                Now write a SECOND script for{' '}
                <code>/api/users</code> (requires auth). Use the{' '}
                <code>setup</code> pattern from this lesson.
              </li>
              <li>Run it. Compare p(95) to the smoke test.</li>
            </ol>
          </>
        }
        hint={
          <>
            For the authed test, create a fixture user in the DB first
            (via seeder or the admin panel) and use those credentials.
            Don&apos;t register a new user inside <code>setup</code>{' '}
            every time — duplicates break.
          </>
        }
        solution={
          <>
            <p>
              You should see <code>/healthz</code> at single-digit ms
              and an authed endpoint at higher latency (DB hit + auth
              check). That delta IS the cost of auth + a database call
              — useful to know.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>The five test types</strong>. We&apos;ll
        build a full suite: load, stress, spike, soak, with
        appropriate thresholds for each.
      </p>
    </>
  )
}
