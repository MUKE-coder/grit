package generate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeManifest fakes a project whose manifest records framework-owned files.
func writeManifest(t *testing.T, root string, files map[string]string) {
	t.Helper()
	type entry struct {
		Generator string `json:"generator"`
		Grit      string `json:"grit"`
		Hash      string `json:"hash"`
	}
	m := map[string]interface{}{
		"format": 1,
		"grit":   "3.196.0",
		"files":  map[string]entry{},
	}
	fm := m["files"].(map[string]entry)
	for path, gen := range files {
		fm[path] = entry{Generator: gen, Grit: "3.196.0", Hash: "sha256:x"}
	}
	dir := filepath.Join(root, ".grit")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(m)
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), b, 0o644); err != nil {
		t.Fatal(err)
	}
}

// A resource name that owns a framework FILE must be refused, even when the
// name is not a reserved MODEL.
//
// These are different sets and that is the whole bug. The audit log's models
// are ActivityLog and UserActivity, both reserved; the helpers that write to
// them live in services/activity.go, which "Activity" claims. So
// `grit generate resource Activity` reported success and took out LogCreate,
// LogUpdate, LogDelete and DiffSummary, which every generated handler calls.
// The build then failed in access_review.go and event_subscribers.go, two
// files the developer had never opened. Found by building a CRM, where
// "Activity" is what an interaction log is called.
func TestGeneratingOverAFrameworkFileIsRefused(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)
	writeManifest(t, root, map[string]string{
		"apps/api/internal/services/activity.go": "scaffold",
		"apps/api/internal/handlers/activity.go": "scaffold",
	})

	def, err := ParseInlineFields("Activity", "body:text")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)

	err = g.checkFileCollisions(g.Names())
	if err == nil {
		t.Fatal("generating Activity over the framework's activity files was allowed")
	}
	for _, want := range []string{
		"internal/services/activity.go",
		"internal/handlers/activity.go",
		"--force",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal does not mention %q:\n%s", want, err)
		}
	}
	// It has to name a way forward, not just say no.
	if !strings.Contains(err.Error(), "grit generate resource ") {
		t.Error("the refusal suggests no alternative name")
	}
}

// CheckReservedName does not catch this, which is why the file check exists.
// If it ever does, this test should be deleted rather than left passing by
// accident.
func TestActivityIsNotCaughtByTheReservedModelList(t *testing.T) {
	if err := CheckReservedName("Activity"); err != nil {
		t.Skip("Activity is now a reserved model name; the file check is belt and braces here")
	}
}

// An ordinary name is untouched.
func TestGeneratingANameTheFrameworkDoesNotOwnIsAllowed(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)
	writeManifest(t, root, map[string]string{
		"apps/api/internal/services/activity.go": "scaffold",
	})

	def, _ := ParseInlineFields("CustomerActivity", "body:text")
	g := newTestGenerator(root, module, def)
	if err := g.checkFileCollisions(g.Names()); err != nil {
		t.Errorf("a free name was refused: %v", err)
	}
}

// Regenerating a resource over its own files is normal and must stay allowed.
func TestRegeneratingYourOwnResourceIsAllowed(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)
	writeManifest(t, root, map[string]string{
		"apps/api/internal/models/deal.go":   "resource:Deal",
		"apps/api/internal/services/deal.go": "resource:Deal",
	})

	def, _ := ParseInlineFields("Deal", "title:string")
	g := newTestGenerator(root, module, def)
	if err := g.checkFileCollisions(g.Names()); err != nil {
		t.Errorf("regenerating an existing resource was refused: %v", err)
	}
}

// A project older than the manifest has nothing to check against, and must
// still be able to generate.
func TestProjectWithoutAManifestStillGenerates(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, _ := ParseInlineFields("Activity", "body:text")
	g := newTestGenerator(root, module, def)
	if err := g.checkFileCollisions(g.Names()); err != nil {
		t.Errorf("a pre-manifest project was refused: %v", err)
	}
}
