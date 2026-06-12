/**
 * Lesson content registry.
 *
 * Each lesson's actual TSX body lives at
 *   /content/lessons/<course>/<chapter>/<lesson>.tsx
 * and exports a default React component (no props).
 *
 * loadLessonContent dynamically imports it; missing files return null
 * so the dynamic lesson page can fall through to <LessonComingSoon />.
 *
 * This means I never have to keep an explicit registry in sync — drop
 * a new lesson TSX in the right place and it just works.
 */

type LessonComponent = React.ComponentType

const moduleMap: Record<string, () => Promise<{ default: LessonComponent }>> = {
  /* ─── Grit Concepts course ─────────────────────────── */
  'concepts/welcome/what-is-grit': () => import('@/content/lessons/concepts/welcome/what-is-grit'),
  'concepts/welcome/when-to-use-grit': () => import('@/content/lessons/concepts/welcome/when-to-use-grit'),
  'concepts/welcome/install-grit': () => import('@/content/lessons/concepts/welcome/install-grit'),
  'concepts/welcome/verify-install': () => import('@/content/lessons/concepts/welcome/verify-install'),

  'concepts/first-project/grit-new': () => import('@/content/lessons/concepts/first-project/grit-new'),
  'concepts/first-project/project-tour': () => import('@/content/lessons/concepts/first-project/project-tour'),
  'concepts/first-project/dev-servers': () => import('@/content/lessons/concepts/first-project/dev-servers'),
  'concepts/first-project/first-look': () => import('@/content/lessons/concepts/first-project/first-look'),

  'concepts/conventions/folder-conventions': () => import('@/content/lessons/concepts/conventions/folder-conventions'),
  'concepts/conventions/naming-conventions': () => import('@/content/lessons/concepts/conventions/naming-conventions'),
  'concepts/conventions/api-response-format': () => import('@/content/lessons/concepts/conventions/api-response-format'),
  'concepts/conventions/error-handling': () => import('@/content/lessons/concepts/conventions/error-handling'),

  'concepts/generators/what-is-a-resource': () => import('@/content/lessons/concepts/generators/what-is-a-resource'),
  'concepts/generators/grit-generate': () => import('@/content/lessons/concepts/generators/grit-generate'),
  'concepts/generators/what-got-generated': () => import('@/content/lessons/concepts/generators/what-got-generated'),
  'concepts/generators/grit-sync': () => import('@/content/lessons/concepts/generators/grit-sync'),

  'concepts/architecture-modes/single-mode': () => import('@/content/lessons/concepts/architecture-modes/single-mode'),
  'concepts/architecture-modes/triple-mode': () => import('@/content/lessons/concepts/architecture-modes/triple-mode'),
  'concepts/architecture-modes/specialized-modes': () => import('@/content/lessons/concepts/architecture-modes/specialized-modes'),
  'concepts/architecture-modes/choosing-a-kit': () => import('@/content/lessons/concepts/architecture-modes/choosing-a-kit'),

  /* ─── Building a Go API course ─────────────────────── */
  'go-api/scaffold-tour/scaffold': () => import('@/content/lessons/go-api/scaffold-tour/scaffold'),
  'go-api/scaffold-tour/project-tour': () => import('@/content/lessons/go-api/scaffold-tour/project-tour'),
  'go-api/scaffold-tour/first-request': () => import('@/content/lessons/go-api/scaffold-tour/first-request'),

  'go-api/models/gorm-basics': () => import('@/content/lessons/go-api/models/gorm-basics'),
  'go-api/models/relations': () => import('@/content/lessons/go-api/models/relations'),
  'go-api/models/migrations': () => import('@/content/lessons/go-api/models/migrations'),

  'go-api/auth/jwt': () => import('@/content/lessons/go-api/auth/jwt'),
  'go-api/auth/oauth': () => import('@/content/lessons/go-api/auth/oauth'),
  'go-api/auth/totp': () => import('@/content/lessons/go-api/auth/totp'),
  'go-api/auth/rbac': () => import('@/content/lessons/go-api/auth/rbac'),

  'go-api/batteries/jobs': () => import('@/content/lessons/go-api/batteries/jobs'),
  'go-api/batteries/mail': () => import('@/content/lessons/go-api/batteries/mail'),
  'go-api/batteries/storage': () => import('@/content/lessons/go-api/batteries/storage'),
  'go-api/batteries/ai': () => import('@/content/lessons/go-api/batteries/ai'),

  'go-api/security-observability/sentinel': () => import('@/content/lessons/go-api/security-observability/sentinel'),
  'go-api/security-observability/pulse': () => import('@/content/lessons/go-api/security-observability/pulse'),
  'go-api/security-observability/audit-log': () => import('@/content/lessons/go-api/security-observability/audit-log'),

  'go-api/deploy/grit-deploy': () => import('@/content/lessons/go-api/deploy/grit-deploy'),
  'go-api/deploy/env-config': () => import('@/content/lessons/go-api/deploy/env-config'),

  /* ─── Mobile course ──────────────────────────────────── */
  'mobile/scaffold/scaffold': () => import('@/content/lessons/mobile/scaffold/scaffold'),
  'mobile/scaffold/expo-tour': () => import('@/content/lessons/mobile/scaffold/expo-tour'),
  'mobile/scaffold/first-run': () => import('@/content/lessons/mobile/scaffold/first-run'),

  'mobile/shared-types/grit-sync-mobile': () => import('@/content/lessons/mobile/shared-types/grit-sync-mobile'),
  'mobile/shared-types/api-client': () => import('@/content/lessons/mobile/shared-types/api-client'),

  'mobile/auth/login-ui': () => import('@/content/lessons/mobile/auth/login-ui'),
  'mobile/auth/secure-storage': () => import('@/content/lessons/mobile/auth/secure-storage'),
  'mobile/auth/refresh': () => import('@/content/lessons/mobile/auth/refresh'),

  'mobile/push-notifications/register': () => import('@/content/lessons/mobile/push-notifications/register'),
  'mobile/push-notifications/send-from-api': () => import('@/content/lessons/mobile/push-notifications/send-from-api'),

  'mobile/ship/eas-build': () => import('@/content/lessons/mobile/ship/eas-build'),
  'mobile/ship/submit': () => import('@/content/lessons/mobile/ship/submit'),
  'mobile/ship/ota': () => import('@/content/lessons/mobile/ship/ota'),

  /* ─── Web (Next.js) course ──────────────────────────── */
  'web-nextjs/scaffold/scaffold': () => import('@/content/lessons/web-nextjs/scaffold/scaffold'),
  'web-nextjs/scaffold/tour': () => import('@/content/lessons/web-nextjs/scaffold/tour'),
  'web-nextjs/scaffold/shared-package': () => import('@/content/lessons/web-nextjs/scaffold/shared-package'),

  'web-nextjs/public-site/landing': () => import('@/content/lessons/web-nextjs/public-site/landing'),
  'web-nextjs/public-site/seo': () => import('@/content/lessons/web-nextjs/public-site/seo'),

  'web-nextjs/dashboard/signup': () => import('@/content/lessons/web-nextjs/dashboard/signup'),
  'web-nextjs/dashboard/dashboard-widgets': () => import('@/content/lessons/web-nextjs/dashboard/dashboard-widgets'),

  'web-nextjs/admin-panel/define-resource': () => import('@/content/lessons/web-nextjs/admin-panel/define-resource'),
  'web-nextjs/admin-panel/datatable': () => import('@/content/lessons/web-nextjs/admin-panel/datatable'),
  'web-nextjs/admin-panel/formbuilder': () => import('@/content/lessons/web-nextjs/admin-panel/formbuilder'),

  'web-nextjs/tenants/tenant-models': () => import('@/content/lessons/web-nextjs/tenants/tenant-models'),
  'web-nextjs/tenants/role-gates': () => import('@/content/lessons/web-nextjs/tenants/role-gates'),
  'web-nextjs/tenants/invitations': () => import('@/content/lessons/web-nextjs/tenants/invitations'),

  /* ─── Desktop (Wails) course ──────────────────────── */
  'desktop/scaffold/scaffold': () => import('@/content/lessons/desktop/scaffold/scaffold'),
  'desktop/scaffold/wails-dev': () => import('@/content/lessons/desktop/scaffold/wails-dev'),
  'desktop/scaffold/first-build': () => import('@/content/lessons/desktop/scaffold/first-build'),

  'desktop/offline/sqlite': () => import('@/content/lessons/desktop/offline/sqlite'),
  'desktop/offline/outbox': () => import('@/content/lessons/desktop/offline/outbox'),
  'desktop/offline/sync': () => import('@/content/lessons/desktop/offline/sync'),

  'desktop/frameless/titlebar': () => import('@/content/lessons/desktop/frameless/titlebar'),
  'desktop/frameless/window-controls': () => import('@/content/lessons/desktop/frameless/window-controls'),

  'desktop/auto-update/updater-go': () => import('@/content/lessons/desktop/auto-update/updater-go'),
  'desktop/auto-update/modal-ui': () => import('@/content/lessons/desktop/auto-update/modal-ui'),
  'desktop/auto-update/release-script': () => import('@/content/lessons/desktop/auto-update/release-script'),

  'desktop/installers/project-nsi': () => import('@/content/lessons/desktop/installers/project-nsi'),
  'desktop/installers/project-slim': () => import('@/content/lessons/desktop/installers/project-slim'),
  'desktop/installers/bitmaps': () => import('@/content/lessons/desktop/installers/bitmaps'),

  /* ─── Multi-platform course ──────────────────────── */
  'multiplatform/foundation/scaffold': () => import('@/content/lessons/multiplatform/foundation/scaffold'),
  'multiplatform/foundation/monorepo-wiring': () => import('@/content/lessons/multiplatform/foundation/monorepo-wiring'),

  'multiplatform/shared-types/grit-sync-multi': () => import('@/content/lessons/multiplatform/shared-types/grit-sync-multi'),
  'multiplatform/shared-types/shared-zod': () => import('@/content/lessons/multiplatform/shared-types/shared-zod'),

  'multiplatform/feature-implementation/pick-feature': () => import('@/content/lessons/multiplatform/feature-implementation/pick-feature'),
  'multiplatform/feature-implementation/web-impl': () => import('@/content/lessons/multiplatform/feature-implementation/web-impl'),
  'multiplatform/feature-implementation/mobile-impl': () => import('@/content/lessons/multiplatform/feature-implementation/mobile-impl'),
  'multiplatform/feature-implementation/desktop-impl': () => import('@/content/lessons/multiplatform/feature-implementation/desktop-impl'),

  'multiplatform/sync/mobile-offline': () => import('@/content/lessons/multiplatform/sync/mobile-offline'),
  'multiplatform/sync/desktop-outbox': () => import('@/content/lessons/multiplatform/sync/desktop-outbox'),
  'multiplatform/sync/conflicts': () => import('@/content/lessons/multiplatform/sync/conflicts'),

  'multiplatform/releases/compat-matrix': () => import('@/content/lessons/multiplatform/releases/compat-matrix'),
  'multiplatform/releases/staggered': () => import('@/content/lessons/multiplatform/releases/staggered'),
}

export async function loadLessonContent(
  courseSlug: string,
  chapterSlug: string,
  lessonSlug: string,
): Promise<LessonComponent | null> {
  const key = `${courseSlug}/${chapterSlug}/${lessonSlug}`
  const loader = moduleMap[key]
  if (!loader) return null
  try {
    const mod = await loader()
    return mod.default
  } catch {
    return null
  }
}
