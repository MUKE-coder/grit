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
