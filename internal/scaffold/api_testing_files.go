package scaffold

// This file scaffolds the testing + supply-chain hygiene assets called
// out in the security & testing course (PHASE 5):
//
//   - tests/k6/{smoke,average-load,stress,spike,soak,breakpoint}.js
//     A ready-to-run k6 suite covering all six load-test types so the
//     project ships with a baseline + breaking-point answer in an
//     afternoon, not a sprint.
//   - .github/workflows/security.yml: govulncheck, pnpm audit and CodeQL
//   - .github/workflows/ci.yml: vet and tests, frontend type-check, tests and build
//     (frontend) + CodeQL.
//   - .github/dependabot.yml — Go modules + npm version + GitHub
//     Actions, weekly schedule. Closes OWASP A03 Supply-Chain Failures.

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeTestingFiles(root string, opts Options) error {
	files := map[string]string{
		filepath.Join(root, "tests", "k6", "README.md"):             k6ReadmeMD(opts),
		filepath.Join(root, "tests", "k6", "lib", "common.js"):      k6CommonJS(),
		filepath.Join(root, "tests", "k6", "smoke.js"):              k6SmokeJS(),
		filepath.Join(root, "tests", "k6", "average-load.js"):       k6AverageLoadJS(),
		filepath.Join(root, "tests", "k6", "stress.js"):             k6StressJS(),
		filepath.Join(root, "tests", "k6", "spike.js"):              k6SpikeJS(),
		filepath.Join(root, "tests", "k6", "soak.js"):               k6SoakJS(),
		filepath.Join(root, "tests", "k6", "breakpoint.js"):         k6BreakpointJS(),
		filepath.Join(root, ".github", "dependabot.yml"):            dependabotYAML(opts),
		filepath.Join(root, ".github", "workflows", "security.yml"): securityCIYAML(opts),
		filepath.Join(root, ".github", "govulncheck-allow.txt"):     govulncheckAllowTXT(),
		filepath.Join(root, ".github", "workflows", "ci.yml"):       ciYAML(opts),
		// Lint config lives beside the Go module it applies to — that is the
		// repo root for --single, apps/api otherwise.
		filepath.Join(opts.APIRoot(root), ".golangci.yml"):         golangciYAML(),
		filepath.Join(root, ".github", "workflows", "lint.yml"):    lintCIYAML(opts),
		filepath.Join(root, ".github", "workflows", "release.yml"): releaseCIYAML(opts),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{PROJECT}}", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func k6ReadmeMD(opts Options) string {
	return `# k6 load tests for ` + opts.ProjectName + `

Six load-test types, each a drop-in entry point sharing one user journey
from ` + "`lib/common.js`" + `. Edit the journey once; reshape the load with the
per-file stages.

| File              | Question it answers                                    | Run                                |
| ----------------- | ------------------------------------------------------ | ---------------------------------- |
| smoke.js          | Does the script + system handle minimal load?          | ` + "`k6 run tests/k6/smoke.js`" + `        |
| average-load.js   | How does it behave under normal expected traffic?      | ` + "`k6 run tests/k6/average-load.js`" + ` |
| stress.js         | How does it behave beyond normal, at its limits?       | ` + "`k6 run tests/k6/stress.js`" + `       |
| spike.js          | Can it survive a sudden, massive surge?                | ` + "`k6 run tests/k6/spike.js`" + `        |
| soak.js           | Does it degrade or leak over long periods?             | ` + "`k6 run tests/k6/soak.js`" + ` (4h)    |
| breakpoint.js     | Exactly where is the breaking point / capacity?        | ` + "`k6 run tests/k6/breakpoint.js`" + `   |

## Setup

` + "```bash" + `
# Install k6 (https://grafana.com/docs/k6/latest/set-up/install-k6/)
# macOS:  brew install k6
# Linux:  sudo apt-get install k6
# Windows: winget install k6

# The API has to be running and migrated: the journey signs in.
grit migrate
k6 run tests/k6/smoke.js
` + "```" + `

## Point it at something worth measuring

Out of the box the journey reads ` + "`/api/health`" + ` and ` + "`/api/profile`" + `. Both
are nearly free, so a run of only those two reports a system faster than
yours is. Name one of your own resources and it reads that too:

` + "```bash" + `
k6 run -e RESOURCE=products tests/k6/average-load.js
` + "```" + `

Three environment variables, all optional:

| Variable                           | Default                  | What it does                              |
| ---------------------------------- | ------------------------ | ----------------------------------------- |
| ` + "`BASE_URL`" + `                       | ` + "`http://localhost:8080`" + ` | Which API to hit                   |
| ` + "`RESOURCE`" + `                       | none                     | Also read ` + "`/api/<resource>`" + ` each iteration |
| ` + "`LOADTEST_EMAIL`" + ` / ` + "`LOADTEST_PASSWORD`" + ` | a throwaway user  | Sign in as an existing account       |

The journey signs in once in ` + "`setup()`" + ` and shares one bearer token across
every virtual user. Use ` + "`LOADTEST_EMAIL`" + ` when the endpoint you care about
needs a role: a resource generated by ` + "`grit generate resource`" + ` wants ADMIN
or its own permission, and a throwaway USER gets 403s that look like
failures but are really a test pointed at the wrong account.

## What thresholds mean

Each test sets ` + "`thresholds`" + ` from the SLO ladder. If a threshold breaches,
k6 exits non-zero, so wiring ` + "`smoke.js`" + ` and ` + "`average-load.js`" + ` into CI
fails the pipeline on a performance regression.

## Reading the result

- **Smoke** must always pass. If it fails, the test script is broken,
  not the system.
- **Average-load** establishes your steady-state p95/p99 baseline.
- **Stress** reveals the failure mode: graceful degradation or collapse.
- **Spike** tests recovery: does p95 return to baseline after the surge?
- **Soak** catches memory leaks and unclosed connections that 5-minute
  tests miss. Watch ` + "`/pulse/ui/`" + ` (or your APM) during the 4h run.
- **Breakpoint** pins exact capacity. The VU count when the error rate
  climbs or p95 breaches your SLO is the real number for capacity
  planning.

Every test pairs with Pulse, Grit's observability dashboard. k6 tells you
p95 moved; it cannot tell you which function. Open ` + "`/pulse/ui/`" + ` in another
tab while the load runs and the bottleneck appears in real time.
`
}

func k6CommonJS() string {
	return `// Shared user journey + helpers for all k6 tests. Editing this single
// file changes the behaviour of every test type: the per-test files
// only differ in their stages (load shape).

import http from 'k6/http'
import { check, sleep } from 'k6'

export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

// SLO thresholds. Adjust these to match your contractual targets.
// They're the gate that turns k6 into a CI regression check.
export const defaultThresholds = {
  http_req_duration: ['p(95)<500', 'p(99)<1000'],
  http_req_failed:   ['rate<0.01'],
  checks:            ['rate>0.99'],
}

// Stress / spike / breakpoint tests intentionally cross thresholds.
// They use looser numbers so the test still finishes (exit code 0)
// when the system is pushed to its real limits.
export const aggressiveThresholds = {
  http_req_duration: ['p(95)<2000'],
  http_req_failed:   ['rate<0.10'],
}

// setup() runs once before the load and hands every virtual user one bearer
// token. A bearer beats the cookie flow here: k6 gives each VU its own cookie
// jar, and the sign-in cookies plus CSRF add a round trip per VU that you did
// not mean to measure.
//
// LOADTEST_EMAIL / LOADTEST_PASSWORD sign in as an account you already have,
// which is how you load-test a route that needs a role: a resource endpoint
// from "grit generate resource" wants ADMIN or its own permission. Without
// them a throwaway USER is registered.
export function setup() {
  const json = { headers: { 'Content-Type': 'application/json' } }
  const existing = Boolean(__ENV.LOADTEST_EMAIL && __ENV.LOADTEST_PASSWORD)
  const email = __ENV.LOADTEST_EMAIL || 'loadtest+' + Date.now() + '@loadtest.local'
  const password = __ENV.LOADTEST_PASSWORD || 'Password123!'

  if (!existing) {
    http.post(BASE_URL + '/api/auth/register', JSON.stringify({
      first_name: 'Load', last_name: 'Test', email: email, password: password,
    }), json)
  }

  const res = http.post(BASE_URL + '/api/auth/login',
    JSON.stringify({ email: email, password: password }), json)
  const token = res.json('data.tokens.access_token')
  if (!token) {
    throw new Error('setup: no access token, login returned ' + res.status +
      '. Is the API running on ' + BASE_URL + ', and has grit migrate been run?')
  }
  return { token: token }
}

// userJourney is one virtual user's loop. Realistic load needs varied
// flows + think time: synthetic hammering finds fake bottlenecks.
export function userJourney(data) {
  // 1. Health check, unauthenticated. The cheapest liveness path, and the
  //    one that tells you the API is still answering under load.
  const health = http.get(BASE_URL + '/api/health')
  check(health, {
    'health: status 200': (r) => r.status === 200,
  })

  // 2. An authenticated read. This is the first request that costs
  //    anything: the auth middleware runs and a row is read.
  const auth = { headers: { Authorization: 'Bearer ' + data.token } }
  const profile = http.get(BASE_URL + '/api/profile', auth)
  check(profile, {
    'profile: status 200': (r) => r.status === 200,
  })

  // 3. Your own resource, named rather than hardcoded:
  //
  //      k6 run -e RESOURCE=products tests/k6/smoke.js
  //
  //    This is the request worth load-testing. The two above are nearly
  //    free, and a journey made of free requests reports a system that is
  //    faster than it is. A hardcoded name would be worse: this file
  //    shipped reading /api/blogs, so every test failed on every project
  //    that had no Blog, and checked for a non-empty body, which an error
  //    response has.
  if (__ENV.RESOURCE) {
    const list = http.get(
      BASE_URL + '/api/' + __ENV.RESOURCE + '?page=1&page_size=20', auth)
    check(list, {
      'resource: status 200': (r) => r.status === 200,
      'resource: returned a data array': (r) => Array.isArray(r.json('data')),
    })
  }

  // Real users pause between actions. Without this, you measure how
  // fast k6 can fire requests, not how the system handles users.
  sleep(1 + Math.random() * 2) // 1 to 3s think time
}
`
}

func k6SmokeJS() string {
	return `// Smoke test — does the script work and does the system handle
// minimal load? PHASE 5 §9.2: run this FIRST before every other type.
import { userJourney, defaultThresholds } from './lib/common.js'

// k6 only runs setup() when the entry script exports it.
export { setup } from './lib/common.js'

export const options = {
  vus: 2,
  duration: '30s',
  thresholds: defaultThresholds,
}

export default userJourney
`
}

func k6AverageLoadJS() string {
	return `// Average-load — how does it behave under normal expected traffic?
// Ramps to 100 VUs, holds 5 minutes. The steady-state baseline you
// compare every future test against.
import { userJourney, defaultThresholds } from './lib/common.js'
export { setup } from './lib/common.js'

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // ramp up
    { duration: '5m', target: 100 },  // hold (steady state)
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: defaultThresholds,
}

export default userJourney
`
}

func k6StressJS() string {
	return `// Stress — how does it behave well beyond normal?
// At 4× normal load you're learning the failure mode, not just the
// limit. Loose thresholds so the run completes even when the system
// degrades under load.
import { userJourney, aggressiveThresholds } from './lib/common.js'
export { setup } from './lib/common.js'

export const options = {
  stages: [
    { duration: '2m', target: 100 },  // normal
    { duration: '5m', target: 400 },  // 4× normal — the stress
    { duration: '2m', target: 0 },    // ramp down
  ],
  thresholds: aggressiveThresholds,
}

export default userJourney
`
}

func k6SpikeJS() string {
	return `// Spike — can it survive a sudden, massive surge (flash sale, viral
// post)? The KEY question is RECOVERY: does p95 return to baseline
// after the surge passes?
import { userJourney, aggressiveThresholds } from './lib/common.js'
export { setup } from './lib/common.js'

export const options = {
  stages: [
    { duration: '1m',  target: 50 },    // baseline
    { duration: '10s', target: 1000 },  // SLAM
    { duration: '1m',  target: 1000 },  // hold the spike
    { duration: '10s', target: 50 },    // drop back
    { duration: '1m',  target: 50 },    // recovery window — watch p95
    { duration: '10s', target: 0 },
  ],
  thresholds: aggressiveThresholds,
}

export default userJourney
`
}

func k6SoakJS() string {
	return `// Soak — does the system leak over hours? A flat latency line that
// slowly tilts upward is the classic memory-leak signature. Run before
// every major release; watch /pulse/ui in another tab.
//
// NOTE: this takes ~4 hours. Schedule it overnight or on a CI runner.
import { userJourney, defaultThresholds } from './lib/common.js'
export { setup } from './lib/common.js'

export const options = {
  stages: [
    { duration: '5m', target: 100 },  // ramp up
    { duration: '4h', target: 100 },  // long hold — the soak
    { duration: '5m', target: 0 },    // ramp down
  ],
  thresholds: defaultThresholds,
}

export default userJourney
`
}

func k6BreakpointJS() string {
	return `// Breakpoint — slowly ramp until the system actually breaks.
// The VU count at which error-rate or p95 first breaches the threshold
// is your REAL capacity number for launch planning.
import { userJourney } from './lib/common.js'
export { setup } from './lib/common.js'

export const options = {
  // Slow climb to 5000 VUs over 1h. Adjust the ceiling for your gear.
  stages: [
    { duration: '1h', target: 5000 },
  ],
  thresholds: {
    // Abort the test the moment thresholds breach. The VU count when
    // that happens is the breaking point.
    http_req_duration: [{ threshold: 'p(95)<1000', abortOnFail: true }],
    http_req_failed:   [{ threshold: 'rate<0.05',  abortOnFail: true }],
  },
}

export default userJourney
`
}
