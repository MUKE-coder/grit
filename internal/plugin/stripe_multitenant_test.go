package plugin

import (
	"regexp"
	"strings"
	"testing"
)

// stripeSource returns one of the plugin's generated files, for a project with
// the given architecture.
func stripeSource(t *testing.T, name string) string {
	t.Helper()
	ctx := Context{
		Root:         t.TempDir(),
		Module:       "app/apps/api",
		APIRoot:      "apps/api",
		Architecture: "triple",
		Frontend:     "next",
	}
	body := stripeTemplate(ctx, name)
	if body == "" {
		t.Fatalf("%s rendered empty", name)
	}
	return body
}

// Billing belongs to the organization, not to the person who clicked Checkout
// (grit#144).
//
// The whole design rests on there being ONE implementation rather than a
// multitenant variant. A variant would mean the security-critical half,
// webhook verification and event idempotency and the conditional state
// updates, exists twice and drifts, with the second copy getting less review.
// It would also make install order load-bearing.
//
// What makes one implementation possible is that the organization is read off
// the gin context by string key. The multitenant middleware does
// c.Set("org_id", ...), so c.GetString("org_id") compiles in a project that
// has no tenant package at all and returns "" there.
func TestBillingScopeNeedsNoTenantPackage(t *testing.T) {
	scope := stripeSource(t, "billing_scope.go.tmpl")

	// The import that would break every single-tenant project. Checked over
	// the code rather than the whole file: the comments explain what the
	// multitenant plugin does, and a bare substring match would read its own
	// prose as a dependency.
	code := stripCommentLines(scope)
	for _, forbidden := range []string{"internal/tenant", "tenant.Owned", "tenant.Unscoped"} {
		if strings.Contains(code, forbidden) {
			t.Errorf("scope.go references %q, so it cannot compile without the multitenant plugin", forbidden)
		}
	}
	if !strings.Contains(scope, `c.GetString(orgContextKey)`) {
		t.Error("scope.go does not read the organization off the request")
	}
	// The key the middleware actually sets. If multitenant ever renames it,
	// billing silently stops scoping and nothing fails: this is the pair that
	// has to agree.
	if !strings.Contains(scope, `orgContextKey = "org_id"`) {
		t.Error("the context key is not org_id, which is what the multitenant middleware sets")
	}
}

// Every read path narrows to the scope. A route that forgets is a route that
// returns another organization's money, and the admin list is the easiest one
// to forget because the operator legitimately has the role.
func TestEveryBillingReadIsScoped(t *testing.T) {
	cases := []struct{ file, what string }{
		{"payments_service.go.tmpl", "sc.Apply("},
		{"subscriptions_service.go.tmpl", "sc.Apply("},
		{"payment_handler.go.tmpl", "payments.ScopeFrom(c)"},
		{"subscription_handler.go.tmpl", "payments.ScopeFrom(c)"},
	}
	for _, tc := range cases {
		if body := stripeSource(t, tc.file); !strings.Contains(body, tc.what) {
			t.Errorf("%s never uses %s, so it reads across organizations", tc.file, tc.what)
		}
	}

	// The two that would be read-across-organizations bugs with a role
	// attached: the admin list and the refund.
	handler := stripeSource(t, "payment_handler.go.tmpl")
	if !strings.Contains(handler, "payments.ScopeFrom(c).Apply(h.DB") {
		t.Error("the admin payment list is not scoped: every organization's payments arrive in one page")
	}
	if !strings.Contains(handler, "h.Payments.Refund(c.Request.Context(), payments.ScopeFrom(c)") {
		t.Error("the refund is not scoped, and a refund cannot be taken back")
	}
}

// A user in two organizations needs two Stripe customers, or both share one
// card, one tax id and one invoice history, and the second organization
// cannot hold a plan at all.
func TestTheStripeCustomerIsPerOrganization(t *testing.T) {
	billing := stripeSource(t, "stripe_billing.go.tmpl")
	if !strings.Contains(billing, "sc.CustomerKey()") {
		t.Error("CreateCustomer does not key its idempotency per scope")
	}
	if strings.Contains(billing, `"customer-"+userID`) {
		t.Error("CreateCustomer still keys one customer per user across every organization")
	}

	scope := stripeSource(t, "billing_scope.go.tmpl")
	if !strings.Contains(scope, `"customer-org-" + s.OrgID`) {
		t.Error("the customer key is not per organization")
	}

	// And the organization has to ride on the subscription, because Stripe
	// sends no organization header: a webhook months later has nothing else
	// to go on.
	if !strings.Contains(billing, `subscription_data[metadata][org_id]`) {
		t.Error("the subscription carries no org_id, so a webhook cannot tell which organization it is for")
	}
}

// Both models carry the column, nullable, so a single-tenant project is
// unchanged and a migration adds a column rather than a NOT NULL.
func TestBothBillingModelsCarryANullableOrg(t *testing.T) {
	for _, file := range []string{"payment_model.go.tmpl", "subscription_model.go.tmpl"} {
		body := stripeSource(t, file)
		// gofmt aligns struct fields, so the spacing between the name and
		// the type is whatever the longest field in the block makes it.
		// Matching on "OrgID *string" with one space passes only by luck.
		if !regexp.MustCompile(`OrgID\s+\*string\s`).MatchString(body) {
			t.Errorf("%s has no nullable OrgID", file)
		}
		if !strings.Contains(body, `json:"org_id,omitempty"`) {
			t.Errorf("%s does not expose org_id", file)
		}
		// not null would refuse every row in a single-tenant project.
		if regexp.MustCompile(`OrgID\s+\*string[^\n]*not null`).MatchString(body) {
			t.Errorf("%s declares org_id not null, which no single-tenant project can satisfy", file)
		}
	}
}

// The webhook settles the right organization without a header, by reading the
// metadata once when the row is created and never moving an existing row.
func TestAWebhookResolvesTheOrganizationFromTheSubscription(t *testing.T) {
	subs := stripeSource(t, "subscriptions_service.go.tmpl")
	if !strings.Contains(subs, `incoming.Metadata["org_id"]`) {
		t.Error("apply never reads the organization from the event")
	}
	if !strings.Contains(subs, "if !known {") {
		t.Error("the organization is read on every event, so a later event could move a subscription between organizations")
	}
}

// The plugin has to deliver the new file, or none of the above exists in a
// project. Two delivery lists is the standing trap here.
func TestTheScopeFileIsDelivered(t *testing.T) {
	ctx := Context{
		Root: t.TempDir(), Module: "app/apps/api", APIRoot: "apps/api",
		Architecture: "triple", Frontend: "next",
	}
	files := stripeFiles(ctx)
	for _, want := range []string{
		"apps/api/internal/payments/scope.go",
		"apps/api/internal/payments/scope_test.go",
	} {
		if _, ok := files[want]; !ok {
			t.Errorf("%s is not in the plugin's file list", want)
		}
	}
}

// stripCommentLines drops // comment lines, so a check for a dependency does
// not match a sentence explaining that dependency. The push-chunking repair
// guarded on a bare token and a comment mentioning it was enough to make the
// repair decide its work was already done.
func stripCommentLines(src string) string {
	var out []string
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
