package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// internal/docx reaches a project, and reaches an upgraded one too.
//
// grit#212: Grit has no way to write a Word document, and a resource's list
// export is not one. A list export streams rows, and what people ask Word for
// is a document with a page size, a named font at a given size, justified
// text, a line spacing and a page break before each chapter. Those four are
// often a marking scheme rather than a preference, so getting them wrong costs
// the person something.
//
// The package goes in writeCodegenRuntimeFiles rather than the new-project map
// because that function runs on upgrade as well. A framework package that only
// new projects get is the shape that shipped media, recovery contacts and
// passkeys half finished.
func TestDocxReachesAProject(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "docxapp", Architecture: ArchAPI}

	if err := writeCodegenRuntimeFiles(root, opts); err != nil {
		t.Fatalf("writing the codegen runtime: %v", err)
	}

	apiRoot := opts.APIRoot(root)
	for _, name := range []string{"docx.go", "report.go", "docx_test.go"} {
		path := filepath.Join(apiRoot, "internal", "docx", name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("internal/docx/%s was not written: %v", name, err)
		}
	}

	src := readFileString(t, filepath.Join(apiRoot, "internal", "docx", "docx.go"))

	// Standard library only. The Go libraries for this are either thinly
	// maintained or AGPL, and an AGPL dependency inside a framework other
	// people build commercial products on is a licence decision nobody would
	// have chosen deliberately.
	for _, banned := range []string{"github.com/", "gitlab.com/", "gopkg.in/"} {
		if strings.Contains(src, banned) {
			t.Errorf("internal/docx imports %s: it is meant to be standard library only", banned)
		}
	}

	// The primitives a submission scheme is written in terms of.
	for _, fn := range []string{
		"func New(",
		"func (d *Doc) Chapter(",
		"func (d *Doc) Heading(",
		"func (d *Doc) Para(",
		"func (d *Doc) Table(",
		"func Attachment(",
		"const ContentType",
	} {
		if !strings.Contains(src, fn) {
			t.Errorf("internal/docx/docx.go has no %s", fn)
		}
	}

	report := readFileString(t, filepath.Join(apiRoot, "internal", "docx", "report.go"))
	if !strings.Contains(report, "func Submission()") {
		t.Error("internal/docx/report.go has no Submission(): the Times New Roman 12pt, " +
			"justified, 1.5 spacing combination is the one schemes ask for")
	}
	if !strings.Contains(report, "func RenderReport(") {
		t.Error("internal/docx/report.go has no RenderReport()")
	}
}

// The package ships its own tests, so a project can run them.
//
// internal/pdf ships pdf_test.go for the same reason: a framework package a
// project cannot test is one the project has to trust.
func TestDocxShipsItsOwnTests(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "docxapp", Architecture: ArchAPI}
	if err := writeCodegenRuntimeFiles(root, opts); err != nil {
		t.Fatalf("writing the codegen runtime: %v", err)
	}

	src := readFileString(t, filepath.Join(opts.APIRoot(root), "internal", "docx", "docx_test.go"))
	for _, name := range []string{
		"func TestEveryRequiredPartIsPresentAndParses(",
		"func TestContentTypesIsTheFirstEntry(",
		"func TestATableCarriesItsGrid(",
		"func TestAttachmentQuotesTheFilename(",
	} {
		if !strings.Contains(src, name) {
			t.Errorf("the shipped docx_test.go has no %s", name)
		}
	}
}

func readFileString(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}
