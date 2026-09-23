package plugin

import (
	"embed"
	"strings"
)

func init() { Register(stripePlugin()) }

//go:embed stripe/*.tmpl
var stripeTemplates embed.FS

// stripePlugin takes payments through Stripe.
//
// Design:
//
//   - No stripe-go. Payments needs four REST calls, and every project already
//     verifies Stripe's webhook signatures in internal/webhooks. stripe-go's
//     own webhook parser also refuses events whose API version differs from
//     the library's, which breaks a working endpoint on the first upgrade.
//
//   - The server decides the amount. There is no route that starts a
//     payment; the app's checkout handler works out the total and calls
//     payments.Service.Create, and the browser only receives a client secret.
//
//   - Stripe's word marks a payment paid, never the browser's: the webhook,
//     or Refresh reading the intent back with the secret key. Every change is
//     a conditional update, so redelivered and out-of-order events are
//     harmless, and OnSucceeded runs once, in the transaction that marks the
//     payment paid.
//
// Found missing building the storefront blueprint, whose blog post did all of
// this by hand.
func stripePlugin() Plugin {
	return Plugin{
		Name:    "stripe",
		Version: "1.1.0",
		Summary: "Stripe payments and subscriptions: server-priced checkout, hosted billing, webhook-settled, refunds",
		Description: `Takes money through Stripe without trusting the browser with the price.

  • payments.Service.Create(Checkout{UserID, Reference, Amount, Currency}) from
    your own checkout handler starts a PaymentIntent and answers the client
    secret; reloading the checkout reuses it, a changed basket cancels it
  • POST /webhooks/stripe (verified with STRIPE_WEBHOOK_SECRET, deduplicated
    by event id) settles it; OnSucceeded marks your order paid, once, in the
    same transaction
  • POST /payments/:id/refresh checks with Stripe, so the return page works
    where no webhook reaches the API, such as on a laptop
  • GET /admin/payments and POST /admin/payments/:id/refund for the admin
  • <StripeCheckout> and <PaymentResult> for the web app, in the app's
    colours; card details go from the browser to Stripe, never to the API

  • Subscriptions: subscriptionService.Start sends the customer to Stripe's
    hosted checkout, the billing portal handles cards, invoices and cancelling,
    and OnActive/OnEnded grant and take away what it buys
  • A failed renewal keeps access while Stripe retries the card, and ends it
    when Stripe gives up: past_due is entitled, unpaid is not

No SDK: plain REST calls, with the API version pinned.`,

		NextSteps: []string{
			"On a project made before v3.311.0, run grit upgrade first: it adds the payment error codes and the CSP list the plugin writes Stripe's origins into",
			"Run the migration:   grit migrate",
			"Install packages:    pnpm install",
			"API .env:            STRIPE_SECRET_KEY=sk_test_...   STRIPE_WEBHOOK_SECRET=whsec_...",
			"Web .env:            NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...",
			"Start a payment in your checkout handler: paymentService.Create(ctx, payments.Checkout{UserID: userID, Reference: \"order:\" + order.ID, Amount: order.TotalCents, Currency: \"usd\"})",
			"Mark the order paid in routes.go: paymentService.OnSucceeded = func(ctx context.Context, tx *gorm.DB, p models.Payment) error { ... }",
			"Webhooks locally:    stripe listen --forward-to localhost:8080/webhooks/stripe (or rely on the return page's refresh)",
			"Subscriptions: make a recurring price in the Stripe dashboard, then POST /subscriptions/checkout {\"price_id\": \"price_...\"} and send the customer to the url it answers",
			"Grant what a subscription buys in routes.go: subscriptionService.OnActive = func(ctx context.Context, tx *gorm.DB, s models.Subscription) error { ... }, and OnEnded to take it away",
			"Where Stripe returns the customer: SITE_URL=https://example.com in .env (it falls back to OAUTH_FRONTEND_URL, then localhost:3000)",
		},

		NodeDeps: []Dependency{
			{Name: "@stripe/stripe-js", Version: "^9.17.0", Workspace: "apps/web"},
			{Name: "@stripe/react-stripe-js", Version: "^6.11.0", Workspace: "apps/web"},
		},

		Files:      stripeFiles,
		Injections: stripeInjections,
	}
}

func stripeTemplate(ctx Context, name string) string {
	raw, err := stripeTemplates.ReadFile("stripe/" + name)
	if err != nil {
		// The templates are compiled in; a missing one is a build mistake.
		panic("stripe plugin template missing: " + name)
	}
	return strings.ReplaceAll(string(raw), "{{MODULE}}", ctx.Module)
}

// hasNextWeb is a Next.js web app, the one frontend the payment form is
// written for: a Vite app's production CSP is in nginx.conf instead.
func hasNextWeb(ctx Context) bool {
	return hasWebApp(ctx) && ctx.Frontend != "tanstack"
}

func stripeFiles(ctx Context) map[string]string {
	files := map[string]string{
		apiPath(ctx, "internal/models/payment.go"):         stripeTemplate(ctx, "payment_model.go.tmpl"),
		apiPath(ctx, "internal/payments/stripe.go"):        stripeTemplate(ctx, "stripe_client.go.tmpl"),
		apiPath(ctx, "internal/payments/service.go"):       stripeTemplate(ctx, "payments_service.go.tmpl"),
		apiPath(ctx, "internal/payments/payments_test.go"): stripeTemplate(ctx, "payments_test.go.tmpl"),
		apiPath(ctx, "internal/handlers/payment.go"):       stripeTemplate(ctx, "payment_handler.go.tmpl"),

		apiPath(ctx, "internal/models/subscription.go"):         stripeTemplate(ctx, "subscription_model.go.tmpl"),
		apiPath(ctx, "internal/payments/billing.go"):            stripeTemplate(ctx, "stripe_billing.go.tmpl"),
		apiPath(ctx, "internal/payments/subscriptions.go"):      stripeTemplate(ctx, "subscriptions_service.go.tmpl"),
		apiPath(ctx, "internal/payments/subscriptions_test.go"): stripeTemplate(ctx, "subscriptions_test.go.tmpl"),
		apiPath(ctx, "internal/handlers/subscription.go"):       stripeTemplate(ctx, "subscription_handler.go.tmpl"),
	}
	if hasNextWeb(ctx) {
		files["apps/web/lib/payments.ts"] = stripeTemplate(ctx, "web_payments.ts.tmpl")
		files["apps/web/hooks/use-payment.ts"] = stripeTemplate(ctx, "web_use_payment.ts.tmpl")
		files["apps/web/components/stripe-checkout.tsx"] = stripeTemplate(ctx, "web_stripe_checkout.tsx.tmpl")
		files["apps/web/lib/subscription.ts"] = stripeTemplate(ctx, "web_subscription.ts.tmpl")
		files["apps/web/hooks/use-subscription.ts"] = stripeTemplate(ctx, "web_use_subscription.ts.tmpl")
	}
	return files
}

func stripeInjections(ctx Context) []Injection {
	injections := []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&Payment{},\n\t\t&Subscription{},",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:imports",
			Code:   "\t\"" + ctx.Module + "/internal/payments\"",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:handlers",
			Code: "\t// Payments through Stripe. Your checkout handler calls\n" +
				"\t// paymentService.Create with the amount it worked out, and\n" +
				"\t// paymentService.OnSucceeded is where the order is marked paid.\n" +
				"\tpaymentService := payments.NewService(db)\n" +
				"\tpaymentService.Register()\n" +
				"\tpaymentHandler := handlers.NewPaymentHandler(db, paymentService)\n" +
				"\t// Subscriptions: a hosted Stripe checkout starts one, Stripe's billing\n" +
				"\t// portal manages it, and the webhooks keep our row in step. Set\n" +
				"\t// OnActive and OnEnded to grant and take away what it buys.\n" +
				"\tsubscriptionService := payments.NewSubscriptions(db, paymentService)\n" +
				"\tsubscriptionService.Register()\n" +
				"\tsubscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService, payments.SiteURL())",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\t// Payments: the caller's own, as recorded or as Stripe has it now.\n" +
				"\t\tprotected.GET(\"/payments/:id\", paymentHandler.Get)\n" +
				"\t\tprotected.POST(\"/payments/:id/refresh\", paymentHandler.Refresh)\n" +
				"\t\t// Subscriptions: the caller's own, and the hosted pages that start\n" +
				"\t\t// and manage them.\n" +
				"\t\tprotected.GET(\"/subscriptions/me\", subscriptionHandler.Current)\n" +
				"\t\tprotected.POST(\"/subscriptions/checkout\", subscriptionHandler.Start)\n" +
				"\t\tprotected.POST(\"/subscriptions/portal\", subscriptionHandler.Portal)\n" +
				"\t\tprotected.POST(\"/subscriptions/refresh\", subscriptionHandler.Refresh)\n" +
				"\t\tprotected.POST(\"/subscriptions/cancel\", subscriptionHandler.Cancel)\n" +
				"\t\tprotected.POST(\"/subscriptions/resume\", subscriptionHandler.Resume)",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:admin",
			Code: "\t\tadmin.GET(\"/admin/payments\", paymentHandler.List)\n" +
				"\t\tadmin.POST(\"/admin/payments/:id/refund\", paymentHandler.Refund)",
		},
	}
	if hasNextWeb(ctx) {
		// Stripe.js is a script, its form is frames, and it calls its API
		// from the browser: all three are refused without these. A config
		// from before the marker is reported with the lines to add.
		injections = append(injections, Injection{
			File:   "apps/web/next.config.ts",
			Marker: "// grit:csp-origins",
			Code: "  [\"script-src\", \"https://js.stripe.com\"],\n" +
				"  [\"frame-src\", \"https://js.stripe.com\"],\n" +
				"  [\"frame-src\", \"https://hooks.stripe.com\"],\n" +
				"  [\"connect-src\", \"https://api.stripe.com\"],\n" +
				"  [\"img-src\", \"https://*.stripe.com\"],",
			Optional: true,
		})
	}
	return injections
}
