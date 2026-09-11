package generate

import (
	"path/filepath"
	"strings"
	"testing"
)

// The database is reached through the request context, and only by the service.
//
// Querying without the request context hands GORM context.Background(), so
// nothing a middleware put on the request can reach a callback. That is why the
// multitenant plugin could not scope a generated resource at all: its
// middleware resolves the active organization onto c.Request, and the scoping
// callback looked somewhere that never had it. Found by building a CRM on the
// plugin and watching Globex read Acme's contacts.
//
// And a handler that runs its own queries is logic nothing else can reuse: it
// did 27 of them while the service beside it was called by nothing.
func contactHandlerAndService(t *testing.T, fields string) (string, string) {
	t.Helper()
	const module = "crm/apps/api"
	root := setupMinimalProject(t, module)
	def, err := ParseInlineFields("Contact", fields)
	if err != nil {
		t.Fatalf("ParseInlineFields: %v", err)
	}
	g := newTestGenerator(root, module, def)
	names := g.Names()
	if err := g.writeGoHandler(names); err != nil {
		t.Fatalf("handler: %v", err)
	}
	if err := g.writeGoService(names); err != nil {
		t.Fatalf("service: %v", err)
	}
	h := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "handlers", "contact.go"))
	s := readTestFile(t, filepath.Join(root, "apps", "api", "internal", "services", "contact.go"))
	mustParse(t, "handlers/contact.go", h)
	mustParse(t, "services/contact.go", s)
	return h, s
}

func TestGeneratedHandlerQueriesThroughTheRequestContext(t *testing.T) {
	h, s := contactHandlerAndService(t, "name:string,company:belongs_to:Company,tags:many_to_many:Tag")

	// The handler hands the service the request's context, with the caller on it.
	if !strings.Contains(h, "func (h *ContactHandler) ctx(c *gin.Context) context.Context") ||
		!strings.Contains(h, "authz.WithActor(c.Request.Context(), authz.ActorOf(c))") {
		t.Fatal("the handler does not pass the request context, with its caller, to the service")
	}

	// It runs no query of its own. Its database is only there to build the service.
	for _, line := range strings.Split(h, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		if strings.Contains(line, "h.DB") && !strings.Contains(line, "&services.ContactService{DB: h.DB") {
			t.Errorf("the handler reaches the database itself:\n  %s", trimmed)
		}
		for _, q := range []string{".Where(", ".First(", ".Find(", ".Transaction(", ".Association(", "FindInBatches(", ".Preload(", ".Scopes(", ".Model("} {
			if strings.Contains(line, q) {
				t.Errorf("the handler still builds a query (%s):\n  %s", q, trimmed)
			}
		}
	}

	// Every service query goes through the context: s.DB only inside db().
	for _, line := range strings.Split(s, "\n") {
		if !strings.Contains(line, "s.DB") || strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		if !strings.Contains(line, "return s.DB.WithContext(ctx)") {
			t.Errorf("this query bypasses the request context, so a scoping "+
				"callback cannot see the request:\n  %s", strings.TrimSpace(line))
		}
	}
}

// The many-to-many writes are built far from the template body, so name them:
// through the transaction on create and update, through the context on patch,
// and never with the error dropped.
func TestComputedSlotsAreScopedToo(t *testing.T) {
	_, s := contactHandlerAndService(t, "name:string,tags:many_to_many:Tag")

	if strings.Contains(s, "s.DB.Model(") {
		t.Error("a write bypasses the request context")
	}
	for _, sig := range []string{"func (s *ContactService) Create(", "func (s *ContactService) Update("} {
		body := method(t, s, sig)
		if !strings.Contains(body, `tx.Model(item).Association("Tags").Replace(relatedTags); err != nil`) {
			t.Errorf("%s does not set the tags in its transaction, checking the error", sig)
		}
		if !strings.Contains(body, "db.Transaction(func(tx *gorm.DB) error {") {
			t.Errorf("%s does not write the row and its links together", sig)
		}
	}
	if !strings.Contains(method(t, s, "func (s *ContactService) Patch("), `db.Model(item).Association("Tags").Replace(related); err != nil`) {
		t.Error("Patch does not replace the tags through the context, checking the error")
	}
}

// A many-to-many id that matches nothing is refused, not dropped.
//
// The lookup found what it could and the set was replaced with that, so a PUT
// whose only tag id was mistyped answered 200 and left the row with no tags.
// The role assignment endpoint already refused an unknown role the same way.
func TestAnUnknownLinkIsRefusedNotDropped(t *testing.T) {
	_, s := contactHandlerAndService(t, "name:string,tags:many_to_many:Tag")
	for _, sig := range []string{"func (s *ContactService) Create(", "func (s *ContactService) Update(", "func (s *ContactService) Patch("} {
		if !strings.Contains(method(t, s, sig), `respond.Rule("tag_ids: %d of the %d ids given do not exist"`) {
			t.Errorf("%s drops an id that matches no tag, so a write naming one empties the set", sig)
		}
	}
}
