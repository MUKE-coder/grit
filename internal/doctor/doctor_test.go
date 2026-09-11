package doctor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A project that is right must produce nothing at all. That is the half of a
// linter that decides whether anyone leaves it switched on, so it is the first
// test, and every case below starts from it and breaks exactly one thing.
func cleanProject() map[string]string {
	return map[string]string{
		".env": `APP_ENV=development
JWT_SECRET=4f8a1d9c2b7e6053f1a8c4d7b2e90563a1c8f4d7b2e905634f8a1d9c2b7e6053
FIELD_ENCRYPTION_KEY=7Gq0W0l5dQ9xT2bV8nK4mJ6pR1sY3zA5cE7hL9oQ2uI=
SENTINEL_PASSWORD=3f9c1a7d5b2e8046
SENTINEL_SECRET_KEY=8b1e4d7a2c9f6053b8e1d4a7c2f90563
SENTINEL_AUDIT_KEY=2c9f6053b8e1d4a78b1e4d7a2c9f6053
PULSE_PASSWORD=c4d7b2e905634f8a
GORM_STUDIO_USERNAME=admin
GORM_STUDIO_PASSWORD=9d2b7e6053f1a8c4
REDIS_HOST=localhost
`,
		"apps/api/go.mod": `module shop/apps/api

go 1.25.0

require (
	github.com/MUKE-coder/gorm-studio v1.1.0
	github.com/MUKE-coder/pulse v1.0.0
	github.com/MUKE-coder/sentinel/v2 v2.5.0
)
`,
		"apps/api/internal/models/invoice.go": "package models\n\n" +
			"type Invoice struct {\n" +
			"\tID     string                 `gorm:\"primaryKey\" json:\"id\"`\n" +
			"\tNumber string                 `gorm:\"size:80\" json:\"number\"`\n" +
			"\tSecret crypto.EncryptedString `gorm:\"type:text\" json:\"secret\"`\n" +
			"\tUserID string                 `gorm:\"size:36;index\" json:\"user_id\"`\n" +
			"\tUser   User                   `gorm:\"foreignKey:UserID\" json:\"user,omitempty\"`\n" +
			"}\n\n" +
			"func (m *Invoice) GetOwnerID() string { return m.UserID }\n",
		"apps/api/internal/services/invoice.go": `package services

var invoiceListConfig = paginate.Config{
	Searchable: []string{"number"},
	Sortable:   map[string]bool{"id": true, "number": true},
	Filterable: map[string]bool{"id": true, "number": true},
}

func (s *InvoiceService) List(ctx context.Context, p paginate.Params, archived string) (paginate.Result[models.Invoice], error) {
	query := s.db(ctx).Model(&models.Invoice{})
	query = authz.ScopeOwned(ctx, query, "user_id")
	return paginate.List[models.Invoice](query, p, invoiceListConfig)
}

func (s *InvoiceService) Export(ctx context.Context, search string, each func(rows []models.Invoice) error) error {
	query := s.db(ctx).Model(&models.Invoice{})
	query = authz.ScopeOwned(ctx, query, "user_id")
	return nil
}

func (s *InvoiceService) GetByID(ctx context.Context, id string) (*models.Invoice, error) {
	var item models.Invoice
	if !authz.Owns(ctx, &item) {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (s *InvoiceService) load(ctx context.Context, id string) (*models.Invoice, error) {
	var item models.Invoice
	if !authz.Owns(ctx, &item) {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (s *InvoiceService) Bulk(ctx context.Context, action string, ids []string, patch map[string]interface{}) (InvoiceBulkResult, error) {
	scope := s.db(ctx).Model(&models.Invoice{}).Where("id IN ?", ids)
	scope = authz.ScopeOwned(ctx, scope, "user_id")
	return InvoiceBulkResult{}, nil
}
`,
		"apps/api/internal/handlers/invoice.go": "package handlers\n\n" +
			"type CreateInvoiceRequest struct {\n" +
			"\tNumber string `json:\"number\" binding:\"required\"`\n" +
			"}\n",
		// A framework table, with the handler and service the scaffold ships for
		// it and none of the generator's marks. It references a user, as a
		// session does, and it must not be read as a user-owned resource: asked
		// the loose way, every project reported four of these.
		"apps/api/internal/models/session.go": "package models\n\n" +
			"type Session struct {\n" +
			"\tID     string `gorm:\"primaryKey\" json:\"id\"`\n" +
			"\tUserID string `gorm:\"size:36;index\" json:\"user_id\"`\n" +
			"\tUser   User   `gorm:\"foreignKey:UserID\" json:\"user,omitempty\"`\n" +
			"}\n",
		"apps/api/internal/handlers/session.go": "package handlers\n\n" +
			"func (h *SessionHandler) List(c *gin.Context) {}\n",
		"apps/api/internal/services/session.go": "package services\n\n" +
			"func (s *SessionService) Revoke(ctx context.Context, id string) error { return nil }\n",
		"apps/api/internal/routes/invoice_routes.go": `package routes

func init() {
	m.Protected.GET("/invoices", h.List)
	m.Protected.POST("/invoices", h.Create)
	m.Protected.PUT("/invoices/:id", h.Update)
	m.Staff.DELETE("/invoices/:id", h.Delete)
}
`,
		"apps/api/internal/routes/routes.go": `package routes

func Setup() {
	sentinel.MountE(r, db, sentinel.Config{
		Counters: redisstore.New(svc.Cache.Client()),
	})
}
`,
	}
}

func write(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for name, body := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// run audits a project built from the clean one with changes applied. A change
// with an empty body deletes the file.
func run(t *testing.T, changes map[string]string) *Report {
	t.Helper()
	files := cleanProject()
	for name, body := range changes {
		if body == "" {
			delete(files, name)
			continue
		}
		files[name] = body
	}
	report, err := Run(write(t, files))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	return report
}

// fired reports the findings of one check.
func fired(r *Report, check string) []Finding {
	var out []Finding
	for _, f := range r.Findings {
		if f.Check == check {
			out = append(out, f)
		}
	}
	return out
}

func TestACorrectProjectHasNothingToReport(t *testing.T) {
	report := run(t, nil)
	if len(report.Findings) != 0 {
		for _, f := range report.Findings {
			t.Errorf("%s: %s %s %s", f.Check, f.Level, f.Resource, f.Message)
		}
		t.Fatalf("a correct project produced %d finding(s); a linter that cries wolf gets turned off", len(report.Findings))
	}
	if report.Checks != len(checks) || len(report.Resources) != 1 || report.Resources[0] != "Invoice" {
		t.Errorf("checks=%d resources=%v", report.Checks, report.Resources)
	}
}

// The bug this came from: the column was readable in Postgres and the API had
// answered 201.
func TestEncryptedFieldWithNoKey(t *testing.T) {
	report := run(t, map[string]string{".env": strings.Replace(cleanProject()[".env"],
		"FIELD_ENCRYPTION_KEY=7Gq0W0l5dQ9xT2bV8nK4mJ6pR1sY3zA5cE7hL9oQ2uI=", "FIELD_ENCRYPTION_KEY=", 1)})
	found := fired(report, "encryption-key-unset")
	if len(found) != 1 || found[0].Level != "error" || !strings.Contains(found[0].Resource, "Invoice") {
		t.Fatalf("want one error naming Invoice, got %+v", found)
	}
}

func TestEncryptedColumnInAWhitelist(t *testing.T) {
	src := cleanProject()["apps/api/internal/services/invoice.go"]
	report := run(t, map[string]string{"apps/api/internal/services/invoice.go": strings.Replace(src,
		`Searchable: []string{"number"}`, `Searchable: []string{"number", "secret"}`, 1)})
	found := fired(report, "encrypted-column-in-whitelist")
	if len(found) != 1 || found[0].Resource != "Invoice.secret" {
		t.Fatalf("want one finding for Invoice.secret, got %+v", found)
	}
}

// Owned and nothing scopes it: every signed-in caller reads every row.
func TestOwnedResourceWithNoScopingAtAll(t *testing.T) {
	// Still a generated resource, with its list config: only the scoping is gone.
	report := run(t, map[string]string{"apps/api/internal/services/invoice.go": `package services

var invoiceListConfig = paginate.Config{Searchable: []string{"number"}}

func (s *InvoiceService) List(ctx context.Context, p paginate.Params) error { return nil }
`})
	found := fired(report, "owned-resource-unscoped")
	if len(found) != 1 || found[0].Level != "error" || !strings.Contains(found[0].Message, "user_id") {
		t.Fatalf("want one error naming the owner column, got %+v", found)
	}
}

// The side door: one method loses its scoping and the rest keep theirs.
func TestOwnedResourceWithOneUnscopedMethod(t *testing.T) {
	src := cleanProject()["apps/api/internal/services/invoice.go"]
	src = strings.Replace(src, `	query = authz.ScopeOwned(ctx, query, "user_id")
	return nil
}`, "\treturn nil\n}", 1) // Export loses it
	report := run(t, map[string]string{"apps/api/internal/services/invoice.go": src})
	found := fired(report, "owned-resource-unscoped")
	if len(found) != 1 || found[0].Resource != "InvoiceService.Export" {
		t.Fatalf("want one finding for Export, got %+v", found)
	}
}

func TestOwnerAcceptedFromTheBody(t *testing.T) {
	report := run(t, map[string]string{"apps/api/internal/handlers/invoice.go": "package handlers\n\n" +
		"type CreateInvoiceRequest struct {\n" +
		"\tNumber string `json:\"number\"`\n" +
		"\tUserID string `json:\"user_id\"`\n" +
		"}\n"})
	found := fired(report, "owner-settable-from-body")
	if len(found) != 1 || found[0].Level != "error" {
		t.Fatalf("want one error, got %+v", found)
	}
}

// A resource that points at a user and is scoped to nobody. A question, not an
// error: a shared catalogue with a created_by is fine.
func TestResourceThatReferencesAUserWithoutOwnership(t *testing.T) {
	model := strings.Replace(cleanProject()["apps/api/internal/models/invoice.go"],
		"func (m *Invoice) GetOwnerID() string { return m.UserID }\n", "", 1)
	report := run(t, map[string]string{
		"apps/api/internal/models/invoice.go": model,
		// Without the owner, the scoping calls are the thing out of place.
		"apps/api/internal/services/invoice.go": "package services\n\nvar invoiceListConfig = paginate.Config{}\n",
	})
	found := fired(report, "resource-could-be-owned")
	if len(found) != 1 || found[0].Level != "warning" {
		t.Fatalf("want one warning, got %+v", found)
	}
}

func TestTenantPluginWithAnUnscopedResource(t *testing.T) {
	report := run(t, map[string]string{"apps/api/internal/tenant/tenant.go": "package tenant\n"})
	found := fired(report, "tenant-shared-resource")
	if len(found) != 1 || found[0].Resource != "Invoice" {
		t.Fatalf("want one finding for Invoice, got %+v", found)
	}
	// With the embed, nothing to say.
	model := strings.Replace(cleanProject()["apps/api/internal/models/invoice.go"],
		"type Invoice struct {\n", "type Invoice struct {\n\ttenant.Owned\n", 1)
	report = run(t, map[string]string{
		"apps/api/internal/tenant/tenant.go":  "package tenant\n",
		"apps/api/internal/models/invoice.go": model,
	})
	if found := fired(report, "tenant-shared-resource"); len(found) != 0 {
		t.Fatalf("a tenant-owned model still reported: %+v", found)
	}
}

func TestPIIColumnLeftInTheClear(t *testing.T) {
	model := strings.Replace(cleanProject()["apps/api/internal/models/invoice.go"],
		"\tUserID string", "\tSSN    string                 `json:\"ssn\"`\n\tUserID string", 1)
	report := run(t, map[string]string{"apps/api/internal/models/invoice.go": model})
	found := fired(report, "pii-column-not-encrypted")
	if len(found) != 1 || found[0].Resource != "Invoice.ssn" {
		t.Fatalf("want one finding for Invoice.ssn, got %+v", found)
	}
}

func TestAppendOnlyResourceThatStillMountsAWrite(t *testing.T) {
	model := cleanProject()["apps/api/internal/models/invoice.go"] +
		"\nfunc init() { appendonly.Register(&Invoice{}) }\n"
	report := run(t, map[string]string{"apps/api/internal/models/invoice.go": model})
	found := fired(report, "append-only-mutable-routes")
	if len(found) != 1 || found[0].Level != "error" {
		t.Fatalf("want one error, got %+v", found)
	}
	// Read and create only: nothing to say.
	report = run(t, map[string]string{
		"apps/api/internal/models/invoice.go": model,
		"apps/api/internal/routes/invoice_routes.go": `package routes

func init() {
	m.Protected.GET("/invoices", h.List)
	m.Protected.POST("/invoices", h.Create)
}
`,
	})
	if found := fired(report, "append-only-mutable-routes"); len(found) != 0 {
		t.Fatalf("an append-only resource with only reads and creates reported: %+v", found)
	}
}

func TestStudioWithNoLoginAndWithTheDefaultPassword(t *testing.T) {
	env := cleanProject()[".env"]
	report := run(t, map[string]string{".env": strings.Replace(env, "GORM_STUDIO_PASSWORD=9d2b7e6053f1a8c4", "GORM_STUDIO_PASSWORD=", 1)})
	if found := fired(report, "studio-unprotected"); len(found) != 1 || found[0].Level != "error" {
		t.Fatalf("no login: want one error, got %+v", found)
	}

	// The default password is a warning in development and an error in production.
	withDefault := strings.Replace(env, "GORM_STUDIO_PASSWORD=9d2b7e6053f1a8c4", "GORM_STUDIO_PASSWORD=studio", 1)
	report = run(t, map[string]string{".env": withDefault})
	if found := fired(report, "studio-unprotected"); len(found) != 1 || found[0].Level != "warning" {
		t.Fatalf("default password in development: want one warning, got %+v", found)
	}
	report = run(t, map[string]string{".env": strings.Replace(withDefault, "APP_ENV=development", "APP_ENV=production", 1)})
	if found := fired(report, "studio-unprotected"); len(found) != 1 || found[0].Level != "error" {
		t.Fatalf("default password in production: want one error, got %+v", found)
	}

	// Switched off, nothing to guard.
	report = run(t, map[string]string{".env": env + "GORM_STUDIO_ENABLED=false\n"})
	if found := fired(report, "studio-unprotected"); len(found) != 0 {
		t.Fatalf("Studio disabled still reported: %+v", found)
	}
}

func TestDefaultCredentialsAndAWeakJWTSecret(t *testing.T) {
	env := cleanProject()[".env"]
	env = strings.Replace(env, "SENTINEL_PASSWORD=3f9c1a7d5b2e8046", "SENTINEL_PASSWORD=sentinel", 1)
	env = strings.Replace(env, "PULSE_PASSWORD=c4d7b2e905634f8a", "PULSE_PASSWORD=pulse", 1)
	env = strings.Replace(env, "JWT_SECRET=4f8a1d9c2b7e6053f1a8c4d7b2e90563a1c8f4d7b2e905634f8a1d9c2b7e6053", "JWT_SECRET=change-me", 1)
	env = strings.Replace(env, "SENTINEL_AUDIT_KEY=2c9f6053b8e1d4a78b1e4d7a2c9f6053", "", 1)
	report := run(t, map[string]string{".env": strings.Replace(env, "APP_ENV=development", "APP_ENV=production", 1)})

	found := fired(report, "default-credentials")
	var named []string
	for _, f := range found {
		named = append(named, f.Resource+":"+f.Level)
	}
	for _, want := range []string{"SENTINEL_PASSWORD:error", "PULSE_PASSWORD:error", "JWT_SECRET:error", "SENTINEL_AUDIT_KEY:warning"} {
		if !strings.Contains(strings.Join(named, " "), want) {
			t.Errorf("missing %s, got %v", want, named)
		}
	}
}

func TestFrameworkLibraryBehindItsFloor(t *testing.T) {
	report := run(t, map[string]string{"apps/api/go.mod": strings.Replace(cleanProject()["apps/api/go.mod"],
		"github.com/MUKE-coder/sentinel/v2 v2.5.0", "github.com/MUKE-coder/sentinel/v2 v2.2.1", 1)})
	found := fired(report, "framework-library-behind")
	if len(found) != 1 || !strings.Contains(found[0].Fix, "grit upgrade") {
		t.Fatalf("want one finding pointing at grit upgrade, got %+v", found)
	}
}

func TestSentinelCountingPerProcess(t *testing.T) {
	report := run(t, map[string]string{"apps/api/internal/routes/routes.go": `package routes

func Setup() {
	sentinel.MountE(r, db, sentinel.Config{})
}
`})
	found := fired(report, "rate-limits-per-process")
	if len(found) != 1 || found[0].Level != "warning" {
		t.Fatalf("want one warning, got %+v", found)
	}

	// No Redis to share them with: nothing to say.
	report = run(t, map[string]string{
		"apps/api/internal/routes/routes.go": "package routes\n\nfunc Setup() {\n\tsentinel.MountE(r, db, sentinel.Config{})\n}\n",
		".env":                               strings.Replace(cleanProject()[".env"], "REDIS_HOST=localhost", "", 1),
	})
	if found := fired(report, "rate-limits-per-process"); len(found) != 0 {
		t.Fatalf("a project with no Redis reported: %+v", found)
	}
}

func TestPublicAllowlistPublishingAHeldBackColumn(t *testing.T) {
	report := run(t, map[string]string{"apps/api/internal/handlers/invoice_public.go": "package handlers\n\n" +
		"type publicInvoice struct {\n" +
		"\tID        string  `json:\"id\"`\n" +
		"\tNumber    string  `json:\"number\"`\n" +
		"\tCostPrice float64 `json:\"cost_price\"`\n" +
		"}\n"})
	found := fired(report, "public-allowlist-sensitive")
	if len(found) != 1 || !strings.Contains(found[0].Resource, "cost_price") {
		t.Fatalf("want one finding for cost_price, got %+v", found)
	}
}

func TestNotAGritProject(t *testing.T) {
	if _, err := Run(t.TempDir()); err == nil {
		t.Fatal("a directory with no API should be refused, not reported as clean")
	}
}
