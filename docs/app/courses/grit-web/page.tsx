import Link from 'next/link'
import { ArrowLeft, ArrowRight, BookOpen, Clock, Trophy, BookMarked, Terminal, Wand2, ShieldCheck, LayoutDashboard, HardDrive, Mail, Sparkles, Rocket } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { LaneFlow } from '@/components/lane-flow'
import { GridFrame } from '@/components/grid-frame'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Grit Web — Building Web Applications with Grit',
  description: '9 self-paced courses teaching you to build full-stack web apps with Grit. From introduction to production deployment.',
}

const courses = [
  {
    number: 0,
    title: 'Introduction to Grit',
    href: '/courses/grit-web/introduction',
    icon: BookMarked,
    description: 'What Grit is, why it exists, the tech stack, architecture modes, philosophy, and comparison with other frameworks.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 1,
    title: 'Your First Grit App',
    href: '/courses/grit-web/first-app',
    icon: Terminal,
    description: 'Install Grit, scaffold a project, understand the structure, start the servers.',
    challenges: 18,
    duration: '30 min',
  },
  {
    number: 2,
    title: 'Code Generator Mastery',
    href: '/courses/grit-web/code-generator',
    icon: Wand2,
    description: 'Generate full-stack CRUD resources with field types, relationships, and modifiers.',
    challenges: 15,
    duration: '30 min',
  },
  {
    number: 3,
    title: 'Authentication & Authorization',
    href: '/courses/grit-web/authentication',
    icon: ShieldCheck,
    description: 'JWT tokens, roles, two-factor auth (TOTP), OAuth2, and custom roles.',
    challenges: 14,
    duration: '30 min',
  },
  {
    number: 4,
    title: 'Admin Panel Customization',
    href: '/courses/grit-web/admin-panel',
    icon: LayoutDashboard,
    description: 'DataTable, FormBuilder, dashboard widgets, multi-step forms, style variants.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 5,
    title: 'File Storage & Uploads',
    href: '/courses/grit-web/file-storage',
    icon: HardDrive,
    description: 'Presigned URL uploads, MinIO, S3, Cloudflare R2, image processing.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 6,
    title: 'Background Jobs & Email',
    href: '/courses/grit-web/jobs-email',
    icon: Mail,
    description: 'Redis job queues, cron scheduling, Resend emails, Mailhog for testing.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 7,
    title: 'AI-Powered Features',
    href: '/courses/grit-web/ai-features',
    icon: Sparkles,
    description: 'Vercel AI Gateway, completions, chat, SSE streaming, building a chat UI.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 8,
    title: 'Deploy to Production',
    href: '/courses/grit-web/deploy',
    icon: Rocket,
    description: 'grit deploy, systemd, Caddy, auto-TLS, maintenance mode, Docker deployment.',
    challenges: 12,
    duration: '30 min',
  },
]

export default function GritWebCourses() {
  const totalChallenges = courses.reduce((sum, c) => sum + c.challenges, 0)

  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Grit Web</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Building Web Applications with Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed mb-6">
            9 self-paced courses that take you from zero to production.
            Each course is ~30 minutes with hands-on challenges you complete yourself.
            No videos, no lectures — just concepts, examples, and practice.
          </p>
          <LaneFlow id="c-gwidx" lanes={['1 · Learn', '2 · Build', '3 · Ship']}
            nodes={[
              { id: 'learn', lane: 0, row: 0, title: 'Introduction', sub: 'what & why', tone: 'primary' },
              { id: 'build', lane: 1, row: 0, title: 'First app → generator → batteries', sub: 'hands-on', tone: 'cyan' },
              { id: 'ship', lane: 2, row: 0, title: 'Deploy', sub: 'production', tone: 'green' },
            ]}
            edges={[{ from: 'learn', to: 'build', tone: 'cyan' }, { from: 'build', to: 'ship', tone: 'green' }]}
            legend={[{ tone: 'primary', label: 'Concepts' }, { tone: 'green', label: 'Deployed app' }]}
            caption="Nine ~30-minute courses take you from concepts to a deployed, production app" />

          {/* Stats */}
          <div className="flex flex-wrap gap-6">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <BookOpen className="h-4 w-4 text-primary" />
              <span><strong className="text-foreground">9</strong> courses</span>
            </div>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4 text-primary" />
              <span><strong className="text-foreground">~4.5 hours</strong> total</span>
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
            href="/courses/grit-web/introduction"
            className="flex items-center gap-2 text-sm font-medium text-primary hover:text-primary/80 transition-colors"
          >
            Start Course 0
            <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
      </main>
    </div>
  )
}
