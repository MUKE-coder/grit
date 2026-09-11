package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The package documented flags.IsEnabled(c, ...) and flags.Variant(c, ...),
// and neither existed: IsEnabled was a method on an engine that lived in a
// local variable of routes.Setup, so application code could not check a flag.
// And a flag could target user IDs only, so "on for the EU business unit"
// meant listing every EU user by hand.
func TestFlagsEngineHasAttributesAndAPackageAPI(t *testing.T) {
	src := apiFlagsGo()
	durableMustFormat(t, "flags.go", src)
	durableMustContain(t, "flags.go", src,
		"func IsEnabled(c *gin.Context, name string) bool",
		"func Variant(c *gin.Context, name string) string",
		"func IsEnabledFor(s Subject, name string) bool",
		"type Subject struct",
		"var AttributesFor = func(c *gin.Context) map[string]string",
		"setDefault(e)",
		"for key, allowed := range rules.Attributes {",
		"func attributeMatches(value string, allowed []string) bool",
	)

	model := apiFeatureFlagModelGo()
	durableMustFormat(t, "feature_flag.go", model)
	durableMustContain(t, "feature_flag.go", model, "Attributes map[string][]string")

	durableMustFormat(t, "flags_test.go", apiFlagsTestGo())
}

// Both files were in the API map, which grit upgrade never rewrites, so no
// fix to flags could reach a project that already existed.
func TestFlagsTravelOnUpgrade(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app", Architecture: ArchAPI}
	if err := writeFrameworkOwnedFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"internal/models/feature_flag.go", "internal/flags/flags.go", "internal/flags/flags_test.go"} {
		if _, err := os.Stat(filepath.Join(opts.APIRoot(root), filepath.FromSlash(rel))); err != nil {
			t.Errorf("%s is not written by the upgrade path: %v", rel, err)
		}
	}
	got, _ := os.ReadFile(filepath.Join(opts.APIRoot(root), "internal", "flags", "flags.go"))
	if strings.Contains(string(got), "{{MODULE}}") {
		t.Error("the module placeholder was not filled in")
	}
}
