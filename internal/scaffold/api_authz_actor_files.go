package scaffold

// The actor: who a query runs for, carried on the context.
//
// Generated services own every query now, and a service has no gin request:
// the same method runs from a handler, a job, a command or a test. Ownership
// scoping used to read the caller straight off the gin context
// (authz.ScopeToOwner(c, ...)), which only a handler has. These helpers read
// it from context.Context instead. The handler puts the caller there with
// WithActor; a job states that it acts as the application with AsSystem; and a
// context with no actor at all sees no owned rows, so forgetting it fails
// closed rather than open.

// APIAuthzActorGo is exported so the generator can add the file to a project
// from before it existed: generated services call authz.ScopeOwned.
func APIAuthzActorGo() string { return authzActorGo() }

func authzActorGo() string {
	return `package authz

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/confidential"
)

// Actor is who a query runs for: the signed-in user, and whether ownership
// scoping should let them through, as it does an ADMIN.
type Actor struct {
	UserID string
	Admin  bool

	// Grants are the permission keys this caller holds, resolved once by the
	// auth middleware. They are carried here rather than read per query
	// because ScopeConfidential runs inside a query builder, where a database
	// read would be a second round trip on every list.
	Grants []string
}

// Holds reports whether the actor was granted permission.
func (a Actor) Holds(permission string) bool {
	for _, g := range a.Grants {
		if g == permission {
			return true
		}
	}
	return false
}

type actorKey struct{}

// ActorOf is the caller of a request, as the auth middleware identified them.
func ActorOf(c *gin.Context) Actor {
	a := Actor{UserID: CurrentUserID(c), Admin: IsAdmin(c)}
	// Resolved once per request by the auth middleware, so this reads what is
	// already there rather than making a second trip to the database.
	if grants, ok := c.Get("user_grants"); ok {
		if list, ok := grants.([]string); ok {
			a.Grants = list
		}
	}
	return a
}

// WithActor returns ctx carrying a.
func WithActor(ctx context.Context, a Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, a)
}

// ActorFrom returns the actor on ctx, if there is one.
func ActorFrom(ctx context.Context) (Actor, bool) {
	a, ok := ctx.Value(actorKey{}).(Actor)
	return a, ok
}

// AsSystem marks ctx as the application acting on its own behalf, for a job or
// a command. Ownership scoping lets it through, as it does an ADMIN.
func AsSystem(ctx context.Context) context.Context {
	return WithActor(ctx, Actor{Admin: true})
}

// UserIDFrom is the acting user's id, or "".
func UserIDFrom(ctx context.Context) string {
	a, _ := ActorFrom(ctx)
	return a.UserID
}

// ScopeOwned narrows a query to the rows the actor on ctx owns. An ADMIN or
// the system gets the query back untouched. A context with no actor, or an
// actor with no user, matches nothing: a scope that opened up when it could
// not tell who was asking would read as protection and give none.
//
// column is a fixed string from the generator, never user input.
func ScopeOwned(ctx context.Context, q *gorm.DB, column string) *gorm.DB {
	a, ok := ActorFrom(ctx)
	switch {
	case ok && a.Admin:
		return q
	case !ok || a.UserID == "":
		return q.Where("1 = 0")
	}
	return q.Where(column+" = ?", a.UserID)
}

// ScopeConfidential narrows a query to rows the actor owns, and does NOT let
// an ADMIN through.
//
// That is the whole difference from ScopeOwned, and it is the point. A product
// whose rows are customer documents, medical records, legal files, coursework,
// wants an operator who can run the system without reading what is in it.
// Being an administrator is not a reason to read somebody's medical record,
// and until now it was the only answer the framework had (grit#203).
//
// Three ways through, and no fourth:
//
//	the owner                      their own rows
//	a holder of permission         everything, and the caller audits the read
//	anything else, ADMIN included  nothing
//
// The system actor is not exempt either. A job that must read confidential
// rows holds the permission like anything else, because "it is a job" is not a
// property anybody can check afterwards and the audit trail should name
// something.
//
// column and permission are fixed strings from the generator, never input.
func ScopeConfidential(ctx context.Context, q *gorm.DB, column, permission string) *gorm.DB {
	a, ok := ActorFrom(ctx)
	switch {
	case !ok:
		// No actor at all: a migration, or a handler that forgot. Nothing,
		// which is the direction a guard has to fail in.
		return q.Where("1 = 0")
	case a.Holds(permission):
		return q
	case a.UserID == "":
		return q.Where("1 = 0")
	}
	return q.Where(column+" = ?", a.UserID)
}

// MayReadEveryConfidentialTable reports whether the actor on ctx holds the
// permission for every table registered as confidential, and names the first
// one they do not.
//
// Every, not any. A whole-database archive contains all of them, so a caller
// who may read patient records and not payroll must not be handed a file with
// payroll in it.
//
// The registry is in internal/confidential rather than here, because the
// models register themselves and this package imports models: a model file
// importing authz is an import cycle.
func MayReadEveryConfidentialTable(ctx context.Context) (string, bool) {
	for table, permission := range confidential.Tables() {
		if !MayReadConfidential(ctx, permission) {
			return table, false
		}
	}
	return "", true
}

// MayReadConfidential reports whether the actor on ctx holds permission, for
// the places that are a yes or a no rather than a query: a file download, an
// export, a backup archive.
func MayReadConfidential(ctx context.Context, permission string) bool {
	a, ok := ActorFrom(ctx)
	return ok && a.Holds(permission)
}

// Owns reports whether the actor on ctx may see row: an ADMIN, the system, or
// the row's owner.
func Owns(ctx context.Context, row Ownable) bool {
	a, ok := ActorFrom(ctx)
	if !ok {
		return false
	}
	return a.Admin || (a.UserID != "" && row.GetOwnerID() == a.UserID)
}

// OwnsConfidential is Owns without the ADMIN exemption: the owner, or a
// holder of permission, and nobody else.
func OwnsConfidential(ctx context.Context, row Ownable, permission string) bool {
	a, ok := ActorFrom(ctx)
	if !ok {
		return false
	}
	return a.Holds(permission) || (a.UserID != "" && row.GetOwnerID() == a.UserID)
}
`
}

func authzActorTestGo() string {
	return `package authz

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type actorTestRow struct{ owner string }

func (r actorTestRow) GetOwnerID() string { return r.owner }

// A service has no request, only a context. Ownership has to hold there too,
// and fail closed when nobody said who is asking.
func TestOwnershipFromTheContext(t *testing.T) {
	bg := context.Background()
	if Owns(bg, actorTestRow{"u1"}) {
		t.Error("a context with no actor saw an owned row")
	}

	alice := WithActor(bg, Actor{UserID: "u1"})
	if !Owns(alice, actorTestRow{"u1"}) {
		t.Error("an owner could not see their own row")
	}
	if Owns(alice, actorTestRow{"u2"}) {
		t.Error("a user saw somebody else's row")
	}
	if !Owns(AsSystem(bg), actorTestRow{"u2"}) {
		t.Error("the system could not see a row")
	}
	if UserIDFrom(alice) != "u1" || UserIDFrom(bg) != "" {
		t.Error("UserIDFrom read the wrong user")
	}

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Set("user_id", "u9")
	c.Set("user_role", "ADMIN")
	if a := ActorOf(c); a.UserID != "u9" || !a.Admin {
		t.Errorf("ActorOf read %+v from the request", a)
	}
}
`
}
