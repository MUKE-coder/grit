package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/scaffold"
)

func richtextNames() Names {
	return Names{Pascal: "Story", Camel: "story", Snake: "story", Plural: "stories", PluralPascal: "Stories"}
}

const wafRoutesFixture = `package routes

func wafExcludedRoutes() []string {
	prefix := "/api/" + APIVersion
	paths := []string{
		"/uploads", "/uploads/*",
		// grit:waf:richtext
	}
	return paths
}
`

func TestHasRichtext(t *testing.T) {
	with := &ResourceDefinition{Fields: []Field{{Name: "title", Type: "string"}, {Name: "body", Type: "richtext"}}}
	if !with.HasRichtext() {
		t.Error("a richtext field was not seen")
	}
	// A textarea is not excluded: every resource has a description, and
	// excluding all of them would leave the WAF inspecting nothing.
	without := &ResourceDefinition{Fields: []Field{{Name: "notes", Type: "text"}}}
	if without.HasRichtext() {
		t.Error("a plain text field counted as richtext")
	}
}

func TestInjectAndRemoveWAFRichtext(t *testing.T) {
	dir := t.TempDir()
	routes := filepath.Join(dir, "routes.go")
	if err := os.WriteFile(routes, []byte(wafRoutesFixture), 0o644); err != nil {
		t.Fatal(err)
	}
	names := richtextNames()

	if !injectWAFRichtext(routes, names) {
		t.Fatal("the exclusion was not injected")
	}
	got := read(t, routes)
	for _, want := range []string{
		`"/stories", "/stories/*",`,
		`"/admin/stories", "/admin/stories/*",`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %s in:\n%s", want, got)
		}
	}
	// The marker has to survive, or the next resource generated is not excluded.
	if !strings.Contains(got, scaffold.WAFRichtextMarker) {
		t.Error("the injection consumed the marker")
	}
	// The entries go before the marker, so the list reads in generation order.
	if strings.Index(got, `"/admin/stories"`) > strings.Index(got, scaffold.WAFRichtextMarker) {
		t.Error("the entries landed after the marker")
	}

	injectWAFRichtext(routes, names)
	if again := read(t, routes); again != got {
		t.Error("generating the resource twice added the exclusion twice")
	}

	if !removeWAFRichtext(routes, names) {
		t.Fatal("grit remove did not take the exclusion out")
	}
	if after := read(t, routes); after != wafRoutesFixture {
		t.Errorf("removal did not restore the file:\n%s", after)
	}
	if removeWAFRichtext(routes, names) {
		t.Error("a second removal reported a change")
	}
}

// An older project has no marker. A generate there must still succeed: the
// alternative is refusing to add a resource over a WAF entry, and grit upgrade
// adds both the marker and the entry.
func TestInjectWAFRichtextSkipsAProjectWithoutTheMarker(t *testing.T) {
	dir := t.TempDir()
	routes := filepath.Join(dir, "routes.go")
	old := strings.Replace(wafRoutesFixture, "\t\t"+scaffold.WAFRichtextMarker+"\n", "", 1)
	if err := os.WriteFile(routes, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	if injectWAFRichtext(routes, richtextNames()) {
		t.Error("an injection was reported with no marker to inject at")
	}
	if got := read(t, routes); got != old {
		t.Error("routes.go was changed anyway")
	}
	if injectWAFRichtext(filepath.Join(dir, "absent.go"), richtextNames()) {
		t.Error("a missing routes.go reported an injection")
	}
}

// grit upgrade writes the block with the line endings the file already had, so
// removal has to match a CRLF copy too.
func TestRemoveWAFRichtextOnCRLF(t *testing.T) {
	dir := t.TempDir()
	routes := filepath.Join(dir, "routes.go")
	crlf := strings.ReplaceAll(wafRoutesFixture, "\n", "\r\n")
	body := strings.ReplaceAll(wafRichtextEntries(richtextNames())+"\n", "\n", "\r\n")
	src := strings.Replace(crlf, "\t\t"+scaffold.WAFRichtextMarker, body+"\t\t"+scaffold.WAFRichtextMarker, 1)
	if err := os.WriteFile(routes, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if !removeWAFRichtext(routes, richtextNames()) {
		t.Fatal("a CRLF exclusion was not removed")
	}
	if got := read(t, routes); got != crlf {
		t.Errorf("removal did not restore the file:\n%q", got)
	}
}
