import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Before diving into four implementations, pick the right
        feature to build. Not every feature ships well on every
        surface — some are web-first, some are mobile-first, some
        only make sense on desktop. The wrong pick wastes the chapter.
      </p>

      <h2>Our pick — Bookmarks</h2>
      <p>
        Bookmarks: let the user save products to a personal list,
        view them, remove them. It&apos;s our running example for
        the next three lessons. Why?
      </p>
      <ul>
        <li>
          <strong>Simple data model.</strong> One table:{' '}
          <code>bookmarks(id, user_id, product_id, created_at)</code>.
        </li>
        <li>
          <strong>Cross-surface obvious.</strong> Users expect to bookmark
          on phone, see it on web. Multi-platform is the value.
        </li>
        <li>
          <strong>Trivial offline story.</strong> Save bookmark offline
          → push when online. Good for chapter 4.
        </li>
        <li>
          <strong>Shows the shared layer.</strong> One type, one Zod
          schema, three frontends.
        </li>
      </ul>

      <h2>The decision rubric</h2>
      <p>For any feature, ask:</p>
      <ul>
        <li>
          <strong>Where do users START?</strong> If 90% of usage is
          on mobile, build mobile first. Don&apos;t force a desktop UI
          for a feature nobody opens desktop for.
        </li>
        <li>
          <strong>Where do users FINISH?</strong> &quot;Quick-add on
          mobile, manage on web&quot; is a common pattern. Lean into
          it; don&apos;t replicate the full UX on every surface.
        </li>
        <li>
          <strong>What HARDWARE matters?</strong> Camera → mobile-first.
          Keyboard shortcuts → desktop/web. Push notifications →
          mobile + web; not desktop (usually).
        </li>
        <li>
          <strong>Can it survive offline?</strong> If yes, mobile +
          desktop need a sync story. If no, fail fast with a clear
          error.
        </li>
      </ul>

      <h2>Examples by category</h2>
      <ul>
        <li>
          <strong>Web-first / mobile read-only:</strong> Admin panels,
          billing, settings. Power-user features. Heavy form work.
        </li>
        <li>
          <strong>Mobile-first / web view-only:</strong> Photo capture,
          location check-in, &quot;quick action&quot; flows.
        </li>
        <li>
          <strong>Desktop-first:</strong> File I/O at scale (batch import
          / export), offline-tolerant tools, long-running tasks (PDF
          generation), shortcuts-heavy workflows.
        </li>
        <li>
          <strong>Equal across all:</strong> Bookmarks, comments,
          notifications, profile. Light data, mostly read.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Don&apos;t force feature parity.</strong> The dangerous
        trap is &quot;every feature on every surface&quot;. That
        triples cost without tripling value. Each surface should have
        its strengths; the API just needs to support all of them.
      </TipBox>

      <h2>What we&apos;ll build over the next 3 lessons</h2>
      <p>
        Bookmarks, end to end:
      </p>
      <ul>
        <li>
          <strong>API</strong> — POST <code>/api/bookmarks</code>, GET{' '}
          <code>/api/bookmarks</code>, DELETE <code>/api/bookmarks/:id</code>.
          (We&apos;ll touch this briefly — the focus is the four frontends.)
        </li>
        <li>
          <strong>Web</strong> — bookmark button on a product page,
          /bookmarks list page.
        </li>
        <li>
          <strong>Mobile</strong> — heart icon on product card, dedicated
          Bookmarks tab.
        </li>
        <li>
          <strong>Desktop</strong> — keyboard shortcut (Cmd/Ctrl+D),
          bookmarks sidebar.
        </li>
      </ul>
      <p>
        Same backend, three different presentations playing to each
        surface&apos;s strengths.
      </p>

      <h2>Quick API scaffold (we&apos;ll trust this and move on)</h2>
      <p>
        Use the resource generator so we can focus on the frontends:
      </p>
      <ul>
        <li>
          Run <code>grit generate resource Bookmark user_id:uint product_id:uint</code>
        </li>
        <li>That gives us model + handler + routes + Zod schemas + TS types.</li>
        <li>
          Add an auth gate so users can only see their own bookmarks
          (we&apos;ll show the snippet next lesson).
        </li>
      </ul>

      <KnowledgeCheck
        question="A teammate proposes adding a 12-step KYC form to the mobile app because 'we should have feature parity'. The same form exists on web. What's the recommended approach?"
        choices={[
          {
            label: 'Build it — feature parity is a customer expectation',
            feedback:
              'Not always. 12-step forms are painful on mobile, and most users complete them once. Pushing them to the web for that step is fine if the data syncs.',
          },
          {
            label: 'Build a mobile-friendly "continue on web" flow — send a deep link, let the user finish on a bigger screen',
            correct: true,
            feedback:
              'Right — pragmatism beats parity. The data is the same; the UX adapts to the surface. Mobile is the entry, web is the workspace.',
          },
          {
            label: 'Build a stripped-down mobile version with fewer fields',
            feedback:
              'Risky — if the trimmed version doesn’t collect required fields, you’re stuck. The “continue on web” approach is cleaner.',
          },
          {
            label: 'Reject the feature on mobile entirely',
            feedback:
              "Too extreme — users WILL hit KYC on mobile if they start there. You just shouldn’t force the whole form there.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Set up for the next three lessons:</p>
            <ol>
              <li>
                Run <code>grit generate resource Bookmark user_id:uint product_id:uint</code>
              </li>
              <li>
                Edit the generated handler so the LIST endpoint only
                returns bookmarks for the authed user (filter by{' '}
                <code>user_id = c.GetUint(&quot;user_id&quot;)</code>).
              </li>
              <li>
                Add a unique index on (user_id, product_id) so a
                user can&apos;t bookmark the same product twice.
              </li>
              <li>
                Run <code>grit sync</code> to update shared types.
              </li>
              <li>
                Use Postman / curl to POST a bookmark and GET the list.
                Confirm the auth gate works.
              </li>
            </ol>
            <p>
              Once that&apos;s green, you&apos;re ready for the
              implementation lessons.
            </p>
          </>
        }
        hint={
          <>
            For the unique index, add{' '}
            <code>uniqueIndex:idx_user_product</code> to the GORM tags
            on both fields. GORM creates a composite index.
          </>
        }
        solution={
          <>
            <p>
              You should now be able to POST{' '}
              <code>/api/bookmarks &#123;product_id: 1&#125;</code> as
              two different users and each sees only their own list.
              The same POST twice returns a 409 (or 422) thanks to the
              unique index.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Web implementation</strong>. Bookmark
        button, list page, optimistic UI updates with React Query.
      </p>
    </>
  )
}
