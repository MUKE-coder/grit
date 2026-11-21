import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Tamper-evident audit logs with hash chains in Postgres',
  description:
    'Chaining entries so a deletion is detectable, keying the hash so it cannot be recomputed, verifying at scale, and pruning without breaking the chain.',
  alternates: { canonical: 'https://gritframework.dev/guides/hash-chained-audit-log' },
  openGraph: {
    title: 'Tamper-evident audit logs with hash chains in Postgres',
    description: 'Tamper-evident, not tamper-proof, and the difference matters.',
    url: 'https://gritframework.dev/guides/hash-chained-audit-log',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="hash-chained-audit-log"
      lead={
        <>
          <p className="mb-4">
            An audit log in a table anybody can write to proves nothing. Whoever could do the
            thing you are auditing can usually also edit the row that records it, and an auditor
            asking &quot;how do you know this is complete&quot; has no answer.
          </p>
          <p>
            A hash chain changes that. Each entry includes the hash of the previous one, so
            altering or deleting any entry invalidates every entry after it, and the mismatch is
            detectable. It is about forty lines, and the subtleties are in verification, keying
            and pruning rather than in the chaining.
          </p>
        </>
      }
    >
      <Outcome>
        <p>An append-only log whose integrity you can check with a command.</p>
        <p>A keyed hash, so somebody with write access cannot simply recompute the chain.</p>
        <p>A pruning strategy that does not destroy the property you built it for.</p>
      </Outcome>

      <Section id="evident" title="Tamper-evident, not tamper-proof">
        <p>
          Be precise about the claim, because it is the first thing a reviewer will test. A hash
          chain does not prevent anybody from changing the data. Somebody with database access
          can rewrite any row they like. What they cannot do is leave the chain consistent while
          doing it, unless they also hold the key and recompute everything after it.
        </p>
        <p>So the property is: a change is detectable by anybody who runs the verification.</p>
        <Checklist
          items={[
            <>
              <strong>It defends against:</strong> a quiet <code>UPDATE</code> to make an action
              look different, a <code>DELETE</code> of an inconvenient entry, and an insert
              backdated into the middle of the history.
            </>,
            <>
              <strong>It does not defend against:</strong> an attacker who holds the key and has
              time to recompute the whole chain, nor against something that was never written
              down in the first place.
            </>,
            <>
              <strong>It is worth having anyway,</strong> because it converts &quot;we trust the
              table&quot; into &quot;here is a check that passes&quot;, which is the difference
              between an assertion and evidence.
            </>,
          ]}
        />
      </Section>

      <Section id="table" title="The table">
        <CodeBlock
          language="sql"
          code={`CREATE TABLE audit_entries (
  id          bigserial   PRIMARY KEY,
  occurred_at timestamptz NOT NULL DEFAULT clock_timestamp(),

  actor_type  text        NOT NULL,       -- user | system | api_key
  actor_id    bigint,
  actor_label text,                       -- who this was AT THE TIME

  action      text        NOT NULL,       -- invoice.void
  entity_type text        NOT NULL,
  entity_id   bigint,
  org_id      bigint,

  changes     jsonb,                      -- before and after, redacted
  ip          inet,
  request_id  text,

  prev_hash   bytea       NOT NULL,
  hash        bytea       NOT NULL
);

CREATE INDEX ON audit_entries (entity_type, entity_id, id DESC);
CREATE INDEX ON audit_entries (actor_id, id DESC);
CREATE INDEX ON audit_entries (org_id, occurred_at DESC);`}
        />
        <Note title="actor_label is denormalised on purpose">
          <p>
            Joining to <code>users</code> to render the log shows who that person is{' '}
            <em>now</em>. The audit question is who they were <em>then</em>: their name, their
            role, their email at the moment of the action. A user renamed or deleted since must
            not change the history, so the label is copied in and never updated.
          </p>
        </Note>
      </Section>

      <Section id="chain" title="Writing a chained entry">
        <CodeBlock
          filename="internal/audit/write.go"
          code={`func (a *Auditor) Write(ctx context.Context, tx *gorm.DB, e Entry) error {
	// The write takes the caller's transaction, so the audit entry and
	// the thing it describes commit together. An audit log written
	// outside the transaction records actions that were rolled back,
	// and misses actions that succeeded after it failed.

	// SELECT ... FOR UPDATE serialises appenders. Without it, two
	// concurrent writes read the same prev_hash and the chain forks.
	var prev []byte
	err := tx.Raw(\`
		SELECT hash FROM audit_entries ORDER BY id DESC LIMIT 1 FOR UPDATE\`).
		Scan(&prev).Error
	if err != nil {
		return fmt.Errorf("read chain head: %w", err)
	}
	if prev == nil {
		prev = genesis    // a fixed, documented value for the first entry
	}

	row := models.AuditEntry{ /* ... */ PrevHash: prev}
	row.Hash = a.digest(prev, row)

	if err := tx.Create(&row).Error; err != nil {
		return fmt.Errorf("append audit entry: %w", err)
	}
	return nil
}`}
        />
        <CodeBlock
          code={`// HMAC, not a bare hash. A plain SHA-256 chain can be recomputed by
// anybody who can read the rows, so somebody with write access can
// alter an entry and rewrite every hash after it. With a key they do
// not have, they cannot.
//
// The key lives wherever your other secrets live, and must NOT be in
// the same database as the log: a single compromise should not yield
// both the rows and the means to rewrite them.
func (a *Auditor) digest(prev []byte, e models.AuditEntry) []byte {
	mac := hmac.New(sha256.New, a.key)
	mac.Write(prev)
	// A canonical, unambiguous encoding. Concatenating fields without
	// separators lets two different entries produce the same input:
	// ("ab", "c") and ("a", "bc") are the same bytes. Length-prefixing
	// or a delimiter that cannot appear in the values prevents that.
	writeField(mac, e.OccurredAt.UTC().Format(time.RFC3339Nano))
	writeField(mac, e.ActorType)
	writeField(mac, strconv.FormatInt(e.ActorID, 10))
	writeField(mac, e.Action)
	writeField(mac, e.EntityType)
	writeField(mac, strconv.FormatInt(e.EntityID, 10))
	writeField(mac, string(canonicalJSON(e.Changes)))
	return mac.Sum(nil)
}

func writeField(w io.Writer, s string) {
	binary.Write(w, binary.BigEndian, uint32(len(s)))
	io.WriteString(w, s)
}`}
        />
        <Warn title="Hash what the database will store, not what you have in memory">
          <p>
            Postgres stores <code>timestamptz</code> at microsecond precision and MySQL at
            millisecond. A Go <code>time.Time</code> has nanoseconds. Hash the nanosecond value,
            read the row back, recompute, and the hash differs, so the chain fails verification
            on data nobody touched. Round or truncate to what the column stores before hashing,
            and verify on the database you actually run rather than on SQLite.
          </p>
        </Warn>
      </Section>

      <Section id="verify" title="Verifying">
        <CodeBlock
          code={`func (a *Auditor) Verify(ctx context.Context, from, to int64) (*Report, error) {
	var prev []byte
	rep := &Report{}

	rows, err := a.db.WithContext(ctx).
		Model(&models.AuditEntry{}).
		Where("id BETWEEN ? AND ?", from, to).
		Order("id").Rows()
	if err != nil {
		return nil, fmt.Errorf("scan entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var e models.AuditEntry
		if err := a.db.ScanRows(rows, &e); err != nil {
			return nil, fmt.Errorf("scan entry: %w", err)
		}
		if prev != nil && !bytes.Equal(e.PrevHash, prev) {
			// A gap: an entry was deleted, or one was inserted here.
			rep.Breaks = append(rep.Breaks, Break{ID: e.ID, Kind: "chain"})
		}
		if !hmac.Equal(a.digest(e.PrevHash, e), e.Hash) {
			// The row's own contents were altered.
			rep.Breaks = append(rep.Breaks, Break{ID: e.ID, Kind: "content"})
		}
		prev = e.Hash
		rep.Checked++
	}
	return rep, rows.Err()
}`}
        />
        <p>
          Two kinds of break, and they mean different things. A content break says this entry was
          edited. A chain break says something was removed or inserted between two entries.
          Reporting them separately is what makes the output actionable rather than just alarming.
        </p>
        <H3>Verifying at scale</H3>
        <p>
          Forty million entries is not a verification you run on every request. Three things make
          it practical:
        </p>
        <Checklist
          items={[
            <>
              <strong>Checkpoints.</strong> Every ten thousand entries, store the hash in a
              separate table. Verification then runs from the last checkpoint rather than from
              the beginning.
            </>,
            <>
              <strong>Nightly, incremental.</strong> Verify only what has been added since the
              last run, and keep the result.
            </>,
            <>
              <strong>Publish the head.</strong> Write the current head hash somewhere outside
              the system: a log sent to a different provider, or an email to the compliance
              inbox. An attacker who controls your database cannot alter what already left it,
              which upgrades the claim considerably.
            </>,
          ]}
        />
      </Section>

      <Section id="write-access" title="Making the table append-only">
        <p>
          The chain detects changes. Stopping them in the first place is a separate, cheap
          control.
        </p>
        <CodeBlock
          language="sql"
          code={`-- The application's role may insert and read, and nothing else.
REVOKE UPDATE, DELETE, TRUNCATE ON audit_entries FROM app_user;
GRANT INSERT, SELECT ON audit_entries TO app_user;

-- And a trigger, which also catches a superuser doing it by hand.
CREATE OR REPLACE FUNCTION audit_is_append_only() RETURNS trigger AS $$
BEGIN
  RAISE EXCEPTION 'audit_entries is append-only';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_no_update BEFORE UPDATE OR DELETE ON audit_entries
  FOR EACH ROW EXECUTE FUNCTION audit_is_append_only();`}
        />
        <p>
          Together these are better than either alone: the grants stop the application, the
          trigger stops an interactive session, and the chain detects anybody who got past both
          by dropping the trigger.
        </p>
      </Section>

      <Section id="what" title="What to record, and what not to">
        <Checklist
          items={[
            <>
              <strong>Record:</strong> authentication events including failures, permission
              denials, role and permission changes, every write to anything financial or
              personal, exports and bulk downloads, administrative impersonation, and
              configuration changes.
            </>,
            <>
              <strong>Record reads too, selectively.</strong> &quot;Who looked at this
              customer&apos;s record&quot; is a question that gets asked, usually at the worst
              moment, and it cannot be answered retroactively. Auditing every read is too much;
              auditing reads of the sensitive entities is not.
            </>,
            <>
              <strong>Never record:</strong> passwords, tokens, session identifiers, full card
              numbers, or anything you would have to delete later for a subject access request.
            </>,
          ]}
        />
        <Warn title="Redact before hashing, not after">
          <p>
            If a value must be removable, it must never enter the hash, because removing it later
            breaks the chain. Hash a redacted form from the start, or hash a reference to the
            value rather than the value. Designing this in is easy; retrofitting it means
            choosing between the regulation and the chain.
          </p>
        </Warn>
      </Section>

      <Section id="pruning" title="Pruning without breaking it">
        <p>
          Retention policies and deletion requests collide with a chain by construction: removing
          an entry breaks every hash after it.
        </p>
        <H3>Segment the chain</H3>
        <CodeBlock
          language="text"
          code={`One chain per period, each sealed when the period ends:

  2026-09  [genesis → ... → head]  sealed, head published
  2026-10  [genesis → ... → head]  sealed, head published
  2026-11  [genesis → ... → current]

Dropping all of 2026-09 removes a whole segment and leaves every other
segment verifiable. The cost: you can no longer prove that the September
segment ever existed, which is exactly what deleting it means.`}
        />
        <H3>Tombstone rather than delete</H3>
        <p>
          For a personal-data erasure request, replace the fields with a marker and record the
          erasure as its own audit entry. The chain stays intact because the entry still exists,
          the personal data is gone, and the fact of the erasure is itself auditable, which is
          usually what the regulation wants.
        </p>
        <p>
          This only works if the hash covered the redacted form from the beginning, which is the
          previous section&apos;s point and the reason it is worth deciding early.
        </p>
      </Section>

      <Section id="compliance" title="What an auditor actually asks">
        <p>
          The technical property is the easy part. These are the questions that follow, and a
          hash chain only answers the first one.
        </p>
        <Checklist
          items={[
            <>
              <strong>&quot;How do you know the log is complete?&quot;</strong> The chain, plus a
              published head.
            </>,
            <>
              <strong>&quot;Who can modify it?&quot;</strong> The grants and the trigger. Show
              that the application role cannot.
            </>,
            <>
              <strong>&quot;How would you know if somebody did?&quot;</strong> The nightly
              verification and what it alerts to. An unmonitored verification is not a control.
            </>,
            <>
              <strong>&quot;How long is it kept, and how is it disposed of?&quot;</strong> The
              segmentation, and the job that enforces it.
            </>,
            <>
              <strong>&quot;Show me every action by this user last quarter.&quot;</strong> Which
              is a question about indexes and a screen, and is the one that actually gets asked
              most often.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="Generated">
        <CodeBlock
          terminal
          code={`grit add audit
grit audit verify                  # the whole chain, or a range
grit audit verify --since 30d`}
        />
        <p>
          That writes the table, the HMAC-keyed chain, the append-only grants and trigger, the
          verification command, the checkpoints and the admin screen, with the audit write taking
          the caller&apos;s transaction so an entry cannot exist for a rolled-back action.
        </p>
        <p>
          The key is read from the environment and is deliberately not stored alongside the log,
          for the reason above: the chain is only worth something if altering it requires
          something the database does not contain.
        </p>
        <p>
          <Link href="/docs/security/audit-log" className="text-accent hover:underline">
            The audit log reference
          </Link>{' '}
          covers the generated schema, the verification and retention.
        </p>
      </Section>

      <Faq
        slug="hash-chained-audit-log"
        items={[
          {
            q: 'Is this a blockchain?',
            a: (
              <>
                It is the one idea blockchains borrowed, a chain of hashes, without the part that
                makes them blockchains, which is distributed consensus about what the chain
                contains. You have one writer and trust your own database, so consensus solves a
                problem you do not have. Calling it a hash chain keeps the conversation honest.
              </>
            ),
            text:
              'It uses the chained-hash idea that blockchains also use, without distributed consensus, which is what makes a blockchain a blockchain. With one writer and your own database, consensus solves a problem you do not have.',
          },
          {
            q: 'Why HMAC rather than a plain SHA-256?',
            a: (
              <>
                Because a plain hash can be recomputed by anybody who can read the table. An
                attacker with write access edits an entry and rewrites every hash after it, and
                the verification passes. With a key they do not hold, they cannot produce valid
                hashes, so the tampering is detectable. Keep the key outside the database.
              </>
            ),
            text:
              'A plain hash can be recomputed by anybody who can read the table, so an attacker with write access can edit an entry and rewrite every hash after it, and verification still passes. A key they do not hold makes that impossible, provided the key is stored outside the database.',
          },
          {
            q: 'Does the chain hurt write throughput?',
            a: (
              <>
                Yes, by design: appending requires reading the current head under a lock, so
                audit writes serialise. In practice a few thousand per second on Postgres, which
                is far more than most applications audit. If it becomes a limit, segment per
                tenant so each tenant has its own chain and the lock is per tenant.
              </>
            ),
            text:
              'Yes, by design: appending reads the current head under a lock, so audit writes serialise. In practice a few thousand per second on Postgres, well beyond what most applications audit. If it becomes a limit, segment per tenant so the lock is per tenant.',
          },
          {
            q: 'Should the audit log live in a separate database?',
            a: (
              <>
                It is stronger, and it costs you the ability to write the entry in the same
                transaction as the action, which reintroduces the dual-write problem: an action
                committed with no audit entry, or an entry for something that rolled back. The
                usual compromise is the same database, with the restricted role and the trigger,
                plus shipping entries to a separate store asynchronously for the
                defence-in-depth copy.
              </>
            ),
            text:
              'It is stronger but costs the ability to write the entry in the same transaction as the action, reintroducing the dual-write problem. The usual compromise is the same database with a restricted role and a trigger, plus asynchronous shipping to a separate store as a second copy.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-authorization-policies-row-rules',
            label: 'Policies and row rules',
            blurb: 'The permission denials worth recording, and why they are worth recording.',
          },
          {
            href: '/guides/owasp-top-10-golang',
            label: 'The OWASP Top 10 in Go',
            blurb: 'Logging and monitoring failures, which is this subject on the list.',
          },
          {
            href: '/docs/security/audit-log',
            label: 'Audit log reference',
            blurb: 'The generated schema, verification command and retention.',
          },
          {
            href: '/docs/backend/append-only',
            label: 'Append-only tables',
            blurb: 'The grants and triggers, as generated.',
          },
        ]}
      />
    </GuidePage>
  )
}
