import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // standalone keeps the Docker image small — Dokploy builds the Dockerfile in
  // this directory and only needs .next/standalone plus static assets.
  output: 'standalone',
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
