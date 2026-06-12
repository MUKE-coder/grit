import Link from 'next/link'
import { notFound } from 'next/navigation'
import type { Metadata } from 'next'
import {
  ArrowRight,
  CheckCircle2,
  Clock,
  GraduationCap,
  ListChecks,
  Trophy,
  Users,
  Sparkles,
  BookOpen,
  Construction,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { findCourse, flatLessons, COURSES, courseTotalMinutes } from '@/config/courses'
import { cn } from '@/lib/utils'
import { NoAITip } from '@/components/course/no-ai-tip'

interface PageProps {
  params: Promise<{ course: string }>
}

export async function generateStaticParams() {
  return COURSES.map((c) => ({ course: c.slug }))
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const { course: slug } = await params
  const course = findCourse(slug)
  if (!course) return {}
  return {
    title: `${course.name} — Grit Courses`,
    description: course.description,
    openGraph: {
      title: course.name,
      description: course.description,
      url: `https://gritframework.dev/courses/${course.slug}`,
      type: 'article',
    },
  }
}

export default async function CourseLandingPage({ params }: PageProps) {
  const { course: slug } = await params
  const course = findCourse(slug)
  if (!course) notFound()

  const totalLessons = flatLessons(course).length
  const totalMinutes = courseTotalMinutes(course)
  const totalHours = Math.round(totalMinutes / 60)
  const totalAssignments = course.chapters.filter((c) => c.assignment).length

  const firstLesson = flatLessons(course)[0]
  const startHref = firstLesson
    ? `/courses/${course.slug}/${firstLesson.chapter.slug}/${firstLesson.lesson.slug}`
    : `/courses/${course.slug}`

  const levelColor = {
    'absolute-beginner': 'text-emerald-500 bg-emerald-500/10 border-emerald-500/20',
    beginner: 'text-sky-500 bg-sky-500/10 border-sky-500/20',
    intermediate: 'text-violet-500 bg-violet-500/10 border-violet-500/20',
    advanced: 'text-rose-500 bg-rose-500/10 border-rose-500/20',
  }[course.level]

  return (
    <div className="max-w-4xl mx-auto py-10 px-6 lg:px-10">
      {/* Hero */}
      <div className={cn('relative rounded-3xl border border-border/40 overflow-hidden mb-12 bg-gradient-to-br', course.accent)}>
        <div className="absolute inset-0 bg-background/50" />
        <div className="relative p-8 md:p-12">
          <div className="flex items-center gap-2.5 mb-4 flex-wrap">
            <span className="inline-flex items-center gap-1 text-[11px] font-mono text-primary uppercase tracking-wider">
              <BookOpen className="h-3 w-3" />
              Learning Path
            </span>
            <span className={cn('inline-flex items-center text-[10px] font-mono uppercase tracking-wider rounded-full px-2 py-0.5 border', levelColor)}>
              {course.level.replace('-', ' ')}
            </span>
            {course.status === 'coming-soon' && (
              <span className="inline-flex items-center gap-1 text-[10px] font-mono uppercase tracking-wider text-amber-500 bg-amber-500/10 border border-amber-500/20 rounded-full px-2 py-0.5">
                <Construction className="h-2.5 w-2.5" />
                In production
              </span>
            )}
          </div>
          <div className="flex items-start gap-4 mb-4">
            <span className="text-5xl md:text-6xl" aria-hidden="true">
              {course.emoji}
            </span>
            <div className="flex-1">
              <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-2">
                {course.name}
              </h1>
              <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
                {course.tagline}
              </p>
            </div>
          </div>

          {/* Quick stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
            <Stat icon={Clock} label="Time" value={`~${totalHours} h`} />
            <Stat icon={ListChecks} label="Lessons" value={String(totalLessons)} />
            <Stat icon={BookOpen} label="Chapters" value={String(course.chapters.length)} />
            <Stat icon={Trophy} label="Assignments" value={String(totalAssignments)} />
          </div>

          {/* CTA */}
          <div className="flex flex-wrap gap-2.5">
            <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground" disabled={course.status === 'coming-soon'}>
              <Link href={startHref}>
                {course.status === 'coming-soon' ? 'Preview the outline' : 'Start course'}
                <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
              </Link>
            </Button>
            <Button asChild variant="outline" className="rounded-full border-border/60">
              <Link href="/courses">All courses</Link>
            </Button>
          </div>
        </div>
      </div>

      {/* No-AI reminder — the goal is to learn, not ship fastest. */}
      <NoAITip />

      {/* What you'll build */}
      <div className="mb-12">
        <h2 className="text-2xl font-semibold tracking-tight mb-2 flex items-center gap-2">
          <Sparkles className="h-5 w-5 text-primary" />
          What you&apos;ll build
        </h2>
        <p className="text-foreground/85 leading-relaxed">{course.whatYoullBuild}</p>
      </div>

      {/* What you'll learn */}
      <div className="mb-12">
        <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
          <GraduationCap className="h-5 w-5 text-primary" />
          What you&apos;ll learn
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-2.5">
          {course.whatYoullLearn.map((item) => (
            <div key={item} className="flex items-start gap-2 text-sm text-foreground/85">
              <CheckCircle2 className="h-4 w-4 text-emerald-500 mt-0.5 shrink-0" />
              <span>{item}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Chapter outline */}
      <div className="mb-12">
        <h2 className="text-2xl font-semibold tracking-tight mb-4 flex items-center gap-2">
          <BookOpen className="h-5 w-5 text-primary" />
          Course outline
        </h2>
        <div className="space-y-3">
          {course.chapters.map((ch) => {
            const chapterLessons = ch.modules.flatMap((m) => m.lessons)
            const chapterMinutes = chapterLessons.reduce((acc, l) => acc + l.minutes, 0)
            return (
              <Link
                key={ch.slug}
                href={`/courses/${course.slug}/${ch.slug}`}
                className="group block rounded-2xl border border-border/50 bg-card/30 p-5 hover:border-primary/40 hover:bg-card/60 transition-all"
              >
                <div className="flex items-start gap-4">
                  <div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center font-mono text-sm font-semibold shrink-0">
                    {ch.number}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors mb-1">
                      {ch.title}
                    </h3>
                    <p className="text-sm text-muted-foreground mb-3 leading-relaxed">{ch.tagline}</p>
                    <div className="flex items-center gap-3 flex-wrap text-[11px] text-muted-foreground/80 font-mono">
                      <span className="inline-flex items-center gap-1">
                        <ListChecks className="h-3 w-3" />
                        {chapterLessons.length} lessons
                      </span>
                      <span className="inline-flex items-center gap-1">
                        <Clock className="h-3 w-3" />
                        {chapterMinutes}m
                      </span>
                      {ch.assignment && (
                        <span className="inline-flex items-center gap-1 text-amber-500">
                          <Trophy className="h-3 w-3" />
                          Assignment
                        </span>
                      )}
                    </div>
                  </div>
                  <ArrowRight className="h-4 w-4 text-muted-foreground/60 group-hover:text-primary group-hover:translate-x-0.5 transition-all shrink-0 mt-3" />
                </div>
              </Link>
            )
          })}
        </div>
      </div>

      {/* Prereqs + who this is for */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-12">
        <div className="rounded-2xl border border-border/50 bg-card/30 p-5">
          <h3 className="text-base font-semibold mb-3 flex items-center gap-2">
            <ListChecks className="h-4 w-4 text-primary" />
            Prerequisites
          </h3>
          <ul className="space-y-2 text-sm text-foreground/85">
            {course.prerequisites.map((p) => (
              <li key={p} className="flex items-start gap-2">
                <span className="text-primary mt-0.5">›</span>
                <span>{p}</span>
              </li>
            ))}
          </ul>
        </div>
        <div className="rounded-2xl border border-border/50 bg-card/30 p-5">
          <h3 className="text-base font-semibold mb-3 flex items-center gap-2">
            <Users className="h-4 w-4 text-primary" />
            Who this is for
          </h3>
          <ul className="space-y-2 text-sm text-foreground/85">
            {course.whoThisIsFor.map((p) => (
              <li key={p} className="flex items-start gap-2">
                <span className="text-primary mt-0.5">›</span>
                <span>{p}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Start CTA */}
      <div className="rounded-2xl border border-primary/20 bg-primary/[0.05] p-6 text-center">
        <h3 className="text-xl font-semibold mb-2">Ready to start?</h3>
        <p className="text-sm text-muted-foreground mb-5 max-w-xl mx-auto">
          {course.status === 'coming-soon'
            ? "The structure's done — content lands chapter by chapter. Start with Grit Concepts to nail the foundation in the meantime."
            : "Lesson 1 takes ~6 minutes. By the end of this hour you'll be writing real code."}
        </p>
        <div className="flex flex-wrap justify-center gap-2.5">
          <Button asChild className="rounded-full bg-primary hover:bg-primary/90 text-primary-foreground">
            <Link href={startHref}>
              {course.status === 'coming-soon' ? 'Preview outline' : 'Start course'}
              <ArrowRight className="ml-1.5 h-3.5 w-3.5" />
            </Link>
          </Button>
        </div>
      </div>
    </div>
  )
}

function Stat({ icon: Icon, label, value }: { icon: React.ComponentType<{ className?: string }>; label: string; value: string }) {
  return (
    <div className="rounded-xl border border-border/40 bg-background/40 backdrop-blur p-3">
      <div className="flex items-center gap-1.5 text-muted-foreground text-[10px] font-mono uppercase tracking-wider mb-1">
        <Icon className="h-3 w-3" />
        {label}
      </div>
      <p className="text-base font-semibold text-foreground">{value}</p>
    </div>
  )
}
