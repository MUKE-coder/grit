package plugin

import (
	"strings"
	"testing"
)

func tenantSource(t *testing.T) string {
	t.Helper()
	return mtTenantPackage(Context{Module: "example.com/app"})
}

// ForEach cannot name models.Organization, and the reason is worth a guard.
//
// models embeds tenant.Owned, so the tenant package importing models is an
// import cycle. The first version of ForEach took a func(ctx, Organization)
// and the project refused to build with "undefined: Organization", which is
// the kind of failure that only shows up when somebody scaffolds a real
// multitenant project: every Grit-side test passed.
//
// So it reads the table by name and hands fn an id. A change back to a row is
// a change that does not compile, and this says why before anybody tries.
func TestForEachDoesNotReferenceTheOrganizationModel(t *testing.T) {
	src := tenantSource(t)

	if !strings.Contains(src, "func ForEach(ctx context.Context, db *gorm.DB, fn func(context.Context, string) error) error") {
		t.Error("ForEach's signature has changed. It must not name a models type: " +
			"models embeds tenant.Owned, so importing it here is an import cycle " +
			"and the project will not build.")
	}
	if strings.Contains(src, "/internal/models") {
		t.Error("the tenant package imports models, which is an import cycle in " +
			"any project whose models embed tenant.Owned, which is all of them")
	}
	if !strings.Contains(src, `Table("organizations")`) {
		t.Error("ForEach no longer reads the organizations table by name, so it " +
			"either names a type it cannot import or reads the wrong table")
	}
}

// A sweep must not process deleted organizations.
//
// A query made by table name does not get GORM's deleted_at scope, so the
// condition has to be written out. Without it a billing sweep invoices every
// organization anybody ever deleted, which is the most expensive version of
// this mistake.
func TestForEachSkipsSoftDeletedOrganizations(t *testing.T) {
	src := tenantSource(t)
	if !strings.Contains(src, `Where("active = ? AND deleted_at IS NULL", true)`) {
		t.Error("ForEach does not exclude inactive and soft-deleted organizations. " +
			"Reading by table name skips GORM's deleted_at scope, so the condition " +
			"must be explicit.")
	}
}

// One organization's failure does not end the sweep.
func TestForEachJoinsFailuresRatherThanReturningTheFirst(t *testing.T) {
	src := tenantSource(t)
	at := strings.Index(src, "func ForEach(")
	if at < 0 {
		t.Fatal("no ForEach")
	}
	body := src[at:]
	if end := strings.Index(body, "\n}\n"); end > 0 {
		body = body[:end]
	}

	if !strings.Contains(body, "errors.Join(failures...)") {
		t.Error("ForEach does not join its failures, so one organization's error " +
			"either stops the sweep or is discarded")
	}
	if !strings.Contains(body, "failures = append(failures,") {
		t.Error("ForEach does not collect failures")
	}
	// The error says which organization. "smtp timeout" with no name sends
	// somebody to read the whole sweep to find out whose.
	if !strings.Contains(body, `fmt.Errorf("organization %s: %w"`) {
		t.Error("a failure does not name the organization it came from")
	}
}

// The plugin wires the notification scope, and removal undoes it.
//
// Per organization, because without it the first organization a sweep reaches
// claims "invoice-overdue:42" and no other organization is ever notified about
// its own invoice 42. That failure is invisible: the log says the notification
// was sent, and it was, to somebody else.
func TestThePluginWiresTheNotificationScope(t *testing.T) {
	injections := multitenantInjections(Context{
		Module: "example.com/app",
		Root:   ".",
	})

	var wiring, imported bool
	for _, inj := range injections {
		if strings.Contains(inj.Code, "notify.Scope = func(ctx context.Context) string") {
			wiring = true
			if inj.Marker != "// grit:handlers" {
				t.Errorf("the scope is assigned at %q, which is not where the "+
					"other startup wiring goes", inj.Marker)
			}
		}
		if strings.Contains(inj.Code, `/internal/notify"`) {
			imported = true
		}
	}
	if !wiring {
		t.Error("the plugin does not set notify.Scope, so a multitenant project " +
			"claims notifications globally")
	}
	if !imported {
		t.Error("the notify import is not injected, so the assignment does not compile")
	}
}
