import Link from 'next/link'
import { ArrowLeft, ArrowRight, BookOpen, Clock, Trophy, Rocket, ShieldCheck, Cloud, Bell, Store } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Grit Mobile — Building Mobile Applications with Grit',
  description: '5 self-paced courses teaching you to build cross-platform mobile apps with Go API + Expo React Native. From scaffolding to the App Store.',
}

const courses = [
  {
    number: 1,
    title: 'Your First Mobile App',
    href: '/courses/grit-mobile/first-app',
    icon: Rocket,
    description: 'Scaffold API + Expo, set up emulator, run on device, understand the project.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 2,
    title: 'Mobile Auth & Navigation',
    href: '/courses/grit-mobile/auth-navigation',
    icon: ShieldCheck,
    description: 'Expo Router, secure token storage, tab navigation, protected routes.',
    challenges: 12,
    duration: '30 min',
  },
  {
    number: 3,
    title: 'API Integration & Offline',
    href: '/courses/grit-mobile/api-offline',
    icon: Cloud,
    description: 'TanStack Query in React Native, pull-to-refresh, infinite scroll, offline cache.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 4,
    title: 'Push Notifications',
    href: '/courses/grit-mobile/notifications',
    icon: Bell,
    description: 'Expo push tokens, grit-notifications plugin, FCM, APNS.',
    challenges: 10,
    duration: '30 min',
  },
  {
    number: 5,
    title: 'Build & App Store',
    href: '/courses/grit-mobile/build-store',
    icon: Store,
    description: 'EAS Build, TestFlight, Play Store, OTA updates.',
    challenges: 10,
    duration: '30 min',
  },
]

export default function GritMobileCourses() {
  const totalChallenges = courses.reduce((sum, c) => sum + c.challenges, 0)

  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Grit Mobile</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Building Mobile Applications with Grit
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed mb-6">
            5 self-paced courses that take you from zero to the App Store.
            Build cross-platform mobile apps with Go API + Expo + React Native.
            Each course is ~30 minutes with hands-on challenges you complete yourself.
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
            href="/courses/grit-mobile/first-app"
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
