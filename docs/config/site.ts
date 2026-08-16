// Current Grit CLI version. Single source of truth for every place the
// docs site displays a version string (header badge, install lesson
// example output, verify-install lesson, animated terminal). Bump this
// when the Go CLI releases; keep it in sync with cmd/grit/main.go's
// `var version` and internal/scaffold.DefaultVersion.
export const GRIT_VERSION = '3.152.0'

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
