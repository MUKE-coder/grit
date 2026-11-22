package scaffold

import (
	"path/filepath"
	"strings"
	"testing"
)

// The plan package is delivered, and it is delivered by the set that travels
// on upgrade.
//
// A project scaffolded before this release gets `grit migrate plan` only
// because plan.go is in writeFrameworkOwnedFiles. Put it in the new-project
// map alone and the command exists in the CLI, fails in every existing
// project, and says the package is missing (grit#58).
func TestThePlanPackageTravelsOnUpgrade(t *testing.T) {
	opts := Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext}
	opts.Normalize()

	src := readScaffoldSource(t, "api_framework_owned_files.go")
	for _, want := range []string{
		`"migrate", "plan.go"):      apiMigratePlanGo()`,
		`"migrate", "plan_test.go"): apiMigratePlanTestGo()`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("writeFrameworkOwnedFiles does not deliver %s", want)
		}
	}
}

// The command the CLI runs and the flag the tool reads have to be the same
// word. They are in different files in different packages and nothing else
// pairs them: a rename on one side gives "flag provided but not defined",
// from a tool the person did not know they were running.
func TestTheMigrateToolTakesTheFlagTheCLISends(t *testing.T) {
	tool := apiMigrateMainGo()
	for _, want := range []string{
		`flag.Bool("plan"`,
		`flag.Bool("strict"`,
		"migrate.BuildPlan(db, models.Models())",
		"migrate.Format(db, p)",
	} {
		if !strings.Contains(tool, want) {
			t.Errorf("cmd/migrate does not contain %q", want)
		}
	}

	cli := readScaffoldSource(t, filepath.Join("..", "..", "cmd", "grit", "main.go"))
	for _, want := range []string{
		`"cmd/migrate/main.go", "--plan"`,
		`argv = append(argv, "--strict")`,
	} {
		if !strings.Contains(cli, want) {
			t.Errorf("grit migrate plan does not send %q", want)
		}
	}
}

// The generated package compiles as Go, which a raw-string template is not
// checked for anywhere else.
//
// Cheap and worth it: the body carries struct tags spliced through a
// sentinel, and a mistake there produces a project that does not build, found
// by whoever scaffolds next rather than here.
func TestThePlanTemplateIsPlausibleGo(t *testing.T) {
	for name, body := range map[string]string{
		"plan.go":      apiMigratePlanGo(),
		"plan_test.go": apiMigratePlanTestGo(),
	} {
		if !strings.HasPrefix(body, "package migrate\n") {
			t.Errorf("%s does not open with its package clause", name)
		}
		if strings.Contains(body, "@@") {
			t.Errorf("%s still contains the backtick sentinel, so the struct "+
				"tags in it are the literal characters @@", name)
		}
		if open, close := strings.Count(body, "{"), strings.Count(body, "}"); open != close {
			t.Errorf("%s has %d { and %d }", name, open, close)
		}
	}

	// The tags survived the sentinel, which is the thing that can silently go
	// wrong about this particular template.
	if !strings.Contains(apiMigratePlanTestGo(), "`gorm:\"primarykey\"`") {
		t.Error("plan_test.go has no backticked struct tag, so the sentinel " +
			"replacement did not happen")
	}
}
