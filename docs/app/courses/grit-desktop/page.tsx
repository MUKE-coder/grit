import Link from 'next/link'
import { ArrowLeft, ArrowRight, BookOpen, Clock, Trophy, Monitor, Database, Palette, FileDown, Package } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Grit Desktop — Building Desktop Applications with Grit',
  description: '5 self-paced courses teaching you to build native desktop apps with Grit and Wails. From your first app to distribution.',
}

const courses = [
  {
    number: 1,
    title: 'Your First Desktop App',
    href: '/courses/grit-desktop/first-app',
    icon: Monitor,
    description: 'Install Wails, scaffold a desktop app, understand Wails bindings, run in dev mode.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 2,
    title: 'Desktop CRUD & Data',
    href: '/courses/grit-desktop/crud-data',
    icon: Database,
    description: 'Generate resources, SQLite database, GORM models, Wails method bindings.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 3,
    title: 'Custom UI & Theming',
    href: '/courses/grit-desktop/custom-ui',
    icon: Palette,
    description: 'Frameless window, custom title bar, sidebar, dark theme, shadcn/ui components.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 4,
    title: 'PDF & Excel Export',
    href: '/courses/grit-desktop/export',
    icon: FileDown,
    description: 'Generate PDFs, Excel spreadsheets, CSV files, download flow.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 5,
    title: 'Build & Distribution',
    href: '/courses/grit-desktop/build',
    icon: Package,
    description: 'Compile native binary, platform targets, app icons, distribution.',
    challenges: 10,
    duration: '30 min',
  },
]

export default function GritDesktopCourses() {
  const totalChallenges = courses.reduce((sum, c) => sum + c.challenges, 0)

  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Grit Desktop</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Building Desktop Applications with Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed mb-6">
            5 self-paced courses that take you from zero to a distributable desktop app.
            Each course is ~30 minutes with hands-on challenges you complete yourself.
            No videos, no lectures — just concepts, examples, and practice.
          </p>

          {/* Stats */}
          <div className="flex flex-wrap gap-6">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <BookOpen className="h-4 w-4 text-primary" />
              <span><strong className="text-foreground">5</strong> courses</span>
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4 text-primary" />
              <span><strong className="text-foreground">~2.5 hours</strong> total</span>
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Trophy className="h-4 w-4 text-primary" />
              <span><strong className="text-foreground">{totalChallenges}</strong> challenges</span>
            </div>
          </div>
        </div>

        <hr className="border-border/40 mb-8" />

        {/* Course List */}
        <div className="space-y-4">
          {courses.map((course) => (
            <Link
              key={course.number}
              href={course.href}
              className="group flex items-start gap-4 p-5 rounded-lg border border-border/40 bg-card/30 hover:border-primary/30 hover:bg-card/50 transition-all"
            >
              <div className="flex items-center justify-center h-10 w-10 rounded-lg bg-primary/10 text-primary shrink-0 group-hover:bg-primary/20 transition-colors">
                <course.icon className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="text-xs font-mono text-muted-foreground">Course {course.number}</span>
                  <span className="text-xs text-muted-foreground/50">•</span>
                  <span className="text-xs text-muted-foreground">{course.duration}</span>
                  <span className="text-xs text-muted-foreground/50">•</span>
                  <span className="text-xs text-muted-foreground">{course.challenges} challenges</span>
                </div>
                <h3 className="font-semibold text-foreground group-hover:text-primary transition-colors">
                  {course.title}
                </h3>
                <p className="text-sm text-muted-foreground mt-1">{course.description}</p>
              </div>
              <ArrowRight className="h-4 w-4 text-muted-foreground group-hover:text-primary transition-colors shrink-0 mt-3" />
            </Link>
          ))}
        </div>

        {/* Navigation */}
        <div className="flex items-center justify-between pt-8 mt-8 border-t border-border/40">
          <Link
            href="/courses"
            className="flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            All Courses
          </Link>
          <Link
            href="/courses/grit-desktop/first-app"
            className="flex items-center gap-2 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
          >
            Start Course 1
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </main>
    </div>
  )
}
