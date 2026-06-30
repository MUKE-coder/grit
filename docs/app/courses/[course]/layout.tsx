import { notFound } from 'next/navigation'
import { SiteHeader } from '@/components/site-header'
import { CourseSidebar, CourseSidebarMobile } from '@/components/course/course-sidebar'
import { findCourse, COURSES } from '@/config/courses'

interface Props {
  params: Promise<{ course: string }>
  children: React.ReactNode
}

// Generate static params so Next.js prerenders all course routes.
export async function generateStaticParams() {
  return COURSES.map((c) => ({ course: c.slug }))
}

/**
 * Shell for every page inside /courses/<course>/* — the course-level
 * landing, every chapter landing, and every lesson. Wraps content
 * with the SiteHeader (for nav) and the CourseSidebar (for in-course
 * navigation). The mobile sidebar is a drawer triggered by a button
 * in the lesson breadcrumbs row.
 */
export default async function CourseLayout({ params, children }: Props) {
  const { course: courseSlug } = await params
  const course = findCourse(courseSlug)
  if (!course) notFound()

  return (
    <div className="min-h-screen bg-background flex flex-col">
      {/* Shared grid backdrop — consistent texture across every in-course page */}
      <div aria-hidden className="pointer-events-none fixed inset-0 -z-10 bg-grit-grid-sm opacity-70 mask-fade-b" />
      <SiteHeader />
      <div className="flex flex-1 min-h-0">
        <CourseSidebar course={course} />
        <main className="flex-1 min-w-0 overflow-x-hidden">
          <div className="lg:hidden px-4 pt-3">
            <CourseSidebarMobile course={course} />
          </div>
          {children}
        </main>
      </div>
    </div>
  )
}
