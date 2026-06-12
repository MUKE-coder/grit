import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        <code>grit new</code> is the command you&apos;ll run the most for the
        first month. This lesson covers what it does, what the prompts mean,
        and the flags that let you skip the prompts when you know what you
        want.
      </p>

      <h2>The simplest form</h2>
      <CodeBlock
        terminal
        code={`grit new my-first-grit`}
      />
      <p>
        That creates a folder called <code>my-first-grit/</code>, prompts
        you for two choices (architecture mode + frontend framework), then
        generates a full project. The interactive flow is friendly — arrow
        keys, enter to confirm, escape to cancel.
      </p>

      <h2>The two prompts</h2>

      <h3>1. Architecture</h3>
      <p>The first prompt asks which kit shape you want:</p>
      <CodeBlock
        language="text"
        code={`? Select architecture:
  > Triple — Web + Admin + API (Turborepo)         [the SaaS shape]
    Double — Web + API (Turborepo)                 [no admin panel]
    Single — Go API + embedded React SPA           [one binary]
    API Only — Go API (no frontend)                [headless backend]
    Mobile — Expo + Go API                         [shared types]`}
      />
      <p>
        For this course, pick <strong>Triple</strong>. We cover every kit in
        chapter 5; for now Triple is the most representative — three apps,
        admin panel, the full experience.
      </p>

      <h3>2. Frontend framework</h3>
      <CodeBlock
        language="text"
        code={`? Select frontend framework:
  > Next.js (App Router)        [SSR, marketing-friendly]
    Vite + TanStack Router      [pure SPA, faster dev]`}
      />
      <p>Pick <strong>Next.js</strong> for now — easier to follow along.</p>

      <h2>Skipping prompts with flags</h2>
      <p>
        Once you know what you want, flags let you scaffold without
        interaction:
      </p>

      <CodeBlock
        terminal
        code={`grit new my-app --triple --next       # Triple + Next.js
grit new my-app --single --vite       # Single binary + Vite SPA
grit new my-app --api                  # API-only
grit new my-app --mobile               # Mobile + API
grit new my-app --triple --desktop     # Triple + Wails desktop`}
      />

      <h2>What &quot;Triple&quot; produces</h2>
      <p>
        After the prompts, you&apos;ll see a printed checklist of what was
        created and how to start:
      </p>

      <CodeBlock
        language="text"
        code={`  ✓ Created my-first-grit/
  ✓ Wrote 84 files
  ✓ Initialized go module

  Next steps:
    cd my-first-grit
    docker compose up -d        # starts Postgres, Redis, MinIO
    pnpm install                # installs JS deps
    grit migrate                # runs DB migrations
    grit seed                   # seeds the admin user
    grit start                  # runs API + web + admin together

  URLs (once running):
    Web:    http://localhost:3000
    Admin:  http://localhost:3001
    API:    http://localhost:8080
    Studio: http://localhost:8080/studio   (DB browser)
    Mail:   http://localhost:8025          (Mailhog UI)`}
      />

      <TipBox tone="success">
        That checklist is your roadmap for the next two lessons. Don&apos;t worry
        about running anything yet — we&apos;ll go file-by-file first.
      </TipBox>

      <h2>Naming rules</h2>
      <p>The project name has to be:</p>
      <ul>
        <li>Lowercase letters, numbers, and dashes only</li>
        <li>Start with a letter</li>
        <li>Not collide with an existing directory</li>
      </ul>
      <p>
        Good: <code>invoice-tracker</code>, <code>order-pos</code>,{' '}
        <code>my-first-grit</code>. Bad: <code>My_App</code>,{' '}
        <code>123-app</code>, <code>app with spaces</code>.
      </p>

      <KnowledgeCheck
        question="You want to scaffold a project called `field-pos` with the Triple kit + Vite frontend. What's the right command?"
        choices={[
          {
            label: 'grit new field-pos --triple --vite',
            correct: true,
            feedback:
              "Right — kit flag (--triple) + frontend flag (--vite) together skip both interactive prompts.",
          },
          {
            label: 'grit new field-pos --triple-vite',
            feedback:
              "Wrong — flags compose individually. There's no --triple-vite combined flag.",
          },
          {
            label: 'grit new field-pos --kit=triple --frontend=vite',
            feedback:
              "Wrong syntax — the long-form flags would be --arch=triple --frontend=vite (with =) but the shorter --triple --vite is the documented usage.",
          },
          {
            label: 'grit new field-pos triple vite',
            feedback:
              "Wrong — positional args are interpreted as the project name. Anything after the first arg without -- is ignored.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Scaffold a project called <code>my-first-grit</code> with the
              triple kit + Next.js frontend.
            </p>
            <p>
              Then{' '}
              <strong>list what was created</strong> by running{' '}
              <code>ls my-first-grit</code> (Mac/Linux) or{' '}
              <code>dir my-first-grit</code> (Windows). Paste the output in{' '}
              <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If you want to skip the prompts, use{' '}
            <code>grit new my-first-grit --triple --next</code>. If you want
            the interactive prompts (recommended for the first time so you
            see them), just <code>grit new my-first-grit</code>.
          </>
        }
        solution={
          <>
            <p>You should see something like:</p>
            <CodeBlock
              terminal
              code={`$ ls my-first-grit
.env             .gitignore       grit.json        README.md
apps/            docker-compose.yml  package.json     turbo.json
docker-compose.prod.yml  packages/  pnpm-workspace.yaml`}
            />
            <p>
              Don&apos;t worry about identifying each file yet — the next lesson
              is dedicated to touring every one. The point of this exercise
              is just to confirm <code>grit new</code> worked.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You have a project on disk. The next lesson is the <strong>guided
        tour</strong> — we walk into <code>apps/</code>,{' '}
        <code>packages/</code>, the root config files, and you&apos;ll be able to
        explain every file to a teammate.
      </p>
    </>
  )
}
