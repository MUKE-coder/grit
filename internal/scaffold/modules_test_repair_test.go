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
	if out, fixed, warn := repairModulesTestSource(src); out != src || len(fixed) > 0 || len(warn) != 1 {
		t.Errorf("fixed %v, warned %v", fixed, warn)
	}
}

func TestFreshTemplatesNeedNoModulesTestRepair(t *testing.T) {
	src := configModulesTestGo()
	if out, fixed, warn := repairModulesTestSource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("the modules test template still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
}
