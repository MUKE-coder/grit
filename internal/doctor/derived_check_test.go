package doctor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// derivedProject writes a project with a manifest, a model, and files
// recorded as generated from it.
func derivedProject(t *testing.T, modelBody string, derived map[string]string) *project {
	t.Helper()
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")

	model := filepath.Join(api, "internal", "models", "widget.go")
	writeAt(t, model, modelBody)

	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	m.RecordFrom("apps/api/internal/models/widget.go", "resource:Widget", "3.0.0", modelBody, nil)

	src := &manifest.Source{
		Path: "apps/api/internal/models/widget.go",
		Hash: manifest.Hash(modelBody),
	}
	for rel, body := range derived {
		writeAt(t, filepath.Join(root, filepath.FromSlash(rel)), body)
		m.RecordFrom(rel, "resource:Widget", "3.0.0", body, src)
	}
	if err := m.Save(root, "3.0.0"); err != nil {
		t.Fatal(err)
	}
	return &project{root: root, apiRoot: api}
}

func writeAt(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// A model nobody has touched since the files came out of it is not a finding.
func TestAModelInStepWithItsFilesIsSilent(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/api/internal/handlers/widget.go": "package handlers\n",
		"apps/api/internal/services/widget.go": "package services\n",
	})
	if findings := checkDerivedFiles(p); len(findings) != 0 {
		t.Errorf("a project in step was reported: %q", findings[0].Message)
	}
}

// Edit the model and the files generated from it are behind. One finding for
// the model, not one per file: a dozen findings about one cause is a wall of
// text that gets scrolled past.
func TestAnEditedModelLeavesItsFilesBehind(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/api/internal/handlers/widget.go": "package handlers\n",
		"apps/api/internal/services/widget.go": "package services\n",
		"packages/shared/types/widget.ts":      "export interface Widget {}\n",
	})
	writeAt(t, filepath.Join(p.root, "apps", "api", "internal", "models", "widget.go"),
		"package models\ntype Widget struct{ Colour string }\n")

	findings := checkDerivedFiles(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %+v", len(findings), findings)
	}
	f := findings[0]
	if f.Resource != "apps/api/internal/models/widget.go" {
		t.Errorf("the finding does not name the model: %q", f.Resource)
	}
	if !strings.Contains(f.Message, "3 files behind it") {
		t.Errorf("the finding does not count the files: %q", f.Message)
	}
	// The fix has to name the commands, including the one for the files
	// grit sync does not touch.
	if !strings.Contains(f.Fix, "grit sync") || !strings.Contains(f.Fix, "grit generate resource Widget") {
		t.Errorf("the fix does not name both commands: %q", f.Fix)
	}
}

// A file the person has edited cannot simply be regenerated, and the finding
// has to say so rather than telling them to run a command that would take
// their work with it.
func TestAnEditedGeneratedFileIsCalledOutSeparately(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/api/internal/handlers/widget.go": "package handlers\n",
		"apps/api/internal/services/widget.go": "package services\n",
	})
	writeAt(t, filepath.Join(p.root, "apps", "api", "internal", "models", "widget.go"),
		"package models\ntype Widget struct{ Colour string }\n")
	// And one of the generated files has been changed since Grit wrote it.
	writeAt(t, filepath.Join(p.root, "apps", "api", "internal", "services", "widget.go"),
		"package services\n\n// mine\n")

	findings := checkDerivedFiles(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1: %+v", len(findings), findings)
	}
	f := findings[0]
	if !strings.Contains(f.Message, "1 file you have edited") {
		t.Errorf("the finding does not separate the edited file: %q", f.Message)
	}
	if !strings.Contains(f.Fix, "overwrites") {
		t.Errorf("the fix tells them to regenerate without warning that it overwrites: %q", f.Fix)
	}
}

// Generated from the model is not the same as invalidated by it. The route
// file, the React Query hooks and the admin page wrappers are generated out
// of the model and contain nothing a new column changes.
func TestTheFilesAColumnCannotInvalidateAreNotReported(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/api/internal/routes/widget_routes.go":             "package routes\n",
		"apps/web/hooks/use-widgets.ts":                         "export const useWidgets = () => {}\n",
		"apps/admin/app/(dashboard)/resources/widgets/page.tsx": "export default function P() {}\n",
	})
	writeAt(t, filepath.Join(p.root, "apps", "api", "internal", "models", "widget.go"),
		"package models\ntype Widget struct{ Colour string }\n")

	if findings := checkDerivedFiles(p); len(findings) != 0 {
		t.Errorf("a file a column cannot invalidate was reported: %q", findings[0].Message)
	}
}

// The handler and the service are named first. The first version sorted
// alphabetically and named three admin page wrappers while the handler was in
// the "and 12 more".
func TestTheMostTellingFilesAreNamedFirst(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/admin/resources/widgets/widgets.ts": "export const r = {}\n",
		"apps/api/internal/factory/widget.go":     "package factory\n",
		"apps/api/internal/handlers/widget.go":    "package handlers\n",
		"apps/api/internal/services/widget.go":    "package services\n",
		"packages/shared/types/widget.ts":         "export interface Widget {}\n",
	})
	writeAt(t, filepath.Join(p.root, "apps", "api", "internal", "models", "widget.go"),
		"package models\ntype Widget struct{ Colour string }\n")

	findings := checkDerivedFiles(p)
	if len(findings) != 1 {
		t.Fatalf("got %d findings, want 1", len(findings))
	}
	named, _, _ := strings.Cut(findings[0].Message, " and ")
	for _, want := range []string{"handlers/widget.go", "services/widget.go"} {
		if !strings.Contains(named, want) {
			t.Errorf("%s is not among the files named first: %q", want, named)
		}
	}
}

// A model that was deleted is grit remove's business, not this check's, and
// two complaints about one event is one too many.
func TestADeletedModelIsNotReported(t *testing.T) {
	p := derivedProject(t, "package models\ntype Widget struct{}\n", map[string]string{
		"apps/api/internal/handlers/widget.go": "package handlers\n",
	})
	if err := os.Remove(filepath.Join(p.root, "apps", "api", "internal", "models", "widget.go")); err != nil {
		t.Fatal(err)
	}
	if findings := checkDerivedFiles(p); len(findings) != 0 {
		t.Errorf("a deleted model was reported: %q", findings[0].Message)
	}
}

// Every project generated before this existed has entries with no source, and
// they have to behave exactly as they did.
func TestEntriesWithNoSourceAreIgnored(t *testing.T) {
	root := t.TempDir()
	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	m.Record("apps/api/internal/handlers/widget.go", "resource:Widget", "3.0.0", "package handlers\n")
	if err := m.Save(root, "3.0.0"); err != nil {
		t.Fatal(err)
	}
	p := &project{root: root, apiRoot: filepath.Join(root, "apps", "api")}
	if findings := checkDerivedFiles(p); len(findings) != 0 {
		t.Errorf("an entry with no source was reported: %q", findings[0].Message)
	}
}

// The manifest a project already has must still load after the field is
// added, and must not gain a "from" key it did not have.
func TestTheManifestFormatStaysBackwardCompatible(t *testing.T) {
	root := t.TempDir()
	old := `{"format":1,"grit":"3.0.0","updatedAt":"now","files":{` +
		`"a.go":{"generator":"scaffold","grit":"3.0.0","hash":"abc"}}}`
	writeAt(t, filepath.Join(root, ".grit", "manifest.json"), old)

	m, err := manifest.Load(root)
	if err != nil {
		t.Fatalf("an older manifest no longer loads: %v", err)
	}
	if entry := m.Files["a.go"]; entry.From != nil {
		t.Error("an entry with no source came back with one")
	}
	if err := m.Save(root, "3.1.0"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, ".grit", "manifest.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"from"`) {
		t.Errorf("saving added a from key to an entry that has no source:\n%s", raw)
	}
	// And it is still valid JSON of the shape the CLI expects.
	var round map[string]any
	if err := json.Unmarshal(raw, &round); err != nil {
		t.Fatalf("the saved manifest is not JSON: %v", err)
	}
}
