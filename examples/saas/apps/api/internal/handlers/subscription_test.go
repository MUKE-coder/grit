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

// mountSubscription builds an engine carrying this resource's routes exactly as routes.go
// mounts them: the auth middleware, then the permission the route names. Wiring
// it any other way would test a router this project does not run.
func mountSubscription(t *testing.T) (*testkit.Client, *gorm.DB, *services.AuthService) {
	t.Helper()

	db := testkit.DB(t)
	auth := testkit.Auth(testkit.Config(), db)
	handler := &SubscriptionHandler{DB: db}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	guarded := r.Group("/api/v1", middleware.Auth(db, auth))
	guarded.GET("/subscriptions", middleware.RequireRole("ADMIN", "perm:subscriptions.view"), handler.List)
	guarded.GET("/subscriptions/:id", middleware.RequireRole("ADMIN", "perm:subscriptions.view"), handler.GetByID)

	return testkit.New(t, r), db, auth
}

// A resource whose routes were injected into the wrong group is reachable by
// anyone, and nothing else in the project would notice.
func TestSubscription_ListRefusesAnonymous(t *testing.T) {
	client, _, _ := mountSubscription(t)

	client.Get("/api/v1/subscriptions").AssertUnauthorized()
}

// The difference between a permission that is declared and one that is checked.
func TestSubscription_ListRefusesAUserWithoutThePermission(t *testing.T) {
	client, db, auth := mountSubscription(t)

	user := testkit.User(t, db, "no-permission@example.com", models.RoleUser)

	client.As(testkit.SignIn(t, db, auth, user)).
		Get("/api/v1/subscriptions").
		AssertForbidden()
}

// The round trip through the model, the service, the handler and the
// serialiser, which is where a wrong column name or a missing json tag shows up.
func TestSubscription_AdminSeesWhatWasCreated(t *testing.T) {
	client, db, auth := mountSubscription(t)

	admin := testkit.User(t, db, "admin@example.com", models.RoleAdmin)
	record := factory.CreateSubscription(t, db)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/subscriptions").
		AssertOK().
		JSON().
		Count("data", 1).
		Where("data.0.id", record.ID)
}

// The most common shape of generated-handler bug, because it depends on how the
// service reports "no rows": a 500 here means the error was passed through
// rather than recognised.
func TestSubscription_UnknownIDIsNotFound(t *testing.T) {
	client, db, auth := mountSubscription(t)

	admin := testkit.User(t, db, "admin-404@example.com", models.RoleAdmin)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/subscriptions/01920000-0000-7000-8000-000000000000").
		AssertNotFound()
}
