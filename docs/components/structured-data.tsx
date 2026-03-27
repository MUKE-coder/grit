interface StructuredDataProps {
  data: Record<string, unknown>
}

export function StructuredData({ data }: StructuredDataProps) {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(data) }}
    />
  )
}

// Organization schema — used in root layout
export function OrganizationSchema() {
  return (
    <StructuredData
      data={{
        '@context': 'https://schema.org',
        '@type': 'Organization',
        name: 'Grit',
        url: 'https://gritframework.dev',
        logo: 'https://gritframework.dev/opengraph-image.png',
        description:
          'The full-stack meta-framework that combines Go, React, and a Filament-like admin panel.',
        sameAs: [
          'https://github.com/MUKE-coder/grit',
          'https://www.youtube.com/@GritFramework',
          'https://www.linkedin.com/company/grit-framework',
          'https://jb.desishub.com',
        ],
      }}
    />
  )
}

// SoftwareApplication schema — used on home page
export function SoftwareApplicationSchema() {
  return (
    <StructuredData
      data={{
        '@context': 'https://schema.org',
        '@type': 'SoftwareApplication',
        name: 'Grit',
        applicationCategory: 'DeveloperApplication',
        operatingSystem: 'Windows, macOS, Linux',
        description:
          'Grit is a full-stack meta-framework that combines Go (Gin + GORM) with React (Next.js) and a Filament-like admin panel. One CLI command scaffolds an entire monorepo with API, web app, admin dashboard, shared types, and Docker setup.',
        url: 'https://gritframework.dev',
        downloadUrl: 'https://github.com/MUKE-coder/grit',
        author: {
          '@type': 'Person',
          name: 'Muke JohnBaptist',
          url: 'https://jb.desishub.com',
        },
        offers: {
          '@type': 'Offer',
          price: '0',
          priceCurrency: 'USD',
        },
        programmingLanguage: ['Go', 'TypeScript', 'JavaScript'],
        softwareRequirements: 'Go 1.21+, Node.js 18+, pnpm',
      }}
    />
  )
}

// FAQPage schema — used on home and docs pages for AEO
export function FAQPageSchema() {
  return (
    <StructuredData
      data={{
        '@context': 'https://schema.org',
        '@type': 'FAQPage',
        mainEntity: [
          {
            '@type': 'Question',
            name: 'What is Grit?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Grit is a full-stack meta-framework that combines Go (Gin + GORM) for the backend with React (Next.js + TypeScript) for the frontend, plus a Filament-like admin panel. One CLI command scaffolds an entire monorepo with API, web app, admin dashboard, shared types, and Docker setup.',
            },
          },
          {
            '@type': 'Question',
            name: 'How do I install Grit?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Install Grit with: go install github.com/MUKE-coder/grit/v3/cmd/grit@latest. Then run "grit new my-project" to scaffold a full-stack project. You need Go 1.21+, Node.js 18+, and pnpm installed.',
            },
          },
          {
            '@type': 'Question',
            name: 'What tech stack does Grit use?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Grit uses Go with Gin web framework and GORM ORM for the backend, Next.js 14+ with App Router and TypeScript for the frontend, Tailwind CSS with shadcn/ui for styling, React Query for data fetching, Zod for validation, Turborepo with pnpm for monorepo management, and Docker for infrastructure.',
            },
          },
          {
            '@type': 'Question',
            name: 'How does Grit compare to Laravel, Django, or Rails?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Unlike Laravel (PHP), Django (Python), or Rails (Ruby) which are backend-only frameworks, Grit is a true full-stack framework. It scaffolds both the Go API backend and the React frontend together in a monorepo with shared types, an admin panel, and code generation. It combines the productivity of these frameworks with the performance of Go and the modern React ecosystem.',
            },
          },
          {
            '@type': 'Question',
            name: 'Can Grit generate full-stack resources?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Yes. Running "grit generate resource Product --fields name:string,price:float,active:bool" creates the Go model, service, handler, routes, Zod schemas, TypeScript types, React Query hooks, and admin panel resource definition — all wired together and ready to use.',
            },
          },
          {
            '@type': 'Question',
            name: 'What batteries does Grit include?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Grit includes Redis caching with middleware, S3-compatible file storage (AWS S3, Cloudflare R2, MinIO), email via Resend with HTML templates, background jobs and cron scheduling via asynq, AI integration (Claude + OpenAI with streaming), security suite with WAF and rate limiting, and a visual database browser (GORM Studio).',
            },
          },
          {
            '@type': 'Question',
            name: 'Does Grit have an admin panel?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Yes. Grit scaffolds a Filament-like admin panel built with Next.js, shadcn/ui, and React Query. It features runtime resource definitions, advanced DataTable with sorting/filtering/pagination, FormBuilder with 10+ field types including multi-step forms, dashboard widgets, dark/light theme, and RBAC-based navigation.',
            },
          },
          {
            '@type': 'Question',
            name: 'Is Grit free and open source?',
            acceptedAnswer: {
              '@type': 'Answer',
              text: 'Yes. Grit is free and open source, available on GitHub at github.com/MUKE-coder/grit. It is built and maintained by Muke JohnBaptist.',
            },
          },
        ],
      }}
    />
  )
}

// TechArticle schema — used on individual doc pages
export function TechArticleSchema({
  title,
  description,
  url,
}: {
  title: string
  description: string
  url: string
}) {
  return (
    <StructuredData
      data={{
        '@context': 'https://schema.org',
        '@type': 'TechArticle',
        headline: title,
        description,
        url,
        author: {
          '@type': 'Person',
          name: 'Muke JohnBaptist',
          url: 'https://jb.desishub.com',
        },
        publisher: {
          '@type': 'Organization',
          name: 'Grit',
          url: 'https://gritframework.dev',
          logo: 'https://gritframework.dev/opengraph-image.png',
        },
        datePublished: '2026-02-01',
        dateModified: new Date().toISOString().split('T')[0],
        mainEntityOfPage: url,
        inLanguage: 'en',
      }}
    />
  )
}
