package plugin

func init() { Register(webhooksPlugin()) }

// webhooksPlugin wires the grit-webhooks package into a project: outbound,
// Standard-Webhooks-signed webhooks with retries, a delivery log, and manual
// replay.
//
// This is the "package + plugin" shape the GoDeps mechanism exists for. The
// signing, retry and delivery logic lives in the versioned grit-webhooks Go
// module (upgraded with `go get -u`); the plugin generates only the thin wiring
// — a wrapper in the handlers package so routes.go needs no new import, and the
// admin UI. Core already VERIFIES incoming webhooks; this SENDS outgoing ones.
//
// Deliberately NOT the built-in realtime module or the incoming-webhook
// verifier — different direction, no overlap.
func webhooksPlugin() Plugin {
	return Plugin{
		Name:    "webhooks",
		Version: "1.0.0",
		Summary: "Outbound webhooks — Standard Webhooks signing, retries, delivery log",
		Description: `Send signed webhooks to your users' endpoints, with retries and a delivery log.

  • Signatures follow the Standard Webhooks spec (standardwebhooks.com), so a
    consumer verifies deliveries with any off-the-shelf library — the content
    signed is "{id}.{timestamp}.{body}", so a captured request can't be replayed
  • Per-subscription secrets (whsec_...), exponential backoff with jitter, and a
    dead-letter after the retry budget is spent
  • An admin page to manage subscriptions, watch the delivery log, and resend a
    failed delivery in one click
  • Dispatch an event from anywhere: handlers.DispatchWebhook("invoice.paid", data)

The runtime lives in the grit-webhooks Go module (go get -u to upgrade); only
the wiring is generated here. Core already verifies INCOMING webhooks — this is
the OUTGOING side.`,

		GoDeps: []Dependency{
			{Name: "github.com/MUKE-coder/grit-plugins/grit-webhooks", Version: "v0.4.1"},
		},

		NextSteps: []string{
			"Rebuild the API:  cd apps/api && go build ./...",
			"Restart the API — subscriptions + delivery_attempts are auto-migrated on boot",
			"Open the admin → System → Webhooks to add a subscription",
			"Fire an event from any handler:  handlers.DispatchWebhook(\"invoice.paid\", gin.H{\"id\": inv.ID})",
		},

		Files:      webhooksFiles,
		Injections: webhooksInjections,
	}
}

func webhooksFiles(ctx Context) map[string]string {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	files := map[string]string{
		p("internal/handlers/webhooks_wire.go"): webhooksWireGo(ctx),
	}

	// Admin page (Next / TanStack share the component via nextToTanStack at
	// scaffold time, but a plugin writes the concrete file, so emit per frontend).
	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		admin := "apps/admin"
		if ctx.Frontend == "tanstack" {
			files[admin+"/src/pages/system/webhooks.tsx"] = webhooksAdminPage()
			files[admin+"/src/routes/_dashboard/system/webhooks.tsx"] = webhooksTanStackRoute()
		} else {
			files[admin+"/app/(dashboard)/system/webhooks/page.tsx"] = webhooksAdminPage()
		}
	}

	return files
}

func webhooksInjections(ctx Context) []Injection {
	api := "apps/api"
	if ctx.Architecture == "single" {
		api = "."
	}
	p := func(rel string) string {
		if api == "." {
			return rel
		}
		return api + "/" + rel
	}

	injections := []Injection{
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:handlers",
			// InitWebhooks migrates the tables and builds the service; the
			// wrapper lives in the handlers package so routes.go needs no new
			// import (it already imports handlers).
			Code: "\twebhookService := handlers.InitWebhooks(db)",
		},
		{
			File:   p("internal/routes/routes.go"),
			Marker: "// grit:routes:protected",
			Code:   "\t\thandlers.RegisterWebhookRoutes(protected, webhookService)",
		},
	}

	// Sidebar nav link (triple/full only).
	if ctx.Architecture == "triple" || ctx.Architecture == "full" {
		admin := "apps/admin"
		injections = append(injections, Injection{
			File:   admin + "/components/chrome/CollapsibleSidebar.tsx",
			Marker: "// grit:nav:system",
			Code:   `  { href: "/system/webhooks", label: "Webhooks", iconKey: "Webhook", adminOnly: true },`,
		})
	}

	return injections
}
