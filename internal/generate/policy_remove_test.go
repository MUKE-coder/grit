package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// generate writes two things for a resource's policy: the file
// internal/policies/<snake>.go, and a register call inside policies.go.
// remove has to take back both, and for a while it took back neither, so
// `grit remove resource Widget` printed success and left a project that did
// not compile: policies/widget.go referenced models.Widget, and policies.go
// still called registerWidgets().
//
// The round-trip test that existed reversed each injection by hand rather
// than calling the remover, so it could not have caught this. These two check
// that the remover's idea of what generate wrote still matches what generate
// writes.
func TestPolicyRemovalMatchesWhatGenerateWrites(t *testing.T) {
	names := MakeNames("Widget")

	// The file. injectPolicyFile writes internal/policies/<snake>.go, which
	// is the path removeResourceFiles has to delete.
	wantFile := filepath.Join("internal", "policies", names.Snake+".go")
	if got := filepath.Join("internal", "policies", "widget.go"); got != wantFile {
		t.Fatalf("the policy file path moved: %s", wantFile)
	}

	// The registration. injectPolicyRegistration adds exactly this line, and
	// the remover matches on the call without its indentation or newline, so
	// one has to contain the other.
	injected := "\tregister" + names.PluralPascal + "()\n"
	removed := "register" + names.PluralPascal + "()"
	if !strings.Contains(injected, removed) {
		t.Errorf("remove looks for %q, which generate's %q does not contain",
			removed, injected)
	}
}

// Both halves have to be in remove.go, because removing one without the other
// swaps one compile error for a different one.
func TestRemoveHandlesBothHalvesOfThePolicy(t *testing.T) {
	src, err := os.ReadFile("remove.go")
	if err != nil {
		t.Fatal(err)
	}
	body := string(src)

	if !strings.Contains(body, `filepath.Join(apiRoot, "internal", "policies", names.Snake+".go")`) {
		t.Error("remove does not delete the resource's policy file, so the " +
			"project stops compiling on models.<Name> the moment the model goes")
	}
	if !strings.Contains(body, `fmt.Sprintf("register%s()", names.PluralPascal)`) {
		t.Error("remove does not take the register call out of policies.go, so " +
			"the project stops compiling on an undefined register<Name>s")
	}
}

// The seeder is written as PluralKebab and was deleted as PluralSnake. Those
// are the same string for a one-word resource and different for every other,
// so `grit remove Widget` was fine and `grit remove BlogPost` left
// blog-posts_seeder.go behind and the project would not build. Every test
// anyone writes uses a one-word name, which is why it lasted.
//
// The general rule: a path built from a name variant has to be deleted using
// the same variant, and a one-word fixture cannot tell two variants apart.
func TestRemoveDeletesTheSeederTheGeneratorActuallyWrote(t *testing.T) {
	names := MakeNames("BlogPost")
	if names.PluralKebab == names.PluralSnake {
		t.Fatal("BlogPost should spell kebab and snake differently; " +
			"this test cannot detect the bug it exists for")
	}

	src, err := os.ReadFile("remove.go")
	if err != nil {
		t.Fatal(err)
	}
	body := string(src)

	// seeder.go writes PluralKebab.
	if !strings.Contains(body, `names.PluralKebab+"_seeder.go"`) {
		t.Error("remove does not delete the kebab-spelled seeder, which is the " +
			"one the generator writes")
	}
}

// The files --tree and a status field add are per-resource too, and each one
// names the model, so leaving any of them behind stops the build the same way.
func TestRemoveDeletesTheTreeAndWorkflowFiles(t *testing.T) {
	src, err := os.ReadFile("remove.go")
	if err != nil {
		t.Fatal(err)
	}
	body := string(src)

	for _, want := range []string{
		`filepath.Join(apiRoot, "internal", "handlers", names.Snake+"_tree.go")`,
		`filepath.Join(apiRoot, "internal", "services", names.Snake+"_tree.go")`,
		`filepath.Join(apiRoot, "internal", "services", names.Snake+"_tree_test.go")`,
		`filepath.Join(apiRoot, "internal", "services", names.Snake+"_workflow.go")`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("remove leaves behind %s", want)
		}
	}
}
