import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Three modes for when you want one slice instead of the full triple:
        API-only, Mobile, and Desktop. Each has its own course in the
        Learning Paths section.
      </p>

      <h2>API-only</h2>
      <CodeBlock terminal code={`grit new my-api --api`} />
      <p>
        Pure Gin + GORM Go backend. No frontend. Auto-served OpenAPI docs
        at <code>/docs</code>. Ships every battery (auth, jobs, mail,
        storage, AI, Sentinel + Pulse) — just no UI.
      </p>

      <p><strong>When to use:</strong></p>
      <ul>
        <li>Backend for a mobile-only product (bring your own iOS / Android client)</li>
        <li>Backend for a Discord / Slack / Telegram bot</li>
        <li>Backend for a native desktop app you&apos;re building separately</li>
        <li>Microservice in a larger system (a Grit project consumed by other services)</li>
      </ul>

      <h2>Mobile (Expo + API)</h2>
      <CodeBlock terminal code={`grit new my-app --mobile`} />
      <p>
        A monorepo with <code>apps/api</code> (Go) and{' '}
        <code>apps/mobile</code> (Expo). Shared types via{' '}
        <code>grit sync</code>. Mobile-friendly auth (refresh tokens in
        SecureStore), EAS Build configured, push-notification scaffolding.
      </p>

      <p><strong>When to use:</strong></p>
      <ul>
        <li>Mobile-first product where the web is secondary or absent</li>
        <li>Field-service apps (drivers, technicians, sales reps)</li>
        <li>Consumer apps for iOS + Android</li>
      </ul>

      <h2>Desktop (Wails + GORM)</h2>
      <CodeBlock terminal code={`grit new-desktop my-app`} />
      <p>
        Standalone Wails v2 app — Go bound to React via Wails. Local SQLite
        (offline-first), local auth, in-app auto-updater, NSIS installers
        (full + slim), branded MUI. Different command (<code>new-desktop</code>{' '}
        instead of <code>new</code>) because it&apos;s structurally different
        from the rest.
      </p>

      <p><strong>When to use:</strong></p>
      <ul>
        <li>POS systems / kiosks that need to work offline</li>
        <li>Field service tools that go to areas with no internet</li>
        <li>Single-user apps where the data should stay on the user&apos;s machine</li>
        <li>Tools that need access to OS-level features (filesystem, native dialogs)</li>
      </ul>

      <TipBox tone="success">
        <strong>Adding a slice to triple:</strong> flags compose.{' '}
        <code>grit new my-app --triple --desktop</code> gives you the SaaS
        shape + a Wails desktop client. <code>--triple --mobile</code> adds
        an Expo app to the triple monorepo. The seventh course
        (&quot;Multi-Platform&quot;) covers exactly this.
      </TipBox>

      <h2>Mode at a glance</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Mode</th>
              <th className="text-left px-3 py-2 font-medium">Output</th>
              <th className="text-left px-3 py-2 font-medium">Best for</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 font-mono">--single</td><td className="text-[12px]">One Go binary + embedded React</td><td className="text-[12px]">Self-hosted, indie, MVP</td></tr>
            <tr><td className="px-3 py-2 font-mono">--single --vite</td><td className="text-[12px]">One Go binary + embedded Vite SPA</td><td className="text-[12px]">Internal dashboards, fast cold starts</td></tr>
            <tr><td className="px-3 py-2 font-mono">--double</td><td className="text-[12px]">apps/web + apps/api</td><td className="text-[12px]">Public product, no admin needed</td></tr>
            <tr><td className="px-3 py-2 font-mono">--triple</td><td className="text-[12px]">web + admin + api</td><td className="text-[12px]">SaaS (the default recommendation)</td></tr>
            <tr><td className="px-3 py-2 font-mono">--api</td><td className="text-[12px]">apps/api only</td><td className="text-[12px]">Bring-your-own-client backends</td></tr>
            <tr><td className="px-3 py-2 font-mono">--mobile</td><td className="text-[12px]">api + Expo mobile</td><td className="text-[12px]">Mobile-first products</td></tr>
            <tr><td className="px-3 py-2 font-mono">new-desktop</td><td className="text-[12px]">Wails app (standalone)</td><td className="text-[12px]">Offline-first / native</td></tr>
          </tbody>
        </table>
      </div>

      <KnowledgeCheck
        question="You're building a delivery driver app. The drivers use phones, dispatch staff uses a web admin, and there's an API for both. Which command?"
        choices={[
          {
            label: 'grit new my-app --triple --mobile',
            correct: true,
            feedback:
              'Right — triple (web + admin + API) + mobile (Expo for drivers) = web admin for dispatch, mobile app for drivers, one API serving both.',
          },
          {
            label: 'grit new my-app --mobile',
            feedback:
              "Wrong — that's mobile + API only. You'd be missing the web admin for dispatch staff.",
          },
          {
            label: 'grit new my-app --single',
            feedback:
              "Wrong — no mobile client, and a single binary isn't the right shape for a driver app on phones.",
          },
          {
            label: 'grit new-desktop my-app',
            feedback:
              "Wrong — desktop doesn't fit a phones-for-drivers product.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              For each of these ideas, write down the right command in{' '}
              <code>notes.md</code>:
            </p>
            <ol>
              <li>
                Open-source feature-flag service (one binary, web UI to
                toggle flags, no admin staff)
              </li>
              <li>
                Restaurant POS that runs on a tablet without internet
              </li>
              <li>
                Discord bot that watches GitHub repos for new issues
              </li>
              <li>
                Fitness tracker with iPhone + Android app
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For each, ask: web surface? admin staff? mobile? offline? Match
            answer to mode.
          </>
        }
        solution={
          <>
            <ol>
              <li>
                <code>grit new flagd --single</code> — one binary, no
                staff/admin, web UI embedded.
              </li>
              <li>
                <code>grit new-desktop restaurant-pos</code> — Wails, local
                SQLite, offline-first.
              </li>
              <li>
                <code>grit new gh-watcher --api</code> — no UI, just a
                backend the bot calls.
              </li>
              <li>
                <code>grit new fitness-app --mobile</code> — Expo mobile +
                Go API.
              </li>
            </ol>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final lesson — a decision tree to pick the right kit for any idea
        without having to memorize the table. After that, the chapter
        assignment is &quot;match 3 of your own product ideas to the right
        kits.&quot;
      </p>
    </>
  )
}
