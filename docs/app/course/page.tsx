import Link from "next/link"
import { ArrowRight, Clock, Code2, BookOpen, Trophy, CheckCircle2, Play, Zap, Users, BarChart3, Shield, Database, Layers, Rocket, GraduationCap, Target, Briefcase } from "lucide-react"
import { Button } from "@/components/ui/button"
import { SiteHeader } from "@/components/site-header"
import { WaitlistForm } from "@/components/waitlist-form"
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Grit Fullstack Course — Go + React from Scratch',
  description: 'Master fullstack development with Grit. Build 5 real-world apps, complete 50+ assignments, and learn Go, React, Next.js, and the Grit framework from zero to production.',
  openGraph: {
    title: 'Grit Fullstack Course — Go + React from Scratch',
    description: 'Master fullstack development with Grit. Build 5 real-world apps, complete 50+ assignments.',
    url: 'https://gritframework.dev/course',
    type: 'website',
  },
}

/* ─── Course Data ─────────────────────────────────────────── */

const stats = [
  { label: "Modules", value: "10", icon: BookOpen },
  { label: "Lessons", value: "60+", icon: Play },
  { label: "Assignments", value: "55+", icon: Code2 },
  { label: "Projects Built", value: "5", icon: Trophy },
]

interface Lesson {
  title: string
  assignment?: string
}

interface Module {
  number: number
  title: string
  description: string
  icon: React.ComponentType<{ className?: string }>
  duration: string
  lessons: Lesson[]
  project?: string
}

const modules: Module[] = [
  {
    number: 1,
    title: "Foundations — Go Essentials",
    description: "Learn the Go programming language from scratch. Variables, types, structs, functions, error handling, pointers, interfaces, and concurrency basics.",
    icon: Code2,
    duration: "8 lessons",
    lessons: [
      { title: "Installing Go & your editor setup", assignment: "Set up Go, VS Code, and run your first program" },
      { title: "Variables, types & constants", assignment: "Build a tip calculator using variables and math" },
      { title: "Control flow — if, switch, for", assignment: "Write a number guessing game" },
      { title: "Functions & multiple return values", assignment: "Create a temperature converter with helper functions" },
      { title: "Structs & methods", assignment: "Model a contact card with struct methods" },
      { title: "Slices, maps & iteration", assignment: "Build an in-memory phonebook (add/search/delete)" },
      { title: "Pointers & interfaces", assignment: "Create a Shape interface with area/perimeter methods" },
      { title: "Error handling & packages", assignment: "Build a file reader that handles errors gracefully" },
    ],
  },
  {
    number: 2,
    title: "Foundations — React & Next.js",
    description: "Learn React fundamentals and Next.js App Router. Components, hooks, state management, routing, data fetching, and TypeScript.",
    icon: Layers,
    duration: "7 lessons",
    lessons: [
      { title: "React components & JSX", assignment: "Build a profile card component" },
      { title: "Props, state & events (useState)", assignment: "Create a counter with increment/decrement/reset" },
      { title: "Lists, conditional rendering & forms", assignment: "Build a todo list with add/delete/toggle" },
      { title: "useEffect, data fetching & loading states", assignment: "Fetch and display users from a public API" },
      { title: "Next.js App Router & file-based routing", assignment: "Create a multi-page portfolio site" },
      { title: "TypeScript essentials for React", assignment: "Convert a JS component to fully typed TypeScript" },
      { title: "Tailwind CSS & shadcn/ui", assignment: "Restyle the todo app with Tailwind + shadcn components" },
    ],
  },
  {
    number: 3,
    title: "Getting Started with Grit",
    description: "Install Grit, scaffold your first project, and understand the monorepo architecture. Docker, configuration, and dev workflow.",
    icon: Rocket,
    duration: "5 lessons",
    lessons: [
      { title: "Installing Grit CLI & Docker", assignment: "Install Grit and verify with grit --version" },
      { title: "Scaffolding your first project", assignment: "Run grit new and explore every generated file" },
      { title: "Project structure deep dive", assignment: "Draw a diagram of the apps/api folder structure" },
      { title: "Docker Compose & infrastructure", assignment: "Start PostgreSQL, Redis, MinIO and connect via Studio" },
      { title: "Environment configuration & dev workflow", assignment: "Configure .env, run the API, web, and admin simultaneously" },
    ],
  },
  {
    number: 4,
    title: "Build App #1 — Contact Manager",
    description: "Your first full-stack Grit app. Generate resources, explore the admin panel, GORM Studio, and auto-generated API docs.",
    icon: Users,
    duration: "6 lessons",
    project: "Contact Manager CRM",
    lessons: [
      { title: "Generate a Group resource", assignment: "Create Group with name, description, color fields" },
      { title: "Generate a Contact resource with belongs_to", assignment: "Create Contact with name, email, phone, notes, belongs_to:Group" },
      { title: "Explore the admin panel", assignment: "Add 10 groups and 20 contacts via the admin UI" },
      { title: "Customize the admin resource definitions", assignment: "Add badges, filters, and search to the contacts table" },
      { title: "GORM Studio & API Documentation", assignment: "Browse data in Studio and test 3 endpoints via /docs" },
      { title: "Build a public contacts page in Next.js", assignment: "Create a searchable contact directory on the web app" },
    ],
  },
  {
    number: 5,
    title: "Go Backend Deep Dive",
    description: "Master Go backend patterns — models, handlers, services, middleware, authentication, JWT, RBAC, and custom business logic.",
    icon: Database,
    duration: "8 lessons",
    lessons: [
      { title: "GORM models — tags, types & relationships", assignment: "Create a Project model with has_many Tasks" },
      { title: "Handlers — request parsing, validation & responses", assignment: "Write a custom handler with query filters" },
      { title: "Services — business logic separation", assignment: "Move task assignment logic to a TaskService" },
      { title: "Authentication — JWT tokens & bcrypt", assignment: "Test the full auth flow: register, login, refresh, me" },
      { title: "Middleware — auth, CORS, logging & custom", assignment: "Write a custom middleware that logs request duration" },
      { title: "RBAC — roles & protected routes", assignment: "Add an EDITOR role and restrict project deletion to ADMIN" },
      { title: "Pagination, sorting & search patterns", assignment: "Implement multi-field search across 3 string columns" },
      { title: "Error handling & validation best practices", assignment: "Add custom validation for email uniqueness and date ranges" },
    ],
  },
  {
    number: 6,
    title: "Build App #2 — Task Manager",
    description: "Build a project & task management app. Practice backend patterns, relationships, role-based access, and custom endpoints.",
    icon: Target,
    duration: "5 lessons",
    project: "Task Manager with RBAC",
    lessons: [
      { title: "Scaffold and generate Project + Task resources", assignment: "Generate both resources with proper relationships" },
      { title: "Add custom endpoints — assign, complete, archive", assignment: "Create 3 custom endpoints beyond standard CRUD" },
      { title: "Implement role-based task assignment", assignment: "Only ADMIN/EDITOR can assign, USER can only view own tasks" },
      { title: "Build task stats endpoint with aggregations", assignment: "Return counts by status, overdue tasks, and completion rate" },
      { title: "Customize admin with badges, filters & widgets", assignment: "Add status badges, priority filter, and task stats widget" },
    ],
  },
  {
    number: 7,
    title: "React Frontend & Admin Mastery",
    description: "Master the frontend — React Query hooks, shared types, forms, DataTable, multi-step forms, dashboard widgets, and theming.",
    icon: BarChart3,
    duration: "7 lessons",
    lessons: [
      { title: "React Query hooks — queries & mutations", assignment: "Create a custom hook with optimistic updates" },
      { title: "Shared package — Zod schemas & TypeScript types", assignment: "Add a new field and sync types across frontend and backend" },
      { title: "DataTable — columns, sorting, filters & search", assignment: "Build a custom column with a progress bar renderer" },
      { title: "FormBuilder — all field types", assignment: "Create a form with text, select, date, toggle, and richtext" },
      { title: "Multi-step forms — modal-steps & page-steps", assignment: "Build a 3-step customer onboarding form" },
      { title: "Dashboard widgets — stats, charts & activity", assignment: "Create a dashboard with 4 stat cards and 2 charts" },
      { title: "Theme customization & style variants", assignment: "Switch between default, modern, minimal, and glass styles" },
    ],
  },
  {
    number: 8,
    title: "Build App #3 — Customer CRM",
    description: "Build a full customer relationship management app with companies, contacts, deals pipeline, activity log, and dashboard analytics.",
    icon: Briefcase,
    duration: "6 lessons",
    project: "Customer CRM with Pipeline",
    lessons: [
      { title: "Design the CRM data model", assignment: "Create Company, Contact, Deal, and Activity resources" },
      { title: "Build a deal pipeline with status workflow", assignment: "Implement deal stages: Lead → Qualified → Proposal → Won/Lost" },
      { title: "Add an activity timeline for contacts", assignment: "Log calls, emails, meetings as activity entries" },
      { title: "Build a CRM dashboard with analytics", assignment: "Show revenue by stage, conversion rate, and top deals" },
      { title: "Implement search across all CRM entities", assignment: "Build a global search that queries companies, contacts, and deals" },
      { title: "Add social login (Google & GitHub)", assignment: "Configure OAuth and test the full social auth flow" },
    ],
  },
  {
    number: 9,
    title: "Batteries — Production Features",
    description: "Add production-grade features: file uploads, email, background jobs, cron, caching, AI integration, security, and observability.",
    icon: Zap,
    duration: "8 lessons",
    lessons: [
      { title: "File storage — upload, resize & serve images", assignment: "Add avatar upload to the user profile" },
      { title: "Email — Resend + HTML templates", assignment: "Send a welcome email on user registration" },
      { title: "Background jobs — async processing with asynq", assignment: "Process image thumbnails in a background worker" },
      { title: "Cron scheduler — recurring tasks", assignment: "Schedule a daily report email and weekly cleanup job" },
      { title: "Redis caching — cache service & middleware", assignment: "Cache the dashboard stats endpoint for 5 minutes" },
      { title: "AI integration — Claude & OpenAI", assignment: "Add an AI-powered 'summarize deal' button" },
      { title: "Security — Sentinel WAF & rate limiting", assignment: "Configure rate limits and test brute-force protection" },
      { title: "Observability — Pulse tracing & metrics", assignment: "Monitor API performance and identify slow endpoints" },
    ],
  },
  {
    number: 10,
    title: "Build App #4 & #5 — Capstone",
    description: "Build two final apps: an inventory management system and a complete SaaS CRM platform. Deploy to production.",
    icon: Trophy,
    duration: "8 lessons",
    project: "Inventory System + SaaS CRM Platform",
    lessons: [
      { title: "App #4 — Inventory: data model & resources", assignment: "Create Product, Category, Warehouse, StockMovement resources" },
      { title: "App #4 — Inventory: stock tracking & alerts", assignment: "Build low-stock alerts and stock movement history" },
      { title: "App #4 — Inventory: reports & export", assignment: "Generate inventory reports with charts and CSV export" },
      { title: "App #5 — SaaS CRM: multi-feature design", assignment: "Plan the full data model for a production CRM" },
      { title: "App #5 — SaaS CRM: pipeline, email & files", assignment: "Combine deals, email templates, and file attachments" },
      { title: "App #5 — SaaS CRM: AI assistant & automation", assignment: "Add AI deal scoring and automated follow-up reminders" },
      { title: "Deployment — Docker, Dokploy & production", assignment: "Deploy your SaaS CRM to a VPS with SSL" },
      { title: "Course wrap-up & next steps", assignment: "Final review: refactor, optimize, and plan your own project" },
    ],
  },
]

/* ─── Total Assignments Count ─────────────────────────────── */
const totalAssignments = modules.reduce((sum, m) => sum + m.lessons.filter(l => l.assignment).length, 0)

/* ─── Projects List ──────────────────────────────────────── */
const projects = [
  { number: 1, name: "Contact Manager", module: 4, description: "CRUD, code generation, admin panel, GORM Studio" },
  { number: 2, name: "Task Manager", module: 6, description: "RBAC, custom endpoints, relationships, aggregations" },
  { number: 3, name: "Customer CRM", module: 8, description: "Pipeline, activity log, dashboard analytics, OAuth" },
  { number: 4, name: "Inventory System", module: 10, description: "Stock tracking, alerts, reports, CSV export" },
  { number: 5, name: "SaaS CRM Platform", module: 10, description: "Full-featured CRM with AI, email, files, deployment" },
]

/* ─── Page Component ──────────────────────────────────────── */

export default function CoursePage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/[0.03] to-transparent" />
          <div className="container max-w-screen-xl relative py-20 px-6">
            <div className="max-w-3xl">
              <div className="flex items-center gap-2 mb-6">
                <span className="inline-flex items-center rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary">
                  Free Course
                </span>
                <span className="inline-flex items-center rounded-full bg-accent/10 border border-accent/20 px-3 py-1 text-xs font-mono font-medium text-muted-foreground">
                  10 Modules &middot; {totalAssignments} Assignments
                </span>
              </div>

              <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1]">
                Grit Fullstack Course
              </h1>
              <p className="text-xl text-muted-foreground leading-relaxed mb-4">
                Go from zero to building production-ready fullstack apps with{" "}
                <span className="text-foreground font-medium">Go + React</span>.
                Learn the Grit framework by building{" "}
                <span className="text-foreground font-medium">5 real-world applications</span>{" "}
                and completing{" "}
                <span className="text-foreground font-medium">{totalAssignments}+ hands-on assignments</span>.
              </p>
              <p className="text-[15px] text-muted-foreground/70 leading-relaxed mb-8">
                No prior Go or backend experience needed. We start from the absolute basics and
                build up to deploying a complete SaaS CRM platform with authentication, file uploads,
                email, background jobs, AI integration, and more.
              </p>

              <div className="max-w-xl">
                <WaitlistForm />
              </div>
            </div>
          </div>
        </section>

        {/* Stats */}
        <section className="border-b border-border/30 bg-card/30">
          <div className="container max-w-screen-xl px-6">
            <div className="grid grid-cols-2 md:grid-cols-4 divide-x divide-border/30">
              {stats.map((stat) => (
                <div key={stat.label} className="py-8 px-6 text-center">
                  <stat.icon className="h-5 w-5 text-primary mx-auto mb-2" />
                  <div className="text-2xl font-bold tracking-tight">{stat.value}</div>
                  <div className="text-xs text-muted-foreground mt-1">{stat.label}</div>
                </div>
              ))}
            </div>
          </div>
        </section>

        {/* What you'll build */}
        <section className="py-16 border-b border-border/30">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                5 Apps You&apos;ll Build
              </h2>
              <p className="text-sm text-muted-foreground mb-8">
                Each project builds on the previous one, progressively introducing new concepts.
              </p>

              <div className="space-y-3">
                {projects.map((p) => (
                  <Link
                    key={p.number}
                    href={`#module-${p.module}`}
                    className="group flex items-start gap-4 rounded-lg border border-border/40 bg-card/50 p-4 hover:border-primary/20 hover:bg-card/80 transition-all"
                  >
                    <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-primary/10 border border-primary/15 text-sm font-mono font-bold text-primary">
                      {p.number}
                    </div>
                    <div className="min-w-0">
                      <div className="font-medium text-[15px] group-hover:text-primary transition-colors">
                        {p.name}
                      </div>
                      <div className="text-xs text-muted-foreground/70 mt-0.5">
                        Module {p.module} &middot; {p.description}
                      </div>
                    </div>
                    <ArrowRight className="ml-auto h-4 w-4 text-muted-foreground/30 group-hover:text-primary/50 transition-colors shrink-0 mt-2.5" />
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* Course Curriculum */}
        <section className="py-16">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-3xl mx-auto">
              <h2 className="text-2xl font-semibold tracking-tight mb-2">
                Course Curriculum
              </h2>
              <p className="text-sm text-muted-foreground mb-10">
                10 modules taking you from absolute beginner to deploying production apps.
              </p>

              <div className="space-y-6">
                {modules.map((mod) => (
                  <div
                    key={mod.number}
                    id={`module-${mod.number}`}
                    className="rounded-xl border border-border/40 bg-card/30 overflow-hidden scroll-mt-20"
                  >
                    {/* Module header */}
                    <div className="flex items-start gap-4 p-5 border-b border-border/30 bg-card/50">
                      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15">
                        <mod.icon className="h-5 w-5 text-primary" />
                      </div>
                      <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span className="text-[10px] font-mono font-semibold text-primary/60 uppercase tracking-wider">
                            Module {mod.number}
                          </span>
                          <span className="text-[10px] text-muted-foreground/50">&middot;</span>
                          <span className="flex items-center gap-1 text-[10px] text-muted-foreground/50">
                            <Clock className="h-3 w-3" />
                            {mod.duration}
                          </span>
                        </div>
                        <h3 className="text-lg font-semibold tracking-tight">
                          {mod.title}
                        </h3>
                        <p className="text-[13px] text-muted-foreground/70 mt-1 leading-relaxed">
                          {mod.description}
                        </p>
                        {mod.project && (
                          <div className="mt-2 inline-flex items-center gap-1.5 rounded-full bg-primary/10 border border-primary/15 px-2.5 py-0.5 text-[10px] font-medium text-primary">
                            <Trophy className="h-3 w-3" />
                            Project: {mod.project}
                          </div>
                        )}
                      </div>
                    </div>

                    {/* Lessons */}
                    <div className="divide-y divide-border/20">
                      {mod.lessons.map((lesson, i) => (
                        <div key={i} className="px-5 py-3 hover:bg-accent/20 transition-colors">
                          <div className="flex items-start gap-3">
                            <div className="flex h-5 w-5 shrink-0 items-center justify-center rounded text-[10px] font-mono text-muted-foreground/50 bg-accent/30 mt-0.5">
                              {i + 1}
                            </div>
                            <div className="min-w-0">
                              <div className="text-[13px] font-medium text-foreground/90">
                                {lesson.title}
                              </div>
                              {lesson.assignment && (
                                <div className="flex items-start gap-1.5 mt-1">
                                  <CheckCircle2 className="h-3 w-3 text-primary/50 shrink-0 mt-0.5" />
                                  <span className="text-[11px] text-muted-foreground/60 leading-relaxed">
                                    {lesson.assignment}
                                  </span>
                                </div>
                              )}
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>

        {/* Waitlist CTA */}
        <section className="py-16 border-t border-border/30">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-xl mx-auto text-center">
              <GraduationCap className="h-8 w-8 text-primary mx-auto mb-4" />
              <h2 className="text-2xl font-semibold tracking-tight mb-3">
                Join the waitlist
              </h2>
              <p className="text-sm text-muted-foreground mb-6">
                Be the first to know when the course launches. Enter your name and
                email below and we&apos;ll notify you on launch day.
              </p>
              <WaitlistForm />
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
