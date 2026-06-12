/**
 * Single source of truth for the 7 kit-based learning paths.
 *
 * Routes are derived from this metadata:
 *   /courses/<course.slug>
 *   /courses/<course.slug>/<chapter.slug>
 *   /courses/<course.slug>/<chapter.slug>/<lesson.slug>
 *   /courses/<course.slug>/<chapter.slug>/assignment
 *
 * Lesson CONTENT lives in
 *   app/courses/<course.slug>/<chapter.slug>/<lesson.slug>/page.tsx
 * — this file only describes structure + metadata so the sidebar, sitemap,
 * search dialog, prev/next nav, and landing pages all stay in sync.
 */

export type Level = 'absolute-beginner' | 'beginner' | 'intermediate' | 'advanced'
export type Difficulty = 'easy' | 'medium' | 'hard'
export type Status = 'available' | 'coming-soon'

export interface Lesson {
  slug: string
  title: string
  tagline: string
  minutes: number
  difficulty: Difficulty
  status: Status
}

export interface Module {
  /** Visual grouping inside a chapter (sidebar header). Not a route. */
  title: string
  lessons: Lesson[]
}

export interface Assignment {
  title: string
  brief: string
  successCriteria: string[]
}

export interface Chapter {
  slug: string
  number: number
  title: string
  tagline: string
  learningGoals: string[]
  modules: Module[]
  assignment?: Assignment
  status: Status
}

export interface Course {
  slug: string
  name: string
  shortName: string
  tagline: string
  description: string
  level: Level
  /** Approx total lesson minutes — derived. */
  estimatedHours: number
  prerequisites: string[]
  whatYoullBuild: string
  whatYoullLearn: string[]
  whoThisIsFor: string[]
  status: Status
  /** Optional emoji shown on cards. */
  emoji: string
  /** Tailwind gradient string for hero accent. */
  accent: string
  chapters: Chapter[]
}

/* ─────────────────────────────────────────────────────────────────── */
/*  Helpers                                                            */
/* ─────────────────────────────────────────────────────────────────── */

export function flatLessons(course: Course): Array<{
  chapter: Chapter
  module: Module
  lesson: Lesson
  /** Sequential index across the whole course (1-based). */
  index: number
}> {
  const out: Array<{ chapter: Chapter; module: Module; lesson: Lesson; index: number }> = []
  let i = 0
  for (const ch of course.chapters) {
    for (const mod of ch.modules) {
      for (const l of mod.lessons) {
        i += 1
        out.push({ chapter: ch, module: mod, lesson: l, index: i })
      }
    }
  }
  return out
}

export function courseTotalMinutes(course: Course): number {
  return flatLessons(course).reduce((acc, x) => acc + x.lesson.minutes, 0)
}

export function findCourse(slug: string): Course | undefined {
  return COURSES.find((c) => c.slug === slug)
}

export function findChapter(courseSlug: string, chapterSlug: string): { course: Course; chapter: Chapter } | undefined {
  const course = findCourse(courseSlug)
  if (!course) return undefined
  const chapter = course.chapters.find((c) => c.slug === chapterSlug)
  if (!chapter) return undefined
  return { course, chapter }
}

export function findLesson(
  courseSlug: string,
  chapterSlug: string,
  lessonSlug: string,
): { course: Course; chapter: Chapter; module: Module; lesson: Lesson; index: number; prev?: { chapterSlug: string; lessonSlug: string }; next?: { chapterSlug: string; lessonSlug: string } } | undefined {
  const course = findCourse(courseSlug)
  if (!course) return undefined
  const all = flatLessons(course)
  const i = all.findIndex((x) => x.chapter.slug === chapterSlug && x.lesson.slug === lessonSlug)
  if (i === -1) return undefined
  const { chapter, module, lesson, index } = all[i]
  return {
    course,
    chapter,
    module,
    lesson,
    index,
    prev: i > 0 ? { chapterSlug: all[i - 1].chapter.slug, lessonSlug: all[i - 1].lesson.slug } : undefined,
    next: i < all.length - 1 ? { chapterSlug: all[i + 1].chapter.slug, lessonSlug: all[i + 1].lesson.slug } : undefined,
  }
}

/* ─────────────────────────────────────────────────────────────────── */
/*  1. Grit Concepts You Need to Know — FULL CONTENT                   */
/* ─────────────────────────────────────────────────────────────────── */

const conceptsCourse: Course = {
  slug: 'concepts',
  name: 'Grit Concepts You Need to Know',
  shortName: 'Grit Concepts',
  tagline: 'The non-negotiables that make every Grit project click — start here.',
  description:
    "Every other course assumes this one. We cover what Grit is, the directory shape every project uses, the conventions you'll see hundreds of times, and the code generator that does 80% of the typing for you. By the end you can scaffold a project, add a resource, and explain what every folder is for.",
  level: 'absolute-beginner',
  estimatedHours: 4,
  prerequisites: [
    'Comfortable in a terminal (cd, ls, running commands)',
    'Have used Git at least once',
    'Know what HTTP and JSON are — no need to be an expert',
  ],
  whatYoullBuild:
    "A working monorepo with Go API + Next.js web + admin panel, a Product resource you scaffolded with one command, and confidence to navigate any Grit codebase.",
  whatYoullLearn: [
    'What Grit is and when to reach for it',
    'How to install Grit and scaffold a project',
    'The folder structure every project uses — and why',
    'Grit\'s naming conventions for files, routes, models, and types',
    'The standard API response envelope',
    'The code generator: grit generate resource',
    'Type sync: grit sync (Go → TypeScript)',
    'The five architecture modes and when to pick which',
  ],
  whoThisIsFor: [
    'Devs new to Grit who want the foundation before specializing',
    'Anyone inheriting a Grit codebase and wanting orientation',
    'Self-taught devs who want a structured walk through a modern full-stack',
  ],
  status: 'available',
  emoji: '🧭',
  accent: 'from-violet-500/30 via-fuchsia-500/20 to-violet-500/10',
  chapters: [
    {
      slug: 'welcome',
      number: 1,
      title: 'What is Grit and Why?',
      tagline: 'The 10-minute orientation: what Grit is, who it\'s for, and what makes it different.',
      learningGoals: [
        'Explain what Grit is in 30 seconds',
        'Identify whether Grit fits your next project',
        'Install Grit and verify it works',
      ],
      status: 'available',
      modules: [
        {
          title: 'The big picture',
          lessons: [
            {
              slug: 'what-is-grit',
              title: 'What is Grit?',
              tagline: 'A full-stack meta-framework for Go + React, batteries included.',
              minutes: 6,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'when-to-use-grit',
              title: 'When to use Grit (and when not to)',
              tagline: 'Honest trade-offs vs. plain Next.js, Rails, Laravel, and others.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
        {
          title: 'Setup',
          lessons: [
            {
              slug: 'install-grit',
              title: 'Installing Grit',
              tagline: 'One-line install with the script, or `go install` if you have Go.',
              minutes: 4,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'verify-install',
              title: 'Verifying your install',
              tagline: 'grit version, grit --help, and how to update later.',
              minutes: 3,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Install Grit and capture proof',
        brief:
          "Install Grit on your machine using the one-line install script. Run `grit version` and paste the output into a Markdown file called `notes.md` you'll grow throughout the course.",
        successCriteria: [
          '`grit version` prints v3.25.x or later',
          '`grit --help` shows the command list',
          'You have a notes.md file with your version captured',
        ],
      },
    },
    {
      slug: 'first-project',
      number: 2,
      title: 'Your First Grit Project',
      tagline: 'Scaffold a real project, run it, and tour every file it produced.',
      learningGoals: [
        'Run `grit new` confidently with the right flags',
        'Identify every top-level folder and what lives in it',
        'Start the dev servers and see the app in your browser',
      ],
      status: 'available',
      modules: [
        {
          title: 'Scaffolding',
          lessons: [
            {
              slug: 'grit-new',
              title: 'Running grit new',
              tagline: 'The interactive prompts, the flags, and what each kit means.',
              minutes: 7,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'project-tour',
              title: 'Tour of your project',
              tagline: 'apps/, packages/, root config — what every folder is for.',
              minutes: 8,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
        {
          title: 'Running it',
          lessons: [
            {
              slug: 'dev-servers',
              title: 'Starting the dev servers',
              tagline: 'docker compose up, the API, web, admin — running everything together.',
              minutes: 6,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'first-look',
              title: 'Your first look',
              tagline: 'Open the URLs, log in with the seeded account, click around.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Scaffold + screenshot tour',
        brief:
          'Scaffold a project called `my-first-grit`, get all three apps running, log into the admin panel, and take screenshots of: the API health endpoint, the web app homepage, and the admin dashboard. Paste them into notes.md.',
        successCriteria: [
          'All three apps are reachable in your browser',
          'You logged into the admin panel successfully',
          'Three screenshots in notes.md',
        ],
      },
    },
    {
      slug: 'conventions',
      number: 3,
      title: 'The Convention Surface',
      tagline: 'The naming, structure, and response patterns Grit assumes — saving you decision fatigue.',
      learningGoals: [
        'Predict where a file should live before you create it',
        'Name things the way Grit expects (and why)',
        'Read any Grit API response without looking up the shape',
      ],
      status: 'available',
      modules: [
        {
          title: 'How files are organized',
          lessons: [
            {
              slug: 'folder-conventions',
              title: 'Folder conventions',
              tagline: 'internal/, packages/shared, apps/ — the rules and the rationale.',
              minutes: 7,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'naming-conventions',
              title: 'Naming conventions',
              tagline: 'snake_case Go, kebab-case TS, plural routes, plural tables.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
        {
          title: 'How data flows',
          lessons: [
            {
              slug: 'api-response-format',
              title: 'The API response envelope',
              tagline: '{ data, message } / { data, meta } / { error } — and HTTP status codes.',
              minutes: 6,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'error-handling',
              title: 'Error handling',
              tagline: 'Go: explicit + wrapped. React: error boundary + toast.',
              minutes: 6,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Add a field to the User model',
        brief:
          'Add a `bio` (string, nullable) field to the User model. Run the migration. Update the seeder. Verify it appears in the API response. Commit the change with a conventional commit message.',
        successCriteria: [
          'User.bio field exists in the DB',
          'GET /api/users returns the bio field',
          'You can edit bio from the admin panel',
        ],
      },
    },
    {
      slug: 'generators',
      number: 4,
      title: 'Code Generation & Type Sync',
      tagline: 'The commands that do 80% of the boilerplate so you focus on the actual product.',
      learningGoals: [
        'Generate a full resource end-to-end with one command',
        'Identify each of the 7 files that get generated and why',
        'Keep TypeScript types in sync with Go models',
      ],
      status: 'available',
      modules: [
        {
          title: 'Resources',
          lessons: [
            {
              slug: 'what-is-a-resource',
              title: 'What is a resource?',
              tagline: 'Model + handler + service + types + routes — the standard slice.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'grit-generate',
              title: 'Generating a Product resource',
              tagline: 'The command, the prompts, what each flag controls.',
              minutes: 7,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'what-got-generated',
              title: 'Touring what got generated',
              tagline: 'Read every file the generator dropped and connect them mentally.',
              minutes: 8,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
        {
          title: 'Keeping types aligned',
          lessons: [
            {
              slug: 'grit-sync',
              title: 'grit sync — Go to TypeScript',
              tagline: 'Why it exists, when to run it, what it does and doesn\'t handle.',
              minutes: 6,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Generate Order + items',
        brief:
          'Generate an `Order` resource with fields: customerId, total, status. Add an `OrderItem` resource that belongs to Order. Run `grit sync` so the frontend types match. Create one Order from the admin panel.',
        successCriteria: [
          'Both resources are generated and migrated',
          'TypeScript types reflect the new models',
          'You created an Order from the admin panel without errors',
        ],
      },
    },
    {
      slug: 'frameworks-patterns',
      number: 5,
      title: 'Frameworks & Patterns',
      tagline: 'The Go side demystified — Gin, GORM, CRUD, and the Handler → Service pattern that every resource follows.',
      learningGoals: [
        'Explain what Gin and GORM are, in one sentence each',
        'Walk a request from URL → handler → service → DB → response',
        'Write the four CRUD operations by hand and reason about why each line is there',
        'Decide whether new logic belongs in the handler or the service',
      ],
      status: 'available',
      modules: [
        {
          title: 'The frameworks',
          lessons: [
            {
              slug: 'gin-basics',
              title: 'What is Gin?',
              tagline: 'The HTTP router behind every Grit API — routes, context, middleware.',
              minutes: 8,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'gorm-basics',
              title: 'What is GORM?',
              tagline: 'The ORM that turns Go structs into SQL — models, queries, relations.',
              minutes: 9,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
        {
          title: 'The pattern',
          lessons: [
            {
              slug: 'handler-service-pattern',
              title: 'Handler → Service pattern',
              tagline: 'Why we split, what goes where, and how to keep handlers thin.',
              minutes: 9,
              difficulty: 'medium',
              status: 'available',
            },
            {
              slug: 'crud-walkthrough',
              title: 'CRUD end-to-end',
              tagline: 'GET / POST / PATCH / DELETE — written by hand, line by line.',
              minutes: 10,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Hand-write a Note resource',
        brief:
          "Without using `grit generate`, hand-write a Note resource: model, service, handler, routes. List notes, create one, update, delete. Then run `grit generate resource Note` in a fresh project and DIFF the two — note where your code differs from the generated. One paragraph in notes.md on what you learned.",
        successCriteria: [
          'All 4 CRUD endpoints work via curl/Postman',
          'Handler is thin (no DB calls)',
          'Service contains all DB work',
          'You can explain every line you wrote',
        ],
      },
    },
    {
      slug: 'batteries',
      number: 6,
      title: 'The Batteries',
      tagline: 'Redis, S3, Mail, Jobs, AI — what each one does, where it lives, and how to modify it.',
      learningGoals: [
        'Identify each battery, what it does, and when you need it',
        'Find the file you would edit to change each battery\'s behaviour',
        'Wire a new use of a battery into your project (e.g., add a cache, send a transactional email, queue a job)',
      ],
      status: 'available',
      modules: [
        {
          title: 'Caching & files',
          lessons: [
            {
              slug: 'batteries-overview',
              title: 'What are the Batteries?',
              tagline: 'The included-out-of-the-box services and what they save you.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'redis-cache',
              title: 'Redis cache',
              tagline: 'Speed up reads, throttle, and store ephemeral state.',
              minutes: 7,
              difficulty: 'medium',
              status: 'available',
            },
            {
              slug: 's3-storage',
              title: 'S3-compatible file storage',
              tagline: 'Upload files, store on R2/MinIO/AWS, serve with signed URLs.',
              minutes: 8,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
        {
          title: 'Async & integrations',
          lessons: [
            {
              slug: 'mail-resend',
              title: 'Mail (Resend)',
              tagline: 'Transactional + marketing email with templates you can edit.',
              minutes: 7,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'async-jobs',
              title: 'Background jobs (asynq)',
              tagline: 'Queue work, retry on failure, schedule with cron — without blocking requests.',
              minutes: 9,
              difficulty: 'medium',
              status: 'available',
            },
            {
              slug: 'ai-integration',
              title: 'AI (Claude + OpenAI)',
              tagline: 'Streaming chat, embeddings, and how the abstraction lets you swap providers.',
              minutes: 8,
              difficulty: 'medium',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Cache + email + job — a real use of three batteries',
        brief:
          "Build a small feature using three batteries: when a new Order is created, (1) cache the order total in Redis with a 5-min TTL, (2) enqueue a job that sends a confirmation email via Resend, (3) the job worker actually sends it. Test the whole loop with a real Resend test key.",
        successCriteria: [
          'Order create returns under 50ms (because email is async)',
          'Email lands in your Resend dashboard',
          'GET on the order pulls from Redis (you can see the X-Cache header)',
          'You can explain WHERE in the code each battery is wired',
        ],
      },
    },
    {
      slug: 'architecture-modes',
      number: 7,
      title: 'Architecture Modes',
      tagline: 'The five shapes a Grit project can take — pick the right one for your idea.',
      learningGoals: [
        'Explain each architecture mode in one sentence',
        'Pick the right kit for a given product idea',
        'Know which mode the rest of the courses cover',
      ],
      status: 'available',
      modules: [
        {
          title: 'The kits',
          lessons: [
            {
              slug: 'single-mode',
              title: 'Single — one Go binary',
              tagline: 'The smallest deploy: Go + embedded React. When less infra wins.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'triple-mode',
              title: 'Triple — web + admin + API',
              tagline: 'The SaaS shape. The recommended starter for most products.',
              minutes: 6,
              difficulty: 'easy',
              status: 'available',
            },
            {
              slug: 'specialized-modes',
              title: 'Specialized: API, Mobile, Desktop',
              tagline: 'When you want one slice — backend only, Expo mobile, or Wails desktop.',
              minutes: 6,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
        {
          title: 'Picking your shape',
          lessons: [
            {
              slug: 'choosing-a-kit',
              title: 'Choosing the right kit',
              tagline: 'A decision tree and the Stack Selector tool.',
              minutes: 5,
              difficulty: 'easy',
              status: 'available',
            },
          ],
        },
      ],
      assignment: {
        title: 'Pick the right kit for three ideas',
        brief:
          'For each of these product ideas, write down which architecture mode you\'d choose and why: (1) a market vendor POS that needs to work offline, (2) a public job-board with admin curation, (3) a backend-only service for a Discord bot. One paragraph each in notes.md.',
        successCriteria: [
          'Three ideas, three picks, three justifications',
          'No idea uses the same mode (challenge yourself)',
        ],
      },
    },
  ],
}

/* ─────────────────────────────────────────────────────────────────── */
/*  Coming-soon courses — structure complete, lesson content TBD        */
/* ─────────────────────────────────────────────────────────────────── */

const goApiCourse: Course = {
  slug: 'go-api',
  name: 'Building a Go API',
  shortName: 'Go API',
  tagline: 'From `grit new --api` to a deployed, secured, observable production API.',
  description:
    'Build a complete Go API on Grit: models, handlers, services, auth, jobs, mail, storage, observability, and deploy. The course that turns Grit Concepts into shippable backend skill.',
  level: 'beginner',
  estimatedHours: 8,
  prerequisites: ['Completed Grit Concepts', 'Comfortable writing a basic Go function'],
  whatYoullBuild: 'A multi-tenant SaaS API with auth, role-based access, file uploads, background jobs, transactional email, and a Pulse + Sentinel admin.',
  whatYoullLearn: [
    'Modeling with GORM — relations, constraints, soft-deletes',
    'JWT auth with refresh tokens + OAuth2 + 2FA',
    'Background jobs with Asynq',
    'File storage with S3 / R2',
    'Transactional mail with Resend',
    'Observability with Pulse + tamper-evident audit log',
    'Deploying to a VPS with grit deploy',
  ],
  whoThisIsFor: ['Devs who finished Grit Concepts', 'Backend devs new to Go but experienced elsewhere'],
  status: 'available',
  emoji: '🦫',
  accent: 'from-sky-500/30 via-cyan-500/20 to-sky-500/10',
  chapters: [
    {
      slug: 'scaffold-tour',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --api and a deep tour of the Go-only project shape.',
      learningGoals: ['Scaffold an API-only project', 'Identify every Go file', 'Run + curl /api/health'],
      status: 'available',
      modules: [
        {
          title: 'Start',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the API', tagline: 'grit new my-api --api walkthrough.', minutes: 5, difficulty: 'easy', status: 'available' },
            { slug: 'project-tour', title: 'Project tour', tagline: 'cmd/server, internal/, the Services pattern.', minutes: 8, difficulty: 'easy', status: 'available' },
            { slug: 'first-request', title: 'Your first request', tagline: 'docker compose up, run the API, hit /api/health.', minutes: 4, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Scaffold, run, and document',
        brief:
          'Scaffold `bench-api --api`, get it running, hit /api/health, and capture a short README in your notes.md describing what each top-level folder is for.',
        successCriteria: [
          '`grit new bench-api --api` ran to completion',
          'GET /api/health returns 200 with the version envelope',
          'notes.md describes apps/api/cmd/server/, internal/, and the root .env in one sentence each',
        ],
      },
    },
    {
      slug: 'models',
      number: 2,
      title: 'Modeling with GORM',
      tagline: 'Define your data layer the Grit way.',
      learningGoals: ['Design relational models', 'Use AutoMigrate', 'Apply constraints and indexes'],
      status: 'available',
      modules: [
        {
          title: 'Core',
          lessons: [
            { slug: 'gorm-basics', title: 'GORM basics', tagline: 'Models, tags, the convention surface.', minutes: 7, difficulty: 'easy', status: 'available' },
            { slug: 'relations', title: 'Relations', tagline: 'HasMany, BelongsTo, ManyToMany.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'migrations', title: 'Migrations', tagline: 'Auto vs. explicit, when to switch.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Model a small invoice domain',
        brief:
          'Add three models: Customer (has many Invoices), Invoice (belongs to Customer, has many LineItems), LineItem (belongs to Invoice). Run `grit migrate` and confirm the foreign keys exist via GORM Studio.',
        successCriteria: [
          'Three models exist with the relations wired',
          'AutoMigrate created the tables + FK columns',
          'GORM Studio shows the relationships',
        ],
      },
    },
    {
      slug: 'auth',
      number: 3,
      title: 'Auth & RBAC',
      tagline: 'JWT, OAuth2, 2FA, and role-based access.',
      learningGoals: ['Add JWT auth', 'Wire OAuth2 (Google + GitHub)', 'Add TOTP 2FA', 'Apply role guards'],
      status: 'available',
      modules: [
        {
          title: 'Auth',
          lessons: [
            { slug: 'jwt', title: 'JWT auth', tagline: 'Login, refresh, middleware.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'oauth', title: 'OAuth2', tagline: 'Google + GitHub sign-in.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'totp', title: 'TOTP 2FA', tagline: 'Authenticator-app codes.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
        {
          title: 'Access control',
          lessons: [
            { slug: 'rbac', title: 'RBAC + invitations', tagline: 'Roles, ownership, team invites.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Lock down a route',
        brief:
          'Add a `/api/admin/stats` endpoint that returns user counts. Protect it with the Auth middleware AND `RequireRoles("admin")`. Verify a regular user gets 404 and an admin gets 200.',
        successCriteria: [
          'Endpoint exists and is reachable from the admin role',
          'Non-admin users receive 404 (not 403 — see lesson 3.4)',
          'You can paste both responses (curl) into notes.md',
        ],
      },
    },
    {
      slug: 'batteries',
      number: 4,
      title: 'Batteries: Jobs, Mail, Storage, AI',
      tagline: 'The included primitives that ship with every Grit API.',
      learningGoals: ['Run background jobs', 'Send transactional email', 'Upload files', 'Call the AI Gateway'],
      status: 'available',
      modules: [
        {
          title: 'Async + IO',
          lessons: [
            { slug: 'jobs', title: 'Background jobs (Asynq)', tagline: 'Enqueue, retry, schedule.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'mail', title: 'Transactional email', tagline: 'Resend + HTML templates.', minutes: 6, difficulty: 'easy', status: 'available' },
            { slug: 'storage', title: 'File storage', tagline: 'S3-compatible: R2, MinIO, B2.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'ai', title: 'AI Gateway', tagline: 'Streaming + 100+ models.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Welcome flow end-to-end',
        brief:
          'When a user signs up: enqueue a welcome-email job (Asynq) that sends an HTML email (Resend → Mailhog in dev), uploads a default avatar to MinIO, and writes one line of a personalised greeting via the AI Gateway. Wire all three.',
        successCriteria: [
          'Mailhog captures a styled welcome email',
          'MinIO console shows the avatar upload',
          'The greeting line shows in the email body',
        ],
      },
    },
    {
      slug: 'security-observability',
      number: 5,
      title: 'Security + Observability',
      tagline: 'Sentinel WAF + Pulse + audit log — the production safety net.',
      learningGoals: ['Configure Sentinel', 'Read Pulse', 'Wire the audit log'],
      status: 'available',
      modules: [
        {
          title: 'Production safety',
          lessons: [
            { slug: 'sentinel', title: 'Sentinel WAF', tagline: 'Rate limit, AuthShield, Anomaly.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'pulse', title: 'Pulse observability', tagline: 'p50/p95/p99 + per-route metrics.', minutes: 7, difficulty: 'easy', status: 'available' },
            { slug: 'audit-log', title: 'Audit log', tagline: 'Tamper-evident SHA-256 chain.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Trip Sentinel + read Pulse',
        brief:
          'Hit `/api/auth/login` with bad credentials 6 times in a row to trigger Sentinel AuthShield. Then open Pulse and screenshot the p95 latency for `/api/health` over the last hour.',
        successCriteria: [
          '6th login attempt is blocked (429 or similar)',
          'Sentinel dashboard shows the AuthShield event',
          'Pulse screenshot includes the p95 line',
        ],
      },
    },
    {
      slug: 'deploy',
      number: 6,
      title: 'Deploy',
      tagline: 'Ship to a VPS with one command.',
      learningGoals: ['Configure production env', 'Run `grit deploy`', 'Set up HTTPS + a domain'],
      status: 'available',
      modules: [
        {
          title: 'Shipping',
          lessons: [
            { slug: 'grit-deploy', title: 'grit deploy', tagline: 'The deploy command, end to end.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'env-config', title: 'Production env vars', tagline: 'Secrets, JWT, database — the must-set list.', minutes: 5, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Ship to a test domain',
        brief:
          'Spin up a $5 VPS, point a subdomain at it, and run `grit deploy --domain api.<your-domain>`. Verify HTTPS works and `/api/health` is reachable from your phone.',
        successCriteria: [
          'Public URL serves /api/health over HTTPS',
          'The Let\'s Encrypt cert is valid (green padlock)',
          'You captured the response from a non-localhost client',
        ],
      },
    },
  ],
}

const mobileCourse: Course = {
  slug: 'mobile',
  name: 'Building Mobile with Go API',
  shortName: 'Mobile',
  tagline: 'Expo + React Native + your own Grit API — shipping to iOS and Android.',
  description:
    'Build a production-grade mobile app powered by a Grit API. Expo Router, type-safe API client from grit sync, secure JWT storage, push notifications, and EAS Build for the App Store + Play Store.',
  level: 'intermediate',
  estimatedHours: 10,
  prerequisites: ['Completed Grit Concepts + Building a Go API', 'Basic React knowledge'],
  whatYoullBuild: 'A two-sided mobile app with auth, real-time updates, push notifications, and offline-first reads.',
  whatYoullLearn: [
    'Scaffolding mobile with `grit new --mobile`',
    'Sharing types between API + mobile via grit sync',
    'Expo Router file-based routing',
    'Secure JWT storage (SecureStore / Keychain)',
    'Push notifications with Expo + APNs/FCM',
    'EAS Build and store submission',
  ],
  whoThisIsFor: ['Web devs adding mobile', 'Mobile devs new to typed shared schemas'],
  status: 'available',
  emoji: '📱',
  accent: 'from-rose-500/30 via-pink-500/20 to-rose-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --mobile and the Expo project layout.',
      learningGoals: ['Scaffold a mobile project', 'Run on simulator + physical device', 'Understand the monorepo shape'],
      status: 'available',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding mobile', tagline: 'The command + flags.', minutes: 5, difficulty: 'easy', status: 'available' },
            { slug: 'expo-tour', title: 'Expo project tour', tagline: 'app/, components/, hooks/.', minutes: 7, difficulty: 'easy', status: 'available' },
            { slug: 'first-run', title: 'First run', tagline: 'Simulator vs. Expo Go vs. dev build.', minutes: 6, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Run the scaffold on a real device',
        brief:
          'Scaffold `my-mobile --mobile`, get it running on either an iOS simulator (macOS) or an Android emulator, AND on a physical device via Expo Go. Screenshot both.',
        successCriteria: [
          'Scaffold completed and `pnpm install` ran clean',
          'Simulator/emulator shows the app',
          'Physical device shows the app via Expo Go',
        ],
      },
    },
    {
      slug: 'shared-types',
      number: 2,
      title: 'Shared Types + API Client',
      tagline: 'grit sync ties API + mobile together type-safely.',
      learningGoals: ['Run grit sync', 'Generate a typed API client', 'Use React Query for data'],
      status: 'available',
      modules: [
        {
          title: 'Integration',
          lessons: [
            { slug: 'grit-sync-mobile', title: 'grit sync for mobile', tagline: 'Types flow from Go structs.', minutes: 6, difficulty: 'medium', status: 'available' },
            { slug: 'api-client', title: 'Typed API client', tagline: 'fetch wrapper + React Query.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Wire a typed fetch from mobile',
        brief:
          'Add a `useUsers()` React Query hook in your mobile app that calls your Grit API. Render the list on a screen. Verify TypeScript autocompletes the User fields.',
        successCriteria: [
          'grit sync runs and generates packages/shared types',
          'Mobile app imports the User type without errors',
          'A screen renders the list from the API',
        ],
      },
    },
    {
      slug: 'auth',
      number: 3,
      title: 'Mobile Auth',
      tagline: 'Login screens + secure token storage.',
      learningGoals: ['Build login UI', 'Store JWT securely', 'Handle refresh + logout'],
      status: 'available',
      modules: [
        {
          title: 'Auth flow',
          lessons: [
            { slug: 'login-ui', title: 'Login + register screens', tagline: 'Forms + validation + errors.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'secure-storage', title: 'SecureStore + Keychain', tagline: 'Where tokens go on iOS + Android.', minutes: 6, difficulty: 'medium', status: 'available' },
            { slug: 'refresh', title: 'Token refresh', tagline: 'Silent refresh on 401.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'End-to-end auth on a device',
        brief:
          'Sign up + log in from your physical phone. Confirm SecureStore holds the refresh token. Kill the app, reopen, you should still be signed in.',
        successCriteria: [
          'Login works on the device, not just the simulator',
          'Refresh token is in SecureStore / Keychain',
          'App reopens signed in after a force-close',
        ],
      },
    },
    {
      slug: 'push-notifications',
      number: 4,
      title: 'Push Notifications',
      tagline: 'Expo Push + your Grit API talking to APNs/FCM.',
      learningGoals: ['Register for push', 'Send from the API', 'Handle taps + deep links'],
      status: 'available',
      modules: [
        {
          title: 'Push',
          lessons: [
            { slug: 'register', title: 'Registering for push', tagline: 'Permissions + token capture.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'send-from-api', title: 'Sending from the API', tagline: 'Grit job worker → Expo Push.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Send yourself a push',
        brief:
          'Wire the mobile app to register an Expo push token. Save it to your Grit API. Add an API endpoint that enqueues an Asynq job which sends a push via Expo. Trigger it from curl.',
        successCriteria: [
          'Push token is persisted on the user row',
          'Curl call enqueues a job',
          'Push arrives on your phone',
        ],
      },
    },
    {
      slug: 'ship',
      number: 5,
      title: 'Ship It',
      tagline: 'EAS Build, app store submission, OTA updates.',
      learningGoals: ['Build for iOS + Android', 'Submit to stores', 'Push OTA updates'],
      status: 'available',
      modules: [
        {
          title: 'Release',
          lessons: [
            { slug: 'eas-build', title: 'EAS Build', tagline: 'Cloud builds without local Xcode pain.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'submit', title: 'Store submission', tagline: 'App Store Connect + Play Console.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'ota', title: 'OTA updates', tagline: 'Push JS updates without a new build.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'EAS build + OTA test',
        brief:
          'Run `eas build --platform android --profile preview` to produce a real APK. Install it on your phone. Then ship a small UI change via `eas update --branch preview`; reopen the app and confirm the change appears without a re-install.',
        successCriteria: [
          'APK installs and runs on a physical Android device',
          'OTA update is published',
          'Reopening the app picks up the OTA without re-installing',
        ],
      },
    },
  ],
}

const webNextCourse: Course = {
  slug: 'web-nextjs',
  name: 'Building Web with Next.js + Go API',
  shortName: 'Web (Next.js)',
  tagline: 'The Triple kit — Next.js public site, Filament-style admin, Grit API.',
  description:
    'Build the full SaaS shape: a marketing site + product surface on Next.js, an admin panel for staff, and a Go API behind both. The combination most products ship with.',
  level: 'intermediate',
  estimatedHours: 10,
  prerequisites: ['Completed Grit Concepts + Building a Go API', 'Comfortable with React + TypeScript'],
  whatYoullBuild: 'A real SaaS: marketing landing, signup, dashboard, admin for staff, billing-ready.',
  whatYoullLearn: [
    'Next.js App Router patterns Grit uses',
    'Server actions + React Query split',
    'Filament-style admin with defineResource()',
    'DataTable + FormBuilder',
    'Tenants, roles, invitations',
    'Building marketing pages that ship fast',
  ],
  whoThisIsFor: ['Devs who completed Grit Concepts', 'React devs new to the App Router or full-stack'],
  status: 'available',
  emoji: '🌐',
  accent: 'from-emerald-500/30 via-green-500/20 to-emerald-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --triple — three apps at once.',
      learningGoals: ['Scaffold the triple kit', 'Understand the monorepo wiring', 'Run all three apps together'],
      status: 'available',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the triple', tagline: 'The command + what it produces.', minutes: 6, difficulty: 'easy', status: 'available' },
            { slug: 'tour', title: 'apps/web + apps/admin + apps/api', tagline: 'Three apps, one monorepo.', minutes: 8, difficulty: 'easy', status: 'available' },
            { slug: 'shared-package', title: 'packages/shared', tagline: 'Where types + schemas live.', minutes: 5, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'All three apps running locally',
        brief:
          'Scaffold `saas --triple`, get web on :3000, admin on :3001, API on :8080. Log in as the seeded admin. Screenshot all three.',
        successCriteria: [
          'web, admin, api respond on their ports',
          'You can log into admin with the seeded credentials',
          'Three screenshots in notes.md',
        ],
      },
    },
    {
      slug: 'public-site',
      number: 2,
      title: 'The Public Site',
      tagline: 'Marketing pages that load fast.',
      learningGoals: ['Build a landing page', 'Add a blog', 'SEO + Open Graph'],
      status: 'available',
      modules: [
        {
          title: 'Marketing',
          lessons: [
            { slug: 'landing', title: 'Landing page', tagline: 'Hero, features, CTA — the standard parts.', minutes: 9, difficulty: 'easy', status: 'available' },
            { slug: 'seo', title: 'SEO + OG', tagline: 'Metadata + sitemap + robots.', minutes: 6, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Ship a marketing site',
        brief:
          'Build a landing page for your product with a hero, three features, and a CTA. Add OG image + metadata + sitemap. Verify with Lighthouse — performance + SEO both 90+.',
        successCriteria: [
          'Landing page renders with hero/features/CTA',
          'Lighthouse SEO score >= 90',
          'OG image renders in the Facebook debugger',
        ],
      },
    },
    {
      slug: 'dashboard',
      number: 3,
      title: 'The User Dashboard',
      tagline: 'Signup → dashboard → settings.',
      learningGoals: ['Wire signup + login', 'Build dashboard widgets', 'Settings page with profile'],
      status: 'available',
      modules: [
        {
          title: 'Auth UI',
          lessons: [
            { slug: 'signup', title: 'Signup + login forms', tagline: 'Server actions + Zod.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'dashboard-widgets', title: 'Dashboard widgets', tagline: 'Stats, charts, activity feed.', minutes: 9, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Working dashboard',
        brief:
          'Sign up a new user. Land them on a dashboard with three stat cards (users, posts, revenue) and an activity feed of recent actions.',
        successCriteria: [
          'Signup creates an account on the API and stores tokens in cookies',
          'Dashboard renders three stat cards',
          'Activity feed shows real DB rows',
        ],
      },
    },
    {
      slug: 'admin-panel',
      number: 4,
      title: 'The Admin Panel',
      tagline: 'Filament-style resource panel for staff.',
      learningGoals: ['Use defineResource()', 'Customize DataTable', 'Build forms with FormBuilder'],
      status: 'available',
      modules: [
        {
          title: 'Admin',
          lessons: [
            { slug: 'define-resource', title: 'defineResource()', tagline: 'The runtime resource definition.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'datatable', title: 'DataTable', tagline: 'Sort, filter, select, paginate.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'formbuilder', title: 'FormBuilder', tagline: '8 field types, validation, dependent fields.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Full CRUD for Product',
        brief:
          'Generate a Product resource. Use defineResource() to build the admin page. List view shows price formatted as money, name searchable, status filterable. Form has name, price, stock, description, isActive.',
        successCriteria: [
          'Product admin lists, filters, sorts, paginates correctly',
          'Form creates + edits products without crashes',
          'Price column shows formatted money',
        ],
      },
    },
    {
      slug: 'tenants',
      number: 5,
      title: 'Tenants + Roles',
      tagline: 'Multi-tenant SaaS the right way.',
      learningGoals: ['Add tenants to your models', 'Wire role-based UI', 'Invitation flow end-to-end'],
      status: 'available',
      modules: [
        {
          title: 'Multi-tenancy',
          lessons: [
            { slug: 'tenant-models', title: 'Tenant models', tagline: 'Tenant-scoped queries on every read.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'role-gates', title: 'Role gates in the UI', tagline: 'Hide actions the user can\'t do.', minutes: 6, difficulty: 'medium', status: 'available' },
            { slug: 'invitations', title: 'Invitation flow', tagline: 'Email invite + accept + role assign.', minutes: 9, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Invite a team member end-to-end',
        brief:
          'As an admin, invite a teammate by email. They click the link, set a password, land in your tenant with `staff` role. They can see customers but not invoices (a role gate).',
        successCriteria: [
          'Invitation email lands in Mailhog',
          'Accepting the invite creates the user with the right role + tenant',
          'A role-gated page hides itself for the staff user',
        ],
      },
    },
  ],
}

const webTanstackCourse: Course = {
  slug: 'web-tanstack',
  name: 'Building Web with TanStack + Go API',
  shortName: 'Web (TanStack)',
  tagline: 'The Vite + TanStack Router path — pure SPA, sub-second cold starts.',
  description:
    'When SEO doesn\'t matter and dashboard performance does. Vite + TanStack Router + Grit API, with file-based routing, type-safe links, and search params as data.',
  level: 'intermediate',
  estimatedHours: 8,
  prerequisites: ['Completed Grit Concepts + Building a Go API'],
  whatYoullBuild: 'A dashboard app with sub-second cold starts, type-safe links, and a search-driven UI.',
  whatYoullLearn: [
    'Why Vite + TanStack vs. Next.js',
    'File-based routing the TanStack way',
    'Type-safe Link, search params as data',
    'Data loading patterns',
    'Auth refresh-on-401',
    'Embedded SPA deploy via go:embed',
  ],
  whoThisIsFor: ['Devs building internal dashboards', 'Teams who want one binary deploy'],
  status: 'coming-soon',
  emoji: '⚡',
  accent: 'from-cyan-500/30 via-sky-500/20 to-cyan-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --single --vite — Go binary + Vite SPA.',
      learningGoals: ['Scaffold the Vite single kit', 'Understand the embed pattern', 'Run dev mode'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding', tagline: 'The Vite single command.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'tour', title: 'Project tour', tagline: 'frontend/ + the Go binary.', minutes: 7, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'routing',
      number: 2,
      title: 'TanStack Router',
      tagline: 'File-based routes, type-safe links, the loader pattern.',
      learningGoals: ['Define routes', 'Use type-safe Link', 'Treat search params as data'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Routing',
          lessons: [
            { slug: 'file-routes', title: 'File-based routes', tagline: 'src/routes/ → URLs.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'type-safe-links', title: 'Type-safe Link', tagline: 'TS yells at you for bad routes.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'search-params', title: 'Search params as data', tagline: 'Filters, sorts, modal state in the URL.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'data',
      number: 3,
      title: 'Data Loading',
      tagline: 'Loaders + React Query for the API.',
      learningGoals: ['Use route loaders', 'Cache with React Query', 'Handle loading + error UI'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Data',
          lessons: [
            { slug: 'loaders', title: 'Route loaders', tagline: 'Fetch on route enter, not on render.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'react-query', title: 'React Query', tagline: 'Cache + invalidation.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'auth',
      number: 4,
      title: 'Auth',
      tagline: 'JWT + the lib/auth refresh-on-401 pattern.',
      learningGoals: ['Build a login screen', 'Auto-refresh expired tokens', 'Protect routes'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Auth',
          lessons: [
            { slug: 'login', title: 'Login screen', tagline: 'Form + redirect.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'refresh-pattern', title: 'Refresh on 401', tagline: 'lib/auth handles it transparently.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'route-protection', title: 'Route protection', tagline: 'beforeLoad guards.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'deploy',
      number: 5,
      title: 'Deploy',
      tagline: 'Embed + ship a single binary.',
      learningGoals: ['Build SPA into Go binary', 'Deploy to a VPS', 'HTTPS + a domain'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Ship',
          lessons: [
            { slug: 'embed-build', title: 'Embed build', tagline: 'go:embed + the SPA fallback pattern.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'deploy', title: 'Deploy', tagline: 'grit deploy + Traefik.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
  ],
}

const desktopCourse: Course = {
  slug: 'desktop',
  name: 'Building Desktop with Go API',
  shortName: 'Desktop',
  tagline: 'Wails v2 + offline-first SQLite + auto-update + branded installers.',
  description:
    'Build a real desktop app: Go backend bound to a React frontend, offline-first SQLite, in-app auto-update via GitHub releases, NSIS installers (full + slim), and a Windows-first ship pipeline.',
  level: 'intermediate',
  estimatedHours: 9,
  prerequisites: ['Completed Grit Concepts', 'Familiar with React + TypeScript'],
  whatYoullBuild: 'A Wails desktop app that runs offline, auto-updates from GitHub, and installs via a branded NSIS installer.',
  whatYoullLearn: [
    'Wails v2 binding Go to React',
    'Local SQLite + GORM',
    'Frameless windows + draggable panels',
    'In-app auto-updater (binary swap)',
    'NSIS installers (full + slim)',
    'WebView2 strategy (offline vs online)',
    'GitHub Actions release pipeline',
  ],
  whoThisIsFor: ['Devs building POS, kiosks, field apps', 'Teams that need offline-first'],
  status: 'available',
  emoji: '🖥️',
  accent: 'from-indigo-500/30 via-violet-500/20 to-indigo-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + First Run',
      tagline: 'grit new-desktop and the Wails dev loop.',
      learningGoals: ['Scaffold a desktop project', 'Run wails dev', 'Build a release binary'],
      status: 'available',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding', tagline: 'grit new-desktop walkthrough.', minutes: 6, difficulty: 'easy', status: 'available' },
            { slug: 'wails-dev', title: 'wails dev loop', tagline: 'Hot reload + dev tools.', minutes: 6, difficulty: 'easy', status: 'available' },
            { slug: 'first-build', title: 'First production build', tagline: 'wails build -nsis.', minutes: 6, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Build + install + uninstall',
        brief:
          'Scaffold `field-pos --desktop`, run `wails dev`, then `wails build -nsis`. Install the resulting Setup.exe, launch the app from the Start menu, uninstall via Settings.',
        successCriteria: [
          'wails dev opens a window in <5s',
          'Setup .exe installed without UAC prompt (per-user install)',
          'App launches from Start menu and uninstalls cleanly',
        ],
      },
    },
    {
      slug: 'offline',
      number: 2,
      title: 'Offline-First',
      tagline: 'SQLite + outbox + sync to the API.',
      learningGoals: ['Use local SQLite', 'Queue local changes', 'Sync when online'],
      status: 'available',
      modules: [
        {
          title: 'Offline',
          lessons: [
            { slug: 'sqlite', title: 'Local SQLite', tagline: 'GORM against a local file.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'outbox', title: 'Outbox pattern', tagline: 'Queue mutations, sync later.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'sync', title: 'Sync engine', tagline: 'Push + pull + conflicts.', minutes: 8, difficulty: 'hard', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Make a sale offline + sync',
        brief:
          'Disconnect from the internet. Record three sales in your app. Reconnect. Confirm all three sync to the server within 30 seconds and appear on the admin panel.',
        successCriteria: [
          'Sales recorded offline persisted to local SQLite',
          'Outbox rows visible while offline',
          'After reconnect, sales appear on the server side',
        ],
      },
    },
    {
      slug: 'frameless',
      number: 3,
      title: 'Frameless Window UI',
      tagline: 'Custom titlebar, drag regions, polished feel.',
      learningGoals: ['Build a custom titlebar', 'Wire drag regions', 'Window controls (min/max/close)'],
      status: 'available',
      modules: [
        {
          title: 'Shell',
          lessons: [
            { slug: 'titlebar', title: 'Custom titlebar', tagline: '--wails-draggable + Tailwind.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'window-controls', title: 'Window controls', tagline: 'Min, max, close, traffic lights.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Polished titlebar with controls',
        brief:
          'Build a custom titlebar with your app icon on the left, title centred, and minimise/maximise/close buttons on the right. Make the empty area draggable.',
        successCriteria: [
          'Dragging anywhere on the empty titlebar moves the window',
          'Min / Max / Close all work',
          'Buttons are NOT in the drag region (don\'t move the window when clicked)',
        ],
      },
    },
    {
      slug: 'auto-update',
      number: 4,
      title: 'In-App Auto-Update',
      tagline: 'Binary swap + GitHub releases + modal UI.',
      learningGoals: ['Wire the updater', 'Build the modal UI', 'Cut a release that auto-distributes'],
      status: 'available',
      modules: [
        {
          title: 'Updater',
          lessons: [
            { slug: 'updater-go', title: 'updater.go', tagline: 'The Wails binding + binary swap.', minutes: 9, difficulty: 'hard', status: 'available' },
            { slug: 'modal-ui', title: 'Modal + banner UI', tagline: 'React components shipping in scaffold.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'release-script', title: 'release-desktop.sh', tagline: 'One-shot pipeline.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Cut v0.1.1 and watch v0.1.0 auto-update',
        brief:
          'Install v0.1.0 of your app. Change a small thing, bump version to 0.1.1, run release-desktop.sh. The running v0.1.0 should detect the new release, prompt the user, and swap to v0.1.1.',
        successCriteria: [
          'GitHub release v0.1.1 published with installer + raw .exe',
          'Running v0.1.0 shows update banner',
          'Clicking Install swaps the binary; relaunch shows v0.1.1',
        ],
      },
    },
    {
      slug: 'installers',
      number: 5,
      title: 'Branded Installers',
      tagline: 'NSIS full + slim, branded MUI bitmaps.',
      learningGoals: ['Customize NSIS', 'Choose between full + slim', 'Generate branded bitmaps'],
      status: 'available',
      modules: [
        {
          title: 'NSIS',
          lessons: [
            { slug: 'project-nsi', title: 'project.nsi (full)', tagline: 'Bundled WebView2 — no internet needed.', minutes: 8, difficulty: 'hard', status: 'available' },
            { slug: 'project-slim', title: 'project-slim.nsi', tagline: 'Online bootstrapper — small download.', minutes: 6, difficulty: 'medium', status: 'available' },
            { slug: 'bitmaps', title: 'Branded bitmaps', tagline: 'header.bmp + welcome.bmp from icon.ico.', minutes: 5, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Two installers, both branded',
        brief:
          'Generate the full installer (~150MB, bundled WebView2) AND the slim installer (~22MB, online bootstrap). Both should have a branded MUI Welcome page with your icon. Install both, uninstall both.',
        successCriteria: [
          'Full installer works on a machine with no internet',
          'Slim installer downloads + installs WebView2 if missing',
          'Both show the branded Welcome bitmap',
        ],
      },
    },
  ],
}

const multiplatformCourse: Course = {
  slug: 'multiplatform',
  name: 'Building Web + Desktop + Mobile',
  shortName: 'Multi-Platform',
  tagline: 'One Go API powering all three surfaces — the kitchen-sink kit.',
  description:
    'Build a product that ships on the web, desktop, and mobile from a single Go API. The triple kit + desktop + mobile, shared types, sync engines, and one release pipeline coordinating it all.',
  level: 'advanced',
  estimatedHours: 14,
  prerequisites: [
    'Completed Building a Go API',
    'Completed at least one of Mobile, Desktop, Web (Next.js or TanStack)',
  ],
  whatYoullBuild: 'A multi-platform SaaS: web app, native desktop, native mobile, all running off one API.',
  whatYoullLearn: [
    'The full triple + desktop + mobile setup',
    'Sharing types across 4 surfaces',
    'Offline sync per platform',
    'Push notifications across web + mobile',
    'Coordinating releases',
    'When to use which surface for which feature',
  ],
  whoThisIsFor: ['Senior devs building multi-surface products', 'Teams scaling from web to mobile + desktop'],
  status: 'available',
  emoji: '🌍',
  accent: 'from-amber-500/30 via-orange-500/20 to-amber-500/10',
  chapters: [
    {
      slug: 'foundation',
      number: 1,
      title: 'The Foundation',
      tagline: 'grit new --triple + adding mobile + desktop on top.',
      learningGoals: ['Scaffold the full stack', 'Understand the monorepo wiring', 'Run all four surfaces'],
      status: 'available',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the full stack', tagline: 'Triple + mobile + desktop.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'monorepo-wiring', title: 'Monorepo wiring', tagline: 'pnpm workspaces, turbo, go module.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Scaffold a four-surface project',
        brief:
          "Run `grit new` with --triple, then layer the mobile and desktop kits on top. Boot all four (api, web, admin, mobile, desktop) and capture a screenshot of each surface showing the seeded data.",
        successCriteria: [
          'API runs on :8080',
          'Web (port 3000) + Admin (port 3001) both load',
          'Expo mobile shows the dev menu',
          'Wails desktop launches a window',
          'All four read from the same DB',
        ],
      },
    },
    {
      slug: 'shared-types',
      number: 2,
      title: 'Shared Types Everywhere',
      tagline: 'One source of truth for 4 surfaces.',
      learningGoals: ['grit sync across surfaces', 'Type-safe API client', 'Shared Zod schemas'],
      status: 'available',
      modules: [
        {
          title: 'Types',
          lessons: [
            { slug: 'grit-sync-multi', title: 'grit sync — multi-surface', tagline: 'One Go model, three TS frontends.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'shared-zod', title: 'Shared Zod schemas', tagline: 'Validate input on every surface.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Add a field and feel the wave',
        brief:
          "Add an `IsFeatured bool` field to the Product model in the Go API. Run `grit sync`. Confirm the TypeScript type and Zod schema both regenerate, and that web / admin / mobile / desktop all type-check against the new field without any per-surface edits.",
        successCriteria: [
          'IsFeatured appears in shared/types/product.ts',
          'IsFeatured appears in shared/zod/product.ts',
          'TS build passes on web + admin + mobile + desktop',
          'No per-surface edits needed beyond using the new field',
        ],
      },
    },
    {
      slug: 'feature-implementation',
      number: 3,
      title: 'Building a Feature Across All Three',
      tagline: 'Pick a real feature — implement it everywhere.',
      learningGoals: ['Build a feature in web', 'Implement on mobile', 'Add to desktop'],
      status: 'available',
      modules: [
        {
          title: 'Cross-surface',
          lessons: [
            { slug: 'pick-feature', title: 'Picking the feature', tagline: 'What ships well across surfaces.', minutes: 5, difficulty: 'easy', status: 'available' },
            { slug: 'web-impl', title: 'Web implementation', tagline: 'Next.js page + admin.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'mobile-impl', title: 'Mobile implementation', tagline: 'Expo screens.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'desktop-impl', title: 'Desktop implementation', tagline: 'Wails screens.', minutes: 9, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Ship Bookmarks on three surfaces',
        brief:
          "Add a Bookmarks feature: pick from your API's catalog, save to a personal list, view + delete from any surface. Implement on web, mobile, and desktop. Each surface should share the same Bookmark type + Zod schema from shared/.",
        successCriteria: [
          'Add/Remove bookmark works on all three surfaces',
          'The list reflects updates within seconds (no full refresh)',
          'You hit the same /api/bookmarks endpoint from all three',
          'No duplicate type / Zod definitions',
        ],
      },
    },
    {
      slug: 'sync',
      number: 4,
      title: 'Sync + Offline',
      tagline: 'Mobile + desktop sync to the same API.',
      learningGoals: ['Mobile offline reads', 'Desktop outbox', 'Resolve conflicts'],
      status: 'available',
      modules: [
        {
          title: 'Offline',
          lessons: [
            { slug: 'mobile-offline', title: 'Mobile offline', tagline: 'React Query persistence.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'desktop-outbox', title: 'Desktop outbox', tagline: 'Queue + push + pull.', minutes: 9, difficulty: 'hard', status: 'available' },
            { slug: 'conflicts', title: 'Conflict resolution', tagline: 'Last-write-wins vs. CRDT-light.', minutes: 8, difficulty: 'hard', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Survive a flaky network',
        brief:
          "Toggle airplane mode on mobile + desktop while the user is mid-flow. Confirm both keep working (offline reads + queued writes). Re-enable network. Confirm both flush their queues and converge to the same state without manual refresh.",
        successCriteria: [
          'Offline reads work on mobile (React Query persistence)',
          'Offline writes queue on desktop (outbox)',
          'Both flush on reconnect and reach the same server state',
          'A conflict between desktop + mobile resolves deterministically',
        ],
      },
    },
    {
      slug: 'releases',
      number: 5,
      title: 'Coordinated Releases',
      tagline: 'Ship to all surfaces together.',
      learningGoals: ['Version compatibility matrix', 'Stagger releases safely', 'Roll back across surfaces'],
      status: 'available',
      modules: [
        {
          title: 'Releases',
          lessons: [
            { slug: 'compat-matrix', title: 'Compatibility matrix', tagline: 'Which API versions support which clients.', minutes: 7, difficulty: 'hard', status: 'available' },
            { slug: 'staggered', title: 'Staggered releases', tagline: 'Web first, then desktop, then mobile.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Cut a coordinated release',
        brief:
          "Ship v1.1.0 across all four surfaces using a staggered plan: API + web first, then desktop, then mobile. Document the compatibility matrix in `notes.md` (API v1.0 supports clients ≥ 1.0; API v1.1 supports clients ≥ 1.0). Roll back ONE surface as a fire drill and confirm the others still work.",
        successCriteria: [
          'A compatibility matrix written down',
          'Web + API deploy first, desktop second, mobile third',
          'A simulated rollback of one surface does not break the others',
          '`notes.md` lists the release order + rollback playbook',
        ],
      },
    },
  ],
}

/* ─────────────────────────────────────────────────────────────────── */
/*  Performance & Quality courses                                       */
/* ─────────────────────────────────────────────────────────────────── */

const loadTestingCourse: Course = {
  slug: 'load-testing',
  name: 'Load Testing with K6',
  shortName: 'K6 Load Testing',
  tagline: 'Prove your API holds up before users find out it doesn\'t.',
  description:
    "Hands-on load testing for a Grit API: write k6 scripts, run smoke / load / stress / spike / soak tests, read p95 latency, find the bottleneck, fix it, prove the fix. The mindset and the tool, with real numbers.",
  level: 'intermediate',
  estimatedHours: 4,
  prerequisites: ['Completed Building a Go API', 'Comfortable in a terminal'],
  whatYoullBuild:
    "A k6 test suite for your Grit API with 6 scenario types, run locally and in CI, that catches regressions before they ship.",
  whatYoullLearn: [
    'When to load-test (and when not)',
    'k6 fundamentals — VUs, iterations, scenarios, thresholds',
    'The 5 test types: smoke, load, stress, spike, soak',
    'Reading p50/p95/p99 and what each tells you',
    'Finding the bottleneck: API, DB, network, or your test rig',
    'Wiring k6 into GitHub Actions so PRs that regress get blocked',
  ],
  whoThisIsFor: [
    'Devs shipping APIs to production',
    'SREs setting up performance gates',
    'Teams who got burned by a surprise outage at launch',
  ],
  status: 'available',
  emoji: '⚡',
  accent: 'from-yellow-500/30 via-amber-500/20 to-yellow-500/10',
  chapters: [
    {
      slug: 'fundamentals',
      number: 1,
      title: 'k6 Fundamentals',
      tagline: 'Install, write your first test, understand the output.',
      learningGoals: ['Install k6', 'Write a 20-line test', 'Read the summary at the end'],
      status: 'available',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'why-load-test', title: 'Why load test?', tagline: 'When it pays for itself — and when it\'s premature.', minutes: 5, difficulty: 'easy', status: 'available' },
            { slug: 'install-k6', title: 'Install k6', tagline: 'One binary, no Node, no deps.', minutes: 3, difficulty: 'easy', status: 'available' },
            { slug: 'first-script', title: 'Your first k6 script', tagline: 'A 20-line script that hits /healthz with 10 VUs.', minutes: 8, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Run your first test',
        brief: 'Install k6, write a 20-line script against your Grit API\'s /healthz endpoint, run it with 10 VUs for 30 seconds, and paste the summary into notes.md. Explain in 2 sentences what p95 means in plain English.',
        successCriteria: ['k6 installed and runs', 'Summary shows http_req_duration metrics', 'You can explain p95 to a non-engineer'],
      },
    },
    {
      slug: 'scenarios',
      number: 2,
      title: 'The Five Test Types',
      tagline: 'Smoke, load, stress, spike, soak — when to run each.',
      learningGoals: ['Pick the right test for the question being asked', 'Define thresholds that fail the build on regression'],
      status: 'available',
      modules: [
        {
          title: 'Scenario types',
          lessons: [
            { slug: 'smoke-load', title: 'Smoke + load tests', tagline: 'Sanity check + steady-state expected traffic.', minutes: 7, difficulty: 'easy', status: 'available' },
            { slug: 'stress-spike', title: 'Stress + spike tests', tagline: 'Push to breaking point; sudden traffic surges.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'soak-test', title: 'Soak test', tagline: 'Run for hours to catch leaks the short tests miss.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'A 6-test suite',
        brief: 'Build a k6 suite in `tests/k6/` with one script per scenario type (smoke, load, stress, spike, soak) plus a registration-flow scenario. Document what each tests in a README.',
        successCriteria: ['All 6 scripts run', 'Each has thresholds that fail on regression', 'README explains when to run each'],
      },
    },
    {
      slug: 'finding-bottlenecks',
      number: 3,
      title: 'Finding the Bottleneck',
      tagline: 'When the test fails, where do you look?',
      learningGoals: ['Use Pulse + DB metrics + system stats to localise the problem', 'Distinguish API CPU, DB CPU, network, and test-rig limits'],
      status: 'available',
      modules: [
        {
          title: 'Diagnosis',
          lessons: [
            { slug: 'reading-metrics', title: 'Reading the numbers', tagline: 'p50 vs p95 vs p99 — what each tells you.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'common-bottlenecks', title: 'Common bottlenecks', tagline: 'DB connection pool, N+1, cold caches, GC pauses.', minutes: 9, difficulty: 'medium', status: 'available' },
            { slug: 'fix-and-prove', title: 'Fix and prove it', tagline: 'Make a change, re-run, compare — disciplined performance work.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Find + fix + prove',
        brief: 'Pick the endpoint with the worst p95 from your suite. Use Pulse + DB logs to identify the bottleneck. Fix it (add an index, add caching, fix N+1). Re-run the same scenario. Paste before/after numbers in notes.md with a 1-paragraph writeup.',
        successCriteria: ['p95 improved by ≥30%', 'You can explain the root cause in one sentence', 'The fix is committed'],
      },
    },
    {
      slug: 'ci-integration',
      number: 4,
      title: 'k6 in CI',
      tagline: 'Block regressions in PRs.',
      learningGoals: ['Run k6 in GitHub Actions', 'Fail the build on threshold breach'],
      status: 'available',
      modules: [
        {
          title: 'Automation',
          lessons: [
            { slug: 'github-actions', title: 'Running k6 in GitHub Actions', tagline: 'One workflow, one Docker image, automatic on PR.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'PR-gating perf check',
        brief: 'Wire k6 smoke + load tests into a GitHub Action that runs on every PR. Open a PR with an intentional regression (e.g., remove a DB index). Confirm the PR check fails.',
        successCriteria: ['CI runs k6 on every PR', 'A regression-PR fails the check', 'The action artifact has the summary'],
      },
    },
  ],
}

const securityCourse: Course = {
  slug: 'security',
  name: 'Security & Pen Testing for Grit APIs',
  shortName: 'Security',
  tagline: 'OWASP Top 10, hands-on — attack your own API, then defend it.',
  description:
    "Walk the OWASP Top 10 against your own Grit API. Each lesson: how the attack works, exploit it on a deliberately-vulnerable endpoint, then ship the fix. Plus the defensive stack Grit gives you out of the box (Sentinel rate-limit, safefetch, authz, CSRF, audit logs).",
  level: 'advanced',
  estimatedHours: 6,
  prerequisites: ['Completed Building a Go API', 'Comfortable with HTTP, JWT, and CORS'],
  whatYoullBuild:
    "A security-hardened Grit API with hands-on confidence in detecting + fixing IDOR, SSRF, broken auth, mass assignment, injection, and the rest of the OWASP Top 10.",
  whatYoullLearn: [
    'How each OWASP Top 10 vulnerability works in practice',
    'Exploit your own endpoint to feel the risk',
    'Ship the fix and verify the exploit no longer works',
    'The Grit defensive stack: Sentinel, safefetch, authz, CSRF, security headers',
    'Threat modelling — what to worry about, what not to',
  ],
  whoThisIsFor: [
    'Devs shipping public APIs',
    'Security-aware founders without a dedicated security team',
    'Anyone tired of mystery 3am pages',
  ],
  status: 'available',
  emoji: '🛡️',
  accent: 'from-rose-500/30 via-red-500/20 to-rose-500/10',
  chapters: [
    {
      slug: 'mindset',
      number: 1,
      title: 'The Attacker\'s Mindset',
      tagline: 'Think like the attacker before you defend like one.',
      learningGoals: ['Frame every endpoint as &quot;what could go wrong?&quot;', 'Map the OWASP Top 10 to real Grit endpoints'],
      status: 'available',
      modules: [
        {
          title: 'Foundation',
          lessons: [
            { slug: 'threat-model', title: 'A threat model in 15 minutes', tagline: 'What you protect, from whom, with what.', minutes: 6, difficulty: 'easy', status: 'available' },
            { slug: 'owasp-tour', title: 'OWASP Top 10 — the speedrun tour', tagline: 'One sentence per category, in plain English.', minutes: 8, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Threat-model your API',
        brief: 'Write a 1-page threat model for your Grit API in notes.md: 5 assets, 5 actors (user, admin, attacker, bot, partner), and 10 threats (one per OWASP category). One paragraph each.',
        successCriteria: ['Document covers all 10 OWASP categories', 'Each threat has a Grit endpoint it applies to', 'You can defend the priorities you set'],
      },
    },
    {
      slug: 'access-control',
      number: 2,
      title: 'Broken Access Control',
      tagline: 'IDOR, missing role checks, and the &quot;just check user_id&quot; rule.',
      learningGoals: ['Spot an IDOR in a code review', 'Add authorization checks the right way'],
      status: 'available',
      modules: [
        {
          title: 'Attack + defend',
          lessons: [
            { slug: 'idor', title: 'IDOR — the most common bug in the wild', tagline: 'Cross-account access by guessing IDs.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'authz-package', title: 'The authz package', tagline: 'Centralized authorization in Grit and how to use it.', minutes: 7, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Exploit + patch an IDOR',
        brief: 'In a fresh Grit project, create two users (A and B), one note owned by A. As B, attempt to PATCH and DELETE A\'s note. Document the exploit, then add the authorization check. Re-test — must return 403.',
        successCriteria: ['Exploit works initially (proves the bug)', 'After patch, returns 403', 'Test added so a regression fails CI'],
      },
    },
    {
      slug: 'injection-ssrf',
      number: 3,
      title: 'Injection & SSRF',
      tagline: 'SQL injection, command injection, and the URL trap.',
      learningGoals: ['Recognise dangerous string-concat patterns', 'Use the safefetch package for any URL the user controls'],
      status: 'available',
      modules: [
        {
          title: 'Attack + defend',
          lessons: [
            { slug: 'sql-injection', title: 'SQL injection in Go', tagline: 'How GORM protects you, and how Sprintf undoes that protection.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'ssrf-safefetch', title: 'SSRF + safefetch', tagline: 'The internal-network attack and how Grit\'s safefetch blocks it.', minutes: 9, difficulty: 'hard', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Try to break your own API',
        brief: 'Use curl to attempt SQL injection on your search endpoints, and SSRF on any URL-fetching endpoint (e.g., webhook validation, OG-image preview). Document each attempt + the response. Add safefetch where missing.',
        successCriteria: ['All injection attempts return clean errors, not data leaks', 'No internal IPs reachable via your API', 'Tests added for the most dangerous endpoints'],
      },
    },
    {
      slug: 'auth-secrets',
      number: 4,
      title: 'Auth + Secret Management',
      tagline: 'JWT pitfalls, session security, and where to put your keys.',
      learningGoals: ['Configure JWT correctly (audience, expiry, rotation)', 'Identify secrets that should never be in code'],
      status: 'available',
      modules: [
        {
          title: 'Hardening',
          lessons: [
            { slug: 'jwt-pitfalls', title: 'JWT pitfalls', tagline: 'alg=none, missing expiry, leaked secrets.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'secrets', title: 'Where do your secrets live?', tagline: '.env, vaults, KMS — what to use when.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Audit your tokens + secrets',
        brief: 'Document where every secret lives. Rotate JWT_SECRET in dev — what breaks? Add audience + expiry validation to your JWT verify. Add one Sentinel rate-limit rule to /api/auth/login.',
        successCriteria: ['Token rotation works without downtime', 'No secrets in committed .env', 'Rate limit verified by k6 test'],
      },
    },
    {
      slug: 'defensive-stack',
      number: 5,
      title: 'The Grit Defensive Stack',
      tagline: 'Sentinel, security headers, CSRF, audit log — wire them all.',
      learningGoals: ['Enable + tune each defence', 'Understand what each blocks and what it doesn\'t'],
      status: 'available',
      modules: [
        {
          title: 'Defence in depth',
          lessons: [
            { slug: 'sentinel-rate-limit', title: 'Sentinel rate limiting', tagline: 'Per-IP, per-user, per-endpoint throttles.', minutes: 7, difficulty: 'medium', status: 'available' },
            { slug: 'csrf-cors', title: 'CSRF + CORS', tagline: 'Why both, when each matters.', minutes: 8, difficulty: 'hard', status: 'available' },
            { slug: 'audit-log', title: 'Audit log', tagline: 'Who did what, when. The black-box recorder for incidents.', minutes: 6, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Full defence audit',
        brief: 'Verify every defence is enabled + tuned: Sentinel limits, CSP headers, CSRF on form endpoints, audit log on sensitive ops. Document each in a `SECURITY.md` at the repo root.',
        successCriteria: ['SECURITY.md exists and lists every active defence', 'curl against deliberately-broken inputs returns the right status', 'Audit log captures admin actions'],
      },
    },
  ],
}

const benchmarkingCourse: Course = {
  slug: 'benchmarking',
  name: 'Benchmarking Your Go Code',
  shortName: 'Benchmarking',
  tagline: 'Go\'s built-in microbenchmarks — measure first, optimise second.',
  description:
    "Go has world-class benchmarking built in. This short course teaches you to write `BenchmarkXxx` functions, read the output, profile with pprof, and refuse to optimise without data. Antidote to vibes-based performance work.",
  level: 'intermediate',
  estimatedHours: 3,
  prerequisites: ['Completed Building a Go API', 'Comfortable writing Go tests'],
  whatYoullBuild:
    "A benchmark suite for your Go service hot paths, with a baseline, a profiling workflow, and one measured optimisation.",
  whatYoullLearn: [
    'Write a Benchmark function in 10 lines',
    'Read ns/op, B/op, allocs/op',
    'Profile with pprof — CPU + heap',
    'Make ONE measured change and prove the win',
    'Compare two implementations with benchstat',
  ],
  whoThisIsFor: [
    'Devs guessing at performance',
    'Anyone tempted to micro-optimise without data',
    'Backend engineers shipping latency-sensitive features',
  ],
  status: 'available',
  emoji: '📊',
  accent: 'from-cyan-500/30 via-sky-500/20 to-cyan-500/10',
  chapters: [
    {
      slug: 'fundamentals',
      number: 1,
      title: 'Microbenchmarks 101',
      tagline: 'BenchmarkXxx, b.N, and what the output means.',
      learningGoals: ['Write your first benchmark', 'Read ns/op and allocs/op without lookup'],
      status: 'available',
      modules: [
        {
          title: 'The basics',
          lessons: [
            { slug: 'first-bench', title: 'Your first benchmark', tagline: 'A `BenchmarkPluralize` in 10 lines.', minutes: 7, difficulty: 'easy', status: 'available' },
            { slug: 'reading-output', title: 'Reading the output', tagline: 'ns/op, B/op, allocs/op — what each tells you.', minutes: 6, difficulty: 'easy', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Benchmark a hot path',
        brief: 'Pick the most-called function in your service (often a serialiser, parser, or auth check). Write a benchmark. Run it 3 times and confirm the numbers are stable. Paste the output in notes.md.',
        successCriteria: ['Benchmark function exists in `_test.go`', 'Numbers stable across 3 runs', 'You can interpret each column'],
      },
    },
    {
      slug: 'profiling',
      number: 2,
      title: 'Profiling with pprof',
      tagline: 'Where is the time really going?',
      learningGoals: ['Capture CPU + heap profiles', 'Read a flame graph'],
      status: 'available',
      modules: [
        {
          title: 'Tools',
          lessons: [
            { slug: 'cpu-profile', title: 'CPU profile', tagline: 'Capture, open in pprof, find the hot frame.', minutes: 8, difficulty: 'medium', status: 'available' },
            { slug: 'heap-profile', title: 'Heap profile', tagline: 'Find allocation sources and reduce GC pressure.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'Profile + diagnose',
        brief: 'Capture a CPU + heap profile of your benchmark. Identify the hottest function. Write 1 paragraph describing what surprises you (almost always there\'s a surprise).',
        successCriteria: ['CPU + heap profiles saved', 'You named the hottest function', 'You explain WHY in 1 paragraph'],
      },
    },
    {
      slug: 'optimise',
      number: 3,
      title: 'Measure → Change → Measure',
      tagline: 'The only legitimate way to optimise.',
      learningGoals: ['Use benchstat to compare runs', 'Make one change and prove the win or roll back'],
      status: 'available',
      modules: [
        {
          title: 'Disciplined optimisation',
          lessons: [
            { slug: 'benchstat', title: 'benchstat — compare two runs', tagline: 'Statistical significance, not vibes.', minutes: 6, difficulty: 'medium', status: 'available' },
            { slug: 'one-optimisation', title: 'Ship one measured optimisation', tagline: 'Pick, change, re-run, decide.', minutes: 8, difficulty: 'medium', status: 'available' },
          ],
        },
      ],
      assignment: {
        title: 'One real, measured win',
        brief: 'Make ONE change to your service\'s hot path. Re-run benchmarks. Use benchstat to verify the change is real (p<0.05). If not, roll back. Paste before/after in notes.md.',
        successCriteria: ['benchstat shows significant improvement', 'OR you rolled back honestly', 'Either way you commit data, not vibes'],
      },
    },
  ],
}

export const COURSES: Course[] = [
  conceptsCourse,
  goApiCourse,
  mobileCourse,
  webNextCourse,
  webTanstackCourse,
  desktopCourse,
  multiplatformCourse,
  loadTestingCourse,
  securityCourse,
  benchmarkingCourse,
]
