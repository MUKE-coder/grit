'use client'

import { useState } from 'react'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ChevronRight, ChevronDown, Check, Circle, Clock, FileText, Trophy, Lock, Menu, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Course } from '@/config/courses'

interface Props {
  course: Course
}

/**
 * CourseSidebar — vertical left rail used inside a course view. Shows
 * every chapter (collapsible), every lesson, the current lesson, and
 * the chapter assignment. Auto-expands the chapter containing the
 * current lesson.
 */
export function CourseSidebar({ course }: Props) {
  const pathname = usePathname()
  const currentChapter = course.chapters.find((ch) =>
    pathname.startsWith(`/courses/${course.slug}/${ch.slug}`),
  )

  // Collapsed by default. Open the current chapter if we're inside one;
  // otherwise open only the first chapter (so users always have a starting
  // point visible without having to fight a wall of sections).
  const initialOpen = currentChapter
    ? [currentChapter.slug]
    : course.chapters[0]
    ? [course.chapters[0].slug]
    : []
  const [openChapters, setOpenChapters] = useState<Set<string>>(new Set(initialOpen))

  const toggle = (slug: string) => {
    setOpenChapters((prev) => {
      const next = new Set(prev)
      if (next.has(slug)) next.delete(slug)
      else next.add(slug)
      return next
    })
  }

  return (
    <aside className="hidden lg:block w-72 shrink-0 border-r border-border/40 bg-sidebar-background/50 overflow-y-auto">
      <div className="sticky top-0 px-5 py-4 border-b border-border/40 bg-sidebar-background/95 backdrop-blur z-10">
        <p className="text-[10px] uppercase tracking-wider text-muted-foreground font-mono mb-1">
          Course
        </p>
        <Link
          href={`/courses/${course.slug}`}
          className="text-sm font-semibold leading-tight hover:text-primary line-clamp-2 block"
        >
          {course.emoji} {course.shortName}
        </Link>
      </div>
      <CourseTreeInner course={course} pathname={pathname} openChapters={openChapters} toggle={toggle} />
    </aside>
  )
}

/* ─── Mobile drawer variant ────────────────────────────────────────── */

export function CourseSidebarMobile({ course }: Props) {
  const [open, setOpen] = useState(false)
  const pathname = usePathname()
  const currentChapter = course.chapters.find((ch) =>
    pathname.startsWith(`/courses/${course.slug}/${ch.slug}`),
  )
  // Mobile drawer also defaults to just the active chapter, or the first
  // when on the course landing.
  const initialOpen = currentChapter
    ? [currentChapter.slug]
    : course.chapters[0]
    ? [course.chapters[0].slug]
    : []
  const [openChapters, setOpenChapters] = useState<Set<string>>(new Set(initialOpen))
  const toggle = (slug: string) => {
    setOpenChapters((prev) => {
      const next = new Set(prev)
      if (next.has(slug)) next.delete(slug)
      else next.add(slug)
      return next
    })
  }

  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="lg:hidden inline-flex items-center gap-1.5 rounded-md border border-border/60 px-2.5 py-1.5 text-xs"
      >
        <Menu className="h-3.5 w-3.5" />
        Course menu
      </button>
      {open && (
        <div className="lg:hidden fixed inset-0 z-50 bg-background">
          <div className="flex items-center justify-between px-4 py-3 border-b border-border/40">
            <p className="text-sm font-semibold">
              {course.emoji} {course.shortName}
            </p>
            <button onClick={() => setOpen(false)} className="p-1">
              <X className="h-4 w-4" />
            </button>
          </div>
          <div className="h-[calc(100vh-49px)] overflow-y-auto">
            <CourseTreeInner
              course={course}
              pathname={pathname}
              openChapters={openChapters}
              toggle={toggle}
              onLessonClick={() => setOpen(false)}
            />
          </div>
        </div>
      )}
    </>
  )
}

/* ─── Shared tree renderer ────────────────────────────────────────── */

function CourseTreeInner({
  course,
  pathname,
  openChapters,
  toggle,
  onLessonClick,
}: {
  course: Course
  pathname: string
  openChapters: Set<string>
  toggle: (slug: string) => void
  onLessonClick?: () => void
}) {
  return (
    <nav className="p-3">
      {course.chapters.map((ch) => {
        const open = openChapters.has(ch.slug)
        const isComing = ch.status === 'coming-soon'
        const chapterBase = `/courses/${course.slug}/${ch.slug}`
        return (
          <div key={ch.slug} className="mb-1">
            <button
              type="button"
              onClick={() => toggle(ch.slug)}
              className={cn(
                'w-full flex items-center gap-2 px-2 py-1.5 rounded-md text-left text-[13px] font-medium hover:bg-accent/30',
                pathname === chapterBase && 'bg-primary/10 text-primary',
              )}
            >
              {open ? (
                <ChevronDown className="h-3 w-3 shrink-0 opacity-60" />
              ) : (
                <ChevronRight className="h-3 w-3 shrink-0 opacity-60" />
              )}
              <span className="text-[10px] font-mono text-muted-foreground shrink-0 w-7">
                ch.{ch.number}
              </span>
              <span className="truncate flex-1">{ch.title}</span>
              {isComing && <Lock className="h-3 w-3 text-muted-foreground/60" />}
            </button>
            {open && (
              <div className="ml-5 pl-2 border-l border-border/30 py-1">
                <Link
                  href={chapterBase}
                  onClick={onLessonClick}
                  className={cn(
                    'block px-2 py-1 rounded text-[12px] hover:bg-accent/30',
                    pathname === chapterBase
                      ? 'text-primary font-medium'
                      : 'text-muted-foreground/80',
                  )}
                >
                  Chapter overview
                </Link>
                {ch.modules.map((mod) => (
                  <div key={mod.title} className="mt-1">
                    <p className="px-2 pt-1.5 pb-1 text-[10px] uppercase tracking-wider text-muted-foreground/70 font-mono">
                      {mod.title}
                    </p>
                    {mod.lessons.map((l) => {
                      const href = `${chapterBase}/${l.slug}`
                      const active = pathname === href
                      const done = false // could wire to localStorage / API later
                      return (
                        <Link
                          key={l.slug}
                          href={href}
                          onClick={onLessonClick}
                          className={cn(
                            'flex items-start gap-2 px-2 py-1 rounded text-[12.5px] leading-snug hover:bg-accent/30 transition-colors',
                            active ? 'bg-primary/10 text-primary font-medium' : 'text-foreground/85',
                            l.status === 'coming-soon' && 'opacity-60',
                          )}
                        >
                          {done ? (
                            <Check className="h-3 w-3 text-emerald-500 shrink-0 mt-0.5" />
                          ) : (
                            <Circle className="h-3 w-3 text-muted-foreground/40 shrink-0 mt-0.5" />
                          )}
                          <span className="flex-1">{l.title}</span>
                          <span className="text-[10px] text-muted-foreground/60 font-mono shrink-0 mt-0.5 flex items-center gap-0.5">
                            <Clock className="h-2.5 w-2.5" />
                            {l.minutes}m
                          </span>
                        </Link>
                      )
                    })}
                  </div>
                ))}
                {ch.assignment && (
                  <Link
                    href={`${chapterBase}/assignment`}
                    onClick={onLessonClick}
                    className={cn(
                      'mt-2 flex items-start gap-2 px-2 py-1.5 rounded text-[12.5px] hover:bg-accent/30 border border-dashed border-amber-500/20',
                      pathname === `${chapterBase}/assignment`
                        ? 'bg-amber-500/10 text-amber-500'
                        : 'text-amber-500/80',
                    )}
                  >
                    <Trophy className="h-3.5 w-3.5 shrink-0 mt-0.5" />
                    <span className="flex-1">Assignment: {ch.assignment.title}</span>
                  </Link>
                )}
              </div>
            )}
          </div>
        )
      })}
    </nav>
  )
}
