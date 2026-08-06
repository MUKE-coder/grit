/*
 * Every benchmark number on this site lives here and nowhere else.
 *
 * The homepage chart, the overview page and the per-framework guides all read
 * from this file, so there is one place to update when the numbers are re-run
 * and no way for the homepage to drift from the methodology page behind it.
 *
 * ── Why these are pairs, and why ratios ───────────────────────────────────
 *
 * Each framework was measured head to head against Grit in its own run, minutes
 * apart, three repetitions, medians reported, zero failed requests.
 *
 * That structure is not decoration. Across runs this machine drifts a long way:
 * Grit's single-row read measured 4,536 req/s in the Bun pair and 8,509 in the
 * Express pair — from an identical binary, in the same sitting, because the
 * machine drifts as write scenarios accumulate in Postgres. So the absolute
 * figures are only meaningful next to the Grit baseline measured beside them,
 * and the RATIO is what survives comparison. Putting Bun's 2,717 on the same
 * axis as Encore's 663 would imply a shared baseline that does not exist.
 *
 * Filled in by hand from `python pair-report.py`. A build step that silently
 * rewrote published performance claims would be worse than an edit you can see
 * in a diff.
 */

export type ScenarioId = 'show' | 'write' | 'list' | 'mixed'

export interface ScenarioResult {
  /** the opponent's median req/s in this pair */
  rps: number
  /** Grit's median req/s measured in the SAME pair */
  gritRps: number
  /** opponent median latency, formatted */
  median: string
  gritMedian: string
  /** app container CPU during the run, % of a 400% allowance */
  appCpu: number
  /** postgres CPU during the same run, % of an 800% allowance */
  dbCpu: number
  gritAppCpu: number
  gritDbCpu: number
  /**
   * Set where repeated measurement could not separate the two frameworks.
   *
   * Both sides queue behind the same saturated Postgres in these rows, so what
   * varies between runs is database state rather than framework speed. Bun's
   * mixed scenario was measured four times and came out 0.91x, 1.21x, 1.19x and
   * 0.91x; across seven repetitions Grit ranged 329 to 643 req/s and Bun 445 to
   * 694. Publishing any single one of those as a result would be picking a
   * number, so the ratio is shown but not claimed.
   */
  inconclusive?: boolean
}

export interface Framework {
  slug: string
  name: string
  version: string
  /** the ORM it uses — every framework here uses one, none use raw SQL */
  orm: string
  tagline: string
  stack: string
  /** under /public/logos — the official mark, used on the chart and the guides */
  logo: string
  /** the framework's own brand colour, so its bar is recognisable at a glance */
  color: string
  /**
   * Some marks are near-black and vanish on a dark background. Setting this
   * inverts the image in dark mode rather than shipping a recoloured logo,
   * which would be someone else's trademark redrawn by us.
   */
  invertOnDark?: boolean
  /**
   * Set while a pair still carries figures measured before v3.134.0, which cut
   * a generated write from seven statements to one. Those runs understate Grit
   * and say so on the page until they are re-run. Delete the flag with the
   * numbers, never before.
   */
  preV3134?: boolean
  results: Record<ScenarioId, ScenarioResult>
}

/** Grit's own side of every pair. */
export const GRIT = {
  name: 'Grit',
  logo: '/logos/grit.png',
  color: '#3BB4F5',
  stack: 'Go 1.22, Gin, GORM',
}

export const SCENARIOS: {
  id: ScenarioId
  label: string
  title: string
  subtitle: string
  request: string
  exercises: string
}[] = [
  {
    id: 'show',
    label: 'Read by ID',
    title: 'GET /products/:id',
    subtitle: 'One indexed lookup and a JSON encode: the cleanest read of framework overhead',
    request: 'GET /api/v1/products/:id',
    exercises: 'routing, a single primary-key lookup, JSON encoding',
  },
  {
    id: 'write',
    label: 'Insert',
    title: 'POST /products',
    subtitle: 'Body parse, validation, one INSERT through the ORM',
    request: 'POST /api/v1/products',
    exercises: 'body parsing, validation, one INSERT',
  },
  {
    id: 'list',
    label: 'Paginated list',
    title: 'GET /products?page=N',
    subtitle: 'Twenty rows hydrated and encoded, plus a COUNT over the table',
    request: 'GET /api/v1/products?page=N&page_size=20',
    exercises: 'routing, hydrating 20 rows, JSON encoding, and a COUNT(*)',
  },
  {
    id: 'mixed',
    label: 'Mixed',
    title: '85% list / 10% read / 5% write',
    subtitle: 'Closer to the shape of a real read-heavy API',
    request: 'a weighted blend of the three above',
    exercises: 'everything above, in the proportions a real API sees',
  },
]

/*
 * Measured 4 August 2026. 50 concurrent users, 30-second runs, 4 CPUs and 2 GB
 * per app container, one shared Postgres on 8 CPUs, the same 10,000 rows
 * restored before every single run, k6 running inside the container network.
 */
export const FRAMEWORKS: Framework[] = [
  {
    slug: 'bun',
    logo: '/logos/bun.png',
    color: '#FBF0DF',
    name: 'Bun',
    version: 'v1.3 + Drizzle',
    orm: 'Drizzle',
    tagline: 'Bun.serve with reusePort workers, Drizzle over Bun’s native SQL client',
    /*
     * Re-measured on v3.134.0, which cut a generated write from seven statements
     * to one. The write scenario reversed: Grit was losing it 0.81x and now wins
     * it 2.23x, going from 2,686 to 4,959 req/s on a machine that was slower
     * overall (its single-row read fell from 10,655 to 4,536 across the same
     * two runs).
     *
     * What changed: an audit row written for unauthenticated requests, a SELECT
     * re-reading the row that had just been inserted, and a transaction wrapped
     * around a single INSERT that Postgres already makes atomic.
     *
     * Note what saturated. On both read and write Bun pins its container near
     * 400% of a 400% allowance while Grit sits at 294-311%, so Bun is at its
     * ceiling there and Grit is not. On list and mixed neither is: Postgres runs
     * at 764-825% and both sides queue behind it, which is why those two ratios
     * sit near 1 and why mixed can land either side of it between runs.
     */
    stack: 'Bun 1.3, Bun.serve, Drizzle ORM',
    results: {
      show:  { rps: 2717, gritRps: 4536, median: '16.2 ms', gritMedian: '8.2 ms',  appCpu: 399, dbCpu: 153, gritAppCpu: 294, gritDbCpu: 249 },
      write: { rps: 2224, gritRps: 4959, median: '18.0 ms', gritMedian: '7.7 ms',  appCpu: 412, dbCpu: 195, gritAppCpu: 311, gritDbCpu: 249 },
      list:  { rps: 590,  gritRps: 615,  median: '70.8 ms', gritMedian: '65.6 ms', appCpu: 267, dbCpu: 798, gritAppCpu: 125, gritDbCpu: 803, inconclusive: true },
      mixed: { rps: 621,  gritRps: 568,  median: '65.8 ms', gritMedian: '70.9 ms', appCpu: 327, dbCpu: 825, gritAppCpu: 125, gritDbCpu: 764, inconclusive: true },
    },
  },
  {
    slug: 'encore',
    logo: '/logos/encore.png',
    color: '#EEEBE0',
    name: 'Encore.ts',
    version: 'v1.57 + Drizzle',
    orm: 'Drizzle',
    tagline: 'Rust HTTP runtime, compiled by Encore’s own CLI, Drizzle over its SQLDatabase',
    /*
     * Worth stating plainly: Encore never saturated anything. Its container sat
     * at 142-195% of a 400% allowance and Postgres at 26-176%, so neither was
     * the limit. Something else bounded it — most likely the Drizzle /
     * node-postgres path rather than the Rust HTTP layer, which is the part
     * Encore is fast at. Its figures here are a floor for Encore, and a tuned
     * Encore setup would very likely do better.
     */
    stack: 'Encore.ts 1.57, Rust runtime, Drizzle ORM',
    results: {
      show:  { rps: 663, gritRps: 6646, median: '71.9 ms',  gritMedian: '5.3 ms',  appCpu: 142, dbCpu: 26,  gritAppCpu: 300, gritDbCpu: 249 },
      write: { rps: 854, gritRps: 4345, median: '51.9 ms',  gritMedian: '8.7 ms',  appCpu: 195, dbCpu: 67,  gritAppCpu: 312, gritDbCpu: 247 },
      list:  { rps: 249, gritRps: 617,  median: '166.0 ms', gritMedian: '67.0 ms', appCpu: 147, dbCpu: 176, gritAppCpu: 170, gritDbCpu: 1084 },
      mixed: { rps: 363, gritRps: 834,  median: '138.0 ms', gritMedian: '50.7 ms', appCpu: 148, dbCpu: 175, gritAppCpu: 112, gritDbCpu: 820 },
    },
  },
  {
    slug: 'express',
    logo: '/logos/express.png',
    color: '#5FA04E',
    name: 'Express',
    version: 'v5 + Prisma',
    orm: 'Prisma',
    tagline: 'Express 5 on Node 22, one cluster worker per CPU, Prisma ORM',
    stack: 'Node 22, Express 5, cluster, Prisma',
    results: {
      show:  { rps: 982, gritRps: 8509, median: '55.1 ms',  gritMedian: '4.5 ms',  appCpu: 420, dbCpu: 45,  gritAppCpu: 288, gritDbCpu: 255 },
      write: { rps: 748, gritRps: 5126, median: '72.8 ms',  gritMedian: '7.5 ms',  appCpu: 416, dbCpu: 46,  gritAppCpu: 288, gritDbCpu: 246 },
      list:  { rps: 280, gritRps: 861,  median: '168.7 ms', gritMedian: '46.8 ms', appCpu: 434, dbCpu: 311, gritAppCpu: 124, gritDbCpu: 820 },
      mixed: { rps: 297, gritRps: 763,  median: '151.5 ms', gritMedian: '51.5 ms', appCpu: 425, dbCpu: 303, gritAppCpu: 114, gritDbCpu: 816 },
    },
  },
  {
    slug: 'nextjs',
    logo: '/logos/nextjs.png',
    color: '#ffffff',
    // Near-black on transparent, so it needs a light backing on this theme.
    invertOnDark: true,
    name: 'Next.js',
    version: 'v15.5 + Prisma',
    orm: 'Prisma',
    tagline: 'App Router route handlers, standalone production build, Prisma over Postgres',
    /*
     * Measured a day after the other pairs, on a machine that had drifted: Grit's
     * single-row read came in at 1,911 here against 4,536 in the Bun pair, from
     * the same binary. Next.js saturated its own container in all four scenarios
     * (406-412% of a 400% allowance), so its figures are genuine ceilings.
     */
    stack: 'Node 22, Next.js 15.5 App Router, standalone output, Prisma',
    results: {
      show:  { rps: 433, gritRps: 6822, median: '103.0 ms', gritMedian: '5.5 ms',  appCpu: 403, dbCpu: 25,  gritAppCpu: 287, gritDbCpu: 246 },
      write: { rps: 499, gritRps: 6157, median: '93.0 ms',  gritMedian: '6.2 ms',  appCpu: 400, dbCpu: 24,  gritAppCpu: 310, gritDbCpu: 260 },
      list:  { rps: 247, gritRps: 1186, median: '189.7 ms', gritMedian: '35.6 ms', appCpu: 405, dbCpu: 206, gritAppCpu: 123, gritDbCpu: 818 },
      mixed: { rps: 271, gritRps: 1053, median: '181.0 ms', gritMedian: '38.5 ms', appCpu: 411, dbCpu: 167, gritAppCpu: 108, gritDbCpu: 799 },
    },
  },
  {
    slug: 'django',
    logo: '/logos/django.svg',
    color: '#44B78B',
    name: 'Django',
    version: 'v5.1 + Django ORM',
    orm: 'Django ORM',
    tagline: 'gunicorn with gevent workers, DRF-shaped views, Django ORM with connection pooling',
    /*
     * Django saturated its own container in all four scenarios (404-426% of a
     * 400% allowance) while Postgres sat at 50-259%, so every figure here is a
     * genuine ceiling for Django rather than a database limit.
     */
    stack: 'Python 3.12, Django 5.1, gunicorn + gevent (9 workers), psycopg pool',
    results: {
      show:  { rps: 811, gritRps: 5983, median: '73.4 ms',  gritMedian: '6.0 ms',  appCpu: 400, dbCpu: 47,  gritAppCpu: 251, gritDbCpu: 211 },
      write: { rps: 913, gritRps: 5901, median: '65.1 ms',  gritMedian: '6.2 ms',  appCpu: 405, dbCpu: 42,  gritAppCpu: 290, gritDbCpu: 242 },
      list:  { rps: 379, gritRps: 1094, median: '118.5 ms', gritMedian: '36.9 ms', appCpu: 408, dbCpu: 290, gritAppCpu: 137, gritDbCpu: 815 },
      mixed: { rps: 481, gritRps: 1163, median: '77.4 ms',  gritMedian: '36.2 ms', appCpu: 398, dbCpu: 289, gritAppCpu: 111, gritDbCpu: 805 },
    },
  },
  {
    slug: 'laravel',
    logo: '/logos/laravel.png',
    color: '#FF2D20',
    name: 'Laravel',
    version: 'v13 + Eloquent',
    orm: 'Eloquent',
    tagline: 'nginx + php-fpm, opcache with a tracing JIT, production-only vendor, Eloquent',
    /*
     * Re-measured after three setup faults were found and fixed. The first
     * published figures had Laravel at 113 req/s on a single-row read; corrected,
     * it does 175. See the guide page for what was wrong — dev dependencies in
     * the autoloader, a fresh Postgres connection per request while every other
     * framework pooled, and an SSL handshake attempt on each of those.
     *
     * Laravel saturated its own container in all four scenarios (428-437% of a
     * 400% allowance) with Postgres at 47-193%, so these are genuine ceilings.
     */
    stack: 'PHP 8.4, Laravel 13, nginx + php-fpm, OPcache + JIT, Eloquent',
    results: {
      show:  { rps: 275, gritRps: 7167, median: '68.8 ms', gritMedian: '5.2 ms',  appCpu: 402, dbCpu: 44,  gritAppCpu: 287, gritDbCpu: 253 },
      write: { rps: 318, gritRps: 4497, median: '23.9 ms', gritMedian: '8.0 ms',  appCpu: 403, dbCpu: 38,  gritAppCpu: 295, gritDbCpu: 248 },
      list:  { rps: 220, gritRps: 1130, median: '71.7 ms', gritMedian: '37.6 ms', appCpu: 404, dbCpu: 189, gritAppCpu: 119, gritDbCpu: 819 },
      mixed: { rps: 210, gritRps: 1180, median: '78.7 ms', gritMedian: '34.7 ms', appCpu: 420, dbCpu: 204, gritAppCpu: 121, gritDbCpu: 820 },
    },
  },
]

/**
 * Frameworks with a guide page but no published numbers, and why. Saying this
 * out loud is the point — a missing column otherwise reads as a framework that
 * was too slow to mention.
 */
export const UNMEASURED: {
  slug: string
  name: string
  logo: string
  color: string
  invertOnDark?: boolean
  reason: string
}[] = [
  // Empty, and that is the goal: every framework with a guide now has a
  // measured pair behind it. Add an entry here rather than quietly shipping a
  // guide with no numbers.
]

export const bySlug = (slug: string) => FRAMEWORKS.find((f) => f.slug === slug)

/** How many times faster Grit was in this pair. Below 1 means Grit lost. */
export const ratio = (r: ScenarioResult) => r.gritRps / r.rps

/** Ordered by Grit's margin, largest first — how the chart reads. */
export function ranked(scenario: ScenarioId): Framework[] {
  return [...FRAMEWORKS].sort(
    (a, b) => ratio(b.results[scenario]) - ratio(a.results[scenario]),
  )
}

/**
 * A number is only a real ceiling for a framework when that framework's own
 * container is what saturated. Where Postgres ran out first while the app still
 * had headroom, the figure is a floor — it would go higher on a bigger database,
 * and quoting it as "X does N req/s" overstates what was measured.
 *
 * A zero means the sample was not captured for that run, so nothing is claimed.
 */
type Cpu = { appCpu: number; dbCpu: number }

export function isAppBound(r: Cpu): boolean {
  if (!r.appCpu && !r.dbCpu) return false
  return r.appCpu > 350
}

export function isDbBound(r: Cpu): boolean {
  return r.dbCpu > 700 && r.appCpu > 0 && r.appCpu < 350
}

/** The same two helpers, read against Grit's side of the pair. */
export const gritCpu = (r: ScenarioResult): Cpu => ({
  appCpu: r.gritAppCpu,
  dbCpu: r.gritDbCpu,
})
