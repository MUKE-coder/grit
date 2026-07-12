import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/backend/webhooks')

export default function WebhooksPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Backend</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Webhooks
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Stripe, GitHub, Twilio, WhatsApp &mdash; anything that pings you.
                Grit&apos;s <code>webhooks</code> package gives you{' '}
                <strong>one route</strong>, <code>POST /webhooks/:provider</code>,
                that verifies the signature, deduplicates retries, persists every
                delivery, and dispatches to a handler you register. Built-in
                verifiers cover Stripe, GitHub, and generic HMAC; adding your own is
                a single function.
              </p>
            </div>

            <div className="prose-grit">
              {/* Mental model */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  The receive pipeline
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The HTTP handler is already wired at{' '}
                  <code>POST /webhooks/:provider</code>. It routes by the{' '}
                  <code>:provider</code> path param to whatever you registered, and
                  runs the same seven steps for every source:
                </p>

                <CodeBlock
                  language="text"
                  filename="POST /webhooks/:provider"
                  code={`  provider ping
       │
       ▼
  1. LookupProvider(":provider")     → 404 if unregistered
  2. read raw body + flatten headers
  3. Provider.Verify(secret, body, headers)   → 401 on signature mismatch
  4. Provider.Extract(body, headers)          → (eventType, externalID)
  5. INSERT webhook_events                     ┐ unique(provider, external_id)
       │  duplicate? → status=skipped, 200 ────┘ (retry becomes a no-op)
  6. webhooks.Dispatch(ctx, event)   → your On() handler runs
  7. UPDATE status = processed | failed (+ handler_error, processed_at)
       │
       ▼
  always 200 to a verified+stored event (so the provider stops retrying)`}
                />

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Why always 200?</strong> Once
                    a webhook is verified and stored, the provider&apos;s job is
                    done &mdash; retrying wouldn&apos;t help. A <em>handler</em>{' '}
                    failure is recorded (<code>status=failed</code>) and recovered via
                    the admin replay endpoint, not by making the provider redeliver
                    forever. Only signature/parse failures return 4xx.
                  </p>
                </div>
              </div>

              {/* Registering */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Registering a provider
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Two calls at app boot: <code>webhooks.Setup(db)</code> once (already
                  done in <code>routes.Setup</code>), then a{' '}
                  <code>webhooks.Register</code> per source plus a{' '}
                  <code>webhooks.On</code> per event you care about. A{' '}
                  <code>Provider</code> is just its secret env var, a verifier, and an
                  extractor:
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/webhooks/setup.go"
                  code={`func SetupProviders() {
    webhooks.Register("stripe", webhooks.Provider{
        SecretEnv: "STRIPE_WEBHOOK_SECRET",
        Verify:    webhooks.StripeVerifier,
        Extract:   webhooks.StripeExtractor,
    })

    // Bind a handler to (provider, eventType).
    webhooks.On("stripe", "invoice.paid", func(ctx context.Context, e *models.WebhookEvent) error {
        // e.Payload is the raw JSON body — unmarshal what you need.
        return fulfillInvoice(ctx, e.Payload)
    })

    // "" is a catch-all: runs for any stripe event without a specific handler.
    webhooks.On("stripe", "", func(ctx context.Context, e *models.WebhookEvent) error {
        log.Printf("unhandled stripe event: %s", e.EventType)
        return nil
    })
}`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  Dispatch prefers an exact <code>(provider, eventType)</code> match
                  and falls back to the catch-all <code>&quot;&quot;</code> handler.
                  If nothing is registered the event is still persisted &mdash; it
                  just sits at <code>status=processed</code> with no side effect, so
                  you never silently lose a delivery.
                </p>
              </div>

              {/* Built-in verifiers */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Built-in verifiers
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Signature verification is the whole point &mdash; it proves the
                  request actually came from the provider and wasn&apos;t forged.
                  Grit ships three verifiers and matching extractors:
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Verifier</th>
                        <th className="px-4 py-2 font-medium">Header &amp; scheme</th>
                        <th className="px-4 py-2 font-medium">Extractor</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr className="border-b border-border/50">
                        <td className="px-4 py-2 font-mono text-[13px]">StripeVerifier</td>
                        <td className="px-4 py-2 text-muted-foreground"><code>Stripe-Signature</code>: <code>t=…,v1=…</code>, HMAC-SHA256 of <code>&quot;{'{ts}'}.{'{body}'}&quot;</code>, 5-min replay tolerance</td>
                        <td className="px-4 py-2 font-mono text-[13px]">StripeExtractor</td>
                      </tr>
                      <tr className="border-b border-border/50">
                        <td className="px-4 py-2 font-mono text-[13px]">GitHubVerifier</td>
                        <td className="px-4 py-2 text-muted-foreground"><code>X-Hub-Signature-256</code>: <code>sha256=…</code>, HMAC-SHA256 of the raw body</td>
                        <td className="px-4 py-2 font-mono text-[13px]">GitHubExtractor</td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2 font-mono text-[13px]">HMACVerifier(header)</td>
                        <td className="px-4 py-2 text-muted-foreground">Any named header: hex HMAC-SHA256 of the raw body</td>
                        <td className="px-4 py-2 font-mono text-[13px]">JSONFieldExtractor(t, id)</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed mb-4">
                  The secret comes from <code>os.Getenv(Provider.SecretEnv)</code> at
                  request time &mdash; set <code>STRIPE_WEBHOOK_SECRET</code>,{' '}
                  <code>GITHUB_WEBHOOK_SECRET</code>, etc. in your environment. An
                  empty secret is a hard verification error, so a misconfigured
                  deploy fails closed rather than accepting unsigned traffic.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/webhooks/verifiers.go (StripeVerifier, excerpt)"
                  code={`// Stripe-Signature is "t=<unix>,v1=<hex>"; v1 = HMAC-SHA256 of
// "<timestamp>.<payload>". 5-minute tolerance guards against replay.
signed := strconv.FormatInt(ts, 10) + "." + string(body)
mac := hmac.New(sha256.New, []byte(secret))
mac.Write([]byte(signed))
expected := hex.EncodeToString(mac.Sum(nil))
for _, s := range sigs {
    if hmac.Equal([]byte(s), []byte(expected)) {
        return nil // valid
    }
}
return fmt.Errorf("webhooks: stripe signature mismatch")`}
                />
              </div>

              {/* Dedupe */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Deduplication
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Providers retry &mdash; Stripe redelivers until it gets a 2xx,
                  GitHub pings on every config save. The <code>ExternalID</code>{' '}
                  (the provider&apos;s own event id) is the idempotency key. A
                  partial unique index on <code>(provider, external_id)</code> makes
                  a duplicate INSERT fail, which the handler catches and turns into
                  an immediate <code>status=skipped</code> 200 &mdash; the handler
                  never runs twice.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/models/webhook_event.go"
                  code={`type WebhookEvent struct {
    ID           string         \`gorm:"primarykey;size:36" json:"id"\`
    Provider     string         \`gorm:"size:50;index;not null" json:"provider"\`
    EventType    string         \`gorm:"size:100;index" json:"event_type"\`
    ExternalID   string         \`gorm:"size:255;index" json:"external_id"\` // provider's event id
    Payload      datatypes.JSON \`gorm:"type:jsonb" json:"payload"\`
    Status       string         \`gorm:"size:20;index;not null;default:pending" json:"status"\`
    HandlerError string         \`gorm:"type:text" json:"handler_error,omitempty"\`
    RetryCount   int            \`gorm:"not null;default:0" json:"retry_count"\`
    ProcessedAt  *time.Time     \`json:"processed_at,omitempty"\`
    CreatedAt    time.Time      \`gorm:"index" json:"created_at"\`
}

// Composite unique index → idempotent receipt.
// CREATE UNIQUE INDEX ... ON webhook_events(provider, external_id) WHERE external_id <> ''`}
                />

                <p className="text-muted-foreground leading-relaxed mt-4">
                  The status field tracks a delivery through its life:{' '}
                  <code>pending</code> (verified, handler not yet run) &rarr;{' '}
                  <code>processed</code> (handler returned nil) or{' '}
                  <code>failed</code> (handler errored &mdash;{' '}
                  <code>HandlerError</code> holds the message), with{' '}
                  <code>skipped</code> for a deduped retry.
                </p>
              </div>

              {/* Admin + replay */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Admin list &amp; replay
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Every delivery is inspectable under <code>/api/admin/webhooks</code>{' '}
                  (admin role). When a handler fails &mdash; a transient outage, or a
                  bug you&apos;ve since deployed a fix for &mdash; replay re-runs it
                  against the stored payload without asking the provider to redeliver:
                </p>

                <div className="overflow-x-auto mb-4">
                  <table className="w-full text-sm border border-border rounded-lg">
                    <thead>
                      <tr className="border-b border-border bg-muted/30 text-left">
                        <th className="px-4 py-2 font-medium">Method &amp; path</th>
                        <th className="px-4 py-2 font-medium">Does</th>
                      </tr>
                    </thead>
                    <tbody className="font-mono text-[13px]">
                      <tr className="border-b border-border/50"><td className="px-4 py-2">GET /api/admin/webhooks</td><td className="px-4 py-2 font-sans text-muted-foreground">List deliveries; filter by <code>?provider=</code> &amp; <code>?status=</code></td></tr>
                      <tr><td className="px-4 py-2">POST /api/admin/webhooks/:id/replay</td><td className="px-4 py-2 font-sans text-muted-foreground">Re-dispatch the stored event; increments <code>retry_count</code></td></tr>
                    </tbody>
                  </table>
                </div>

                <p className="text-muted-foreground leading-relaxed">
                  Replay bumps <code>retry_count</code> atomically with{' '}
                  <code>gorm.Expr(&quot;retry_count + ?&quot;, 1)</code>, so two
                  concurrent replays of the same event each add one instead of
                  clobbering each other, then records the new{' '}
                  <code>processed</code>/<code>failed</code> outcome.
                </p>
              </div>

              {/* Custom verifier */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding a custom verifier
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  A verifier is just a <code>VerifyFunc</code> &mdash;{' '}
                  <code>func(secret string, body []byte, headers map[string]string) error</code>{' '}
                  &mdash; returning <code>nil</code> when the signature is valid. For a
                  provider whose scheme isn&apos;t covered, write one and register it.
                  For a simple hex-HMAC provider you don&apos;t even need to: reach for{' '}
                  <code>HMACVerifier(&quot;X-Signature&quot;)</code>.
                </p>

                <CodeBlock
                  language="go"
                  filename="internal/webhooks/setup.go"
                  code={`// A partner that base64-encodes the HMAC in "X-Partner-Signature".
func partnerVerify(secret string, body []byte, headers map[string]string) error {
    got := headers["X-Partner-Signature"]
    if got == "" {
        return fmt.Errorf("missing X-Partner-Signature")
    }
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write(body)
    expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
    if !hmac.Equal([]byte(got), []byte(expected)) {
        return fmt.Errorf("signature mismatch")
    }
    return nil
}

func init() {
    webhooks.Register("partner", webhooks.Provider{
        SecretEnv: "PARTNER_WEBHOOK_SECRET",
        Verify:    partnerVerify,
        Extract:   webhooks.JSONFieldExtractor("type", "id"), // reads body.type + body.id
    })
    webhooks.On("partner", "order.created", handlePartnerOrder)
}`}
                />

                <div className="rounded-lg border border-primary/20 bg-primary/5 p-4 mt-4">
                  <p className="text-sm text-muted-foreground leading-relaxed">
                    <strong className="text-foreground">Extractors decide idempotency.</strong>{' '}
                    <code>ExtractFunc</code> returns <code>(eventType, externalID)</code>.
                    The event type drives handler dispatch; the external id drives
                    dedupe. Header-based providers (GitHub reads{' '}
                    <code>X-GitHub-Event</code> + <code>X-GitHub-Delivery</code>) use a
                    custom extractor; JSON-envelope providers use{' '}
                    <code>JSONFieldExtractor</code>.
                  </p>
                </div>
              </div>

              {/* Go deeper callout */}
              <div className="mb-12">
                <div className="rounded-lg border border-primary/20 bg-primary/5 p-5">
                  <h3 className="text-base font-semibold text-foreground mb-2">
                    Go deeper
                  </h3>
                  <p className="text-sm text-muted-foreground leading-relaxed mb-3">
                    Build a production Stripe webhook receiver from scratch &mdash;
                    signature verification, idempotent fulfilment, failure replay, and
                    testing with the Stripe CLI.
                  </p>
                  <Link
                    href="/courses/webhook-receiver"
                    className="text-sm text-primary hover:underline font-medium"
                  >
                    Course: Building a Webhook Receiver &rarr;
                  </Link>
                </div>
              </div>

              {/* Prev / Next */}
              <div className="flex items-center justify-between border-t border-border pt-8 mt-12">
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/feature-flags" className="gap-2">
                    <ArrowLeft className="h-4 w-4" />
                    Feature Flags
                  </Link>
                </Button>
                <Button variant="ghost" asChild>
                  <Link href="/docs/backend/realtime" className="gap-2">
                    Realtime (WebSockets)
                    <ArrowRight className="h-4 w-4" />
                  </Link>
                </Button>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
