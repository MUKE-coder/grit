import Link from 'next/link'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/changelog')

export default function ChangelogPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Release History</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Changelog
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                All notable changes to Grit are documented here. Each release includes new features,
                bug fixes, and any breaking changes you need to be aware of.
              </p>
            </div>

            {/* v3.20.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.20.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Webhook receiver framework (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/57" target="_blank" rel="noopener noreferrer">#57</a>).
                  Wiring up Stripe / GitHub / WhatsApp / any HMAC-signed inbound webhook
                  is now {'<'}10 lines of app code. Signature verification, idempotency,
                  failed-handler replay — all framework concerns now.
                </p>

                <h3>The shape</h3>
                <pre><code>{`// In your app boot (e.g. internal/webhooks/handlers.go)
func init() {
    webhooks.Register("stripe", webhooks.Provider{
        SecretEnv: "STRIPE_WEBHOOK_SECRET",
        Verify:    webhooks.StripeVerifier,
        Extract:   webhooks.StripeExtractor,
    })

    webhooks.On("stripe", "invoice.paid", func(ctx context.Context, e *models.WebhookEvent) error {
        // … process the event
        return nil
    })
}`}</code></pre>
                <p>
                  The framework already mounted <code>POST /webhooks/:provider</code> in
                  routes — the path param picks the registered provider. No per-provider
                  routing code in your app.
                </p>

                <h3>Pipeline</h3>
                <ol>
                  <li>Route hits → look up provider (404 if unknown).</li>
                  <li>Read raw body + headers.</li>
                  <li>
                    <code>provider.Verify(secret, body, headers)</code> — 401 on signature
                    mismatch.
                  </li>
                  <li>
                    <code>provider.Extract(body, headers)</code> returns{' '}
                    <code>(eventType, externalID)</code>.
                  </li>
                  <li>
                    Insert into <code>webhook_events</code> — UNIQUE on{' '}
                    <code>(provider, external_id)</code> means duplicate deliveries
                    become <code>status=skipped</code> no-ops.
                  </li>
                  <li>
                    <code>webhooks.Dispatch(ctx, event)</code> runs the registered handler
                    for <code>(provider, eventType)</code>, falling back to a catch-all{' '}
                    <code>""</code> handler if no specific match.
                  </li>
                  <li>
                    Handler success → <code>status=processed</code>; handler error →{' '}
                    <code>status=failed</code> + <code>handler_error</code> recorded.
                    Provider always gets <code>200</code> once we persisted the event,
                    so retries don't hammer.
                  </li>
                </ol>

                <h3>Shipped verifiers</h3>
                <ul>
                  <li>
                    <code>HMACVerifier(header)</code> — generic hex HMAC-SHA256 in a
                    named header. Most simple partners use this.
                  </li>
                  <li>
                    <code>StripeVerifier</code> — Stripe's <code>t=...,v1=...</code>{' '}
                    scheme with 5-minute replay tolerance.
                  </li>
                  <li>
                    <code>GitHubVerifier</code> — GitHub's <code>X-Hub-Signature-256: sha256=...</code>{' '}
                    header.
                  </li>
                  <li>
                    Roll your own <code>VerifyFunc</code> for anything else — it's just{' '}
                    <code>func(secret string, body []byte, headers map[string]string) error</code>.
                  </li>
                </ul>

                <h3>Shipped extractors</h3>
                <ul>
                  <li>
                    <code>JSONFieldExtractor("type", "id")</code> — pulls top-level fields
                    from the JSON body. The most common shape (Stripe-style envelopes).
                  </li>
                  <li>
                    <code>GitHubExtractor</code> — reads <code>X-GitHub-Event</code> +{' '}
                    <code>X-GitHub-Delivery</code> headers.
                  </li>
                </ul>

                <h3>Admin endpoints</h3>
                <ul>
                  <li>
                    <code>GET /api/admin/webhooks?provider=stripe&status=failed</code> —
                    paginated list with the standard envelope. Filters: provider, status.
                  </li>
                  <li>
                    <code>POST /api/admin/webhooks/:id/replay</code> — re-runs the
                    handler for an existing event. Increments <code>retry_count</code>{' '}
                    and records the new outcome. Use this after a deploy fixes a
                    handler bug.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Pairs naturally with the v3.10 idempotency middleware — both are
                  &quot;safe replay&quot; primitives, just on different sides of the
                  network. Outbound retries reuse <code>Idempotency-Key</code>; inbound
                  duplicates dedupe on <code>(provider, external_id)</code>.
                </p>
              </div>
            </div>

            {/* v3.19.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.19.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Tamper-evident audit log via append-only hash chain (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/48" target="_blank" rel="noopener noreferrer">#48</a>).
                  Builds on the v3.16 <code>ActivityLog</code> — every row now carries{' '}
                  <code>PrevHash</code> + <code>Hash</code> columns where{' '}
                  <code>Hash = SHA-256(PrevHash || canonical(row))</code>. Mutating any
                  row breaks the chain on the next verification pass.
                </p>

                <h3>The chain</h3>
                <ul>
                  <li>
                    <strong>Genesis row</strong> has <code>PrevHash = ""</code>; every
                    subsequent row references the previous row's <code>Hash</code>.
                  </li>
                  <li>
                    Hash input is the <em>stable canonical form</em> of the
                    audit-relevant fields (user_id, method, path, status, payload digest,
                    IP, UA, duration, created_at unix-nano). ID + PrevHash + Hash
                    themselves are <em>not</em> in the canonical form — they're either
                    random (ID) or derived (Hash, PrevHash).
                  </li>
                  <li>
                    Insert uses <code>FOR UPDATE</code> lock on the latest row inside
                    the same transaction, so concurrent writes serialize cleanly without
                    forking the chain.
                  </li>
                </ul>

                <h3>The package</h3>
                <p>
                  New <code>internal/audit</code> ships these:
                </p>
                <ul>
                  <li>
                    <code>audit.Canonical(entry)</code> — stable JSON bytes for hashing.
                  </li>
                  <li>
                    <code>audit.ComputeHash(prevHash, canonical)</code> — runs SHA-256
                    over <code>prevHash || canonical</code>; returns hex.
                  </li>
                  <li>
                    <code>audit.AppendChained(db, entry)</code> — atomic insert with
                    chain lock. The middleware uses this; you can call it from anywhere.
                  </li>
                  <li>
                    <code>audit.VerifyChain(db)</code> — walks every row in{' '}
                    <code>(created_at, id)</code> order and recomputes hashes. Returns{' '}
                    <code>ChainStatus</code> with the first mismatch (broken_at_id +
                    expected vs got + message).
                  </li>
                </ul>

                <h3>The endpoint</h3>
                <pre><code>{`GET /api/admin/activity/integrity

→ { "valid": true, "total_entries": 12345 }

→ { "valid": false, "broken_at": 47, "broken_at_id": "uuid",
    "expected": "abc123...", "got": "def456...",
    "message": "hash mismatch — row was modified, deleted, or inserted out of order" }`}</code></pre>
                <p>
                  Wire this to a nightly cron + alerting webhook for free SOC2-ish audit
                  monitoring. Run it on-demand from a settings page when staff need the
                  current chain state.
                </p>

                <h3>What this defends against</h3>
                <ul>
                  <li>
                    Direct SQL <code>UPDATE</code> / <code>DELETE</code> on{' '}
                    <code>activity_logs</code> — the most common attack vector
                    (DBA covering tracks).
                  </li>
                  <li>
                    Out-of-band insertion of forged history.
                  </li>
                </ul>

                <h3>What it does NOT defend against</h3>
                <ul>
                  <li>
                    Compromise of the running server itself — an attacker with code
                    execution can rewrite the entire chain.
                  </li>
                  <li>
                    External anchoring (publishing the daily root hash to a public
                    ledger like a tweet, a transaction, or a Sigstore log) is the
                    follow-up — flagged in #48 as bonus material, not shipped here.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Verification cost: O(n) — about 2–3 seconds per million rows on a warm
                  cache. The middleware insert is still fire-and-forget so audit DB
                  latency never blocks the response path; chain failures log instead of
                  cascading.
                </p>
              </div>
            </div>

            {/* v3.18.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.18.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  PDF generation module (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/13" target="_blank" rel="noopener noreferrer">#13</a>).
                  Every scaffolded API ships <code>internal/pdf/</code> with Grit-styled
                  section helpers + a worked <code>RenderInvoice</code> template. Pure Go,
                  no Chromium / wkhtmltopdf native dependencies.
                </p>

                <h3>The Doc primitives</h3>
                <p>
                  <code>pdf.New()</code> returns a <code>*Doc</code> preconfigured with
                  Helvetica + 20mm margins + A4 portrait + Grit blue accent. Embeds the
                  underlying <code>*fpdf.Fpdf</code> so the full library is available when
                  helpers don't fit.
                </p>
                <ul>
                  <li>
                    <code>Header(title, subtitle)</code> — accent-colored 22pt title +
                    muted-gray subtitle line.
                  </li>
                  <li>
                    <code>KV(label, value)</code> + <code>TwoColumnKV(...)</code> — small-caps
                    label + body value pairs.
                  </li>
                  <li>
                    <code>Table(headers, rows, widths, aligns)</code> — light gray header
                    row, plain data rows, configurable widths + alignment per column.
                  </li>
                  <li>
                    <code>Totals([]TotalLine)</code> — right-aligned totals stack; the bold
                    line gets accent coloring + slightly larger size for the grand total.
                  </li>
                  <li>
                    <code>Notes(text)</code> — labeled multiline section, skipped when empty.
                  </li>
                  <li>
                    <code>Footer(text)</code> — centered italic 25mm above the page bottom.
                  </li>
                  <li>
                    <code>d.Bytes()</code> finalizes and returns the PDF byte slice ready
                    to stream to <code>c.Data(200, "application/pdf", b)</code>.
                  </li>
                </ul>

                <h3>RenderInvoice — worked example</h3>
                <pre><code>{`pdf.RenderInvoice(pdf.Invoice{
    Number:    "INV-202605-0001",
    IssueDate: time.Now(),
    DueDate:   time.Now().Add(14 * 24 * time.Hour),
    BillTo:    pdf.Party{Name: "Abu Seal", Contact: "abu@example.com"},
    Items: []pdf.LineItem{
        {Description: "Office rent — June", Quantity: 1, UnitPrice: 1500000, Total: 1500000},
        {Description: "Service charge",      Quantity: 1, UnitPrice:  120000, Total:  120000},
    },
    Subtotal: 1620000, Total: 1620000,
    Currency: "UGX",
    Notes:    "Pay by mobile money: +256...",
})`}</code></pre>
                <p>
                  Returns <code>([]byte, error)</code> — wire it to a handler:
                </p>
                <pre><code>{`func (h *InvoiceHandler) PDF(c *gin.Context) {
    inv, _ := h.Service.GetByID(c.Param("id"))
    bytes, err := pdf.RenderInvoice(toInvoice(inv))
    if err != nil { respond.Internal(c, err); return }
    c.Header("Content-Disposition", \`attachment; filename="\` + inv.Number + \`.pdf"\`)
    c.Data(200, "application/pdf", bytes)
}`}</code></pre>
                <p className="text-sm text-muted-foreground">
                  Copy <code>invoice.go</code> as a starting point for receipts, leases,
                  statements, quotes — the same primitives compose all of them. Add
                  <code> github.com/go-pdf/fpdf v0.9.0</code> dependency lands automatically
                  in scaffolded <code>go.mod</code>.
                </p>
              </div>
            </div>

            {/* v3.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.17.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Quality-of-life bundle. Four GitHub issues closed:
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer"> #12</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/31" target="_blank" rel="noopener noreferrer">#31</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/35" target="_blank" rel="noopener noreferrer">#35</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/43" target="_blank" rel="noopener noreferrer">#43</a>.
                </p>

                <h3><code>grit init</code> — #35</h3>
                <p>
                  New CLI command writes <code>CLAUDE.md</code> + <code>AGENTS.md</code>{' '}
                  to the current directory. Both files carry the framework's hard rules
                  (Forms / Frontend stdlib / Data / Backend / Resources / Sync / Auth)
                  so contributors and AI assistants get the conventions right on first
                  PR. Skips existing files unless <code>--force</code> is passed; re-run
                  with <code>--force</code> after a major framework upgrade to refresh.
                </p>

                <h3>Verbose AutoMigrate — #31</h3>
                <p>
                  <code>Migrate()</code> now snapshots <code>ColumnTypes</code> before
                  and after each <code>AutoMigrate</code> call and logs a diff:
                </p>
                <pre><code>{`================================================================
DATABASE MIGRATION — 8 model(s) registered
================================================================
  + created models.Building
  ~ models.User — added 2 column(s): is_vip, vip_notes
----------------------------------------------------------------
Migration done — 1 created, 1 altered (+2 column), 6 unchanged.
================================================================`}</code></pre>
                <p>
                  Silent migrations are gone. Also fixes a pre-existing bug where{' '}
                  <code>Migrate</code> <em>skipped</em> already-existing tables — so
                  columns added to a model never actually landed in the DB. Now they do.
                </p>

                <h3>Cursor-based pagination — #43</h3>
                <ul>
                  <li>
                    <code>paginate.List</code> gains opt-in cursor mode via{' '}
                    <code>Config.CursorMode: true</code>. Response carries{' '}
                    <code>Meta.NextCursor</code> + <code>Meta.HasMore</code>{' '}
                    instead of <code>Page</code>/<code>Pages</code>.
                  </li>
                  <li>
                    Detects <code>HasMore</code> by fetching <code>PageSize + 1</code>{' '}
                    rows — no separate count query needed.
                  </li>
                  <li>
                    Cursor is opaque base64 of <code>(sort_value, id)</code> so pages
                    stay stable when rows insert mid-pagination. Works with any sort
                    field; extracts the value via reflection on the last row.
                  </li>
                  <li>
                    Total count opt-in via <code>Config.IncludeTotal</code> — costs an
                    extra <code>COUNT(*)</code>, leave off unless your UI shows a
                    &quot;X of Y&quot; indicator.
                  </li>
                  <li>
                    Offset mode stays the default for back-compat; new resources can
                    flip the flag.
                  </li>
                </ul>

                <h3>Generator quality — #12</h3>
                <p>The remaining tag-default heuristics from issue #12:</p>
                <ul>
                  <li>
                    <strong>URL fields</strong> (suffix <code>_url</code> + named{' '}
                    <code>url</code> / <code>image</code> / <code>avatar</code> /{' '}
                    <code>thumbnail</code> / <code>logo</code> / <code>cover</code> /{' '}
                    <code>icon</code> / <code>banner</code> / <code>photo</code>) get{' '}
                    <code>size:500</code> instead of <code>size:255</code>. UTM-tagged
                    links and signed S3 URLs blow past 255 in the wild.
                  </li>
                  <li>
                    <strong>Long-text fields</strong> named <code>description</code> /{' '}
                    <code>notes</code> / <code>content</code> / <code>body</code> /{' '}
                    <code>summary</code> / <code>bio</code> / <code>details</code> /{' '}
                    <code>comment</code> / <code>comments</code> / <code>message</code>{' '}
                    get <code>type:text</code>.
                  </li>
                  <li>
                    <strong>Money fields</strong> on <code>float</code> type (suffix{' '}
                    <code>_amount</code> / <code>_price</code> / <code>_total</code> /{' '}
                    <code>_cost</code> / <code>_fee</code> / <code>_balance</code> /{' '}
                    <code>_rent</code> / <code>_salary</code> / <code>_wage</code> /{' '}
                    <code>_value</code> / <code>_revenue</code> / <code>_deposit</code>{' '}
                    + named <code>amount</code> / <code>price</code> / <code>total</code> /{' '}
                    <code>cost</code> / <code>fee</code> / <code>balance</code> /{' '}
                    <code>subtotal</code>) get <code>type:decimal(12,2)</code> for
                    fixed-precision storage. No more <code>1.99 + 0.01 = 1.9999999</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.16.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Three coherent admin-operations features at once:
                  CSV/Excel export per resource (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/15" target="_blank" rel="noopener noreferrer">#15</a>),
                  activity audit log middleware (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/32" target="_blank" rel="noopener noreferrer">#32</a>),
                  and the <code>apiErrorMessage</code> frontend helper (
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/27" target="_blank" rel="noopener noreferrer">#27</a>).
                </p>

                <h3>CSV / Excel export per resource — #15</h3>
                <ul>
                  <li>
                    New <code>internal/export</code> package: <code>CSV(w, items, opts)</code>{' '}
                    and <code>XLSX(w, items, opts)</code> with a typed{' '}
                    <code>Column{'{Header, Field, Format}'}</code> config. Field uses
                    dot-notation for associations (<code>"Tenant.Name"</code>).
                  </li>
                  <li>
                    Format strings: <code>"currency:UGX"</code>, <code>"date:2006-01-02"</code>,{' '}
                    <code>"datetime"</code>, <code>"bool"</code>. Empty string falls back to{' '}
                    <code>fmt.Sprintf("%v")</code>.
                  </li>
                  <li>
                    Resource generator now emits an <code>Export(c *gin.Context)</code>{' '}
                    handler method on every new resource, with columns derived from the
                    field list. Routes inject <code>GET /api/{'<plural>'}/export</code>{' '}
                    automatically.
                  </li>
                  <li>
                    Honours the same <code>search</code> param as List, so users can
                    export a filtered subset.
                  </li>
                  <li>
                    Adds <code>github.com/xuri/excelize/v2 v2.8.1</code> to scaffolded{' '}
                    <code>go.mod</code>.
                  </li>
                </ul>

                <h3>Activity audit log — #32</h3>
                <ul>
                  <li>
                    New <code>models.ActivityLog</code> with user_id + method + path +
                    status + payload digest (sha256, not raw body) + IP + user-agent +
                    duration. UUID PK; <code>created_at</code> indexed for time-range queries.
                  </li>
                  <li>
                    New <code>middleware.ActivityLogger(db)</code> mounted on every
                    protected mutation route. Skips safe methods + non-2xx responses +
                    unauthenticated requests.
                  </li>
                  <li>
                    Insert is fire-and-forget (goroutine). Audit DB latency never
                    blocks the response path; if the DB is down the entry drops
                    rather than failing the request.
                  </li>
                  <li>
                    New endpoint <code>GET /api/admin/activity</code> (admin-only) with{' '}
                    <code>paginate.List</code> filtering by <code>user_id</code>,{' '}
                    <code>method</code>, and <code>path</code> prefix. Drop in any audit-log UI.
                  </li>
                </ul>

                <h3>apiErrorMessage helper — #27</h3>
                <ul>
                  <li>
                    Three helpers in <code>packages/shared/types/api.ts</code>:{' '}
                    <code>apiErrorMessage(err, fallback?)</code>,{' '}
                    <code>apiErrorCode(err)</code>,{' '}
                    <code>apiErrorFields(err)</code>.
                  </li>
                  <li>
                    Walks the standard envelope chain
                    (<code>response.data.error.message</code>) plus axios{' '}
                    <code>err.message</code> plus a fallback so{' '}
                    <code>toast.error(apiErrorMessage(err))</code> is always meaningful.
                  </li>
                  <li>
                    <code>apiErrorCode</code> returns the envelope's <code>code</code>{' '}
                    string (<code>"VALIDATION_ERROR"</code>,{' '}
                    <code>"VERSION_CONFLICT"</code>, etc.) for branching logic.
                  </li>
                  <li>
                    <code>apiErrorFields</code> surfaces per-field validation details so
                    forms can highlight specific inputs.
                  </li>
                  <li>
                    New <code>internal/respond</code> package on the server side too:{' '}
                    <code>respond.NotFound / Validation / Forbidden / Conflict /
                    Internal</code> for handlers, replacing ad-hoc inline{' '}
                    <code>c.JSON(500, gin.H{'{...}'})</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.15.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Frontend stdlib + form primitives. Closes seven GitHub issues at once
                  (<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/19" target="_blank" rel="noopener noreferrer">#19</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/20" target="_blank" rel="noopener noreferrer">#20</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/21" target="_blank" rel="noopener noreferrer">#21</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/22" target="_blank" rel="noopener noreferrer">#22</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/23" target="_blank" rel="noopener noreferrer">#23</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/33" target="_blank" rel="noopener noreferrer">#33</a>,{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/34" target="_blank" rel="noopener noreferrer">#34</a>).
                  Every primitive lifted from real Grit-built business apps.
                </p>

                <h3>Format helpers (<code>lib/format.ts</code>) — #33</h3>
                <ul>
                  <li>
                    <code>formatCurrency(amount, currency?)</code> — locale-aware,
                    no-decimal mode for UGX / JPY / KRW / RWF / TZS / VND.
                  </li>
                  <li>
                    <code>formatDate(value, fmt?)</code> — token formatter (yyyy / MMMM /
                    MMM / MM / dd / HH / mm / ss). Default <code>"MMM d, yyyy"</code>.
                  </li>
                  <li>
                    <code>formatDateTime(value)</code> — &quot;May 2, 2026 · 2:30 PM&quot;.
                  </li>
                  <li>
                    <code>humanize("checked_in")</code> → &quot;Checked in&quot;.
                  </li>
                  <li>
                    <code>initials("Abu Seal")</code> → &quot;AS&quot;.
                  </li>
                  <li>
                    <code>setFormatConfig({"{ locale, currency }"})</code> at boot to
                    override.
                  </li>
                </ul>

                <h3><code>{'<CurrencyField>'}</code> — #19</h3>
                <p>
                  Live comma formatting as the user types
                  (&quot;3000&quot; → &quot;3,000&quot;), paste-friendly
                  (&quot;$1,234.56&quot; works), emits raw <code>number</code> to{' '}
                  <code>onChange</code>. Optional <code>prefix</code> slot for currency
                  code. Auto-toggles between formatted display (blur) and raw digits
                  (focus) so editing isn't a fight.
                </p>

                <h3><code>{'<SearchableSelect>'}</code> — #20</h3>
                <p>
                  Combobox with typeahead, ↑/↓/Enter/Esc keyboard nav, portaled
                  dropdown (escapes overflow:hidden ancestors), optional clear button.
                  Replaces native <code>{'<select>'}</code> for FK fields and any enum
                  with &gt; 5 values.
                </p>

                <h3><code>{'<DateField>'}</code> + <code>{'<DateRangeFilter>'}</code> — #21</h3>
                <ul>
                  <li>
                    <code>{'<DateField>'}</code> wraps native{' '}
                    <code>{'<input type="date">'}</code> with the standard label/hint/
                    error chrome — picked the native one for a11y, RTL, and i18n.
                  </li>
                  <li>
                    <code>{'<DateRangeFilter>'}</code> — preset chip bar (Today / Last 7 /
                    Last 30 / This month / Last month / Last 90 / This year / All) +
                    custom-range fallback.
                  </li>
                  <li>
                    <code>presetRange("last90")</code> helper exposed for non-UI uses.
                  </li>
                </ul>

                <h3><code>{'<Drawer>'}</code> — #22</h3>
                <p>
                  Right-edge slide-in panel. Closes on Esc + backdrop + X button.
                  Configurable widths (<code>sm</code>/<code>md</code>/<code>lg</code>/<code>xl</code>).
                  Optional sticky footer slot for the typical Cancel/Save row. Pair
                  with <code>{'<FormGrid>'}</code> + <code>{'<FormActions>'}</code>{' '}
                  from v3.11 for the standard create/edit experience.
                </p>

                <h3><code>{'<StatusBadge>'}</code> — #34</h3>
                <p>
                  Status string → coloured pill. Default map covers
                  paid / active / completed / pending / overdue / cancelled / draft / archived /
                  checked_in / in_progress and friends. Override or extend per app:
                </p>
                <pre><code>{`setStatusVariants({
  shipped: "info",
  on_hold: "warning",
});`}</code></pre>

                <h3><code>{'<AppShell>'}</code> + grouped sidebar — #23</h3>
                <ul>
                  <li>
                    New <code>lib/nav-config.ts</code> — single source of truth for
                    sidebar sections. Adding a new section is a one-line config edit.
                  </li>
                  <li>
                    <code>components/layout/sidebar.tsx</code> rewritten as a
                    config-driven grouped sidebar (section title + items with icons +
                    optional badges).
                  </li>
                  <li>
                    New <code>components/layout/app-shell.tsx</code> — bundles
                    TitleBar + Sidebar + Topbar + scrollable content + Cmd/Ctrl-K
                    command palette in one component. Wrap your dashboard{' '}
                    <code>Outlet</code> with this.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.14.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  <strong>Offline-first foundation.</strong> Git-style sync model — work
                  locally, click Sync explicitly, resolve conflicts per-field, push
                  one-by-one. Every scaffolded API now has Version-tracked rows + the
                  POST /api/sync/push and GET /api/sync/pull endpoints; every desktop
                  scaffold ships a local SQLite mirror, an outbox with squash semantics,
                  and a title-bar Sync button + conflict-resolution dialog.
                </p>

                <h3>Server: versioning + sync endpoints</h3>
                <ul>
                  <li>
                    <code>Version int</code> column added to User, Upload, Blog. A{' '}
                    <code>BeforeUpdate</code> GORM hook auto-increments on every
                    server-side write. The resource generator emits both on every new
                    model.
                  </li>
                  <li>
                    <code>POST /api/sync/push</code> accepts a batch of changes; each
                    entry includes the version the client believes the server has. On
                    mismatch the response contains <code>VERSION_CONFLICT</code> + the
                    current server state, so the client can drive a merge UI.
                  </li>
                  <li>
                    <code>GET /api/sync/pull?model=X&since=cursor</code> returns every
                    row in the table updated after the cursor, paginated, with a new
                    cursor in the response.
                  </li>
                  <li>
                    New <code>internal/sync/registry.go</code> maps logical table names
                    (e.g. <code>"buildings"</code>) to <code>reflect.Type</code> so the
                    handler decodes dynamic payloads. New resources auto-register via{' '}
                    <code>// grit:sync</code> marker.
                  </li>
                </ul>

                <h3>Desktop: sync engine</h3>
                <ul>
                  <li>
                    New <code>apps/desktop/sync/</code> Go package. Opens a local SQLite
                    file under the OS user-config dir on app boot.
                  </li>
                  <li>
                    Three tables: <code>sync_records</code> (local mirror — reads come
                    from here), <code>sync_outbox</code> (pending changes; UNIQUE on
                    (model, entity_id) for squash), <code>sync_cursors</code>{' '}
                    (incremental pull positions).
                  </li>
                  <li>
                    Squash semantics: edit a record three times offline → one outbox
                    entry with the final state. delete-after-create cancels both
                    locally without ever hitting the network.
                  </li>
                  <li>
                    <code>Engine.Sync()</code> runs Pull then Push. Push posts the whole
                    outbox in one HTTP call; the response drives per-entry result
                    handling — successes clear from the outbox, conflicts get the
                    server state stashed for the merge UI.
                  </li>
                </ul>

                <h3>Wails bindings</h3>
                <p>The frontend talks to the engine through these Wails-bound methods on App:</p>
                <ul>
                  <li>
                    <code>LocalCreate</code> / <code>LocalUpdate</code> / <code>LocalDelete</code> —
                    write-through to local SQLite + outbox.
                  </li>
                  <li>
                    <code>LocalGet</code> / <code>LocalList</code> — read from the local mirror.
                  </li>
                  <li>
                    <code>Sync(tables)</code> — pull listed tables then push the outbox.
                    Returns counts.
                  </li>
                  <li>
                    <code>PendingCount</code>, <code>GetPendingChanges</code> — drive the
                    title-bar badge and the review panel.
                  </li>
                  <li>
                    <code>ResolveConflict(table, entityID, mergedData, serverVersion)</code> —
                    accepts the user's merge for a conflicted entry.
                  </li>
                </ul>

                <h3>UI</h3>
                <ul>
                  <li>
                    <strong>Title-bar Sync button</strong> with a pending-count badge.
                    Green refresh icon when clean; amber alert + count when there are
                    pending changes.
                  </li>
                  <li>
                    <strong>PendingChangesPanel</strong> — right-edge drawer listing
                    every outbox entry, split into &quot;Needs review&quot; (conflicts)
                    and &quot;Ready to push&quot;. <code>Sync now</code> button at the bottom.
                  </li>
                  <li>
                    <strong>ConflictDialog</strong> — field-level merge UI. Three columns
                    (Field / Local / Server v_N), per-field click to choose. Apply
                    builds the merged record and calls <code>ResolveConflict</code>.
                  </li>
                </ul>

                <h3>React hooks</h3>
                <ul>
                  <li>
                    <code>usePendingCount()</code> — polls every 2s for the badge.
                  </li>
                  <li>
                    <code>usePendingChanges()</code> — full outbox + refresh function.
                  </li>
                  <li>
                    <code>useSyncMutation(tables)</code> — kicks off a Sync, exposes
                    running/result/error state.
                  </li>
                  <li>
                    <code>useResolveConflict()</code> — applies one merge and refreshes.
                  </li>
                </ul>

                <h3>Wire format</h3>
                <pre><code>{`POST /api/sync/push
{ "changes": [
    { "op": "create", "model": "buildings", "id": "uuid", "version": 0, "data": {...} },
    { "op": "update", "model": "tenants",   "id": "uuid", "version": 5, "data": {...} },
    { "op": "delete", "model": "leases",    "id": "uuid", "version": 3 }
] }

→ { "results": [
    { "ok": true, "new_version": 1 },
    { "ok": false, "code": "VERSION_CONFLICT", "server_version": 7, "server_data": {...} },
    { "ok": true }
] }`}</code></pre>

                <p className="text-sm text-muted-foreground">
                  Deferred to v3.14.1: React Query offline-aware data hooks
                  (<code>useOfflineList</code>, <code>useOfflineGet</code>,{' '}
                  <code>useOfflineMutation</code>) and the resource generator emitting
                  offline-aware frontend hooks. The engine and primitives ship now;
                  ergonomics layer next.
                </p>
              </div>
            </div>

            {/* v3.13.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.13.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  New <code>grit generate sequence</code> command produces atomic, gap-free
                  sequential numbers like <code>INV-202605-0001</code>. Pattern lifted from a
                  real Grit-built rental management app — invoice / receipt / order numbering
                  is now a one-liner.
                </p>

                <h3>What it generates</h3>
                <ul>
                  <li>
                    <strong>First invocation only:</strong>{' '}
                    <code>internal/sequence/sequence.go</code> — a generic counter package
                    with <code>Counter</code> (the GORM-backed row),{' '}
                    <code>Config</code> (name + prefix + reset + width), and an atomic{' '}
                    <code>Next(db, cfg, t)</code> helper.
                  </li>
                  <li>
                    <strong>Every invocation:</strong>{' '}
                    <code>internal/services/&lt;name&gt;_sequence.go</code> — a typed
                    convenience wrapper. Handlers call e.g.{' '}
                    <code>services.NextInvoiceNumber(h.DB, time.Now())</code> without
                    knowing the prefix or reset cadence.
                  </li>
                  <li>
                    Auto-injects <code>&sequence.Counter{'{}'}</code> into the{' '}
                    <code>Models()</code> migration slice (idempotent).
                  </li>
                </ul>

                <h3>Mechanics</h3>
                <ul>
                  <li>
                    Counter rows keyed by <code>(name, bucket)</code> where bucket is{' '}
                    <code>"YYYYMM"</code> for monthly resets, <code>"YYYY"</code> for yearly,
                    or empty for never. So a monthly counter automatically restarts at 1
                    on the first call of each new month.
                  </li>
                  <li>
                    Atomic via row-level <code>SELECT FOR UPDATE</code> on Postgres
                    (concurrent callers serialize on the counter row). SQLite serializes
                    writes globally so it's also safe.
                  </li>
                </ul>

                <h3>Usage</h3>
                <pre><code>{`grit generate sequence Invoice
grit generate sequence Order --prefix ORD --reset yearly --width 6
grit generate sequence Receipt --reset never`}</code></pre>

                <p>Flags:</p>
                <ul>
                  <li>
                    <code>--prefix</code> — alphabetic prefix (default: first 3 chars of
                    the name, uppercased)
                  </li>
                  <li>
                    <code>--reset</code> — when the counter resets:{' '}
                    <code>monthly</code> (default), <code>yearly</code>, or <code>never</code>
                  </li>
                  <li>
                    <code>--width</code> — zero-padded width of the numeric portion
                    (default 4)
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  The <code>grit generate report</code> generator (Recharts tabs page +
                  Go ReportService) is deferred to a future release — it needs more design
                  work for the React chart layer than fits a same-day release.
                </p>
              </div>
            </div>

            {/* v3.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.12.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Realtime WebSocket hub baked into every API + a desktop client + hooks
                  for subscribing. And a sweep of every remaining numeric ID — UUIDs are
                  now the canonical ID type everywhere in the framework.
                </p>

                <h3>Realtime hub (API)</h3>
                <ul>
                  <li>
                    New package: <code>internal/realtime/hub.go</code>. One <code>Hub</code> per
                    process; each user can have multiple connections (desktop + mobile + web).
                  </li>
                  <li>
                    <code>Hub.SendToUser(userID, evt)</code>, <code>SendToUsers(ids, evt)</code>,
                    and <code>Broadcast(evt)</code> let any handler or service push events.
                  </li>
                  <li>
                    Slow-client safe: per-connection 32-message send buffer; when full, that
                    one client's message is dropped — never blocks the entire hub. Slow
                    clients resync on their next REST refetch.
                  </li>
                  <li>
                    New handler: <code>internal/handlers/realtime.go</code> upgrades the
                    request to a WebSocket and registers the client with the hub.
                  </li>
                  <li>
                    Mounted at <code>GET /api/ws?token=&lt;jwt&gt;</code> — query-string auth
                    because browsers can't set custom headers on the WS handshake.
                  </li>
                  <li>
                    Wire format: <code>{'{ type: "<topic>", payload: {...} }'}</code>. Suggested
                    topics: <code>chat.message.new</code>, <code>notification.new</code>,{' '}
                    <code>system.connected</code>, or your own{' '}
                    <code>{'resource.<name>.<verb>'}</code> namespace.
                  </li>
                  <li>
                    Dependency added: <code>github.com/gorilla/websocket v1.5.3</code>.
                  </li>
                </ul>

                <h3>Realtime client (desktop)</h3>
                <ul>
                  <li>
                    New file: <code>frontend/src/lib/realtime.ts</code>. Singleton client with
                    auto-reconnect via exponential backoff (1s, 2s, 4s, 8s, capped at 15s).
                  </li>
                  <li>
                    Global <code>realtimeBus</code> EventTarget — any component can subscribe.
                  </li>
                  <li>
                    Start from <code>AuthProvider</code> after tokens land, stop on logout.
                  </li>
                  <li>
                    New hook: <code>useRealtimeEvent&lt;T&gt;(type, callback)</code> subscribes to a
                    typed topic and unsubscribes on unmount. Plus{' '}
                    <code>useRealtimeAny()</code> for catch-all handlers (debug, toast bar).
                  </li>
                </ul>

                <h3>ID consistency sweep — UUIDs everywhere</h3>
                <p>
                  v3.9.1 standardized the User model on string UUID PKs but a long tail of
                  numeric IDs remained in the framework. v3.12.0 cleans them all up.
                </p>
                <ul>
                  <li>
                    <strong>Go scaffold:</strong> the prebuilt <code>Blog</code> model in{' '}
                    <code>api_blog_files.go</code> swaps from <code>gorm.Model</code> (auto-incr
                    uint) to a string UUID PK with a <code>BeforeCreate</code> hook. Service
                    signatures (<code>GetByID</code>, <code>Update</code>, <code>Delete</code>) and
                    handler param parsing all switch from <code>uint</code> to <code>string</code>.
                  </li>
                  <li>
                    <strong>Standalone desktop scaffold</strong> (<code>grit new-desktop</code>):
                    User, Blog, and Contact models all switch to string UUID PKs with{' '}
                    <code>BeforeCreate</code> hooks. All Wails-bound App methods and underlying
                    service signatures use <code>id string</code>. Frontend mutation typings
                    follow.
                  </li>
                  <li>
                    <strong>Shared TS types:</strong> <code>User</code>, <code>Upload</code>, and{' '}
                    <code>Blog</code> interfaces all use <code>id: string</code>. URL builders in{' '}
                    <code>API_ROUTES</code> take <code>id: string</code>. The <code>BlogSchema</code>{' '}
                    Zod schema uses <code>z.string()</code>.
                  </li>
                  <li>
                    <strong>Admin TS:</strong> <code>DataTable</code> selection state, generic{' '}
                    <code>useResourceItem</code> / <code>useUpdateResource</code> /{' '}
                    <code>useDeleteResource</code> mutation typings,{' '}
                    <code>RelationshipSelectField</code> single value,{' '}
                    <code>MultiRelationshipSelectField</code> array values, and <code>handleDelete</code>{' '}
                    callbacks all switch from <code>number</code> to <code>string</code>.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  Net effect: UUID is the canonical ID type across the entire framework. Any
                  resource generator output, any scaffolded type, any Go signature — all
                  string UUIDs. No more <code>id: number</code> hiding in some corner.
                </p>
              </div>
            </div>

            {/* v3.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.11.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Three new desktop primitive files ship with every <code>--desktop</code>{' '}
                  scaffold, lifted from a real Grit-built rental management app. They cover
                  the master-detail layout, form chrome, and filter chips that every CRUD
                  page reinvents — saving ~200 LOC per resource.
                </p>

                <h3><code>components/two-pane.tsx</code> — master-detail layout</h3>
                <ul>
                  <li>
                    <code>TwoPane</code> — outer flex container with overflow handling.
                  </li>
                  <li>
                    <code>ListPane</code> — fixed-width (352px) left pane with title +
                    count + new button + searchbar + optional filters slot + scrollable body
                    + optional footer. Toolbar slot for refresh buttons or other actions.
                  </li>
                  <li>
                    <code>ListRow</code> — icon/avatar + title + subtitle + right-side
                    meta. Selected state shows a 2px accent bar on the left edge.
                  </li>
                  <li>
                    <code>DetailPane</code> — right pane with optional header + scrollable
                    content. <code>empty=true</code> renders an <code>EmptyState</code>{' '}
                    with the configured title/hint instead.
                  </li>
                  <li>
                    <code>EmptyState</code>, <code>DetailSection</code> (small caps section
                    header), and <code>DetailField</code> (labelled value rows for read views).
                  </li>
                </ul>

                <h3><code>components/form.tsx</code> — form chrome</h3>
                <ul>
                  <li>
                    <code>TextField</code>, <code>TextAreaField</code>, <code>SelectField</code> —
                    forwarded refs, consistent label/hint/error layout, focus ring,
                    disabled styling. Plug straight into <code>react-hook-form</code>.
                  </li>
                  <li>
                    <code>FormGrid</code> — 1, 2, or 3 columns on <code>{'>'}=sm</code>; stacks
                    on small screens.
                  </li>
                  <li>
                    <code>FormSection</code> — small caps title + optional description
                    over a stack of fields.
                  </li>
                  <li>
                    <code>FormActions</code> — Cancel + Submit pair with{' '}
                    <code>isPending</code> support (button disables, label flips to
                    &quot;Saving...&quot;).
                  </li>
                </ul>

                <h3><code>components/filter-chip.tsx</code> — filter chips</h3>
                <ul>
                  <li>
                    <code>FilterChip</code> — toggleable pill, active state shows accent
                    background; optional <code>onClear</code> renders an X to clear a single
                    filter; optional <code>count</code> renders a small count badge.
                  </li>
                  <li>
                    <code>FilterBar</code> — horizontal scrollable wrapper. Drop into
                    <code>ListPane</code>&apos;s <code>filters</code> slot.
                  </li>
                </ul>

                <h3>Tailwind tokens</h3>
                <ul>
                  <li>
                    Added <code>listpane</code> spacing token (<code>22rem</code> / 352px)
                    to the desktop Tailwind config so <code>w-listpane</code> works.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.10.0
                </span>
                <span className="text-sm text-muted-foreground">May 2, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Foundation release for upcoming offline-first work. Every scaffolded API
                  now ships with idempotent-retry semantics; every scaffolded client
                  auto-attaches an <code>Idempotency-Key</code> on mutations; and the
                  desktop scaffold gains a connection-status indicator backed by an
                  API heartbeat.
                </p>

                <h3>Idempotency middleware (API)</h3>
                <ul>
                  <li>
                    New file: <code>internal/middleware/idempotency.go</code> — wired into{' '}
                    <code>routes.Setup</code> as a global middleware.
                  </li>
                  <li>
                    Activates only when the request carries an <code>Idempotency-Key</code> header
                    and the method is <code>POST</code>/<code>PUT</code>/<code>PATCH</code>/<code>DELETE</code>.
                  </li>
                  <li>
                    First 2xx response is cached in Redis for <strong>24 hours</strong>, keyed by{' '}
                    <code>(method, path, key)</code>. Subsequent requests with the same key replay
                    the cached response instead of re-executing the handler.
                  </li>
                  <li>
                    Errors (4xx/5xx) are intentionally <strong>not</strong> cached — clients can
                    retry transient failures with the same key.
                  </li>
                  <li>
                    Sets <code>Idempotent-Replayed: true</code> response header on cache hits so
                    clients can distinguish replays from fresh executions.
                  </li>
                </ul>

                <h3>Client-side header injection</h3>
                <ul>
                  <li>
                    Desktop, Expo, web, and admin clients all auto-attach a UUIDv4{' '}
                    <code>Idempotency-Key</code> on unsafe methods via the request interceptor.
                  </li>
                  <li>
                    The 401-refresh path now reuses the same key when re-issuing a request after a
                    token refresh — so a token expiring mid-write can never double-create.
                  </li>
                </ul>

                <h3>Online-status hook (desktop)</h3>
                <ul>
                  <li>
                    New hook: <code>useOnlineStatus()</code> at{' '}
                    <code>frontend/src/hooks/use-online-status.ts</code>.
                  </li>
                  <li>
                    Combines <code>navigator.onLine</code> (cheap pre-check) with a 15-second
                    heartbeat to <code>/api/health</code> (the truth signal). Returns{' '}
                    <code>{'{ isOnline, lastCheckedAt }'}</code>.
                  </li>
                  <li>
                    Heartbeat times out after 5s so a sleeping laptop surfaces as offline
                    instantly on wake.
                  </li>
                  <li>
                    The title-bar gains a <code>ConnectionIndicator</code> — small green/amber
                    dot reflecting API reachability. Hover for last-checked timestamp.
                  </li>
                </ul>

                <p className="text-sm text-muted-foreground">
                  This is the foundation for the offline-first scaffold landing in a later
                  v3.x release — write-queues, optimistic updates, and last-write-wins
                  conflict resolution all need stable idempotency keys to be safe.
                </p>
              </div>
            </div>

            {/* v3.9.2 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.2
                </span>
                <span className="text-sm text-muted-foreground">April 25, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Every <code>grit generate resource</code> run now emits a List handler that is
                  ~15 lines instead of ~55. The page / sort / search boilerplate moved into a shared{' '}
                  <code>internal/paginate</code> package that ships with every scaffolded API —
                  one source of truth for clamping, whitelisting, and search. Addresses{' '}
                  <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/14" target="_blank" rel="noopener noreferrer">issue #14</a>.
                </p>

                <h3>New <code>paginate</code> package</h3>
                <ul>
                  <li>
                    <code>paginate.List[T](query, paginate.Bind(c), paginate.Config{'{'}...{'}'})</code> —
                    typed, generic helper that runs search, sort, filter, and pagination against
                    any <code>*gorm.DB</code> query.
                  </li>
                  <li>
                    <code>paginate.Bind(c)</code> reads <code>page</code>, <code>page_size</code>,
                    {' '}<code>search</code>, <code>sort_by</code>, <code>sort_order</code> from the
                    Gin query, clamps <code>page</code> to ≥ 1 and <code>page_size</code> to
                    {' '}<code>[1, 100]</code>.
                  </li>
                  <li>
                    <code>paginate.Config</code> whitelists sortable columns and declares the
                    searchable column set — requests for columns outside the whitelist fall back
                    to <code>created_at desc</code>.
                  </li>
                  <li>
                    <code>paginate.Result[T]</code> returns the canonical{' '}
                    <code>{'{'} data, meta: {'{'} total, page, page_size, pages {'}'} {'}'}</code>{' '}
                    envelope — matches the existing API response format exactly.
                  </li>
                </ul>

                <h3>Generator update</h3>
                <ul>
                  <li>
                    The emitted List handler now delegates to <code>paginate.List</code>. Every
                    generated resource gets the same clamping, whitelisting, and UUID-safe search
                    behavior — no per-resource drift.
                  </li>
                  <li>
                    Searchable column selection uses <code>IsSearchable()</code> (text / string /
                    slug / richtext only), so FK UUID columns are no longer accidentally included
                    in <code>ILIKE</code> search — a leftover rough edge from{' '}
                    <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer">issue #12</a>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.9.1 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.1
                </span>
                <span className="text-sm text-muted-foreground">April 24, 2026</span>
              </div>

              <div className="prose-grit">
                <p>
                  Patch release fixing compilation and consistency bugs in v3.9.0. Every
                  freshly scaffolded project (including <code>--mobile --desktop</code>)
                  and every <code>grit generate resource</code> run now produces Go code that builds
                  cleanly on the first try. Thanks to <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/9" target="_blank" rel="noopener noreferrer">issue #9</a>,
                  {' '}<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/10" target="_blank" rel="noopener noreferrer">#10</a>,
                  {' '}<a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/11" target="_blank" rel="noopener noreferrer">#11</a>,
                  {' '}and <a className="text-primary hover:underline" href="https://github.com/MUKE-coder/grit/issues/12" target="_blank" rel="noopener noreferrer">#12</a>.
                </p>

                <h3>Scaffold fixes</h3>
                <ul>
                  <li>
                    <strong>Missing imports</strong>: added <code>&quot;log&quot;</code> to <code>config.go</code>,{' '}
                    <code>&quot;gorm.io/gorm/logger&quot;</code> to <code>user.go</code>,{' '}
                    <code>&quot;net/http&quot;</code> to <code>middleware/logger.go</code>.
                  </li>
                  <li>
                    <strong>Stray package prefix</strong>: removed <code>handlers.</code> qualifier on{' '}
                    <code>IsTrustedDevice</code> (same-package call).
                  </li>
                  <li>
                    <strong>User ID type consistency</strong>: normalized <code>UserID</code> and{' '}
                    <code>UploadID</code> to <code>string</code> UUIDs across 2FA models, auth service,
                    TOTP handler (<code>c.GetString(&quot;user_id&quot;)</code> replaces{' '}
                    <code>c.GetUint</code>), jobs package, and upload handler.
                  </li>
                </ul>

                <h3>Desktop scaffold fixes</h3>
                <ul>
                  <li>
                    <code>keychain.go</code> moved from <code>internal/</code> to the top level (the
                    subdirectory file was declaring <code>package main</code>, which Go rejects).
                  </li>
                  <li>
                    <code>go.mod</code> module path fixed from{' '}
                    <code>&lt;project&gt;/apps/api/apps/desktop</code> to{' '}
                    <code>&lt;project&gt;/apps/desktop</code>.
                  </li>
                </ul>

                <h3>Resource generator fixes</h3>
                <ul>
                  <li>
                    Service signatures take <code>id string</code> instead of <code>id uint</code> --
                    matches the UUID string PK the models have always emitted.
                  </li>
                  <li>
                    Handler FK fields, handler M2M arrays, TS interface FK fields, and TanStack hook
                    ID types all switched to <code>string</code> (were <code>uint</code> /{' '}
                    <code>number</code>).
                  </li>
                  <li>
                    <strong>Initialism-aware</strong>{' '}
                    <code>toPascalCase</code> / <code>toSnakeCase</code>:{' '}
                    <code>owner_id</code> → <code>OwnerID</code> (was <code>OwnerId</code>),{' '}
                    <code>image_url</code> → <code>ImageURL</code> (was <code>ImageUrl</code>),{' '}
                    <code>api_key</code> → <code>APIKey</code>. Round-trips correctly
                    (snake → pascal → snake).
                  </li>
                  <li>
                    <strong>Zod schemas</strong> now emit snake_case field names matching the Go
                    handler&apos;s JSON tags (previously emitted camelCase, causing validation and
                    <code>ShouldBindJSON</code> mismatches).
                  </li>
                  <li>
                    Zod FK and M2M validators use <code>z.string().uuid()</code> instead of{' '}
                    <code>z.number().int()</code>.
                  </li>
                  <li>
                    FK columns generate with <code>gorm:&quot;size:36;index&quot;</code> (matches UUID PK width).
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.9.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.9.0
                </span>
                <span className="text-sm text-muted-foreground">April 15, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>New: <code>--desktop</code> flag</h3>
                <ul>
                  <li>
                    <strong>Desktop + mobile + API in one monorepo</strong> &mdash;{' '}
                    <code>grit new myapp --mobile --desktop</code> scaffolds a complete multi-client
                    SaaS: Go API shared by an Expo mobile app AND a Wails desktop app. All three
                    share the same <code>packages/shared</code> types and schemas.
                  </li>
                  <li>
                    <strong>Wails as a thin client</strong> &mdash; The new desktop app is a
                    frameless Wails window that calls the shared API over HTTP. No embedded Go
                    business logic, no local SQLite. Wails bindings are used only for native OS
                    features: window controls, file dialogs, and OS keychain (macOS Keychain,
                    Windows Credential Manager, Linux Secret Service) for JWT storage.
                  </li>
                  <li>
                    <strong>Distinct from <code>grit new-desktop</code></strong> &mdash; The
                    standalone offline-first desktop scaffold (<code>grit new-desktop</code>) is
                    unchanged. <code>--desktop</code> is a new, separate capability for
                    always-online multi-client apps.
                  </li>
                </ul>

                <h3>Premium Desktop UX</h3>
                <ul>
                  <li>
                    <strong>Platform-aware window chrome</strong> &mdash; macOS traffic lights on
                    the left, Windows/Linux controls on the right. Detected at runtime via
                    <code>GetPlatform()</code> Wails binding.
                  </li>
                  <li>
                    <strong>Command palette</strong> (<code>⌘K</code>) &mdash; every scaffolded
                    desktop app ships with a Raycast/Linear-style command palette. Searchable
                    navigation + actions with keyboard-first UX.
                  </li>
                  <li>
                    <strong>Fixed 240px sidebar</strong> &mdash; not collapsible. Desktop windows
                    are wide enough; collapse toggles are a web pattern.
                  </li>
                  <li>
                    <strong>Global keyboard shortcuts</strong> &mdash;{' '}
                    <code>useShortcuts()</code> hook with defaults: <code>⌘K</code> palette,{' '}
                    <code>⌘,</code> settings, <code>⌘L</code> logout, <code>Esc</code> to close.
                  </li>
                  <li>
                    <strong>More negative space</strong> &mdash; content padding is <code>32px</code>
                    (vs web&apos;s <code>24px</code>) for long focus sessions. Subtler shadows
                    (OS chrome already provides elevation).
                  </li>
                </ul>

                <h3>Style Guide</h3>
                <ul>
                  <li>
                    New <strong>§14.5 Desktop App Patterns</strong> section in
                    <code>GRIT_STYLE_GUIDE.md</code> covering window chrome, sidebar (not
                    collapsible), topbar, command palette, keyboard shortcuts, OS keychain
                    integration, typography (tighter than web), and do&apos;s &amp; don&apos;ts
                    (no breadcrumbs, no header banners, no web-style autoplay).
                  </li>
                </ul>

                <h3>Usage</h3>
                <pre><code>{`grit new myapp --mobile --desktop --next
# apps/api + apps/web + apps/expo + apps/desktop

grit new myapp --desktop --triple
# apps/api + apps/web + apps/admin + apps/desktop

grit new myapp --api --desktop
# apps/api + apps/desktop (minimal)`}</code></pre>
              </div>
            </div>

            {/* v3.8.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.8.0
                </span>
                <span className="text-sm text-muted-foreground">April 11, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Design System</h3>
                <ul>
                  <li>
                    <strong>GRIT_STYLE_GUIDE.md</strong> &mdash; First official style guide for all Grit-scaffolded
                    projects. Premium Minimal aesthetic (Linear / Vercel school), Grit purple{' '}
                    <code>#6C5CE7</code> primary, Onest font. Covers typography, color palette, spacing,
                    shadows, every component spec (buttons, inputs, cards, tables, modals), auth page rules,
                    CLI scaffolding design, admin panel patterns, email templates.
                  </li>
                </ul>

                <h3>Admin Layout</h3>
                <ul>
                  <li>
                    <strong>Topbar refactor</strong> &mdash; Moved sidebar collapse toggle to top-left of the
                    topbar (next to mobile menu button). Moved theme toggle, notifications bell, and enhanced
                    user menu to the top-right cluster alongside search. The sidebar now contains only
                    navigation. Matches modern dashboard patterns (Linear, Vercel, Raycast).
                  </li>
                  <li>
                    <strong>Enhanced user menu</strong> &mdash; Dropdown now shows User Activity, Settings,
                    Billing, and Log out sections with a user name/email header.
                  </li>
                </ul>

                <h3>PageHeader Component</h3>
                <ul>
                  <li>
                    <strong>Consistent page headers</strong> &mdash; New <code>&lt;PageHeader /&gt;</code>{' '}
                    component at <code>components/layout/page-header.tsx</code> with title, description,
                    breadcrumbs, actions slot, and a 4-card stats grid. Every generated resource page
                    auto-includes it.
                  </li>
                  <li>
                    <strong>Auto-generated stats cards</strong> &mdash; Resource pages now ship with 4 default
                    stat cards (Total, This Week, This Month, Updated Recently) fetched from the API. Override
                    via <code>defineResource({`{ stats: { cards: [...] } }`})</code> or disable with{' '}
                    <code>stats: false</code>.
                  </li>
                </ul>

                <h3>Auth Pages</h3>
                <ul>
                  <li>
                    <strong>New centered auth variant</strong> &mdash; <code>grit new myapp --style centered</code>{' '}
                    scaffolds Linear-school single-card auth pages (login, sign-up, forgot-password). ~420px
                    card on a subtle radial gradient background. The original split-screen design remains the
                    default (unchanged).
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.7.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.7.0
                </span>
                <span className="text-sm text-muted-foreground">April 3, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Security</h3>
                <ul>
                  <li>
                    <strong>Security headers middleware</strong> &mdash; New <code>SecurityHeaders()</code>{' '}
                    middleware adds X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Referrer-Policy,
                    Permissions-Policy, and HSTS (when HTTPS detected) on every response.
                  </li>
                  <li>
                    <strong>Max body size middleware</strong> &mdash; 10MB default limit returns 413 on exceed.
                  </li>
                  <li>
                    <strong>JWT secret validation</strong> &mdash; Warns if <code>JWT_SECRET</code> is shorter
                    than 32 characters.
                  </li>
                  <li>
                    <strong>Sentinel WAF</strong> &mdash; Now runs in <code>ModeBlock</code> in production{' '}
                    (was always <code>ModeLog</code>). Development keeps <code>ModeLog</code>.
                  </li>
                </ul>

                <h3>Performance</h3>
                <ul>
                  <li>
                    <strong>GORM AutoMigrate silence</strong> &mdash; Migration now uses a{' '}
                    <code>logger.Silent</code> session to suppress schema inspection SQL noise. Fixes{' '}
                    <a href="https://github.com/MUKE-coder/grit/issues/8" className="text-primary hover:underline"
                       target="_blank" rel="noopener noreferrer">issue #8</a>.
                  </li>
                </ul>

                <h3>Web App Auth</h3>
                <ul>
                  <li>
                    <strong>Auth pages for the web app</strong> &mdash; The web app (<code>apps/web</code>) now
                    ships with its own auth pages: login, register, forgot-password, OAuth callback.
                    Previously only the admin panel had auth. This is critical for e-commerce and SaaS where
                    end users log in on the web app, not the admin.
                  </li>
                  <li>
                    <strong><code>useAuth()</code> hook</strong> &mdash; React Query + js-cookie token
                    management with <code>AuthProvider</code> context wrapping the web app.
                  </li>
                </ul>

                <h3>Mobile (Expo)</h3>
                <ul>
                  <li>
                    <strong>Major Expo scaffold upgrade</strong> &mdash; 4 tabs (Home, Explore, Profile,
                    Settings) instead of 2. All forms use react-hook-form + zod. Home screen with stat cards
                    and pull-to-refresh. Explore screen with search and category discovery. Settings with
                    SectionList. Profile with display/edit mode.
                  </li>
                  <li>
                    <strong>OAuth in mobile</strong> &mdash; Google OAuth via <code>expo-web-browser</code>{' '}
                    with deep-link callback handling.
                  </li>
                  <li>
                    <strong>New Expo dependencies</strong> &mdash; react-hook-form, @hookform/resolvers, zod,
                    expo-image, expo-haptics, expo-web-browser. Splash screen config in{' '}
                    <code>app.json</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.6.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v3.6.0
                </span>
                <span className="text-sm text-muted-foreground">March 27, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Scaffold into current directory</strong> &mdash; <code>grit new .</code> and{' '}
                    <code>grit new ./</code> now scaffold into the current directory instead of creating a
                    subfolder. Infers the project name from the folder name. Also auto-detects when the
                    current directory name matches the project name.
                  </li>
                  <li>
                    <strong><code>--force</code> flag</strong> &mdash; Allows scaffolding into non-empty
                    directories. Useful when a repo was cloned first (with README, .git, LICENSE) before
                    scaffolding: <code>grit new . --triple --vite --force</code>.
                  </li>
                  <li>
                    <strong><code>--here</code> flag</strong> &mdash; Explicit alternative to{' '}
                    <code>grit new .</code> for in-place scaffolding.
                  </li>
                  <li>
                    <strong>30 standalone courses</strong> &mdash; Added 20 new courses to the learning
                    platform (42 total across 3 tracks + 20 standalone). Topics include testing, GORM
                    mastery, WebSockets, Stripe payments, blog/CMS, CI/CD, middleware, and the 100-component
                    UI registry.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Flags now skip interactive prompt</strong> &mdash; Running{' '}
                    <code>grit new myapp --triple --vite</code> no longer shows the architecture/frontend
                    selection prompt. Flags act as true shortcuts for non-interactive setup.
                  </li>
                  <li>
                    <strong>Module path upgrade to /v3</strong> &mdash; Fixed{' '}
                    <code>go install ...@latest</code> downloading v2.9.0 instead of v3.x.{' '}
                    All import paths updated from <code>grit/v2</code> to <code>grit/v3</code>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.5.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.5.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Documentation</h3>
                <ul>
                  <li>
                    <strong>Full docs redesign</strong> &mdash; Rebuilt the documentation site with a
                    Tailwind CSS-inspired aesthetic. New dark theme (<code>#0b1120</code>), sky-blue accents,
                    cleaner header with backdrop blur, redesigned code blocks with file tabs and line
                    highlighting, and new <code>StepWithCode</code> component for two-column step-by-step
                    guides (text left, code right).
                  </li>
                  <li>
                    <strong>Installation page redesigned</strong> &mdash; Step-numbered sections (01-04)
                    with the new two-column layout, system requirements table, architecture shortcuts,
                    and services grid.
                  </li>
                  <li>
                    <strong>Architecture Modes page</strong> &mdash; Visual cards for all 5 architectures
                    (single, double, triple, API only, mobile) with directory structure trees, features
                    list, ideal use cases, and frontend framework comparison.
                  </li>
                  <li>
                    <strong>TanStack Router guide</strong> &mdash; Complete guide for the TanStack Router
                    frontend option: project structure, routing patterns, comparison table with Next.js,
                    route examples, and admin panel auth guards.
                  </li>
                  <li>
                    <strong>New CLI Commands page</strong> &mdash; Documents <code>grit routes</code>,
                    <code>grit down/up</code> (maintenance mode), and <code>grit deploy</code>. Includes
                    complete command reference table for all 21 CLI commands.
                  </li>
                  <li>
                    <strong>Deploy Command guide</strong> &mdash; Step-by-step deployment pipeline with
                    systemd service unit and Caddyfile examples, flags table.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>Updated skill file with all v3.x architecture modes, frontend options, and new CLI commands.</li>
                  <li>Updated sidebar with new pages: Architecture Modes, New CLI Commands, TanStack Router, Deploy Command.</li>
                  <li>Frontend sidebar section renamed from &ldquo;Frontend (Next.js)&rdquo; to &ldquo;Frontend&rdquo; to reflect multi-framework support.</li>
                </ul>
              </div>
            </div>

            {/* v3.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-architecture code generator</strong> &mdash; <code>grit generate resource</code>
                    now works for all 5 architecture modes and both frontend frameworks. Generates Go model,
                    service, and handler at the correct path (<code>internal/</code> for single app,
                    <code>apps/api/internal/</code> for monorepo). Generates React Query hooks and admin
                    resource pages for both Next.js and TanStack Router.
                  </li>
                  <li>
                    <strong><code>grit.json</code> project manifest</strong> &mdash; Every scaffolded project
                    now includes a <code>grit.json</code> file at the root with <code>architecture</code> and
                    <code>frontend</code> fields. The generator reads this to determine correct file paths
                    and template variants, eliminating fragile filesystem heuristics.
                  </li>
                  <li>
                    <strong>TanStack Router resource generation</strong> &mdash; When generating resources
                    in a TanStack Router project, creates route files at
                    <code>src/routes/_dashboard/resources/</code> using <code>createFileRoute</code> instead
                    of Next.js <code>app/(dashboard)/resources/</code> page convention.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.3.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features (Goravel-Inspired)</h3>
                <ul>
                  <li>
                    <strong><code>grit routes</code></strong> &mdash; List all registered API routes in a
                    formatted table. Parses <code>routes.go</code> and shows method, path, handler, and
                    middleware group (public/protected/admin). Works for both monorepo and single app projects.
                  </li>
                  <li>
                    <strong><code>grit down</code> / <code>grit up</code></strong> &mdash; Maintenance mode.
                    <code>grit down</code> creates a <code>.maintenance</code> file that triggers the new
                    maintenance middleware, returning 503 for all requests. <code>grit up</code> removes
                    it and resumes normal operation.
                  </li>
                  <li>
                    <strong><code>grit deploy</code></strong> &mdash; One-command production deployment.
                    Cross-compiles for Linux, builds frontend, uploads binary via SCP, configures a systemd
                    service, and optionally sets up Caddy reverse proxy with auto-TLS. Supports
                    <code>--host</code>, <code>--domain</code>, <code>--key</code> flags or
                    <code>DEPLOY_HOST</code>/<code>DEPLOY_DOMAIN</code>/<code>DEPLOY_KEY_FILE</code> env vars.
                  </li>
                  <li>
                    <strong>Maintenance middleware</strong> &mdash; All scaffolded projects now include a
                    <code>Maintenance()</code> Gin middleware that checks for a <code>.maintenance</code>
                    file on every request. Runs as the first global middleware.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.2.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.2.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Single app architecture</strong> &mdash; <code>grit new my-app --single</code> creates
                    a single Go binary that serves both the API and an embedded React SPA. Uses <code>go:embed</code>
                    to bake the built frontend into the binary at compile time. One file to deploy. Dev mode runs
                    Go on <code>:8080</code> and Vite on <code>:5173</code> with API proxy.
                  </li>
                  <li>
                    <strong>Parameterized API paths</strong> &mdash; All Go API file generators now use
                    <code>opts.APIRoot()</code> and <code>opts.Module()</code> helpers, enabling the same
                    template functions to generate files for both monorepo (<code>apps/api/</code>) and
                    single app (project root) architectures.
                  </li>
                </ul>

                <h3>Single App Structure</h3>
                <ul>
                  <li><code>cmd/server/main.go</code> &mdash; Entry point with <code>go:embed frontend/dist/*</code> and SPA fallback routing</li>
                  <li><code>internal/</code> &mdash; Full Go backend (same as monorepo API)</li>
                  <li><code>frontend/</code> &mdash; React + Vite + TanStack Router SPA</li>
                  <li><code>Makefile</code> &mdash; <code>make dev</code> (parallel servers), <code>make build</code> (single binary)</li>
                </ul>
              </div>
            </div>

            {/* v3.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.1.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>TanStack Router frontend scaffold</strong> &mdash; When selecting
                    TanStack Router (Vite) as your frontend, both the web app and admin panel are
                    now fully scaffolded with Vite + TanStack Router + React Query + Tailwind CSS.
                    Includes file-based routing via <code>@tanstack/router-vite-plugin</code>,
                    API proxy in dev mode, and all the same features as the Next.js scaffold.
                  </li>
                  <li>
                    <strong>TanStack Router admin panel</strong> &mdash; Complete admin panel with
                    TanStack Router: auth pages (login, sign-up, forgot password), dashboard layout
                    with sidebar, resource management (users, blogs) via ResourcePage component,
                    system pages (jobs, files, cron, mail, security), profile page. All existing
                    React components (DataTable, FormBuilder, widgets) are reused with automatic
                    <code>&quot;use client&quot;</code> directive stripping.
                  </li>
                </ul>
              </div>
            </div>

            {/* v3.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v3.0.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Interactive project creation</strong> &mdash; <code>grit new my-app</code> now
                    launches an interactive prompt to select your architecture and frontend framework.
                    Power users can skip with flags: <code>--single --vite</code>, <code>--triple --next</code>,
                    <code>--api</code>, etc.
                  </li>
                  <li>
                    <strong>5 architecture modes</strong> &mdash; Choose the project structure that fits your team:
                    <strong>Single</strong> (Go API + embedded React SPA, one binary),
                    <strong>Double</strong> (Web + API Turborepo),
                    <strong>Triple</strong> (Web + Admin + API Turborepo),
                    <strong>API Only</strong> (Go backend, no frontend),
                    <strong>Mobile</strong> (API + Expo React Native).
                  </li>
                  <li>
                    <strong>Frontend framework choice</strong> &mdash; Pick between <strong>Next.js</strong> (SSR,
                    App Router) and <strong>TanStack Router</strong> (Vite, fast builds, small bundle, SPA).
                    Available for all architecture modes that include a frontend.
                  </li>
                </ul>

                <h3>Breaking Changes</h3>
                <ul>
                  <li>
                    <strong>Options struct refactored</strong> &mdash; The internal <code>Options</code> struct
                    now uses <code>Architecture</code> and <code>Frontend</code> enum fields instead of boolean
                    flags. Legacy flags (<code>--api</code>, <code>--mobile</code>, <code>--full</code>) still
                    work via the <code>Normalize()</code> migration layer.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.9.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.9.0
                </span>
                <span className="text-sm text-muted-foreground">March 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Two-Factor Authentication (TOTP)</strong> &mdash; Every <code>grit new</code> project
                    now includes a complete 2FA system with authenticator app support (Google Authenticator,
                    Authy, 1Password, etc.). Zero-dependency RFC 6238 implementation with HMAC-SHA1.
                    Includes setup flow with QR code URI generation, 6-digit code verification with
                    &plusmn;1 window clock skew tolerance, and seamless integration with the existing
                    JWT login flow.
                  </li>
                  <li>
                    <strong>Backup Codes</strong> &mdash; 10 one-time-use recovery codes generated when
                    enabling 2FA. Each code is individually bcrypt-hashed for storage. Codes can be
                    regenerated at any time (invalidates previous set). Use during login as an alternative
                    to the authenticator app.
                  </li>
                  <li>
                    <strong>Trusted Devices</strong> &mdash; &ldquo;Remember this device&rdquo; option
                    during TOTP verification. Sets an HttpOnly cookie with a SHA-256 hashed token stored
                    in the database. Trusted devices last 30 days with sliding expiry (refreshed on each use).
                    Users can revoke all trusted devices from their account.
                  </li>
                </ul>

                <h3>New Endpoints</h3>
                <ul>
                  <li><code>POST /api/auth/totp/setup</code> &mdash; Generate TOTP secret + QR URI (authenticated)</li>
                  <li><code>POST /api/auth/totp/enable</code> &mdash; Verify initial code and activate 2FA</li>
                  <li><code>POST /api/auth/totp/verify</code> &mdash; Verify TOTP code during login (public, uses pending token)</li>
                  <li><code>POST /api/auth/totp/backup-codes/verify</code> &mdash; Use backup code during login</li>
                  <li><code>POST /api/auth/totp/disable</code> &mdash; Disable 2FA (requires password)</li>
                  <li><code>GET /api/auth/totp/status</code> &mdash; Check 2FA status, remaining backup codes, trusted device count</li>
                  <li><code>POST /api/auth/totp/backup-codes</code> &mdash; Regenerate backup codes</li>
                  <li><code>DELETE /api/auth/totp/trusted-devices</code> &mdash; Revoke all trusted devices</li>
                </ul>
              </div>
            </div>

            {/* v2.8.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.8.0
                </span>
                <span className="text-sm text-muted-foreground">March 16, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Vercel AI Gateway integration</strong> &mdash; Replaced the multi-provider AI service
                    (Claude, OpenAI, Gemini with separate API implementations) with{' '}
                    <a href="https://vercel.com/ai-gateway" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Vercel AI Gateway</a>.
                    One API key now gives access to hundreds of models from all major providers through a single
                    OpenAI-compatible endpoint. Models use the <code>provider/model</code> format
                    (e.g. <code>anthropic/claude-sonnet-4-6</code>, <code>openai/gpt-5.4</code>,{' '}
                    <code>google/gemini-2.5-pro</code>). Includes automatic retries, fallbacks,
                    spend monitoring, and zero markup on tokens.
                  </li>
                </ul>

                <h3>Breaking Changes</h3>
                <ul>
                  <li>
                    <strong>AI environment variables</strong> &mdash; <code>AI_PROVIDER</code>,{' '}
                    <code>AI_API_KEY</code>, and <code>AI_MODEL</code> have been replaced with{' '}
                    <code>AI_GATEWAY_API_KEY</code>, <code>AI_GATEWAY_MODEL</code>, and{' '}
                    <code>AI_GATEWAY_URL</code>. Update your <code>.env</code> file accordingly.
                    Get your API key from{' '}
                    <a href="https://vercel.com/ai-gateway" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">vercel.com/ai-gateway</a>.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.7.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.7.0
                </span>
                <span className="text-sm text-muted-foreground">March 10, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>10 Official Plugins</strong> &mdash; New <code>grit-plugins</code> ecosystem with
                    drop-in Go packages for common functionality: WebSockets (<code>grit-websockets</code>),
                    Stripe payments (<code>grit-stripe</code>), OAuth social login (<code>grit-oauth</code>),
                    notifications (<code>grit-notifications</code>), full-text search (<code>grit-search</code>),
                    video processing (<code>grit-video</code>), WebRTC conferencing (<code>grit-conference</code>),
                    outgoing webhooks (<code>grit-webhooks</code>), i18n translations (<code>grit-i18n</code>),
                    and PDF/Excel/CSV export (<code>grit-export</code>). Each plugin includes a Claude Code
                    skill file for AI-assisted integration.
                  </li>
                  <li>
                    <strong>Claude Code Skills format</strong> &mdash; Updated the scaffolded AI skill file
                    from a monolithic <code>GRIT_SKILL.md</code> to the official Claude Code skills directory
                    structure (<code>.claude/skills/grit/SKILL.md</code> + <code>reference.md</code>) with
                    YAML frontmatter. AI assistants can now discover and use Grit conventions automatically.
                  </li>
                  <li>
                    <strong>Grit UI component registry (100 components)</strong> &mdash; Expanded from 91 to
                    100 pre-built components across 5 categories: marketing (21), auth (10), SaaS (30),
                    ecommerce (20), and layout (20).
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/plugins" className="text-primary hover:underline">Plugins</Link>{' '}
                    page &mdash; overview of all 10 plugins with installation, environment setup,
                    quick start code, features, and use cases for each.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.6.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.6.0
                </span>
                <span className="text-sm text-muted-foreground">March 6, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>GORM Studio (Desktop)</strong> &mdash; Replaced the broken custom HTML studio with the real{' '}
                    <code>gorm-studio</code> package. Desktop studio now runs on port 8080 at <code>/studio</code> using
                    Gin + gorm-studio, matching the web scaffold. Auto-opens browser on launch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.5.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v2.5.0
                </span>
                <span className="text-sm text-muted-foreground">March 6, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>GRIT_SKILL.md</strong> &mdash; Desktop scaffolds now include a <code>GRIT_SKILL.md</code> file
                    in the project root. This is a comprehensive AI reference (12 sections) covering architecture,
                    CLI commands, resource generation, field types, code markers, golden rules, and common LLM mistakes
                    &mdash; so AI assistants can work with the project correctly out of the box.
                  </li>
                  <li>
                    <strong>Comprehensive README</strong> &mdash; The scaffolded <code>README.md</code> now includes a
                    full project walkthrough, &ldquo;Adding a New Module&rdquo; guide, supported field types table,
                    customization section (window size, title bar, database, app name), code markers reference,
                    and a ready-to-use AI prompt for building a Task Manager app.
                  </li>
                </ul>

                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Dashboard stats cache</strong> &mdash; Dashboard statistics now update immediately after
                    creating a blog or contact. Changed query keys from <code>[&quot;blogs-stats&quot;]</code> to{' '}
                    <code>[&quot;blogs&quot;, &quot;stats&quot;]</code> so TanStack Query{`'`}s prefix matching
                    invalidates dashboard queries when resources are created or deleted.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-accent/15 px-3 py-1 text-sm font-semibold text-primary">
                  v2.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 5, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Window controls on auth pages</strong> &mdash; Login and register pages now include
                    minimize, maximize, and close buttons with a draggable title area, so users can move and
                    manage the window before signing in.
                  </li>
                  <li>
                    <strong>Show/hide password toggle</strong> &mdash; All password fields on login and register
                    pages now have an eye icon toggle to reveal or hide the password text.
                  </li>
                </ul>

                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop build script</strong> &mdash; Removed <code>tsc</code> from the frontend
                    build script. TanStack Router{`'`}s Vite plugin generates <code>routeTree.gen.ts</code> during
                    the Vite build, so running <code>tsc</code> before Vite caused{' '}
                    <code>Cannot find module {`'`}./routeTree.gen{`'`}</code> errors.
                  </li>
                  <li>
                    <strong>Title bar import path</strong> &mdash; Fixed the Wails binding import in{' '}
                    <code>title-bar.tsx</code> from a 2-level to 3-level relative path.
                  </li>
                  <li>
                    <strong>Auth hook file extension</strong> &mdash; Renamed <code>use-auth.ts</code> to{' '}
                    <code>use-auth.tsx</code> so TypeScript handles the JSX correctly.
                  </li>
                  <li>
                    <strong>Create resource cache refresh</strong> &mdash; Blog and contact create pages now
                    invalidate the React Query cache before navigating back, so new records appear in the
                    table immediately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.2.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.2.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>Desktop auth hook file extension</strong> &mdash; Renamed the scaffolded{' '}
                    <code>use-auth.ts</code> to <code>use-auth.tsx</code> so TypeScript correctly handles
                    the JSX in <code>&lt;AuthContext.Provider&gt;</code>. Previously, <code>grit new-desktop</code>{' '}
                    projects would fail to compile with <code>TS1005: {`'>'`} expected</code> errors.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added Desktop Handbook PDF download links to all 8 desktop documentation pages.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.1.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>TanStack Router for desktop</strong> &mdash; Migrated the desktop frontend from
                    React Router to{' '}
                    <a href="https://tanstack.com/router" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">TanStack Router</a>{' '}
                    with file-based routing. Routes are auto-discovered by the Vite plugin &mdash; no centralized
                    route registry. Uses <code>createHashHistory()</code> for Wails compatibility and{' '}
                    <code>Route.useParams()</code> for type-safe params. Resource generation now creates 5 files
                    (list, new, edit routes + model + service) and performs 10 injections (down from 12).
                  </li>
                  <li>
                    <strong>Mobile navigation</strong> &mdash; Added a hamburger menu to the docs site header,
                    visible below the <code>lg</code> breakpoint. Opens a Sheet sidebar with all navigation links.
                    Auto-closes on link click.
                  </li>
                  <li>
                    <strong>CGO-free SQLite</strong> &mdash; Replaced <code>gorm.io/driver/sqlite</code> (requires
                    CGO) with <code>github.com/glebarez/sqlite</code> (pure Go) in all scaffold templates. Desktop
                    apps now build and run without CGO or a C compiler.
                  </li>
                  <li>
                    <strong>20 Desktop Project Ideas</strong> &mdash; New{' '}
                    <Link href="/docs/desktop/project-ideas" className="text-primary hover:underline">project ideas page</Link>{' '}
                    with 20 ready-to-build desktop app ideas across business, education, healthcare, logistics,
                    and more. Each includes resources, field definitions, and <code>grit generate</code> commands.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    Added TanStack Router explanations to all desktop doc pages: overview, getting started,
                    first app, resource generation, and POS app.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/desktop/llm-reference" className="text-primary hover:underline">LLM Reference</Link>,{' '}
                    <Link href="/docs/ai-skill" className="text-primary hover:underline">GRIT_SKILL.md</Link>, and
                    database docs to reflect TanStack Router and CGO-free SQLite changes.
                  </li>
                </ul>
              </div>
            </div>

            {/* v2.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v2.0.0
                </span>
                <span className="text-sm text-muted-foreground">March 4, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Native desktop apps (Wails)</strong> &mdash; New <code>grit new-desktop</code> command
                    scaffolds a complete desktop application with Go backend, React frontend (Vite + TanStack Router +
                    TanStack Query), SQLite database, JWT authentication, blog and contact CRUD, PDF/Excel export,
                    custom title bar, dark theme, and GORM Studio. Compiles to a single native executable for
                    Windows, macOS, and Linux. See{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop docs</Link>.
                  </li>
                  <li>
                    <strong>Desktop resource generation</strong> &mdash; <code>grit generate resource</code> now
                    works inside desktop projects. Generates Go model, service, and TanStack Router route files
                    (list, new, edit), then injects code into 10 locations (db.go, main.go, app.go, types.go,
                    sidebar.tsx, studio/main.go) using <code>grit:</code> markers. See{' '}
                    <Link href="/docs/desktop/resource-generation" className="text-primary hover:underline">Desktop Resource Generation</Link>.
                  </li>
                  <li>
                    <strong>Project type auto-detection</strong> &mdash; All CLI commands now auto-detect whether
                    you are inside a web (Turborepo) or desktop (Wails) project. No flags needed.
                  </li>
                  <li>
                    <strong><code>grit start</code> for desktop</strong> &mdash; Running <code>grit start</code>{' '}
                    inside a desktop project launches <code>wails dev</code> with hot-reload for both Go and React.
                  </li>
                  <li>
                    <strong><code>grit compile</code></strong> &mdash; New command that runs <code>wails build</code>{' '}
                    to produce a distributable native binary.
                  </li>
                  <li>
                    <strong><code>grit studio</code></strong> &mdash; New command that launches GORM Studio. For
                    desktop projects it starts a standalone server on port 4000. For web projects it opens the
                    browser to the embedded Studio route.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code> for desktop</strong> &mdash; Removes a previously
                    generated desktop resource, deleting files and reversing all 10 marker injections.
                  </li>
                  <li>
                    <strong>Grit UI component registry (91 components)</strong> &mdash; Every scaffolded web
                    project now includes a shadcn-compatible component registry with 91 pre-built components
                    across 5 categories: marketing (14), auth (10), SaaS (30), ecommerce (20), and layout (18).
                    Install via <code>npx shadcn@latest add</code> from <code>/r</code> endpoints.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/desktop" className="text-primary hover:underline">Desktop (Wails)</Link>{' '}
                    section &mdash; 8 pages covering overview, getting started, first app tutorial, POS app
                    tutorial, resource generation, building/distribution, project ideas, and LLM reference.
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">LLM Reference</Link>{' '}
                    with complete desktop section: project structure, CLI commands, markers, and architecture comparison.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.4.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.4.0
                </span>
                <span className="text-sm text-muted-foreground">March 2, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Gzip response compression</strong> &mdash; All API responses are now
                    compressed automatically via a custom <code>Gzip()</code> middleware using the Go
                    standard library <code>compress/gzip</code> at <code>BestSpeed</code>.
                    JSON payloads shrink by 60–80%, reducing bandwidth on paginated list endpoints
                    with zero external dependencies.
                  </li>
                  <li>
                    <strong>Request ID tracing</strong> &mdash; A <code>RequestID()</code> middleware
                    injects a unique <code>X-Request-ID</code> header on every request (echoes the
                    upstream header or generates a nanosecond-based ID). The ID is stored in Gin
                    context and included in every structured log line for end-to-end request tracing.
                  </li>
                  <li>
                    <strong>Database connection pool tuning</strong> &mdash; The scaffold now sets
                    four GORM pool parameters: <code>MaxIdleConns(10)</code>,{' '}
                    <code>MaxOpenConns(100)</code>, <code>ConnMaxLifetime(30m)</code>, and{' '}
                    <code>ConnMaxIdleTime(10m)</code>. This prevents stale connections after network
                    interruptions and avoids connection exhaustion under load.
                  </li>
                  <li>
                    <strong>Cache-Control headers on public blog endpoints</strong> &mdash; The{' '}
                    <code>ListPublished</code> handler now returns{' '}
                    <code>Cache-Control: public, max-age=300</code> (5 minutes) and{' '}
                    <code>GetBySlug</code> returns <code>Cache-Control: public, max-age=3600</code>{' '}
                    (1 hour). CDNs and edge caches can now serve public blog content without hitting
                    the Go API.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/concepts/performance" className="text-primary hover:underline">Performance</Link>{' '}
                    page &mdash; comprehensive guide to all backend (Go/API) and frontend (Next.js)
                    performance optimisations that ship with every Grit project out of the box.
                    Covers Gzip, Request ID, connection pool, Cache-Control, presigned uploads,
                    background jobs, Redis caching, Server Components, ISR, React Query, next/image,
                    Turborepo, and code splitting.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/ai-skill/llm-guide" className="text-primary hover:underline">Complete LLM Reference</Link>{' '}
                    page &mdash; a dedicated machine-readable guide that teaches AI assistants
                    everything about Grit: project structure, all CLI commands, every field type,
                    code patterns, API response format, code markers, naming conventions, all
                    batteries, performance features, and the golden rules that must never be broken.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.3.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.3.0
                </span>
                <span className="text-sm text-muted-foreground">February 26, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Presigned URL uploads</strong> &mdash; File uploads now bypass the API server entirely.
                    The browser gets a presigned PUT URL, uploads directly to S3/R2/MinIO, then records the upload
                    in the database. This fixes file uploads breaking behind reverse proxies (Dokploy/Traefik/Nginx)
                    due to request body size limits and timeouts. Includes progress tracking via XHR.
                  </li>
                  <li>
                    <strong>Error pages for scaffolded apps</strong> &mdash; New <code>grit new</code> projects now include{' '}
                    <code>error.tsx</code>, <code>not-found.tsx</code>, and <code>global-error.tsx</code> for both
                    admin and web apps. Errors are displayed with styled UI instead of the default Next.js error page.
                  </li>
                  <li>
                    <strong>Production-ready Docker config</strong> &mdash; <code>docker-compose.prod.yml</code> now uses{' '}
                    <code>expose</code> instead of <code>ports</code>, <code>env_file</code> for secrets, MinIO service,
                    named bridge network, build args for <code>NEXT_PUBLIC_API_URL</code>, and Go 1.24.
                  </li>
                  <li>
                    <strong>Sentinel ExcludePaths</strong> &mdash; Pulse, GORM Studio, Sentinel, and API docs paths are
                    now excluded from rate limiting by default, fixing Pulse health checks triggering rate limits.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/create-without-docker" className="text-primary hover:underline">Create without Docker</Link>{' '}
                    guide &mdash; set up a Grit project using Neon, Upstash, Cloudflare R2, and Resend instead of Docker.
                  </li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>
                    Scaffold Dockerfile updated from Go 1.23 to Go 1.24
                  </li>
                  <li>
                    Next.js Dockerfile now accepts <code>NEXT_PUBLIC_API_URL</code> as a build argument
                  </li>
                  <li>
                    <code>.env</code> template includes Docker Compose production variables (<code>POSTGRES_USER</code>,{' '}
                    <code>POSTGRES_PASSWORD</code>, <code>POSTGRES_DB</code>, <code>API_URL</code>)
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.1.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.1.0
                </span>
                <span className="text-sm text-muted-foreground">February 25, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Default font changed to Onest</strong> &mdash; New projects scaffolded with{' '}
                    <code>grit new</code> now use the{' '}
                    <a href="https://fonts.google.com/specimen/Onest" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Onest</a>{' '}
                    Google Font for all UI text instead of DM Sans. JetBrains Mono remains the code font.
                    The font is loaded via <code>next/font/google</code> with weights 400, 500, 600, and 700.
                  </li>
                  <li>
                    <strong>Hire Us page</strong> &mdash; New{' '}
                    <Link href="/hire" className="text-primary hover:underline">/hire</Link>{' '}
                    page for professional Grit development services. Includes service offerings,
                    tech stack overview, and contact CTA.
                  </li>
                  <li>
                    <strong>Monetization banners</strong> &mdash; Docs sidebar now shows promotional cards for{' '}
                    <a href="https://gritcms.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">GritCMS</a>,
                    developer hiring services, and{' '}
                    <Link href="/donate" className="text-primary hover:underline">donations</Link>{' '}
                    &mdash; visible on every documentation page.
                  </li>
                  <li>
                    <strong>Grit Fullstack Course page</strong> &mdash; New{' '}
                    <Link href="/course" className="text-primary hover:underline">/course</Link>{' '}
                    page with a 10-module curriculum covering Go, React, Next.js, and the full Grit stack.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    Top navigation now includes GritCMS, Hire Us, and a Sponsor heart icon for quick access
                    to all revenue channels.
                  </li>
                  <li>
                    <code>richtext</code> added to the FieldType union for better type safety in the code generator.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>OAuth callback fix</strong> &mdash; Fixed <code>TokenPair</code> struct field access
                    in the social login callback handler (was using map indexing instead of struct fields).
                  </li>
                  <li>
                    <strong>Course waitlist fix</strong> &mdash; Fixed Google Sheets submission to use
                    form-encoded data instead of JSON.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/getting-started/cli-cheatsheet" className="text-primary hover:underline">CLI Cheatsheet</Link>{' '}
                    page &mdash; complete reference for all Grit CLI commands with flags, field types,
                    generated files, common workflows, and full command tree.
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/oauth" className="text-primary hover:underline">Social Login (OAuth2)</Link>{' '}
                    setup guide for Google and GitHub authentication.
                  </li>
                  <li>
                    Updated Docker Cheat Sheet with force remove commands for containers and volumes.
                  </li>
                  <li>
                    Updated AI skill guide with social login (OAuth2) section.
                  </li>
                </ul>
              </div>
            </div>

            {/* v1.0.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v1.0.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Social Login (Google + GitHub)</strong> &mdash; Every <code>grit new</code> project now
                    includes OAuth2 social authentication via{' '}
                    <a href="https://github.com/markbates/goth" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Gothic</a>.
                    Users can sign in with Google or GitHub on all auth pages (login, register, admin).
                    Accounts are linked by email &mdash; existing users who sign in with a social provider are automatically connected.
                    Configurable via <code>GOOGLE_CLIENT_ID</code>, <code>GITHUB_CLIENT_ID</code> environment variables.
                  </li>
                  <li>
                    <strong>GORM Studio v1.0.1</strong> &mdash; Updated to the first stable tagged release of GORM Studio.
                  </li>
                </ul>

                <h3>Improvements</h3>
                <ul>
                  <li>
                    User model now includes <code>Provider</code>, <code>GoogleID</code>, and <code>GithubID</code> fields
                    for social account linking. Password field is now nullable to support OAuth-only accounts.
                  </li>
                  <li>
                    Admin users table shows Provider column with badges (Email, Google, GitHub) and new filter option.
                  </li>
                  <li>
                    Social login buttons (Google + GitHub) appear on all 4 admin style variants (default, modern, minimal, glass).
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.19.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.19.0
                </span>
                <span className="text-sm text-muted-foreground">February 24, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Fixes</h3>
                <ul>
                  <li>
                    <strong>gin-docs AuthConfig</strong> &mdash; Updated scaffold template to use the new{' '}
                    <code>gindocs.AuthConfig</code> struct instead of the deprecated <code>gindocs.AuthBearer</code> constant,
                    fixing compilation errors in newly scaffolded projects.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/tutorials/contact-app" className="text-primary hover:underline">Your First App</Link>{' '}
                    tutorial &mdash; step-by-step Contact Manager guide covering project setup, resource generation, and CRUD
                  </li>
                  <li>
                    New{' '}
                    <Link href="/docs/deployment/dokploy" className="text-primary hover:underline">Dokploy Deployment</Link>{' '}
                    guide with Dockerfile examples
                  </li>
                  <li>
                    Improved terminal blocks across all tutorials with copy buttons and horizontal scroll
                  </li>
                  <li>
                    Updated{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link>{' '}
                    page to reflect the new <code>AuthConfig</code> struct format
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.18.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.18.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Pulse (Observability)</strong> &mdash; Every <code>grit new</code> project now includes{' '}
                    <a href="https://github.com/MUKE-coder/pulse" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">Pulse</a>,
                    a self-hosted observability SDK. Provides request tracing, database monitoring, runtime metrics,
                    error tracking, health checks, alerting, Prometheus export, and an embedded React dashboard
                    at <code>/pulse</code>. Enabled by default, configurable via <code>PULSE_ENABLED</code>.
                    See <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse docs</Link>.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/pulse" className="text-primary hover:underline">Pulse (Observability)</Link> page
                    covering configuration, endpoints, health checks, alerting, Prometheus metrics, and data storage
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.17.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.17.0
                </span>
                <span className="text-sm text-muted-foreground">February 22, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>API Documentation (gin-docs)</strong> &mdash; Replaced hand-written Scalar/OpenAPI
                    spec with{' '}
                    <a href="https://github.com/MUKE-coder/gin-docs" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">gin-docs</a>,
                    a zero-annotation API documentation generator. Routes and GORM models are introspected
                    automatically to produce an OpenAPI 3.1 spec with interactive Scalar or Swagger UI,
                    plus Postman and Insomnia export.
                  </li>
                  <li>
                    <strong>Dark/Light mode for Go Playground</strong> &mdash; The playground now follows the
                    site-wide theme toggle, switching between VS Code dark and light CodeMirror themes.
                  </li>
                  <li>
                    <strong>Umami Analytics</strong> &mdash; Optional visitor analytics via self-hosted Umami,
                    configured with <code>NEXT_PUBLIC_UMAMI_WEBSITE_ID</code> environment variable.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    New{' '}
                    <Link href="/docs/backend/api-docs" className="text-primary hover:underline">API Documentation</Link> page
                    covering gin-docs configuration, GORM model schemas, route customization, UI switching, and spec export
                  </li>
                  <li>Full SEO + AEO implementation: sitemap, robots.txt, JSON-LD structured data, per-page metadata</li>
                </ul>

                <h3>Infrastructure</h3>
                <ul>
                  <li>Added Dockerfile for docs site deployment (Next.js standalone output)</li>
                  <li>Google Search Console verification</li>
                </ul>
              </div>
            </div>

            {/* v0.16.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.16.0
                </span>
                <span className="text-sm text-muted-foreground">February 21, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Go Playground</strong> &mdash; Interactive code editor at{' '}
                    <Link href="/playground" className="text-primary hover:underline">/playground</Link> with
                    Go syntax highlighting, code execution via the official Go Playground API, example snippets,
                    share links, and keyboard shortcuts (Ctrl+Enter to run).
                  </li>
                  <li>
                    <strong>GORM Studio updated</strong> &mdash; Updated to latest version with raw SQL editor,
                    schema export (SQL/JSON/YAML/DBML/ERD), data import/export (JSON/CSV/SQL/XLSX),
                    and Go model generation from database schema.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>
                    <Link href="/docs/prerequisites/golang" className="text-primary hover:underline">Go for Grit Developers</Link> &mdash;
                    comprehensive rewrite with 22 sections covering methods, Gin routing, middleware, CORS,
                    handler/service architecture, GORM CRUD, migrations, seeding, JWT auth flow, and RBAC
                  </li>
                  <li>Fixed right-side table of contents for the Go prerequisites page</li>
                  <li>New Middleware and CORS sections added to Go guide</li>
                </ul>
              </div>
            </div>

            {/* v0.15.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.15.0
                </span>
                <span className="text-sm text-muted-foreground">February 20, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Security (Sentinel)</strong> &mdash; Every <code>grit new</code> project now ships with
                    a production-grade security suite powered by{' '}
                    <Link href="https://github.com/MUKE-coder/sentinel" className="text-primary hover:underline">Sentinel</Link>.
                    Includes WAF, rate limiting, brute-force protection, anomaly detection, IP geolocation,
                    security headers, and a real-time threat dashboard at <code>/sentinel/ui</code>.
                    See <Link href="/docs/batteries/security" className="text-primary hover:underline">Security docs</Link>.
                  </li>
                  <li>
                    <strong>Admin security page</strong> &mdash; New System &rarr; Security page in the admin panel
                    embeds the Sentinel dashboard for monitoring threats without leaving the admin UI.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/batteries/security" className="text-primary hover:underline">Security (Sentinel)</Link> documentation page</li>
                  <li>Migrated getting-started pages (Installation, Quick Start, Troubleshooting) to use CodeBlock component</li>
                  <li>Added prerequisite learning pages for Go, Next.js, and Docker</li>
                </ul>
              </div>
            </div>

            {/* v0.14.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.14.0
                </span>
                <span className="text-sm text-muted-foreground">February 18, 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Multi-step forms</strong> &mdash; New <code>formView: &quot;modal-steps&quot;</code> and{' '}
                    <code>&quot;page-steps&quot;</code> variants with horizontal/vertical step indicators,
                    per-step validation, progress bar, and clickable step navigation.
                    See <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link>.
                  </li>
                  <li>
                    <strong>Standalone component usage</strong> &mdash; FormBuilder, FormStepper, and DataTable
                    can now be used on any page in both web and admin apps without the resource system.
                    See <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link>.
                  </li>
                  <li>
                    <strong>Richtext field type</strong> &mdash; New <code>richtext</code> field with Tiptap WYSIWYG
                    editor (bold, italic, headings, lists, code blocks, links, undo/redo).
                  </li>
                  <li>
                    <strong><code>string_array</code> field type</strong> &mdash; Store arrays of strings
                    using <code>datatypes.JSONSlice[string]</code>. Works with PostgreSQL and SQLite.
                    Maps to <code>string[]</code> in TypeScript and <code>z.array(z.string())</code> in Zod.
                  </li>
                  <li>
                    <strong>Built-in blog example</strong> &mdash; <code>grit new</code> now scaffolds a complete
                    blog with model, service, handler, seed data, public web pages, and admin resource definition.
                  </li>
                  <li>
                    <strong>Sidebar user avatar</strong> &mdash; Admin sidebar shows the current user&apos;s avatar
                    with a dropdown menu for profile and logout.
                  </li>
                  <li>
                    <strong>Profile avatar upload</strong> &mdash; Profile page now supports avatar image upload.
                  </li>
                  <li>
                    <strong><code>react-hook-form</code> in web app</strong> &mdash; Web app scaffold now includes{' '}
                    <code>react-hook-form</code> as a dependency, enabling standalone FormBuilder usage out of the box.
                  </li>
                </ul>

                <h3>Bug Fixes</h3>
                <ul>
                  <li>
                    <strong>Scalar API docs crash</strong> &mdash; Fixed <code>c.String</code> treating HTML as
                    a format string. Now uses <code>c.Data</code> to avoid panics when Scalar HTML
                    contains <code>%</code> characters in CSS/JS.
                  </li>
                  <li>
                    <strong>Blog route conflict</strong> &mdash; Admin blog CRUD routes moved
                    from <code>/api/blogs</code> to <code>/api/admin/blogs</code> to avoid conflict
                    with public blog routes.
                  </li>
                  <li>
                    <strong>Select dropdown styling</strong> &mdash; Fixed relationship select dropdown
                    rendering behind modals using portal-based positioning.
                  </li>
                </ul>

                <h3>Documentation</h3>
                <ul>
                  <li>New: <Link href="/docs/tutorials/product-catalog" className="text-primary hover:underline">Build a Product Catalog</Link> tutorial &mdash; resource generation, multi-step forms, standalone DataTable &amp; FormBuilder</li>
                  <li>New: <Link href="/docs/admin/multi-step-forms" className="text-primary hover:underline">Multi-Step Forms</Link> guide</li>
                  <li>New: <Link href="/docs/admin/standalone-usage" className="text-primary hover:underline">Standalone Usage</Link> guide</li>
                  <li>New: Changelog page</li>
                  <li>Updated CLI Commands, Code Generation, Quick Start, Resources, Shared Package, Web App, Seeders, and Forms pages</li>
                </ul>
              </div>
            </div>

            {/* v0.12.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.12.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Relationship support</strong> &mdash; New <code>belongs_to</code> and{' '}
                    <code>many_to_many</code> field types for the code generator. Automatically creates
                    foreign keys, junction tables, and relationship-aware form fields.
                  </li>
                  <li>
                    <strong>Relationship select fields</strong> &mdash; New <code>relationship-select</code> and{' '}
                    <code>multi-relationship-select</code> form field components with search, portal-based
                    dropdowns, and tag-based multi-select.
                  </li>
                  <li>
                    <strong>Beginner tutorial</strong> &mdash; &quot;Learn Grit Step by Step&quot; tutorial
                    walking through building a full-stack app from scratch.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.11.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.11.0
                </span>
                <span className="text-sm text-muted-foreground">February 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Full-page form view</strong> &mdash; New <code>formView: &quot;page&quot;</code> option
                    renders forms as dedicated pages instead of modals.
                  </li>
                  <li>
                    <strong><code>slug</code> field type</strong> &mdash; Auto-generates URL-friendly slugs with
                    unique suffixes. Excluded from create/update forms and Zod schemas.
                  </li>
                  <li>
                    <strong>DataTable column customization</strong> &mdash; Hide/show columns, column visibility
                    toggle in table toolbar.
                  </li>
                  <li>
                    <strong><code>grit start</code> commands</strong> &mdash; <code>grit start client</code> and{' '}
                    <code>grit start server</code> for running frontend and API separately.
                  </li>
                </ul>
              </div>
            </div>

            {/* v0.10.0 */}
            <div className="mb-12">
              <div className="flex items-center gap-3 mb-4">
                <span className="inline-flex items-center rounded-lg bg-muted/50 px-3 py-1 text-sm font-semibold text-muted-foreground">
                  v0.10.0
                </span>
                <span className="text-sm text-muted-foreground">January 2026</span>
              </div>

              <div className="prose-grit">
                <h3>Features</h3>
                <ul>
                  <li>
                    <strong>Style variants</strong> &mdash; <code>--style</code> flag for <code>grit new</code> with
                    4 admin panel styles: default, modern, minimal, and glass.
                  </li>
                  <li>
                    <strong>Air hot reloading</strong> &mdash; Go API development with automatic rebuild on file
                    changes using Air.
                  </li>
                  <li>
                    <strong><code>grit remove resource</code></strong> &mdash; Remove a generated resource and
                    clean up all injected code (model, handler, routes, schemas, types, hooks, admin pages).
                  </li>
                  <li>
                    <strong>AI workflow docs</strong> &mdash; Guides for using Grit with Claude and Antigravity AI assistants.
                  </li>
                </ul>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Introduction
                </Link>
              </Button>
              <div />
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
