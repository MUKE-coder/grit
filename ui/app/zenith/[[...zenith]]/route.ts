import { createZenithRoute } from 'zenith-analytics/next'

import { ZENITH_CONFIG, zenithDashboardReady } from '@/lib/analytics'

// Without this, Next could statically render the route at build time and serve
// every visitor the same cached page — fatal for a password gate.
export const dynamic = 'force-dynamic'

// Deliberately outside the (site) and (preview) route groups: this is a route
// handler, so it renders no layout, and it must not pull in the site chrome or
// the tracking snippet. Zenith serves the whole dashboard itself.

const notConfigured = () =>
  new Response('Zenith dashboard is not configured on this deployment.', {
    status: 503,
  })

// Secrets present → the real dashboard proxy. Absent (local dev, CI) → a plain
// 503 instead of createZenithRoute's intentional module-load throw, which would
// otherwise fail the build on any machine without the secrets.
const handlers = zenithDashboardReady()
  ? createZenithRoute(ZENITH_CONFIG)
  : { GET: async () => notConfigured(), POST: async () => notConfigured() }

export const { GET, POST } = handlers
