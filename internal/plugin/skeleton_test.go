package plugin

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

// writeSkeleton puts the generated plugin on disk, the way the command does.
func writeSkeleton(t *testing.T, name string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), name)
	for rel, body := range SkeletonFiles(name) {
		full := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// grit#143 asked for a template repository or this command and argued for the
// command, because a repository drifts from the manifest schema the first time
// the schema gains a field. This is the other half of that argument: the
// skeleton is loaded with LoadDir, which is what `grit plugin add ./dir`
// calls, so a schema change the skeleton does not follow fails here.
//
// A test that only checked the JSON parses would not do: LoadDir validates
// that every `from` exists, that no path climbs out, and that the plugin does
// something.
func TestTheSkeletonLoads(t *testing.T) {
	dir := writeSkeleton(t, "audit-trail")

	p, err := LoadDir(dir)
	if err != nil {
		t.Fatalf("the skeleton this CLI writes is not a plugin this CLI can install: %v", err)
	}
	if p.Name != "audit-trail" {
		t.Errorf("loaded as %q", p.Name)
	}
	if p.Summary == "" {
		t.Error("no summary, which is what grit plugin list shows")
	}
}

// The list the issue asked the example to exercise: one file, one injection,
// one when clause, one Go test, one admin file. Checked against what LoadDir
// yields for a real project rather than against the manifest text.
func TestTheSkeletonExercisesTheManifest(t *testing.T) {
	dir := writeSkeleton(t, "audit-trail")
	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	triple := Context{
		Root:         filepath.Join("tmp", "testbed"),
		Module:       "testbed/apps/api",
		APIRoot:      "apps/api",
		Architecture: "triple",
		Frontend:     "next",
	}

	files := p.Files(triple)
	if len(files) < 3 {
		t.Errorf("writes %d files, want the Go package, its test and the admin card: %v", len(files), keysOf(files))
	}
	var sawGo, sawTest, sawAdmin bool
	for to := range files {
		switch {
		case strings.HasSuffix(to, "_test.go"):
			sawTest = true
		case strings.HasSuffix(to, ".go"):
			sawGo = true
		case strings.HasSuffix(to, ".tsx"):
			sawAdmin = true
		}
	}
	if !sawGo || !sawTest || !sawAdmin {
		t.Errorf("go=%v test=%v admin=%v, want all three: %v", sawGo, sawTest, sawAdmin, keysOf(files))
	}

	if n := len(p.Injections(triple)); n != 2 {
		t.Errorf("%d injections, want the import and the route: one without the other does not compile", n)
	}

	// Every path lands under the API root or the admin app, and none of them
	// still holds a placeholder.
	for to, body := range files {
		if strings.Contains(to, "{{") {
			t.Errorf("%s still has a placeholder in its path", to)
		}
		if strings.Contains(body, "{{") {
			t.Errorf("%s still has a placeholder in its body", to)
		}
	}
}

// The when clause is the one in the list that is easy to get wrong by leaving
// it out, and an --api project is where leaving it out shows: a .tsx file
// lands in a directory nothing compiles, and nobody notices until a build
// somewhere else.
func TestTheAdminFileIsSkippedWhereThereIsNoAdmin(t *testing.T) {
	dir := writeSkeleton(t, "audit-trail")
	p, err := LoadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	apiOnly := Context{
		Root:         filepath.Join("tmp", "testbed"),
		Module:       "testbed/apps/api",
		APIRoot:      "apps/api",
		Architecture: "api",
	}
	for to := range p.Files(apiOnly) {
		if strings.HasSuffix(to, ".tsx") {
			t.Errorf("wrote %s into an --api project, which has no admin app", to)
		}
	}
	// And the Go half still installs, or the when clause is doing too much.
	if len(p.Files(apiOnly)) == 0 {
		t.Error("an --api project gets no files at all")
	}
}

// Two things that have to agree, with nothing comparing them until now: the
// markers the skeleton injects at, and the markers the scaffold emits.
//
// This is the failure mode the issue names and the one with no error message.
// An injection at a marker that is not in the file is not refused at install
// time; it simply does not happen, and the first sign is a route answering
// 404. A skeleton that taught a marker typo would teach it to everybody.
func TestTheSkeletonsMarkersAreMarkersGritEmits(t *testing.T) {
	routes := scaffold.APIRoutesSource()
	if !strings.Contains(routes, "func ") {
		t.Fatal("the routes template is not what this test thinks it is")
	}
	for _, marker := range SkeletonMarkers {
		if !strings.Contains(routes, marker) {
			t.Errorf("the skeleton injects at %q and routes.go has no such marker", marker)
		}
	}

	// And the manifest uses those and only those.
	dir := writeSkeleton(t, "audit-trail")
	raw, err := os.ReadFile(filepath.Join(dir, "plugin.json"))
	if err != nil {
		t.Fatal(err)
	}
	var m struct {
		Injections []struct{ Marker string } `json:"injections"`
	}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	for _, in := range m.Injections {
		if !containsString(SkeletonMarkers, in.Marker) {
			t.Errorf("plugin.json injects at %q, which is not in SkeletonMarkers, so nothing checks it against routes.go", in.Marker)
		}
	}
	if len(m.Injections) != len(SkeletonMarkers) {
		t.Errorf("%d injections and %d declared markers", len(m.Injections), len(SkeletonMarkers))
	}
}

// The README documents the markers an author may inject at, and a list in
// prose is a list that goes stale. Each one it names has to be real.
func TestTheReadmeOnlyNamesRealMarkers(t *testing.T) {
	readme := SkeletonFiles("audit-trail")["README.md"]
	routes := scaffold.APIRoutesSource()

	var named int
	for _, line := range strings.Split(readme, "\n") {
		idx := strings.Index(line, "// grit:")
		if idx < 0 {
			continue
		}
		marker := strings.Fields(line[idx:])[0]
		named++
		// grit:models and grit:handlers are in routes.go too; anything the
		// README names has to be somewhere a plugin can reach.
		if !strings.Contains(routes, marker) {
			t.Errorf("the README names %q, which is not in routes.go", marker)
		}
	}
	if named < 5 {
		t.Errorf("the README names %d markers; the test is reading the wrong thing", named)
	}
}

// Names that would install and then not come back out, or shadow a built-in.
func TestPluginNameValidation(t *testing.T) {
	for _, ok := range []string{"audit-trail", "hello", "a1", "one-two-three"} {
		if err := ValidatePluginName(ok); err != nil {
			t.Errorf("%q was refused: %v", ok, err)
		}
	}
	for _, bad := range []string{
		"", "  ", "Audit", "audit_trail", "-audit", "audit-", "audit--trail",
		"audit/trail", `audit\trail`, "1audit",
	} {
		if err := ValidatePluginName(bad); err == nil {
			t.Errorf("%q was accepted", bad)
		}
	}
	// A built-in, whichever ones exist, so this does not need updating when
	// one is added.
	if built := All(); len(built) > 0 {
		if err := ValidatePluginName(built[0].Name); err == nil {
			t.Errorf("%q is a built-in and was accepted", built[0].Name)
		}
	}
}

// The Go package name has to be an identifier, and the plugin name is
// hyphenated, so the two cannot be the same string.
func TestTheGoPackageNameIsAnIdentifier(t *testing.T) {
	files := SkeletonFiles("audit-trail")

	service, ok := files[filepath.ToSlash(filepath.Join("files", "api", "audittrail.go"))]
	if !ok {
		t.Fatalf("no Go file at files/api/audittrail.go: %v", keysOf(files))
	}
	if !strings.Contains(service, "package audittrail") {
		t.Error("the Go file does not declare a package named after the plugin without its hyphens")
	}
	if strings.Contains(service, "package audit-trail") {
		t.Error("the package name has a hyphen in it, which is not legal Go")
	}
}

func containsString(list []string, want string) bool {
	for _, s := range list {
		if s == want {
			return true
		}
	}
	return false
}
