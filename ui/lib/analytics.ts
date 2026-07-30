import type { ZenithConfig } from 'zenith-analytics'

/**
 * Analytics config for the Grit UI site.
 *
 * Deliberately OPTIONAL, unlike the docs site. The docs use <Analytics required />
 * so a misconfigured deploy fails loudly, which is right there — analytics is
 * expected on that site and a silent miss is a real loss. Here the site is
 * useful with or without measurement, so a missing key must not fail the build;
 * it just means no events. Failing a deploy over an optional key is hostile.
 */
export const ZENITH = {
  backendUrl: process.env.ZENITH_URL ?? '',
  siteKey: process.env.ZENITH_SITE_KEY ?? '',
}

/** True only when both halves are present, so <Analytics> is never rendered half-configured. */
export const analyticsEnabled = Boolean(ZENITH.backendUrl && ZENITH.siteKey)

/**
 * The public pair plus the three secrets, for the /zenith dashboard route.
 *
 * SERVER ONLY. Never import this into a client component — `apiKey` can read all
 * analytics and `jwtSecret` signs the dashboard session cookie.
 *
 * No fallbacks on the secrets: absent from the environment has to stay
 * `undefined` so zenithDashboardReady() can tell configured from not. A `?? ''`
 * here would pass every truthiness check and mount a dashboard that cannot
 * authenticate anyone — the same trap that was hit on the docs site.
 */
export const ZENITH_CONFIG: Partial<ZenithConfig> = {
  ...ZENITH,
  apiKey: process.env.ZENITH_API_KEY, // reads analytics
  dashboardPath: '/zenith',
  protected: true,
  passwordHash: process.env.ZENITH_PW_HASH, // gates the dashboard
  jwtSecret: process.env.ZENITH_JWT_SECRET, // signs the session cookie
  siteDomain: 'ui.gritframework.dev',
}

/**
 * createZenithRoute validates its config at module load and throws on missing
 * secrets — correct for a production deploy, fatal for a local build or a CI run
 * with no env vars. The route mounts the real handler only when this is true.
 */
export function zenithDashboardReady(): boolean {
  return Boolean(
    ZENITH_CONFIG.backendUrl &&
      ZENITH_CONFIG.apiKey &&
      ZENITH_CONFIG.passwordHash &&
      ZENITH_CONFIG.jwtSecret,
  )
}
