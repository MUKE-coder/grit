import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Real domains have shapes — a customer has many invoices, an invoice
        belongs to a customer, a tag belongs to many products. GORM has one
        struct tag per relation kind. This lesson covers all three.
      </p>

      <h2>1. BelongsTo (the foreign key sits on this row)</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/invoice.go"
        code={`type Invoice struct {
    ID         uuid.UUID \`gorm:"type:uuid;primaryKey" json:"id"\`
    CustomerID uuid.UUID \`gorm:"type:uuid;not null;index" json:"customer_id"\`
    Customer   Customer  \`gorm:"foreignKey:CustomerID" json:"customer,omitempty"\`
    Total      decimal.Decimal \`gorm:"type:decimal(10,2);not null" json:"total"\`
    CreatedAt  time.Time \`json:"created_at"\`
}`}
      />
      <p>
        Two pieces — the FK column (<code>CustomerID</code>) and the relation
        field (<code>Customer Customer</code>). Use <code>Preload</code> to
        eager-load:
      </p>
      <CodeBlock
        language="go"
        code={`var invoice Invoice
db.Preload("Customer").First(&invoice, "id = ?", id)
// invoice.Customer.Name is now populated`}
      />

      <h2>2. HasMany (the other table owns the FK)</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/customer.go"
        code={`type Customer struct {
    ID       uuid.UUID  \`gorm:"type:uuid;primaryKey" json:"id"\`
    Name     string     \`gorm:"not null" json:"name"\`
    Invoices []Invoice  \`gorm:"foreignKey:CustomerID" json:"invoices,omitempty"\`
}`}
      />
      <p>
        <code>Invoices []Invoice</code> means &quot;this customer has many
        invoices&quot;. The FK lives on the Invoice row. Eager-load the same
        way:
      </p>
      <CodeBlock
        language="go"
        code={`var c Customer
db.Preload("Invoices").First(&c, "id = ?", id)
// c.Invoices is a slice`}
      />

      <h2>3. Many-to-Many (a join table)</h2>
      <CodeBlock
        language="go"
        code={`type Product struct {
    ID   uuid.UUID \`gorm:"type:uuid;primaryKey"\`
    Tags []Tag     \`gorm:"many2many:product_tags;" json:"tags,omitempty"\`
}
type Tag struct {
    ID   uuid.UUID \`gorm:"type:uuid;primaryKey"\`
    Name string    \`gorm:"uniqueIndex;not null"\`
}`}
      />
      <p>
        GORM creates the <code>product_tags</code> join table for you. Add a
        tag:
      </p>
      <CodeBlock
        language="go"
        code={`db.Model(&product).Association("Tags").Append(&tag)`}
      />

      <TipBox tone="success">
        <strong>Naming convention:</strong> the FK is always{' '}
        <code>{`<RelationName>ID`}</code> in the model and <code>{`<relation_name>_id`}</code>{' '}
        in the DB (snake_case). GORM follows this automatically; you only
        need <code>foreignKey:</code> when overriding.
      </TipBox>

      <h2>OnDelete behaviour</h2>
      <p>
        What happens when you delete a customer with invoices? Three
        choices:
      </p>
      <CodeBlock
        language="go"
        code={`// Cascade — delete invoices too
Customer Customer \`gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"\`

// Set null — keep invoices but null out CustomerID (need nullable FK)
Customer Customer \`gorm:"foreignKey:CustomerID;constraint:OnDelete:SET NULL"\`

// Default: restrict — refuse to delete the customer if invoices exist`}
      />
      <p>
        Production-grade Grit projects tend to use <code>RESTRICT</code> +
        soft-delete the parent. Cascade is irreversible — pick it
        consciously.
      </p>

      <KnowledgeCheck
        question={<>You declared <code>{`Customer Customer \`gorm:"foreignKey:CustomerID"\``}</code> but{' '}<code>{`db.Preload("Customer").First(...)`}</code> returns Customer with all-zero fields. What&apos;s most likely wrong?</>}
        choices={[
          {
            label: 'You forgot to declare CustomerID — Grit can\'t find the FK column.',
            correct: true,
            feedback:
              "Right — `Customer Customer` is the relation field, but you also need `CustomerID uuid.UUID` (the actual column). Without the FK column, the relation has nothing to join on.",
          },
          {
            label: 'You need to call `db.Joins(\"Customer\")` instead of Preload.',
            feedback:
              "Different mechanism — Joins INNER JOINs, Preload does a follow-up query. Both work if the FK exists; if it doesn't, both return zero values.",
          },
          {
            label: 'GORM doesn\'t support Preload on uuid PKs.',
            feedback:
              "Wrong — GORM handles uuid PKs fine. The constraint is having the FK column declared.",
          },
          {
            label: 'You need to add `gorm:\"references:ID\"` to the relation tag.',
            feedback:
              "Optional, not required. GORM defaults to `references:ID`; you only need it when the parent's PK isn't named ID.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Add the invoice domain to your bench-api:
            </p>
            <ol>
              <li>
                <code>Customer</code> (HasMany Invoices)
              </li>
              <li>
                <code>Invoice</code> (BelongsTo Customer, HasMany LineItems)
              </li>
              <li>
                <code>LineItem</code> (BelongsTo Invoice)
              </li>
            </ol>
            <p>
              Wire each with the right FK + relation fields. Run{' '}
              <code>grit migrate</code> and confirm via GORM Studio that the
              FK columns exist. Paste the schema in <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            Each model needs (a) the relation field (<code>Invoices []Invoice</code>{' '}
            on Customer) AND (b) the FK column on the child{' '}
            (<code>CustomerID uuid.UUID</code> on Invoice). Forget one and
            preload returns zero values.
          </>
        }
        solution={
          <>
            <p>Three models, six fields each (PK, FK, parent, plus
            domain-specific):</p>
            <CodeBlock
              language="go"
              code={`type Customer struct {
    ID       uuid.UUID
    Name     string
    Invoices []Invoice \`gorm:"foreignKey:CustomerID"\`
}

type Invoice struct {
    ID         uuid.UUID
    CustomerID uuid.UUID
    Customer   Customer   \`gorm:"foreignKey:CustomerID"\`
    LineItems  []LineItem \`gorm:"foreignKey:InvoiceID"\`
}

type LineItem struct {
    ID        uuid.UUID
    InvoiceID uuid.UUID
    Invoice   Invoice \`gorm:"foreignKey:InvoiceID"\`
    Quantity  int
    Price     decimal.Decimal
}`}
            />
            <p>
              That&apos;s chapter 2&apos;s assignment done. AutoMigrate creates the
              tables + FK columns + indexes on boot.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Models compile, AutoMigrate runs. But should you use AutoMigrate
        forever? Last lesson of this chapter: when AutoMigrate is fine and
        when to switch to explicit migrations.
      </p>
    </>
  )
}
