import Link from 'next/link'
import { notFound } from 'next/navigation'
import type { Metadata } from 'next'
import {
  ArrowLeft,
  ArrowRight,
  Trophy,
  CheckCircle2,
  ListChecks,
  Home,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { findChapter, COURSES } from '@/config/courses'

interface PageProps {
  params: Promise<{ course: string; chapter: string }>
}

export async function generateStaticParams() {
  const out: Array<{ course: string; chapter: string }> = []
  for (const c of COURSES) {
    for (const ch of c.chapters) {
      if (ch.assignment) out.push({ course: c.slug, chapter: ch.slug })
    }
  }
  return out
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { course, chapter } = await params
  const found = findChapter(course, chapter)
  if (!found || !found.chapter.assignment) return {}
  return {
    title: `Assignment: ${found.chapter.assignment.title} — ${found.course.shortName}`,
    description: found.chapter.assignment.brief.slice(0, 160),
  }
}

export default async function AssignmentPage({ params }: PageProps) {
  const { course: courseSlug, chapter: chapterSlug } = await params
  const found = findChapter(courseSlug, chapterSlug)
  if (!found || !found.chapter.assignment) notFound()
  const { course, chapter } = found
  const a = chapter.assignment!

  // Find the next chapter for "Continue" CTA
  const idx = course.chapters.findIndex((c) => c.slug === chapter.slug)
  const nextChapter = idx < course.chapters.length - 1 ? course.chapters[idx + 1] : null

  return (
    <article className="max-w-3xl mx-auto py-10 px-6 lg:px-10">
      {/* Breadcrumb */}
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

      {/* Hero */}
      <div className="mb-10">
        <p className="text-[11px] uppercase tracking-wider text-amber-500 font-mono mb-2 inline-flex items-center gap-1.5">
          <Trophy className="h-3 w-3" />
          Chapter {chapter.number} Assignment
        </p>
        <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-3">
          {a.title}
        </h1>
      </div>

      {/* Brief */}
      <div className="mb-8 rounded-2xl border border-amber-500/30 bg-amber-500/[0.04] p-5">
        <p className="text-[10px] uppercase tracking-wider text-amber-500 font-mono mb-2">
          The brief
        </p>
        <p className="text-base text-foreground/90 leading-relaxed">{a.brief}</p>
      </div>

      {/* Success criteria */}
      <div className="mb-10">
        <h2 className="text-xl font-semibold tracking-tight mb-4 flex items-center gap-2">
          <ListChecks className="h-5 w-5 text-primary" />
          You&apos;ve completed this when
        </h2>
        <ul className="space-y-2.5">
          {a.successCriteria.map((c) => (
            <li
              key={c}
              className="flex items-start gap-2.5 rounded-xl border border-border/50 bg-card/30 p-3.5"
            >
              <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5 shrink-0" />
              <span className="text-sm text-foreground/90 leading-relaxed">{c}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Done CTA */}
      <div className="rounded-2xl border border-primary/20 bg-primary/[0.05] p-5 mb-10 text-center">
        <p className="text-sm text-foreground/85 mb-2">Worked through every criterion?</p>
        <p className="text-xs text-muted-foreground mb-4">
          Push your code to GitHub, paste the link in your notes.md, and move on.
        </p>
        {nextChapter && (
          <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
            <Link href={`/courses/${course.slug}/${nextChapter.slug}`}>
              Continue to ch.{nextChapter.number}: {nextChapter.title}
              <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
            </Link>
          </Button>
        )}
        {!nextChapter && (
          <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
            <Link href={`/courses/${course.slug}`}>
              Finish course
              <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
            </Link>
          </Button>
        )}
      </div>

      {/* Nav */}
      <div className="flex items-center justify-between pt-6 border-t border-border/30">
        <Button
          variant="ghost"
          size="sm"
          asChild
          className="text-muted-foreground/80 hover:text-foreground"
        >
          <Link href={`/courses/${course.slug}/${chapter.slug}`} className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            Back to chapter
          </Link>
        </Button>
      </div>
    </article>
  )
}
