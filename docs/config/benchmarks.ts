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
 * Grit's single-row read measured 4,536 req/s in the Bun pair, 4,392 in the
 * Encore pair, 1,911 in the Next.js pair and 1,635 in the Express pair — from an
 * identical binary, as hours of write scenarios accumulated in Postgres. So the absolute figures are
 * only meaningful next to the Grit baseline measured beside them, and the RATIO
 * is what survives comparison. Putting Bun's 2,717 on the same axis as Encore's
 * 438 would imply a shared baseline that does not exist.
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
      list:  { rps: 590,  gritRps: 615,  median: '70.8 ms', gritMedian: '65.6 ms', appCpu: 267, dbCpu: 798, gritAppCpu: 125, gritDbCpu: 803 },
      mixed: { rps: 621,  gritRps: 568,  median: '65.8 ms', gritMedian: '70.9 ms', appCpu: 327, dbCpu: 825, gritAppCpu: 125, gritDbCpu: 764 },
    },
  },
  {
    slug: 'encore',
    preV3134: true,
    logo: '/logos/encore.png',
    color: '#EEEBE0',
    name: 'Encore.ts',
    version: 'v1.57 + Drizzle',
    orm: 'Drizzle',
    tagline: 'Rust HTTP runtime, compiled by Encore’s own CLI, Drizzle over its SQLDatabase',
    /*
     * Worth stating plainly: Encore never saturated anything. Its container sat
     * at 155-211% of a 400% allowance and Postgres at 30-200%, so neither was
     * the limit. Something else bounded it — most likely the Drizzle /
     * node-postgres path rather than the Rust HTTP layer, which is the part
     * Encore is fast at. Its figures here are a floor for Encore, and a tuned
     * Encore setup would very likely do better.
     */
    stack: 'Encore.ts 1.57, Rust runtime, Drizzle ORM',
    results: {
      show:  { rps: 438, gritRps: 4392, median: '103.8 ms', gritMedian: '8.7 ms',  appCpu: 165, dbCpu: 32,  gritAppCpu: 317, gritDbCpu: 248 },
      write: { rps: 633, gritRps: 1063, median: '73.7 ms',  gritMedian: '41.5 ms', appCpu: 206, dbCpu: 70,  gritAppCpu: 334, gritDbCpu: 514 },
      list:  { rps: 247, gritRps: 548,  median: '184.7 ms', gritMedian: '78.4 ms', appCpu: 159, dbCpu: 190, gritAppCpu: 156, gritDbCpu: 897 },
      mixed: { rps: 257, gritRps: 523,  median: '190.0 ms', gritMedian: '80.3 ms', appCpu: 162, dbCpu: 177, gritAppCpu: 146, gritDbCpu: 975 },
    },
  },
  {
    slug: 'express',
    preV3134: true,
    logo: '/logos/express.png',
    color: '#5FA04E',
    name: 'Express',
    version: 'v5 + Prisma',
    orm: 'Prisma',
    tagline: 'Express 5 on Node 22, one cluster worker per CPU, Prisma ORM',
    stack: 'Node 22, Express 5, cluster, Prisma',
    results: {
      show:  { rps: 306, gritRps: 1635, median: '154.9 ms', gritMedian: '22.6 ms', appCpu: 476, dbCpu: 59,  gritAppCpu: 401, gritDbCpu: 324 },
      write: { rps: 415, gritRps: 513,  median: '105.1 ms', gritMedian: '88.5 ms', appCpu: 435, dbCpu: 51,  gritAppCpu: 482, gritDbCpu: 629 },
      list:  { rps: 233, gritRps: 448,  median: '205.7 ms', gritMedian: '95.3 ms', appCpu: 452, dbCpu: 374, gritAppCpu: 138, gritDbCpu: 940 },
      mixed: { rps: 177, gritRps: 458,  median: '272.9 ms', gritMedian: '93.5 ms', appCpu: 458, dbCpu: 312, gritAppCpu: 127, gritDbCpu: 900 },
    },
  },
  {
    slug: 'nextjs',
    preV3134: true,
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
      show:  { rps: 183, gritRps: 1911, median: '258.7 ms', gritMedian: '18.6 ms', appCpu: 411, dbCpu: 21,  gritAppCpu: 296, gritDbCpu: 235 },
      write: { rps: 336, gritRps: 1142, median: '132.7 ms', gritMedian: '38.0 ms', appCpu: 406, dbCpu: 19,  gritAppCpu: 314, gritDbCpu: 482 },
      list:  { rps: 169, gritRps: 550,  median: '284.0 ms', gritMedian: '74.8 ms', appCpu: 412, dbCpu: 172, gritAppCpu: 114, gritDbCpu: 846 },
      mixed: { rps: 186, gritRps: 533,  median: '258.6 ms', gritMedian: '75.5 ms', appCpu: 408, dbCpu: 125, gritAppCpu: 128, gritDbCpu: 842 },
    },
  },
  {
    slug: 'django',
    preV3134: true,
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
      show:  { rps: 608, gritRps: 4651, median: '87.0 ms',  gritMedian: '8.0 ms',  appCpu: 404, dbCpu: 51,  gritAppCpu: 279, gritDbCpu: 238 },
      write: { rps: 797, gritRps: 1583, median: '68.1 ms',  gritMedian: '26.8 ms', appCpu: 407, dbCpu: 58,  gritAppCpu: 303, gritDbCpu: 498 },
      list:  { rps: 274, gritRps: 736,  median: '147.8 ms', gritMedian: '56.0 ms', appCpu: 407, dbCpu: 259, gritAppCpu: 122, gritDbCpu: 680 },
      mixed: { rps: 249, gritRps: 662,  median: '191.4 ms', gritMedian: '64.7 ms', appCpu: 426, dbCpu: 250, gritAppCpu: 121, gritDbCpu: 814 },
    },
  },
  {
    slug: 'laravel',
    preV3134: true,
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
      show:  { rps: 175, gritRps: 4340, median: '92.1 ms',  gritMedian: '8.6 ms',  appCpu: 434, dbCpu: 47,  gritAppCpu: 300, gritDbCpu: 248 },
      write: { rps: 146, gritRps: 1118, median: '56.7 ms',  gritMedian: '40.5 ms', appCpu: 428, dbCpu: 47,  gritAppCpu: 316, gritDbCpu: 531 },
      list:  { rps: 125, gritRps: 523,  median: '93.3 ms',  gritMedian: '81.4 ms', appCpu: 437, dbCpu: 193, gritAppCpu: 173, gritDbCpu: 950 },
      mixed: { rps: 101, gritRps: 585,  median: '151.7 ms', gritMedian: '70.4 ms', appCpu: 433, dbCpu: 183, gritAppCpu: 178, gritDbCpu: 981 },
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
