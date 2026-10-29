// Current Grit CLI version. Single source of truth for every place the
// docs site displays a version string (header badge, install lesson
// example output, verify-install lesson, animated terminal). Bump this
// when the Go CLI releases; keep it in sync with cmd/grit/main.go's
// `var version` and internal/scaffold.DefaultVersion.
export const GRIT_VERSION = '3.438.0'

/**
 * Which channel the version above belongs to.
 *
 * A release is stable unless its version carries a prerelease suffix, and the
 * suffix is the only thing that decides it: `3.398.0-rc.1` is a release
 * candidate, `3.398.0` is stable. That is not just a naming convention here.
 * The release workflow marks a tag containing -rc, -beta or -alpha as a GitHub
 * prerelease, and `grit update` asks GitHub for the latest release, which by
 * definition excludes prereleases. So a prerelease is something you install on
 * purpose and never something you are moved onto.
 *
 * Derived rather than configured, so the badge cannot disagree with the version
 * it sits next to.
 */
export type ReleaseChannel = 'stable' | 'rc' | 'beta' | 'alpha'

export function releaseChannel(version: string = GRIT_VERSION): ReleaseChannel {
  const suffix = version.split('-')[1]?.toLowerCase() ?? ''
  if (suffix.startsWith('rc')) return 'rc'
  if (suffix.startsWith('beta')) return 'beta'
  if (suffix.startsWith('alpha')) return 'alpha'
  return 'stable'
}

/** What the channel is called where a reader sees it. */
export const RELEASE_CHANNEL_LABEL: Record<ReleaseChannel, string> = {
  stable: 'stable',
  rc: 'release candidate',
  beta: 'beta',
  alpha: 'alpha',
}

export const siteConfig = {
  name: 'Grit',
  // "Go + React Full-Stack Framework" describes the ingredients, which a dozen
  // stacks share. The title is what shows in a search result and a shared link,
  // so it states the thing only this does: one command, the whole vertical slice.
  title: 'Grit: Full-stack Go apps from one command',
  version: GRIT_VERSION,
  description:
    'Describe a resource; Grit writes the Go model, API, migrations, TypeScript types, React hooks and admin screen. Auth, RBAC, jobs, storage, realtime, observability and one-command deploy included.',
  url: 'https://gritframework.dev',
  ogImage: 'https://gritframework.dev/opengraph-image.png',
  creator: 'MUKE-coder',
  author: 'Muke JohnBaptist',
  github: 'https://github.com/MUKE-coder/grit',
  youtube: 'https://www.youtube.com/@GritFramework',
  linkedin: 'https://www.linkedin.com/company/grit-framework',
  // The one place the invite is written. community-cta re-exports it as
  // DISCORD_INVITE_URL for the components, so a new invite is a one-line
  // change here rather than a search for the old code.
  discord: 'https://discord.gg/6nJpTgN7U',
  // Same rule as the Discord invite above: written once, here.
  whatsapp: 'https://chat.whatsapp.com/HXsOWlp5W6o9lV4YhsY376',
  website: 'https://jb.desishub.com',
  keywords: [
    'Go framework',
    'React framework',
    'full-stack framework',
    'Grit',
    'Gin',
    'GORM',
    'Next.js',
    'admin panel',
    'code generator',
    'monorepo',
    'Go + React',
    'TypeScript',
    'Tailwind CSS',
    'shadcn/ui',
    'REST API',
    'RBAC',
    'scaffolding',
    'full-stack Go',
    'Go web framework',
    'React admin dashboard',
  ],
}
