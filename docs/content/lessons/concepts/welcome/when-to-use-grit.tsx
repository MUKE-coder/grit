import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Every tool is right for some jobs and wrong for others. Grit&apos;s no
        exception. By the end of this lesson you&apos;ll be able to look at a
        project idea and say <em>yes-Grit</em> or <em>no-Grit</em> with
        confidence — and know <em>why</em>.
      </p>

      <h2>Grit shines when</h2>
      <ul>
        <li>
          <strong>You want to ship a SaaS, dashboard, or admin-driven
          product.</strong> The triple kit (web + admin + API) is purpose-built
          for this shape.
        </li>
        <li>
          <strong>You like Go for the backend.</strong> Performance, single
          binary, no runtime — and you want type safety from DB → API →
          frontend without manually keeping types in sync.
        </li>
        <li>
          <strong>You&apos;re a small team (1–5 people).</strong> The conventions
          mean every team member is productive in any part of the code; no
          tribal knowledge.
        </li>
        <li>
          <strong>You want all surfaces from one codebase.</strong> Web,
          desktop (Wails), and mobile (Expo) — all backed by the same Go API
          with shared types.
        </li>
        <li>
          <strong>You don&apos;t want to wire-up auth, jobs, mail, storage, AI
          from scratch.</strong> They ship pre-wired.
        </li>
      </ul>

      <h2>Grit is overkill (or wrong) when</h2>
      <ul>
        <li>
          <strong>It&apos;s a single static landing page.</strong> Use Astro,
          plain HTML, or a Next.js single-page if SEO matters.
        </li>
        <li>
          <strong>You need Python or Node on the backend</strong> for a
          specific library (e.g., heavy ML, a particular SDK with no Go port).
        </li>
        <li>
          <strong>You&apos;re building a CMS-driven content site</strong> —
          Sanity / Payload / WordPress will be a faster fit.
        </li>
        <li>
          <strong>You&apos;re shipping a CLI or a daemon</strong> — write straight
          Go; the meta-framework adds nothing.
        </li>
      </ul>

      <h2>Honest comparison to common alternatives</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm border-collapse">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">If you&apos;d normally use</th>
              <th className="text-left px-3 py-2 font-medium">Reach for Grit when…</th>
              <th className="text-left px-3 py-2 font-medium">Stick with it when…</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2.5 font-mono text-[12px] align-top">Next.js (plain)</td>
              <td className="px-3 py-2.5 align-top">You want a separate Go API + a real admin panel without re-implementing both.</td>
              <td className="px-3 py-2.5 align-top">Pure marketing/blog/SSR-only site with no API surface.</td>
            </tr>
            <tr>
              <td className="px-3 py-2.5 font-mono text-[12px] align-top">Rails / Laravel</td>
              <td className="px-3 py-2.5 align-top">You want Go performance + a typed React frontend instead of server-rendered views.</td>
              <td className="px-3 py-2.5 align-top">You&apos;ve already invested deeply in the ecosystem and have a working team.</td>
            </tr>
            <tr>
              <td className="px-3 py-2.5 font-mono text-[12px] align-top">Supabase / Firebase</td>
              <td className="px-3 py-2.5 align-top">You want to own the backend code (auth, DB, business logic) without vendor lock-in.</td>
              <td className="px-3 py-2.5 align-top">You&apos;re fine with the BaaS pricing and want zero backend code.</td>
            </tr>
            <tr>
              <td className="px-3 py-2.5 font-mono text-[12px] align-top">Plain Gin + GORM</td>
              <td className="px-3 py-2.5 align-top">You&apos;d otherwise spend a week wiring auth/jobs/mail/admin yourself.</td>
              <td className="px-3 py-2.5 align-top">You want full control over every dependency.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <TipBox tone="success">
        <strong>Rule of thumb:</strong> if you can describe your project as
        &quot;X for Y&quot; (the Slack for restaurants, the Notion for warehouses, the
        Stripe Dashboard for vets), Grit is probably right. SaaS-shaped
        products are where it pays off hardest.
      </TipBox>

      <KnowledgeCheck
        question="Which of these projects is the WORST fit for Grit?"
        choices={[
          {
            label: 'A multi-tenant invoicing SaaS with staff admin and customer portal',
            feedback:
              'Strong fit — that\'s exactly the triple kit\'s shape (web + admin + API).',
          },
          {
            label: 'A field-service POS that needs to work offline',
            feedback:
              'Strong fit — the desktop kit ships with offline-first SQLite + sync engine for exactly this.',
          },
          {
            label: 'A blog about your dog with three posts',
            correct: true,
            feedback:
              "Right — Grit is overkill for a static blog. Reach for Astro, a Next.js single-page, or even plain HTML. Grit shines when you have an API, an admin, or multi-tenancy.",
          },
          {
            label: 'A motorcycle dealership management system with loans and repayments',
            feedback:
              "Strong fit — multi-tenant admin + customer portal + complex domain model. This is exactly the shape the Grit demo (Grit Motors) is built around.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Pick three real product ideas (yours, your friend&apos;s, or pulled
              from <a href="https://news.ycombinator.com/show" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">HN Show</a>).
              For each, write one paragraph in <code>notes.md</code>:
            </p>
            <ul>
              <li>What is it (one sentence)</li>
              <li>Would Grit be a good fit? (yes / no)</li>
              <li>Why (one sentence — the deciding factor)</li>
            </ul>
            <p>Try to find at least one example where Grit is the wrong fit.</p>
          </>
        }
        hint={
          <>
            For each idea, ask: <em>is there an admin panel? Multi-tenant? Need
            a typed API surface for multiple clients? Long-running web app
            rather than a single page?</em> Yes to most of those → Grit fits.
          </>
        }
        solution={
          <>
            <p>Example answers:</p>
            <ol>
              <li>
                <strong>Recipe-sharing app</strong> — social, user-generated
                content, search, comments. <em>Yes-Grit</em>: triple kit. The
                admin panel runs moderation, the API serves the public site +
                a future mobile app.
              </li>
              <li>
                <strong>Personal portfolio with two project pages</strong> —{' '}
                <em>No-Grit</em>: static, no admin needed, no backend logic.
                Use a static site generator.
              </li>
              <li>
                <strong>A Discord bot that watches a Telegram channel</strong>{' '}
                — <em>No-Grit</em>: no frontend, no admin, just a daemon. Write
                straight Go (no scaffolder helping here).
              </li>
            </ol>
            <p>
              If at least one of your three said &quot;no&quot; with a specific reason,
              you&apos;ve internalised the trade-off.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You can pitch Grit and you can spot when it fits. The next two lessons
        get hands-on: install Grit on your machine and verify the install
        works.
      </p>
    </>
  )
}
