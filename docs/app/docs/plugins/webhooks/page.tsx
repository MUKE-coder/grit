import Link from 'next/link'
import { ArrowLeft, ArrowRight, Webhook } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { PluginSteps, RelatedPlugins } from '@/components/plugin-detail'
import type { PluginStep, RelatedPlugin } from '@/components/plugin-detail'

export const metadata = {
  title: 'Webhooks plugin',
  description:
    'Outbound webhooks that follow the Standard Webhooks spec: per-subscription secrets, replay-proof signatures, exponential backoff with a dead letter, and a delivery log you can resend from.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/webhooks' },
}

const STEPS: PluginStep[] = [
  {
    title: 'Install it',
    run: 'grit plugin add webhooks',
    then: 'Writes the wiring into internal/handlers and the admin page. The runtime itself is the grit-webhooks Go module, so upgrades are go get rather than a regenerate.',
  },
  {
    title: 'Rebuild and restart the API',
    run: 'cd apps/api && go build ./...',
    then: 'The subscriptions and delivery_attempts tables are migrated on boot, so there is no separate migration to run.',
  },
  {
    title: 'Add a subscription',
    then: 'Admin, then System, then Webhooks. A subscription is a URL, the events it wants, and a generated secret that starts whsec_. Show the secret to the consumer once: it is theirs, not yours.',
  },
  {
    title: 'Dispatch an event from your own code',
    run: 'handlers.DispatchWebhook("invoice.paid", gin.H{"id": inv.ID, "total": inv.Total})',
    language: 'go',
    then: 'Returns immediately. Delivery, retries and the log all happen behind it, so a slow consumer never slows your handler down.',
  },
  {
    title: 'Watch it arrive, then break it on purpose',
    then: 'The delivery log shows every attempt with its status. Point a subscription at a URL that 500s and watch the backoff, the dead letter when the budget runs out, and the resend button putting it back in the queue.',
  },
]

const RELATED: RelatedPlugin[] = [
  {
    name: 'Stripe',
    href: '/docs/plugins/stripe',
    why: 'The incoming direction. Stripe sends you webhooks; this plugin sends yours to other people.',
  },
  {
    name: 'Realtime',
    href: '/docs/backend/realtime',
    why: 'For your own frontend, where a websocket beats an HTTP callback you then have to secure.',
  },
]

export default function WebhooksPluginPage() {
  return (
    <DocsShell>
          <div className="mb-3 flex items-center gap-2">
            <Webhook className="h-5 w-5 text-accent" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Plugin</span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Webhooks</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Tell other people&apos;s systems when something happened in yours. Signed to the{' '}
            <Link href="https://www.standardwebhooks.com" target="_blank" rel="noreferrer" className="text-accent hover:underline">
              Standard Webhooks
            </Link>{' '}
            spec, retried with backoff, and logged so you can prove what you sent.
          </p>
          <p className="mb-8 leading-relaxed text-muted-foreground">
            This is the <strong>outgoing</strong> direction. For verifying webhooks other people send{' '}
            <em>you</em>, which the core already handles, see{' '}
            <Link href="/docs/backend/webhooks" className="text-accent hover:underline">
              receiving webhooks
            </Link>
            .
          </p>

          <CodeBlock language="bash" code={`grit plugin add webhooks
cd apps/api && go build ./...`} />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Sending one</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            One call, from anywhere a handler can reach. It does not block on the network: the dispatch is
            queued, and a consumer that takes nine seconds to answer is the queue&apos;s problem rather than
            your request&apos;s.
          </p>
          <CodeBlock
            language="go"
            code={`func (h *InvoiceHandler) MarkPaid(c *gin.Context) {
    inv, err := h.Invoices.MarkPaid(c.Request.Context(), c.Param("id"))
    if err != nil {
        respond.WriteError(c, err, "Could not mark the invoice paid")
        return
    }

    // Everyone subscribed to invoice.paid hears about it.
    handlers.DispatchWebhook("invoice.paid", gin.H{
        "id":    inv.ID,
        "total": inv.Total,
    })

    respond.OK(c, inv)
}`}
          />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Why the signature is shaped like that</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The signed content is <code>{'{id}.{timestamp}.{body}'}</code>, not the body alone. Signing only the
            body produces a signature that stays valid for ever: anybody who captures one request can send it
            again tomorrow, and it verifies. With the id and timestamp inside the signature, a replay is a
            message the consumer has already seen, or one whose timestamp is outside its tolerance.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Following the spec also means your consumers do not have to trust your crypto. They verify with any
            off-the-shelf Standard Webhooks library, in whatever language they use, and the integration
            instructions you owe them are one paragraph.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What happens when it fails</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Consumers go down. Each attempt is recorded with its status, retries back off exponentially with
            jitter so a recovering endpoint is not knocked over by everything at once, and a delivery that
            exhausts its budget lands in the dead letter rather than disappearing. The admin page lists them and
            resends in one click.
          </p>
          <p className="mb-8 leading-relaxed text-muted-foreground">
            That log is the part worth having before you need it. &quot;We never got the callback&quot; is
            unanswerable without it, and unarguable with it.
          </p>

          <PluginSteps steps={STEPS} intro={'Standing up outbound webhooks on a real project, and then breaking one on purpose, which is the only way to find out whether the retry path works.'} />

          <RelatedPlugins items={RELATED} />

          <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
            <Link href="/docs/plugins/video">
              <Button variant="ghost" className="gap-2">
                <ArrowLeft className="h-4 w-4" />
                Video
              </Button>
            </Link>
            <Link href="/docs/plugins/marketplace">
              <Button variant="ghost" className="gap-2">
                Marketplace
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
    </DocsShell>
  )
}
