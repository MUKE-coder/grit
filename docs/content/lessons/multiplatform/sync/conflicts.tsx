import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Mobile edits a note offline. Desktop edits the same note
        offline. Both reconnect. Whose change wins? Without a policy,
        you get data loss or silent corruption. This lesson is the
        policy.
      </p>

      <h2>The shape of the conflict</h2>
      <Diagram label="The two-actor conflict" caption="Both clients changed the same row while disconnected. Reconciliation needs a deterministic rule.">
{`   time
   ────────────────────────────────────────────►

   Mobile:   read note (v1) ── offline edit ──► push (v2-mobile)
                                                       │
                                                       ▼
   Server:   v1 ─────────────────────────────────── ? ─────────►
                                                       ▲
                                                       │
   Desktop:  read note (v1) ── offline edit ──► push (v2-desktop)`}
      </Diagram>

      <h2>Three policies, ranked by complexity</h2>

      <h3>1. Last-write-wins (LWW)</h3>
      <p>
        Simplest, cheapest, lossy. Whichever push arrives second
        clobbers the first. Defensible when:
      </p>
      <ul>
        <li>The same user owns both devices (low collision).</li>
        <li>The data has low value (preferences, bookmarks).</li>
        <li>You can&apos;t reasonably merge two states.</li>
      </ul>
      <p>
        Default for Grit. The bookmark feature uses LWW because: it&apos;s
        a personal toggle, conflicts are rare, and the user wouldn&apos;t
        notice if the wrong &quot;last edit&quot; wins.
      </p>

      <h3>2. Server-wins / Client-wins</h3>
      <p>
        A constant rule: the server&apos;s version always wins (the
        client&apos;s pending changes are dropped) OR the client
        always wins (the server takes whatever the client sends).
      </p>
      <ul>
        <li>
          <strong>Server-wins:</strong> safe for read-mostly data
          (catalogs, settings managed in admin).
        </li>
        <li>
          <strong>Client-wins:</strong> simple for user-owned data
          (notes, drafts) but loses if two clients race.
        </li>
      </ul>

      <h3>3. Version + 409 + manual merge</h3>
      <p>
        Each row has a <code>version</code> column. The client sends
        the version it&apos;s editing. The server compares; if newer
        than expected, returns 409 with the latest. The client merges
        (or asks the user).
      </p>
      <p>
        Highest fidelity, highest cost. Use it for high-stakes data —
        invoice totals, contract text, anything where silent data loss
        is unacceptable.
      </p>

      <h2>Implementing LWW with timestamps</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/bookmark.go"
        code={`// Client sends If-Unmodified-Since with the timestamp it read.
// Server compares and accepts only if the row hasn't changed.

func (h *BookmarkHandler) Update(c *gin.Context) {
  var input UpdateInput
  if err := c.ShouldBindJSON(&input); err != nil { /* 400 */ }

  var existing models.Bookmark
  h.db.First(&existing, c.Param("id"))

  clientTs, _ := time.Parse(time.RFC3339, c.GetHeader("If-Unmodified-Since"))
  if existing.UpdatedAt.After(clientTs) {
    // The server has a newer version. Accept the new write anyway (LWW)
    // BUT log it so we can see how often it happens.
    auditConflict(existing, input)
  }

  existing.Note = input.Note
  h.db.Save(&existing)
  c.JSON(200, gin.H{"data": existing})
}`}
      />
      <p>
        Pure LWW: the most recent write wins. The conflict is logged
        but doesn&apos;t change the outcome. Use this until your audit
        log shows conflicts often enough to need 409.
      </p>

      <h2>Implementing 409 + manual merge</h2>
      <CodeBlock
        language="go"
        code={`func (h *NoteHandler) Update(c *gin.Context) {
  var input UpdateNoteInput
  c.ShouldBindJSON(&input)

  var existing models.Note
  h.db.First(&existing, c.Param("id"))

  if existing.Version != input.BaseVersion {
    c.JSON(409, gin.H{
      "error": "version_conflict",
      "server_version": existing,
    })
    return
  }
  existing.Body = input.Body
  existing.Version++
  h.db.Save(&existing)
  c.JSON(200, gin.H{"data": existing})
}`}
      />
      <CodeBlock
        language="tsx"
        filename="apps/desktop/frontend (snippet)"
        code={`async function saveNote(note: Note, body: string) {
  const res = await fetch(\`/api/notes/\${note.id}\`, {
    method: 'PUT',
    body: JSON.stringify({ body, base_version: note.version }),
  })
  if (res.status === 409) {
    const { server_version } = await res.json()
    return showMergeUI(local: body, server: server_version.body)
  }
}`}
      />
      <p>
        <code>showMergeUI</code> is up to you — a side-by-side diff,
        a &quot;keep mine / keep theirs&quot; toggle, or an
        auto-merge for non-overlapping fields.
      </p>

      <h2>Per-field policies</h2>
      <p>
        Real apps mix policies per field within the same record:
      </p>
      <ul>
        <li><code>name</code> — LWW. Cheap, low stakes.</li>
        <li><code>price_cents</code> — 409. High stakes, requires care.</li>
        <li><code>updated_at</code> — server-owned. Client cannot edit.</li>
        <li><code>tags</code> — merge (set-union). No data loss.</li>
      </ul>
      <p>
        That last one is poor man&apos;s CRDT: any field that&apos;s a
        set, list, or counter can usually be merged commutatively
        (union, max, sum). Use this where you can.
      </p>

      <TipBox tone="warning">
        <strong>Conflicts you can&apos;t resolve are bugs.</strong> If
        the user can&apos;t tell the system how to reconcile, neither
        can your code. Don&apos;t pretend by silently picking one. Log
        it, surface it (&quot;Two versions of this note — pick
        one&quot;), or prevent it (lock-on-edit).
      </TipBox>

      <h2>Lock-on-edit — the cheap escape hatch</h2>
      <p>
        If conflicts are rare AND high-stakes, lock the record while
        someone&apos;s editing. Google Docs solved this differently
        (real-time collab), but smaller apps can use a stale lock:
      </p>
      <CodeBlock
        language="go"
        code={`// When user opens for edit:
POST /api/notes/123/lock        → 200 if lock acquired, 409 if locked by someone else

// Lock auto-releases after 5 minutes of inactivity.`}
      />
      <p>
        Pessimistic, but for low-frequency edits across desktop + web
        + mobile (e.g., a shared invoice), it&apos;s the simplest
        thing that works.
      </p>

      <h2>The audit log</h2>
      <p>
        Regardless of policy, log every conflict. The Grit API has
        an audit_logs table — use it:
      </p>
      <CodeBlock
        language="go"
        code={`auditLog.Create(&AuditLog{
  Entity: "bookmark", EntityID: existing.ID,
  Action: "conflict_lww", UserID: userID,
  Meta: jsonOf(map[string]any{
    "server_ts": existing.UpdatedAt, "client_ts": clientTs,
    "client_origin": c.GetHeader("X-Client"),
  }),
})`}
      />
      <p>
        Now you can answer &quot;how often do conflicts happen?&quot;
        with data. If the answer is &quot;5 per month&quot;, LWW is
        fine. If it&apos;s 500/day, time for 409 + merge.
      </p>

      <KnowledgeCheck
        question="Your app shows the user's notes across web, mobile, desktop. Mobile + desktop both edit the same note offline. Which conflict policy is the best STARTING point?"
        choices={[
          {
            label: 'CRDT-based real-time merge',
            feedback:
              'Overkill for a notes app — months of work for a problem that may be rare. Start simpler, measure, then upgrade.',
          },
          {
            label: 'Last-write-wins + audit log; upgrade to 409+merge ONLY if the audit log shows conflicts are frequent',
            correct: true,
            feedback:
              "Right — LWW is cheap to ship and the audit log tells you whether to invest in something heavier. Don't build 409+merge speculatively; build it when data says you should.",
          },
          {
            label: 'Lock-on-edit',
            feedback:
              'Heavy-handed for notes — users hate "someone else is editing, try later" on their own data on their own devices.',
          },
          {
            label: 'Force the user to pick a winner every time',
            feedback:
              "Loud and annoying. The user shouldn't have to think for every save. Surface it only for actual conflicts.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Implement the LWW policy + audit log on bookmarks:</p>
            <ol>
              <li>
                Pick one model in your API (e.g., <code>Note</code>) and add
                a <code>Body</code> field if it doesn&apos;t have one.
              </li>
              <li>
                In the handler, accept <code>If-Unmodified-Since</code>{' '}
                from the client. Apply the new value either way (LWW),
                but log conflicts.
              </li>
              <li>
                On desktop: queue a note update via the outbox while
                offline.
              </li>
              <li>
                On mobile: edit the same note while offline.
              </li>
              <li>
                Reconnect both. Confirm whichever pushed last is what
                shows up, and that the audit log has ONE conflict
                entry.
              </li>
              <li>
                Bonus: render conflict count on an admin dashboard widget.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For the racing test, send mobile first, wait a beat,
            then send desktop. If both arrive simultaneously, your
            DB transactional isolation matters — Postgres&apos; default
            (READ COMMITTED) is enough for LWW.
          </>
        }
        solution={
          <>
            <p>
              After the dust settles: both clients see the
              latest-writer&apos;s version on next refetch.
              You&apos;ve made an explicit choice (vs. silently losing
              data) and you have metrics to revisit it later.
            </p>
            <p>
              For most consumer apps, LWW + audit is the right answer
              for years. By the time conflicts justify 409+merge,
              you&apos;ll know — and you&apos;ll have data to design
              the merge UI well.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 5 — <strong>Coordinated releases</strong>. Now that
        all four surfaces talk to one API and survive offline, we
        face the last problem: shipping changes without breaking old
        clients still in the wild.
      </p>
    </>
  )
}
