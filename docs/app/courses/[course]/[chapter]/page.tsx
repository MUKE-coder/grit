import Link from 'next/link'
import { notFound } from 'next/navigation'
import type { Metadata } from 'next'
import {
  ArrowRight,
  ArrowLeft,
  CheckCircle2,
  Clock,
  Trophy,
  BookOpen,
  Target,
  Sparkles,
  Zap,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { findChapter, findCourse, COURSES } from '@/config/courses'
import { cn } from '@/lib/utils'

interface PageProps {
  params: Promise<{ course: string; chapter: string }>
}

export async function generateStaticParams() {
  const out: Array<{ course: string; chapter: string }> = []
  for (const c of COURSES) {
    for (const ch of c.chapters) {
      out.push({ course: c.slug, chapter: ch.slug })
    }
  }
  return out
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { course: courseSlug, chapter: chapterSlug } = await params
  const found = findChapter(courseSlug, chapterSlug)
  if (!found) return {}
  return {
    title: `${found.chapter.title} — ${found.course.shortName}`,
    description: found.chapter.tagline,
    openGraph: {
      title: `${found.chapter.title} — ${found.course.name}`,
      description: found.chapter.tagline,
      url: `https://gritframework.dev/courses/${found.course.slug}/${found.chapter.slug}`,
      type: 'article',
    },
  }
}

export default async function ChapterLandingPage({ params }: PageProps) {
  const { course: courseSlug, chapter: chapterSlug } = await params
  const found = findChapter(courseSlug, chapterSlug)
  if (!found) notFound()
  const { course, chapter } = found

  // Find prev / next chapters
  const chapters = course.chapters
  const idx = chapters.findIndex((c) => c.slug === chapter.slug)
  const prevChapter = idx > 0 ? chapters[idx - 1] : null
  const nextChapter = idx < chapters.length - 1 ? chapters[idx + 1] : null

  const allLessons = chapter.modules.flatMap((m) => m.lessons)
  const firstLesson = allLessons[0]
  const totalMinutes = allLessons.reduce((acc, l) => acc + l.minutes, 0)

  return (
    <div className="max-w-3xl mx-auto py-10 px-6 lg:px-10">
      {/* Breadcrumb */}
      <nav className="flex items-center gap-1.5 text-xs text-muted-foreground/80 mb-6">
        <Link href="/courses" className="hover:text-foreground">
          Courses
        </Link>
        <span className="opacity-50">/</span>
        <Link href={`/courses/${course.slug}`} className="hover:text-foreground">
          {course.shortName}
        </Link>
      </nav>

      {/* Hero */}
      <div className="mb-10">
        <p className="text-[11px] uppercase tracking-wider text-primary font-mono mb-2 inline-flex items-center gap-1.5">
          <BookOpen className="h-3 w-3" />
          Chapter {chapter.number}
        </p>
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-3">
          {chapter.title}
        </h1>
        <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
          {chapter.tagline}
        </p>
        <div className="flex items-center gap-2 mt-4 flex-wrap text-[11px] font-mono text-muted-foreground">
          <span className="inline-flex items-center gap-1 rounded-full border border-border/40 px-2 py-0.5">
            <Clock className="h-3 w-3" />
            ~{totalMinutes} min total
          </span>
          <span className="inline-flex items-center gap-1 rounded-full border border-border/40 px-2 py-0.5">
            {allLessons.length} lessons
          </span>
          {chapter.assignment && (
            <span className="inline-flex items-center gap-1 rounded-full border border-amber-500/30 text-amber-500 bg-amber-500/5 px-2 py-0.5">
              <Trophy className="h-3 w-3" />
              Assignment
            </span>
          )}
        </div>
      </div>

      {/* Learning goals */}
      <div className="mb-10 rounded-2xl border border-primary/20 bg-primary/[0.04] p-5">
        <h2 className="text-sm font-semibold mb-3 flex items-center gap-2 text-primary">
          <Target className="h-4 w-4" />
          By the end of this chapter you&apos;ll be able to
        </h2>
        <ul className="space-y-1.5 text-sm text-foreground/90">
          {chapter.learningGoals.map((g) => (
            <li key={g} className="flex items-start gap-2">
              <CheckCircle2 className="h-3.5 w-3.5 text-emerald-500 mt-0.5 shrink-0" />
              <span>{g}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Modules + lessons */}
      <div className="mb-10">
        <h2 className="text-xl font-semibold tracking-tight mb-4 flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-primary" />
          Lessons
        </h2>
        <div className="space-y-5">
          {chapter.modules.map((mod) => (
            <div key={mod.title}>
              <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-mono mb-2 pl-1">
                Module: {mod.title}
              </p>
              <div className="space-y-2">
                {mod.lessons.map((l, i) => {
                  const diffColor = {
                    easy: 'text-emerald-500',
                    medium: 'text-amber-500',
                    hard: 'text-rose-500',
                  }[l.difficulty]
                  return (
                    <Link
                      key={l.slug}
                      href={`/courses/${course.slug}/${chapter.slug}/${l.slug}`}
                      className="group block rounded-xl border border-border/50 bg-card/30 p-4 hover:border-primary/40 hover:bg-card/60 transition-all"
                    >
                      <div className="flex items-start gap-3">
                        <span className="text-[10px] font-mono text-muted-foreground/70 mt-1 w-6 shrink-0">
                          {String(i + 1).padStart(2, '0')}
                        </span>
                        <div className="flex-1 min-w-0">
                          <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors mb-0.5">
                            {l.title}
                          </h3>
                          <p className="text-xs text-muted-foreground leading-relaxed">
                            {l.tagline}
                          </p>
                          <div className="mt-2 flex items-center gap-2 text-[10px] font-mono text-muted-foreground/70">
                            <span className="inline-flex items-center gap-0.5">
                              <Clock className="h-2.5 w-2.5" />
                              {l.minutes}m
                            </span>
                            <span className={cn('inline-flex items-center gap-0.5', diffColor)}>
                              <Zap className="h-2.5 w-2.5" />
                              {l.difficulty}
                            </span>
                          </div>
                        </div>
                        <ArrowRight className="h-3.5 w-3.5 text-muted-foreground/60 group-hover:text-primary group-hover:translate-x-0.5 transition-all shrink-0 mt-1" />
                      </div>
                    </Link>
                  )
                })}
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Assignment preview */}
      {chapter.assignment && (
        <div className="mb-10 rounded-2xl border border-amber-500/30 bg-amber-500/[0.04] p-5">
          <div className="flex items-center gap-2 mb-2">
            <Trophy className="h-4 w-4 text-amber-500" />
            <p className="text-[10px] uppercase tracking-wider text-amber-500 font-mono">
              Chapter assignment
            </p>
          </div>
          <h3 className="text-base font-semibold mb-2">{chapter.assignment.title}</h3>
          <p className="text-sm text-foreground/85 leading-relaxed mb-3">
            {chapter.assignment.brief}
          </p>
          <Link
            href={`/courses/${course.slug}/${chapter.slug}/assignment`}
            className="text-xs font-medium text-amber-500 hover:underline inline-flex items-center gap-1"
          >
            See success criteria
            <ArrowRight className="h-3 w-3" />
          </Link>
        </div>
      )}

      {/* Start CTA */}
      {firstLesson && (
        <div className="rounded-2xl border border-primary/20 bg-primary/[0.05] p-5 mb-10 text-center">
          <p className="text-sm text-muted-foreground mb-4">Lesson 1 takes ~{firstLesson.minutes} min.</p>
          <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
            <Link href={`/courses/${course.slug}/${chapter.slug}/${firstLesson.slug}`}>
              Start chapter
              <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
            </Link>
          </Button>
        </div>
      )}

      {/* Prev / Next chapter nav */}
      <div className="flex items-center justify-between pt-6 border-t border-border/30">
        {prevChapter ? (
          <Button variant="ghost" size="sm" asChild className="text-muted-foreground/80 hover:text-foreground">
            <Link href={`/courses/${course.slug}/${prevChapter.slug}`} className="gap-1.5">
              <ArrowLeft className="h-3.5 w-3.5" />
              ch.{prevChapter.number}: {prevChapter.title}
            </Link>
          </Button>
        ) : (
          <span />
        )}
        {nextChapter && (
          <Button variant="ghost" size="sm" asChild className="text-muted-foreground/80 hover:text-foreground">
            <Link href={`/courses/${course.slug}/${nextChapter.slug}`} className="gap-1.5">
              ch.{nextChapter.number}: {nextChapter.title}
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>
          </Button>
        )}
      </div>
    </div>
  )
}
