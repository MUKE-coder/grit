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
      slug: 'architecture-modes',
      number: 5,
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
  status: 'coming-soon',
  emoji: '🦫',
  accent: 'from-sky-500/30 via-cyan-500/20 to-sky-500/10',
  chapters: [
    {
      slug: 'scaffold-tour',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --api and a deep tour of the Go-only project shape.',
      learningGoals: ['Scaffold an API-only project', 'Identify every Go file', 'Run + curl /api/health'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Start',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the API', tagline: 'grit new my-api --api walkthrough.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'project-tour', title: 'Project tour', tagline: 'cmd/server, internal/, the Services pattern.', minutes: 8, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'first-request', title: 'Your first request', tagline: 'docker compose up, run the API, hit /api/health.', minutes: 4, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'models',
      number: 2,
      title: 'Modeling with GORM',
      tagline: 'Define your data layer the Grit way.',
      learningGoals: ['Design relational models', 'Use AutoMigrate', 'Apply constraints and indexes'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Core',
          lessons: [
            { slug: 'gorm-basics', title: 'GORM basics', tagline: 'Models, tags, the convention surface.', minutes: 7, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'relations', title: 'Relations', tagline: 'HasMany, BelongsTo, ManyToMany.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'migrations', title: 'Migrations', tagline: 'Auto vs. explicit, when to switch.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'auth',
      number: 3,
      title: 'Auth & RBAC',
      tagline: 'JWT, OAuth2, 2FA, and role-based access.',
      learningGoals: ['Add JWT auth', 'Wire OAuth2 (Google + GitHub)', 'Add TOTP 2FA', 'Apply role guards'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Auth',
          lessons: [
            { slug: 'jwt', title: 'JWT auth', tagline: 'Login, refresh, middleware.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'oauth', title: 'OAuth2', tagline: 'Google + GitHub sign-in.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'totp', title: 'TOTP 2FA', tagline: 'Authenticator-app codes.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
        {
          title: 'Access control',
          lessons: [
            { slug: 'rbac', title: 'RBAC + invitations', tagline: 'Roles, ownership, team invites.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'batteries',
      number: 4,
      title: 'Batteries: Jobs, Mail, Storage, AI',
      tagline: 'The included primitives that ship with every Grit API.',
      learningGoals: ['Run background jobs', 'Send transactional email', 'Upload files', 'Call the AI Gateway'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Async + IO',
          lessons: [
            { slug: 'jobs', title: 'Background jobs (Asynq)', tagline: 'Enqueue, retry, schedule.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'mail', title: 'Transactional email', tagline: 'Resend + HTML templates.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'storage', title: 'File storage', tagline: 'S3-compatible: R2, MinIO, B2.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'ai', title: 'AI Gateway', tagline: 'Streaming + 100+ models.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'security-observability',
      number: 5,
      title: 'Security + Observability',
      tagline: 'Sentinel WAF + Pulse + audit log — the production safety net.',
      learningGoals: ['Configure Sentinel', 'Read Pulse', 'Wire the audit log'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Production safety',
          lessons: [
            { slug: 'sentinel', title: 'Sentinel WAF', tagline: 'Rate limit, AuthShield, Anomaly.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'pulse', title: 'Pulse observability', tagline: 'p50/p95/p99 + per-route metrics.', minutes: 7, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'audit-log', title: 'Audit log', tagline: 'Tamper-evident SHA-256 chain.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'deploy',
      number: 6,
      title: 'Deploy',
      tagline: 'Ship to a VPS with one command.',
      learningGoals: ['Configure production env', 'Run `grit deploy`', 'Set up HTTPS + a domain'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Shipping',
          lessons: [
            { slug: 'grit-deploy', title: 'grit deploy', tagline: 'The deploy command, end to end.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'env-config', title: 'Production env vars', tagline: 'Secrets, JWT, database — the must-set list.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
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
  status: 'coming-soon',
  emoji: '📱',
  accent: 'from-rose-500/30 via-pink-500/20 to-rose-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --mobile and the Expo project layout.',
      learningGoals: ['Scaffold a mobile project', 'Run on simulator + physical device', 'Understand the monorepo shape'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding mobile', tagline: 'The command + flags.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'expo-tour', title: 'Expo project tour', tagline: 'app/, components/, hooks/.', minutes: 7, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'first-run', title: 'First run', tagline: 'Simulator vs. Expo Go vs. dev build.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'shared-types',
      number: 2,
      title: 'Shared Types + API Client',
      tagline: 'grit sync ties API + mobile together type-safely.',
      learningGoals: ['Run grit sync', 'Generate a typed API client', 'Use React Query for data'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Integration',
          lessons: [
            { slug: 'grit-sync-mobile', title: 'grit sync for mobile', tagline: 'Types flow from Go structs.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'api-client', title: 'Typed API client', tagline: 'fetch wrapper + React Query.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'auth',
      number: 3,
      title: 'Mobile Auth',
      tagline: 'Login screens + secure token storage.',
      learningGoals: ['Build login UI', 'Store JWT securely', 'Handle refresh + logout'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Auth flow',
          lessons: [
            { slug: 'login-ui', title: 'Login + register screens', tagline: 'Forms + validation + errors.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'secure-storage', title: 'SecureStore + Keychain', tagline: 'Where tokens go on iOS + Android.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'refresh', title: 'Token refresh', tagline: 'Silent refresh on 401.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'push-notifications',
      number: 4,
      title: 'Push Notifications',
      tagline: 'Expo Push + your Grit API talking to APNs/FCM.',
      learningGoals: ['Register for push', 'Send from the API', 'Handle taps + deep links'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Push',
          lessons: [
            { slug: 'register', title: 'Registering for push', tagline: 'Permissions + token capture.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'send-from-api', title: 'Sending from the API', tagline: 'Grit job worker → Expo Push.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'ship',
      number: 5,
      title: 'Ship It',
      tagline: 'EAS Build, app store submission, OTA updates.',
      learningGoals: ['Build for iOS + Android', 'Submit to stores', 'Push OTA updates'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Release',
          lessons: [
            { slug: 'eas-build', title: 'EAS Build', tagline: 'Cloud builds without local Xcode pain.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'submit', title: 'Store submission', tagline: 'App Store Connect + Play Console.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'ota', title: 'OTA updates', tagline: 'Push JS updates without a new build.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
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
  status: 'coming-soon',
  emoji: '🌐',
  accent: 'from-emerald-500/30 via-green-500/20 to-emerald-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + Tour',
      tagline: 'grit new --triple — three apps at once.',
      learningGoals: ['Scaffold the triple kit', 'Understand the monorepo wiring', 'Run all three apps together'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the triple', tagline: 'The command + what it produces.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'tour', title: 'apps/web + apps/admin + apps/api', tagline: 'Three apps, one monorepo.', minutes: 8, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'shared-package', title: 'packages/shared', tagline: 'Where types + schemas live.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'public-site',
      number: 2,
      title: 'The Public Site',
      tagline: 'Marketing pages that load fast.',
      learningGoals: ['Build a landing page', 'Add a blog', 'SEO + Open Graph'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Marketing',
          lessons: [
            { slug: 'landing', title: 'Landing page', tagline: 'Hero, features, CTA — the standard parts.', minutes: 9, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'seo', title: 'SEO + OG', tagline: 'Metadata + sitemap + robots.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'dashboard',
      number: 3,
      title: 'The User Dashboard',
      tagline: 'Signup → dashboard → settings.',
      learningGoals: ['Wire signup + login', 'Build dashboard widgets', 'Settings page with profile'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Auth UI',
          lessons: [
            { slug: 'signup', title: 'Signup + login forms', tagline: 'Server actions + Zod.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'dashboard-widgets', title: 'Dashboard widgets', tagline: 'Stats, charts, activity feed.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'admin-panel',
      number: 4,
      title: 'The Admin Panel',
      tagline: 'Filament-style resource panel for staff.',
      learningGoals: ['Use defineResource()', 'Customize DataTable', 'Build forms with FormBuilder'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Admin',
          lessons: [
            { slug: 'define-resource', title: 'defineResource()', tagline: 'The runtime resource definition.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'datatable', title: 'DataTable', tagline: 'Sort, filter, select, paginate.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'formbuilder', title: 'FormBuilder', tagline: '8 field types, validation, dependent fields.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'tenants',
      number: 5,
      title: 'Tenants + Roles',
      tagline: 'Multi-tenant SaaS the right way.',
      learningGoals: ['Add tenants to your models', 'Wire role-based UI', 'Invitation flow end-to-end'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Multi-tenancy',
          lessons: [
            { slug: 'tenant-models', title: 'Tenant models', tagline: 'Tenant-scoped queries on every read.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'role-gates', title: 'Role gates in the UI', tagline: 'Hide actions the user can\'t do.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'invitations', title: 'Invitation flow', tagline: 'Email invite + accept + role assign.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
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
  status: 'coming-soon',
  emoji: '🖥️',
  accent: 'from-indigo-500/30 via-violet-500/20 to-indigo-500/10',
  chapters: [
    {
      slug: 'scaffold',
      number: 1,
      title: 'Scaffold + First Run',
      tagline: 'grit new-desktop and the Wails dev loop.',
      learningGoals: ['Scaffold a desktop project', 'Run wails dev', 'Build a release binary'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding', tagline: 'grit new-desktop walkthrough.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'wails-dev', title: 'wails dev loop', tagline: 'Hot reload + dev tools.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'first-build', title: 'First production build', tagline: 'wails build -nsis.', minutes: 6, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'offline',
      number: 2,
      title: 'Offline-First',
      tagline: 'SQLite + outbox + sync to the API.',
      learningGoals: ['Use local SQLite', 'Queue local changes', 'Sync when online'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Offline',
          lessons: [
            { slug: 'sqlite', title: 'Local SQLite', tagline: 'GORM against a local file.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'outbox', title: 'Outbox pattern', tagline: 'Queue mutations, sync later.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'sync', title: 'Sync engine', tagline: 'Push + pull + conflicts.', minutes: 8, difficulty: 'hard', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'frameless',
      number: 3,
      title: 'Frameless Window UI',
      tagline: 'Custom titlebar, drag regions, polished feel.',
      learningGoals: ['Build a custom titlebar', 'Wire drag regions', 'Window controls (min/max/close)'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Shell',
          lessons: [
            { slug: 'titlebar', title: 'Custom titlebar', tagline: '--wails-draggable + Tailwind.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'window-controls', title: 'Window controls', tagline: 'Min, max, close, traffic lights.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'auto-update',
      number: 4,
      title: 'In-App Auto-Update',
      tagline: 'Binary swap + GitHub releases + modal UI.',
      learningGoals: ['Wire the updater', 'Build the modal UI', 'Cut a release that auto-distributes'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Updater',
          lessons: [
            { slug: 'updater-go', title: 'updater.go', tagline: 'The Wails binding + binary swap.', minutes: 9, difficulty: 'hard', status: 'coming-soon' },
            { slug: 'modal-ui', title: 'Modal + banner UI', tagline: 'React components shipping in scaffold.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'release-script', title: 'release-desktop.sh', tagline: 'One-shot pipeline.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'installers',
      number: 5,
      title: 'Branded Installers',
      tagline: 'NSIS full + slim, branded MUI bitmaps.',
      learningGoals: ['Customize NSIS', 'Choose between full + slim', 'Generate branded bitmaps'],
      status: 'coming-soon',
      modules: [
        {
          title: 'NSIS',
          lessons: [
            { slug: 'project-nsi', title: 'project.nsi (full)', tagline: 'Bundled WebView2 — no internet needed.', minutes: 8, difficulty: 'hard', status: 'coming-soon' },
            { slug: 'project-slim', title: 'project-slim.nsi', tagline: 'Online bootstrapper — small download.', minutes: 6, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'bitmaps', title: 'Branded bitmaps', tagline: 'header.bmp + welcome.bmp from icon.ico.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
          ],
        },
      ],
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
  status: 'coming-soon',
  emoji: '🌍',
  accent: 'from-amber-500/30 via-orange-500/20 to-amber-500/10',
  chapters: [
    {
      slug: 'foundation',
      number: 1,
      title: 'The Foundation',
      tagline: 'grit new --triple + adding mobile + desktop on top.',
      learningGoals: ['Scaffold the full stack', 'Understand the monorepo wiring', 'Run all four surfaces'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Setup',
          lessons: [
            { slug: 'scaffold', title: 'Scaffolding the full stack', tagline: 'Triple + mobile + desktop.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'monorepo-wiring', title: 'Monorepo wiring', tagline: 'pnpm workspaces, turbo, go module.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'shared-types',
      number: 2,
      title: 'Shared Types Everywhere',
      tagline: 'One source of truth for 4 surfaces.',
      learningGoals: ['grit sync across surfaces', 'Type-safe API client', 'Shared Zod schemas'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Types',
          lessons: [
            { slug: 'grit-sync-multi', title: 'grit sync — multi-surface', tagline: 'One Go model, three TS frontends.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'shared-zod', title: 'Shared Zod schemas', tagline: 'Validate input on every surface.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'feature-implementation',
      number: 3,
      title: 'Building a Feature Across All Three',
      tagline: 'Pick a real feature — implement it everywhere.',
      learningGoals: ['Build a feature in web', 'Implement on mobile', 'Add to desktop'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Cross-surface',
          lessons: [
            { slug: 'pick-feature', title: 'Picking the feature', tagline: 'What ships well across surfaces.', minutes: 5, difficulty: 'easy', status: 'coming-soon' },
            { slug: 'web-impl', title: 'Web implementation', tagline: 'Next.js page + admin.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'mobile-impl', title: 'Mobile implementation', tagline: 'Expo screens.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'desktop-impl', title: 'Desktop implementation', tagline: 'Wails screens.', minutes: 9, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'sync',
      number: 4,
      title: 'Sync + Offline',
      tagline: 'Mobile + desktop sync to the same API.',
      learningGoals: ['Mobile offline reads', 'Desktop outbox', 'Resolve conflicts'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Offline',
          lessons: [
            { slug: 'mobile-offline', title: 'Mobile offline', tagline: 'React Query persistence.', minutes: 8, difficulty: 'medium', status: 'coming-soon' },
            { slug: 'desktop-outbox', title: 'Desktop outbox', tagline: 'Queue + push + pull.', minutes: 9, difficulty: 'hard', status: 'coming-soon' },
            { slug: 'conflicts', title: 'Conflict resolution', tagline: 'Last-write-wins vs. CRDT-light.', minutes: 8, difficulty: 'hard', status: 'coming-soon' },
          ],
        },
      ],
    },
    {
      slug: 'releases',
      number: 5,
      title: 'Coordinated Releases',
      tagline: 'Ship to all surfaces together.',
      learningGoals: ['Version compatibility matrix', 'Stagger releases safely', 'Roll back across surfaces'],
      status: 'coming-soon',
      modules: [
        {
          title: 'Releases',
          lessons: [
            { slug: 'compat-matrix', title: 'Compatibility matrix', tagline: 'Which API versions support which clients.', minutes: 7, difficulty: 'hard', status: 'coming-soon' },
            { slug: 'staggered', title: 'Staggered releases', tagline: 'Web first, then desktop, then mobile.', minutes: 7, difficulty: 'medium', status: 'coming-soon' },
          ],
        },
      ],
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
]
