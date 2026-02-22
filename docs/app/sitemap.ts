import type { MetadataRoute } from 'next'

const baseUrl = 'https://gritframework.dev'

export default function sitemap(): MetadataRoute.Sitemap {
  const routes = [
    // Home
    { path: '/', priority: 1.0, changeFrequency: 'weekly' as const },
    { path: '/playground', priority: 0.8, changeFrequency: 'monthly' as const },

    // Getting Started
    { path: '/docs', priority: 0.9, changeFrequency: 'weekly' as const },
    { path: '/docs/getting-started/installation', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/quick-start', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/configuration', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/philosophy', priority: 0.5, changeFrequency: 'yearly' as const },
    { path: '/docs/getting-started/project-structure', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/troubleshooting', priority: 0.6, changeFrequency: 'monthly' as const },

    // Prerequisites
    { path: '/docs/prerequisites/golang', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/prerequisites/nextjs', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/prerequisites/docker', priority: 0.7, changeFrequency: 'monthly' as const },

    // Concepts
    { path: '/docs/concepts/architecture', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/concepts/cli', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/concepts/code-generation', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/concepts/naming-conventions', priority: 0.5, changeFrequency: 'yearly' as const },
    { path: '/docs/concepts/styles', priority: 0.5, changeFrequency: 'monthly' as const },
    { path: '/docs/concepts/type-system', priority: 0.6, changeFrequency: 'monthly' as const },

    // Backend
    { path: '/docs/backend/authentication', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/handlers', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/middleware', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/migrations', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/models', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/rbac', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/response-format', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/seeders', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/services', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/backend/api-docs', priority: 0.7, changeFrequency: 'monthly' as const },

    // Frontend
    { path: '/docs/frontend/hooks', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/frontend/shared-package', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/frontend/web-app', priority: 0.6, changeFrequency: 'monthly' as const },

    // Admin Panel
    { path: '/docs/admin/overview', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/resources', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/datatable', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/forms', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/multi-step-forms', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/relationships', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/widgets', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/admin/standalone-usage', priority: 0.6, changeFrequency: 'monthly' as const },

    // Batteries
    { path: '/docs/batteries/ai', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/caching', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/cron', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/email', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/jobs', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/security', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/batteries/storage', priority: 0.6, changeFrequency: 'monthly' as const },

    // Infrastructure
    { path: '/docs/infrastructure/database', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/infrastructure/deployment', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/infrastructure/docker', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/infrastructure/docker-cheatsheet', priority: 0.5, changeFrequency: 'monthly' as const },

    // Design
    { path: '/docs/design/theme', priority: 0.5, changeFrequency: 'monthly' as const },

    // Tutorials
    { path: '/docs/tutorials/blog', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/tutorials/ecommerce', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/tutorials/learn', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/docs/tutorials/product-catalog', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/tutorials/saas', priority: 0.7, changeFrequency: 'monthly' as const },

    // AI Workflows
    { path: '/docs/ai-workflows/claude', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/ai-workflows/antigravity', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/ai-skill', priority: 0.5, changeFrequency: 'monthly' as const },

    // Changelog
    { path: '/docs/changelog', priority: 0.6, changeFrequency: 'weekly' as const },
  ]

  return routes.map((route) => ({
    url: `${baseUrl}${route.path}`,
    lastModified: new Date(),
    changeFrequency: route.changeFrequency,
    priority: route.priority,
  }))
}
