import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import Link from 'next/link'

export default function Lesson() {
  return (
    <>
      <p>
        You&apos;ve seen every mode. Now the framework for picking — a decision
        tree you can run any new idea through and land on the right command
        without memorizing the table.
      </p>

      <h2>The decision tree</h2>
      <CodeBlock
        language="text"
        code={`Will your product run offline (no internet)?
├── YES → grit new-desktop my-app   [Wails desktop]
└── NO  →
    Does it have a frontend?
    ├── NO  → grit new my-app --api
    └── YES →
        Is the frontend mobile?
        ├── YES  → grit new my-app --mobile
        └── NO  (web frontend) →
            Do you need a separate admin panel for staff?
            ├── YES → grit new my-app --triple
            └── NO  →
                Single binary deploy?
                ├── YES → grit new my-app --single
                └── NO  → grit new my-app --double`}
      />

      <h2>Three quick rules of thumb</h2>
      <ul>
        <li>
          <strong>Default to triple.</strong> When in doubt, use{' '}
          <code>--triple</code>. You can ignore the admin app if you don&apos;t
          need it today, and it&apos;s already wired when you do.
        </li>
        <li>
          <strong>Skip the admin only when you&apos;re sure.</strong> Most products
          eventually need one — support staff, ops, finance. Pre-wiring it
          beats retrofitting.
        </li>
        <li>
          <strong>Mobile + desktop are surfaces, not modes.</strong> Add{' '}
          <code>--mobile</code> or <code>--desktop</code> to triple if you
          need both. Don&apos;t treat them as exclusive picks.
        </li>
      </ul>

      <h2>The Stack Selector tool</h2>
      <p>
        Don&apos;t want to memorize the tree? Use the{' '}
        <Link href="/docs/stack-selector" className="text-primary hover:underline">
          Stack Selector
        </Link>{' '}
        on the docs site. Answer 4 questions, get the exact command. Same
        logic as the decision tree above, in a UI.
      </p>

      <TipBox tone="success">
        <strong>Pro tip:</strong> Once you have a Grit project, you can
        change modes by re-scaffolding — you just don&apos;t want to. Pick
        consciously the first time.
      </TipBox>

      <h2>Real-world examples</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Product</th>
              <th className="text-left px-3 py-2 font-medium">Command</th>
              <th className="text-left px-3 py-2 font-medium">Why</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2">Linear-like task tracker</td>
              <td className="font-mono text-[12px]">--triple</td>
              <td className="text-[12px]">SaaS with admin (billing, support)</td>
            </tr>
            <tr>
              <td className="px-3 py-2">Self-hosted Plausible alternative</td>
              <td className="font-mono text-[12px]">--single</td>
              <td className="text-[12px]">One binary, customer-managed</td>
            </tr>
            <tr>
              <td className="px-3 py-2">Coffee-shop offline POS</td>
              <td className="font-mono text-[12px]">new-desktop</td>
              <td className="text-[12px]">No internet during rush</td>
            </tr>
            <tr>
              <td className="px-3 py-2">Stripe-like payments backend</td>
              <td className="font-mono text-[12px]">--api</td>
              <td className="text-[12px]">SDKs in many languages, no web UI from us</td>
            </tr>
            <tr>
              <td className="px-3 py-2">Habit-tracking app</td>
              <td className="font-mono text-[12px]">--mobile</td>
              <td className="text-[12px]">Phone-first, no web companion</td>
            </tr>
            <tr>
              <td className="px-3 py-2">Boda-boda dealership (Grit demo)</td>
              <td className="font-mono text-[12px]">--single --vite</td>
              <td className="text-[12px]">One binary deploy, dashboard-only UI</td>
            </tr>
          </tbody>
        </table>
      </div>

      <KnowledgeCheck
        question="A friend is starting an event-ticketing platform. They want a public site to buy tickets, an admin for the venue owner to manage events, and they want eventually to ship a mobile app for scanning tickets at the door. What do you tell them?"
        choices={[
          {
            label: 'Start with --single and refactor later',
            feedback:
              "Wrong — single doesn't have an admin app, and they explicitly want one. They'd be re-scaffolding within a month.",
          },
          {
            label: 'Start with --triple now; add --mobile when the mobile scanner work begins',
            correct: true,
            feedback:
              "Right — triple covers the public site + admin + API immediately. They can scaffold --triple --mobile in one shot if mobile is near-term, OR add the mobile slice later by generating the apps/mobile structure.",
          },
          {
            label: 'Three separate Grit projects: one for web, one for admin, one for mobile',
            feedback:
              "Wrong — they\'d have three sets of types to keep aligned manually. The monorepo + shared package exists exactly to prevent this.",
          },
          {
            label: 'Start with --api and let them build the UI in any framework they want',
            feedback:
              "Possible, but punts the question. They have a clear UI vision (web + admin); --triple gives them both pre-wired.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              The chapter assignment is waiting — &quot;Pick the right kit for
              three ideas.&quot; Open it from the sidebar and tackle it now.
              Pick three real product ideas (yours, your friend&apos;s, anything
              from HN), trace each through the decision tree, and write
              one paragraph justifying your pick.
            </p>
            <p>
              The trick: try to pick three DIFFERENT modes. That forces you
              to find ideas that genuinely justify different shapes.
            </p>
          </>
        }
        hint={
          <>
            If you can&apos;t find three different modes, the constraint is
            telling you something: you mostly think in triple-shaped
            products. That&apos;s fine — most SaaS is triple. But challenge
            yourself.
          </>
        }
        solution={
          <>
            <p>
              No solution here — the assignment page on the sidebar
              evaluates your three picks. Common pitfalls:
            </p>
            <ul>
              <li>Picking triple for everything (over-engineering)</li>
              <li>Picking single for something that needs an admin (under-engineering)</li>
              <li>Picking desktop for a product that doesn&apos;t need offline (over-spec)</li>
            </ul>
            <p>
              If your three are <em>genuinely different</em> shapes, you&apos;ve
              internalised the kit choices. That&apos;s the goal.
            </p>
          </>
        }
      />

      <h2>You&apos;ve finished Grit Concepts 🎉</h2>
      <p>
        Five chapters, ~20 lessons, five assignments. You can now:
      </p>
      <ul>
        <li>Explain what Grit is in 30 seconds</li>
        <li>Install + verify the CLI and keep it updated</li>
        <li>Scaffold a project, navigate every folder, run the dev servers</li>
        <li>Recognise the naming conventions and API response shapes</li>
        <li>Generate a resource end-to-end and keep types in sync</li>
        <li>Pick the right architecture mode for any product idea</li>
      </ul>

      <p>
        Next, pick one of the specialist courses based on the kit you want
        to dive deeper into:{' '}
        <Link href="/courses/go-api" className="text-primary hover:underline">
          Building a Go API
        </Link>
        ,{' '}
        <Link href="/courses/web-nextjs" className="text-primary hover:underline">
          Web with Next.js
        </Link>
        ,{' '}
        <Link href="/courses/mobile" className="text-primary hover:underline">
          Mobile
        </Link>
        , or{' '}
        <Link href="/courses/desktop" className="text-primary hover:underline">
          Desktop
        </Link>
        . All assume the foundation you just finished.
      </p>
    </>
  )
}
