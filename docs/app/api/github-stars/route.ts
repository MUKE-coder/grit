import { NextResponse } from 'next/server'

// Cache the GitHub star count server-side for an hour. Unauthenticated
// GitHub API allows 60 req/hr/IP — caching keeps us well under that and
// shields the client from CORS / rate-limit flakiness. The client island
// (<GitHubStars />) just reads this endpoint.
export const revalidate = 3600

const REPO = 'MUKE-coder/grit'

export async function GET() {
  try {
    const res = await fetch(`https://api.github.com/repos/${REPO}`, {
      headers: { Accept: 'application/vnd.github+json' },
      next: { revalidate: 3600 },
    })
    if (!res.ok) throw new Error(`GitHub responded ${res.status}`)
    const data = await res.json()
    const stars = typeof data?.stargazers_count === 'number' ? data.stargazers_count : null
    return NextResponse.json(
      { stars },
      { headers: { 'Cache-Control': 'public, s-maxage=3600, stale-while-revalidate=86400' } },
    )
  } catch {
    // Soft-fail: the badge falls back to its icon-only state.
    return NextResponse.json({ stars: null })
  }
}
