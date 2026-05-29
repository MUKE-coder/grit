package scaffold

// This file scaffolds the testing + supply-chain hygiene assets called
// out in the security & testing course (PHASE 5):
//
//   - tests/k6/{smoke,average-load,stress,spike,soak,breakpoint}.js
//     A ready-to-run k6 suite covering all six load-test types so the
//     project ships with a baseline + breaking-point answer in an
//     afternoon, not a sprint.
//   - .github/workflows/security.yml — govulncheck (Go) + npm audit
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
		filepath.Join(root, "tests", "k6", "README.md"):           k6ReadmeMD(opts),
		filepath.Join(root, "tests", "k6", "lib", "common.js"):    k6CommonJS(),
		filepath.Join(root, "tests", "k6", "smoke.js"):            k6SmokeJS(),
		filepath.Join(root, "tests", "k6", "average-load.js"):     k6AverageLoadJS(),
		filepath.Join(root, "tests", "k6", "stress.js"):           k6StressJS(),
		filepath.Join(root, "tests", "k6", "spike.js"):            k6SpikeJS(),
		filepath.Join(root, "tests", "k6", "soak.js"):             k6SoakJS(),
		filepath.Join(root, "tests", "k6", "breakpoint.js"):       k6BreakpointJS(),
		filepath.Join(root, ".github", "dependabot.yml"):          dependabotYAML(opts),
		filepath.Join(root, ".github", "workflows", "security.yml"): securityCIYAML(),
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

The six load-test types from PHASE 5 §9.2 — each a drop-in entry point
sharing the same user journey from ` + "`lib/common.js`" + `. Edit the journey
once; reshape via the per-file stages.

| File              | Question it answers                                    | Run                                |
| ----------------- | ------------------------------------------------------ | ---------------------------------- |
| smoke.js          | Does the script + system handle minimal load?          | ` + "`k6 run tests/k6/smoke.js`" + `        |
| average-load.js   | How does it behave under normal expected traffic?      | ` + "`k6 run tests/k6/average-load.js`" + ` |
| stress.js         | How does it behave beyond normal — at its limits?      | ` + "`k6 run tests/k6/stress.js`" + `       |
| spike.js          | Can it survive a sudden, massive surge?                | ` + "`k6 run tests/k6/spike.js`" + `        |
| soak.js           | Does it degrade or leak over long periods?             | ` + "`k6 run tests/k6/soak.js`" + ` (4h)    |
| breakpoint.js     | Exactly where is the breaking point / capacity?        | ` + "`k6 run tests/k6/breakpoint.js`" + `   |

## Setup

` + "```bash" + `
# Install k6 (https://grafana.com/docs/k6/latest/set-up/install-k6/)
# macOS:  brew install k6
# Linux:  sudo apt-get install k6
# Windows: winget install k6

# Point the suite at your local server
export BASE_URL=http://localhost:8080

# Run any test
k6 run tests/k6/average-load.js
` + "```" + `

## What thresholds mean

Each test sets ` + "`thresholds`" + ` from the SLO ladder. If a threshold breaches,
k6 exits non-zero — wire ` + "`smoke.js`" + ` and ` + "`average-load.js`" + ` into CI and
the pipeline fails on a performance regression. See ` + "`.github/workflows/`" + `
in this repo (if performance-CI is enabled) for the pattern.

## Reading the result

- **Smoke** must always pass. If it fails, the test script is broken,
  not the system.
- **Average-load** establishes your steady-state p95/p99 baseline.
- **Stress** reveals failure mode — graceful degradation vs collapse.
- **Spike** tests recovery: does p95 return to baseline after the surge?
- **Soak** catches memory leaks and unclosed connections that 5-min
  tests miss. Watch ` + "`/pulse/ui/`" + ` (or your APM) during the 4h run.
- **Breakpoint** pins exact capacity. The VU count when error-rate
  climbs or p95 breaches your SLO is the real number for capacity
  planning.

Every test pairs naturally with Pulse (Grit's observability dashboard)
or any external APM — open the dashboard in another tab while the load
runs and watch the bottleneck appear in real time.
`
}

func k6CommonJS() string {
	return `// Shared user journey + helpers for all k6 tests. Editing this single
// file changes the behaviour of every test type — the per-test files
// only differ in their stages (load shape).

import http from 'k6/http'
import { check, sleep } from 'k6'

export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080'

// SLO thresholds — adjust these to match your contractual targets.
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

// userJourney is one virtual user's loop. Realistic load needs varied
// flows + think time — synthetic hammering finds fake bottlenecks
// (PHASE 5 §9.3 "Model realistic traffic").
export function userJourney() {
  // 1. Public health check (cheap)
  const health = http.get(BASE_URL + '/health')
  check(health, {
    'health: status 200': (r) => r.status === 200,
  })

  // 2. List blogs (typical read-heavy page)
  const blogs = http.get(BASE_URL + '/api/blogs')
  check(blogs, {
    'blogs: status 200':   (r) => r.status === 200,
    'blogs: body non-empty': (r) => r.body && r.body.length > 0,
  })

  // Real users pause between actions. Without this, you measure how
  // fast k6 can fire requests, not how the system handles users.
  sleep(1 + Math.random() * 2) // 1–3s think time
}
`
}

func k6SmokeJS() string {
	return `// Smoke test — does the script work and does the system handle
// minimal load? PHASE 5 §9.2: run this FIRST before every other type.
import { userJourney } from './lib/common.js'
import { defaultThresholds } from './lib/common.js'

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

// dependabotYAML produces the Dependabot config that closes OWASP A03
// Supply-Chain Failures — automated PRs for vulnerable Go modules, npm
// packages, and GitHub Actions updates.
func dependabotYAML(opts Options) string {
	return `# Dependabot — automated dependency updates.
# Closes OWASP Top 10:2025 A03 Supply-Chain Failures by surfacing
# vulnerable dependencies as PRs you can review and merge.

version: 2

updates:
  # Go modules
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    labels:
      - dependencies
      - go

  # Frontend (Vite single-app layout; adjust if your tree differs)
  - package-ecosystem: npm
    directory: "/frontend"
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    labels:
      - dependencies
      - javascript

  # GitHub Actions workflows
  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: weekly
    labels:
      - dependencies
      - ci
`
}

func securityCIYAML() string {
	return `name: security

# Security scans that should run on every PR and on a weekly schedule.
# Covers OWASP Top 10:2025 A03 (Supply Chain) and A02 (Misconfiguration)
# via static analysis. Pair with Dependabot for full coverage.

on:
  push:
    branches: [main, master]
  pull_request:
    branches: [main, master]
  schedule:
    # Weekly Monday 06:00 UTC — catches newly-disclosed CVEs in
    # dependencies that didn't change in your code.
    - cron: '0 6 * * 1'

jobs:
  govulncheck:
    name: Go vulnerability scan (govulncheck)
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      - name: Run govulncheck
        run: govulncheck ./...

  npm-audit:
    name: Frontend npm audit
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
      - name: Install pnpm
        run: npm install -g pnpm@9.15.0
      - name: Install deps
        working-directory: frontend
        run: pnpm install --frozen-lockfile
      # Fail only on high/critical so noisy mediums don't block PRs.
      - name: Run npm audit (high+)
        working-directory: frontend
        run: pnpm audit --audit-level=high

  codeql:
    name: CodeQL static analysis
    runs-on: ubuntu-latest
    permissions:
      actions: read
      contents: read
      security-events: write
    strategy:
      fail-fast: false
      matrix:
        language: ['go', 'javascript']
    steps:
      - uses: actions/checkout@v4
      - uses: github/codeql-action/init@v3
        with:
          languages: ${{ matrix.language }}
      - uses: github/codeql-action/autobuild@v3
      - uses: github/codeql-action/analyze@v3
        with:
          category: "/language:${{ matrix.language }}"
`
}
