// Security headers for the docs site itself.
//
// Deliberately NOT the same policy the scaffold ships to generated apps
// (internal/scaffold/next_security_headers.go). A generated app is self-
// contained, so it can run a strict Content-Security-Policy. This site is not:
// it loads Umami analytics from analytics.gritframework.dev, an AI widget,
// uploadthing assets and a Google Apps Script form endpoint. A CSP here needs
// every one of those origins allow-listed and verified in a browser, because a
// missed origin fails silently (blocked script, no HTTP error).
//
// So: ship the six headers that carry no breakage risk now — they take the site
// from an F to an A on securityheaders.com — and add CSP separately once it's
// been exercised against the real third-party set.
const securityHeaders = [
  { key: 'X-Content-Type-Options', value: 'nosniff' },
  { key: 'X-Frame-Options', value: 'DENY' },
  { key: 'Referrer-Policy', value: 'strict-origin-when-cross-origin' },
  {
    key: 'Permissions-Policy',
    value: 'camera=(), microphone=(), geolocation=(), payment=(), usb=()',
  },
  // Browsers ignore HSTS over plain http, so this is safe in dev too.
  {
    key: 'Strict-Transport-Security',
    value: 'max-age=63072000; includeSubDomains; preload',
  },
  { key: 'Cross-Origin-Opener-Policy', value: 'same-origin' },
]

/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  // Don't advertise the framework + version.
  poweredByHeader: false,
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
  async headers() {
    return [{ source: '/:path*', headers: securityHeaders }]
  },
  // Deployment used to be split between /docs/infrastructure and
  // /docs/deployment, which meant two places to look and a Dokploy guide in
  // each. Everything about shipping now lives under /docs/deployment. These
  // are permanent: the old URLs are indexed and linked from outside the site.
  async redirects() {
    return [
      {
        source: '/docs/infrastructure/deployment',
        destination: '/docs/deployment/vps',
        permanent: true,
      },
      {
        source: '/docs/infrastructure/deploy-command',
        destination: '/docs/deployment/deploy-command',
        permanent: true,
      },
      {
        source: '/docs/infrastructure/dokploy',
        destination: '/docs/deployment/dokploy',
        permanent: true,
      },
    ]
  },
}

export default nextConfig
