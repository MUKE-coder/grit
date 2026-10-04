package handlers

import (
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"library/apps/api/internal/authz"
	"library/apps/api/internal/factory"
	"library/apps/api/internal/jsontime"
	"library/apps/api/internal/middleware"
	"library/apps/api/internal/models"
	"library/apps/api/internal/policies"
	"library/apps/api/internal/services"
	"library/apps/api/internal/testkit"
)

// The rule in internal/policies/loan.go, end to end.
//
// This is the test the policy layer exists for, and it is the reason this
// example is checked in rather than described: the claim that a returned loan
// cannot be edited is either true of the running code or it is marketing, and
// the only way to tell is to make the request.
func TestAReturnedLoanCannotBeEdited(t *testing.T) {
	// The rules are registered by main in a real run. A test registers them
	// itself so the file under test is the file being tested, rather than
	// whatever a previous test happened to define.
	authz.ResetForTest()
	t.Cleanup(authz.ResetForTest)
	policies.Register()

	client, db, auth := mountLoans(t)
	librarian := testkit.User(t, db, "librarian@example.com", models.RoleAdmin)
	token := testkit.SignIn(t, db, auth, librarian)

	author := factory.CreateAuthor(t, db)
	book := factory.CreateBook(t, db, func(m *models.Book) { m.AuthorID = author.ID })

	// A loan still out: editable, which is the half that proves the rule is a
	// rule and not a blanket refusal.
	out := factory.CreateLoan(t, db, func(m *models.Loan) {
		m.BookID = book.ID
		m.BorrowerID = librarian.ID
		m.ReturnedAt = nil
	})
	client.As(token).
		Put("/api/v1/loans/"+out.ID, map[string]any{"status": "overdue"}).
		AssertOK()

	// A loan that came back: refused, with words a person can act on.
	returned := factory.CreateLoan(t, db, func(m *models.Loan) {
		m.BookID = book.ID
		m.BorrowerID = librarian.ID
		m.ReturnedAt = &jsontime.DateTime{Time: time.Date(2026, 3, 2, 10, 0, 0, 0, time.UTC)}
	})
	client.As(token).
		Put("/api/v1/loans/"+returned.ID, map[string]any{"status": "out"}).
		AssertForbidden().
		AssertCode("FORBIDDEN").
		JSON().
		Where("error.message",
			"This loan was returned on 2 March 2026, and a returned loan is a record of "+
				"what happened rather than something to edit. Create a new loan if the book "+
				"has gone out again.")
}

// Deleting is a separate ability, and the rules differ: a loan still out can be
// deleted, a returned one cannot. A single "loans.write" could not say that.
func TestDeletingIsJudgedSeparatelyFromEditing(t *testing.T) {
	authz.ResetForTest()
	t.Cleanup(authz.ResetForTest)
	policies.Register()

	client, db, auth := mountLoans(t)
	librarian := testkit.User(t, db, "librarian2@example.com", models.RoleAdmin)
	token := testkit.SignIn(t, db, auth, librarian)

	author := factory.CreateAuthor(t, db)
	book := factory.CreateBook(t, db, func(m *models.Book) { m.AuthorID = author.ID })

	returned := factory.CreateLoan(t, db, func(m *models.Loan) {
		m.BookID = book.ID
		m.BorrowerID = librarian.ID
		m.ReturnedAt = &jsontime.DateTime{Time: time.Date(2026, 3, 2, 10, 0, 0, 0, time.UTC)}
	})
	client.As(token).Delete("/api/v1/loans/" + returned.ID).AssertForbidden()

	out := factory.CreateLoan(t, db, func(m *models.Loan) {
		m.BookID = book.ID
		m.BorrowerID = librarian.ID
		m.ReturnedAt = nil
	})
	client.As(token).Delete("/api/v1/loans/" + out.ID).AssertOK()
}

// mountLoans builds an engine carrying the loan routes as routes.go mounts
// them: the auth middleware, then the permission each route names.
func mountLoans(t *testing.T) (*testkit.Client, *gorm.DB, *services.AuthService) {
	t.Helper()

	db := testkit.DB(t)
	auth := testkit.Auth(testkit.Config(), db)
	handler := &LoanHandler{DB: db}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	guarded := r.Group("/api/v1", middleware.Auth(db, auth))
	guarded.PUT("/loans/:id", middleware.RequireRole("ADMIN", "perm:loans.edit"), handler.Update)
	guarded.DELETE("/loans/:id", middleware.RequireRole("ADMIN", "perm:loans.delete"), handler.Delete)

	return testkit.New(t, r), db, auth
}
