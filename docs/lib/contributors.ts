// The contributor list, from GitHub, cached for an hour.
//
// Fetched rather than kept in a data file, because a hand-maintained roster of
// contributors is a roster that is wrong the first time somebody lands a PR and
// nobody remembers to edit it. An hour is plenty: this page changes a few times
// a year.
//
// The fallback is the point of the shape. Unauthenticated GitHub allows 60
// requests an hour per IP, and this page should show the people who built the
// thing even when that budget is gone or GitHub is having a morning. See
// app/api/github-stars/route.ts, which soft-fails the same way.

export interface Contributor {
  login: string
  avatarUrl: string
  profileUrl: string
  commits: number
}

const REPO = 'MUKE-coder/grit'

/**
 * Last known good, used when GitHub cannot be reached.
 *
 * Deliberately not a complete mirror of the API: it carries the people, not
 * their exact counts, because a stale count presented as live is worse than an
 * obviously-rounded one. `commits` here is only used for ordering.
 */
const FALLBACK: Contributor[] = [
  {
    login: 'MUKE-coder',
    avatarUrl: 'https://avatars.githubusercontent.com/u/64189841?v=4',
    profileUrl: 'https://github.com/MUKE-coder',
    commits: 1259,
  },
  {
    login: 'mosespace',
    avatarUrl: 'https://avatars.githubusercontent.com/u/111237257?v=4',
    profileUrl: 'https://github.com/mosespace',
    commits: 6,
  },
  {
    login: 'daljo25',
    avatarUrl: 'https://avatars.githubusercontent.com/u/7244602?v=4',
    profileUrl: 'https://github.com/daljo25',
    commits: 1,
  },
]

interface GitHubContributor {
  login?: string
  type?: string
  avatar_url?: string
  html_url?: string
  contributions?: number
}

export async function getContributors(): Promise<{
  contributors: Contributor[]
  live: boolean
}> {
  try {
    const res = await fetch(
      `https://api.github.com/repos/${REPO}/contributors?per_page=100&anon=0`,
      {
        headers: { Accept: 'application/vnd.github+json' },
        next: { revalidate: 3600 },
      },
    )
    if (!res.ok) throw new Error(`GitHub responded ${res.status}`)

    const data: unknown = await res.json()
    if (!Array.isArray(data)) throw new Error('GitHub returned something that is not a list')

    const contributors = (data as GitHubContributor[])
      // Bots commit. They are not contributors, and a page that lists one
      // beside the people who wrote the framework is a page nobody trusts.
      .filter((c) => c.type === 'User' && typeof c.login === 'string')
      .map((c) => ({
        login: c.login as string,
        avatarUrl: c.avatar_url ?? '',
        profileUrl: c.html_url ?? `https://github.com/${c.login}`,
        commits: typeof c.contributions === 'number' ? c.contributions : 0,
      }))
      .sort((a, b) => b.commits - a.commits)

    if (contributors.length === 0) throw new Error('GitHub returned an empty list')
    return { contributors, live: true }
  } catch {
    return { contributors: FALLBACK, live: false }
  }
}
