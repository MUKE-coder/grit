import type { NextConfig } from 'next'

const nextConfig: NextConfig = {
  // Standalone emits a self-contained server with only the modules it needs,
  // which is what a production Next deploy ships. Running `next start` against
  // a full node_modules tree would measure a different thing.
  output: 'standalone',
  // Nothing here renders; the app is four route handlers. Turning off the
  // powered-by header keeps responses byte-comparable with the others.
  poweredByHeader: false,
}

export default nextConfig
