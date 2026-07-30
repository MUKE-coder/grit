import { ZENITH } from './analytics'

/**
 * Install counts per block, keyed by registry name.
 *
 * Reading is deliberately separated from the UI behind one function, so the data
 * source can change without touching a component. Right now that source is the
 * self-hosted Zenith instance; it could later be a store of our own.
 *
 * TWO THINGS THIS DOES ON PURPOSE:
 *
 * It fails to an empty map, never to an error and never to a zero. A missing map
 * hides the badges entirely; a map full of zeros would state, on every block,
 * that nobody has ever installed it. A wrong number is worse than no number.
 *
 * It is cached with a revalidate window rather than fetched per request. These
 * pages are statically generated, and making them dynamic to render a
 * soft metric would couple every page view to the analytics service being up.
 */

const REVALIDATE_SECONDS = 900 // 15 minutes; installs are not a live metric.

type CountMap = Record<string, number>

export async function getInstallCounts(): Promise<CountMap> {
  const apiKey = process.env.ZENITH_API_KEY
  if (!apiKey || !ZENITH.backendUrl) return {}

  try {
    const url = new URL('/api/stats/events', ZENITH.backendUrl)
    url.searchParams.set('event', 'block_copy')
    url.searchParams.set('group_by', 'block')
    url.searchParams.set('range', 'all')

    const res = await fetch(url, {
      headers: { 'X-Zenith-API-Key': apiKey },
      next: { revalidate: REVALIDATE_SECONDS },
    })

    if (!res.ok) return {}

    return parseCounts(await res.json())
  } catch {
    // Analytics being unreachable must never fail a page build.
    return {}
  }
}

/**
 * Normalises whatever shape the stats endpoint returns into name -> count.
 *
 * Written defensively across the few plausible shapes rather than against one
 * assumed schema, because an aggregate endpoint that changes its envelope would
 * otherwise silently produce zeros — and zeros render as real numbers.
 */
function parseCounts(payload: unknown): CountMap {
  const out: CountMap = {}
  if (!payload || typeof payload !== 'object') return out

  const rows =
    (payload as { rows?: unknown[] }).rows ??
    (payload as { data?: unknown[] }).data ??
    (payload as { results?: unknown[] }).results ??
    (Array.isArray(payload) ? payload : [])

  if (!Array.isArray(rows)) return out

  for (const row of rows) {
    if (!row || typeof row !== 'object') continue
    const r = row as Record<string, unknown>

    const key = r.block ?? r.value ?? r.name ?? r.label
    const raw = r.count ?? r.total ?? r.events ?? r.value

    if (typeof key !== 'string') continue
    const n = typeof raw === 'number' ? raw : Number(raw)
    if (!Number.isFinite(n) || n < 0) continue

    // Same block can appear once per `kind`; sum them into one install figure.
    out[key] = (out[key] ?? 0) + Math.floor(n)
  }

  return out
}
