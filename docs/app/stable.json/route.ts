import { STABLE, STABLE_MIN_AGE_DAYS } from '@/config/stable'
import { GRIT_VERSION } from '@/config/site'

/**
 * The channel file the CLI reads for `grit update --channel stable`.
 *
 * A JSON endpoint rather than a GitHub tag. A moving tag means force-pushing
 * one, and a force-pushed tag is a thing nobody can audit afterwards: this is
 * a file in the repository, written by scripts/promote-stable.py, reviewable
 * in the diff that moved it (grit#245).
 *
 * `edge` is here too so one request answers both questions. A CLI that had to
 * ask GitHub for the latest and this for the stable would report them from
 * two moments, and the pair would occasionally disagree about which is newer.
 */
export const dynamic = 'force-static'

export function GET() {
  return Response.json(
    {
      stable: STABLE.version,
      edge: GRIT_VERSION,
      promoted: STABLE.promoted,
      released: STABLE.released,
      canaryRuns: STABLE.canaryRuns,
      minAgeDays: STABLE_MIN_AGE_DAYS,
      rule: 'the newest release at least 7 days old with every scheduled canary run in its first 7 days green',
      docs: 'https://gritframework.dev/docs/versioning#channels',
    },
    {
      headers: {
        // An hour. The promoter runs daily, and a CLI reading an hour-old
        // answer installs the same release as one reading a fresh one.
        'Cache-Control': 'public, max-age=0, s-maxage=3600, stale-while-revalidate=86400',
      },
    },
  )
}
