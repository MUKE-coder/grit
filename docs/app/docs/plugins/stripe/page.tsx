import Link from 'next/link'
import { ArrowLeft, ArrowRight, CreditCard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Stripe plugin | Grit',
  description:
    'Take payments through Stripe with the price decided on the server, the order settled by a verified webhook exactly once, refunds from the admin, and a payment form for the web app.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/stripe' },
}

export default function StripePluginPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <CreditCard className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Plugin</span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Stripe</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Take money for an order without trusting the browser with the price. Your server works out what is owed,
            Stripe&apos;s form takes the card, and the order is marked paid when Stripe says so, once, however many
            times Stripe says it.
          </p>

          <CodeBlock language="bash" code={`grit plugin add stripe
grit migrate
pnpm install`} />

          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Then set the keys from the Stripe dashboard: <code>STRIPE_SECRET_KEY</code> and{' '}
            <code>STRIPE_WEBHOOK_SECRET</code> for the API, <code>NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY</code> for the
            web app. Without a secret key the payment routes answer <code>PAYMENTS_UNAVAILABLE</code> rather than
            failing somewhere less obvious. On a project made before v3.311.0, run <code>grit upgrade</code> first.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Start a payment from your own handler</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            There is deliberately no route that starts a payment. Your checkout handler already knows the order, so
            it works out the total and asks for that amount. The browser receives a client secret that can confirm
            this one payment and nothing else.
          </p>
          <CodeBlock language="go" code={`// In your checkout handler, after the order is priced on the server.
started, err := h.Payments.Create(ctx, payments.Checkout{
    UserID:    userID,
    Reference: "order:" + order.ID,
    Amount:    order.TotalCents, // cents, as Stripe counts
    Currency:  "usd",
})
// started.ClientSecret goes to the browser; started.Payment is the row.`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            A customer who reloads the checkout gets the same payment back, not a second one. If the basket changed,
            the payment for the old amount is canceled at Stripe first, so its secret cannot be used to pay the old
            price.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Mark the order paid</h2>
          <CodeBlock language="go" code={`// routes.go, after paymentService is made.
paymentService.OnSucceeded = func(ctx context.Context, tx *gorm.DB, p models.Payment) error {
    orderID, ok := strings.CutPrefix(p.Reference, "order:")
    if !ok {
        return nil
    }
    return tx.Model(&models.Order{}).
        Where("id = ? AND total_cents = ?", orderID, p.Amount).
        Update("status", "paid").Error
}`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            It runs inside the transaction that marks the payment succeeded. If it returns an error, neither
            change is kept, the webhook is recorded as failed, and replaying it from the admin runs both again.
            <code> OnRefunded</code> works the same way for refunds.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The payment form</h2>
          <CodeBlock language="tsx" code={`import { StripeCheckout, PaymentResult } from "@/components/stripe-checkout";

// The checkout page: started is what your endpoint answered.
<StripeCheckout started={started} returnPath="/orders/thanks" />

// The return page: reads ?payment= and checks with Stripe.
<PaymentResult id={searchParams.get("payment") ?? undefined} />`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Stripe&apos;s Payment Element, in the app&apos;s colours, with whatever payment methods are switched on
            in the dashboard. Card details go from the browser to Stripe and never reach the API. The plugin adds
            Stripe.js to the web app&apos;s Content-Security-Policy: its script, its frames and its API, each of
            which the browser otherwise refuses with nothing but a console message.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What you get</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left">
                  <th className="py-2 pr-4 font-medium text-foreground">Endpoint</th>
                  <th className="py-2 font-medium text-foreground">Does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">POST /webhooks/stripe</td>
                  <td className="py-2">
                    Stripe&apos;s events, verified with <code>STRIPE_WEBHOOK_SECRET</code> and stored once per event
                    id. Settles payments, cancellations and refunds.
                  </td>
                </tr>
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">GET /api/v1/payments/:id</td>
                  <td className="py-2">One of your payments, as recorded.</td>
                </tr>
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/payments/:id/refresh</td>
                  <td className="py-2">
                    Asks Stripe how it stands and records it, so the return page works where no webhook reaches the
                    API, such as on a laptop.
                  </td>
                </tr>
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">GET /api/v1/admin/payments</td>
                  <td className="py-2">Every payment, filterable by status, reference and user.</td>
                </tr>
                <tr>
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/admin/payments/:id/refund</td>
                  <td className="py-2">Refunds all of a payment, or the amount given.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The details that matter</h2>
          <ul className="mb-8 list-disc space-y-3 pl-6 leading-relaxed text-muted-foreground">
            <li>
              <strong className="text-foreground">Stripe&apos;s word, never the browser&apos;s.</strong> A redirect
              back to your site proves nothing. A payment is marked paid by a signed webhook, or by the API reading
              the intent back from Stripe with the secret key, and only when the amount and currency Stripe reports
              are the ones the payment was for.
            </li>
            <li>
              <strong className="text-foreground">At least once, in any order.</strong> Stripe retries and does not
              promise order. Every change is a conditional update from the states it may come from, so a second
              delivery changes nothing and a late <code>processing</code> cannot undo a <code>succeeded</code>.
            </li>
            <li>
              <strong className="text-foreground">Retries do not charge twice.</strong> Each intent is created with
              the payment&apos;s own id as its idempotency key, and each refund with a key made of what has been
              refunded so far, so a double click makes one refund.
            </li>
            <li>
              <strong className="text-foreground">No SDK.</strong> Four REST calls with the API version pinned.
              stripe-go&apos;s webhook parser refuses any event whose API version differs from the library&apos;s,
              which breaks a working endpoint on the first upgrade; the signature check every Grit project already
              has does not care.
            </li>
          </ul>

          <p className="mb-8 leading-relaxed text-muted-foreground">
            Built for the storefront blueprint, and checked against a running API on a new project: a wrong signature
            was refused with 401, a redelivered event was skipped as a duplicate, a late <code>processing</code> left a
            paid payment paid, an intent reporting 1 cent against a payment of 1000 was recorded as a failed event
            with that reason and left unpaid, and a partial refund was recorded. The shipped tests run the rest against
            a fake Stripe.
          </p>

          <p className="mb-8 leading-relaxed text-muted-foreground">
            Testing webhooks locally: <code>stripe listen --forward-to localhost:8080/webhooks/stripe</code> prints
            the signing secret to use. Or skip it: the return page&apos;s refresh settles card payments on its own.
          </p>

          <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
            <Link href="/docs/plugins/video">
              <Button variant="ghost" className="gap-2">
                <ArrowLeft className="h-4 w-4" />
                Video
              </Button>
            </Link>
            <Link href="/docs/plugins/authoring">
              <Button variant="ghost" className="gap-2">
                Writing a plugin
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
