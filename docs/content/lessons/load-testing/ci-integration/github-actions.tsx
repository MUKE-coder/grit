import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        A load test you remember to run is a load test you run twice. The
        value is in running it on every pull request, so the regression is
        attributed to the change that caused it rather than discovered three
        weeks later with forty commits to bisect.
      </p>

      <h2>The workflow</h2>
      <CodeBlock
        language="yaml"
        filename=".github/workflows/loadtest.yml"
        code={`name: Load test

on:
  pull_request:
  push:
    branches: [main]

permissions:
  contents: read

jobs:
  k6:
    name: smoke + average load
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:17
        env:
          POSTGRES_USER: grit
          POSTGRES_PASSWORD: grit
          POSTGRES_DB: grit
        ports: ['5432:5432']
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: apps/api/go.mod
          cache-dependency-path: apps/api/go.sum

      # What "grit migrate" shells out to, from the same directory.
      - name: Migrate
        working-directory: apps/api
        env:
          DB_DRIVER: postgres
          DATABASE_URL: postgres://grit:grit@localhost:5432/grit?sslmode=disable
        run: go run cmd/migrate/main.go

      - name: Start the API
        working-directory: apps/api
        env:
          DB_DRIVER: postgres
          DATABASE_URL: postgres://grit:grit@localhost:5432/grit?sslmode=disable
          APP_ENV: production
        run: |
          go run cmd/server/main.go &
          for i in $(seq 1 60); do
            curl -sf http://localhost:8080/api/health && break
            sleep 1
          done

      - uses: grafana/setup-k6-action@v1
      - name: Smoke
        run: k6 run tests/k6/smoke.js
      - name: Average load
        run: k6 run tests/k6/average-load.js`}
      />

      <p>
        Smoke first, as its own step. When the average load step fails you
        want the summary to say which one broke, and a green smoke step
        means the script and the deployment are fine and the finding is
        real.
      </p>

      <h2>The health check loop is not optional</h2>
      <p>
        <code>go run</code> returns as soon as the process starts, not when
        it is serving. Without the loop, k6 begins while the API is still
        connecting to Postgres and running migrations, and you get a run of
        connection refused that looks exactly like a system that cannot
        handle two users.
      </p>

      <TipBox tone="warning">
        <p>
          Poll <code>/api/health</code>, not <code>/health</code>. The
          unversioned health route is <code>/api/health</code> and the
          versioned one is <code>/api/v1/health</code>; a bare{' '}
          <code>/health</code> 404s, and <code>curl -sf</code> treats a 404
          as a failure, so the loop runs its full sixty seconds and then
          starts k6 anyway.
        </p>
      </TipBox>

      <h2>Gate on errors, not on latency</h2>
      <p>
        This is the part people get wrong, and the project&apos;s own CI
        already says why about benchmarks: timing a shared runner tells you
        about the runner. A GitHub hosted runner is a noisy neighbour
        environment, and its p95 varies by more between two identical runs
        than most real regressions do.
      </p>

      <CodeBlock
        language="javascript"
        filename="a CI-specific threshold set"
        code={`// In CI, these two are trustworthy and the latency one is not.
thresholds: {
  http_req_failed: ['rate<0.01'],
  checks:          ['rate>0.99'],
  // Kept deliberately loose: it catches an order-of-magnitude regression
  // and ignores the runner being busy.
  http_req_duration: ['p(95)<2000'],
}`}
      />

      <p>
        An error rate is a property of your code. A p95 on a shared runner
        is a property of whoever else is on the machine. Gate the merge on
        the first, and keep the tight latency threshold for the runs you do
        on known hardware.
      </p>

      <h2>Which tests belong in CI</h2>
      <CodeBlock
        language="text"
        code={`smoke.js          every PR          30 seconds
average-load.js   every PR          9 minutes
stress.js         nightly           9 minutes, and it is meant to fail things
spike.js          before a launch   4 minutes
soak.js           before a release  4 hours
breakpoint.js     when capacity planning, by hand`}
      />

      <p>
        Nine minutes on every pull request is a real cost. If that is too
        much, shorten the hold rather than dropping the test: two minutes at
        a hundred users still catches an N+1 that was not there yesterday.
      </p>

      <TipBox tone="info">
        <p>
          Keep the k6 summary as an artifact with{' '}
          <code>--summary-export=summary.json</code> and upload it. A
          regression is much easier to argue about when you can diff this
          week&apos;s numbers against the merge that caused them.
        </p>
      </TipBox>

      <Exercise
        prompt={
          <>
            <p>
              <strong>One.</strong> Add the workflow above to your project
              and push it on a branch. Confirm both k6 steps run green.
            </p>
            <p>
              <strong>Two.</strong> Break it on purpose: add a{' '}
              <code>time.Sleep(3 * time.Second)</code> to a handler the
              journey hits, and push. Confirm the workflow fails and read
              which threshold it names.
            </p>
            <p>
              <strong>Three.</strong> Remove the sleep. Then delete the
              health check loop and push. Note what the failure looks like
              now, and how long it takes to understand it.
            </p>
          </>
        }
        hint={
          <p>
            In step three, read the k6 output rather than the step name. The
            distinction worth seeing is between a request that was answered
            slowly and a request that was never answered.
          </p>
        }
        solution={
          <>
            <p>
              Step two fails on <code>http_req_duration</code> with a p95 in
              the seconds, and the summary points straight at the endpoint.
              That is the regression gate working.
            </p>
            <p>
              Step three fails on <code>http_req_failed</code> at 100% with
              connection refused, which reads as a total outage and is
              actually a race in the workflow. Both are red; only one is
              about your code. Knowing the difference on sight is most of
              what makes a CI load test worth having rather than something
              people start ignoring.
            </p>
          </>
        }
      />

      <KnowledgeCheck
        question="Your CI load test fails about one run in four on p95, and the failures are not reproducible locally. What is the right fix?"
        choices={[
          {
            label: 'Loosen the latency threshold in CI and gate on error rate and checks instead',
            correct: true,
            feedback:
              'A shared runner varies more than most regressions. Gate on what is a property of your code, and keep the tight latency number for known hardware.',
          },
          {
            label: 'Retry the job automatically until it passes',
            feedback:
              'That is the same as deleting the gate, with extra minutes. A test that is retried until green stops being a test.',
          },
          {
            label: 'Reduce the virtual user count until it stops failing',
            feedback:
              'It changes what you measure to make the number convenient, and a smaller run is also noisier per request.',
          },
          {
            label: 'Move the load test out of CI entirely',
            feedback:
              'Then regressions get found by users. The error rate and check gates are reliable on a shared runner; keep those.',
          },
        ]}
      />
    </>
  )
}
