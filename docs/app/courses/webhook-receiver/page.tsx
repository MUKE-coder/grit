import Link from 'next/link'
import { Cable, ShieldCheck, Repeat, Layers } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Webhook Receiver — Stripe, GitHub & HMAC Verifiers in Grit',
  description:
    'Receive third-party webhooks safely in Grit: verify Stripe / GitHub / generic HMAC signatures, dedup on (provider, external_id), and replay events for safe retries.',
}

const learn = [
  { icon: ShieldCheck, title: 'Verify signatures', body: 'Reject forged calls — Stripe, GitHub, and generic HMAC verifiers ship in the box.' },
  { icon: Layers, title: 'Dedup events', body: 'A unique (provider, external_id) index makes duplicate deliveries a no-op.' },
  { icon: Repeat, title: 'Replay safely', body: 'Re-process a stored event without asking the provider to resend it.' },
  { icon: Cable, title: 'Handle the payload', body: 'Map a verified event to your domain logic — fulfil orders, sync repos, notify users.' },
]

export default function WebhookReceiverCourse() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Webhook Receiver</span>
        </div>

        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">4 challenges</span>
          </div>
          <h1 className="font-display text-3xl md:text-4xl font-bold text-foreground mb-4">
            Webhook Receiver
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Grit ships a hardened webhook receiver: <strong className="text-foreground">signature
            verification</strong> for Stripe, GitHub, and generic HMAC providers, automatic
            <strong className="text-foreground"> deduplication</strong> on (provider, external_id),
            and a <strong className="text-foreground">replay</strong> path so retries are always safe.
          </p>
        </div>

        <div className="grid sm:grid-cols-2 rounded-xl border border-foreground/15 overflow-hidden mb-12">
          {learn.map(({ icon: Icon, title, body }) => (
            <div key={title} className="border-b border-r border-foreground/15 p-5">
              <Icon className="h-5 w-5 text-primary mb-3" />
              <h3 className="font-semibold text-foreground text-sm mb-1">{title}</h3>
              <p className="text-sm text-muted-foreground leading-relaxed">{body}</p>
            </div>
          ))}
        </div>

        {/* 1 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">What is a webhook?</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            A webhook is an HTTP request a third party sends <em>to you</em> when something
            happens on their side — a Stripe payment succeeds, a GitHub branch is pushed, a
            Twilio message is delivered. Instead of you polling their API, they push the event.
          </p>
          <Definition term="Why webhooks need hardening">
            The endpoint is public, so anyone can POST to it. Providers also retry on timeouts,
            so the same event can arrive several times, and the network can deliver them out of
            order. A safe receiver must <strong>authenticate</strong>, <strong>deduplicate</strong>,
            and be <strong>idempotent</strong>.
          </Definition>
        </section>

        {/* 2 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Verifying the signature</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Providers sign each request with a shared secret. Grit recomputes the signature over
            the <em>raw</em> request body and compares it in constant time. A mismatch is
            rejected with <Code>401</Code> before any handler runs.
          </p>
          <Definition term="HMAC">
            Hash-based Message Authentication Code — <Code>HMAC(secret, body)</Code>. Only someone
            who knows the secret can produce a valid code for a given body, so a matching HMAC
            proves the request is authentic and unmodified.
          </Definition>
          <CodeBlock filename="conceptual — verification step">
{`// Grit verifies before dispatch; you just register the secret.
// Stripe:  Stripe-Signature header, secret = whsec_...
// GitHub:  X-Hub-Signature-256 header, secret = your webhook secret
// Generic: X-Signature header = hex(hmac_sha256(secret, rawBody))
if !verify(provider, rawBody, signatureHeader, secret) {
    c.JSON(401, gin.H{"error": "invalid signature"})
    return
}`}
          </CodeBlock>
          <Note>
            Signatures are computed over the <strong>raw bytes</strong>. If middleware re-encodes
            the JSON before verification, the signature won&apos;t match — Grit captures the raw
            body for exactly this reason.
          </Note>
        </section>

        {/* 3 */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Dedup &amp; replay</h2>
          <p className="text-muted-foreground leading-relaxed mb-4">
            Every verified event is stored with the provider&apos;s own event id. A unique index
            on <Code>(provider, external_id)</Code> means a duplicate delivery is silently
            ignored — the handler runs exactly once.
          </p>
          <Definition term="Idempotency">
            Processing the same event N times has the same effect as processing it once. Combined
            with the unique index, this is what makes provider retries harmless.
          </Definition>
          <CodeBlock filename="the dedup guarantee">
{`// First delivery  -> insert row, run handler
// Retry / duplicate -> insert hits the unique (provider, external_id)
//                      index, is skipped, handler does NOT run again
//
// Replay (manual)  -> re-run the handler against the STORED payload,
//                      no call back to the provider required`}
          </CodeBlock>
          <Tip>
            Store the raw payload, not just the parsed fields. Replay reads from that stored copy,
            so you can re-drive a fix through old events without asking Stripe or GitHub to resend.
          </Tip>
        </section>

        {/* challenges */}
        <section className="mb-12">
          <h2 className="font-display text-2xl font-bold text-foreground mb-4">Practice</h2>
          <Challenge number={1} title="Send a forged request">
            <p>POST a fake JSON body to your webhook endpoint with no (or a wrong) signature
            header. What status code do you get? Confirm no handler logic ran.</p>
          </Challenge>
          <Challenge number={2} title="Use the Stripe CLI">
            <p>Run <Code>stripe listen --forward-to localhost:8080/api/webhooks/stripe</Code> and
            trigger a test event. Does it verify and process? Find the stored event row in GORM Studio.</p>
          </Challenge>
          <Challenge number={3} title="Force a duplicate">
            <p>Replay the same Stripe event id twice. Confirm the handler ran only once. Which
            column enforces that?</p>
          </Challenge>
          <Challenge number={4} title="Replay a stored event">
            <p>Trigger a replay of a previously received event. Did your handler re-run against the
            stored payload? Why is replaying from storage safer than asking the provider to resend?</p>
          </Challenge>
        </section>

        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses/stripe-payments', label: 'Stripe Payments' }}
          />
        </div>
      </main>
    </div>
  )
}
