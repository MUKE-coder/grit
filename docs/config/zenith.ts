import type { ZenithConfig } from 'zenith-analytics'

// PUBLIC by design: both values ship inside the tracking snippet on every page,
// and the site key only authorizes writing events.
//
// Deliberately NO fallback values. These are read at build time and baked into
// every prerendered page, so a default here cannot be corrected later — it just
// produces a site that looks instrumented and reports nothing. A placeholder
// default also silently defeats the `required` prop on <Analytics>, which only
// throws when a value is genuinely missing.
//
// Local builds need these in docs/.env.local (see .env.example).
export const ZENITH_PUBLIC = {
  backendUrl: process.env.ZENITH_URL,
  siteKey: process.env.ZENITH_SITE_KEY,
}

// The public half plus the three secrets. Server-side only. Note the secrets
// have no fallback values — absent from the environment means undefined.
export const ZENITH_CONFIG: Partial<ZenithConfig> = {
  ...ZENITH_PUBLIC,
  apiKey: process.env.ZENITH_API_KEY, // reads analytics
  dashboardPath: '/zenith',
  protected: true,
  passwordHash: process.env.ZENITH_PW_HASH, // gates the dashboard
  jwtSecret: process.env.ZENITH_JWT_SECRET, // signs the session cookie
  siteDomain: 'gritframework.dev',
}

// createZenithRoute validates its config at module load and throws on missing
// secrets — correct for a production deploy, fatal for a local build without
// env vars. The dashboard route mounts the real handler only when this is true.
export function zenithDashboardReady(): boolean {
  return Boolean(
    ZENITH_CONFIG.apiKey && ZENITH_CONFIG.passwordHash && ZENITH_CONFIG.jwtSecret,
  )
}
