package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func TestRepairModulesTest(t *testing.T) {
	fresh := configModulesTestGo()
	old := strings.Replace(fresh, modulesTestEnvLines, "", 1)
	if old == fresh {
		t.Fatal("could not rebuild the modules test from before the repair")
	}
	out, fixed, warn := repairModulesTestSource(old)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if out != fresh {
		t.Errorf("the repaired test differs from a fresh one:\n%s", out)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairModulesTestSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed modules_test.go again")
	}
}

func TestRepairModulesTestWarnsOnAnEditedFile(t *testing.T) {
	src := "package config\n\nfunc TestModules_DefaultOn(t *testing.T) {}\n"
	// Two warnings, because this file is recognisable as the module test and
	// matches neither of the things the repair knows how to fix: the APP_ENV
	// line and the SSO exception. Both are worth saying, and neither is
	// guessed at: the file comes back exactly as it went in.
	out, fixed, warn := repairModulesTestSource(src)
	if out != src {
		t.Error("the repair rewrote a file it does not recognise")
	}
	if len(fixed) > 0 {
		t.Errorf("the repair claimed to fix something: %v", fixed)
	}
	if len(warn) != 2 {
		t.Errorf("want a warning for each check that could not run, got %d: %v", len(warn), warn)
	}
}

func TestFreshTemplatesNeedNoModulesTestRepair(t *testing.T) {
	src := configModulesTestGo()
	if out, fixed, warn := repairModulesTestSource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("the modules test template still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
}

// The SSO exception, for a project that already exists.
//
// writeSecurityFiles writes config/modules_test.go and its only callers are
// the two new-project paths, so the template change reaches a new project and
// nothing else. Without this repair, every existing project's upgrade gives it
// an SSO flag its own test says must not exist, and `go test ./...` fails on
// two assertions in a file the reader never wrote.
func TestRepairModulesTestNamesTheSSOException(t *testing.T) {
	// The file as the previous release wrote it, then the APP_ENV repair.
	before, _, _ := repairModulesTestSource(configModulesTestGo())
	_ = before

	// A project on the old template: strip what this change added.
	old := configModulesTestGo()
	for _, e := range modulesTestSSOEdits {
		if strings.Count(old, e.new) != 1 {
			t.Fatalf("the template does not contain the repair's result, so one of them has drifted:\n%s", e.new)
		}
		old = strings.Replace(old, e.new, e.old, 1)
	}
	if strings.Contains(old, "defaultOff") {
		t.Fatal("could not build the pre-SSO file, so this test asserts nothing")
	}

	got, fixed, warnings := repairModulesTestSource(old)
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if len(fixed) == 0 {
		t.Error("the repair changed the file without saying what it did")
	}
	if got != configModulesTestGo() {
		t.Error("a repaired modules_test.go is not the generated one, so a repaired project " +
			"and a new one disagree about the same test")
	}

	// And again: a second upgrade finds nothing to do.
	twice, _, _ := repairModulesTestSource(got)
	if twice != got {
		t.Error("a second upgrade changed the file again")
	}
}
