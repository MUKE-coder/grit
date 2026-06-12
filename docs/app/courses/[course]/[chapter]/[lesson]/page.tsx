import { notFound } from 'next/navigation'
import type { Metadata } from 'next'
import { LessonShell } from '@/components/course/lesson-shell'
import { LessonComingSoon } from '@/components/course/lesson-coming-soon'
import { findLesson, COURSES } from '@/config/courses'
import { loadLessonContent } from '@/lib/lesson-content'

interface PageProps {
  params: Promise<{ course: string; chapter: string; lesson: string }>
}

export async function generateStaticParams() {
  const out: Array<{ course: string; chapter: string; lesson: string }> = []
  for (const c of COURSES) {
    for (const ch of c.chapters) {
      for (const mod of ch.modules) {
        for (const l of mod.lessons) {
          out.push({ course: c.slug, chapter: ch.slug, lesson: l.slug })
        }
      }
    }
  }
  return out
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { course: courseSlug, chapter: chapterSlug, lesson: lessonSlug } = await params
  const found = findLesson(courseSlug, chapterSlug, lessonSlug)
  if (!found) return {}
  return {
    title: `${found.lesson.title} — ${found.course.shortName}`,
    description: found.lesson.tagline,
    openGraph: {
      title: `${found.lesson.title} — ${found.course.name}`,
      description: found.lesson.tagline,
      url: `https://gritframework.dev/courses/${courseSlug}/${chapterSlug}/${lessonSlug}`,
      type: 'article',
    },
  }
}

export default async function LessonPage({ params }: PageProps) {
  const { course: courseSlug, chapter: chapterSlug, lesson: lessonSlug } = await params
  const found = findLesson(courseSlug, chapterSlug, lessonSlug)
  if (!found) notFound()
  const { course, chapter, lesson, prev, next } = found

  const Content = await loadLessonContent(courseSlug, chapterSlug, lessonSlug)

  return (
    <LessonShell course={course} chapter={chapter} lesson={lesson} prev={prev} next={next}>
      {Content ? (
        <Content />
      ) : (
        <LessonComingSoon
          courseSlug={course.slug}
          courseName={course.name}
          lessonTitle={lesson.title}
        />
      )}
    </LessonShell>
  )
}
