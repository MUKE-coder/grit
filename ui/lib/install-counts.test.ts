/**
 * Run with:  npx tsx lib/install-counts.test.ts
 *
 * Stubs fetch rather than calling the live Zenith instance, so this asserts the
 * two things that can silently break: the exact request (path, query, auth
 * header) and the mapping of the documented response onto name -> count.
 *
 * The payloads below match https://zenith.gritframework.dev/docs/api in shape.
 * If that response ever changes, this is where it shows up — the alternative is
 * a page that renders no badges and reports no error.
 *
 * Everything runs inside main(): this package has no "type": "module", so a
 * top-level await here is a transform error, not a runtime one.
 */
import assert from 'node:assert/strict'

process.env.ZENITH_URL = 'https://analytics.example.dev'
process.env.ZENITH_SITE_KEY = 'zk_test'
process.env.ZENITH_API_KEY = 'zsk_test'

type Captured = { url: string; headers: Record<string, string> }

function stub(payload: unknown, ok = true): Captured {
  const captured: Captured = { url: '', headers: {} }
  globalThis.fetch = (async (input: string | URL | Request, init?: RequestInit) => {
    captured.url = String(input)
    captured.headers = (init?.headers ?? {}) as Record<string, string>
    return {
      ok,
      status: ok ? 200 : 401,
      json: async () => payload,
    } as Response
  }) as typeof fetch
  return captured
}

async function main() {
  // Imported after the environment is set: the module reads ZENITH_URL at load.
  const { getInstallCounts } = await import('./install-counts')

  // ── The documented happy path ──────────────────────────────────────────────
  {
    const captured = stub({
      events: [{ name: 'block_copy', count: 1204, visitors: 980 }],
      props: [
        // Rows for the property we want.
        { key: 'block', value: 'marketing-hero-sections-simple-centered', count: 812 },
        {
          key: 'block',
          value: 'marketing-feature-sections-three-column-icons',
          count: 41,
        },
        // Rows for other properties on the same event must be ignored, not
        // mistaken for blocks — "command" is a kind, not a block name.
        { key: 'kind', value: 'command', count: 700 },
        { key: 'kind', value: 'code', count: 504 },
        { key: 'category', value: 'marketing', count: 1204 },
      ],
    })

    const counts = await getInstallCounts()

    assert.equal(
      captured.url,
      'https://analytics.example.dev/api/stats/events?name=block_copy',
      `request URL was ${captured.url}`,
    )
    assert.equal(captured.headers['X-Zenith-API-Key'], 'zsk_test', 'API key header missing')

    assert.deepEqual(counts, {
      'marketing-hero-sections-simple-centered': 812,
      'marketing-feature-sections-three-column-icons': 41,
    })
  }

  // ── Everything that must fail to {} rather than to zeros ───────────────────
  {
    // Non-2xx (bad or revoked key).
    stub({ props: [{ key: 'block', value: 'x', count: 9 }] }, false)
    assert.deepEqual(await getInstallCounts(), {}, 'a 401 must yield no counts')
  }
  {
    // No props at all — what the endpoint returns when `name` is omitted.
    stub({ events: [{ name: 'block_copy', count: 1204 }] })
    assert.deepEqual(await getInstallCounts(), {}, 'missing props must yield no counts')
  }
  {
    // A changed envelope must not throw and must not invent numbers.
    stub({ data: [{ block: 'x', count: 9 }] })
    assert.deepEqual(await getInstallCounts(), {}, 'unknown shape must yield no counts')
  }
  {
    stub(null)
    assert.deepEqual(await getInstallCounts(), {}, 'null body must yield no counts')
  }
  {
    globalThis.fetch = (async () => {
      throw new Error('ECONNREFUSED')
    }) as typeof fetch
    assert.deepEqual(await getInstallCounts(), {}, 'a network failure must never throw')
  }

  // ── Junk rows are dropped individually, not fatally ────────────────────────
  {
    stub({
      props: [
        { key: 'block', value: 'good', count: 5 },
        { key: 'block', value: 'stringly', count: '7' }, // coerced
        { key: 'block', value: 'nan', count: 'abc' }, // dropped, never NaN
        { key: 'block', value: '', count: 3 }, // no name, dropped
        { key: 'block', count: 3 }, // no value, dropped
        { key: 'block', value: 'zero', count: 0 }, // 0 renders nothing anyway
        { key: 'block', value: 'negative', count: -4 }, // dropped
        null,
        'nonsense',
      ],
    })

    const counts = await getInstallCounts()
    assert.deepEqual(counts, { good: 5, stringly: 7 })
    assert.ok(!('nan' in counts), 'NaN must never reach the page')
  }

  // ── No key configured means no request at all ──────────────────────────────
  {
    delete process.env.ZENITH_API_KEY
    let called = false
    globalThis.fetch = (async () => {
      called = true
      return { ok: true, json: async () => ({}) } as Response
    }) as typeof fetch

    assert.deepEqual(await getInstallCounts(), {})
    assert.equal(called, false, 'must not call Zenith without an API key')
  }

  console.log('install-counts: all assertions passed')
}

main().catch((err) => {
  console.error(err)
  process.exit(1)
})
