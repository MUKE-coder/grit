import http from 'k6/http'
import { check } from 'k6'
import { Trend, Rate } from 'k6/metrics'
import { SharedArray } from 'k6/data'

/*
 * One script, four scenarios, pointed at whichever base URL you pass:
 *
 *   k6 run -e BASE=http://localhost:8091 -e SCENARIO=list  k6/bench.js
 *   k6 run -e BASE=http://localhost:8092 -e SCENARIO=list  k6/bench.js
 *
 * Same script both sides on purpose. Two scripts drift, and then you are
 * comparing the scripts.
 *
 * Scenarios:
 *   list   GET /products?page=N        paginated read, the common case
 *   show   GET /products/:id           single-row read by primary key
 *   write  POST /products              insert with validation
 *   mixed  85% read / 10% show / 5% write, closer to a real workload
 */

const BASE = __ENV.BASE
const SCENARIO = __ENV.SCENARIO || 'list'
const VUS = parseInt(__ENV.VUS || '50', 10)
const DURATION = __ENV.DURATION || '30s'

if (!BASE) throw new Error('set -e BASE=http://host:port')

// Loaded once and shared across VUs rather than per-VU, so the ids are not
// re-parsed 50 times and memory stays flat.
const IDS = new SharedArray('ids', () => JSON.parse(open('../seed/ids.json')))

const ttfb = new Trend('ttfb', true)
const failures = new Rate('failed_requests')

export const options = {
  scenarios: {
    [SCENARIO]: {
      executor: 'constant-vus',
      vus: VUS,
      duration: DURATION,
      exec: SCENARIO,
      // Long enough for in-flight requests to land; without it a slow tail
      // gets counted as an error at the cut-off.
      gracefulStop: '10s',
    },
  },
  // Discard the connection-reuse warm-up rather than letting it widen p99.
  discardResponseBodies: false,
  thresholds: {
    failed_requests: ['rate<0.01'],
  },
}

const JSON_HEADERS = { 'Content-Type': 'application/json', Accept: 'application/json' }

function record(res, expected) {
  ttfb.add(res.timings.waiting)
  const ok = check(res, {
    'status ok': (r) => r.status === expected,
    'has data': (r) => r.body && r.body.length > 2,
  })
  failures.add(!ok)
  return ok
}

export function list() {
  // Walk pages rather than hammering page 1, so neither side gets to serve the
  // same buffer-cached page for the whole run.
  const page = (__ITER % 50) + 1
  const res = http.get(`${BASE}/api/v1/products?page=${page}&page_size=20`, { headers: JSON_HEADERS })
  record(res, 200)
}

export function show() {
  const id = IDS[__ITER % IDS.length]
  const res = http.get(`${BASE}/api/v1/products/${id}`, { headers: JSON_HEADERS })
  record(res, 200)
}

export function write() {
  // __VU and __ITER keep the SKU unique without a shared counter, which would
  // serialise the VUs and measure the counter instead of the framework.
  const n = `${__VU}-${__ITER}`
  const res = http.post(
    `${BASE}/api/v1/products`,
    JSON.stringify({
      name: `Bench Product ${n}`,
      sku: `BENCH-${n}`,
      description: 'Created by the k6 write scenario.',
      price: 19.99,
      stock: 42,
      active: true,
    }),
    { headers: JSON_HEADERS },
  )
  record(res, 201)
}

export function mixed() {
  const roll = __ITER % 20
  if (roll < 17) return list()
  if (roll < 19) return show()
  return write()
}
