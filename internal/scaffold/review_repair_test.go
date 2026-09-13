package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// contactRoutesV239 is contact_routes.go as v3.239.0 generated it, copied from
// the reviewed project.
const contactRoutesV239 = `package routes

import (
	"contact-app/apps/api/internal/handlers"
	"contact-app/apps/api/internal/middleware"
)

func init() {
	RegisterRoutes(func(m *Mount) {
		h := &handlers.ContactHandler{
			DB: m.DB,
		}

		m.Protected.GET("/contacts", h.List)
		m.Protected.GET("/contacts/export", h.Export)
		m.Protected.POST("/contacts/import", h.Import)
		m.Protected.GET("/contacts/import/template", h.Template)
		m.Protected.GET("/contacts/:id", h.GetByID)
		m.Protected.GET("/contacts/:id/pdf", h.PDF)
		m.Protected.POST("/contacts", h.Create)
		m.Protected.PUT("/contacts/:id", h.Update)
		m.Protected.PATCH("/contacts/:id", h.Patch)

		m.Staff.DELETE("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.delete"), h.Delete)
		m.Staff.POST("/contacts/bulk", middleware.RequireRole("ADMIN", "perm:contacts.delete"), h.Bulk)
	})
}
`

// C3: the reviewed project's address book was readable and writable by an
// account registered a minute earlier.
func TestRepairPutsSharedRoutesBehindAPermission(t *testing.T) {
	out, fixed, warn := repairRouteSource(contactRoutesV239, map[string]bool{}, true, "contact-app/apps/api")
	if len(warn) > 0 {
		t.Errorf("warnings on an untouched generated file: %v", warn)
	}
	if len(fixed) != 1 || !strings.Contains(fixed[0], "9 shared routes") {
		t.Errorf("reported %v, want the 9 routes it moved", fixed)
	}
	mustFormat(t, "the route file", out)
	for _, want := range []string{
		`m.Staff.GET("/contacts", middleware.RequireRole("ADMIN", "perm:contacts.view"), h.List)`,
		`m.Staff.POST("/contacts/import", middleware.RequireRole("ADMIN", "perm:contacts.create"), h.Import)`,
		`m.Staff.GET("/contacts/:id/pdf", middleware.RequireRole("ADMIN", "perm:contacts.view"), h.PDF)`,
		`m.Staff.PATCH("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.edit"), h.Patch)`,
		`m.Staff.DELETE("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.delete"), h.Delete)`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("the repaired file is missing %s", want)
		}
	}
	if strings.Contains(out, "m.Protected.") {
		t.Errorf("a route is still open to any signed-in account:\n%s", out)
	}
	if again, fixed, _ := repairRouteSource(out, map[string]bool{}, true, "contact-app/apps/api"); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the file again")
	}
}

// An owned resource's handler scopes every row to the caller, so its protected
// routes are the design, not the bug.
func TestRepairLeavesAnOwnedResourcesRoutes(t *testing.T) {
	out, fixed, _ := repairRouteSource(contactRoutesV239, map[string]bool{"Contact": true}, true, "contact-app/apps/api")
	if out != contactRoutesV239 || len(fixed) > 0 {
		t.Error("an owned resource's routes were moved")
	}
}

// A project older than the staff group has nowhere to put a permission, so
// shared data goes to ADMIN, and no import is added that nothing would use.
func TestRepairWithoutAStaffGroupUsesTheAdminGroup(t *testing.T) {
	src := strings.Replace(contactRoutesV239, "\t\"contact-app/apps/api/internal/middleware\"\n", "", 1)
	src = strings.ReplaceAll(src, `m.Staff.DELETE("/contacts/:id", middleware.RequireRole("ADMIN", "perm:contacts.delete"), h.Delete)`, `m.Admin.DELETE("/contacts/:id", h.Delete)`)
	src = strings.ReplaceAll(src, `m.Staff.POST("/contacts/bulk", middleware.RequireRole("ADMIN", "perm:contacts.delete"), h.Bulk)`, `m.Admin.POST("/contacts/bulk", h.Bulk)`)

	out, _, _ := repairRouteSource(src, map[string]bool{}, false, "contact-app/apps/api")
	mustFormat(t, "the route file", out)
	if !strings.Contains(out, `m.Admin.GET("/contacts", h.List)`) || strings.Contains(out, "m.Protected.") {
		t.Errorf("shared routes did not move to the admin group:\n%s", out)
	}
	if strings.Contains(out, "internal/middleware") {
		t.Error("added a middleware import nothing uses")
	}
}

// A staff project whose route file never imported middleware gets the import.
func TestRepairAddsTheMiddlewareImport(t *testing.T) {
	src := strings.Replace(contactRoutesV239, "\t\"contact-app/apps/api/internal/middleware\"\n", "", 1)
	src = strings.Replace(src, "\t\tm.Staff.DELETE", "\t\tm.Admin.DELETE", 1)
	out, _, _ := repairRouteSource(src, map[string]bool{}, true, "contact-app/apps/api")
	if !strings.Contains(out, `"contact-app/apps/api/internal/middleware"`) {
		t.Errorf("the middleware import was not added:\n%s", out)
	}
}

// Owned and tenant-owned resources are both scoped, so neither is moved.
func TestScopedResourcesAreFoundByHandlerAndModel(t *testing.T) {
	api := t.TempDir()
	for rel, body := range map[string]string{
		"internal/handlers/note.go": "package handlers\n\nfunc (h *NoteHandler) List(c *gin.Context) {\n\tauthz.ScopeToOwner(c, query, \"user_id\")\n}\n",
		"internal/handlers/deal.go": "package handlers\n\nfunc (h *DealHandler) List(c *gin.Context) {}\n",
		"internal/handlers/lot.go":  "package handlers\n\nfunc (h *LotHandler) List(c *gin.Context) {}\n",
		"internal/models/deal.go":   "package models\n\ntype Deal struct {\n\ttenant.Owned\n\tTitle string\n}\n",
		"internal/models/lot.go":    "package models\n\ntype Lot struct {\n\tID    string\n\tTitle string\n}\n",
	} {
		path := filepath.Join(api, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	owned, err := ownedResourceHandlers(filepath.Join(api, "internal", "handlers"))
	if err != nil {
		t.Fatal(err)
	}
	if !owned["Note"] || !owned["Deal"] || owned["Lot"] {
		t.Errorf("scoped resources: %v, want Note and Deal only", owned)
	}
}

// Routes the developer added are theirs.
func TestRepairLeavesHandWrittenRoutes(t *testing.T) {
	src := strings.Replace(contactRoutesV239, "\t\tm.Protected.GET(\"/contacts\", h.List)\n",
		"\t\tm.Protected.GET(\"/contacts\", h.List)\n\t\tm.Protected.GET(\"/contacts/mine\", h.Mine)\n", 1)
	out, _, _ := repairRouteSource(src, map[string]bool{}, true, "contact-app/apps/api")
	if !strings.Contains(out, `m.Protected.GET("/contacts/mine", h.Mine)`) {
		t.Error("a hand-written route was moved")
	}
}

// C1: register, then push {"op":"update","model":"users","data":{"role":"ADMIN"}}.
func TestRepairTakesUsersAndUploadsOutOfSync(t *testing.T) {
	src := "package routes\n\nfunc setup() {\n" +
		"\tsyncRegistry.Register(\"users\", &models.User{})\n" +
		"\tsyncRegistry.Register(\"uploads\", &models.Upload{})\n" +
		"\tsyncRegistry.Register(\"blogs\", &models.Blog{})\n" +
		"\tsyncRegistry.Register(\"contacts\", &models.Contact{})\n}\n"
	out, fixed, _ := repairSyncRegistrySource(src)
	if strings.Contains(out, `"users"`) || strings.Contains(out, `"uploads"`) {
		t.Errorf("users or uploads are still synced:\n%s", out)
	}
	if !strings.Contains(out, `syncRegistry.Register("contacts", &models.Contact{})`) || len(fixed) != 1 {
		t.Error("the project's own synced models were touched")
	}
}

// H8: .env.example held the JWT secret, database password and dashboard
// passwords that were in .env.
func TestScrubEnvExampleReplacesGeneratedSecrets(t *testing.T) {
	root := t.TempDir()
	secret := strings.Repeat("ab12", 16)
	src := "APP_NAME=contact-app\r\nJWT_SECRET=" + secret + "\r\nPOSTGRES_PASSWORD=" + secret + secret + "\r\n" +
		"FIELD_ENCRYPTION_KEY=q83vEjRWeJq83vEjRWeJq83vEjRWeJq83vEjRWeJ=\r\nMINIO_ACCESS_KEY=minioadmin\r\n"
	path := filepath.Join(root, ".env.example")
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := scrubEnvExample(root); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(raw)
	if strings.Contains(got, secret) || strings.Contains(got, "q83vEj") {
		t.Errorf("a secret survived:\n%s", got)
	}
	want := "APP_NAME=contact-app\r\nJWT_SECRET=CHANGE_ME\r\nPOSTGRES_PASSWORD=CHANGE_ME\r\n" +
		"FIELD_ENCRYPTION_KEY=CHANGE_ME\r\nMINIO_ACCESS_KEY=minioadmin\r\n"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}

// The desktop app's background sync would fail every round on the two models
// the API no longer syncs.
func TestRepairDesktopSyncTables(t *testing.T) {
	src := "package main\n\nvar syncTables = []string{\n\t\"users\",\n\t\"uploads\",\n\t\"contacts\",\n\t// grit:sync-tables\n}\n\nvar other = []string{\"users\",\n}\n"
	out, fixed, _ := repairDesktopSyncTablesSource(src)
	mustFormat(t, "app.go", out)
	want := "package main\n\nvar syncTables = []string{\n\t\"contacts\",\n\t// grit:sync-tables\n}\n\nvar other = []string{\"users\",\n}\n"
	if out != want || len(fixed) != 1 {
		t.Errorf("got %q (%v)", out, fixed)
	}
}

func TestIgnoreLocalDatabasesOnce(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, ".gitignore")
	if err := os.WriteFile(path, []byte("node_modules\nsentinel.db\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if err := ignoreLocalDatabases(root); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(raw), "\n*.db\n"); n != 1 {
		t.Errorf("*.db appears %d times, want 1:\n%s", n, raw)
	}
}
