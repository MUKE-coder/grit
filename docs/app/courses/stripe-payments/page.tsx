import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Stripe Payments & Subscriptions for SaaS',
  description: 'Integrate Stripe payments into your Grit SaaS app — checkout sessions, subscription plans, webhooks, customer portal, and billing management.',
}

export default function StripePaymentsCourse() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Stripe Payments & Subscriptions</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Stripe Payments & Subscriptions for SaaS
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Every SaaS needs payments. Stripe is the industry standard — it handles credit cards,
            subscriptions, invoices, and compliance so you don{"'"}t have to. In this course, you{"'"}ll
            integrate Stripe into your Grit app: create checkout sessions, manage subscription plans,
            handle webhooks, and build a complete billing system with a customer portal.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is Stripe? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is Stripe?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Accepting payments on the internet is surprisingly complex. You need to securely collect
            credit card numbers, comply with PCI-DSS regulations, handle failed charges, manage refunds,
            send invoices, and deal with subscriptions. Stripe does all of this for you. You never see
            or store a credit card number — Stripe handles the sensitive parts.
          </p>

          <Definition term="Stripe">
            A payment processing platform that handles credit card payments, subscriptions, invoices,
            and financial compliance. Instead of building payment infrastructure yourself (which requires
            PCI-DSS certification), you integrate Stripe{"'"}s API. Stripe takes a fee per transaction
            (typically 2.9% + 30 cents) and handles everything else.
          </Definition>

          <Definition term="Payment Gateway">
            The intermediary between your application and the credit card networks (Visa, Mastercard, etc.).
            When a customer enters their card, the payment gateway encrypts it, sends it to the card network,
            gets approval or rejection, and returns the result to your app. Stripe is both a payment gateway
            and a full payment platform.
          </Definition>

          <Definition term="Subscription Billing">
            Automatically charging a customer on a recurring schedule (monthly, yearly). Stripe manages
            the entire lifecycle: initial charge, recurring billing, failed payment retries, plan upgrades
            and downgrades, cancellations, and prorated refunds. You define the plans, Stripe collects
            the money.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The payment flow with Stripe:
          </p>

          <CodeBlock filename="Payment Flow">
{`1. User clicks "Subscribe to Pro Plan" on your site
2. Your API creates a Stripe Checkout Session
3. User is redirected to Stripe's hosted payment page
4. User enters credit card on Stripe's page (you never see the card)
5. Stripe charges the card
6. Stripe redirects user back to your success URL
7. Stripe sends a webhook to your API: "payment succeeded"
8. Your API activates the user's subscription`}
          </CodeBlock>

          <Note>
            You never handle credit card numbers. The user enters their card on Stripe{"'"}s hosted page,
            not on your site. This means you don{"'"}t need PCI-DSS certification (a complex, expensive
            security audit). Stripe is PCI-DSS Level 1 certified — the highest level.
          </Note>

          <Challenge number={1} title="Create a Stripe Account">
            <p>Go to <Code>stripe.com</Code> and create a free account. Once logged in, find the
            Dashboard. Switch to {'"'}Test Mode{'"'} (toggle in the top bar). Find your test API keys
            under Developers → API Keys. You should see a Publishable key (starts with <Code>pk_test_</Code>)
            and a Secret key (starts with <Code>sk_test_</Code>). Never share the secret key.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: The grit-stripe Plugin ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The grit-stripe Plugin</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>grit-stripe</Code> plugin wraps Stripe{"'"}s Go SDK with opinionated helpers designed
            for SaaS applications. It provides checkout session creation, subscription management, customer
            portal access, and webhook handling — all pre-configured for the common SaaS billing patterns.
          </p>

          <CodeBlock filename="Terminal">
{`# Install the plugin
go get github.com/MUKE-coder/grit-plugins/grit-stripe`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            What the plugin provides:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CreateCheckoutSession</strong> — Redirect users to Stripe{"'"}s hosted payment page</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CreatePortalSession</strong> — Let users manage their subscription (change plan, update card, cancel)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">HandleWebhook</strong> — Process Stripe events (payment succeeded, subscription cancelled, etc.)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">GetSubscription</strong> — Check a customer{"'"}s current subscription status</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">CancelSubscription</strong> — Cancel at period end or immediately</li>
          </ul>

          <Challenge number={2} title="Install the Plugin">
            <p>Install the <Code>grit-stripe</Code> plugin in your Grit project by running
            <Code>go get github.com/MUKE-coder/grit-plugins/grit-stripe</Code> from your
            <Code>apps/api/</Code> directory. Verify it appears in your <Code>go.mod</Code> and
            <Code>go.sum</Code> files.</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Configuration ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Configuration</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Stripe uses API keys for authentication. You need three keys in your <Code>.env</Code> file:
            the secret key (for server-side API calls), the publishable key (for client-side Stripe.js),
            and the webhook secret (for verifying webhook signatures).
          </p>

          <CodeBlock filename=".env">
{`# Stripe Configuration (TEST MODE — safe for development)
STRIPE_SECRET_KEY=sk_test_51...your_secret_key
STRIPE_PUBLISHABLE_KEY=pk_test_51...your_publishable_key
STRIPE_WEBHOOK_SECRET=whsec_...your_webhook_secret`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Initialize the plugin in your API:
          </p>

          <CodeBlock filename="internal/config/stripe.go">
{`import gritstripe "github.com/MUKE-coder/grit-plugins/grit-stripe"

func InitStripe() {
    gritstripe.Init(gritstripe.Config{
        SecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
        WebhookSecret:  os.Getenv("STRIPE_WEBHOOK_SECRET"),
    })
}`}
          </CodeBlock>

          <Tip>
            Stripe test keys (starting with <Code>sk_test_</Code> and <Code>pk_test_</Code>) work with
            test credit card numbers. Use <Code>4242 4242 4242 4242</Code> with any future expiry and
            any CVC to simulate a successful payment. Use <Code>4000 0000 0000 0002</Code> to simulate
            a declined card. No real money is charged in test mode.
          </Tip>

          <Challenge number={3} title="Configure Stripe Keys">
            <p>Add your Stripe test keys to your <Code>.env</Code> file. Initialize the grit-stripe
            plugin in your config. Start your API and verify there are no errors related to Stripe
            initialization. The secret key starts with <Code>sk_test_</Code> — if it starts with
            <Code>sk_live_</Code>, you{"'"}re using production keys (switch to test mode in the Stripe Dashboard).</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Creating a Checkout Session ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Creating a Checkout Session</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When a user clicks {'"'}Subscribe{'"'} on your pricing page, your API creates a Checkout Session —
            a temporary URL that takes the user to Stripe{"'"}s hosted payment page. The user enters their
            card on Stripe{"'"}s page, and Stripe redirects them back to your app after payment.
          </p>

          <Definition term="Checkout Session">
            A Stripe-hosted payment page created on-demand for each purchase. You specify what the customer
            is buying (a price/product), where to redirect after success or cancellation, and optionally the
            customer{"'"}s email. Stripe handles the entire payment form, card validation, 3D Secure
            authentication, and error handling. The session expires after 24 hours if unused.
          </Definition>

          <CodeBlock filename="internal/handler/billing.go">
{`import gritstripe "github.com/MUKE-coder/grit-plugins/grit-stripe"

func (h *Handler) CreateCheckout(c *gin.Context) {
    user := middleware.GetUser(c) // From auth middleware

    var input struct {
        PriceID string // json:"price_id"
    }
    c.ShouldBindJSON(&input)

    session := gritstripe.CreateCheckoutSession(gritstripe.CheckoutParams{
        PriceID:       input.PriceID,
        SuccessURL:    "https://myapp.com/billing/success",
        CancelURL:     "https://myapp.com/billing",
        CustomerEmail: user.Email,
        Metadata: map[string]string{
            "user_id": fmt.Sprint(user.ID),
        },
    })

    // Return the Stripe-hosted payment page URL
    c.JSON(200, gin.H{"url": session.URL})
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the frontend, when the user clicks a plan{"'"}s {'"'}Subscribe{'"'} button:
          </p>

          <CodeBlock filename="Frontend Flow">
{`// 1. Call your API to create a checkout session
// POST /api/billing/checkout with { price_id: "price_xxx" }

// 2. API returns { url: "https://checkout.stripe.com/c/pay/..." }

// 3. Redirect the user to Stripe's payment page
// window.location.href = response.url

// 4. User enters card on Stripe's page
// 5. On success: Stripe redirects to your success URL
// 6. On cancel: Stripe redirects to your cancel URL`}
          </CodeBlock>

          <Challenge number={4} title="Create a Checkout Session">
            <p>Add a <Code>POST /api/billing/checkout</Code> endpoint to your API. It should accept
            a <Code>price_id</Code> in the request body and return a Stripe Checkout URL. Test it with
            a tool like curl or Postman. Copy the returned URL and open it in your browser — do you
            see Stripe{"'"}s payment page? Use test card <Code>4242 4242 4242 4242</Code> to complete payment.</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Subscription Plans ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Subscription Plans</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before you can charge customers, you need to create Products and Prices in your Stripe
            Dashboard. A Product represents what you{"'"}re selling (e.g., {'"'}Pro Plan{'"'}). A Price represents
            how much it costs and how often (e.g., {'"'}$20/month{'"'}).
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A typical SaaS has 3 tiers:
          </p>

          <CodeBlock filename="Subscription Tiers">
{`Free Plan
  - Price: $0/month (no Stripe integration needed)
  - Features: 1 project, 100 API calls/day, community support
  - Stripe Price ID: none (default for all users)

Pro Plan
  - Price: $20/month
  - Features: 10 projects, 10,000 API calls/day, email support
  - Stripe Price ID: price_pro_monthly

Team Plan
  - Price: $50/month
  - Features: Unlimited projects, 100,000 API calls/day, priority support, team members
  - Stripe Price ID: price_team_monthly`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Create these in the Stripe Dashboard:
          </p>

          <ol className="space-y-3 text-muted-foreground mb-6">
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">1.</span>
              <span>Go to Products in the Stripe Dashboard</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">2.</span>
              <span>Click {'"'}Add Product{'"'}</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">3.</span>
              <span>Name: {'"'}Pro Plan{'"'}, Pricing: Recurring, Amount: $20.00, Billing period: Monthly</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">4.</span>
              <span>Save — Stripe generates a Price ID (starts with <Code>price_</Code>)</span>
            </li>
            <li className="flex gap-3">
              <span className="text-primary font-mono font-bold shrink-0">5.</span>
              <span>Repeat for Team Plan ($50/month)</span>
            </li>
          </ol>

          <Tip>
            Store Price IDs in your <Code>.env</Code> file or a config table, not hardcoded in your Go
            code. This way you can change plans or prices without redeploying. You can also add yearly
            prices (e.g., $200/year for Pro) and let users choose monthly or yearly billing.
          </Tip>

          <Challenge number={5} title="Create Stripe Products">
            <p>In your Stripe Dashboard (test mode), create 2 products: {'"'}Pro Plan{'"'} at $20/month and
            {'"'}Team Plan{'"'} at $50/month. Copy the Price IDs (they start with <Code>price_</Code>). Add
            them to your <Code>.env</Code> file as <Code>STRIPE_PRO_PRICE_ID</Code> and
            <Code>STRIPE_TEAM_PRICE_ID</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Webhook Handling ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Webhook Handling</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s the problem: the user pays on Stripe{"'"}s page and gets redirected to your success
            URL. But what if they close the browser before the redirect? What if the network drops?
            You can{"'"}t rely on the redirect to activate their subscription. Instead, Stripe sends
            webhooks — HTTP POST requests to your API — for every important event.
          </p>

          <Definition term="Webhook">
            An HTTP POST request that Stripe sends to your API when something happens. Instead of your
            API polling Stripe ({'"'}did the payment go through yet?{'"'}), Stripe tells your API ({'"'}payment
            succeeded!{'"'}). Webhooks are signed with a secret so you can verify they actually came from
            Stripe and not an attacker.
          </Definition>

          <CodeBlock filename="internal/handler/webhook.go">
{`func (h *Handler) HandleStripeWebhook(c *gin.Context) {
    gritstripe.HandleWebhook(c, gritstripe.WebhookHandlers{

        // Payment successful — activate subscription
        "checkout.session.completed": func(event stripe.Event) {
            // Extract user_id from session metadata
            // Update user's plan in database
            // Set subscription_status = "active"
            // Set stripe_customer_id on user record
        },

        // Subscription renewed (monthly payment)
        "invoice.payment_succeeded": func(event stripe.Event) {
            // Update subscription period end date
            // Log the payment
        },

        // Payment failed (card declined, expired, etc.)
        "invoice.payment_failed": func(event stripe.Event) {
            // Send email: "Your payment failed"
            // Grace period: keep subscription active for 3 days
            // After grace period: downgrade to free
        },

        // User cancelled — subscription ends at period end
        "customer.subscription.deleted": func(event stripe.Event) {
            // Downgrade user to free plan
            // Set subscription_status = "cancelled"
        },

        // Subscription updated (plan change)
        "customer.subscription.updated": func(event stripe.Event) {
            // Update user's plan in database
            // If upgraded: activate new features immediately
            // If downgraded: keep current plan until period end
        },
    })
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Register the webhook route (no auth middleware — Stripe needs to access it):
          </p>

          <CodeBlock filename="routes.go">
{`// Webhook endpoint — no auth middleware!
// Stripe authenticates via webhook signature, not JWT
r.POST("/api/webhooks/stripe", h.HandleStripeWebhook)`}
          </CodeBlock>

          <Note>
            The webhook endpoint must NOT have auth middleware. Stripe sends the request, not a logged-in
            user. Authentication is handled by the webhook signature — Stripe signs each request with
            your <Code>STRIPE_WEBHOOK_SECRET</Code>, and the plugin verifies the signature automatically.
          </Note>

          <Challenge number={6} title="List Important Stripe Events">
            <p>For a SaaS application, list 5 Stripe webhook events you should handle and describe what
            your API should do for each one. Think about: successful payments, failed payments,
            cancellations, plan changes, and trial expirations. Which event is the most critical
            and why?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Customer Portal ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Customer Portal</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Users need to manage their subscription: update their credit card, change plans, view
            invoices, or cancel. Building all of this yourself is weeks of work. Stripe{"'"}s Customer
            Portal gives you a hosted page where users can do all of this — with zero UI code from you.
          </p>

          <CodeBlock filename="internal/handler/billing.go">
{`func (h *Handler) CreatePortalSession(c *gin.Context) {
    user := middleware.GetUser(c)

    // Get the user's Stripe customer ID from your database
    stripeCustomerID := user.StripeCustomerID

    portalSession := gritstripe.CreatePortalSession(
        stripeCustomerID,
        "https://myapp.com/billing", // Return URL after portal
    )

    c.JSON(200, gin.H{"url": portalSession.URL})
}

// Route: POST /api/billing/portal`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The Customer Portal lets users:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Update payment method</strong> — change credit card, add backup card</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">View invoices</strong> — download PDF invoices for each billing period</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Change plan</strong> — upgrade or downgrade (with prorated billing)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Cancel subscription</strong> — cancel at end of billing period</li>
          </ul>

          <Tip>
            You need to configure the Customer Portal in the Stripe Dashboard under Settings → Customer
            Portal. Enable the features you want (plan changes, cancellation, invoice history). You can
            also customize the branding to match your app.
          </Tip>

          <Challenge number={7} title="Create a Portal Session">
            <p>Add a <Code>POST /api/billing/portal</Code> endpoint. After completing a test checkout
            (from challenge 4), use the Stripe customer ID to create a portal session. Open the returned
            URL. Can you see the subscription? Can you update the payment method? Can you cancel?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Building a Pricing Page ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Building a Pricing Page</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The pricing page is where users decide which plan to buy. It typically shows 3 plan cards
            side by side, with features listed for each plan and a {'"'}Subscribe{'"'} button that creates
            a checkout session.
          </p>

          <CodeBlock filename="Pricing Page Structure">
{`// PricingPage component
//
// Layout: 3 cards in a row (responsive: stack on mobile)
//
// Free Plan Card:
//   - Title: "Free"
//   - Price: "$0/month"
//   - Features: checkmark list (1 project, 100 API calls, etc.)
//   - Button: "Get Started" → link to /register
//
// Pro Plan Card (highlighted — "Most Popular"):
//   - Title: "Pro"
//   - Price: "$20/month"
//   - Features: checkmark list (10 projects, 10K API calls, etc.)
//   - Button: "Subscribe" → POST /api/billing/checkout
//     with price_id = STRIPE_PRO_PRICE_ID
//
// Team Plan Card:
//   - Title: "Team"
//   - Price: "$50/month"
//   - Features: checkmark list (unlimited, priority support, etc.)
//   - Button: "Subscribe" → POST /api/billing/checkout
//     with price_id = STRIPE_TEAM_PRICE_ID
//
// Each Subscribe button:
//   1. Calls POST /api/billing/checkout
//   2. Gets back { url: "https://checkout.stripe.com/..." }
//   3. Redirects: window.location.href = url`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Design tips for an effective pricing page:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Highlight the recommended plan</strong> with a border, badge ({'"'}Most Popular{'"'}), or different background</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Show savings for yearly billing</strong> — add a toggle (Monthly/Yearly) with a {'"'}Save 20%{'"'} badge on yearly</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Use checkmarks for features</strong> — green checkmarks for included, gray X or dash for excluded</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Keep it simple</strong> — 3 plans maximum, clear differentiation between each tier</li>
          </ul>

          <Challenge number={8} title="Build a Pricing Page">
            <p>Build a pricing page with 3 plan cards: Free ($0), Pro ($20/mo), and Team ($50/mo).
            Each paid plan{"'"}s button should call your checkout endpoint with the correct Price ID
            and redirect to Stripe. Style the Pro plan as the recommended option. Add at least
            4 features per plan with checkmark icons.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Checking Subscription Status ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Checking Subscription Status</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Once users can subscribe, you need to check their subscription status on protected routes.
            Pro features should only be accessible to Pro subscribers. Team features should only be
            accessible to Team subscribers.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Store subscription data on the User model:
          </p>

          <CodeBlock filename="models/user.go">
{`type User struct {
    gorm.Model
    Name               string
    Email              string
    Password           string
    Role               string // "USER", "ADMIN"

    // Stripe fields
    StripeCustomerID   string // Stripe customer ID
    SubscriptionPlan   string // "free", "pro", "team"
    SubscriptionStatus string // "active", "cancelled", "past_due"
    PeriodEnd          time.Time // When current billing period ends
}`}
          </CodeBlock>

          <CodeBlock filename="middleware/subscription.go">
{`func RequireSubscription(plans ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := GetUser(c)

        // Check if user's plan is in the required plans
        allowed := false
        for _, plan := range plans {
            if user.SubscriptionPlan == plan {
                allowed = true
                break
            }
        }

        if !allowed {
            c.JSON(403, gin.H{
                "error": "This feature requires a " + plans[0] + " subscription",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Usage in routes:
// r.GET("/api/analytics", RequireSubscription("pro", "team"), h.GetAnalytics)
// r.GET("/api/team/members", RequireSubscription("team"), h.GetTeamMembers)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            On the frontend, show or hide features based on the user{"'"}s plan:
          </p>

          <CodeBlock filename="Plan Check on Frontend">
{`// After login, the user object includes subscription info
// user.subscription_plan = "pro"
// user.subscription_status = "active"

// Show/hide features based on plan:
// If plan is "free": show upgrade banner
// If plan is "pro": show Pro features, hide Team features
// If plan is "team": show all features

// On the dashboard:
// Free users see: "Upgrade to Pro for advanced analytics"
// Pro users see: the analytics dashboard
// Team users see: analytics + team management`}
          </CodeBlock>

          <Challenge number={9} title="Add Subscription Checking">
            <p>Add the <Code>SubscriptionPlan</Code> and <Code>SubscriptionStatus</Code> fields to your
            User model. Create a <Code>RequireSubscription</Code> middleware. Protect a route (e.g.,
            <Code>/api/analytics</Code>) with <Code>RequireSubscription({'"'}pro{'"'}, {'"'}team{'"'})</Code>. Test:
            (1) access the route as a free user (should get 403), (2) manually set a user{"'"}s plan to
            {'"'}pro{'"'} in the database, (3) access the route again (should get 200).</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary + Final Challenge ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            You{"'"}ve learned everything needed to monetize your Grit SaaS application:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Stripe fundamentals</strong> — payment gateway, PCI compliance, test vs live mode</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">grit-stripe plugin</strong> — checkout sessions, portal sessions, webhook handling</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Configuration</strong> — API keys, webhook secrets, test card numbers</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Checkout Sessions</strong> — Stripe-hosted payment pages, no PCI burden</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Subscription Plans</strong> — Products and Prices in the Stripe Dashboard</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Webhooks</strong> — server-to-server event notifications for payment lifecycle</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Customer Portal</strong> — Stripe-hosted subscription management page</span></li>
            <li className="flex gap-2"><CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" /> <span><strong className="text-foreground">Subscription middleware</strong> — protect routes based on user{"'"}s plan</span></li>
          </ul>

          <Challenge number={10} title="Final Challenge: Pricing Page + Checkout">
            <p>Build the complete checkout flow:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Create a pricing page with Free, Pro ($20/mo), and Team ($50/mo) tiers</li>
              <li>Each paid plan{"'"}s button calls <Code>POST /api/billing/checkout</Code></li>
              <li>User is redirected to Stripe, completes payment with test card</li>
              <li>Success URL shows a {'"'}Welcome to Pro!{'"'} page</li>
            </ol>
          </Challenge>

          <Challenge number={11} title="Final Challenge: Webhooks + Activation">
            <p>Complete the webhook handling:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Handle <Code>checkout.session.completed</Code> — activate the subscription in your database</li>
              <li>Handle <Code>invoice.payment_failed</Code> — send a notification (log it for now)</li>
              <li>Handle <Code>customer.subscription.deleted</Code> — downgrade to free plan</li>
              <li>Test: complete a checkout, verify the user{"'"}s plan changed in the database</li>
            </ol>
          </Challenge>

          <Challenge number={12} title="Final Challenge: Full Billing System">
            <p>Build the complete billing system:</p>
            <ol className="mt-2 space-y-2 list-decimal list-inside">
              <li>Pricing page with 3 plans and checkout flow</li>
              <li>Webhook handling for activation, failure, and cancellation</li>
              <li>Customer portal link on the billing settings page</li>
              <li>Subscription status displayed on the user dashboard</li>
              <li><Code>RequireSubscription</Code> middleware protecting Pro and Team routes</li>
              <li>Upgrade banner for free users on restricted pages</li>
            </ol>
          </Challenge>
        </section>

        {/* ═══ Footer ═══ */}
        <CourseFooter />

        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
