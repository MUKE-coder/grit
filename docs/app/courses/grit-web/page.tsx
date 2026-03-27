import {
  BookOpen,
  Clock,
  Trophy,
  Terminal,
  Wand2,
  ShieldCheck,
  LayoutDashboard,
  HardDrive,
  Mail,
  Sparkles,
  Rocket,
  ChevronDown,
  ArrowLeft,
} from "lucide-react"
import Link from "next/link"
import { SiteHeader } from "@/components/site-header"
import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Grit Web — Building Web Applications with Grit",
  description:
    "8 self-paced mini-courses (~4 hours total) covering everything from scaffolding your first Grit app to deploying to production.",
  openGraph: {
    title: "Grit Web — Building Web Applications with Grit",
    description:
      "8 self-paced mini-courses (~4 hours total) covering everything from scaffolding your first Grit app to deploying to production.",
    url: "https://gritframework.dev/courses/grit-web",
  },
}

/* -------------------------------------------------------------------------- */
/*  Types                                                                     */
/* -------------------------------------------------------------------------- */

interface Lesson {
  name: string
  explanation: string
  keyPoints: string[]
}

interface MiniCourse {
  number: number
  title: string
  duration: string
  icon: React.ComponentType<{ className?: string }>
  overview: string
  lessons: Lesson[]
  challenge: {
    title: string
    description: string
  }
}

/* -------------------------------------------------------------------------- */
/*  Course data                                                               */
/* -------------------------------------------------------------------------- */

const courses: MiniCourse[] = [
  {
    number: 1,
    title: "Your First Grit App",
    duration: "30 min",
    icon: Terminal,
    overview:
      "Install Grit, scaffold your first project, understand the project structure, and run the development servers.",
    lessons: [
      {
        name: "Installing Grit",
        explanation:
          "Set up every prerequisite and install the Grit CLI globally.",
        keyPoints: [
          "Prerequisites: Go 1.21+, Node 18+, pnpm, Docker",
          "go install github.com/MUKE-coder/grit@latest",
          "Verify with grit version",
        ],
      },
      {
        name: "Scaffolding a Project",
        explanation:
          "Use the grit new command to generate a complete full-stack project from interactive prompts.",
        keyPoints: [
          "grit new myapp",
          "Interactive prompts: architecture, frontend, style",
          "Project created with all boilerplate ready",
        ],
      },
      {
        name: "Project Structure",
        explanation:
          "Understand every top-level directory and where each piece of your application lives.",
        keyPoints: [
          "apps/api — Go backend (Gin + GORM)",
          "apps/web — Next.js frontend",
          "apps/admin — Admin panel",
          "packages/shared — TypeScript types & Zod schemas",
        ],
      },
      {
        name: "Running the App",
        explanation:
          "Start Docker services, install dependencies, and launch every dev server.",
        keyPoints: [
          "docker compose up -d",
          "pnpm install && pnpm dev",
          "Visit localhost:3000 (web), localhost:3001 (admin), localhost:8080 (API)",
        ],
      },
      {
        name: "Exploring Built-in Tools",
        explanation:
          "Tour the developer tools that come pre-installed with every Grit project.",
        keyPoints: [
          "GORM Studio (/studio) — visual database browser",
          "API Docs (/docs) — auto-generated endpoint reference",
          "Pulse (/pulse/ui) — real-time event monitor",
          "Sentinel (/sentinel/ui) — rate limiting dashboard",
        ],
      },
    ],
    challenge: {
      title: "Build a Bookstore",
      description:
        'Scaffold a project called "bookstore", start all services, register an account, and explore every built-in tool.',
    },
  },
  {
    number: 2,
    title: "Code Generator Mastery",
    duration: "30 min",
    icon: Wand2,
    overview:
      "Master the grit generate command to create full-stack CRUD resources with various field types and relationships.",
    lessons: [
      {
        name: "Basic Generation",
        explanation:
          "Generate a complete CRUD resource from a single CLI command.",
        keyPoints: [
          'grit generate resource Book --fields "title:string,author:string,published:bool"',
          "Creates Go model, service, handler, routes, Zod schema, TS types, React hooks, and admin page",
        ],
      },
      {
        name: "Field Types",
        explanation:
          "Learn when to use each supported field type for accurate data modeling.",
        keyPoints: [
          "string — short text (titles, names)",
          "text — longer text, richtext — HTML content",
          "int, float — numbers; bool — true/false",
          "datetime, date — timestamps and calendar dates",
        ],
      },
      {
        name: "Modifiers",
        explanation:
          "Apply constraints and behaviors to fields using modifier suffixes.",
        keyPoints: [
          ":unique — adds a unique database index",
          ":optional — makes the field nullable",
          "Combine them: email:string:unique:optional",
        ],
      },
      {
        name: "Relationships",
        explanation:
          "Model connections between resources using foreign keys and junction tables.",
        keyPoints: [
          "belongs_to:Category — adds category_id foreign key",
          "many_to_many:Tag — creates a junction table",
        ],
      },
      {
        name: "Slug Fields",
        explanation:
          "Auto-generate URL-friendly slugs from any source field.",
        keyPoints: [
          ":slug:source — e.g. slug:slug:title",
          "Automatically generates URL-friendly strings",
          "Unique by default for safe URLs",
        ],
      },
      {
        name: "Interactive Mode",
        explanation:
          "Build resources step-by-step with a guided field builder.",
        keyPoints: [
          "grit generate resource -i",
          "Prompts for each field name, type, and modifiers",
          "Great for learning or complex resources",
        ],
      },
      {
        name: "Removing Resources",
        explanation:
          "Cleanly remove a generated resource and reverse all code injections.",
        keyPoints: [
          "grit remove resource Book",
          "Deletes all generated files",
          "Reverses route, import, and registry injections",
        ],
      },
      {
        name: "Type Sync",
        explanation:
          "Regenerate TypeScript types whenever Go models change.",
        keyPoints: [
          "grit sync",
          "Parses Go structs and outputs Zod schemas + TS interfaces",
          "Keeps frontend and backend perfectly aligned",
        ],
      },
    ],
    challenge: {
      title: "Build a Blog System",
      description:
        "Create a blog system with Category (name, slug), Post (title, slug:title, content:richtext, belongs_to:Category, published:bool), and Tag (name) with Post many_to_many Tag.",
    },
  },
  {
    number: 3,
    title: "Authentication & Authorization",
    duration: "30 min",
    icon: ShieldCheck,
    overview:
      "Understand the complete auth system: JWT tokens, roles, 2FA with TOTP, OAuth2, and custom roles.",
    lessons: [
      {
        name: "JWT Flow",
        explanation:
          "Walk through the full authentication lifecycle from registration to token refresh.",
        keyPoints: [
          "Register -> Login -> access token + refresh token",
          "Tokens stored in HTTP-only cookies",
          "Automatic refresh when access token expires",
        ],
      },
      {
        name: "Role-Based Access",
        explanation:
          "Restrict endpoints to specific user roles using built-in middleware.",
        keyPoints: [
          "Built-in roles: ADMIN, EDITOR, USER",
          "middleware.RequireRole() for route protection",
          "Role checked on every authenticated request",
        ],
      },
      {
        name: "Two-Factor Auth (TOTP)",
        explanation:
          "Add an extra layer of security with time-based one-time passwords.",
        keyPoints: [
          "Setup endpoint returns QR code for authenticator app",
          "Verification with 6-digit TOTP code",
          "Backup codes for recovery, trusted device support",
        ],
      },
      {
        name: "OAuth2 Social Login",
        explanation:
          "Let users sign in with their existing Google or GitHub accounts.",
        keyPoints: [
          "Supports Google + GitHub out of the box",
          "Configure client IDs in .env",
          "Automatic account linking by email address",
        ],
      },
      {
        name: "Custom Roles",
        explanation:
          "Define new roles beyond the built-in three using the CLI.",
        keyPoints: [
          "grit add role MODERATOR",
          "Updates Go constants, middleware, and TypeScript types",
          "Immediately usable in route protection",
        ],
      },
      {
        name: "Protecting Routes",
        explanation:
          "Apply role restrictions during resource generation or manually in routes.go.",
        keyPoints: [
          "--roles flag during grit generate",
          "Role middleware applied in routes.go",
          "Both API and admin panel respect the same roles",
        ],
      },
    ],
    challenge: {
      title: "Secure Your App",
      description:
        "Set up 2FA on your account using Google Authenticator, test backup codes, add a MODERATOR role, and create a resource that only ADMIN and MODERATOR can access.",
    },
  },
  {
    number: 4,
    title: "Admin Panel Customization",
    duration: "30 min",
    icon: LayoutDashboard,
    overview:
      "Customize the admin panel: dashboard widgets, DataTable configuration, FormBuilder, multi-step forms, and style variants.",
    lessons: [
      {
        name: "Resource Definitions",
        explanation:
          "Configure how each resource appears in the admin panel using defineResource().",
        keyPoints: [
          "defineResource() pattern with columns, form fields, filters",
          "Declarative configuration — no custom UI code needed",
          "Supports all field types and relationships",
        ],
      },
      {
        name: "DataTable",
        explanation:
          "Configure the data listing table with sorting, filtering, and export.",
        keyPoints: [
          "Pagination, sorting, and filtering built-in",
          "Column visibility toggles",
          "CSV and JSON export with one click",
        ],
      },
      {
        name: "FormBuilder",
        explanation:
          "Define forms declaratively with support for many field types.",
        keyPoints: [
          "Field types: text, select, textarea, date, richtext, file, relationship-select",
          "Validation rules defined alongside fields",
          "Automatic layout and responsive design",
        ],
      },
      {
        name: "Multi-Step Forms",
        explanation:
          "Split long forms into steps with per-step validation.",
        keyPoints: [
          'formView: "modal-steps" for modal wizards',
          'formView: "page-steps" for full-page wizards',
          "Each step validates independently before advancing",
        ],
      },
      {
        name: "Dashboard Widgets",
        explanation:
          "Add stats cards, charts, and activity feeds to the admin dashboard.",
        keyPoints: [
          "Stats cards with trends and icons",
          "Line charts and bar charts for time-series data",
          "Activity feed for recent actions",
        ],
      },
      {
        name: "Style Variants",
        explanation:
          "Switch between pre-built visual themes for the admin panel.",
        keyPoints: [
          "Variants: default, modern, minimal, glass",
          "One config change to switch styles globally",
          "Consistent with the Grit design system",
        ],
      },
      {
        name: "Standalone Usage",
        explanation:
          "Use DataTable and FormBuilder outside the resource system for custom pages.",
        keyPoints: [
          "Import DataTable and FormBuilder directly",
          "Pass data and column/field configs as props",
          "Full flexibility for custom admin pages",
        ],
      },
    ],
    challenge: {
      title: "Build an Invoice System",
      description:
        "Create an Invoice resource with a multi-step form (Step 1: client info, Step 2: line items, Step 3: review). Add dashboard stats cards showing total revenue and invoice count.",
    },
  },
  {
    number: 5,
    title: "File Storage & Uploads",
    duration: "30 min",
    icon: HardDrive,
    overview:
      "Configure file storage (MinIO, S3, R2), implement presigned URL uploads, and process images with background jobs.",
    lessons: [
      {
        name: "Storage Configuration",
        explanation:
          "Set up the storage driver and credentials in your environment.",
        keyPoints: [
          "STORAGE_DRIVER env var (minio, s3, r2)",
          "Endpoint, bucket, access key, and secret key",
          "Shared config across all upload endpoints",
        ],
      },
      {
        name: "MinIO for Dev",
        explanation:
          "Use MinIO as a local S3-compatible storage service during development.",
        keyPoints: [
          "Included in docker-compose.yml",
          "Web console at localhost:9001",
          "Identical API to AWS S3",
        ],
      },
      {
        name: "Presigned URL Flow",
        explanation:
          "Upload files directly from the browser to S3 without proxying through your API.",
        keyPoints: [
          "Browser requests a presigned URL from the API",
          "Browser uploads directly to S3/MinIO",
          "API records the file metadata in the database",
        ],
      },
      {
        name: "Upload Progress",
        explanation:
          "Show real-time upload progress to users in the frontend.",
        keyPoints: [
          "XHR progress event tracking",
          "Progress bar component included",
          "Works with any S3-compatible storage",
        ],
      },
      {
        name: "Image Processing",
        explanation:
          "Automatically generate thumbnails using background jobs.",
        keyPoints: [
          "Upload triggers a background job",
          "Thumbnail generated at configured sizes",
          "Original and thumbnail URLs stored together",
        ],
      },
      {
        name: "Cloud Storage",
        explanation:
          "Switch from local MinIO to production cloud storage.",
        keyPoints: [
          "Change STORAGE_DRIVER and credentials in .env",
          "Supports Cloudflare R2 and AWS S3",
          "Zero code changes required",
        ],
      },
    ],
    challenge: {
      title: "Build a Photo Gallery",
      description:
        "Create a Photo Gallery resource with title and image upload. Upload 5 images, verify thumbnails are generated, then switch storage driver to Cloudflare R2.",
    },
  },
  {
    number: 6,
    title: "Background Jobs & Email",
    duration: "30 min",
    icon: Mail,
    overview:
      "Use Redis-backed background jobs for async processing and send transactional emails with Resend.",
    lessons: [
      {
        name: "asynq Overview",
        explanation:
          "Understand the Redis-based job queue that powers async processing in Grit.",
        keyPoints: [
          "Redis job queue with worker pools",
          "Automatic retry logic with exponential backoff",
          "Dashboard for monitoring job status",
        ],
      },
      {
        name: "Built-in Jobs",
        explanation:
          "Tour the jobs that come pre-configured in every Grit project.",
        keyPoints: [
          "Email sending (transactional and bulk)",
          "Image processing (thumbnail generation)",
          "Cleanup workers (expired tokens, temp files)",
        ],
      },
      {
        name: "Creating Custom Jobs",
        explanation:
          "Define, register, and dispatch your own background jobs.",
        keyPoints: [
          "Define a task type and payload struct",
          "Register a handler function in the worker",
          "Dispatch from any service with asynq client",
        ],
      },
      {
        name: "Cron Scheduler",
        explanation:
          "Schedule recurring tasks using cron expressions.",
        keyPoints: [
          "Standard cron expression syntax",
          'e.g. "0 0 * * *" for daily at midnight',
          "Managed through the admin cron dashboard",
        ],
      },
      {
        name: "Email Service",
        explanation:
          "Send transactional emails using the Resend API.",
        keyPoints: [
          "Configure MAIL_FROM and RESEND_API_KEY in .env",
          "Send emails from any service or job handler",
          "HTML templates with dynamic data",
        ],
      },
      {
        name: "Email Templates",
        explanation:
          "Use and customize the built-in email templates.",
        keyPoints: [
          "Welcome, password reset, verification, notification",
          "HTML templates with inline styles",
          "Easy to add custom templates",
        ],
      },
      {
        name: "Mailhog for Dev",
        explanation:
          "Catch and preview all emails locally during development.",
        keyPoints: [
          "Included in docker-compose.yml",
          "Web UI at localhost:8025",
          "All emails intercepted — nothing sent externally",
        ],
      },
      {
        name: "Admin Dashboard",
        explanation:
          "Monitor jobs and scheduled tasks from the admin panel.",
        keyPoints: [
          "Job status: pending, active, completed, failed",
          "Retry counts and failure reasons",
          "Cron task schedule and history",
        ],
      },
    ],
    challenge: {
      title: "Daily Report Job",
      description:
        "Create a custom background job that generates a daily report email. Use cron to schedule it at midnight. Test with Mailhog.",
    },
  },
  {
    number: 7,
    title: "AI-Powered Features",
    duration: "30 min",
    icon: Sparkles,
    overview:
      "Add AI to your app using Vercel AI Gateway. Build completions, chat, and streaming endpoints.",
    lessons: [
      {
        name: "Vercel AI Gateway",
        explanation:
          "Use one API key to access hundreds of AI models through a unified gateway.",
        keyPoints: [
          "Single API key for all providers",
          "Provider/model format (e.g. anthropic/claude-sonnet-4-6)",
          "Automatic fallback and load balancing",
        ],
      },
      {
        name: "Configuration",
        explanation:
          "Set up your AI gateway credentials and default model.",
        keyPoints: [
          "AI_GATEWAY_API_KEY in .env",
          "AI_GATEWAY_MODEL for the default model",
          "AI_GATEWAY_URL for the gateway endpoint",
        ],
      },
      {
        name: "Completion Endpoint",
        explanation:
          "Build a simple prompt-in, response-out endpoint.",
        keyPoints: [
          "POST /api/ai/complete",
          "Send a prompt, receive a full response",
          "Good for one-shot tasks like summarization",
        ],
      },
      {
        name: "Chat Endpoint",
        explanation:
          "Build multi-turn conversation endpoints with message history.",
        keyPoints: [
          "POST /api/ai/chat",
          "Send message array with roles (user, assistant, system)",
          "Maintains conversation context",
        ],
      },
      {
        name: "Streaming",
        explanation:
          "Stream AI responses in real-time using Server-Sent Events.",
        keyPoints: [
          "POST /api/ai/stream",
          "Server-Sent Events (SSE) for real-time tokens",
          "Much better UX for long responses",
        ],
      },
      {
        name: "Switching Models",
        explanation:
          "Swap AI providers by changing a single configuration string.",
        keyPoints: [
          "Change AI_GATEWAY_MODEL in .env",
          "e.g. anthropic/claude-sonnet-4-6 -> openai/gpt-5.4",
          "No code changes required",
        ],
      },
      {
        name: "Building a Chat UI",
        explanation:
          "Create a React component that consumes an SSE stream for real-time chat.",
        keyPoints: [
          "EventSource API or fetch with ReadableStream",
          "Token-by-token rendering in the UI",
          "Loading states and error handling",
        ],
      },
    ],
    challenge: {
      title: "Product Description Generator",
      description:
        'Build a product description generator — user enters product name and features, AI generates a marketing description. Add a "regenerate" button that tries a different model.',
    },
  },
  {
    number: 8,
    title: "Deploy to Production",
    duration: "30 min",
    icon: Rocket,
    overview:
      "Deploy your Grit app to a VPS with one command using grit deploy. Configure systemd, Caddy, and HTTPS.",
    lessons: [
      {
        name: "Prerequisites",
        explanation:
          "What you need before deploying: a server and a domain.",
        keyPoints: [
          "VPS from DigitalOcean, Hetzner, or similar ($5/mo)",
          "SSH access configured with your public key",
          "Domain name pointed to your server IP",
        ],
      },
      {
        name: "The Deploy Command",
        explanation:
          "Deploy your entire app with a single CLI command.",
        keyPoints: [
          "grit deploy --host <ip> --domain <domain>",
          "Handles build, upload, and service configuration",
          "Idempotent — safe to run multiple times",
        ],
      },
      {
        name: "Build Pipeline",
        explanation:
          "How Grit builds your app for production.",
        keyPoints: [
          "Cross-compiles Go binary for linux/amd64",
          "Builds Next.js frontend with optimizations",
          "Single binary + static assets for deployment",
        ],
      },
      {
        name: "Upload",
        explanation:
          "How the built artifacts get to your server.",
        keyPoints: [
          "SCP binary to /opt/myapp/",
          "Frontend assets served by Caddy or embedded",
          "Fast incremental updates on re-deploy",
        ],
      },
      {
        name: "systemd Service",
        explanation:
          "Run your app as a managed system service with auto-restart.",
        keyPoints: [
          "Auto-restart on crash or reboot",
          "Environment file for secrets",
          "Security hardening (read-only filesystem, no escalation)",
        ],
      },
      {
        name: "Caddy Reverse Proxy",
        explanation:
          "Automatic HTTPS and reverse proxy with Caddy.",
        keyPoints: [
          "Automatic Let's Encrypt TLS certificates",
          "Security headers (HSTS, CSP, etc.)",
          "Access logging and request buffering",
        ],
      },
      {
        name: "Environment Variables",
        explanation:
          "Configure your production environment securely.",
        keyPoints: [
          "Production .env with real credentials",
          "Database URL, Redis URL, storage keys",
          "Secrets never committed to version control",
        ],
      },
      {
        name: "Updates",
        explanation:
          "Re-deploy updates with zero downtime.",
        keyPoints: [
          "Same grit deploy command for updates",
          "systemd graceful restart — no dropped connections",
          "Rollback by deploying a previous version",
        ],
      },
      {
        name: "Docker Alternative",
        explanation:
          "Deploy with Docker Compose instead of bare metal.",
        keyPoints: [
          "docker-compose.prod.yml included in project",
          "All services containerized",
          "Good for teams already using Docker in production",
        ],
      },
      {
        name: "Maintenance Mode",
        explanation:
          "Take your app offline gracefully during major updates.",
        keyPoints: [
          "grit down before deploy",
          "Shows maintenance page to users",
          "grit up to restore service",
        ],
      },
    ],
    challenge: {
      title: "Deploy Your Bookstore",
      description:
        "Get a $5 VPS, point a domain to it, and deploy your bookstore app. Verify HTTPS works, test maintenance mode during a re-deploy.",
    },
  },
]

/* -------------------------------------------------------------------------- */
/*  Components                                                                */
/* -------------------------------------------------------------------------- */

function StatBadge({
  icon: Icon,
  label,
  value,
}: {
  icon: React.ComponentType<{ className?: string }>
  label: string
  value: string
}) {
  return (
    <div className="flex items-center gap-2 rounded-lg bg-white/[0.04] border border-white/[0.06] px-4 py-2.5">
      <Icon className="h-4 w-4 text-primary" />
      <div className="text-sm">
        <span className="font-semibold text-foreground">{value}</span>{" "}
        <span className="text-muted-foreground">{label}</span>
      </div>
    </div>
  )
}

function LessonItem({
  index,
  lesson,
}: {
  index: number
  lesson: Lesson
}) {
  return (
    <div className="flex gap-4">
      <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-md bg-primary/10 text-xs font-mono font-semibold text-primary mt-0.5">
        {index}
      </div>
      <div className="flex-1 min-w-0">
        <h4 className="text-sm font-semibold text-foreground mb-1">
          {lesson.name}
        </h4>
        <p className="text-sm text-muted-foreground leading-relaxed mb-2">
          {lesson.explanation}
        </p>
        <ul className="space-y-1">
          {lesson.keyPoints.map((point, i) => (
            <li
              key={i}
              className="flex items-start gap-2 text-[13px] text-muted-foreground/80"
            >
              <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-primary/40" />
              <span>{point}</span>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}

function CourseSection({ course }: { course: MiniCourse }) {
  const Icon = course.icon

  return (
    <details className="group rounded-xl border border-border/40 bg-card/50 overflow-hidden">
      <summary className="flex cursor-pointer items-center gap-4 px-6 py-5 hover:bg-white/[0.02] transition-colors list-none [&::-webkit-details-marker]:hidden">
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15">
          <Icon className="h-5 w-5 text-primary" />
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 flex-wrap">
            <h3 className="text-base font-semibold text-foreground">
              {course.number}. {course.title}
            </h3>
            <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 border border-primary/20 px-2.5 py-0.5 text-xs font-mono font-medium text-primary">
              <Clock className="h-3 w-3" />
              {course.duration}
            </span>
          </div>
          <p className="text-sm text-muted-foreground mt-1 line-clamp-1">
            {course.overview}
          </p>
        </div>

        <ChevronDown className="h-5 w-5 shrink-0 text-muted-foreground transition-transform group-open:rotate-180" />
      </summary>

      <div className="border-t border-border/30 px-6 py-6">
        {/* Overview */}
        <p className="text-sm text-muted-foreground leading-relaxed mb-6">
          {course.overview}
        </p>

        {/* Lessons */}
        <div className="space-y-5 mb-8">
          {course.lessons.map((lesson, i) => (
            <LessonItem key={i} index={i + 1} lesson={lesson} />
          ))}
        </div>

        {/* DIY Challenge */}
        <div className="rounded-lg bg-primary/10 border border-primary/20 p-5">
          <div className="flex items-center gap-2 mb-2">
            <Trophy className="h-4 w-4 text-primary" />
            <h4 className="text-sm font-semibold text-primary">
              DIY Challenge: {course.challenge.title}
            </h4>
          </div>
          <p className="text-sm text-muted-foreground leading-relaxed">
            {course.challenge.description}
          </p>
        </div>
      </div>
    </details>
  )
}

/* -------------------------------------------------------------------------- */
/*  Page                                                                      */
/* -------------------------------------------------------------------------- */

export default function GritWebCoursePage() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main>
        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/[0.04] to-transparent" />
          <div className="container max-w-screen-xl relative py-20 px-6">
            <div className="max-w-3xl mx-auto text-center">
              <Link
                href="/courses"
                className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-6"
              >
                <ArrowLeft className="h-4 w-4" />
                All Courses
              </Link>

              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6 ml-4">
                <BookOpen className="h-3.5 w-3.5" />
                Grit Web
              </div>

              <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1] text-foreground">
                Building Web Applications with Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl mx-auto mb-8">
                8 self-paced mini-courses that take you from zero to a deployed,
                production-ready web application. Each course is{" "}
                <span className="text-foreground font-medium">~30 minutes</span>{" "}
                and ends with a hands-on challenge.
              </p>

              {/* Stats */}
              <div className="flex flex-wrap items-center justify-center gap-3">
                <StatBadge icon={BookOpen} value="8" label="courses" />
                <StatBadge icon={Clock} value="~4 hrs" label="total" />
                <StatBadge icon={Trophy} value="8" label="challenges" />
              </div>
            </div>
          </div>
        </section>

        {/* Course list */}
        <section className="py-16">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-3xl mx-auto space-y-4">
              {courses.map((course) => (
                <CourseSection key={course.number} course={course} />
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
