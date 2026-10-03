import { llmsIndex } from '@/lib/llms'

// Static: the index is built from files in this repo, so it changes when the
// site is rebuilt and never per request.
export const dynamic = 'force-static'

export function GET() {
  return new Response(llmsIndex(), {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'Cache-Control': 'public, max-age=0, s-maxage=3600, stale-while-revalidate=86400',
    },
  })
}
