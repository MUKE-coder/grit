import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Last lesson. With the compatibility matrix from the previous
        lesson in hand, you ship a coordinated release. The order
        matters — backwards from how you think, sometimes.
      </p>

      <h2>The staggered order</h2>
      <Diagram label="Staggered release timeline" caption="API first, then frontends. Mobile last because it&apos;s the slowest to roll out.">
{`   Time ────────────────────────────────────────────────►

   Step 1   ┌──────────┐
            │   API    │  v1.5  →  Postgres migration runs FIRST
            │  (b/c    │
            │ compat)  │
            └──────────┘

   Step 2          ┌──────────┐
                   │   Web    │  Deploy to Vercel. Instant.
                   │  Admin   │
                   └──────────┘

   Step 3                 ┌──────────┐
                          │ Desktop  │  Auto-updater pushes to existing
                          │ release  │  installs. Slow trickle.
                          └──────────┘

   Step 4                       ┌──────────┐
                                │  Mobile  │  Submit to stores. Days of review.
                                │  release │  Days of user updates.
                                └──────────┘`}
      </Diagram>

      <h2>Why API first</h2>
      <p>
        The API is the surface that EVERYTHING depends on. If you
        update mobile first, mobile v1.5 calls a v1.5 endpoint that
        doesn&apos;t exist yet. 500s. So:
      </p>
      <ol>
        <li>Backwards-compatible API change (the matrix says: existing clients still work).</li>
        <li>Database migration runs as part of API deploy.</li>
        <li>Old API endpoints kept alive for clients still on old versions.</li>
      </ol>
      <p>
        The constraint: every API deploy must be backwards-compatible
        with the OLDEST supported client. The matrix tells you what
        that is.
      </p>

      <h2>Why web second</h2>
      <p>
        Web is instant — push to Vercel, every visitor on next refresh
        gets the new version. So you ship API + web on the same day,
        often within minutes. The user experience: load the page →
        new features.
      </p>

      <h2>Why desktop third</h2>
      <p>
        Desktop has the auto-updater (from the Desktop course), so
        it pushes faster than mobile but slower than web. Typical
        rollout:
      </p>
      <ul>
        <li>Day 0: API + web deployed.</li>
        <li>Day 1: cut desktop release. Auto-updater nudges users.</li>
        <li>Day 3: ~50% of active desktop users on new version.</li>
        <li>Day 7: ~85% on new version.</li>
      </ul>
      <p>
        Why not day 0 for desktop? Because if you discover an API
        bug, you only have to roll back the API + web. Pulling
        desktop is more complex (existing installs need to revert).
      </p>

      <h2>Why mobile last</h2>
      <p>
        Mobile is two layers of slow:
      </p>
      <ol>
        <li>App Store / Play Store review (1-7 days).</li>
        <li>User update behavior (1-30 days).</li>
      </ol>
      <p>
        Implication: mobile is shipping NEXT week&apos;s API change
        THIS week, so that by the time your code is in users&apos;
        hands, the API supports it. Mobile builds are always 1-2
        sprints ahead in their assumed API.
      </p>

      <TipBox tone="warning">
        <strong>Expo OTA is the cheat code.</strong> For
        JavaScript-only changes (no new native modules), Expo
        Updates lets you push a JS bundle directly to existing
        installs in minutes — no store review. Use OTA aggressively
        for bug fixes; reserve full builds for native changes. The
        mobile course covers this.
      </TipBox>

      <h2>The rollback plan</h2>
      <p>
        For every release, write a 1-page rollback plan BEFORE you
        deploy:
      </p>
      <CodeBlock
        language="markdown"
        filename="releases/v1.5-rollback.md"
        code={`# v1.5 Rollback plan

## What ships
- API: adds /api/featured-products endpoint
- Web: featured-products carousel on homepage
- Desktop: not affected
- Mobile: not affected

## If broken
1. API rollback: \`fly deploy --image api:v1.4\`  (DB migration is additive — safe to leave)
2. Web rollback: in Vercel dashboard → promote previous deployment
3. Desktop: no action
4. Mobile: no action

## Smoke tests after each step
- API:    curl /healthz returns 200
- Web:    homepage loads
- Test user can log in`}
      />
      <p>
        Boring document. Saves you when it&apos;s 2am.
      </p>

      <h2>Coordinating the people side</h2>
      <p>
        Release coordination is half technical, half social:
      </p>
      <ul>
        <li>
          <strong>Release notes per surface.</strong> Don&apos;t conflate.
          Web users don&apos;t care about your iOS-only fix.
        </li>
        <li>
          <strong>One channel for the whole release.</strong> A Slack
          thread per release where status lands. Easier to scan than
          15 PRs.
        </li>
        <li>
          <strong>Freeze window.</strong> Friday afternoon is not when
          you deploy. Sunday before a holiday is not when you deploy.
        </li>
        <li>
          <strong>Post-deploy checklist.</strong> Smoke-test each
          surface within an hour. Page someone if anything looks off.
        </li>
      </ul>

      <h2>Pulse + Sentinel earn their keep</h2>
      <p>
        After deploy, Pulse (your Grit observability) tells you
        request volumes per surface. If web traffic drops 80% in 5
        min, something&apos;s broken. Sentinel rate-limits keep a
        misbehaving client from taking down the API while you
        investigate.
      </p>

      <KnowledgeCheck
        question="You release v1.5 — API and web on Monday, desktop on Wednesday, mobile next week. Wednesday afternoon a user reports a bug that only repros on the API change. What's the cleanest rollback?"
        choices={[
          {
            label: 'Roll back API + web + desktop together',
            feedback:
              "Web doesn't strictly need rolling back if the API rollback is backwards-compatible. Don't ripple changes you don't have to.",
          },
          {
            label: "If the API change is backwards-compat-safe with v1.4 clients, roll back ONLY the API. Web/desktop on v1.5 still work because they don't use the broken behavior heavily.",
            correct: true,
            feedback:
              "Right — that's the value of additive API changes. Existing clients survive on the older API. Roll back only what's broken; minimise blast radius.",
          },
          {
            label: 'Roll forward with a fix instead',
            feedback:
              "Sometimes — depends on severity and confidence. For most production bugs, roll back THEN write the fix calmly.",
          },
          {
            label: 'Wait until mobile ships to fix everything at once',
            feedback:
              "Letting users hit a known bug while you wait is bad UX and erodes trust.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For chapter 5&apos;s assignment, cut a coordinated release:</p>
            <ol>
              <li>
                Pick a small but real feature (e.g., add{' '}
                <code>is_featured</code> to products + show on each
                surface).
              </li>
              <li>Update <code>docs/compat-matrix.md</code>.</li>
              <li>Write a 1-page <code>releases/vX.Y-rollback.md</code> before deploying.</li>
              <li>Deploy API + web (Vercel-style).</li>
              <li>Wait 24h, then cut desktop with the auto-updater.</li>
              <li>
                Cut a mobile EAS build, OR ship as an OTA update if
                the change is JS-only.
              </li>
              <li>
                Practice the rollback: pretend the API change broke
                something. Walk the rollback steps in your doc. Time
                it. Aim for &lt; 10 minutes.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the rollback drill, use a staging environment. Don&apos;t
            actually break prod. The point is muscle memory.
          </>
        }
        solution={
          <>
            <p>
              By the end of this assignment you&apos;ll have shipped a
              feature across 4 surfaces in a controlled, staggered
              way, with a written rollback plan. That&apos;s
              production-grade release engineering — what separates
              hobby projects from products.
            </p>
          </>
        }
      />

      <h2>You finished Building Web + Desktop + Mobile 🎉</h2>
      <p>
        Five chapters, 13 lessons. You can now:
      </p>
      <ul>
        <li>
          Scaffold a four-surface project: API + web + admin + mobile
          + desktop on one monorepo
        </li>
        <li>Share types and Zod schemas across all surfaces with one command</li>
        <li>
          Build a feature with surface-appropriate UX on each platform,
          backed by one API
        </li>
        <li>Handle offline reads + writes on mobile (persistence) and desktop (outbox)</li>
        <li>Resolve cross-surface conflicts with an explicit policy</li>
        <li>
          Coordinate staggered releases without breaking older clients,
          with a written rollback plan
        </li>
      </ul>

      <h2>Where to go from here</h2>
      <p>
        You&apos;ve completed the full Grit kit course set. Some
        directions:
      </p>
      <ul>
        <li>
          Ship a real product. Honestly — the courses gave you the
          map; shipping closes the loop.
        </li>
        <li>
          Open-source contribution: PRs to{' '}
          <a href="https://github.com/MUKE-coder/grit" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">github.com/MUKE-coder/grit</a>{' '}
          welcome — generators, kits, docs.
        </li>
        <li>
          Build a Grit UI component for the registry — the catalog
          keeps growing.
        </li>
        <li>
          Share your story: tweet, write a post, record a video. The
          community grows from real builders showing real work.
        </li>
      </ul>
      <p>
        Thanks for finishing the course. Go ship something.
      </p>
    </>
  )
}
