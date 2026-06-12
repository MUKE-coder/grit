import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        &quot;Who changed this and when?&quot; is a question every production app
        eventually answers. Grit&apos;s audit log answers it tamper-evidently —
        each row stores a SHA-256 hash of the previous row + its canonical
        fields. Delete a row, the chain breaks at verify time.
      </p>

      <h2>The audit log model</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/activity_log.go"
        code={`type ActivityLog struct {
    ID         uuid.UUID  \`gorm:"type:uuid;primaryKey"\`
    SequenceID int64      \`gorm:"uniqueIndex;autoIncrement"\`
    UserID     *uuid.UUID \`gorm:"type:uuid"\`        // nullable — system events
    Event      string     \`gorm:"not null;index"\`   // 'user.login', 'order.refund', ...
    Subject    string     \`gorm:"index"\`            // resource id this is about
    Metadata   datatypes.JSON
    IP         string
    UserAgent  string
    CreatedAt  time.Time

    PrevHash   string \`gorm:"index"\`   // SHA-256 of the previous row's CanonicalHash
    CanonicalHash string                // SHA-256 of THIS row's fields, including PrevHash
}`}
      />
      <p>
        Two hash columns. <code>PrevHash</code> chains back to the previous
        row; <code>CanonicalHash</code> is this row&apos;s identity. Hash
        construction:
      </p>
      <CodeBlock
        language="text"
        code={`CanonicalHash = SHA256(
    SequenceID + UserID + Event + Subject + Metadata + CreatedAt + PrevHash
)`}
      />

      <h2>Why this is tamper-evident</h2>
      <p>
        Suppose an attacker has DB access and tries to hide their actions:
      </p>
      <ul>
        <li>
          <strong>Delete a row</strong> — the next row&apos;s <code>PrevHash</code>{' '}
          doesn&apos;t match any existing row. Chain broken.
        </li>
        <li>
          <strong>Edit a row&apos;s fields</strong> — the row&apos;s{' '}
          <code>CanonicalHash</code> no longer matches what its content
          hashes to. Chain broken.
        </li>
        <li>
          <strong>Edit and recompute</strong> — possible, but every row
          AFTER the edit has stale <code>PrevHash</code> values that no longer
          match. To stay consistent, they&apos;d have to rebuild every
          subsequent row&apos;s hash — and the verifier still notices the
          mismatch with a hash you exported earlier.
        </li>
      </ul>
      <p>
        Same pattern as git commits and blockchain blocks. Cryptographic
        chaining, no central authority.
      </p>

      <h2>Writing audit events</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/order.go (excerpt)"
        code={`func (h *OrderHandler) Refund(c *gin.Context) {
    order, _ := h.svc.RefundByID(ctx, orderID)

    middleware.LogSecurityEvent(c, h.db, middleware.AuditEvent{
        Event:   "order.refund",
        Subject: order.ID.String(),
        Metadata: map[string]any{
            "amount": order.Total,
            "reason": "customer requested",
        },
    })
    respond.OK(c, order, "Refund processed")
}`}
      />
      <p>
        The middleware grabs <code>user_id</code> from context, computes the
        <code>PrevHash</code> from the last row, computes the new{' '}
        <code>CanonicalHash</code>, writes the row.
      </p>

      <h2>Standard events Grit pre-wires</h2>
      <ul>
        <li><code>user.login</code>, <code>user.logout</code></li>
        <li><code>user.register</code>, <code>user.role_changed</code></li>
        <li><code>user.password_changed</code>, <code>user.totp_enabled</code></li>
        <li><code>account.locked</code> (AuthShield triggered)</li>
        <li><code>authz.denied</code> (RequireRoles rejected)</li>
        <li><code>sentinel.suspicious</code> (WAF / Anomaly fired)</li>
      </ul>
      <p>
        Add your domain events alongside: <code>order.created</code>,{' '}
        <code>refund.issued</code>, <code>invoice.voided</code>. Anything
        material to your business.
      </p>

      <TipBox tone="success">
        <strong>Compliance gold:</strong> SOC 2 Type 2 auditors specifically
        ask about audit-log integrity. &quot;Show me you can&apos;t silently edit
        history&quot;. Grit&apos;s hash chain answers this in 20 seconds.
      </TipBox>

      <h2>Verifying the chain</h2>
      <p>
        Run periodically (cron job, weekly batch):
      </p>
      <CodeBlock
        language="go"
        code={`func VerifyAuditChain(db *gorm.DB) error {
    var rows []ActivityLog
    db.Order("sequence_id ASC").Find(&rows)

    var prevHash string
    for _, row := range rows {
        if row.PrevHash != prevHash {
            return fmt.Errorf("chain break at seq %d", row.SequenceID)
        }
        if row.CanonicalHash != computeHash(row, prevHash) {
            return fmt.Errorf("row tampered at seq %d", row.SequenceID)
        }
        prevHash = row.CanonicalHash
    }
    return nil
}`}
      />
      <p>
        Alert / page when this returns an error. Now you&apos;re informed about
        DB-level tampering.
      </p>

      <h2>Performance</h2>
      <p>
        Audit writes are synchronous (you don&apos;t want to lose audit events
        in a worker queue). Each write is one SQL INSERT + one SHA-256 hash —
        sub-millisecond on modern hardware. For very high-write systems,
        you can move this to the worker queue with an idempotency key, at
        the cost of weakened guarantees.
      </p>

      <KnowledgeCheck
        question="A disgruntled employee with DB access deletes their own `user.role_changed` audit row to hide that they made themselves admin. What happens at the next chain verification?"
        choices={[
          {
            label: 'Nothing — the DB doesn\'t notice deletes',
            feedback:
              "Wrong at the chain level. The DB itself doesn't notice, but the verify job does — the next row's PrevHash references a row that no longer exists.",
          },
          {
            label: 'The next row\'s PrevHash points to a row that no longer exists — chain break detected',
            correct: true,
            feedback:
              "Right — the audit ID + hash of the deleted row was the PrevHash for the row that came after. With the row gone, the chain has a hole the verifier immediately finds.",
          },
          {
            label: 'GORM throws an error at delete time',
            feedback:
              "Wrong — GORM does not enforce hash chain integrity. The DB happily deletes the row; the chain check is application-level (run on a schedule).",
          },
          {
            label: 'The hashes auto-recompute',
            feedback:
              "Wrong — hashes are stored at write time, not computed on the fly. That's the whole point — the stored hashes are the integrity record.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Trigger an audit event and inspect the chain:
            </p>
            <ol>
              <li>Log into your API — that fires <code>user.login</code>.</li>
              <li>
                Open GORM Studio &amp; look at the <code>activity_logs</code>{' '}
                table.
              </li>
              <li>
                Find the most recent row. Note its <code>PrevHash</code>{' '}
                and <code>CanonicalHash</code>.
              </li>
              <li>
                In <code>notes.md</code>: paste them, and explain what would
                happen if you edited <code>Event</code> from{' '}
                <code>user.login</code> to <code>user.logout</code>.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            Editing Event changes the inputs to{' '}
            <code>CanonicalHash</code>. The stored hash would no longer
            match the computed-from-current-content hash. The verify job
            detects this immediately.
          </>
        }
        solution={
          <>
            <p>
              Two columns you should capture:
            </p>
            <CodeBlock
              language="text"
              code={`PrevHash:      sha256(...) of previous row
CanonicalHash: sha256(seq=1234 + user_id=... + event="user.login" + ...) = abcd...`}
            />
            <p>
              An edit to <code>Event</code> would change the inputs to{' '}
              <code>CanonicalHash</code>. The stored hash (<code>abcd...</code>){' '}
              would no longer match the recomputed hash. Verify job
              raises an alert.
            </p>
            <p>That completes chapter 5&apos;s assignment.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final chapter — <strong>Deploy</strong>. Take what you built to a
        VPS with one command and a domain name.
      </p>
    </>
  )
}
