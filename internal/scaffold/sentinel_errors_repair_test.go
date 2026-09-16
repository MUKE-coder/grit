package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// M35: no template compares a sentinel error with == any more.
//
// The check is over the scaffold's own source, because the comparison is
// emitted from a dozen different template functions and naming them one by one
// is how the thirteenth gets missed.
func TestNoTemplateComparesASentinelWithEquals(t *testing.T) {
	bad := regexp.MustCompile(`\b(?:err|e|cerr|derr|ferr|serr|terr|txErr|rerr|werr|result\.Error|res\.Error|tx\.Error|db\.Error) (?:==|!=) [a-z]\w*\.Err[A-Z]\w*`)
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		// The repair itself quotes the comparison it removes.
		if name == "sentinel_errors_repair.go" {
			continue
		}
		raw, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		for _, hit := range bad.FindAllString(string(raw), -1) {
			t.Errorf("%s: %q compares a sentinel error with ==; a %%w wrap turns that check off silently. Use errors.Is", name, hit)
		}
	}
}

func TestRepairSentinelErrorSource(t *testing.T) {
	src := `package handlers

import (
	"net/http"

	"gorm.io/gorm"
)

func load(db *gorm.DB) error {
	err := db.First(nil).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	if serve() != http.ErrServerClosed {
		return nil
	}
	return nil
}

func serve() error { return nil }
`
	out, fixed, warn := repairSentinelErrorSource(src)
	if len(warn) != 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v warnings %v", fixed, warn)
	}
	if !strings.Contains(out, "if errors.Is(err, gorm.ErrRecordNotFound) {") {
		t.Error("== was not rewritten")
	}
	if !strings.Contains(out, "if !errors.Is(err, gorm.ErrRecordNotFound) {") {
		t.Error("!= was not rewritten")
	}
	// A call is not a variable that holds an error; leave it alone rather than
	// produce errors.Is(serve(), …) from a shape the repair cannot reason about.
	if !strings.Contains(out, "if serve() != http.ErrServerClosed {") {
		t.Error("the repair rewrote something that is not one of the named error variables")
	}
	if !strings.Contains(out, `"errors"`) {
		t.Error("the errors import was not added")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Errorf("the repaired file is not valid Go: %v", err)
	}
	if again, fixed, _ := repairSentinelErrorSource(out); again != out || len(fixed) != 0 {
		t.Error("the repair is not idempotent")
	}
}

// A switch on err compares with ==, so it is the same defect with different
// punctuation. A switch holding anything the repair cannot reason about is
// left whole rather than half-rewritten.
func TestRepairRewritesSwitchOnError(t *testing.T) {
	src := `package handlers

import (
	"gorm.io/gorm"
)

func code(err error) string {
	out := "BAD_REQUEST"
	switch err {
	case ErrReviewClosed:
		out = "REVIEW_CLOSED"
	case gorm.ErrRecordNotFound:
		out = "NOT_FOUND"
	}
	return out
}

func kind(err error) string {
	switch err {
	case ErrReviewClosed:
		return "closed"
	case wrap(err):
		return "wrapped"
	}
	return ""
}

func wrap(err error) error { return err }

var ErrReviewClosed = wrap(nil)
`
	out, fixed, warn := repairSentinelErrorSource(src)
	if len(warn) != 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v warnings %v", fixed, warn)
	}
	for _, want := range []string{
		"\tswitch {\n\tcase errors.Is(err, ErrReviewClosed):",
		"case errors.Is(err, gorm.ErrRecordNotFound):",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the rewritten switch is missing %q", want)
		}
	}
	if !strings.Contains(out, "\tswitch err {\n\tcase ErrReviewClosed:\n\t\treturn \"closed\"") {
		t.Error("a switch with a case the repair cannot reason about should be left whole")
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Errorf("the repaired file is not valid Go: %v", err)
	}
}

// A project already on errors.Is is left exactly as it is.
func TestRepairSentinelErrorsLeavesACleanProjectAlone(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "apps", "api", "internal", "handlers")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	src := "package handlers\n\nimport (\n\t\"errors\"\n\n\t\"gorm.io/gorm\"\n)\n\nvar _ = errors.Is(gorm.ErrRecordNotFound, gorm.ErrRecordNotFound)\n"
	path := filepath.Join(dir, "thing.go")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := repairSentinelErrors(root, Options{ProjectName: "app"}); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != src {
		t.Error("a file with no == comparison was rewritten")
	}
}
