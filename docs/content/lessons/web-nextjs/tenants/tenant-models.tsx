import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Multi-tenant SaaS = many isolated customers in one DB. Each tenant
        sees only their own data. Done right, it&apos;s invisible; done wrong,
        you leak customer A&apos;s data to customer B. This lesson covers the
        Grit pattern.
      </p>

      <h2>The model</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/tenant.go"
        code={`type Tenant struct {
    ID        uuid.UUID \`gorm:"type:uuid;primaryKey" json:"id"\`
    Name      string    \`gorm:"not null" json:"name"\`
    Slug      string    \`gorm:"uniqueIndex;not null" json:"slug"\`
    PlanID    string    \`json:"plan_id"\`
    CreatedAt time.Time \`json:"created_at"\`
}

// Every row in a tenant-scoped table gets a TenantID
type Invoice struct {
    ID         uuid.UUID \`gorm:"type:uuid;primaryKey" json:"id"\`
    TenantID   uuid.UUID \`gorm:"type:uuid;not null;index" json:"tenant_id"\`  // ← the scope
    CustomerID uuid.UUID \`gorm:"type:uuid;not null" json:"customer_id"\`
    Total      decimal.Decimal
    CreatedAt  time.Time
}`}
      />
      <p>
        Two important pieces: a <code>tenants</code> table, and a{' '}
        <code>TenantID</code> column on every tenant-scoped row (orders,
        invoices, products — anything that &quot;belongs to&quot; a tenant).
      </p>

      <h2>The middleware that scopes queries</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/middleware/tenant.go"
        code={`func TenantScope(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, _ := c.Get("user_id")
        var user models.User
        db.First(&user, "id = ?", userID)
        c.Set("tenant_id", user.TenantID.String())
        c.Next()
    }
}`}
      />
      <p>
        Every protected request gets <code>tenant_id</code> on its context
        from the authenticated user&apos;s tenant.
      </p>

      <h2>Using it in handlers</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/invoice.go"
        code={`func (h *InvoiceHandler) List(c *gin.Context) {
    tenantID := c.GetString("tenant_id")

    var invoices []models.Invoice
    h.db.Where("tenant_id = ?", tenantID).Find(&invoices)

    c.JSON(200, gin.H{"data": invoices})
}

func (h *InvoiceHandler) Create(c *gin.Context) {
    tenantID, _ := uuid.Parse(c.GetString("tenant_id"))

    var inv models.Invoice
    c.ShouldBindJSON(&inv)
    inv.TenantID = tenantID  // ← force the tenant, never trust input

    h.db.Create(&inv)
    c.JSON(201, gin.H{"data": inv})
}`}
      />
      <p>
        Every read filters by tenant; every write forces the tenant from
        context, never from the request body. <strong>Never accept{' '}
        <code>tenant_id</code> in the request body</strong> — that&apos;d let
        the user write to another tenant.
      </p>

      <TipBox tone="danger">
        <strong>The #1 multi-tenant bug:</strong> forgetting{' '}
        <code>WHERE tenant_id = ?</code> on a query. The query returns
        EVERY tenant&apos;s rows. Code review for it; add a lint rule that
        flags <code>db.Find()</code> without a Where clause in tenant-
        scoped handlers.
      </TipBox>

      <h2>Row-level security as belt-and-braces</h2>
      <p>
        Postgres supports row-level security (RLS) — DB-level filtering by
        a session var. If a handler forgets the filter, the DB itself
        rejects. Set up:
      </p>
      <CodeBlock
        language="sql"
        code={`ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON invoices
  USING (tenant_id = current_setting('app.current_tenant')::uuid);

-- Then in the middleware:
db.Exec("SET LOCAL app.current_tenant = ?", tenantID)`}
      />
      <p>
        Now even a buggy handler can&apos;t leak cross-tenant data. Worth the
        ~30 minutes to set up for high-stakes SaaS.
      </p>

      <h2>The Tenant model on the web side</h2>
      <p>
        On signup, create a tenant + assign the user as owner:
      </p>
      <CodeBlock
        language="go"
        code={`tenant := &Tenant{Name: input.CompanyName, Slug: slug.Make(input.CompanyName)}
db.Create(tenant)

user := &User{
    Email:    input.Email,
    TenantID: tenant.ID,
    Role:     "owner",
}
db.Create(user)`}
      />

      <KnowledgeCheck
        question="A handler does `db.Find(&invoices)` without filtering by tenant. What's the impact?"
        choices={[
          {
            label: 'Every authenticated user sees every tenant\'s invoices — catastrophic data leak',
            correct: true,
            feedback:
              "Right — the worst kind of multi-tenant bug. Pre-launch this gets you into the news. Always WHERE tenant_id, always force the value from context, never from input.",
          },
          {
            label: 'Slow query — but no data leak',
            feedback:
              "Slower AND a data leak. The latter is the showstopper.",
          },
          {
            label: 'GORM refuses to run the query',
            feedback:
              "Wrong — GORM happily runs unscoped queries. The DB doesn't know it's supposed to filter.",
          },
          {
            label: 'Postgres returns an empty list as default',
            feedback:
              "Wrong — without RLS, Postgres returns everything.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Make one of your existing models tenant-scoped:</p>
            <ol>
              <li>
                Add <code>TenantID uuid.UUID</code> to the Product model.
              </li>
              <li>Run <code>grit migrate</code>.</li>
              <li>
                Update the Product handler&apos;s List + Create to filter by /
                force tenant.
              </li>
              <li>
                Create two test tenants. Make sure tenant A&apos;s admin sees
                only tenant A&apos;s products.
              </li>
            </ol>
            <p>Paste the WHERE clause + the force-from-context line in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            For testing two tenants, you can manually set{' '}
            <code>tenant_id</code> on two seeded users via{' '}
            <code>UPDATE users SET tenant_id = ...</code> in GORM Studio.
          </>
        }
        solution={
          <>
            <CodeBlock
              language="go"
              code={`// List
h.db.Where("tenant_id = ?", c.GetString("tenant_id")).Find(&products)

// Create
product.TenantID, _ = uuid.Parse(c.GetString("tenant_id"))
h.db.Create(&product)`}
            />
            <p>Two lines, multi-tenant correctness.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Tenants done. Next — <strong>role gates in the UI</strong>. Don&apos;t
        just block on the server; hide buttons the user can&apos;t use.
      </p>
    </>
  )
}
