import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { DocsShell } from '@/components/docs/shell'
import { Button } from '@/components/ui/button'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/expert-skills')

// Every row is the skill's own description, taken from its SKILL.md
// frontmatter rather than written here. A page that paraphrases twenty-eight
// skills is twenty-eight chances to describe one of them wrongly, and it
// drifts the first time one changes.
const GROUPS = [
  {
    title: 'Start here',
    skills: [
      { name: 'grit-build', does: 'The router for a Grit project. Reads the project, interviews you once, then assigns the work to the right experts in the right order and runs them.', when: 'build X, add this feature, what should I do next, set up this project, or any request big enough to touch more than one layer.' },
      { name: 'grit-core', does: 'The Grit framework itself: conventions, the API contract, the folder structure, what belongs in the framework versus what you build, and how to model a domain before any code is generated.', when: 'how does Grit do X, where should this live, model this domain, is this the right shape, or any decision that the rest of the build will be derived from.' },
      { name: 'grit-setup', does: 'Install and check what the Grit experts depend on: the external skills they use for design research and browser testing, the GitHub CLI for filing issues, and the Grit binary itself.', when: 'set up the grit experts, install the dependencies, or when another expert says a skill is missing.' },
    ],
  },
  {
    title: 'The architecture you chose',
    skills: [
      { name: 'grit-triple', does: 'The triple architecture: the Go API, a customer-facing Next.js web app, and a separate Next.js admin, in one pnpm monorepo with a shared package.', when: 'grit.json says triple, for scaffold a triple project, wire this resource into web and admin, or any monorepo or shared-package problem.' },
      { name: 'grit-double', does: 'The double architecture: the Go API plus one frontend, with the admin as a route group inside it.', when: 'grit.json says double, for scaffold a double project, where does the admin live here, or when deciding between double and triple.' },
      { name: 'grit-single', does: 'The single architecture: one app that is both the frontend and the project root, with the Go API beside it and the admin as a section.', when: 'grit.json says single, for scaffold a single project, where is the shared package here, or Vite and TanStack questions in a Grit project.' },
    ],
  },
  {
    title: 'Building',
    skills: [
      { name: 'grit-resource', does: 'Run Grit\'s generator correctly: field types, the flags, inline fields versus a YAML definition, conditional fields, relationships, and removing a resource cleanly.', when: 'generate a resource, add a field, what field type should this be, grit generate, grit remove.' },
      { name: 'grit-data', does: 'GORM, migrations, indexes, seeders, queries and Postgres in a Grit project.', when: 'the query is slow, add an index, write a migration, seed some data, this hook is not firing, import a spreadsheet, or anything where the database is the thing being changed.' },
      { name: 'grit-admin', does: 'The Grit admin panel: resource definitions, the data table, the form builder, dashboard widgets, the sidebar, bulk actions and the system pages.', when: 'add an admin screen, the admin table, a dashboard widget, the form is wrong, customise the admin, the sidebar.' },
      { name: 'grit-ui-ux', does: 'Design the customer-facing surfaces: research a pattern before building it, hold the result to the anti-slop standard, and get motion and accessibility right.', when: 'design this page, make it look better, it looks generic, design review, add animation, a landing page, or any screen a customer sees.' },
      { name: 'grit-plugins', does: 'Grit plugins: what already exists so you do not rebuild it, whether a feature should be a plugin at all, writing one from scratch, and proposing one for Grit to ship officially.', when: 'is there a plugin for this, write a plugin, should this be a plugin, grit plugin add, package this up, suggest a plugin, multi-tenancy, a kanban board, a calendar.' },
      { name: 'grit-blueprints', does: 'Grit blueprints: what complete applications already exist to start from or read, whether the thing being built should become one, working inside a blueprint repository, and proposing a new one.', when: 'is there an example of this, what should I start from, how did they do chat, make this a blueprint, grit blueprint, suggest a blueprint.' },
    ],
  },
  {
    title: 'Other surfaces',
    skills: [
      { name: 'grit-mobile', does: 'The Expo app: screens against the shared API and types, offline field use, uploads, images, haptics, and getting through App Store and Play review.', when: 'add a mobile app, the Expo app, field mode, submit to the store, my app was rejected, the image will not load on my phone.' },
      { name: 'grit-desktop', does: 'The Wails desktop app: offline-first with a local SQLite mirror and a sync outbox, Go bindings to React, packaging, and shipping updates through GitHub Releases.', when: 'add a desktop app, make it work offline, the desktop build, auto-update, the installer.' },
      { name: 'grit-offline', does: 'Design and build the offline engine for a desktop or mobile app: the local mirror, the write outbox, sync and replay, conflict rules, and what must never be queued.', when: 'make it work offline, offline first, field mode, sync when the connection comes back, the outbox, it works but not on site.' },
    ],
  },
  {
    title: 'Before you ship',
    skills: [
      { name: 'grit-reviewer', does: 'Read the diff and find the bugs before a user does, with the failure modes a Grit project actually has.', when: 'review this, check my changes, did I break anything, review before I commit, or as the last step of any build.' },
      { name: 'grit-critic', does: 'The independent review of a feature somebody has marked done: runs it, reads the backend for bottlenecks, profiles it, looks at the UI, checks the tests actually test it, and says what to fix in priority order.', when: 'review this feature, is this really done, quality check, find the bottlenecks, it is marked complete, sign this off, or before any demo.' },
      { name: 'grit-design-review', does: 'Assess a page against Grit\'s visual contract in DESIGN.md: measures the budgets (type styles, radii, rhythm, shadows, contrast, measure, target size), diffs it against the designer\'s mockup, then looks at it, then says what to fix in order.', when: 'review this page, does this meet our design standard, design review, is this page consistent, check the landing page, does this match the design, before I ship this page, or after building any new page or section.' },
      { name: 'grit-security', does: 'Auth, authorization, the WAF and the scanner in a Grit project: sessions, RBAC, per-resource policies, Sentinel, Scout, API keys and the route table.', when: 'is this secure, add roles, who can do this, set up the WAF, run a security scan, a pentest found this, lock down this endpoint.' },
      { name: 'grit-perf', does: 'Make a Grit app fast, starting from a measurement: Pulse for profiling and N+1 detection, k6 for load, caching, pagination and query work.', when: 'it is slow, load test this, find the N+1, why is this endpoint slow, add caching, performance review.' },
    ],
  },
  {
    title: 'Operations',
    skills: [
      { name: 'grit-deploy', does: 'Get a Grit project onto a server and keep it there: Docker Compose, Dokploy on a VPS, environment variables, storage, migrations on deploy, CI, and the desktop release.', when: 'deploy this, set up the VPS, dokploy, production env vars, the build fails in CI, ship a release.' },
      { name: 'grit-adopt', does: 'Bring an existing Grit project up to the current framework without breaking what already works, and decide which new features are worth adopting afterwards.', when: 'Also checks whether the core batteries (Sentinel, Pulse, GORM Studio, Scout) are in step. Use for upgrade this project, we are on an old Grit, is this project out of date, what would upgrading cost, are our versions in sync, how do we adopt the new features, migrate this project.' },
    ],
  },
  {
    title: 'About the work itself',
    skills: [
      { name: 'grit-report', does: 'Report on how the experts are actually performing on this project: route accuracy, whether they verified, how much they asked, framework friction and rework.', when: 'how are the experts doing, grit report, was that efficient, which expert is wasting time, or at the end of a phase.' },
      { name: 'grit-tokens', does: 'What this project has spent on tokens, where it went and who the big spender is, read from Claude Code\'s own transcripts. Alerts at a threshold and reports every increment, and writes an HTML page with the full breakdown.', when: 'how many tokens have we used, what is this costing, token report, why is my context filling up, which tool is expensive, or set it running to be told as the number grows.' },
      { name: 'grit-feedback', does: 'File what Grit got wrong as a GitHub issue, and check later whether the fix shipped.', when: 'an expert hits a framework bug, a missing feature, or has to hand-write something Grit should have given it, and for has that issue been fixed, check my grit issues, what did we file.' },
      { name: 'grit-advocate', does: 'Watches the work and files what is worth telling people about as an issue in MUKE-coder/grit-stories, each with the evidence, the channel it suits and a draft ready to post.', when: 'is there a story in this, write that up, what should I post about, turn this into a LinkedIn post, case study, what have we got to share, or at the end of a build.' },
      { name: 'grit-teacher', does: 'Turns a feature that just shipped into a tutorial a beginner can follow in a fresh project, with the real-world scenario it is for and a checkpoint at every step. Writes one markdown file per lesson into .grit/tutorials/.', when: 'write a tutorial for this, teach this to a beginner, how would someone new build this, explain this for a learner, turn this into a lesson, document how to do it from scratch, or after a feature is verified and shipped.' },
    ],
  },
  {
    title: 'Read by the others',
    skills: [
      { name: 'grit-references', does: 'Reference material the Grit experts read on demand.', when: 'Not a skill to invoke on its own: the other grit-* skills name the file they need and read it from this directory. If you are looking for what to do, use grit-build.' },
    ],
  },
]

const CALLED_A_SKILL = [
  {
    what: 'this pack',
    is: 'The twenty-seven above, plus the reference material. Installed by grit skills install into the configuration your agent reads, at user level or per project.',
  },
  {
    what: 'skills/grit/SKILL.md',
    is: 'One general skill, in the Grit repository. Smaller, and not the same thing as the pack.',
  },
  {
    what: '.claude/skills/grit/SKILL.md',
    is: 'Written into your project by grit new, and tailored to the architecture you chose.',
  },
  {
    what: 'GRIT_SKILL.md',
    is: 'A single file of framework context to paste into an assistant that has no skill mechanism.',
  },
]

export default function ExpertSkillsPage() {
  return (
    <DocsShell>
      <Link
        href="/docs"
        className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-8"
      >
        <ArrowLeft className="h-3.5 w-3.5" />
        Docs
      </Link>

      <h1 className="text-3xl font-bold tracking-tight mb-4">The expert skills</h1>

      <p className="text-muted-foreground leading-relaxed mb-6">
        Twenty-seven skills that teach an AI assistant how to build with Grit,
        plus the reference material they read. They are not a prompt: each one
        reads the docs out of the <code className="text-xs">grit</code> binary
        rather than from memory, interviews you once and remembers the answers,
        and files what Grit got wrong as an issue while it works.
      </p>

      <CodeBlock
        language="bash"
        code={`grit skills install          the pack, at user level
grit skills install --all    the pack and the five skills it works better with
grit skills status           what is installed`}
      />

      <p className="text-muted-foreground leading-relaxed mb-10">
        The pack is{' '}
        <a
          href="https://github.com/MUKE-coder/grit-expert-skills"
          className="underline underline-offset-4 hover:text-foreground"
        >
          MUKE-coder/grit-expert-skills
        </a>
        , and it is public, so you can read every skill before installing it.
        That is worth doing: a skill tells an agent how to write code in your
        project.
      </p>

      <div className="rounded-xl border border-border/50 bg-card/40 p-5 mb-10">
        <h2 className="font-semibold mb-2">Start with one of them</h2>
        <p className="text-sm text-muted-foreground leading-relaxed">
          <code className="text-xs">grit-build</code> is the router. It reads
          the project, works out which architecture it is, asks you the
          questions it cannot answer from the files, and then calls the others
          in order. You do not have to pick from the list below: the list is
          there so you can see what it will reach for, and so you can name one
          directly when you already know which you want.
        </p>
      </div>

      {GROUPS.map((group) => (
        <section key={group.title} className="mb-10">
          <h2 className="text-xl font-bold tracking-tight mb-4">{group.title}</h2>
          <div className="overflow-x-auto rounded-xl border border-border/50">
            <table className="w-full text-sm">
              <thead className="bg-card/60">
                <tr className="text-left">
                  <th className="px-4 py-3 font-semibold whitespace-nowrap">Skill</th>
                  <th className="px-4 py-3 font-semibold">What it does</th>
                  <th className="px-4 py-3 font-semibold">When it fires</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/40">
                {group.skills.map((skill) => (
                  <tr key={skill.name}>
                    <td className="px-4 py-3 font-mono text-[13px] align-top whitespace-nowrap">
                      {skill.name}
                    </td>
                    <td className="px-4 py-3 text-muted-foreground leading-relaxed align-top">
                      {skill.does}
                    </td>
                    <td className="px-4 py-3 text-muted-foreground leading-relaxed align-top">
                      {skill.when || <span className="opacity-60">called by the others</span>}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>
      ))}

      <h2 className="text-xl font-bold tracking-tight mb-4">Four things get called a Grit skill</h2>
      <p className="text-muted-foreground leading-relaxed mb-4">
        They do different jobs, and knowing which one you are reading matters
        more than it sounds. Only the first is this pack:
      </p>
      <div className="overflow-x-auto rounded-xl border border-border/50 mb-4">
        <table className="w-full text-sm">
          <tbody className="divide-y divide-border/40">
            {CALLED_A_SKILL.map((row) => (
              <tr key={row.what}>
                <td className="px-4 py-3 font-mono text-[13px] align-top whitespace-nowrap">
                  {row.what}
                </td>
                <td className="px-4 py-3 text-muted-foreground leading-relaxed">{row.is}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <p className="text-muted-foreground leading-relaxed mb-10">
        The last of those is for an assistant with no skill mechanism at all, and{' '}
        <Link href="/docs/ai-skill" className="underline underline-offset-4 hover:text-foreground">
          the LLM skill guide
        </Link>{' '}
        covers it.
      </p>

      <div className="flex items-center justify-between border-t border-border/40 pt-6">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/ai-workflows/llms-txt" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            llms.txt
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/ai-skill" className="gap-1.5">
            LLM Skill Guide
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
