package handlers

import (
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"saas/apps/api/internal/factory"
	"saas/apps/api/internal/middleware"
	"saas/apps/api/internal/models"
	"saas/apps/api/internal/services"
	"saas/apps/api/internal/testkit"
)

// What this example exists to prove.
//
// A SaaS has two front doors over one users table: the customer portal in
// apps/web, and the operator console in apps/admin. They share a login, a JWT
// and an API. The only thing standing between one customer and another
// customer's billing history is that `grit generate resource --owned-by user`
// put authz.ScopeOwned on the list query and authz.Owns on the read.
//
// That is a security property, and a security property described in a README
// is a claim. These are the requests that make it a fact. If a future change to
// the generator drops the scope, this file fails rather than somebody's
// customer seeing a stranger's invoices.

// mountPortal mounts the two resources a customer portal reads, exactly as
// routes.go mounts them: the auth middleware, then the permission the route
// names. The portal reaches the same endpoints the admin does, which is the
// whole point of the test.
func mountPortal(t *testing.T) (*testkit.Client, *gorm.DB, *services.AuthService) {
	t.Helper()

	db := testkit.DB(t)
	auth := testkit.Auth(testkit.Config(), db)
	subscriptions := &SubscriptionHandler{DB: db}
	invoices := &InvoiceHandler{DB: db}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	guarded := r.Group("/api/v1", middleware.Auth(db, auth))
	guarded.GET("/subscriptions", middleware.RequireRole("ADMIN", "perm:subscriptions.view"), subscriptions.List)
	guarded.GET("/subscriptions/:id", middleware.RequireRole("ADMIN", "perm:subscriptions.view"), subscriptions.GetByID)
	guarded.GET("/invoices", middleware.RequireRole("ADMIN", "perm:invoices.view"), invoices.List)
	guarded.GET("/invoices/:id", middleware.RequireRole("ADMIN", "perm:invoices.view"), invoices.GetByID)

	return testkit.New(t, r), db, auth
}

// customer makes a signed-in account with the grants the portal gives its
// customers, and a subscription with one invoice against it.
//
// The grants matter. A customer is a plain USER who has been given
// subscriptions.view and invoices.view, and nothing else. They are not an
// operator with a narrower screen: they reach the same endpoint the operator
// does and the row scope is what confines them.
func customer(t *testing.T, db *gorm.DB, auth *services.AuthService, email string) (*models.User, string, *models.Subscription, *models.Invoice) {
	t.Helper()

	user := testkit.User(t, db, email, models.RoleUser)

	role := models.Role{Name: "Customer " + email}
	if err := role.SetGrants([]string{"subscriptions.view", "invoices.view"}); err != nil {
		t.Fatalf("granting the portal permissions: %v", err)
	}
	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("creating the customer role: %v", err)
	}
	if err := db.Create(&models.UserRole{UserID: user.ID, RoleID: role.ID}).Error; err != nil {
		t.Fatalf("assigning the customer role: %v", err)
	}

	plan := factory.CreatePlan(t, db)
	subscription := factory.CreateSubscription(t, db, func(m *models.Subscription) {
		m.UserID = user.ID
		m.PlanID = plan.ID
		m.Status = "active"
	})
	invoice := factory.CreateInvoice(t, db, func(m *models.Invoice) {
		m.UserID = user.ID
		m.SubscriptionID = subscription.ID
		m.Status = "paid"
	})

	return user, testkit.SignIn(t, db, auth, user), subscription, invoice
}

// A customer listing their subscriptions sees theirs and nobody else's.
func TestThePortalListsOnlyYourOwnSubscription(t *testing.T) {
	client, db, auth := mountPortal(t)

	_, ada, adasSub, _ := customer(t, db, auth, "ada@example.com")
	_, _, basilsSub, _ := customer(t, db, auth, "basil@example.com")

	body := client.As(ada).Get("/api/v1/subscriptions").AssertOK().Body()

	if !strings.Contains(body, adasSub.ID) {
		t.Error("Ada cannot see her own subscription, so the portal shows her nothing")
	}
	if strings.Contains(body, basilsSub.ID) {
		t.Fatal("Ada can see Basil's subscription: the owner scope is not on the list query")
	}
}

// And the same for the invoices, which are the ones that carry money.
func TestThePortalListsOnlyYourOwnInvoices(t *testing.T) {
	client, db, auth := mountPortal(t)

	_, ada, _, adasInvoice := customer(t, db, auth, "ada@example.com")
	_, _, _, basilsInvoice := customer(t, db, auth, "basil@example.com")

	body := client.As(ada).Get("/api/v1/invoices").AssertOK().Body()

	if !strings.Contains(body, adasInvoice.ID) {
		t.Error("Ada cannot see her own invoice")
	}
	if strings.Contains(body, basilsInvoice.ID) {
		t.Fatal("Ada can see Basil's invoice: the owner scope is not on the list query")
	}
}

// Knowing the id is not enough.
//
// A list that filters and a read that does not is the usual shape of this bug:
// the screen looks right and the record is one guessed URL away. The answer is
// 404 and not 403, because a 403 confirms the row exists, which on an endpoint
// keyed by id is the leak in miniature.
func TestAGuessedIDIsStillNotYours(t *testing.T) {
	client, db, auth := mountPortal(t)

	_, ada, _, _ := customer(t, db, auth, "ada@example.com")
	_, _, basilsSub, basilsInvoice := customer(t, db, auth, "basil@example.com")

	client.As(ada).Get("/api/v1/subscriptions/" + basilsSub.ID).AssertNotFound()
	client.As(ada).Get("/api/v1/invoices/" + basilsInvoice.ID).AssertNotFound()
}

// The operator console is the other front door, and it is supposed to see
// everything. Without this half the test above passes on an API that simply
// returns nothing to anybody.
func TestTheOperatorConsoleSeesEveryCustomer(t *testing.T) {
	client, db, auth := mountPortal(t)

	_, _, adasSub, _ := customer(t, db, auth, "ada@example.com")
	_, _, basilsSub, _ := customer(t, db, auth, "basil@example.com")

	admin := testkit.User(t, db, "operator@example.com", models.RoleAdmin)
	body := client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/subscriptions").
		AssertOK().
		Body()

	if !strings.Contains(body, adasSub.ID) || !strings.Contains(body, basilsSub.ID) {
		t.Error("the operator cannot see both customers, so the admin panel is as confined as the portal")
	}
}

// A customer who has not been granted the portal's permissions gets nothing,
// which is what makes the grant the thing that opens the door rather than the
// role name.
func TestASignedInStrangerIsNotACustomer(t *testing.T) {
	client, db, auth := mountPortal(t)

	stranger := testkit.User(t, db, "stranger@example.com", models.RoleUser)

	client.As(testkit.SignIn(t, db, auth, stranger)).
		Get("/api/v1/subscriptions").
		AssertForbidden()
}
