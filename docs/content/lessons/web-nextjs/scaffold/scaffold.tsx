import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { PrereqLinks } from '@/components/course/prereq-links'

export default function Lesson() {
  return (
    <>
      <p>
        Welcome to <strong>Building Web with Next.js + Go API</strong>.
        We&apos;re using the triple kit — three apps in one monorepo (web,
        admin, api). By the end of this lesson, you&apos;ll have scaffolded a
        SaaS-shaped project ready for the rest of the course.
      </p>

      <PrereqLinks
        prereqs={['nextjs', 'golang', 'docker']}
        intro={
          <>
            New to App Router or Go? Read the relevant primers first —
            this course assumes basic familiarity with both, plus
            running Docker containers locally.
          </>
        }
      />

      <h2>The command</h2>
      <CodeBlock terminal code={`grit new saas --triple`} />
      <p>
        Pick &quot;Next.js&quot; when prompted for the frontend framework (or use{' '}
        <code>--next</code> to skip the prompt). You&apos;ll end up with{' '}
        <code>apps/web</code>, <code>apps/admin</code>, and{' '}
        <code>apps/api</code> all wired together.
      </p>

      <h2>The shape</h2>
      <CodeBlock
        language="text"
        filename="saas/"
        code={`saas/
├── apps/
│   ├── api/              Go (Gin + GORM) backend
│   ├── web/              Public Next.js site + customer dashboard
│   └── admin/            Filament-style admin panel (Next.js)
├── packages/shared/      Zod schemas + TS types
├── docker-compose.yml
├── turbo.json
├── pnpm-workspace.yaml
└── grit.json`}
      />

      <h2>Three different audiences</h2>
      <ul>
        <li>
          <strong>apps/web</strong> — your customers. Marketing pages
          (public) + dashboard (logged-in).
        </li>
        <li>
          <strong>apps/admin</strong> — your staff. Different auth scope,
          different look, behind its own subdomain in prod.
        </li>
        <li>
          <strong>apps/api</strong> — same Go API serves both. Same code,
          different route groups for different roles.
        </li>
      </ul>

      <TipBox tone="info">
        <strong>Why three apps?</strong> Different audiences, different
        security posture (you can IP-allowlist admin), different SEO needs
        (marketing optimised, admin isn&apos;t indexed). Separating them at
        build time is cleaner than trying to gate one app by role.
      </TipBox>

      <h2>First run</h2>
      <CodeBlock
        terminal
        code={`cd saas
docker compose up -d        # postgres, redis, minio, mailhog
pnpm install                # installs every workspace's deps
grit migrate                # creates DB tables
grit seed                   # creates seeded admin user
grit start                  # API + web + admin together`}
      />
      <p>
        Three URLs come up:{' '}
        <code>http://localhost:3000</code> (web),{' '}
        <code>http://localhost:3001</code> (admin),{' '}
        <code>http://localhost:8080</code> (API).
      </p>

      <KnowledgeCheck
        question="Why does Grit ship the admin as a SEPARATE Next.js app instead of an admin section inside apps/web?"
        choices={[
          {
            label: 'Performance — admin pages are heavier',
            feedback:
              "Possible but minor. The deeper reason is audience separation.",
          },
          {
            label: 'Different audiences, different auth scopes, different deploy + SEO needs',
            correct: true,
            feedback:
              "Right — staff and customers should never share an auth surface, the marketing pages should be SEO-optimised while admin shouldn't be indexed, and admin can sit behind a VPN. Two apps make all three trivial.",
          },
          {
            label: 'Next.js can\'t handle both in one app',
            feedback:
              "Wrong — it can. Grit just doesn't.",
          },
          {
            label: 'apps/admin lets you use a different framework',
            feedback:
              "Both apps are Next.js. The split is about audience, not framework choice.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Scaffold and run the full triple stack:
            </p>
            <ol>
              <li><code>grit new saas --triple --next</code></li>
              <li><code>cd saas && docker compose up -d</code></li>
              <li><code>pnpm install</code></li>
              <li><code>grit migrate && grit seed</code></li>
              <li><code>grit start</code></li>
              <li>
                Open all three URLs in different browser tabs. Paste each
                tab&apos;s title bar text in <code>notes.md</code>.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If port 3000 or 3001 is busy, Next.js auto-shifts to the next
            available port. Look at the <code>grit start</code> output for
            the actual URLs.
          </>
        }
        solution={
          <>
            <p>Three browser titles to paste:</p>
            <CodeBlock
              language="text"
              code={`Web    — "saas | Public Site"
Admin  — "saas | Admin"
API    — JSON envelope from /api/health`}
            />
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Three apps running. Next lesson — tour each app&apos;s file
        structure so you know where to put new code.
      </p>
    </>
  )
}
