package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// Two people saving the same record: the second silently overwrote the first.
// Models carried a version that every update incremented, and nothing checked
// it. Update and Patch now scope the write to the version the client sent.
func TestUpdatesHonourIfMatch(t *testing.T) {
	const module = "shop/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Lot", "title:string,bid:int")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	if err := g.writeGoHandler(g.Names()); err != nil {
		t.Fatalf("handler: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "lot.go"))
	mustParse(t, "lot.go", h)

	for _, sig := range []string{"func (h *LotHandler) Update(", "func (h *LotHandler) Patch("} {
		body := method(t, h, sig)
		if !strings.Contains(body, "Scopes(concurrency.IfMatch(c))") {
			t.Errorf("%s does not scope its write to the version the client read", sig)
		}
		if !strings.Contains(body, "concurrency.Conflicted(c, written)") {
			t.Errorf("%s does not answer a stale version with a conflict", sig)
		}
		if !strings.Contains(body, `c.Header("ETag", concurrency.Tag(item.Version))`) {
			t.Errorf("%s does not return the new version", sig)
		}
	}
	if !strings.Contains(method(t, h, "func (h *LotHandler) GetByID("), `c.Header("ETag", concurrency.Tag(item.Version))`) {
		t.Error("a read does not say which version to send back")
	}
	if !strings.Contains(h, `"shop/apps/api/internal/concurrency"`) {
		t.Error("the concurrency import is missing")
	}
}
