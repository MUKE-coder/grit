import type { MetadataRoute } from 'next'
import { COURSES } from '@/config/courses'
import { getAllPosts } from '@/lib/blog'

const baseUrl = 'https://gritframework.dev'

export default function sitemap(): MetadataRoute.Sitemap {
  const routes = [
    // Home
    { path: '/', priority: 1.0, changeFrequency: 'weekly' as const },
    { path: '/playground', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/hire', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/showcase', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/donate', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/pitch', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/blog', priority: 0.8, changeFrequency: 'daily' as const },

    // Getting Started
    { path: '/docs', priority: 0.9, changeFrequency: 'weekly' as const },
    { path: '/docs/getting-started/installation', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/quick-start', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/configuration', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/philosophy', priority: 0.5, changeFrequency: 'yearly' as const },
    { path: '/docs/getting-started/project-structure', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/troubleshooting', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/getting-started/create-without-docker', priority: 0.7, changeFrequency: 'monthly' as const },

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
    { path: '/docs/backend/pulse', priority: 0.7, changeFrequency: 'monthly' as const },

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
    { path: '/docs/infrastructure/dokploy', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/docs/infrastructure/docker', priority: 0.6, changeFrequency: 'monthly' as const },
    { path: '/docs/infrastructure/docker-cheatsheet', priority: 0.5, changeFrequency: 'monthly' as const },

    // Design
    { path: '/docs/design/theme', priority: 0.5, changeFrequency: 'monthly' as const },

    // Tutorials
    { path: '/docs/tutorials/contact-app', priority: 0.9, changeFrequency: 'monthly' as const },
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

    // Plugins
    { path: '/docs/plugins', priority: 0.7, changeFrequency: 'monthly' as const },

    // Courses
    { path: '/courses', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/introduction', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/first-app', priority: 0.9, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/code-generator', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/authentication', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/admin-panel', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/file-storage', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/jobs-email', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/ai-features', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-web/deploy', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop/first-app', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop/crud-data', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop/custom-ui', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop/export', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-desktop/build', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile/first-app', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile/auth-navigation', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile/api-offline', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile/notifications', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-mobile/build-store', priority: 0.7, changeFrequency: 'monthly' as const },

    // Standalone Courses
    { path: '/courses/batteries', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/api-masterclass', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/mobile-fitness-app', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/ecommerce-spa', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/api-docs-scalar', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/security-deep-dive', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/pulse-analytics', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/gorm-studio', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/react-vite-go', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/deployment-guide', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/saas-with-ai', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/invoice-desktop', priority: 0.7, changeFrequency: 'monthly' as const },

    { path: '/courses/testing', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/gorm-mastery', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/realtime-chat', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/stripe-payments', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/blog-cms', priority: 0.8, changeFrequency: 'monthly' as const },
    { path: '/courses/cicd-github', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/custom-middleware', priority: 0.7, changeFrequency: 'monthly' as const },
    { path: '/courses/grit-ui-components', priority: 0.7, changeFrequency: 'monthly' as const },

    // Course (legacy single course page)
    { path: '/course', priority: 0.5, changeFrequency: 'monthly' as const },
  ]

  // Learning Paths — auto-generate routes for every course, chapter, lesson,
  // and assignment from courses.ts so a sitemap entry exists for each.
  for (const course of COURSES) {
    routes.push({ path: `/courses/${course.slug}`, priority: 0.9, changeFrequency: 'weekly' as const })
    for (const chapter of course.chapters) {
      routes.push({
        path: `/courses/${course.slug}/${chapter.slug}`,
        priority: 0.7,
        changeFrequency: 'monthly' as const,
      })
      for (const module of chapter.modules) {
        for (const lesson of module.lessons) {
          routes.push({
            path: `/courses/${course.slug}/${chapter.slug}/${lesson.slug}`,
            priority: lesson.status === 'available' ? 0.7 : 0.4,
            changeFrequency: 'monthly' as const,
          })
        }
      }
      if (chapter.assignment) {
        routes.push({
          path: `/courses/${course.slug}/${chapter.slug}/assignment`,
          priority: 0.6,
          changeFrequency: 'monthly' as const,
        })
      }
    }
  }

  // Blog — one entry per published post.
  for (const post of getAllPosts()) {
    routes.push({ path: `/blog/${post.slug}`, priority: 0.7, changeFrequency: 'weekly' as const })
  }

  return routes.map((route) => ({
    url: `${baseUrl}${route.path}`,
    lastModified: new Date(),
    changeFrequency: route.changeFrequency,
    priority: route.priority,
  }))
}
