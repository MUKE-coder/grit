import { createZenithRoute } from 'zenith-analytics/next'

import { ZENITH_CONFIG, zenithDashboardReady } from '@/config/zenith'

// Without this, Next could statically render the route at build time and
// serve every visitor the same cached page — fatal for a password gate.
export const dynamic = 'force-dynamic'

const notConfigured = () =>
  new Response('Zenith dashboard is not configured on this deployment.', {
    status: 503,
  })

// Secrets present → the real dashboard proxy. Absent (local dev, CI) → a
// plain 503 instead of createZenithRoute's intentional startup throw.
const handlers = zenithDashboardReady()
  ? createZenithRoute(ZENITH_CONFIG)
  : { GET: async () => notConfigured(), POST: async () => notConfigured() }

export const { GET, POST } = handlers
