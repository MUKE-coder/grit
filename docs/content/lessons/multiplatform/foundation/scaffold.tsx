import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        This course is the kitchen sink: one Go API powering four
        frontends — a public website, an admin panel, a mobile app, and
        a desktop app. The first job is scaffolding all of them so
        they can talk to the same database from day one.
      </p>

      <h2>What we&apos;re building</h2>
      <Diagram label="Four surfaces, one API" caption="Every frontend points at the same Go API. The API is the single source of truth.">
{`                  ┌──────────────────────────────┐
                  │   Go API  (apps/api)         │
                  │   :8080  → Postgres          │
                  └────────────┬─────────────────┘
                               │
        ┌──────────┬───────────┼──────────┬──────────────┐
        │          │           │          │              │
   ┌────▼───┐ ┌────▼────┐ ┌────▼────┐ ┌───▼─────┐  ┌─────▼─────┐
   │  Web   │ │  Admin  │ │ Mobile  │ │ Desktop │  │ Shared    │
   │ (3000) │ │ (3001)  │ │ (Expo)  │ │ (Wails) │  │ types/zod │
   └────────┘ └─────────┘ └─────────┘ └─────────┘  └───────────┘`}
      </Diagram>
      <p>
        Same API, same data, four different surfaces. Each surface
        owns presentation; nothing owns business logic except the API.
      </p>

      <h2>Step 1 — Scaffold the triple kit</h2>
      <p>
        Triple kit means: API + web + admin in one monorepo. This is
        the foundation; mobile and desktop bolt on next.
      </p>
      <CodeBlock
        language="bash"
        code={`grit new my-saas --kit=triple
cd my-saas
pnpm install`}
      />
      <p>That gives you:</p>
      <ul>
        <li><code>apps/api/</code> — Go API</li>
        <li><code>apps/web/</code> — Next.js public + dashboard</li>
        <li><code>apps/admin/</code> — Next.js admin panel</li>
        <li><code>packages/shared/</code> — Zod schemas + TS types</li>
      </ul>

      <h2>Step 2 — Add the mobile app</h2>
      <p>
        Mobile lives at <code>apps/mobile/</code>. The Grit CLI knows
        how to add it on top of an existing project.
      </p>
      <CodeBlock
        language="bash"
        code={`grit add mobile
# scaffolds apps/mobile/ as an Expo app
# wires it into pnpm-workspace.yaml
# points API_URL at http://localhost:8080`}
      />
      <p>
        Now the workspace has 4 frontends. Mobile reads from{' '}
        <code>packages/shared</code> just like the web apps do.
      </p>

      <h2>Step 3 — Add the desktop app</h2>
      <CodeBlock
        language="bash"
        code={`grit add desktop
# scaffolds apps/desktop/ as a Wails project
# wires it into the workspace
# points API_URL at http://localhost:8080`}
      />
      <p>
        Wails is Go + React shipped as a single native binary. Same
        shared types, same API client patterns — different runtime.
      </p>

      <TipBox tone="info">
        <strong>Why scaffold in this order?</strong> Triple first, then
        add mobile + desktop. Grit&apos;s shared package is most
        battle-tested in the web direction — mobile and desktop slot
        in by reusing it. Doing it the other way around (mobile first,
        try to retrofit web) usually leads to duplicated types.
      </TipBox>

      <h2>What the final tree looks like</h2>
      <CodeBlock
        language="text"
        code={`my-saas/
├── apps/
│   ├── api/          Go API (Gin + GORM)
│   ├── web/          Next.js public site + dashboard
│   ├── admin/        Next.js admin panel
│   ├── mobile/       Expo + React Native
│   └── desktop/      Wails (Go + React)
├── packages/
│   └── shared/       Zod schemas + TS types (shared by web/admin/mobile/desktop)
├── grit.config.ts    Project config (kits enabled, etc.)
├── pnpm-workspace.yaml
└── turbo.json`}
      />

      <h2>Boot the whole thing</h2>
      <CodeBlock
        language="bash"
        code={`# Terminal 1 — infra
docker compose up -d

# Terminal 2 — API
cd apps/api && go run cmd/server/main.go

# Terminal 3 — web + admin (turbo handles both)
pnpm dev

# Terminal 4 — mobile
cd apps/mobile && pnpm start

# Terminal 5 — desktop
cd apps/desktop && wails dev`}
      />
      <p>
        Five terminals feels like a lot — and it is. Most days you
        won&apos;t run all five. You boot the ones for the surface
        you&apos;re editing. The point is they ALL CAN run against
        the same dev DB.
      </p>

      <KnowledgeCheck
        question="Why scaffold mobile + desktop AFTER triple, not as part of one big command?"
        choices={[
          {
            label: 'The CLI is faster doing it in steps',
            feedback:
              'Speed isn’t the reason — the CLI runs in seconds either way.',
          },
          {
            label: 'Shared types stabilize in the web direction first, then mobile + desktop reuse them',
            correct: true,
            feedback:
              'Right — packages/shared evolves with the web first. Mobile + desktop don’t reshape it; they consume it. Doing it the other way often leads to per-surface duplicate types.',
          },
          {
            label: 'Mobile + desktop need different Go API versions',
            feedback:
              'No — all surfaces hit the same API. There’s only one API.',
          },
          {
            label: 'Required by Apple App Store rules',
            feedback: 'Apple has no opinion on your monorepo layout.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 1&apos;s assignment, get all four surfaces up:</p>
            <ol>
              <li>Run <code>grit new my-saas --kit=triple</code>.</li>
              <li>Run <code>grit add mobile</code> and <code>grit add desktop</code>.</li>
              <li>Boot the API, web, admin, mobile, and desktop.</li>
              <li>Screenshot each — even just the default scaffold pages.</li>
              <li>Confirm all 4 frontends list the same seeded user from <code>/api/users</code>.</li>
            </ol>
            <p>
              The point: prove the loop closes. One API, four frontends, one DB.
            </p>
          </>
        }
        hint={
          <>
            If mobile can&apos;t reach the API on a phone, you&apos;ll
            need your machine&apos;s LAN IP instead of localhost (e.g.,{' '}
            <code>192.168.1.x:8080</code>). The Expo dev tools show
            you the URL it expects.
          </>
        }
        solution={
          <>
            <p>
              You should be able to register a user via the web,
              log into the admin and see them, log into mobile and
              see the same name, and into desktop and see them too.
              The DB has ONE users table. That&apos;s the point.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Monorepo wiring</strong>. The boring
        plumbing that makes four frontends behave like one project:
        pnpm workspaces, turbo pipelines, the Go module layout.
      </p>
    </>
  )
}
