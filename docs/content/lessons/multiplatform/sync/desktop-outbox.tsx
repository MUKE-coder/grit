import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'

export default function Lesson() {
  return (
    <>
      <p>
        Desktop users keep the app open for hours. They expect to
        bookmark something, close the lid, open it tomorrow on a
        flight, and have everything &quot;just work&quot;. The outbox
        pattern is how you deliver on that promise.
      </p>

      <h2>What an outbox is</h2>
      <Diagram label="Outbox flow" caption="Writes go to a local SQLite outbox first. A worker drains it to the server. Restart-safe, retry-safe.">
{`   User clicks "Bookmark" (offline)
              │
              ▼
   ┌────────────────────────┐
   │ Local SQLite           │
   │ ┌────────────────────┐ │
   │ │ bookmarks (cache)  │◄┼──── update for instant UI
   │ └────────────────────┘ │
   │ ┌────────────────────┐ │
   │ │ outbox             │◄┼──── enqueue: POST /api/bookmarks
   │ │ (id, op, payload,  │ │
   │ │  status, attempts) │ │
   │ └────────────────────┘ │
   └─────────┬──────────────┘
             │  (later, online)
             ▼
   Worker drains outbox →  HTTP POST → API
                                       │
                              200      │      4xx/5xx
                          mark done    │   retry / dead-letter`}
      </Diagram>

      <h2>The schema</h2>
      <CodeBlock
        language="go"
        filename="apps/desktop/internal/db/outbox.go"
        code={`package db

import "time"

type OutboxItem struct {
  ID         uint      \`gorm:"primaryKey"\`
  Op         string    // "create_bookmark", "delete_bookmark"
  Payload    string    // JSON
  Status     string    // "pending", "done", "failed"
  Attempts   int
  LastError  string
  CreatedAt  time.Time
  UpdatedAt  time.Time
}`}
      />

      <h2>Enqueue on every mutation</h2>
      <CodeBlock
        language="go"
        filename="apps/desktop/internal/services/bookmark_service.go"
        code={`func (s *BookmarkService) Create(productID uint) error {
  // 1. Apply locally so the UI is instant
  local := models.Bookmark{ProductID: productID, UserID: s.currentUserID()}
  if err := s.db.Create(&local).Error; err != nil {
    return err
  }
  // 2. Enqueue the server sync
  payload, _ := json.Marshal(map[string]any{"product_id": productID})
  return s.db.Create(&db.OutboxItem{
    Op: "create_bookmark", Payload: string(payload), Status: "pending",
  }).Error
}`}
      />
      <p>
        Critical detail: BOTH the local insert and the outbox enqueue
        happen in the same DB. Wrap them in a transaction so they
        commit atomically — either both, or neither. Otherwise you
        can &quot;create a bookmark&quot; that never reaches the
        server.
      </p>

      <h2>The drain worker</h2>
      <CodeBlock
        language="go"
        filename="apps/desktop/internal/sync/outbox_worker.go"
        code={`func (w *OutboxWorker) Tick(ctx context.Context) error {
  var items []db.OutboxItem
  w.db.Where("status = ?", "pending").
       Order("id ASC").Limit(20).Find(&items)

  for _, item := range items {
    err := w.send(ctx, item)
    if err == nil {
      w.db.Model(&item).Updates(map[string]any{
        "status": "done", "updated_at": time.Now(),
      })
      continue
    }
    item.Attempts++
    item.LastError = err.Error()
    if item.Attempts >= 5 {
      item.Status = "failed"  // dead letter
    }
    w.db.Save(&item)
  }
  return nil
}

func (w *OutboxWorker) Run(ctx context.Context) {
  t := time.NewTicker(10 * time.Second)
  defer t.Stop()
  for {
    select {
    case <-ctx.Done(): return
    case <-t.C: _ = w.Tick(ctx)
    }
  }
}`}
      />
      <p>
        Every 10 seconds the worker picks up 20 items, tries to send
        them, marks done or retries. Five strikes and it&apos;s dead-lettered.
      </p>

      <h2>Network awareness</h2>
      <p>
        The Tick should bail fast if offline — no point hammering
        the network stack. Wails apps can detect connectivity from
        the Go side:
      </p>
      <CodeBlock
        language="go"
        code={`func (w *OutboxWorker) isOnline() bool {
  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  req, _ := http.NewRequestWithContext(ctx, "HEAD", w.apiURL+"/healthz", nil)
  resp, err := http.DefaultClient.Do(req)
  if err != nil { return false }
  defer resp.Body.Close()
  return resp.StatusCode == 200
}`}
      />

      <TipBox tone="info">
        <strong>Idempotency keys.</strong> If the worker sends a
        POST but the network drops before the response comes back,
        it&apos;ll retry — and create a duplicate bookmark. Solution:
        generate a UUID on the desktop, send it as{' '}
        <code>Idempotency-Key</code>, have the API dedupe.
        Without this, retries cause duplicates. Don&apos;t skip.
      </TipBox>

      <h2>Showing queue state in the UI</h2>
      <CodeBlock
        language="tsx"
        code={`<div className="flex items-center gap-2 text-xs text-text-muted">
  <span className={cn('h-2 w-2 rounded-full', queueCount === 0 ? 'bg-emerald-500' : 'bg-amber-500')} />
  {queueCount === 0 ? 'All synced' : queueCount + ' pending sync'}
</div>`}
      />
      <p>
        Tiny status indicator in the status bar. Users want to know
        if their work is safely sent. A green dot is reassuring.
      </p>

      <h2>What about reads?</h2>
      <p>
        Reads come from the local SQLite cache directly — no outbox
        needed. A separate puller pings the API every few minutes
        and updates the local cache. That&apos;s the next lesson&apos;s
        topic when we talk about conflicts: what happens when the
        server has data the local cache doesn&apos;t (or vice versa)?
      </p>

      <KnowledgeCheck
        question="The user creates a bookmark while offline. The desktop crashes before the worker syncs. The user reopens the app online. What happens?"
        choices={[
          {
            label: 'The bookmark is lost — it was only in memory',
            feedback:
              'Lost only if you don\'t persist the local DB or outbox. Both are SQLite tables on disk → survive crashes.',
          },
          {
            label: 'The local bookmark is preserved AND the outbox entry is still pending, so the worker syncs it on first tick',
            correct: true,
            feedback:
              "Right — both rows are persisted in SQLite at commit. Crash-safety is the WHOLE point of using a DB-backed outbox vs an in-memory queue.",
          },
          {
            label: 'The user has to re-create the bookmark',
            feedback: 'Defeats the purpose — outbox is exactly to prevent this.',
          },
          {
            label: 'The bookmark syncs but is marked as "recovered"',
            feedback:
              'No special state — it’s just a pending outbox entry that gets drained normally.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Build the outbox end to end:</p>
            <ol>
              <li>Add the OutboxItem model + auto-migrate.</li>
              <li>Wrap your local Create/Delete bookmark logic with the outbox enqueue.</li>
              <li>Wire the worker to run via <code>app.OnStartup</code>.</li>
              <li>Add an idempotency key (UUID) to every payload.</li>
              <li>
                Test: disable network, bookmark 3 products. Force-kill
                the app, re-launch online — confirm all 3 bookmarks
                show up server-side within 15s.
              </li>
              <li>
                Stress test: disable network, bookmark something, KILL
                THE PROCESS (not graceful shutdown). Reopen online —
                bookmark still syncs.
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For network simulation on Wails, you can flip your
            machine&apos;s wifi off and on. Or use a tool like
            <code> Network Link Conditioner</code> on Mac for
            controlled latency / packet loss.
          </>
        }
        solution={
          <>
            <p>
              Outbox count should drain to 0 within a tick or two of
              reconnection. The user sees their bookmarks appear on
              web within seconds without any manual intervention.
            </p>
            <p>
              This is the desktop feature that earns its keep over
              a web app. &quot;Offline-first&quot; is real, not a
              marketing line.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Conflict resolution</strong>. Mobile
        edits a bookmark while desktop has the old version. Whose
        change wins? We&apos;ll set the policy.
      </p>
    </>
  )
}
