import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Grit&apos;s default migration strategy is GORM <code>AutoMigrate</code>:
        run it on boot, GORM reconciles the schema with your struct
        definitions. Magical when it works. This lesson covers when to keep
        it and when to switch.
      </p>

      <h2>What AutoMigrate does</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/database/database.go (excerpt)"
        code={`db.AutoMigrate(
    &models.User{},
    &models.Customer{},
    &models.Invoice{},
    &models.LineItem{},
)`}
      />
      <p>For each struct, GORM:</p>
      <ul>
        <li>Creates the table if it doesn&apos;t exist</li>
        <li>Adds new columns when you add struct fields</li>
        <li>Adds new indexes when you add new <code>gorm:&quot;index&quot;</code> tags</li>
        <li>Adds new foreign keys when you add relations</li>
      </ul>

      <h2>What AutoMigrate <em>doesn&apos;t</em> do</h2>
      <ul>
        <li>
          <strong>Remove columns</strong> — you can delete a struct field but
          the column lingers. Manual <code>ALTER TABLE DROP COLUMN</code>.
        </li>
        <li>
          <strong>Rename columns</strong> — it would see the rename as &quot;drop
          A, add B&quot; and lose data. Use raw SQL with a manual rename.
        </li>
        <li>
          <strong>Change column types</strong> — only certain narrowings
          work; widening usually does, narrowing rarely does.
        </li>
        <li>
          <strong>Data migrations</strong> — backfilling values when a new
          column lands. AutoMigrate handles schema; YOU handle data.
        </li>
      </ul>

      <TipBox tone="warning">
        AutoMigrate is great for development and for products that don&apos;t
        run on multiple workers simultaneously. For production with multiple
        API replicas, run migrations as a separate, one-time step BEFORE
        boot — not during.
      </TipBox>

      <h2>Running migrations explicitly</h2>
      <p>
        Grit ships <code>grit migrate</code> as a CLI command that calls
        AutoMigrate on the same model list. In production:
      </p>
      <CodeBlock
        terminal
        code={`# in your deploy script, BEFORE starting the new API binary
grit migrate

# then start the API
./api-binary`}
      />
      <p>
        This way the migration only runs once per release (not per worker
        on boot), and a failed migration aborts the deploy before serving
        traffic.
      </p>

      <h2>When to switch to explicit migrations</h2>
      <ul>
        <li>You need to do <strong>data migrations</strong> (backfills, splits, joins)</li>
        <li>You need to <strong>rename columns</strong> without dropping data</li>
        <li>You ship to customers who self-host (they need versioned migrations they can roll back)</li>
        <li>You need a <strong>review trail</strong> — explicit SQL files in git that&apos;s reviewable</li>
      </ul>
      <p>
        The graduation path: add <code>golang-migrate/migrate</code> or{' '}
        <code>pressly/goose</code>, put SQL files in{' '}
        <code>internal/database/migrations/</code>, run them from a CLI. Grit
        doesn&apos;t prescribe one — pick what fits your team.
      </p>

      <h2>The pragmatic middle: AutoMigrate + occasional raw SQL</h2>
      <p>
        Most Grit projects start with AutoMigrate, then run an occasional
        manual <code>ALTER TABLE</code> for renames + drops. That&apos;s fine for
        teams of 1–5 with a clear deploy process.
      </p>

      <KnowledgeCheck
        question="You renamed a field from `Total` to `TotalCents` in your Invoice struct. You ran `grit migrate`. What's in the DB?"
        choices={[
          {
            label: 'Total renamed to TotalCents with data preserved',
            feedback:
              "Wrong — AutoMigrate can't see a rename. It sees 'Total is gone, TotalCents is new'.",
          },
          {
            label: 'Both `total` (with the old data) AND `total_cents` (empty) exist',
            correct: true,
            feedback:
              "Right — AutoMigrate added `total_cents` because the struct field is new. It DIDN'T drop `total` because AutoMigrate never drops. You now have a dead column with the old data.",
          },
          {
            label: 'The migration fails — GORM rejects the rename',
            feedback:
              "Wrong — GORM doesn't try to detect renames. It just adds new columns silently.",
          },
          {
            label: 'Only `total_cents` exists with the data copied over',
            feedback:
              "Wrong — there's no automatic data copy. You'd have to do that with an explicit migration.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Trigger the rename gotcha on your bench-api:
            </p>
            <ol>
              <li>
                In your <code>Invoice</code> model, rename a field — e.g.,{' '}
                <code>Total</code> → <code>TotalAmount</code>.
              </li>
              <li>Run <code>grit migrate</code>.</li>
              <li>
                Open GORM Studio and look at the invoices table. You&apos;ll see
                both columns.
              </li>
              <li>
                In <code>notes.md</code>, write the SQL you&apos;d run to
                consolidate: copy data + drop old column.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            The fix is two SQL statements: <code>UPDATE</code> to copy from
            the old to the new column, then <code>ALTER TABLE ... DROP COLUMN</code>{' '}
            the old one.
          </>
        }
        solution={
          <>
            <CodeBlock
              language="sql"
              code={`UPDATE invoices SET total_amount = total WHERE total_amount IS NULL;
ALTER TABLE invoices DROP COLUMN total;`}
            />
            <p>
              Run via <code>psql</code> or GORM Studio&apos;s Query tab. Future
              boots see one column, with the correct data. This is the kind of
              manual step you graduate from when you switch to explicit
              migrations.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>Auth + RBAC</strong>. JWT, OAuth2, TOTP 2FA, and
        role-based gating. The heaviest chapter, the highest leverage.
      </p>
    </>
  )
}
