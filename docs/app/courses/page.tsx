import Link from "next/link"
import { Globe, Monitor, Smartphone, Clock, ArrowRight, BookOpen } from "lucide-react"
import { SiteHeader } from "@/components/site-header"
import type { Metadata } from "next"

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
    courseCount: 8,
    totalTime: "~4 hours",
    courses: [
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
    courseCount: 5,
    totalTime: "~2.5 hours",
    courses: [
      "Your First Desktop App",
      "Desktop CRUD & Data",
      "Custom UI & Theming",
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

/* -- Page component ------------------------------------------------------- */

export default function CoursesPage() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/[0.04] to-transparent" />
          <div className="container max-w-screen-xl relative py-20 px-6">
            <div className="max-w-3xl mx-auto text-center">
              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
                <BookOpen className="h-3.5 w-3.5" />
                DIY Self-Paced Courses
              </div>

              <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1] text-foreground">
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

        {/* Course categories */}
        <section className="py-16">
          <div className="container max-w-screen-xl px-6">
            <div className="grid gap-6 md:grid-cols-3">
              {categories.map((cat) => (
                <Link
                  key={cat.title}
                  href={cat.href}
                  className="group relative flex flex-col rounded-xl border border-border/40 bg-card/50 p-6 hover:border-primary/30 hover:bg-card/70 transition-all"
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
          </div>
        </section>
      </main>
    </div>
  )
}
