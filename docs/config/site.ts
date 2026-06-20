// Current Grit CLI version. Single source of truth for every place the
// docs site displays a version string (header badge, install lesson
// example output, verify-install lesson, animated terminal). Bump this
// when the Go CLI releases; keep it in sync with cmd/grit/main.go's
// `var version` and internal/scaffold.DefaultVersion.
export const GRIT_VERSION = '3.31.5'

export const siteConfig = {
  name: 'Grit',
  title: 'Grit — Go + React Full-Stack Framework',
  version: GRIT_VERSION,
  description:
    'The full-stack meta-framework that combines Go, React, and a Filament-like admin panel. Scaffold entire projects, generate full-stack resources, and ship fast.',
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
