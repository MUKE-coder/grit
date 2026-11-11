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

	// Every belongs_to parent, made first and passed in.
	//
	// The factory deliberately leaves the column empty and its doc comment
	// says to pass the parent in; this test did not, so the row went in with
	// an empty foreign key. SQLite accepted it because the test database did
	// not enforce foreign keys, and Postgres, which is what the project runs
	// on, refused it. Green where nobody ships, red where everybody does.
	// See grit#137.
	parentFixtures, parentOverride, secondRowOverride := g.parentFixtures(names)

	// Who signs in. An ADMIN for every resource but a --confidential one,
	// where an ADMIN sees nothing on purpose and a test written that way
	// asserts the opposite of the flag (grit#205).
	actorLabel, actorFixture, actorVar := g.actorFor(names, "", "view")
	// Each test inserts its own user, so each needs its own email: two tests
	// in one package creating the same address collide on the unique index.
	// The variable name is the same in all of them, which is what lets one
	// override line serve every fixture row.
	writerFixture := func(suffix string) string {
		_, fixture, _ := g.actorFor(names, suffix, "view", "edit")
		return fixture
	}
	writerVar := actorVar
	if g.Definition.Confidential {
		parentFixtures, parentOverride = g.confidentialOwnerOverride(names, parentFixtures, parentOverride, actorVar)
		_, secondRowOverride = g.confidentialOwnerOverride(names, "", secondRowOverride, actorVar)
	}

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
	guarded.PUT("/%s/:id", middleware.RequireRole("ADMIN", "perm:%s.edit"), handler.Update)
	guarded.POST("/%s/bulk-edit", middleware.RequireRole("ADMIN", "perm:%s.edit"), handler.BulkEdit)

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
func Test%s_%sSeesWhatWasCreated(t *testing.T) {
	client, db, auth := mount%s(t)

%s%s	record := factory.Create%s(t, db%s)

	client.As(testkit.SignIn(t, db, auth, %s)).
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
		// 6, 7, 8: mount's doc comment, mount, the handler type.
		names.Pascal, names.Pascal, names.Pascal,
		// 9 to 16: the four routes.
		names.PluralKebab, names.PluralSnake,
		names.PluralKebab, names.PluralSnake,
		names.PluralKebab, names.PluralSnake,
		names.PluralKebab, names.PluralSnake,
		// 13 to 15, 16 to 18: the two refusal tests.
		names.Pascal, names.Pascal, names.PluralKebab,
		names.Pascal, names.Pascal, names.PluralKebab,
		// 19 to 27: the round-trip test. Its caller is an ADMIN, or for a
		// --confidential resource the row's owner, because an ADMIN sees
		// nothing there on purpose (grit#205). Its fixture line creates every
		// belongs_to parent first (grit#137).
		names.Pascal, actorLabel, names.Pascal,
		actorFixture, parentFixtures, names.Pascal, parentOverride, actorVar,
		names.PluralKebab,
		// 25 to 27: the 404 test.
		names.Pascal, names.Pascal, names.PluralKebab,
	)

	// The two tiers the flag adds, which the tests above cannot cover: they
	// sign in as the owner, and these are about everybody else.
	//
	// The row belongs to somebody neither caller is, so "the list is empty"
	// means the scope refused it rather than that nothing was created. Pointed
	// at rowOwner rather than at the signed-in user, which is the difference
	// between asserting the rule and asserting an empty table.
	if g.Definition.Confidential {
		hiddenFixtures, hiddenOverride := g.confidentialOwnerOverride(
			names, parentFixtures, parentOverride, "rowOwner")
		fmt.Fprintf(&b, confidentialTiersTest,
			// An ADMIN sees nothing.
			names.Pascal, names.Pascal,
			hiddenFixtures, names.Pascal, hiddenOverride,
			names.PluralKebab, names.PluralKebab,
			// A grant holder sees everything.
			names.PluralSnake,
			names.Pascal, names.Pascal,
			hiddenFixtures, names.Pascal, hiddenOverride,
			names.PluralSnake, names.PluralSnake,
			names.PluralKebab,
		)
	}

	// The grid test needs a patch with something in it. BulkEdit skips an edit
	// whose writable columns are all unchanged, which is correct (touching a
	// cell and putting it back is a thing people do) and means an empty patch
	// never reaches the version check: the first version of this test asserted
	// a 422 and got a 200 for that reason. A resource with no plain text column
	// gets the stale-save test and not this one, rather than a test that cannot
	// assert what it says.
	// An append-only resource has no update to be stale about. It has no PUT
	// route, the GORM guard refuses the write, and the database trigger
	// refuses it under that: a version conflict cannot happen because no
	// version can change. These tests passed until v3.494.0 only because the
	// test database had no triggers, which is the hole grit#194 closed.
	if g.Definition.AppendOnly {
		return b.String()
	}

	patchKey, patchValue := g.patchableTextField()
	if patchKey == "" {
		// No column this test could change, so the grid half is left out
		// rather than emitted asserting a 422 it cannot cause.
		fmt.Fprintf(&b, staleSaveTest,
			names.Pascal, names.Pascal,
			writerFixture("version"), parentFixtures, names.Pascal, parentOverride, writerVar,
			names.PluralKebab, names.PluralKebab, names.PluralKebab, names.PluralKebab,
		)
		return b.String()
	}

	// A second Fprintf rather than more of the first: that one is already
	// twenty-seven positional arguments, and inserting into the middle of it
	// renumbers every comment that explains which is which.
	fmt.Fprintf(&b, staleSaveTest,
		names.Pascal, names.Pascal,
		writerFixture("version"), parentFixtures, names.Pascal, parentOverride, writerVar,
		names.PluralKebab, names.PluralKebab, names.PluralKebab, names.PluralKebab,
	)
	fmt.Fprintf(&b, gridStaleRowTest,
		names.Pascal, names.Pascal,
		writerFixture("grid"),
		// The parents are created once and both rows share them: two of them
		// would redeclare the same variables, which is what v3.483.0 did.
		// The second row drops any one_to_one, whose foreign key is unique and
		// cannot be pointed at twice (grit#195).
		parentFixtures, names.Pascal, parentOverride,
		names.Pascal, secondRowOverride, writerVar,
		names.PluralKebab,
		names.PluralKebab, patchKey, patchValue, patchKey, patchValue,
		names.PluralKebab,
		names.PluralKebab, patchKey, patchValue,
	)

	return b.String()
}

// parentFixtures returns the lines that create this resource's belongs_to
// parents, and the override that points the record at them.
//
// A self-reference is skipped: a resource whose parent is its own kind can be
// created with none, and making one first would need the same override and
// recurse.
func (g *Generator) parentFixtures(names Names) (fixtures, override, secondOverride string) {
	var created []string
	var assignments []string
	// The same, without the relations whose foreign key is unique. A second row
	// cannot point at the same parent through a one_to_one, so it points at one
	// of its own: the factory makes it, because the column is left empty
	// (grit#195).
	var shareable []string
	seen := map[string]bool{}

	for _, f := range g.Definition.Fields {
		if !f.IsBelongsTo() {
			continue
		}
		related := f.RelatedModelName()
		if related == "" || related == names.Pascal {
			continue
		}
		// An optional parent is left nil, like the factory leaves it. The
		// column is a *string, so assigning parent.ID does not compile; and a
		// row without the relation is a valid row, which is what the column
		// being optional means, so the test does not need the parent at all
		// (grit#193).
		if g.relOptional(f) {
			continue
		}
		variable := toCamelCase(f.Name)
		if seen[variable] {
			continue
		}
		seen[variable] = true

		created = append(created, fmt.Sprintf("\t%s := factory.Create%s(t, db)\n",
			variable, related))
		assignment := fmt.Sprintf("\t\tm.%s = %s.ID\n",
			toPascalCase(f.FKColumnName()), variable)
		assignments = append(assignments, assignment)
		if !f.IsOneToOne() {
			shareable = append(shareable, assignment)
		}
	}

	if len(created) == 0 {
		return "", "", ""
	}

	fixtures = strings.Join(created, "") + "\n"
	override = fmt.Sprintf(", func(m *models.%s) {\n%s\t}",
		names.Pascal, strings.Join(assignments, ""))
	secondOverride = override
	if len(shareable) != len(assignments) {
		if len(shareable) == 0 {
			// Every relation this resource has is unique, so the second row
			// shares nothing and the factory fills all of it.
			secondOverride = ""
		} else {
			secondOverride = fmt.Sprintf(", func(m *models.%s) {\n%s\t}",
				names.Pascal, strings.Join(shareable, ""))
		}
	}
	return fixtures, override, secondOverride
}

// patchableTextField is a column the grid test can change: the first plain text
// field that is not unique.
//
// Not unique, because writing one value to two rows of a unique column is
// either a constraint violation or, worse, not one. Plain text, because the
// value has to be writable as a JSON string without knowing the column's
// format: a money field is two columns, an enum rejects anything off its list,
// and a date wants a layout.
func (g *Generator) patchableTextField() (key, value string) {
	for _, f := range g.Definition.Fields {
		if f.Unique || f.ServerSet || f.IsBelongsTo() {
			continue
		}
		switch FieldType(f.Type) {
		case FieldString, FieldText, FieldRichtext:
			return toSnakeCase(f.Name), "edited by the grid test"
		}
	}
	return "", ""
}

// staleSaveTest is the If-Match round trip, for every resource.
const staleSaveTest = `

// Two people editing the same record.
//
// The API has carried a version column, an ETag and a 409 for this since long
// before any client sent the precondition, so the whole mechanism was inert and
// looked exactly like a working one (grit#182). It is only visible from the
// request: hence a test on the header rather than on the handler.
func Test%s_StaleSaveIsRefused(t *testing.T) {
	client, db, auth := mount%s(t)

%s%s	record := factory.Create%s(t, db%s)
	client = client.As(testkit.SignIn(t, db, auth, %s))

	// The version the record is at, and the ETag the client is meant to send
	// back. Both, because a detail response missing the ETag is a client with
	// nothing to send.
	client.Get("/api/v1/%s/" + record.ID).
		AssertOK().
		AssertHeader("ETag", ` + "`" + `W/"1"` + "`" + `).
		JSON().
		Where("data.version", float64(1))

	// The first save, with the version that was read.
	client.With("If-Match", ` + "`" + `W/"1"` + "`" + `).
		Put("/api/v1/%s/" + record.ID, map[string]any{}).
		AssertOK()

	// The second, with the same version, which the record has moved past. This
	// is somebody saving a form they opened before the first save landed.
	client.With("If-Match", ` + "`" + `W/"1"` + "`" + `).
		Put("/api/v1/%s/" + record.ID, map[string]any{}).
		AssertConflict().
		AssertCode("VERSION_CONFLICT")

	// And a client that sends no precondition is unaffected. Making the header
	// compulsory would break every script and every integration that writes to
	// this API.
	client.Put("/api/v1/%s/" + record.ID, map[string]any{}).AssertOK()
}
`

// gridStaleRowTest is the same for bulk-edit, emitted only for a resource with
// a plain text column the test can change.
const gridStaleRowTest = `
// The same, for the grid, which is where two people editing at once is likeliest.
//
// One If-Match header cannot carry a version per row, so each item in a
// bulk-edit carries its own, and a stale row is reported by its index so the
// grid can put the message on that line.
func Test%s_BulkEditRefusesAStaleRow(t *testing.T) {
	client, db, auth := mount%s(t)

%s%s	first := factory.Create%s(t, db%s)
	second := factory.Create%s(t, db%s)
	client = client.As(testkit.SignIn(t, db, auth, %s))

	// Somebody else saves the second row, so the grid's copy of it is stale.
	client.With("If-Match", ` + "`" + `W/"1"` + "`" + `).
		Put("/api/v1/%s/" + second.ID, map[string]any{}).
		AssertOK()

	client.Post("/api/v1/%s/bulk-edit", map[string]any{
		"items": []map[string]any{
			{"id": first.ID, "patch": map[string]any{%q: %q}, "version": 1},
			{"id": second.ID, "patch": map[string]any{%q: %q}, "version": 1},
		},
	}).
		AssertUnprocessable().
		JSON().
		Where("error.details.rows.0.index", float64(1))

	// The first row is untouched: the whole grid is refused rather than the
	// half of it that was still current. A spreadsheet half saved is worse than
	// one not saved.
	client.Get("/api/v1/%s/" + first.ID).
		AssertOK().
		JSON().
		Where("data.version", float64(1))

	// With no versions at all, as a client written before this did, the save
	// goes through.
	client.Post("/api/v1/%s/bulk-edit", map[string]any{
		"items": []map[string]any{
			{"id": first.ID, "patch": map[string]any{%q: %q}},
		},
	}).AssertOK()
}
`

// --------------------------------------------------------------- the actor
//
// Who the generated tests sign in as.
//
// For every resource but one it is an ADMIN, because an ADMIN can reach
// everything and the tests are about the round trip rather than about
// authorization. For a --confidential resource an ADMIN can reach nothing,
// deliberately, so tests written that way assert the opposite of the flag's
// contract: three of them failed on a freshly generated resource, and the
// obvious reading was that the flag was broken (grit#205).
//
// So a confidential resource's tests act as the row's owner, holding the
// permission the route names, and the row is created belonging to them.

// actorFor is the caller of the tests that expect to succeed, and the one
// fixture line that creates them.
//
// suffix distinguishes the users of tests that run in the same package: every
// test here uses its own email, because testkit.User inserts a row and two
// tests inserting "admin@example.com" collide on the unique index.
func (g *Generator) actorFor(names Names, suffix string, abilities ...string) (label, fixture, variable string) {
	email := "admin"
	if suffix != "" {
		email += "-" + suffix
	}
	if !g.Definition.Confidential {
		return "Admin", fmt.Sprintf("\tadmin := testkit.User(t, db, %q, models.RoleAdmin)\n",
			email+"@example.com"), "admin"
	}

	// A caller holding the route's permission and nothing else: not an ADMIN,
	// whose exemption is the thing --confidential takes away.
	var grants []string
	for _, a := range abilities {
		grants = append(grants, fmt.Sprintf("%q", names.PluralSnake+"."+a))
	}
	email = "owner"
	if suffix != "" {
		email += "-" + suffix
	}
	return "Owner", fmt.Sprintf("\towner := testkit.Granted(t, db, %q, %s)\n",
		email+"@example.com", strings.Join(grants, ", ")), "owner"
}

// confidentialOwnerOverride points the fixture row at the signed-in owner.
//
// The owner is itself a belongs_to field, so parentFixtures already creates a
// user for it and assigns it. That user is not the one signing in, which would
// make the owner-sees-their-own-rows test assert nothing: it would read as an
// empty list either way. This replaces that one assignment and drops the
// parent fixture line that is now unused.
func (g *Generator) confidentialOwnerOverride(names Names, fixtures, override, ownerVar string) (string, string) {
	owner := g.Definition.OwnerField()
	if owner == nil {
		return fixtures, override
	}
	variable := toCamelCase(owner.Name)
	field := toPascalCase(owner.FKColumnName())

	// The parent the override pointed at, which nothing needs now.
	fixtures = strings.Replace(fixtures,
		fmt.Sprintf("\t%s := factory.Create%s(t, db)\n", variable, owner.RelatedModelName()), "", 1)
	if strings.TrimSpace(fixtures) == "" {
		fixtures = ""
	}

	assignment := fmt.Sprintf("\t\tm.%s = %s.ID\n", field, variable)
	mine := fmt.Sprintf("\t\tm.%s = %s.ID\n", field, ownerVar)
	if strings.Contains(override, assignment) {
		return fixtures, strings.Replace(override, assignment, mine, 1)
	}
	// No belongs_to override at all, which is a resource whose only relation
	// is the owner and it was skipped as optional. Build one.
	return fixtures, fmt.Sprintf(", func(m *models.%s) {\n%s\t}", names.Pascal, mine)
}

// confidentialTiersTest is the two assertions the flag exists for, and which a
// generated project had nowhere: an administrator sees nothing, and a holder of
// the grant sees everything.
//
// Only for a --confidential resource. For any other, "an ADMIN sees nothing"
// is false and the test would be asserting a bug.
const confidentialTiersTest = `
// --confidential: the operator. An empty list and a 404, which is the whole
// point of the flag and is otherwise asserted nowhere in a generated project.
//
// A 404 rather than a 403, deliberately: a 403 confirms the row exists, which
// is itself something an administrator is not meant to learn here.
func Test%s_AnAdminSeesNothing(t *testing.T) {
	client, db, auth := mount%s(t)

	rowOwner := testkit.User(t, db, "owner-hidden@example.com", models.RoleUser)
%s	record := factory.Create%s(t, db%s)
	admin := testkit.User(t, db, "admin-blind@example.com", models.RoleAdmin)
	client = client.As(testkit.SignIn(t, db, auth, admin))

	client.Get("/api/v1/%s").AssertOK().JSON().Count("data", 0)
	client.Get("/api/v1/%s/" + record.ID).AssertNotFound()
}

// --confidential: the third tier. A named person holding %s.confidential
// reads rows they do not own, and every read is in the activity log, which is
// why the grant goes to a person rather than to a role.
func Test%s_AGrantHolderSeesEveryRow(t *testing.T) {
	client, db, auth := mount%s(t)

	rowOwner := testkit.User(t, db, "owner-seen@example.com", models.RoleUser)
%s	record := factory.Create%s(t, db%s)
	reader := testkit.Granted(t, db, "confidential-reader@example.com",
		"%s.view", "%s.confidential")

	client.As(testkit.SignIn(t, db, auth, reader)).
		Get("/api/v1/%s").
		AssertOK().
		JSON().
		Count("data", 1).
		Where("data.0.id", record.ID)
}
`
