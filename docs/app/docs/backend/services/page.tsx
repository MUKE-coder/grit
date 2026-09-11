import Link from 'next/link'
import { ArrowLeft, ArrowRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { LaneFlow } from '@/components/lane-flow'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/services')

export default function ServicesPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />
      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend (Go API)</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Services
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                A service owns every database read and write for its resource. The one{' '}
                <code>grit generate resource</code> writes is complete: list, export, get, create,
                update, patch, delete and bulk, with the resource&apos;s whitelists, ownership rules
                and transactions. The <Link href="/docs/backend/handlers">handler</Link> calls it, and
                so can a job, a command or a test.
              </p>
            </div>

            <div className="prose-grit">
              {/* ── Handler or service ─────────────────────────────── */}
              <h2 id="when-to-use">Handler or Service?</h2>
              <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                <table className="w-full text-sm">
                  <thead>
                    <tr className="border-b border-border/30 bg-accent/20">
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Concern</th>
                      <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Where</th>
                    </tr>
                  </thead>
                  <tbody className="text-muted-foreground">
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Binding and validating the request body</td>
                      <td className="px-4 py-2.5">Handler</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Status codes, headers, rendering CSV or PDF</td>
                      <td className="px-4 py-2.5">Handler</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Any query, however simple</td>
                      <td className="px-4 py-2.5">Service</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Several writes that must land together</td>
                      <td className="px-4 py-2.5">Service, in a transaction</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Who may see or change a row</td>
                      <td className="px-4 py-2.5">Service, from the caller on the context</td>
                    </tr>
                    <tr className="border-b border-border/20">
                      <td className="px-4 py-2.5">Business rules beyond binding tags</td>
                      <td className="px-4 py-2.5">Service, returning <code>respond.Rule</code></td>
                    </tr>
                    <tr>
                      <td className="px-4 py-2.5">External calls: email, storage, AI</td>
                      <td className="px-4 py-2.5">A dedicated service</td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <LaneFlow
                id="svc-layers"
                lanes={['Callers', 'Business layer', 'Data & external']}
                nodes={[
                  { id: 'handler', lane: 0, row: 0, title: 'Handler', sub: 'HTTP request', tone: 'cyan' },
                  { id: 'job', lane: 0, row: 1, title: 'Job · Command', sub: 'no request', tone: 'cyan' },
                  { id: 'test', lane: 0, row: 2, title: 'Test', sub: 'no request', tone: 'cyan' },
                  { id: 'service', lane: 1, row: 1, title: 'Service', sub: 'rules · transactions', tone: 'primary' },
                  { id: 'gorm', lane: 2, row: 0, title: 'GORM', sub: 'queries', tone: 'green' },
                  { id: 'ext', lane: 2, row: 1, title: 'Email · Storage', sub: 'external APIs', tone: 'violet' },
                  { id: 'jobs', lane: 2, row: 2, title: 'AI · Jobs', sub: 'async work', tone: 'amber' },
                ]}
                edges={[
                  { from: 'handler', to: 'service', label: 'ctx', tone: 'primary' },
                  { from: 'job', to: 'service', label: 'ctx', tone: 'primary' },
                  { from: 'test', to: 'service', label: 'ctx', tone: 'primary' },
                  { from: 'service', to: 'gorm', label: 'query', tone: 'green' },
                  { from: 'service', to: 'ext', label: 'call', tone: 'violet' },
                  { from: 'service', to: 'jobs', label: 'enqueue', tone: 'amber' },
                ]}
                legend={[
                  { tone: 'cyan', label: 'Callers' },
                  { tone: 'primary', label: 'Service (logic)' },
                  { tone: 'green', label: 'Data & external' },
                ]}
                caption="Every caller gets the same rules, because the rules are in the service"
              />

              {/* ── The generated service ─────────────────────────────── */}
              <h2 id="service-pattern">The Generated Service</h2>
              <p>
                Services live in <code>apps/api/internal/services/</code>, one file per resource. The
                top of a generated one:
              </p>
              <CodeBlock language="go" filename="apps/api/internal/services/post.go" code={`// PostService owns every database read and write for posts.
type PostService struct {
    DB *gorm.DB
}

// postListConfig is what a client may search, sort and filter posts
// by. Whitelisted, because each name ends up in SQL.
var postListConfig = paginate.Config{
    Searchable: []string{"title", "body"},
    Sortable:   map[string]bool{"id": true, "created_at": true, "title": true, "published": true},
    Filterable: map[string]bool{"id": true, "title": true, "published": true},
}

// writablePost is every column Patch and Bulk may write. id, the
// timestamps and the version are the framework's, and are dropped.
var writablePost = map[string]bool{
    "title":     true,
    "body":      true,
    "published": true,
}

// db binds the database to ctx, so whatever a middleware put there reaches
// GORM's callbacks: the multitenant plugin scopes by it.
func (s *PostService) db(ctx context.Context) *gorm.DB {
    return s.DB.WithContext(ctx)
}`} />
              <p>And the methods it comes with:</p>
              <CodeBlock language="go" filename="services/post.go (methods)" code={`List(ctx, p paginate.Params, archived string) (paginate.Result[models.Post], error)
Export(ctx, search string, each func(rows []models.Post) error) error
GetByID(ctx, id string) (*models.Post, error)
Create(ctx, item *models.Post) error
Update(ctx, id string, updates map[string]interface{}, pre *concurrency.Precondition) (*models.Post, error)
Patch(ctx, id string, body map[string]interface{}, pre *concurrency.Precondition) (*models.Post, map[string]interface{}, error)
Delete(ctx, id string) (*models.Post, error)
Bulk(ctx, action string, ids []string, patch map[string]interface{}) (PostBulkResult, error)`} />
              <p>
                <code>grit generate field</code> adds a new column to the whitelists as well as the
                model and the handler, so a field added later can be filtered and patched like one
                generated with the resource.
              </p>
              <p>The flags add their own queries, in the same place:</p>
              <ul>
                <li>
                  <strong><code>--public</code></strong> adds <code>ListPublic</code>,{' '}
                  <code>GetPublic</code>, <code>RelatedPublic</code> when the resource has a parent,
                  and <code>PublicSubtreeIDs</code> and <code>PublicTreeRows</code> for a tree. What
                  a caller may filter by stays in the public handler, where you edit it, and is
                  handed to <code>ListPublic</code>.
                </li>
                <li>
                  <strong>The CSV import</strong> lives in <code>services/post_import.go</code>:{' '}
                  <code>StartImport</code> records the job, and <code>ImportCSV</code> reads the rows,
                  resolves their relations and writes them in batches. It runs after the response,
                  so the handler gives it <code>context.WithoutCancel(h.ctx(c))</code>: the caller
                  and the organization, without the cancellation.
                </li>
                <li>
                  <strong><code>--tree</code></strong> has its hierarchy on a{' '}
                  <code>PostTreeService</code> in <code>services/post_tree.go</code>, whose methods
                  take the context too. <code>Move</code> with a nil parent keeps the one the node
                  has, which is how a reorder is sent.
                </li>
              </ul>

              {/* ── The context ─────────────────────────────── */}
              <h2 id="context">The Context Carries the Caller</h2>
              <p>
                Every method takes a <code>context.Context</code> first. It carries who is asking,
                which an owned resource (<code>--owned-by</code>) scopes its rows by, as well as the
                organization the multitenant plugin resolved and the cancellation when a client goes
                away.
              </p>
              <CodeBlock language="go" filename="three ways to call a service" code={`// In a handler. The generated ctx helper does exactly this.
ctx := authz.WithActor(c.Request.Context(), authz.ActorOf(c))

// In a job or a command, acting for the system: every row, like ADMIN.
ctx := authz.AsSystem(context.Background())

// Acting for one user, so owned rows are scoped to them.
ctx := authz.WithActor(context.Background(), authz.Actor{UserID: userID})`} />
              <p>
                On an owned resource a context with no caller on it matches nothing. That is on
                purpose: a job that forgot to say who it acts for gets an empty list, not every
                user&apos;s rows. Somebody else&apos;s row comes back as{' '}
                <code>gorm.ErrRecordNotFound</code>, so a wrong guess at an id cannot be told from a
                right one.
              </p>
              <CodeBlock language="go" filename="jobs/overdue.go" code={`svc := &services.InvoiceService{DB: db}
ctx := authz.AsSystem(context.Background())

// nil precondition: no version check. A job is not editing a copy it read earlier.
if _, err := svc.Update(ctx, id, map[string]interface{}{"status": "overdue"}, nil); err != nil {
    return fmt.Errorf("marking invoice %s overdue: %w", id, err)
}`} />

              {/* ── Versions ─────────────────────────────── */}
              <h2 id="versions">Writes and Versions</h2>
              <p>
                <code>Update</code> and <code>Patch</code> take a <code>*concurrency.Precondition</code>.
                A handler builds it from the request&apos;s <code>If-Match</code> header with{' '}
                <code>concurrency.FromRequest(c)</code>; with one, the write lands only if the row is
                still at that version, and otherwise the method returns a{' '}
                <code>*concurrency.ErrConflict</code> naming the version it is at. <code>nil</code>{' '}
                means no check.
              </p>
              <p>
                <code>Update</code> writes the columns in the map it is given. The generated handler
                builds that map from its typed request; if you call it yourself, use column names.{' '}
                <code>Patch</code> takes a raw body and keeps only the writable columns.
              </p>

              {/* ── Your own methods ─────────────────────────────── */}
              <h2 id="business-logic">Adding Your Own Methods</h2>
              <p>
                Put them in a file of your own in the same package, such as{' '}
                <code>services/invoice_billing.go</code>. It is the same type, so your methods can use{' '}
                <code>s.db(ctx)</code> and the owner-checked <code>s.load</code>, and regenerating the
                resource never touches your file.
              </p>
              <CodeBlock language="go" filename="services/invoice_billing.go" code={`// MarkPaid records a payment. The rule lives here, so the route, the
// nightly reconciliation job and a test all get the same answer.
func (s *InvoiceService) MarkPaid(ctx context.Context, id string) (*models.Invoice, error) {
    item, err := s.load(ctx, id) // not found if the caller may not see it
    if err != nil {
        return nil, err
    }
    if item.PaidAt != nil {
        return nil, respond.Rule("invoice %s is already paid", item.Number) // 422
    }
    now := time.Now()
    if err := s.db(ctx).Model(item).Update("paid_at", now).Error; err != nil {
        return nil, fmt.Errorf("marking invoice paid: %w", err)
    }
    item.PaidAt = &now
    return item, nil
}`} />
              <p>The handler for it is four lines of HTTP:</p>
              <CodeBlock language="go" filename="handlers/invoice_billing.go" code={`func (h *InvoiceHandler) MarkPaid(c *gin.Context) {
    item, err := h.service().MarkPaid(h.ctx(c), c.Param("id"))
    if err != nil {
        h.fail(c, err, "Failed to mark invoice paid")
        return
    }
    c.JSON(http.StatusOK, gin.H{"data": item, "message": "Invoice marked paid"})
}`} />

              {/* ── Transaction Handling ─────────────────────────────── */}
              <h2 id="transactions">Transactions</h2>
              <p>
                When a method makes several writes that must succeed or fail together, run them in a
                GORM transaction. The generated service does this wherever a write has more than one
                statement: a row and its many-to-many links, a row and its line items, a bulk action.
              </p>
              <CodeBlock language="go" filename="services/order.go (CreateOrder)" code={`// CreateOrder creates an order and decrements product stock atomically.
func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error {
    return s.db(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(order).Error; err != nil {
            return fmt.Errorf("creating order: %w", err)
        }

        for i := range items {
            items[i].OrderID = order.ID
            if err := tx.Create(&items[i]).Error; err != nil {
                return fmt.Errorf("creating order item: %w", err)
            }

            result := tx.Model(&models.Product{}).
                Where("id = ? AND stock >= ?", items[i].ProductID, items[i].Quantity).
                Update("stock", gorm.Expr("stock - ?", items[i].Quantity))
            if result.Error != nil {
                return fmt.Errorf("updating stock: %w", result.Error)
            }
            if result.RowsAffected == 0 {
                return respond.Rule("not enough stock for product %s", items[i].ProductID)
            }
        }
        return nil // commit
    })
}`} />
              <ul>
                <li>Use <code>tx</code> for every query inside the callback, not <code>s.DB</code>.</li>
                <li>Returning <code>nil</code> commits; returning an error rolls everything back.</li>
                <li>A panic inside the callback is recovered and rolled back.</li>
                <li>Start from <code>s.db(ctx)</code>, so the transaction carries the request&apos;s context.</li>
              </ul>

              {/* ── Errors ─────────────────────────────── */}
              <h2 id="errors">Errors</h2>
              <p>
                Return Go errors, never HTTP codes; the handler&apos;s <code>fail</code> picks the
                status. Three have a meaning of their own:
              </p>
              <ul>
                <li><code>gorm.ErrRecordNotFound</code> becomes a 404.</li>
                <li><code>*concurrency.ErrConflict</code> becomes a 409 naming the current version.</li>
                <li>
                  <code>respond.Rule(&quot;...&quot;)</code> becomes a 422 carrying your message: use it
                  for a rule the caller broke. Anything else is logged and answered with a plain 500,
                  so wrap it with <code>fmt.Errorf(&quot;context: %w&quot;, err)</code> for the log.
                </li>
              </ul>

              {/* ── Best Practices ─────────────────────────────── */}
              <h2 id="best-practices">Best Practices</h2>
              <ul>
                <li>
                  <strong>One service per resource</strong>, each in its own file, with your additions
                  in files of your own beside it.
                </li>
                <li>
                  <strong>Context first, always.</strong> Query through <code>s.db(ctx)</code>; a query
                  on <code>s.DB</code> directly cannot see the caller, the organization or a
                  cancelled request.
                </li>
                <li>
                  <strong>Say who a job acts for.</strong> <code>authz.AsSystem</code> or{' '}
                  <code>authz.WithActor</code>; a bare <code>context.Background()</code> sees nothing
                  of an owned resource.
                </li>
                <li>
                  <strong>Keep services independent.</strong> When two need to collaborate, the caller
                  orchestrates them.
                </li>
              </ul>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 mt-10 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/handlers" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Handlers
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/backend/middleware" className="gap-1.5">
                  Middleware
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
