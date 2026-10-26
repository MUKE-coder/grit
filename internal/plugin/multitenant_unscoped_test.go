package plugin

import (
	"strings"
	"testing"
)

// Unscoped returned GORM's chainable instance, which shares one Statement
// across every call made on it. A handle kept in a variable therefore carried
// the previous call's WHERE clause and its error into the next (grit#103):
//
//	tx = tenant.Unscoped(tx)
//	tx.Where("slug = ?", slug).Take(&existing)  // ErrRecordNotFound
//	return tx.Create(&org).Error                // also "record not found"
//
// The one-shot form in the doc comment was always fine, and keeping the handle
// is the natural next step, so this found people who had done nothing wrong
// and the failure pointed nowhere near scoping.
//
// The behaviour is covered by a test the plugin ships into every project. This
// is the cheap guard on the line itself, because the failure mode is somebody
// tidying away a Session call that looks redundant.
func TestUnscopedReturnsASession(t *testing.T) {
	src := mtTenantPackage(Context{Module: "example/apps/api"})

	start := strings.Index(src, "func Unscoped(")
	if start < 0 {
		t.Fatal("the tenant package has no Unscoped")
	}
	body := src[start:]
	if end := strings.Index(body, "\n}"); end > 0 {
		body = body[:end]
	}

	if !strings.Contains(body, "Session(&gorm.Session{})") {
		t.Error("Unscoped must return a session: without it a kept handle " +
			"shares one Statement and leaks conditions and errors between calls")
	}
}

// And the test that proves the behaviour has to actually reach the project,
// or the guard above is the only thing standing between this and a repeat.
func TestTheLeakTestShipsWithThePlugin(t *testing.T) {
	tests := mtTenantTests(Context{Module: "example/apps/api"})
	for _, want := range []string{
		"func TestUnscoped_HandleDoesNotLeakBetweenStatements",
		"the create inherited the previous call's error",
		`"errors"`,
	} {
		if !strings.Contains(tests, want) {
			t.Errorf("the emitted tenant tests are missing %q", want)
		}
	}
}
