# k6 load tests for commerce

The six load-test types from PHASE 5 §9.2 — each a drop-in entry point
sharing the same user journey from `lib/common.js`. Edit the journey
once; reshape via the per-file stages.

| File              | Question it answers                                    | Run                                |
| ----------------- | ------------------------------------------------------ | ---------------------------------- |
| smoke.js          | Does the script + system handle minimal load?          | `k6 run tests/k6/smoke.js`        |
| average-load.js   | How does it behave under normal expected traffic?      | `k6 run tests/k6/average-load.js` |
| stress.js         | How does it behave beyond normal — at its limits?      | `k6 run tests/k6/stress.js`       |
| spike.js          | Can it survive a sudden, massive surge?                | `k6 run tests/k6/spike.js`        |
| soak.js           | Does it degrade or leak over long periods?             | `k6 run tests/k6/soak.js` (4h)    |
| breakpoint.js     | Exactly where is the breaking point / capacity?        | `k6 run tests/k6/breakpoint.js`   |

## Setup

```bash
# Install k6 (https://grafana.com/docs/k6/latest/set-up/install-k6/)
# macOS:  brew install k6
# Linux:  sudo apt-get install k6
# Windows: winget install k6

# Point the suite at your local server
export BASE_URL=http://localhost:8080

# Run any test
k6 run tests/k6/average-load.js
```

## What thresholds mean

Each test sets `thresholds` from the SLO ladder. If a threshold breaches,
k6 exits non-zero — wire `smoke.js` and `average-load.js` into CI and
the pipeline fails on a performance regression. See `.github/workflows/`
in this repo (if performance-CI is enabled) for the pattern.

## Reading the result

- **Smoke** must always pass. If it fails, the test script is broken,
  not the system.
- **Average-load** establishes your steady-state p95/p99 baseline.
- **Stress** reveals failure mode — graceful degradation vs collapse.
- **Spike** tests recovery: does p95 return to baseline after the surge?
- **Soak** catches memory leaks and unclosed connections that 5-min
  tests miss. Watch `/pulse/ui/` (or your APM) during the 4h run.
- **Breakpoint** pins exact capacity. The VU count when error-rate
  climbs or p95 breaches your SLO is the real number for capacity
  planning.

Every test pairs naturally with Pulse (Grit's observability dashboard)
or any external APM — open the dashboard in another tab while the load
runs and watch the bottleneck appear in real time.
