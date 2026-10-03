package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// writeResourceTest emits internal/handlers/<snake>_test.go: the tests that come
// with a generated resource.
//
// # Why the generator writes tests at all
//
// Until now `grit generate resource` wrote twelve files and no test, and the
// only thing standing between a generator change and a broken resource was
// somebody scaffolding a project by hand and looking at it. A test per resource
// means a generator change is checked by the thing it generates, in the
// project's own CI, on the developer's machine.
//
// # What it asserts, and why these four
//
// Not the CRUD mechanics: those are the same code for every resource and are
// already covered by the framework's own suite. What varies per resource, and
// is therefore worth generating a test for, is the wiring:
//
//  1. The route exists and refuses an anonymous caller. A resource whose routes
//     were injected into the wrong group is reachable by anyone, and nothing
//     else in the project would notice.
//  2. The permission is enforced. A signed-in user without widgets.view gets a
//     403, which is the difference between a permission that is declared and one
//     that is checked.
//  3. A created row comes back from the list. This is the round trip through the
//     model, the service, the handler and the serialiser, which is where a wrong
//     column name or a missing json tag shows up.
//  4. An unknown id is a 404 and not a 500. The most common shape of
//     generated-handler bug, because it depends on how the service reports "no
//     rows".
func (g *Generator) writeResourceTest(names Names) error {
	path := filepath.Join(g.APIRoot(), "internal", "handlers", names.Snake+"_test.go")
	return writeFileWithDirs(path, g.resourceTestSource(names))
}

func (g *Generator) resourceTestSource(names Names) string {
	var b strings.Builder

	fmt.Fprintf(&b, `package handlers

import (
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	%q
	%q
	%q
	%q
	%q
)

// mount%s builds an engine carrying this resource's routes exactly as routes.go
// mounts them: the auth middleware, then the permission the route names. Wiring
// it any other way would test a router this project does not run.
func mount%s(t *testing.T) (*testkit.Client, *gorm.DB, *services.AuthService) {
	t.Helper()

	db := testkit.DB(t)
	auth := testkit.Auth(testkit.Config(), db)
	handler := &%sHandler{DB: db}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	guarded := r.Group("/api/v1", middleware.Auth(db, auth))
	guarded.GET("/%s", middleware.RequireRole("ADMIN", "perm:%s.view"), handler.List)
	guarded.GET("/%s/:id", middleware.RequireRole("ADMIN", "perm:%s.view"), handler.GetByID)

	return testkit.New(t, r), db, auth
}

// A resource whose routes were injected into the wrong group is reachable by
// anyone, and nothing else in the project would notice.
func Test%s_ListRefusesAnonymous(t *testing.T) {
	client, _, _ := mount%s(t)

	client.Get("/api/v1/%s").AssertUnauthorized()
}

// The difference between a permission that is declared and one that is checked.
func Test%s_ListRefusesAUserWithoutThePermission(t *testing.T) {
	client, db, auth := mount%s(t)

	user := testkit.User(t, db, "no-permission@example.com", models.RoleUser)

	client.As(testkit.SignIn(t, db, auth, user)).
		Get("/api/v1/%s").
		AssertForbidden()
}

// The round trip through the model, the service, the handler and the
// serialiser, which is where a wrong column name or a missing json tag shows up.
func Test%s_AdminSeesWhatWasCreated(t *testing.T) {
	client, db, auth := mount%s(t)

	admin := testkit.User(t, db, "admin@example.com", models.RoleAdmin)
	record := factory.Create%s(t, db)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/%s").
		AssertOK().
		JSON().
		Count("data", 1).
		Where("data.0.id", record.ID)
}

// The most common shape of generated-handler bug, because it depends on how the
// service reports "no rows": a 500 here means the error was passed through
// rather than recognised.
func Test%s_UnknownIDIsNotFound(t *testing.T) {
	client, db, auth := mount%s(t)

	admin := testkit.User(t, db, "admin-404@example.com", models.RoleAdmin)

	client.As(testkit.SignIn(t, db, auth, admin)).
		Get("/api/v1/%s/01920000-0000-7000-8000-000000000000").
		AssertNotFound()
}
`,
		g.Module+"/internal/factory",
		g.Module+"/internal/middleware",
		g.Module+"/internal/models",
		g.Module+"/internal/services",
		g.Module+"/internal/testkit",
		names.Pascal, names.Pascal, names.Pascal,
		names.PluralKebab, names.PluralSnake,
		names.PluralKebab, names.PluralSnake,
		names.Pascal, names.Pascal, names.PluralKebab,
		names.Pascal, names.Pascal, names.PluralKebab,
		names.Pascal, names.Pascal, names.Pascal, names.PluralKebab,
		names.Pascal, names.Pascal, names.PluralKebab,
	)

	return b.String()
}
