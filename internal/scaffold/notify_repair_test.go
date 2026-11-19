package scaffold

import (
	"strings"
	"testing"
)

// The repair's output is the template, exactly.
//
// Two things that have to agree: the models/user.go a new project gets and the
// one an upgraded project is repaired into. A test over either alone passes
// while the pair is broken, and here the break is quiet in the worst way: the
// package compiles, the import resolves, and the first notify.Once in a
// scheduled sweep fails against a table that was never created, in the one
// place nobody is watching.
func TestTheNotifyRepairProducesTheTemplate(t *testing.T) {
	fresh := apiUserModelGo()
	if !strings.Contains(fresh, "&NotificationClaim{}") {
		t.Fatal("the claims table is not in the template's model list, so a new " +
			"project's migrate never creates it")
	}

	before := strings.Replace(fresh, notifyModelRegistration, "", 1)
	if before == fresh || strings.Contains(before, "&NotificationClaim{}") {
		t.Fatal("could not rebuild the model list from before the registration")
	}

	out, fixed, warn := repairNotifyModelSource(before)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Error("the repaired model list is not the template's. The two have to be " +
			"the same bytes, or a new project and an upgraded one migrate different " +
			"schemas from the same version.")
	}
}

// Running it twice changes nothing.
func TestTheNotifyRepairIsIdempotent(t *testing.T) {
	fresh := apiUserModelGo()
	out, fixed, warn := repairNotifyModelSource(fresh)
	if out != fresh || len(fixed) != 0 || len(warn) != 0 {
		t.Errorf("a second pass changed the file: fixed %v, warned %v", fixed, warn)
	}
}

// A model list with no marker is reported rather than edited by guesswork.
func TestTheNotifyRepairWarnsWithoutTheMarker(t *testing.T) {
	out, fixed, warn := repairNotifyModelSource("package models\n\nfunc Models() []interface{} { return nil }\n")
	if len(fixed) != 0 {
		t.Errorf("claimed to fix a file with no marker: %v", fixed)
	}
	if len(warn) != 1 || !strings.Contains(warn[0], "&NotificationClaim{}") {
		t.Errorf("the warning does not say what to add by hand: %v", warn)
	}
	if strings.Contains(out, "NotificationClaim") {
		t.Error("it edited the file anyway")
	}
}

// notify does not import the tenant package.
//
// It cannot: notify reads internal/models, and in a multitenant project models
// embeds tenant.Owned, so tenant importing notify or notify importing tenant
// is an import cycle either way. That is why the scope is a variable the
// wiring assigns rather than a call, and the tempting simplification is
// exactly the one that does not compile.
func TestNotifyDoesNotImportTenant(t *testing.T) {
	body := apiNotifyGo()
	if strings.Contains(body, "/internal/tenant") {
		t.Error("internal/notify imports the tenant package, which is an import " +
			"cycle in any multitenant project: models embeds tenant.Owned and " +
			"notify reads models")
	}
	if !strings.Contains(body, "var Scope = func(ctx context.Context) string") {
		t.Error("the scope hook is gone, so a multitenant project claims " +
			"notifications globally: the first organization swept claims the key " +
			"and no other organization is ever notified")
	}
}

// The claim is written before the work and released if the work fails.
//
// Both halves matter and they fail in opposite directions. Writing the claim
// after the send leaves a window where two replicas both send; never releasing
// a failed claim means the notification never sends again. The order is the
// whole design of the package, so it is worth a guard that reads it.
func TestTheClaimIsTakenBeforeTheWorkAndReleasedOnFailure(t *testing.T) {
	body := apiNotifyGo()

	claimAt := strings.Index(body, "claimed, err := claim(ctx, db, key)")
	callAt := strings.Index(body, "if err := fn(); err != nil {")
	if claimAt < 0 || callAt < 0 {
		t.Fatal("Once no longer claims then calls, so this guard cannot read it")
	}
	if claimAt > callAt {
		t.Error("the work runs before the claim, so two replicas can both do it")
	}

	// The release is inside the failure branch, not after it.
	failure := body[callAt:]
	end := strings.Index(failure, "return nil")
	if end < 0 {
		t.Fatal("could not find the end of Once")
	}
	if !strings.Contains(failure[:end], "Delete(&models.NotificationClaim{})") {
		t.Error("a failed send does not release its claim, so that notification " +
			"never sends again and nothing says so")
	}
}
