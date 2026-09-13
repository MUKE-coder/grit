package scaffold

import (
	"encoding/json"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The Blog model as v3.241.0 generated it.
const blogModelV241 = "package models\n\ntype Blog struct {\n" +
	"\tID          string         `gorm:\"primarykey;size:36\" json:\"id\"`\n" +
	"\tContent     string         `gorm:\"type:text\" json:\"content\"`\n" +
	"\tExcerpt     string         `gorm:\"size:500\" json:\"excerpt\"`\n}\n"

func TestRepairTagsBlogContentForTheSanitiser(t *testing.T) {
	out, fixed, warn := repairBlogContentSource(blogModelV241)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	if !strings.Contains(out, "`gorm:\"type:text\" json:\"content\" sanitize:\"html\"`") {
		t.Errorf("Content was not tagged:\n%s", out)
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairBlogContentSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the model again")
	}
}

func TestRepairNamesARewrittenBlogModel(t *testing.T) {
	src := strings.Replace(blogModelV241, "gorm:\"type:text\" json:\"content\"", "json:\"content\"", 1)
	out, _, warn := repairBlogContentSource(src)
	if out != src || len(warn) != 1 {
		t.Errorf("a rewritten Content line should be left alone and named, got warn %v", warn)
	}
}

// The scaffold template and the repair must agree, or upgrade would keep
// finding something to fix in a fresh project.
func TestScaffoldBlogModelIsAlreadyTagged(t *testing.T) {
	if !strings.Contains(blogModelGo(), `json:"content" sanitize:"html"`) {
		t.Error("the scaffolded Blog model does not tag Content for the sanitiser")
	}
	if out, _, _ := repairBlogContentSource(blogModelGo()); out != blogModelGo() {
		t.Error("the repair would change a freshly scaffolded Blog model")
	}
}

func TestEnsureDependency(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		"lf":    "{\n  \"dependencies\": {\n    \"react\": \"19.2.7\"\n  }\n}\n",
		"crlf":  "{\r\n  \"dependencies\": {\r\n    \"react\": \"19.2.7\"\r\n  }\r\n}\r\n",
		"root":  "{\n  \"devDependencies\": {\n    \"turbo\": \"^2\"\n  }\n}\n",
		"empty": "{\n  \"peer\": {\"react\": \"19\"},\n  \"dependencies\": {}\n}\n",
	}
	for name, src := range cases {
		path := filepath.Join(dir, name+".json")
		if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
			t.Fatal(err)
		}
		added, err := ensureDependency(path, "dompurify", "^3.4.15")
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		raw, _ := os.ReadFile(path)
		if name == "root" {
			if added || string(raw) != src {
				t.Errorf("root: a package.json without React was changed")
			}
			continue
		}
		var parsed map[string]map[string]string
		if err := json.Unmarshal(raw, &parsed); err != nil {
			t.Fatalf("%s: not valid JSON after the edit: %v\n%s", name, err, raw)
		}
		if !added || parsed["dependencies"]["dompurify"] != "^3.4.15" {
			t.Errorf("%s: dompurify not added:\n%s", name, raw)
		}
		if name == "crlf" && strings.Contains(strings.ReplaceAll(string(raw), "\r\n", ""), "\n") {
			t.Errorf("crlf: mixed line endings:\n%q", raw)
		}
		if again, _ := ensureDependency(path, "dompurify", "^3.4.15"); again {
			t.Errorf("%s: added twice", name)
		}
	}
}

// The package that ships into every project must at least be Go.
func TestSanitizeTemplatesAreGo(t *testing.T) {
	for name, src := range map[string]string{"html.go": apiSanitizeHTMLGo(), "html_test.go": apiSanitizeHTMLTestGo()} {
		if _, err := format.Source([]byte(src)); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	if !strings.Contains(apiGoMod(Options{ProjectName: "x"}), bluemondayModule+" "+bluemondayVersion) {
		t.Error("the scaffold go.mod does not require bluemonday")
	}
}
