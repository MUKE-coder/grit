import Link from 'next/link'
import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Single mode is the smallest possible Grit deploy: one Go binary with
        the React frontend embedded inside it. No separate frontend server,
        no reverse proxy juggling — drop the binary on a VPS and run it.
      </p>

      <h2>The shape</h2>
      <CodeBlock
        terminal
        code={`grit new my-app --single        # Next.js frontend embedded
grit new my-app --single --vite # Vite + TanStack frontend embedded`}
      />

      <p>
        You get ONE directory (not <code>apps/</code>). The Go code lives at
        the root; the React app lives in <code>frontend/</code> and gets
        compiled into the Go binary at build time via <code>go:embed</code>.
      </p>

      <CodeBlock
        language="text"
        filename="my-app/"
        code={`my-app/
├── main.go              # Go entry point
├── internal/            # Handlers, services, models
├── frontend/            # React source
│   ├── src/
│   └── dist/            # Built at build time, embedded into the binary
├── grit.json
└── go.mod`}
      />

      <h2>Why this is special</h2>
      <ul>
        <li>
          <strong>One file deploy.</strong>{' '}
          <code>scp myapp user@host:</code> + <code>./myapp</code> on the box.
          Done.
        </li>
        <li>
          <strong>No CORS.</strong> The frontend is served by the same Go
          process as the API — same origin, no preflight requests.
        </li>
        <li>
          <strong>No reverse proxy required.</strong> Optional, of course —
          Caddy / Traefik for TLS — but not architecturally needed.
        </li>
        <li>
          <strong>Smallest infra surface.</strong> Postgres + one binary.
          That&apos;s it.
        </li>
      </ul>

      <h2>When it shines</h2>
      <ul>
        <li>Solo / indie products</li>
        <li>Internal tools deployed to a single VPS</li>
        <li>Self-hosted apps you ship to customers (one binary they can run)</li>
        <li>MVPs you want to iterate on without infra complexity</li>
      </ul>

      <h2>When to skip it</h2>
      <ul>
        <li>
          You need a separate admin panel for non-engineers (no admin
          surface in single mode by default — use triple)
        </li>
        <li>
          The web and the API will be deployed independently (different
          regions, different scaling)
        </li>
        <li>You&apos;re scaling to multiple frontend regions but one backend</li>
      </ul>

      <TipBox tone="success">
        <strong>Single + Vite</strong> (<code>--single --vite</code>) is the
        Grit demo at <Link href="/docs/demo" className="text-primary hover:underline">demo.gritframework.dev</Link>{' '}
        — sub-second cold starts, the same single-binary deploy. Worth
        considering when SEO doesn&apos;t matter.
      </TipBox>

      <KnowledgeCheck
        question="You're shipping a self-hosted invoicing app to customers (each runs it on their own VPS). Which mode?"
        choices={[
          {
            label: 'Triple — gives them a real admin panel',
            feedback:
              "Possible, but heavy. Three apps + Postgres + Redis + Turborepo is a lot for a self-hosted box.",
          },
          {
            label: 'Single — they run one binary',
            correct: true,
            feedback:
              "Right — self-hosted means the customer wants minimum infra. Single mode = one binary + Postgres. Much easier to support.",
          },
          {
            label: 'API-only — they bring their own UI',
            feedback:
              "Wrong for this case — customers want a working product, not an API to integrate with.",
          },
          {
            label: 'Desktop — Wails native app',
            feedback:
              "Possible if they want desktop-only, but a typical 'self-hosted' invoicing app is web-based so users can access from anywhere on the LAN.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In <code>notes.md</code>, list three real apps you&apos;d ship as
              single mode. For each, write one sentence on why single
              beats triple for that case.
            </p>
          </>
        }
        hint={
          <>
            Self-hosted, solo, indie, internal-tool, MVP — all flags for
            &quot;single mode is probably right&quot;.
          </>
        }
        solution={
          <>
            <p>Examples:</p>
            <ol>
              <li>
                A LAN-only POS for a coffee shop — single binary on their
                router box, no separate admin needed.
              </li>
              <li>
                A personal Read-It-Later clone — solo project, you&apos;re the
                only user, no admin needed.
              </li>
              <li>
                A self-hosted feature-flag service for a small team — one
                binary on a tiny VPS.
              </li>
            </ol>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Single is the smallest. Next is its opposite — <strong>triple</strong>:
        the kitchen-sink SaaS shape with web + admin + API.
      </p>
    </>
  )
}
