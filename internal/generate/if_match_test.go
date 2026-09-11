package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// Two people saving the same record: the second silently overwrote the first.
// Models carried a version that every update incremented, and nothing checked
// it. The handler now reads the version the client sent, the service scopes
// the write to it, and a miss comes back as a 409 naming the current one.
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
	if err := g.writeGoService(g.Names()); err != nil {
		t.Fatalf("service: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "lot.go"))
	s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "lot.go"))
	mustParse(t, "handlers/lot.go", h)
	mustParse(t, "services/lot.go", s)

	for _, sig := range []string{"func (h *LotHandler) Update(", "func (h *LotHandler) Patch("} {
		body := method(t, h, sig)
		if !strings.Contains(body, "concurrency.FromRequest(c)") {
			t.Errorf("%s does not pass on the version the client read", sig)
		}
		if !strings.Contains(body, `c.Header("ETag", concurrency.Tag(item.Version))`) {
			t.Errorf("%s does not return the new version", sig)
		}
	}
	for _, sig := range []string{"func (s *LotService) Update(", "func (s *LotService) Patch("} {
		body := method(t, s, sig)
		if !strings.Contains(body, "Scopes(pre.Scope)") {
			t.Errorf("%s does not scope its write to the version the client read", sig)
		}
		if !strings.Contains(body, "pre.Missed(written)") || !strings.Contains(body, "s.conflict(ctx, item.ID)") {
			t.Errorf("%s does not answer a stale version with a conflict", sig)
		}
	}
	if !strings.Contains(method(t, s, "func (s *LotService) conflict("), "&concurrency.ErrConflict{Current: current.Version}") {
		t.Error("the conflict does not carry the version the record is at")
	}
	if !strings.Contains(method(t, h, "func (h *LotHandler) fail("), "concurrency.WriteConflict(c, conflict.Current)") {
		t.Error("the handler does not answer a conflict with 409 and the current version")
	}
	if !strings.Contains(method(t, h, "func (h *LotHandler) GetByID("), `c.Header("ETag", concurrency.Tag(item.Version))`) {
		t.Error("a read does not say which version to send back")
	}
	for name, src := range map[string]string{"handler": h, "service": s} {
		if !strings.Contains(src, `"shop/apps/api/internal/concurrency"`) {
			t.Errorf("the %s does not import concurrency", name)
		}
	}
}
