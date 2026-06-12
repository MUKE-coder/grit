import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        The sync engine drains the outbox to the server, pulls down server
        changes, and handles conflicts. This lesson covers the loop +
        the conflict-resolution policy you&apos;ll choose for your project.
      </p>

      <h2>The loop</h2>
      <Diagram label="Sync engine loop" caption="Runs every 30s + on reconnect. Push first, then pull. Conflicts surface to the user; never auto-overwrite their work.">
{`   ┌─────────────────────────────────────────────────┐
   │   ┌────────┐                                    │
   │   │  loop  │  every 30s + on reconnect          │
   │   └───┬────┘                                    │
   │       │                                         │
   │       ▼                                         │
   │   ┌──────────────────────┐                      │
   │   │  online?  no  → wait │                      │
   │   └─────┬────────────────┘                      │
   │         │ yes                                   │
   │         ▼                                       │
   │   ┌─────────────────────┐                       │
   │   │  squash outbox      │  reduce churn         │
   │   └─────────┬───────────┘                       │
   │             ▼                                   │
   │   ┌─────────────────────┐                       │
   │   │  POST /api/sync/push│  send local changes   │
   │   │   {batch:[…]}        │                      │
   │   └─────────┬───────────┘                       │
   │             ▼                                   │
   │   ┌─────────────────────┐                       │
   │   │  delete sent rows   │  from outbox          │
   │   │  mark sale.ServerID │                       │
   │   └─────────┬───────────┘                       │
   │             ▼                                   │
   │   ┌─────────────────────┐                       │
   │   │  GET /api/sync/pull │  fetch server changes │
   │   │   ?since=<cursor>   │  since last pull      │
   │   └─────────┬───────────┘                       │
   │             ▼                                   │
   │   ┌─────────────────────┐                       │
   │   │  apply to local DB  │  conflict? → resolve  │
   │   │  update cursor      │                       │
   │   └─────────────────────┘                       │
   │                                                 │
   └─────────────────────────────────────────────────┘`}
      </Diagram>

      <h2>The push call</h2>
      <CodeBlock
        language="go"
        filename="sync/engine.go"
        code={`func (e *Engine) Push(ctx context.Context) error {
    e.Squash()

    var entries []OutboxEntry
    e.db.Order("id ASC").Limit(100).Find(&entries)
    if len(entries) == 0 { return nil }

    res, err := e.client.Post("/api/sync/push", entries)
    if err != nil {
        return fmt.Errorf("network: %w", err)
    }

    // Per-entry result lets us track partial successes
    for _, r := range res.Results {
        if r.Status == "ok" {
            e.db.Delete(&OutboxEntry{ID: r.LocalID})
            e.db.Model(&Sale{}).Where("id = ?", r.EntityID).Update("server_id", r.ServerID)
        } else if r.Status == "conflict" {
            e.handleConflict(r)
        } else {
            e.db.Model(&OutboxEntry{ID: r.LocalID}).Update("attempts", gorm.Expr("attempts + 1"))
            e.db.Model(&OutboxEntry{ID: r.LocalID}).Update("last_error", r.Error)
        }
    }
    return nil
}`}
      />

      <h2>The pull cursor</h2>
      <p>
        For pull, the server needs to know &quot;changes since when&quot;. Grit
        uses a per-device cursor (a server-side timestamp or sequence
        number) stored in a one-row <code>Cursor</code> table:
      </p>
      <CodeBlock
        language="go"
        code={`type Cursor struct {
    Entity string \`gorm:"primaryKey"\`   // 'sale'
    Value  string                         // server-issued cursor
}

// Pull
func (e *Engine) Pull(ctx context.Context) error {
    var cur Cursor
    e.db.First(&cur, "entity = 'sale'")
    res, err := e.client.Get("/api/sync/pull?since=" + cur.Value)
    if err != nil { return err }

    for _, row := range res.Changes {
        // Upsert into local DB
        e.db.Save(&row)
    }
    cur.Value = res.NextCursor
    e.db.Save(&cur)
    return nil
}`}
      />

      <h2>Conflict resolution — the three strategies</h2>
      <ul>
        <li>
          <strong>Last-write-wins (LWW)</strong> — newer timestamp clobbers
          older. Simple, lossy. Fine for non-financial data (settings,
          activity logs).
        </li>
        <li>
          <strong>Server-wins</strong> — local change is dropped if the
          server&apos;s version is newer. Lossy on the client side. Use when
          the server is the source of truth (inventory authoritative).
        </li>
        <li>
          <strong>Manual resolution</strong> — mark the row{' '}
          <code>has_conflict=true</code>, surface in UI for the user to
          choose. The right answer for financial / customer-facing data.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>Never silently merge financial data.</strong> Sales,
        invoices, refunds — if the server says $100 and the client says
        $120, picking either silently is wrong. Surface the conflict; let
        the user decide.
      </TipBox>

      <h2>The Wails wiring</h2>
      <CodeBlock
        language="go"
        filename="app.go (excerpt)"
        code={`func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Background sync loop — every 30s, and on online events
    go a.sync.Loop(ctx)
}

// Bound for the React side to force a sync (button press)
func (a *App) SyncNow() error {
    if err := a.sync.Push(a.ctx); err != nil { return err }
    return a.sync.Pull(a.ctx)
}`}
      />
      <p>
        Auto-runs in the background; the user can also press a Sync button
        to force one. Both call the same engine.
      </p>

      <KnowledgeCheck
        question="Two cashiers edit the same product offline. Cashier A changes the price to $10; Cashier B changes the description. They both sync. What's the safest behaviour?"
        choices={[
          {
            label: 'Last write wins — whoever syncs last overwrites',
            feedback:
              "Wrong — that loses one of the legitimate edits. Cashier A\'s price change OR B\'s description change is silently dropped.",
          },
          {
            label: 'Field-level merge — A\'s price + B\'s description coexist',
            correct: true,
            feedback:
              "Right — they edited DIFFERENT fields. A proper sync engine merges field-by-field; the row ends up with the new price AND the new description. No conflict.",
          },
          {
            label: 'Reject both, ask the user',
            feedback:
              "Too pessimistic — there\'s no actual conflict. Both edits are compatible.",
          },
          {
            label: 'Random pick',
            feedback:
              "Catastrophically wrong. Don\'t flip a coin with customer data.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>For the chapter assignment, simulate the full offline → online cycle:</p>
            <ol>
              <li>Disconnect from the internet (turn off Wi-Fi).</li>
              <li>Record 3 notes in your app.</li>
              <li>
                Watch the &quot;pending sync&quot; counter — should be 3.
              </li>
              <li>Reconnect to the internet.</li>
              <li>
                Within 30 seconds (or after clicking Sync), the counter
                drops to 0.
              </li>
              <li>
                Check your server — the 3 notes should appear.
              </li>
            </ol>
            <p>Paste a video or screenshot sequence in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            For testing without a real server, point the sync engine at{' '}
            <code>localhost:8080</code> with a Grit API running. The Grit
            API kit includes a <code>/api/sync/push</code> endpoint you
            can scaffold.
          </>
        }
        solution={
          <>
            <p>
              The dance — record offline, queue, reconnect, drain — is the
              entire offline-first proposition. Once you have it working,
              your customers can record sales during 4-hour internet
              outages and never lose a transaction.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 3 — <strong>Frameless window UI</strong>. Custom titlebar,
        drag regions, the polish that makes a Wails app feel native.
      </p>
    </>
  )
}
