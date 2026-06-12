import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Five surfaces in one repo means five build systems pretending
        to be one. This lesson is the plumbing — how pnpm, Turbo, and
        the Go module work together so you don&apos;t end up running
        every command in every folder.
      </p>

      <h2>pnpm-workspace.yaml — declare your TS packages</h2>
      <CodeBlock
        language="yaml"
        filename="pnpm-workspace.yaml"
        code={`packages:
  - "apps/web"
  - "apps/admin"
  - "apps/mobile"
  - "apps/desktop/frontend"
  - "packages/*"`}
      />
      <p>
        Note: desktop is special. Wails has its own Go module at{' '}
        <code>apps/desktop/</code> AND a React frontend at{' '}
        <code>apps/desktop/frontend/</code>. Only the frontend is in
        the pnpm workspace.
      </p>

      <h2>Importing shared types</h2>
      <p>
        Every TS package declares <code>@my-saas/shared</code> as a
        dependency. pnpm symlinks it so changes are picked up
        immediately:
      </p>
      <CodeBlock
        language="json"
        filename="apps/web/package.json"
        code={`{
  "name": "@my-saas/web",
  "dependencies": {
    "@my-saas/shared": "workspace:*"
  }
}`}
      />
      <p>Now in any TS file:</p>
      <CodeBlock
        language="ts"
        code={`import { Product, ProductSchema } from '@my-saas/shared'`}
      />
      <p>
        Same import, same path, in web / admin / mobile / desktop.
        The shared package becomes muscle memory.
      </p>

      <h2>turbo.json — the build orchestrator</h2>
      <CodeBlock
        language="json"
        filename="turbo.json"
        code={`{
  "$schema": "https://turbo.build/schema.json",
  "pipeline": {
    "dev": {
      "cache": false,
      "persistent": true
    },
    "build": {
      "dependsOn": ["^build"],
      "outputs": ["dist/**", ".next/**", "!.next/cache/**"]
    },
    "lint": {},
    "type-check": {}
  }
}`}
      />
      <p>
        Now <code>pnpm dev</code> at the root runs every package&apos;s
        dev script in parallel, and <code>pnpm build</code> builds them
        in the right order based on dependencies.
      </p>

      <h2>The Go side — separate modules</h2>
      <p>
        Go doesn&apos;t play in the pnpm workspace. There are two Go
        modules:
      </p>
      <ul>
        <li><code>apps/api/go.mod</code> — the API</li>
        <li><code>apps/desktop/go.mod</code> — the Wails desktop wrapper</li>
      </ul>
      <p>
        They&apos;re independent on purpose — the API ships as a Docker
        image; desktop ships as a native binary. Pinning their Go
        deps separately keeps them clean.
      </p>

      <TipBox tone="info">
        <strong>Why not one big go.mod?</strong> The API and the
        desktop wrapper need different deps — the API uses Gin + GORM;
        desktop uses Wails. Sharing a go.mod would force you to vendor
        all of Wails into your server image. Separate modules keep each
        binary small.
      </TipBox>

      <h2>Running it all from the root</h2>
      <CodeBlock
        language="bash"
        code={`# Boots web + admin + mobile + desktop frontend — Turbo handles parallelism
pnpm dev

# API runs separately (Turbo doesn't manage Go)
make api          # or: cd apps/api && go run cmd/server/main.go

# Or use the included Procfile + overmind:
overmind start`}
      />
      <p>
        The Procfile that ships with the scaffold defines every
        service. Install <a href="https://github.com/DarthSim/overmind" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">overmind</a> once and you boot the
        whole monorepo with one command.
      </p>

      <h2>The .env split</h2>
      <p>
        One <code>.env</code> at the root holds shared values
        (DATABASE_URL, JWT_SECRET). Each app reads what it needs.
        Mobile is special — it can&apos;t read root .env files at
        runtime, so its values are baked at build time via{' '}
        <code>EXPO_PUBLIC_*</code> vars.
      </p>
      <CodeBlock
        language="env"
        filename=".env"
        code={`# Shared by API and frontends
DATABASE_URL=postgres://...
JWT_SECRET=...

# Frontend-side overrides
NEXT_PUBLIC_API_URL=http://localhost:8080
EXPO_PUBLIC_API_URL=http://192.168.1.10:8080`}
      />

      <KnowledgeCheck
        question="You add a new field to packages/shared/types/product.ts. What needs to happen for web / admin / mobile / desktop to see it?"
        choices={[
          {
            label: 'Run pnpm install in each app',
            feedback:
              'No — pnpm workspaces symlink packages/shared, so changes propagate without reinstalling.',
          },
          {
            label: 'Nothing — the import resolves to the live file, every dev server picks it up on next reload',
            correct: true,
            feedback:
              'Right — that’s the whole point of the workspace symlink. Save in packages/shared, see it in every consuming app within seconds.',
          },
          {
            label: 'Rebuild each app',
            feedback:
              'Not at dev time — Vite / Next / Metro hot-reload the new types. You rebuild for production.',
          },
          {
            label: 'Re-run grit sync',
            feedback:
              'Sync only runs when you change the Go model. Editing TS directly doesn’t need sync.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Prove the workspace works:</p>
            <ol>
              <li>
                Edit <code>packages/shared/types/user.ts</code> — add{' '}
                <code>nickname: string</code> to the User type.
              </li>
              <li>
                Save. In the web app, type <code>user.</code> in any
                tsx file — autocomplete should show <code>nickname</code>.
              </li>
              <li>
                Same check in admin, mobile, desktop. All four see it
                immediately.
              </li>
              <li>
                Revert the change (we&apos;ll add fields properly via{' '}
                <code>grit sync</code> in the next chapter).
              </li>
            </ol>
          </>
        }
        hint={
          <>
            If autocomplete doesn&apos;t show up, your TS server may be
            stuck on a cached package. Restart it in VS Code (Cmd/Ctrl+
            Shift+P → &quot;TypeScript: Restart TS Server&quot;).
          </>
        }
        solution={
          <>
            <p>
              All four IDEs should see the new field within a second of
              saving. That&apos;s the workspace symlink at work — there
              are no copies, just one source.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Shared Types Everywhere</strong>. We
        stop hand-editing the shared package and use <code>grit sync</code>{' '}
        to generate it from Go.
      </p>
    </>
  )
}
