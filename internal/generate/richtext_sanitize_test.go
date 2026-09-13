package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A richtext field is HTML that somebody else renders, so the model tags it for
// internal/sanitize, which cleans it on every write.
func TestRichtextFieldsAreTaggedForTheSanitiser(t *testing.T) {
	const module = "blog/apps/api"
	root := setupMinimalProject(t, module)
	g := newTestGenerator(root, module, mustFields(t, "Article", "title:string,body:richtext"))
	if err := g.writeGoModel(g.Names()); err != nil {
		t.Fatalf("model: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "models", "article.go"))
	mustParse(t, "article.go", src)

	var body, title string
	for _, line := range strings.Split(src, "\n") {
		switch {
		case strings.Contains(line, `json:"body"`):
			body = line
		case strings.Contains(line, `json:"title"`):
			title = line
		}
	}
	if !strings.Contains(body, `sanitize:"html"`) {
		t.Errorf("the richtext field is not tagged: %q", body)
	}
	if strings.Contains(title, "sanitize") {
		t.Errorf("a plain string field was tagged: %q", title)
	}
}
