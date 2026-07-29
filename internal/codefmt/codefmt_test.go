package codefmt

import "testing"

// The two properties that matter for generated models: struct tags line up, and
// imports inside a group get sorted. Both are things the string-concatenating
// generators cannot reliably do themselves.
func TestGoAlignsStructFieldsAndSortsImports(t *testing.T) {
	src := `package models

import (
	"time"

	"gorm.io/gorm"

	"app/internal/ids"
	"app/internal/files"
)

type Asset struct {
	ID string ` + "`" + `json:"id"` + "`" + `
	Title string ` + "`" + `json:"title"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at"` + "`" + `
	DeletedAt gorm.DeletedAt ` + "`" + `json:"-"` + "`" + `
	Cover *files.FileRef ` + "`" + `json:"cover"` + "`" + `
	Ref string ` + "`" + `json:"ref"` + "`" + `
}

func (m *Asset) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" { m.ID = ids.New() }
	return nil
}
`
	got := Go(src)

	if got == src {
		t.Fatal("expected formatting to change the source")
	}
	// files sorts before ids; the generator appended them in that order but a
	// future one might not.
	filesAt, idsAt := indexOf(got, `"app/internal/files"`), indexOf(got, `"app/internal/ids"`)
	if filesAt == -1 || idsAt == -1 {
		t.Fatal("imports were dropped by formatting")
	}
	if filesAt > idsAt {
		t.Errorf("imports not sorted: files at %d, ids at %d", filesAt, idsAt)
	}
	// Aligned tags mean the struct block contains runs of padding.
	if !contains(got, "\tTitle     string") {
		t.Errorf("struct fields not aligned:\n%s", got)
	}
	// The one-line if must be split across lines.
	if contains(got, `if m.ID == "" { m.ID = ids.New() }`) {
		t.Error("statement body not expanded onto its own line")
	}
}

// A parse failure must return the input untouched. This is what keeps the
// formatter from being able to regress scaffolding: broken input is written
// exactly as it would have been before codefmt existed, and the compiler
// reports it against the generated project.
func TestGoReturnsInputUnchangedWhenItDoesNotParse(t *testing.T) {
	broken := "package models\n\nfunc Oops( {\n"
	if got := Go(broken); got != broken {
		t.Errorf("expected unparseable source returned verbatim, got:\n%s", got)
	}
}

// Formatting must be idempotent, or every regeneration would churn the diff.
func TestGoIsIdempotent(t *testing.T) {
	src := "package m\n\nimport \"fmt\"\n\nfunc F() { fmt.Println(1) }\n"
	once := Go(src)
	if twice := Go(once); twice != once {
		t.Errorf("formatting is not idempotent:\nonce:\n%s\ntwice:\n%s", once, twice)
	}
}

// Non-Go files must pass through byte-for-byte — TypeScript, JSON and Markdown
// all go through the same writers.
func TestFileOnlyTouchesGoFiles(t *testing.T) {
	ts := "export const a = {  b : 1 }\n"
	for _, path := range []string{"a.ts", "a.tsx", "a.json", "a.md", "go.mod", "a.gohtml"} {
		if got := File(path, ts); got != ts {
			t.Errorf("%s: expected untouched, got %q", path, got)
		}
	}

	goSrc := "package m\nfunc F()  {}\n"
	if got := File("a.go", goSrc); got == goSrc {
		t.Error("a.go: expected formatting to apply")
	}
	if got := File("A.GO", goSrc); got == goSrc {
		t.Error("A.GO: extension match must be case-insensitive")
	}
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

func contains(haystack, needle string) bool { return indexOf(haystack, needle) != -1 }
