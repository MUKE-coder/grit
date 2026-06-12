import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'
import { Diagram } from '@/components/course/diagram'
import { PrereqLinks } from '@/components/course/prereq-links'

export default function Lesson() {
  return (
    <>
      <p>
        Gin gets the request to your handler. The handler calls a
        service. The service calls <strong>GORM</strong> to actually
        touch the database. This lesson is GORM — the ORM that turns
        Go structs into SQL.
      </p>

      <PrereqLinks prereqs={['golang']} />

      <h2>What is an ORM?</h2>
      <p>
        An ORM (Object-Relational Mapper) is a library that lets you
        work with database rows as if they were native language objects.
        Instead of:
      </p>
      <CodeBlock
        language="sql"
        code={`SELECT id, name, price_cents FROM products WHERE id = 42;`}
      />
      <p>You write:</p>
      <CodeBlock
        language="go"
        code={`var p Product
db.First(&p, 42)`}
      />
      <p>
        GORM generates the SQL, executes it, and unmarshals the row
        into your <code>Product</code> struct. You get type safety,
        less boilerplate, and code that reads like your domain.
      </p>

      <h2>The Model — your domain in a struct</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/product.go"
        code={`package models

import "time"

type Product struct {
    ID         uint      \`gorm:"primaryKey"\`
    Name       string    \`gorm:"not null;index"\`
    Slug       string    \`gorm:"uniqueIndex;not null"\`
    PriceCents int       \`gorm:"not null"\`
    InStock    bool      \`gorm:"default:true"\`
    UserID     uint      // foreign key — points to User
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt \`gorm:"index"\`  // soft delete
}`}
      />
      <p>
        Four things to notice:
      </p>
      <ul>
        <li>
          <strong>It&apos;s just a Go struct.</strong> No special base
          class, no &quot;active record&quot; magic. The TYPE is your
          domain; the tags tell GORM how to map it.
        </li>
        <li>
          <strong><code>gorm:</code> struct tags</strong> control SQL
          behaviour: primary key, indexes, uniqueness, defaults, not-null.
        </li>
        <li>
          <strong><code>CreatedAt</code> / <code>UpdatedAt</code></strong>{' '}
          are managed by GORM automatically — set on insert, updated on save.
        </li>
        <li>
          <strong><code>DeletedAt</code></strong> turns on soft delete.
          A row never disappears; it&apos;s marked deleted with a
          timestamp. Queries auto-filter it out.
        </li>
      </ul>

      <h2>From struct to table — AutoMigrate</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/db/db.go"
        code={`func Migrate(db *gorm.DB) error {
  return db.AutoMigrate(
    &models.User{},
    &models.Product{},
    &models.Order{},
  )
}`}
      />
      <p>
        AutoMigrate reads your struct definitions and either creates the
        table or adds missing columns. It will NOT drop columns or
        change types — by design. For schema removal you write a
        migration explicitly.
      </p>

      <h2>The five queries you&apos;ll write 95% of the time</h2>
      <CodeBlock
        language="go"
        code={`// 1. Find one by primary key
var p Product
db.First(&p, 42)                                   // SELECT * FROM products WHERE id = 42

// 2. Find one by a condition
db.Where("slug = ?", "fancy-thing").First(&p)      // ... WHERE slug = ?

// 3. List many
var ps []Product
db.Where("in_stock = ?", true).Limit(20).Find(&ps) // ... WHERE in_stock = true LIMIT 20

// 4. Create
db.Create(&Product{Name: "New", PriceCents: 999})  // INSERT ...

// 5. Update
db.Model(&p).Update("price_cents", 1299)           // UPDATE ... SET price_cents = ?`}
      />
      <p>
        Notice the <code>?</code> placeholders — GORM parameterises
        every value, so you can&apos;t SQL-inject by passing a string.
        Never build queries with <code>fmt.Sprintf</code> — always{' '}
        <code>?</code>.
      </p>

      <h2>Relationships — the &quot;hasMany&quot; / &quot;belongsTo&quot; pair</h2>
      <CodeBlock
        language="go"
        code={`type User struct {
  ID       uint
  Email    string
  Products []Product   // hasMany: a user can own many products
}

type Product struct {
  ID     uint
  Name   string
  UserID uint        // belongsTo: a product belongs to a user
  User   User        // optional — for preload
}`}
      />
      <p>
        Two-way: the &quot;has many&quot; side has a slice; the &quot;belongs to&quot;
        side has the foreign key (<code>UserID</code>) and optionally a
        struct field (<code>User</code>) for preloading.
      </p>

      <h3>Preloading — joining without joining</h3>
      <CodeBlock
        language="go"
        code={`var u User
db.Preload("Products").First(&u, 7)
// Executes TWO queries:
//   SELECT * FROM users WHERE id = 7
//   SELECT * FROM products WHERE user_id = 7
// Then GORM stitches them together into u.Products`}
      />
      <p>
        Preload runs a separate query rather than a JOIN. Two small
        queries usually beat one big JOIN with duplicated user data.
        Use <code>Joins</code> when you need a true SQL JOIN (e.g.,
        filtering by a related field).
      </p>

      <TipBox tone="warning">
        <strong>N+1 is the classic ORM trap.</strong> If you loop over
        users and call <code>db.First(&p, u.ProductID)</code> inside the
        loop, you&apos;ll fire one query per user. 1000 users = 1001
        queries. The fix is <code>Preload</code> or{' '}
        <code>Find(&products, ids)</code> outside the loop. We&apos;ll
        revisit this when you build your first list endpoint.
      </TipBox>

      <h2>How the DB connection reaches your service</h2>
      <Diagram label="DB injection in Grit" caption="GORM is opened once at startup and passed by pointer wherever a service needs it.">
{`   main.go
     │
     │  db, _ := gorm.Open(postgres.Open(dsn))   // 1. open once
     │
     ▼
   internal/db.Migrate(db)                       // 2. ensure tables exist
     │
     ▼
   services.New(db)                              // 3. inject into services
     │     (every service holds a *gorm.DB pointer)
     │
     ▼
   handlers.New(services)                        // 4. handlers call services
     │
     ▼
   routes.Register(r, handlers)                  // 5. routes call handlers`}
      </Diagram>
      <p>
        The <code>*gorm.DB</code> is a long-lived value with its own
        connection pool. Pass the pointer down; never open a new one
        per request.
      </p>

      <h2>The real list query — pagination + filter + count</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/services/product_service.go (typical)"
        code={`func (s *ProductService) List(ctx context.Context, q ListQuery) ([]Product, int64, error) {
  tx := s.db.WithContext(ctx).Model(&Product{})

  if q.Search != "" {
    tx = tx.Where("name ILIKE ?", "%"+q.Search+"%")
  }
  if q.InStockOnly {
    tx = tx.Where("in_stock = ?", true)
  }

  var total int64
  if err := tx.Count(&total).Error; err != nil { return nil, 0, err }

  var items []Product
  if err := tx.
    Order("created_at DESC").
    Offset((q.Page - 1) * q.PageSize).
    Limit(q.PageSize).
    Find(&items).Error; err != nil { return nil, 0, err }

  return items, total, nil
}`}
      />
      <p>
        Three observations:
      </p>
      <ul>
        <li>
          <strong><code>WithContext</code></strong> — passes the request
          context so a cancelled request also cancels the DB query.
          Always do this in services.
        </li>
        <li>
          <strong>Chaining</strong> — each <code>.Where</code> returns
          a new <code>*gorm.DB</code>, so you build up the query
          step-by-step. Reading bottom-up: ORDER, OFFSET, LIMIT, then
          execute.
        </li>
        <li>
          <strong>Count BEFORE limit.</strong> <code>Count(&total)</code>{' '}
          ignores LIMIT/OFFSET (GORM is smart) so you get the true row
          count for pagination.
        </li>
      </ul>

      <KnowledgeCheck
        question={<>You write <code>{`db.Where(fmt.Sprintf("name = '%s'", input)).First(&p)`}</code>. Why is this a bug?</>}
        choices={[
          {
            label: "Because GORM doesn't support Sprintf",
            feedback:
              "It does — but you should never use it for query values. The real issue is security.",
          },
          {
            label: "SQL injection — the user can pass `'; DROP TABLE products; --` and run arbitrary SQL",
            correct: true,
            feedback:
              "Right — string-concatenating user input into SQL is the classic injection vector. Always pass values via `?` placeholders so GORM parameterises them. You CAN safely Sprintf-build the QUERY STRUCTURE (column names) — never values.",
          },
          {
            label: "It's slower than a placeholder",
            feedback:
              "Marginal. The real reason is security; performance is the secondary issue.",
          },
          {
            label: "GORM caches Sprintf badly",
            feedback: "Not the issue. Security is.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Practice the five queries:</p>
            <ol>
              <li>
                In a fresh Grit project, ensure the User model exists.
                Use <code>gorm.io/playground</code> or a unit test (the
                generated <code>auth_test.go</code> has a SQLite test
                DB you can borrow).
              </li>
              <li>Insert 3 users with <code>db.Create</code>.</li>
              <li>Find one by email with <code>.Where(&quot;email = ?&quot;, ...)</code>.</li>
              <li>
                List all users created in the last hour using{' '}
                <code>.Where(&quot;created_at &gt; ?&quot;, time.Now().Add(-time.Hour))</code>.
              </li>
              <li>
                Update one user&apos;s name with{' '}
                <code>db.Model(&u).Update(&quot;name&quot;, ...)</code>.
              </li>
              <li>
                Soft-delete a user with <code>db.Delete(&u)</code> —
                then confirm a normal <code>Find</code> no longer
                returns them (GORM auto-filters soft-deleted rows).
              </li>
            </ol>
          </>
        }
        hint={
          <>
            For ad-hoc experimentation, you can use{' '}
            <code>gorm.io/driver/sqlite</code> with a temp file —{' '}
            <code>gorm.Open(sqlite.Open(&quot;:memory:&quot;), &amp;gorm.Config&#123;&#125;)</code>{' '}
            gives you a pure-Go in-memory DB. No Postgres needed.
          </>
        }
        solution={
          <>
            <p>
              You should be comfortable doing CRUD against any model
              without looking up GORM syntax. The next lesson uses this
              fluency to walk the full Handler → Service pattern.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — <strong>Handler → Service pattern</strong>. Now
        that you know how Gin gets the request to your code AND how
        GORM gets it to the DB, the pattern in between is what every
        Grit endpoint follows.
      </p>
    </>
  )
}
