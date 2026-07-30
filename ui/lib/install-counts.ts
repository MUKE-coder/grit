import { ZENITH } from './analytics'

/**
 * Install counts per block, keyed by REGISTRY NAME, read from Zenith.
 *
 *   GET /api/stats/events?name=block_copy
 *   X-Zenith-API-Key: <secret>
 *
 *   {
 *     "events": [{ "name": "block_copy", "count": 1204, "visitors": 980 }],
 *     "props":  [{ "key": "block", "value": "marketing-hero-sections-simple-centered",
 *                  "count": 812 }]
 *   }
 *
 * `props` is a flat breakdown across every property the event carries — one row
 * per (key, value) pair — so the per-block figures are the rows where
 * key === "block". Zenith has already aggregated them, which means each block
 * appears exactly once with its total across both copy kinds. Do not sum
 * duplicates here; there are none, and summing would double-count if the
 * endpoint ever returned a `kind` dimension alongside.
 *
 * Registry name, not block slug: slugs are only unique within a subcategory
 * ("stats" exists under both Marketing and Application UI), so counting by slug
 * would silently merge two different blocks into one number.
 *
 * TWO THINGS THIS DOES ON PURPOSE:
 *
 * It fails to an empty map, never to an error and never to a zero. An empty map
 * hides the badges; a map of zeros would state on every block that nobody has
 * ever installed it. A wrong number is worse than no number.
 *
 * It is cached with a revalidate window rather than fetched per request. These
 * pages are statically generated, and going dynamic to render a soft metric
 * would couple every page view to the analytics service being up. With no API
 * key at build time the first render simply has no badges, and the numbers
 * appear on the first revalidation once the key is in the runtime environment.
 */

const EVENT_NAME = 'block_copy'
const PROP_KEY = 'block'
const REVALIDATE_SECONDS = 900 // 15 minutes; installs are not a live metric.

type CountMap = Record<string, number>

export async function getInstallCounts(): Promise<CountMap> {
  const apiKey = process.env.ZENITH_API_KEY
  if (!apiKey || !ZENITH.backendUrl) return {}

  try {
    const url = new URL('/api/stats/events', ZENITH.backendUrl)
    // `name` is what drills into the property breakdown. Without it the response
    // carries only the top-level per-event totals and no `props` array at all.
    url.searchParams.set('name', EVENT_NAME)

    const res = await fetch(url, {
      headers: { 'X-Zenith-API-Key': apiKey },
      next: { revalidate: REVALIDATE_SECONDS },
    })

    if (!res.ok) return {}

    return parseProps(await res.json())
  } catch {
    // Analytics being unreachable must never fail a page build.
    return {}
  }
}

/**
 * Picks the block rows out of the flat property breakdown.
 *
 * Validates every field rather than trusting the shape: this runs during the
 * build, and a `count` that arrives as a string would otherwise render as
 * "NaN installs" on a live page.
 */
function parseProps(payload: unknown): CountMap {
  const out: CountMap = {}
  if (!payload || typeof payload !== 'object') return out

  const props = (payload as { props?: unknown }).props
  if (!Array.isArray(props)) return out

  for (const row of props) {
    if (!row || typeof row !== 'object') continue
    const r = row as Record<string, unknown>

    if (r.key !== PROP_KEY) continue
    if (typeof r.value !== 'string' || r.value === '') continue

    const n = typeof r.count === 'number' ? r.count : Number(r.count)
    if (!Number.isFinite(n) || n <= 0) continue

    out[r.value] = Math.floor(n)
  }

  return out
}
