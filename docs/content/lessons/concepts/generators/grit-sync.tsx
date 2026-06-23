import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        One command, one purpose: keep the TypeScript types in sync with the
        Go structs. <code>grit sync</code> reads every model and re-writes
        <code> packages/shared/src/types/</code>. Run it any time the API
        shape changes.
      </p>

      <h2>What it does</h2>
      <CodeBlock
        terminal
        code={`grit sync`}
      />

      <p>You&apos;ll see:</p>

      <CodeBlock
        language="text"
        code={`  Syncing Go types → TypeScript...

  ✓ packages/shared/types/contact.ts
  ✓ packages/shared/schemas/contact.ts
  ✓ apps/admin/resources/contacts.ts  (added 1 field to columns + form)
  ✓ packages/shared/types/group.ts
  ✓ packages/shared/schemas/group.ts

  ✅ Synced 2 model(s) to TypeScript + Zod
  ✅ Auto-added 1 field to admin resource files`}
      />

      <TipBox tone="success">
        <strong>v3.31.16:</strong> sync now also walks{' '}
        <code>apps/admin/resources/&lt;plural&gt;.ts</code> and
        auto-appends any new model fields between the{' '}
        <code>// grit:cols:auto-end</code> and{' '}
        <code>// grit:fields:auto-end</code> markers. Customised
        entries are never touched. Resources scaffolded before
        v3.31.16 don&apos;t have the markers — sync prints a
        per-resource warning telling you what to add.
      </TipBox>

      <h2>Why it matters</h2>
      <p>
        Without sync, every time you add a field to a Go struct you also
        have to remember to add it to the TS type. That&apos;s how schema drift
        happens — your API returns <code>is_active</code> but the frontend
        types still say it doesn&apos;t exist, so VS Code says &quot;property does not
        exist&quot; while the data IS being sent.
      </p>
      <p>
        With sync, you change Go once and run one command. TS picks up the
        new field. Build errors point you at every place to update.
      </p>

      <h2>When to run it</h2>
      <ul>
        <li>After every <code>grit generate resource ...</code></li>
        <li>After you manually edit a Go struct</li>
        <li>After you pull from main and someone else changed models</li>
        <li>In CI before <code>pnpm build</code></li>
      </ul>

      <TipBox tone="info">
        Many teams add <code>grit sync</code> to a pre-commit hook or to the
        first step of <code>turbo dev</code>. Make it muscle memory and you&apos;ll
        never ship a drifted type again.
      </TipBox>

      <h2>What it doesn&apos;t do</h2>
      <p>
        Sync handles the shape — the TS type matches the Go struct. It
        doesn&apos;t handle <strong>behaviour</strong> — e.g., custom validation
        rules, derived fields, or computed properties. Those still need
        manual work in the Zod schema.
      </p>

      <KnowledgeCheck
        question="You add a `discount` field to Product in Go. You re-run `grit migrate` but not `grit sync`. What breaks?"
        choices={[
          {
            label: 'The migration fails — the DB column wasn\'t created.',
            feedback:
              'Wrong — `grit migrate` runs GORM AutoMigrate, which created the column.',
          },
          {
            label: 'The admin panel breaks — its form is missing discount.',
            feedback:
              "Half right — the admin will likely build, but it will be missing the field. The deeper problem: the TypeScript types in apps/admin and apps/web no longer match what the API returns.",
          },
          {
            label: "TypeScript build fails because the type doesn't include `discount` but the API returns it.",
            correct: true,
            feedback:
              "Right — your frontend code that does `product.discount` won\'t compile because the TS type doesn\'t have the field. `grit sync` regenerates the type and the build goes green.",
          },
          {
            label: 'Nothing — TS types and Go structs are independent.',
            feedback:
              "Wrong direction — they're supposed to stay aligned. That's the whole point of having a sync command.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Trigger schema drift, then fix it:</p>
            <ol>
              <li>
                Add{' '}
                <code>
                  IsActive bool {`\`gorm:"default:true" json:"is_active"\``}
                </code>{' '}
                to <code>Product</code> in{' '}
                <code>apps/api/internal/models/product.go</code>.
              </li>
              <li>Run <code>grit migrate</code>.</li>
              <li>
                Try to access <code>product.is_active</code> somewhere in{' '}
                <code>apps/web/</code> — TypeScript will complain.
              </li>
              <li>
                Run <code>grit sync</code>. The TS error vanishes.
              </li>
            </ol>
            <p>Paste before/after TS errors in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            The simplest place to trigger the TS error: open{' '}
            <code>apps/web/hooks/use-products.ts</code> and add{' '}
            <code>const _ = data[0].is_active</code> somewhere. Save. TS
            screams. Run sync. TS smiles.
          </>
        }
        solution={
          <>
            <p>Before sync:</p>
            <CodeBlock language="text" code={`Property 'is_active' does not exist on type 'Product'.`} />
            <p>After sync:</p>
            <CodeBlock language="text" code={`(no errors)`} />
            <p>
              That round-trip — change Go, run sync, TypeScript catches up —
              is the loop you&apos;ll do dozens of times per week.
            </p>
          </>
        }
      />

      <h2>You&apos;ve finished chapter 4</h2>
      <p>
        Generate writes 8 files; sync keeps the types aligned. Together they
        let you add a resource to a Grit project in under a minute. The
        chapter assignment — generate an <code>Order</code> + an{' '}
        <code>OrderItem</code> + sync — is waiting in the sidebar.
      </p>

      <h2>What&apos;s next</h2>
      <p>
        Chapter 5 — the last chapter. We zoom out from individual resources
        and talk about the 5 architecture modes and how to pick the right
        one for your idea.
      </p>
    </>
  )
}
