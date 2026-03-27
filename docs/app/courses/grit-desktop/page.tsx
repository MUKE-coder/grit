import Link from "next/link"
import {
  Monitor,
  Clock,
  BookOpen,
  ArrowLeft,
  Trophy,
  ChevronDown,
  Folder,
  Play,
  Paintbrush,
  FileOutput,
  Package,
} from "lucide-react"
import { SiteHeader } from "@/components/site-header"
import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Grit Desktop Course — Building Desktop Applications with Grit",
  description:
    "Free, self-paced course covering Wails + Go + React desktop app development. 5 mini-courses, ~2.5 hours, 15 challenges.",
  openGraph: {
    title: "Grit Desktop Course — Building Desktop Applications with Grit",
    description:
      "Free, self-paced course covering Wails + Go + React desktop app development. 5 mini-courses, ~2.5 hours, 15 challenges.",
    url: "https://gritframework.dev/courses/grit-desktop",
  },
}

/* -- Course data ----------------------------------------------------------- */

interface Lesson {
  title: string
  description: string
}

interface MiniCourse {
  number: number
  title: string
  duration: string
  overview: string
  icon: React.ComponentType<{ className?: string }>
  lessons: Lesson[]
  challenge: {
    title: string
    description: string
  }
}

const courses: MiniCourse[] = [
  {
    number: 1,
    title: "Your First Desktop App",
    duration: "30 min",
    overview:
      "Install Wails, scaffold a desktop app with Grit, understand the project structure, run in dev mode.",
    icon: Play,
    lessons: [
      {
        title: "Prerequisites",
        description:
          "Go, Node, Wails CLI (go install github.com/wailsapp/wails/v2/cmd/wails@latest).",
      },
      {
        title: "Scaffolding",
        description:
          "grit new-desktop myapp — what gets created and where files live.",
      },
      {
        title: "Project Structure",
        description:
          "main.go (Wails bootstrap), app.go (bound methods), frontend/ (React + Vite + TanStack Router).",
      },
      {
        title: "Dev Mode",
        description:
          "wails dev or grit start — hot reload for both Go and React simultaneously.",
      },
      {
        title: "Tour",
        description:
          "Login page, dashboard, blog CRUD, contact CRUD — everything scaffolded out of the box.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        'Scaffold a desktop app called "notes-app", run it in dev mode, register an account, create 3 blog posts.',
    },
  },
  {
    number: 2,
    title: "Desktop CRUD & Data",
    duration: "30 min",
    overview:
      "Generate resources in desktop projects, understand Wails bindings, GORM with SQLite.",
    icon: Folder,
    lessons: [
      {
        title: "Resource Generation",
        description:
          "grit generate resource works for desktop too — same workflow, desktop output.",
      },
      {
        title: "Wails Bindings",
        description:
          "Go functions callable from React directly — no HTTP layer needed.",
      },
      {
        title: "SQLite Database",
        description:
          "Local file-based database. No Docker, no server, CGO-free driver.",
      },
      {
        title: "GORM Models",
        description:
          "Same patterns as web projects — struct tags, migrations, relationships.",
      },
      {
        title: "Service Layer",
        description:
          "Business logic lives in Go services, called via Wails bindings from the frontend.",
      },
      {
        title: "TanStack Router",
        description:
          "File-based routes in frontend/src/routes/ — familiar React routing.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Create a Task resource with title, description, priority (int), and done (bool). Add filtering by priority in the UI.",
    },
  },
  {
    number: 3,
    title: "Custom UI & Theming",
    duration: "30 min",
    overview:
      "Customize the desktop app UI: title bar, sidebar, theme colors, layouts.",
    icon: Paintbrush,
    lessons: [
      {
        title: "Frameless Window",
        description:
          "Custom title bar with minimize, maximize, and close buttons — no native chrome.",
      },
      {
        title: "Draggable Regions",
        description:
          "CSS -webkit-app-region for window dragging on the custom title bar.",
      },
      {
        title: "Sidebar Navigation",
        description:
          "Auto-generated from resources. Collapsible with icons and labels.",
      },
      {
        title: "Dark Theme",
        description:
          "Grit theme CSS variables — bg-background, bg-elevated, accent, and more.",
      },
      {
        title: "shadcn/ui Components",
        description:
          "Using UI primitives (Button, Input, Dialog, etc.) in desktop apps.",
      },
      {
        title: "Responsive Panels",
        description:
          "Resizable split panels for master-detail layouts.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Change the accent color from purple to sky blue, add a custom dashboard card that shows task completion percentage, customize the title bar with your app name.",
    },
  },
  {
    number: 4,
    title: "PDF & Excel Export",
    duration: "30 min",
    overview:
      "Export data from your desktop app to PDF and Excel files.",
    icon: FileOutput,
    lessons: [
      {
        title: "Export Service",
        description:
          "Built-in export functionality — Go generates files, React triggers downloads.",
      },
      {
        title: "PDF Generation",
        description:
          "Generating formatted reports from your data with Go PDF libraries.",
      },
      {
        title: "Excel Spreadsheets",
        description:
          "Structured data export to .xlsx format with headers and styling.",
      },
      {
        title: "CSV Export",
        description:
          "Simple comma-separated output for data portability.",
      },
      {
        title: "Download Flow",
        description:
          "Go generates the file on disk, React triggers the download via Wails runtime.",
      },
      {
        title: "Custom Templates",
        description:
          "Styling PDF output with headers, footers, logos, and table layouts.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Add an export button to your tasks list that generates a PDF report showing all tasks grouped by priority, with completion statistics at the top.",
    },
  },
  {
    number: 5,
    title: "Build & Distribution",
    duration: "30 min",
    overview:
      "Build your desktop app for Windows, macOS, and Linux. Package for distribution.",
    icon: Package,
    lessons: [
      {
        title: "Build Command",
        description:
          "grit compile or wails build — single binary with embedded frontend.",
      },
      {
        title: "Platform Targets",
        description:
          "Windows (.exe), macOS (.app), Linux (binary) — all from one codebase.",
      },
      {
        title: "App Icon",
        description:
          "Customizing the application icon for each platform.",
      },
      {
        title: "Cross-Compilation",
        description:
          "Building for other platforms from your development machine.",
      },
      {
        title: "Distribution",
        description:
          "Packaging strategies, installers (NSIS, DMG), and code signing overview.",
      },
      {
        title: "Auto-Update",
        description:
          "Strategies for pushing updates to users of your desktop app.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Build your notes app for your current platform, test the compiled binary, share it with a friend on a different OS.",
    },
  },
]

/* -- Page component -------------------------------------------------------- */

export default function GritDesktopCoursePage() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main>
        {/* Back link */}
        <div className="container max-w-screen-xl px-6 pt-8">
          <Link
            href="/courses"
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to All Courses
          </Link>
        </div>

        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/[0.04] to-transparent" />
          <div className="container max-w-screen-xl relative py-16 px-6">
            <div className="max-w-3xl mx-auto text-center">
              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
                <Monitor className="h-3.5 w-3.5" />
                Desktop Track
              </div>

              <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1] text-foreground">
                Building Desktop Applications with Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl mx-auto mb-8">
                Learn to build cross-platform desktop apps with{" "}
                <span className="text-foreground font-medium">
                  Go + Wails + React
                </span>
                . From scaffolding to distribution in five self-paced
                mini-courses.
              </p>

              {/* Stats */}
              <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
                <span className="flex items-center gap-1.5">
                  <BookOpen className="h-4 w-4 text-primary/60" />
                  5 courses
                </span>
                <span className="flex items-center gap-1.5">
                  <Clock className="h-4 w-4 text-primary/60" />
                  ~2.5 hours
                </span>
                <span className="flex items-center gap-1.5">
                  <Trophy className="h-4 w-4 text-primary/60" />
                  5 challenges
                </span>
              </div>
            </div>
          </div>
        </section>

        {/* Courses */}
        <section className="py-16">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-3xl mx-auto space-y-6">
              {courses.map((course) => (
                <details
                  key={course.number}
                  className="group rounded-xl border border-border/40 bg-card/50 overflow-hidden"
                >
                  <summary className="flex items-center gap-4 p-6 cursor-pointer select-none hover:bg-card/70 transition-colors list-none [&::-webkit-details-marker]:hidden">
                    {/* Number badge */}
                    <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary">
                      {course.number}
                    </span>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h2 className="text-lg font-semibold tracking-tight text-foreground">
                          {course.title}
                        </h2>
                      </div>
                      <p className="text-sm text-muted-foreground mt-0.5">
                        {course.overview}
                      </p>
                    </div>

                    <div className="flex items-center gap-3 shrink-0">
                      <span className="hidden sm:flex items-center gap-1 text-xs text-muted-foreground">
                        <Clock className="h-3.5 w-3.5" />
                        {course.duration}
                      </span>
                      <ChevronDown className="h-5 w-5 text-muted-foreground transition-transform group-open:rotate-180" />
                    </div>
                  </summary>

                  <div className="px-6 pb-6 pt-2 border-t border-border/20">
                    {/* Lessons */}
                    <h3 className="text-xs font-mono font-medium text-muted-foreground uppercase tracking-wider mb-4">
                      Lessons
                    </h3>
                    <ul className="space-y-3 mb-6">
                      {course.lessons.map((lesson, i) => (
                        <li
                          key={i}
                          className="flex items-start gap-3 text-sm"
                        >
                          <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded bg-primary/10 text-[11px] font-mono font-medium text-primary/70 mt-px">
                            {i + 1}
                          </span>
                          <div>
                            <span className="font-medium text-foreground">
                              {lesson.title}
                            </span>
                            <span className="text-muted-foreground">
                              {" "}
                              &mdash; {lesson.description}
                            </span>
                          </div>
                        </li>
                      ))}
                    </ul>

                    {/* DIY Challenge */}
                    <div className="bg-primary/10 border border-primary/20 rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-2">
                        <Trophy className="h-4 w-4 text-primary" />
                        <h4 className="text-sm font-semibold text-foreground">
                          {course.challenge.title}
                        </h4>
                      </div>
                      <p className="text-sm text-muted-foreground leading-relaxed">
                        {course.challenge.description}
                      </p>
                    </div>
                  </div>
                </details>
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
