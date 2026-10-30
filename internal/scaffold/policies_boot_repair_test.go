package scaffold

import (
	"strings"
	"testing"
)

// policies.Register() had no call site anywhere in a generated project, so
// every rule in internal/policies was registered into a map nothing filled
// and no rule ran on any path. A Patch on a row a rule denied answered 200.
//
// Found while testing the enterprise review's 2.4, which reported the rules
// as evadable through a bulk action. They were evadable through everything.

const mainWithoutPolicies = `package main

import (
	"log"

	"example/apps/api/internal/config"
	"example/apps/api/internal/routes"
)

func main() {
	cfg := config.Load()
	db := connect()

	// Setup router
	router := routes.Setup(db, cfg, svc)
	log.Fatal(router.Run())
}
`

func TestTheUpgradeWiresThePoliciesIntoMain(t *testing.T) {
	out, changes, warnings := repairPoliciesBootSource(mainWithoutPolicies, "example/apps/api")
	if len(warnings) > 0 {
		t.Errorf("warnings on a main.go the repair should handle: %v", warnings)
	}
	if len(changes) == 0 {
		t.Fatal("nothing changed, so every rule in internal/policies stays dead")
	}

	if !strings.Contains(out, "policies.Register()") {
		t.Error("main still does not register the policies")
	}
	if !strings.Contains(out, `"example/apps/api/internal/policies"`) {
		t.Error("the policies import was not added, so main.go does not compile")
	}
	// Before the router, because a rule consulted after the server is serving
	// is a rule that did not apply to the first requests.
	if strings.Index(out, "policies.Register()") > strings.Index(out, "routes.Setup") {
		t.Error("Register is called after the router is built")
	}

	// Upgrading twice must not register twice: Register is not idempotent in
	// any interesting way, but two calls is a smell and a growing main.go.
	again, moreChanges, _ := repairPoliciesBootSource(out, "example/apps/api")
	if again != out || len(moreChanges) > 0 {
		t.Error("a second upgrade added the call again")
	}
	if n := strings.Count(out, "policies.Register()"); n != 1 {
		t.Errorf("Register is called %d times", n)
	}
}

// A main.go somebody restructured is reported, not half-edited. A silent
// skip here means the project keeps believing its rules run.
func TestThePoliciesRepairWarnsWhenItCannotFindTheSpot(t *testing.T) {
	const rewritten = "package main\n\nfunc main() { serve() }\n"

	out, changes, warnings := repairPoliciesBootSource(rewritten, "example/apps/api")
	if out != rewritten || len(changes) > 0 {
		t.Error("a main.go the repair does not understand was edited anyway")
	}
	if len(warnings) == 0 {
		t.Fatal("no warning, so the project is left believing its policies run")
	}
	if !strings.Contains(warnings[0], "stays dead") {
		t.Errorf("the warning does not say what the consequence is: %v", warnings)
	}
}
