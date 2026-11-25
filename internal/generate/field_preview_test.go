package generate

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// update rewrites the committed copy instead of comparing against it:
//
//	go test ./internal/generate -run TestFieldPreviewDataIsCurrent -update
var update = flag.Bool("update", false, "rewrite docs/config/field-preview.generated.ts")

// previewPath is the committed copy, relative to this package.
const previewPath = "../../docs/config/field-preview.generated.ts"

// The website's field builder describes what the generator produces, and this
// is what stops it describing a framework that no longer exists.
//
// The builder shows, per field, the Go column, the TypeScript, the Zod rule
// and the admin control. Written out by hand in a TSX file, that is a second
// description of the generator, and the day somebody changes GoType or
// FormFieldType the page would be wrong with nothing failing. Nobody reads a
// marketing page against the source.
//
// So the data is computed from the same methods the generator calls, and this
// compares the committed copy against them.
func TestFieldPreviewDataIsCurrent(t *testing.T) {
	want := FieldPreviewTS()

	if *update {
		if err := os.WriteFile(previewPath, []byte(want), 0o644); err != nil { // #nosec G306 -- generated docs data
			t.Fatalf("writing %s: %v", previewPath, err)
		}
		t.Logf("rewrote %s", filepath.Clean(previewPath))
		return
	}

	got, err := os.ReadFile(previewPath) // #nosec G304 -- a path in this repository
	if err != nil {
		t.Fatalf("reading %s: %v\n\nRegenerate it with:\n  go test ./internal/generate -run TestFieldPreviewDataIsCurrent -update", previewPath, err)
	}

	// Compared with line endings normalised: this repository checks out CRLF
	// on Windows, and a test that fails on that would fail for everybody on
	// one platform and nobody on the others.
	if normaliseEOL(string(got)) != normaliseEOL(want) {
		t.Errorf("docs/config/field-preview.generated.ts is behind internal/generate.\n"+
			"The field builder would describe field types the generator no longer produces.\n\n"+
			"Regenerate it with:\n  go test ./internal/generate -run TestFieldPreviewDataIsCurrent -update\n\n%s",
			firstDifference(normaliseEOL(string(got)), normaliseEOL(want)))
	}
}

// Every field type a person can write is in the table, or the builder offers a
// type it cannot describe.
func TestEveryFieldTypeHasAPreview(t *testing.T) {
	previews := map[string]FieldPreview{}
	for _, p := range FieldPreviews() {
		previews[p.Type] = p
	}

	for _, name := range ValidFieldTypes() {
		p, ok := previews[name]
		if !ok {
			t.Errorf("%s has no preview, so the field builder cannot say what it produces", name)
			continue
		}
		// A type with no Go type is one the sample did not fill in enough to
		// render, which would show as a blank cell on the page.
		if p.Base.Go == "" {
			t.Errorf("%s previews an empty Go type", name)
		}
		if p.Base.TS == "" {
			t.Errorf("%s previews an empty TypeScript type", name)
		}
	}
}

// The sample fields are filled in enough to mean something.
//
// A select with no choices, a belongs_to with no model and a vector with no
// dimension are not fields anybody could write, and asking the generator
// about them describes a shape it never emits. This checks the three that
// need more than a name actually got it.
func TestPreviewSamplesAreRealisticEnough(t *testing.T) {
	if got := previewSample("select"); len(got.Options) == 0 {
		t.Error("the select sample has no choices, so its preview describes a dropdown with nothing in it")
	}
	if got := previewSample("belongs_to"); got.RelatedModel == "" {
		t.Error("the belongs_to sample names no model, so its preview cannot show the foreign key")
	}
	if got := previewSample("vector"); got.VectorDim() == 0 {
		t.Error("the vector sample has no dimension, so its preview describes a column the generator does not emit")
	}
	// money is the one people are most often surprised by, so its note has to
	// survive any reshuffle of the table.
	for _, p := range FieldPreviews() {
		if p.Type == "money" && !strings.Contains(p.Note, "minor units") {
			t.Errorf("the money note no longer mentions minor units: %q", p.Note)
		}
	}
}

func normaliseEOL(s string) string { return strings.ReplaceAll(s, "\r\n", "\n") }
