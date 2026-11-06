package generate

import (
	"strings"
	"testing"
)

// The precondition, on both sides of the wire.
//
// grit#182: the API has carried a version column, an incrementing hook, an ETag
// and a 409 for a long time, and no generated client ever sent the header. The
// mechanism was inert, and inert is indistinguishable from working until two
// people save at once. These read the emitted source, because a claim about
// what reaches the API is a claim about what the generator wrote.

// generatorFor is a generator over a definition with the given fields.
func generatorFor(t *testing.T, fields []Field) *Generator {
	t.Helper()
	g, err := NewGenerator(&ResourceDefinition{Name: "Widget", Fields: fields})
	if err != nil {
		t.Fatalf("NewGenerator: %v", err)
	}
	g.Module = "example/apps/api"
	return g
}

// The grid sends a version per row, so the service checks one per row. One
// If-Match header cannot carry forty.
func TestBulkEditChecksAVersionPerRow(t *testing.T) {
	g := generatorFor(t, []Field{{Name: "name", Type: "string"}})
	src := g.serviceSource(BuildNames(g.Definition))

	for _, want := range []string{
		// The field it arrives in.
		`Version *int `,
		// The readable per-row refusal, before anything is written.
		"somebody else saved this row since you loaded it",
		// And the WHERE that closes the gap between that check and the write.
		`q = q.Where("version = ?", *w.version)`,
		"&concurrency.ErrConflict{Current: *w.version}",
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the generated service is missing %q, so a grid save overwrites "+
				"a row somebody else edited", want)
		}
	}

	// A row with no version is written as it always was. Making the
	// precondition compulsory would break every script that writes to this API.
	if !strings.Contains(src, "if w.version != nil {") {
		t.Error("the version is applied unconditionally, so a client that sends none is refused")
	}
}

// Every resource ships the test for a stale save, and the grid test only where
// there is a column it can actually change.
func TestGeneratedResourceTestsCoverTheStaleSave(t *testing.T) {
	withText := generatorFor(t, []Field{{Name: "name", Type: "string"}})
	src := withText.resourceTestSource(BuildNames(withText.Definition))

	for _, want := range []string{
		"func TestWidget_StaleSaveIsRefused",
		"func TestWidget_BulkEditRefusesAStaleRow",
		`AssertHeader("ETag"`,
		`AssertCode("VERSION_CONFLICT")`,
		// The routes the two tests need, mounted with the permission routes.go
		// gives them. A test that mounts PUT behind the view permission is
		// testing a router this project does not run.
		`guarded.PUT("/widgets/:id", middleware.RequireRole("ADMIN", "perm:widgets.edit"), handler.Update)`,
		`guarded.POST("/widgets/bulk-edit", middleware.RequireRole("ADMIN", "perm:widgets.edit"), handler.BulkEdit)`,
		// The patch has to change something: BulkEdit skips an edit whose
		// writable columns are all unchanged, so an empty one never reaches the
		// version check and the test would assert a 422 it cannot cause.
		`"name": "edited by the grid test"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("the generated test file is missing %q", want)
		}
	}

	// A resource with no plain text column gets the stale-save test and not the
	// grid one, rather than a grid test that cannot assert what it says.
	noText := generatorFor(t, []Field{{Name: "total", Type: "int"}, {Name: "active", Type: "toggle"}})
	bare := noText.resourceTestSource(BuildNames(noText.Definition))
	if !strings.Contains(bare, "func TestWidget_StaleSaveIsRefused") {
		t.Error("a resource with no text column lost the stale-save test too")
	}
	if strings.Contains(bare, "func TestWidget_BulkEditRefusesAStaleRow") {
		t.Error("a resource with no column the grid test can change still got the grid test")
	}
}

// patchableTextField picks a column the grid test may write.
func TestPatchableTextField(t *testing.T) {
	for name, tc := range map[string]struct {
		fields []Field
		want   string
	}{
		"the first plain string": {
			fields: []Field{{Name: "title", Type: "string"}, {Name: "body", Type: "text"}},
			want:   "title",
		},
		"text counts": {
			fields: []Field{{Name: "total", Type: "int"}, {Name: "notes", Type: "text"}},
			want:   "notes",
		},
		// Writing one value to two rows of a unique column is either a
		// constraint violation or, worse, not one.
		"a unique column does not": {
			fields: []Field{{Name: "sku", Type: "string", Unique: true}, {Name: "label", Type: "string"}},
			want:   "label",
		},
		// A column the server fills is not the caller's to send.
		"a server-set column does not": {
			fields: []Field{{Name: "slug_cache", Type: "string", ServerSet: true}, {Name: "label", Type: "string"}},
			want:   "label",
		},
		// money is two columns, an enum rejects anything off its list, and a
		// date wants a layout. None of them is writable as a bare string.
		"nothing writable as a plain string": {
			fields: []Field{{Name: "price", Type: "money"}, {Name: "due", Type: "date"}},
			want:   "",
		},
		"snake-cased, as the json tag spells it": {
			fields: []Field{{Name: "DisplayName", Type: "string"}},
			want:   "display_name",
		},
	} {
		t.Run(name, func(t *testing.T) {
			g := generatorFor(t, tc.fields)
			key, value := g.patchableTextField()
			if key != tc.want {
				t.Errorf("patchableTextField() = %q, want %q", key, tc.want)
			}
			if tc.want != "" && value == "" {
				t.Error("a key with no value to write")
			}
		})
	}
}
