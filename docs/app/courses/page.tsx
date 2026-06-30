import Link from "next/link"
import { Globe, Monitor, Smartphone, Clock, ArrowRight, BookOpen, Zap, Server, Dumbbell, ShoppingCart, FileText, Shield, Activity, Database, Code2, Rocket, Bot, Receipt, TestTube2, HardDrive, MessageSquare, CreditCard, Newspaper, GitBranch, Wrench, Palette, Wifi, ShieldCheck, Flag, Cable, Layers, FileSpreadsheet, GraduationCap, Trophy, ListChecks, Construction, Sparkles } from "lucide-react"
import { SiteHeader } from "@/components/site-header"
import { GridFrame } from "@/components/grid-frame"
import type { Metadata } from "next"
import { COURSES, flatLessons, courseTotalMinutes } from "@/config/courses"

export const metadata: Metadata = {
  title: "Grit Courses",
  description:
    "Free, self-paced 30-minute courses for building web, desktop, and mobile applications with the Grit framework.",
  openGraph: {
    title: "Grit Courses",
    description:
      "Free, self-paced 30-minute courses for building web, desktop, and mobile applications with the Grit framework.",
    url: "https://gritframework.dev/courses",
  },
}

/* -- Course category data ------------------------------------------------- */

interface CourseCategory {
  title: string
  subtitle: string
  icon: React.ComponentType<{ className?: string }>
  href: string
  courseCount: number
  totalTime: string
  courses: string[]
}

const categories: CourseCategory[] = [
  {
    title: "Grit Web",
    subtitle: "Building Web Applications",
    icon: Globe,
    href: "/courses/grit-web",
    courseCount: 9,
    totalTime: "~4.5 hours",
    courses: [
      "Introduction to Grit",
      "Your First Grit App",
      "Code Generator Mastery",
      "Authentication & Authorization",
      "Admin Panel Customization",
      "File Storage & Uploads",
      "Background Jobs & Email",
      "AI-Powered Features",
      "Deploy to Production",
    ],
  },
  {
    title: "Grit Desktop",
    subtitle: "Building Desktop Applications",
    icon: Monitor,
    href: "/courses/grit-desktop",
    courseCount: 6,
    totalTime: "~3 hours",
    courses: [
      "Your First Desktop App",
      "Desktop CRUD & Data",
      "Custom UI & Theming",
      "Offline-First with Sync Engine",
      "PDF & Excel Export",
      "Build & Distribution",
    ],
  },
  {
    title: "Grit Mobile",
    subtitle: "Building Mobile Applications",
    icon: Smartphone,
    href: "/courses/grit-mobile",
    courseCount: 5,
    totalTime: "~2.5 hours",
    courses: [
      "Your First Mobile App",
      "Mobile Auth & Navigation",
      "API Integration & Offline",
      "Push Notifications",
      "Build & App Store",
    ],
  },
]

const standaloneCourses = [
  { title: "Offline-First Desktop", subtitle: "Local SQLite + outbox + Git-style sync", href: "/docs/desktop/offline", icon: Wifi, duration: "30 min" },
  { title: "Audit Log + Hash Chain", subtitle: "Tamper-evident activity tracking for SOC2", href: "/courses/audit-log", icon: ShieldCheck, duration: "30 min" },
  { title: "Feature Flags & A/B Testing", subtitle: "Sticky bucketing, percentage rollouts, realtime push", href: "/courses/feature-flags", icon: Flag, duration: "30 min" },
  { title: "Webhook Receiver", subtitle: "Stripe / GitHub / HMAC verifiers + replay", href: "/courses/webhook-receiver", icon: Cable, duration: "30 min" },
  { title: "CSV / Excel Export", subtitle: "Auto-generated per resource via grit generate", href: "/courses/export", icon: FileSpreadsheet, duration: "30 min" },
  { title: "Realtime + WebSocket Hub", subtitle: "SendToUser, Broadcast, useRealtimeEvent", href: "/courses/realtime-chat", icon: Layers, duration: "30 min" },
  { title: "Batteries Included", subtitle: "Every feature that ships with Grit", href: "/courses/batteries", icon: Zap, duration: "30 min" },
  { title: "API-Only Masterclass", subtitle: "Build & deploy a REST API with Go", href: "/courses/api-masterclass", icon: Server, duration: "30 min" },
  { title: "Build a Fitness App", subtitle: "Go API + Expo React Native", href: "/courses/mobile-fitness-app", icon: Dumbbell, duration: "30 min" },
  { title: "E-commerce Store", subtitle: "Single app architecture with Vite", href: "/courses/ecommerce-spa", icon: ShoppingCart, duration: "30 min" },
  { title: "API Docs: Scalar & Swagger", subtitle: "Auto-generated API documentation", href: "/courses/api-docs-scalar", icon: FileText, duration: "30 min" },
  { title: "Security Deep Dive", subtitle: "Auth, 2FA, WAF & rate limiting", href: "/courses/security-deep-dive", icon: Shield, duration: "30 min" },
  { title: "Pulse Analytics", subtitle: "Tracing, metrics & monitoring", href: "/courses/pulse-analytics", icon: Activity, duration: "30 min" },
  { title: "GORM Studio", subtitle: "The visual database browser", href: "/courses/gorm-studio", icon: Database, duration: "30 min" },
  { title: "React + Vite + Go", subtitle: "Building with TanStack Router", href: "/courses/react-vite-go", icon: Code2, duration: "30 min" },
  { title: "Deployment Guide", subtitle: "Dokploy, Orbita, VPS & Vercel", href: "/courses/deployment-guide", icon: Rocket, duration: "30 min" },
  { title: "SaaS with Claude Code", subtitle: "AI-assisted SaaS development", href: "/courses/saas-with-ai", icon: Bot, duration: "30 min" },
  { title: "Invoice Generator", subtitle: "Desktop app with Wails + PDF export", href: "/courses/invoice-desktop", icon: Receipt, duration: "30 min" },
  { title: "Testing Your Grit App", subtitle: "Go, Vitest & Playwright", href: "/courses/testing", icon: TestTube2, duration: "30 min" },
  { title: "Database Mastery", subtitle: "GORM models, migrations & queries", href: "/courses/gorm-mastery", icon: HardDrive, duration: "30 min" },
  { title: "Real-Time Chat", subtitle: "WebSockets with grit-websockets", href: "/courses/realtime-chat", icon: MessageSquare, duration: "30 min" },
  { title: "Stripe Payments", subtitle: "Subscriptions & billing for SaaS", href: "/courses/stripe-payments", icon: CreditCard, duration: "30 min" },
  { title: "Blog & CMS", subtitle: "Complete content management system", href: "/courses/blog-cms", icon: Newspaper, duration: "30 min" },
  { title: "CI/CD with GitHub Actions", subtitle: "Automated testing & deployment", href: "/courses/cicd-github", icon: GitBranch, duration: "30 min" },
  { title: "Custom Middleware", subtitle: "Extending Grit with hooks", href: "/courses/custom-middleware", icon: Wrench, duration: "30 min" },
  { title: "Grit UI Components", subtitle: "Using the 100-component registry", href: "/courses/grit-ui-components", icon: Palette, duration: "30 min" },
]

/* -- Page component ------------------------------------------------------- */

export default function CoursesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div
            className="absolute inset-0 -z-10"
            style={{
              background:
                'radial-gradient(ellipse 70% 70% at 50% -10%, hsl(var(--primary) / 0.12), transparent 60%)',
            }}
          />
          <span className="crosshair absolute top-16 left-[16%] text-foreground/20 hidden md:block" style={{ width: 14, height: 14 }} />
          <span className="crosshair absolute top-28 right-[18%] text-primary/30 hidden md:block" style={{ width: 14, height: 14 }} />
          <div className="max-w-6xl mx-auto relative py-24 px-6">
            <div className="max-w-3xl mx-auto text-center">
              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
                <BookOpen className="h-3.5 w-3.5" />
                DIY Self-Paced Courses
              </div>

              <h1 className="font-display text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1] text-foreground">
                Learn Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl mx-auto">
                Free, self-paced <span className="text-foreground font-medium">30-minute courses</span> that teach you
                how to build web, desktop, and mobile applications with the Grit framework.
                Pick a track and start building.
              </p>
            </div>
          </div>
        </section>

        {/* Learning Paths — kit-based multi-chapter courses */}
        <section className="py-20 md:py-24 border-b border-border/30">
          <div className="max-w-6xl mx-auto px-6">
            <div className="text-center mb-14">
              <div className="inline-flex items-center gap-2 rounded-full bg-violet-500/10 border border-violet-500/20 px-3 py-1 text-xs font-mono font-medium text-violet-400 mb-4">
                <GraduationCap className="h-3.5 w-3.5" />
                Learning Paths
              </div>
              <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground tracking-tight mb-4">
                Pick your kit, follow a path
              </h2>
              <p className="text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">
                Multi-chapter courses built around the seven Grit tech kits. Each has
                its own chapters, lessons, exercises, and end-of-chapter assignments —
                tuned for beginners and designed to be worked through in order.
              </p>
            </div>

            <div className="grid md:grid-cols-2 lg:grid-cols-3 rounded-xl border-[3px] border-double border-foreground/20 overflow-hidden mb-4">
              {COURSES.map((course) => {
                const lessons = flatLessons(course).length
                const hours = Math.round(courseTotalMinutes(course) / 60)
                const isComing = course.status === "coming-soon"
                return (
                  <Link
                    key={course.slug}
                    href={`/courses/${course.slug}`}
                    className={`group relative overflow-hidden border-b-[3px] border-r-[3px] border-double border-foreground/20 bg-card/20 p-6 hover:bg-card/50 transition-colors flex flex-col`}
                  >
                    <div className={`absolute -top-12 -right-12 h-32 w-32 rounded-full bg-gradient-to-br ${course.accent} blur-2xl opacity-50 pointer-events-none`} />
                    <div className="relative flex items-start gap-3 mb-3">
                      <span className="text-3xl shrink-0" aria-hidden="true">{course.emoji}</span>
                      <div className="flex-1 min-w-0">
                        <h3 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors leading-tight mb-1">
                          {course.name}
                        </h3>
                        <p className="text-xs text-muted-foreground leading-relaxed line-clamp-2">
                          {course.tagline}
                        </p>
                      </div>
                    </div>
                    <div className="relative flex items-center gap-3 text-[10px] font-mono text-muted-foreground/80 mb-4">
                      <span className="inline-flex items-center gap-1">
                        <BookOpen className="h-2.5 w-2.5" />
                        {course.chapters.length} ch
                      </span>
                      <span className="inline-flex items-center gap-1">
                        <ListChecks className="h-2.5 w-2.5" />
                        {lessons} lessons
                      </span>
                      <span className="inline-flex items-center gap-1">
                        <Clock className="h-2.5 w-2.5" />
                        ~{hours}h
                      </span>
                    </div>
                    <div className="relative mt-auto flex items-center justify-between">
                      <span className="inline-flex items-center gap-1.5 text-xs font-medium text-primary group-hover:gap-2.5 transition-all">
                        {isComing ? "Preview outline" : "Start course"}
                        <ArrowRight className="h-3.5 w-3.5" />
                      </span>
                      {isComing && (
                        <span className="inline-flex items-center gap-1 text-[9px] font-mono text-amber-500 bg-amber-500/10 border border-amber-500/20 rounded-full px-1.5 py-0.5">
                          <Construction className="h-2.5 w-2.5" />
                          In production
                        </span>
                      )}
                    </div>
                  </Link>
                )
              })}
            </div>

            <p className="text-center text-xs text-muted-foreground/70 mt-6">
              <Sparkles className="inline h-3 w-3 mr-1 text-primary" />
              Start with <Link href="/courses/concepts" className="text-primary hover:underline">Grit Concepts</Link> if you&apos;re new — every other path assumes it.
            </p>
          </div>
        </section>

        {/* Focused Tutorials (legacy categories) */}
        <section className="py-20 md:py-24">
          <div className="max-w-6xl mx-auto px-6">
            <div className="text-center mb-14">
              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-4">
                <Zap className="h-3.5 w-3.5" />
                Focused Tutorials
              </div>
              <h2 className="font-display text-3xl md:text-4xl font-bold text-foreground tracking-tight mb-4">
                Bite-sized walk-throughs
              </h2>
              <p className="text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">
                Shorter, focused tutorials on specific topics — pick one, finish in
                ~30 minutes, build something concrete.
              </p>
            </div>
            <div className="grid md:grid-cols-3 rounded-xl border-[3px] border-double border-foreground/20 overflow-hidden">
              {categories.map((cat) => (
                <Link
                  key={cat.title}
                  href={cat.href}
                  className="group relative flex flex-col border-b-[3px] border-r-[3px] border-double border-foreground/20 bg-card/20 p-7 hover:bg-card/50 transition-colors"
                >
                  {/* Icon + meta */}
                  <div className="flex items-center gap-3 mb-4">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15">
                      <cat.icon className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <h2 className="text-lg font-semibold tracking-tight text-foreground group-hover:text-primary transition-colors">
                        {cat.title}
                      </h2>
                      <p className="text-xs text-muted-foreground">{cat.subtitle}</p>
                    </div>
                  </div>

                  {/* Stats row */}
                  <div className="flex items-center gap-4 mb-5 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1">
                      <BookOpen className="h-3.5 w-3.5 text-primary/60" />
                      {cat.courseCount} courses
                    </span>
                    <span className="flex items-center gap-1">
                      <Clock className="h-3.5 w-3.5 text-primary/60" />
                      {cat.totalTime}
                    </span>
                  </div>

                  {/* Mini-course list */}
                  <ul className="flex-1 space-y-2 mb-6">
                    {cat.courses.map((course, i) => (
                      <li key={i} className="flex items-start gap-2.5 text-[13px] text-muted-foreground/80">
                        <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded bg-primary/10 text-[10px] font-mono font-medium text-primary/70 mt-px">
                          {i + 1}
                        </span>
                        {course}
                      </li>
                    ))}
                  </ul>

                  {/* CTA */}
                  <div className="flex items-center gap-1.5 text-sm font-medium text-primary group-hover:gap-2.5 transition-all">
                    Start Learning
                    <ArrowRight className="h-4 w-4" />
                  </div>
                </Link>
              ))}
            </div>

            {/* Standalone Courses */}
            <div className="mt-24">
              <div className="text-center mb-14">
                <h2 className="font-display text-2xl md:text-3xl font-bold text-foreground mb-3">Standalone Courses</h2>
                <p className="text-muted-foreground">Deep dives, practical builds, and specialized topics</p>
              </div>

              <div className="grid sm:grid-cols-2 lg:grid-cols-3 rounded-xl border-[3px] border-double border-foreground/20 overflow-hidden">
                {standaloneCourses.map((course) => (
                  <Link
                    key={course.href}
                    href={course.href}
                    className="group flex items-start gap-3 p-5 border-b-[3px] border-r-[3px] border-double border-foreground/20 bg-card/20 hover:bg-card/50 transition-colors"
                  >
                    <div className="flex items-center justify-center h-9 w-9 rounded-lg bg-primary/10 text-primary shrink-0 group-hover:bg-primary/20 transition-colors">
                      <course.icon className="h-4 w-4" />
                    </div>
                    <div className="min-w-0">
                      <h3 className="text-sm font-semibold text-foreground group-hover:text-primary transition-colors">
                        {course.title}
                      </h3>
                      <p className="text-xs text-muted-foreground mt-0.5">{course.subtitle}</p>
                      <span className="text-[10px] text-muted-foreground/60 mt-1 block">{course.duration}</span>
                    </div>
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
