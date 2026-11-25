package doctor

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/MUKE-coder/grit/v3/internal/generate"
	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Realistic generated models, because the shape is computed by the same
// rebuild `grit generate resource <Name>` uses and that needs real columns to
// read. A struct with no fields rebuilds into nothing and has no shape.
const usageRecordModel = "package models\n" +
	"\n" +
	"import \"time\"\n" +
	"\n" +
	"type UsageRecord struct {\n" +
	"\tID        string    `gorm:\"primarykey;size:36\" json:\"id\"`\n" +
	"\tLabel     string    `gorm:\"size:255\" json:\"label\" binding:\"required\"`\n" +
	"\tFxRate    float64   `json:\"fx_rate\"`\n" +
	"\tCreatedAt time.Time `json:\"created_at\"`\n" +
	"}\n"

// The same model with the gorm tag that makes AutoMigrate keep the column's
// precision. This is the edit grit#307 is about: the only place it can go, and
// nothing a generator reads.
const usageRecordModelPinned = "package models\n" +
	"\n" +
	"import \"time\"\n" +
	"\n" +
	"type UsageRecord struct {\n" +
	"\tID        string    `gorm:\"primarykey;size:36\" json:\"id\"`\n" +
	"\tLabel     string    `gorm:\"size:255\" json:\"label\" binding:\"required\"`\n" +
	"\tFxRate    float64   `gorm:\"type:numeric(18,8)\" json:\"fx_rate\"`\n" +
	"\tCreatedAt time.Time `json:\"created_at\"`\n" +
	"}\n"

// And a real change: a new column, which every downstream file is genuinely
// behind.
const usageRecordModelNewColumn = "package models\n" +
	"\n" +
	"import \"time\"\n" +
	"\n" +
	"type UsageRecord struct {\n" +
	"\tID        string    `gorm:\"primarykey;size:36\" json:\"id\"`\n" +
	"\tLabel     string    `gorm:\"size:255\" json:\"label\" binding:\"required\"`\n" +
	"\tFxRate    float64   `json:\"fx_rate\"`\n" +
	"\tMemo      string    `gorm:\"size:255\" json:\"memo\"`\n" +
	"\tCreatedAt time.Time `json:\"created_at\"`\n" +
	"}\n"

const usageRecordRel = "apps/api/internal/models/usage_record.go"

// shapedProject records the derived files with both a byte hash and a shape,
// which is what the generator now does, and then leaves the model on disk as
// onDisk.
func shapedProject(t *testing.T, generated, onDisk string) *project {
	t.Helper()
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")
	modelPath := filepath.Join(root, filepath.FromSlash(usageRecordRel))

	// A project root marker and a go.mod, which the rebuild looks for.
	writeAt(t, filepath.Join(root, "docker-compose.yml"), "version: '3'\n")
	writeAt(t, filepath.Join(api, "go.mod"), "module pigeon/apps/api\n\ngo 1.21\n")

	// Record against the model as generated.
	writeAt(t, modelPath, generated)
	shape := generate.ModelShape(modelPath)
	if shape == "" {
		t.Fatal("the fixture model has no shape, so this test would prove nothing")
	}

	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	m.RecordFrom(usageRecordRel, "resource:UsageRecord", "3.0.0", generated, nil)
	src := &manifest.Source{Path: usageRecordRel, Hash: manifest.Hash(generated), Shape: shape}
	for _, derived := range []string{
		"apps/api/internal/handlers/usage_record.go",
		"apps/api/internal/services/usage_record.go",
		"apps/admin/lib/resources/usage-records.ts",
	} {
		body := "// generated\n"
		writeAt(t, filepath.Join(root, filepath.FromSlash(derived)), body)
		m.RecordFrom(derived, "resource:UsageRecord", "3.0.0", body, src)
	}
	if err := m.Save(root, "3.0.0"); err != nil {
		t.Fatal(err)
	}

	// Then the edit.
	writeAt(t, modelPath, onDisk)
	return &project{root: root, apiRoot: api}
}

// A gorm tag that pins a column's type is not drift.
//
// grit#307: the manifest recorded the model's bytes, so this edit reported six
// files behind the model on every run, for good. The only way to clear it was
// `grit generate resource UsageRecord`, which rewrites the model and drops the
// tag, so the column reverted to double precision on the next migrate and the
// rate stopped round-tripping, with nothing to say why.
func TestAGormTagThatChangesNoGeneratedFileIsNotDrift(t *testing.T) {
	p := shapedProject(t, usageRecordModel, usageRecordModelPinned)
	if findings := checkDerivedFiles(p); len(findings) != 0 {
		t.Errorf("a gorm type pin was reported as drift: %q", findings[0].Message)
	}
}

// The check still fires on a change the downstream files really are behind.
func TestANewColumnIsStillDrift(t *testing.T) {
	p := shapedProject(t, usageRecordModel, usageRecordModelNewColumn)
	findings := checkDerivedFiles(p)
	if len(findings) != 1 {
		t.Fatalf("a new column produced %d findings, want 1", len(findings))
	}
	if !strings.Contains(findings[0].Message, "3 files behind it") {
		t.Errorf("message = %q", findings[0].Message)
	}
}

// An entry written before Shape existed keeps the old behaviour: compare the
// bytes, report the drift. Doing anything else would silence the check for
// every project generated before this release.
func TestAnEntryWithNoShapeStillComparesBytes(t *testing.T) {
	root := t.TempDir()
	api := filepath.Join(root, "apps", "api")
	writeAt(t, filepath.Join(root, "docker-compose.yml"), "version: '3'\n")
	writeAt(t, filepath.Join(api, "go.mod"), "module pigeon/apps/api\n\ngo 1.21\n")
	writeAt(t, filepath.Join(root, filepath.FromSlash(usageRecordRel)), usageRecordModelPinned)

	m := &manifest.Manifest{Format: manifest.FormatVersion, Files: map[string]manifest.Entry{}}
	// No Shape, as every entry written before this release has.
	src := &manifest.Source{Path: usageRecordRel, Hash: manifest.Hash(usageRecordModel)}
	body := "// generated\n"
	writeAt(t, filepath.Join(api, "internal", "handlers", "usage_record.go"), body)
	m.RecordFrom("apps/api/internal/handlers/usage_record.go", "resource:UsageRecord", "3.0.0", body, src)
	if err := m.Save(root, "3.0.0"); err != nil {
		t.Fatal(err)
	}

	findings := checkDerivedFiles(&project{root: root, apiRoot: api})
	if len(findings) != 1 {
		t.Fatalf("an entry with no shape produced %d findings, want 1", len(findings))
	}
	// And the remedy now says what regenerating costs, which it did not when
	// no derived file had been edited: that is the case where the advice reads
	// as safe and reverts the tag.
	if !strings.Contains(findings[0].Fix, "rewrites the model too") {
		t.Errorf("the remedy does not warn that regenerating rewrites the model: %q", findings[0].Fix)
	}
}

// The shape is stable: the same model hashes the same way twice, or every
// doctor run would report drift against its own last recording.
func TestTheShapeOfOneModelIsStable(t *testing.T) {
	root := t.TempDir()
	writeAt(t, filepath.Join(root, "docker-compose.yml"), "version: '3'\n")
	writeAt(t, filepath.Join(root, "apps", "api", "go.mod"), "module pigeon/apps/api\n\ngo 1.21\n")
	path := filepath.Join(root, "apps", "api", "internal", "models", "usage_record.go")
	writeAt(t, path, usageRecordModel)

	first := generate.ModelShape(path)
	second := generate.ModelShape(path)
	if first == "" {
		t.Fatal("no shape")
	}
	if first != second {
		t.Errorf("the shape is not stable: %q then %q", first, second)
	}
	// The pinned model shapes the same; the one with a new column does not.
	writeAt(t, path, usageRecordModelPinned)
	if pinned := generate.ModelShape(path); pinned != first {
		t.Errorf("a gorm type pin changed the shape: %q vs %q", pinned, first)
	}
	writeAt(t, path, usageRecordModelNewColumn)
	if added := generate.ModelShape(path); added == first {
		t.Error("a new column did not change the shape")
	}
}
