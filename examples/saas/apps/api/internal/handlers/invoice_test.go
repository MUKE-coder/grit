package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/factory"
	"saas/apps/api/internal/middleware"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/services"
	"saas/apps/api/internal/testkit"
)

// mountInvoice builds an engine carrying this resource's routes exactly as routes.go
// mounts them: the auth middleware, then the permission the route names. Wiring
// it any other way would test a router this project does not run.
func mountInvoice(t *testing.T) (*testkit.Client, *gorm.DB, *services.AuthService) {
	t.Helper()

	db := testkit.DB(t)
	auth := testkit.Auth(testkit.Config(), db)
	handler := &InvoiceHandler{DB: db}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	guarded := r.Group("/api/v1", middleware.Auth(db, auth))
	guarded.GET("/invoices", middleware.RequireRole("ADMIN", "perm:invoices.view"), handler.List)
	guarded.GET("/invoices/:id", middleware.RequireRole("ADMIN", "perm:invoices.view"), handler.GetByID)

	return testkit.New(t, r), db, auth
}

// A resource whose routes were injected into the wrong group is reachable by
// anyone, and nothing else in the project would notice.
func TestInvoice_ListRefusesAnonymous(t *testing.T) {
	client, _, _ := mountInvoice(t)

	client.Get("/api/v1/invoices").AssertUnauthorized()
}

// The difference between a permission that is declared and one that is checked.
func TestInvoice_ListRefusesAUserWithoutThePermission(t *testing.T) {
	client, db, auth := mountInvoice(t)

	user := testkit.User(t, db, "no-permission@example.com", models.RoleUser)

	client.As(testkit.SignIn(t, db, auth, user)).
		Get("/api/v1/invoices").
		AssertForbidden()
}

// The round trip through the model, the service, the handler and the
// serialiser, which is where a wrong column name or a missing json tag shows up.
func TestInvoice_AdminSeesWhatWasCreated(t *testing.T) {
	client, db, auth := mountInvoice(t)

	admin := testkit.User(t, db, "admin@example.com", models.RoleAdmin)
	record := factory.CreateInvoice(t, db)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/invoices").
		AssertOK().
		JSON().
		Count("data", 1).
		Where("data.0.id", record.ID)
}

// The most common shape of generated-handler bug, because it depends on how the
// service reports "no rows": a 500 here means the error was passed through
// rather than recognised.
func TestInvoice_UnknownIDIsNotFound(t *testing.T) {
	client, db, auth := mountInvoice(t)

	admin := testkit.User(t, db, "admin-404@example.com", models.RoleAdmin)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/invoices/01920000-0000-7000-8000-000000000000").
		AssertNotFound()
}
