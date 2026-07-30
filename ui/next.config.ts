import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // standalone keeps the Docker image small — Dokploy builds the Dockerfile in
  // this directory and only needs .next/standalone plus static assets.
  output: 'standalone',
  experimental: {
    // Required for app/global-not-found.tsx to be used at all. This project has
    // two root layouts — (site) and (preview) — so an unmatched URL matches
    // neither group and a route-group not-found.tsx is silently ignored in
    // favour of Next's built-in 404 page.
    globalNotFound: true,
  },
  images: {
    remotePatterns: [{ protocol: 'https', hostname: 'picsum.photos' }],
  },
  async headers() {
    return [
      {
        // The registry is meant to be fetched by other people's tooling —
        // `npx shadcn add`, CI jobs, the Grit CLI — so it has to be readable
        // cross-origin. Everything it serves is public source code.
        source: '/r/:path*',
        headers: [
          { key: 'Access-Control-Allow-Origin', value: '*' },
          { key: 'Access-Control-Allow-Methods', value: 'GET, OPTIONS' },
        ],
      },
    ]
  },
}

export default nextConfig
