package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// A generated handler must talk to the database through the request context.
//
// Querying h.DB directly hands GORM context.Background(), so nothing a
// middleware put on the request can reach a callback. That is why the
// multitenant plugin could not scope a generated resource at all: its
// middleware resolves the active organization onto c.Request, and the scoping
// callback looked somewhere that never had it. Unregistered, that returned
// every tenant's rows to every tenant; registered, every query failed closed
// with a 500. Found by building a CRM on the plugin and watching Globex read
// Acme's contacts.
//
// It is also what makes a cancelled request cancel its query instead of
// holding a connection to build a response nobody will read.
func TestGeneratedHandlerQueriesThroughTheRequestContext(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, err := ParseInlineFields("Contact", "name:string,company:belongs_to:Company,tags:many_to_many:Tag")
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "contact.go"))

	if !strings.Contains(src, "func (h *ContactHandler) scoped(c *gin.Context) *gorm.DB") {
		t.Fatal("no scoped() helper: handlers have no way to reach the request context")
	}

	// Exactly one h.DB may survive, inside the helper itself. Everything else,
	// including the code built from computed slots (the create call, the
	// reload, the many-to-many association writes), has to go through it.
	for _, line := range strings.Split(src, "\n") {
		if !strings.Contains(line, "h.DB") {
			continue
		}
		if strings.Contains(line, "return h.DB.WithContext(c.Request.Context())") {
			continue
		}
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		t.Errorf("this query bypasses the request context, so a scoping "+
			"callback cannot see the request:\n  %s", strings.TrimSpace(line))
	}
}

// The slot-built database calls are the ones most easily missed, so name them.
func TestComputedSlotsAreScopedToo(t *testing.T) {
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)

	def, _ := ParseInlineFields("Contact", "name:string,tags:many_to_many:Tag")
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	src := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "contact.go"))

	// The many-to-many association writer comes from m2mCreateCode, built far
	// from the template body.
	if strings.Contains(src, "h.DB.Model(&item).Association") {
		t.Error("the many-to-many association write bypasses the request context")
	}
	if !strings.Contains(src, "h.scoped(c).Model(&item).Association") {
		t.Error("expected the association write to go through scoped()")
	}
}
