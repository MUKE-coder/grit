import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Zero-downtime Postgres migrations with GORM',
  description:
    'Expand and contract, the DDL that takes a lock and blocks every read, CREATE INDEX CONCURRENTLY, and why AutoMigrate stops at development.',
  alternates: { canonical: 'https://gritframework.dev/guides/gorm-zero-downtime-migrations' },
  openGraph: {
    title: 'Zero-downtime Postgres migrations with GORM',
    description: 'Which statements are safe, which take a lock, and the pattern that avoids both.',
    url: 'https://gritframework.dev/guides/gorm-zero-downtime-migrations',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="gorm-zero-downtime-migrations"
      lead={
        <>
          <p className="mb-4">
            A migration takes down a production service in one of two ways: it holds a lock that
            blocks every query against a table, or it changes a shape the currently running code
            still expects. Both are avoidable, and the techniques are not complicated once you
            know which statements are dangerous.
          </p>
          <p>
            This covers which Postgres operations are safe, which are not and what to do instead,
            and the expand-and-contract sequence that makes a schema change survive a rolling
            deploy.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A list of the DDL that is safe to run on a live table, and the list that is not.</p>
        <p>
          A migration sequence that works while two versions of your application are running at
          once, which is what a rolling deploy means.
        </p>
        <p>A rollback plan that is not &quot;restore the backup&quot;.</p>
      </Outcome>

      <Section id="automigrate" title="Where AutoMigrate stops">
        <p>
          <code>db.AutoMigrate(&amp;User{})</code> compares structs to the database and adds what
          is missing. It is excellent for development and unusable in production, for four
          specific reasons rather than as a matter of taste.
        </p>
        <Checklist
          items={[
            <>
              <strong>It never drops or renames.</strong> Remove a field and the column stays.
              Rename one and you get a new column beside the old one, with the data in the old
              one.
            </>,
            <>
              <strong>It cannot be reviewed.</strong> There is no artefact in the pull request.
              Nobody can see that this deploy adds a NOT NULL column to a table with forty million
              rows until it is running.
            </>,
            <>
              <strong>There is no down.</strong> No rollback, no record of what ran, no way to
              reproduce the schema at a given release.
            </>,
            <>
              <strong>It takes whatever lock the generated statement needs.</strong> You did not
              choose the statement, so you cannot know what it will hold.
            </>,
          ]}
        />
        <p>
          So: <code>AutoMigrate</code> locally for speed, versioned SQL files everywhere else. The
          structs remain the source of truth for the Go side; the migration files are the record
          of how the database got here.
        </p>
      </Section>

      <Section id="locks" title="Which statements are safe">
        <p>
          The question is always which lock a statement takes and for how long. An{' '}
          <code>ACCESS EXCLUSIVE</code> lock blocks every read and every write on that table,
          including plain <code>SELECT</code>. Taking one for ten milliseconds is invisible.
          Taking one for four minutes is an outage.
        </p>
        <H3>Safe: instant, regardless of table size</H3>
        <CodeBlock
          language="sql"
          code={`ALTER TABLE invoices ADD COLUMN note text;                  -- nullable, no default
ALTER TABLE invoices ADD COLUMN status text DEFAULT 'draft'; -- PG 11+: metadata only
ALTER TABLE invoices DROP COLUMN note;                       -- metadata only
ALTER TABLE invoices RENAME COLUMN a TO b;                   -- instant, and see below
ALTER TABLE invoices ALTER COLUMN note DROP NOT NULL;
CREATE INDEX CONCURRENTLY idx_invoices_status ON invoices (status);`}
        />
        <p>
          Adding a column with a default stopped rewriting the table in Postgres 11. On an older
          version it rewrites every row while holding the lock, which is the single most common
          cause of a migration outage in the wild.
        </p>
        <H3>Dangerous: a full table scan or rewrite, under a lock</H3>
        <CodeBlock
          language="sql"
          code={`CREATE INDEX idx ON invoices (status);        -- blocks writes for the build
ALTER TABLE invoices ADD COLUMN x int NOT NULL;  -- must check every row
ALTER TABLE invoices ALTER COLUMN total TYPE bigint;  -- rewrites the table
ALTER TABLE invoices ADD CONSTRAINT ... CHECK (...);  -- validates every row
ALTER TABLE invoices ADD FOREIGN KEY ...;             -- validates every row
ALTER TABLE invoices SET NOT NULL;                    -- scans every row`}
        />
        <H3>The safe equivalents</H3>
        <CodeBlock
          language="sql"
          code={`-- Index: build it without blocking writes. It cannot run inside a
-- transaction, which matters for your migration runner.
CREATE INDEX CONCURRENTLY idx_invoices_status ON invoices (status);

-- A check constraint, in two steps: add it unvalidated, which takes a
-- brief lock, then validate it, which takes a weak one.
ALTER TABLE invoices ADD CONSTRAINT total_positive
  CHECK (total >= 0) NOT VALID;
ALTER TABLE invoices VALIDATE CONSTRAINT total_positive;

-- NOT NULL on Postgres 12+: a validated CHECK lets the scan be skipped.
ALTER TABLE invoices ADD CONSTRAINT x_not_null CHECK (x IS NOT NULL) NOT VALID;
ALTER TABLE invoices VALIDATE CONSTRAINT x_not_null;
ALTER TABLE invoices ALTER COLUMN x SET NOT NULL;   -- now instant
ALTER TABLE invoices DROP CONSTRAINT x_not_null;`}
        />
        <Warn title="A lock queues behind, and in front of, everything">
          <p>
            A migration waiting for an <code>ACCESS EXCLUSIVE</code> lock does not wait quietly.
            Every query that arrives after it queues behind it, including reads, because the lock
            request is ahead of them. One long-running <code>SELECT</code> can therefore cause a
            complete outage on a table: the migration waits for the select, and everything else
            waits for the migration.
          </p>
          <p>Always bound the wait, so a stuck migration fails fast instead of taking the site:</p>
          <CodeBlock
            language="sql"
            code={`SET lock_timeout = '3s';
SET statement_timeout = '30s';`}
          />
        </Warn>
      </Section>

      <Section id="expand" title="Expand and contract">
        <p>
          During a rolling deploy, the old and new versions of your application run at the same
          time against the same database. Any migration that both versions cannot live with is an
          outage, however fast it runs.
        </p>
        <p>The sequence is four deploys, and the discipline is that each one is reversible.</p>
        <CodeBlock
          language="text"
          code={`Goal: rename invoices.total (int) to total_cents (bigint).

1. EXPAND     add total_cents, nullable. Old code ignores it.
2. BACKFILL   copy in batches. Both columns exist; both versions work.
   DUAL WRITE new code writes both, reads total. Deployable, reversible.
3. SWITCH     new code reads total_cents, still writes both.
4. CONTRACT   stop writing total. A later deploy drops it.`}
        />
        <H3>1. Expand</H3>
        <CodeBlock
          language="sql"
          filename="migrations/0042_add_total_cents.up.sql"
          code={`ALTER TABLE invoices ADD COLUMN total_cents bigint;`}
        />
        <H3>2. Backfill, in batches</H3>
        <CodeBlock
          code={`// Never a single UPDATE over the whole table: it holds row locks on
// every row it has touched until it commits, and a failure at 90% rolls
// back the entire thing.
func backfillTotalCents(ctx context.Context, db *gorm.DB) error {
	for {
		res := db.WithContext(ctx).Exec(\`
			UPDATE invoices SET total_cents = total::bigint
			WHERE id IN (
				SELECT id FROM invoices
				WHERE total_cents IS NULL
				ORDER BY id
				LIMIT 5000
				FOR UPDATE SKIP LOCKED
			)\`)
		if res.Error != nil {
			return fmt.Errorf("backfill batch: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return nil
		}
		// Breathe. Replication lag is the thing that bites here: a tight
		// loop can push a read replica minutes behind.
		time.Sleep(100 * time.Millisecond)
	}
}`}
        />
        <p>
          <code>SKIP LOCKED</code> means the backfill steps around rows your application is
          currently updating rather than waiting for them, so the two do not fight. Restartable by
          construction: it selects only rows that are still null.
        </p>
        <H3>3. Dual write, then switch</H3>
        <CodeBlock
          code={`// Deploy A: write both, read the old. Safe to roll back, because the
// old column is still current.
inv.Total = cents
inv.TotalCents = int64(cents)

// Deploy B, after the backfill has finished: read the new, keep writing
// both. Still safe to roll back.

// Deploy C: stop writing the old one.
// Only now is dropping the column safe, and it is a separate migration.`}
        />
        <Note title="Why not just rename it">
          <p>
            <code>ALTER TABLE ... RENAME COLUMN</code> is instant and is still an outage during a
            rolling deploy: the moment it commits, every running instance of the old version is
            querying a column that no longer exists. The rename is cheap for the database and
            catastrophic for the application. Expand and contract exists for the application, not
            for the lock.
          </p>
        </Note>
      </Section>

      <Section id="runner" title="The migration runner">
        <p>Four properties matter, and most hand-rolled runners miss at least one.</p>
        <Checklist
          items={[
            <>
              <strong>A lock, so two instances cannot migrate at once.</strong> Rolling deploys
              start several pods simultaneously. An advisory lock in Postgres is the simplest
              correct answer.
            </>,
            <>
              <strong>A record of what has run,</strong> with a checksum, so an edited migration
              that already ran is an error rather than a silent divergence.
            </>,
            <>
              <strong>Each migration in its own transaction,</strong> with an escape hatch for the
              statements that cannot be, which is <code>CREATE INDEX CONCURRENTLY</code> and a few
              others.
            </>,
            <>
              <strong>A timeout,</strong> so a blocked migration fails rather than queueing the
              whole application behind it.
            </>,
          ]}
        />
        <CodeBlock
          code={`// One instance migrates; the others wait and then find nothing to do.
// 4242 is an arbitrary constant, shared by every instance of this app.
if err := db.Exec("SELECT pg_advisory_lock(4242)").Error; err != nil {
	return fmt.Errorf("acquire migration lock: %w", err)
}
defer db.Exec("SELECT pg_advisory_unlock(4242)")`}
        />
        <Warn title="Migrating on application boot is a trap">
          <p>
            It is convenient and it couples deployment to schema change: a migration that takes
            four minutes delays every pod&apos;s readiness, and a failed migration crash-loops
            the application rather than failing one job. Run migrations as a separate step that
            must succeed before the new version is rolled out.
          </p>
        </Warn>
      </Section>

      <Section id="rollback" title="Rolling back">
        <p>
          Down migrations are worth writing and are not a plan. A down that drops a column you
          spent a deploy backfilling destroys data, and you will not run it in the incident, so
          it provides comfort rather than safety.
        </p>
        <p>The honest position:</p>
        <Checklist
          items={[
            <>
              <strong>Write the down anyway,</strong> for local development and for the cases
              where it is genuinely reversible, such as dropping an index you just created.
            </>,
            <>
              <strong>Make every forward migration safe for the previous version of the code.</strong>{' '}
              That is what the expand-and-contract shape buys: the rollback is a code rollback,
              which is fast and tested, rather than a schema rollback, which is neither.
            </>,
            <>
              <strong>Never combine a destructive change with a release you might reverse.</strong>{' '}
              Drop the column a week later, in its own deploy, when the previous version is no
              longer something you would roll back to.
            </>,
            <>
              <strong>Take a snapshot before a destructive migration,</strong> and know how long a
              restore actually takes. An untested restore is an assumption.
            </>,
          ]}
        />
      </Section>

      <Section id="testing" title="Testing a migration before it meets production">
        <p>
          Running it against an empty development database tells you the syntax is valid and
          nothing else. Three checks that find real problems:
        </p>
        <Checklist
          items={[
            <>
              <strong>Run it against a copy with production-like volume.</strong> Lock duration is
              a function of row count, and every dangerous statement in this page is instant on a
              thousand rows.
            </>,
            <>
              <strong>Run it while the old version is serving traffic.</strong> A load generator
              against the old code during the migration is the only way to observe the lock you
              did not expect.
            </>,
            <>
              <strong>Check what the migration actually did,</strong> not what the struct says.
              A test schema built by <code>AutoMigrate</code> and a production schema built by
              migrations are different schemas, and a constraint that exists in one and not the
              other fails silently in the accepting direction: the test inserts the row and
              passes.
            </>,
          ]}
        />
        <CodeBlock
          language="sql"
          code={`-- During the migration, from another session: what is blocked, and by
-- what. The most useful query in an incident.
SELECT pid, state, wait_event_type, left(query, 80) AS query,
       now() - query_start AS duration
FROM pg_stat_activity
WHERE datname = current_database() AND state <> 'idle'
ORDER BY query_start;`}
        />
      </Section>

      <Section id="with-grit" title="Generated migrations">
        <CodeBlock
          terminal
          code={`grit generate migration add_total_cents_to_invoices
grit generate migration add_status_index --index invoices:status   # CONCURRENTLY
grit migrate                                        # advisory lock, checksums
grit migrate status                                 # what each run changed
grit migrate down                                   # undo the last run`}
        />
        <p>
          The <code>--index</code> form writes <code>CREATE INDEX CONCURRENTLY</code> and marks
          the migration as not transactional, because those two facts always travel together and
          forgetting the second turns the first into an error.
        </p>
        <p>
          <code>grit index</code> reports tables whose queries would benefit from an index and
          does not create one, because an index is a write cost as well as a read benefit and
          that is a decision rather than a fix.
        </p>
        <p>
          <Link href="/docs/backend/migrations" className="text-accent hover:underline">
            The migrations reference
          </Link>{' '}
          covers the runner, the checksums and the rollback.
        </p>
      </Section>

      <Faq
        slug="gorm-zero-downtime-migrations"
        items={[
          {
            q: 'Is AutoMigrate ever acceptable in production?',
            a: (
              <>
                For a single-instance internal tool where a few seconds of downtime is irrelevant,
                it is a reasonable trade. For anything with users it is not, mostly because it
                cannot be reviewed: nobody sees the statement before it runs, so nobody can notice
                that it is about to rewrite a large table.
              </>
            ),
            text:
              'For a single-instance internal tool where seconds of downtime are irrelevant, it is a reasonable trade. For anything with users it is not, mainly because the statement cannot be reviewed before it runs, so nobody notices when it is about to rewrite a large table.',
          },
          {
            q: 'How long is too long to hold a lock?',
            a: (
              <>
                Any <code>ACCESS EXCLUSIVE</code> lock over about a second on a busy table is
                visible to users, because every query queues behind it. Set{' '}
                <code>lock_timeout</code> to a few seconds so a migration that cannot get the lock
                fails rather than building a queue, and retry it when the table is quieter.
              </>
            ),
            text:
              'Any ACCESS EXCLUSIVE lock held for more than about a second on a busy table is user-visible, because queries queue behind it. Set lock_timeout to a few seconds so the migration fails rather than building a queue, then retry when the table is quieter.',
          },
          {
            q: 'Why does CREATE INDEX CONCURRENTLY fail inside my migration?',
            a: (
              <>
                It cannot run inside a transaction, and most migration runners wrap each file in
                one. The runner needs a way to mark a migration as non-transactional. The second
                thing to know is that a concurrent build can fail and leave an invalid index
                behind, which must be dropped before retrying; it does not clean up after itself.
              </>
            ),
            text:
              'It cannot run inside a transaction, and most runners wrap each migration in one, so the runner needs a non-transactional mode. A concurrent build can also fail and leave an invalid index behind, which must be dropped before retrying.',
          },
          {
            q: 'Do I need expand and contract for every change?',
            a: (
              <>
                No. Adding a nullable column, adding an index concurrently and dropping an unused
                index are all single-step and safe. Expand and contract is for changes where the
                old code and the new code disagree about the shape: renames, type changes,
                splitting a column, and making something required.
              </>
            ),
            text:
              'No. Adding a nullable column, adding an index concurrently and dropping an unused index are single-step and safe. Expand and contract is for changes where old and new code disagree about the shape: renames, type changes, splitting a column and making something required.',
          },
          {
            q: 'Should migrations run on application startup?',
            a: (
              <>
                No. It couples schema change to deployment, delays readiness for every instance,
                and turns a failed migration into a crash loop. Run them as a separate step that
                must succeed before the new version rolls out, so a failure is one visible failed
                job rather than an application that will not start.
              </>
            ),
            text:
              'No. It couples schema change to deployment, delays readiness for every instance and turns a failed migration into a crash loop. Run them as a separate step that must succeed before the new version rolls out.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/gorm-tutorial-relationships',
            label: 'GORM relationships',
            blurb: 'The models these migrations are changing.',
          },
          {
            href: '/guides/go-multi-tenancy',
            label: 'Multi-tenant SaaS in Go',
            blurb: 'Why schema-per-tenant makes every migration a job with a progress bar.',
          },
          {
            href: '/docs/backend/migrations',
            label: 'Migrations reference',
            blurb: 'The generated runner, checksums, the advisory lock and rollback.',
          },
          {
            href: '/docs/scaling',
            label: 'Scaling',
            blurb: 'Read replicas, connection pooling, and what a migration does to both.',
          },
        ]}
      />
    </GuidePage>
  )
}
