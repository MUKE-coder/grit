import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        GORM is the ORM that ships with Grit. It speaks Postgres and SQLite
        out of the box; you write Go structs, it writes SQL. This lesson is
        the 7-minute crash course you need to read every Grit model.
      </p>

      <h2>A model is a Go struct + tags</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/customer.go"
        code={`package models

import (
    "time"
    "github.com/google/uuid"
)

type Customer struct {
    ID        uuid.UUID \`gorm:"type:uuid;primaryKey" json:"id"\`
    Email     string    \`gorm:"uniqueIndex;not null" json:"email"\`
    Name      string    \`gorm:"not null" json:"name"\`
    CreatedAt time.Time \`json:"created_at"\`
    UpdatedAt time.Time \`json:"updated_at"\`
}`}
      />

      <h2>The tags that matter</h2>
      <ul>
        <li>
          <code>gorm:&quot;primaryKey&quot;</code> — this field is the PK
        </li>
        <li>
          <code>gorm:&quot;type:uuid&quot;</code> — column type (use uuid for IDs;
          Grit defaults to UUID PKs)
        </li>
        <li>
          <code>gorm:&quot;not null&quot;</code> — DB-level NOT NULL constraint
        </li>
        <li>
          <code>gorm:&quot;uniqueIndex&quot;</code> — unique index on this column
        </li>
        <li>
          <code>gorm:&quot;default:value&quot;</code> — column default
        </li>
        <li>
          <code>gorm:&quot;-&quot;</code> — skip this field entirely (never write
          to DB)
        </li>
        <li>
          <code>json:&quot;snake_name&quot;</code> — how it appears in the API
          response
        </li>
      </ul>

      <h2>UUID, not auto-increment</h2>
      <p>
        Grit uses UUIDs for primary keys, not <code>uint</code> auto-increment.
        Two reasons: (1) IDOR defence — attackers can&apos;t guess the next ID;
        (2) globally unique — you can generate them client-side and merge
        across DBs without conflict.
      </p>
      <CodeBlock
        language="go"
        code={`// Auto-generated in a BeforeCreate hook
func (c *Customer) BeforeCreate(tx *gorm.DB) error {
    if c.ID == uuid.Nil {
        c.ID = uuid.New()
    }
    return nil
}`}
      />
      <p>
        You don&apos;t need to write this hook yourself if you use{' '}
        <code>grit generate resource</code> — it&apos;s already in the template.
      </p>

      <h2>CreatedAt + UpdatedAt — free</h2>
      <p>
        GORM auto-sets <code>CreatedAt</code> on insert and{' '}
        <code>UpdatedAt</code> on every save. No work on your part — just
        declare them with the right names and types.
      </p>

      <TipBox tone="info">
        Soft deletes — adding{' '}
        <code>DeletedAt gorm.DeletedAt {`\`gorm:"index"\``}</code> turns on
        soft delete. Calling <code>db.Delete(&amp;customer)</code> sets the
        column; queries automatically exclude soft-deleted rows. Useful for
        audit-able domains.
      </TipBox>

      <h2>Reading + writing</h2>
      <CodeBlock
        language="go"
        code={`// Create
c := &Customer{Email: "alex@example.com", Name: "Alex"}
db.Create(c)

// Read by ID
var c Customer
db.First(&c, "id = ?", id)

// Update one field
db.Model(&c).Update("name", "Alex Doe")

// Query with WHERE
var customers []Customer
db.Where("email LIKE ?", "%@example.com").Find(&customers)`}
      />
      <p>
        Always pass arguments as the second+ parameter — never concatenate
        strings (covered as SQLi defence in the Concepts course).
      </p>

      <KnowledgeCheck
        question={<>A teammate writes <code>{`db.Raw("SELECT * FROM users WHERE email = '" + email + "'").Scan(&u)`}</code>. What&apos;s wrong?</>}
        choices={[
          {
            label: 'Nothing — Raw queries are fine.',
            feedback:
              "Wrong — string concatenation in queries IS SQL injection. An attacker who controls `email` can rewrite the query.",
          },
          {
            label: 'Should be `db.Raw(\"SELECT * FROM users WHERE email = ?\", email).Scan(&u)` with ? and a separate argument.',
            correct: true,
            feedback:
              "Right — GORM passes the argument as a parameterized query under the hood. SQLi is closed at the driver level. This was a core Defender's Handbook lesson too.",
          },
          {
            label: 'Should be `db.First(&u, \"email = \" + email)`.',
            feedback:
              "Wrong direction — still concatenating. The fix is to use ? + a separate argument.",
          },
          {
            label: 'Switch to ORM-style db.Where; raw queries should be avoided.',
            feedback:
              "Avoidance works, but isn't the precise issue. The real fix is parameterized queries — works for both Raw and Where.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Add a Customer model to your bench-api. Then write three lines
              of code in <code>main.go</code> (or a temporary debug handler)
              that create one and read it back.
            </p>
            <CodeBlock
              language="go"
              code={`db.AutoMigrate(&Customer{})
c := &Customer{Email: "alex@example.com", Name: "Alex"}
db.Create(c)
var found Customer
db.First(&found, "email = ?", "alex@example.com")
log.Println("found:", found.ID, found.Name)`}
            />
            <p>Paste the log output in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            Add the model in <code>internal/models/customer.go</code>. Wire
            <code> AutoMigrate(&amp;models.Customer{`{}`})</code> in the
            database setup so the table gets created on boot.
          </>
        }
        solution={
          <>
            <p>You should see something like:</p>
            <CodeBlock
              language="text"
              code={`found: 9b4d... Alex`}
            />
            <p>
              UUID + name printed. GORM did the table creation, the insert,
              and the read for you.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Single tables are easy. Real apps have <em>relations</em> — Customer
        has many Invoices; an Invoice belongs to a Customer. Next lesson:
        the three relation kinds and how Grit defines each.
      </p>
    </>
  )
}
