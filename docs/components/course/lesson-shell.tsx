import Link from 'next/link'
import { ArrowLeft, ArrowRight, Clock, Zap, Home, Github, MessageSquarePlus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import type { Course, Chapter, Lesson } from '@/config/courses'

interface Props {
  course: Course
  chapter: Chapter
  lesson: Lesson
  prev?: { chapterSlug: string; lessonSlug: string }
  next?: { chapterSlug: string; lessonSlug: string }
  children: React.ReactNode
}

/**
 * LessonShell — wraps the actual lesson content with breadcrumbs,
 * the title + meta strip (time / difficulty), and the prev/next nav
 * at the bottom. Lesson pages just render their content as children.
 */
export function LessonShell({ course, chapter, lesson, prev, next, children }: Props) {
  const diffColor = {
    easy: 'text-emerald-500 bg-emerald-500/10 border-emerald-500/20',
    medium: 'text-amber-500 bg-amber-500/10 border-amber-500/20',
    hard: 'text-rose-500 bg-rose-500/10 border-rose-500/20',
  }[lesson.difficulty]

  return (
    <article className="max-w-3xl mx-auto py-10 px-6 lg:px-10">
      {/* Breadcrumbs */}
      <nav className="flex items-center gap-1.5 text-xs text-muted-foreground/80 mb-6 overflow-x-auto whitespace-nowrap">
        <Link href="/courses" className="hover:text-foreground inline-flex items-center gap-1">
          <Home className="h-3 w-3" />
          Courses
        </Link>
        <span className="opacity-50">/</span>
        <Link href={`/courses/${course.slug}`} className="hover:text-foreground">
          {course.shortName}
        </Link>
        <span className="opacity-50">/</span>
        <Link
          href={`/courses/${course.slug}/${chapter.slug}`}
          className="hover:text-foreground"
        >
          ch.{chapter.number}: {chapter.title}
        </Link>
      </nav>

      {/* Title + meta */}
      <div className="mb-8">
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-3">
          {lesson.title}
        </h1>
        <p className="text-base text-muted-foreground leading-relaxed mb-4">
          {lesson.tagline}
        </p>
        <div className="flex items-center gap-2 flex-wrap">
          <span className="inline-flex items-center gap-1 text-[11px] font-mono text-muted-foreground bg-card/50 border border-border/40 rounded-full px-2 py-0.5">
            <Clock className="h-3 w-3" />
            {lesson.minutes} min
          </span>
          <span className={cn('inline-flex items-center gap-1 text-[11px] font-mono rounded-full px-2 py-0.5 border', diffColor)}>
            <Zap className="h-3 w-3" />
            {lesson.difficulty}
          </span>
        </div>
      </div>

      {/* Lesson content */}
      <div className="prose-grit max-w-none">{children}</div>

      {/* Suggest improvement — GitHub issue, lesson-scoped. */}
      <div className="mt-12 rounded-2xl border border-border/60 bg-card/30 p-5">
        <div className="flex items-start gap-4">
          <div className="h-10 w-10 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center shrink-0">
            <MessageSquarePlus className="h-5 w-5 text-primary" />
          </div>
          <div className="flex-1 min-w-0">
            <h3 className="text-base font-semibold mb-1">
              Spot a typo? Have an idea?
            </h3>
            <p className="text-sm text-muted-foreground leading-relaxed mb-3">
              Help us improve this lesson. One click opens a GitHub issue with
              the lesson URL pre-filled — suggest clearer wording, report a
              bug, or request more depth. The course keeps improving thanks
              to learners like you.
            </p>
            <a
              href={`https://github.com/MUKE-coder/grit/issues/new?template=course-feedback.yml&title=${encodeURIComponent(
                `[Course feedback] ${course.shortName} · ${lesson.title}`,
              )}&chapter-url=${encodeURIComponent(
                `https://gritframework.dev/courses/${course.slug}/${chapter.slug}/${lesson.slug}`,
              )}`}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-1.5 rounded-full bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-medium px-4 py-2 transition-colors"
            >
              <Github className="h-3.5 w-3.5" />
              Suggest an improvement on GitHub
              <ArrowRight className="h-3 w-3" />
            </a>
          </div>
        </div>
      </div>

      {/* Prev / Next nav */}
      <div className="mt-14 pt-6 border-t border-border/30 flex items-center justify-between gap-4">
        {prev ? (
          <Button variant="ghost" size="sm" asChild className="text-muted-foreground/80 hover:text-foreground">
            <Link href={`/courses/${course.slug}/${prev.chapterSlug}/${prev.lessonSlug}`} className="gap-1.5">
              <ArrowLeft className="h-3.5 w-3.5" />
              Previous lesson
            </Link>
          </Button>
        ) : (
          <span />
        )}
        {next ? (
          <Button variant="ghost" size="sm" asChild className="text-muted-foreground/80 hover:text-foreground">
            <Link href={`/courses/${course.slug}/${next.chapterSlug}/${next.lessonSlug}`} className="gap-1.5">
              Next lesson
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </Button>
        ) : (
          <Button variant="ghost" size="sm" asChild className="text-primary hover:bg-primary/10">
            <Link href={`/courses/${course.slug}`} className="gap-1.5">
              Finish course
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </Button>
        )}
      </div>
    </article>
  )
}
