package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func driftProject(t *testing.T, generator string) string {
	t.Helper()
	root := t.TempDir()

	form := filepath.Join(root, "apps", "expo", "components", "resource-forms")
	if err := os.MkdirAll(form, 0o755); err != nil {
		t.Fatal(err)
	}
	// A form that knows about title and done, but not priority.
	if err := os.WriteFile(filepath.Join(form, "tasks-form.tsx"),
		[]byte(`<Input value={values.title} /><Toggle value={values.done} />`), 0o644); err != nil {
		t.Fatal(err)
	}

	dir := filepath.Join(root, ".grit")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{
		"format": 1,
		"files": map[string]map[string]string{
			"apps/api/internal/models/task.go": {"generator": generator, "grit": "3.200.0", "hash": "x"},
		},
	}
	b, _ := json.Marshal(m)
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

func taskStruct() GoStruct {
	return GoStruct{Name: "Task", Fields: []GoField{
		{Name: "ID", JSONName: "id"},
		{Name: "Title", JSONName: "title"},
		{Name: "Done", JSONName: "done"},
		{Name: "Priority", JSONName: "priority"},
		{Name: "Version", JSONName: "version"},
		{Name: "ArchivedAt", JSONName: "archived_at"},
		{Name: "CreatedAt", JSONName: "created_at"},
	}}
}

// A field added to a Go model never reaches the mobile or desktop screens.
//
// grit sync updates the shared type, the Zod schema and the admin resource
// definition, which SyncAdminResource patches through markers. The mobile
// forms and screens and the desktop columns are only written when the resource
// is generated, and neither sync nor `grit generate field` touches them.
//
// Verified on an emulator: a column added mid-project was in the API response,
// in the shared type, and absent from the mobile create form, with nothing
// anywhere saying so.
func TestSyncReportsFieldsMissingFromClientScreens(t *testing.T) {
	root := driftProject(t, "resource:Task")

	drift := findClientUIDrift(root, taskStruct())
	if len(drift) == 0 {
		t.Fatal("no drift reported for a form that does not mention priority")
	}
	var missing []string
	for _, d := range drift {
		missing = append(missing, d.Missing...)
	}
	joined := strings.Join(missing, ",")

	if !strings.Contains(joined, "priority") {
		t.Errorf("priority was not reported as missing, got %q", joined)
	}
	// Framework-owned columns are never rendered and reporting them was most
	// of the first version's output. A warning that is mostly wrong is one
	// people learn to scroll past.
	for _, noise := range []string{"version", "archived_at", "id", "created_at"} {
		if strings.Contains(joined, noise) {
			t.Errorf("%q was reported as missing; it is framework-owned and no "+
				"screen is expected to show it", noise)
		}
	}
}

// A built-in resource must not be reported.
//
// Its screens are hand-written and were never derived from the model, so the
// fields are not missing, and suggesting a regeneration would be advice that
// destroys them.
func TestSyncIgnoresResourcesItDidNotGenerate(t *testing.T) {
	root := driftProject(t, "scaffold")

	if drift := findClientUIDrift(root, taskStruct()); len(drift) != 0 {
		t.Errorf("drift reported for a resource the generator did not write: %+v", drift)
	}
}

// A project older than the manifest gets silence rather than a guess.
func TestSyncSaysNothingWithoutAManifest(t *testing.T) {
	root := t.TempDir()
	form := filepath.Join(root, "apps", "expo", "components", "resource-forms")
	_ = os.MkdirAll(form, 0o755)
	_ = os.WriteFile(filepath.Join(form, "tasks-form.tsx"), []byte("<Input />"), 0o644)

	if drift := findClientUIDrift(root, taskStruct()); len(drift) != 0 {
		t.Errorf("drift reported without a manifest to attribute the resource: %+v", drift)
	}
}
