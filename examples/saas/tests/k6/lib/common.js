// Shared user journey + helpers for all k6 tests. Editing this single
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
