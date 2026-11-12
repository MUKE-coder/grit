package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// desktopScreenWithMarkers is the shape the generator writes: a fenced COLUMNS
// array and a fenced SEARCH_KEYS line.
const desktopScreenWithMarkers = `import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/app/invoices/")({
  component: InvoicesPage,
});

const COLUMNS: DataColumn[] = [
  // grit:cols:auto-start
  { key: "reference", label: "Reference", format: "text" },
  { key: "created_at", label: "Created", format: "relative" },
  // grit:cols:auto-end
];

// grit:search-keys:start
const SEARCH_KEYS = ["reference"];
// grit:search-keys:end

function InvoicesPage() {
  return null;
}
`

// writeDesktopScreen lays out a project with a desktop screen and the admin
// resource file that marks the resource as generated.
func writeDesktopScreen(t *testing.T, body string) string {
	t.Helper()
	root := t.TempDir()

	screen := filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app")
	if err := os.MkdirAll(screen, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(screen, "invoices.index.tsx"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	// generatedResource looks for the model the generator wrote.
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join("testdata", "invoice_model.go.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(models, "invoice.go"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	// generatedResource asks the manifest who wrote the resource, and answers
	// false for a project that has none: a pre-manifest project gets the drift
	// report and no edits, because a wrong guess here rewrites hand-written
	// screens. So the fixture has to record what the generator would record.
	recordGeneratedResource(t, root, "Invoice")
	return root
}

// recordGeneratedResource writes the manifest entry `grit generate resource`
// leaves behind.
func recordGeneratedResource(t *testing.T, root, pascal string) {
	t.Helper()
	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	m.Record("apps/api/internal/models/"+toSnakeCase(pascal)+".go", "resource:"+pascal, "3.509.0", "")
	if err := m.Save(root, "3.509.0"); err != nil {
		t.Fatal(err)
	}
}

// invoiceStruct is the parsed model, read from the same fixture the rebuild
// tests use so both read the shape a real generator writes.
func invoiceStruct(t *testing.T, root string) GoStruct {
	t.Helper()
	structs, err := parseGoStructs(filepath.Join(root, "apps", "api", "internal", "models", "invoice.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range structs {
		if s.Name == "Invoice" {
			return s
		}
	}
	t.Fatal("no Invoice struct in the fixture")
	return GoStruct{}
}

// The gap grit#87 reported: a field is in the API, in the shared types and in
// the admin, and absent from the screen somebody is looking at on a laptop.
func TestANewFieldReachesTheDesktopScreen(t *testing.T) {
	root := writeDesktopScreen(t, desktopScreenWithMarkers)
	added, warnings, err := SyncDesktopScreen(root, invoiceStruct(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) > 0 {
		t.Errorf("warnings on a screen with markers: %v", warnings)
	}
	if added == 0 {
		t.Fatal("nothing was added to a screen that is missing most of the model")
	}

	raw, err := os.ReadFile(filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app", "invoices.index.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)

	// Each of these renders differently, so getting the format wrong is
	// visible rather than theoretical.
	for _, want := range []string{
		`{ key: "notes", label: "Notes", format: "text" },`,
		`{ key: "total", label: "Total", format: "money" },`,
		`{ key: "paid", label: "Paid", format: "boolean" },`,
		`{ key: "status", label: "Status", format: "badge" },`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the screen is missing %s", want)
		}
	}

	// The association field is not a column, and neither is anything the
	// framework owns.
	for _, unwanted := range []string{`key: "client"`, `key: "version"`, `key: "deleted_at"`, `key: "archived_at"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("the screen gained %s, which is not a column somebody declared", unwanted)
		}
	}

	// A column you can see and cannot find is half a feature.
	if !strings.Contains(got, `"notes"`) || !strings.Contains(got, "const SEARCH_KEYS = [") {
		t.Error("SEARCH_KEYS did not gain the new text columns")
	}
	// But not the ones where a substring match finds nothing.
	if strings.Contains(got, `SEARCH_KEYS = ["reference", "paid"`) {
		t.Error("a boolean was added to SEARCH_KEYS, where it can only make search look broken")
	}

	// The fences survive, or a second sync has nowhere to write.
	for _, marker := range []string{desktopColsStart, desktopColsEnd, desktopSearchStart, desktopSearchEnd} {
		if !strings.Contains(got, marker) {
			t.Errorf("the %s marker was consumed", marker)
		}
	}
}

// Insert-only, and quiet when there is nothing to do. A sync that rewrites an
// unchanged file shows up as a change on every run forever.
func TestASecondSyncOfTheDesktopScreenChangesNothing(t *testing.T) {
	root := writeDesktopScreen(t, desktopScreenWithMarkers)
	s := invoiceStruct(t, root)
	if _, _, err := SyncDesktopScreen(root, s); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app", "invoices.index.tsx")
	first, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	added, _, err := SyncDesktopScreen(root, s)
	if err != nil {
		t.Fatal(err)
	}
	if added != 0 {
		t.Errorf("a second sync added %d more columns", added)
	}
	second, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(first) != string(second) {
		t.Error("a second sync rewrote the file")
	}
}

// A row somebody customised is left exactly as it is, because the field is
// already mentioned. That is what makes this safe to leave running.
func TestACustomisedColumnIsNotDuplicated(t *testing.T) {
	customised := strings.Replace(desktopScreenWithMarkers,
		`  { key: "created_at", label: "Created", format: "relative" },`,
		`  { key: "created_at", label: "Created", format: "relative" },
  { key: "total", label: "Amount due", format: "money", align: "right" },`, 1)

	root := writeDesktopScreen(t, customised)
	if _, _, err := SyncDesktopScreen(root, invoiceStruct(t, root)); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app", "invoices.index.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)

	if strings.Count(got, `key: "total"`) != 1 {
		t.Errorf("total appears %d times; a customised column was duplicated", strings.Count(got, `key: "total"`))
	}
	if !strings.Contains(got, `label: "Amount due"`) {
		t.Error("the customised label was overwritten")
	}
}

// A screen generated before the markers existed cannot be edited in place, and
// has to say so and name the command that does work rather than guessing at
// where the array ends.
func TestAScreenWithNoMarkersSaysSoAndChangesNothing(t *testing.T) {
	unfenced := strings.NewReplacer(
		"  "+desktopColsStart+"\n", "",
		"  "+desktopColsEnd+"\n", "",
	).Replace(desktopScreenWithMarkers)

	root := writeDesktopScreen(t, unfenced)
	path := filepath.Join(root, "apps", "desktop", "frontend", "src", "routes", "app", "invoices.index.tsx")
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	added, warnings, err := SyncDesktopScreen(root, invoiceStruct(t, root))
	if err != nil {
		t.Fatal(err)
	}
	if added != 0 {
		t.Errorf("added %d columns to a file with no fence to add them inside", added)
	}
	if len(warnings) != 1 {
		t.Fatalf("said %d things about a screen it could not update: %v", len(warnings), warnings)
	}
	if !strings.Contains(warnings[0], "--force") {
		t.Errorf("the warning does not name the command that works: %q", warnings[0])
	}

	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("a file it could not update was written to anyway")
	}
}

// A project with no desktop app is silent, not an error.
func TestNoDesktopScreenIsNotAnError(t *testing.T) {
	root := t.TempDir()
	models := filepath.Join(root, "apps", "api", "internal", "models")
	if err := os.MkdirAll(models, 0o755); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(filepath.Join("testdata", "invoice_model.go.txt"))
	if err := os.WriteFile(filepath.Join(models, "invoice.go"), raw, 0o644); err != nil {
		t.Fatal(err)
	}

	added, warnings, err := SyncDesktopScreen(root, invoiceStruct(t, root))
	if err != nil {
		t.Errorf("a project with no desktop app returned %v", err)
	}
	if added != 0 || len(warnings) != 0 {
		t.Errorf("reported %d added and %v for a project with no desktop app", added, warnings)
	}
}
