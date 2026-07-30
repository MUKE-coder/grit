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
