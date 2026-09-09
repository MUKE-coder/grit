package plugin

import (
	"strings"
	"testing"
)

// The plugin has to turn its own scoping on.
//
// It installed a correct tenant package and a correct middleware and wired
// neither. tenant.RegisterScoping was never called, so the GORM callbacks that
// add the org filter and stamp OrgID were never installed; middleware.Tenant
// was never mounted, so nothing put the active organization on the request
// context for them to read.
//
// Following the docs exactly — embed tenant.Owned, send X-Organization-ID —
// produced a multi-tenant app with no isolation at all. Verified on a CRM:
// one org read another org's contacts, and org_id came back empty on every
// insert. The plugin's own middleware comment says "the whole isolation model
// rests on this check", and nothing called it.
func TestMultitenantTurnsItsScopingOn(t *testing.T) {
	ctx := Context{Root: ".", Module: "crm/apps/api", Architecture: "triple", Frontend: "next"}
	injections := multitenantInjections(ctx)

	var joined string
	markers := map[string]bool{}
	for _, inj := range injections {
		joined += inj.Code + "\n"
		markers[inj.Marker] = true
	}

	if !strings.Contains(joined, "tenant.RegisterScoping(db)") {
		t.Error("RegisterScoping is never called: tenant.Owned is just a column " +
			"and every tenant-owned model is read and written unscoped")
	}
	if !strings.Contains(joined, "middleware.Tenant(db)") {
		t.Error("the tenant middleware is never mounted: nothing resolves the " +
			"active organization onto the request context")
	}
	if !strings.Contains(joined, `"crm/apps/api/internal/tenant"`) {
		t.Error("routes.go never imports the tenant package, so the wiring " +
			"above does not compile")
	}
	if !markers["// grit:middleware:protected"] {
		t.Error("the middleware is not mounted at the protected-group marker, " +
			"where it can read the authenticated user")
	}
}

// Injection code is written verbatim, so a template placeholder would land in
// the source as those literal characters.
func TestMultitenantInjectionsCarryNoTemplatePlaceholders(t *testing.T) {
	ctx := Context{Root: ".", Module: "crm/apps/api", Architecture: "triple", Frontend: "next"}
	for _, inj := range multitenantInjections(ctx) {
		if strings.Contains(inj.Code, "{{") {
			t.Errorf("injection into %s carries an unsubstituted placeholder:\n%s",
				inj.File, inj.Code)
		}
	}
}
