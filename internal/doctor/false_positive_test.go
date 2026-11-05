package doctor

import (
	"strings"
	"testing"
)

// grit#178: two warnings that were wrong on a real multitenant project.
//
// The rule for a check in this package is that it stays quiet on a project
// that is right, because a linter that cries wolf is a linter people turn off.
// Both of these broke it on the shape most likely to hit them.

// A tenant-owned model is already scoped: the callers who can see a row are
// the ones in its organization, not everybody signed in. A user column on such
// a model is usually somebody the row NAMES, an account manager or an
// assignee, so "regenerate with --owned-by user" is wrong advice.
func TestTenantOwnedResourcesAreNotToldToScopeByUser(t *testing.T) {
	tenantOwned := &model{
		Name:        "Company",
		TenantOwned: true,
		Fields: []field{
			{Go: "OwnerID", Type: "string", JSON: "owner_id"},
			{Go: "Owner", Type: "User", JSON: "owner"},
		},
	}
	tenantOwned.byGo = map[string]field{}
	for _, f := range tenantOwned.Fields {
		tenantOwned.byGo[f.Go] = f
	}

	p := &project{
		resources: []string{"Company"},
		models:    map[string]*model{"Company": tenantOwned},
		hasTenant: true,
	}

	if findings := checkUnownedUserResource(p); len(findings) != 0 {
		t.Errorf("a tenant-owned model was told to scope by user: %q\n\n"+
			"It is already scoped to its organization, and its user column is "+
			"usually the person the row names rather than the person it belongs to.",
			findings[0].Message)
	}

	// And the check still fires on a model that genuinely is unscoped, or it
	// would be a check that passes by doing nothing.
	shared := &model{
		Name:   "Note",
		Fields: tenantOwned.Fields,
		byGo:   tenantOwned.byGo,
	}
	p.models["Note"] = shared
	p.resources = []string{"Note"}
	if findings := checkUnownedUserResource(p); len(findings) != 1 {
		t.Errorf("a shared model referencing a user produced %d findings, want 1: "+
			"the check has been switched off rather than narrowed", len(findings))
	}
}

// A project that registers services.OnSSOLogin has decided its membership
// policy. Warning at somebody who has already written the thing you are asking
// for is how a linter gets switched off.
func TestTheSSOWarningIsQuietWhenTheHookIsRegistered(t *testing.T) {
	const routes = `	sso := r.Group("/auth/sso")`

	// Every project with the plugin DECLARES the hook, so the declaration
	// alone must not silence the warning.
	declaredOnly := `
func OnSSOLogin(fn func(ctx context.Context, user *models.User) error) {
	ssoLoginHooks = append(ssoLoginHooks, fn)
}
`
	p := &project{routes: routes, hasTenant: true, allGo: declaredOnly}
	if findings := checkTenantSSOProvisioning(p); len(findings) != 1 {
		t.Errorf("declaring OnSSOLogin silenced the warning (%d findings). Every project "+
			"with the plugin declares it, so this would silence it for everybody.",
			len(findings))
	}

	// A call site is the project saying it has handled this.
	registered := declaredOnly + "\nfunc init() { services.OnSSOLogin(svc.Join) }\n"
	p.allGo = registered
	if findings := checkTenantSSOProvisioning(p); len(findings) != 0 {
		t.Errorf("a project that registers services.OnSSOLogin was still warned: %q",
			findings[0].Message)
	}
}

// The check's own name has to stay findable, since it is what a project
// silences in CI.
func TestTheFalsePositiveChecksAreStillRegistered(t *testing.T) {
	var names []string
	for _, c := range checks {
		names = append(names, c.name)
	}
	joined := strings.Join(names, " ")
	for _, want := range []string{"resource-could-be-owned", "tenant-sso-no-organization"} {
		if !strings.Contains(joined, want) {
			t.Errorf("the %s check is gone: it was meant to be narrowed, not removed", want)
		}
	}
}
