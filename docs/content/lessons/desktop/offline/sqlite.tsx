import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Local SQLite is the heart of an offline-first app. Your desktop
        binary writes to a local file; the file survives launches; the
        sync engine pushes those rows to your server later. This lesson
        covers the SQLite + GORM setup.
      </p>

      <h2>Where the file lives</h2>
      <CodeBlock
        language="text"
        code={`Windows:  %APPDATA%\\field-pos\\local.db
macOS:    ~/Library/Application Support/field-pos/local.db
Linux:    ~/.config/field-pos/local.db`}
      />
      <p>
        Grit uses{' '}
        <code>github.com/adrg/xdg</code> to resolve the right per-OS path.
        The file persists across launches and survives app updates.
      </p>

      <h2>The connection</h2>
      <CodeBlock
        language="go"
        filename="internal/db/db.go"
        code={`import (
    "github.com/glebarez/sqlite"
    "gorm.io/gorm"
)

func Open(path string) (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open(path+"?_pragma=journal_mode(WAL)"), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    db.AutoMigrate(&User{}, &Sale{}, &Product{})
    return db, nil
}`}
      />
      <p>
        Note: <strong>pure-Go SQLite</strong> (<code>glebarez/sqlite</code>),
        not CGO. Wails builds work on Windows without a C toolchain. WAL
        mode prevents readers from blocking writers — critical for a
        responsive UI that polls.
      </p>

      <h2>Models</h2>
      <CodeBlock
        language="go"
        filename="internal/models/sale.go"
        code={`type Sale struct {
    ID         string    \`gorm:"primaryKey"\`              // UUID — generated client-side
    LocalSeq   int64     \`gorm:"autoIncrement;index"\`     // monotonic local ordering
    UserID     string    \`gorm:"not null"\`
    Total      decimal.Decimal
    Status     string    \`gorm:"default:'completed'"\`
    ServerID   *string   \`gorm:"index"\`                   // nil until synced
    CreatedAt  time.Time
}`}
      />
      <p>Two ID columns matter:</p>
      <ul>
        <li>
          <strong><code>ID</code></strong> — UUID generated on the device.
          Stable forever; survives sync. Server adopts this same ID.
        </li>
        <li>
          <strong><code>ServerID</code></strong> — only set after the row
          syncs successfully. Useful for &quot;synced ✓&quot; vs &quot;pending&quot; UI.
        </li>
      </ul>

      <TipBox tone="success">
        <strong>UUIDs not auto-increment.</strong> If the same desktop is
        used by two cashiers offline, both can write rows without ID
        collisions when they later sync. Auto-increment IDs would
        collide. UUIDs are the universal sync-safe choice.
      </TipBox>

      <h2>Reading + writing from Go (called by React)</h2>
      <CodeBlock
        language="go"
        filename="app.go"
        code={`func (a *App) RecordSale(items []SaleItem) (*Sale, error) {
    sale := &Sale{
        ID:     uuid.NewString(),
        UserID: a.currentUserID,
        Total:  computeTotal(items),
        Items:  items,
    }
    if err := a.db.Create(sale).Error; err != nil {
        return nil, fmt.Errorf("save sale: %w", err)
    }
    // Enqueue into outbox (lesson 2.2)
    a.outbox.Enqueue(sale)
    return sale, nil
}

func (a *App) ListSales(limit int) ([]Sale, error) {
    var sales []Sale
    a.db.Order("local_seq DESC").Limit(limit).Find(&sales)
    return sales, nil
}`}
      />
      <p>
        Both bound to React. React calls{' '}
        <code>RecordSale([...items])</code>, gets the saved row back, and
        re-renders the list. No network involved — pure local DB.
      </p>

      <h2>Backups + corruption recovery</h2>
      <p>
        SQLite is incredibly reliable, but for an offline POS that&apos;s a
        cashier&apos;s only proof of today&apos;s sales, ship a daily backup:
      </p>
      <CodeBlock
        language="go"
        code={`// On app startup, copy local.db to local.db.YYYYMMDD.bak
// Keep last 7 days; rotate older ones out`}
      />
      <p>
        Belt-and-braces. Disk corruption is rare but not impossible. A
        7-day rolling backup costs ~kilobytes and saves a customer&apos;s
        Tuesday when their Wednesday goes wrong.
      </p>

      <KnowledgeCheck
        question="A user takes their desktop app offline for 3 days, records 200 sales, then reconnects. Two cashiers used the same device simultaneously. What collides?"
        choices={[
          {
            label: 'Sale IDs — they\'re client-generated',
            feedback:
              "Wrong if the IDs are UUIDs (Grit's default). UUID collisions are vanishingly rare. If the project uses auto-increment IDs, two cashiers would create sale #501 each — that's why Grit uses UUIDs.",
          },
          {
            label: 'Nothing — UUIDs are globally unique by design',
            correct: true,
            feedback:
              "Right — that's exactly why UUIDs are the default for sync-aware models. Server receives 200 sales with unique IDs and inserts them all without conflict.",
          },
          {
            label: 'LocalSeq — both cashiers would have seq 501',
            feedback:
              "LocalSeq is auto-increment per device. It's only used for local ordering, not as a globally-unique key. So no collision.",
          },
          {
            label: 'The server rejects the batch as duplicates',
            feedback:
              "Only if the server itself can't tell the rows apart. With UUIDs as the PK, the server sees 200 unique rows and inserts them all.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Wire local SQLite + a write:</p>
            <ol>
              <li>
                Confirm your scaffolded project has <code>internal/db/</code>{' '}
                with a <code>Connect()</code> function.
              </li>
              <li>
                Add a <code>Note</code> model + an <code>AddNote(text string)</code>{' '}
                method on App.
              </li>
              <li>
                Wire a React form that calls <code>AddNote</code>.
              </li>
              <li>
                Add a few notes. Quit the app. Reopen — your notes should
                persist.
              </li>
            </ol>
            <p>Paste a screenshot of the notes list after a restart in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If the notes disappear on restart, the DB path is probably
            wrong (writing to <code>./</code> which gets nuked). Use{' '}
            <code>xdg.DataFile()</code> from <code>github.com/adrg/xdg</code>{' '}
            to resolve a stable per-user path.
          </>
        }
        solution={
          <>
            <p>
              You should see the same three notes after relaunching. That
              proves the SQLite file is at a stable persistent path —
              survives kill + relaunch.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Local writes are saved. Next lesson — the{' '}
        <strong>outbox pattern</strong> for sending them to the server
        when online.
      </p>
    </>
  )
}
