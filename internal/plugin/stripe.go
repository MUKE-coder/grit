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
		Version: "1.0.0",
		Summary: "Stripe payments: server-priced checkout, webhook-settled, refunds, and a payment form for the web app",
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

No SDK: four REST calls, with the API version pinned.`,

		NextSteps: []string{
			"On a project made before v3.311.0, run grit upgrade first: it adds the payment error codes and the CSP list the plugin writes Stripe's origins into",
			"Run the migration:   grit migrate",
			"Install packages:    pnpm install",
			"API .env:            STRIPE_SECRET_KEY=sk_test_...   STRIPE_WEBHOOK_SECRET=whsec_...",
			"Web .env:            NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY=pk_test_...",
			"Start a payment in your checkout handler: paymentService.Create(ctx, payments.Checkout{UserID: userID, Reference: \"order:\" + order.ID, Amount: order.TotalCents, Currency: \"usd\"})",
			"Mark the order paid in routes.go: paymentService.OnSucceeded = func(ctx context.Context, tx *gorm.DB, p models.Payment) error { ... }",
			"Webhooks locally:    stripe listen --forward-to localhost:8080/webhooks/stripe (or rely on the return page's refresh)",
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
	}
	if hasNextWeb(ctx) {
		files["apps/web/lib/payments.ts"] = stripeTemplate(ctx, "web_payments.ts.tmpl")
		files["apps/web/hooks/use-payment.ts"] = stripeTemplate(ctx, "web_use_payment.ts.tmpl")
		files["apps/web/components/stripe-checkout.tsx"] = stripeTemplate(ctx, "web_stripe_checkout.tsx.tmpl")
	}
	return files
}

func stripeInjections(ctx Context) []Injection {
	injections := []Injection{
		{
			File:   apiPath(ctx, "internal/models/user.go"),
			Marker: "// grit:models",
			Code:   "\t\t&Payment{},",
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
				"\tpaymentHandler := handlers.NewPaymentHandler(db, paymentService)",
		},
		{
			File:   apiPath(ctx, "internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code: "\t\t// Payments: the caller's own, as recorded or as Stripe has it now.\n" +
				"\t\tprotected.GET(\"/payments/:id\", paymentHandler.Get)\n" +
				"\t\tprotected.POST(\"/payments/:id/refresh\", paymentHandler.Refresh)",
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
