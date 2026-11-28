package plugin

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The customer-facing half of a SaaS: where the organization and billing
// screens land, and the two API gaps closed with them.
//
// These read the emitted source rather than running it, as the rest of this
// package's tests do: the handler is a template this binary does not link.

// projectWithCustomerArea is a root that has app/(app), which is what
// grit add web-auth writes and what both plugins look for.
func projectWithCustomerArea(t *testing.T, has bool) Context {
	t.Helper()
	root := t.TempDir()
	if has {
		if err := os.MkdirAll(filepath.Join(root, "apps", "web", "app", "(app)"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return Context{
		Root:         root,
		Module:       "keep/apps/api",
		APIRoot:      "apps/api",
		Architecture: "triple",
		Frontend:     "next",
	}
}

// The organization screen belongs inside the customer area.
//
// app/account/organization is a SIBLING of app/(app)/account, not a child of
// it, so the screen rendered with no sidebar, no header and none of the
// layout's guard: a page that looked like a different application, reached
// from a URL nothing linked to.
func TestTheOrganizationScreenLandsInsideTheCustomerArea(t *testing.T) {
	files := multitenantFiles(projectWithCustomerArea(t, true))

	if _, ok := files["apps/web/app/(app)/organization/page.tsx"]; !ok {
		t.Error("the organization screen is outside the customer shell, so it has no sidebar and nothing links to it")
	}
	if _, ok := files["apps/web/app/account/organization/page.tsx"]; ok {
		t.Error("the organization screen was written twice, at two URLs")
	}
}

// And a project without one keeps the old location, because there is no
// shell to put it in.
func TestWithoutACustomerAreaTheOrganizationScreenStaysWhereItWas(t *testing.T) {
	files := multitenantFiles(projectWithCustomerArea(t, false))

	if _, ok := files["apps/web/app/account/organization/page.tsx"]; !ok {
		t.Error("a project with no customer area lost its organization screen entirely")
	}
	if _, ok := files["apps/web/app/(app)/organization/page.tsx"]; ok {
		t.Error("a page was written into a route group this project does not have")
	}
}

// Both halves of a sidebar entry or neither. An entry naming an icon the
// layout does not import does not compile, and that file is the one every
// signed-in page renders through.
func TestTheOrganizationSidebarEntryBringsItsIcon(t *testing.T) {
	var entry, icon bool
	for _, inj := range multitenantInjections(projectWithCustomerArea(t, true)) {
		if inj.File != "apps/web/app/(app)/layout.tsx" {
			continue
		}
		switch inj.Marker {
		case "// grit:nav":
			entry = true
			if !strings.Contains(inj.Code, "Building2") {
				t.Errorf("the entry uses an icon the import does not name: %q", inj.Code)
			}
		case "// grit:nav-icons":
			icon = true
		}
		if !inj.Optional {
			t.Error("a sidebar injection must be optional: a project with no customer area has no marker, and that is not a failed install")
		}
	}
	if !entry || !icon {
		t.Errorf("organization sidebar: entry=%v icon=%v, want both", entry, icon)
	}
}

// Billing is a screen, not only a hook. The plugin shipped use-subscription
// and left the one page every paying product needs to be written by hand.
func TestStripeWritesTheBillingScreen(t *testing.T) {
	ctx := projectWithCustomerArea(t, true)
	files := stripeFiles(ctx)

	page, ok := files["apps/web/app/(app)/billing/page.tsx"]
	if !ok {
		t.Fatal("no billing screen: the customer has hooks and nowhere to use them")
	}
	for _, want := range []string{"useSubscription", "useBillingPortal", "useCancelSubscription"} {
		if !strings.Contains(page, want) {
			t.Errorf("the billing screen does not use %s", want)
		}
	}
	// Invoices are Stripe's portal's job. A table here would be a copy that
	// drifts, and saying so in the file is what stops somebody adding one.
	if strings.Contains(page, "invoices.map") {
		t.Error("the billing screen rebuilt the invoice list that the portal already serves")
	}

	var entry, icon bool
	for _, inj := range stripeInjections(ctx) {
		if inj.File == "apps/web/app/(app)/layout.tsx" {
			if inj.Marker == "// grit:nav" {
				entry = true
			}
			if inj.Marker == "// grit:nav-icons" {
				icon = true
			}
		}
	}
	if !entry || !icon {
		t.Errorf("billing sidebar: entry=%v icon=%v, want both", entry, icon)
	}
}

func TestWithoutACustomerAreaStripeWritesNoBillingScreen(t *testing.T) {
	if _, ok := stripeFiles(projectWithCustomerArea(t, false))["apps/web/app/(app)/billing/page.tsx"]; ok {
		t.Error("a billing page was written into a route group this project does not have")
	}
}

// The member list used to return every row, which is fine for the five-person
// organization everybody tests with and wrong for the thousand-person one.
func TestTheMemberListTakesAPage(t *testing.T) {
	src := mtHandler(Context{Module: "keep/apps/api"})

	start := strings.Index(src, "func (h *OrganizationHandler) Members")
	if start < 0 {
		t.Fatal("Members is gone; this test checks nothing")
	}
	body := src[start:]
	if end := strings.Index(body, "\ntype addMemberInput"); end > 0 {
		body = body[:end]
	}

	for _, want := range []string{"page_size", "Limit(pageSize)", "Offset(", "Count(&total)", "\"meta\""} {
		if !strings.Contains(body, want) {
			t.Errorf("the member list has no %s, so it still returns the whole table", want)
		}
	}
	// The sort column goes into ORDER BY as text, so it comes from a map and
	// never from the query string.
	if !strings.Contains(src, "memberSortColumns") {
		t.Error("the sort column is not resolved through a fixed map")
	}
	if strings.Contains(body, `Order(c.Query(`) {
		t.Error("a query parameter is being spliced into ORDER BY")
	}
	// The count has to be taken on the same conditions as the page, or the
	// total answers a different question from the rows.
	if !strings.Contains(body, "query.Session(&gorm.Session{}).Count") {
		t.Error("the total is not counted on the filtered query")
	}
}

// Hooks are what let a project do something on create without editing the
// handler, which is the file grit upgrade maintains.
func TestTheOrganizationLifecycleHasHooks(t *testing.T) {
	src := mtHandler(Context{Module: "keep/apps/api"})

	for _, want := range []string{
		"type OrganizationHooks struct",
		"BeforeCreate func(ctx context.Context, tx *gorm.DB, org *models.Organization) error",
		"AfterCreate func(ctx context.Context, tx *gorm.DB, org *models.Organization) error",
		"AfterMemberAdded func(ctx context.Context, tx *gorm.DB, orgID, userID, role string) error",
		"AfterMemberRemoved func(ctx context.Context, tx *gorm.DB, orgID, userID string) error",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("missing %q", want)
		}
	}

	// Inside the transaction, or a hook that fails leaves an organization
	// claiming to have something it does not.
	create := src[strings.Index(src, "func (h *OrganizationHandler) Create"):]
	if end := strings.Index(create, "\n// Get returns"); end > 0 {
		create = create[:end]
	}
	txAt := strings.Index(create, "h.DB.Transaction(")
	beforeAt := strings.Index(create, "h.Hooks.BeforeCreate")
	afterAt := strings.Index(create, "h.Hooks.AfterCreate")
	if txAt < 0 || beforeAt < txAt || afterAt < txAt {
		t.Error("the create hooks run outside the transaction, so a failed hook leaves the organization behind")
	}
	// Every hook is nil until a project sets one.
	if strings.Count(create, "h.Hooks.BeforeCreate != nil") != 1 {
		t.Error("BeforeCreate is called without checking that it is set")
	}
}
