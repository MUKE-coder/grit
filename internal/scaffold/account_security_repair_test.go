package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

// The templates carry the account security text, and each repair turns the
// text Grit generated before into the same result.

func TestTemplatesCarryAccountSecurity(t *testing.T) {
	checks := map[string][]string{
		"idempotency middleware": {apiIdempotencyMiddlewareGo(), idempotencyKeyNew, idempotencyScopeFunc},
		"OAuth handler":          {apiAuthOAuthGo(), oauthEmailGuard, oauthUnverifiedLink},
		"user handler":           {apiUserHandlerGo(), profileChecks, profileEmailNew, profileRevokeNew, "CurrentPassword string"},
	}
	for name, parts := range checks {
		for _, want := range parts[1:] {
			if !strings.Contains(parts[0], want) {
				t.Errorf("the %s template is missing %.60q", name, want)
			}
		}
	}
	if strings.Contains(apiIdempotencyMiddlewareGo(), idempotencyKeyOld) {
		t.Error("the idempotency template still keys the cache without the caller")
	}
	// The admin's Update handler sets an email the same way, legitimately: only
	// UpdateProfile, a user changing their own address, needs the checks.
	user := apiUserHandlerGo()
	start := strings.Index(user, "func (h *UserHandler) UpdateProfile(c *gin.Context) {")
	if start < 0 {
		t.Fatal("UpdateProfile is missing from the user handler template")
	}
	profile := user[start : start+strings.Index(user[start:], "\n}\n")]
	if strings.Contains(profile, profileEmailOld) {
		t.Error("the profile handler still takes a new email unchecked")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(user, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("the user handler template is not valid Go: %v", err)
	}
}

func TestRepairIdempotencyScope(t *testing.T) {
	fresh := strings.ReplaceAll(apiIdempotencyMiddlewareGo(), "{{MODULE}}", "example.com/app")
	if out, c, w := repairIdempotencySource(fresh); out != fresh || len(c)+len(w) != 0 {
		t.Errorf("a fresh file needs no repair: %v %v", c, w)
	}
	old := strings.Replace(fresh, idempotencyKeyNew, idempotencyKeyOld, 1)
	old = strings.Replace(old, idempotencyScopeFunc, "", 1)
	for _, imp := range []string{"\t\"crypto/sha256\"\n", "\t\"encoding/hex\"\n", "\t\"strings\"\n"} {
		old = strings.Replace(old, imp, "", 1)
	}
	out, changes, warnings := repairIdempotencySource(old)
	if len(warnings) != 0 || len(changes) != 1 || !strings.Contains(out, "idempotencyScope(c)") {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	mustBeValidGo(t, out)
}

func TestRepairOAuthLink(t *testing.T) {
	fresh := strings.ReplaceAll(apiAuthOAuthGo(), "{{MODULE}}", "example.com/app")
	old := strings.Replace(fresh, oauthEmailGuard, "", 1)
	old = strings.Replace(old, oauthUnverifiedLink, oauthLinkAnchor, 1)
	out, changes, warnings := repairOAuthLinkSource(old)
	if len(warnings) != 0 || len(changes) != 1 || out != fresh {
		t.Fatalf("the repaired handler differs from a fresh one (changes %v warnings %v)", changes, warnings)
	}
}

func TestRepairUserByIDRoute(t *testing.T) {
	src := "package routes\n\nfunc r() {\n" + userByIDRouteOld + "\t\tother()\n" + userListRoute + "}\n"
	out, changes, warnings := repairUserByIDRouteSource(src)
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	if strings.Contains(out, userByIDRouteOld) || !strings.Contains(out, userByIDStaffRoute) {
		t.Errorf("the route did not move:\n%s", out)
	}
}

func TestRepairProfileChange(t *testing.T) {
	fresh := strings.ReplaceAll(apiUserHandlerGo(), "{{MODULE}}", "example.com/app")
	if out, c, w := repairProfileChangeSource(fresh); out != fresh || len(c)+len(w) != 0 {
		t.Errorf("a fresh file needs no repair: %v %v", c, w)
	}
	old := strings.Replace(fresh, profileRequestNew, profileRequestAnchor, 1)
	old = strings.Replace(old, profileChecks, "", 1)
	old = strings.Replace(old, profileEmailNew, profileEmailOld, 1)
	old = strings.Replace(old, profileRevokeNew, profileRevokeOld, 1)
	old = strings.Replace(old, "\t\"strings\"\n", "", 1)
	out, changes, warnings := repairProfileChangeSource(old)
	if len(warnings) != 0 || len(changes) != 1 {
		t.Fatalf("changes %v warnings %v", changes, warnings)
	}
	for _, want := range []string{profileChecks, profileEmailNew, profileRevokeNew, "CurrentPassword string"} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired handler is missing %.60q", want)
		}
	}
	mustBeValidGo(t, out)
}

func TestRepairGothicStore(t *testing.T) {
	src := "package main\n\nimport (\n\t\"fmt\"\n\n\t\"github.com/gorilla/sessions\"\n\t\"github.com/markbates/goth/gothic\"\n)\n\nfunc main() {\n" + gothicStoreOld + "\tfmt.Println(gothic.Store)\n}\n"
	src = strings.Replace(src, "func main() {\n", "func main() {\n\tvar cfg struct{ JWTSecret, AppURL string }\n", 1)
	out, changes, warnings := repairGothicStoreSource(src)
	if len(warnings) != 0 || len(changes) != 1 || strings.Contains(out, gothicStoreOld) {
		t.Fatalf("changes %v warnings %v:\n%s", changes, warnings, out)
	}
	mustBeValidGo(t, out)
}

func mustBeValidGo(t *testing.T, src string) {
	t.Helper()
	if _, err := format.Source([]byte(src)); err != nil {
		t.Errorf("not valid Go: %v", err)
	}
}
