import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build a SaaS with Grit + Claude Code — Grit Course',
  description: 'Build a project management SaaS using Grit scaffolding and Claude Code for AI-assisted development. Teams, projects, tasks, billing, and deployment.',
}

export default function SaasWithAICourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Build a SaaS with Grit + Claude Code</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build a SaaS with Grit + Claude Code
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            In this course, you will build a real project management SaaS — teams, projects, tasks,
            members, and a billing page — using Grit for the foundation and Claude Code as your
            AI-powered pair programmer for feature development.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What We're Building ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What We{"'"}re Building</h2>

          <Definition term="SaaS (Software as a Service)">
            A web application that users pay to access, typically on a monthly or yearly subscription.
            The app runs on your servers, users access it through their browser. Think Slack, Notion,
            Linear, or Jira. The business model: you build it once, host it once, and sell access to
            many customers.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Our project management SaaS will include:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Teams</strong> — organizations that group members together</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Projects</strong> — belong to a team, have a name, description, and status</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Tasks</strong> — belong to a project, assigned to team members, have priority and due dates</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Members</strong> — users who belong to a team with different roles</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Billing page</strong> — pricing tiers (Free, Pro, Team) for the SaaS model</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The twist: instead of building every feature by hand, you will use Claude Code to accelerate
            development. Grit provides the structure, Claude Code fills in the features.
          </p>

          <Challenge number={1} title="Name 5 SaaS Products">
            <p>Name 5 SaaS products you use daily. For each, identify: what problem it solves, whether it{"'"}s free or paid, and what makes it sticky (why you keep using it instead of switching to an alternative).</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Why Claude Code + Grit? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Why Claude Code + Grit?</h2>

          <Definition term="Claude Code">
            Anthropic{"'"}s AI coding assistant that runs in your terminal. Unlike chat-based AI tools,
            Claude Code can read your entire project, understand its structure, create and modify files,
            run commands, and make changes across multiple files at once. It works with your codebase
            directly — not through copy-paste.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every Grit project ships with a skill file that teaches Claude Code how your project works:
          </p>

          <Definition term="Skill File">
            A markdown file at <Code>.claude/skills/grit/SKILL.md</Code> that teaches AI assistants about
            your project{"'"}s conventions, patterns, and structure. It tells Claude Code things like: where
            models live, how routes are structured, what code markers to use for injections, and how
            the admin panel works. Think of it as an onboarding document — but for AI.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Because Grit follows strict conventions (file locations, naming patterns, code markers), Claude
            Code can work with your project reliably. It knows:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Where to add new models (<Code>internal/models/</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Where to add new handlers (<Code>internal/handlers/</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> How to register routes (using code markers)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> How the admin panel resource system works</li>
            <li className="flex gap-2"><span className="text-primary">•</span> What the API response format looks like</li>
          </ul>

          <Challenge number={2} title="Find the Skill File">
            <p>After scaffolding your project (next section), find the <Code>.claude/skills/</Code> directory. Read the <Code>SKILL.md</Code> file. What conventions does it teach Claude Code about? List at least 5.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Scaffold the SaaS ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold the SaaS</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Start with a Triple architecture scaffold — Go API, Next.js web app, and admin panel.
            This is the standard SaaS setup:
          </p>

          <CodeBlock filename="Terminal">
{`grit new saas-pm --triple --next --style modern`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This gives you:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Go API</strong> — authentication, file storage, email, background jobs, AI service</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Next.js Web</strong> — public-facing app where users interact with projects and tasks</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Admin Panel</strong> — internal dashboard for managing users, data, and system health</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Shared Package</strong> — Zod schemas and TypeScript types shared between frontends</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Docker Compose</strong> — PostgreSQL, Redis, MinIO ready to go</li>
          </ul>

          <Tip>
            The <Code>--style modern</Code> flag gives you a clean, minimal design that works well for
            SaaS products. You can always customize it later.
          </Tip>

          <Challenge number={3} title="Scaffold and Explore">
            <p>Scaffold the project with the command above. Start the dev server with <Code>grit dev</Code>. Visit the web app, admin panel, and API health endpoint. Everything should work out of the box.</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Design the Data Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Design the Data Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Now comes the power of Grit + Claude Code together. Instead of generating one resource at a
            time, ask Claude Code to generate all three in sequence:
          </p>

          <CodeBlock filename="Claude Code Prompt">
{`Generate resources for a project management SaaS:

1. Team (name, slug)
2. Project (name, description, team_id:belongs_to:Team, status)
3. Task (title, description, priority:int, status,
   assignee_id:belongs_to:User,
   project_id:belongs_to:Project,
   due_date:date:optional)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Claude Code reads the skill file and knows to run <Code>grit generate resource</Code> for each
            one. For each resource, Grit creates:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Go model with GORM tags</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Service layer with CRUD operations</li>
            <li className="flex gap-2"><span className="text-primary">•</span> HTTP handlers with validation</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Zod schemas for frontend validation</li>
            <li className="flex gap-2"><span className="text-primary">•</span> TypeScript types</li>
            <li className="flex gap-2"><span className="text-primary">•</span> React Query hooks for data fetching</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Admin panel resource definition</li>
          </ul>

          <Note>
            The <Code>belongs_to</Code> syntax creates foreign key relationships. <Code>team_id:belongs_to:Team</Code> means
            the Project model has a <Code>TeamID</Code> field that references the Team model, with GORM
            automatically setting up the foreign key constraint.
          </Note>

          <Challenge number={4} title="Generate the Resources">
            <p>Ask Claude Code to generate all three resources. Review what it creates — check the Go models, the API routes, and the admin panel. Are the relationships correct? Does <Code>belongs_to</Code> create the expected foreign keys?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: AI-Assisted Feature Development ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">AI-Assisted Feature Development</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With the scaffolding done, use Claude Code to build features that go beyond basic CRUD.
            Here are examples of prompts that work well with Grit projects:
          </p>

          <CodeBlock filename="Claude Code Prompt — Dashboard">
{`Add a dashboard page at /dashboard that shows:
- Total tasks by status (pie chart or count cards)
- My assigned tasks (filtered to current user)
- Recent activity across all projects`}
          </CodeBlock>

          <CodeBlock filename="Claude Code Prompt — Board View">
{`Create a project board view at /projects/:id/board that shows
tasks in columns by status (Todo, In Progress, Review, Done).
Use drag-and-drop to move tasks between columns.`}
          </CodeBlock>

          <CodeBlock filename="Claude Code Prompt — Invitations">
{`Add team member invitation:
- API endpoint POST /api/teams/:id/invite (accepts email)
- Send invitation email using the Resend service
- Create a /join/:token page that accepts the invitation`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The key to effective prompts with Claude Code:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Be specific about URLs</strong> — tell it exactly which route to create</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Describe the UI</strong> — mention cards, tables, charts, columns</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Reference existing patterns</strong> — {'"'}use the same style as the users page{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Review the output</strong> — always read what Claude Code generates before accepting</li>
          </ul>

          <Challenge number={5} title="Build a Feature with AI">
            <p>Choose one of the three prompts above (dashboard, board view, or invitations) and give it to Claude Code. Review the generated code. Does it follow Grit conventions? Did it put files in the right locations? Make at least one refinement request.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Custom Dashboard ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Custom Dashboard</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every SaaS needs a dashboard — the first thing users see after logging in. Use Claude Code
            to build a data-rich dashboard with multiple widgets:
          </p>

          <CodeBlock filename="Claude Code Prompt">
{`Build a custom dashboard at /dashboard with these widgets:

1. Stats row: Total Projects, Active Tasks, Completed This Week,
   Overdue Tasks (4 cards with icons and counts)
2. Tasks by Status: bar chart showing task distribution
3. Tasks by Priority: colored badges showing High/Medium/Low counts
4. My Tasks: table of tasks assigned to me, sorted by due date
5. Team Activity: recent actions (task created, status changed, etc.)

Use the existing API endpoints. Add new endpoints if needed.`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Claude Code will create the frontend components and any custom API endpoints needed for
            aggregated data. It knows to use TanStack Query for data fetching and shadcn/ui components
            for the UI because the skill file tells it to.
          </p>

          <Challenge number={6} title="Build a Dashboard">
            <p>Ask Claude Code to build a dashboard with at least 3 different widgets. After it generates the code, check: Does the data look correct? Are the API calls efficient (not N+1 queries)? Does the layout look good on both desktop and mobile?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: API Customization ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">API Customization</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit generates standard CRUD endpoints, but a real SaaS needs custom business logic.
            Use Claude Code to add endpoints that go beyond CRUD:
          </p>

          <CodeBlock filename="Custom Endpoints">
{`# Tasks for a specific project
GET /api/projects/:id/tasks

# Assign a task to a team member
PUT /api/tasks/:id/assign
Body: { "assignee_id": "uuid-here" }

# Dashboard statistics (aggregated)
GET /api/dashboard/stats
Response: {
  "total_projects": 12,
  "active_tasks": 45,
  "completed_this_week": 8,
  "overdue_count": 3
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When asking Claude Code to create custom endpoints, be explicit about the request and
            response format. Grit{"'"}s API follows a consistent pattern (data envelope, meta for lists,
            error format), and Claude Code knows this from the skill file.
          </p>

          <CodeBlock filename="Claude Code Prompt">
{`Add a custom endpoint GET /api/dashboard/stats that returns:
- total_projects: count of all projects for the current user's team
- active_tasks: count of tasks with status != "done"
- completed_this_week: count of tasks completed in the last 7 days
- overdue_count: count of tasks past due_date with status != "done"

Follow the existing API response format. Add the route
to the authenticated group.`}
          </CodeBlock>

          <Challenge number={7} title="Add a Custom Endpoint">
            <p>Ask Claude Code to create the <Code>GET /api/dashboard/stats</Code> endpoint. Test it with <Code>curl</Code> or your API client. Does it return accurate counts? Is it behind authentication?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Frontend Pages ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Frontend Pages</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With the API in place, build the key frontend pages. These are where Claude Code really
            shines — it can create complex React components much faster than writing them by hand:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Project Board</strong> — Kanban-style columns (Todo, In Progress, Review, Done) with tasks as cards</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Task Detail</strong> — full task view with description, assignee, comments, status changes</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Team Settings</strong> — manage members, roles, and team information</li>
          </ul>

          <CodeBlock filename="Claude Code Prompt">
{`Build a project board page at /projects/[id]/board:
- Fetch tasks for this project from GET /api/projects/:id/tasks
- Display 4 columns: Todo, In Progress, Review, Done
- Each task shows: title, priority badge, assignee avatar, due date
- Clicking a task opens the task detail page
- Use the existing shadcn/ui Card component for task cards`}
          </CodeBlock>

          <Tip>
            When Claude Code generates a complex page, test it immediately. If something doesn{"'"}t look
            right, describe the issue and ask for a fix. Iterating with AI is much faster than rewriting
            from scratch.
          </Tip>

          <Challenge number={8} title="Build a Custom Page">
            <p>Choose one of the three pages above and ask Claude Code to build it. Test the page — does it fetch data correctly? Does the layout work? Ask Claude Code to fix at least one issue you find.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Billing Page ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Billing Page (UI Only)</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Every SaaS needs a pricing page. For this course, you will build the UI — the visual layout
            of pricing tiers. Actual payment processing (Stripe, LemonSqueezy) would be added later
            via a Grit plugin.
          </p>

          <CodeBlock filename="Claude Code Prompt">
{`Create a pricing page at /pricing with 3 plan cards:

Free Plan ($0/month):
- Up to 3 projects
- Up to 10 tasks per project
- 1 team member
- Basic support

Pro Plan ($12/month) — highlighted as "Most Popular":
- Unlimited projects
- Unlimited tasks
- Up to 10 team members
- Priority support
- File attachments

Team Plan ($29/month):
- Everything in Pro
- Unlimited team members
- Admin controls
- Audit log
- Custom integrations
- Dedicated support

Include a monthly/yearly toggle (yearly = 20% discount).
Use shadcn/ui Card components with the Grit dark theme.`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The pricing page is purely frontend — no API calls needed. It{"'"}s a static page that
            shows the plans and (eventually) links to a payment provider.
          </p>

          <Challenge number={9} title="Build a Pricing Page">
            <p>Ask Claude Code to build the pricing page. Review the design: Does the {'"'}Most Popular{'"'} card stand out? Does the monthly/yearly toggle work? Is the dark theme consistent with the rest of the app?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Deploy ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Deploy Your SaaS</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With your SaaS built, deploy it to a server. Grit makes this a single command:
          </p>

          <CodeBlock filename="Terminal">
{`grit deploy --host deploy@server.com --domain pm.myapp.com`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This builds the Go binary, uploads it, configures systemd and Caddy, and starts the service.
            Your SaaS is live at <Code>https://pm.myapp.com</Code> with automatic HTTPS.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before deploying, make sure your production environment variables are set:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>APP_ENV=production</Code></li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>JWT_SECRET</Code> — a strong random string</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>DB_PASSWORD</Code> — a strong database password</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>RESEND_API_KEY</Code> — for real email delivery</li>
          </ul>

          <Note>
            For a more comprehensive deployment guide covering Dokploy, Orbita, and Vercel, see the
            <Link href="/courses/deployment-guide" className="text-primary hover:underline ml-1">Deployment Guide course</Link>.
          </Note>

          <Challenge number={10} title="Deploy Your SaaS">
            <p>Deploy your project management SaaS to a VPS. Create an account, log in, create a team, add a project, and create a few tasks. Does everything work in production exactly as it did in development?</p>
          </Challenge>
        </section>

        {/* ═══ Section 11: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What You Learned</h2>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> What a SaaS is and how to structure one with Grit</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> How Claude Code + the Grit skill file accelerate development</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Generating related resources with <Code>belongs_to</Code> relationships</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Using Claude Code to build features beyond CRUD</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building custom dashboards with data widgets</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Adding custom API endpoints for business logic</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Creating complex frontend pages with AI assistance</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Building a SaaS pricing page</li>
            <li className="flex gap-2"><CheckCircle2 className="h-4 w-4 text-green-500 shrink-0 mt-1" /> Deploying a SaaS to production</li>
          </ul>

          <Challenge number={11} title="Add 3 More Features">
            <p>Using Claude Code, add these three features to your SaaS:</p>
            <ol className="list-decimal list-inside space-y-1 mt-2">
              <li><strong className="text-foreground">Team member roles</strong> — owner, admin, member. Only owners and admins can invite members or delete projects.</li>
              <li><strong className="text-foreground">Task comments with mentions</strong> — users can comment on tasks and @mention team members.</li>
              <li><strong className="text-foreground">Email notifications</strong> — send an email when a task is assigned to someone, using Grit{"'"}s Resend integration.</li>
            </ol>
            <p className="mt-2">For each feature, write a clear prompt, review the output, and test it end-to-end.</p>
          </Challenge>

          <Challenge number={12} title="Ship It">
            <p>Deploy the final version with all features. Create a demo account and walk through the entire workflow: sign up, create team, invite a member, create project, add tasks, assign tasks, comment on a task. Time yourself — how long does the full workflow take? Is the UX smooth?</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
